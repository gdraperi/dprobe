// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// typeInfo holds details for the xml representation of a type.
type typeInfo struct ***REMOVED***
	xmlname *fieldInfo
	fields  []fieldInfo
***REMOVED***

// fieldInfo holds details for the xml representation of a single field.
type fieldInfo struct ***REMOVED***
	idx     []int
	name    string
	xmlns   string
	flags   fieldFlags
	parents []string
***REMOVED***

type fieldFlags int

const (
	fElement fieldFlags = 1 << iota
	fAttr
	fCharData
	fInnerXml
	fComment
	fAny

	fOmitEmpty

	fMode = fElement | fAttr | fCharData | fInnerXml | fComment | fAny
)

var tinfoMap = make(map[reflect.Type]*typeInfo)
var tinfoLock sync.RWMutex

var nameType = reflect.TypeOf(Name***REMOVED******REMOVED***)

// getTypeInfo returns the typeInfo structure with details necessary
// for marshalling and unmarshalling typ.
func getTypeInfo(typ reflect.Type) (*typeInfo, error) ***REMOVED***
	tinfoLock.RLock()
	tinfo, ok := tinfoMap[typ]
	tinfoLock.RUnlock()
	if ok ***REMOVED***
		return tinfo, nil
	***REMOVED***
	tinfo = &typeInfo***REMOVED******REMOVED***
	if typ.Kind() == reflect.Struct && typ != nameType ***REMOVED***
		n := typ.NumField()
		for i := 0; i < n; i++ ***REMOVED***
			f := typ.Field(i)
			if f.PkgPath != "" || f.Tag.Get("xml") == "-" ***REMOVED***
				continue // Private field
			***REMOVED***

			// For embedded structs, embed its fields.
			if f.Anonymous ***REMOVED***
				t := f.Type
				if t.Kind() == reflect.Ptr ***REMOVED***
					t = t.Elem()
				***REMOVED***
				if t.Kind() == reflect.Struct ***REMOVED***
					inner, err := getTypeInfo(t)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					if tinfo.xmlname == nil ***REMOVED***
						tinfo.xmlname = inner.xmlname
					***REMOVED***
					for _, finfo := range inner.fields ***REMOVED***
						finfo.idx = append([]int***REMOVED***i***REMOVED***, finfo.idx...)
						if err := addFieldInfo(typ, tinfo, &finfo); err != nil ***REMOVED***
							return nil, err
						***REMOVED***
					***REMOVED***
					continue
				***REMOVED***
			***REMOVED***

			finfo, err := structFieldInfo(typ, &f)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if f.Name == "XMLName" ***REMOVED***
				tinfo.xmlname = finfo
				continue
			***REMOVED***

			// Add the field if it doesn't conflict with other fields.
			if err := addFieldInfo(typ, tinfo, finfo); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	tinfoLock.Lock()
	tinfoMap[typ] = tinfo
	tinfoLock.Unlock()
	return tinfo, nil
***REMOVED***

// structFieldInfo builds and returns a fieldInfo for f.
func structFieldInfo(typ reflect.Type, f *reflect.StructField) (*fieldInfo, error) ***REMOVED***
	finfo := &fieldInfo***REMOVED***idx: f.Index***REMOVED***

	// Split the tag from the xml namespace if necessary.
	tag := f.Tag.Get("xml")
	if i := strings.Index(tag, " "); i >= 0 ***REMOVED***
		finfo.xmlns, tag = tag[:i], tag[i+1:]
	***REMOVED***

	// Parse flags.
	tokens := strings.Split(tag, ",")
	if len(tokens) == 1 ***REMOVED***
		finfo.flags = fElement
	***REMOVED*** else ***REMOVED***
		tag = tokens[0]
		for _, flag := range tokens[1:] ***REMOVED***
			switch flag ***REMOVED***
			case "attr":
				finfo.flags |= fAttr
			case "chardata":
				finfo.flags |= fCharData
			case "innerxml":
				finfo.flags |= fInnerXml
			case "comment":
				finfo.flags |= fComment
			case "any":
				finfo.flags |= fAny
			case "omitempty":
				finfo.flags |= fOmitEmpty
			***REMOVED***
		***REMOVED***

		// Validate the flags used.
		valid := true
		switch mode := finfo.flags & fMode; mode ***REMOVED***
		case 0:
			finfo.flags |= fElement
		case fAttr, fCharData, fInnerXml, fComment, fAny:
			if f.Name == "XMLName" || tag != "" && mode != fAttr ***REMOVED***
				valid = false
			***REMOVED***
		default:
			// This will also catch multiple modes in a single field.
			valid = false
		***REMOVED***
		if finfo.flags&fMode == fAny ***REMOVED***
			finfo.flags |= fElement
		***REMOVED***
		if finfo.flags&fOmitEmpty != 0 && finfo.flags&(fElement|fAttr) == 0 ***REMOVED***
			valid = false
		***REMOVED***
		if !valid ***REMOVED***
			return nil, fmt.Errorf("xml: invalid tag in field %s of type %s: %q",
				f.Name, typ, f.Tag.Get("xml"))
		***REMOVED***
	***REMOVED***

	// Use of xmlns without a name is not allowed.
	if finfo.xmlns != "" && tag == "" ***REMOVED***
		return nil, fmt.Errorf("xml: namespace without name in field %s of type %s: %q",
			f.Name, typ, f.Tag.Get("xml"))
	***REMOVED***

	if f.Name == "XMLName" ***REMOVED***
		// The XMLName field records the XML element name. Don't
		// process it as usual because its name should default to
		// empty rather than to the field name.
		finfo.name = tag
		return finfo, nil
	***REMOVED***

	if tag == "" ***REMOVED***
		// If the name part of the tag is completely empty, get
		// default from XMLName of underlying struct if feasible,
		// or field name otherwise.
		if xmlname := lookupXMLName(f.Type); xmlname != nil ***REMOVED***
			finfo.xmlns, finfo.name = xmlname.xmlns, xmlname.name
		***REMOVED*** else ***REMOVED***
			finfo.name = f.Name
		***REMOVED***
		return finfo, nil
	***REMOVED***

	if finfo.xmlns == "" && finfo.flags&fAttr == 0 ***REMOVED***
		// If it's an element no namespace specified, get the default
		// from the XMLName of enclosing struct if possible.
		if xmlname := lookupXMLName(typ); xmlname != nil ***REMOVED***
			finfo.xmlns = xmlname.xmlns
		***REMOVED***
	***REMOVED***

	// Prepare field name and parents.
	parents := strings.Split(tag, ">")
	if parents[0] == "" ***REMOVED***
		parents[0] = f.Name
	***REMOVED***
	if parents[len(parents)-1] == "" ***REMOVED***
		return nil, fmt.Errorf("xml: trailing '>' in field %s of type %s", f.Name, typ)
	***REMOVED***
	finfo.name = parents[len(parents)-1]
	if len(parents) > 1 ***REMOVED***
		if (finfo.flags & fElement) == 0 ***REMOVED***
			return nil, fmt.Errorf("xml: %s chain not valid with %s flag", tag, strings.Join(tokens[1:], ","))
		***REMOVED***
		finfo.parents = parents[:len(parents)-1]
	***REMOVED***

	// If the field type has an XMLName field, the names must match
	// so that the behavior of both marshalling and unmarshalling
	// is straightforward and unambiguous.
	if finfo.flags&fElement != 0 ***REMOVED***
		ftyp := f.Type
		xmlname := lookupXMLName(ftyp)
		if xmlname != nil && xmlname.name != finfo.name ***REMOVED***
			return nil, fmt.Errorf("xml: name %q in tag of %s.%s conflicts with name %q in %s.XMLName",
				finfo.name, typ, f.Name, xmlname.name, ftyp)
		***REMOVED***
	***REMOVED***
	return finfo, nil
***REMOVED***

// lookupXMLName returns the fieldInfo for typ's XMLName field
// in case it exists and has a valid xml field tag, otherwise
// it returns nil.
func lookupXMLName(typ reflect.Type) (xmlname *fieldInfo) ***REMOVED***
	for typ.Kind() == reflect.Ptr ***REMOVED***
		typ = typ.Elem()
	***REMOVED***
	if typ.Kind() != reflect.Struct ***REMOVED***
		return nil
	***REMOVED***
	for i, n := 0, typ.NumField(); i < n; i++ ***REMOVED***
		f := typ.Field(i)
		if f.Name != "XMLName" ***REMOVED***
			continue
		***REMOVED***
		finfo, err := structFieldInfo(typ, &f)
		if finfo.name != "" && err == nil ***REMOVED***
			return finfo
		***REMOVED***
		// Also consider errors as a non-existent field tag
		// and let getTypeInfo itself report the error.
		break
	***REMOVED***
	return nil
***REMOVED***

func min(a, b int) int ***REMOVED***
	if a <= b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

// addFieldInfo adds finfo to tinfo.fields if there are no
// conflicts, or if conflicts arise from previous fields that were
// obtained from deeper embedded structures than finfo. In the latter
// case, the conflicting entries are dropped.
// A conflict occurs when the path (parent + name) to a field is
// itself a prefix of another path, or when two paths match exactly.
// It is okay for field paths to share a common, shorter prefix.
func addFieldInfo(typ reflect.Type, tinfo *typeInfo, newf *fieldInfo) error ***REMOVED***
	var conflicts []int
Loop:
	// First, figure all conflicts. Most working code will have none.
	for i := range tinfo.fields ***REMOVED***
		oldf := &tinfo.fields[i]
		if oldf.flags&fMode != newf.flags&fMode ***REMOVED***
			continue
		***REMOVED***
		if oldf.xmlns != "" && newf.xmlns != "" && oldf.xmlns != newf.xmlns ***REMOVED***
			continue
		***REMOVED***
		minl := min(len(newf.parents), len(oldf.parents))
		for p := 0; p < minl; p++ ***REMOVED***
			if oldf.parents[p] != newf.parents[p] ***REMOVED***
				continue Loop
			***REMOVED***
		***REMOVED***
		if len(oldf.parents) > len(newf.parents) ***REMOVED***
			if oldf.parents[len(newf.parents)] == newf.name ***REMOVED***
				conflicts = append(conflicts, i)
			***REMOVED***
		***REMOVED*** else if len(oldf.parents) < len(newf.parents) ***REMOVED***
			if newf.parents[len(oldf.parents)] == oldf.name ***REMOVED***
				conflicts = append(conflicts, i)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if newf.name == oldf.name ***REMOVED***
				conflicts = append(conflicts, i)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Without conflicts, add the new field and return.
	if conflicts == nil ***REMOVED***
		tinfo.fields = append(tinfo.fields, *newf)
		return nil
	***REMOVED***

	// If any conflict is shallower, ignore the new field.
	// This matches the Go field resolution on embedding.
	for _, i := range conflicts ***REMOVED***
		if len(tinfo.fields[i].idx) < len(newf.idx) ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	// Otherwise, if any of them is at the same depth level, it's an error.
	for _, i := range conflicts ***REMOVED***
		oldf := &tinfo.fields[i]
		if len(oldf.idx) == len(newf.idx) ***REMOVED***
			f1 := typ.FieldByIndex(oldf.idx)
			f2 := typ.FieldByIndex(newf.idx)
			return &TagPathError***REMOVED***typ, f1.Name, f1.Tag.Get("xml"), f2.Name, f2.Tag.Get("xml")***REMOVED***
		***REMOVED***
	***REMOVED***

	// Otherwise, the new field is shallower, and thus takes precedence,
	// so drop the conflicting fields from tinfo and append the new one.
	for c := len(conflicts) - 1; c >= 0; c-- ***REMOVED***
		i := conflicts[c]
		copy(tinfo.fields[i:], tinfo.fields[i+1:])
		tinfo.fields = tinfo.fields[:len(tinfo.fields)-1]
	***REMOVED***
	tinfo.fields = append(tinfo.fields, *newf)
	return nil
***REMOVED***

// A TagPathError represents an error in the unmarshalling process
// caused by the use of field tags with conflicting paths.
type TagPathError struct ***REMOVED***
	Struct       reflect.Type
	Field1, Tag1 string
	Field2, Tag2 string
***REMOVED***

func (e *TagPathError) Error() string ***REMOVED***
	return fmt.Sprintf("%s field %q with tag %q conflicts with field %q with tag %q", e.Struct, e.Field1, e.Tag1, e.Field2, e.Tag2)
***REMOVED***

// value returns v's field value corresponding to finfo.
// It's equivalent to v.FieldByIndex(finfo.idx), but initializes
// and dereferences pointers as necessary.
func (finfo *fieldInfo) value(v reflect.Value) reflect.Value ***REMOVED***
	for i, x := range finfo.idx ***REMOVED***
		if i > 0 ***REMOVED***
			t := v.Type()
			if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct ***REMOVED***
				if v.IsNil() ***REMOVED***
					v.Set(reflect.New(v.Type().Elem()))
				***REMOVED***
				v = v.Elem()
			***REMOVED***
		***REMOVED***
		v = v.Field(x)
	***REMOVED***
	return v
***REMOVED***
