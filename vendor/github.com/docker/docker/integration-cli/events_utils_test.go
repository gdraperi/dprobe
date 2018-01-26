package main

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	eventstestutils "github.com/docker/docker/daemon/events/testutils"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/go-check/check"
	"github.com/sirupsen/logrus"
)

// eventMatcher is a function that tries to match an event input.
// It returns true if the event matches and a map with
// a set of key/value to identify the match.
type eventMatcher func(text string) (map[string]string, bool)

// eventMatchProcessor is a function to handle an event match.
// It receives a map of key/value with the information extracted in a match.
type eventMatchProcessor func(matches map[string]string)

// eventObserver runs an events commands and observes its output.
type eventObserver struct ***REMOVED***
	buffer             *bytes.Buffer
	command            *exec.Cmd
	scanner            *bufio.Scanner
	startTime          string
	disconnectionError error
***REMOVED***

// newEventObserver creates the observer and initializes the command
// without running it. Users must call `eventObserver.Start` to start the command.
func newEventObserver(c *check.C, args ...string) (*eventObserver, error) ***REMOVED***
	since := daemonTime(c).Unix()
	return newEventObserverWithBacklog(c, since, args...)
***REMOVED***

// newEventObserverWithBacklog creates a new observer changing the start time of the backlog to return.
func newEventObserverWithBacklog(c *check.C, since int64, args ...string) (*eventObserver, error) ***REMOVED***
	startTime := strconv.FormatInt(since, 10)
	cmdArgs := []string***REMOVED***"events", "--since", startTime***REMOVED***
	if len(args) > 0 ***REMOVED***
		cmdArgs = append(cmdArgs, args...)
	***REMOVED***
	eventsCmd := exec.Command(dockerBinary, cmdArgs...)
	stdout, err := eventsCmd.StdoutPipe()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &eventObserver***REMOVED***
		buffer:    new(bytes.Buffer),
		command:   eventsCmd,
		scanner:   bufio.NewScanner(stdout),
		startTime: startTime,
	***REMOVED***, nil
***REMOVED***

// Start starts the events command.
func (e *eventObserver) Start() error ***REMOVED***
	return e.command.Start()
***REMOVED***

// Stop stops the events command.
func (e *eventObserver) Stop() ***REMOVED***
	e.command.Process.Kill()
	e.command.Wait()
***REMOVED***

// Match tries to match the events output with a given matcher.
func (e *eventObserver) Match(match eventMatcher, process eventMatchProcessor) ***REMOVED***
	for e.scanner.Scan() ***REMOVED***
		text := e.scanner.Text()
		e.buffer.WriteString(text)
		e.buffer.WriteString("\n")

		if matches, ok := match(text); ok ***REMOVED***
			process(matches)
		***REMOVED***
	***REMOVED***

	err := e.scanner.Err()
	if err == nil ***REMOVED***
		err = io.EOF
	***REMOVED***

	logrus.Debugf("EventObserver scanner loop finished: %v", err)
	e.disconnectionError = err
***REMOVED***

func (e *eventObserver) CheckEventError(c *check.C, id, event string, match eventMatcher) ***REMOVED***
	var foundEvent bool
	scannerOut := e.buffer.String()

	if e.disconnectionError != nil ***REMOVED***
		until := daemonUnixTime(c)
		out, _ := dockerCmd(c, "events", "--since", e.startTime, "--until", until)
		events := strings.Split(strings.TrimSpace(out), "\n")
		for _, e := range events ***REMOVED***
			if _, ok := match(e); ok ***REMOVED***
				foundEvent = true
				break
			***REMOVED***
		***REMOVED***
		scannerOut = out
	***REMOVED***
	if !foundEvent ***REMOVED***
		c.Fatalf("failed to observe event `%s` for %s. Disconnection error: %v\nout:\n%v", event, id, e.disconnectionError, scannerOut)
	***REMOVED***
***REMOVED***

// matchEventLine matches a text with the event regular expression.
// It returns the matches and true if the regular expression matches with the given id and event type.
// It returns an empty map and false if there is no match.
func matchEventLine(id, eventType string, actions map[string]chan bool) eventMatcher ***REMOVED***
	return func(text string) (map[string]string, bool) ***REMOVED***
		matches := eventstestutils.ScanMap(text)
		if len(matches) == 0 ***REMOVED***
			return matches, false
		***REMOVED***

		if matchIDAndEventType(matches, id, eventType) ***REMOVED***
			if _, ok := actions[matches["action"]]; ok ***REMOVED***
				return matches, true
			***REMOVED***
		***REMOVED***
		return matches, false
	***REMOVED***
***REMOVED***

// processEventMatch closes an action channel when an event line matches the expected action.
func processEventMatch(actions map[string]chan bool) eventMatchProcessor ***REMOVED***
	return func(matches map[string]string) ***REMOVED***
		if ch, ok := actions[matches["action"]]; ok ***REMOVED***
			ch <- true
		***REMOVED***
	***REMOVED***
***REMOVED***

// parseEventAction parses an event text and returns the action.
// It fails if the text is not in the event format.
func parseEventAction(c *check.C, text string) string ***REMOVED***
	matches := eventstestutils.ScanMap(text)
	return matches["action"]
***REMOVED***

// eventActionsByIDAndType returns the actions for a given id and type.
// It fails if the text is not in the event format.
func eventActionsByIDAndType(c *check.C, events []string, id, eventType string) []string ***REMOVED***
	var filtered []string
	for _, event := range events ***REMOVED***
		matches := eventstestutils.ScanMap(event)
		c.Assert(matches, checker.Not(checker.IsNil))
		if matchIDAndEventType(matches, id, eventType) ***REMOVED***
			filtered = append(filtered, matches["action"])
		***REMOVED***
	***REMOVED***
	return filtered
***REMOVED***

// matchIDAndEventType returns true if an event matches a given id and type.
// It also resolves names in the event attributes if the id doesn't match.
func matchIDAndEventType(matches map[string]string, id, eventType string) bool ***REMOVED***
	return matchEventID(matches, id) && matches["eventType"] == eventType
***REMOVED***

func matchEventID(matches map[string]string, id string) bool ***REMOVED***
	matchID := matches["id"] == id || strings.HasPrefix(matches["id"], id)
	if !matchID && matches["attributes"] != "" ***REMOVED***
		// try matching a name in the attributes
		attributes := map[string]string***REMOVED******REMOVED***
		for _, a := range strings.Split(matches["attributes"], ", ") ***REMOVED***
			kv := strings.Split(a, "=")
			attributes[kv[0]] = kv[1]
		***REMOVED***
		matchID = attributes["name"] == id
	***REMOVED***
	return matchID
***REMOVED***

func parseEvents(c *check.C, out, match string) ***REMOVED***
	events := strings.Split(strings.TrimSpace(out), "\n")
	for _, event := range events ***REMOVED***
		matches := eventstestutils.ScanMap(event)
		matched, err := regexp.MatchString(match, matches["action"])
		c.Assert(err, checker.IsNil)
		c.Assert(matched, checker.True, check.Commentf("Matcher: %s did not match %s", match, matches["action"]))
	***REMOVED***
***REMOVED***

func parseEventsWithID(c *check.C, out, match, id string) ***REMOVED***
	events := strings.Split(strings.TrimSpace(out), "\n")
	for _, event := range events ***REMOVED***
		matches := eventstestutils.ScanMap(event)
		c.Assert(matchEventID(matches, id), checker.True)

		matched, err := regexp.MatchString(match, matches["action"])
		c.Assert(err, checker.IsNil)
		c.Assert(matched, checker.True, check.Commentf("Matcher: %s did not match %s", match, matches["action"]))
	***REMOVED***
***REMOVED***
