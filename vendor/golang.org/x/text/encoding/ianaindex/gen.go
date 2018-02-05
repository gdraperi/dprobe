// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/internal/gen"
)

type registry struct ***REMOVED***
	XMLName  xml.Name `xml:"registry"`
	Updated  string   `xml:"updated"`
	Registry []struct ***REMOVED***
		ID     string `xml:"id,attr"`
		Record []struct ***REMOVED***
			Name string `xml:"name"`
			Xref []struct ***REMOVED***
				Type string `xml:"type,attr"`
				Data string `xml:"data,attr"`
			***REMOVED*** `xml:"xref"`
			Desc struct ***REMOVED***
				Data string `xml:",innerxml"`
			***REMOVED*** `xml:"description,"`
			MIB   string   `xml:"value"`
			Alias []string `xml:"alias"`
			MIME  string   `xml:"preferred_alias"`
		***REMOVED*** `xml:"record"`
	***REMOVED*** `xml:"registry"`
***REMOVED***

func main() ***REMOVED***
	r := gen.OpenIANAFile("assignments/character-sets/character-sets.xml")
	reg := &registry***REMOVED******REMOVED***
	if err := xml.NewDecoder(r).Decode(&reg); err != nil && err != io.EOF ***REMOVED***
		log.Fatalf("Error decoding charset registry: %v", err)
	***REMOVED***
	if len(reg.Registry) == 0 || reg.Registry[0].ID != "character-sets-1" ***REMOVED***
		log.Fatalf("Unexpected ID %s", reg.Registry[0].ID)
	***REMOVED***

	x := &indexInfo***REMOVED******REMOVED***

	for _, rec := range reg.Registry[0].Record ***REMOVED***
		mib := identifier.MIB(parseInt(rec.MIB))
		x.addEntry(mib, rec.Name)
		for _, a := range rec.Alias ***REMOVED***
			a = strings.Split(a, " ")[0] // strip comments.
			x.addAlias(a, mib)
			// MIB name aliases are prefixed with a "cs" (character set) in the
			// registry to identify them as display names and to ensure that
			// the name starts with a lowercase letter in case it is used as
			// an identifier. We remove it to be left with a nice clean name.
			if strings.HasPrefix(a, "cs") ***REMOVED***
				x.setName(2, a[2:])
			***REMOVED***
		***REMOVED***
		if rec.MIME != "" ***REMOVED***
			x.addAlias(rec.MIME, mib)
			x.setName(1, rec.MIME)
		***REMOVED***
	***REMOVED***

	w := gen.NewCodeWriter()

	fmt.Fprintln(w, `import "golang.org/x/text/encoding/internal/identifier"`)

	writeIndex(w, x)

	w.WriteGoFile("tables.go", "ianaindex")
***REMOVED***

type alias struct ***REMOVED***
	name string
	mib  identifier.MIB
***REMOVED***

type indexInfo struct ***REMOVED***
	// compacted index from code to MIB
	codeToMIB []identifier.MIB
	alias     []alias
	names     [][3]string
***REMOVED***

func (ii *indexInfo) Len() int ***REMOVED***
	return len(ii.codeToMIB)
***REMOVED***

func (ii *indexInfo) Less(a, b int) bool ***REMOVED***
	return ii.codeToMIB[a] < ii.codeToMIB[b]
***REMOVED***

func (ii *indexInfo) Swap(a, b int) ***REMOVED***
	ii.codeToMIB[a], ii.codeToMIB[b] = ii.codeToMIB[b], ii.codeToMIB[a]
	// Co-sort the names.
	ii.names[a], ii.names[b] = ii.names[b], ii.names[a]
***REMOVED***

func (ii *indexInfo) setName(i int, name string) ***REMOVED***
	ii.names[len(ii.names)-1][i] = name
***REMOVED***

func (ii *indexInfo) addEntry(mib identifier.MIB, name string) ***REMOVED***
	ii.names = append(ii.names, [3]string***REMOVED***name, name, name***REMOVED***)
	ii.addAlias(name, mib)
	ii.codeToMIB = append(ii.codeToMIB, mib)
***REMOVED***

func (ii *indexInfo) addAlias(name string, mib identifier.MIB) ***REMOVED***
	// Don't add duplicates for the same mib. Adding duplicate aliases for
	// different MIBs will cause the compiler to barf on an invalid map: great!.
	for i := len(ii.alias) - 1; i >= 0 && ii.alias[i].mib == mib; i-- ***REMOVED***
		if ii.alias[i].name == name ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	ii.alias = append(ii.alias, alias***REMOVED***name, mib***REMOVED***)
	lower := strings.ToLower(name)
	if lower != name ***REMOVED***
		ii.addAlias(lower, mib)
	***REMOVED***
***REMOVED***

const maxMIMENameLen = '0' - 1 // officially 40, but we leave some buffer.

func writeIndex(w *gen.CodeWriter, x *indexInfo) ***REMOVED***
	sort.Stable(x)

	// Write constants.
	fmt.Fprintln(w, "const (")
	for i, m := range x.codeToMIB ***REMOVED***
		if i == 0 ***REMOVED***
			fmt.Fprintf(w, "enc%d = iota\n", m)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(w, "enc%d\n", m)
		***REMOVED***
	***REMOVED***
	fmt.Fprintln(w, "numIANA")
	fmt.Fprintln(w, ")")

	w.WriteVar("ianaToMIB", x.codeToMIB)

	var ianaNames, mibNames []string
	for _, names := range x.names ***REMOVED***
		n := names[0]
		if names[0] != names[1] ***REMOVED***
			// MIME names are mostly identical to IANA names. We share the
			// tables by setting the first byte of the string to an index into
			// the string itself (< maxMIMENameLen) to the IANA name. The MIME
			// name immediately follows the index.
			x := len(names[1]) + 1
			if x > maxMIMENameLen ***REMOVED***
				log.Fatalf("MIME name length (%d) > %d", x, maxMIMENameLen)
			***REMOVED***
			n = string(x) + names[1] + names[0]
		***REMOVED***
		ianaNames = append(ianaNames, n)
		mibNames = append(mibNames, names[2])
	***REMOVED***

	w.WriteVar("ianaNames", ianaNames)
	w.WriteVar("mibNames", mibNames)

	w.WriteComment(`
	TODO: Instead of using a map, we could use binary search strings doing
	on-the fly lower-casing per character. This allows to always avoid
	allocation and will be considerably more compact.`)
	fmt.Fprintln(w, "var ianaAliases = map[string]int***REMOVED***")
	for _, a := range x.alias ***REMOVED***
		fmt.Fprintf(w, "%q: enc%d,\n", a.name, a.mib)
	***REMOVED***
	fmt.Fprintln(w, "***REMOVED***")
***REMOVED***

func parseInt(s string) int ***REMOVED***
	x, err := strconv.ParseInt(s, 10, 64)
	if err != nil ***REMOVED***
		log.Fatalf("Could not parse integer: %v", err)
	***REMOVED***
	return int(x)
***REMOVED***
