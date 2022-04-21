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

// SslCertDetails The encoded SslCert of ClientVpn which is used in LDAP.
type SslCertDetails interface {
}

type sslcertdetails struct {
	JsonData    []byte
	ContentType string `json:"contentType"`
}

// UnmarshalJSON unmarshals json
func (m *sslcertdetails) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalersslcertdetails sslcertdetails
	s := struct {
		Model Unmarshalersslcertdetails
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.ContentType = s.Model.ContentType

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *sslcertdetails) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.ContentType {
	case "BASE64ENCODED":
		mm := Base64SslCertDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

func (m sslcertdetails) String() string {
	return common.PointerString(m)
}

// SslCertDetailsContentTypeEnum Enum with underlying type: string
type SslCertDetailsContentTypeEnum string

// Set of constants representing the allowable values for SslCertDetailsContentTypeEnum
const (
	SslCertDetailsContentTypeBase64encoded SslCertDetailsContentTypeEnum = "BASE64ENCODED"
)

var mappingSslCertDetailsContentType = map[string]SslCertDetailsContentTypeEnum{
	"BASE64ENCODED": SslCertDetailsContentTypeBase64encoded,
}

// GetSslCertDetailsContentTypeEnumValues Enumerates the set of values for SslCertDetailsContentTypeEnum
func GetSslCertDetailsContentTypeEnumValues() []SslCertDetailsContentTypeEnum {
	values := make([]SslCertDetailsContentTypeEnum, 0)
	for _, v := range mappingSslCertDetailsContentType {
		values = append(values, v)
	}
	return values
}
