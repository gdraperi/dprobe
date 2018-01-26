package opts

import (
	"fmt"
	"net"
	"path"
	"regexp"
	"strings"

	units "github.com/docker/go-units"
)

var (
	alphaRegexp  = regexp.MustCompile(`[a-zA-Z]`)
	domainRegexp = regexp.MustCompile(`^(:?(:?[a-zA-Z0-9]|(:?[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9]))(:?\.(:?[a-zA-Z0-9]|(:?[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])))*)\.?\s*$`)
)

// ListOpts holds a list of values and a validation function.
type ListOpts struct ***REMOVED***
	values    *[]string
	validator ValidatorFctType
***REMOVED***

// NewListOpts creates a new ListOpts with the specified validator.
func NewListOpts(validator ValidatorFctType) ListOpts ***REMOVED***
	var values []string
	return *NewListOptsRef(&values, validator)
***REMOVED***

// NewListOptsRef creates a new ListOpts with the specified values and validator.
func NewListOptsRef(values *[]string, validator ValidatorFctType) *ListOpts ***REMOVED***
	return &ListOpts***REMOVED***
		values:    values,
		validator: validator,
	***REMOVED***
***REMOVED***

func (opts *ListOpts) String() string ***REMOVED***
	if len(*opts.values) == 0 ***REMOVED***
		return ""
	***REMOVED***
	return fmt.Sprintf("%v", *opts.values)
***REMOVED***

// Set validates if needed the input value and adds it to the
// internal slice.
func (opts *ListOpts) Set(value string) error ***REMOVED***
	if opts.validator != nil ***REMOVED***
		v, err := opts.validator(value)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		value = v
	***REMOVED***
	(*opts.values) = append((*opts.values), value)
	return nil
***REMOVED***

// Delete removes the specified element from the slice.
func (opts *ListOpts) Delete(key string) ***REMOVED***
	for i, k := range *opts.values ***REMOVED***
		if k == key ***REMOVED***
			(*opts.values) = append((*opts.values)[:i], (*opts.values)[i+1:]...)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// GetMap returns the content of values in a map in order to avoid
// duplicates.
func (opts *ListOpts) GetMap() map[string]struct***REMOVED******REMOVED*** ***REMOVED***
	ret := make(map[string]struct***REMOVED******REMOVED***)
	for _, k := range *opts.values ***REMOVED***
		ret[k] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	return ret
***REMOVED***

// GetAll returns the values of slice.
func (opts *ListOpts) GetAll() []string ***REMOVED***
	return (*opts.values)
***REMOVED***

// GetAllOrEmpty returns the values of the slice
// or an empty slice when there are no values.
func (opts *ListOpts) GetAllOrEmpty() []string ***REMOVED***
	v := *opts.values
	if v == nil ***REMOVED***
		return make([]string, 0)
	***REMOVED***
	return v
***REMOVED***

// Get checks the existence of the specified key.
func (opts *ListOpts) Get(key string) bool ***REMOVED***
	for _, k := range *opts.values ***REMOVED***
		if k == key ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Len returns the amount of element in the slice.
func (opts *ListOpts) Len() int ***REMOVED***
	return len((*opts.values))
***REMOVED***

// Type returns a string name for this Option type
func (opts *ListOpts) Type() string ***REMOVED***
	return "list"
***REMOVED***

// WithValidator returns the ListOpts with validator set.
func (opts *ListOpts) WithValidator(validator ValidatorFctType) *ListOpts ***REMOVED***
	opts.validator = validator
	return opts
***REMOVED***

// NamedOption is an interface that list and map options
// with names implement.
type NamedOption interface ***REMOVED***
	Name() string
***REMOVED***

// NamedListOpts is a ListOpts with a configuration name.
// This struct is useful to keep reference to the assigned
// field name in the internal configuration struct.
type NamedListOpts struct ***REMOVED***
	name string
	ListOpts
***REMOVED***

var _ NamedOption = &NamedListOpts***REMOVED******REMOVED***

// NewNamedListOptsRef creates a reference to a new NamedListOpts struct.
func NewNamedListOptsRef(name string, values *[]string, validator ValidatorFctType) *NamedListOpts ***REMOVED***
	return &NamedListOpts***REMOVED***
		name:     name,
		ListOpts: *NewListOptsRef(values, validator),
	***REMOVED***
***REMOVED***

// Name returns the name of the NamedListOpts in the configuration.
func (o *NamedListOpts) Name() string ***REMOVED***
	return o.name
***REMOVED***

// MapOpts holds a map of values and a validation function.
type MapOpts struct ***REMOVED***
	values    map[string]string
	validator ValidatorFctType
***REMOVED***

// Set validates if needed the input value and add it to the
// internal map, by splitting on '='.
func (opts *MapOpts) Set(value string) error ***REMOVED***
	if opts.validator != nil ***REMOVED***
		v, err := opts.validator(value)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		value = v
	***REMOVED***
	vals := strings.SplitN(value, "=", 2)
	if len(vals) == 1 ***REMOVED***
		(opts.values)[vals[0]] = ""
	***REMOVED*** else ***REMOVED***
		(opts.values)[vals[0]] = vals[1]
	***REMOVED***
	return nil
***REMOVED***

// GetAll returns the values of MapOpts as a map.
func (opts *MapOpts) GetAll() map[string]string ***REMOVED***
	return opts.values
***REMOVED***

func (opts *MapOpts) String() string ***REMOVED***
	return fmt.Sprintf("%v", opts.values)
***REMOVED***

// Type returns a string name for this Option type
func (opts *MapOpts) Type() string ***REMOVED***
	return "map"
***REMOVED***

// NewMapOpts creates a new MapOpts with the specified map of values and a validator.
func NewMapOpts(values map[string]string, validator ValidatorFctType) *MapOpts ***REMOVED***
	if values == nil ***REMOVED***
		values = make(map[string]string)
	***REMOVED***
	return &MapOpts***REMOVED***
		values:    values,
		validator: validator,
	***REMOVED***
***REMOVED***

// NamedMapOpts is a MapOpts struct with a configuration name.
// This struct is useful to keep reference to the assigned
// field name in the internal configuration struct.
type NamedMapOpts struct ***REMOVED***
	name string
	MapOpts
***REMOVED***

var _ NamedOption = &NamedMapOpts***REMOVED******REMOVED***

// NewNamedMapOpts creates a reference to a new NamedMapOpts struct.
func NewNamedMapOpts(name string, values map[string]string, validator ValidatorFctType) *NamedMapOpts ***REMOVED***
	return &NamedMapOpts***REMOVED***
		name:    name,
		MapOpts: *NewMapOpts(values, validator),
	***REMOVED***
***REMOVED***

// Name returns the name of the NamedMapOpts in the configuration.
func (o *NamedMapOpts) Name() string ***REMOVED***
	return o.name
***REMOVED***

// ValidatorFctType defines a validator function that returns a validated string and/or an error.
type ValidatorFctType func(val string) (string, error)

// ValidatorFctListType defines a validator function that returns a validated list of string and/or an error
type ValidatorFctListType func(val string) ([]string, error)

// ValidateIPAddress validates an Ip address.
func ValidateIPAddress(val string) (string, error) ***REMOVED***
	var ip = net.ParseIP(strings.TrimSpace(val))
	if ip != nil ***REMOVED***
		return ip.String(), nil
	***REMOVED***
	return "", fmt.Errorf("%s is not an ip address", val)
***REMOVED***

// ValidateDNSSearch validates domain for resolvconf search configuration.
// A zero length domain is represented by a dot (.).
func ValidateDNSSearch(val string) (string, error) ***REMOVED***
	if val = strings.Trim(val, " "); val == "." ***REMOVED***
		return val, nil
	***REMOVED***
	return validateDomain(val)
***REMOVED***

func validateDomain(val string) (string, error) ***REMOVED***
	if alphaRegexp.FindString(val) == "" ***REMOVED***
		return "", fmt.Errorf("%s is not a valid domain", val)
	***REMOVED***
	ns := domainRegexp.FindSubmatch([]byte(val))
	if len(ns) > 0 && len(ns[1]) < 255 ***REMOVED***
		return string(ns[1]), nil
	***REMOVED***
	return "", fmt.Errorf("%s is not a valid domain", val)
***REMOVED***

// ValidateLabel validates that the specified string is a valid label, and returns it.
// Labels are in the form on key=value.
func ValidateLabel(val string) (string, error) ***REMOVED***
	if strings.Count(val, "=") < 1 ***REMOVED***
		return "", fmt.Errorf("bad attribute format: %s", val)
	***REMOVED***
	return val, nil
***REMOVED***

// ValidateSingleGenericResource validates that a single entry in the
// generic resource list is valid.
// i.e 'GPU=UID1' is valid however 'GPU:UID1' or 'UID1' isn't
func ValidateSingleGenericResource(val string) (string, error) ***REMOVED***
	if strings.Count(val, "=") < 1 ***REMOVED***
		return "", fmt.Errorf("invalid node-generic-resource format `%s` expected `name=value`", val)
	***REMOVED***
	return val, nil
***REMOVED***

// ParseLink parses and validates the specified string as a link format (name:alias)
func ParseLink(val string) (string, string, error) ***REMOVED***
	if val == "" ***REMOVED***
		return "", "", fmt.Errorf("empty string specified for links")
	***REMOVED***
	arr := strings.Split(val, ":")
	if len(arr) > 2 ***REMOVED***
		return "", "", fmt.Errorf("bad format for links: %s", val)
	***REMOVED***
	if len(arr) == 1 ***REMOVED***
		return val, val, nil
	***REMOVED***
	// This is kept because we can actually get a HostConfig with links
	// from an already created container and the format is not `foo:bar`
	// but `/foo:/c1/bar`
	if strings.HasPrefix(arr[0], "/") ***REMOVED***
		_, alias := path.Split(arr[1])
		return arr[0][1:], alias, nil
	***REMOVED***
	return arr[0], arr[1], nil
***REMOVED***

// MemBytes is a type for human readable memory bytes (like 128M, 2g, etc)
type MemBytes int64

// String returns the string format of the human readable memory bytes
func (m *MemBytes) String() string ***REMOVED***
	// NOTE: In spf13/pflag/flag.go, "0" is considered as "zero value" while "0 B" is not.
	// We return "0" in case value is 0 here so that the default value is hidden.
	// (Sometimes "default 0 B" is actually misleading)
	if m.Value() != 0 ***REMOVED***
		return units.BytesSize(float64(m.Value()))
	***REMOVED***
	return "0"
***REMOVED***

// Set sets the value of the MemBytes by passing a string
func (m *MemBytes) Set(value string) error ***REMOVED***
	val, err := units.RAMInBytes(value)
	*m = MemBytes(val)
	return err
***REMOVED***

// Type returns the type
func (m *MemBytes) Type() string ***REMOVED***
	return "bytes"
***REMOVED***

// Value returns the value in int64
func (m *MemBytes) Value() int64 ***REMOVED***
	return int64(*m)
***REMOVED***

// UnmarshalJSON is the customized unmarshaler for MemBytes
func (m *MemBytes) UnmarshalJSON(s []byte) error ***REMOVED***
	if len(s) <= 2 || s[0] != '"' || s[len(s)-1] != '"' ***REMOVED***
		return fmt.Errorf("invalid size: %q", s)
	***REMOVED***
	val, err := units.RAMInBytes(string(s[1 : len(s)-1]))
	*m = MemBytes(val)
	return err
***REMOVED***
