// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"

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
				// Any []struct ***REMOVED***
				// 	Data string `xml:",chardata"`
				// ***REMOVED*** `xml:",any"`
				// Data string `xml:",chardata"`
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

	w := &bytes.Buffer***REMOVED******REMOVED***
	fmt.Fprintf(w, "const (\n")
	for _, rec := range reg.Registry[0].Record ***REMOVED***
		constName := ""
		for _, a := range rec.Alias ***REMOVED***
			if strings.HasPrefix(a, "cs") && strings.IndexByte(a, '-') == -1 ***REMOVED***
				// Some of the constant definitions have comments in them. Strip those.
				constName = strings.Title(strings.SplitN(a[2:], "\n", 2)[0])
			***REMOVED***
		***REMOVED***
		if constName == "" ***REMOVED***
			switch rec.MIB ***REMOVED***
			case "2085":
				constName = "HZGB2312" // Not listed as alias for some reason.
			default:
				log.Fatalf("No cs alias defined for %s.", rec.MIB)
			***REMOVED***
		***REMOVED***
		if rec.MIME != "" ***REMOVED***
			rec.MIME = fmt.Sprintf(" (MIME: %s)", rec.MIME)
		***REMOVED***
		fmt.Fprintf(w, "// %s is the MIB identifier with IANA name %s%s.\n//\n", constName, rec.Name, rec.MIME)
		if len(rec.Desc.Data) > 0 ***REMOVED***
			fmt.Fprint(w, "// ")
			d := xml.NewDecoder(strings.NewReader(rec.Desc.Data))
			inElem := true
			attr := ""
			for ***REMOVED***
				t, err := d.Token()
				if err != nil ***REMOVED***
					if err != io.EOF ***REMOVED***
						log.Fatal(err)
					***REMOVED***
					break
				***REMOVED***
				switch x := t.(type) ***REMOVED***
				case xml.CharData:
					attr = "" // Don't need attribute info.
					a := bytes.Split([]byte(x), []byte("\n"))
					for i, b := range a ***REMOVED***
						if b = bytes.TrimSpace(b); len(b) != 0 ***REMOVED***
							if !inElem && i > 0 ***REMOVED***
								fmt.Fprint(w, "\n// ")
							***REMOVED***
							inElem = false
							fmt.Fprintf(w, "%s ", string(b))
						***REMOVED***
					***REMOVED***
				case xml.StartElement:
					if x.Name.Local == "xref" ***REMOVED***
						inElem = true
						use := false
						for _, a := range x.Attr ***REMOVED***
							if a.Name.Local == "type" ***REMOVED***
								use = use || a.Value != "person"
							***REMOVED***
							if a.Name.Local == "data" && use ***REMOVED***
								attr = a.Value + " "
							***REMOVED***
						***REMOVED***
					***REMOVED***
				case xml.EndElement:
					inElem = false
					fmt.Fprint(w, attr)
				***REMOVED***
			***REMOVED***
			fmt.Fprint(w, "\n")
		***REMOVED***
		for _, x := range rec.Xref ***REMOVED***
			switch x.Type ***REMOVED***
			case "rfc":
				fmt.Fprintf(w, "// Reference: %s\n", strings.ToUpper(x.Data))
			case "uri":
				fmt.Fprintf(w, "// Reference: %s\n", x.Data)
			***REMOVED***
		***REMOVED***
		fmt.Fprintf(w, "%s MIB = %s\n", constName, rec.MIB)
		fmt.Fprintln(w)
	***REMOVED***
	fmt.Fprintln(w, ")")

	gen.WriteGoFile("mib.go", "identifier", w.Bytes())
***REMOVED***
