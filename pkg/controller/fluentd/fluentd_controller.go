/*
Copyright 2019 Platform9 Systems, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
