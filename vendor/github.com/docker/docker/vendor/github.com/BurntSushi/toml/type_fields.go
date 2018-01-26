package toml

// Struct field handling is adapted from code in encoding/json:
//
// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the Go distribution.

import (
	"reflect"
	"sort"
	"sync"
)

// A field represents a single field found in a struct.
type field struct ***REMOVED***
	name  string       // the name of the field (`toml` tag included)
	tag   bool         // whether field has a `toml` tag
	index []int        // represents the depth of an anonymous field
	typ   reflect.Type // the type of the field
***REMOVED***

// byName sorts field by name, breaking ties with depth,
// then breaking ties with "name came from toml tag", then
// breaking ties with index sequence.
type byName []field

func (x byName) Len() int ***REMOVED*** return len(x) ***REMOVED***

func (x byName) Swap(i, j int) ***REMOVED*** x[i], x[j] = x[j], x[i] ***REMOVED***

func (x byName) Less(i, j int) bool ***REMOVED***
	if x[i].name != x[j].name ***REMOVED***
		return x[i].name < x[j].name
	***REMOVED***
	if len(x[i].index) != len(x[j].index) ***REMOVED***
		return len(x[i].index) < len(x[j].index)
	***REMOVED***
	if x[i].tag != x[j].tag ***REMOVED***
		return x[i].tag
	***REMOVED***
	return byIndex(x).Less(i, j)
***REMOVED***

// byIndex sorts field by index sequence.
type byIndex []field

func (x byIndex) Len() int ***REMOVED*** return len(x) ***REMOVED***

func (x byIndex) Swap(i, j int) ***REMOVED*** x[i], x[j] = x[j], x[i] ***REMOVED***

func (x byIndex) Less(i, j int) bool ***REMOVED***
	for k, xik := range x[i].index ***REMOVED***
		if k >= len(x[j].index) ***REMOVED***
			return false
		***REMOVED***
		if xik != x[j].index[k] ***REMOVED***
			return xik < x[j].index[k]
		***REMOVED***
	***REMOVED***
	return len(x[i].index) < len(x[j].index)
***REMOVED***

// typeFields returns a list of fields that TOML should recognize for the given
// type. The algorithm is breadth-first search over the set of structs to
// include - the top struct and then any reachable anonymous structs.
func typeFields(t reflect.Type) []field ***REMOVED***
	// Anonymous fields to explore at the current level and the next.
	current := []field***REMOVED******REMOVED***
	next := []field***REMOVED******REMOVED***typ: t***REMOVED******REMOVED***

	// Count of queued names for current level and the next.
	count := map[reflect.Type]int***REMOVED******REMOVED***
	nextCount := map[reflect.Type]int***REMOVED******REMOVED***

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool***REMOVED******REMOVED***

	// Fields found.
	var fields []field

	for len(next) > 0 ***REMOVED***
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int***REMOVED******REMOVED***

		for _, f := range current ***REMOVED***
			if visited[f.typ] ***REMOVED***
				continue
			***REMOVED***
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ ***REMOVED***
				sf := f.typ.Field(i)
				if sf.PkgPath != "" ***REMOVED*** // unexported
					continue
				***REMOVED***
				name := sf.Tag.Get("toml")
				if name == "-" ***REMOVED***
					continue
				***REMOVED***
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr ***REMOVED***
					// Follow pointer.
					ft = ft.Elem()
				***REMOVED***

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct ***REMOVED***
					tagged := name != ""
					if name == "" ***REMOVED***
						name = sf.Name
					***REMOVED***
					fields = append(fields, field***REMOVED***name, tagged, index, ft***REMOVED***)
					if count[f.typ] > 1 ***REMOVED***
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					***REMOVED***
					continue
				***REMOVED***

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 ***REMOVED***
					f := field***REMOVED***name: ft.Name(), index: index, typ: ft***REMOVED***
					next = append(next, f)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	sort.Sort(byName(fields))

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with TOML tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance ***REMOVED***
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ ***REMOVED***
			fj := fields[i+advance]
			if fj.name != name ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if advance == 1 ***REMOVED*** // Only one field with this name
			out = append(out, fi)
			continue
		***REMOVED***
		dominant, ok := dominantField(fields[i : i+advance])
		if ok ***REMOVED***
			out = append(out, dominant)
		***REMOVED***
	***REMOVED***

	fields = out
	sort.Sort(byIndex(fields))

	return fields
***REMOVED***

// dominantField looks through the fields, all of which are known to
// have the same name, to find the single field that dominates the
// others using Go's embedding rules, modified by the presence of
// TOML tags. If there are multiple top-level fields, the boolean
// will be false: This condition is an error in Go and we skip all
// the fields.
func dominantField(fields []field) (field, bool) ***REMOVED***
	// The fields are sorted in increasing index-length order. The winner
	// must therefore be one with the shortest index length. Drop all
	// longer entries, which is easy: just truncate the slice.
	length := len(fields[0].index)
	tagged := -1 // Index of first tagged field.
	for i, f := range fields ***REMOVED***
		if len(f.index) > length ***REMOVED***
			fields = fields[:i]
			break
		***REMOVED***
		if f.tag ***REMOVED***
			if tagged >= 0 ***REMOVED***
				// Multiple tagged fields at the same level: conflict.
				// Return no field.
				return field***REMOVED******REMOVED***, false
			***REMOVED***
			tagged = i
		***REMOVED***
	***REMOVED***
	if tagged >= 0 ***REMOVED***
		return fields[tagged], true
	***REMOVED***
	// All remaining fields have the same length. If there's more than one,
	// we have a conflict (two fields named "X" at the same level) and we
	// return no field.
	if len(fields) > 1 ***REMOVED***
		return field***REMOVED******REMOVED***, false
	***REMOVED***
	return fields[0], true
***REMOVED***

var fieldCache struct ***REMOVED***
	sync.RWMutex
	m map[reflect.Type][]field
***REMOVED***

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field ***REMOVED***
	fieldCache.RLock()
	f := fieldCache.m[t]
	fieldCache.RUnlock()
	if f != nil ***REMOVED***
		return f
	***REMOVED***

	// Compute fields without lock.
	// Might duplicate effort but won't hold other computations back.
	f = typeFields(t)
	if f == nil ***REMOVED***
		f = []field***REMOVED******REMOVED***
	***REMOVED***

	fieldCache.Lock()
	if fieldCache.m == nil ***REMOVED***
		fieldCache.m = map[reflect.Type][]field***REMOVED******REMOVED***
	***REMOVED***
	fieldCache.m[t] = f
	fieldCache.Unlock()
	return f
***REMOVED***
