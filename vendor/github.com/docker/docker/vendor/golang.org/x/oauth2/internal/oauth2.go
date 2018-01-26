// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal contains support packages for oauth2 package.
package internal

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ParseKey converts the binary contents of a private key file
// to an *rsa.PrivateKey. It detects whether the private key is in a
// PEM container or not. If so, it extracts the the private key
// from PEM container before conversion. It only supports PEM
// containers with no passphrase.
func ParseKey(key []byte) (*rsa.PrivateKey, error) ***REMOVED***
	block, _ := pem.Decode(key)
	if block != nil ***REMOVED***
		key = block.Bytes
	***REMOVED***
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil ***REMOVED***
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("private key should be a PEM or plain PKSC1 or PKCS8; parse error: %v", err)
		***REMOVED***
	***REMOVED***
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok ***REMOVED***
		return nil, errors.New("private key is invalid")
	***REMOVED***
	return parsed, nil
***REMOVED***

func ParseINI(ini io.Reader) (map[string]map[string]string, error) ***REMOVED***
	result := map[string]map[string]string***REMOVED***
		"": map[string]string***REMOVED******REMOVED***, // root section
	***REMOVED***
	scanner := bufio.NewScanner(ini)
	currentSection := ""
	for scanner.Scan() ***REMOVED***
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, ";") ***REMOVED***
			// comment.
			continue
		***REMOVED***
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") ***REMOVED***
			currentSection = strings.TrimSpace(line[1 : len(line)-1])
			result[currentSection] = map[string]string***REMOVED******REMOVED***
			continue
		***REMOVED***
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] != "" ***REMOVED***
			result[currentSection][strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		***REMOVED***
	***REMOVED***
	if err := scanner.Err(); err != nil ***REMOVED***
		return nil, fmt.Errorf("error scanning ini: %v", err)
	***REMOVED***
	return result, nil
***REMOVED***

func CondVal(v string) []string ***REMOVED***
	if v == "" ***REMOVED***
		return nil
	***REMOVED***
	return []string***REMOVED***v***REMOVED***
***REMOVED***
