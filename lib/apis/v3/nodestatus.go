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

package v3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KindNodeStatus     = "NodeStatus"
	KindNodeStatusList = "NodeStatusList"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeStatus contains the configuration for any BGP routing.
type NodeStatus struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the NodeStatus.
	Spec NodeStatusSpec `json:"spec,omitempty"`
}

// NodeStatusSpec contains the values of the BGP configuration.
type NodeStatusSpec struct {
	// Status is the status is the node status that would be available from calicoctl.
	Status string `json:"status,omitempty" validate:"omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeStatusList contains a list of NodeStatus resources.
type NodeStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []NodeStatus `json:"items"`
}

// New NodeStatus creates a new (zeroed) NodeStatus struct with the TypeMetadata
// initialized to the current version.
func NewNodeStatus() *NodeStatus {
	return &NodeStatus{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindNodeStatus,
			APIVersion: GroupVersionCurrent,
		},
	}
}

// NewNodeStatusList creates a new zeroed) NodeStatusList struct with the TypeMetadata
// initialized to the current version.
func NewNodeStatusList() *NodeStatusList {
	return &NodeStatusList{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindNodeStatusList,
			APIVersion: GroupVersionCurrent,
		},
	}
}
