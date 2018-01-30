// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"bufio"
	"bytes"
	"encoding"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

const (
	// A generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

// Marshal returns the XML encoding of v.
//
// Marshal handles an array or slice by marshalling each of the elements.
// Marshal handles a pointer by marshalling the value it points at or, if the
// pointer is nil, by writing nothing. Marshal handles an interface value by
// marshalling the value it contains or, if the interface value is nil, by
// writing nothing. Marshal handles all other data by writing one or more XML
// elements containing the data.
//
// The name for the XML elements is taken from, in order of preference:
//     - the tag on the XMLName field, if the data is a struct
//     - the value of the XMLName field of type xml.Name
//     - the tag of the struct field used to obtain the data
//     - the name of the struct field used to obtain the data
//     - the name of the marshalled type
//
// The XML element for a struct contains marshalled elements for each of the
// exported fields of the struct, with these exceptions:
//     - the XMLName field, described above, is omitted.
//     - a field with tag "-" is omitted.
//     - a field with tag "name,attr" becomes an attribute with
//       the given name in the XML element.
//     - a field with tag ",attr" becomes an attribute with the
//       field name in the XML element.
//     - a field with tag ",chardata" is written as character data,
//       not as an XML element.
//     - a field with tag ",innerxml" is written verbatim, not subject
//       to the usual marshalling procedure.
//     - a field with tag ",comment" is written as an XML comment, not
//       subject to the usual marshalling procedure. It must not contain
//       the "--" string within it.
//     - a field with a tag including the "omitempty" option is omitted
//       if the field value is empty. The empty values are false, 0, any
//       nil pointer or interface value, and any array, slice, map, or
//       string of length zero.
//     - an anonymous struct field is handled as if the fields of its
//       value were part of the outer struct.
//
// If a field uses a tag "a>b>c", then the element c will be nested inside
// parent elements a and b. Fields that appear next to each other that name
// the same parent will be enclosed in one XML element.
//
// See MarshalIndent for an example.
//
// Marshal will return an error if asked to marshal a channel, function, or map.
func Marshal(v interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	var b bytes.Buffer
	if err := NewEncoder(&b).Encode(v); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b.Bytes(), nil
***REMOVED***

// Marshaler is the interface implemented by objects that can marshal
// themselves into valid XML elements.
//
// MarshalXML encodes the receiver as zero or more XML elements.
// By convention, arrays or slices are typically encoded as a sequence
// of elements, one per entry.
// Using start as the element tag is not required, but doing so
// will enable Unmarshal to match the XML elements to the correct
// struct field.
// One common implementation strategy is to construct a separate
// value with a layout corresponding to the desired XML and then
// to encode it using e.EncodeElement.
// Another common strategy is to use repeated calls to e.EncodeToken
// to generate the XML output one token at a time.
// The sequence of encoded tokens must make up zero or more valid
// XML elements.
type Marshaler interface ***REMOVED***
	MarshalXML(e *Encoder, start StartElement) error
***REMOVED***

// MarshalerAttr is the interface implemented by objects that can marshal
// themselves into valid XML attributes.
//
// MarshalXMLAttr returns an XML attribute with the encoded value of the receiver.
// Using name as the attribute name is not required, but doing so
// will enable Unmarshal to match the attribute to the correct
// struct field.
// If MarshalXMLAttr returns the zero attribute Attr***REMOVED******REMOVED***, no attribute
// will be generated in the output.
// MarshalXMLAttr is used only for struct fields with the
// "attr" option in the field tag.
type MarshalerAttr interface ***REMOVED***
	MarshalXMLAttr(name Name) (Attr, error)
***REMOVED***

// MarshalIndent works like Marshal, but each XML element begins on a new
// indented line that starts with prefix and is followed by one or more
// copies of indent according to the nesting depth.
func MarshalIndent(v interface***REMOVED******REMOVED***, prefix, indent string) ([]byte, error) ***REMOVED***
	var b bytes.Buffer
	enc := NewEncoder(&b)
	enc.Indent(prefix, indent)
	if err := enc.Encode(v); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b.Bytes(), nil
***REMOVED***

// An Encoder writes XML data to an output stream.
type Encoder struct ***REMOVED***
	p printer
***REMOVED***

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder ***REMOVED***
	e := &Encoder***REMOVED***printer***REMOVED***Writer: bufio.NewWriter(w)***REMOVED******REMOVED***
	e.p.encoder = e
	return e
***REMOVED***

// Indent sets the encoder to generate XML in which each element
// begins on a new indented line that starts with prefix and is followed by
// one or more copies of indent according to the nesting depth.
func (enc *Encoder) Indent(prefix, indent string) ***REMOVED***
	enc.p.prefix = prefix
	enc.p.indent = indent
***REMOVED***

// Encode writes the XML encoding of v to the stream.
//
// See the documentation for Marshal for details about the conversion
// of Go values to XML.
//
// Encode calls Flush before returning.
func (enc *Encoder) Encode(v interface***REMOVED******REMOVED***) error ***REMOVED***
	err := enc.p.marshalValue(reflect.ValueOf(v), nil, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return enc.p.Flush()
***REMOVED***

// EncodeElement writes the XML encoding of v to the stream,
// using start as the outermost tag in the encoding.
//
// See the documentation for Marshal for details about the conversion
// of Go values to XML.
//
// EncodeElement calls Flush before returning.
func (enc *Encoder) EncodeElement(v interface***REMOVED******REMOVED***, start StartElement) error ***REMOVED***
	err := enc.p.marshalValue(reflect.ValueOf(v), nil, &start)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return enc.p.Flush()
***REMOVED***

var (
	begComment   = []byte("<!--")
	endComment   = []byte("-->")
	endProcInst  = []byte("?>")
	endDirective = []byte(">")
)

// EncodeToken writes the given XML token to the stream.
// It returns an error if StartElement and EndElement tokens are not
// properly matched.
//
// EncodeToken does not call Flush, because usually it is part of a
// larger operation such as Encode or EncodeElement (or a custom
// Marshaler's MarshalXML invoked during those), and those will call
// Flush when finished. Callers that create an Encoder and then invoke
// EncodeToken directly, without using Encode or EncodeElement, need to
// call Flush when finished to ensure that the XML is written to the
// underlying writer.
//
// EncodeToken allows writing a ProcInst with Target set to "xml" only
// as the first token in the stream.
//
// When encoding a StartElement holding an XML namespace prefix
// declaration for a prefix that is not already declared, contained
// elements (including the StartElement itself) will use the declared
// prefix when encoding names with matching namespace URIs.
func (enc *Encoder) EncodeToken(t Token) error ***REMOVED***

	p := &enc.p
	switch t := t.(type) ***REMOVED***
	case StartElement:
		if err := p.writeStart(&t); err != nil ***REMOVED***
			return err
		***REMOVED***
	case EndElement:
		if err := p.writeEnd(t.Name); err != nil ***REMOVED***
			return err
		***REMOVED***
	case CharData:
		escapeText(p, t, false)
	case Comment:
		if bytes.Contains(t, endComment) ***REMOVED***
			return fmt.Errorf("xml: EncodeToken of Comment containing --> marker")
		***REMOVED***
		p.WriteString("<!--")
		p.Write(t)
		p.WriteString("-->")
		return p.cachedWriteError()
	case ProcInst:
		// First token to be encoded which is also a ProcInst with target of xml
		// is the xml declaration. The only ProcInst where target of xml is allowed.
		if t.Target == "xml" && p.Buffered() != 0 ***REMOVED***
			return fmt.Errorf("xml: EncodeToken of ProcInst xml target only valid for xml declaration, first token encoded")
		***REMOVED***
		if !isNameString(t.Target) ***REMOVED***
			return fmt.Errorf("xml: EncodeToken of ProcInst with invalid Target")
		***REMOVED***
		if bytes.Contains(t.Inst, endProcInst) ***REMOVED***
			return fmt.Errorf("xml: EncodeToken of ProcInst containing ?> marker")
		***REMOVED***
		p.WriteString("<?")
		p.WriteString(t.Target)
		if len(t.Inst) > 0 ***REMOVED***
			p.WriteByte(' ')
			p.Write(t.Inst)
		***REMOVED***
		p.WriteString("?>")
	case Directive:
		if !isValidDirective(t) ***REMOVED***
			return fmt.Errorf("xml: EncodeToken of Directive containing wrong < or > markers")
		***REMOVED***
		p.WriteString("<!")
		p.Write(t)
		p.WriteString(">")
	default:
		return fmt.Errorf("xml: EncodeToken of invalid token type")

	***REMOVED***
	return p.cachedWriteError()
***REMOVED***

// isValidDirective reports whether dir is a valid directive text,
// meaning angle brackets are matched, ignoring comments and strings.
func isValidDirective(dir Directive) bool ***REMOVED***
	var (
		depth     int
		inquote   uint8
		incomment bool
	)
	for i, c := range dir ***REMOVED***
		switch ***REMOVED***
		case incomment:
			if c == '>' ***REMOVED***
				if n := 1 + i - len(endComment); n >= 0 && bytes.Equal(dir[n:i+1], endComment) ***REMOVED***
					incomment = false
				***REMOVED***
			***REMOVED***
			// Just ignore anything in comment
		case inquote != 0:
			if c == inquote ***REMOVED***
				inquote = 0
			***REMOVED***
			// Just ignore anything within quotes
		case c == '\'' || c == '"':
			inquote = c
		case c == '<':
			if i+len(begComment) < len(dir) && bytes.Equal(dir[i:i+len(begComment)], begComment) ***REMOVED***
				incomment = true
			***REMOVED*** else ***REMOVED***
				depth++
			***REMOVED***
		case c == '>':
			if depth == 0 ***REMOVED***
				return false
			***REMOVED***
			depth--
		***REMOVED***
	***REMOVED***
	return depth == 0 && inquote == 0 && !incomment
***REMOVED***

// Flush flushes any buffered XML to the underlying writer.
// See the EncodeToken documentation for details about when it is necessary.
func (enc *Encoder) Flush() error ***REMOVED***
	return enc.p.Flush()
***REMOVED***

type printer struct ***REMOVED***
	*bufio.Writer
	encoder    *Encoder
	seq        int
	indent     string
	prefix     string
	depth      int
	indentedIn bool
	putNewline bool
	defaultNS  string
	attrNS     map[string]string // map prefix -> name space
	attrPrefix map[string]string // map name space -> prefix
	prefixes   []printerPrefix
	tags       []Name
***REMOVED***

// printerPrefix holds a namespace undo record.
// When an element is popped, the prefix record
// is set back to the recorded URL. The empty
// prefix records the URL for the default name space.
//
// The start of an element is recorded with an element
// that has mark=true.
type printerPrefix struct ***REMOVED***
	prefix string
	url    string
	mark   bool
***REMOVED***

func (p *printer) prefixForNS(url string, isAttr bool) string ***REMOVED***
	// The "http://www.w3.org/XML/1998/namespace" name space is predefined as "xml"
	// and must be referred to that way.
	// (The "http://www.w3.org/2000/xmlns/" name space is also predefined as "xmlns",
	// but users should not be trying to use that one directly - that's our job.)
	if url == xmlURL ***REMOVED***
		return "xml"
	***REMOVED***
	if !isAttr && url == p.defaultNS ***REMOVED***
		// We can use the default name space.
		return ""
	***REMOVED***
	return p.attrPrefix[url]
***REMOVED***

// defineNS pushes any namespace definition found in the given attribute.
// If ignoreNonEmptyDefault is true, an xmlns="nonempty"
// attribute will be ignored.
func (p *printer) defineNS(attr Attr, ignoreNonEmptyDefault bool) error ***REMOVED***
	var prefix string
	if attr.Name.Local == "xmlns" ***REMOVED***
		if attr.Name.Space != "" && attr.Name.Space != "xml" && attr.Name.Space != xmlURL ***REMOVED***
			return fmt.Errorf("xml: cannot redefine xmlns attribute prefix")
		***REMOVED***
	***REMOVED*** else if attr.Name.Space == "xmlns" && attr.Name.Local != "" ***REMOVED***
		prefix = attr.Name.Local
		if attr.Value == "" ***REMOVED***
			// Technically, an empty XML namespace is allowed for an attribute.
			// From http://www.w3.org/TR/xml-names11/#scoping-defaulting:
			//
			// 	The attribute value in a namespace declaration for a prefix may be
			//	empty. This has the effect, within the scope of the declaration, of removing
			//	any association of the prefix with a namespace name.
			//
			// However our namespace prefixes here are used only as hints. There's
			// no need to respect the removal of a namespace prefix, so we ignore it.
			return nil
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Ignore: it's not a namespace definition
		return nil
	***REMOVED***
	if prefix == "" ***REMOVED***
		if attr.Value == p.defaultNS ***REMOVED***
			// No need for redefinition.
			return nil
		***REMOVED***
		if attr.Value != "" && ignoreNonEmptyDefault ***REMOVED***
			// We have an xmlns="..." value but
			// it can't define a name space in this context,
			// probably because the element has an empty
			// name space. In this case, we just ignore
			// the name space declaration.
			return nil
		***REMOVED***
	***REMOVED*** else if _, ok := p.attrPrefix[attr.Value]; ok ***REMOVED***
		// There's already a prefix for the given name space,
		// so use that. This prevents us from
		// having two prefixes for the same name space
		// so attrNS and attrPrefix can remain bijective.
		return nil
	***REMOVED***
	p.pushPrefix(prefix, attr.Value)
	return nil
***REMOVED***

// createNSPrefix creates a name space prefix attribute
// to use for the given name space, defining a new prefix
// if necessary.
// If isAttr is true, the prefix is to be created for an attribute
// prefix, which means that the default name space cannot
// be used.
func (p *printer) createNSPrefix(url string, isAttr bool) ***REMOVED***
	if _, ok := p.attrPrefix[url]; ok ***REMOVED***
		// We already have a prefix for the given URL.
		return
	***REMOVED***
	switch ***REMOVED***
	case !isAttr && url == p.defaultNS:
		// We can use the default name space.
		return
	case url == "":
		// The only way we can encode names in the empty
		// name space is by using the default name space,
		// so we must use that.
		if p.defaultNS != "" ***REMOVED***
			// The default namespace is non-empty, so we
			// need to set it to empty.
			p.pushPrefix("", "")
		***REMOVED***
		return
	case url == xmlURL:
		return
	***REMOVED***
	// TODO If the URL is an existing prefix, we could
	// use it as is. That would enable the
	// marshaling of elements that had been unmarshaled
	// and with a name space prefix that was not found.
	// although technically it would be incorrect.

	// Pick a name. We try to use the final element of the path
	// but fall back to _.
	prefix := strings.TrimRight(url, "/")
	if i := strings.LastIndex(prefix, "/"); i >= 0 ***REMOVED***
		prefix = prefix[i+1:]
	***REMOVED***
	if prefix == "" || !isName([]byte(prefix)) || strings.Contains(prefix, ":") ***REMOVED***
		prefix = "_"
	***REMOVED***
	if strings.HasPrefix(prefix, "xml") ***REMOVED***
		// xmlanything is reserved.
		prefix = "_" + prefix
	***REMOVED***
	if p.attrNS[prefix] != "" ***REMOVED***
		// Name is taken. Find a better one.
		for p.seq++; ; p.seq++ ***REMOVED***
			if id := prefix + "_" + strconv.Itoa(p.seq); p.attrNS[id] == "" ***REMOVED***
				prefix = id
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	p.pushPrefix(prefix, url)
***REMOVED***

// writeNamespaces writes xmlns attributes for all the
// namespace prefixes that have been defined in
// the current element.
func (p *printer) writeNamespaces() ***REMOVED***
	for i := len(p.prefixes) - 1; i >= 0; i-- ***REMOVED***
		prefix := p.prefixes[i]
		if prefix.mark ***REMOVED***
			return
		***REMOVED***
		p.WriteString(" ")
		if prefix.prefix == "" ***REMOVED***
			// Default name space.
			p.WriteString(`xmlns="`)
		***REMOVED*** else ***REMOVED***
			p.WriteString("xmlns:")
			p.WriteString(prefix.prefix)
			p.WriteString(`="`)
		***REMOVED***
		EscapeText(p, []byte(p.nsForPrefix(prefix.prefix)))
		p.WriteString(`"`)
	***REMOVED***
***REMOVED***

// pushPrefix pushes a new prefix on the prefix stack
// without checking to see if it is already defined.
func (p *printer) pushPrefix(prefix, url string) ***REMOVED***
	p.prefixes = append(p.prefixes, printerPrefix***REMOVED***
		prefix: prefix,
		url:    p.nsForPrefix(prefix),
	***REMOVED***)
	p.setAttrPrefix(prefix, url)
***REMOVED***

// nsForPrefix returns the name space for the given
// prefix. Note that this is not valid for the
// empty attribute prefix, which always has an empty
// name space.
func (p *printer) nsForPrefix(prefix string) string ***REMOVED***
	if prefix == "" ***REMOVED***
		return p.defaultNS
	***REMOVED***
	return p.attrNS[prefix]
***REMOVED***

// markPrefix marks the start of an element on the prefix
// stack.
func (p *printer) markPrefix() ***REMOVED***
	p.prefixes = append(p.prefixes, printerPrefix***REMOVED***
		mark: true,
	***REMOVED***)
***REMOVED***

// popPrefix pops all defined prefixes for the current
// element.
func (p *printer) popPrefix() ***REMOVED***
	for len(p.prefixes) > 0 ***REMOVED***
		prefix := p.prefixes[len(p.prefixes)-1]
		p.prefixes = p.prefixes[:len(p.prefixes)-1]
		if prefix.mark ***REMOVED***
			break
		***REMOVED***
		p.setAttrPrefix(prefix.prefix, prefix.url)
	***REMOVED***
***REMOVED***

// setAttrPrefix sets an attribute name space prefix.
// If url is empty, the attribute is removed.
// If prefix is empty, the default name space is set.
func (p *printer) setAttrPrefix(prefix, url string) ***REMOVED***
	if prefix == "" ***REMOVED***
		p.defaultNS = url
		return
	***REMOVED***
	if url == "" ***REMOVED***
		delete(p.attrPrefix, p.attrNS[prefix])
		delete(p.attrNS, prefix)
		return
	***REMOVED***
	if p.attrPrefix == nil ***REMOVED***
		// Need to define a new name space.
		p.attrPrefix = make(map[string]string)
		p.attrNS = make(map[string]string)
	***REMOVED***
	// Remove any old prefix value. This is OK because we maintain a
	// strict one-to-one mapping between prefix and URL (see
	// defineNS)
	delete(p.attrPrefix, p.attrNS[prefix])
	p.attrPrefix[url] = prefix
	p.attrNS[prefix] = url
***REMOVED***

var (
	marshalerType     = reflect.TypeOf((*Marshaler)(nil)).Elem()
	marshalerAttrType = reflect.TypeOf((*MarshalerAttr)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

// marshalValue writes one or more XML elements representing val.
// If val was obtained from a struct field, finfo must have its details.
func (p *printer) marshalValue(val reflect.Value, finfo *fieldInfo, startTemplate *StartElement) error ***REMOVED***
	if startTemplate != nil && startTemplate.Name.Local == "" ***REMOVED***
		return fmt.Errorf("xml: EncodeElement of StartElement with missing name")
	***REMOVED***

	if !val.IsValid() ***REMOVED***
		return nil
	***REMOVED***
	if finfo != nil && finfo.flags&fOmitEmpty != 0 && isEmptyValue(val) ***REMOVED***
		return nil
	***REMOVED***

	// Drill into interfaces and pointers.
	// This can turn into an infinite loop given a cyclic chain,
	// but it matches the Go 1 behavior.
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr ***REMOVED***
		if val.IsNil() ***REMOVED***
			return nil
		***REMOVED***
		val = val.Elem()
	***REMOVED***

	kind := val.Kind()
	typ := val.Type()

	// Check for marshaler.
	if val.CanInterface() && typ.Implements(marshalerType) ***REMOVED***
		return p.marshalInterface(val.Interface().(Marshaler), p.defaultStart(typ, finfo, startTemplate))
	***REMOVED***
	if val.CanAddr() ***REMOVED***
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(marshalerType) ***REMOVED***
			return p.marshalInterface(pv.Interface().(Marshaler), p.defaultStart(pv.Type(), finfo, startTemplate))
		***REMOVED***
	***REMOVED***

	// Check for text marshaler.
	if val.CanInterface() && typ.Implements(textMarshalerType) ***REMOVED***
		return p.marshalTextInterface(val.Interface().(encoding.TextMarshaler), p.defaultStart(typ, finfo, startTemplate))
	***REMOVED***
	if val.CanAddr() ***REMOVED***
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(textMarshalerType) ***REMOVED***
			return p.marshalTextInterface(pv.Interface().(encoding.TextMarshaler), p.defaultStart(pv.Type(), finfo, startTemplate))
		***REMOVED***
	***REMOVED***

	// Slices and arrays iterate over the elements. They do not have an enclosing tag.
	if (kind == reflect.Slice || kind == reflect.Array) && typ.Elem().Kind() != reflect.Uint8 ***REMOVED***
		for i, n := 0, val.Len(); i < n; i++ ***REMOVED***
			if err := p.marshalValue(val.Index(i), finfo, startTemplate); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	tinfo, err := getTypeInfo(typ)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Create start element.
	// Precedence for the XML element name is:
	// 0. startTemplate
	// 1. XMLName field in underlying struct;
	// 2. field name/tag in the struct field; and
	// 3. type name
	var start StartElement

	// explicitNS records whether the element's name space has been
	// explicitly set (for example an XMLName field).
	explicitNS := false

	if startTemplate != nil ***REMOVED***
		start.Name = startTemplate.Name
		explicitNS = true
		start.Attr = append(start.Attr, startTemplate.Attr...)
	***REMOVED*** else if tinfo.xmlname != nil ***REMOVED***
		xmlname := tinfo.xmlname
		if xmlname.name != "" ***REMOVED***
			start.Name.Space, start.Name.Local = xmlname.xmlns, xmlname.name
		***REMOVED*** else if v, ok := xmlname.value(val).Interface().(Name); ok && v.Local != "" ***REMOVED***
			start.Name = v
		***REMOVED***
		explicitNS = true
	***REMOVED***
	if start.Name.Local == "" && finfo != nil ***REMOVED***
		start.Name.Local = finfo.name
		if finfo.xmlns != "" ***REMOVED***
			start.Name.Space = finfo.xmlns
			explicitNS = true
		***REMOVED***
	***REMOVED***
	if start.Name.Local == "" ***REMOVED***
		name := typ.Name()
		if name == "" ***REMOVED***
			return &UnsupportedTypeError***REMOVED***typ***REMOVED***
		***REMOVED***
		start.Name.Local = name
	***REMOVED***

	// defaultNS records the default name space as set by a xmlns="..."
	// attribute. We don't set p.defaultNS because we want to let
	// the attribute writing code (in p.defineNS) be solely responsible
	// for maintaining that.
	defaultNS := p.defaultNS

	// Attributes
	for i := range tinfo.fields ***REMOVED***
		finfo := &tinfo.fields[i]
		if finfo.flags&fAttr == 0 ***REMOVED***
			continue
		***REMOVED***
		attr, err := p.fieldAttr(finfo, val)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if attr.Name.Local == "" ***REMOVED***
			continue
		***REMOVED***
		start.Attr = append(start.Attr, attr)
		if attr.Name.Space == "" && attr.Name.Local == "xmlns" ***REMOVED***
			defaultNS = attr.Value
		***REMOVED***
	***REMOVED***
	if !explicitNS ***REMOVED***
		// Historic behavior: elements use the default name space
		// they are contained in by default.
		start.Name.Space = defaultNS
	***REMOVED***
	// Historic behaviour: an element that's in a namespace sets
	// the default namespace for all elements contained within it.
	start.setDefaultNamespace()

	if err := p.writeStart(&start); err != nil ***REMOVED***
		return err
	***REMOVED***

	if val.Kind() == reflect.Struct ***REMOVED***
		err = p.marshalStruct(tinfo, val)
	***REMOVED*** else ***REMOVED***
		s, b, err1 := p.marshalSimple(typ, val)
		if err1 != nil ***REMOVED***
			err = err1
		***REMOVED*** else if b != nil ***REMOVED***
			EscapeText(p, b)
		***REMOVED*** else ***REMOVED***
			p.EscapeString(s)
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := p.writeEnd(start.Name); err != nil ***REMOVED***
		return err
	***REMOVED***

	return p.cachedWriteError()
***REMOVED***

// fieldAttr returns the attribute of the given field.
// If the returned attribute has an empty Name.Local,
// it should not be used.
// The given value holds the value containing the field.
func (p *printer) fieldAttr(finfo *fieldInfo, val reflect.Value) (Attr, error) ***REMOVED***
	fv := finfo.value(val)
	name := Name***REMOVED***Space: finfo.xmlns, Local: finfo.name***REMOVED***
	if finfo.flags&fOmitEmpty != 0 && isEmptyValue(fv) ***REMOVED***
		return Attr***REMOVED******REMOVED***, nil
	***REMOVED***
	if fv.Kind() == reflect.Interface && fv.IsNil() ***REMOVED***
		return Attr***REMOVED******REMOVED***, nil
	***REMOVED***
	if fv.CanInterface() && fv.Type().Implements(marshalerAttrType) ***REMOVED***
		attr, err := fv.Interface().(MarshalerAttr).MarshalXMLAttr(name)
		return attr, err
	***REMOVED***
	if fv.CanAddr() ***REMOVED***
		pv := fv.Addr()
		if pv.CanInterface() && pv.Type().Implements(marshalerAttrType) ***REMOVED***
			attr, err := pv.Interface().(MarshalerAttr).MarshalXMLAttr(name)
			return attr, err
		***REMOVED***
	***REMOVED***
	if fv.CanInterface() && fv.Type().Implements(textMarshalerType) ***REMOVED***
		text, err := fv.Interface().(encoding.TextMarshaler).MarshalText()
		if err != nil ***REMOVED***
			return Attr***REMOVED******REMOVED***, err
		***REMOVED***
		return Attr***REMOVED***name, string(text)***REMOVED***, nil
	***REMOVED***
	if fv.CanAddr() ***REMOVED***
		pv := fv.Addr()
		if pv.CanInterface() && pv.Type().Implements(textMarshalerType) ***REMOVED***
			text, err := pv.Interface().(encoding.TextMarshaler).MarshalText()
			if err != nil ***REMOVED***
				return Attr***REMOVED******REMOVED***, err
			***REMOVED***
			return Attr***REMOVED***name, string(text)***REMOVED***, nil
		***REMOVED***
	***REMOVED***
	// Dereference or skip nil pointer, interface values.
	switch fv.Kind() ***REMOVED***
	case reflect.Ptr, reflect.Interface:
		if fv.IsNil() ***REMOVED***
			return Attr***REMOVED******REMOVED***, nil
		***REMOVED***
		fv = fv.Elem()
	***REMOVED***
	s, b, err := p.marshalSimple(fv.Type(), fv)
	if err != nil ***REMOVED***
		return Attr***REMOVED******REMOVED***, err
	***REMOVED***
	if b != nil ***REMOVED***
		s = string(b)
	***REMOVED***
	return Attr***REMOVED***name, s***REMOVED***, nil
***REMOVED***

// defaultStart returns the default start element to use,
// given the reflect type, field info, and start template.
func (p *printer) defaultStart(typ reflect.Type, finfo *fieldInfo, startTemplate *StartElement) StartElement ***REMOVED***
	var start StartElement
	// Precedence for the XML element name is as above,
	// except that we do not look inside structs for the first field.
	if startTemplate != nil ***REMOVED***
		start.Name = startTemplate.Name
		start.Attr = append(start.Attr, startTemplate.Attr...)
	***REMOVED*** else if finfo != nil && finfo.name != "" ***REMOVED***
		start.Name.Local = finfo.name
		start.Name.Space = finfo.xmlns
	***REMOVED*** else if typ.Name() != "" ***REMOVED***
		start.Name.Local = typ.Name()
	***REMOVED*** else ***REMOVED***
		// Must be a pointer to a named type,
		// since it has the Marshaler methods.
		start.Name.Local = typ.Elem().Name()
	***REMOVED***
	// Historic behaviour: elements use the name space of
	// the element they are contained in by default.
	if start.Name.Space == "" ***REMOVED***
		start.Name.Space = p.defaultNS
	***REMOVED***
	start.setDefaultNamespace()
	return start
***REMOVED***

// marshalInterface marshals a Marshaler interface value.
func (p *printer) marshalInterface(val Marshaler, start StartElement) error ***REMOVED***
	// Push a marker onto the tag stack so that MarshalXML
	// cannot close the XML tags that it did not open.
	p.tags = append(p.tags, Name***REMOVED******REMOVED***)
	n := len(p.tags)

	err := val.MarshalXML(p.encoder, start)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Make sure MarshalXML closed all its tags. p.tags[n-1] is the mark.
	if len(p.tags) > n ***REMOVED***
		return fmt.Errorf("xml: %s.MarshalXML wrote invalid XML: <%s> not closed", receiverType(val), p.tags[len(p.tags)-1].Local)
	***REMOVED***
	p.tags = p.tags[:n-1]
	return nil
***REMOVED***

// marshalTextInterface marshals a TextMarshaler interface value.
func (p *printer) marshalTextInterface(val encoding.TextMarshaler, start StartElement) error ***REMOVED***
	if err := p.writeStart(&start); err != nil ***REMOVED***
		return err
	***REMOVED***
	text, err := val.MarshalText()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	EscapeText(p, text)
	return p.writeEnd(start.Name)
***REMOVED***

// writeStart writes the given start element.
func (p *printer) writeStart(start *StartElement) error ***REMOVED***
	if start.Name.Local == "" ***REMOVED***
		return fmt.Errorf("xml: start tag with no name")
	***REMOVED***

	p.tags = append(p.tags, start.Name)
	p.markPrefix()
	// Define any name spaces explicitly declared in the attributes.
	// We do this as a separate pass so that explicitly declared prefixes
	// will take precedence over implicitly declared prefixes
	// regardless of the order of the attributes.
	ignoreNonEmptyDefault := start.Name.Space == ""
	for _, attr := range start.Attr ***REMOVED***
		if err := p.defineNS(attr, ignoreNonEmptyDefault); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// Define any new name spaces implied by the attributes.
	for _, attr := range start.Attr ***REMOVED***
		name := attr.Name
		// From http://www.w3.org/TR/xml-names11/#defaulting
		// "Default namespace declarations do not apply directly
		// to attribute names; the interpretation of unprefixed
		// attributes is determined by the element on which they
		// appear."
		// This means we don't need to create a new namespace
		// when an attribute name space is empty.
		if name.Space != "" && !name.isNamespace() ***REMOVED***
			p.createNSPrefix(name.Space, true)
		***REMOVED***
	***REMOVED***
	p.createNSPrefix(start.Name.Space, false)

	p.writeIndent(1)
	p.WriteByte('<')
	p.writeName(start.Name, false)
	p.writeNamespaces()
	for _, attr := range start.Attr ***REMOVED***
		name := attr.Name
		if name.Local == "" || name.isNamespace() ***REMOVED***
			// Namespaces have already been written by writeNamespaces above.
			continue
		***REMOVED***
		p.WriteByte(' ')
		p.writeName(name, true)
		p.WriteString(`="`)
		p.EscapeString(attr.Value)
		p.WriteByte('"')
	***REMOVED***
	p.WriteByte('>')
	return nil
***REMOVED***

// writeName writes the given name. It assumes
// that p.createNSPrefix(name) has already been called.
func (p *printer) writeName(name Name, isAttr bool) ***REMOVED***
	if prefix := p.prefixForNS(name.Space, isAttr); prefix != "" ***REMOVED***
		p.WriteString(prefix)
		p.WriteByte(':')
	***REMOVED***
	p.WriteString(name.Local)
***REMOVED***

func (p *printer) writeEnd(name Name) error ***REMOVED***
	if name.Local == "" ***REMOVED***
		return fmt.Errorf("xml: end tag with no name")
	***REMOVED***
	if len(p.tags) == 0 || p.tags[len(p.tags)-1].Local == "" ***REMOVED***
		return fmt.Errorf("xml: end tag </%s> without start tag", name.Local)
	***REMOVED***
	if top := p.tags[len(p.tags)-1]; top != name ***REMOVED***
		if top.Local != name.Local ***REMOVED***
			return fmt.Errorf("xml: end tag </%s> does not match start tag <%s>", name.Local, top.Local)
		***REMOVED***
		return fmt.Errorf("xml: end tag </%s> in namespace %s does not match start tag <%s> in namespace %s", name.Local, name.Space, top.Local, top.Space)
	***REMOVED***
	p.tags = p.tags[:len(p.tags)-1]

	p.writeIndent(-1)
	p.WriteByte('<')
	p.WriteByte('/')
	p.writeName(name, false)
	p.WriteByte('>')
	p.popPrefix()
	return nil
***REMOVED***

func (p *printer) marshalSimple(typ reflect.Type, val reflect.Value) (string, []byte, error) ***REMOVED***
	switch val.Kind() ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10), nil, nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()), nil, nil
	case reflect.String:
		return val.String(), nil, nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil, nil
	case reflect.Array:
		if typ.Elem().Kind() != reflect.Uint8 ***REMOVED***
			break
		***REMOVED***
		// [...]byte
		var bytes []byte
		if val.CanAddr() ***REMOVED***
			bytes = val.Slice(0, val.Len()).Bytes()
		***REMOVED*** else ***REMOVED***
			bytes = make([]byte, val.Len())
			reflect.Copy(reflect.ValueOf(bytes), val)
		***REMOVED***
		return "", bytes, nil
	case reflect.Slice:
		if typ.Elem().Kind() != reflect.Uint8 ***REMOVED***
			break
		***REMOVED***
		// []byte
		return "", val.Bytes(), nil
	***REMOVED***
	return "", nil, &UnsupportedTypeError***REMOVED***typ***REMOVED***
***REMOVED***

var ddBytes = []byte("--")

func (p *printer) marshalStruct(tinfo *typeInfo, val reflect.Value) error ***REMOVED***
	s := parentStack***REMOVED***p: p***REMOVED***
	for i := range tinfo.fields ***REMOVED***
		finfo := &tinfo.fields[i]
		if finfo.flags&fAttr != 0 ***REMOVED***
			continue
		***REMOVED***
		vf := finfo.value(val)

		// Dereference or skip nil pointer, interface values.
		switch vf.Kind() ***REMOVED***
		case reflect.Ptr, reflect.Interface:
			if !vf.IsNil() ***REMOVED***
				vf = vf.Elem()
			***REMOVED***
		***REMOVED***

		switch finfo.flags & fMode ***REMOVED***
		case fCharData:
			if err := s.setParents(&noField, reflect.Value***REMOVED******REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
			if vf.CanInterface() && vf.Type().Implements(textMarshalerType) ***REMOVED***
				data, err := vf.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				Escape(p, data)
				continue
			***REMOVED***
			if vf.CanAddr() ***REMOVED***
				pv := vf.Addr()
				if pv.CanInterface() && pv.Type().Implements(textMarshalerType) ***REMOVED***
					data, err := pv.Interface().(encoding.TextMarshaler).MarshalText()
					if err != nil ***REMOVED***
						return err
					***REMOVED***
					Escape(p, data)
					continue
				***REMOVED***
			***REMOVED***
			var scratch [64]byte
			switch vf.Kind() ***REMOVED***
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				Escape(p, strconv.AppendInt(scratch[:0], vf.Int(), 10))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				Escape(p, strconv.AppendUint(scratch[:0], vf.Uint(), 10))
			case reflect.Float32, reflect.Float64:
				Escape(p, strconv.AppendFloat(scratch[:0], vf.Float(), 'g', -1, vf.Type().Bits()))
			case reflect.Bool:
				Escape(p, strconv.AppendBool(scratch[:0], vf.Bool()))
			case reflect.String:
				if err := EscapeText(p, []byte(vf.String())); err != nil ***REMOVED***
					return err
				***REMOVED***
			case reflect.Slice:
				if elem, ok := vf.Interface().([]byte); ok ***REMOVED***
					if err := EscapeText(p, elem); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED***
			continue

		case fComment:
			if err := s.setParents(&noField, reflect.Value***REMOVED******REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
			k := vf.Kind()
			if !(k == reflect.String || k == reflect.Slice && vf.Type().Elem().Kind() == reflect.Uint8) ***REMOVED***
				return fmt.Errorf("xml: bad type for comment field of %s", val.Type())
			***REMOVED***
			if vf.Len() == 0 ***REMOVED***
				continue
			***REMOVED***
			p.writeIndent(0)
			p.WriteString("<!--")
			dashDash := false
			dashLast := false
			switch k ***REMOVED***
			case reflect.String:
				s := vf.String()
				dashDash = strings.Index(s, "--") >= 0
				dashLast = s[len(s)-1] == '-'
				if !dashDash ***REMOVED***
					p.WriteString(s)
				***REMOVED***
			case reflect.Slice:
				b := vf.Bytes()
				dashDash = bytes.Index(b, ddBytes) >= 0
				dashLast = b[len(b)-1] == '-'
				if !dashDash ***REMOVED***
					p.Write(b)
				***REMOVED***
			default:
				panic("can't happen")
			***REMOVED***
			if dashDash ***REMOVED***
				return fmt.Errorf(`xml: comments must not contain "--"`)
			***REMOVED***
			if dashLast ***REMOVED***
				// "--->" is invalid grammar. Make it "- -->"
				p.WriteByte(' ')
			***REMOVED***
			p.WriteString("-->")
			continue

		case fInnerXml:
			iface := vf.Interface()
			switch raw := iface.(type) ***REMOVED***
			case []byte:
				p.Write(raw)
				continue
			case string:
				p.WriteString(raw)
				continue
			***REMOVED***

		case fElement, fElement | fAny:
			if err := s.setParents(finfo, vf); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if err := p.marshalValue(vf, finfo, nil); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err := s.setParents(&noField, reflect.Value***REMOVED******REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	return p.cachedWriteError()
***REMOVED***

var noField fieldInfo

// return the bufio Writer's cached write error
func (p *printer) cachedWriteError() error ***REMOVED***
	_, err := p.Write(nil)
	return err
***REMOVED***

func (p *printer) writeIndent(depthDelta int) ***REMOVED***
	if len(p.prefix) == 0 && len(p.indent) == 0 ***REMOVED***
		return
	***REMOVED***
	if depthDelta < 0 ***REMOVED***
		p.depth--
		if p.indentedIn ***REMOVED***
			p.indentedIn = false
			return
		***REMOVED***
		p.indentedIn = false
	***REMOVED***
	if p.putNewline ***REMOVED***
		p.WriteByte('\n')
	***REMOVED*** else ***REMOVED***
		p.putNewline = true
	***REMOVED***
	if len(p.prefix) > 0 ***REMOVED***
		p.WriteString(p.prefix)
	***REMOVED***
	if len(p.indent) > 0 ***REMOVED***
		for i := 0; i < p.depth; i++ ***REMOVED***
			p.WriteString(p.indent)
		***REMOVED***
	***REMOVED***
	if depthDelta > 0 ***REMOVED***
		p.depth++
		p.indentedIn = true
	***REMOVED***
***REMOVED***

type parentStack struct ***REMOVED***
	p       *printer
	xmlns   string
	parents []string
***REMOVED***

// setParents sets the stack of current parents to those found in finfo.
// It only writes the start elements if vf holds a non-nil value.
// If finfo is &noField, it pops all elements.
func (s *parentStack) setParents(finfo *fieldInfo, vf reflect.Value) error ***REMOVED***
	xmlns := s.p.defaultNS
	if finfo.xmlns != "" ***REMOVED***
		xmlns = finfo.xmlns
	***REMOVED***
	commonParents := 0
	if xmlns == s.xmlns ***REMOVED***
		for ; commonParents < len(finfo.parents) && commonParents < len(s.parents); commonParents++ ***REMOVED***
			if finfo.parents[commonParents] != s.parents[commonParents] ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Pop off any parents that aren't in common with the previous field.
	for i := len(s.parents) - 1; i >= commonParents; i-- ***REMOVED***
		if err := s.p.writeEnd(Name***REMOVED***
			Space: s.xmlns,
			Local: s.parents[i],
		***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	s.parents = finfo.parents
	s.xmlns = xmlns
	if commonParents >= len(s.parents) ***REMOVED***
		// No new elements to push.
		return nil
	***REMOVED***
	if (vf.Kind() == reflect.Ptr || vf.Kind() == reflect.Interface) && vf.IsNil() ***REMOVED***
		// The element is nil, so no need for the start elements.
		s.parents = s.parents[:commonParents]
		return nil
	***REMOVED***
	// Push any new parents required.
	for _, name := range s.parents[commonParents:] ***REMOVED***
		start := &StartElement***REMOVED***
			Name: Name***REMOVED***
				Space: s.xmlns,
				Local: name,
			***REMOVED***,
		***REMOVED***
		// Set the default name space for parent elements
		// to match what we do with other elements.
		if s.xmlns != s.p.defaultNS ***REMOVED***
			start.setDefaultNamespace()
		***REMOVED***
		if err := s.p.writeStart(start); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// A MarshalXMLError is returned when Marshal encounters a type
// that cannot be converted into XML.
type UnsupportedTypeError struct ***REMOVED***
	Type reflect.Type
***REMOVED***

func (e *UnsupportedTypeError) Error() string ***REMOVED***
	return "xml: unsupported type: " + e.Type.String()
***REMOVED***

func isEmptyValue(v reflect.Value) bool ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	***REMOVED***
	return false
***REMOVED***
