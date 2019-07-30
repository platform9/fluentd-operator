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
		Data: map[string][]byte{
			"fake-key": []byte("fake-val"),
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
		Data: map[string][]byte{
			"user":     []byte("fake-user"),
			"password": []byte("fake-password"),
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
					ValueFrom: v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "user",
					},
				},
				v1alpha1.Param{
					Name: "password",
					ValueFrom: v1alpha1.ValueFrom{
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
		Data: map[string][]byte{
			"user":     []byte("fake-user"),
			"password": []byte("fake-password"),
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
					ValueFrom: v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "aws_key_id",
					},
				},
				v1alpha1.Param{
					Name: "aws_sec_key",
					ValueFrom: v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "aws_sec_key",
					},
				},
				v1alpha1.Param{
					Name: "s3_bucket",
					ValueFrom: v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "s3_bucket",
					},
				},
				v1alpha1.Param{
					Name: "s3_region",
					ValueFrom: v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "s3_region",
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
		Data: map[string][]byte{
			"aws_key_id":  []byte("fake-id"),
			"aws_sec_key": []byte("fake-secret"),
			"s3_bucket":   []byte("fake-bucket"),
			"s3_region":   []byte("fake-region"),
		},
	}

	o := NewOutput(fake.NewFakeClient(&secret), &obj)

	assert.NotNil(t, o)

	params, err := o.getS3Params()

	assert.Nil(t, err)

	keys := []string{"aws_key_id", "aws_sec_key", "s3_bucket", "s3_region"}

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
					ValueFrom: v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "aws_key_id",
					},
				},
				v1alpha1.Param{
					Name: "aws_sec_key",
					ValueFrom: v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "aws_sec_key",
					},
				},
				v1alpha1.Param{
					Name: "s3_bucket",
					ValueFrom: v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "s3_bucket",
					},
				},
				v1alpha1.Param{
					Name: "s3_region",
					ValueFrom: v1alpha1.ValueFrom{
						Name: "fake-secret",
						Key:  "s3_region",
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
		Data: map[string][]byte{
			"aws_key_id":  []byte("fake-id"),
			"aws_sec_key": []byte("fake-secret"),
			"s3_bucket":   []byte("fake-bucket"),
			"s3_region":   []byte("fake-endpoint"),
		},
	}

	o := NewOutput(fake.NewFakeClient(&secret), &obj)

	assert.NotNil(t, o)

	_, err := o.Render()

	assert.Nil(t, err)

}
