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
	"github.com/oracle/oci-go-sdk/v49/common"
)

// ChangeVnicAttachmentsCompartmentRequest This structure is used when changing vnic attachment compartment.
type ChangeVnicAttachmentsCompartmentRequest struct {

	// List of VNICs whose attachments need to move to the destination compartment
	VnicIds []string `mandatory:"false" json:"vnicIds"`

	// The ID of the parent resource that these VNICs are attached to
	ParentResourceId *string `mandatory:"false" json:"parentResourceId"`

	// The destination compartment ID
	DestinationCompartmentId *string `mandatory:"false" json:"destinationCompartmentId"`

	// The compartment ID to create the work request in. This is NOT the customer's compartment
	// but one that belongs to the calling service. Typically this is not the same as the destination
	// compartment ID
	WorkRequestCompartmentId *string `mandatory:"false" json:"workRequestCompartmentId"`
}

func (m ChangeVnicAttachmentsCompartmentRequest) String() string {
	return common.PointerString(m)
}
