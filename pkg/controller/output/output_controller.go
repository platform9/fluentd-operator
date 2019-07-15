package output

import (
	"bytes"
	"context"
	"fmt"

	"github.com/platform9/fluentd-operator/pkg/fluentd"

	"github.com/platform9/fluentd-operator/pkg/resources"

	loggingv1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_output")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Output Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileOutput{
		client:  mgr.GetClient(),
		scheme:  mgr.GetScheme(),
		fluentd: fluentd.New(mgr),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("output-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Output
	err = c.Watch(&source.Kind{Type: &loggingv1alpha1.Output{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Output
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &loggingv1alpha1.Output{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileOutput implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileOutput{}

// ReconcileOutput reconciles a Output object
type ReconcileOutput struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client  client.Client
	scheme  *runtime.Scheme
	fluentd *fluentd.Reconciler
}

// Reconcile reads that state of the cluster for a Output object and makes changes based on the state read
// and what is in the Output.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileOutput) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Name", request.Name)
	reqLogger.Info("Reconciling Output")

	// Fetch the Output instance
	instance := &loggingv1alpha1.Output{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	renderers := []resources.Resource{
		resources.NewSource(), resources.NewOutput(r.client, instance),
	}

	var buff []byte
	var newline bytes.Buffer
	fmt.Fprintf(&newline, "\n\n")
	for _, r := range renderers {
		out, err := r.Render()
		if err != nil {
			return reconcile.Result{}, err
		}
		buff = append(buff, out...)
		buff = append(buff, newline.Bytes()...)
	}

	// Update configmap for fluentd
	log.Info("Refreshing fluentd...")
	return reconcile.Result{}, r.fluentd.Refresh(buff)
}
