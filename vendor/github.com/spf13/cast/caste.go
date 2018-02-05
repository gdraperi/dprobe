// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var errNegativeNotAllowed = errors.New("unable to cast negative value")

// ToTimeE casts an interface to a time.Time type.
func ToTimeE(i interface***REMOVED******REMOVED***) (tim time.Time, err error) ***REMOVED***
	i = indirect(i)

	switch v := i.(type) ***REMOVED***
	case time.Time:
		return v, nil
	case string:
		return StringToDate(v)
	case int:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(v, 0), nil
	case int32:
		return time.Unix(int64(v), 0), nil
	case uint:
		return time.Unix(int64(v), 0), nil
	case uint64:
		return time.Unix(int64(v), 0), nil
	case uint32:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
	***REMOVED***
***REMOVED***

// ToDurationE casts an interface to a time.Duration type.
func ToDurationE(i interface***REMOVED******REMOVED***) (d time.Duration, err error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case time.Duration:
		return s, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		d = time.Duration(ToInt64(s))
		return
	case float32, float64:
		d = time.Duration(ToFloat64(s))
		return
	case string:
		if strings.ContainsAny(s, "nsuµmh") ***REMOVED***
			d, err = time.ParseDuration(s)
		***REMOVED*** else ***REMOVED***
			d, err = time.ParseDuration(s + "ns")
		***REMOVED***
		return
	default:
		err = fmt.Errorf("unable to cast %#v of type %T to Duration", i, i)
		return
	***REMOVED***
***REMOVED***

// ToBoolE casts an interface to a bool type.
func ToBoolE(i interface***REMOVED******REMOVED***) (bool, error) ***REMOVED***
	i = indirect(i)

	switch b := i.(type) ***REMOVED***
	case bool:
		return b, nil
	case nil:
		return false, nil
	case int:
		if i.(int) != 0 ***REMOVED***
			return true, nil
		***REMOVED***
		return false, nil
	case string:
		return strconv.ParseBool(i.(string))
	default:
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	***REMOVED***
***REMOVED***

// ToFloat64E casts an interface to a float64 type.
func ToFloat64E(i interface***REMOVED******REMOVED***) (float64, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	case int:
		return float64(s), nil
	case int64:
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int8:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case uint64:
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	case string:
		v, err := strconv.ParseFloat(s, 64)
		if err == nil ***REMOVED***
			return v, nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	***REMOVED***
***REMOVED***

// ToFloat32E casts an interface to a float32 type.
func ToFloat32E(i interface***REMOVED******REMOVED***) (float32, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case float64:
		return float32(s), nil
	case float32:
		return s, nil
	case int:
		return float32(s), nil
	case int64:
		return float32(s), nil
	case int32:
		return float32(s), nil
	case int16:
		return float32(s), nil
	case int8:
		return float32(s), nil
	case uint:
		return float32(s), nil
	case uint64:
		return float32(s), nil
	case uint32:
		return float32(s), nil
	case uint16:
		return float32(s), nil
	case uint8:
		return float32(s), nil
	case string:
		v, err := strconv.ParseFloat(s, 32)
		if err == nil ***REMOVED***
			return float32(v), nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v of type %T to float32", i, i)
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float32", i, i)
	***REMOVED***
***REMOVED***

// ToInt64E casts an interface to an int64 type.
func ToInt64E(i interface***REMOVED******REMOVED***) (int64, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case int:
		return int64(s), nil
	case int64:
		return s, nil
	case int32:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int8:
		return int64(s), nil
	case uint:
		return int64(s), nil
	case uint64:
		return int64(s), nil
	case uint32:
		return int64(s), nil
	case uint16:
		return int64(s), nil
	case uint8:
		return int64(s), nil
	case float64:
		return int64(s), nil
	case float32:
		return int64(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil ***REMOVED***
			return v, nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	***REMOVED***
***REMOVED***

// ToInt32E casts an interface to an int32 type.
func ToInt32E(i interface***REMOVED******REMOVED***) (int32, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case int:
		return int32(s), nil
	case int64:
		return int32(s), nil
	case int32:
		return s, nil
	case int16:
		return int32(s), nil
	case int8:
		return int32(s), nil
	case uint:
		return int32(s), nil
	case uint64:
		return int32(s), nil
	case uint32:
		return int32(s), nil
	case uint16:
		return int32(s), nil
	case uint8:
		return int32(s), nil
	case float64:
		return int32(s), nil
	case float32:
		return int32(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil ***REMOVED***
			return int32(v), nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v of type %T to int32", i, i)
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int32", i, i)
	***REMOVED***
***REMOVED***

// ToInt16E casts an interface to an int16 type.
func ToInt16E(i interface***REMOVED******REMOVED***) (int16, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case int:
		return int16(s), nil
	case int64:
		return int16(s), nil
	case int32:
		return int16(s), nil
	case int16:
		return s, nil
	case int8:
		return int16(s), nil
	case uint:
		return int16(s), nil
	case uint64:
		return int16(s), nil
	case uint32:
		return int16(s), nil
	case uint16:
		return int16(s), nil
	case uint8:
		return int16(s), nil
	case float64:
		return int16(s), nil
	case float32:
		return int16(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil ***REMOVED***
			return int16(v), nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v of type %T to int16", i, i)
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int16", i, i)
	***REMOVED***
***REMOVED***

// ToInt8E casts an interface to an int8 type.
func ToInt8E(i interface***REMOVED******REMOVED***) (int8, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case int:
		return int8(s), nil
	case int64:
		return int8(s), nil
	case int32:
		return int8(s), nil
	case int16:
		return int8(s), nil
	case int8:
		return s, nil
	case uint:
		return int8(s), nil
	case uint64:
		return int8(s), nil
	case uint32:
		return int8(s), nil
	case uint16:
		return int8(s), nil
	case uint8:
		return int8(s), nil
	case float64:
		return int8(s), nil
	case float32:
		return int8(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil ***REMOVED***
			return int8(v), nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v of type %T to int8", i, i)
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int8", i, i)
	***REMOVED***
***REMOVED***

// ToIntE casts an interface to an int type.
func ToIntE(i interface***REMOVED******REMOVED***) (int, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case int:
		return s, nil
	case int64:
		return int(s), nil
	case int32:
		return int(s), nil
	case int16:
		return int(s), nil
	case int8:
		return int(s), nil
	case uint:
		return int(s), nil
	case uint64:
		return int(s), nil
	case uint32:
		return int(s), nil
	case uint16:
		return int(s), nil
	case uint8:
		return int(s), nil
	case float64:
		return int(s), nil
	case float32:
		return int(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil ***REMOVED***
			return int(v), nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	***REMOVED***
***REMOVED***

// ToUintE casts an interface to a uint type.
func ToUintE(i interface***REMOVED******REMOVED***) (uint, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case string:
		v, err := strconv.ParseUint(s, 0, 0)
		if err == nil ***REMOVED***
			return uint(v), nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v to uint: %s", i, err)
	case int:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint(s), nil
	case int64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint(s), nil
	case int32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint(s), nil
	case int16:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint(s), nil
	case int8:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint(s), nil
	case uint:
		return s, nil
	case uint64:
		return uint(s), nil
	case uint32:
		return uint(s), nil
	case uint16:
		return uint(s), nil
	case uint8:
		return uint(s), nil
	case float64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint(s), nil
	case float32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint(s), nil
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint", i, i)
	***REMOVED***
***REMOVED***

// ToUint64E casts an interface to a uint64 type.
func ToUint64E(i interface***REMOVED******REMOVED***) (uint64, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case string:
		v, err := strconv.ParseUint(s, 0, 64)
		if err == nil ***REMOVED***
			return v, nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v to uint64: %s", i, err)
	case int:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint64(s), nil
	case int64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint64(s), nil
	case int32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint64(s), nil
	case int16:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint64(s), nil
	case int8:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint64(s), nil
	case uint:
		return uint64(s), nil
	case uint64:
		return s, nil
	case uint32:
		return uint64(s), nil
	case uint16:
		return uint64(s), nil
	case uint8:
		return uint64(s), nil
	case float32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint64(s), nil
	case float64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint64(s), nil
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint64", i, i)
	***REMOVED***
***REMOVED***

// ToUint32E casts an interface to a uint32 type.
func ToUint32E(i interface***REMOVED******REMOVED***) (uint32, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case string:
		v, err := strconv.ParseUint(s, 0, 32)
		if err == nil ***REMOVED***
			return uint32(v), nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v to uint32: %s", i, err)
	case int:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint32(s), nil
	case int64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint32(s), nil
	case int32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint32(s), nil
	case int16:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint32(s), nil
	case int8:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint32(s), nil
	case uint:
		return uint32(s), nil
	case uint64:
		return uint32(s), nil
	case uint32:
		return s, nil
	case uint16:
		return uint32(s), nil
	case uint8:
		return uint32(s), nil
	case float64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint32(s), nil
	case float32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint32(s), nil
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint32", i, i)
	***REMOVED***
***REMOVED***

// ToUint16E casts an interface to a uint16 type.
func ToUint16E(i interface***REMOVED******REMOVED***) (uint16, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case string:
		v, err := strconv.ParseUint(s, 0, 16)
		if err == nil ***REMOVED***
			return uint16(v), nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v to uint16: %s", i, err)
	case int:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint16(s), nil
	case int64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint16(s), nil
	case int32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint16(s), nil
	case int16:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint16(s), nil
	case int8:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint16(s), nil
	case uint:
		return uint16(s), nil
	case uint64:
		return uint16(s), nil
	case uint32:
		return uint16(s), nil
	case uint16:
		return s, nil
	case uint8:
		return uint16(s), nil
	case float64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint16(s), nil
	case float32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint16(s), nil
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint16", i, i)
	***REMOVED***
***REMOVED***

// ToUint8E casts an interface to a uint type.
func ToUint8E(i interface***REMOVED******REMOVED***) (uint8, error) ***REMOVED***
	i = indirect(i)

	switch s := i.(type) ***REMOVED***
	case string:
		v, err := strconv.ParseUint(s, 0, 8)
		if err == nil ***REMOVED***
			return uint8(v), nil
		***REMOVED***
		return 0, fmt.Errorf("unable to cast %#v to uint8: %s", i, err)
	case int:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint8(s), nil
	case int64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint8(s), nil
	case int32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint8(s), nil
	case int16:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint8(s), nil
	case int8:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint8(s), nil
	case uint:
		return uint8(s), nil
	case uint64:
		return uint8(s), nil
	case uint32:
		return uint8(s), nil
	case uint16:
		return uint8(s), nil
	case uint8:
		return s, nil
	case float64:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint8(s), nil
	case float32:
		if s < 0 ***REMOVED***
			return 0, errNegativeNotAllowed
		***REMOVED***
		return uint8(s), nil
	case bool:
		if s ***REMOVED***
			return 1, nil
		***REMOVED***
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint8", i, i)
	***REMOVED***
***REMOVED***

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if a == nil ***REMOVED***
		return nil
	***REMOVED***
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr ***REMOVED***
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	***REMOVED***
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() ***REMOVED***
		v = v.Elem()
	***REMOVED***
	return v.Interface()
***REMOVED***

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirectToStringerOrError returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of fmt.Stringer
// or error,
func indirectToStringerOrError(a interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if a == nil ***REMOVED***
		return nil
	***REMOVED***

	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() ***REMOVED***
		v = v.Elem()
	***REMOVED***
	return v.Interface()
***REMOVED***

// ToStringE casts an interface to a string type.
func ToStringE(i interface***REMOVED******REMOVED***) (string, error) ***REMOVED***
	i = indirectToStringerOrError(i)

	switch s := i.(type) ***REMOVED***
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatInt(int64(s), 10), nil
	case uint64:
		return strconv.FormatInt(int64(s), 10), nil
	case uint32:
		return strconv.FormatInt(int64(s), 10), nil
	case uint16:
		return strconv.FormatInt(int64(s), 10), nil
	case uint8:
		return strconv.FormatInt(int64(s), 10), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	***REMOVED***
***REMOVED***

// ToStringMapStringE casts an interface to a map[string]string type.
func ToStringMapStringE(i interface***REMOVED******REMOVED***) (map[string]string, error) ***REMOVED***
	var m = map[string]string***REMOVED******REMOVED***

	switch v := i.(type) ***REMOVED***
	case map[string]string:
		return v, nil
	case map[string]interface***REMOVED******REMOVED***:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = ToString(val)
		***REMOVED***
		return m, nil
	case map[interface***REMOVED******REMOVED***]string:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = ToString(val)
		***REMOVED***
		return m, nil
	case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = ToString(val)
		***REMOVED***
		return m, nil
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]string", i, i)
	***REMOVED***
***REMOVED***

// ToStringMapStringSliceE casts an interface to a map[string][]string type.
func ToStringMapStringSliceE(i interface***REMOVED******REMOVED***) (map[string][]string, error) ***REMOVED***
	var m = map[string][]string***REMOVED******REMOVED***

	switch v := i.(type) ***REMOVED***
	case map[string][]string:
		return v, nil
	case map[string][]interface***REMOVED******REMOVED***:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = ToStringSlice(val)
		***REMOVED***
		return m, nil
	case map[string]string:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = []string***REMOVED***val***REMOVED***
		***REMOVED***
	case map[string]interface***REMOVED******REMOVED***:
		for k, val := range v ***REMOVED***
			switch vt := val.(type) ***REMOVED***
			case []interface***REMOVED******REMOVED***:
				m[ToString(k)] = ToStringSlice(vt)
			case []string:
				m[ToString(k)] = vt
			default:
				m[ToString(k)] = []string***REMOVED***ToString(val)***REMOVED***
			***REMOVED***
		***REMOVED***
		return m, nil
	case map[interface***REMOVED******REMOVED***][]string:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = ToStringSlice(val)
		***REMOVED***
		return m, nil
	case map[interface***REMOVED******REMOVED***]string:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = ToStringSlice(val)
		***REMOVED***
		return m, nil
	case map[interface***REMOVED******REMOVED***][]interface***REMOVED******REMOVED***:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = ToStringSlice(val)
		***REMOVED***
		return m, nil
	case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
		for k, val := range v ***REMOVED***
			key, err := ToStringE(k)
			if err != nil ***REMOVED***
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			***REMOVED***
			value, err := ToStringSliceE(val)
			if err != nil ***REMOVED***
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			***REMOVED***
			m[key] = value
		***REMOVED***
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
	***REMOVED***
	return m, nil
***REMOVED***

// ToStringMapBoolE casts an interface to a map[string]bool type.
func ToStringMapBoolE(i interface***REMOVED******REMOVED***) (map[string]bool, error) ***REMOVED***
	var m = map[string]bool***REMOVED******REMOVED***

	switch v := i.(type) ***REMOVED***
	case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = ToBool(val)
		***REMOVED***
		return m, nil
	case map[string]interface***REMOVED******REMOVED***:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = ToBool(val)
		***REMOVED***
		return m, nil
	case map[string]bool:
		return v, nil
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]bool", i, i)
	***REMOVED***
***REMOVED***

// ToStringMapE casts an interface to a map[string]interface***REMOVED******REMOVED*** type.
func ToStringMapE(i interface***REMOVED******REMOVED***) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	var m = map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***

	switch v := i.(type) ***REMOVED***
	case map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***:
		for k, val := range v ***REMOVED***
			m[ToString(k)] = val
		***REMOVED***
		return m, nil
	case map[string]interface***REMOVED******REMOVED***:
		return v, nil
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]interface***REMOVED******REMOVED***", i, i)
	***REMOVED***
***REMOVED***

// ToSliceE casts an interface to a []interface***REMOVED******REMOVED*** type.
func ToSliceE(i interface***REMOVED******REMOVED***) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	var s []interface***REMOVED******REMOVED***

	switch v := i.(type) ***REMOVED***
	case []interface***REMOVED******REMOVED***:
		return append(s, v...), nil
	case []map[string]interface***REMOVED******REMOVED***:
		for _, u := range v ***REMOVED***
			s = append(s, u)
		***REMOVED***
		return s, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []interface***REMOVED******REMOVED***", i, i)
	***REMOVED***
***REMOVED***

// ToBoolSliceE casts an interface to a []bool type.
func ToBoolSliceE(i interface***REMOVED******REMOVED***) ([]bool, error) ***REMOVED***
	if i == nil ***REMOVED***
		return []bool***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
	***REMOVED***

	switch v := i.(type) ***REMOVED***
	case []bool:
		return v, nil
	***REMOVED***

	kind := reflect.TypeOf(i).Kind()
	switch kind ***REMOVED***
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]bool, s.Len())
		for j := 0; j < s.Len(); j++ ***REMOVED***
			val, err := ToBoolE(s.Index(j).Interface())
			if err != nil ***REMOVED***
				return []bool***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
			***REMOVED***
			a[j] = val
		***REMOVED***
		return a, nil
	default:
		return []bool***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
	***REMOVED***
***REMOVED***

// ToStringSliceE casts an interface to a []string type.
func ToStringSliceE(i interface***REMOVED******REMOVED***) ([]string, error) ***REMOVED***
	var a []string

	switch v := i.(type) ***REMOVED***
	case []interface***REMOVED******REMOVED***:
		for _, u := range v ***REMOVED***
			a = append(a, ToString(u))
		***REMOVED***
		return a, nil
	case []string:
		return v, nil
	case string:
		return strings.Fields(v), nil
	case interface***REMOVED******REMOVED***:
		str, err := ToStringE(v)
		if err != nil ***REMOVED***
			return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
		***REMOVED***
		return []string***REMOVED***str***REMOVED***, nil
	default:
		return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
	***REMOVED***
***REMOVED***

// ToIntSliceE casts an interface to a []int type.
func ToIntSliceE(i interface***REMOVED******REMOVED***) ([]int, error) ***REMOVED***
	if i == nil ***REMOVED***
		return []int***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	***REMOVED***

	switch v := i.(type) ***REMOVED***
	case []int:
		return v, nil
	***REMOVED***

	kind := reflect.TypeOf(i).Kind()
	switch kind ***REMOVED***
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int, s.Len())
		for j := 0; j < s.Len(); j++ ***REMOVED***
			val, err := ToIntE(s.Index(j).Interface())
			if err != nil ***REMOVED***
				return []int***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
			***REMOVED***
			a[j] = val
		***REMOVED***
		return a, nil
	default:
		return []int***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	***REMOVED***
***REMOVED***

// ToDurationSliceE casts an interface to a []time.Duration type.
func ToDurationSliceE(i interface***REMOVED******REMOVED***) ([]time.Duration, error) ***REMOVED***
	if i == nil ***REMOVED***
		return []time.Duration***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
	***REMOVED***

	switch v := i.(type) ***REMOVED***
	case []time.Duration:
		return v, nil
	***REMOVED***

	kind := reflect.TypeOf(i).Kind()
	switch kind ***REMOVED***
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]time.Duration, s.Len())
		for j := 0; j < s.Len(); j++ ***REMOVED***
			val, err := ToDurationE(s.Index(j).Interface())
			if err != nil ***REMOVED***
				return []time.Duration***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
			***REMOVED***
			a[j] = val
		***REMOVED***
		return a, nil
	default:
		return []time.Duration***REMOVED******REMOVED***, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
	***REMOVED***
***REMOVED***

// StringToDate attempts to parse a string into a time.Time type using a
// predefined list of formats.  If no suitable format is found, an error is
// returned.
func StringToDate(s string) (time.Time, error) ***REMOVED***
	return parseDateWith(s, []string***REMOVED***
		time.RFC3339,
		"2006-01-02T15:04:05", // iso8601 without timezone
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
		"2006-01-02",
		"02 Jan 2006",
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05Z07:00", // RFC3339 without T
		"2006-01-02 15:04:05",
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	***REMOVED***)
***REMOVED***

func parseDateWith(s string, dates []string) (d time.Time, e error) ***REMOVED***
	for _, dateType := range dates ***REMOVED***
		if d, e = time.Parse(dateType, s); e == nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return d, fmt.Errorf("unable to parse date: %s", s)
***REMOVED***
