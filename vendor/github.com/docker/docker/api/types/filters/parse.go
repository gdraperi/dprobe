/*Package filters provides tools for encoding a mapping of keys to a set of
multiple values.
*/
package filters

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types/versions"
)

// Args stores a mapping of keys to a set of multiple values.
type Args struct ***REMOVED***
	fields map[string]map[string]bool
***REMOVED***

// KeyValuePair are used to initialize a new Args
type KeyValuePair struct ***REMOVED***
	Key   string
	Value string
***REMOVED***

// Arg creates a new KeyValuePair for initializing Args
func Arg(key, value string) KeyValuePair ***REMOVED***
	return KeyValuePair***REMOVED***Key: key, Value: value***REMOVED***
***REMOVED***

// NewArgs returns a new Args populated with the initial args
func NewArgs(initialArgs ...KeyValuePair) Args ***REMOVED***
	args := Args***REMOVED***fields: map[string]map[string]bool***REMOVED******REMOVED******REMOVED***
	for _, arg := range initialArgs ***REMOVED***
		args.Add(arg.Key, arg.Value)
	***REMOVED***
	return args
***REMOVED***

// ParseFlag parses a key=value string and adds it to an Args.
//
// Deprecated: Use Args.Add()
func ParseFlag(arg string, prev Args) (Args, error) ***REMOVED***
	filters := prev
	if len(arg) == 0 ***REMOVED***
		return filters, nil
	***REMOVED***

	if !strings.Contains(arg, "=") ***REMOVED***
		return filters, ErrBadFormat
	***REMOVED***

	f := strings.SplitN(arg, "=", 2)

	name := strings.ToLower(strings.TrimSpace(f[0]))
	value := strings.TrimSpace(f[1])

	filters.Add(name, value)

	return filters, nil
***REMOVED***

// ErrBadFormat is an error returned when a filter is not in the form key=value
//
// Deprecated: this error will be removed in a future version
var ErrBadFormat = errors.New("bad format of filter (expected name=value)")

// ToParam encodes the Args as args JSON encoded string
//
// Deprecated: use ToJSON
func ToParam(a Args) (string, error) ***REMOVED***
	return ToJSON(a)
***REMOVED***

// MarshalJSON returns a JSON byte representation of the Args
func (args Args) MarshalJSON() ([]byte, error) ***REMOVED***
	if len(args.fields) == 0 ***REMOVED***
		return []byte***REMOVED******REMOVED***, nil
	***REMOVED***
	return json.Marshal(args.fields)
***REMOVED***

// ToJSON returns the Args as a JSON encoded string
func ToJSON(a Args) (string, error) ***REMOVED***
	if a.Len() == 0 ***REMOVED***
		return "", nil
	***REMOVED***
	buf, err := json.Marshal(a)
	return string(buf), err
***REMOVED***

// ToParamWithVersion encodes Args as a JSON string. If version is less than 1.22
// then the encoded format will use an older legacy format where the values are a
// list of strings, instead of a set.
//
// Deprecated: Use ToJSON
func ToParamWithVersion(version string, a Args) (string, error) ***REMOVED***
	if a.Len() == 0 ***REMOVED***
		return "", nil
	***REMOVED***

	if version != "" && versions.LessThan(version, "1.22") ***REMOVED***
		buf, err := json.Marshal(convertArgsToSlice(a.fields))
		return string(buf), err
	***REMOVED***

	return ToJSON(a)
***REMOVED***

// FromParam decodes a JSON encoded string into Args
//
// Deprecated: use FromJSON
func FromParam(p string) (Args, error) ***REMOVED***
	return FromJSON(p)
***REMOVED***

// FromJSON decodes a JSON encoded string into Args
func FromJSON(p string) (Args, error) ***REMOVED***
	args := NewArgs()

	if p == "" ***REMOVED***
		return args, nil
	***REMOVED***

	raw := []byte(p)
	err := json.Unmarshal(raw, &args)
	if err == nil ***REMOVED***
		return args, nil
	***REMOVED***

	// Fallback to parsing arguments in the legacy slice format
	deprecated := map[string][]string***REMOVED******REMOVED***
	if legacyErr := json.Unmarshal(raw, &deprecated); legacyErr != nil ***REMOVED***
		return args, err
	***REMOVED***

	args.fields = deprecatedArgs(deprecated)
	return args, nil
***REMOVED***

// UnmarshalJSON populates the Args from JSON encode bytes
func (args Args) UnmarshalJSON(raw []byte) error ***REMOVED***
	if len(raw) == 0 ***REMOVED***
		return nil
	***REMOVED***
	return json.Unmarshal(raw, &args.fields)
***REMOVED***

// Get returns the list of values associated with the key
func (args Args) Get(key string) []string ***REMOVED***
	values := args.fields[key]
	if values == nil ***REMOVED***
		return make([]string, 0)
	***REMOVED***
	slice := make([]string, 0, len(values))
	for key := range values ***REMOVED***
		slice = append(slice, key)
	***REMOVED***
	return slice
***REMOVED***

// Add a new value to the set of values
func (args Args) Add(key, value string) ***REMOVED***
	if _, ok := args.fields[key]; ok ***REMOVED***
		args.fields[key][value] = true
	***REMOVED*** else ***REMOVED***
		args.fields[key] = map[string]bool***REMOVED***value: true***REMOVED***
	***REMOVED***
***REMOVED***

// Del removes a value from the set
func (args Args) Del(key, value string) ***REMOVED***
	if _, ok := args.fields[key]; ok ***REMOVED***
		delete(args.fields[key], value)
		if len(args.fields[key]) == 0 ***REMOVED***
			delete(args.fields, key)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Len returns the number of keys in the mapping
func (args Args) Len() int ***REMOVED***
	return len(args.fields)
***REMOVED***

// MatchKVList returns true if all the pairs in sources exist as key=value
// pairs in the mapping at key, or if there are no values at key.
func (args Args) MatchKVList(key string, sources map[string]string) bool ***REMOVED***
	fieldValues := args.fields[key]

	//do not filter if there is no filter set or cannot determine filter
	if len(fieldValues) == 0 ***REMOVED***
		return true
	***REMOVED***

	if len(sources) == 0 ***REMOVED***
		return false
	***REMOVED***

	for value := range fieldValues ***REMOVED***
		testKV := strings.SplitN(value, "=", 2)

		v, ok := sources[testKV[0]]
		if !ok ***REMOVED***
			return false
		***REMOVED***
		if len(testKV) == 2 && testKV[1] != v ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// Match returns true if any of the values at key match the source string
func (args Args) Match(field, source string) bool ***REMOVED***
	if args.ExactMatch(field, source) ***REMOVED***
		return true
	***REMOVED***

	fieldValues := args.fields[field]
	for name2match := range fieldValues ***REMOVED***
		match, err := regexp.MatchString(name2match, source)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		if match ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// ExactMatch returns true if the source matches exactly one of the values.
func (args Args) ExactMatch(key, source string) bool ***REMOVED***
	fieldValues, ok := args.fields[key]
	//do not filter if there is no filter set or cannot determine filter
	if !ok || len(fieldValues) == 0 ***REMOVED***
		return true
	***REMOVED***

	// try to match full name value to avoid O(N) regular expression matching
	return fieldValues[source]
***REMOVED***

// UniqueExactMatch returns true if there is only one value and the source
// matches exactly the value.
func (args Args) UniqueExactMatch(key, source string) bool ***REMOVED***
	fieldValues := args.fields[key]
	//do not filter if there is no filter set or cannot determine filter
	if len(fieldValues) == 0 ***REMOVED***
		return true
	***REMOVED***
	if len(args.fields[key]) != 1 ***REMOVED***
		return false
	***REMOVED***

	// try to match full name value to avoid O(N) regular expression matching
	return fieldValues[source]
***REMOVED***

// FuzzyMatch returns true if the source matches exactly one value,  or the
// source has one of the values as a prefix.
func (args Args) FuzzyMatch(key, source string) bool ***REMOVED***
	if args.ExactMatch(key, source) ***REMOVED***
		return true
	***REMOVED***

	fieldValues := args.fields[key]
	for prefix := range fieldValues ***REMOVED***
		if strings.HasPrefix(source, prefix) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Include returns true if the key exists in the mapping
//
// Deprecated: use Contains
func (args Args) Include(field string) bool ***REMOVED***
	_, ok := args.fields[field]
	return ok
***REMOVED***

// Contains returns true if the key exists in the mapping
func (args Args) Contains(field string) bool ***REMOVED***
	_, ok := args.fields[field]
	return ok
***REMOVED***

type invalidFilter string

func (e invalidFilter) Error() string ***REMOVED***
	return "Invalid filter '" + string(e) + "'"
***REMOVED***

func (invalidFilter) InvalidParameter() ***REMOVED******REMOVED***

// Validate compared the set of accepted keys against the keys in the mapping.
// An error is returned if any mapping keys are not in the accepted set.
func (args Args) Validate(accepted map[string]bool) error ***REMOVED***
	for name := range args.fields ***REMOVED***
		if !accepted[name] ***REMOVED***
			return invalidFilter(name)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// WalkValues iterates over the list of values for a key in the mapping and calls
// op() for each value. If op returns an error the iteration stops and the
// error is returned.
func (args Args) WalkValues(field string, op func(value string) error) error ***REMOVED***
	if _, ok := args.fields[field]; !ok ***REMOVED***
		return nil
	***REMOVED***
	for v := range args.fields[field] ***REMOVED***
		if err := op(v); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func deprecatedArgs(d map[string][]string) map[string]map[string]bool ***REMOVED***
	m := map[string]map[string]bool***REMOVED******REMOVED***
	for k, v := range d ***REMOVED***
		values := map[string]bool***REMOVED******REMOVED***
		for _, vv := range v ***REMOVED***
			values[vv] = true
		***REMOVED***
		m[k] = values
	***REMOVED***
	return m
***REMOVED***

func convertArgsToSlice(f map[string]map[string]bool) map[string][]string ***REMOVED***
	m := map[string][]string***REMOVED******REMOVED***
	for k, v := range f ***REMOVED***
		values := []string***REMOVED******REMOVED***
		for kk := range v ***REMOVED***
			if v[kk] ***REMOVED***
				values = append(values, kk)
			***REMOVED***
		***REMOVED***
		m[k] = values
	***REMOVED***
	return m
***REMOVED***
