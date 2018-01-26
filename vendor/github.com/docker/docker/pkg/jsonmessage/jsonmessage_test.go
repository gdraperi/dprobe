package jsonmessage

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/pkg/term"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) ***REMOVED***
	je := JSONError***REMOVED***404, "Not found"***REMOVED***
	if je.Error() != "Not found" ***REMOVED***
		t.Fatalf("Expected 'Not found' got '%s'", je.Error())
	***REMOVED***
***REMOVED***

func TestProgress(t *testing.T) ***REMOVED***
	termsz, err := term.GetWinsize(0)
	if err != nil ***REMOVED***
		// we can safely ignore the err here
		termsz = nil
	***REMOVED***
	jp := JSONProgress***REMOVED******REMOVED***
	if jp.String() != "" ***REMOVED***
		t.Fatalf("Expected empty string, got '%s'", jp.String())
	***REMOVED***

	expected := "      1B"
	jp2 := JSONProgress***REMOVED***Current: 1***REMOVED***
	if jp2.String() != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, jp2.String())
	***REMOVED***

	expectedStart := "[==========>                                        ]      20B/100B"
	if termsz != nil && termsz.Width <= 110 ***REMOVED***
		expectedStart = "    20B/100B"
	***REMOVED***
	jp3 := JSONProgress***REMOVED***Current: 20, Total: 100, Start: time.Now().Unix()***REMOVED***
	// Just look at the start of the string
	// (the remaining time is really hard to test -_-)
	if jp3.String()[:len(expectedStart)] != expectedStart ***REMOVED***
		t.Fatalf("Expected to start with %q, got %q", expectedStart, jp3.String())
	***REMOVED***

	expected = "[=========================>                         ]      50B/100B"
	if termsz != nil && termsz.Width <= 110 ***REMOVED***
		expected = "    50B/100B"
	***REMOVED***
	jp4 := JSONProgress***REMOVED***Current: 50, Total: 100***REMOVED***
	if jp4.String() != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, jp4.String())
	***REMOVED***

	// this number can't be negative gh#7136
	expected = "[==================================================>]      50B"
	if termsz != nil && termsz.Width <= 110 ***REMOVED***
		expected = "    50B"
	***REMOVED***
	jp5 := JSONProgress***REMOVED***Current: 50, Total: 40***REMOVED***
	if jp5.String() != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, jp5.String())
	***REMOVED***

	expected = "[=========================>                         ] 50/100 units"
	if termsz != nil && termsz.Width <= 110 ***REMOVED***
		expected = "    50/100 units"
	***REMOVED***
	jp6 := JSONProgress***REMOVED***Current: 50, Total: 100, Units: "units"***REMOVED***
	if jp6.String() != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, jp6.String())
	***REMOVED***

	// this number can't be negative
	expected = "[==================================================>] 50 units"
	if termsz != nil && termsz.Width <= 110 ***REMOVED***
		expected = "    50 units"
	***REMOVED***
	jp7 := JSONProgress***REMOVED***Current: 50, Total: 40, Units: "units"***REMOVED***
	if jp7.String() != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, jp7.String())
	***REMOVED***

	expected = "[=========================>                         ] "
	if termsz != nil && termsz.Width <= 110 ***REMOVED***
		expected = ""
	***REMOVED***
	jp8 := JSONProgress***REMOVED***Current: 50, Total: 100, HideCounts: true***REMOVED***
	if jp8.String() != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, jp8.String())
	***REMOVED***
***REMOVED***

func TestJSONMessageDisplay(t *testing.T) ***REMOVED***
	now := time.Now()
	messages := map[JSONMessage][]string***REMOVED***
		// Empty
		***REMOVED******REMOVED***: ***REMOVED***"\n", "\n"***REMOVED***,
		// Status
		***REMOVED***
			Status: "status",
		***REMOVED***: ***REMOVED***
			"status\n",
			"status\n",
		***REMOVED***,
		// General
		***REMOVED***
			Time:   now.Unix(),
			ID:     "ID",
			From:   "From",
			Status: "status",
		***REMOVED***: ***REMOVED***
			fmt.Sprintf("%v ID: (from From) status\n", time.Unix(now.Unix(), 0).Format(RFC3339NanoFixed)),
			fmt.Sprintf("%v ID: (from From) status\n", time.Unix(now.Unix(), 0).Format(RFC3339NanoFixed)),
		***REMOVED***,
		// General, with nano precision time
		***REMOVED***
			TimeNano: now.UnixNano(),
			ID:       "ID",
			From:     "From",
			Status:   "status",
		***REMOVED***: ***REMOVED***
			fmt.Sprintf("%v ID: (from From) status\n", time.Unix(0, now.UnixNano()).Format(RFC3339NanoFixed)),
			fmt.Sprintf("%v ID: (from From) status\n", time.Unix(0, now.UnixNano()).Format(RFC3339NanoFixed)),
		***REMOVED***,
		// General, with both times Nano is preferred
		***REMOVED***
			Time:     now.Unix(),
			TimeNano: now.UnixNano(),
			ID:       "ID",
			From:     "From",
			Status:   "status",
		***REMOVED***: ***REMOVED***
			fmt.Sprintf("%v ID: (from From) status\n", time.Unix(0, now.UnixNano()).Format(RFC3339NanoFixed)),
			fmt.Sprintf("%v ID: (from From) status\n", time.Unix(0, now.UnixNano()).Format(RFC3339NanoFixed)),
		***REMOVED***,
		// Stream over status
		***REMOVED***
			Status: "status",
			Stream: "stream",
		***REMOVED***: ***REMOVED***
			"stream",
			"stream",
		***REMOVED***,
		// With progress message
		***REMOVED***
			Status:          "status",
			ProgressMessage: "progressMessage",
		***REMOVED***: ***REMOVED***
			"status progressMessage",
			"status progressMessage",
		***REMOVED***,
		// With progress, stream empty
		***REMOVED***
			Status:   "status",
			Stream:   "",
			Progress: &JSONProgress***REMOVED***Current: 1***REMOVED***,
		***REMOVED***: ***REMOVED***
			"",
			fmt.Sprintf("%c[1K%c[K\rstatus       1B\r", 27, 27),
		***REMOVED***,
	***REMOVED***

	// The tests :)
	for jsonMessage, expectedMessages := range messages ***REMOVED***
		// Without terminal
		data := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
		if err := jsonMessage.Display(data, nil); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if data.String() != expectedMessages[0] ***REMOVED***
			t.Fatalf("Expected %q,got %q", expectedMessages[0], data.String())
		***REMOVED***
		// With terminal
		data = bytes.NewBuffer([]byte***REMOVED******REMOVED***)
		if err := jsonMessage.Display(data, &noTermInfo***REMOVED******REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if data.String() != expectedMessages[1] ***REMOVED***
			t.Fatalf("\nExpected %q\n     got %q", expectedMessages[1], data.String())
		***REMOVED***
	***REMOVED***
***REMOVED***

// Test JSONMessage with an Error. It will return an error with the text as error, not the meaning of the HTTP code.
func TestJSONMessageDisplayWithJSONError(t *testing.T) ***REMOVED***
	data := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	jsonMessage := JSONMessage***REMOVED***Error: &JSONError***REMOVED***404, "Can't find it"***REMOVED******REMOVED***

	err := jsonMessage.Display(data, &noTermInfo***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Can't find it" ***REMOVED***
		t.Fatalf("Expected a JSONError 404, got %q", err)
	***REMOVED***

	jsonMessage = JSONMessage***REMOVED***Error: &JSONError***REMOVED***401, "Anything"***REMOVED******REMOVED***
	err = jsonMessage.Display(data, &noTermInfo***REMOVED******REMOVED***)
	assert.EqualError(t, err, "authentication is required")
***REMOVED***

func TestDisplayJSONMessagesStreamInvalidJSON(t *testing.T) ***REMOVED***
	var (
		inFd uintptr
	)
	data := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	reader := strings.NewReader("This is not a 'valid' JSON []")
	inFd, _ = term.GetFdInfo(reader)

	if err := DisplayJSONMessagesStream(reader, data, inFd, false, nil); err == nil && err.Error()[:17] != "invalid character" ***REMOVED***
		t.Fatalf("Should have thrown an error (invalid character in ..), got %q", err)
	***REMOVED***
***REMOVED***

func TestDisplayJSONMessagesStream(t *testing.T) ***REMOVED***
	var (
		inFd uintptr
	)

	messages := map[string][]string***REMOVED***
		// empty string
		"": ***REMOVED***
			"",
			""***REMOVED***,
		// Without progress & ID
		"***REMOVED*** \"status\": \"status\" ***REMOVED***": ***REMOVED***
			"status\n",
			"status\n",
		***REMOVED***,
		// Without progress, with ID
		"***REMOVED*** \"id\": \"ID\",\"status\": \"status\" ***REMOVED***": ***REMOVED***
			"ID: status\n",
			fmt.Sprintf("ID: status\n"),
		***REMOVED***,
		// With progress
		"***REMOVED*** \"id\": \"ID\", \"status\": \"status\", \"progress\": \"ProgressMessage\" ***REMOVED***": ***REMOVED***
			"ID: status ProgressMessage",
			fmt.Sprintf("\n%c[%dAID: status ProgressMessage%c[%dB", 27, 1, 27, 1),
		***REMOVED***,
		// With progressDetail
		"***REMOVED*** \"id\": \"ID\", \"status\": \"status\", \"progressDetail\": ***REMOVED*** \"Current\": 1***REMOVED*** ***REMOVED***": ***REMOVED***
			"", // progressbar is disabled in non-terminal
			fmt.Sprintf("\n%c[%dA%c[1K%c[K\rID: status       1B\r%c[%dB", 27, 1, 27, 27, 27, 1),
		***REMOVED***,
	***REMOVED***

	// Use $TERM which is unlikely to exist, forcing DisplayJSONMessageStream to
	// (hopefully) use &noTermInfo.
	origTerm := os.Getenv("TERM")
	os.Setenv("TERM", "xyzzy-non-existent-terminfo")

	for jsonMessage, expectedMessages := range messages ***REMOVED***
		data := bytes.NewBuffer([]byte***REMOVED******REMOVED***)
		reader := strings.NewReader(jsonMessage)
		inFd, _ = term.GetFdInfo(reader)

		// Without terminal
		if err := DisplayJSONMessagesStream(reader, data, inFd, false, nil); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if data.String() != expectedMessages[0] ***REMOVED***
			t.Fatalf("Expected an %q, got %q", expectedMessages[0], data.String())
		***REMOVED***

		// With terminal
		data = bytes.NewBuffer([]byte***REMOVED******REMOVED***)
		reader = strings.NewReader(jsonMessage)
		if err := DisplayJSONMessagesStream(reader, data, inFd, true, nil); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if data.String() != expectedMessages[1] ***REMOVED***
			t.Fatalf("\nExpected %q\n     got %q", expectedMessages[1], data.String())
		***REMOVED***
	***REMOVED***
	os.Setenv("TERM", origTerm)

***REMOVED***
