package awsutil

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// StringValue returns the string representation of a value.
func StringValue(i interface***REMOVED******REMOVED***) string ***REMOVED***
	var buf bytes.Buffer
	stringValue(reflect.ValueOf(i), 0, &buf)
	return buf.String()
***REMOVED***

func stringValue(v reflect.Value, indent int, buf *bytes.Buffer) ***REMOVED***
	for v.Kind() == reflect.Ptr ***REMOVED***
		v = v.Elem()
	***REMOVED***

	switch v.Kind() ***REMOVED***
	case reflect.Struct:
		buf.WriteString("***REMOVED***\n")

		names := []string***REMOVED******REMOVED***
		for i := 0; i < v.Type().NumField(); i++ ***REMOVED***
			name := v.Type().Field(i).Name
			f := v.Field(i)
			if name[0:1] == strings.ToLower(name[0:1]) ***REMOVED***
				continue // ignore unexported fields
			***REMOVED***
			if (f.Kind() == reflect.Ptr || f.Kind() == reflect.Slice) && f.IsNil() ***REMOVED***
				continue // ignore unset fields
			***REMOVED***
			names = append(names, name)
		***REMOVED***

		for i, n := range names ***REMOVED***
			val := v.FieldByName(n)
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(n + ": ")
			stringValue(val, indent+2, buf)

			if i < len(names)-1 ***REMOVED***
				buf.WriteString(",\n")
			***REMOVED***
		***REMOVED***

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "***REMOVED***")
	case reflect.Slice:
		nl, id, id2 := "", "", ""
		if v.Len() > 3 ***REMOVED***
			nl, id, id2 = "\n", strings.Repeat(" ", indent), strings.Repeat(" ", indent+2)
		***REMOVED***
		buf.WriteString("[" + nl)
		for i := 0; i < v.Len(); i++ ***REMOVED***
			buf.WriteString(id2)
			stringValue(v.Index(i), indent+2, buf)

			if i < v.Len()-1 ***REMOVED***
				buf.WriteString("," + nl)
			***REMOVED***
		***REMOVED***

		buf.WriteString(nl + id + "]")
	case reflect.Map:
		buf.WriteString("***REMOVED***\n")

		for i, k := range v.MapKeys() ***REMOVED***
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(k.String() + ": ")
			stringValue(v.MapIndex(k), indent+2, buf)

			if i < v.Len()-1 ***REMOVED***
				buf.WriteString(",\n")
			***REMOVED***
		***REMOVED***

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "***REMOVED***")
	default:
		format := "%v"
		switch v.Interface().(type) ***REMOVED***
		case string:
			format = "%q"
		***REMOVED***
		fmt.Fprintf(buf, format, v.Interface())
	***REMOVED***
***REMOVED***
