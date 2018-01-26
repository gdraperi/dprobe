package xmlutil

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// UnmarshalXML deserializes an xml.Decoder into the container v. V
// needs to match the shape of the XML expected to be decoded.
// If the shape doesn't match unmarshaling will fail.
func UnmarshalXML(v interface***REMOVED******REMOVED***, d *xml.Decoder, wrapper string) error ***REMOVED***
	n, err := XMLToStruct(d, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n.Children != nil ***REMOVED***
		for _, root := range n.Children ***REMOVED***
			for _, c := range root ***REMOVED***
				if wrappedChild, ok := c.Children[wrapper]; ok ***REMOVED***
					c = wrappedChild[0] // pull out wrapped element
				***REMOVED***

				err = parse(reflect.ValueOf(v), c, "")
				if err != nil ***REMOVED***
					if err == io.EOF ***REMOVED***
						return nil
					***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	return nil
***REMOVED***

// parse deserializes any value from the XMLNode. The type tag is used to infer the type, or reflect
// will be used to determine the type from r.
func parse(r reflect.Value, node *XMLNode, tag reflect.StructTag) error ***REMOVED***
	rtype := r.Type()
	if rtype.Kind() == reflect.Ptr ***REMOVED***
		rtype = rtype.Elem() // check kind of actual element type
	***REMOVED***

	t := tag.Get("type")
	if t == "" ***REMOVED***
		switch rtype.Kind() ***REMOVED***
		case reflect.Struct:
			t = "structure"
		case reflect.Slice:
			t = "list"
		case reflect.Map:
			t = "map"
		***REMOVED***
	***REMOVED***

	switch t ***REMOVED***
	case "structure":
		if field, ok := rtype.FieldByName("_"); ok ***REMOVED***
			tag = field.Tag
		***REMOVED***
		return parseStruct(r, node, tag)
	case "list":
		return parseList(r, node, tag)
	case "map":
		return parseMap(r, node, tag)
	default:
		return parseScalar(r, node, tag)
	***REMOVED***
***REMOVED***

// parseStruct deserializes a structure and its fields from an XMLNode. Any nested
// types in the structure will also be deserialized.
func parseStruct(r reflect.Value, node *XMLNode, tag reflect.StructTag) error ***REMOVED***
	t := r.Type()
	if r.Kind() == reflect.Ptr ***REMOVED***
		if r.IsNil() ***REMOVED*** // create the structure if it's nil
			s := reflect.New(r.Type().Elem())
			r.Set(s)
			r = s
		***REMOVED***

		r = r.Elem()
		t = t.Elem()
	***REMOVED***

	// unwrap any payloads
	if payload := tag.Get("payload"); payload != "" ***REMOVED***
		field, _ := t.FieldByName(payload)
		return parseStruct(r.FieldByName(payload), node, field.Tag)
	***REMOVED***

	for i := 0; i < t.NumField(); i++ ***REMOVED***
		field := t.Field(i)
		if c := field.Name[0:1]; strings.ToLower(c) == c ***REMOVED***
			continue // ignore unexported fields
		***REMOVED***

		// figure out what this field is called
		name := field.Name
		if field.Tag.Get("flattened") != "" && field.Tag.Get("locationNameList") != "" ***REMOVED***
			name = field.Tag.Get("locationNameList")
		***REMOVED*** else if locName := field.Tag.Get("locationName"); locName != "" ***REMOVED***
			name = locName
		***REMOVED***

		// try to find the field by name in elements
		elems := node.Children[name]

		if elems == nil ***REMOVED*** // try to find the field in attributes
			if val, ok := node.findElem(name); ok ***REMOVED***
				elems = []*XMLNode***REMOVED******REMOVED***Text: val***REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***

		member := r.FieldByName(field.Name)
		for _, elem := range elems ***REMOVED***
			err := parse(member, elem, field.Tag)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// parseList deserializes a list of values from an XML node. Each list entry
// will also be deserialized.
func parseList(r reflect.Value, node *XMLNode, tag reflect.StructTag) error ***REMOVED***
	t := r.Type()

	if tag.Get("flattened") == "" ***REMOVED*** // look at all item entries
		mname := "member"
		if name := tag.Get("locationNameList"); name != "" ***REMOVED***
			mname = name
		***REMOVED***

		if Children, ok := node.Children[mname]; ok ***REMOVED***
			if r.IsNil() ***REMOVED***
				r.Set(reflect.MakeSlice(t, len(Children), len(Children)))
			***REMOVED***

			for i, c := range Children ***REMOVED***
				err := parse(r.Index(i), c, "")
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED*** // flattened list means this is a single element
		if r.IsNil() ***REMOVED***
			r.Set(reflect.MakeSlice(t, 0, 0))
		***REMOVED***

		childR := reflect.Zero(t.Elem())
		r.Set(reflect.Append(r, childR))
		err := parse(r.Index(r.Len()-1), node, "")
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// parseMap deserializes a map from an XMLNode. The direct children of the XMLNode
// will also be deserialized as map entries.
func parseMap(r reflect.Value, node *XMLNode, tag reflect.StructTag) error ***REMOVED***
	if r.IsNil() ***REMOVED***
		r.Set(reflect.MakeMap(r.Type()))
	***REMOVED***

	if tag.Get("flattened") == "" ***REMOVED*** // look at all child entries
		for _, entry := range node.Children["entry"] ***REMOVED***
			parseMapEntry(r, entry, tag)
		***REMOVED***
	***REMOVED*** else ***REMOVED*** // this element is itself an entry
		parseMapEntry(r, node, tag)
	***REMOVED***

	return nil
***REMOVED***

// parseMapEntry deserializes a map entry from a XML node.
func parseMapEntry(r reflect.Value, node *XMLNode, tag reflect.StructTag) error ***REMOVED***
	kname, vname := "key", "value"
	if n := tag.Get("locationNameKey"); n != "" ***REMOVED***
		kname = n
	***REMOVED***
	if n := tag.Get("locationNameValue"); n != "" ***REMOVED***
		vname = n
	***REMOVED***

	keys, ok := node.Children[kname]
	values := node.Children[vname]
	if ok ***REMOVED***
		for i, key := range keys ***REMOVED***
			keyR := reflect.ValueOf(key.Text)
			value := values[i]
			valueR := reflect.New(r.Type().Elem()).Elem()

			parse(valueR, value, "")
			r.SetMapIndex(keyR, valueR)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// parseScaller deserializes an XMLNode value into a concrete type based on the
// interface type of r.
//
// Error is returned if the deserialization fails due to invalid type conversion,
// or unsupported interface type.
func parseScalar(r reflect.Value, node *XMLNode, tag reflect.StructTag) error ***REMOVED***
	switch r.Interface().(type) ***REMOVED***
	case *string:
		r.Set(reflect.ValueOf(&node.Text))
		return nil
	case []byte:
		b, err := base64.StdEncoding.DecodeString(node.Text)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Set(reflect.ValueOf(b))
	case *bool:
		v, err := strconv.ParseBool(node.Text)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Set(reflect.ValueOf(&v))
	case *int64:
		v, err := strconv.ParseInt(node.Text, 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Set(reflect.ValueOf(&v))
	case *float64:
		v, err := strconv.ParseFloat(node.Text, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Set(reflect.ValueOf(&v))
	case *time.Time:
		const ISO8601UTC = "2006-01-02T15:04:05Z"
		t, err := time.Parse(ISO8601UTC, node.Text)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Set(reflect.ValueOf(&t))
	default:
		return fmt.Errorf("unsupported value: %v (%s)", r.Interface(), r.Type())
	***REMOVED***
	return nil
***REMOVED***
