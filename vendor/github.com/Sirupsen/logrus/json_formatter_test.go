package logrus

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestErrorNotLost(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED******REMOVED***

	b, err := formatter.Format(WithField("error", errors.New("wild walrus")))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***

	entry := make(map[string]interface***REMOVED******REMOVED***)
	err = json.Unmarshal(b, &entry)
	if err != nil ***REMOVED***
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	***REMOVED***

	if entry["error"] != "wild walrus" ***REMOVED***
		t.Fatal("Error field not set")
	***REMOVED***
***REMOVED***

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED******REMOVED***

	b, err := formatter.Format(WithField("omg", errors.New("wild walrus")))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***

	entry := make(map[string]interface***REMOVED******REMOVED***)
	err = json.Unmarshal(b, &entry)
	if err != nil ***REMOVED***
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	***REMOVED***

	if entry["omg"] != "wild walrus" ***REMOVED***
		t.Fatal("Error field not set")
	***REMOVED***
***REMOVED***

func TestFieldClashWithTime(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED******REMOVED***

	b, err := formatter.Format(WithField("time", "right now!"))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***

	entry := make(map[string]interface***REMOVED******REMOVED***)
	err = json.Unmarshal(b, &entry)
	if err != nil ***REMOVED***
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	***REMOVED***

	if entry["fields.time"] != "right now!" ***REMOVED***
		t.Fatal("fields.time not set to original time field")
	***REMOVED***

	if entry["time"] != "0001-01-01T00:00:00Z" ***REMOVED***
		t.Fatal("time field not set to current time, was: ", entry["time"])
	***REMOVED***
***REMOVED***

func TestFieldClashWithMsg(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED******REMOVED***

	b, err := formatter.Format(WithField("msg", "something"))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***

	entry := make(map[string]interface***REMOVED******REMOVED***)
	err = json.Unmarshal(b, &entry)
	if err != nil ***REMOVED***
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	***REMOVED***

	if entry["fields.msg"] != "something" ***REMOVED***
		t.Fatal("fields.msg not set to original msg field")
	***REMOVED***
***REMOVED***

func TestFieldClashWithLevel(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED******REMOVED***

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***

	entry := make(map[string]interface***REMOVED******REMOVED***)
	err = json.Unmarshal(b, &entry)
	if err != nil ***REMOVED***
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	***REMOVED***

	if entry["fields.level"] != "something" ***REMOVED***
		t.Fatal("fields.level not set to original level field")
	***REMOVED***
***REMOVED***

func TestJSONEntryEndsWithNewline(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED******REMOVED***

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***

	if b[len(b)-1] != '\n' ***REMOVED***
		t.Fatal("Expected JSON log entry to end with a newline")
	***REMOVED***
***REMOVED***

func TestJSONMessageKey(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED***
		FieldMap: FieldMap***REMOVED***
			FieldKeyMsg: "message",
		***REMOVED***,
	***REMOVED***

	b, err := formatter.Format(&Entry***REMOVED***Message: "oh hai"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***
	s := string(b)
	if !(strings.Contains(s, "message") && strings.Contains(s, "oh hai")) ***REMOVED***
		t.Fatal("Expected JSON to format message key")
	***REMOVED***
***REMOVED***

func TestJSONLevelKey(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED***
		FieldMap: FieldMap***REMOVED***
			FieldKeyLevel: "somelevel",
		***REMOVED***,
	***REMOVED***

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***
	s := string(b)
	if !strings.Contains(s, "somelevel") ***REMOVED***
		t.Fatal("Expected JSON to format level key")
	***REMOVED***
***REMOVED***

func TestJSONTimeKey(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED***
		FieldMap: FieldMap***REMOVED***
			FieldKeyTime: "timeywimey",
		***REMOVED***,
	***REMOVED***

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***
	s := string(b)
	if !strings.Contains(s, "timeywimey") ***REMOVED***
		t.Fatal("Expected JSON to format time key")
	***REMOVED***
***REMOVED***

func TestJSONDisableTimestamp(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED***
		DisableTimestamp: true,
	***REMOVED***

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***
	s := string(b)
	if strings.Contains(s, FieldKeyTime) ***REMOVED***
		t.Error("Did not prevent timestamp", s)
	***REMOVED***
***REMOVED***

func TestJSONEnableTimestamp(t *testing.T) ***REMOVED***
	formatter := &JSONFormatter***REMOVED******REMOVED***

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil ***REMOVED***
		t.Fatal("Unable to format entry: ", err)
	***REMOVED***
	s := string(b)
	if !strings.Contains(s, FieldKeyTime) ***REMOVED***
		t.Error("Timestamp not present", s)
	***REMOVED***
***REMOVED***
