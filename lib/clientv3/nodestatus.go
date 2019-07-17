// Copyright (c) 2019 Tigera, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clientv3

import (
	"context"

	apiv3 "github.com/projectcalico/libcalico-go/lib/apis/v3"
	cerrors "github.com/projectcalico/libcalico-go/lib/errors"
	"github.com/projectcalico/libcalico-go/lib/options"
	validator "github.com/projectcalico/libcalico-go/lib/validator/v3"
	"github.com/projectcalico/libcalico-go/lib/watch"
)

// NodeStatusInterface has methods to work with NodeStatus resources.
type NodeStatusInterface interface {
	Create(ctx context.Context, res *apiv3.NodeStatus, opts options.SetOptions) (*apiv3.NodeStatus, error)
	Update(ctx context.Context, res *apiv3.NodeStatus, opts options.SetOptions) (*apiv3.NodeStatus, error)
	Delete(ctx context.Context, name string, opts options.DeleteOptions) (*apiv3.NodeStatus, error)
	Get(ctx context.Context, name string, opts options.GetOptions) (*apiv3.NodeStatus, error)
	List(ctx context.Context, opts options.ListOptions) (*apiv3.NodeStatusList, error)
	Watch(ctx context.Context, opts options.ListOptions) (watch.Interface, error)
}

// nodeStatus implements NodeStatusInterface
type nodeStatus struct {
	client client
}

// Create takes the representation of a NodeStatus and creates it.
// Returns the stored representation of the NodeStatus, and an error
// if there is any.
func (r nodeStatus) Create(ctx context.Context, res *apiv3.NodeStatus, opts options.SetOptions) (*apiv3.NodeStatus, error) {
	if err := validator.Validate(res); err != nil {
		return nil, err
	}

	out, err := r.client.resources.Create(ctx, opts, apiv3.KindNodeStatus, res)
	if out != nil {
		return out.(*apiv3.NodeStatus), err
	}
	return nil, err
}

// Update takes the representation of a NodeStatus and updates it.
// Returns the stored representation of the NodeStatus, and an error
// if there is any.
func (r nodeStatus) Update(ctx context.Context, res *apiv3.NodeStatus, opts options.SetOptions) (*apiv3.NodeStatus, error) {
	if err := validator.Validate(res); err != nil {
		return nil, err
	}

	out, err := r.client.resources.Update(ctx, opts, apiv3.KindNodeStatus, res)
	if out != nil {
		return out.(*apiv3.NodeStatus), err
	}
	return nil, err
}

// Delete takes name of the NodeStatus and deletes it. Returns an
// error if one occurs.
func (r nodeStatus) Delete(ctx context.Context, name string, opts options.DeleteOptions) (*apiv3.NodeStatus, error) {
	out, err := r.client.resources.Delete(ctx, opts, apiv3.KindNodeStatus, noNamespace, name)
	if out != nil {
		return out.(*apiv3.NodeStatus), err
	}
	return nil, err
}

// Get takes name of the NodeStatus, and returns the corresponding
// NodeStatus object, and an error if there is any.
func (r nodeStatus) Get(ctx context.Context, name string, opts options.GetOptions) (*apiv3.NodeStatus, error) {
	out, err := r.client.resources.Get(ctx, opts, apiv3.KindNodeStatus, noNamespace, name)
	if out != nil {
		return out.(*apiv3.NodeStatus), err
	}
	return nil, err
}

// List returns the list of NodeStatus objects that match the supplied options.
func (r nodeStatus) List(ctx context.Context, opts options.ListOptions) (*apiv3.NodeStatusList, error) {
	res := &apiv3.NodeStatusList{}
	if err := r.client.resources.List(ctx, opts, apiv3.KindNodeStatus, apiv3.KindNodeStatusList, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Watch returns a watch.Interface that watches the NodeStatus that
// match the supplied options.
func (r nodeStatus) Watch(ctx context.Context, opts options.ListOptions) (watch.Interface, error) {
	return r.client.resources.Watch(ctx, opts, apiv3.KindNodeStatus, nil)
}
