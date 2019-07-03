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
