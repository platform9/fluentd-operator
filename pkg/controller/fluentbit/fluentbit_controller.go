package fluentbit

import (
	"github.com/platform9/fluentd-operator/pkg/fluentbit"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_fluentbit")

// Add creates a new Fluentbit Controller and adds it to the Manager.
func Add(mgr manager.Manager) error {
	rc := newReconciler(mgr)
	return add(mgr, rc)
}

// newReconciler returns a new fluentbit.Reconciler
func newReconciler(mgr manager.Manager) *fluentbit.Reconciler {
	return fluentbit.New(mgr)
}

// add adds a new Controller to mgr
func add(mgr manager.Manager, r *fluentbit.Reconciler) error {
	// Create a new controller
	c, err := controller.New("fluentbit-controller", mgr, controller.Options{
		Reconciler: r,
	})

	if err != nil {
		return err
	}

	if err := r.CreateIfNeeded(); err != nil {
		log.Error(err, "Error creating fluent-bit")
		return err
	}

	return c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appsv1.DaemonSet{},
	})
}

var _ reconcile.Reconciler = &fluentbit.Reconciler{}
