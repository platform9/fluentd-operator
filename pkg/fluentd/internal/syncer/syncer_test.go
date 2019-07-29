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
