// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml_test

import (
	"encoding/xml"
	"fmt"
	"os"
)

func ExampleMarshalIndent() ***REMOVED***
	type Address struct ***REMOVED***
		City, State string
	***REMOVED***
	type Person struct ***REMOVED***
		XMLName   xml.Name `xml:"person"`
		Id        int      `xml:"id,attr"`
		FirstName string   `xml:"name>first"`
		LastName  string   `xml:"name>last"`
		Age       int      `xml:"age"`
		Height    float32  `xml:"height,omitempty"`
		Married   bool
		Address
		Comment string `xml:",comment"`
	***REMOVED***

	v := &Person***REMOVED***Id: 13, FirstName: "John", LastName: "Doe", Age: 42***REMOVED***
	v.Comment = " Need more details. "
	v.Address = Address***REMOVED***"Hanga Roa", "Easter Island"***REMOVED***

	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil ***REMOVED***
		fmt.Printf("error: %v\n", err)
	***REMOVED***

	os.Stdout.Write(output)
	// Output:
	//   <person id="13">
	//       <name>
	//           <first>John</first>
	//           <last>Doe</last>
	//       </name>
	//       <age>42</age>
	//       <Married>false</Married>
	//       <City>Hanga Roa</City>
	//       <State>Easter Island</State>
	//       <!-- Need more details. -->
	//   </person>
***REMOVED***

func ExampleEncoder() ***REMOVED***
	type Address struct ***REMOVED***
		City, State string
	***REMOVED***
	type Person struct ***REMOVED***
		XMLName   xml.Name `xml:"person"`
		Id        int      `xml:"id,attr"`
		FirstName string   `xml:"name>first"`
		LastName  string   `xml:"name>last"`
		Age       int      `xml:"age"`
		Height    float32  `xml:"height,omitempty"`
		Married   bool
		Address
		Comment string `xml:",comment"`
	***REMOVED***

	v := &Person***REMOVED***Id: 13, FirstName: "John", LastName: "Doe", Age: 42***REMOVED***
	v.Comment = " Need more details. "
	v.Address = Address***REMOVED***"Hanga Roa", "Easter Island"***REMOVED***

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil ***REMOVED***
		fmt.Printf("error: %v\n", err)
	***REMOVED***

	// Output:
	//   <person id="13">
	//       <name>
	//           <first>John</first>
	//           <last>Doe</last>
	//       </name>
	//       <age>42</age>
	//       <Married>false</Married>
	//       <City>Hanga Roa</City>
	//       <State>Easter Island</State>
	//       <!-- Need more details. -->
	//   </person>
***REMOVED***

// This example demonstrates unmarshaling an XML excerpt into a value with
// some preset fields. Note that the Phone field isn't modified and that
// the XML <Company> element is ignored. Also, the Groups field is assigned
// considering the element path provided in its tag.
func ExampleUnmarshal() ***REMOVED***
	type Email struct ***REMOVED***
		Where string `xml:"where,attr"`
		Addr  string
	***REMOVED***
	type Address struct ***REMOVED***
		City, State string
	***REMOVED***
	type Result struct ***REMOVED***
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"FullName"`
		Phone   string
		Email   []Email
		Groups  []string `xml:"Group>Value"`
		Address
	***REMOVED***
	v := Result***REMOVED***Name: "none", Phone: "none"***REMOVED***

	data := `
		<Person>
			<FullName>Grace R. Emlin</FullName>
			<Company>Example Inc.</Company>
			<Email where="home">
				<Addr>gre@example.com</Addr>
			</Email>
			<Email where='work'>
				<Addr>gre@work.com</Addr>
			</Email>
			<Group>
				<Value>Friends</Value>
				<Value>Squash</Value>
			</Group>
			<City>Hanga Roa</City>
			<State>Easter Island</State>
		</Person>
	`
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil ***REMOVED***
		fmt.Printf("error: %v", err)
		return
	***REMOVED***
	fmt.Printf("XMLName: %#v\n", v.XMLName)
	fmt.Printf("Name: %q\n", v.Name)
	fmt.Printf("Phone: %q\n", v.Phone)
	fmt.Printf("Email: %v\n", v.Email)
	fmt.Printf("Groups: %v\n", v.Groups)
	fmt.Printf("Address: %v\n", v.Address)
	// Output:
	// XMLName: xml.Name***REMOVED***Space:"", Local:"Person"***REMOVED***
	// Name: "Grace R. Emlin"
	// Phone: "none"
	// Email: [***REMOVED***home gre@example.com***REMOVED*** ***REMOVED***work gre@work.com***REMOVED***]
	// Groups: [Friends Squash]
	// Address: ***REMOVED***Hanga Roa Easter Island***REMOVED***
***REMOVED***
