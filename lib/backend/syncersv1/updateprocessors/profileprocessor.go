// Copyright (c) 2017 Tigera, Inc. All rights reserved.

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

package updateprocessors

import (
	"errors"
	"strings"

	"github.com/projectcalico/libcalico-go/lib/apiv2"
	"github.com/projectcalico/libcalico-go/lib/backend/model"
	"github.com/projectcalico/libcalico-go/lib/backend/watchersyncer"
	cnet "github.com/projectcalico/libcalico-go/lib/net"
)

// Create a new SyncerUpdateProcessor to sync Profile data in v1 format for
// consumption by Felix.
func NewProfileUpdateProcessor() watchersyncer.SyncerUpdateProcessor {
	return NewGeneralUpdateProcessor(apiv2.KindProfile, convertProfileV2ToV1)
}

// Convert v2 KVPair to the equivalent v1 KVPair.
func convertProfileV2ToV1(kvp *model.KVPair) (*model.KVPair, error) {
	// Validate against incorrect key/value kinds.  This indicates a code bug rather
	// than a user error.
	v2key, ok := kvp.Key.(model.ResourceKey)
	if !ok || v2key.Kind != apiv2.KindProfile {
		return nil, errors.New("Key is not a valid Profile resource key")
	}

	if kvp.Value == nil {
		return nil, errors.New("Deletion attempted without enough information to form a v1 Profile key")
	}

	v2res, ok := kvp.Value.(*apiv2.Profile)
	if !ok {
		return nil, errors.New("Value is not a valid Profile resource value")
	}

	v1key, err := convertProfileV2ToV1Key(v2res)
	if err != nil {
		return nil, err
	}

	v1value, err := convertProfileV2ToV1Value(v2res)
	if err != nil {
		// Currently ignore the error so that incorrect values get skipped instead of erroring out
		return &model.KVPair{
			Key: v1key,
		}, nil
	}

	return &model.KVPair{
		Key:      v1key,
		Value:    v1value,
		Revision: kvp.Revision,
	}, nil
}

func convertProfileV2ToV1Key(v2res *apiv2.Profile) (model.ProfileKey, error) {
	if v2res.GetName() == "" {
		return model.ProfileKey{}, errors.New("Missing Name field to create a v1 Profile Key")
	}
	return model.ProfileKey{
		Name: v2res.GetName(),
	}, nil

}

func convertProfileV2ToV1Value(v2res *apiv2.Profile) (*model.Profile, error) {
	var v1value *model.Profile
	// Deletion operations will have empty values so skip if empty
	if !v1ProfileFieldsEmpty(v2res) {
		var irules []model.Rule
		for _, irule := range v2res.Spec.IngressRules {
			irules = append(irules, RuleAPIV2ToBackend(irule))
		}

		var erules []model.Rule
		for _, erule := range v2res.Spec.EgressRules {
			erules = append(erules, RuleAPIV2ToBackend(erule))
		}

		rules := model.ProfileRules{
			InboundRules:  irules,
			OutboundRules: erules,
		}

		v1value = &model.Profile{
			Rules:  rules,
			Labels: v2res.Spec.LabelsToApply,
		}
	}

	return v1value, nil
}

func v1ProfileFieldsEmpty(v2res *apiv2.Profile) bool {
	empty := true
	empty = empty && len(v2res.Spec.IngressRules) == 0
	empty = empty && len(v2res.Spec.EgressRules) == 0
	empty = empty && len(v2res.Spec.LabelsToApply) == 0
	return empty
}

func RulesAPIV2ToBackend(ars []apiv2.Rule) []model.Rule {
	if ars == nil {
		return []model.Rule{}
	}

	brs := make([]model.Rule, len(ars))
	for idx, ar := range ars {
		brs[idx] = RuleAPIV2ToBackend(ar)
	}
	return brs
}

// RuleAPIToBackend converts an API Rule structure to a Backend Rule structure.
func RuleAPIV2ToBackend(ar apiv2.Rule) model.Rule {
	var icmpCode, icmpType, notICMPCode, notICMPType *int
	if ar.ICMP != nil {
		icmpCode = ar.ICMP.Code
		icmpType = ar.ICMP.Type
	}

	if ar.NotICMP != nil {
		notICMPCode = ar.NotICMP.Code
		notICMPType = ar.NotICMP.Type
	}

	return model.Rule{
		Action:      ruleActionAPIV2ToBackend(ar.Action),
		IPVersion:   ar.IPVersion,
		Protocol:    ar.Protocol,
		ICMPCode:    icmpCode,
		ICMPType:    icmpType,
		NotProtocol: ar.NotProtocol,
		NotICMPCode: notICMPCode,
		NotICMPType: notICMPType,

		SrcTag:      ar.Source.Tag,
		SrcNets:     convertStringsToNets(ar.Source.Nets),
		SrcSelector: ar.Source.Selector,
		SrcPorts:    ar.Source.Ports,
		DstTag:      ar.Destination.Tag,
		DstNets:     normalizeIPNets(ar.Destination.Nets),
		DstSelector: ar.Destination.Selector,
		DstPorts:    ar.Destination.Ports,

		NotSrcTag:      ar.Source.NotTag,
		NotSrcNets:     convertStringsToNets(ar.Source.NotNets),
		NotSrcSelector: ar.Source.NotSelector,
		NotSrcPorts:    ar.Source.NotPorts,
		NotDstTag:      ar.Destination.NotTag,
		NotDstNets:     normalizeIPNets(ar.Destination.NotNets),
		NotDstSelector: ar.Destination.NotSelector,
		NotDstPorts:    ar.Destination.NotPorts,
	}
}

// normalizeIPNet converts an IPNet to a network by ensuring the IP address is correctly masked.
func normalizeIPNet(n string) *cnet.IPNet {
	if n == "" {
		return nil
	}
	_, ipn, err := cnet.ParseCIDROrIP(n)
	if err != nil {
		return nil
	}
	return ipn.Network()
}

// normalizeIPNets converts an []*IPNet to a slice of networks by ensuring the IP addresses
// are correctly masked.
func normalizeIPNets(nets []string) []*cnet.IPNet {
	if len(nets) == 0 {
		return nil
	}
	out := make([]*cnet.IPNet, len(nets))
	for i, n := range nets {
		out[i] = normalizeIPNet(n)
	}
	return out
}

// ruleActionAPIV2ToBackend converts the rule action field value from the API
// value to the equivalent backend value.
func ruleActionAPIV2ToBackend(action apiv2.Action) string {
	if action == apiv2.Pass {
		// TODO: why does this say tiers?
		return "next-tier"
	}
	return strings.ToLower(string(action))
}

func convertStringsToNets(strs []string) []*cnet.IPNet {
	var nets []*cnet.IPNet
	for _, str := range strs {
		_, ipn, err := cnet.ParseCIDROrIP(str)
		if err != nil {
			continue
		}
		nets = append(nets, ipn)
	}
	return nets
}
