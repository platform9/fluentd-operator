package resource

import "sigs.k8s.io/controller-runtime/pkg/reconcile"
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Reconciler defines interface for reconciling resources controlled by fluentd-operator
type Reconciler interface {
	Reconcile(obj metav1.Object) (reconcile.Result, error)
}
