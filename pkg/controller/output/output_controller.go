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

package output

import (
	"bytes"
	"context"
	"fmt"

	"github.com/platform9/fluentd-operator/pkg/fluentd"

	loggingv1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	"github.com/platform9/fluentd-operator/pkg/resources"
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
		log.Error(err, "Error adding watch")
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

	err := cl.List(context.TODO(), &lo, instances)

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
