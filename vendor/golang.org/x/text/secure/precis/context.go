// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import "errors"

// This file contains tables and code related to context rules.

type catBitmap uint16

const (
	// These bits, once set depending on the current value, are never unset.
	bJapanese catBitmap = 1 << iota
	bArabicIndicDigit
	bExtendedArabicIndicDigit

	// These bits are set on each iteration depending on the current value.
	bJoinStart
	bJoinMid
	bJoinEnd
	bVirama
	bLatinSmallL
	bGreek
	bHebrew

	// These bits indicated which of the permanent bits need to be set at the
	// end of the checks.
	bMustHaveJapn

	permanent = bJapanese | bArabicIndicDigit | bExtendedArabicIndicDigit | bMustHaveJapn
)

const finalShift = 10

var errContext = errors.New("precis: contextual rule violated")

func init() ***REMOVED***
	// Programmatically set these required bits as, manually setting them seems
	// too error prone.
	for i, ct := range categoryTransitions ***REMOVED***
		categoryTransitions[i].keep |= permanent
		categoryTransitions[i].accept |= ct.term
	***REMOVED***
***REMOVED***

var categoryTransitions = []struct ***REMOVED***
	keep catBitmap // mask selecting which bits to keep from the previous state
	set  catBitmap // mask for which bits to set for this transition

	// These bitmaps are used for rules that require lookahead.
	// term&accept == term must be true, which is enforced programmatically.
	term   catBitmap // bits accepted as termination condition
	accept catBitmap // bits that pass, but not sufficient as termination

	// The rule function cannot take a *context as an argument, as it would
	// cause the context to escape, adding significant overhead.
	rule func(beforeBits catBitmap) (doLookahead bool, err error)
***REMOVED******REMOVED***
	joiningL:          ***REMOVED***set: bJoinStart***REMOVED***,
	joiningD:          ***REMOVED***set: bJoinStart | bJoinEnd***REMOVED***,
	joiningT:          ***REMOVED***keep: bJoinStart, set: bJoinMid***REMOVED***,
	joiningR:          ***REMOVED***set: bJoinEnd***REMOVED***,
	viramaModifier:    ***REMOVED***set: bVirama***REMOVED***,
	viramaJoinT:       ***REMOVED***set: bVirama | bJoinMid***REMOVED***,
	latinSmallL:       ***REMOVED***set: bLatinSmallL***REMOVED***,
	greek:             ***REMOVED***set: bGreek***REMOVED***,
	greekJoinT:        ***REMOVED***set: bGreek | bJoinMid***REMOVED***,
	hebrew:            ***REMOVED***set: bHebrew***REMOVED***,
	hebrewJoinT:       ***REMOVED***set: bHebrew | bJoinMid***REMOVED***,
	japanese:          ***REMOVED***set: bJapanese***REMOVED***,
	katakanaMiddleDot: ***REMOVED***set: bMustHaveJapn***REMOVED***,

	zeroWidthNonJoiner: ***REMOVED***
		term:   bJoinEnd,
		accept: bJoinMid,
		rule: func(before catBitmap) (doLookAhead bool, err error) ***REMOVED***
			if before&bVirama != 0 ***REMOVED***
				return false, nil
			***REMOVED***
			if before&bJoinStart == 0 ***REMOVED***
				return false, errContext
			***REMOVED***
			return true, nil
		***REMOVED***,
	***REMOVED***,
	zeroWidthJoiner: ***REMOVED***
		rule: func(before catBitmap) (doLookAhead bool, err error) ***REMOVED***
			if before&bVirama == 0 ***REMOVED***
				err = errContext
			***REMOVED***
			return false, err
		***REMOVED***,
	***REMOVED***,
	middleDot: ***REMOVED***
		term: bLatinSmallL,
		rule: func(before catBitmap) (doLookAhead bool, err error) ***REMOVED***
			if before&bLatinSmallL == 0 ***REMOVED***
				return false, errContext
			***REMOVED***
			return true, nil
		***REMOVED***,
	***REMOVED***,
	greekLowerNumeralSign: ***REMOVED***
		set:  bGreek,
		term: bGreek,
		rule: func(before catBitmap) (doLookAhead bool, err error) ***REMOVED***
			return true, nil
		***REMOVED***,
	***REMOVED***,
	hebrewPreceding: ***REMOVED***
		set: bHebrew,
		rule: func(before catBitmap) (doLookAhead bool, err error) ***REMOVED***
			if before&bHebrew == 0 ***REMOVED***
				err = errContext
			***REMOVED***
			return false, err
		***REMOVED***,
	***REMOVED***,
	arabicIndicDigit: ***REMOVED***
		set: bArabicIndicDigit,
		rule: func(before catBitmap) (doLookAhead bool, err error) ***REMOVED***
			if before&bExtendedArabicIndicDigit != 0 ***REMOVED***
				err = errContext
			***REMOVED***
			return false, err
		***REMOVED***,
	***REMOVED***,
	extendedArabicIndicDigit: ***REMOVED***
		set: bExtendedArabicIndicDigit,
		rule: func(before catBitmap) (doLookAhead bool, err error) ***REMOVED***
			if before&bArabicIndicDigit != 0 ***REMOVED***
				err = errContext
			***REMOVED***
			return false, err
		***REMOVED***,
	***REMOVED***,
***REMOVED***
