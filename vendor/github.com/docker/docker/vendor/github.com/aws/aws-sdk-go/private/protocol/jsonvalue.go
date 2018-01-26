package protocol

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
)

// EscapeMode is the mode that should be use for escaping a value
type EscapeMode uint

// The modes for escaping a value before it is marshaled, and unmarshaled.
const (
	NoEscape EscapeMode = iota
	Base64Escape
	QuotedEscape
)

// EncodeJSONValue marshals the value into a JSON string, and optionally base64
// encodes the string before returning it.
//
// Will panic if the escape mode is unknown.
func EncodeJSONValue(v aws.JSONValue, escape EscapeMode) (string, error) ***REMOVED***
	b, err := json.Marshal(v)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	switch escape ***REMOVED***
	case NoEscape:
		return string(b), nil
	case Base64Escape:
		return base64.StdEncoding.EncodeToString(b), nil
	case QuotedEscape:
		return strconv.Quote(string(b)), nil
	***REMOVED***

	panic(fmt.Sprintf("EncodeJSONValue called with unknown EscapeMode, %v", escape))
***REMOVED***

// DecodeJSONValue will attempt to decode the string input as a JSONValue.
// Optionally decoding base64 the value first before JSON unmarshaling.
//
// Will panic if the escape mode is unknown.
func DecodeJSONValue(v string, escape EscapeMode) (aws.JSONValue, error) ***REMOVED***
	var b []byte
	var err error

	switch escape ***REMOVED***
	case NoEscape:
		b = []byte(v)
	case Base64Escape:
		b, err = base64.StdEncoding.DecodeString(v)
	case QuotedEscape:
		var u string
		u, err = strconv.Unquote(v)
		b = []byte(u)
	default:
		panic(fmt.Sprintf("DecodeJSONValue called with unknown EscapeMode, %v", escape))
	***REMOVED***

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	m := aws.JSONValue***REMOVED******REMOVED***
	err = json.Unmarshal(b, &m)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return m, nil
***REMOVED***
