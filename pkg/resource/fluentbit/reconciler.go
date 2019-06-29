package fluentbit

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	"github.com/platform9/fluentd-operator/pkg/options"
	"github.com/platform9/fluentd-operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("fluentbit_reconciler")

var volumePaths = map[string]string{
	"varlog":                 "/var/log",
	"varlibdockercontainers": "/var/lib/docker/containers",
	"fluent-bit-config":      "/fluent-bit/etc/fluent-bit.conf",
}

var labels = map[string]string{
	"k8s-app":                       "fluent-bit",
	"kubernetes.io/cluster-service": "true",
	"created_by":                    "fluentd-operator",
}

// Reconciler reconciles fluentbit daemonset
type Reconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// New returns new instance of reconciler
func New(client client.Client, scheme *runtime.Scheme) *Reconciler {
	return &Reconciler{
		client: client,
		scheme: scheme,
	}
}

// Reconcile reads state of cluster for DaemonSet objects and makes
// changes per how the controller definition needs to be
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	// TODO: Watch daemonsets only in current namespace
	instance := &appsv1.DaemonSet{}

	// Only interested in configured namespace for fluentbit daemonset
	if request.Namespace != options.LogNs {
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

	if !utils.CheckSubset(labels, objLabels) {
		reqLogger.Info("Is not interesting")
		return reconcile.Result{}, nil
	}

	desired := getDesired()
	desired.ResourceVersion = current.ResourceVersion
	if err := r.client.Update(context.TODO(), desired); err != nil {
		reqLogger.Error(err, "Error updating")
		return reconcile.Result{},
			emperror.WrapWith(err, "update failed", "type", desired.GetObjectKind().GroupVersionKind())
	}

	return reconcile.Result{}, nil
}

// CreateIfNeeded creates fluentbit daemonset if needed
func (r *Reconciler) CreateIfNeeded() error {

	desired := getDesired()
	if err := r.client.Create(context.TODO(), desired); err != nil {
		if errors.IsAlreadyExists(err) {
			log.Info("fluentbit daemonset already exists, skipping...")
			return nil
		}

		return emperror.WrapWith(err, "create failed", "type", desired.GetObjectKind().GroupVersionKind())
	}

	return nil
}

// getDesired returns the desired resource spec for fluentbit
func getDesired() *appsv1.DaemonSet {
	annotations := map[string]string{
		"prometheus.io/scrape": "true",
		"prometheus.io/port":   "2020",
		"prometheus.io/path":   "/api/v1/metrics/prometheus",
	}

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluent-bit",
			Namespace: options.LogNs,
			Labels:    labels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: *getPodSpec(),
			},
		},
	}
}

func getPodSpec() *corev1.PodSpec {
	return &corev1.PodSpec{
		Tolerations: []corev1.Toleration{
			{
				Key:    "node-role.kubernetes.io/master",
				Effect: "NoSchedule",
			},
		},
		Containers: []corev1.Container{
			{
				Name:            "fluent-bit",
				Image:           "fluent/fluent-bit:1.0.6", // TODO: customize
				ImagePullPolicy: "IfNotPresent",
				Ports: []corev1.ContainerPort{{
					Name:          "prometheus", // TODO: customize
					ContainerPort: 2020,
				}},
				ReadinessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/api/v1/metrics/prometheus",
							Port: intstr.IntOrString{
								IntVal: 2020,
							},
						},
					},
				},
				LivenessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/api/v1/metrics/prometheus",
							Port: intstr.IntOrString{
								IntVal: 2020,
							},
						},
					},
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						"cpu":    resource.MustParse("5m"),
						"memory": resource.MustParse("10Mi"),
					},
					Limits: corev1.ResourceList{
						"cpu":    resource.MustParse("50m"),
						"memory": resource.MustParse("60Mi"),
					},
				},
				VolumeMounts: getVolumeMounts(),
			},
		},
		Volumes:            getVolumes(),
		ServiceAccountName: options.SvcAcct,
	}
}

func getVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "varlog",
			MountPath: volumePaths["varlog"],
			ReadOnly:  true,
		},
		{
			Name:      "varlibdockercontainers",
			MountPath: volumePaths["varlibdockercontainers"],
			ReadOnly:  true,
		},
		{
			Name:      "fluent-bit-config",
			MountPath: volumePaths["fluent-bit-config"],
		},
	}
}

func getVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "varlog",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: volumePaths["varlog"],
				},
			},
		},
		{
			Name: "varlibdockercontainers",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: volumePaths["varlibdockercontainers"],
				},
			},
		},
		{
			Name: "fluent-bit-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "fluent-bit-config",
					},
				},
			},
		},
	}
}
