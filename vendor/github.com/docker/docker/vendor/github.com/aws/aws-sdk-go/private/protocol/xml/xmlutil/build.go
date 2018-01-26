// Package xmlutil provides XML serialization of AWS requests and responses.
package xmlutil

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/private/protocol"
)

// BuildXML will serialize params into an xml.Encoder.
// Error will be returned if the serialization of any of the params or nested values fails.
func BuildXML(params interface***REMOVED******REMOVED***, e *xml.Encoder) error ***REMOVED***
	b := xmlBuilder***REMOVED***encoder: e, namespaces: map[string]string***REMOVED******REMOVED******REMOVED***
	root := NewXMLElement(xml.Name***REMOVED******REMOVED***)
	if err := b.buildValue(reflect.ValueOf(params), root, ""); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, c := range root.Children ***REMOVED***
		for _, v := range c ***REMOVED***
			return StructToXML(e, v, false)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Returns the reflection element of a value, if it is a pointer.
func elemOf(value reflect.Value) reflect.Value ***REMOVED***
	for value.Kind() == reflect.Ptr ***REMOVED***
		value = value.Elem()
	***REMOVED***
	return value
***REMOVED***

// A xmlBuilder serializes values from Go code to XML
type xmlBuilder struct ***REMOVED***
	encoder    *xml.Encoder
	namespaces map[string]string
***REMOVED***

// buildValue generic XMLNode builder for any type. Will build value for their specific type
// struct, list, map, scalar.
//
// Also takes a "type" tag value to set what type a value should be converted to XMLNode as. If
// type is not provided reflect will be used to determine the value's type.
func (b *xmlBuilder) buildValue(value reflect.Value, current *XMLNode, tag reflect.StructTag) error ***REMOVED***
	value = elemOf(value)
	if !value.IsValid() ***REMOVED*** // no need to handle zero values
		return nil
	***REMOVED*** else if tag.Get("location") != "" ***REMOVED*** // don't handle non-body location values
		return nil
	***REMOVED***

	t := tag.Get("type")
	if t == "" ***REMOVED***
		switch value.Kind() ***REMOVED***
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
		if field, ok := value.Type().FieldByName("_"); ok ***REMOVED***
			tag = tag + reflect.StructTag(" ") + field.Tag
		***REMOVED***
		return b.buildStruct(value, current, tag)
	case "list":
		return b.buildList(value, current, tag)
	case "map":
		return b.buildMap(value, current, tag)
	default:
		return b.buildScalar(value, current, tag)
	***REMOVED***
***REMOVED***

// buildStruct adds a struct and its fields to the current XMLNode. All fields any any nested
// types are converted to XMLNodes also.
func (b *xmlBuilder) buildStruct(value reflect.Value, current *XMLNode, tag reflect.StructTag) error ***REMOVED***
	if !value.IsValid() ***REMOVED***
		return nil
	***REMOVED***

	fieldAdded := false

	// unwrap payloads
	if payload := tag.Get("payload"); payload != "" ***REMOVED***
		field, _ := value.Type().FieldByName(payload)
		tag = field.Tag
		value = elemOf(value.FieldByName(payload))

		if !value.IsValid() ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	child := NewXMLElement(xml.Name***REMOVED***Local: tag.Get("locationName")***REMOVED***)

	// there is an xmlNamespace associated with this struct
	if prefix, uri := tag.Get("xmlPrefix"), tag.Get("xmlURI"); uri != "" ***REMOVED***
		ns := xml.Attr***REMOVED***
			Name:  xml.Name***REMOVED***Local: "xmlns"***REMOVED***,
			Value: uri,
		***REMOVED***
		if prefix != "" ***REMOVED***
			b.namespaces[prefix] = uri // register the namespace
			ns.Name.Local = "xmlns:" + prefix
		***REMOVED***

		child.Attr = append(child.Attr, ns)
	***REMOVED***

	t := value.Type()
	for i := 0; i < value.NumField(); i++ ***REMOVED***
		member := elemOf(value.Field(i))
		field := t.Field(i)

		if field.PkgPath != "" ***REMOVED***
			continue // ignore unexported fields
		***REMOVED***
		if field.Tag.Get("ignore") != "" ***REMOVED***
			continue
		***REMOVED***

		mTag := field.Tag
		if mTag.Get("location") != "" ***REMOVED*** // skip non-body members
			continue
		***REMOVED***

		if protocol.CanSetIdempotencyToken(value.Field(i), field) ***REMOVED***
			token := protocol.GetIdempotencyToken()
			member = reflect.ValueOf(token)
		***REMOVED***

		memberName := mTag.Get("locationName")
		if memberName == "" ***REMOVED***
			memberName = field.Name
			mTag = reflect.StructTag(string(mTag) + ` locationName:"` + memberName + `"`)
		***REMOVED***
		if err := b.buildValue(member, child, mTag); err != nil ***REMOVED***
			return err
		***REMOVED***

		fieldAdded = true
	***REMOVED***

	if fieldAdded ***REMOVED*** // only append this child if we have one ore more valid members
		current.AddChild(child)
	***REMOVED***

	return nil
***REMOVED***

// buildList adds the value's list items to the current XMLNode as children nodes. All
// nested values in the list are converted to XMLNodes also.
func (b *xmlBuilder) buildList(value reflect.Value, current *XMLNode, tag reflect.StructTag) error ***REMOVED***
	if value.IsNil() ***REMOVED*** // don't build omitted lists
		return nil
	***REMOVED***

	// check for unflattened list member
	flattened := tag.Get("flattened") != ""

	xname := xml.Name***REMOVED***Local: tag.Get("locationName")***REMOVED***
	if flattened ***REMOVED***
		for i := 0; i < value.Len(); i++ ***REMOVED***
			child := NewXMLElement(xname)
			current.AddChild(child)
			if err := b.buildValue(value.Index(i), child, ""); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		list := NewXMLElement(xname)
		current.AddChild(list)

		for i := 0; i < value.Len(); i++ ***REMOVED***
			iname := tag.Get("locationNameList")
			if iname == "" ***REMOVED***
				iname = "member"
			***REMOVED***

			child := NewXMLElement(xml.Name***REMOVED***Local: iname***REMOVED***)
			list.AddChild(child)
			if err := b.buildValue(value.Index(i), child, ""); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// buildMap adds the value's key/value pairs to the current XMLNode as children nodes. All
// nested values in the map are converted to XMLNodes also.
//
// Error will be returned if it is unable to build the map's values into XMLNodes
func (b *xmlBuilder) buildMap(value reflect.Value, current *XMLNode, tag reflect.StructTag) error ***REMOVED***
	if value.IsNil() ***REMOVED*** // don't build omitted maps
		return nil
	***REMOVED***

	maproot := NewXMLElement(xml.Name***REMOVED***Local: tag.Get("locationName")***REMOVED***)
	current.AddChild(maproot)
	current = maproot

	kname, vname := "key", "value"
	if n := tag.Get("locationNameKey"); n != "" ***REMOVED***
		kname = n
	***REMOVED***
	if n := tag.Get("locationNameValue"); n != "" ***REMOVED***
		vname = n
	***REMOVED***

	// sorting is not required for compliance, but it makes testing easier
	keys := make([]string, value.Len())
	for i, k := range value.MapKeys() ***REMOVED***
		keys[i] = k.String()
	***REMOVED***
	sort.Strings(keys)

	for _, k := range keys ***REMOVED***
		v := value.MapIndex(reflect.ValueOf(k))

		mapcur := current
		if tag.Get("flattened") == "" ***REMOVED*** // add "entry" tag to non-flat maps
			child := NewXMLElement(xml.Name***REMOVED***Local: "entry"***REMOVED***)
			mapcur.AddChild(child)
			mapcur = child
		***REMOVED***

		kchild := NewXMLElement(xml.Name***REMOVED***Local: kname***REMOVED***)
		kchild.Text = k
		vchild := NewXMLElement(xml.Name***REMOVED***Local: vname***REMOVED***)
		mapcur.AddChild(kchild)
		mapcur.AddChild(vchild)

		if err := b.buildValue(v, vchild, ""); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// buildScalar will convert the value into a string and append it as a attribute or child
// of the current XMLNode.
//
// The value will be added as an attribute if tag contains a "xmlAttribute" attribute value.
//
// Error will be returned if the value type is unsupported.
func (b *xmlBuilder) buildScalar(value reflect.Value, current *XMLNode, tag reflect.StructTag) error ***REMOVED***
	var str string
	switch converted := value.Interface().(type) ***REMOVED***
	case string:
		str = converted
	case []byte:
		if !value.IsNil() ***REMOVED***
			str = base64.StdEncoding.EncodeToString(converted)
		***REMOVED***
	case bool:
		str = strconv.FormatBool(converted)
	case int64:
		str = strconv.FormatInt(converted, 10)
	case int:
		str = strconv.Itoa(converted)
	case float64:
		str = strconv.FormatFloat(converted, 'f', -1, 64)
	case float32:
		str = strconv.FormatFloat(float64(converted), 'f', -1, 32)
	case time.Time:
		const ISO8601UTC = "2006-01-02T15:04:05Z"
		str = converted.UTC().Format(ISO8601UTC)
	default:
		return fmt.Errorf("unsupported value for param %s: %v (%s)",
			tag.Get("locationName"), value.Interface(), value.Type().Name())
	***REMOVED***

	xname := xml.Name***REMOVED***Local: tag.Get("locationName")***REMOVED***
	if tag.Get("xmlAttribute") != "" ***REMOVED*** // put into current node's attribute list
		attr := xml.Attr***REMOVED***Name: xname, Value: str***REMOVED***
		current.Attr = append(current.Attr, attr)
	***REMOVED*** else ***REMOVED*** // regular text node
		current.AddChild(&XMLNode***REMOVED***Name: xname, Text: str***REMOVED***)
	***REMOVED***
	return nil
***REMOVED***
