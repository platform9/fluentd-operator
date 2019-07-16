package resources

import (
	"testing"

	"github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestValueFrom(t *testing.T) {
	obj := v1alpha1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "fake",
		},
	}

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake-secret",
			Namespace: "fake",
		},
		StringData: map[string]string{
			"fake-key": "fake-val",
		},
	}

	o := NewOutput(fake.NewFakeClient(&secret), &obj)

	assert.NotNil(t, o)

	vf := v1alpha1.ValueFrom{
		Name: "fake-secret",
		Key:  "fake-key",
	}

	val, err := o.getValueFrom(&vf)
	assert.Nil(t, err)
	assert.Equal(t, "fake-val", val)
}

func TestEsParams(t *testing.T) {
	obj := v1alpha1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "fake",
		},
		Spec: v1alpha1.OutputSpec{
			Type: "elasticsearch",
		},
	}

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake-secret",
			Namespace: "fake",
		},
		StringData: map[string]string{
			"user":     "fake-user",
			"password": "fake-password",
		},
	}

	o := NewOutput(fake.NewFakeClient(&secret), &obj)

	assert.NotNil(t, o)

	params, err := o.getEsParams()

	assert.Nil(t, err)

	keys := []string{"index_name", "host", "port", "scheme"}

	for _, k := range keys {
		_, ok := params[k]
		assert.True(t, ok)
	}
}

func TestEsRender(t *testing.T) {
	obj := v1alpha1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "fake",
		},
		Spec: v1alpha1.OutputSpec{
			Type: "elasticsearch",
			Params: []v1alpha1.Param{
				v1alpha1.Param{
					Name: "user",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "user",
					},
				},
				v1alpha1.Param{
					Name: "password",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "password",
					},
				},
			},
		},
	}

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake-secret",
			Namespace: "fake",
		},
		StringData: map[string]string{
			"user":     "fake-user",
			"password": "fake-password",
		},
	}

	o := NewOutput(fake.NewFakeClient(&secret), &obj)

	assert.NotNil(t, o)

	_, err := o.Render()

	assert.Nil(t, err)

}

func TestS3Params(t *testing.T) {
	obj := v1alpha1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "fake",
		},
		Spec: v1alpha1.OutputSpec{
			Type: "s3",
			Params: []v1alpha1.Param{
				v1alpha1.Param{
					Name: "aws_key_id",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "aws_key_id",
					},
				},
				v1alpha1.Param{
					Name: "aws_sec_key",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "aws_sec_key",
					},
				},
				v1alpha1.Param{
					Name: "s3_bucket",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "s3_bucket",
					},
				},
				v1alpha1.Param{
					Name: "s3_endpoint",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "s3_endpoint",
					},
				},
			},
		},
	}

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake-secret",
			Namespace: "fake",
		},
		StringData: map[string]string{
			"aws_key_id":  "fake-id",
			"aws_sec_key": "fake-secret",
			"s3_bucket":   "fake-bucket",
			"s3_endpoint": "fake-endpoint",
		},
	}

	o := NewOutput(fake.NewFakeClient(&secret), &obj)

	assert.NotNil(t, o)

	params, err := o.getS3Params()

	assert.Nil(t, err)

	keys := []string{"aws_key_id", "aws_sec_key", "s3_bucket", "s3_endpoint"}

	for _, k := range keys {
		_, ok := params[k]
		assert.True(t, ok)
	}
}

func TestS3Render(t *testing.T) {
	obj := v1alpha1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "fake",
		},
		Spec: v1alpha1.OutputSpec{
			Type: "s3",
			Params: []v1alpha1.Param{
				v1alpha1.Param{
					Name: "aws_key_id",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "aws_key_id",
					},
				},
				v1alpha1.Param{
					Name: "aws_sec_key",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "aws_sec_key",
					},
				},
				v1alpha1.Param{
					Name: "s3_bucket",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "s3_bucket",
					},
				},
				v1alpha1.Param{
					Name: "s3_endpoint",
					ValueFrom: &v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "s3_endpoint",
					},
				},
			},
		},
	}

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake-secret",
			Namespace: "fake",
		},
		StringData: map[string]string{
			"aws_key_id":  "fake-id",
			"aws_sec_key": "fake-secret",
			"s3_bucket":   "fake-bucket",
			"s3_endpoint": "fake-endpoint",
		},
	}

	o := NewOutput(fake.NewFakeClient(&secret), &obj)

	assert.NotNil(t, o)

	_, err := o.Render()

	assert.Nil(t, err)

}
