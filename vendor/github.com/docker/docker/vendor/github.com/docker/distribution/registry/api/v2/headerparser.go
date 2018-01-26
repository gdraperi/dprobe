package v2

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	// according to rfc7230
	reToken            = regexp.MustCompile(`^[^"(),/:;<=>?@[\]***REMOVED******REMOVED***[:space:][:cntrl:]]+`)
	reQuotedValue      = regexp.MustCompile(`^[^\\"]+`)
	reEscapedCharacter = regexp.MustCompile(`^[[:blank:][:graph:]]`)
)

// parseForwardedHeader is a benevolent parser of Forwarded header defined in rfc7239. The header contains
// a comma-separated list of forwarding key-value pairs. Each list element is set by single proxy. The
// function parses only the first element of the list, which is set by the very first proxy. It returns a map
// of corresponding key-value pairs and an unparsed slice of the input string.
//
// Examples of Forwarded header values:
//
//  1. Forwarded: For=192.0.2.43; Proto=https,For="[2001:db8:cafe::17]",For=unknown
//  2. Forwarded: for="192.0.2.43:443"; host="registry.example.org", for="10.10.05.40:80"
//
// The first will be parsed into ***REMOVED***"for": "192.0.2.43", "proto": "https"***REMOVED*** while the second into
// ***REMOVED***"for": "192.0.2.43:443", "host": "registry.example.org"***REMOVED***.
func parseForwardedHeader(forwarded string) (map[string]string, string, error) ***REMOVED***
	// Following are states of forwarded header parser. Any state could transition to a failure.
	const (
		// terminating state; can transition to Parameter
		stateElement = iota
		// terminating state; can transition to KeyValueDelimiter
		stateParameter
		// can transition to Value
		stateKeyValueDelimiter
		// can transition to one of ***REMOVED*** QuotedValue, PairEnd ***REMOVED***
		stateValue
		// can transition to one of ***REMOVED*** EscapedCharacter, PairEnd ***REMOVED***
		stateQuotedValue
		// can transition to one of ***REMOVED*** QuotedValue ***REMOVED***
		stateEscapedCharacter
		// terminating state; can transition to one of ***REMOVED*** Parameter, Element ***REMOVED***
		statePairEnd
	)

	var (
		parameter string
		value     string
		parse     = forwarded[:]
		res       = map[string]string***REMOVED******REMOVED***
		state     = stateElement
	)

Loop:
	for ***REMOVED***
		// skip spaces unless in quoted value
		if state != stateQuotedValue && state != stateEscapedCharacter ***REMOVED***
			parse = strings.TrimLeftFunc(parse, unicode.IsSpace)
		***REMOVED***

		if len(parse) == 0 ***REMOVED***
			if state != stateElement && state != statePairEnd && state != stateParameter ***REMOVED***
				return nil, parse, fmt.Errorf("unexpected end of input")
			***REMOVED***
			// terminating
			break
		***REMOVED***

		switch state ***REMOVED***
		// terminate at list element delimiter
		case stateElement:
			if parse[0] == ',' ***REMOVED***
				parse = parse[1:]
				break Loop
			***REMOVED***
			state = stateParameter

		// parse parameter (the key of key-value pair)
		case stateParameter:
			match := reToken.FindString(parse)
			if len(match) == 0 ***REMOVED***
				return nil, parse, fmt.Errorf("failed to parse token at position %d", len(forwarded)-len(parse))
			***REMOVED***
			parameter = strings.ToLower(match)
			parse = parse[len(match):]
			state = stateKeyValueDelimiter

		// parse '='
		case stateKeyValueDelimiter:
			if parse[0] != '=' ***REMOVED***
				return nil, parse, fmt.Errorf("expected '=', not '%c' at position %d", parse[0], len(forwarded)-len(parse))
			***REMOVED***
			parse = parse[1:]
			state = stateValue

		// parse value or quoted value
		case stateValue:
			if parse[0] == '"' ***REMOVED***
				parse = parse[1:]
				state = stateQuotedValue
			***REMOVED*** else ***REMOVED***
				value = reToken.FindString(parse)
				if len(value) == 0 ***REMOVED***
					return nil, parse, fmt.Errorf("failed to parse value at position %d", len(forwarded)-len(parse))
				***REMOVED***
				if _, exists := res[parameter]; exists ***REMOVED***
					return nil, parse, fmt.Errorf("duplicate parameter %q at position %d", parameter, len(forwarded)-len(parse))
				***REMOVED***
				res[parameter] = value
				parse = parse[len(value):]
				value = ""
				state = statePairEnd
			***REMOVED***

		// parse a part of quoted value until the first backslash
		case stateQuotedValue:
			match := reQuotedValue.FindString(parse)
			value += match
			parse = parse[len(match):]
			switch ***REMOVED***
			case len(parse) == 0:
				return nil, parse, fmt.Errorf("unterminated quoted string")
			case parse[0] == '"':
				res[parameter] = value
				value = ""
				parse = parse[1:]
				state = statePairEnd
			case parse[0] == '\\':
				parse = parse[1:]
				state = stateEscapedCharacter
			***REMOVED***

		// parse escaped character in a quoted string, ignore the backslash
		// transition back to QuotedValue state
		case stateEscapedCharacter:
			c := reEscapedCharacter.FindString(parse)
			if len(c) == 0 ***REMOVED***
				return nil, parse, fmt.Errorf("invalid escape sequence at position %d", len(forwarded)-len(parse)-1)
			***REMOVED***
			value += c
			parse = parse[1:]
			state = stateQuotedValue

		// expect either a new key-value pair, new list or end of input
		case statePairEnd:
			switch parse[0] ***REMOVED***
			case ';':
				parse = parse[1:]
				state = stateParameter
			case ',':
				state = stateElement
			default:
				return nil, parse, fmt.Errorf("expected ',' or ';', not %c at position %d", parse[0], len(forwarded)-len(parse))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return res, parse, nil
***REMOVED***
