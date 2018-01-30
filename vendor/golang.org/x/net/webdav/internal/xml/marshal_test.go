// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type DriveType int

const (
	HyperDrive DriveType = iota
	ImprobabilityDrive
)

type Passenger struct ***REMOVED***
	Name   []string `xml:"name"`
	Weight float32  `xml:"weight"`
***REMOVED***

type Ship struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"spaceship"`

	Name      string       `xml:"name,attr"`
	Pilot     string       `xml:"pilot,attr"`
	Drive     DriveType    `xml:"drive"`
	Age       uint         `xml:"age"`
	Passenger []*Passenger `xml:"passenger"`
	secret    string
***REMOVED***

type NamedType string

type Port struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"port"`
	Type    string   `xml:"type,attr,omitempty"`
	Comment string   `xml:",comment"`
	Number  string   `xml:",chardata"`
***REMOVED***

type Domain struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"domain"`
	Country string   `xml:",attr,omitempty"`
	Name    []byte   `xml:",chardata"`
	Comment []byte   `xml:",comment"`
***REMOVED***

type Book struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"book"`
	Title   string   `xml:",chardata"`
***REMOVED***

type Event struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"event"`
	Year    int      `xml:",chardata"`
***REMOVED***

type Movie struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"movie"`
	Length  uint     `xml:",chardata"`
***REMOVED***

type Pi struct ***REMOVED***
	XMLName       struct***REMOVED******REMOVED*** `xml:"pi"`
	Approximation float32  `xml:",chardata"`
***REMOVED***

type Universe struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"universe"`
	Visible float64  `xml:",chardata"`
***REMOVED***

type Particle struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"particle"`
	HasMass bool     `xml:",chardata"`
***REMOVED***

type Departure struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED***  `xml:"departure"`
	When    time.Time `xml:",chardata"`
***REMOVED***

type SecretAgent struct ***REMOVED***
	XMLName   struct***REMOVED******REMOVED*** `xml:"agent"`
	Handle    string   `xml:"handle,attr"`
	Identity  string
	Obfuscate string `xml:",innerxml"`
***REMOVED***

type NestedItems struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"result"`
	Items   []string `xml:">item"`
	Item1   []string `xml:"Items>item1"`
***REMOVED***

type NestedOrder struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"result"`
	Field1  string   `xml:"parent>c"`
	Field2  string   `xml:"parent>b"`
	Field3  string   `xml:"parent>a"`
***REMOVED***

type MixedNested struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"result"`
	A       string   `xml:"parent1>a"`
	B       string   `xml:"b"`
	C       string   `xml:"parent1>parent2>c"`
	D       string   `xml:"parent1>d"`
***REMOVED***

type NilTest struct ***REMOVED***
	A interface***REMOVED******REMOVED*** `xml:"parent1>parent2>a"`
	B interface***REMOVED******REMOVED*** `xml:"parent1>b"`
	C interface***REMOVED******REMOVED*** `xml:"parent1>parent2>c"`
***REMOVED***

type Service struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"service"`
	Domain  *Domain  `xml:"host>domain"`
	Port    *Port    `xml:"host>port"`
	Extra1  interface***REMOVED******REMOVED***
	Extra2  interface***REMOVED******REMOVED*** `xml:"host>extra2"`
***REMOVED***

var nilStruct *Ship

type EmbedA struct ***REMOVED***
	EmbedC
	EmbedB EmbedB
	FieldA string
***REMOVED***

type EmbedB struct ***REMOVED***
	FieldB string
	*EmbedC
***REMOVED***

type EmbedC struct ***REMOVED***
	FieldA1 string `xml:"FieldA>A1"`
	FieldA2 string `xml:"FieldA>A2"`
	FieldB  string
	FieldC  string
***REMOVED***

type NameCasing struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"casing"`
	Xy      string
	XY      string
	XyA     string `xml:"Xy,attr"`
	XYA     string `xml:"XY,attr"`
***REMOVED***

type NamePrecedence struct ***REMOVED***
	XMLName     Name              `xml:"Parent"`
	FromTag     XMLNameWithoutTag `xml:"InTag"`
	FromNameVal XMLNameWithoutTag
	FromNameTag XMLNameWithTag
	InFieldName string
***REMOVED***

type XMLNameWithTag struct ***REMOVED***
	XMLName Name   `xml:"InXMLNameTag"`
	Value   string `xml:",chardata"`
***REMOVED***

type XMLNameWithNSTag struct ***REMOVED***
	XMLName Name   `xml:"ns InXMLNameWithNSTag"`
	Value   string `xml:",chardata"`
***REMOVED***

type XMLNameWithoutTag struct ***REMOVED***
	XMLName Name
	Value   string `xml:",chardata"`
***REMOVED***

type NameInField struct ***REMOVED***
	Foo Name `xml:"ns foo"`
***REMOVED***

type AttrTest struct ***REMOVED***
	Int   int     `xml:",attr"`
	Named int     `xml:"int,attr"`
	Float float64 `xml:",attr"`
	Uint8 uint8   `xml:",attr"`
	Bool  bool    `xml:",attr"`
	Str   string  `xml:",attr"`
	Bytes []byte  `xml:",attr"`
***REMOVED***

type OmitAttrTest struct ***REMOVED***
	Int   int     `xml:",attr,omitempty"`
	Named int     `xml:"int,attr,omitempty"`
	Float float64 `xml:",attr,omitempty"`
	Uint8 uint8   `xml:",attr,omitempty"`
	Bool  bool    `xml:",attr,omitempty"`
	Str   string  `xml:",attr,omitempty"`
	Bytes []byte  `xml:",attr,omitempty"`
***REMOVED***

type OmitFieldTest struct ***REMOVED***
	Int   int           `xml:",omitempty"`
	Named int           `xml:"int,omitempty"`
	Float float64       `xml:",omitempty"`
	Uint8 uint8         `xml:",omitempty"`
	Bool  bool          `xml:",omitempty"`
	Str   string        `xml:",omitempty"`
	Bytes []byte        `xml:",omitempty"`
	Ptr   *PresenceTest `xml:",omitempty"`
***REMOVED***

type AnyTest struct ***REMOVED***
	XMLName  struct***REMOVED******REMOVED***  `xml:"a"`
	Nested   string    `xml:"nested>value"`
	AnyField AnyHolder `xml:",any"`
***REMOVED***

type AnyOmitTest struct ***REMOVED***
	XMLName  struct***REMOVED******REMOVED***   `xml:"a"`
	Nested   string     `xml:"nested>value"`
	AnyField *AnyHolder `xml:",any,omitempty"`
***REMOVED***

type AnySliceTest struct ***REMOVED***
	XMLName  struct***REMOVED******REMOVED***    `xml:"a"`
	Nested   string      `xml:"nested>value"`
	AnyField []AnyHolder `xml:",any"`
***REMOVED***

type AnyHolder struct ***REMOVED***
	XMLName Name
	XML     string `xml:",innerxml"`
***REMOVED***

type RecurseA struct ***REMOVED***
	A string
	B *RecurseB
***REMOVED***

type RecurseB struct ***REMOVED***
	A *RecurseA
	B string
***REMOVED***

type PresenceTest struct ***REMOVED***
	Exists *struct***REMOVED******REMOVED***
***REMOVED***

type IgnoreTest struct ***REMOVED***
	PublicSecret string `xml:"-"`
***REMOVED***

type MyBytes []byte

type Data struct ***REMOVED***
	Bytes  []byte
	Attr   []byte `xml:",attr"`
	Custom MyBytes
***REMOVED***

type Plain struct ***REMOVED***
	V interface***REMOVED******REMOVED***
***REMOVED***

type MyInt int

type EmbedInt struct ***REMOVED***
	MyInt
***REMOVED***

type Strings struct ***REMOVED***
	X []string `xml:"A>B,omitempty"`
***REMOVED***

type PointerFieldsTest struct ***REMOVED***
	XMLName  Name    `xml:"dummy"`
	Name     *string `xml:"name,attr"`
	Age      *uint   `xml:"age,attr"`
	Empty    *string `xml:"empty,attr"`
	Contents *string `xml:",chardata"`
***REMOVED***

type ChardataEmptyTest struct ***REMOVED***
	XMLName  Name    `xml:"test"`
	Contents *string `xml:",chardata"`
***REMOVED***

type MyMarshalerTest struct ***REMOVED***
***REMOVED***

var _ Marshaler = (*MyMarshalerTest)(nil)

func (m *MyMarshalerTest) MarshalXML(e *Encoder, start StartElement) error ***REMOVED***
	e.EncodeToken(start)
	e.EncodeToken(CharData([]byte("hello world")))
	e.EncodeToken(EndElement***REMOVED***start.Name***REMOVED***)
	return nil
***REMOVED***

type MyMarshalerAttrTest struct***REMOVED******REMOVED***

var _ MarshalerAttr = (*MyMarshalerAttrTest)(nil)

func (m *MyMarshalerAttrTest) MarshalXMLAttr(name Name) (Attr, error) ***REMOVED***
	return Attr***REMOVED***name, "hello world"***REMOVED***, nil
***REMOVED***

type MyMarshalerValueAttrTest struct***REMOVED******REMOVED***

var _ MarshalerAttr = MyMarshalerValueAttrTest***REMOVED******REMOVED***

func (m MyMarshalerValueAttrTest) MarshalXMLAttr(name Name) (Attr, error) ***REMOVED***
	return Attr***REMOVED***name, "hello world"***REMOVED***, nil
***REMOVED***

type MarshalerStruct struct ***REMOVED***
	Foo MyMarshalerAttrTest `xml:",attr"`
***REMOVED***

type MarshalerValueStruct struct ***REMOVED***
	Foo MyMarshalerValueAttrTest `xml:",attr"`
***REMOVED***

type InnerStruct struct ***REMOVED***
	XMLName Name `xml:"testns outer"`
***REMOVED***

type OuterStruct struct ***REMOVED***
	InnerStruct
	IntAttr int `xml:"int,attr"`
***REMOVED***

type OuterNamedStruct struct ***REMOVED***
	InnerStruct
	XMLName Name `xml:"outerns test"`
	IntAttr int  `xml:"int,attr"`
***REMOVED***

type OuterNamedOrderedStruct struct ***REMOVED***
	XMLName Name `xml:"outerns test"`
	InnerStruct
	IntAttr int `xml:"int,attr"`
***REMOVED***

type OuterOuterStruct struct ***REMOVED***
	OuterStruct
***REMOVED***

type NestedAndChardata struct ***REMOVED***
	AB       []string `xml:"A>B"`
	Chardata string   `xml:",chardata"`
***REMOVED***

type NestedAndComment struct ***REMOVED***
	AB      []string `xml:"A>B"`
	Comment string   `xml:",comment"`
***REMOVED***

type XMLNSFieldStruct struct ***REMOVED***
	Ns   string `xml:"xmlns,attr"`
	Body string
***REMOVED***

type NamedXMLNSFieldStruct struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"testns test"`
	Ns      string   `xml:"xmlns,attr"`
	Body    string
***REMOVED***

type XMLNSFieldStructWithOmitEmpty struct ***REMOVED***
	Ns   string `xml:"xmlns,attr,omitempty"`
	Body string
***REMOVED***

type NamedXMLNSFieldStructWithEmptyNamespace struct ***REMOVED***
	XMLName struct***REMOVED******REMOVED*** `xml:"test"`
	Ns      string   `xml:"xmlns,attr"`
	Body    string
***REMOVED***

type RecursiveXMLNSFieldStruct struct ***REMOVED***
	Ns   string                     `xml:"xmlns,attr"`
	Body *RecursiveXMLNSFieldStruct `xml:",omitempty"`
	Text string                     `xml:",omitempty"`
***REMOVED***

func ifaceptr(x interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	return &x
***REMOVED***

var (
	nameAttr     = "Sarah"
	ageAttr      = uint(12)
	contentsAttr = "lorem ipsum"
)

// Unless explicitly stated as such (or *Plain), all of the
// tests below are two-way tests. When introducing new tests,
// please try to make them two-way as well to ensure that
// marshalling and unmarshalling are as symmetrical as feasible.
var marshalTests = []struct ***REMOVED***
	Value         interface***REMOVED******REMOVED***
	ExpectXML     string
	MarshalOnly   bool
	UnmarshalOnly bool
***REMOVED******REMOVED***
	// Test nil marshals to nothing
	***REMOVED***Value: nil, ExpectXML: ``, MarshalOnly: true***REMOVED***,
	***REMOVED***Value: nilStruct, ExpectXML: ``, MarshalOnly: true***REMOVED***,

	// Test value types
	***REMOVED***Value: &Plain***REMOVED***true***REMOVED***, ExpectXML: `<Plain><V>true</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***false***REMOVED***, ExpectXML: `<Plain><V>false</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***int(42)***REMOVED***, ExpectXML: `<Plain><V>42</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***int8(42)***REMOVED***, ExpectXML: `<Plain><V>42</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***int16(42)***REMOVED***, ExpectXML: `<Plain><V>42</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***int32(42)***REMOVED***, ExpectXML: `<Plain><V>42</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***uint(42)***REMOVED***, ExpectXML: `<Plain><V>42</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***uint8(42)***REMOVED***, ExpectXML: `<Plain><V>42</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***uint16(42)***REMOVED***, ExpectXML: `<Plain><V>42</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***uint32(42)***REMOVED***, ExpectXML: `<Plain><V>42</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***float32(1.25)***REMOVED***, ExpectXML: `<Plain><V>1.25</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***float64(1.25)***REMOVED***, ExpectXML: `<Plain><V>1.25</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***uintptr(0xFFDD)***REMOVED***, ExpectXML: `<Plain><V>65501</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***"gopher"***REMOVED***, ExpectXML: `<Plain><V>gopher</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***[]byte("gopher")***REMOVED***, ExpectXML: `<Plain><V>gopher</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***"</>"***REMOVED***, ExpectXML: `<Plain><V>&lt;/&gt;</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***[]byte("</>")***REMOVED***, ExpectXML: `<Plain><V>&lt;/&gt;</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***[3]byte***REMOVED***'<', '/', '>'***REMOVED******REMOVED***, ExpectXML: `<Plain><V>&lt;/&gt;</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***NamedType("potato")***REMOVED***, ExpectXML: `<Plain><V>potato</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***[]int***REMOVED***1, 2, 3***REMOVED******REMOVED***, ExpectXML: `<Plain><V>1</V><V>2</V><V>3</V></Plain>`***REMOVED***,
	***REMOVED***Value: &Plain***REMOVED***[3]int***REMOVED***1, 2, 3***REMOVED******REMOVED***, ExpectXML: `<Plain><V>1</V><V>2</V><V>3</V></Plain>`***REMOVED***,
	***REMOVED***Value: ifaceptr(true), MarshalOnly: true, ExpectXML: `<bool>true</bool>`***REMOVED***,

	// Test time.
	***REMOVED***
		Value:     &Plain***REMOVED***time.Unix(1e9, 123456789).UTC()***REMOVED***,
		ExpectXML: `<Plain><V>2001-09-09T01:46:40.123456789Z</V></Plain>`,
	***REMOVED***,

	// A pointer to struct***REMOVED******REMOVED*** may be used to test for an element's presence.
	***REMOVED***
		Value:     &PresenceTest***REMOVED***new(struct***REMOVED******REMOVED***)***REMOVED***,
		ExpectXML: `<PresenceTest><Exists></Exists></PresenceTest>`,
	***REMOVED***,
	***REMOVED***
		Value:     &PresenceTest***REMOVED******REMOVED***,
		ExpectXML: `<PresenceTest></PresenceTest>`,
	***REMOVED***,

	// A pointer to struct***REMOVED******REMOVED*** may be used to test for an element's presence.
	***REMOVED***
		Value:     &PresenceTest***REMOVED***new(struct***REMOVED******REMOVED***)***REMOVED***,
		ExpectXML: `<PresenceTest><Exists></Exists></PresenceTest>`,
	***REMOVED***,
	***REMOVED***
		Value:     &PresenceTest***REMOVED******REMOVED***,
		ExpectXML: `<PresenceTest></PresenceTest>`,
	***REMOVED***,

	// A []byte field is only nil if the element was not found.
	***REMOVED***
		Value:         &Data***REMOVED******REMOVED***,
		ExpectXML:     `<Data></Data>`,
		UnmarshalOnly: true,
	***REMOVED***,
	***REMOVED***
		Value:         &Data***REMOVED***Bytes: []byte***REMOVED******REMOVED***, Custom: MyBytes***REMOVED******REMOVED***, Attr: []byte***REMOVED******REMOVED******REMOVED***,
		ExpectXML:     `<Data Attr=""><Bytes></Bytes><Custom></Custom></Data>`,
		UnmarshalOnly: true,
	***REMOVED***,

	// Check that []byte works, including named []byte types.
	***REMOVED***
		Value:     &Data***REMOVED***Bytes: []byte("ab"), Custom: MyBytes("cd"), Attr: []byte***REMOVED***'v'***REMOVED******REMOVED***,
		ExpectXML: `<Data Attr="v"><Bytes>ab</Bytes><Custom>cd</Custom></Data>`,
	***REMOVED***,

	// Test innerxml
	***REMOVED***
		Value: &SecretAgent***REMOVED***
			Handle:    "007",
			Identity:  "James Bond",
			Obfuscate: "<redacted/>",
		***REMOVED***,
		ExpectXML:   `<agent handle="007"><Identity>James Bond</Identity><redacted/></agent>`,
		MarshalOnly: true,
	***REMOVED***,
	***REMOVED***
		Value: &SecretAgent***REMOVED***
			Handle:    "007",
			Identity:  "James Bond",
			Obfuscate: "<Identity>James Bond</Identity><redacted/>",
		***REMOVED***,
		ExpectXML:     `<agent handle="007"><Identity>James Bond</Identity><redacted/></agent>`,
		UnmarshalOnly: true,
	***REMOVED***,

	// Test structs
	***REMOVED***Value: &Port***REMOVED***Type: "ssl", Number: "443"***REMOVED***, ExpectXML: `<port type="ssl">443</port>`***REMOVED***,
	***REMOVED***Value: &Port***REMOVED***Number: "443"***REMOVED***, ExpectXML: `<port>443</port>`***REMOVED***,
	***REMOVED***Value: &Port***REMOVED***Type: "<unix>"***REMOVED***, ExpectXML: `<port type="&lt;unix&gt;"></port>`***REMOVED***,
	***REMOVED***Value: &Port***REMOVED***Number: "443", Comment: "https"***REMOVED***, ExpectXML: `<port><!--https-->443</port>`***REMOVED***,
	***REMOVED***Value: &Port***REMOVED***Number: "443", Comment: "add space-"***REMOVED***, ExpectXML: `<port><!--add space- -->443</port>`, MarshalOnly: true***REMOVED***,
	***REMOVED***Value: &Domain***REMOVED***Name: []byte("google.com&friends")***REMOVED***, ExpectXML: `<domain>google.com&amp;friends</domain>`***REMOVED***,
	***REMOVED***Value: &Domain***REMOVED***Name: []byte("google.com"), Comment: []byte(" &friends ")***REMOVED***, ExpectXML: `<domain>google.com<!-- &friends --></domain>`***REMOVED***,
	***REMOVED***Value: &Book***REMOVED***Title: "Pride & Prejudice"***REMOVED***, ExpectXML: `<book>Pride &amp; Prejudice</book>`***REMOVED***,
	***REMOVED***Value: &Event***REMOVED***Year: -3114***REMOVED***, ExpectXML: `<event>-3114</event>`***REMOVED***,
	***REMOVED***Value: &Movie***REMOVED***Length: 13440***REMOVED***, ExpectXML: `<movie>13440</movie>`***REMOVED***,
	***REMOVED***Value: &Pi***REMOVED***Approximation: 3.14159265***REMOVED***, ExpectXML: `<pi>3.1415927</pi>`***REMOVED***,
	***REMOVED***Value: &Universe***REMOVED***Visible: 9.3e13***REMOVED***, ExpectXML: `<universe>9.3e+13</universe>`***REMOVED***,
	***REMOVED***Value: &Particle***REMOVED***HasMass: true***REMOVED***, ExpectXML: `<particle>true</particle>`***REMOVED***,
	***REMOVED***Value: &Departure***REMOVED***When: ParseTime("2013-01-09T00:15:00-09:00")***REMOVED***, ExpectXML: `<departure>2013-01-09T00:15:00-09:00</departure>`***REMOVED***,
	***REMOVED***Value: atomValue, ExpectXML: atomXml***REMOVED***,
	***REMOVED***
		Value: &Ship***REMOVED***
			Name:  "Heart of Gold",
			Pilot: "Computer",
			Age:   1,
			Drive: ImprobabilityDrive,
			Passenger: []*Passenger***REMOVED***
				***REMOVED***
					Name:   []string***REMOVED***"Zaphod", "Beeblebrox"***REMOVED***,
					Weight: 7.25,
				***REMOVED***,
				***REMOVED***
					Name:   []string***REMOVED***"Trisha", "McMillen"***REMOVED***,
					Weight: 5.5,
				***REMOVED***,
				***REMOVED***
					Name:   []string***REMOVED***"Ford", "Prefect"***REMOVED***,
					Weight: 7,
				***REMOVED***,
				***REMOVED***
					Name:   []string***REMOVED***"Arthur", "Dent"***REMOVED***,
					Weight: 6.75,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		ExpectXML: `<spaceship name="Heart of Gold" pilot="Computer">` +
			`<drive>` + strconv.Itoa(int(ImprobabilityDrive)) + `</drive>` +
			`<age>1</age>` +
			`<passenger>` +
			`<name>Zaphod</name>` +
			`<name>Beeblebrox</name>` +
			`<weight>7.25</weight>` +
			`</passenger>` +
			`<passenger>` +
			`<name>Trisha</name>` +
			`<name>McMillen</name>` +
			`<weight>5.5</weight>` +
			`</passenger>` +
			`<passenger>` +
			`<name>Ford</name>` +
			`<name>Prefect</name>` +
			`<weight>7</weight>` +
			`</passenger>` +
			`<passenger>` +
			`<name>Arthur</name>` +
			`<name>Dent</name>` +
			`<weight>6.75</weight>` +
			`</passenger>` +
			`</spaceship>`,
	***REMOVED***,

	// Test a>b
	***REMOVED***
		Value: &NestedItems***REMOVED***Items: nil, Item1: nil***REMOVED***,
		ExpectXML: `<result>` +
			`<Items>` +
			`</Items>` +
			`</result>`,
	***REMOVED***,
	***REMOVED***
		Value: &NestedItems***REMOVED***Items: []string***REMOVED******REMOVED***, Item1: []string***REMOVED******REMOVED******REMOVED***,
		ExpectXML: `<result>` +
			`<Items>` +
			`</Items>` +
			`</result>`,
		MarshalOnly: true,
	***REMOVED***,
	***REMOVED***
		Value: &NestedItems***REMOVED***Items: nil, Item1: []string***REMOVED***"A"***REMOVED******REMOVED***,
		ExpectXML: `<result>` +
			`<Items>` +
			`<item1>A</item1>` +
			`</Items>` +
			`</result>`,
	***REMOVED***,
	***REMOVED***
		Value: &NestedItems***REMOVED***Items: []string***REMOVED***"A", "B"***REMOVED***, Item1: nil***REMOVED***,
		ExpectXML: `<result>` +
			`<Items>` +
			`<item>A</item>` +
			`<item>B</item>` +
			`</Items>` +
			`</result>`,
	***REMOVED***,
	***REMOVED***
		Value: &NestedItems***REMOVED***Items: []string***REMOVED***"A", "B"***REMOVED***, Item1: []string***REMOVED***"C"***REMOVED******REMOVED***,
		ExpectXML: `<result>` +
			`<Items>` +
			`<item>A</item>` +
			`<item>B</item>` +
			`<item1>C</item1>` +
			`</Items>` +
			`</result>`,
	***REMOVED***,
	***REMOVED***
		Value: &NestedOrder***REMOVED***Field1: "C", Field2: "B", Field3: "A"***REMOVED***,
		ExpectXML: `<result>` +
			`<parent>` +
			`<c>C</c>` +
			`<b>B</b>` +
			`<a>A</a>` +
			`</parent>` +
			`</result>`,
	***REMOVED***,
	***REMOVED***
		Value: &NilTest***REMOVED***A: "A", B: nil, C: "C"***REMOVED***,
		ExpectXML: `<NilTest>` +
			`<parent1>` +
			`<parent2><a>A</a></parent2>` +
			`<parent2><c>C</c></parent2>` +
			`</parent1>` +
			`</NilTest>`,
		MarshalOnly: true, // Uses interface***REMOVED******REMOVED***
	***REMOVED***,
	***REMOVED***
		Value: &MixedNested***REMOVED***A: "A", B: "B", C: "C", D: "D"***REMOVED***,
		ExpectXML: `<result>` +
			`<parent1><a>A</a></parent1>` +
			`<b>B</b>` +
			`<parent1>` +
			`<parent2><c>C</c></parent2>` +
			`<d>D</d>` +
			`</parent1>` +
			`</result>`,
	***REMOVED***,
	***REMOVED***
		Value:     &Service***REMOVED***Port: &Port***REMOVED***Number: "80"***REMOVED******REMOVED***,
		ExpectXML: `<service><host><port>80</port></host></service>`,
	***REMOVED***,
	***REMOVED***
		Value:     &Service***REMOVED******REMOVED***,
		ExpectXML: `<service></service>`,
	***REMOVED***,
	***REMOVED***
		Value: &Service***REMOVED***Port: &Port***REMOVED***Number: "80"***REMOVED***, Extra1: "A", Extra2: "B"***REMOVED***,
		ExpectXML: `<service>` +
			`<host><port>80</port></host>` +
			`<Extra1>A</Extra1>` +
			`<host><extra2>B</extra2></host>` +
			`</service>`,
		MarshalOnly: true,
	***REMOVED***,
	***REMOVED***
		Value: &Service***REMOVED***Port: &Port***REMOVED***Number: "80"***REMOVED***, Extra2: "example"***REMOVED***,
		ExpectXML: `<service>` +
			`<host><port>80</port></host>` +
			`<host><extra2>example</extra2></host>` +
			`</service>`,
		MarshalOnly: true,
	***REMOVED***,
	***REMOVED***
		Value: &struct ***REMOVED***
			XMLName struct***REMOVED******REMOVED*** `xml:"space top"`
			A       string   `xml:"x>a"`
			B       string   `xml:"x>b"`
			C       string   `xml:"space x>c"`
			C1      string   `xml:"space1 x>c"`
			D1      string   `xml:"space1 x>d"`
			E1      string   `xml:"x>e"`
		***REMOVED******REMOVED***
			A:  "a",
			B:  "b",
			C:  "c",
			C1: "c1",
			D1: "d1",
			E1: "e1",
		***REMOVED***,
		ExpectXML: `<top xmlns="space">` +
			`<x><a>a</a><b>b</b><c>c</c></x>` +
			`<x xmlns="space1">` +
			`<c>c1</c>` +
			`<d>d1</d>` +
			`</x>` +
			`<x>` +
			`<e>e1</e>` +
			`</x>` +
			`</top>`,
	***REMOVED***,
	***REMOVED***
		Value: &struct ***REMOVED***
			XMLName Name
			A       string `xml:"x>a"`
			B       string `xml:"x>b"`
			C       string `xml:"space x>c"`
			C1      string `xml:"space1 x>c"`
			D1      string `xml:"space1 x>d"`
		***REMOVED******REMOVED***
			XMLName: Name***REMOVED***
				Space: "space0",
				Local: "top",
			***REMOVED***,
			A:  "a",
			B:  "b",
			C:  "c",
			C1: "c1",
			D1: "d1",
		***REMOVED***,
		ExpectXML: `<top xmlns="space0">` +
			`<x><a>a</a><b>b</b></x>` +
			`<x xmlns="space"><c>c</c></x>` +
			`<x xmlns="space1">` +
			`<c>c1</c>` +
			`<d>d1</d>` +
			`</x>` +
			`</top>`,
	***REMOVED***,
	***REMOVED***
		Value: &struct ***REMOVED***
			XMLName struct***REMOVED******REMOVED*** `xml:"top"`
			B       string   `xml:"space x>b"`
			B1      string   `xml:"space1 x>b"`
		***REMOVED******REMOVED***
			B:  "b",
			B1: "b1",
		***REMOVED***,
		ExpectXML: `<top>` +
			`<x xmlns="space"><b>b</b></x>` +
			`<x xmlns="space1"><b>b1</b></x>` +
			`</top>`,
	***REMOVED***,

	// Test struct embedding
	***REMOVED***
		Value: &EmbedA***REMOVED***
			EmbedC: EmbedC***REMOVED***
				FieldA1: "", // Shadowed by A.A
				FieldA2: "", // Shadowed by A.A
				FieldB:  "A.C.B",
				FieldC:  "A.C.C",
			***REMOVED***,
			EmbedB: EmbedB***REMOVED***
				FieldB: "A.B.B",
				EmbedC: &EmbedC***REMOVED***
					FieldA1: "A.B.C.A1",
					FieldA2: "A.B.C.A2",
					FieldB:  "", // Shadowed by A.B.B
					FieldC:  "A.B.C.C",
				***REMOVED***,
			***REMOVED***,
			FieldA: "A.A",
		***REMOVED***,
		ExpectXML: `<EmbedA>` +
			`<FieldB>A.C.B</FieldB>` +
			`<FieldC>A.C.C</FieldC>` +
			`<EmbedB>` +
			`<FieldB>A.B.B</FieldB>` +
			`<FieldA>` +
			`<A1>A.B.C.A1</A1>` +
			`<A2>A.B.C.A2</A2>` +
			`</FieldA>` +
			`<FieldC>A.B.C.C</FieldC>` +
			`</EmbedB>` +
			`<FieldA>A.A</FieldA>` +
			`</EmbedA>`,
	***REMOVED***,

	// Test that name casing matters
	***REMOVED***
		Value:     &NameCasing***REMOVED***Xy: "mixed", XY: "upper", XyA: "mixedA", XYA: "upperA"***REMOVED***,
		ExpectXML: `<casing Xy="mixedA" XY="upperA"><Xy>mixed</Xy><XY>upper</XY></casing>`,
	***REMOVED***,

	// Test the order in which the XML element name is chosen
	***REMOVED***
		Value: &NamePrecedence***REMOVED***
			FromTag:     XMLNameWithoutTag***REMOVED***Value: "A"***REMOVED***,
			FromNameVal: XMLNameWithoutTag***REMOVED***XMLName: Name***REMOVED***Local: "InXMLName"***REMOVED***, Value: "B"***REMOVED***,
			FromNameTag: XMLNameWithTag***REMOVED***Value: "C"***REMOVED***,
			InFieldName: "D",
		***REMOVED***,
		ExpectXML: `<Parent>` +
			`<InTag>A</InTag>` +
			`<InXMLName>B</InXMLName>` +
			`<InXMLNameTag>C</InXMLNameTag>` +
			`<InFieldName>D</InFieldName>` +
			`</Parent>`,
		MarshalOnly: true,
	***REMOVED***,
	***REMOVED***
		Value: &NamePrecedence***REMOVED***
			XMLName:     Name***REMOVED***Local: "Parent"***REMOVED***,
			FromTag:     XMLNameWithoutTag***REMOVED***XMLName: Name***REMOVED***Local: "InTag"***REMOVED***, Value: "A"***REMOVED***,
			FromNameVal: XMLNameWithoutTag***REMOVED***XMLName: Name***REMOVED***Local: "FromNameVal"***REMOVED***, Value: "B"***REMOVED***,
			FromNameTag: XMLNameWithTag***REMOVED***XMLName: Name***REMOVED***Local: "InXMLNameTag"***REMOVED***, Value: "C"***REMOVED***,
			InFieldName: "D",
		***REMOVED***,
		ExpectXML: `<Parent>` +
			`<InTag>A</InTag>` +
			`<FromNameVal>B</FromNameVal>` +
			`<InXMLNameTag>C</InXMLNameTag>` +
			`<InFieldName>D</InFieldName>` +
			`</Parent>`,
		UnmarshalOnly: true,
	***REMOVED***,

	// xml.Name works in a plain field as well.
	***REMOVED***
		Value:     &NameInField***REMOVED***Name***REMOVED***Space: "ns", Local: "foo"***REMOVED******REMOVED***,
		ExpectXML: `<NameInField><foo xmlns="ns"></foo></NameInField>`,
	***REMOVED***,
	***REMOVED***
		Value:         &NameInField***REMOVED***Name***REMOVED***Space: "ns", Local: "foo"***REMOVED******REMOVED***,
		ExpectXML:     `<NameInField><foo xmlns="ns"><ignore></ignore></foo></NameInField>`,
		UnmarshalOnly: true,
	***REMOVED***,

	// Marshaling zero xml.Name uses the tag or field name.
	***REMOVED***
		Value:       &NameInField***REMOVED******REMOVED***,
		ExpectXML:   `<NameInField><foo xmlns="ns"></foo></NameInField>`,
		MarshalOnly: true,
	***REMOVED***,

	// Test attributes
	***REMOVED***
		Value: &AttrTest***REMOVED***
			Int:   8,
			Named: 9,
			Float: 23.5,
			Uint8: 255,
			Bool:  true,
			Str:   "str",
			Bytes: []byte("byt"),
		***REMOVED***,
		ExpectXML: `<AttrTest Int="8" int="9" Float="23.5" Uint8="255"` +
			` Bool="true" Str="str" Bytes="byt"></AttrTest>`,
	***REMOVED***,
	***REMOVED***
		Value: &AttrTest***REMOVED***Bytes: []byte***REMOVED******REMOVED******REMOVED***,
		ExpectXML: `<AttrTest Int="0" int="0" Float="0" Uint8="0"` +
			` Bool="false" Str="" Bytes=""></AttrTest>`,
	***REMOVED***,
	***REMOVED***
		Value: &OmitAttrTest***REMOVED***
			Int:   8,
			Named: 9,
			Float: 23.5,
			Uint8: 255,
			Bool:  true,
			Str:   "str",
			Bytes: []byte("byt"),
		***REMOVED***,
		ExpectXML: `<OmitAttrTest Int="8" int="9" Float="23.5" Uint8="255"` +
			` Bool="true" Str="str" Bytes="byt"></OmitAttrTest>`,
	***REMOVED***,
	***REMOVED***
		Value:     &OmitAttrTest***REMOVED******REMOVED***,
		ExpectXML: `<OmitAttrTest></OmitAttrTest>`,
	***REMOVED***,

	// pointer fields
	***REMOVED***
		Value:       &PointerFieldsTest***REMOVED***Name: &nameAttr, Age: &ageAttr, Contents: &contentsAttr***REMOVED***,
		ExpectXML:   `<dummy name="Sarah" age="12">lorem ipsum</dummy>`,
		MarshalOnly: true,
	***REMOVED***,

	// empty chardata pointer field
	***REMOVED***
		Value:       &ChardataEmptyTest***REMOVED******REMOVED***,
		ExpectXML:   `<test></test>`,
		MarshalOnly: true,
	***REMOVED***,

	// omitempty on fields
	***REMOVED***
		Value: &OmitFieldTest***REMOVED***
			Int:   8,
			Named: 9,
			Float: 23.5,
			Uint8: 255,
			Bool:  true,
			Str:   "str",
			Bytes: []byte("byt"),
			Ptr:   &PresenceTest***REMOVED******REMOVED***,
		***REMOVED***,
		ExpectXML: `<OmitFieldTest>` +
			`<Int>8</Int>` +
			`<int>9</int>` +
			`<Float>23.5</Float>` +
			`<Uint8>255</Uint8>` +
			`<Bool>true</Bool>` +
			`<Str>str</Str>` +
			`<Bytes>byt</Bytes>` +
			`<Ptr></Ptr>` +
			`</OmitFieldTest>`,
	***REMOVED***,
	***REMOVED***
		Value:     &OmitFieldTest***REMOVED******REMOVED***,
		ExpectXML: `<OmitFieldTest></OmitFieldTest>`,
	***REMOVED***,

	// Test ",any"
	***REMOVED***
		ExpectXML: `<a><nested><value>known</value></nested><other><sub>unknown</sub></other></a>`,
		Value: &AnyTest***REMOVED***
			Nested: "known",
			AnyField: AnyHolder***REMOVED***
				XMLName: Name***REMOVED***Local: "other"***REMOVED***,
				XML:     "<sub>unknown</sub>",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		Value: &AnyTest***REMOVED***Nested: "known",
			AnyField: AnyHolder***REMOVED***
				XML:     "<unknown/>",
				XMLName: Name***REMOVED***Local: "AnyField"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		ExpectXML: `<a><nested><value>known</value></nested><AnyField><unknown/></AnyField></a>`,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<a><nested><value>b</value></nested></a>`,
		Value: &AnyOmitTest***REMOVED***
			Nested: "b",
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<a><nested><value>b</value></nested><c><d>e</d></c><g xmlns="f"><h>i</h></g></a>`,
		Value: &AnySliceTest***REMOVED***
			Nested: "b",
			AnyField: []AnyHolder***REMOVED***
				***REMOVED***
					XMLName: Name***REMOVED***Local: "c"***REMOVED***,
					XML:     "<d>e</d>",
				***REMOVED***,
				***REMOVED***
					XMLName: Name***REMOVED***Space: "f", Local: "g"***REMOVED***,
					XML:     "<h>i</h>",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<a><nested><value>b</value></nested></a>`,
		Value: &AnySliceTest***REMOVED***
			Nested: "b",
		***REMOVED***,
	***REMOVED***,

	// Test recursive types.
	***REMOVED***
		Value: &RecurseA***REMOVED***
			A: "a1",
			B: &RecurseB***REMOVED***
				A: &RecurseA***REMOVED***"a2", nil***REMOVED***,
				B: "b1",
			***REMOVED***,
		***REMOVED***,
		ExpectXML: `<RecurseA><A>a1</A><B><A><A>a2</A></A><B>b1</B></B></RecurseA>`,
	***REMOVED***,

	// Test ignoring fields via "-" tag
	***REMOVED***
		ExpectXML: `<IgnoreTest></IgnoreTest>`,
		Value:     &IgnoreTest***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML:   `<IgnoreTest></IgnoreTest>`,
		Value:       &IgnoreTest***REMOVED***PublicSecret: "can't tell"***REMOVED***,
		MarshalOnly: true,
	***REMOVED***,
	***REMOVED***
		ExpectXML:     `<IgnoreTest><PublicSecret>ignore me</PublicSecret></IgnoreTest>`,
		Value:         &IgnoreTest***REMOVED******REMOVED***,
		UnmarshalOnly: true,
	***REMOVED***,

	// Test escaping.
	***REMOVED***
		ExpectXML: `<a><nested><value>dquote: &#34;; squote: &#39;; ampersand: &amp;; less: &lt;; greater: &gt;;</value></nested><empty></empty></a>`,
		Value: &AnyTest***REMOVED***
			Nested:   `dquote: "; squote: '; ampersand: &; less: <; greater: >;`,
			AnyField: AnyHolder***REMOVED***XMLName: Name***REMOVED***Local: "empty"***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<a><nested><value>newline: &#xA;; cr: &#xD;; tab: &#x9;;</value></nested><AnyField></AnyField></a>`,
		Value: &AnyTest***REMOVED***
			Nested:   "newline: \n; cr: \r; tab: \t;",
			AnyField: AnyHolder***REMOVED***XMLName: Name***REMOVED***Local: "AnyField"***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: "<a><nested><value>1\r2\r\n3\n\r4\n5</value></nested></a>",
		Value: &AnyTest***REMOVED***
			Nested: "1\n2\n3\n\n4\n5",
		***REMOVED***,
		UnmarshalOnly: true,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<EmbedInt><MyInt>42</MyInt></EmbedInt>`,
		Value: &EmbedInt***REMOVED***
			MyInt: 42,
		***REMOVED***,
	***REMOVED***,
	// Test omitempty with parent chain; see golang.org/issue/4168.
	***REMOVED***
		ExpectXML: `<Strings><A></A></Strings>`,
		Value:     &Strings***REMOVED******REMOVED***,
	***REMOVED***,
	// Custom marshalers.
	***REMOVED***
		ExpectXML: `<MyMarshalerTest>hello world</MyMarshalerTest>`,
		Value:     &MyMarshalerTest***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<MarshalerStruct Foo="hello world"></MarshalerStruct>`,
		Value:     &MarshalerStruct***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<MarshalerValueStruct Foo="hello world"></MarshalerValueStruct>`,
		Value:     &MarshalerValueStruct***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<outer xmlns="testns" int="10"></outer>`,
		Value:     &OuterStruct***REMOVED***IntAttr: 10***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<test xmlns="outerns" int="10"></test>`,
		Value:     &OuterNamedStruct***REMOVED***XMLName: Name***REMOVED***Space: "outerns", Local: "test"***REMOVED***, IntAttr: 10***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<test xmlns="outerns" int="10"></test>`,
		Value:     &OuterNamedOrderedStruct***REMOVED***XMLName: Name***REMOVED***Space: "outerns", Local: "test"***REMOVED***, IntAttr: 10***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<outer xmlns="testns" int="10"></outer>`,
		Value:     &OuterOuterStruct***REMOVED***OuterStruct***REMOVED***IntAttr: 10***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<NestedAndChardata><A><B></B><B></B></A>test</NestedAndChardata>`,
		Value:     &NestedAndChardata***REMOVED***AB: make([]string, 2), Chardata: "test"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<NestedAndComment><A><B></B><B></B></A><!--test--></NestedAndComment>`,
		Value:     &NestedAndComment***REMOVED***AB: make([]string, 2), Comment: "test"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<XMLNSFieldStruct xmlns="http://example.com/ns"><Body>hello world</Body></XMLNSFieldStruct>`,
		Value:     &XMLNSFieldStruct***REMOVED***Ns: "http://example.com/ns", Body: "hello world"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<testns:test xmlns:testns="testns" xmlns="http://example.com/ns"><Body>hello world</Body></testns:test>`,
		Value:     &NamedXMLNSFieldStruct***REMOVED***Ns: "http://example.com/ns", Body: "hello world"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<testns:test xmlns:testns="testns"><Body>hello world</Body></testns:test>`,
		Value:     &NamedXMLNSFieldStruct***REMOVED***Ns: "", Body: "hello world"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<XMLNSFieldStructWithOmitEmpty><Body>hello world</Body></XMLNSFieldStructWithOmitEmpty>`,
		Value:     &XMLNSFieldStructWithOmitEmpty***REMOVED***Body: "hello world"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		// The xmlns attribute must be ignored because the <test>
		// element is in the empty namespace, so it's not possible
		// to set the default namespace to something non-empty.
		ExpectXML:   `<test><Body>hello world</Body></test>`,
		Value:       &NamedXMLNSFieldStructWithEmptyNamespace***REMOVED***Ns: "foo", Body: "hello world"***REMOVED***,
		MarshalOnly: true,
	***REMOVED***,
	***REMOVED***
		ExpectXML: `<RecursiveXMLNSFieldStruct xmlns="foo"><Body xmlns=""><Text>hello world</Text></Body></RecursiveXMLNSFieldStruct>`,
		Value: &RecursiveXMLNSFieldStruct***REMOVED***
			Ns: "foo",
			Body: &RecursiveXMLNSFieldStruct***REMOVED***
				Text: "hello world",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestMarshal(t *testing.T) ***REMOVED***
	for idx, test := range marshalTests ***REMOVED***
		if test.UnmarshalOnly ***REMOVED***
			continue
		***REMOVED***
		data, err := Marshal(test.Value)
		if err != nil ***REMOVED***
			t.Errorf("#%d: marshal(%#v): %s", idx, test.Value, err)
			continue
		***REMOVED***
		if got, want := string(data), test.ExpectXML; got != want ***REMOVED***
			if strings.Contains(want, "\n") ***REMOVED***
				t.Errorf("#%d: marshal(%#v):\nHAVE:\n%s\nWANT:\n%s", idx, test.Value, got, want)
			***REMOVED*** else ***REMOVED***
				t.Errorf("#%d: marshal(%#v):\nhave %#q\nwant %#q", idx, test.Value, got, want)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type AttrParent struct ***REMOVED***
	X string `xml:"X>Y,attr"`
***REMOVED***

type BadAttr struct ***REMOVED***
	Name []string `xml:"name,attr"`
***REMOVED***

var marshalErrorTests = []struct ***REMOVED***
	Value interface***REMOVED******REMOVED***
	Err   string
	Kind  reflect.Kind
***REMOVED******REMOVED***
	***REMOVED***
		Value: make(chan bool),
		Err:   "xml: unsupported type: chan bool",
		Kind:  reflect.Chan,
	***REMOVED***,
	***REMOVED***
		Value: map[string]string***REMOVED***
			"question": "What do you get when you multiply six by nine?",
			"answer":   "42",
		***REMOVED***,
		Err:  "xml: unsupported type: map[string]string",
		Kind: reflect.Map,
	***REMOVED***,
	***REMOVED***
		Value: map[*Ship]bool***REMOVED***nil: false***REMOVED***,
		Err:   "xml: unsupported type: map[*xml.Ship]bool",
		Kind:  reflect.Map,
	***REMOVED***,
	***REMOVED***
		Value: &Domain***REMOVED***Comment: []byte("f--bar")***REMOVED***,
		Err:   `xml: comments must not contain "--"`,
	***REMOVED***,
	// Reject parent chain with attr, never worked; see golang.org/issue/5033.
	***REMOVED***
		Value: &AttrParent***REMOVED******REMOVED***,
		Err:   `xml: X>Y chain not valid with attr flag`,
	***REMOVED***,
	***REMOVED***
		Value: BadAttr***REMOVED***[]string***REMOVED***"X", "Y"***REMOVED******REMOVED***,
		Err:   `xml: unsupported type: []string`,
	***REMOVED***,
***REMOVED***

var marshalIndentTests = []struct ***REMOVED***
	Value     interface***REMOVED******REMOVED***
	Prefix    string
	Indent    string
	ExpectXML string
***REMOVED******REMOVED***
	***REMOVED***
		Value: &SecretAgent***REMOVED***
			Handle:    "007",
			Identity:  "James Bond",
			Obfuscate: "<redacted/>",
		***REMOVED***,
		Prefix:    "",
		Indent:    "\t",
		ExpectXML: fmt.Sprintf("<agent handle=\"007\">\n\t<Identity>James Bond</Identity><redacted/>\n</agent>"),
	***REMOVED***,
***REMOVED***

func TestMarshalErrors(t *testing.T) ***REMOVED***
	for idx, test := range marshalErrorTests ***REMOVED***
		data, err := Marshal(test.Value)
		if err == nil ***REMOVED***
			t.Errorf("#%d: marshal(%#v) = [success] %q, want error %v", idx, test.Value, data, test.Err)
			continue
		***REMOVED***
		if err.Error() != test.Err ***REMOVED***
			t.Errorf("#%d: marshal(%#v) = [error] %v, want %v", idx, test.Value, err, test.Err)
		***REMOVED***
		if test.Kind != reflect.Invalid ***REMOVED***
			if kind := err.(*UnsupportedTypeError).Type.Kind(); kind != test.Kind ***REMOVED***
				t.Errorf("#%d: marshal(%#v) = [error kind] %s, want %s", idx, test.Value, kind, test.Kind)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Do invertibility testing on the various structures that we test
func TestUnmarshal(t *testing.T) ***REMOVED***
	for i, test := range marshalTests ***REMOVED***
		if test.MarshalOnly ***REMOVED***
			continue
		***REMOVED***
		if _, ok := test.Value.(*Plain); ok ***REMOVED***
			continue
		***REMOVED***
		vt := reflect.TypeOf(test.Value)
		dest := reflect.New(vt.Elem()).Interface()
		err := Unmarshal([]byte(test.ExpectXML), dest)

		switch fix := dest.(type) ***REMOVED***
		case *Feed:
			fix.Author.InnerXML = ""
			for i := range fix.Entry ***REMOVED***
				fix.Entry[i].Author.InnerXML = ""
			***REMOVED***
		***REMOVED***

		if err != nil ***REMOVED***
			t.Errorf("#%d: unexpected error: %#v", i, err)
		***REMOVED*** else if got, want := dest, test.Value; !reflect.DeepEqual(got, want) ***REMOVED***
			t.Errorf("#%d: unmarshal(%q):\nhave %#v\nwant %#v", i, test.ExpectXML, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMarshalIndent(t *testing.T) ***REMOVED***
	for i, test := range marshalIndentTests ***REMOVED***
		data, err := MarshalIndent(test.Value, test.Prefix, test.Indent)
		if err != nil ***REMOVED***
			t.Errorf("#%d: Error: %s", i, err)
			continue
		***REMOVED***
		if got, want := string(data), test.ExpectXML; got != want ***REMOVED***
			t.Errorf("#%d: MarshalIndent:\nGot:%s\nWant:\n%s", i, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

type limitedBytesWriter struct ***REMOVED***
	w      io.Writer
	remain int // until writes fail
***REMOVED***

func (lw *limitedBytesWriter) Write(p []byte) (n int, err error) ***REMOVED***
	if lw.remain <= 0 ***REMOVED***
		println("error")
		return 0, errors.New("write limit hit")
	***REMOVED***
	if len(p) > lw.remain ***REMOVED***
		p = p[:lw.remain]
		n, _ = lw.w.Write(p)
		lw.remain = 0
		return n, errors.New("write limit hit")
	***REMOVED***
	n, err = lw.w.Write(p)
	lw.remain -= n
	return n, err
***REMOVED***

func TestMarshalWriteErrors(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	const writeCap = 1024
	w := &limitedBytesWriter***REMOVED***&buf, writeCap***REMOVED***
	enc := NewEncoder(w)
	var err error
	var i int
	const n = 4000
	for i = 1; i <= n; i++ ***REMOVED***
		err = enc.Encode(&Passenger***REMOVED***
			Name:   []string***REMOVED***"Alice", "Bob"***REMOVED***,
			Weight: 5,
		***REMOVED***)
		if err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if err == nil ***REMOVED***
		t.Error("expected an error")
	***REMOVED***
	if i == n ***REMOVED***
		t.Errorf("expected to fail before the end")
	***REMOVED***
	if buf.Len() != writeCap ***REMOVED***
		t.Errorf("buf.Len() = %d; want %d", buf.Len(), writeCap)
	***REMOVED***
***REMOVED***

func TestMarshalWriteIOErrors(t *testing.T) ***REMOVED***
	enc := NewEncoder(errWriter***REMOVED******REMOVED***)

	expectErr := "unwritable"
	err := enc.Encode(&Passenger***REMOVED******REMOVED***)
	if err == nil || err.Error() != expectErr ***REMOVED***
		t.Errorf("EscapeTest = [error] %v, want %v", err, expectErr)
	***REMOVED***
***REMOVED***

func TestMarshalFlush(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.EncodeToken(CharData("hello world")); err != nil ***REMOVED***
		t.Fatalf("enc.EncodeToken: %v", err)
	***REMOVED***
	if buf.Len() > 0 ***REMOVED***
		t.Fatalf("enc.EncodeToken caused actual write: %q", buf.Bytes())
	***REMOVED***
	if err := enc.Flush(); err != nil ***REMOVED***
		t.Fatalf("enc.Flush: %v", err)
	***REMOVED***
	if buf.String() != "hello world" ***REMOVED***
		t.Fatalf("after enc.Flush, buf.String() = %q, want %q", buf.String(), "hello world")
	***REMOVED***
***REMOVED***

var encodeElementTests = []struct ***REMOVED***
	desc      string
	value     interface***REMOVED******REMOVED***
	start     StartElement
	expectXML string
***REMOVED******REMOVED******REMOVED***
	desc:  "simple string",
	value: "hello",
	start: StartElement***REMOVED***
		Name: Name***REMOVED***Local: "a"***REMOVED***,
	***REMOVED***,
	expectXML: `<a>hello</a>`,
***REMOVED***, ***REMOVED***
	desc:  "string with added attributes",
	value: "hello",
	start: StartElement***REMOVED***
		Name: Name***REMOVED***Local: "a"***REMOVED***,
		Attr: []Attr***REMOVED******REMOVED***
			Name:  Name***REMOVED***Local: "x"***REMOVED***,
			Value: "y",
		***REMOVED***, ***REMOVED***
			Name:  Name***REMOVED***Local: "foo"***REMOVED***,
			Value: "bar",
		***REMOVED******REMOVED***,
	***REMOVED***,
	expectXML: `<a x="y" foo="bar">hello</a>`,
***REMOVED***, ***REMOVED***
	desc: "start element with default name space",
	value: struct ***REMOVED***
		Foo XMLNameWithNSTag
	***REMOVED******REMOVED***
		Foo: XMLNameWithNSTag***REMOVED***
			Value: "hello",
		***REMOVED***,
	***REMOVED***,
	start: StartElement***REMOVED***
		Name: Name***REMOVED***Space: "ns", Local: "a"***REMOVED***,
		Attr: []Attr***REMOVED******REMOVED***
			Name: Name***REMOVED***Local: "xmlns"***REMOVED***,
			// "ns" is the name space defined in XMLNameWithNSTag
			Value: "ns",
		***REMOVED******REMOVED***,
	***REMOVED***,
	expectXML: `<a xmlns="ns"><InXMLNameWithNSTag>hello</InXMLNameWithNSTag></a>`,
***REMOVED***, ***REMOVED***
	desc: "start element in name space with different default name space",
	value: struct ***REMOVED***
		Foo XMLNameWithNSTag
	***REMOVED******REMOVED***
		Foo: XMLNameWithNSTag***REMOVED***
			Value: "hello",
		***REMOVED***,
	***REMOVED***,
	start: StartElement***REMOVED***
		Name: Name***REMOVED***Space: "ns2", Local: "a"***REMOVED***,
		Attr: []Attr***REMOVED******REMOVED***
			Name: Name***REMOVED***Local: "xmlns"***REMOVED***,
			// "ns" is the name space defined in XMLNameWithNSTag
			Value: "ns",
		***REMOVED******REMOVED***,
	***REMOVED***,
	expectXML: `<ns2:a xmlns:ns2="ns2" xmlns="ns"><InXMLNameWithNSTag>hello</InXMLNameWithNSTag></ns2:a>`,
***REMOVED***, ***REMOVED***
	desc:  "XMLMarshaler with start element with default name space",
	value: &MyMarshalerTest***REMOVED******REMOVED***,
	start: StartElement***REMOVED***
		Name: Name***REMOVED***Space: "ns2", Local: "a"***REMOVED***,
		Attr: []Attr***REMOVED******REMOVED***
			Name: Name***REMOVED***Local: "xmlns"***REMOVED***,
			// "ns" is the name space defined in XMLNameWithNSTag
			Value: "ns",
		***REMOVED******REMOVED***,
	***REMOVED***,
	expectXML: `<ns2:a xmlns:ns2="ns2" xmlns="ns">hello world</ns2:a>`,
***REMOVED******REMOVED***

func TestEncodeElement(t *testing.T) ***REMOVED***
	for idx, test := range encodeElementTests ***REMOVED***
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		err := enc.EncodeElement(test.value, test.start)
		if err != nil ***REMOVED***
			t.Fatalf("enc.EncodeElement: %v", err)
		***REMOVED***
		err = enc.Flush()
		if err != nil ***REMOVED***
			t.Fatalf("enc.Flush: %v", err)
		***REMOVED***
		if got, want := buf.String(), test.expectXML; got != want ***REMOVED***
			t.Errorf("#%d(%s): EncodeElement(%#v, %#v):\nhave %#q\nwant %#q", idx, test.desc, test.value, test.start, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkMarshal(b *testing.B) ***REMOVED***
	b.ReportAllocs()
	for i := 0; i < b.N; i++ ***REMOVED***
		Marshal(atomValue)
	***REMOVED***
***REMOVED***

func BenchmarkUnmarshal(b *testing.B) ***REMOVED***
	b.ReportAllocs()
	xml := []byte(atomXml)
	for i := 0; i < b.N; i++ ***REMOVED***
		Unmarshal(xml, &Feed***REMOVED******REMOVED***)
	***REMOVED***
***REMOVED***

// golang.org/issue/6556
func TestStructPointerMarshal(t *testing.T) ***REMOVED***
	type A struct ***REMOVED***
		XMLName string `xml:"a"`
		B       []interface***REMOVED******REMOVED***
	***REMOVED***
	type C struct ***REMOVED***
		XMLName Name
		Value   string `xml:"value"`
	***REMOVED***

	a := new(A)
	a.B = append(a.B, &C***REMOVED***
		XMLName: Name***REMOVED***Local: "c"***REMOVED***,
		Value:   "x",
	***REMOVED***)

	b, err := Marshal(a)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if x := string(b); x != "<a><c><value>x</value></c></a>" ***REMOVED***
		t.Fatal(x)
	***REMOVED***
	var v A
	err = Unmarshal(b, &v)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

var encodeTokenTests = []struct ***REMOVED***
	desc string
	toks []Token
	want string
	err  string
***REMOVED******REMOVED******REMOVED***
	desc: "start element with name space",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "local"***REMOVED***, nil***REMOVED***,
	***REMOVED***,
	want: `<space:local xmlns:space="space">`,
***REMOVED***, ***REMOVED***
	desc: "start element with no name",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", ""***REMOVED***, nil***REMOVED***,
	***REMOVED***,
	err: "xml: start tag with no name",
***REMOVED***, ***REMOVED***
	desc: "end element with no name",
	toks: []Token***REMOVED***
		EndElement***REMOVED***Name***REMOVED***"space", ""***REMOVED******REMOVED***,
	***REMOVED***,
	err: "xml: end tag with no name",
***REMOVED***, ***REMOVED***
	desc: "char data",
	toks: []Token***REMOVED***
		CharData("foo"),
	***REMOVED***,
	want: `foo`,
***REMOVED***, ***REMOVED***
	desc: "char data with escaped chars",
	toks: []Token***REMOVED***
		CharData(" \t\n"),
	***REMOVED***,
	want: " &#x9;\n",
***REMOVED***, ***REMOVED***
	desc: "comment",
	toks: []Token***REMOVED***
		Comment("foo"),
	***REMOVED***,
	want: `<!--foo-->`,
***REMOVED***, ***REMOVED***
	desc: "comment with invalid content",
	toks: []Token***REMOVED***
		Comment("foo-->"),
	***REMOVED***,
	err: "xml: EncodeToken of Comment containing --> marker",
***REMOVED***, ***REMOVED***
	desc: "proc instruction",
	toks: []Token***REMOVED***
		ProcInst***REMOVED***"Target", []byte("Instruction")***REMOVED***,
	***REMOVED***,
	want: `<?Target Instruction?>`,
***REMOVED***, ***REMOVED***
	desc: "proc instruction with empty target",
	toks: []Token***REMOVED***
		ProcInst***REMOVED***"", []byte("Instruction")***REMOVED***,
	***REMOVED***,
	err: "xml: EncodeToken of ProcInst with invalid Target",
***REMOVED***, ***REMOVED***
	desc: "proc instruction with bad content",
	toks: []Token***REMOVED***
		ProcInst***REMOVED***"", []byte("Instruction?>")***REMOVED***,
	***REMOVED***,
	err: "xml: EncodeToken of ProcInst with invalid Target",
***REMOVED***, ***REMOVED***
	desc: "directive",
	toks: []Token***REMOVED***
		Directive("foo"),
	***REMOVED***,
	want: `<!foo>`,
***REMOVED***, ***REMOVED***
	desc: "more complex directive",
	toks: []Token***REMOVED***
		Directive("DOCTYPE doc [ <!ELEMENT doc '>'> <!-- com>ment --> ]"),
	***REMOVED***,
	want: `<!DOCTYPE doc [ <!ELEMENT doc '>'> <!-- com>ment --> ]>`,
***REMOVED***, ***REMOVED***
	desc: "directive instruction with bad name",
	toks: []Token***REMOVED***
		Directive("foo>"),
	***REMOVED***,
	err: "xml: EncodeToken of Directive containing wrong < or > markers",
***REMOVED***, ***REMOVED***
	desc: "end tag without start tag",
	toks: []Token***REMOVED***
		EndElement***REMOVED***Name***REMOVED***"foo", "bar"***REMOVED******REMOVED***,
	***REMOVED***,
	err: "xml: end tag </bar> without start tag",
***REMOVED***, ***REMOVED***
	desc: "mismatching end tag local name",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, nil***REMOVED***,
		EndElement***REMOVED***Name***REMOVED***"", "bar"***REMOVED******REMOVED***,
	***REMOVED***,
	err:  "xml: end tag </bar> does not match start tag <foo>",
	want: `<foo>`,
***REMOVED***, ***REMOVED***
	desc: "mismatching end tag namespace",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, nil***REMOVED***,
		EndElement***REMOVED***Name***REMOVED***"another", "foo"***REMOVED******REMOVED***,
	***REMOVED***,
	err:  "xml: end tag </foo> in namespace another does not match start tag <foo> in namespace space",
	want: `<space:foo xmlns:space="space">`,
***REMOVED***, ***REMOVED***
	desc: "start element with explicit namespace",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "local"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "x"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<x:local xmlns:x="space" x:foo="value">`,
***REMOVED***, ***REMOVED***
	desc: "start element with explicit namespace and colliding prefix",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "local"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "x"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, "value"***REMOVED***,
			***REMOVED***Name***REMOVED***"x", "bar"***REMOVED***, "other"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<x:local xmlns:x_1="x" xmlns:x="space" x:foo="value" x_1:bar="other">`,
***REMOVED***, ***REMOVED***
	desc: "start element using previously defined namespace",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"", "local"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "x"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"space", "x"***REMOVED***, "y"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<local xmlns:x="space"><x:foo x:x="y">`,
***REMOVED***, ***REMOVED***
	desc: "nested name space with same prefix",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "x"***REMOVED***, "space1"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "x"***REMOVED***, "space2"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"space1", "a"***REMOVED***, "space1 value"***REMOVED***,
			***REMOVED***Name***REMOVED***"space2", "b"***REMOVED***, "space2 value"***REMOVED***,
		***REMOVED******REMOVED***,
		EndElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED******REMOVED***,
		EndElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"space1", "a"***REMOVED***, "space1 value"***REMOVED***,
			***REMOVED***Name***REMOVED***"space2", "b"***REMOVED***, "space2 value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns:x="space1"><foo xmlns:x="space2"><foo xmlns:space1="space1" space1:a="space1 value" x:b="space2 value"></foo></foo><foo xmlns:space2="space2" x:a="space1 value" space2:b="space2 value">`,
***REMOVED***, ***REMOVED***
	desc: "start element defining several prefixes for the same name space",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "a"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"xmlns", "b"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"space", "x"***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<a:foo xmlns:a="space" a:x="value">`,
***REMOVED***, ***REMOVED***
	desc: "nested element redefines name space",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "x"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "y"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"space", "a"***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns:x="space"><x:foo x:a="value">`,
***REMOVED***, ***REMOVED***
	desc: "nested element creates alias for default name space",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "y"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"space", "a"***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns="space"><foo xmlns:y="space" y:a="value">`,
***REMOVED***, ***REMOVED***
	desc: "nested element defines default name space with existing prefix",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "x"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"space", "a"***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns:x="space"><foo xmlns="space" x:a="value">`,
***REMOVED***, ***REMOVED***
	desc: "nested element uses empty attribute name space when default ns defined",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "attr"***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns="space"><foo attr="value">`,
***REMOVED***, ***REMOVED***
	desc: "redefine xmlns",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"foo", "xmlns"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	err: `xml: cannot redefine xmlns attribute prefix`,
***REMOVED***, ***REMOVED***
	desc: "xmlns with explicit name space #1",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xml", "xmlns"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns="space">`,
***REMOVED***, ***REMOVED***
	desc: "xmlns with explicit name space #2",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***xmlURL, "xmlns"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns="space">`,
***REMOVED***, ***REMOVED***
	desc: "empty name space declaration is ignored",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"xmlns", "foo"***REMOVED***, ""***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo>`,
***REMOVED***, ***REMOVED***
	desc: "attribute with no name is ignored",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", ""***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo>`,
***REMOVED***, ***REMOVED***
	desc: "namespace URL with non-valid name",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"/34", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"/34", "x"***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<_:foo xmlns:_="/34" _:x="value">`,
***REMOVED***, ***REMOVED***
	desc: "nested element resets default namespace to empty",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, ""***REMOVED***,
			***REMOVED***Name***REMOVED***"", "x"***REMOVED***, "value"***REMOVED***,
			***REMOVED***Name***REMOVED***"space", "x"***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns="space"><foo xmlns:space="space" xmlns="" x="value" space:x="value">`,
***REMOVED***, ***REMOVED***
	desc: "nested element requires empty default name space",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, nil***REMOVED***,
	***REMOVED***,
	want: `<foo xmlns="space"><foo xmlns="">`,
***REMOVED***, ***REMOVED***
	desc: "attribute uses name space from xmlns",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"some/space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "attr"***REMOVED***, "value"***REMOVED***,
			***REMOVED***Name***REMOVED***"some/space", "other"***REMOVED***, "other value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<space:foo xmlns:space="some/space" attr="value" space:other="other value">`,
***REMOVED***, ***REMOVED***
	desc: "default name space should not be used by attributes",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"xmlns", "bar"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"space", "baz"***REMOVED***, "foo"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"space", "baz"***REMOVED***, nil***REMOVED***,
		EndElement***REMOVED***Name***REMOVED***"space", "baz"***REMOVED******REMOVED***,
		EndElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns:bar="space" xmlns="space" bar:baz="foo"><baz></baz></foo>`,
***REMOVED***, ***REMOVED***
	desc: "default name space not used by attributes, not explicitly defined",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "space"***REMOVED***,
			***REMOVED***Name***REMOVED***"space", "baz"***REMOVED***, "foo"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"space", "baz"***REMOVED***, nil***REMOVED***,
		EndElement***REMOVED***Name***REMOVED***"space", "baz"***REMOVED******REMOVED***,
		EndElement***REMOVED***Name***REMOVED***"space", "foo"***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo xmlns:space="space" xmlns="space" space:baz="foo"><baz></baz></foo>`,
***REMOVED***, ***REMOVED***
	desc: "impossible xmlns declaration",
	toks: []Token***REMOVED***
		StartElement***REMOVED***Name***REMOVED***"", "foo"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "space"***REMOVED***,
		***REMOVED******REMOVED***,
		StartElement***REMOVED***Name***REMOVED***"space", "bar"***REMOVED***, []Attr***REMOVED***
			***REMOVED***Name***REMOVED***"space", "attr"***REMOVED***, "value"***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***,
	want: `<foo><space:bar xmlns:space="space" space:attr="value">`,
***REMOVED******REMOVED***

func TestEncodeToken(t *testing.T) ***REMOVED***
loop:
	for i, tt := range encodeTokenTests ***REMOVED***
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		var err error
		for j, tok := range tt.toks ***REMOVED***
			err = enc.EncodeToken(tok)
			if err != nil && j < len(tt.toks)-1 ***REMOVED***
				t.Errorf("#%d %s token #%d: %v", i, tt.desc, j, err)
				continue loop
			***REMOVED***
		***REMOVED***
		errorf := func(f string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
			t.Errorf("#%d %s token #%d:%s", i, tt.desc, len(tt.toks)-1, fmt.Sprintf(f, a...))
		***REMOVED***
		switch ***REMOVED***
		case tt.err != "" && err == nil:
			errorf(" expected error; got none")
			continue
		case tt.err == "" && err != nil:
			errorf(" got error: %v", err)
			continue
		case tt.err != "" && err != nil && tt.err != err.Error():
			errorf(" error mismatch; got %v, want %v", err, tt.err)
			continue
		***REMOVED***
		if err := enc.Flush(); err != nil ***REMOVED***
			errorf(" %v", err)
			continue
		***REMOVED***
		if got := buf.String(); got != tt.want ***REMOVED***
			errorf("\ngot  %v\nwant %v", got, tt.want)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestProcInstEncodeToken(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.EncodeToken(ProcInst***REMOVED***"xml", []byte("Instruction")***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("enc.EncodeToken: expected to be able to encode xml target ProcInst as first token, %s", err)
	***REMOVED***

	if err := enc.EncodeToken(ProcInst***REMOVED***"Target", []byte("Instruction")***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("enc.EncodeToken: expected to be able to add non-xml target ProcInst")
	***REMOVED***

	if err := enc.EncodeToken(ProcInst***REMOVED***"xml", []byte("Instruction")***REMOVED***); err == nil ***REMOVED***
		t.Fatalf("enc.EncodeToken: expected to not be allowed to encode xml target ProcInst when not first token")
	***REMOVED***
***REMOVED***

func TestDecodeEncode(t *testing.T) ***REMOVED***
	var in, out bytes.Buffer
	in.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<?Target Instruction?>
<root>
</root>	
`)
	dec := NewDecoder(&in)
	enc := NewEncoder(&out)
	for tok, err := dec.Token(); err == nil; tok, err = dec.Token() ***REMOVED***
		err = enc.EncodeToken(tok)
		if err != nil ***REMOVED***
			t.Fatalf("enc.EncodeToken: Unable to encode token (%#v), %v", tok, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Issue 9796. Used to fail with GORACE="halt_on_error=1" -race.
func TestRace9796(t *testing.T) ***REMOVED***
	type A struct***REMOVED******REMOVED***
	type B struct ***REMOVED***
		C []A `xml:"X>Y"`
	***REMOVED***
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			Marshal(B***REMOVED***[]A***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***)
			wg.Done()
		***REMOVED***()
	***REMOVED***
	wg.Wait()
***REMOVED***

func TestIsValidDirective(t *testing.T) ***REMOVED***
	testOK := []string***REMOVED***
		"<>",
		"< < > >",
		"<!DOCTYPE '<' '>' '>' <!--nothing-->>",
		"<!DOCTYPE doc [ <!ELEMENT doc ANY> <!ELEMENT doc ANY> ]>",
		"<!DOCTYPE doc [ <!ELEMENT doc \"ANY> '<' <!E\" LEMENT '>' doc ANY> ]>",
		"<!DOCTYPE doc <!-- just>>>> a < comment --> [ <!ITEM anything> ] >",
	***REMOVED***
	testKO := []string***REMOVED***
		"<",
		">",
		"<!--",
		"-->",
		"< > > < < >",
		"<!dummy <!-- > -->",
		"<!DOCTYPE doc '>",
		"<!DOCTYPE doc '>'",
		"<!DOCTYPE doc <!--comment>",
	***REMOVED***
	for _, s := range testOK ***REMOVED***
		if !isValidDirective(Directive(s)) ***REMOVED***
			t.Errorf("Directive %q is expected to be valid", s)
		***REMOVED***
	***REMOVED***
	for _, s := range testKO ***REMOVED***
		if isValidDirective(Directive(s)) ***REMOVED***
			t.Errorf("Directive %q is expected to be invalid", s)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Issue 11719. EncodeToken used to silently eat tokens with an invalid type.
func TestSimpleUseOfEncodeToken(t *testing.T) ***REMOVED***
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.EncodeToken(&StartElement***REMOVED***Name: Name***REMOVED***"", "object1"***REMOVED******REMOVED***); err == nil ***REMOVED***
		t.Errorf("enc.EncodeToken: pointer type should be rejected")
	***REMOVED***
	if err := enc.EncodeToken(&EndElement***REMOVED***Name: Name***REMOVED***"", "object1"***REMOVED******REMOVED***); err == nil ***REMOVED***
		t.Errorf("enc.EncodeToken: pointer type should be rejected")
	***REMOVED***
	if err := enc.EncodeToken(StartElement***REMOVED***Name: Name***REMOVED***"", "object2"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Errorf("enc.EncodeToken: StartElement %s", err)
	***REMOVED***
	if err := enc.EncodeToken(EndElement***REMOVED***Name: Name***REMOVED***"", "object2"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Errorf("enc.EncodeToken: EndElement %s", err)
	***REMOVED***
	if err := enc.EncodeToken(Universe***REMOVED******REMOVED***); err == nil ***REMOVED***
		t.Errorf("enc.EncodeToken: invalid type not caught")
	***REMOVED***
	if err := enc.Flush(); err != nil ***REMOVED***
		t.Errorf("enc.Flush: %s", err)
	***REMOVED***
	if buf.Len() == 0 ***REMOVED***
		t.Errorf("enc.EncodeToken: empty buffer")
	***REMOVED***
	want := "<object2></object2>"
	if buf.String() != want ***REMOVED***
		t.Errorf("enc.EncodeToken: expected %q; got %q", want, buf.String())
	***REMOVED***
***REMOVED***
