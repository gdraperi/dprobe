// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkcs12

import (
	"crypto/x509"
	"encoding/asn1"
	"errors"
)

var (
	// see https://tools.ietf.org/html/rfc7292#appendix-D
	oidCertTypeX509Certificate = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 9, 22, 1***REMOVED***)
	oidPKCS8ShroundedKeyBag    = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 12, 10, 1, 2***REMOVED***)
	oidCertBag                 = asn1.ObjectIdentifier([]int***REMOVED***1, 2, 840, 113549, 1, 12, 10, 1, 3***REMOVED***)
)

type certBag struct ***REMOVED***
	Id   asn1.ObjectIdentifier
	Data []byte `asn1:"tag:0,explicit"`
***REMOVED***

func decodePkcs8ShroudedKeyBag(asn1Data, password []byte) (privateKey interface***REMOVED******REMOVED***, err error) ***REMOVED***
	pkinfo := new(encryptedPrivateKeyInfo)
	if err = unmarshal(asn1Data, pkinfo); err != nil ***REMOVED***
		return nil, errors.New("pkcs12: error decoding PKCS#8 shrouded key bag: " + err.Error())
	***REMOVED***

	pkData, err := pbDecrypt(pkinfo, password)
	if err != nil ***REMOVED***
		return nil, errors.New("pkcs12: error decrypting PKCS#8 shrouded key bag: " + err.Error())
	***REMOVED***

	ret := new(asn1.RawValue)
	if err = unmarshal(pkData, ret); err != nil ***REMOVED***
		return nil, errors.New("pkcs12: error unmarshaling decrypted private key: " + err.Error())
	***REMOVED***

	if privateKey, err = x509.ParsePKCS8PrivateKey(pkData); err != nil ***REMOVED***
		return nil, errors.New("pkcs12: error parsing PKCS#8 private key: " + err.Error())
	***REMOVED***

	return privateKey, nil
***REMOVED***

func decodeCertBag(asn1Data []byte) (x509Certificates []byte, err error) ***REMOVED***
	bag := new(certBag)
	if err := unmarshal(asn1Data, bag); err != nil ***REMOVED***
		return nil, errors.New("pkcs12: error decoding cert bag: " + err.Error())
	***REMOVED***
	if !bag.Id.Equal(oidCertTypeX509Certificate) ***REMOVED***
		return nil, NotImplementedError("only X509 certificates are supported")
	***REMOVED***
	return bag.Data, nil
***REMOVED***
