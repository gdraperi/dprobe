package docker

import (
	"net/http"
	"sort"
	"strings"
)

type authenticationScheme byte

const (
	basicAuth  authenticationScheme = 1 << iota // Defined in RFC 7617
	digestAuth                                  // Defined in RFC 7616
	bearerAuth                                  // Defined in RFC 6750
)

// challenge carries information from a WWW-Authenticate response header.
// See RFC 2617.
type challenge struct ***REMOVED***
	// scheme is the auth-scheme according to RFC 2617
	scheme authenticationScheme

	// parameters are the auth-params according to RFC 2617
	parameters map[string]string
***REMOVED***

type byScheme []challenge

func (bs byScheme) Len() int      ***REMOVED*** return len(bs) ***REMOVED***
func (bs byScheme) Swap(i, j int) ***REMOVED*** bs[i], bs[j] = bs[j], bs[i] ***REMOVED***

// Sort in priority order: token > digest > basic
func (bs byScheme) Less(i, j int) bool ***REMOVED*** return bs[i].scheme > bs[j].scheme ***REMOVED***

// Octet types from RFC 2616.
type octetType byte

var octetTypes [256]octetType

const (
	isToken octetType = 1 << iota
	isSpace
)

func init() ***REMOVED***
	// OCTET      = <any 8-bit sequence of data>
	// CHAR       = <any US-ASCII character (octets 0 - 127)>
	// CTL        = <any US-ASCII control character (octets 0 - 31) and DEL (127)>
	// CR         = <US-ASCII CR, carriage return (13)>
	// LF         = <US-ASCII LF, linefeed (10)>
	// SP         = <US-ASCII SP, space (32)>
	// HT         = <US-ASCII HT, horizontal-tab (9)>
	// <">        = <US-ASCII double-quote mark (34)>
	// CRLF       = CR LF
	// LWS        = [CRLF] 1*( SP | HT )
	// TEXT       = <any OCTET except CTLs, but including LWS>
	// separators = "(" | ")" | "<" | ">" | "@" | "," | ";" | ":" | "\" | <">
	//              | "/" | "[" | "]" | "?" | "=" | "***REMOVED***" | "***REMOVED***" | SP | HT
	// token      = 1*<any CHAR except CTLs or separators>
	// qdtext     = <any TEXT except <">>

	for c := 0; c < 256; c++ ***REMOVED***
		var t octetType
		isCtl := c <= 31 || c == 127
		isChar := 0 <= c && c <= 127
		isSeparator := strings.IndexRune(" \t\"(),/:;<=>?@[]\\***REMOVED******REMOVED***", rune(c)) >= 0
		if strings.IndexRune(" \t\r\n", rune(c)) >= 0 ***REMOVED***
			t |= isSpace
		***REMOVED***
		if isChar && !isCtl && !isSeparator ***REMOVED***
			t |= isToken
		***REMOVED***
		octetTypes[c] = t
	***REMOVED***
***REMOVED***

func parseAuthHeader(header http.Header) []challenge ***REMOVED***
	challenges := []challenge***REMOVED******REMOVED***
	for _, h := range header[http.CanonicalHeaderKey("WWW-Authenticate")] ***REMOVED***
		v, p := parseValueAndParams(h)
		var s authenticationScheme
		switch v ***REMOVED***
		case "basic":
			s = basicAuth
		case "digest":
			s = digestAuth
		case "bearer":
			s = bearerAuth
		default:
			continue
		***REMOVED***
		challenges = append(challenges, challenge***REMOVED***scheme: s, parameters: p***REMOVED***)
	***REMOVED***
	sort.Stable(byScheme(challenges))
	return challenges
***REMOVED***

func parseValueAndParams(header string) (value string, params map[string]string) ***REMOVED***
	params = make(map[string]string)
	value, s := expectToken(header)
	if value == "" ***REMOVED***
		return
	***REMOVED***
	value = strings.ToLower(value)
	for ***REMOVED***
		var pkey string
		pkey, s = expectToken(skipSpace(s))
		if pkey == "" ***REMOVED***
			return
		***REMOVED***
		if !strings.HasPrefix(s, "=") ***REMOVED***
			return
		***REMOVED***
		var pvalue string
		pvalue, s = expectTokenOrQuoted(s[1:])
		if pvalue == "" ***REMOVED***
			return
		***REMOVED***
		pkey = strings.ToLower(pkey)
		params[pkey] = pvalue
		s = skipSpace(s)
		if !strings.HasPrefix(s, ",") ***REMOVED***
			return
		***REMOVED***
		s = s[1:]
	***REMOVED***
***REMOVED***

func skipSpace(s string) (rest string) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		if octetTypes[s[i]]&isSpace == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return s[i:]
***REMOVED***

func expectToken(s string) (token, rest string) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		if octetTypes[s[i]]&isToken == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return s[:i], s[i:]
***REMOVED***

func expectTokenOrQuoted(s string) (value string, rest string) ***REMOVED***
	if !strings.HasPrefix(s, "\"") ***REMOVED***
		return expectToken(s)
	***REMOVED***
	s = s[1:]
	for i := 0; i < len(s); i++ ***REMOVED***
		switch s[i] ***REMOVED***
		case '"':
			return s[:i], s[i+1:]
		case '\\':
			p := make([]byte, len(s)-1)
			j := copy(p, s[:i])
			escape := true
			for i = i + 1; i < len(s); i++ ***REMOVED***
				b := s[i]
				switch ***REMOVED***
				case escape:
					escape = false
					p[j] = b
					j++
				case b == '\\':
					escape = true
				case b == '"':
					return string(p[:j]), s[i+1:]
				default:
					p[j] = b
					j++
				***REMOVED***
			***REMOVED***
			return "", ""
		***REMOVED***
	***REMOVED***
	return "", ""
***REMOVED***
