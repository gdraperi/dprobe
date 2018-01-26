package dbus

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// Variant represents the D-Bus variant type.
type Variant struct ***REMOVED***
	sig   Signature
	value interface***REMOVED******REMOVED***
***REMOVED***

// MakeVariant converts the given value to a Variant. It panics if v cannot be
// represented as a D-Bus type.
func MakeVariant(v interface***REMOVED******REMOVED***) Variant ***REMOVED***
	return Variant***REMOVED***SignatureOf(v), v***REMOVED***
***REMOVED***

// ParseVariant parses the given string as a variant as described at
// https://developer.gnome.org/glib/unstable/gvariant-text.html. If sig is not
// empty, it is taken to be the expected signature for the variant.
func ParseVariant(s string, sig Signature) (Variant, error) ***REMOVED***
	tokens := varLex(s)
	p := &varParser***REMOVED***tokens: tokens***REMOVED***
	n, err := varMakeNode(p)
	if err != nil ***REMOVED***
		return Variant***REMOVED******REMOVED***, err
	***REMOVED***
	if sig.str == "" ***REMOVED***
		sig, err = varInfer(n)
		if err != nil ***REMOVED***
			return Variant***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***
	v, err := n.Value(sig)
	if err != nil ***REMOVED***
		return Variant***REMOVED******REMOVED***, err
	***REMOVED***
	return MakeVariant(v), nil
***REMOVED***

// format returns a formatted version of v and whether this string can be parsed
// unambigously.
func (v Variant) format() (string, bool) ***REMOVED***
	switch v.sig.str[0] ***REMOVED***
	case 'b', 'i':
		return fmt.Sprint(v.value), true
	case 'n', 'q', 'u', 'x', 't', 'd', 'h':
		return fmt.Sprint(v.value), false
	case 's':
		return strconv.Quote(v.value.(string)), true
	case 'o':
		return strconv.Quote(string(v.value.(ObjectPath))), false
	case 'g':
		return strconv.Quote(v.value.(Signature).str), false
	case 'v':
		s, unamb := v.value.(Variant).format()
		if !unamb ***REMOVED***
			return "<@" + v.value.(Variant).sig.str + " " + s + ">", true
		***REMOVED***
		return "<" + s + ">", true
	case 'y':
		return fmt.Sprintf("%#x", v.value.(byte)), false
	***REMOVED***
	rv := reflect.ValueOf(v.value)
	switch rv.Kind() ***REMOVED***
	case reflect.Slice:
		if rv.Len() == 0 ***REMOVED***
			return "[]", false
		***REMOVED***
		unamb := true
		buf := bytes.NewBuffer([]byte("["))
		for i := 0; i < rv.Len(); i++ ***REMOVED***
			// TODO: slooow
			s, b := MakeVariant(rv.Index(i).Interface()).format()
			unamb = unamb && b
			buf.WriteString(s)
			if i != rv.Len()-1 ***REMOVED***
				buf.WriteString(", ")
			***REMOVED***
		***REMOVED***
		buf.WriteByte(']')
		return buf.String(), unamb
	case reflect.Map:
		if rv.Len() == 0 ***REMOVED***
			return "***REMOVED******REMOVED***", false
		***REMOVED***
		unamb := true
		var buf bytes.Buffer
		kvs := make([]string, rv.Len())
		for i, k := range rv.MapKeys() ***REMOVED***
			s, b := MakeVariant(k.Interface()).format()
			unamb = unamb && b
			buf.Reset()
			buf.WriteString(s)
			buf.WriteString(": ")
			s, b = MakeVariant(rv.MapIndex(k).Interface()).format()
			unamb = unamb && b
			buf.WriteString(s)
			kvs[i] = buf.String()
		***REMOVED***
		buf.Reset()
		buf.WriteByte('***REMOVED***')
		sort.Strings(kvs)
		for i, kv := range kvs ***REMOVED***
			if i > 0 ***REMOVED***
				buf.WriteString(", ")
			***REMOVED***
			buf.WriteString(kv)
		***REMOVED***
		buf.WriteByte('***REMOVED***')
		return buf.String(), unamb
	***REMOVED***
	return `"INVALID"`, true
***REMOVED***

// Signature returns the D-Bus signature of the underlying value of v.
func (v Variant) Signature() Signature ***REMOVED***
	return v.sig
***REMOVED***

// String returns the string representation of the underlying value of v as
// described at https://developer.gnome.org/glib/unstable/gvariant-text.html.
func (v Variant) String() string ***REMOVED***
	s, unamb := v.format()
	if !unamb ***REMOVED***
		return "@" + v.sig.str + " " + s
	***REMOVED***
	return s
***REMOVED***

// Value returns the underlying value of v.
func (v Variant) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return v.value
***REMOVED***
