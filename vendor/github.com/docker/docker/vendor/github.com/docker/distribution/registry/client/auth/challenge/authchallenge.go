package challenge

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Challenge carries information from a WWW-Authenticate response header.
// See RFC 2617.
type Challenge struct ***REMOVED***
	// Scheme is the auth-scheme according to RFC 2617
	Scheme string

	// Parameters are the auth-params according to RFC 2617
	Parameters map[string]string
***REMOVED***

// Manager manages the challenges for endpoints.
// The challenges are pulled out of HTTP responses. Only
// responses which expect challenges should be added to
// the manager, since a non-unauthorized request will be
// viewed as not requiring challenges.
type Manager interface ***REMOVED***
	// GetChallenges returns the challenges for the given
	// endpoint URL.
	GetChallenges(endpoint url.URL) ([]Challenge, error)

	// AddResponse adds the response to the challenge
	// manager. The challenges will be parsed out of
	// the WWW-Authenicate headers and added to the
	// URL which was produced the response. If the
	// response was authorized, any challenges for the
	// endpoint will be cleared.
	AddResponse(resp *http.Response) error
***REMOVED***

// NewSimpleManager returns an instance of
// Manger which only maps endpoints to challenges
// based on the responses which have been added the
// manager. The simple manager will make no attempt to
// perform requests on the endpoints or cache the responses
// to a backend.
func NewSimpleManager() Manager ***REMOVED***
	return &simpleManager***REMOVED***
		Challanges: make(map[string][]Challenge),
	***REMOVED***
***REMOVED***

type simpleManager struct ***REMOVED***
	sync.RWMutex
	Challanges map[string][]Challenge
***REMOVED***

func normalizeURL(endpoint *url.URL) ***REMOVED***
	endpoint.Host = strings.ToLower(endpoint.Host)
	endpoint.Host = canonicalAddr(endpoint)
***REMOVED***

func (m *simpleManager) GetChallenges(endpoint url.URL) ([]Challenge, error) ***REMOVED***
	normalizeURL(&endpoint)

	m.RLock()
	defer m.RUnlock()
	challenges := m.Challanges[endpoint.String()]
	return challenges, nil
***REMOVED***

func (m *simpleManager) AddResponse(resp *http.Response) error ***REMOVED***
	challenges := ResponseChallenges(resp)
	if resp.Request == nil ***REMOVED***
		return fmt.Errorf("missing request reference")
	***REMOVED***
	urlCopy := url.URL***REMOVED***
		Path:   resp.Request.URL.Path,
		Host:   resp.Request.URL.Host,
		Scheme: resp.Request.URL.Scheme,
	***REMOVED***
	normalizeURL(&urlCopy)

	m.Lock()
	defer m.Unlock()
	m.Challanges[urlCopy.String()] = challenges
	return nil
***REMOVED***

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

// ResponseChallenges returns a list of authorization challenges
// for the given http Response. Challenges are only checked if
// the response status code was a 401.
func ResponseChallenges(resp *http.Response) []Challenge ***REMOVED***
	if resp.StatusCode == http.StatusUnauthorized ***REMOVED***
		// Parse the WWW-Authenticate Header and store the challenges
		// on this endpoint object.
		return parseAuthHeader(resp.Header)
	***REMOVED***

	return nil
***REMOVED***

func parseAuthHeader(header http.Header) []Challenge ***REMOVED***
	challenges := []Challenge***REMOVED******REMOVED***
	for _, h := range header[http.CanonicalHeaderKey("WWW-Authenticate")] ***REMOVED***
		v, p := parseValueAndParams(h)
		if v != "" ***REMOVED***
			challenges = append(challenges, Challenge***REMOVED***Scheme: v, Parameters: p***REMOVED***)
		***REMOVED***
	***REMOVED***
	return challenges
***REMOVED***

func parseValueAndParams(header string) (value string, params map[string]string) ***REMOVED***
	params = make(map[string]string)
	value, s := expectToken(header)
	if value == "" ***REMOVED***
		return
	***REMOVED***
	value = strings.ToLower(value)
	s = "," + skipSpace(s)
	for strings.HasPrefix(s, ",") ***REMOVED***
		var pkey string
		pkey, s = expectToken(skipSpace(s[1:]))
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
	***REMOVED***
	return
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
