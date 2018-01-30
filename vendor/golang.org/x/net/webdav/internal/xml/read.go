// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// BUG(rsc): Mapping between XML elements and data structures is inherently flawed:
// an XML element is an order-dependent collection of anonymous
// values, while a data structure is an order-independent collection
// of named values.
// See package json for a textual representation more suitable
// to data structures.

// Unmarshal parses the XML-encoded data and stores the result in
// the value pointed to by v, which must be an arbitrary struct,
// slice, or string. Well-formed data that does not fit into v is
// discarded.
//
// Because Unmarshal uses the reflect package, it can only assign
// to exported (upper case) fields. Unmarshal uses a case-sensitive
// comparison to match XML element names to tag values and struct
// field names.
//
// Unmarshal maps an XML element to a struct using the following rules.
// In the rules, the tag of a field refers to the value associated with the
// key 'xml' in the struct field's tag (see the example above).
//
//   * If the struct has a field of type []byte or string with tag
//      ",innerxml", Unmarshal accumulates the raw XML nested inside the
//      element in that field. The rest of the rules still apply.
//
//   * If the struct has a field named XMLName of type xml.Name,
//      Unmarshal records the element name in that field.
//
//   * If the XMLName field has an associated tag of the form
//      "name" or "namespace-URL name", the XML element must have
//      the given name (and, optionally, name space) or else Unmarshal
//      returns an error.
//
//   * If the XML element has an attribute whose name matches a
//      struct field name with an associated tag containing ",attr" or
//      the explicit name in a struct field tag of the form "name,attr",
//      Unmarshal records the attribute value in that field.
//
//   * If the XML element contains character data, that data is
//      accumulated in the first struct field that has tag ",chardata".
//      The struct field may have type []byte or string.
//      If there is no such field, the character data is discarded.
//
//   * If the XML element contains comments, they are accumulated in
//      the first struct field that has tag ",comment".  The struct
//      field may have type []byte or string. If there is no such
//      field, the comments are discarded.
//
//   * If the XML element contains a sub-element whose name matches
//      the prefix of a tag formatted as "a" or "a>b>c", unmarshal
//      will descend into the XML structure looking for elements with the
//      given names, and will map the innermost elements to that struct
//      field. A tag starting with ">" is equivalent to one starting
//      with the field name followed by ">".
//
//   * If the XML element contains a sub-element whose name matches
//      a struct field's XMLName tag and the struct field has no
//      explicit name tag as per the previous rule, unmarshal maps
//      the sub-element to that struct field.
//
//   * If the XML element contains a sub-element whose name matches a
//      field without any mode flags (",attr", ",chardata", etc), Unmarshal
//      maps the sub-element to that struct field.
//
//   * If the XML element contains a sub-element that hasn't matched any
//      of the above rules and the struct has a field with tag ",any",
//      unmarshal maps the sub-element to that struct field.
//
//   * An anonymous struct field is handled as if the fields of its
//      value were part of the outer struct.
//
//   * A struct field with tag "-" is never unmarshalled into.
//
// Unmarshal maps an XML element to a string or []byte by saving the
// concatenation of that element's character data in the string or
// []byte. The saved []byte is never nil.
//
// Unmarshal maps an attribute value to a string or []byte by saving
// the value in the string or slice.
//
// Unmarshal maps an XML element to a slice by extending the length of
// the slice and mapping the element to the newly created value.
//
// Unmarshal maps an XML element or attribute value to a bool by
// setting it to the boolean value represented by the string.
//
// Unmarshal maps an XML element or attribute value to an integer or
// floating-point field by setting the field to the result of
// interpreting the string value in decimal. There is no check for
// overflow.
//
// Unmarshal maps an XML element to an xml.Name by recording the
// element name.
//
// Unmarshal maps an XML element to a pointer by setting the pointer
// to a freshly allocated value and then mapping the element to that value.
//
func Unmarshal(data []byte, v interface***REMOVED******REMOVED***) error ***REMOVED***
	return NewDecoder(bytes.NewReader(data)).Decode(v)
***REMOVED***

// Decode works like xml.Unmarshal, except it reads the decoder
// stream to find the start element.
func (d *Decoder) Decode(v interface***REMOVED******REMOVED***) error ***REMOVED***
	return d.DecodeElement(v, nil)
***REMOVED***

// DecodeElement works like xml.Unmarshal except that it takes
// a pointer to the start XML element to decode into v.
// It is useful when a client reads some raw XML tokens itself
// but also wants to defer to Unmarshal for some elements.
func (d *Decoder) DecodeElement(v interface***REMOVED******REMOVED***, start *StartElement) error ***REMOVED***
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr ***REMOVED***
		return errors.New("non-pointer passed to Unmarshal")
	***REMOVED***
	return d.unmarshal(val.Elem(), start)
***REMOVED***

// An UnmarshalError represents an error in the unmarshalling process.
type UnmarshalError string

func (e UnmarshalError) Error() string ***REMOVED*** return string(e) ***REMOVED***

// Unmarshaler is the interface implemented by objects that can unmarshal
// an XML element description of themselves.
//
// UnmarshalXML decodes a single XML element
// beginning with the given start element.
// If it returns an error, the outer call to Unmarshal stops and
// returns that error.
// UnmarshalXML must consume exactly one XML element.
// One common implementation strategy is to unmarshal into
// a separate value with a layout matching the expected XML
// using d.DecodeElement,  and then to copy the data from
// that value into the receiver.
// Another common strategy is to use d.Token to process the
// XML object one token at a time.
// UnmarshalXML may not use d.RawToken.
type Unmarshaler interface ***REMOVED***
	UnmarshalXML(d *Decoder, start StartElement) error
***REMOVED***

// UnmarshalerAttr is the interface implemented by objects that can unmarshal
// an XML attribute description of themselves.
//
// UnmarshalXMLAttr decodes a single XML attribute.
// If it returns an error, the outer call to Unmarshal stops and
// returns that error.
// UnmarshalXMLAttr is used only for struct fields with the
// "attr" option in the field tag.
type UnmarshalerAttr interface ***REMOVED***
	UnmarshalXMLAttr(attr Attr) error
***REMOVED***

// receiverType returns the receiver type to use in an expression like "%s.MethodName".
func receiverType(val interface***REMOVED******REMOVED***) string ***REMOVED***
	t := reflect.TypeOf(val)
	if t.Name() != "" ***REMOVED***
		return t.String()
	***REMOVED***
	return "(" + t.String() + ")"
***REMOVED***

// unmarshalInterface unmarshals a single XML element into val.
// start is the opening tag of the element.
func (p *Decoder) unmarshalInterface(val Unmarshaler, start *StartElement) error ***REMOVED***
	// Record that decoder must stop at end tag corresponding to start.
	p.pushEOF()

	p.unmarshalDepth++
	err := val.UnmarshalXML(p, *start)
	p.unmarshalDepth--
	if err != nil ***REMOVED***
		p.popEOF()
		return err
	***REMOVED***

	if !p.popEOF() ***REMOVED***
		return fmt.Errorf("xml: %s.UnmarshalXML did not consume entire <%s> element", receiverType(val), start.Name.Local)
	***REMOVED***

	return nil
***REMOVED***

// unmarshalTextInterface unmarshals a single XML element into val.
// The chardata contained in the element (but not its children)
// is passed to the text unmarshaler.
func (p *Decoder) unmarshalTextInterface(val encoding.TextUnmarshaler, start *StartElement) error ***REMOVED***
	var buf []byte
	depth := 1
	for depth > 0 ***REMOVED***
		t, err := p.Token()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch t := t.(type) ***REMOVED***
		case CharData:
			if depth == 1 ***REMOVED***
				buf = append(buf, t...)
			***REMOVED***
		case StartElement:
			depth++
		case EndElement:
			depth--
		***REMOVED***
	***REMOVED***
	return val.UnmarshalText(buf)
***REMOVED***

// unmarshalAttr unmarshals a single XML attribute into val.
func (p *Decoder) unmarshalAttr(val reflect.Value, attr Attr) error ***REMOVED***
	if val.Kind() == reflect.Ptr ***REMOVED***
		if val.IsNil() ***REMOVED***
			val.Set(reflect.New(val.Type().Elem()))
		***REMOVED***
		val = val.Elem()
	***REMOVED***

	if val.CanInterface() && val.Type().Implements(unmarshalerAttrType) ***REMOVED***
		// This is an unmarshaler with a non-pointer receiver,
		// so it's likely to be incorrect, but we do what we're told.
		return val.Interface().(UnmarshalerAttr).UnmarshalXMLAttr(attr)
	***REMOVED***
	if val.CanAddr() ***REMOVED***
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(unmarshalerAttrType) ***REMOVED***
			return pv.Interface().(UnmarshalerAttr).UnmarshalXMLAttr(attr)
		***REMOVED***
	***REMOVED***

	// Not an UnmarshalerAttr; try encoding.TextUnmarshaler.
	if val.CanInterface() && val.Type().Implements(textUnmarshalerType) ***REMOVED***
		// This is an unmarshaler with a non-pointer receiver,
		// so it's likely to be incorrect, but we do what we're told.
		return val.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(attr.Value))
	***REMOVED***
	if val.CanAddr() ***REMOVED***
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(textUnmarshalerType) ***REMOVED***
			return pv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(attr.Value))
		***REMOVED***
	***REMOVED***

	copyValue(val, []byte(attr.Value))
	return nil
***REMOVED***

var (
	unmarshalerType     = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
	unmarshalerAttrType = reflect.TypeOf((*UnmarshalerAttr)(nil)).Elem()
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

// Unmarshal a single XML element into val.
func (p *Decoder) unmarshal(val reflect.Value, start *StartElement) error ***REMOVED***
	// Find start element if we need it.
	if start == nil ***REMOVED***
		for ***REMOVED***
			tok, err := p.Token()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if t, ok := tok.(StartElement); ok ***REMOVED***
				start = &t
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Load value from interface, but only if the result will be
	// usefully addressable.
	if val.Kind() == reflect.Interface && !val.IsNil() ***REMOVED***
		e := val.Elem()
		if e.Kind() == reflect.Ptr && !e.IsNil() ***REMOVED***
			val = e
		***REMOVED***
	***REMOVED***

	if val.Kind() == reflect.Ptr ***REMOVED***
		if val.IsNil() ***REMOVED***
			val.Set(reflect.New(val.Type().Elem()))
		***REMOVED***
		val = val.Elem()
	***REMOVED***

	if val.CanInterface() && val.Type().Implements(unmarshalerType) ***REMOVED***
		// This is an unmarshaler with a non-pointer receiver,
		// so it's likely to be incorrect, but we do what we're told.
		return p.unmarshalInterface(val.Interface().(Unmarshaler), start)
	***REMOVED***

	if val.CanAddr() ***REMOVED***
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(unmarshalerType) ***REMOVED***
			return p.unmarshalInterface(pv.Interface().(Unmarshaler), start)
		***REMOVED***
	***REMOVED***

	if val.CanInterface() && val.Type().Implements(textUnmarshalerType) ***REMOVED***
		return p.unmarshalTextInterface(val.Interface().(encoding.TextUnmarshaler), start)
	***REMOVED***

	if val.CanAddr() ***REMOVED***
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(textUnmarshalerType) ***REMOVED***
			return p.unmarshalTextInterface(pv.Interface().(encoding.TextUnmarshaler), start)
		***REMOVED***
	***REMOVED***

	var (
		data         []byte
		saveData     reflect.Value
		comment      []byte
		saveComment  reflect.Value
		saveXML      reflect.Value
		saveXMLIndex int
		saveXMLData  []byte
		saveAny      reflect.Value
		sv           reflect.Value
		tinfo        *typeInfo
		err          error
	)

	switch v := val; v.Kind() ***REMOVED***
	default:
		return errors.New("unknown type " + v.Type().String())

	case reflect.Interface:
		// TODO: For now, simply ignore the field. In the near
		//       future we may choose to unmarshal the start
		//       element on it, if not nil.
		return p.Skip()

	case reflect.Slice:
		typ := v.Type()
		if typ.Elem().Kind() == reflect.Uint8 ***REMOVED***
			// []byte
			saveData = v
			break
		***REMOVED***

		// Slice of element values.
		// Grow slice.
		n := v.Len()
		if n >= v.Cap() ***REMOVED***
			ncap := 2 * n
			if ncap < 4 ***REMOVED***
				ncap = 4
			***REMOVED***
			new := reflect.MakeSlice(typ, n, ncap)
			reflect.Copy(new, v)
			v.Set(new)
		***REMOVED***
		v.SetLen(n + 1)

		// Recur to read element into slice.
		if err := p.unmarshal(v.Index(n), start); err != nil ***REMOVED***
			v.SetLen(n)
			return err
		***REMOVED***
		return nil

	case reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.String:
		saveData = v

	case reflect.Struct:
		typ := v.Type()
		if typ == nameType ***REMOVED***
			v.Set(reflect.ValueOf(start.Name))
			break
		***REMOVED***

		sv = v
		tinfo, err = getTypeInfo(typ)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Validate and assign element name.
		if tinfo.xmlname != nil ***REMOVED***
			finfo := tinfo.xmlname
			if finfo.name != "" && finfo.name != start.Name.Local ***REMOVED***
				return UnmarshalError("expected element type <" + finfo.name + "> but have <" + start.Name.Local + ">")
			***REMOVED***
			if finfo.xmlns != "" && finfo.xmlns != start.Name.Space ***REMOVED***
				e := "expected element <" + finfo.name + "> in name space " + finfo.xmlns + " but have "
				if start.Name.Space == "" ***REMOVED***
					e += "no name space"
				***REMOVED*** else ***REMOVED***
					e += start.Name.Space
				***REMOVED***
				return UnmarshalError(e)
			***REMOVED***
			fv := finfo.value(sv)
			if _, ok := fv.Interface().(Name); ok ***REMOVED***
				fv.Set(reflect.ValueOf(start.Name))
			***REMOVED***
		***REMOVED***

		// Assign attributes.
		// Also, determine whether we need to save character data or comments.
		for i := range tinfo.fields ***REMOVED***
			finfo := &tinfo.fields[i]
			switch finfo.flags & fMode ***REMOVED***
			case fAttr:
				strv := finfo.value(sv)
				// Look for attribute.
				for _, a := range start.Attr ***REMOVED***
					if a.Name.Local == finfo.name && (finfo.xmlns == "" || finfo.xmlns == a.Name.Space) ***REMOVED***
						if err := p.unmarshalAttr(strv, a); err != nil ***REMOVED***
							return err
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***

			case fCharData:
				if !saveData.IsValid() ***REMOVED***
					saveData = finfo.value(sv)
				***REMOVED***

			case fComment:
				if !saveComment.IsValid() ***REMOVED***
					saveComment = finfo.value(sv)
				***REMOVED***

			case fAny, fAny | fElement:
				if !saveAny.IsValid() ***REMOVED***
					saveAny = finfo.value(sv)
				***REMOVED***

			case fInnerXml:
				if !saveXML.IsValid() ***REMOVED***
					saveXML = finfo.value(sv)
					if p.saved == nil ***REMOVED***
						saveXMLIndex = 0
						p.saved = new(bytes.Buffer)
					***REMOVED*** else ***REMOVED***
						saveXMLIndex = p.savedOffset()
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Find end element.
	// Process sub-elements along the way.
Loop:
	for ***REMOVED***
		var savedOffset int
		if saveXML.IsValid() ***REMOVED***
			savedOffset = p.savedOffset()
		***REMOVED***
		tok, err := p.Token()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch t := tok.(type) ***REMOVED***
		case StartElement:
			consumed := false
			if sv.IsValid() ***REMOVED***
				consumed, err = p.unmarshalPath(tinfo, sv, nil, &t)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				if !consumed && saveAny.IsValid() ***REMOVED***
					consumed = true
					if err := p.unmarshal(saveAny, &t); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if !consumed ***REMOVED***
				if err := p.Skip(); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***

		case EndElement:
			if saveXML.IsValid() ***REMOVED***
				saveXMLData = p.saved.Bytes()[saveXMLIndex:savedOffset]
				if saveXMLIndex == 0 ***REMOVED***
					p.saved = nil
				***REMOVED***
			***REMOVED***
			break Loop

		case CharData:
			if saveData.IsValid() ***REMOVED***
				data = append(data, t...)
			***REMOVED***

		case Comment:
			if saveComment.IsValid() ***REMOVED***
				comment = append(comment, t...)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if saveData.IsValid() && saveData.CanInterface() && saveData.Type().Implements(textUnmarshalerType) ***REMOVED***
		if err := saveData.Interface().(encoding.TextUnmarshaler).UnmarshalText(data); err != nil ***REMOVED***
			return err
		***REMOVED***
		saveData = reflect.Value***REMOVED******REMOVED***
	***REMOVED***

	if saveData.IsValid() && saveData.CanAddr() ***REMOVED***
		pv := saveData.Addr()
		if pv.CanInterface() && pv.Type().Implements(textUnmarshalerType) ***REMOVED***
			if err := pv.Interface().(encoding.TextUnmarshaler).UnmarshalText(data); err != nil ***REMOVED***
				return err
			***REMOVED***
			saveData = reflect.Value***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	if err := copyValue(saveData, data); err != nil ***REMOVED***
		return err
	***REMOVED***

	switch t := saveComment; t.Kind() ***REMOVED***
	case reflect.String:
		t.SetString(string(comment))
	case reflect.Slice:
		t.Set(reflect.ValueOf(comment))
	***REMOVED***

	switch t := saveXML; t.Kind() ***REMOVED***
	case reflect.String:
		t.SetString(string(saveXMLData))
	case reflect.Slice:
		t.Set(reflect.ValueOf(saveXMLData))
	***REMOVED***

	return nil
***REMOVED***

func copyValue(dst reflect.Value, src []byte) (err error) ***REMOVED***
	dst0 := dst

	if dst.Kind() == reflect.Ptr ***REMOVED***
		if dst.IsNil() ***REMOVED***
			dst.Set(reflect.New(dst.Type().Elem()))
		***REMOVED***
		dst = dst.Elem()
	***REMOVED***

	// Save accumulated data.
	switch dst.Kind() ***REMOVED***
	case reflect.Invalid:
		// Probably a comment.
	default:
		return errors.New("cannot unmarshal into " + dst0.Type().String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		itmp, err := strconv.ParseInt(string(src), 10, dst.Type().Bits())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		dst.SetInt(itmp)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		utmp, err := strconv.ParseUint(string(src), 10, dst.Type().Bits())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		dst.SetUint(utmp)
	case reflect.Float32, reflect.Float64:
		ftmp, err := strconv.ParseFloat(string(src), dst.Type().Bits())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		dst.SetFloat(ftmp)
	case reflect.Bool:
		value, err := strconv.ParseBool(strings.TrimSpace(string(src)))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		dst.SetBool(value)
	case reflect.String:
		dst.SetString(string(src))
	case reflect.Slice:
		if len(src) == 0 ***REMOVED***
			// non-nil to flag presence
			src = []byte***REMOVED******REMOVED***
		***REMOVED***
		dst.SetBytes(src)
	***REMOVED***
	return nil
***REMOVED***

// unmarshalPath walks down an XML structure looking for wanted
// paths, and calls unmarshal on them.
// The consumed result tells whether XML elements have been consumed
// from the Decoder until start's matching end element, or if it's
// still untouched because start is uninteresting for sv's fields.
func (p *Decoder) unmarshalPath(tinfo *typeInfo, sv reflect.Value, parents []string, start *StartElement) (consumed bool, err error) ***REMOVED***
	recurse := false
Loop:
	for i := range tinfo.fields ***REMOVED***
		finfo := &tinfo.fields[i]
		if finfo.flags&fElement == 0 || len(finfo.parents) < len(parents) || finfo.xmlns != "" && finfo.xmlns != start.Name.Space ***REMOVED***
			continue
		***REMOVED***
		for j := range parents ***REMOVED***
			if parents[j] != finfo.parents[j] ***REMOVED***
				continue Loop
			***REMOVED***
		***REMOVED***
		if len(finfo.parents) == len(parents) && finfo.name == start.Name.Local ***REMOVED***
			// It's a perfect match, unmarshal the field.
			return true, p.unmarshal(finfo.value(sv), start)
		***REMOVED***
		if len(finfo.parents) > len(parents) && finfo.parents[len(parents)] == start.Name.Local ***REMOVED***
			// It's a prefix for the field. Break and recurse
			// since it's not ok for one field path to be itself
			// the prefix for another field path.
			recurse = true

			// We can reuse the same slice as long as we
			// don't try to append to it.
			parents = finfo.parents[:len(parents)+1]
			break
		***REMOVED***
	***REMOVED***
	if !recurse ***REMOVED***
		// We have no business with this element.
		return false, nil
	***REMOVED***
	// The element is not a perfect match for any field, but one
	// or more fields have the path to this element as a parent
	// prefix. Recurse and attempt to match these.
	for ***REMOVED***
		var tok Token
		tok, err = p.Token()
		if err != nil ***REMOVED***
			return true, err
		***REMOVED***
		switch t := tok.(type) ***REMOVED***
		case StartElement:
			consumed2, err := p.unmarshalPath(tinfo, sv, parents, &t)
			if err != nil ***REMOVED***
				return true, err
			***REMOVED***
			if !consumed2 ***REMOVED***
				if err := p.Skip(); err != nil ***REMOVED***
					return true, err
				***REMOVED***
			***REMOVED***
		case EndElement:
			return true, nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// Skip reads tokens until it has consumed the end element
// matching the most recent start element already consumed.
// It recurs if it encounters a start element, so it can be used to
// skip nested structures.
// It returns nil if it finds an end element matching the start
// element; otherwise it returns an error describing the problem.
func (d *Decoder) Skip() error ***REMOVED***
	for ***REMOVED***
		tok, err := d.Token()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch tok.(type) ***REMOVED***
		case StartElement:
			if err := d.Skip(); err != nil ***REMOVED***
				return err
			***REMOVED***
		case EndElement:
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***
