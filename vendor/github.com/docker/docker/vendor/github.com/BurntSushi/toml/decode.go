package toml

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"strings"
	"time"
)

var e = fmt.Errorf

// Unmarshaler is the interface implemented by objects that can unmarshal a
// TOML description of themselves.
type Unmarshaler interface ***REMOVED***
	UnmarshalTOML(interface***REMOVED******REMOVED***) error
***REMOVED***

// Unmarshal decodes the contents of `p` in TOML format into a pointer `v`.
func Unmarshal(p []byte, v interface***REMOVED******REMOVED***) error ***REMOVED***
	_, err := Decode(string(p), v)
	return err
***REMOVED***

// Primitive is a TOML value that hasn't been decoded into a Go value.
// When using the various `Decode*` functions, the type `Primitive` may
// be given to any value, and its decoding will be delayed.
//
// A `Primitive` value can be decoded using the `PrimitiveDecode` function.
//
// The underlying representation of a `Primitive` value is subject to change.
// Do not rely on it.
//
// N.B. Primitive values are still parsed, so using them will only avoid
// the overhead of reflection. They can be useful when you don't know the
// exact type of TOML data until run time.
type Primitive struct ***REMOVED***
	undecoded interface***REMOVED******REMOVED***
	context   Key
***REMOVED***

// DEPRECATED!
//
// Use MetaData.PrimitiveDecode instead.
func PrimitiveDecode(primValue Primitive, v interface***REMOVED******REMOVED***) error ***REMOVED***
	md := MetaData***REMOVED***decoded: make(map[string]bool)***REMOVED***
	return md.unify(primValue.undecoded, rvalue(v))
***REMOVED***

// PrimitiveDecode is just like the other `Decode*` functions, except it
// decodes a TOML value that has already been parsed. Valid primitive values
// can *only* be obtained from values filled by the decoder functions,
// including this method. (i.e., `v` may contain more `Primitive`
// values.)
//
// Meta data for primitive values is included in the meta data returned by
// the `Decode*` functions with one exception: keys returned by the Undecoded
// method will only reflect keys that were decoded. Namely, any keys hidden
// behind a Primitive will be considered undecoded. Executing this method will
// update the undecoded keys in the meta data. (See the example.)
func (md *MetaData) PrimitiveDecode(primValue Primitive, v interface***REMOVED******REMOVED***) error ***REMOVED***
	md.context = primValue.context
	defer func() ***REMOVED*** md.context = nil ***REMOVED***()
	return md.unify(primValue.undecoded, rvalue(v))
***REMOVED***

// Decode will decode the contents of `data` in TOML format into a pointer
// `v`.
//
// TOML hashes correspond to Go structs or maps. (Dealer's choice. They can be
// used interchangeably.)
//
// TOML arrays of tables correspond to either a slice of structs or a slice
// of maps.
//
// TOML datetimes correspond to Go `time.Time` values.
//
// All other TOML types (float, string, int, bool and array) correspond
// to the obvious Go types.
//
// An exception to the above rules is if a type implements the
// encoding.TextUnmarshaler interface. In this case, any primitive TOML value
// (floats, strings, integers, booleans and datetimes) will be converted to
// a byte string and given to the value's UnmarshalText method. See the
// Unmarshaler example for a demonstration with time duration strings.
//
// Key mapping
//
// TOML keys can map to either keys in a Go map or field names in a Go
// struct. The special `toml` struct tag may be used to map TOML keys to
// struct fields that don't match the key name exactly. (See the example.)
// A case insensitive match to struct names will be tried if an exact match
// can't be found.
//
// The mapping between TOML values and Go values is loose. That is, there
// may exist TOML values that cannot be placed into your representation, and
// there may be parts of your representation that do not correspond to
// TOML values. This loose mapping can be made stricter by using the IsDefined
// and/or Undecoded methods on the MetaData returned.
//
// This decoder will not handle cyclic types. If a cyclic type is passed,
// `Decode` will not terminate.
func Decode(data string, v interface***REMOVED******REMOVED***) (MetaData, error) ***REMOVED***
	p, err := parse(data)
	if err != nil ***REMOVED***
		return MetaData***REMOVED******REMOVED***, err
	***REMOVED***
	md := MetaData***REMOVED***
		p.mapping, p.types, p.ordered,
		make(map[string]bool, len(p.ordered)), nil,
	***REMOVED***
	return md, md.unify(p.mapping, rvalue(v))
***REMOVED***

// DecodeFile is just like Decode, except it will automatically read the
// contents of the file at `fpath` and decode it for you.
func DecodeFile(fpath string, v interface***REMOVED******REMOVED***) (MetaData, error) ***REMOVED***
	bs, err := ioutil.ReadFile(fpath)
	if err != nil ***REMOVED***
		return MetaData***REMOVED******REMOVED***, err
	***REMOVED***
	return Decode(string(bs), v)
***REMOVED***

// DecodeReader is just like Decode, except it will consume all bytes
// from the reader and decode it for you.
func DecodeReader(r io.Reader, v interface***REMOVED******REMOVED***) (MetaData, error) ***REMOVED***
	bs, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return MetaData***REMOVED******REMOVED***, err
	***REMOVED***
	return Decode(string(bs), v)
***REMOVED***

// unify performs a sort of type unification based on the structure of `rv`,
// which is the client representation.
//
// Any type mismatch produces an error. Finding a type that we don't know
// how to handle produces an unsupported type error.
func (md *MetaData) unify(data interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***

	// Special case. Look for a `Primitive` value.
	if rv.Type() == reflect.TypeOf((*Primitive)(nil)).Elem() ***REMOVED***
		// Save the undecoded data and the key context into the primitive
		// value.
		context := make(Key, len(md.context))
		copy(context, md.context)
		rv.Set(reflect.ValueOf(Primitive***REMOVED***
			undecoded: data,
			context:   context,
		***REMOVED***))
		return nil
	***REMOVED***

	// Special case. Unmarshaler Interface support.
	if rv.CanAddr() ***REMOVED***
		if v, ok := rv.Addr().Interface().(Unmarshaler); ok ***REMOVED***
			return v.UnmarshalTOML(data)
		***REMOVED***
	***REMOVED***

	// Special case. Handle time.Time values specifically.
	// TODO: Remove this code when we decide to drop support for Go 1.1.
	// This isn't necessary in Go 1.2 because time.Time satisfies the encoding
	// interfaces.
	if rv.Type().AssignableTo(rvalue(time.Time***REMOVED******REMOVED***).Type()) ***REMOVED***
		return md.unifyDatetime(data, rv)
	***REMOVED***

	// Special case. Look for a value satisfying the TextUnmarshaler interface.
	if v, ok := rv.Interface().(TextUnmarshaler); ok ***REMOVED***
		return md.unifyText(data, v)
	***REMOVED***
	// BUG(burntsushi)
	// The behavior here is incorrect whenever a Go type satisfies the
	// encoding.TextUnmarshaler interface but also corresponds to a TOML
	// hash or array. In particular, the unmarshaler should only be applied
	// to primitive TOML values. But at this point, it will be applied to
	// all kinds of values and produce an incorrect error whenever those values
	// are hashes or arrays (including arrays of tables).

	k := rv.Kind()

	// laziness
	if k >= reflect.Int && k <= reflect.Uint64 ***REMOVED***
		return md.unifyInt(data, rv)
	***REMOVED***
	switch k ***REMOVED***
	case reflect.Ptr:
		elem := reflect.New(rv.Type().Elem())
		err := md.unify(data, reflect.Indirect(elem))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		rv.Set(elem)
		return nil
	case reflect.Struct:
		return md.unifyStruct(data, rv)
	case reflect.Map:
		return md.unifyMap(data, rv)
	case reflect.Array:
		return md.unifyArray(data, rv)
	case reflect.Slice:
		return md.unifySlice(data, rv)
	case reflect.String:
		return md.unifyString(data, rv)
	case reflect.Bool:
		return md.unifyBool(data, rv)
	case reflect.Interface:
		// we only support empty interfaces.
		if rv.NumMethod() > 0 ***REMOVED***
			return e("Unsupported type '%s'.", rv.Kind())
		***REMOVED***
		return md.unifyAnything(data, rv)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return md.unifyFloat64(data, rv)
	***REMOVED***
	return e("Unsupported type '%s'.", rv.Kind())
***REMOVED***

func (md *MetaData) unifyStruct(mapping interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	tmap, ok := mapping.(map[string]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return mismatch(rv, "map", mapping)
	***REMOVED***

	for key, datum := range tmap ***REMOVED***
		var f *field
		fields := cachedTypeFields(rv.Type())
		for i := range fields ***REMOVED***
			ff := &fields[i]
			if ff.name == key ***REMOVED***
				f = ff
				break
			***REMOVED***
			if f == nil && strings.EqualFold(ff.name, key) ***REMOVED***
				f = ff
			***REMOVED***
		***REMOVED***
		if f != nil ***REMOVED***
			subv := rv
			for _, i := range f.index ***REMOVED***
				subv = indirect(subv.Field(i))
			***REMOVED***
			if isUnifiable(subv) ***REMOVED***
				md.decoded[md.context.add(key).String()] = true
				md.context = append(md.context, key)
				if err := md.unify(datum, subv); err != nil ***REMOVED***
					return e("Type mismatch for '%s.%s': %s",
						rv.Type().String(), f.name, err)
				***REMOVED***
				md.context = md.context[0 : len(md.context)-1]
			***REMOVED*** else if f.name != "" ***REMOVED***
				// Bad user! No soup for you!
				return e("Field '%s.%s' is unexported, and therefore cannot "+
					"be loaded with reflection.", rv.Type().String(), f.name)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (md *MetaData) unifyMap(mapping interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	tmap, ok := mapping.(map[string]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return badtype("map", mapping)
	***REMOVED***
	if rv.IsNil() ***REMOVED***
		rv.Set(reflect.MakeMap(rv.Type()))
	***REMOVED***
	for k, v := range tmap ***REMOVED***
		md.decoded[md.context.add(k).String()] = true
		md.context = append(md.context, k)

		rvkey := indirect(reflect.New(rv.Type().Key()))
		rvval := reflect.Indirect(reflect.New(rv.Type().Elem()))
		if err := md.unify(v, rvval); err != nil ***REMOVED***
			return err
		***REMOVED***
		md.context = md.context[0 : len(md.context)-1]

		rvkey.SetString(k)
		rv.SetMapIndex(rvkey, rvval)
	***REMOVED***
	return nil
***REMOVED***

func (md *MetaData) unifyArray(data interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	datav := reflect.ValueOf(data)
	if datav.Kind() != reflect.Slice ***REMOVED***
		return badtype("slice", data)
	***REMOVED***
	sliceLen := datav.Len()
	if sliceLen != rv.Len() ***REMOVED***
		return e("expected array length %d; got TOML array of length %d",
			rv.Len(), sliceLen)
	***REMOVED***
	return md.unifySliceArray(datav, rv)
***REMOVED***

func (md *MetaData) unifySlice(data interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	datav := reflect.ValueOf(data)
	if datav.Kind() != reflect.Slice ***REMOVED***
		return badtype("slice", data)
	***REMOVED***
	sliceLen := datav.Len()
	if rv.IsNil() ***REMOVED***
		rv.Set(reflect.MakeSlice(rv.Type(), sliceLen, sliceLen))
	***REMOVED***
	return md.unifySliceArray(datav, rv)
***REMOVED***

func (md *MetaData) unifySliceArray(data, rv reflect.Value) error ***REMOVED***
	sliceLen := data.Len()
	for i := 0; i < sliceLen; i++ ***REMOVED***
		v := data.Index(i).Interface()
		sliceval := indirect(rv.Index(i))
		if err := md.unify(v, sliceval); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (md *MetaData) unifyDatetime(data interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	if _, ok := data.(time.Time); ok ***REMOVED***
		rv.Set(reflect.ValueOf(data))
		return nil
	***REMOVED***
	return badtype("time.Time", data)
***REMOVED***

func (md *MetaData) unifyString(data interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	if s, ok := data.(string); ok ***REMOVED***
		rv.SetString(s)
		return nil
	***REMOVED***
	return badtype("string", data)
***REMOVED***

func (md *MetaData) unifyFloat64(data interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	if num, ok := data.(float64); ok ***REMOVED***
		switch rv.Kind() ***REMOVED***
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			rv.SetFloat(num)
		default:
			panic("bug")
		***REMOVED***
		return nil
	***REMOVED***
	return badtype("float", data)
***REMOVED***

func (md *MetaData) unifyInt(data interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	if num, ok := data.(int64); ok ***REMOVED***
		if rv.Kind() >= reflect.Int && rv.Kind() <= reflect.Int64 ***REMOVED***
			switch rv.Kind() ***REMOVED***
			case reflect.Int, reflect.Int64:
				// No bounds checking necessary.
			case reflect.Int8:
				if num < math.MinInt8 || num > math.MaxInt8 ***REMOVED***
					return e("Value '%d' is out of range for int8.", num)
				***REMOVED***
			case reflect.Int16:
				if num < math.MinInt16 || num > math.MaxInt16 ***REMOVED***
					return e("Value '%d' is out of range for int16.", num)
				***REMOVED***
			case reflect.Int32:
				if num < math.MinInt32 || num > math.MaxInt32 ***REMOVED***
					return e("Value '%d' is out of range for int32.", num)
				***REMOVED***
			***REMOVED***
			rv.SetInt(num)
		***REMOVED*** else if rv.Kind() >= reflect.Uint && rv.Kind() <= reflect.Uint64 ***REMOVED***
			unum := uint64(num)
			switch rv.Kind() ***REMOVED***
			case reflect.Uint, reflect.Uint64:
				// No bounds checking necessary.
			case reflect.Uint8:
				if num < 0 || unum > math.MaxUint8 ***REMOVED***
					return e("Value '%d' is out of range for uint8.", num)
				***REMOVED***
			case reflect.Uint16:
				if num < 0 || unum > math.MaxUint16 ***REMOVED***
					return e("Value '%d' is out of range for uint16.", num)
				***REMOVED***
			case reflect.Uint32:
				if num < 0 || unum > math.MaxUint32 ***REMOVED***
					return e("Value '%d' is out of range for uint32.", num)
				***REMOVED***
			***REMOVED***
			rv.SetUint(unum)
		***REMOVED*** else ***REMOVED***
			panic("unreachable")
		***REMOVED***
		return nil
	***REMOVED***
	return badtype("integer", data)
***REMOVED***

func (md *MetaData) unifyBool(data interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	if b, ok := data.(bool); ok ***REMOVED***
		rv.SetBool(b)
		return nil
	***REMOVED***
	return badtype("boolean", data)
***REMOVED***

func (md *MetaData) unifyAnything(data interface***REMOVED******REMOVED***, rv reflect.Value) error ***REMOVED***
	rv.Set(reflect.ValueOf(data))
	return nil
***REMOVED***

func (md *MetaData) unifyText(data interface***REMOVED******REMOVED***, v TextUnmarshaler) error ***REMOVED***
	var s string
	switch sdata := data.(type) ***REMOVED***
	case TextMarshaler:
		text, err := sdata.MarshalText()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		s = string(text)
	case fmt.Stringer:
		s = sdata.String()
	case string:
		s = sdata
	case bool:
		s = fmt.Sprintf("%v", sdata)
	case int64:
		s = fmt.Sprintf("%d", sdata)
	case float64:
		s = fmt.Sprintf("%f", sdata)
	default:
		return badtype("primitive (string-like)", data)
	***REMOVED***
	if err := v.UnmarshalText([]byte(s)); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// rvalue returns a reflect.Value of `v`. All pointers are resolved.
func rvalue(v interface***REMOVED******REMOVED***) reflect.Value ***REMOVED***
	return indirect(reflect.ValueOf(v))
***REMOVED***

// indirect returns the value pointed to by a pointer.
// Pointers are followed until the value is not a pointer.
// New values are allocated for each nil pointer.
//
// An exception to this rule is if the value satisfies an interface of
// interest to us (like encoding.TextUnmarshaler).
func indirect(v reflect.Value) reflect.Value ***REMOVED***
	if v.Kind() != reflect.Ptr ***REMOVED***
		if v.CanAddr() ***REMOVED***
			pv := v.Addr()
			if _, ok := pv.Interface().(TextUnmarshaler); ok ***REMOVED***
				return pv
			***REMOVED***
		***REMOVED***
		return v
	***REMOVED***
	if v.IsNil() ***REMOVED***
		v.Set(reflect.New(v.Type().Elem()))
	***REMOVED***
	return indirect(reflect.Indirect(v))
***REMOVED***

func isUnifiable(rv reflect.Value) bool ***REMOVED***
	if rv.CanSet() ***REMOVED***
		return true
	***REMOVED***
	if _, ok := rv.Interface().(TextUnmarshaler); ok ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func badtype(expected string, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return e("Expected %s but found '%T'.", expected, data)
***REMOVED***

func mismatch(user reflect.Value, expected string, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return e("Type mismatch for %s. Expected %s but found '%T'.",
		user.Type().String(), expected, data)
***REMOVED***
