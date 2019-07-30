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

package utils_test

import (
	"os"
	"testing"

	"github.com/platform9/fluentd-operator/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	err := os.Chdir("../../")
	assert.Nil(t, err)
}

func getKeys(d map[string][]byte) []string {
	keys := []string{}

	for k := range d {
		keys = append(keys, k)
	}

	return keys
}
func TestGetConfigForFluentbit(t *testing.T) {
	d, err := utils.GetCfgMapData("fluent-bit")
	assert.Nil(t, err)
	assert.Equal(t, 6, len(d))
	keys := getKeys(d)
	assert.ElementsMatch(t, []string{"filter.conf", "fluent-bit.conf", "input.conf", "output.conf",
		"parsers.conf", "null.conf"}, keys)
}

func TestGetConfigForFluentd(t *testing.T) {
	d, err := utils.GetCfgMapData("fluentd")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(d))

	keys := getKeys(d)
	assert.ElementsMatch(t, []string{"fluent.conf"}, keys)
}
