package jsonmessage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	gotty "github.com/Nvveen/Gotty"
	"github.com/docker/docker/pkg/term"
	units "github.com/docker/go-units"
)

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time isalways the same number of characters.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

// JSONError wraps a concrete Code and Message, `Code` is
// is an integer error code, `Message` is the error message.
type JSONError struct ***REMOVED***
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
***REMOVED***

func (e *JSONError) Error() string ***REMOVED***
	return e.Message
***REMOVED***

// JSONProgress describes a Progress. terminalFd is the fd of the current terminal,
// Start is the initial value for the operation. Current is the current status and
// value of the progress made towards Total. Total is the end value describing when
// we made 100% progress for an operation.
type JSONProgress struct ***REMOVED***
	terminalFd uintptr
	Current    int64 `json:"current,omitempty"`
	Total      int64 `json:"total,omitempty"`
	Start      int64 `json:"start,omitempty"`
	// If true, don't show xB/yB
	HideCounts bool   `json:"hidecounts,omitempty"`
	Units      string `json:"units,omitempty"`
***REMOVED***

func (p *JSONProgress) String() string ***REMOVED***
	var (
		width       = 200
		pbBox       string
		numbersBox  string
		timeLeftBox string
	)

	ws, err := term.GetWinsize(p.terminalFd)
	if err == nil ***REMOVED***
		width = int(ws.Width)
	***REMOVED***

	if p.Current <= 0 && p.Total <= 0 ***REMOVED***
		return ""
	***REMOVED***
	if p.Total <= 0 ***REMOVED***
		switch p.Units ***REMOVED***
		case "":
			current := units.HumanSize(float64(p.Current))
			return fmt.Sprintf("%8v", current)
		default:
			return fmt.Sprintf("%d %s", p.Current, p.Units)
		***REMOVED***
	***REMOVED***

	percentage := int(float64(p.Current)/float64(p.Total)*100) / 2
	if percentage > 50 ***REMOVED***
		percentage = 50
	***REMOVED***
	if width > 110 ***REMOVED***
		// this number can't be negative gh#7136
		numSpaces := 0
		if 50-percentage > 0 ***REMOVED***
			numSpaces = 50 - percentage
		***REMOVED***
		pbBox = fmt.Sprintf("[%s>%s] ", strings.Repeat("=", percentage), strings.Repeat(" ", numSpaces))
	***REMOVED***

	switch ***REMOVED***
	case p.HideCounts:
	case p.Units == "": // no units, use bytes
		current := units.HumanSize(float64(p.Current))
		total := units.HumanSize(float64(p.Total))

		numbersBox = fmt.Sprintf("%8v/%v", current, total)

		if p.Current > p.Total ***REMOVED***
			// remove total display if the reported current is wonky.
			numbersBox = fmt.Sprintf("%8v", current)
		***REMOVED***
	default:
		numbersBox = fmt.Sprintf("%d/%d %s", p.Current, p.Total, p.Units)

		if p.Current > p.Total ***REMOVED***
			// remove total display if the reported current is wonky.
			numbersBox = fmt.Sprintf("%d %s", p.Current, p.Units)
		***REMOVED***
	***REMOVED***

	if p.Current > 0 && p.Start > 0 && percentage < 50 ***REMOVED***
		fromStart := time.Now().UTC().Sub(time.Unix(p.Start, 0))
		perEntry := fromStart / time.Duration(p.Current)
		left := time.Duration(p.Total-p.Current) * perEntry
		left = (left / time.Second) * time.Second

		if width > 50 ***REMOVED***
			timeLeftBox = " " + left.String()
		***REMOVED***
	***REMOVED***
	return pbBox + numbersBox + timeLeftBox
***REMOVED***

// JSONMessage defines a message struct. It describes
// the created time, where it from, status, ID of the
// message. It's used for docker events.
type JSONMessage struct ***REMOVED***
	Stream          string        `json:"stream,omitempty"`
	Status          string        `json:"status,omitempty"`
	Progress        *JSONProgress `json:"progressDetail,omitempty"`
	ProgressMessage string        `json:"progress,omitempty"` //deprecated
	ID              string        `json:"id,omitempty"`
	From            string        `json:"from,omitempty"`
	Time            int64         `json:"time,omitempty"`
	TimeNano        int64         `json:"timeNano,omitempty"`
	Error           *JSONError    `json:"errorDetail,omitempty"`
	ErrorMessage    string        `json:"error,omitempty"` //deprecated
	// Aux contains out-of-band data, such as digests for push signing and image id after building.
	Aux *json.RawMessage `json:"aux,omitempty"`
***REMOVED***

/* Satisfied by gotty.TermInfo as well as noTermInfo from below */
type termInfo interface ***REMOVED***
	Parse(attr string, params ...interface***REMOVED******REMOVED***) (string, error)
***REMOVED***

type noTermInfo struct***REMOVED******REMOVED*** // canary used when no terminfo.

func (ti *noTermInfo) Parse(attr string, params ...interface***REMOVED******REMOVED***) (string, error) ***REMOVED***
	return "", fmt.Errorf("noTermInfo")
***REMOVED***

func clearLine(out io.Writer, ti termInfo) ***REMOVED***
	// el2 (clear whole line) is not exposed by terminfo.

	// First clear line from beginning to cursor
	if attr, err := ti.Parse("el1"); err == nil ***REMOVED***
		fmt.Fprintf(out, "%s", attr)
	***REMOVED*** else ***REMOVED***
		fmt.Fprintf(out, "\x1b[1K")
	***REMOVED***
	// Then clear line from cursor to end
	if attr, err := ti.Parse("el"); err == nil ***REMOVED***
		fmt.Fprintf(out, "%s", attr)
	***REMOVED*** else ***REMOVED***
		fmt.Fprintf(out, "\x1b[K")
	***REMOVED***
***REMOVED***

func cursorUp(out io.Writer, ti termInfo, l int) ***REMOVED***
	if l == 0 ***REMOVED*** // Should never be the case, but be tolerant
		return
	***REMOVED***
	if attr, err := ti.Parse("cuu", l); err == nil ***REMOVED***
		fmt.Fprintf(out, "%s", attr)
	***REMOVED*** else ***REMOVED***
		fmt.Fprintf(out, "\x1b[%dA", l)
	***REMOVED***
***REMOVED***

func cursorDown(out io.Writer, ti termInfo, l int) ***REMOVED***
	if l == 0 ***REMOVED*** // Should never be the case, but be tolerant
		return
	***REMOVED***
	if attr, err := ti.Parse("cud", l); err == nil ***REMOVED***
		fmt.Fprintf(out, "%s", attr)
	***REMOVED*** else ***REMOVED***
		fmt.Fprintf(out, "\x1b[%dB", l)
	***REMOVED***
***REMOVED***

// Display displays the JSONMessage to `out`. `termInfo` is non-nil if `out`
// is a terminal. If this is the case, it will erase the entire current line
// when displaying the progressbar.
func (jm *JSONMessage) Display(out io.Writer, termInfo termInfo) error ***REMOVED***
	if jm.Error != nil ***REMOVED***
		if jm.Error.Code == 401 ***REMOVED***
			return fmt.Errorf("authentication is required")
		***REMOVED***
		return jm.Error
	***REMOVED***
	var endl string
	if termInfo != nil && jm.Stream == "" && jm.Progress != nil ***REMOVED***
		clearLine(out, termInfo)
		endl = "\r"
		fmt.Fprintf(out, endl)
	***REMOVED*** else if jm.Progress != nil && jm.Progress.String() != "" ***REMOVED*** //disable progressbar in non-terminal
		return nil
	***REMOVED***
	if jm.TimeNano != 0 ***REMOVED***
		fmt.Fprintf(out, "%s ", time.Unix(0, jm.TimeNano).Format(RFC3339NanoFixed))
	***REMOVED*** else if jm.Time != 0 ***REMOVED***
		fmt.Fprintf(out, "%s ", time.Unix(jm.Time, 0).Format(RFC3339NanoFixed))
	***REMOVED***
	if jm.ID != "" ***REMOVED***
		fmt.Fprintf(out, "%s: ", jm.ID)
	***REMOVED***
	if jm.From != "" ***REMOVED***
		fmt.Fprintf(out, "(from %s) ", jm.From)
	***REMOVED***
	if jm.Progress != nil && termInfo != nil ***REMOVED***
		fmt.Fprintf(out, "%s %s%s", jm.Status, jm.Progress.String(), endl)
	***REMOVED*** else if jm.ProgressMessage != "" ***REMOVED*** //deprecated
		fmt.Fprintf(out, "%s %s%s", jm.Status, jm.ProgressMessage, endl)
	***REMOVED*** else if jm.Stream != "" ***REMOVED***
		fmt.Fprintf(out, "%s%s", jm.Stream, endl)
	***REMOVED*** else ***REMOVED***
		fmt.Fprintf(out, "%s%s\n", jm.Status, endl)
	***REMOVED***
	return nil
***REMOVED***

// DisplayJSONMessagesStream displays a json message stream from `in` to `out`, `isTerminal`
// describes if `out` is a terminal. If this is the case, it will print `\n` at the end of
// each line and move the cursor while displaying.
func DisplayJSONMessagesStream(in io.Reader, out io.Writer, terminalFd uintptr, isTerminal bool, auxCallback func(*json.RawMessage)) error ***REMOVED***
	var (
		dec = json.NewDecoder(in)
		ids = make(map[string]int)
	)

	var termInfo termInfo

	if isTerminal ***REMOVED***
		term := os.Getenv("TERM")
		if term == "" ***REMOVED***
			term = "vt102"
		***REMOVED***

		var err error
		if termInfo, err = gotty.OpenTermInfo(term); err != nil ***REMOVED***
			termInfo = &noTermInfo***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		diff := 0
		var jm JSONMessage
		if err := dec.Decode(&jm); err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			return err
		***REMOVED***

		if jm.Aux != nil ***REMOVED***
			if auxCallback != nil ***REMOVED***
				auxCallback(jm.Aux)
			***REMOVED***
			continue
		***REMOVED***

		if jm.Progress != nil ***REMOVED***
			jm.Progress.terminalFd = terminalFd
		***REMOVED***
		if jm.ID != "" && (jm.Progress != nil || jm.ProgressMessage != "") ***REMOVED***
			line, ok := ids[jm.ID]
			if !ok ***REMOVED***
				// NOTE: This approach of using len(id) to
				// figure out the number of lines of history
				// only works as long as we clear the history
				// when we output something that's not
				// accounted for in the map, such as a line
				// with no ID.
				line = len(ids)
				ids[jm.ID] = line
				if termInfo != nil ***REMOVED***
					fmt.Fprintf(out, "\n")
				***REMOVED***
			***REMOVED***
			diff = len(ids) - line
			if termInfo != nil ***REMOVED***
				cursorUp(out, termInfo, diff)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// When outputting something that isn't progress
			// output, clear the history of previous lines. We
			// don't want progress entries from some previous
			// operation to be updated (for example, pull -a
			// with multiple tags).
			ids = make(map[string]int)
		***REMOVED***
		err := jm.Display(out, termInfo)
		if jm.ID != "" && termInfo != nil ***REMOVED***
			cursorDown(out, termInfo, diff)
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type stream interface ***REMOVED***
	io.Writer
	FD() uintptr
	IsTerminal() bool
***REMOVED***

// DisplayJSONMessagesToStream prints json messages to the output stream
func DisplayJSONMessagesToStream(in io.Reader, stream stream, auxCallback func(*json.RawMessage)) error ***REMOVED***
	return DisplayJSONMessagesStream(in, stream, stream.FD(), stream.IsTerminal(), auxCallback)
***REMOVED***
