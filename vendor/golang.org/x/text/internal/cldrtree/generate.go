// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldrtree

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/internal/gen"
)

func generate(b *Builder, t *Tree, w *gen.CodeWriter) error ***REMOVED***
	fmt.Fprintln(w, `import "golang.org/x/text/internal/cldrtree"`)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "var tree = &cldrtree.Tree***REMOVED***locales, indices, buckets***REMOVED***\n\n")

	w.WriteComment("Path values:\n" + b.stats())
	fmt.Fprintln(w)

	// Generate enum types.
	for _, e := range b.enums ***REMOVED***
		// Build enum types.
		w.WriteComment("%s specifies a property of a CLDR field.", e.name)
		fmt.Fprintf(w, "type %s uint16\n", e.name)
	***REMOVED***

	d, err := getEnumData(b)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	fmt.Fprintln(w, "const (")
	for i, k := range d.keys ***REMOVED***
		fmt.Fprintf(w, "%s %s = %d // %s\n", toCamel(k), d.enums[i], d.m[k], k)
	***REMOVED***
	fmt.Fprintln(w, ")")

	w.WriteVar("locales", t.Locales)
	w.WriteVar("indices", t.Indices)

	// Generate string buckets.
	fmt.Fprintln(w, "var buckets = []string***REMOVED***")
	for i := range t.Buckets ***REMOVED***
		fmt.Fprintf(w, "bucket%d,\n", i)
	***REMOVED***
	fmt.Fprint(w, "***REMOVED***\n\n")
	w.Size += int(reflect.TypeOf("").Size()) * len(t.Buckets)

	// Generate string buckets.
	for i, bucket := range t.Buckets ***REMOVED***
		w.WriteVar(fmt.Sprint("bucket", i), bucket)
	***REMOVED***
	return nil
***REMOVED***

func generateTestData(b *Builder, w *gen.CodeWriter) error ***REMOVED***
	d, err := getEnumData(b)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	fmt.Fprintln(w)
	fmt.Fprintln(w, "var enumMap = map[string]uint16***REMOVED***")
	fmt.Fprintln(w, `"": 0,`)
	for _, k := range d.keys ***REMOVED***
		fmt.Fprintf(w, "%q: %d,\n", k, d.m[k])
	***REMOVED***
	fmt.Fprintln(w, "***REMOVED***")
	return nil
***REMOVED***

func toCamel(s string) string ***REMOVED***
	p := strings.Split(s, "-")
	for i, s := range p[1:] ***REMOVED***
		p[i+1] = strings.Title(s)
	***REMOVED***
	return strings.Replace(strings.Join(p, ""), "/", "", -1)
***REMOVED***

func (b *Builder) stats() string ***REMOVED***
	w := &bytes.Buffer***REMOVED******REMOVED***

	b.rootMeta.validate()
	for _, es := range b.enums ***REMOVED***
		fmt.Fprintf(w, "<%s>\n", es.name)
		printEnumValues(w, es, 1, nil)
	***REMOVED***
	fmt.Fprintln(w)
	printEnums(w, b.rootMeta.typeInfo, 0)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Nr elem:           ", len(b.strToBucket))
	fmt.Fprintln(w, "uniqued size:      ", b.size)
	fmt.Fprintln(w, "total string size: ", b.sizeAll)
	fmt.Fprintln(w, "bucket waste:      ", b.bucketWaste)

	return w.String()
***REMOVED***

func printEnums(w io.Writer, s *typeInfo, indent int) ***REMOVED***
	idStr := strings.Repeat("  ", indent) + "- "
	e := s.enum
	if e == nil ***REMOVED***
		if len(s.entries) > 0 ***REMOVED***
			panic(fmt.Errorf("has entries but no enum values: %#v", s.entries))
		***REMOVED***
		return
	***REMOVED***
	if e.name != "" ***REMOVED***
		fmt.Fprintf(w, "%s<%s>\n", idStr, e.name)
	***REMOVED*** else ***REMOVED***
		printEnumValues(w, e, indent, s)
	***REMOVED***
	if s.sharedKeys() ***REMOVED***
		for _, v := range s.entries ***REMOVED***
			printEnums(w, v, indent+1)
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func printEnumValues(w io.Writer, e *enum, indent int, info *typeInfo) ***REMOVED***
	idStr := strings.Repeat("  ", indent) + "- "
	for i := 0; i < len(e.keys); i++ ***REMOVED***
		fmt.Fprint(w, idStr)
		k := e.keys[i]
		if u, err := strconv.ParseUint(k, 10, 16); err == nil ***REMOVED***
			fmt.Fprintf(w, "%s", k)
			// Skip contiguous integers
			var v, last uint64
			for i++; i < len(e.keys); i++ ***REMOVED***
				k = e.keys[i]
				if v, err = strconv.ParseUint(k, 10, 16); err != nil ***REMOVED***
					break
				***REMOVED***
				last = v
			***REMOVED***
			if u < last ***REMOVED***
				fmt.Fprintf(w, `..%d`, last)
			***REMOVED***
			fmt.Fprintln(w)
			if err != nil ***REMOVED***
				fmt.Fprintf(w, "%s%s\n", idStr, k)
			***REMOVED***
		***REMOVED*** else if k == "" ***REMOVED***
			fmt.Fprintln(w, `""`)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(w, "%s\n", k)
		***REMOVED***
		if info != nil && !info.sharedKeys() ***REMOVED***
			if e := info.entries[enumIndex(i)]; e != nil ***REMOVED***
				printEnums(w, e, indent+1)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func getEnumData(b *Builder) (*enumData, error) ***REMOVED***
	d := &enumData***REMOVED***m: map[string]int***REMOVED******REMOVED******REMOVED***
	if errStr := d.insert(b.rootMeta.typeInfo); errStr != "" ***REMOVED***
		// TODO: consider returning the error.
		return nil, fmt.Errorf("cldrtree: %s", errStr)
	***REMOVED***
	return d, nil
***REMOVED***

type enumData struct ***REMOVED***
	m     map[string]int
	keys  []string
	enums []string
***REMOVED***

func (d *enumData) insert(t *typeInfo) (errStr string) ***REMOVED***
	e := t.enum
	if e == nil ***REMOVED***
		return ""
	***REMOVED***
	for i, k := range e.keys ***REMOVED***
		if _, err := strconv.ParseUint(k, 10, 16); err == nil ***REMOVED***
			// We don't include any enum that has integer values.
			break
		***REMOVED***
		if v, ok := d.m[k]; ok ***REMOVED***
			if v != i ***REMOVED***
				return fmt.Sprintf("%q has value %d and %d", k, i, v)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			d.m[k] = i
			if k != "" ***REMOVED***
				d.keys = append(d.keys, k)
				d.enums = append(d.enums, e.name)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for i := range t.enum.keys ***REMOVED***
		if e := t.entries[enumIndex(i)]; e != nil ***REMOVED***
			if errStr := d.insert(e); errStr != "" ***REMOVED***
				return fmt.Sprintf("%q>%v", t.enum.keys[i], errStr)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***
