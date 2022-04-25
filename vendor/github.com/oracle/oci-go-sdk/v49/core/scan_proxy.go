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

// ScanProxy Details pertaining to a scan proxy instance created for a scan listener FQDN/IPs
type ScanProxy struct {

	// An Oracle-assigned unique identifier for a scanProxy within a privateEndpoint. You specify
	// this ID when you want to get or update or delete a scanProxy
	Id *string `mandatory:"true" json:"id"`

	// Type indicating whether Scan listener is specified by its FQDN or list of IPs
	ScanListenerType ScanProxyScanListenerTypeEnum `mandatory:"true" json:"scanListenerType"`

	// The FQDN/IPs and port information of customer's Real Application Cluster (RAC)'s SCAN
	// listeners.
	ScanListenerInfo []ScanListenerInfo `mandatory:"true" json:"scanListenerInfo"`

	// The port to which service DB client has to connect on scan proxy to initiate scan
	// connections.
	ScanProxyPort *int `mandatory:"true" json:"scanProxyPort"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the private endpoint
	// associated with the reverse connection.
	PrivateEndpointId *string `mandatory:"true" json:"privateEndpointId"`

	// The protocol used for communication between client, scanProxy and RAC's scan
	// listeners
	Protocol ScanProxyProtocolEnum `mandatory:"false" json:"protocol,omitempty"`

	// The scan proxy instance's current state.
	LifecycleState ScanProxyLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The date and time the scan proxy instance was created, in the format defined
	// by RFC3339 (https://tools.ietf.org/html/rfc3339).
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	ScanListenerWallet *WalletInfo `mandatory:"false" json:"scanListenerWallet"`
}

func (m ScanProxy) String() string {
	return common.PointerString(m)
}

// ScanProxyScanListenerTypeEnum Enum with underlying type: string
type ScanProxyScanListenerTypeEnum string

// Set of constants representing the allowable values for ScanProxyScanListenerTypeEnum
const (
	ScanProxyScanListenerTypeFqdn ScanProxyScanListenerTypeEnum = "FQDN"
	ScanProxyScanListenerTypeIp   ScanProxyScanListenerTypeEnum = "IP"
)

var mappingScanProxyScanListenerType = map[string]ScanProxyScanListenerTypeEnum{
	"FQDN": ScanProxyScanListenerTypeFqdn,
	"IP":   ScanProxyScanListenerTypeIp,
}

// GetScanProxyScanListenerTypeEnumValues Enumerates the set of values for ScanProxyScanListenerTypeEnum
func GetScanProxyScanListenerTypeEnumValues() []ScanProxyScanListenerTypeEnum {
	values := make([]ScanProxyScanListenerTypeEnum, 0)
	for _, v := range mappingScanProxyScanListenerType {
		values = append(values, v)
	}
	return values
}

// ScanProxyProtocolEnum Enum with underlying type: string
type ScanProxyProtocolEnum string

// Set of constants representing the allowable values for ScanProxyProtocolEnum
const (
	ScanProxyProtocolTcp  ScanProxyProtocolEnum = "TCP"
	ScanProxyProtocolTcps ScanProxyProtocolEnum = "TCPS"
)

var mappingScanProxyProtocol = map[string]ScanProxyProtocolEnum{
	"TCP":  ScanProxyProtocolTcp,
	"TCPS": ScanProxyProtocolTcps,
}

// GetScanProxyProtocolEnumValues Enumerates the set of values for ScanProxyProtocolEnum
func GetScanProxyProtocolEnumValues() []ScanProxyProtocolEnum {
	values := make([]ScanProxyProtocolEnum, 0)
	for _, v := range mappingScanProxyProtocol {
		values = append(values, v)
	}
	return values
}

// ScanProxyLifecycleStateEnum Enum with underlying type: string
type ScanProxyLifecycleStateEnum string

// Set of constants representing the allowable values for ScanProxyLifecycleStateEnum
const (
	ScanProxyLifecycleStateProvisioning ScanProxyLifecycleStateEnum = "PROVISIONING"
	ScanProxyLifecycleStateAvailable    ScanProxyLifecycleStateEnum = "AVAILABLE"
	ScanProxyLifecycleStateUpdating     ScanProxyLifecycleStateEnum = "UPDATING"
	ScanProxyLifecycleStateTerminating  ScanProxyLifecycleStateEnum = "TERMINATING"
	ScanProxyLifecycleStateTerminated   ScanProxyLifecycleStateEnum = "TERMINATED"
	ScanProxyLifecycleStateFailed       ScanProxyLifecycleStateEnum = "FAILED"
)

var mappingScanProxyLifecycleState = map[string]ScanProxyLifecycleStateEnum{
	"PROVISIONING": ScanProxyLifecycleStateProvisioning,
	"AVAILABLE":    ScanProxyLifecycleStateAvailable,
	"UPDATING":     ScanProxyLifecycleStateUpdating,
	"TERMINATING":  ScanProxyLifecycleStateTerminating,
	"TERMINATED":   ScanProxyLifecycleStateTerminated,
	"FAILED":       ScanProxyLifecycleStateFailed,
}

// GetScanProxyLifecycleStateEnumValues Enumerates the set of values for ScanProxyLifecycleStateEnum
func GetScanProxyLifecycleStateEnumValues() []ScanProxyLifecycleStateEnum {
	values := make([]ScanProxyLifecycleStateEnum, 0)
	for _, v := range mappingScanProxyLifecycleState {
		values = append(values, v)
	}
	return values
}
