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
	v1alpha1 "github.com/platform9/fluentd-operator/pkg/apis/logging/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeOutputs implements OutputInterface
type FakeOutputs struct {
	Fake *FakeLoggingV1alpha1
}

var outputsResource = schema.GroupVersionResource{Group: "logging.pf9.io", Version: "v1alpha1", Resource: "outputs"}

var outputsKind = schema.GroupVersionKind{Group: "logging.pf9.io", Version: "v1alpha1", Kind: "Output"}

// Get takes name of the output, and returns the corresponding output object, and an error if there is any.
func (c *FakeOutputs) Get(name string, options v1.GetOptions) (result *v1alpha1.Output, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(outputsResource, name), &v1alpha1.Output{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Output), err
}

// List takes label and field selectors, and returns the list of Outputs that match those selectors.
func (c *FakeOutputs) List(opts v1.ListOptions) (result *v1alpha1.OutputList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(outputsResource, outputsKind, opts), &v1alpha1.OutputList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.OutputList{ListMeta: obj.(*v1alpha1.OutputList).ListMeta}
	for _, item := range obj.(*v1alpha1.OutputList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested outputs.
func (c *FakeOutputs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(outputsResource, opts))
}

// Create takes the representation of a output and creates it.  Returns the server's representation of the output, and an error, if there is any.
func (c *FakeOutputs) Create(output *v1alpha1.Output) (result *v1alpha1.Output, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(outputsResource, output), &v1alpha1.Output{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Output), err
}

// Update takes the representation of a output and updates it. Returns the server's representation of the output, and an error, if there is any.
func (c *FakeOutputs) Update(output *v1alpha1.Output) (result *v1alpha1.Output, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(outputsResource, output), &v1alpha1.Output{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Output), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeOutputs) UpdateStatus(output *v1alpha1.Output) (*v1alpha1.Output, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(outputsResource, "status", output), &v1alpha1.Output{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Output), err
}

// Delete takes name of the output and deletes it. Returns an error if one occurs.
func (c *FakeOutputs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(outputsResource, name), &v1alpha1.Output{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeOutputs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(outputsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.OutputList{})
	return err
}

// Patch applies the patch and returns the patched output.
func (c *FakeOutputs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Output, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(outputsResource, name, pt, data, subresources...), &v1alpha1.Output{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Output), err
}
