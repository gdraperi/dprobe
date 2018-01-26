package dbus

import (
	"errors"
	"reflect"
	"strings"
)

var (
	byteType        = reflect.TypeOf(byte(0))
	boolType        = reflect.TypeOf(false)
	uint8Type       = reflect.TypeOf(uint8(0))
	int16Type       = reflect.TypeOf(int16(0))
	uint16Type      = reflect.TypeOf(uint16(0))
	int32Type       = reflect.TypeOf(int32(0))
	uint32Type      = reflect.TypeOf(uint32(0))
	int64Type       = reflect.TypeOf(int64(0))
	uint64Type      = reflect.TypeOf(uint64(0))
	float64Type     = reflect.TypeOf(float64(0))
	stringType      = reflect.TypeOf("")
	signatureType   = reflect.TypeOf(Signature***REMOVED***""***REMOVED***)
	objectPathType  = reflect.TypeOf(ObjectPath(""))
	variantType     = reflect.TypeOf(Variant***REMOVED***Signature***REMOVED***""***REMOVED***, nil***REMOVED***)
	interfacesType  = reflect.TypeOf([]interface***REMOVED******REMOVED******REMOVED******REMOVED***)
	unixFDType      = reflect.TypeOf(UnixFD(0))
	unixFDIndexType = reflect.TypeOf(UnixFDIndex(0))
)

// An InvalidTypeError signals that a value which cannot be represented in the
// D-Bus wire format was passed to a function.
type InvalidTypeError struct ***REMOVED***
	Type reflect.Type
***REMOVED***

func (e InvalidTypeError) Error() string ***REMOVED***
	return "dbus: invalid type " + e.Type.String()
***REMOVED***

// Store copies the values contained in src to dest, which must be a slice of
// pointers. It converts slices of interfaces from src to corresponding structs
// in dest. An error is returned if the lengths of src and dest or the types of
// their elements don't match.
func Store(src []interface***REMOVED******REMOVED***, dest ...interface***REMOVED******REMOVED***) error ***REMOVED***
	if len(src) != len(dest) ***REMOVED***
		return errors.New("dbus.Store: length mismatch")
	***REMOVED***

	for i := range src ***REMOVED***
		if err := store(src[i], dest[i]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func store(src, dest interface***REMOVED******REMOVED***) error ***REMOVED***
	if reflect.TypeOf(dest).Elem() == reflect.TypeOf(src) ***REMOVED***
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(src))
		return nil
	***REMOVED*** else if hasStruct(dest) ***REMOVED***
		rv := reflect.ValueOf(dest).Elem()
		switch rv.Kind() ***REMOVED***
		case reflect.Struct:
			vs, ok := src.([]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return errors.New("dbus.Store: type mismatch")
			***REMOVED***
			t := rv.Type()
			ndest := make([]interface***REMOVED******REMOVED***, 0, rv.NumField())
			for i := 0; i < rv.NumField(); i++ ***REMOVED***
				field := t.Field(i)
				if field.PkgPath == "" && field.Tag.Get("dbus") != "-" ***REMOVED***
					ndest = append(ndest, rv.Field(i).Addr().Interface())
				***REMOVED***
			***REMOVED***
			if len(vs) != len(ndest) ***REMOVED***
				return errors.New("dbus.Store: type mismatch")
			***REMOVED***
			err := Store(vs, ndest...)
			if err != nil ***REMOVED***
				return errors.New("dbus.Store: type mismatch")
			***REMOVED***
		case reflect.Slice:
			sv := reflect.ValueOf(src)
			if sv.Kind() != reflect.Slice ***REMOVED***
				return errors.New("dbus.Store: type mismatch")
			***REMOVED***
			rv.Set(reflect.MakeSlice(rv.Type(), sv.Len(), sv.Len()))
			for i := 0; i < sv.Len(); i++ ***REMOVED***
				if err := store(sv.Index(i).Interface(), rv.Index(i).Addr().Interface()); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		case reflect.Map:
			sv := reflect.ValueOf(src)
			if sv.Kind() != reflect.Map ***REMOVED***
				return errors.New("dbus.Store: type mismatch")
			***REMOVED***
			keys := sv.MapKeys()
			rv.Set(reflect.MakeMap(sv.Type()))
			for _, key := range keys ***REMOVED***
				v := reflect.New(sv.Type().Elem())
				if err := store(v, sv.MapIndex(key).Interface()); err != nil ***REMOVED***
					return err
				***REMOVED***
				rv.SetMapIndex(key, v.Elem())
			***REMOVED***
		default:
			return errors.New("dbus.Store: type mismatch")
		***REMOVED***
		return nil
	***REMOVED*** else ***REMOVED***
		return errors.New("dbus.Store: type mismatch")
	***REMOVED***
***REMOVED***

func hasStruct(v interface***REMOVED******REMOVED***) bool ***REMOVED***
	t := reflect.TypeOf(v)
	for ***REMOVED***
		switch t.Kind() ***REMOVED***
		case reflect.Struct:
			return true
		case reflect.Slice, reflect.Ptr, reflect.Map:
			t = t.Elem()
		default:
			return false
		***REMOVED***
	***REMOVED***
***REMOVED***

// An ObjectPath is an object path as defined by the D-Bus spec.
type ObjectPath string

// IsValid returns whether the object path is valid.
func (o ObjectPath) IsValid() bool ***REMOVED***
	s := string(o)
	if len(s) == 0 ***REMOVED***
		return false
	***REMOVED***
	if s[0] != '/' ***REMOVED***
		return false
	***REMOVED***
	if s[len(s)-1] == '/' && len(s) != 1 ***REMOVED***
		return false
	***REMOVED***
	// probably not used, but technically possible
	if s == "/" ***REMOVED***
		return true
	***REMOVED***
	split := strings.Split(s[1:], "/")
	for _, v := range split ***REMOVED***
		if len(v) == 0 ***REMOVED***
			return false
		***REMOVED***
		for _, c := range v ***REMOVED***
			if !isMemberChar(c) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// A UnixFD is a Unix file descriptor sent over the wire. See the package-level
// documentation for more information about Unix file descriptor passsing.
type UnixFD int32

// A UnixFDIndex is the representation of a Unix file descriptor in a message.
type UnixFDIndex uint32

// alignment returns the alignment of values of type t.
func alignment(t reflect.Type) int ***REMOVED***
	switch t ***REMOVED***
	case variantType:
		return 1
	case objectPathType:
		return 4
	case signatureType:
		return 1
	case interfacesType: // sometimes used for structs
		return 8
	***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Uint8:
		return 1
	case reflect.Uint16, reflect.Int16:
		return 2
	case reflect.Uint32, reflect.Int32, reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return 4
	case reflect.Uint64, reflect.Int64, reflect.Float64, reflect.Struct:
		return 8
	case reflect.Ptr:
		return alignment(t.Elem())
	***REMOVED***
	return 1
***REMOVED***

// isKeyType returns whether t is a valid type for a D-Bus dict.
func isKeyType(t reflect.Type) bool ***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float64,
		reflect.String:

		return true
	***REMOVED***
	return false
***REMOVED***

// isValidInterface returns whether s is a valid name for an interface.
func isValidInterface(s string) bool ***REMOVED***
	if len(s) == 0 || len(s) > 255 || s[0] == '.' ***REMOVED***
		return false
	***REMOVED***
	elem := strings.Split(s, ".")
	if len(elem) < 2 ***REMOVED***
		return false
	***REMOVED***
	for _, v := range elem ***REMOVED***
		if len(v) == 0 ***REMOVED***
			return false
		***REMOVED***
		if v[0] >= '0' && v[0] <= '9' ***REMOVED***
			return false
		***REMOVED***
		for _, c := range v ***REMOVED***
			if !isMemberChar(c) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// isValidMember returns whether s is a valid name for a member.
func isValidMember(s string) bool ***REMOVED***
	if len(s) == 0 || len(s) > 255 ***REMOVED***
		return false
	***REMOVED***
	i := strings.Index(s, ".")
	if i != -1 ***REMOVED***
		return false
	***REMOVED***
	if s[0] >= '0' && s[0] <= '9' ***REMOVED***
		return false
	***REMOVED***
	for _, c := range s ***REMOVED***
		if !isMemberChar(c) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func isMemberChar(c rune) bool ***REMOVED***
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') || c == '_'
***REMOVED***
