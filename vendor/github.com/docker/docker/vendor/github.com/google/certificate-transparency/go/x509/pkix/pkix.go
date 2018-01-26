// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pkix contains shared, low level structures used for ASN.1 parsing
// and serialization of X.509 certificates, CRL and OCSP.
package pkix

import (
	// START CT CHANGES
	"github.com/google/certificate-transparency/go/asn1"
	// END CT CHANGES
	"math/big"
	"time"
)

// AlgorithmIdentifier represents the ASN.1 structure of the same name. See RFC
// 5280, section 4.1.1.2.
type AlgorithmIdentifier struct ***REMOVED***
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
***REMOVED***

type RDNSequence []RelativeDistinguishedNameSET

type RelativeDistinguishedNameSET []AttributeTypeAndValue

// AttributeTypeAndValue mirrors the ASN.1 structure of the same name in
// http://tools.ietf.org/html/rfc5280#section-4.1.2.4
type AttributeTypeAndValue struct ***REMOVED***
	Type  asn1.ObjectIdentifier
	Value interface***REMOVED******REMOVED***
***REMOVED***

// Extension represents the ASN.1 structure of the same name. See RFC
// 5280, section 4.2.
type Extension struct ***REMOVED***
	Id       asn1.ObjectIdentifier
	Critical bool `asn1:"optional"`
	Value    []byte
***REMOVED***

// Name represents an X.509 distinguished name. This only includes the common
// elements of a DN.  Additional elements in the name are ignored.
type Name struct ***REMOVED***
	Country, Organization, OrganizationalUnit []string
	Locality, Province                        []string
	StreetAddress, PostalCode                 []string
	SerialNumber, CommonName                  string

	Names []AttributeTypeAndValue
***REMOVED***

func (n *Name) FillFromRDNSequence(rdns *RDNSequence) ***REMOVED***
	for _, rdn := range *rdns ***REMOVED***
		if len(rdn) == 0 ***REMOVED***
			continue
		***REMOVED***
		atv := rdn[0]
		n.Names = append(n.Names, atv)
		value, ok := atv.Value.(string)
		if !ok ***REMOVED***
			continue
		***REMOVED***

		t := atv.Type
		if len(t) == 4 && t[0] == 2 && t[1] == 5 && t[2] == 4 ***REMOVED***
			switch t[3] ***REMOVED***
			case 3:
				n.CommonName = value
			case 5:
				n.SerialNumber = value
			case 6:
				n.Country = append(n.Country, value)
			case 7:
				n.Locality = append(n.Locality, value)
			case 8:
				n.Province = append(n.Province, value)
			case 9:
				n.StreetAddress = append(n.StreetAddress, value)
			case 10:
				n.Organization = append(n.Organization, value)
			case 11:
				n.OrganizationalUnit = append(n.OrganizationalUnit, value)
			case 17:
				n.PostalCode = append(n.PostalCode, value)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var (
	oidCountry            = []int***REMOVED***2, 5, 4, 6***REMOVED***
	oidOrganization       = []int***REMOVED***2, 5, 4, 10***REMOVED***
	oidOrganizationalUnit = []int***REMOVED***2, 5, 4, 11***REMOVED***
	oidCommonName         = []int***REMOVED***2, 5, 4, 3***REMOVED***
	oidSerialNumber       = []int***REMOVED***2, 5, 4, 5***REMOVED***
	oidLocality           = []int***REMOVED***2, 5, 4, 7***REMOVED***
	oidProvince           = []int***REMOVED***2, 5, 4, 8***REMOVED***
	oidStreetAddress      = []int***REMOVED***2, 5, 4, 9***REMOVED***
	oidPostalCode         = []int***REMOVED***2, 5, 4, 17***REMOVED***
)

// appendRDNs appends a relativeDistinguishedNameSET to the given RDNSequence
// and returns the new value. The relativeDistinguishedNameSET contains an
// attributeTypeAndValue for each of the given values. See RFC 5280, A.1, and
// search for AttributeTypeAndValue.
func appendRDNs(in RDNSequence, values []string, oid asn1.ObjectIdentifier) RDNSequence ***REMOVED***
	if len(values) == 0 ***REMOVED***
		return in
	***REMOVED***

	s := make([]AttributeTypeAndValue, len(values))
	for i, value := range values ***REMOVED***
		s[i].Type = oid
		s[i].Value = value
	***REMOVED***

	return append(in, s)
***REMOVED***

func (n Name) ToRDNSequence() (ret RDNSequence) ***REMOVED***
	ret = appendRDNs(ret, n.Country, oidCountry)
	ret = appendRDNs(ret, n.Organization, oidOrganization)
	ret = appendRDNs(ret, n.OrganizationalUnit, oidOrganizationalUnit)
	ret = appendRDNs(ret, n.Locality, oidLocality)
	ret = appendRDNs(ret, n.Province, oidProvince)
	ret = appendRDNs(ret, n.StreetAddress, oidStreetAddress)
	ret = appendRDNs(ret, n.PostalCode, oidPostalCode)
	if len(n.CommonName) > 0 ***REMOVED***
		ret = appendRDNs(ret, []string***REMOVED***n.CommonName***REMOVED***, oidCommonName)
	***REMOVED***
	if len(n.SerialNumber) > 0 ***REMOVED***
		ret = appendRDNs(ret, []string***REMOVED***n.SerialNumber***REMOVED***, oidSerialNumber)
	***REMOVED***

	return ret
***REMOVED***

// CertificateList represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1. Use Certificate.CheckCRLSignature to verify the
// signature.
type CertificateList struct ***REMOVED***
	TBSCertList        TBSCertificateList
	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     asn1.BitString
***REMOVED***

// HasExpired reports whether now is past the expiry time of certList.
func (certList *CertificateList) HasExpired(now time.Time) bool ***REMOVED***
	return now.After(certList.TBSCertList.NextUpdate)
***REMOVED***

// TBSCertificateList represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1.
type TBSCertificateList struct ***REMOVED***
	Raw                 asn1.RawContent
	Version             int `asn1:"optional,default:2"`
	Signature           AlgorithmIdentifier
	Issuer              RDNSequence
	ThisUpdate          time.Time
	NextUpdate          time.Time
	RevokedCertificates []RevokedCertificate `asn1:"optional"`
	Extensions          []Extension          `asn1:"tag:0,optional,explicit"`
***REMOVED***

// RevokedCertificate represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1.
type RevokedCertificate struct ***REMOVED***
	SerialNumber   *big.Int
	RevocationTime time.Time
	Extensions     []Extension `asn1:"optional"`
***REMOVED***
