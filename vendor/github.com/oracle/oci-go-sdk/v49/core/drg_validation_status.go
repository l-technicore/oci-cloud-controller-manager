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

// DrgValidationStatus Status of DRG route tables after migration or upgrade. isValid == true indicates that the routes tables before and after migration or upgrade match up as we would expect.
type DrgValidationStatus struct {

	// The `drgId` of the upgraded DRG.
	DrgId *string `mandatory:"true" json:"drgId"`

	// isValid == true indicates that the routes tables before and after migration or upgrade match up as we would expect.
	IsValid *bool `mandatory:"true" json:"isValid"`
}

func (m DrgValidationStatus) String() string {
	return common.PointerString(m)
}
