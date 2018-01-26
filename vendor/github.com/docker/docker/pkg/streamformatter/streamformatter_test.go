package streamformatter

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawProgressFormatterFormatStatus(t *testing.T) ***REMOVED***
	sf := rawProgressFormatter***REMOVED******REMOVED***
	res := sf.formatStatus("ID", "%s%d", "a", 1)
	assert.Equal(t, "a1\r\n", string(res))
***REMOVED***

func TestRawProgressFormatterFormatProgress(t *testing.T) ***REMOVED***
	sf := rawProgressFormatter***REMOVED******REMOVED***
	jsonProgress := &jsonmessage.JSONProgress***REMOVED***
		Current: 15,
		Total:   30,
		Start:   1,
	***REMOVED***
	res := sf.formatProgress("id", "action", jsonProgress, nil)
	out := string(res)
	assert.True(t, strings.HasPrefix(out, "action [===="))
	assert.Contains(t, out, "15B/30B")
	assert.True(t, strings.HasSuffix(out, "\r"))
***REMOVED***

func TestFormatStatus(t *testing.T) ***REMOVED***
	res := FormatStatus("ID", "%s%d", "a", 1)
	expected := `***REMOVED***"status":"a1","id":"ID"***REMOVED***` + streamNewline
	assert.Equal(t, expected, string(res))
***REMOVED***

func TestFormatError(t *testing.T) ***REMOVED***
	res := FormatError(errors.New("Error for formatter"))
	expected := `***REMOVED***"errorDetail":***REMOVED***"message":"Error for formatter"***REMOVED***,"error":"Error for formatter"***REMOVED***` + "\r\n"
	assert.Equal(t, expected, string(res))
***REMOVED***

func TestFormatJSONError(t *testing.T) ***REMOVED***
	err := &jsonmessage.JSONError***REMOVED***Code: 50, Message: "Json error"***REMOVED***
	res := FormatError(err)
	expected := `***REMOVED***"errorDetail":***REMOVED***"code":50,"message":"Json error"***REMOVED***,"error":"Json error"***REMOVED***` + streamNewline
	assert.Equal(t, expected, string(res))
***REMOVED***

func TestJsonProgressFormatterFormatProgress(t *testing.T) ***REMOVED***
	sf := &jsonProgressFormatter***REMOVED******REMOVED***
	jsonProgress := &jsonmessage.JSONProgress***REMOVED***
		Current: 15,
		Total:   30,
		Start:   1,
	***REMOVED***
	res := sf.formatProgress("id", "action", jsonProgress, &AuxFormatter***REMOVED***Writer: &bytes.Buffer***REMOVED******REMOVED******REMOVED***)
	msg := &jsonmessage.JSONMessage***REMOVED******REMOVED***

	require.NoError(t, json.Unmarshal(res, msg))
	assert.Equal(t, "id", msg.ID)
	assert.Equal(t, "action", msg.Status)

	// jsonProgress will always be in the format of:
	// [=========================>                         ]      15B/30B 412910h51m30s
	// The last entry '404933h7m11s' is the timeLeftBox.
	// However, the timeLeftBox field may change as jsonProgress.String() depends on time.Now().
	// Therefore, we have to strip the timeLeftBox from the strings to do the comparison.

	// Compare the jsonProgress strings before the timeLeftBox
	expectedProgress := "[=========================>                         ]      15B/30B"
	// if terminal column is <= 110, expectedProgressShort is expected.
	expectedProgressShort := "      15B/30B"
	if !(strings.HasPrefix(msg.ProgressMessage, expectedProgress) ||
		strings.HasPrefix(msg.ProgressMessage, expectedProgressShort)) ***REMOVED***
		t.Fatalf("ProgressMessage without the timeLeftBox must be %s or %s, got: %s",
			expectedProgress, expectedProgressShort, msg.ProgressMessage)
	***REMOVED***

	assert.Equal(t, jsonProgress, msg.Progress)
***REMOVED***

func TestJsonProgressFormatterFormatStatus(t *testing.T) ***REMOVED***
	sf := jsonProgressFormatter***REMOVED******REMOVED***
	res := sf.formatStatus("ID", "%s%d", "a", 1)
	assert.Equal(t, `***REMOVED***"status":"a1","id":"ID"***REMOVED***`+streamNewline, string(res))
***REMOVED***

func TestNewJSONProgressOutput(t *testing.T) ***REMOVED***
	b := bytes.Buffer***REMOVED******REMOVED***
	b.Write(FormatStatus("id", "Downloading"))
	_ = NewJSONProgressOutput(&b, false)
	assert.Equal(t, `***REMOVED***"status":"Downloading","id":"id"***REMOVED***`+streamNewline, b.String())
***REMOVED***

func TestAuxFormatterEmit(t *testing.T) ***REMOVED***
	b := bytes.Buffer***REMOVED******REMOVED***
	aux := &AuxFormatter***REMOVED***Writer: &b***REMOVED***
	sampleAux := &struct ***REMOVED***
		Data string
	***REMOVED******REMOVED***"Additional data"***REMOVED***
	err := aux.Emit(sampleAux)
	require.NoError(t, err)
	assert.Equal(t, `***REMOVED***"aux":***REMOVED***"Data":"Additional data"***REMOVED******REMOVED***`+streamNewline, b.String())
***REMOVED***
