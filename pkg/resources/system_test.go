package resources_test

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/platform9/fluentd-operator/pkg/options"
	"github.com/platform9/fluentd-operator/pkg/resources"
	"github.com/stretchr/testify/assert"
)

type TestSystem struct {
	XMLName  xml.Name `xml:"system"`
	Endpoint string   `xml:",innerxml"`
}

func TestSystemRender(t *testing.T) {
	s := resources.NewSystem()

	buf, err := s.Render()
	assert.Nil(t, err)

	var found TestSystem
	assert.Nil(t, xml.Unmarshal(buf, &found))

	rpcVal := strings.Trim(found.Endpoint, " \n")
	assert.Equal(t, fmt.Sprintf("rpc_endpoint 0.0.0.0:%d", *(options.ReloadPort)), rpcVal)

}
