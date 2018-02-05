package toml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type tomlOpts struct ***REMOVED***
	name      string
	comment   string
	commented bool
	include   bool
	omitempty bool
***REMOVED***

type encOpts struct ***REMOVED***
	quoteMapKeys            bool
	arraysOneElementPerLine bool
***REMOVED***

var encOptsDefaults = encOpts***REMOVED***
	quoteMapKeys: false,
***REMOVED***

var timeType = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
var marshalerType = reflect.TypeOf(new(Marshaler)).Elem()

// Check if the given marshall type maps to a Tree primitive
func isPrimitive(mtype reflect.Type) bool ***REMOVED***
	switch mtype.Kind() ***REMOVED***
	case reflect.Ptr:
		return isPrimitive(mtype.Elem())
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	case reflect.Struct:
		return mtype == timeType || isCustomMarshaler(mtype)
	default:
		return false
	***REMOVED***
***REMOVED***

// Check if the given marshall type maps to a Tree slice
func isTreeSlice(mtype reflect.Type) bool ***REMOVED***
	switch mtype.Kind() ***REMOVED***
	case reflect.Slice:
		return !isOtherSlice(mtype)
	default:
		return false
	***REMOVED***
***REMOVED***

// Check if the given marshall type maps to a non-Tree slice
func isOtherSlice(mtype reflect.Type) bool ***REMOVED***
	switch mtype.Kind() ***REMOVED***
	case reflect.Ptr:
		return isOtherSlice(mtype.Elem())
	case reflect.Slice:
		return isPrimitive(mtype.Elem()) || isOtherSlice(mtype.Elem())
	default:
		return false
	***REMOVED***
***REMOVED***

// Check if the given marshall type maps to a Tree
func isTree(mtype reflect.Type) bool ***REMOVED***
	switch mtype.Kind() ***REMOVED***
	case reflect.Map:
		return true
	case reflect.Struct:
		return !isPrimitive(mtype)
	default:
		return false
	***REMOVED***
***REMOVED***

func isCustomMarshaler(mtype reflect.Type) bool ***REMOVED***
	return mtype.Implements(marshalerType)
***REMOVED***

func callCustomMarshaler(mval reflect.Value) ([]byte, error) ***REMOVED***
	return mval.Interface().(Marshaler).MarshalTOML()
***REMOVED***

// Marshaler is the interface implemented by types that
// can marshal themselves into valid TOML.
type Marshaler interface ***REMOVED***
	MarshalTOML() ([]byte, error)
***REMOVED***

/*
Marshal returns the TOML encoding of v.  Behavior is similar to the Go json
encoder, except that there is no concept of a Marshaler interface or MarshalTOML
function for sub-structs, and currently only definite types can be marshaled
(i.e. no `interface***REMOVED******REMOVED***`).

The following struct annotations are supported:

  toml:"Field"      Overrides the field's name to output.
  omitempty         When set, empty values and groups are not emitted.
  comment:"comment" Emits a # comment on the same line. This supports new lines.
  commented:"true"  Emits the value as commented.

Note that pointers are automatically assigned the "omitempty" option, as TOML
explicitly does not handle null values (saying instead the label should be
dropped).

Tree structural types and corresponding marshal types:

  *Tree                            (*)struct, (*)map[string]interface***REMOVED******REMOVED***
  []*Tree                          (*)[](*)struct, (*)[](*)map[string]interface***REMOVED******REMOVED***
  []interface***REMOVED******REMOVED*** (as interface***REMOVED******REMOVED***)   (*)[]primitive, (*)[]([]interface***REMOVED******REMOVED***)
  interface***REMOVED******REMOVED***                      (*)primitive

Tree primitive types and corresponding marshal types:

  uint64     uint, uint8-uint64, pointers to same
  int64      int, int8-uint64, pointers to same
  float64    float32, float64, pointers to same
  string     string, pointers to same
  bool       bool, pointers to same
  time.Time  time.Time***REMOVED******REMOVED***, pointers to same
*/
func Marshal(v interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return NewEncoder(nil).marshal(v)
***REMOVED***

// Encoder writes TOML values to an output stream.
type Encoder struct ***REMOVED***
	w io.Writer
	encOpts
***REMOVED***

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder ***REMOVED***
	return &Encoder***REMOVED***
		w:       w,
		encOpts: encOptsDefaults,
	***REMOVED***
***REMOVED***

// Encode writes the TOML encoding of v to the stream.
//
// See the documentation for Marshal for details.
func (e *Encoder) Encode(v interface***REMOVED******REMOVED***) error ***REMOVED***
	b, err := e.marshal(v)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := e.w.Write(b); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// QuoteMapKeys sets up the encoder to encode
// maps with string type keys with quoted TOML keys.
//
// This relieves the character limitations on map keys.
func (e *Encoder) QuoteMapKeys(v bool) *Encoder ***REMOVED***
	e.quoteMapKeys = v
	return e
***REMOVED***

// ArraysWithOneElementPerLine sets up the encoder to encode arrays
// with more than one element on multiple lines instead of one.
//
// For example:
//
//   A = [1,2,3]
//
// Becomes
//
//   A = [
//     1,
//     2,
//     3
//   ]
func (e *Encoder) ArraysWithOneElementPerLine(v bool) *Encoder ***REMOVED***
	e.arraysOneElementPerLine = v
	return e
***REMOVED***

func (e *Encoder) marshal(v interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	mtype := reflect.TypeOf(v)
	if mtype.Kind() != reflect.Struct ***REMOVED***
		return []byte***REMOVED******REMOVED***, errors.New("Only a struct can be marshaled to TOML")
	***REMOVED***
	sval := reflect.ValueOf(v)
	if isCustomMarshaler(mtype) ***REMOVED***
		return callCustomMarshaler(sval)
	***REMOVED***
	t, err := e.valueToTree(mtype, sval)
	if err != nil ***REMOVED***
		return []byte***REMOVED******REMOVED***, err
	***REMOVED***

	var buf bytes.Buffer
	_, err = t.writeTo(&buf, "", "", 0, e.arraysOneElementPerLine)

	return buf.Bytes(), err
***REMOVED***

// Convert given marshal struct or map value to toml tree
func (e *Encoder) valueToTree(mtype reflect.Type, mval reflect.Value) (*Tree, error) ***REMOVED***
	if mtype.Kind() == reflect.Ptr ***REMOVED***
		return e.valueToTree(mtype.Elem(), mval.Elem())
	***REMOVED***
	tval := newTree()
	switch mtype.Kind() ***REMOVED***
	case reflect.Struct:
		for i := 0; i < mtype.NumField(); i++ ***REMOVED***
			mtypef, mvalf := mtype.Field(i), mval.Field(i)
			opts := tomlOptions(mtypef)
			if opts.include && (!opts.omitempty || !isZero(mvalf)) ***REMOVED***
				val, err := e.valueToToml(mtypef.Type, mvalf)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				tval.SetWithComment(opts.name, opts.comment, opts.commented, val)
			***REMOVED***
		***REMOVED***
	case reflect.Map:
		for _, key := range mval.MapKeys() ***REMOVED***
			mvalf := mval.MapIndex(key)
			val, err := e.valueToToml(mtype.Elem(), mvalf)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if e.quoteMapKeys ***REMOVED***
				keyStr, err := tomlValueStringRepresentation(key.String(), "", e.arraysOneElementPerLine)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				tval.SetPath([]string***REMOVED***keyStr***REMOVED***, val)
			***REMOVED*** else ***REMOVED***
				tval.Set(key.String(), val)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return tval, nil
***REMOVED***

// Convert given marshal slice to slice of Toml trees
func (e *Encoder) valueToTreeSlice(mtype reflect.Type, mval reflect.Value) ([]*Tree, error) ***REMOVED***
	tval := make([]*Tree, mval.Len(), mval.Len())
	for i := 0; i < mval.Len(); i++ ***REMOVED***
		val, err := e.valueToTree(mtype.Elem(), mval.Index(i))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		tval[i] = val
	***REMOVED***
	return tval, nil
***REMOVED***

// Convert given marshal slice to slice of toml values
func (e *Encoder) valueToOtherSlice(mtype reflect.Type, mval reflect.Value) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	tval := make([]interface***REMOVED******REMOVED***, mval.Len(), mval.Len())
	for i := 0; i < mval.Len(); i++ ***REMOVED***
		val, err := e.valueToToml(mtype.Elem(), mval.Index(i))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		tval[i] = val
	***REMOVED***
	return tval, nil
***REMOVED***

// Convert given marshal value to toml value
func (e *Encoder) valueToToml(mtype reflect.Type, mval reflect.Value) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if mtype.Kind() == reflect.Ptr ***REMOVED***
		return e.valueToToml(mtype.Elem(), mval.Elem())
	***REMOVED***
	switch ***REMOVED***
	case isCustomMarshaler(mtype):
		return callCustomMarshaler(mval)
	case isTree(mtype):
		return e.valueToTree(mtype, mval)
	case isTreeSlice(mtype):
		return e.valueToTreeSlice(mtype, mval)
	case isOtherSlice(mtype):
		return e.valueToOtherSlice(mtype, mval)
	default:
		switch mtype.Kind() ***REMOVED***
		case reflect.Bool:
			return mval.Bool(), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return mval.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return mval.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return mval.Float(), nil
		case reflect.String:
			return mval.String(), nil
		case reflect.Struct:
			return mval.Interface().(time.Time), nil
		default:
			return nil, fmt.Errorf("Marshal can't handle %v(%v)", mtype, mtype.Kind())
		***REMOVED***
	***REMOVED***
***REMOVED***

// Unmarshal attempts to unmarshal the Tree into a Go struct pointed by v.
// Neither Unmarshaler interfaces nor UnmarshalTOML functions are supported for
// sub-structs, and only definite types can be unmarshaled.
func (t *Tree) Unmarshal(v interface***REMOVED******REMOVED***) error ***REMOVED***
	d := Decoder***REMOVED***tval: t***REMOVED***
	return d.unmarshal(v)
***REMOVED***

// Marshal returns the TOML encoding of Tree.
// See Marshal() documentation for types mapping table.
func (t *Tree) Marshal() ([]byte, error) ***REMOVED***
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(t)
	return buf.Bytes(), err
***REMOVED***

// Unmarshal parses the TOML-encoded data and stores the result in the value
// pointed to by v. Behavior is similar to the Go json encoder, except that there
// is no concept of an Unmarshaler interface or UnmarshalTOML function for
// sub-structs, and currently only definite types can be unmarshaled to (i.e. no
// `interface***REMOVED******REMOVED***`).
//
// The following struct annotations are supported:
//
//   toml:"Field" Overrides the field's name to map to.
//
// See Marshal() documentation for types mapping table.
func Unmarshal(data []byte, v interface***REMOVED******REMOVED***) error ***REMOVED***
	t, err := LoadReader(bytes.NewReader(data))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return t.Unmarshal(v)
***REMOVED***

// Decoder reads and decodes TOML values from an input stream.
type Decoder struct ***REMOVED***
	r    io.Reader
	tval *Tree
	encOpts
***REMOVED***

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder ***REMOVED***
	return &Decoder***REMOVED***
		r:       r,
		encOpts: encOptsDefaults,
	***REMOVED***
***REMOVED***

// Decode reads a TOML-encoded value from it's input
// and unmarshals it in the value pointed at by v.
//
// See the documentation for Marshal for details.
func (d *Decoder) Decode(v interface***REMOVED******REMOVED***) error ***REMOVED***
	var err error
	d.tval, err = LoadReader(d.r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return d.unmarshal(v)
***REMOVED***

func (d *Decoder) unmarshal(v interface***REMOVED******REMOVED***) error ***REMOVED***
	mtype := reflect.TypeOf(v)
	if mtype.Kind() != reflect.Ptr || mtype.Elem().Kind() != reflect.Struct ***REMOVED***
		return errors.New("Only a pointer to struct can be unmarshaled from TOML")
	***REMOVED***

	sval, err := d.valueFromTree(mtype.Elem(), d.tval)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	reflect.ValueOf(v).Elem().Set(sval)
	return nil
***REMOVED***

// Convert toml tree to marshal struct or map, using marshal type
func (d *Decoder) valueFromTree(mtype reflect.Type, tval *Tree) (reflect.Value, error) ***REMOVED***
	if mtype.Kind() == reflect.Ptr ***REMOVED***
		return d.unwrapPointer(mtype, tval)
	***REMOVED***
	var mval reflect.Value
	switch mtype.Kind() ***REMOVED***
	case reflect.Struct:
		mval = reflect.New(mtype).Elem()
		for i := 0; i < mtype.NumField(); i++ ***REMOVED***
			mtypef := mtype.Field(i)
			opts := tomlOptions(mtypef)
			if opts.include ***REMOVED***
				baseKey := opts.name
				keysToTry := []string***REMOVED***baseKey, strings.ToLower(baseKey), strings.ToTitle(baseKey)***REMOVED***
				for _, key := range keysToTry ***REMOVED***
					exists := tval.Has(key)
					if !exists ***REMOVED***
						continue
					***REMOVED***
					val := tval.Get(key)
					mvalf, err := d.valueFromToml(mtypef.Type, val)
					if err != nil ***REMOVED***
						return mval, formatError(err, tval.GetPosition(key))
					***REMOVED***
					mval.Field(i).Set(mvalf)
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	case reflect.Map:
		mval = reflect.MakeMap(mtype)
		for _, key := range tval.Keys() ***REMOVED***
			// TODO: path splits key
			val := tval.GetPath([]string***REMOVED***key***REMOVED***)
			mvalf, err := d.valueFromToml(mtype.Elem(), val)
			if err != nil ***REMOVED***
				return mval, formatError(err, tval.GetPosition(key))
			***REMOVED***
			mval.SetMapIndex(reflect.ValueOf(key), mvalf)
		***REMOVED***
	***REMOVED***
	return mval, nil
***REMOVED***

// Convert toml value to marshal struct/map slice, using marshal type
func (d *Decoder) valueFromTreeSlice(mtype reflect.Type, tval []*Tree) (reflect.Value, error) ***REMOVED***
	mval := reflect.MakeSlice(mtype, len(tval), len(tval))
	for i := 0; i < len(tval); i++ ***REMOVED***
		val, err := d.valueFromTree(mtype.Elem(), tval[i])
		if err != nil ***REMOVED***
			return mval, err
		***REMOVED***
		mval.Index(i).Set(val)
	***REMOVED***
	return mval, nil
***REMOVED***

// Convert toml value to marshal primitive slice, using marshal type
func (d *Decoder) valueFromOtherSlice(mtype reflect.Type, tval []interface***REMOVED******REMOVED***) (reflect.Value, error) ***REMOVED***
	mval := reflect.MakeSlice(mtype, len(tval), len(tval))
	for i := 0; i < len(tval); i++ ***REMOVED***
		val, err := d.valueFromToml(mtype.Elem(), tval[i])
		if err != nil ***REMOVED***
			return mval, err
		***REMOVED***
		mval.Index(i).Set(val)
	***REMOVED***
	return mval, nil
***REMOVED***

// Convert toml value to marshal value, using marshal type
func (d *Decoder) valueFromToml(mtype reflect.Type, tval interface***REMOVED******REMOVED***) (reflect.Value, error) ***REMOVED***
	if mtype.Kind() == reflect.Ptr ***REMOVED***
		return d.unwrapPointer(mtype, tval)
	***REMOVED***

	switch tval.(type) ***REMOVED***
	case *Tree:
		if isTree(mtype) ***REMOVED***
			return d.valueFromTree(mtype, tval.(*Tree))
		***REMOVED***
		return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to a tree", tval, tval)
	case []*Tree:
		if isTreeSlice(mtype) ***REMOVED***
			return d.valueFromTreeSlice(mtype, tval.([]*Tree))
		***REMOVED***
		return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to trees", tval, tval)
	case []interface***REMOVED******REMOVED***:
		if isOtherSlice(mtype) ***REMOVED***
			return d.valueFromOtherSlice(mtype, tval.([]interface***REMOVED******REMOVED***))
		***REMOVED***
		return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to a slice", tval, tval)
	default:
		switch mtype.Kind() ***REMOVED***
		case reflect.Bool, reflect.Struct:
			val := reflect.ValueOf(tval)
			// if this passes for when mtype is reflect.Struct, tval is a time.Time
			if !val.Type().ConvertibleTo(mtype) ***REMOVED***
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			***REMOVED***

			return val.Convert(mtype), nil
		case reflect.String:
			val := reflect.ValueOf(tval)
			// stupidly, int64 is convertible to string. So special case this.
			if !val.Type().ConvertibleTo(mtype) || val.Kind() == reflect.Int64 ***REMOVED***
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			***REMOVED***

			return val.Convert(mtype), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val := reflect.ValueOf(tval)
			if !val.Type().ConvertibleTo(mtype) ***REMOVED***
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			***REMOVED***
			if reflect.Indirect(reflect.New(mtype)).OverflowInt(val.Int()) ***REMOVED***
				return reflect.ValueOf(nil), fmt.Errorf("%v(%T) would overflow %v", tval, tval, mtype.String())
			***REMOVED***

			return val.Convert(mtype), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			val := reflect.ValueOf(tval)
			if !val.Type().ConvertibleTo(mtype) ***REMOVED***
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			***REMOVED***
			if val.Int() < 0 ***REMOVED***
				return reflect.ValueOf(nil), fmt.Errorf("%v(%T) is negative so does not fit in %v", tval, tval, mtype.String())
			***REMOVED***
			if reflect.Indirect(reflect.New(mtype)).OverflowUint(uint64(val.Int())) ***REMOVED***
				return reflect.ValueOf(nil), fmt.Errorf("%v(%T) would overflow %v", tval, tval, mtype.String())
			***REMOVED***

			return val.Convert(mtype), nil
		case reflect.Float32, reflect.Float64:
			val := reflect.ValueOf(tval)
			if !val.Type().ConvertibleTo(mtype) ***REMOVED***
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			***REMOVED***
			if reflect.Indirect(reflect.New(mtype)).OverflowFloat(val.Float()) ***REMOVED***
				return reflect.ValueOf(nil), fmt.Errorf("%v(%T) would overflow %v", tval, tval, mtype.String())
			***REMOVED***

			return val.Convert(mtype), nil
		default:
			return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v(%v)", tval, tval, mtype, mtype.Kind())
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *Decoder) unwrapPointer(mtype reflect.Type, tval interface***REMOVED******REMOVED***) (reflect.Value, error) ***REMOVED***
	val, err := d.valueFromToml(mtype.Elem(), tval)
	if err != nil ***REMOVED***
		return reflect.ValueOf(nil), err
	***REMOVED***
	mval := reflect.New(mtype.Elem())
	mval.Elem().Set(val)
	return mval, nil
***REMOVED***

func tomlOptions(vf reflect.StructField) tomlOpts ***REMOVED***
	tag := vf.Tag.Get("toml")
	parse := strings.Split(tag, ",")
	var comment string
	if c := vf.Tag.Get("comment"); c != "" ***REMOVED***
		comment = c
	***REMOVED***
	commented, _ := strconv.ParseBool(vf.Tag.Get("commented"))
	result := tomlOpts***REMOVED***name: vf.Name, comment: comment, commented: commented, include: true, omitempty: false***REMOVED***
	if parse[0] != "" ***REMOVED***
		if parse[0] == "-" && len(parse) == 1 ***REMOVED***
			result.include = false
		***REMOVED*** else ***REMOVED***
			result.name = strings.Trim(parse[0], " ")
		***REMOVED***
	***REMOVED***
	if vf.PkgPath != "" ***REMOVED***
		result.include = false
	***REMOVED***
	if len(parse) > 1 && strings.Trim(parse[1], " ") == "omitempty" ***REMOVED***
		result.omitempty = true
	***REMOVED***
	if vf.Type.Kind() == reflect.Ptr ***REMOVED***
		result.omitempty = true
	***REMOVED***
	return result
***REMOVED***

func isZero(val reflect.Value) bool ***REMOVED***
	switch val.Type().Kind() ***REMOVED***
	case reflect.Map:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return val.Len() == 0
	default:
		return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	***REMOVED***
***REMOVED***

func formatError(err error, pos Position) error ***REMOVED***
	if err.Error()[0] == '(' ***REMOVED*** // Error already contains position information
		return err
	***REMOVED***
	return fmt.Errorf("%s: %s", pos, err)
***REMOVED***
