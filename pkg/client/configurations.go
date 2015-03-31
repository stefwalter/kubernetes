/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package client

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

type ConfigurationsInterface interface {
	Configurations() ConfigurationInterface
}

type ConfigurationInterface interface {
	Get(name string) (*api.Configuration, error)
	Create(configuration *api.Configuration) (*api.Configuration, error)
	List(label labels.Selector, field fields.Selector) (*api.ConfigurationList, error)
	Delete(name string) error
	Update(configuration *api.Configuration) (*api.Configuration, error)
	Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error)
}

type configurations struct {
	r *Client
}

func newConfigurations(c *Client) *configurations {
	return &configurations{c}
}

func (conf *configurations) Create(configuration *api.Configuration) (*api.Configuration, error) {
	result := &api.Configuration{}
	err := conf.r.Post().Resource("configurations").Body(configuration).Do().Into(result)
	return result, err
}

func (conf *configurations) List(label labels.Selector, field fields.Selector) (*api.ConfigurationList, error) {
	result := &api.ConfigurationList{}
	err := conf.r.Get().Resource("configurations").
		LabelsSelectorParam(label).
		FieldsSelectorParam(field).
		Do().
		Into(result)

	return result, err
}

func (conf *configurations) Get(name string) (*api.Configuration, error) {
	result := &api.Configuration{}
	err := conf.r.Get().
		Resource("configurations").
		Name(name).
		Do().
		Into(result)

	return result, err
}

func (conf *configurations) Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return conf.r.Get().
		Prefix("watch").
		Resource("configurations").
		Param("resourceVersion", resourceVersion).
		LabelsSelectorParam(label).
		FieldsSelectorParam(field).
		Watch()
}

func (conf *configurations) Delete(name string) error {
	return conf.r.Delete().Resource("configurations").Name(name).Do().Error()
}

func (conf *configurations) Update(configuration *api.Configuration) (*api.Configuration, error) {
	result := &api.Configuration{}
	err := conf.r.Put().
		Resource("configurations").
		Name(configuration.Name).
		Body(configuration).
		Do().
		Into(result)

	return result, err
}
