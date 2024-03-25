// Copyright 2017 Oracle and/or its affiliates. All rights reserved.
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

package oci

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/oracle/oci-cloud-controller-manager/pkg/oci/client"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	listersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/utils/net"
)

// Protects Load Balancers against multiple updates in parallel
type loadBalancerLocks struct {
	locks sets.String
	mux   sync.Mutex
}

func NewLoadBalancerLocks() *loadBalancerLocks {
	return &loadBalancerLocks{
		locks: sets.NewString(),
	}
}

func (lbl *loadBalancerLocks) TryAcquire(lbname string) bool {
	lbl.mux.Lock()
	defer lbl.mux.Unlock()
	if lbl.locks.Has(lbname) {
		return false
	}
	lbl.locks.Insert(lbname)
	return true
}

func (lbl *loadBalancerLocks) Release(lbname string) {
	lbl.mux.Lock()
	defer lbl.mux.Unlock()
	lbl.locks.Delete(lbname)
}

// MapProviderIDToResourceID parses the provider id and returns the instance ocid.
func MapProviderIDToResourceID(providerID string) (string, error) {
	if providerID == "" {
		return providerID, errors.New("provider ID is empty")
	}
	if strings.HasPrefix(providerID, providerPrefix) {
		return strings.TrimPrefix(providerID, providerPrefix), nil
	}
	return providerID, nil
}

// NodeInternalIP returns the nodes internal ip
// A node managed by the CCM will always have an internal ip
// since it's not possible to deploy an instance without a private ip.
func NodeInternalIP(node *api.Node) client.IpAddresses {
	ipAddresses := client.IpAddresses{
		V4: "",
		V6: "",
	}
	for _, addr := range node.Status.Addresses {
		if addr.Type == api.NodeInternalIP {
			if net.IsIPv6String(addr.Address) {
				ipAddresses.V6 = addr.Address
			} else {
				ipAddresses.V4 = addr.Address
			}
		}
	}
	return ipAddresses
}

// NodeExternalIp returns the nodes external ip
// A node managed by the CCM may have an external ip
// in case of IPv6, it could be possible that a compute instance has only GUA IPv6
func NodeExternalIp(node *api.Node) client.IpAddresses {
	ipAddresses := client.IpAddresses{
		V4: "",
		V6: "",
	}
	for _, addr := range node.Status.Addresses {
		if addr.Type == api.NodeExternalIP {
			if net.IsIPv6String(addr.Address) {
				ipAddresses.V6 = addr.Address
			} else {
				ipAddresses.V4 = addr.Address
			}
		}
	}
	return ipAddresses
}

// RemoveDuplicatesFromList takes Slice and returns new Slice with no duplicate elements
// (e.g. if given list is {"a", "b", "a"}, function returns new slice with {"a", "b"}
func RemoveDuplicatesFromList(list []string) []string {
	return sets.NewString(list...).List()
}

// DeepEqualLists diffs two slices and returns bool if the slices are equal/not-equal.
// the duplicates and order of items in both lists is ignored.
func DeepEqualLists(listA, listB []string) bool {
	return sets.NewString(listA...).Equal(sets.NewString(listB...))
}

// GetNodeMap returns a map of nodes in the cluster indexed by node name
func GetNodesMap(nodeLister listersv1.NodeLister) (nodeMap map[string]*api.Node, err error) {
	nodeList, err := nodeLister.List(labels.Everything())
	if err != nil {
		return
	}
	nodeMap = make(map[string]*api.Node)
	for _, node := range nodeList {
		nodeMap[node.Name] = node
	}
	return
}

func GetIsFeatureEnabledFromEnv(logger *zap.SugaredLogger, featureName string, defaultValue bool) bool {
	enableFeature := defaultValue
	enableFeatureEnvVar, ok := os.LookupEnv(featureName)
	if ok {
		var err error
		enableFeature, err = strconv.ParseBool(enableFeatureEnvVar)
		if err != nil {
			logger.With(zap.Error(err)).Errorf("failed to parse %s envvar, defaulting to %t", featureName, defaultValue)
			return defaultValue
		}
	}
	return enableFeature
}
