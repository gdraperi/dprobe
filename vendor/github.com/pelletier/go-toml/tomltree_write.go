package toml

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// encodes a string to a TOML-compliant string value
func encodeTomlString(value string) string ***REMOVED***
	var b bytes.Buffer

	for _, rr := range value ***REMOVED***
		switch rr ***REMOVED***
		case '\b':
			b.WriteString(`\b`)
		case '\t':
			b.WriteString(`\t`)
		case '\n':
			b.WriteString(`\n`)
		case '\f':
			b.WriteString(`\f`)
		case '\r':
			b.WriteString(`\r`)
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		default:
			intRr := uint16(rr)
			if intRr < 0x001F ***REMOVED***
				b.WriteString(fmt.Sprintf("\\u%0.4X", intRr))
			***REMOVED*** else ***REMOVED***
				b.WriteRune(rr)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return b.String()
***REMOVED***

func tomlValueStringRepresentation(v interface***REMOVED******REMOVED***, indent string, arraysOneElementPerLine bool) (string, error) ***REMOVED***
	switch value := v.(type) ***REMOVED***
	case uint64:
		return strconv.FormatUint(value, 10), nil
	case int64:
		return strconv.FormatInt(value, 10), nil
	case float64:
		// Ensure a round float does contain a decimal point. Otherwise feeding
		// the output back to the parser would convert to an integer.
		if math.Trunc(value) == value ***REMOVED***
			return strings.ToLower(strconv.FormatFloat(value, 'f', 1, 32)), nil
		***REMOVED***
		return strings.ToLower(strconv.FormatFloat(value, 'f', -1, 32)), nil
	case string:
		return "\"" + encodeTomlString(value) + "\"", nil
	case []byte:
		b, _ := v.([]byte)
		return tomlValueStringRepresentation(string(b), indent, arraysOneElementPerLine)
	case bool:
		if value ***REMOVED***
			return "true", nil
		***REMOVED***
		return "false", nil
	case time.Time:
		return value.Format(time.RFC3339), nil
	case nil:
		return "", nil
	***REMOVED***

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Slice ***REMOVED***
		var values []string
		for i := 0; i < rv.Len(); i++ ***REMOVED***
			item := rv.Index(i).Interface()
			itemRepr, err := tomlValueStringRepresentation(item, indent, arraysOneElementPerLine)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			values = append(values, itemRepr)
		***REMOVED***
		if arraysOneElementPerLine && len(values) > 1 ***REMOVED***
			stringBuffer := bytes.Buffer***REMOVED******REMOVED***
			valueIndent := indent + `  ` // TODO: move that to a shared encoder state

			stringBuffer.WriteString("[\n")

			for i, value := range values ***REMOVED***
				stringBuffer.WriteString(valueIndent)
				stringBuffer.WriteString(value)
				if i != len(values)-1 ***REMOVED***
					stringBuffer.WriteString(`,`)
				***REMOVED***
				stringBuffer.WriteString("\n")
			***REMOVED***

			stringBuffer.WriteString(indent + "]")

			return stringBuffer.String(), nil
		***REMOVED***
		return "[" + strings.Join(values, ",") + "]", nil
	***REMOVED***
	return "", fmt.Errorf("unsupported value type %T: %v", v, v)
***REMOVED***

func (t *Tree) writeTo(w io.Writer, indent, keyspace string, bytesCount int64, arraysOneElementPerLine bool) (int64, error) ***REMOVED***
	simpleValuesKeys := make([]string, 0)
	complexValuesKeys := make([]string, 0)

	for k := range t.values ***REMOVED***
		v := t.values[k]
		switch v.(type) ***REMOVED***
		case *Tree, []*Tree:
			complexValuesKeys = append(complexValuesKeys, k)
		default:
			simpleValuesKeys = append(simpleValuesKeys, k)
		***REMOVED***
	***REMOVED***

	sort.Strings(simpleValuesKeys)
	sort.Strings(complexValuesKeys)

	for _, k := range simpleValuesKeys ***REMOVED***
		v, ok := t.values[k].(*tomlValue)
		if !ok ***REMOVED***
			return bytesCount, fmt.Errorf("invalid value type at %s: %T", k, t.values[k])
		***REMOVED***

		repr, err := tomlValueStringRepresentation(v.value, indent, arraysOneElementPerLine)
		if err != nil ***REMOVED***
			return bytesCount, err
		***REMOVED***

		if v.comment != "" ***REMOVED***
			comment := strings.Replace(v.comment, "\n", "\n"+indent+"#", -1)
			start := "# "
			if strings.HasPrefix(comment, "#") ***REMOVED***
				start = ""
			***REMOVED***
			writtenBytesCountComment, errc := writeStrings(w, "\n", indent, start, comment, "\n")
			bytesCount += int64(writtenBytesCountComment)
			if errc != nil ***REMOVED***
				return bytesCount, errc
			***REMOVED***
		***REMOVED***

		var commented string
		if v.commented ***REMOVED***
			commented = "# "
		***REMOVED***
		writtenBytesCount, err := writeStrings(w, indent, commented, k, " = ", repr, "\n")
		bytesCount += int64(writtenBytesCount)
		if err != nil ***REMOVED***
			return bytesCount, err
		***REMOVED***
	***REMOVED***

	for _, k := range complexValuesKeys ***REMOVED***
		v := t.values[k]

		combinedKey := k
		if keyspace != "" ***REMOVED***
			combinedKey = keyspace + "." + combinedKey
		***REMOVED***
		var commented string
		if t.commented ***REMOVED***
			commented = "# "
		***REMOVED***

		switch node := v.(type) ***REMOVED***
		// node has to be of those two types given how keys are sorted above
		case *Tree:
			tv, ok := t.values[k].(*Tree)
			if !ok ***REMOVED***
				return bytesCount, fmt.Errorf("invalid value type at %s: %T", k, t.values[k])
			***REMOVED***
			if tv.comment != "" ***REMOVED***
				comment := strings.Replace(tv.comment, "\n", "\n"+indent+"#", -1)
				start := "# "
				if strings.HasPrefix(comment, "#") ***REMOVED***
					start = ""
				***REMOVED***
				writtenBytesCountComment, errc := writeStrings(w, "\n", indent, start, comment)
				bytesCount += int64(writtenBytesCountComment)
				if errc != nil ***REMOVED***
					return bytesCount, errc
				***REMOVED***
			***REMOVED***
			writtenBytesCount, err := writeStrings(w, "\n", indent, commented, "[", combinedKey, "]\n")
			bytesCount += int64(writtenBytesCount)
			if err != nil ***REMOVED***
				return bytesCount, err
			***REMOVED***
			bytesCount, err = node.writeTo(w, indent+"  ", combinedKey, bytesCount, arraysOneElementPerLine)
			if err != nil ***REMOVED***
				return bytesCount, err
			***REMOVED***
		case []*Tree:
			for _, subTree := range node ***REMOVED***
				writtenBytesCount, err := writeStrings(w, "\n", indent, commented, "[[", combinedKey, "]]\n")
				bytesCount += int64(writtenBytesCount)
				if err != nil ***REMOVED***
					return bytesCount, err
				***REMOVED***

				bytesCount, err = subTree.writeTo(w, indent+"  ", combinedKey, bytesCount, arraysOneElementPerLine)
				if err != nil ***REMOVED***
					return bytesCount, err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return bytesCount, nil
***REMOVED***

func writeStrings(w io.Writer, s ...string) (int, error) ***REMOVED***
	var n int
	for i := range s ***REMOVED***
		b, err := io.WriteString(w, s[i])
		n += b
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
	***REMOVED***
	return n, nil
***REMOVED***

// WriteTo encode the Tree as Toml and writes it to the writer w.
// Returns the number of bytes written in case of success, or an error if anything happened.
func (t *Tree) WriteTo(w io.Writer) (int64, error) ***REMOVED***
	return t.writeTo(w, "", "", 0, false)
***REMOVED***

// ToTomlString generates a human-readable representation of the current tree.
// Output spans multiple lines, and is suitable for ingest by a TOML parser.
// If the conversion cannot be performed, ToString returns a non-nil error.
func (t *Tree) ToTomlString() (string, error) ***REMOVED***
	var buf bytes.Buffer
	_, err := t.WriteTo(&buf)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return buf.String(), nil
***REMOVED***

// String generates a human-readable representation of the current tree.
// Alias of ToString. Present to implement the fmt.Stringer interface.
func (t *Tree) String() string ***REMOVED***
	result, _ := t.ToTomlString()
	return result
***REMOVED***

// ToMap recursively generates a representation of the tree using Go built-in structures.
// The following types are used:
//
//	* bool
//	* float64
//	* int64
//	* string
//	* uint64
//	* time.Time
//	* map[string]interface***REMOVED******REMOVED*** (where interface***REMOVED******REMOVED*** is any of this list)
//	* []interface***REMOVED******REMOVED*** (where interface***REMOVED******REMOVED*** is any of this list)
func (t *Tree) ToMap() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	result := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***

	for k, v := range t.values ***REMOVED***
		switch node := v.(type) ***REMOVED***
		case []*Tree:
			var array []interface***REMOVED******REMOVED***
			for _, item := range node ***REMOVED***
				array = append(array, item.ToMap())
			***REMOVED***
			result[k] = array
		case *Tree:
			result[k] = node.ToMap()
		case *tomlValue:
			result[k] = node.value
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***
