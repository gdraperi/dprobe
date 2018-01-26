// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pkcs12 implements some of PKCS#12.
//
// This implementation is distilled from https://tools.ietf.org/html/rfc7292
// and referenced documents. It is intended for decoding P12/PFX-stored
// certificates and keys for use with the crypto/tls package.
package pkcs12

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

var (
	oidDataContentType          = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 7, 1***REMOVED***)
	oidEncryptedDataContentType = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 7, 6***REMOVED***)

	oidFriendlyName     = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 9, 20***REMOVED***)
	oidLocalKeyID       = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 9, 21***REMOVED***)
	oidMicrosoftCSPName = asn1.ObjectIdentifier([]int***REMOVED***1, 3, 6, 1, 4, 1, 311, 17, 1***REMOVED***)
)

type pfxPdu struct ***REMOVED***
	Version  int
	AuthSafe contentInfo
	MacData  macData `asn1:"optional"`
***REMOVED***

type contentInfo struct ***REMOVED***
	ContentType asn1.ObjectIdentifier
	Content     asn1.RawValue `asn1:"tag:0,explicit,optional"`
***REMOVED***

type encryptedData struct ***REMOVED***
	Version              int
	EncryptedContentInfo encryptedContentInfo
***REMOVED***

type encryptedContentInfo struct ***REMOVED***
	ContentType                asn1.ObjectIdentifier
	ContentEncryptionAlgorithm pkix.AlgorithmIdentifier
	EncryptedContent           []byte `asn1:"tag:0,optional"`
***REMOVED***

func (i encryptedContentInfo) Algorithm() pkix.AlgorithmIdentifier ***REMOVED***
	return i.ContentEncryptionAlgorithm
***REMOVED***

func (i encryptedContentInfo) Data() []byte ***REMOVED*** return i.EncryptedContent ***REMOVED***

type safeBag struct ***REMOVED***
	Id         asn1.ObjectIdentifier
	Value      asn1.RawValue     `asn1:"tag:0,explicit"`
	Attributes []pkcs12Attribute `asn1:"set,optional"`
***REMOVED***

type pkcs12Attribute struct ***REMOVED***
	Id    asn1.ObjectIdentifier
	Value asn1.RawValue `asn1:"set"`
***REMOVED***

type encryptedPrivateKeyInfo struct ***REMOVED***
	AlgorithmIdentifier pkix.AlgorithmIdentifier
	EncryptedData       []byte
***REMOVED***

func (i encryptedPrivateKeyInfo) Algorithm() pkix.AlgorithmIdentifier ***REMOVED***
	return i.AlgorithmIdentifier
***REMOVED***

func (i encryptedPrivateKeyInfo) Data() []byte ***REMOVED***
	return i.EncryptedData
***REMOVED***

// PEM block types
const (
	certificateType = "CERTIFICATE"
	privateKeyType  = "PRIVATE KEY"
)

// unmarshal calls asn1.Unmarshal, but also returns an error if there is any
// trailing data after unmarshaling.
func unmarshal(in []byte, out interface***REMOVED******REMOVED***) error ***REMOVED***
	trailing, err := asn1.Unmarshal(in, out)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(trailing) != 0 ***REMOVED***
		return errors.New("pkcs12: trailing data found")
	***REMOVED***
	return nil
***REMOVED***

// ConvertToPEM converts all "safe bags" contained in pfxData to PEM blocks.
func ToPEM(pfxData []byte, password string) ([]*pem.Block, error) ***REMOVED***
	encodedPassword, err := bmpString(password)
	if err != nil ***REMOVED***
		return nil, ErrIncorrectPassword
	***REMOVED***

	bags, encodedPassword, err := getSafeContents(pfxData, encodedPassword)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	blocks := make([]*pem.Block, 0, len(bags))
	for _, bag := range bags ***REMOVED***
		block, err := convertBag(&bag, encodedPassword)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		blocks = append(blocks, block)
	***REMOVED***

	return blocks, nil
***REMOVED***

func convertBag(bag *safeBag, password []byte) (*pem.Block, error) ***REMOVED***
	block := &pem.Block***REMOVED***
		Headers: make(map[string]string),
	***REMOVED***

	for _, attribute := range bag.Attributes ***REMOVED***
		k, v, err := convertAttribute(&attribute)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		block.Headers[k] = v
	***REMOVED***

	switch ***REMOVED***
	case bag.Id.Equal(oidCertBag):
		block.Type = certificateType
		certsData, err := decodeCertBag(bag.Value.Bytes)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		block.Bytes = certsData
	case bag.Id.Equal(oidPKCS8ShroundedKeyBag):
		block.Type = privateKeyType

		key, err := decodePkcs8ShroudedKeyBag(bag.Value.Bytes, password)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch key := key.(type) ***REMOVED***
		case *rsa.PrivateKey:
			block.Bytes = x509.MarshalPKCS1PrivateKey(key)
		case *ecdsa.PrivateKey:
			block.Bytes, err = x509.MarshalECPrivateKey(key)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		default:
			return nil, errors.New("found unknown private key type in PKCS#8 wrapping")
		***REMOVED***
	default:
		return nil, errors.New("don't know how to convert a safe bag of type " + bag.Id.String())
	***REMOVED***
	return block, nil
***REMOVED***

func convertAttribute(attribute *pkcs12Attribute) (key, value string, err error) ***REMOVED***
	isString := false

	switch ***REMOVED***
	case attribute.Id.Equal(oidFriendlyName):
		key = "friendlyName"
		isString = true
	case attribute.Id.Equal(oidLocalKeyID):
		key = "localKeyId"
	case attribute.Id.Equal(oidMicrosoftCSPName):
		// This key is chosen to match OpenSSL.
		key = "Microsoft CSP Name"
		isString = true
	default:
		return "", "", errors.New("pkcs12: unknown attribute with OID " + attribute.Id.String())
	***REMOVED***

	if isString ***REMOVED***
		if err := unmarshal(attribute.Value.Bytes, &attribute.Value); err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
		if value, err = decodeBMPString(attribute.Value.Bytes); err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var id []byte
		if err := unmarshal(attribute.Value.Bytes, &id); err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
		value = hex.EncodeToString(id)
	***REMOVED***

	return key, value, nil
***REMOVED***

// Decode extracts a certificate and private key from pfxData. This function
// assumes that there is only one certificate and only one private key in the
// pfxData.
func Decode(pfxData []byte, password string) (privateKey interface***REMOVED******REMOVED***, certificate *x509.Certificate, err error) ***REMOVED***
	encodedPassword, err := bmpString(password)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	bags, encodedPassword, err := getSafeContents(pfxData, encodedPassword)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if len(bags) != 2 ***REMOVED***
		err = errors.New("pkcs12: expected exactly two safe bags in the PFX PDU")
		return
	***REMOVED***

	for _, bag := range bags ***REMOVED***
		switch ***REMOVED***
		case bag.Id.Equal(oidCertBag):
			if certificate != nil ***REMOVED***
				err = errors.New("pkcs12: expected exactly one certificate bag")
			***REMOVED***

			certsData, err := decodeCertBag(bag.Value.Bytes)
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
			certs, err := x509.ParseCertificates(certsData)
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
			if len(certs) != 1 ***REMOVED***
				err = errors.New("pkcs12: expected exactly one certificate in the certBag")
				return nil, nil, err
			***REMOVED***
			certificate = certs[0]

		case bag.Id.Equal(oidPKCS8ShroundedKeyBag):
			if privateKey != nil ***REMOVED***
				err = errors.New("pkcs12: expected exactly one key bag")
			***REMOVED***

			if privateKey, err = decodePkcs8ShroudedKeyBag(bag.Value.Bytes, encodedPassword); err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if certificate == nil ***REMOVED***
		return nil, nil, errors.New("pkcs12: certificate missing")
	***REMOVED***
	if privateKey == nil ***REMOVED***
		return nil, nil, errors.New("pkcs12: private key missing")
	***REMOVED***

	return
***REMOVED***

func getSafeContents(p12Data, password []byte) (bags []safeBag, updatedPassword []byte, err error) ***REMOVED***
	pfx := new(pfxPdu)
	if err := unmarshal(p12Data, pfx); err != nil ***REMOVED***
		return nil, nil, errors.New("pkcs12: error reading P12 data: " + err.Error())
	***REMOVED***

	if pfx.Version != 3 ***REMOVED***
		return nil, nil, NotImplementedError("can only decode v3 PFX PDU's")
	***REMOVED***

	if !pfx.AuthSafe.ContentType.Equal(oidDataContentType) ***REMOVED***
		return nil, nil, NotImplementedError("only password-protected PFX is implemented")
	***REMOVED***

	// unmarshal the explicit bytes in the content for type 'data'
	if err := unmarshal(pfx.AuthSafe.Content.Bytes, &pfx.AuthSafe.Content); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if len(pfx.MacData.Mac.Algorithm.Algorithm) == 0 ***REMOVED***
		return nil, nil, errors.New("pkcs12: no MAC in data")
	***REMOVED***

	if err := verifyMac(&pfx.MacData, pfx.AuthSafe.Content.Bytes, password); err != nil ***REMOVED***
		if err == ErrIncorrectPassword && len(password) == 2 && password[0] == 0 && password[1] == 0 ***REMOVED***
			// some implementations use an empty byte array
			// for the empty string password try one more
			// time with empty-empty password
			password = nil
			err = verifyMac(&pfx.MacData, pfx.AuthSafe.Content.Bytes, password)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
	***REMOVED***

	var authenticatedSafe []contentInfo
	if err := unmarshal(pfx.AuthSafe.Content.Bytes, &authenticatedSafe); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if len(authenticatedSafe) != 2 ***REMOVED***
		return nil, nil, NotImplementedError("expected exactly two items in the authenticated safe")
	***REMOVED***

	for _, ci := range authenticatedSafe ***REMOVED***
		var data []byte

		switch ***REMOVED***
		case ci.ContentType.Equal(oidDataContentType):
			if err := unmarshal(ci.Content.Bytes, &data); err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
		case ci.ContentType.Equal(oidEncryptedDataContentType):
			var encryptedData encryptedData
			if err := unmarshal(ci.Content.Bytes, &encryptedData); err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
			if encryptedData.Version != 0 ***REMOVED***
				return nil, nil, NotImplementedError("only version 0 of EncryptedData is supported")
			***REMOVED***
			if data, err = pbDecrypt(encryptedData.EncryptedContentInfo, password); err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
		default:
			return nil, nil, NotImplementedError("only data and encryptedData content types are supported in authenticated safe")
		***REMOVED***

		var safeContents []safeBag
		if err := unmarshal(data, &safeContents); err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		bags = append(bags, safeContents...)
	***REMOVED***

	return bags, password, nil
***REMOVED***
