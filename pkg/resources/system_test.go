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
