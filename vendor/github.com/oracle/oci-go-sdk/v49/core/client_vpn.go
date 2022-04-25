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

// ClientVpn The ClientVPNEnpoint is a server endpoint that allow customer get access to openVPN service.
type ClientVpn struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment that user sent request.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of clientVPNEndpoint.
	Id *string `mandatory:"true" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the attachedSubnet (VNIC) in customer tenancy.
	SubnetId *string `mandatory:"true" json:"subnetId"`

	// The IP address in attached subnet.
	IpAddressInAttachedSubnet *string `mandatory:"true" json:"ipAddressInAttachedSubnet"`

	// Whether re-route Internet traffic or not.
	IsRerouteEnabled *bool `mandatory:"true" json:"isRerouteEnabled"`

	// A list of accessible subnets from this clientVPNEndpoint.
	AccessibleSubnetCidrs []string `mandatory:"true" json:"accessibleSubnetCidrs"`

	DnsConfig *DnsConfigDetails `mandatory:"true" json:"dnsConfig"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The current state of the ClientVPNenpoint.
	LifecycleState ClientVpnLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// A limit that allows the maximum number of VPN concurrent connections.
	MaxConnections *int `mandatory:"false" json:"maxConnections"`

	// A subnet for openVPN clients to access servers.
	ClientSubnetCidr *string `mandatory:"false" json:"clientSubnetCidr"`

	// Allowed values:
	//   * `NAT`: NAT mode supports the one-way access. In NAT mode, client can access the Internet from server endpoint
	//   but server endpoint cannot access the Internet from client.
	//   * `ROUTING`: ROUTING mode supports the two-way access. In ROUTING mode, client and server endpoint can access the
	//   Internet to each other.
	AddressingMode ClientVpnAddressingModeEnum `mandatory:"false" json:"addressingMode,omitempty"`

	// Allowed values:
	//   * `LOCAL`: Local authentication mode that applies users and password to get authentication through the server.
	//   * `RADIUS`: RADIUS authentication mode applies users and password to get authentication through the RADIUS server.
	//   * `LDAP`: LDAP authentication mode that applies users and passwords to get authentication through the LDAP server.
	AuthenticationMode ClientVpnAuthenticationModeEnum `mandatory:"false" json:"authenticationMode,omitempty"`

	// The time ClientVpn was created, in the format defined by RFC3339 (https://tools.ietf.org/html/rfc3339).
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	RadiusConfig *RadiusConfigDetails `mandatory:"false" json:"radiusConfig"`

	LdapConfig *LdapConfigDetails `mandatory:"false" json:"ldapConfig"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m ClientVpn) String() string {
	return common.PointerString(m)
}

// ClientVpnLifecycleStateEnum Enum with underlying type: string
type ClientVpnLifecycleStateEnum string

// Set of constants representing the allowable values for ClientVpnLifecycleStateEnum
const (
	ClientVpnLifecycleStateCreating ClientVpnLifecycleStateEnum = "CREATING"
	ClientVpnLifecycleStateActive   ClientVpnLifecycleStateEnum = "ACTIVE"
	ClientVpnLifecycleStateInactive ClientVpnLifecycleStateEnum = "INACTIVE"
	ClientVpnLifecycleStateFailed   ClientVpnLifecycleStateEnum = "FAILED"
	ClientVpnLifecycleStateDeleted  ClientVpnLifecycleStateEnum = "DELETED"
	ClientVpnLifecycleStateDeleting ClientVpnLifecycleStateEnum = "DELETING"
	ClientVpnLifecycleStateUpdating ClientVpnLifecycleStateEnum = "UPDATING"
)

var mappingClientVpnLifecycleState = map[string]ClientVpnLifecycleStateEnum{
	"CREATING": ClientVpnLifecycleStateCreating,
	"ACTIVE":   ClientVpnLifecycleStateActive,
	"INACTIVE": ClientVpnLifecycleStateInactive,
	"FAILED":   ClientVpnLifecycleStateFailed,
	"DELETED":  ClientVpnLifecycleStateDeleted,
	"DELETING": ClientVpnLifecycleStateDeleting,
	"UPDATING": ClientVpnLifecycleStateUpdating,
}

// GetClientVpnLifecycleStateEnumValues Enumerates the set of values for ClientVpnLifecycleStateEnum
func GetClientVpnLifecycleStateEnumValues() []ClientVpnLifecycleStateEnum {
	values := make([]ClientVpnLifecycleStateEnum, 0)
	for _, v := range mappingClientVpnLifecycleState {
		values = append(values, v)
	}
	return values
}

// ClientVpnAddressingModeEnum Enum with underlying type: string
type ClientVpnAddressingModeEnum string

// Set of constants representing the allowable values for ClientVpnAddressingModeEnum
const (
	ClientVpnAddressingModeNat     ClientVpnAddressingModeEnum = "NAT"
	ClientVpnAddressingModeRouting ClientVpnAddressingModeEnum = "ROUTING"
)

var mappingClientVpnAddressingMode = map[string]ClientVpnAddressingModeEnum{
	"NAT":     ClientVpnAddressingModeNat,
	"ROUTING": ClientVpnAddressingModeRouting,
}

// GetClientVpnAddressingModeEnumValues Enumerates the set of values for ClientVpnAddressingModeEnum
func GetClientVpnAddressingModeEnumValues() []ClientVpnAddressingModeEnum {
	values := make([]ClientVpnAddressingModeEnum, 0)
	for _, v := range mappingClientVpnAddressingMode {
		values = append(values, v)
	}
	return values
}

// ClientVpnAuthenticationModeEnum Enum with underlying type: string
type ClientVpnAuthenticationModeEnum string

// Set of constants representing the allowable values for ClientVpnAuthenticationModeEnum
const (
	ClientVpnAuthenticationModeLocal  ClientVpnAuthenticationModeEnum = "LOCAL"
	ClientVpnAuthenticationModeRadius ClientVpnAuthenticationModeEnum = "RADIUS"
	ClientVpnAuthenticationModeLdap   ClientVpnAuthenticationModeEnum = "LDAP"
)

var mappingClientVpnAuthenticationMode = map[string]ClientVpnAuthenticationModeEnum{
	"LOCAL":  ClientVpnAuthenticationModeLocal,
	"RADIUS": ClientVpnAuthenticationModeRadius,
	"LDAP":   ClientVpnAuthenticationModeLdap,
}

// GetClientVpnAuthenticationModeEnumValues Enumerates the set of values for ClientVpnAuthenticationModeEnum
func GetClientVpnAuthenticationModeEnumValues() []ClientVpnAuthenticationModeEnum {
	values := make([]ClientVpnAuthenticationModeEnum, 0)
	for _, v := range mappingClientVpnAuthenticationMode {
		values = append(values, v)
	}
	return values
}
