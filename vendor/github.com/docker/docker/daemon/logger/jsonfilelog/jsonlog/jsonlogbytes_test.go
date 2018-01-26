package jsonlog

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONLogsMarshalJSONBuf(t *testing.T) ***REMOVED***
	logs := map[*JSONLogs]string***REMOVED***
		***REMOVED***Log: []byte(`"A log line with \\"`)***REMOVED***:                  `^***REMOVED***\"log\":\"\\\"A log line with \\\\\\\\\\\"\",\"time\":`,
		***REMOVED***Log: []byte("A log line")***REMOVED***:                            `^***REMOVED***\"log\":\"A log line\",\"time\":`,
		***REMOVED***Log: []byte("A log line with \r")***REMOVED***:                    `^***REMOVED***\"log\":\"A log line with \\r\",\"time\":`,
		***REMOVED***Log: []byte("A log line with & < >")***REMOVED***:                 `^***REMOVED***\"log\":\"A log line with \\u0026 \\u003c \\u003e\",\"time\":`,
		***REMOVED***Log: []byte("A log line with utf8 : 🚀 ψ ω β")***REMOVED***:        `^***REMOVED***\"log\":\"A log line with utf8 : 🚀 ψ ω β\",\"time\":`,
		***REMOVED***Stream: "stdout"***REMOVED***:                                     `^***REMOVED***\"stream\":\"stdout\",\"time\":`,
		***REMOVED***Stream: "stdout", Log: []byte("A log line")***REMOVED***:          `^***REMOVED***\"log\":\"A log line\",\"stream\":\"stdout\",\"time\":`,
		***REMOVED***Created: time.Date(2017, 9, 1, 1, 1, 1, 1, time.UTC)***REMOVED***: `^***REMOVED***\"time\":"2017-09-01T01:01:01.000000001Z"***REMOVED***$`,

		***REMOVED******REMOVED***: `^***REMOVED***\"time\":"0001-01-01T00:00:00Z"***REMOVED***$`,
		// These ones are a little weird
		***REMOVED***Log: []byte("\u2028 \u2029")***REMOVED***: `^***REMOVED***\"log\":\"\\u2028 \\u2029\",\"time\":`,
		***REMOVED***Log: []byte***REMOVED***0xaF***REMOVED******REMOVED***:            `^***REMOVED***\"log\":\"\\ufffd\",\"time\":`,
		***REMOVED***Log: []byte***REMOVED***0x7F***REMOVED******REMOVED***:            `^***REMOVED***\"log\":\"\x7f\",\"time\":`,
		// with raw attributes
		***REMOVED***Log: []byte("A log line"), RawAttrs: []byte(`***REMOVED***"hello":"world","value":1234***REMOVED***`)***REMOVED***: `^***REMOVED***\"log\":\"A log line\",\"attrs\":***REMOVED***\"hello\":\"world\",\"value\":1234***REMOVED***,\"time\":`,
		// with Tag set
		***REMOVED***Log: []byte("A log line with tag"), RawAttrs: []byte(`***REMOVED***"hello":"world","value":1234***REMOVED***`)***REMOVED***: `^***REMOVED***\"log\":\"A log line with tag\",\"attrs\":***REMOVED***\"hello\":\"world\",\"value\":1234***REMOVED***,\"time\":`,
	***REMOVED***
	for jsonLog, expression := range logs ***REMOVED***
		var buf bytes.Buffer
		err := jsonLog.MarshalJSONBuf(&buf)
		require.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile(expression), buf.String())
		assert.NoError(t, json.Unmarshal(buf.Bytes(), &map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***))
	***REMOVED***
***REMOVED***
