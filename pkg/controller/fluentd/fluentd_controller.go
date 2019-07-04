package fluentd

import (
	"github.com/platform9/fluentd-operator/pkg/fluentd"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_fluentd")

// Add creates a new Fluentd Controller and adds it to the Manager.
func Add(mgr manager.Manager) error {
	rc := newReconciler(mgr)
	return add(mgr, rc)
}

// newReconciler returns a new fluentbit.Reconciler
func newReconciler(mgr manager.Manager) *fluentd.Reconciler {
	return fluentd.New(mgr)
}

// add adds a new Controller to mgr
func add(mgr manager.Manager, r *fluentd.Reconciler) error {
	// Create a new controller
	c, err := controller.New("fluentd-controller", mgr, controller.Options{
		Reconciler: r,
	})

	if err != nil {
		return err
	}

	if err := r.CreateIfNeeded(); err != nil {
		log.Error(err, "Error creating fluentd")
		return err
	}

	if err := c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}); err != nil {
		log.Error(err, "Error adding watch")
		return err
	}

	if err := c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForObject{}); err != nil {
		log.Error(err, "Error adding watch")
		return err
	}

	return c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &fluentd.Reconciler{}
