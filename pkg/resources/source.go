package resources

import (
	"bytes"
	"fmt"

	"github.com/platform9/fluentd-operator/pkg/options"
)

// Source represents implementation of fluentd source configuration.
type Source struct {
	port int
}

// NewSource returns a new source object
func NewSource() *Source {
	return &Source{
		port: *(options.ForwardPort),
	}
}

// Render returns byte array representing fluentd configuration of a source
func (s *Source) Render() ([]byte, error) {
	var ret bytes.Buffer
	fmt.Fprintf(&ret, "<source>")
	fmt.Fprintf(&ret, "\n    @type forward")
	fmt.Fprintf(&ret, "\n    port %d", s.port)
	fmt.Fprintf(&ret, "\n    bind 0.0.0.0")
	fmt.Fprintf(&ret, "\n</source>")

	return ret.Bytes(), nil
}
