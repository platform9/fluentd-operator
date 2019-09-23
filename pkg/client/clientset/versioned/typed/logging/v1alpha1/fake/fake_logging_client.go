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

package fake

import (
	v1alpha1 "github.com/platform9/fluentd-operator/pkg/client/clientset/versioned/typed/logging/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeLoggingV1alpha1 struct {
	*testing.Fake
}

func (c *FakeLoggingV1alpha1) Outputs() v1alpha1.OutputInterface {
	return &FakeOutputs{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeLoggingV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
