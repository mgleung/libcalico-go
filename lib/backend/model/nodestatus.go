// Copyright (c) 2019 Tigera, Inc. All rights reserved.
//
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

package model

import (
	"fmt"
	"reflect"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/projectcalico/libcalico-go/lib/errors"
)

var (
	matchNodeStatus = regexp.MustCompile("^/?calico/node/v1/status/([^/]+)/(.+)$")
	typeNodeStatus  = rawStringType
)

type NodeStatusKey struct {
	// The hostname for the host specific node status
	Nodename string `json:"-" validate:"required,name"`

	// The name of the host specific node status key.
	Name string `json:"-" validate:"required,name"`
}

func (key NodeStatusKey) defaultPath() (string, error) {
	return key.defaultDeletePath()
}

func (key NodeStatusKey) defaultDeletePath() (string, error) {
	if key.Nodename == "" {
		return "", errors.ErrorInsufficientIdentifiers{Name: "node"}
	}
	if key.Name == "" {
		return "", errors.ErrorInsufficientIdentifiers{Name: "name"}
	}
	e := fmt.Sprintf("/calico/node/v1/status/%s/%s", key.Nodename, key.Name)
	return e, nil
}

func (key NodeStatusKey) defaultDeleteParentPaths() ([]string, error) {
	return nil, nil
}

func (key NodeStatusKey) valueType() (reflect.Type, error) {
	return typeNodeStatus, nil
}

func (key NodeStatusKey) String() string {
	return fmt.Sprintf("NodeStatus(node=%s; name=%s)", key.Nodename, key.Name)
}

type NodeStatusListOptions struct {
	Nodename string
	Name     string
}

func (options NodeStatusListOptions) defaultPathRoot() string {
	k := "/calico/node/v1/status/%s"
	if options.Nodename == "" {
		return k
	}
	k = k + fmt.Sprintf("/%s", options.Nodename)
	if options.Name == "" {
		return k
	}
	k = k + fmt.Sprintf("/%s", options.Name)
	return k
}

func (options NodeStatusListOptions) KeyFromDefaultPath(path string) Key {
	log.Debugf("Get NodeStatus key from %s", path)
	r := matchNodeStatus.FindAllStringSubmatch(path, -1)
	if len(r) != 1 {
		log.Debugf("Didn't match regex")
		return nil
	}
	nodename := r[0][1]
	name := r[0][2]
	if options.Nodename != "" && nodename != options.Nodename {
		log.Debugf("Didn't match nodename %s != %s", options.Nodename, nodename)
		return nil
	}
	if options.Name != "" && name != options.Name {
		log.Debugf("Didn't match name %s != %s", options.Name, name)
		return nil
	}
	return NodeStatusKey{Nodename: nodename, Name: name}
}
