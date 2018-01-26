package dbus

import (
	"fmt"
	"reflect"
	"strings"
)

var sigToType = map[byte]reflect.Type***REMOVED***
	'y': byteType,
	'b': boolType,
	'n': int16Type,
	'q': uint16Type,
	'i': int32Type,
	'u': uint32Type,
	'x': int64Type,
	't': uint64Type,
	'd': float64Type,
	's': stringType,
	'g': signatureType,
	'o': objectPathType,
	'v': variantType,
	'h': unixFDIndexType,
***REMOVED***

// Signature represents a correct type signature as specified by the D-Bus
// specification. The zero value represents the empty signature, "".
type Signature struct ***REMOVED***
	str string
***REMOVED***

// SignatureOf returns the concatenation of all the signatures of the given
// values. It panics if one of them is not representable in D-Bus.
func SignatureOf(vs ...interface***REMOVED******REMOVED***) Signature ***REMOVED***
	var s string
	for _, v := range vs ***REMOVED***
		s += getSignature(reflect.TypeOf(v))
	***REMOVED***
	return Signature***REMOVED***s***REMOVED***
***REMOVED***

// SignatureOfType returns the signature of the given type. It panics if the
// type is not representable in D-Bus.
func SignatureOfType(t reflect.Type) Signature ***REMOVED***
	return Signature***REMOVED***getSignature(t)***REMOVED***
***REMOVED***

// getSignature returns the signature of the given type and panics on unknown types.
func getSignature(t reflect.Type) string ***REMOVED***
	// handle simple types first
	switch t.Kind() ***REMOVED***
	case reflect.Uint8:
		return "y"
	case reflect.Bool:
		return "b"
	case reflect.Int16:
		return "n"
	case reflect.Uint16:
		return "q"
	case reflect.Int32:
		if t == unixFDType ***REMOVED***
			return "h"
		***REMOVED***
		return "i"
	case reflect.Uint32:
		if t == unixFDIndexType ***REMOVED***
			return "h"
		***REMOVED***
		return "u"
	case reflect.Int64:
		return "x"
	case reflect.Uint64:
		return "t"
	case reflect.Float64:
		return "d"
	case reflect.Ptr:
		return getSignature(t.Elem())
	case reflect.String:
		if t == objectPathType ***REMOVED***
			return "o"
		***REMOVED***
		return "s"
	case reflect.Struct:
		if t == variantType ***REMOVED***
			return "v"
		***REMOVED*** else if t == signatureType ***REMOVED***
			return "g"
		***REMOVED***
		var s string
		for i := 0; i < t.NumField(); i++ ***REMOVED***
			field := t.Field(i)
			if field.PkgPath == "" && field.Tag.Get("dbus") != "-" ***REMOVED***
				s += getSignature(t.Field(i).Type)
			***REMOVED***
		***REMOVED***
		return "(" + s + ")"
	case reflect.Array, reflect.Slice:
		return "a" + getSignature(t.Elem())
	case reflect.Map:
		if !isKeyType(t.Key()) ***REMOVED***
			panic(InvalidTypeError***REMOVED***t***REMOVED***)
		***REMOVED***
		return "a***REMOVED***" + getSignature(t.Key()) + getSignature(t.Elem()) + "***REMOVED***"
	***REMOVED***
	panic(InvalidTypeError***REMOVED***t***REMOVED***)
***REMOVED***

// ParseSignature returns the signature represented by this string, or a
// SignatureError if the string is not a valid signature.
func ParseSignature(s string) (sig Signature, err error) ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return
	***REMOVED***
	if len(s) > 255 ***REMOVED***
		return Signature***REMOVED***""***REMOVED***, SignatureError***REMOVED***s, "too long"***REMOVED***
	***REMOVED***
	sig.str = s
	for err == nil && len(s) != 0 ***REMOVED***
		err, s = validSingle(s, 0)
	***REMOVED***
	if err != nil ***REMOVED***
		sig = Signature***REMOVED***""***REMOVED***
	***REMOVED***

	return
***REMOVED***

// ParseSignatureMust behaves like ParseSignature, except that it panics if s
// is not valid.
func ParseSignatureMust(s string) Signature ***REMOVED***
	sig, err := ParseSignature(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return sig
***REMOVED***

// Empty retruns whether the signature is the empty signature.
func (s Signature) Empty() bool ***REMOVED***
	return s.str == ""
***REMOVED***

// Single returns whether the signature represents a single, complete type.
func (s Signature) Single() bool ***REMOVED***
	err, r := validSingle(s.str, 0)
	return err != nil && r == ""
***REMOVED***

// String returns the signature's string representation.
func (s Signature) String() string ***REMOVED***
	return s.str
***REMOVED***

// A SignatureError indicates that a signature passed to a function or received
// on a connection is not a valid signature.
type SignatureError struct ***REMOVED***
	Sig    string
	Reason string
***REMOVED***

func (e SignatureError) Error() string ***REMOVED***
	return fmt.Sprintf("dbus: invalid signature: %q (%s)", e.Sig, e.Reason)
***REMOVED***

// Try to read a single type from this string. If it was successfull, err is nil
// and rem is the remaining unparsed part. Otherwise, err is a non-nil
// SignatureError and rem is "". depth is the current recursion depth which may
// not be greater than 64 and should be given as 0 on the first call.
func validSingle(s string, depth int) (err error, rem string) ***REMOVED***
	if s == "" ***REMOVED***
		return SignatureError***REMOVED***Sig: s, Reason: "empty signature"***REMOVED***, ""
	***REMOVED***
	if depth > 64 ***REMOVED***
		return SignatureError***REMOVED***Sig: s, Reason: "container nesting too deep"***REMOVED***, ""
	***REMOVED***
	switch s[0] ***REMOVED***
	case 'y', 'b', 'n', 'q', 'i', 'u', 'x', 't', 'd', 's', 'g', 'o', 'v', 'h':
		return nil, s[1:]
	case 'a':
		if len(s) > 1 && s[1] == '***REMOVED***' ***REMOVED***
			i := findMatching(s[1:], '***REMOVED***', '***REMOVED***')
			if i == -1 ***REMOVED***
				return SignatureError***REMOVED***Sig: s, Reason: "unmatched '***REMOVED***'"***REMOVED***, ""
			***REMOVED***
			i++
			rem = s[i+1:]
			s = s[2:i]
			if err, _ = validSingle(s[:1], depth+1); err != nil ***REMOVED***
				return err, ""
			***REMOVED***
			err, nr := validSingle(s[1:], depth+1)
			if err != nil ***REMOVED***
				return err, ""
			***REMOVED***
			if nr != "" ***REMOVED***
				return SignatureError***REMOVED***Sig: s, Reason: "too many types in dict"***REMOVED***, ""
			***REMOVED***
			return nil, rem
		***REMOVED***
		return validSingle(s[1:], depth+1)
	case '(':
		i := findMatching(s, '(', ')')
		if i == -1 ***REMOVED***
			return SignatureError***REMOVED***Sig: s, Reason: "unmatched ')'"***REMOVED***, ""
		***REMOVED***
		rem = s[i+1:]
		s = s[1:i]
		for err == nil && s != "" ***REMOVED***
			err, s = validSingle(s, depth+1)
		***REMOVED***
		if err != nil ***REMOVED***
			rem = ""
		***REMOVED***
		return
	***REMOVED***
	return SignatureError***REMOVED***Sig: s, Reason: "invalid type character"***REMOVED***, ""
***REMOVED***

func findMatching(s string, left, right rune) int ***REMOVED***
	n := 0
	for i, v := range s ***REMOVED***
		if v == left ***REMOVED***
			n++
		***REMOVED*** else if v == right ***REMOVED***
			n--
		***REMOVED***
		if n == 0 ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// typeFor returns the type of the given signature. It ignores any left over
// characters and panics if s doesn't start with a valid type signature.
func typeFor(s string) (t reflect.Type) ***REMOVED***
	err, _ := validSingle(s, 0)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	if t, ok := sigToType[s[0]]; ok ***REMOVED***
		return t
	***REMOVED***
	switch s[0] ***REMOVED***
	case 'a':
		if s[1] == '***REMOVED***' ***REMOVED***
			i := strings.LastIndex(s, "***REMOVED***")
			t = reflect.MapOf(sigToType[s[2]], typeFor(s[3:i]))
		***REMOVED*** else ***REMOVED***
			t = reflect.SliceOf(typeFor(s[1:]))
		***REMOVED***
	case '(':
		t = interfacesType
	***REMOVED***
	return
***REMOVED***
