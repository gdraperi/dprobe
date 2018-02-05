// Parsing keys handling both bare and quoted keys.

package toml

import (
	"bytes"
	"errors"
	"fmt"
	"unicode"
)

// Convert the bare key group string to an array.
// The input supports double quotation to allow "." inside the key name,
// but escape sequences are not supported. Lexers must unescape them beforehand.
func parseKey(key string) ([]string, error) ***REMOVED***
	groups := []string***REMOVED******REMOVED***
	var buffer bytes.Buffer
	inQuotes := false
	wasInQuotes := false
	ignoreSpace := true
	expectDot := false

	for _, char := range key ***REMOVED***
		if ignoreSpace ***REMOVED***
			if char == ' ' ***REMOVED***
				continue
			***REMOVED***
			ignoreSpace = false
		***REMOVED***
		switch char ***REMOVED***
		case '"':
			if inQuotes ***REMOVED***
				groups = append(groups, buffer.String())
				buffer.Reset()
				wasInQuotes = true
			***REMOVED***
			inQuotes = !inQuotes
			expectDot = false
		case '.':
			if inQuotes ***REMOVED***
				buffer.WriteRune(char)
			***REMOVED*** else ***REMOVED***
				if !wasInQuotes ***REMOVED***
					if buffer.Len() == 0 ***REMOVED***
						return nil, errors.New("empty table key")
					***REMOVED***
					groups = append(groups, buffer.String())
					buffer.Reset()
				***REMOVED***
				ignoreSpace = true
				expectDot = false
				wasInQuotes = false
			***REMOVED***
		case ' ':
			if inQuotes ***REMOVED***
				buffer.WriteRune(char)
			***REMOVED*** else ***REMOVED***
				expectDot = true
			***REMOVED***
		default:
			if !inQuotes && !isValidBareChar(char) ***REMOVED***
				return nil, fmt.Errorf("invalid bare character: %c", char)
			***REMOVED***
			if !inQuotes && expectDot ***REMOVED***
				return nil, errors.New("what?")
			***REMOVED***
			buffer.WriteRune(char)
			expectDot = false
		***REMOVED***
	***REMOVED***
	if inQuotes ***REMOVED***
		return nil, errors.New("mismatched quotes")
	***REMOVED***
	if buffer.Len() > 0 ***REMOVED***
		groups = append(groups, buffer.String())
	***REMOVED***
	if len(groups) == 0 ***REMOVED***
		return nil, errors.New("empty key")
	***REMOVED***
	return groups, nil
***REMOVED***

func isValidBareChar(r rune) bool ***REMOVED***
	return isAlphanumeric(r) || r == '-' || unicode.IsNumber(r)
***REMOVED***
