package resources

import (
	"bytes"
	"fmt"

	"github.com/platform9/fluentd-operator/pkg/options"
)

// System represents implementation of fluentd System configuration.
type System struct {
	port int
}

// NewSystem returns a new System object
func NewSystem() *System {
	return &System{
		port: *(options.ReloadPort),
	}
}

// Render returns byte array representing fluentd configuration of a System
func (s *System) Render() ([]byte, error) {
	var ret bytes.Buffer
	fmt.Fprintf(&ret, "<system>")
	fmt.Fprintf(&ret, fmt.Sprintf("\n    rpc_endpoint 0.0.0.0:%d", s.port))
	fmt.Fprintf(&ret, "\n</system>")

	return ret.Bytes(), nil
}
