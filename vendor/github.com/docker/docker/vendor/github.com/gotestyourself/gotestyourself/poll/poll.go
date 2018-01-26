/*Package poll provides tools for testing asynchronous code.
 */
package poll

import (
	"fmt"
	"time"
)

// TestingT is the subset of testing.T used by WaitOn
type TestingT interface ***REMOVED***
	LogT
	Fatalf(format string, args ...interface***REMOVED******REMOVED***)
***REMOVED***

// LogT is a logging interface that is passed to the WaitOn check function
type LogT interface ***REMOVED***
	Log(args ...interface***REMOVED******REMOVED***)
	Logf(format string, args ...interface***REMOVED******REMOVED***)
***REMOVED***

type helperT interface ***REMOVED***
	Helper()
***REMOVED***

// Settings are used to configure the behaviour of WaitOn
type Settings struct ***REMOVED***
	// Timeout is the maximum time to wait for the condition. Defaults to 10s
	Timeout time.Duration
	// Delay is the time to sleep between checking the condition. Detaults to
	// 1ms
	Delay time.Duration
***REMOVED***

func defaultConfig() *Settings ***REMOVED***
	return &Settings***REMOVED***Timeout: 10 * time.Second, Delay: time.Millisecond***REMOVED***
***REMOVED***

// SettingOp is a function which accepts and modifies Settings
type SettingOp func(config *Settings)

// WithDelay sets the delay to wait between polls
func WithDelay(delay time.Duration) SettingOp ***REMOVED***
	return func(config *Settings) ***REMOVED***
		config.Delay = delay
	***REMOVED***
***REMOVED***

// WithTimeout sets the timeout
func WithTimeout(timeout time.Duration) SettingOp ***REMOVED***
	return func(config *Settings) ***REMOVED***
		config.Timeout = timeout
	***REMOVED***
***REMOVED***

// Result of a check performed by WaitOn
type Result interface ***REMOVED***
	// Error indicates that the check failed and polling should stop, and the
	// the has failed
	Error() error
	// Done indicates that polling should stop, and the test should proceed
	Done() bool
	// Message provides the most recent state when polling has not completed
	Message() string
***REMOVED***

type result struct ***REMOVED***
	done    bool
	message string
	err     error
***REMOVED***

func (r result) Done() bool ***REMOVED***
	return r.done
***REMOVED***

func (r result) Message() string ***REMOVED***
	return r.message
***REMOVED***

func (r result) Error() error ***REMOVED***
	return r.err
***REMOVED***

// Continue returns a Result that indicates to WaitOn that it should continue
// polling. The message text will be used as the failure message if the timeout
// is reached.
func Continue(message string, args ...interface***REMOVED******REMOVED***) Result ***REMOVED***
	return result***REMOVED***message: fmt.Sprintf(message, args...)***REMOVED***
***REMOVED***

// Success returns a Result where Done() returns true, which indicates to WaitOn
// that it should stop polling and exit without an error.
func Success() Result ***REMOVED***
	return result***REMOVED***done: true***REMOVED***
***REMOVED***

// Error returns a Result that indicates to WaitOn that it should fail the test
// and stop polling.
func Error(err error) Result ***REMOVED***
	return result***REMOVED***err: err***REMOVED***
***REMOVED***

// WaitOn a condition or until a timeout. Poll by calling check and exit when
// check returns a done Result. To fail a test and exit polling with an error
// return a error result.
func WaitOn(t TestingT, check func(t LogT) Result, pollOps ...SettingOp) ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	config := defaultConfig()
	for _, pollOp := range pollOps ***REMOVED***
		pollOp(config)
	***REMOVED***

	var lastMessage string
	after := time.After(config.Timeout)
	chResult := make(chan Result)
	for ***REMOVED***
		go func() ***REMOVED***
			chResult <- check(t)
		***REMOVED***()
		select ***REMOVED***
		case <-after:
			if lastMessage == "" ***REMOVED***
				lastMessage = "first check never completed"
			***REMOVED***
			t.Fatalf("timeout hit after %s: %s", config.Timeout, lastMessage)
		case result := <-chResult:
			switch ***REMOVED***
			case result.Error() != nil:
				t.Fatalf("polling check failed: %s", result.Error())
			case result.Done():
				return
			***REMOVED***
			time.Sleep(config.Delay)
			lastMessage = result.Message()
		***REMOVED***
	***REMOVED***
***REMOVED***
