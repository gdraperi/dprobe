package opts

// QuotedString is a string that may have extra quotes around the value. The
// quotes are stripped from the value.
type QuotedString struct ***REMOVED***
	value *string
***REMOVED***

// Set sets a new value
func (s *QuotedString) Set(val string) error ***REMOVED***
	*s.value = trimQuotes(val)
	return nil
***REMOVED***

// Type returns the type of the value
func (s *QuotedString) Type() string ***REMOVED***
	return "string"
***REMOVED***

func (s *QuotedString) String() string ***REMOVED***
	return *s.value
***REMOVED***

func trimQuotes(value string) string ***REMOVED***
	lastIndex := len(value) - 1
	for _, char := range []byte***REMOVED***'\'', '"'***REMOVED*** ***REMOVED***
		if value[0] == char && value[lastIndex] == char ***REMOVED***
			return value[1:lastIndex]
		***REMOVED***
	***REMOVED***
	return value
***REMOVED***

// NewQuotedString returns a new quoted string option
func NewQuotedString(value *string) *QuotedString ***REMOVED***
	return &QuotedString***REMOVED***value: value***REMOVED***
***REMOVED***
