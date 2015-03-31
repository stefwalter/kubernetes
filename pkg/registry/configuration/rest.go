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

package config

import (
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/validation"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

// REST provides a config registry for accessing apiserver's RESTStorage model.
type REST struct {
	registry generic.Registry
}

// NewStorage returns a new REST. You must use a registry created by
// NewEtcdRegistry unless you're testing.
func NewStorage(registry generic.Registry) *REST {
	return &REST{
		registry: registry,
	}
}

func (*REST) New() runtime.Object {
	return &api.Configuration{}
}

func (*REST) NewList() runtime.Object {
	return &api.ConfigurationList{}
}
func (rs *REST) Create(ctx api.Context, obj runtime.Object) (runtime.Object, error) {
	config, ok := obj.(*api.Configuration)
	if !ok {
		return nil, fmt.Errorf("invalid object type")
	}

	if config.Namespace != "" {
		return nil, false, errors.NewConflict("configuration", config.Namespace, fmt.Errorf("Configuration objects should not have a namespace."))
	}

	if len(config.Name) == 0 {
		return nil, false, errors.NewConflict("configuration", config.Name, fmt.Errorf("Configuration objects must have a name."))
	}

	if errs := validation.ValidateConfiguration(config); len(errs) > 0 {
		return nil, errors.NewInvalid("configuration", config.Name, errs)
	}
	api.FillObjectMetaSystemFields(ctx, &config.ObjectMeta)

	err := rs.registry.CreateWithName(ctx, config.Name, config)
	if err != nil {
		return nil, err
	}
	return rs.registry.Get(ctx, config.Name)
}

func (rs *REST) Update(ctx api.Context, obj runtime.Object) (runtime.Object, bool, error) {
	config, ok := obj.(*api.Config)
	if !ok {
		return nil, false, fmt.Errorf("invalid configuration object: %#v", obj)
	}

	if config.Namespace != "" {
		return nil, false, errors.NewConflict("configuration", config.Namespace, fmt.Errorf("Configuration objects should not have a namespace."))
	}

	oldObj, err := rs.registry.Get(ctx, config.Name)
	if err != nil {
		return nil, false, err
	}

	editConfig := oldObj.(*api.Configuration)

	// set the editable fields on the existing object
	editConfig.Labels = config.Labels
	editConfig.ResourceVersion = config.ResourceVersion
	editConfig.Annotations = config.Annotations
	editConfig.Flags = config.Flags

	if errs := validation.ValidateConfiguration(editConfig); len(errs) > 0 {
		return nil, false, errors.NewInvalid("configuration", editConfig.Name, errs)
	}

	err = rs.registry.UpdateWithName(ctx, editConfig.Name, editConfig)
	if err != nil {
		return nil, false, err
	}
	out, err := rs.registry.Get(ctx, editConfig.Name)
	return out, false, err
}

func (rs *REST) Delete(ctx api.Context, name string) (runtime.Object, error) {
	obj, err := rs.registry.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	_, ok := obj.(*api.Configuration)
	if !ok {
		return nil, fmt.Errorf("invalid object type")
	}

	return rs.registry.Delete(ctx, name, nil)
}

func (rs *REST) Get(ctx api.Context, name string) (runtime.Object, error) {
	obj, err := rs.registry.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	config, ok := obj.(*api.Config)
	if !ok {
		return nil, fmt.Errorf("invalid object type")
	}
	return config, err
}

func (rs *REST) getAttrs(obj runtime.Object) (objLabels labels.Set, objFields fields.Set, err error) {
	config, ok := obj.(*api.Configuration)
	if !ok {
		return nil, nil, fmt.Errorf("invalid object type")
	}

	return labels.Set{}, fields.Set{
		"type": string(config.Type),
	}, nil
}

func (rs *REST) List(ctx api.Context, label labels.Selector, field fields.Selector) (runtime.Object, error) {
	return rs.registry.ListPredicate(ctx, &generic.SelectionPredicate{label, field, rs.getAttrs})
}

func (rs *REST) Watch(ctx api.Context, label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return rs.registry.WatchPredicate(ctx, &generic.SelectionPredicate{label, field, rs.getAttrs}, resourceVersion)
}
