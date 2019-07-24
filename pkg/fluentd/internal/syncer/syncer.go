package syncer

import (
	"github.com/imdario/mergo"
	"github.com/platform9/fluentd-operator/pkg/options"
	"github.com/platform9/fluentd-operator/pkg/utils"
	"github.com/presslabs/controller-util/mergo/transformers"
	"github.com/presslabs/controller-util/syncer"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("fluentd_syncer")

const (
	cfgMapName = "fluentd-config"
	svcName    = "fluentd"
)

var volumePaths = map[string]string{
	cfgMapName: "/fluentd/etc/",
}

// Labels defines operator enforced labels for fluentd deployment
var Labels = map[string]string{
	"k8s-app":                       "fluentd",
	"kubernetes.io/cluster-service": "true",
	"created_by":                    "fluentd-operator",
}

type fdSyncer struct {
}

type fdCfgMapSyncer struct {
	data []byte
}

type fdSvcSyncer struct {
}

func getLabels() labels.Set {
	return Labels
}

// NewFluentdSyncer returns a sync interface compliant implementation for fluentd
func NewFluentdSyncer(c client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd",
			Namespace: *(options.LogNs),
		},
	}

	sync := &fdSyncer{}

	return syncer.NewObjectSyncer("Deployment", nil, obj, c, scheme, sync.SyncFn)
}

// NewFluentdCfgMapSyncer returns a sync interface compliant implementation for fluentd configmap
func NewFluentdCfgMapSyncer(c client.Client, scheme *runtime.Scheme, params ...[]byte) syncer.Interface {
	obj := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfgMapName,
			Namespace: *(options.LogNs),
		},
	}

	sync := &fdCfgMapSyncer{}

	if len(params) > 0 {
		sync.data = params[0]
	}

	return syncer.NewObjectSyncer("ConfigMap", nil, obj, c, scheme, sync.SyncFn)
}

// NewFluentdSvcSyncer returns a sync interface compliant implementation for fluentd service
func NewFluentdSvcSyncer(c client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: svcName,
			Namespace: *(options.LogNs),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "forwarder",
					Port:       62073,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(62073),
				},
				corev1.ServicePort{
					Name:       "webhook",
					Port:       45550,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(45550),
				},
			},
		},
	}

	sync := &fdSvcSyncer{}

	return syncer.NewObjectSyncer("Service", nil, obj, c, scheme, sync.SyncFn)
}

// SyncFn sync the Fluentd service per spec
func (s *fdSvcSyncer) SyncFn(in runtime.Object) error {
	out := in.(*corev1.Service)
	if len(out.ObjectMeta.Labels) == 0 {
		out.ObjectMeta.Labels = map[string]string{}
	}

	if len(out.Spec.Selector) == 0 {
		out.Spec.Selector = map[string]string{}
	}
	for k, v := range Labels {
		out.ObjectMeta.Labels[k] = v
		out.Spec.Selector[k] = v
	}

	return nil
}

// SyncFn syncs the Fluentd config map per spec
func (s *fdCfgMapSyncer) SyncFn(in runtime.Object) error {
	out := in.(*corev1.ConfigMap)
	if len(out.ObjectMeta.Labels) == 0 {
		out.ObjectMeta.Labels = map[string]string{}
	}

	for k, v := range Labels {
		out.ObjectMeta.Labels[k] = v
	}

	if len(s.data) > 0 {
		out.BinaryData = map[string][]byte{
			"fluent.conf": s.data,
		}
	} else if len(out.BinaryData) == 0 {
		if d, err := utils.GetCfgMapData("fluentd"); err != nil {
			return err
		} else {
			out.BinaryData = d
		}
	}

	return nil
}

// SyncFn syncs the Fluentd cluster object with operator spec
func (s *fdSyncer) SyncFn(in runtime.Object) error {
	annotations := map[string]string{
		"prometheus.io/scrape": "true",
		"prometheus.io/port":   "2020",
		"prometheus.io/path":   "/api/v1/metrics/prometheus",
	}

	out := in.(*appsv1.Deployment)

	var replicas int32 = 1
	out.ObjectMeta.Labels = Labels
	out.Spec.Replicas = &replicas // TODO: Use HPA
	out.Spec.Selector = metav1.SetAsLabelSelector(getLabels())

	if len(out.Spec.Template.ObjectMeta.Annotations) == 0 {
		out.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}

	for k, v := range annotations {
		out.Spec.Template.ObjectMeta.Annotations[k] = v
	}

	if len(out.Spec.Template.ObjectMeta.Labels) == 0 {
		out.Spec.Template.ObjectMeta.Labels = map[string]string{}
	}

	for k, v := range Labels {
		out.Spec.Template.ObjectMeta.Labels[k] = v
	}

	return mergo.Merge(&out.Spec.Template.Spec, getPodSpec(), mergo.WithTransformers(transformers.PodSpec))
}

func getPodSpec() corev1.PodSpec {
	return corev1.PodSpec{
		Tolerations: []corev1.Toleration{
			{
				Key:    "node-role.kubernetes.io/master",
				Effect: "NoSchedule",
			},
		},
		Containers: []corev1.Container{
			{
				Name:            "fluentd",
				Image:           *(options.FluentdImage),
				ImagePullPolicy: "IfNotPresent",
				Ports: []corev1.ContainerPort{{
					Name:          "prometheus", // TODO: customize
					ContainerPort: 2020,
				}, {
					Name:          "source",
					ContainerPort: 62073,
				}},
				Env: []corev1.EnvVar{{
					Name:  "FLUENT_ELASTICSEARCH_SED_DISABLE",
					Value: "1",
				}},
				// TODO: Customize
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						"cpu":    resource.MustParse("100m"),
						"memory": resource.MustParse("200Mi"),
					},
					Limits: corev1.ResourceList{
						"cpu":    resource.MustParse("1"),
						"memory": resource.MustParse("200Mi"),
					},
				},
				VolumeMounts: getVolumeMounts(),
			},
		},
		Volumes:            getVolumes(),
		ServiceAccountName: *(options.SvcAcct),
	}
}

func getVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      cfgMapName,
			MountPath: volumePaths[cfgMapName],
			ReadOnly:  true,
		},
	}
}

func getVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: cfgMapName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cfgMapName,
					},
				},
			},
		},
	}
}
