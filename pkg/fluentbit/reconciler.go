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

package fluentbit

import (
	"context"

	"k8s.io/client-go/tools/record"

	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/go-logr/logr"
	fbsyncer "github.com/platform9/fluentd-operator/pkg/fluentbit/internal/syncer"
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

var log = logf.Log.WithName("fluentbit_reconciler")

// Reconciler reconciles fluentbit daemonset
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
		recorder: mgr.GetEventRecorderFor("controller.fluentbit"),
	}
}

// Reconcile reads state of cluster for DaemonSet objects and makes
// changes per how the controller definition needs to be
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	// TODO: Watch daemonsets only in current namespace
	instance := &appsv1.DaemonSet{}

	// Only interested in configured namespace for fluentbit daemonset
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

// Reconcile ensures fluentbit deployment is according to operator definition
func (r *Reconciler) reconcile(reqLogger logr.Logger, current *appsv1.DaemonSet) (reconcile.Result, error) {
	// Compare labels
	objLabels := current.GetLabels()

	if !utils.CheckSubset(fbsyncer.Labels, objLabels) {
		reqLogger.Info("Is not interesting")
		return reconcile.Result{}, nil
	}

	syncers := []syncer.Interface{fbsyncer.NewFluentbitSyncer(r.client, r.scheme),
		fbsyncer.NewFluentbitCfgMapSyncer(r.client, r.scheme),
	}

	for _, sync := range syncers {
		if err := syncer.Sync(context.TODO(), sync, r.recorder); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// CreateIfNeeded creates fluentbit daemonset if needed
func (r *Reconciler) CreateIfNeeded() error {
	syncers := []syncer.Interface{fbsyncer.NewFluentbitSyncer(r.client, r.scheme),
		fbsyncer.NewFluentbitCfgMapSyncer(r.client, r.scheme),
	}

	for _, sync := range syncers {
		if err := syncer.Sync(context.TODO(), sync, r.recorder); err != nil {
			if errors.IsAlreadyExists(err) {
				log.Info("fluentbit object already exists, skipping...")
				continue
			}
			return err
		}
	}
	return nil
}
