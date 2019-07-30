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

package fluentd

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/platform9/fluentd-operator/pkg/options"
	"github.com/stretchr/testify/assert"
	api_rt "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type FakeManager struct {
}

func (fm *FakeManager) GetClient() client.Client {
	return fake.NewFakeClient()
}

func (fm *FakeManager) GetScheme() *api_rt.Scheme {
	return &api_rt.Scheme{}
}

func (fm *FakeManager) GetRecorder(name string) record.EventRecorder {
	return record.NewFakeRecorder(128)
}
func TestMain(t *testing.T) {
	_, fn, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	*(options.CfgDir) = filepath.Join(filepath.Dir(fn), "../../etc/conf")
}

func TestCreate(t *testing.T) {
	err := createIfNeeded(fake.NewFakeClient(), &api_rt.Scheme{}, record.NewFakeRecorder(128))

	assert.Nil(t, err)
}

func TestReconcile(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/api/config.reload", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	ts := httptest.NewServer(r)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.Nil(t, err)

	*(options.ReloadHost) = u.Hostname()
	*(options.ReloadPort), err = strconv.Atoi(u.Port())

	assert.Nil(t, err)

	var data []byte
	err = refresh(fake.NewFakeClient(), &api_rt.Scheme{}, record.NewFakeRecorder(128), data)

	assert.Nil(t, err)
}
