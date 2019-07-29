package fluentd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/go-logr/logr"
	fdsyncer "github.com/platform9/fluentd-operator/pkg/fluentd/internal/syncer"
	"github.com/platform9/fluentd-operator/pkg/options"
	"github.com/platform9/fluentd-operator/pkg/utils"
	"github.com/presslabs/controller-util/syncer"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("fluentd_reconciler")

// Reconciler reconciles fluentd deployment
type Reconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// New returns new instance of reconciler
func New(mgr manager.Manager) *Reconciler {
	return &Reconciler{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetRecorder("controller.fluentd"),
	}
}

// Reconcile reads state of cluster for Deployment objects and makes
// changes per how the controller definition needs to be
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	// TODO: Watch deployments only in current namespace
	instance := &appsv1.Deployment{}

	// Only interested in configured namespace for fluentd deployment
	if request.Namespace != *(options.LogNs) {
		return reconcile.Result{}, nil
	}

	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, r.CreateIfNeeded()
		}
		return reconcile.Result{}, err
	}

	return r.reconcile(reqLogger, instance)
}

// Reconcile ensures fluentd deployment is according to operator definition
func (r *Reconciler) reconcile(reqLogger logr.Logger, current *appsv1.Deployment) (reconcile.Result, error) {
	// Compare labels
	objLabels := current.GetLabels()

	if !utils.CheckSubset(fdsyncer.Labels, objLabels) {
		reqLogger.Info("Is not interesting")
		return reconcile.Result{}, nil
	}

	syncers := []syncer.Interface{fdsyncer.NewFluentdSyncer(r.client, r.scheme),
		fdsyncer.NewFluentdCfgMapSyncer(r.client, r.scheme),
		fdsyncer.NewFluentdSvcSyncer(r.client, r.scheme),
	}

	for _, sync := range syncers {
		if err := syncer.Sync(context.TODO(), sync, r.recorder); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// CreateIfNeeded creates fluentd deployment if needed
func (r *Reconciler) CreateIfNeeded() error {
	return createIfNeeded(r.client, r.scheme, r.recorder)
}

func createIfNeeded(c client.Client, s *runtime.Scheme, e record.EventRecorder) error {
	syncers := []syncer.Interface{fdsyncer.NewFluentdSyncer(c, s),
		fdsyncer.NewFluentdCfgMapSyncer(c, s),
		fdsyncer.NewFluentdSvcSyncer(c, s),
	}

	for _, sync := range syncers {
		if err := syncer.Sync(context.TODO(), sync, e); err != nil {
			if errors.IsAlreadyExists(err) {
				log.Info("fluentd object already exists, skipping...")
				continue
			}
			return err
		}
	}
	return nil
}

// Refresh changes the fluentd configmap and reload it
func (r *Reconciler) Refresh(data []byte) error {
	return refresh(r.client, r.scheme, r.recorder, data)
}

func refresh(c client.Client, s *runtime.Scheme, e record.EventRecorder, data []byte) error {
	syncers := []syncer.Interface{
		fdsyncer.NewFluentdCfgMapSyncer(c, s, data),
	}

	for _, sync := range syncers {
		if err := syncer.Sync(context.TODO(), sync, e); err != nil {
			return err
		}
	}

	// Reload service, if needed
	svcURL := fmt.Sprintf("http://%s:%d/api/config.reload", *(options.ReloadHost), *(options.ReloadPort))
	req, err := http.NewRequest("POST", svcURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Info("When reloading fluentd", "error", err)
	} else {
		respStr, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			log.Info("fluentd reload response", "status", resp.StatusCode, "message", string(respStr))
		} else {
			log.Error(err, "when reading response")
		}
		resp.Body.Close()
	}

	return err
}
