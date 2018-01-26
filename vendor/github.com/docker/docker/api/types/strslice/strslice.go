package strslice

import "encoding/json"

// StrSlice represents a string or an array of strings.
// We need to override the json decoder to accept both options.
type StrSlice []string

// UnmarshalJSON decodes the byte slice whether it's a string or an array of
// strings. This method is needed to implement json.Unmarshaler.
func (e *StrSlice) UnmarshalJSON(b []byte) error ***REMOVED***
	if len(b) == 0 ***REMOVED***
		// With no input, we preserve the existing value by returning nil and
		// leaving the target alone. This allows defining default values for
		// the type.
		return nil
	***REMOVED***

	p := make([]string, 0, 1)
	if err := json.Unmarshal(b, &p); err != nil ***REMOVED***
		var s string
		if err := json.Unmarshal(b, &s); err != nil ***REMOVED***
			return err
		***REMOVED***
		p = append(p, s)
	***REMOVED***

	*e = p
	return nil
***REMOVED***
