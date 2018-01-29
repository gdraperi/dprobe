// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bcrypt

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBcryptingIsEasy(t *testing.T) ***REMOVED***
	pass := []byte("mypassword")
	hp, err := GenerateFromPassword(pass, 0)
	if err != nil ***REMOVED***
		t.Fatalf("GenerateFromPassword error: %s", err)
	***REMOVED***

	if CompareHashAndPassword(hp, pass) != nil ***REMOVED***
		t.Errorf("%v should hash %s correctly", hp, pass)
	***REMOVED***

	notPass := "notthepass"
	err = CompareHashAndPassword(hp, []byte(notPass))
	if err != ErrMismatchedHashAndPassword ***REMOVED***
		t.Errorf("%v and %s should be mismatched", hp, notPass)
	***REMOVED***
***REMOVED***

func TestBcryptingIsCorrect(t *testing.T) ***REMOVED***
	pass := []byte("allmine")
	salt := []byte("XajjQvNhvvRt5GSeFk1xFe")
	expectedHash := []byte("$2a$10$XajjQvNhvvRt5GSeFk1xFeyqRrsxkhBkUiQeg0dt.wU1qD4aFDcga")

	hash, err := bcrypt(pass, 10, salt)
	if err != nil ***REMOVED***
		t.Fatalf("bcrypt blew up: %v", err)
	***REMOVED***
	if !bytes.HasSuffix(expectedHash, hash) ***REMOVED***
		t.Errorf("%v should be the suffix of %v", hash, expectedHash)
	***REMOVED***

	h, err := newFromHash(expectedHash)
	if err != nil ***REMOVED***
		t.Errorf("Unable to parse %s: %v", string(expectedHash), err)
	***REMOVED***

	// This is not the safe way to compare these hashes. We do this only for
	// testing clarity. Use bcrypt.CompareHashAndPassword()
	if err == nil && !bytes.Equal(expectedHash, h.Hash()) ***REMOVED***
		t.Errorf("Parsed hash %v should equal %v", h.Hash(), expectedHash)
	***REMOVED***
***REMOVED***

func TestVeryShortPasswords(t *testing.T) ***REMOVED***
	key := []byte("k")
	salt := []byte("XajjQvNhvvRt5GSeFk1xFe")
	_, err := bcrypt(key, 10, salt)
	if err != nil ***REMOVED***
		t.Errorf("One byte key resulted in error: %s", err)
	***REMOVED***
***REMOVED***

func TestTooLongPasswordsWork(t *testing.T) ***REMOVED***
	salt := []byte("XajjQvNhvvRt5GSeFk1xFe")
	// One byte over the usual 56 byte limit that blowfish has
	tooLongPass := []byte("012345678901234567890123456789012345678901234567890123456")
	tooLongExpected := []byte("$2a$10$XajjQvNhvvRt5GSeFk1xFe5l47dONXg781AmZtd869sO8zfsHuw7C")
	hash, err := bcrypt(tooLongPass, 10, salt)
	if err != nil ***REMOVED***
		t.Fatalf("bcrypt blew up on long password: %v", err)
	***REMOVED***
	if !bytes.HasSuffix(tooLongExpected, hash) ***REMOVED***
		t.Errorf("%v should be the suffix of %v", hash, tooLongExpected)
	***REMOVED***
***REMOVED***

type InvalidHashTest struct ***REMOVED***
	err  error
	hash []byte
***REMOVED***

var invalidTests = []InvalidHashTest***REMOVED***
	***REMOVED***ErrHashTooShort, []byte("$2a$10$fooo")***REMOVED***,
	***REMOVED***ErrHashTooShort, []byte("$2a")***REMOVED***,
	***REMOVED***HashVersionTooNewError('3'), []byte("$3a$10$sssssssssssssssssssssshhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")***REMOVED***,
	***REMOVED***InvalidHashPrefixError('%'), []byte("%2a$10$sssssssssssssssssssssshhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")***REMOVED***,
	***REMOVED***InvalidCostError(32), []byte("$2a$32$sssssssssssssssssssssshhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")***REMOVED***,
***REMOVED***

func TestInvalidHashErrors(t *testing.T) ***REMOVED***
	check := func(name string, expected, err error) ***REMOVED***
		if err == nil ***REMOVED***
			t.Errorf("%s: Should have returned an error", name)
		***REMOVED***
		if err != nil && err != expected ***REMOVED***
			t.Errorf("%s gave err %v but should have given %v", name, err, expected)
		***REMOVED***
	***REMOVED***
	for _, iht := range invalidTests ***REMOVED***
		_, err := newFromHash(iht.hash)
		check("newFromHash", iht.err, err)
		err = CompareHashAndPassword(iht.hash, []byte("anything"))
		check("CompareHashAndPassword", iht.err, err)
	***REMOVED***
***REMOVED***

func TestUnpaddedBase64Encoding(t *testing.T) ***REMOVED***
	original := []byte***REMOVED***101, 201, 101, 75, 19, 227, 199, 20, 239, 236, 133, 32, 30, 109, 243, 30***REMOVED***
	encodedOriginal := []byte("XajjQvNhvvRt5GSeFk1xFe")

	encoded := base64Encode(original)

	if !bytes.Equal(encodedOriginal, encoded) ***REMOVED***
		t.Errorf("Encoded %v should have equaled %v", encoded, encodedOriginal)
	***REMOVED***

	decoded, err := base64Decode(encodedOriginal)
	if err != nil ***REMOVED***
		t.Fatalf("base64Decode blew up: %s", err)
	***REMOVED***

	if !bytes.Equal(decoded, original) ***REMOVED***
		t.Errorf("Decoded %v should have equaled %v", decoded, original)
	***REMOVED***
***REMOVED***

func TestCost(t *testing.T) ***REMOVED***
	suffix := "XajjQvNhvvRt5GSeFk1xFe5l47dONXg781AmZtd869sO8zfsHuw7C"
	for _, vers := range []string***REMOVED***"2a", "2"***REMOVED*** ***REMOVED***
		for _, cost := range []int***REMOVED***4, 10***REMOVED*** ***REMOVED***
			s := fmt.Sprintf("$%s$%02d$%s", vers, cost, suffix)
			h := []byte(s)
			actual, err := Cost(h)
			if err != nil ***REMOVED***
				t.Errorf("Cost, error: %s", err)
				continue
			***REMOVED***
			if actual != cost ***REMOVED***
				t.Errorf("Cost, expected: %d, actual: %d", cost, actual)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	_, err := Cost([]byte("$a$a$" + suffix))
	if err == nil ***REMOVED***
		t.Errorf("Cost, malformed but no error returned")
	***REMOVED***
***REMOVED***

func TestCostValidationInHash(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	pass := []byte("mypassword")

	for c := 0; c < MinCost; c++ ***REMOVED***
		p, _ := newFromPassword(pass, c)
		if p.cost != DefaultCost ***REMOVED***
			t.Errorf("newFromPassword should default costs below %d to %d, but was %d", MinCost, DefaultCost, p.cost)
		***REMOVED***
	***REMOVED***

	p, _ := newFromPassword(pass, 14)
	if p.cost != 14 ***REMOVED***
		t.Errorf("newFromPassword should default cost to 14, but was %d", p.cost)
	***REMOVED***

	hp, _ := newFromHash(p.Hash())
	if p.cost != hp.cost ***REMOVED***
		t.Errorf("newFromHash should maintain the cost at %d, but was %d", p.cost, hp.cost)
	***REMOVED***

	_, err := newFromPassword(pass, 32)
	if err == nil ***REMOVED***
		t.Fatalf("newFromPassword: should return a cost error")
	***REMOVED***
	if err != InvalidCostError(32) ***REMOVED***
		t.Errorf("newFromPassword: should return cost error, got %#v", err)
	***REMOVED***
***REMOVED***

func TestCostReturnsWithLeadingZeroes(t *testing.T) ***REMOVED***
	hp, _ := newFromPassword([]byte("abcdefgh"), 7)
	cost := hp.Hash()[4:7]
	expected := []byte("07$")

	if !bytes.Equal(expected, cost) ***REMOVED***
		t.Errorf("single digit costs in hash should have leading zeros: was %v instead of %v", cost, expected)
	***REMOVED***
***REMOVED***

func TestMinorNotRequired(t *testing.T) ***REMOVED***
	noMinorHash := []byte("$2$10$XajjQvNhvvRt5GSeFk1xFeyqRrsxkhBkUiQeg0dt.wU1qD4aFDcga")
	h, err := newFromHash(noMinorHash)
	if err != nil ***REMOVED***
		t.Fatalf("No minor hash blew up: %s", err)
	***REMOVED***
	if h.minor != 0 ***REMOVED***
		t.Errorf("Should leave minor version at 0, but was %d", h.minor)
	***REMOVED***

	if !bytes.Equal(noMinorHash, h.Hash()) ***REMOVED***
		t.Errorf("Should generate hash %v, but created %v", noMinorHash, h.Hash())
	***REMOVED***
***REMOVED***

func BenchmarkEqual(b *testing.B) ***REMOVED***
	b.StopTimer()
	passwd := []byte("somepasswordyoulike")
	hash, _ := GenerateFromPassword(passwd, 10)
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		CompareHashAndPassword(hash, passwd)
	***REMOVED***
***REMOVED***

func BenchmarkGeneration(b *testing.B) ***REMOVED***
	b.StopTimer()
	passwd := []byte("mylongpassword1234")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		GenerateFromPassword(passwd, 10)
	***REMOVED***
***REMOVED***

// See Issue https://github.com/golang/go/issues/20425.
func TestNoSideEffectsFromCompare(t *testing.T) ***REMOVED***
	source := []byte("passw0rd123456")
	password := source[:len(source)-6]
	token := source[len(source)-6:]
	want := make([]byte, len(source))
	copy(want, source)

	wantHash := []byte("$2a$10$LK9XRuhNxHHCvjX3tdkRKei1QiCDUKrJRhZv7WWZPuQGRUM92rOUa")
	_ = CompareHashAndPassword(wantHash, password)

	got := bytes.Join([][]byte***REMOVED***password, token***REMOVED***, []byte(""))
	if !bytes.Equal(got, want) ***REMOVED***
		t.Errorf("got=%q want=%q", got, want)
	***REMOVED***
***REMOVED***
