package memdb

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
)

// Indexer is an interface used for defining indexes
type Indexer interface ***REMOVED***
	// ExactFromArgs is used to build an exact index lookup
	// based on arguments
	FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error)
***REMOVED***

// SingleIndexer is an interface used for defining indexes
// generating a single entry per object
type SingleIndexer interface ***REMOVED***
	// FromObject is used to extract an index value from an
	// object or to indicate that the index value is missing.
	FromObject(raw interface***REMOVED******REMOVED***) (bool, []byte, error)
***REMOVED***

// MultiIndexer is an interface used for defining indexes
// generating multiple entries per object
type MultiIndexer interface ***REMOVED***
	// FromObject is used to extract index values from an
	// object or to indicate that the index value is missing.
	FromObject(raw interface***REMOVED******REMOVED***) (bool, [][]byte, error)
***REMOVED***

// PrefixIndexer can optionally be implemented for any
// indexes that support prefix based iteration. This may
// not apply to all indexes.
type PrefixIndexer interface ***REMOVED***
	// PrefixFromArgs returns a prefix that should be used
	// for scanning based on the arguments
	PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error)
***REMOVED***

// StringFieldIndex is used to extract a field from an object
// using reflection and builds an index on that field.
type StringFieldIndex struct ***REMOVED***
	Field     string
	Lowercase bool
***REMOVED***

func (s *StringFieldIndex) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Dereference the pointer if any

	fv := v.FieldByName(s.Field)
	if !fv.IsValid() ***REMOVED***
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", s.Field, obj)
	***REMOVED***

	val := fv.String()
	if val == "" ***REMOVED***
		return false, nil, nil
	***REMOVED***

	if s.Lowercase ***REMOVED***
		val = strings.ToLower(val)
	***REMOVED***

	// Add the null character as a terminator
	val += "\x00"
	return true, []byte(val), nil
***REMOVED***

func (s *StringFieldIndex) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***
	arg, ok := args[0].(string)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("argument must be a string: %#v", args[0])
	***REMOVED***
	if s.Lowercase ***REMOVED***
		arg = strings.ToLower(arg)
	***REMOVED***
	// Add the null character as a terminator
	arg += "\x00"
	return []byte(arg), nil
***REMOVED***

func (s *StringFieldIndex) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	val, err := s.FromArgs(args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Strip the null terminator, the rest is a prefix
	n := len(val)
	if n > 0 ***REMOVED***
		return val[:n-1], nil
	***REMOVED***
	return val, nil
***REMOVED***

// StringSliceFieldIndex is used to extract a field from an object
// using reflection and builds an index on that field.
type StringSliceFieldIndex struct ***REMOVED***
	Field     string
	Lowercase bool
***REMOVED***

func (s *StringSliceFieldIndex) FromObject(obj interface***REMOVED******REMOVED***) (bool, [][]byte, error) ***REMOVED***
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Dereference the pointer if any

	fv := v.FieldByName(s.Field)
	if !fv.IsValid() ***REMOVED***
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", s.Field, obj)
	***REMOVED***

	if fv.Kind() != reflect.Slice || fv.Type().Elem().Kind() != reflect.String ***REMOVED***
		return false, nil, fmt.Errorf("field '%s' is not a string slice", s.Field)
	***REMOVED***

	length := fv.Len()
	vals := make([][]byte, 0, length)
	for i := 0; i < fv.Len(); i++ ***REMOVED***
		val := fv.Index(i).String()
		if val == "" ***REMOVED***
			continue
		***REMOVED***

		if s.Lowercase ***REMOVED***
			val = strings.ToLower(val)
		***REMOVED***

		// Add the null character as a terminator
		val += "\x00"
		vals = append(vals, []byte(val))
	***REMOVED***
	if len(vals) == 0 ***REMOVED***
		return false, nil, nil
	***REMOVED***
	return true, vals, nil
***REMOVED***

func (s *StringSliceFieldIndex) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***
	arg, ok := args[0].(string)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("argument must be a string: %#v", args[0])
	***REMOVED***
	if s.Lowercase ***REMOVED***
		arg = strings.ToLower(arg)
	***REMOVED***
	// Add the null character as a terminator
	arg += "\x00"
	return []byte(arg), nil
***REMOVED***

func (s *StringSliceFieldIndex) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	val, err := s.FromArgs(args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Strip the null terminator, the rest is a prefix
	n := len(val)
	if n > 0 ***REMOVED***
		return val[:n-1], nil
	***REMOVED***
	return val, nil
***REMOVED***

// UUIDFieldIndex is used to extract a field from an object
// using reflection and builds an index on that field by treating
// it as a UUID. This is an optimization to using a StringFieldIndex
// as the UUID can be more compactly represented in byte form.
type UUIDFieldIndex struct ***REMOVED***
	Field string
***REMOVED***

func (u *UUIDFieldIndex) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Dereference the pointer if any

	fv := v.FieldByName(u.Field)
	if !fv.IsValid() ***REMOVED***
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", u.Field, obj)
	***REMOVED***

	val := fv.String()
	if val == "" ***REMOVED***
		return false, nil, nil
	***REMOVED***

	buf, err := u.parseString(val, true)
	return true, buf, err
***REMOVED***

func (u *UUIDFieldIndex) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***
	switch arg := args[0].(type) ***REMOVED***
	case string:
		return u.parseString(arg, true)
	case []byte:
		if len(arg) != 16 ***REMOVED***
			return nil, fmt.Errorf("byte slice must be 16 characters")
		***REMOVED***
		return arg, nil
	default:
		return nil,
			fmt.Errorf("argument must be a string or byte slice: %#v", args[0])
	***REMOVED***
***REMOVED***

func (u *UUIDFieldIndex) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***
	switch arg := args[0].(type) ***REMOVED***
	case string:
		return u.parseString(arg, false)
	case []byte:
		return arg, nil
	default:
		return nil,
			fmt.Errorf("argument must be a string or byte slice: %#v", args[0])
	***REMOVED***
***REMOVED***

// parseString parses a UUID from the string. If enforceLength is false, it will
// parse a partial UUID. An error is returned if the input, stripped of hyphens,
// is not even length.
func (u *UUIDFieldIndex) parseString(s string, enforceLength bool) ([]byte, error) ***REMOVED***
	// Verify the length
	l := len(s)
	if enforceLength && l != 36 ***REMOVED***
		return nil, fmt.Errorf("UUID must be 36 characters")
	***REMOVED*** else if l > 36 ***REMOVED***
		return nil, fmt.Errorf("Invalid UUID length. UUID have 36 characters; got %d", l)
	***REMOVED***

	hyphens := strings.Count(s, "-")
	if hyphens > 4 ***REMOVED***
		return nil, fmt.Errorf(`UUID should have maximum of 4 "-"; got %d`, hyphens)
	***REMOVED***

	// The sanitized length is the length of the original string without the "-".
	sanitized := strings.Replace(s, "-", "", -1)
	sanitizedLength := len(sanitized)
	if sanitizedLength%2 != 0 ***REMOVED***
		return nil, fmt.Errorf("Input (without hyphens) must be even length")
	***REMOVED***

	dec, err := hex.DecodeString(sanitized)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Invalid UUID: %v", err)
	***REMOVED***

	return dec, nil
***REMOVED***

// FieldSetIndex is used to extract a field from an object using reflection and
// builds an index on whether the field is set by comparing it against its
// type's nil value.
type FieldSetIndex struct ***REMOVED***
	Field string
***REMOVED***

func (f *FieldSetIndex) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Dereference the pointer if any

	fv := v.FieldByName(f.Field)
	if !fv.IsValid() ***REMOVED***
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", f.Field, obj)
	***REMOVED***

	if fv.Interface() == reflect.Zero(fv.Type()).Interface() ***REMOVED***
		return true, []byte***REMOVED***0***REMOVED***, nil
	***REMOVED***

	return true, []byte***REMOVED***1***REMOVED***, nil
***REMOVED***

func (f *FieldSetIndex) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromBoolArgs(args)
***REMOVED***

// ConditionalIndex builds an index based on a condition specified by a passed
// user function. This function may examine the passed object and return a
// boolean to encapsulate an arbitrarily complex conditional.
type ConditionalIndex struct ***REMOVED***
	Conditional ConditionalIndexFunc
***REMOVED***

// ConditionalIndexFunc is the required function interface for a
// ConditionalIndex.
type ConditionalIndexFunc func(obj interface***REMOVED******REMOVED***) (bool, error)

func (c *ConditionalIndex) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	// Call the user's function
	res, err := c.Conditional(obj)
	if err != nil ***REMOVED***
		return false, nil, fmt.Errorf("ConditionalIndexFunc(%#v) failed: %v", obj, err)
	***REMOVED***

	if res ***REMOVED***
		return true, []byte***REMOVED***1***REMOVED***, nil
	***REMOVED***

	return true, []byte***REMOVED***0***REMOVED***, nil
***REMOVED***

func (c *ConditionalIndex) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromBoolArgs(args)
***REMOVED***

// fromBoolArgs is a helper that expects only a single boolean argument and
// returns a single length byte array containing either a one or zero depending
// on whether the passed input is true or false respectively.
func fromBoolArgs(args []interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***

	if val, ok := args[0].(bool); !ok ***REMOVED***
		return nil, fmt.Errorf("argument must be a boolean type: %#v", args[0])
	***REMOVED*** else if val ***REMOVED***
		return []byte***REMOVED***1***REMOVED***, nil
	***REMOVED***

	return []byte***REMOVED***0***REMOVED***, nil
***REMOVED***

// CompoundIndex is used to build an index using multiple sub-indexes
// Prefix based iteration is supported as long as the appropriate prefix
// of indexers support it. All sub-indexers are only assumed to expect
// a single argument.
type CompoundIndex struct ***REMOVED***
	Indexes []Indexer

	// AllowMissing results in an index based on only the indexers
	// that return data. If true, you may end up with 2/3 columns
	// indexed which might be useful for an index scan. Otherwise,
	// the CompoundIndex requires all indexers to be satisfied.
	AllowMissing bool
***REMOVED***

func (c *CompoundIndex) FromObject(raw interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	var out []byte
	for i, idxRaw := range c.Indexes ***REMOVED***
		idx, ok := idxRaw.(SingleIndexer)
		if !ok ***REMOVED***
			return false, nil, fmt.Errorf("sub-index %d error: %s", i, "sub-index must be a SingleIndexer")
		***REMOVED***
		ok, val, err := idx.FromObject(raw)
		if err != nil ***REMOVED***
			return false, nil, fmt.Errorf("sub-index %d error: %v", i, err)
		***REMOVED***
		if !ok ***REMOVED***
			if c.AllowMissing ***REMOVED***
				break
			***REMOVED*** else ***REMOVED***
				return false, nil, nil
			***REMOVED***
		***REMOVED***
		out = append(out, val...)
	***REMOVED***
	return true, out, nil
***REMOVED***

func (c *CompoundIndex) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != len(c.Indexes) ***REMOVED***
		return nil, fmt.Errorf("less arguments than index fields")
	***REMOVED***
	var out []byte
	for i, arg := range args ***REMOVED***
		val, err := c.Indexes[i].FromArgs(arg)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("sub-index %d error: %v", i, err)
		***REMOVED***
		out = append(out, val...)
	***REMOVED***
	return out, nil
***REMOVED***

func (c *CompoundIndex) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) > len(c.Indexes) ***REMOVED***
		return nil, fmt.Errorf("more arguments than index fields")
	***REMOVED***
	var out []byte
	for i, arg := range args ***REMOVED***
		if i+1 < len(args) ***REMOVED***
			val, err := c.Indexes[i].FromArgs(arg)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("sub-index %d error: %v", i, err)
			***REMOVED***
			out = append(out, val...)
		***REMOVED*** else ***REMOVED***
			prefixIndexer, ok := c.Indexes[i].(PrefixIndexer)
			if !ok ***REMOVED***
				return nil, fmt.Errorf("sub-index %d does not support prefix scanning", i)
			***REMOVED***
			val, err := prefixIndexer.PrefixFromArgs(arg)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("sub-index %d error: %v", i, err)
			***REMOVED***
			out = append(out, val...)
		***REMOVED***
	***REMOVED***
	return out, nil
***REMOVED***
