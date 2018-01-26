// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package cmp

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	// Path is a list of PathSteps describing the sequence of operations to get
	// from some root type to the current position in the value tree.
	// The first Path element is always an operation-less PathStep that exists
	// simply to identify the initial type.
	//
	// When traversing structs with embedded structs, the embedded struct will
	// always be accessed as a field before traversing the fields of the
	// embedded struct themselves. That is, an exported field from the
	// embedded struct will never be accessed directly from the parent struct.
	Path []PathStep

	// PathStep is a union-type for specific operations to traverse
	// a value's tree structure. Users of this package never need to implement
	// these types as values of this type will be returned by this package.
	PathStep interface ***REMOVED***
		String() string
		Type() reflect.Type // Resulting type after performing the path step
		isPathStep()
	***REMOVED***

	// SliceIndex is an index operation on a slice or array at some index Key.
	SliceIndex interface ***REMOVED***
		PathStep
		Key() int // May return -1 if in a split state

		// SplitKeys returns the indexes for indexing into slices in the
		// x and y values, respectively. These indexes may differ due to the
		// insertion or removal of an element in one of the slices, causing
		// all of the indexes to be shifted. If an index is -1, then that
		// indicates that the element does not exist in the associated slice.
		//
		// Key is guaranteed to return -1 if and only if the indexes returned
		// by SplitKeys are not the same. SplitKeys will never return -1 for
		// both indexes.
		SplitKeys() (x int, y int)

		isSliceIndex()
	***REMOVED***
	// MapIndex is an index operation on a map at some index Key.
	MapIndex interface ***REMOVED***
		PathStep
		Key() reflect.Value
		isMapIndex()
	***REMOVED***
	// TypeAssertion represents a type assertion on an interface.
	TypeAssertion interface ***REMOVED***
		PathStep
		isTypeAssertion()
	***REMOVED***
	// StructField represents a struct field access on a field called Name.
	StructField interface ***REMOVED***
		PathStep
		Name() string
		Index() int
		isStructField()
	***REMOVED***
	// Indirect represents pointer indirection on the parent type.
	Indirect interface ***REMOVED***
		PathStep
		isIndirect()
	***REMOVED***
	// Transform is a transformation from the parent type to the current type.
	Transform interface ***REMOVED***
		PathStep
		Name() string
		Func() reflect.Value
		isTransform()
	***REMOVED***
)

func (pa *Path) push(s PathStep) ***REMOVED***
	*pa = append(*pa, s)
***REMOVED***

func (pa *Path) pop() ***REMOVED***
	*pa = (*pa)[:len(*pa)-1]
***REMOVED***

// Last returns the last PathStep in the Path.
// If the path is empty, this returns a non-nil PathStep that reports a nil Type.
func (pa Path) Last() PathStep ***REMOVED***
	if len(pa) > 0 ***REMOVED***
		return pa[len(pa)-1]
	***REMOVED***
	return pathStep***REMOVED******REMOVED***
***REMOVED***

// String returns the simplified path to a node.
// The simplified path only contains struct field accesses.
//
// For example:
//	MyMap.MySlices.MyField
func (pa Path) String() string ***REMOVED***
	var ss []string
	for _, s := range pa ***REMOVED***
		if _, ok := s.(*structField); ok ***REMOVED***
			ss = append(ss, s.String())
		***REMOVED***
	***REMOVED***
	return strings.TrimPrefix(strings.Join(ss, ""), ".")
***REMOVED***

// GoString returns the path to a specific node using Go syntax.
//
// For example:
//	(*root.MyMap["key"].(*mypkg.MyStruct).MySlices)[2][3].MyField
func (pa Path) GoString() string ***REMOVED***
	var ssPre, ssPost []string
	var numIndirect int
	for i, s := range pa ***REMOVED***
		var nextStep PathStep
		if i+1 < len(pa) ***REMOVED***
			nextStep = pa[i+1]
		***REMOVED***
		switch s := s.(type) ***REMOVED***
		case *indirect:
			numIndirect++
			pPre, pPost := "(", ")"
			switch nextStep.(type) ***REMOVED***
			case *indirect:
				continue // Next step is indirection, so let them batch up
			case *structField:
				numIndirect-- // Automatic indirection on struct fields
			case nil:
				pPre, pPost = "", "" // Last step; no need for parenthesis
			***REMOVED***
			if numIndirect > 0 ***REMOVED***
				ssPre = append(ssPre, pPre+strings.Repeat("*", numIndirect))
				ssPost = append(ssPost, pPost)
			***REMOVED***
			numIndirect = 0
			continue
		case *transform:
			ssPre = append(ssPre, s.trans.name+"(")
			ssPost = append(ssPost, ")")
			continue
		case *typeAssertion:
			// Elide type assertions immediately following a transform to
			// prevent overly verbose path printouts.
			// Some transforms return interface***REMOVED******REMOVED*** because of Go's lack of
			// generics, but typically take in and return the exact same
			// concrete type. Other times, the transform creates an anonymous
			// struct, which will be very verbose to print.
			if _, ok := nextStep.(*transform); ok ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		ssPost = append(ssPost, s.String())
	***REMOVED***
	for i, j := 0, len(ssPre)-1; i < j; i, j = i+1, j-1 ***REMOVED***
		ssPre[i], ssPre[j] = ssPre[j], ssPre[i]
	***REMOVED***
	return strings.Join(ssPre, "") + strings.Join(ssPost, "")
***REMOVED***

type (
	pathStep struct ***REMOVED***
		typ reflect.Type
	***REMOVED***

	sliceIndex struct ***REMOVED***
		pathStep
		xkey, ykey int
	***REMOVED***
	mapIndex struct ***REMOVED***
		pathStep
		key reflect.Value
	***REMOVED***
	typeAssertion struct ***REMOVED***
		pathStep
	***REMOVED***
	structField struct ***REMOVED***
		pathStep
		name string
		idx  int

		// These fields are used for forcibly accessing an unexported field.
		// pvx, pvy, and field are only valid if unexported is true.
		unexported bool
		force      bool                // Forcibly allow visibility
		pvx, pvy   reflect.Value       // Parent values
		field      reflect.StructField // Field information
	***REMOVED***
	indirect struct ***REMOVED***
		pathStep
	***REMOVED***
	transform struct ***REMOVED***
		pathStep
		trans *transformer
	***REMOVED***
)

func (ps pathStep) Type() reflect.Type ***REMOVED*** return ps.typ ***REMOVED***
func (ps pathStep) String() string ***REMOVED***
	if ps.typ == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	s := ps.typ.String()
	if s == "" || strings.ContainsAny(s, "***REMOVED******REMOVED***\n") ***REMOVED***
		return "root" // Type too simple or complex to print
	***REMOVED***
	return fmt.Sprintf("***REMOVED***%s***REMOVED***", s)
***REMOVED***

func (si sliceIndex) String() string ***REMOVED***
	switch ***REMOVED***
	case si.xkey == si.ykey:
		return fmt.Sprintf("[%d]", si.xkey)
	case si.ykey == -1:
		// [5->?] means "I don't know where X[5] went"
		return fmt.Sprintf("[%d->?]", si.xkey)
	case si.xkey == -1:
		// [?->3] means "I don't know where Y[3] came from"
		return fmt.Sprintf("[?->%d]", si.ykey)
	default:
		// [5->3] means "X[5] moved to Y[3]"
		return fmt.Sprintf("[%d->%d]", si.xkey, si.ykey)
	***REMOVED***
***REMOVED***
func (mi mapIndex) String() string      ***REMOVED*** return fmt.Sprintf("[%#v]", mi.key) ***REMOVED***
func (ta typeAssertion) String() string ***REMOVED*** return fmt.Sprintf(".(%v)", ta.typ) ***REMOVED***
func (sf structField) String() string   ***REMOVED*** return fmt.Sprintf(".%s", sf.name) ***REMOVED***
func (in indirect) String() string      ***REMOVED*** return "*" ***REMOVED***
func (tf transform) String() string     ***REMOVED*** return fmt.Sprintf("%s()", tf.trans.name) ***REMOVED***

func (si sliceIndex) Key() int ***REMOVED***
	if si.xkey != si.ykey ***REMOVED***
		return -1
	***REMOVED***
	return si.xkey
***REMOVED***
func (si sliceIndex) SplitKeys() (x, y int) ***REMOVED*** return si.xkey, si.ykey ***REMOVED***
func (mi mapIndex) Key() reflect.Value      ***REMOVED*** return mi.key ***REMOVED***
func (sf structField) Name() string         ***REMOVED*** return sf.name ***REMOVED***
func (sf structField) Index() int           ***REMOVED*** return sf.idx ***REMOVED***
func (tf transform) Name() string           ***REMOVED*** return tf.trans.name ***REMOVED***
func (tf transform) Func() reflect.Value    ***REMOVED*** return tf.trans.fnc ***REMOVED***

func (pathStep) isPathStep()           ***REMOVED******REMOVED***
func (sliceIndex) isSliceIndex()       ***REMOVED******REMOVED***
func (mapIndex) isMapIndex()           ***REMOVED******REMOVED***
func (typeAssertion) isTypeAssertion() ***REMOVED******REMOVED***
func (structField) isStructField()     ***REMOVED******REMOVED***
func (indirect) isIndirect()           ***REMOVED******REMOVED***
func (transform) isTransform()         ***REMOVED******REMOVED***

var (
	_ SliceIndex    = sliceIndex***REMOVED******REMOVED***
	_ MapIndex      = mapIndex***REMOVED******REMOVED***
	_ TypeAssertion = typeAssertion***REMOVED******REMOVED***
	_ StructField   = structField***REMOVED******REMOVED***
	_ Indirect      = indirect***REMOVED******REMOVED***
	_ Transform     = transform***REMOVED******REMOVED***

	_ PathStep = sliceIndex***REMOVED******REMOVED***
	_ PathStep = mapIndex***REMOVED******REMOVED***
	_ PathStep = typeAssertion***REMOVED******REMOVED***
	_ PathStep = structField***REMOVED******REMOVED***
	_ PathStep = indirect***REMOVED******REMOVED***
	_ PathStep = transform***REMOVED******REMOVED***
)

// isExported reports whether the identifier is exported.
func isExported(id string) bool ***REMOVED***
	r, _ := utf8.DecodeRuneInString(id)
	return unicode.IsUpper(r)
***REMOVED***

// isValid reports whether the identifier is valid.
// Empty and underscore-only strings are not valid.
func isValid(id string) bool ***REMOVED***
	ok := id != "" && id != "_"
	for j, c := range id ***REMOVED***
		ok = ok && (j > 0 || !unicode.IsDigit(c))
		ok = ok && (c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c))
	***REMOVED***
	return ok
***REMOVED***
