// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agent

import "testing"

func addTestKey(t *testing.T, a Agent, keyName string) ***REMOVED***
	err := a.Add(AddedKey***REMOVED***
		PrivateKey: testPrivateKeys[keyName],
		Comment:    keyName,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("failed to add key %q: %v", keyName, err)
	***REMOVED***
***REMOVED***

func removeTestKey(t *testing.T, a Agent, keyName string) ***REMOVED***
	err := a.Remove(testPublicKeys[keyName])
	if err != nil ***REMOVED***
		t.Fatalf("failed to remove key %q: %v", keyName, err)
	***REMOVED***
***REMOVED***

func validateListedKeys(t *testing.T, a Agent, expectedKeys []string) ***REMOVED***
	listedKeys, err := a.List()
	if err != nil ***REMOVED***
		t.Fatalf("failed to list keys: %v", err)
		return
	***REMOVED***
	actualKeys := make(map[string]bool)
	for _, key := range listedKeys ***REMOVED***
		actualKeys[key.Comment] = true
	***REMOVED***

	matchedKeys := make(map[string]bool)
	for _, expectedKey := range expectedKeys ***REMOVED***
		if !actualKeys[expectedKey] ***REMOVED***
			t.Fatalf("expected key %q, but was not found", expectedKey)
		***REMOVED*** else ***REMOVED***
			matchedKeys[expectedKey] = true
		***REMOVED***
	***REMOVED***

	for actualKey := range actualKeys ***REMOVED***
		if !matchedKeys[actualKey] ***REMOVED***
			t.Fatalf("key %q was found, but was not expected", actualKey)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestKeyringAddingAndRemoving(t *testing.T) ***REMOVED***
	keyNames := []string***REMOVED***"dsa", "ecdsa", "rsa", "user"***REMOVED***

	// add all test private keys
	k := NewKeyring()
	for _, keyName := range keyNames ***REMOVED***
		addTestKey(t, k, keyName)
	***REMOVED***
	validateListedKeys(t, k, keyNames)

	// remove a key in the middle
	keyToRemove := keyNames[1]
	keyNames = append(keyNames[:1], keyNames[2:]...)

	removeTestKey(t, k, keyToRemove)
	validateListedKeys(t, k, keyNames)

	// remove all keys
	err := k.RemoveAll()
	if err != nil ***REMOVED***
		t.Fatalf("failed to remove all keys: %v", err)
	***REMOVED***
	validateListedKeys(t, k, []string***REMOVED******REMOVED***)
***REMOVED***
