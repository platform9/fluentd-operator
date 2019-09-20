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

package v1alpha1

import (
	v1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	scheme "github.com/platform9/fluentd-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// OutputsGetter has a method to return a OutputInterface.
// A group's client should implement this interface.
type OutputsGetter interface {
	Outputs() OutputInterface
}

// OutputInterface has methods to work with Output resources.
type OutputInterface interface {
	Create(*v1alpha1.Output) (*v1alpha1.Output, error)
	Update(*v1alpha1.Output) (*v1alpha1.Output, error)
	UpdateStatus(*v1alpha1.Output) (*v1alpha1.Output, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Output, error)
	List(opts v1.ListOptions) (*v1alpha1.OutputList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Output, err error)
	OutputExpansion
}

// outputs implements OutputInterface
type outputs struct {
	client rest.Interface
}

// newOutputs returns a Outputs
func newOutputs(c *LoggingV1alpha1Client) *outputs {
	return &outputs{
		client: c.RESTClient(),
	}
}

// Get takes name of the output, and returns the corresponding output object, and an error if there is any.
func (c *outputs) Get(name string, options v1.GetOptions) (result *v1alpha1.Output, err error) {
	result = &v1alpha1.Output{}
	err = c.client.Get().
		Resource("outputs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Outputs that match those selectors.
func (c *outputs) List(opts v1.ListOptions) (result *v1alpha1.OutputList, err error) {
	result = &v1alpha1.OutputList{}
	err = c.client.Get().
		Resource("outputs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested outputs.
func (c *outputs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("outputs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a output and creates it.  Returns the server's representation of the output, and an error, if there is any.
func (c *outputs) Create(output *v1alpha1.Output) (result *v1alpha1.Output, err error) {
	result = &v1alpha1.Output{}
	err = c.client.Post().
		Resource("outputs").
		Body(output).
		Do().
		Into(result)
	return
}

// Update takes the representation of a output and updates it. Returns the server's representation of the output, and an error, if there is any.
func (c *outputs) Update(output *v1alpha1.Output) (result *v1alpha1.Output, err error) {
	result = &v1alpha1.Output{}
	err = c.client.Put().
		Resource("outputs").
		Name(output.Name).
		Body(output).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *outputs) UpdateStatus(output *v1alpha1.Output) (result *v1alpha1.Output, err error) {
	result = &v1alpha1.Output{}
	err = c.client.Put().
		Resource("outputs").
		Name(output.Name).
		SubResource("status").
		Body(output).
		Do().
		Into(result)
	return
}

// Delete takes name of the output and deletes it. Returns an error if one occurs.
func (c *outputs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("outputs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *outputs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("outputs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched output.
func (c *outputs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Output, err error) {
	result = &v1alpha1.Output{}
	err = c.client.Patch(pt).
		Resource("outputs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
