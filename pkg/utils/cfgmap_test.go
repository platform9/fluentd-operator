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
