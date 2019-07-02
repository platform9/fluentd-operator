package syncer

import (
	"fmt"
	"io/ioutil"

	"github.com/platform9/fluentd-operator/pkg/options"
	"github.com/presslabs/controller-util/syncer"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	cfgMapName = "fluent-bit-config"
)

// Labels defines operator enforced labels for fluentbit daemonset
var Labels = map[string]string{
	"k8s-app":                       "fluent-bit",
	"kubernetes.io/cluster-service": "true",
	"created_by":                    "fluentd-operator",
}

var volumePaths = map[string]string{
	"varlog":                 "/var/log",
	"varlibdockercontainers": "/var/lib/docker/containers",
	cfgMapName:               "/fluent-bit/etc/",
}

type fbSyncer struct {
}

type fbCfgMapSyncer struct {
}

// NewFluentbitSyncer returns a sync interface compliant implementation for fluentbit
func NewFluentbitSyncer(c client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluent-bit",
			Namespace: *(options.LogNs),
		},
	}

	sync := &fbSyncer{}

	return syncer.NewObjectSyncer("DaemonSet", nil, obj, c, scheme, sync.SyncFn)
}

// NewFluentbitCfgMapSyncer returns a sync interface compliant implementation for fluentbit configmap
func NewFluentbitCfgMapSyncer(c client.Client, scheme *runtime.Scheme) syncer.Interface {
	obj := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfgMapName,
			Namespace: *(options.LogNs),
		},
	}

	sync := &fbCfgMapSyncer{}

	return syncer.NewObjectSyncer("ConfigMap", nil, obj, c, scheme, sync.SyncFn)
}

// SyncFn syncs the Fluentbit config map per spec
func (s *fbCfgMapSyncer) SyncFn(in runtime.Object) error {
	out := in.(*corev1.ConfigMap)
	out.ObjectMeta.Labels = Labels
	if len(out.Data) == 0 {
		if d, err := getCfgMapData(); err != nil {
			return err
		} else {
			out.Data = d
		}
	}
	return nil
}

// SyncFn syncs the Fluentbit cluster object with operator spec
func (s *fbSyncer) SyncFn(in runtime.Object) error {
	annotations := map[string]string{
		"prometheus.io/scrape": "true",
		"prometheus.io/port":   "2020",
		"prometheus.io/path":   "/api/v1/metrics/prometheus",
	}

	out := in.(*appsv1.DaemonSet)

	out.ObjectMeta.Labels = Labels
	out.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: Labels,
	}
	out.Spec.Template.ObjectMeta.Annotations = annotations
	out.Spec.Template.ObjectMeta.Labels = Labels
	out.Spec.Template.Spec = *getPodSpec()
	return nil
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
				Image:           *(options.FluentbitImage),
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
		ServiceAccountName: *(options.SvcAcct),
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
			Name:      cfgMapName,
			MountPath: volumePaths[cfgMapName],
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

func getCfgMapData() (map[string]string, error) {
	subDir := fmt.Sprintf("%s/conf/fluent-bit", *(options.CfgDir))
	data := map[string]string{}
	files, err := ioutil.ReadDir(subDir)
	if err != nil {
		return data, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", subDir, file.Name()))
		if err != nil {
			return data, err
		}

		data[file.Name()] = string(content)
	}

	return data, nil
}
