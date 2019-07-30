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

package syncer_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/platform9/fluentd-operator/pkg/fluentd/internal/syncer"
	"github.com/platform9/fluentd-operator/pkg/options"
	"github.com/stretchr/testify/assert"
	api_rt "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMain(t *testing.T) {
	_, fn, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	*(options.CfgDir) = filepath.Join(filepath.Dir(fn), "../../../../etc/conf")
}

func TestSvcSyncer(t *testing.T) {
	s := syncer.NewFluentdSvcSyncer(fake.NewFakeClient(), &api_rt.Scheme{})
	_, err := s.Sync(context.TODO())
	assert.Nil(t, err)
}

func TestCfgMapSyncer(t *testing.T) {
	c := syncer.NewFluentdCfgMapSyncer(fake.NewFakeClient(), &api_rt.Scheme{})
	_, err := c.Sync(context.TODO())
	assert.Nil(t, err)
}

func TestFluentdSyncer(t *testing.T) {
	f := syncer.NewFluentdSyncer(fake.NewFakeClient(), &api_rt.Scheme{})
	_, err := f.Sync(context.TODO())
	assert.Nil(t, err)
}
