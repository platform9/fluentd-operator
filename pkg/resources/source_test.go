package resources_test

import (
	"encoding/xml"
	"testing"

	"github.com/platform9/fluentd-operator/pkg/resources"
	"github.com/stretchr/testify/assert"
)

type TestSource struct {
	XMLName xml.Name `xml:"source"`
	Data    string   `xml:",innerxml"`
}

func TestSourceRender(t *testing.T) {
	s := resources.NewSource()

	buf, err := s.Render()
	assert.Nil(t, err)

	var found TestSource
	assert.Nil(t, xml.Unmarshal(buf, &found))

}
