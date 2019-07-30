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
