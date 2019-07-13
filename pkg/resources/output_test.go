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

	o := NewOutput(fake.NewFakeClient(&secret), obj)

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

	o := NewOutput(fake.NewFakeClient(&secret), obj)

	assert.NotNil(t, o)

	params, err := o.getEsParams()

	assert.Nil(t, err)

	keys := []string{"index_name", "host", "port", "scheme"}

	for _, k := range keys {
		_, ok := params[k]
		assert.True(t, ok)
	}
}
