package awsutil

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmespath/go-jmespath"
)

var indexRe = regexp.MustCompile(`(.+)\[(-?\d+)?\]$`)

// rValuesAtPath returns a slice of values found in value v. The values
// in v are explored recursively so all nested values are collected.
func rValuesAtPath(v interface***REMOVED******REMOVED***, path string, createPath, caseSensitive, nilTerm bool) []reflect.Value ***REMOVED***
	pathparts := strings.Split(path, "||")
	if len(pathparts) > 1 ***REMOVED***
		for _, pathpart := range pathparts ***REMOVED***
			vals := rValuesAtPath(v, pathpart, createPath, caseSensitive, nilTerm)
			if len(vals) > 0 ***REMOVED***
				return vals
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	values := []reflect.Value***REMOVED***reflect.Indirect(reflect.ValueOf(v))***REMOVED***
	components := strings.Split(path, ".")
	for len(values) > 0 && len(components) > 0 ***REMOVED***
		var index *int64
		var indexStar bool
		c := strings.TrimSpace(components[0])
		if c == "" ***REMOVED*** // no actual component, illegal syntax
			return nil
		***REMOVED*** else if caseSensitive && c != "*" && strings.ToLower(c[0:1]) == c[0:1] ***REMOVED***
			// TODO normalize case for user
			return nil // don't support unexported fields
		***REMOVED***

		// parse this component
		if m := indexRe.FindStringSubmatch(c); m != nil ***REMOVED***
			c = m[1]
			if m[2] == "" ***REMOVED***
				index = nil
				indexStar = true
			***REMOVED*** else ***REMOVED***
				i, _ := strconv.ParseInt(m[2], 10, 32)
				index = &i
				indexStar = false
			***REMOVED***
		***REMOVED***

		nextvals := []reflect.Value***REMOVED******REMOVED***
		for _, value := range values ***REMOVED***
			// pull component name out of struct member
			if value.Kind() != reflect.Struct ***REMOVED***
				continue
			***REMOVED***

			if c == "*" ***REMOVED*** // pull all members
				for i := 0; i < value.NumField(); i++ ***REMOVED***
					if f := reflect.Indirect(value.Field(i)); f.IsValid() ***REMOVED***
						nextvals = append(nextvals, f)
					***REMOVED***
				***REMOVED***
				continue
			***REMOVED***

			value = value.FieldByNameFunc(func(name string) bool ***REMOVED***
				if c == name ***REMOVED***
					return true
				***REMOVED*** else if !caseSensitive && strings.ToLower(name) == strings.ToLower(c) ***REMOVED***
					return true
				***REMOVED***
				return false
			***REMOVED***)

			if nilTerm && value.Kind() == reflect.Ptr && len(components[1:]) == 0 ***REMOVED***
				if !value.IsNil() ***REMOVED***
					value.Set(reflect.Zero(value.Type()))
				***REMOVED***
				return []reflect.Value***REMOVED***value***REMOVED***
			***REMOVED***

			if createPath && value.Kind() == reflect.Ptr && value.IsNil() ***REMOVED***
				// TODO if the value is the terminus it should not be created
				// if the value to be set to its position is nil.
				value.Set(reflect.New(value.Type().Elem()))
				value = value.Elem()
			***REMOVED*** else ***REMOVED***
				value = reflect.Indirect(value)
			***REMOVED***

			if value.Kind() == reflect.Slice || value.Kind() == reflect.Map ***REMOVED***
				if !createPath && value.IsNil() ***REMOVED***
					value = reflect.ValueOf(nil)
				***REMOVED***
			***REMOVED***

			if value.IsValid() ***REMOVED***
				nextvals = append(nextvals, value)
			***REMOVED***
		***REMOVED***
		values = nextvals

		if indexStar || index != nil ***REMOVED***
			nextvals = []reflect.Value***REMOVED******REMOVED***
			for _, valItem := range values ***REMOVED***
				value := reflect.Indirect(valItem)
				if value.Kind() != reflect.Slice ***REMOVED***
					continue
				***REMOVED***

				if indexStar ***REMOVED*** // grab all indices
					for i := 0; i < value.Len(); i++ ***REMOVED***
						idx := reflect.Indirect(value.Index(i))
						if idx.IsValid() ***REMOVED***
							nextvals = append(nextvals, idx)
						***REMOVED***
					***REMOVED***
					continue
				***REMOVED***

				// pull out index
				i := int(*index)
				if i >= value.Len() ***REMOVED*** // check out of bounds
					if createPath ***REMOVED***
						// TODO resize slice
					***REMOVED*** else ***REMOVED***
						continue
					***REMOVED***
				***REMOVED*** else if i < 0 ***REMOVED*** // support negative indexing
					i = value.Len() + i
				***REMOVED***
				value = reflect.Indirect(value.Index(i))

				if value.Kind() == reflect.Slice || value.Kind() == reflect.Map ***REMOVED***
					if !createPath && value.IsNil() ***REMOVED***
						value = reflect.ValueOf(nil)
					***REMOVED***
				***REMOVED***

				if value.IsValid() ***REMOVED***
					nextvals = append(nextvals, value)
				***REMOVED***
			***REMOVED***
			values = nextvals
		***REMOVED***

		components = components[1:]
	***REMOVED***
	return values
***REMOVED***

// ValuesAtPath returns a list of values at the case insensitive lexical
// path inside of a structure.
func ValuesAtPath(i interface***REMOVED******REMOVED***, path string) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	result, err := jmespath.Search(path, i)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	v := reflect.ValueOf(result)
	if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) ***REMOVED***
		return nil, nil
	***REMOVED***
	if s, ok := result.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return s, err
	***REMOVED***
	if v.Kind() == reflect.Map && v.Len() == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	if v.Kind() == reflect.Slice ***REMOVED***
		out := make([]interface***REMOVED******REMOVED***, v.Len())
		for i := 0; i < v.Len(); i++ ***REMOVED***
			out[i] = v.Index(i).Interface()
		***REMOVED***
		return out, nil
	***REMOVED***

	return []interface***REMOVED******REMOVED******REMOVED***result***REMOVED***, nil
***REMOVED***

// SetValueAtPath sets a value at the case insensitive lexical path inside
// of a structure.
func SetValueAtPath(i interface***REMOVED******REMOVED***, path string, v interface***REMOVED******REMOVED***) ***REMOVED***
	if rvals := rValuesAtPath(i, path, true, false, v == nil); rvals != nil ***REMOVED***
		for _, rval := range rvals ***REMOVED***
			if rval.Kind() == reflect.Ptr && rval.IsNil() ***REMOVED***
				continue
			***REMOVED***
			setValue(rval, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func setValue(dstVal reflect.Value, src interface***REMOVED******REMOVED***) ***REMOVED***
	if dstVal.Kind() == reflect.Ptr ***REMOVED***
		dstVal = reflect.Indirect(dstVal)
	***REMOVED***
	srcVal := reflect.ValueOf(src)

	if !srcVal.IsValid() ***REMOVED*** // src is literal nil
		if dstVal.CanAddr() ***REMOVED***
			// Convert to pointer so that pointer's value can be nil'ed
			//                     dstVal = dstVal.Addr()
		***REMOVED***
		dstVal.Set(reflect.Zero(dstVal.Type()))

	***REMOVED*** else if srcVal.Kind() == reflect.Ptr ***REMOVED***
		if srcVal.IsNil() ***REMOVED***
			srcVal = reflect.Zero(dstVal.Type())
		***REMOVED*** else ***REMOVED***
			srcVal = reflect.ValueOf(src).Elem()
		***REMOVED***
		dstVal.Set(srcVal)
	***REMOVED*** else ***REMOVED***
		dstVal.Set(srcVal)
	***REMOVED***

***REMOVED***
