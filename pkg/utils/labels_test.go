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
	"testing"

	"github.com/platform9/fluentd-operator/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestWithSameDicts(t *testing.T) {
	var a = map[string]string{"1": "", "2": "", "3": ""}
	var b = map[string]string{"1": "", "2": "", "3": ""}
	assert.True(t, utils.CheckSubset(a, b))
}

func TestWithDisjointDicts(t *testing.T) {
	var a = map[string]string{"1": "", "2": "", "3": "", "11": ""}
	var b = map[string]string{"1": "", "2": "", "3": "", "15": ""}
	assert.False(t, utils.CheckSubset(a, b))
}

func TestWithEmptyDicts(t *testing.T) {
	var a = map[string]string{}
	var b = map[string]string{"1": "", "2": "", "3": ""}
	assert.True(t, utils.CheckSubset(a, b))
}

func TestWithDiffValues(t *testing.T) {
	var a = map[string]string{"1": "", "2": "", "3": "true"}
	var b = map[string]string{"1": "", "2": "", "3": "false"}
	assert.False(t, utils.CheckSubset(a, b))
}
