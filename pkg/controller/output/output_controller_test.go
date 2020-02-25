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
	"context"
	"encoding/json"
	"testing"

	loggingv1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TestClient struct {
	sw client.StatusWriter
}

type TestStatusWriter struct {
}

func NewTestClient() client.Client {
	return &TestClient{}
}

func (t *TestClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	return nil
}

func (t *TestClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return nil
}

func (t *TestClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	return nil
}

func (t *TestClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	return nil
}

func (t *TestClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	return nil
}

func (t *TestClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}

func (t *TestClient) List(ctx context.Context, list runtime.Object, opt ...client.ListOption) error {
	decoder := scheme.Codecs.UniversalDecoder()
	var objs runtime.Object
	objs = &loggingv1alpha1.OutputList{
		Items: []loggingv1alpha1.Output{
			loggingv1alpha1.Output{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ES-Output",
				},
				Spec: loggingv1alpha1.OutputSpec{
					Type: "elasticsearch",
					Params: []loggingv1alpha1.Param{
						loggingv1alpha1.Param{
							Name:  "user",
							Value: "fake-user",
						},
						loggingv1alpha1.Param{
							Name:  "password",
							Value: "fake-password",
						},
					},
				},
			},
		},
	}

	j, err := json.Marshal(objs)
	if err != nil {
		return err
	}

	decoder.Decode(j, nil, list)
	return nil
}

func (t *TestClient) Status() client.StatusWriter {
	return t.sw
}

func TestFluentdConfig(t *testing.T) {
	cl := NewTestClient()
	buf, err := getFluentdConfig(cl)
	t.Log(err)
	assert.Nil(t, err)
	assert.NotEmpty(t, buf)
}
