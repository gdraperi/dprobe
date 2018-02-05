// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"unicode/utf16"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// Input holds an input string in both UTF-8 and UTF-16 format.
type Input struct ***REMOVED***
	index int // used for restoring to original random order
	UTF8  []byte
	UTF16 []uint16
	key   []byte // used for sorting
***REMOVED***

func (i Input) String() string ***REMOVED***
	return string(i.UTF8)
***REMOVED***

func makeInput(s8 []byte, s16 []uint16) Input ***REMOVED***
	return Input***REMOVED***UTF8: s8, UTF16: s16***REMOVED***
***REMOVED***

func makeInputString(s string) Input ***REMOVED***
	return Input***REMOVED***
		UTF8:  []byte(s),
		UTF16: utf16.Encode([]rune(s)),
	***REMOVED***
***REMOVED***

// Collator is an interface for architecture-specific implementations of collation.
type Collator interface ***REMOVED***
	// Key generates a sort key for the given input.  Implemenations
	// may return nil if a collator does not support sort keys.
	Key(s Input) []byte

	// Compare returns -1 if a < b, 1 if a > b and 0 if a == b.
	Compare(a, b Input) int
***REMOVED***

// CollatorFactory creates a Collator for a given language tag.
type CollatorFactory struct ***REMOVED***
	name        string
	makeFn      func(tag string) (Collator, error)
	description string
***REMOVED***

var collators = []CollatorFactory***REMOVED******REMOVED***

// AddFactory registers f as a factory for an implementation of Collator.
func AddFactory(f CollatorFactory) ***REMOVED***
	collators = append(collators, f)
***REMOVED***

func getCollator(name, locale string) Collator ***REMOVED***
	for _, f := range collators ***REMOVED***
		if f.name == name ***REMOVED***
			col, err := f.makeFn(locale)
			if err != nil ***REMOVED***
				log.Fatal(err)
			***REMOVED***
			return col
		***REMOVED***
	***REMOVED***
	log.Fatalf("collator of type %q not found", name)
	return nil
***REMOVED***

// goCollator is an implemention of Collator using go's own collator.
type goCollator struct ***REMOVED***
	c   *collate.Collator
	buf collate.Buffer
***REMOVED***

func init() ***REMOVED***
	AddFactory(CollatorFactory***REMOVED***"go", newGoCollator, "Go's native collator implementation."***REMOVED***)
***REMOVED***

func newGoCollator(loc string) (Collator, error) ***REMOVED***
	c := &goCollator***REMOVED***c: collate.New(language.Make(loc))***REMOVED***
	return c, nil
***REMOVED***

func (c *goCollator) Key(b Input) []byte ***REMOVED***
	return c.c.Key(&c.buf, b.UTF8)
***REMOVED***

func (c *goCollator) Compare(a, b Input) int ***REMOVED***
	return c.c.Compare(a.UTF8, b.UTF8)
***REMOVED***
