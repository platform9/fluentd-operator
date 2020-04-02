package output

import (
	"bytes"
	"context"
	"fmt"

	loggingv1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	"github.com/platform9/fluentd-operator/pkg/fluentd"
	"github.com/platform9/fluentd-operator/pkg/resources"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("output_reconciler")

// Reconciler reconciles a Output object
type Reconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client  client.Client
	scheme  *runtime.Scheme
	fluentd *fluentd.Reconciler
}

// New returns new instance of reconciler
func New(mgr manager.Manager) *Reconciler {
	return &Reconciler{
		client:  mgr.GetClient(),
		scheme:  mgr.GetScheme(),
		fluentd: fluentd.New(mgr),
	}
}

// Reconcile reads that state of the cluster for a Output object and makes changes based on the state read
// and what is in the Output.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Name", request.Name)
	reqLogger.Info("Reconciling Output")

	// Fetch the Output instance
	instance := &loggingv1alpha1.Output{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if !errors.IsNotFound(err) {
			// Error reading object, requeue
			return reconcile.Result{}, err
		}
	}

	buff, err := getFluentdConfig(r.client)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update configmap for fluentd
	log.Info("Refreshing fluentd...")
	return reconcile.Result{}, r.fluentd.Refresh(buff)
}

func getFluentdConfig(cl client.Client) ([]byte, error) {
	// Simple algorithm to render all outputs once one changes. This lets us keep thing simple and write entire config
	// as one.
	instances := &loggingv1alpha1.OutputList{}
	lo := client.ListOptions{}

	err := cl.List(context.TODO(), instances, &lo)

	if err != nil {
		return []byte{}, err
	}

	// Source rendering is not configurable yet.
	renderers := []resources.Resource{
		resources.NewSystem(),
		resources.NewSource(),
	}

	for i := range instances.Items {
		renderers = append(renderers, resources.NewOutput(cl, &instances.Items[i]))
	}

	var buff []byte
	var newline bytes.Buffer
	fmt.Fprintf(&newline, "\n\n")
	for _, r := range renderers {
		out, err := r.Render()
		if err != nil {
			return []byte{}, err
		}
		buff = append(buff, out...)
		buff = append(buff, newline.Bytes()...)
	}

	return buff, nil
}
