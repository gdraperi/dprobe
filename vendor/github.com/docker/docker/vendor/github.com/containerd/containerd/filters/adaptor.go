package filters

// Adaptor specifies the mapping of fieldpaths to a type. For the given field
// path, the value and whether it is present should be returned. The mapping of
// the fieldpath to a field is deferred to the adaptor implementation, but
// should generally follow protobuf field path/mask semantics.
type Adaptor interface ***REMOVED***
	Field(fieldpath []string) (value string, present bool)
***REMOVED***

// AdapterFunc allows implementation specific matching of fieldpaths
type AdapterFunc func(fieldpath []string) (string, bool)

// Field returns the field name and true if it exists
func (fn AdapterFunc) Field(fieldpath []string) (string, bool) ***REMOVED***
	return fn(fieldpath)
***REMOVED***
