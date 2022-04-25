// Copyright (c) 2016, 2018, 2021, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/v49/common"
)

// VolumeGroupSourceFromVolumeGroupReplicaDetails Specifies the source volume replica which the volume group will be created from.
// The volume group replica shoulbe be in the same availability domain as the volume group.
// Only one volume group can be created from a volume group replica at the same time.
type VolumeGroupSourceFromVolumeGroupReplicaDetails struct {

	// The OCID of the volume group replica.
	VolumeGroupReplicaId *string `mandatory:"true" json:"volumeGroupReplicaId"`
}

func (m VolumeGroupSourceFromVolumeGroupReplicaDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m VolumeGroupSourceFromVolumeGroupReplicaDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeVolumeGroupSourceFromVolumeGroupReplicaDetails VolumeGroupSourceFromVolumeGroupReplicaDetails
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeVolumeGroupSourceFromVolumeGroupReplicaDetails
	}{
		"volumeGroupReplicaId",
		(MarshalTypeVolumeGroupSourceFromVolumeGroupReplicaDetails)(m),
	}

	return json.Marshal(&s)
}
