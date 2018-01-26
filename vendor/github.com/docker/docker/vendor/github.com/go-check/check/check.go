// Package check is a rich testing extension for Go's testing package.
//
// For details about the project, see:
//
//     http://labix.org/gocheck
//
package check

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// -----------------------------------------------------------------------
// Internal type which deals with suite method calling.

const (
	fixtureKd = iota
	testKd
)

type funcKind int

const (
	succeededSt = iota
	failedSt
	skippedSt
	panickedSt
	fixturePanickedSt
	missedSt
)

type funcStatus uint32

// A method value can't reach its own Method structure.
type methodType struct ***REMOVED***
	reflect.Value
	Info reflect.Method
***REMOVED***

func newMethod(receiver reflect.Value, i int) *methodType ***REMOVED***
	return &methodType***REMOVED***receiver.Method(i), receiver.Type().Method(i)***REMOVED***
***REMOVED***

func (method *methodType) PC() uintptr ***REMOVED***
	return method.Info.Func.Pointer()
***REMOVED***

func (method *methodType) suiteName() string ***REMOVED***
	t := method.Info.Type.In(0)
	if t.Kind() == reflect.Ptr ***REMOVED***
		t = t.Elem()
	***REMOVED***
	return t.Name()
***REMOVED***

func (method *methodType) String() string ***REMOVED***
	return method.suiteName() + "." + method.Info.Name
***REMOVED***

func (method *methodType) matches(re *regexp.Regexp) bool ***REMOVED***
	return (re.MatchString(method.Info.Name) ||
		re.MatchString(method.suiteName()) ||
		re.MatchString(method.String()))
***REMOVED***

type C struct ***REMOVED***
	method    *methodType
	kind      funcKind
	testName  string
	_status   funcStatus
	logb      *logger
	logw      io.Writer
	done      chan *C
	reason    string
	mustFail  bool
	tempDir   *tempDir
	benchMem  bool
	startTime time.Time
	timer
***REMOVED***

func (c *C) status() funcStatus ***REMOVED***
	return funcStatus(atomic.LoadUint32((*uint32)(&c._status)))
***REMOVED***

func (c *C) setStatus(s funcStatus) ***REMOVED***
	atomic.StoreUint32((*uint32)(&c._status), uint32(s))
***REMOVED***

func (c *C) stopNow() ***REMOVED***
	runtime.Goexit()
***REMOVED***

// logger is a concurrency safe byte.Buffer
type logger struct ***REMOVED***
	sync.Mutex
	writer bytes.Buffer
***REMOVED***

func (l *logger) Write(buf []byte) (int, error) ***REMOVED***
	l.Lock()
	defer l.Unlock()
	return l.writer.Write(buf)
***REMOVED***

func (l *logger) WriteTo(w io.Writer) (int64, error) ***REMOVED***
	l.Lock()
	defer l.Unlock()
	return l.writer.WriteTo(w)
***REMOVED***

func (l *logger) String() string ***REMOVED***
	l.Lock()
	defer l.Unlock()
	return l.writer.String()
***REMOVED***

// -----------------------------------------------------------------------
// Handling of temporary files and directories.

type tempDir struct ***REMOVED***
	sync.Mutex
	path    string
	counter int
***REMOVED***

func (td *tempDir) newPath() string ***REMOVED***
	td.Lock()
	defer td.Unlock()
	if td.path == "" ***REMOVED***
		var err error
		for i := 0; i != 100; i++ ***REMOVED***
			path := fmt.Sprintf("%s%ccheck-%d", os.TempDir(), os.PathSeparator, rand.Int())
			if err = os.Mkdir(path, 0700); err == nil ***REMOVED***
				td.path = path
				break
			***REMOVED***
		***REMOVED***
		if td.path == "" ***REMOVED***
			panic("Couldn't create temporary directory: " + err.Error())
		***REMOVED***
	***REMOVED***
	result := filepath.Join(td.path, strconv.Itoa(td.counter))
	td.counter += 1
	return result
***REMOVED***

func (td *tempDir) removeAll() ***REMOVED***
	td.Lock()
	defer td.Unlock()
	if td.path != "" ***REMOVED***
		err := os.RemoveAll(td.path)
		if err != nil ***REMOVED***
			fmt.Fprintf(os.Stderr, "WARNING: Error cleaning up temporaries: "+err.Error())
		***REMOVED***
	***REMOVED***
***REMOVED***

// Create a new temporary directory which is automatically removed after
// the suite finishes running.
func (c *C) MkDir() string ***REMOVED***
	path := c.tempDir.newPath()
	if err := os.Mkdir(path, 0700); err != nil ***REMOVED***
		panic(fmt.Sprintf("Couldn't create temporary directory %s: %s", path, err.Error()))
	***REMOVED***
	return path
***REMOVED***

// -----------------------------------------------------------------------
// Low-level logging functions.

func (c *C) log(args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.writeLog([]byte(fmt.Sprint(args...) + "\n"))
***REMOVED***

func (c *C) logf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.writeLog([]byte(fmt.Sprintf(format+"\n", args...)))
***REMOVED***

func (c *C) logNewLine() ***REMOVED***
	c.writeLog([]byte***REMOVED***'\n'***REMOVED***)
***REMOVED***

func (c *C) writeLog(buf []byte) ***REMOVED***
	c.logb.Write(buf)
	if c.logw != nil ***REMOVED***
		c.logw.Write(buf)
	***REMOVED***
***REMOVED***

func hasStringOrError(x interface***REMOVED******REMOVED***) (ok bool) ***REMOVED***
	_, ok = x.(fmt.Stringer)
	if ok ***REMOVED***
		return
	***REMOVED***
	_, ok = x.(error)
	return
***REMOVED***

func (c *C) logValue(label string, value interface***REMOVED******REMOVED***) ***REMOVED***
	if label == "" ***REMOVED***
		if hasStringOrError(value) ***REMOVED***
			c.logf("... %#v (%q)", value, value)
		***REMOVED*** else ***REMOVED***
			c.logf("... %#v", value)
		***REMOVED***
	***REMOVED*** else if value == nil ***REMOVED***
		c.logf("... %s = nil", label)
	***REMOVED*** else ***REMOVED***
		if hasStringOrError(value) ***REMOVED***
			fv := fmt.Sprintf("%#v", value)
			qv := fmt.Sprintf("%q", value)
			if fv != qv ***REMOVED***
				c.logf("... %s %s = %s (%s)", label, reflect.TypeOf(value), fv, qv)
				return
			***REMOVED***
		***REMOVED***
		if s, ok := value.(string); ok && isMultiLine(s) ***REMOVED***
			c.logf(`... %s %s = "" +`, label, reflect.TypeOf(value))
			c.logMultiLine(s)
		***REMOVED*** else ***REMOVED***
			c.logf("... %s %s = %#v", label, reflect.TypeOf(value), value)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *C) logMultiLine(s string) ***REMOVED***
	b := make([]byte, 0, len(s)*2)
	i := 0
	n := len(s)
	for i < n ***REMOVED***
		j := i + 1
		for j < n && s[j-1] != '\n' ***REMOVED***
			j++
		***REMOVED***
		b = append(b, "...     "...)
		b = strconv.AppendQuote(b, s[i:j])
		if j < n ***REMOVED***
			b = append(b, " +"...)
		***REMOVED***
		b = append(b, '\n')
		i = j
	***REMOVED***
	c.writeLog(b)
***REMOVED***

func isMultiLine(s string) bool ***REMOVED***
	for i := 0; i+1 < len(s); i++ ***REMOVED***
		if s[i] == '\n' ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (c *C) logString(issue string) ***REMOVED***
	c.log("... ", issue)
***REMOVED***

func (c *C) logCaller(skip int) ***REMOVED***
	// This is a bit heavier than it ought to be.
	skip += 1 // Our own frame.
	pc, callerFile, callerLine, ok := runtime.Caller(skip)
	if !ok ***REMOVED***
		return
	***REMOVED***
	var testFile string
	var testLine int
	testFunc := runtime.FuncForPC(c.method.PC())
	if runtime.FuncForPC(pc) != testFunc ***REMOVED***
		for ***REMOVED***
			skip += 1
			if pc, file, line, ok := runtime.Caller(skip); ok ***REMOVED***
				// Note that the test line may be different on
				// distinct calls for the same test.  Showing
				// the "internal" line is helpful when debugging.
				if runtime.FuncForPC(pc) == testFunc ***REMOVED***
					testFile, testLine = file, line
					break
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if testFile != "" && (testFile != callerFile || testLine != callerLine) ***REMOVED***
		c.logCode(testFile, testLine)
	***REMOVED***
	c.logCode(callerFile, callerLine)
***REMOVED***

func (c *C) logCode(path string, line int) ***REMOVED***
	c.logf("%s:%d:", nicePath(path), line)
	code, err := printLine(path, line)
	if code == "" ***REMOVED***
		code = "..." // XXX Open the file and take the raw line.
		if err != nil ***REMOVED***
			code += err.Error()
		***REMOVED***
	***REMOVED***
	c.log(indent(code, "    "))
***REMOVED***

var valueGo = filepath.Join("reflect", "value.go")
var asmGo = filepath.Join("runtime", "asm_")

func (c *C) logPanic(skip int, value interface***REMOVED******REMOVED***) ***REMOVED***
	skip++ // Our own frame.
	initialSkip := skip
	for ; ; skip++ ***REMOVED***
		if pc, file, line, ok := runtime.Caller(skip); ok ***REMOVED***
			if skip == initialSkip ***REMOVED***
				c.logf("... Panic: %s (PC=0x%X)\n", value, pc)
			***REMOVED***
			name := niceFuncName(pc)
			path := nicePath(file)
			if strings.Contains(path, "/gopkg.in/check.v") ***REMOVED***
				continue
			***REMOVED***
			if name == "Value.call" && strings.HasSuffix(path, valueGo) ***REMOVED***
				continue
			***REMOVED***
			if (name == "call16" || name == "call32") && strings.Contains(path, asmGo) ***REMOVED***
				continue
			***REMOVED***
			c.logf("%s:%d\n  in %s", nicePath(file), line, name)
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *C) logSoftPanic(issue string) ***REMOVED***
	c.log("... Panic: ", issue)
***REMOVED***

func (c *C) logArgPanic(method *methodType, expectedType string) ***REMOVED***
	c.logf("... Panic: %s argument should be %s",
		niceFuncName(method.PC()), expectedType)
***REMOVED***

// -----------------------------------------------------------------------
// Some simple formatting helpers.

var initWD, initWDErr = os.Getwd()

func init() ***REMOVED***
	if initWDErr == nil ***REMOVED***
		initWD = strings.Replace(initWD, "\\", "/", -1) + "/"
	***REMOVED***
***REMOVED***

func nicePath(path string) string ***REMOVED***
	if initWDErr == nil ***REMOVED***
		if strings.HasPrefix(path, initWD) ***REMOVED***
			return path[len(initWD):]
		***REMOVED***
	***REMOVED***
	return path
***REMOVED***

func niceFuncPath(pc uintptr) string ***REMOVED***
	function := runtime.FuncForPC(pc)
	if function != nil ***REMOVED***
		filename, line := function.FileLine(pc)
		return fmt.Sprintf("%s:%d", nicePath(filename), line)
	***REMOVED***
	return "<unknown path>"
***REMOVED***

func niceFuncName(pc uintptr) string ***REMOVED***
	function := runtime.FuncForPC(pc)
	if function != nil ***REMOVED***
		name := path.Base(function.Name())
		if i := strings.Index(name, "."); i > 0 ***REMOVED***
			name = name[i+1:]
		***REMOVED***
		if strings.HasPrefix(name, "(*") ***REMOVED***
			if i := strings.Index(name, ")"); i > 0 ***REMOVED***
				name = name[2:i] + name[i+1:]
			***REMOVED***
		***REMOVED***
		if i := strings.LastIndex(name, ".*"); i != -1 ***REMOVED***
			name = name[:i] + "." + name[i+2:]
		***REMOVED***
		if i := strings.LastIndex(name, "·"); i != -1 ***REMOVED***
			name = name[:i] + "." + name[i+2:]
		***REMOVED***
		return name
	***REMOVED***
	return "<unknown function>"
***REMOVED***

// -----------------------------------------------------------------------
// Result tracker to aggregate call results.

type Result struct ***REMOVED***
	Succeeded        int
	Failed           int
	Skipped          int
	Panicked         int
	FixturePanicked  int
	ExpectedFailures int
	Missed           int    // Not even tried to run, related to a panic in the fixture.
	RunError         error  // Houston, we've got a problem.
	WorkDir          string // If KeepWorkDir is true
***REMOVED***

type resultTracker struct ***REMOVED***
	result          Result
	_lastWasProblem bool
	_waiting        int
	_missed         int
	_expectChan     chan *C
	_doneChan       chan *C
	_stopChan       chan bool
***REMOVED***

func newResultTracker() *resultTracker ***REMOVED***
	return &resultTracker***REMOVED***_expectChan: make(chan *C), // Synchronous
		_doneChan: make(chan *C, 32), // Asynchronous
		_stopChan: make(chan bool)***REMOVED***   // Synchronous
***REMOVED***

func (tracker *resultTracker) start() ***REMOVED***
	go tracker._loopRoutine()
***REMOVED***

func (tracker *resultTracker) waitAndStop() ***REMOVED***
	<-tracker._stopChan
***REMOVED***

func (tracker *resultTracker) expectCall(c *C) ***REMOVED***
	tracker._expectChan <- c
***REMOVED***

func (tracker *resultTracker) callDone(c *C) ***REMOVED***
	tracker._doneChan <- c
***REMOVED***

func (tracker *resultTracker) _loopRoutine() ***REMOVED***
	for ***REMOVED***
		var c *C
		if tracker._waiting > 0 ***REMOVED***
			// Calls still running. Can't stop.
			select ***REMOVED***
			// XXX Reindent this (not now to make diff clear)
			case c = <-tracker._expectChan:
				tracker._waiting += 1
			case c = <-tracker._doneChan:
				tracker._waiting -= 1
				switch c.status() ***REMOVED***
				case succeededSt:
					if c.kind == testKd ***REMOVED***
						if c.mustFail ***REMOVED***
							tracker.result.ExpectedFailures++
						***REMOVED*** else ***REMOVED***
							tracker.result.Succeeded++
						***REMOVED***
					***REMOVED***
				case failedSt:
					tracker.result.Failed++
				case panickedSt:
					if c.kind == fixtureKd ***REMOVED***
						tracker.result.FixturePanicked++
					***REMOVED*** else ***REMOVED***
						tracker.result.Panicked++
					***REMOVED***
				case fixturePanickedSt:
					// Track it as missed, since the panic
					// was on the fixture, not on the test.
					tracker.result.Missed++
				case missedSt:
					tracker.result.Missed++
				case skippedSt:
					if c.kind == testKd ***REMOVED***
						tracker.result.Skipped++
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// No calls.  Can stop, but no done calls here.
			select ***REMOVED***
			case tracker._stopChan <- true:
				return
			case c = <-tracker._expectChan:
				tracker._waiting += 1
			case c = <-tracker._doneChan:
				panic("Tracker got an unexpected done call.")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// -----------------------------------------------------------------------
// The underlying suite runner.

type suiteRunner struct ***REMOVED***
	suite                     interface***REMOVED******REMOVED***
	setUpSuite, tearDownSuite *methodType
	setUpTest, tearDownTest   *methodType
	onTimeout                 *methodType
	tests                     []*methodType
	tracker                   *resultTracker
	tempDir                   *tempDir
	keepDir                   bool
	output                    *outputWriter
	reportedProblemLast       bool
	benchTime                 time.Duration
	benchMem                  bool
	checkTimeout              time.Duration
***REMOVED***

type RunConf struct ***REMOVED***
	Output        io.Writer
	Stream        bool
	Verbose       bool
	Filter        string
	Benchmark     bool
	BenchmarkTime time.Duration // Defaults to 1 second
	BenchmarkMem  bool
	KeepWorkDir   bool
	CheckTimeout  time.Duration
***REMOVED***

// Create a new suiteRunner able to run all methods in the given suite.
func newSuiteRunner(suite interface***REMOVED******REMOVED***, runConf *RunConf) *suiteRunner ***REMOVED***
	var conf RunConf
	if runConf != nil ***REMOVED***
		conf = *runConf
	***REMOVED***
	if conf.Output == nil ***REMOVED***
		conf.Output = os.Stdout
	***REMOVED***
	if conf.Benchmark ***REMOVED***
		conf.Verbose = true
	***REMOVED***

	suiteType := reflect.TypeOf(suite)
	suiteNumMethods := suiteType.NumMethod()
	suiteValue := reflect.ValueOf(suite)

	runner := &suiteRunner***REMOVED***
		suite:        suite,
		output:       newOutputWriter(conf.Output, conf.Stream, conf.Verbose),
		tracker:      newResultTracker(),
		benchTime:    conf.BenchmarkTime,
		benchMem:     conf.BenchmarkMem,
		tempDir:      &tempDir***REMOVED******REMOVED***,
		keepDir:      conf.KeepWorkDir,
		tests:        make([]*methodType, 0, suiteNumMethods),
		checkTimeout: conf.CheckTimeout,
	***REMOVED***
	if runner.benchTime == 0 ***REMOVED***
		runner.benchTime = 1 * time.Second
	***REMOVED***

	var filterRegexp *regexp.Regexp
	if conf.Filter != "" ***REMOVED***
		if regexp, err := regexp.Compile(conf.Filter); err != nil ***REMOVED***
			msg := "Bad filter expression: " + err.Error()
			runner.tracker.result.RunError = errors.New(msg)
			return runner
		***REMOVED*** else ***REMOVED***
			filterRegexp = regexp
		***REMOVED***
	***REMOVED***

	for i := 0; i != suiteNumMethods; i++ ***REMOVED***
		method := newMethod(suiteValue, i)
		switch method.Info.Name ***REMOVED***
		case "SetUpSuite":
			runner.setUpSuite = method
		case "TearDownSuite":
			runner.tearDownSuite = method
		case "SetUpTest":
			runner.setUpTest = method
		case "TearDownTest":
			runner.tearDownTest = method
		case "OnTimeout":
			runner.onTimeout = method
		default:
			prefix := "Test"
			if conf.Benchmark ***REMOVED***
				prefix = "Benchmark"
			***REMOVED***
			if !strings.HasPrefix(method.Info.Name, prefix) ***REMOVED***
				continue
			***REMOVED***
			if filterRegexp == nil || method.matches(filterRegexp) ***REMOVED***
				runner.tests = append(runner.tests, method)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return runner
***REMOVED***

// Run all methods in the given suite.
func (runner *suiteRunner) run() *Result ***REMOVED***
	if runner.tracker.result.RunError == nil && len(runner.tests) > 0 ***REMOVED***
		runner.tracker.start()
		if runner.checkFixtureArgs() ***REMOVED***
			c := runner.runFixture(runner.setUpSuite, "", nil)
			if c == nil || c.status() == succeededSt ***REMOVED***
				for i := 0; i != len(runner.tests); i++ ***REMOVED***
					c := runner.runTest(runner.tests[i])
					if c.status() == fixturePanickedSt ***REMOVED***
						runner.skipTests(missedSt, runner.tests[i+1:])
						break
					***REMOVED***
				***REMOVED***
			***REMOVED*** else if c != nil && c.status() == skippedSt ***REMOVED***
				runner.skipTests(skippedSt, runner.tests)
			***REMOVED*** else ***REMOVED***
				runner.skipTests(missedSt, runner.tests)
			***REMOVED***
			runner.runFixture(runner.tearDownSuite, "", nil)
		***REMOVED*** else ***REMOVED***
			runner.skipTests(missedSt, runner.tests)
		***REMOVED***
		runner.tracker.waitAndStop()
		if runner.keepDir ***REMOVED***
			runner.tracker.result.WorkDir = runner.tempDir.path
		***REMOVED*** else ***REMOVED***
			runner.tempDir.removeAll()
		***REMOVED***
	***REMOVED***
	return &runner.tracker.result
***REMOVED***

// Create a call object with the given suite method, and fork a
// goroutine with the provided dispatcher for running it.
func (runner *suiteRunner) forkCall(method *methodType, kind funcKind, testName string, logb *logger, dispatcher func(c *C)) *C ***REMOVED***
	var logw io.Writer
	if runner.output.Stream ***REMOVED***
		logw = runner.output
	***REMOVED***
	if logb == nil ***REMOVED***
		logb = new(logger)
	***REMOVED***
	c := &C***REMOVED***
		method:    method,
		kind:      kind,
		testName:  testName,
		logb:      logb,
		logw:      logw,
		tempDir:   runner.tempDir,
		done:      make(chan *C, 1),
		timer:     timer***REMOVED***benchTime: runner.benchTime***REMOVED***,
		startTime: time.Now(),
		benchMem:  runner.benchMem,
	***REMOVED***
	runner.tracker.expectCall(c)
	go (func() ***REMOVED***
		runner.reportCallStarted(c)
		defer runner.callDone(c)
		dispatcher(c)
	***REMOVED***)()
	return c
***REMOVED***

type timeoutErr struct ***REMOVED***
	method *methodType
	t      time.Duration
***REMOVED***

func (e timeoutErr) Error() string ***REMOVED***
	return fmt.Sprintf("%s test timed out after %v", e.method.String(), e.t)
***REMOVED***

func isTimeout(e error) bool ***REMOVED***
	if e == nil ***REMOVED***
		return false
	***REMOVED***
	_, ok := e.(timeoutErr)
	return ok
***REMOVED***

// Same as forkCall(), but wait for call to finish before returning.
func (runner *suiteRunner) runFunc(method *methodType, kind funcKind, testName string, logb *logger, dispatcher func(c *C)) *C ***REMOVED***
	var timeout <-chan time.Time
	if runner.checkTimeout != 0 ***REMOVED***
		timeout = time.After(runner.checkTimeout)
	***REMOVED***
	c := runner.forkCall(method, kind, testName, logb, dispatcher)
	select ***REMOVED***
	case <-c.done:
	case <-timeout:
		if runner.onTimeout != nil ***REMOVED***
			// run the OnTimeout callback, allowing the suite to collect any sort of debug information it can
			// `runFixture` is syncronous, so run this in a separate goroutine with a timeout
			cChan := make(chan *C)
			go func() ***REMOVED***
				cChan <- runner.runFixture(runner.onTimeout, c.testName, c.logb)
			***REMOVED***()
			select ***REMOVED***
			case <-cChan:
			case <-time.After(runner.checkTimeout):
			***REMOVED***
		***REMOVED***
		panic(timeoutErr***REMOVED***method, runner.checkTimeout***REMOVED***)
	***REMOVED***
	return c
***REMOVED***

// Handle a finished call.  If there were any panics, update the call status
// accordingly.  Then, mark the call as done and report to the tracker.
func (runner *suiteRunner) callDone(c *C) ***REMOVED***
	value := recover()
	if value != nil ***REMOVED***
		switch v := value.(type) ***REMOVED***
		case *fixturePanic:
			if v.status == skippedSt ***REMOVED***
				c.setStatus(skippedSt)
			***REMOVED*** else ***REMOVED***
				c.logSoftPanic("Fixture has panicked (see related PANIC)")
				c.setStatus(fixturePanickedSt)
			***REMOVED***
		default:
			c.logPanic(1, value)
			c.setStatus(panickedSt)
		***REMOVED***
	***REMOVED***
	if c.mustFail ***REMOVED***
		switch c.status() ***REMOVED***
		case failedSt:
			c.setStatus(succeededSt)
		case succeededSt:
			c.setStatus(failedSt)
			c.logString("Error: Test succeeded, but was expected to fail")
			c.logString("Reason: " + c.reason)
		***REMOVED***
	***REMOVED***

	runner.reportCallDone(c)
	c.done <- c
***REMOVED***

// Runs a fixture call synchronously.  The fixture will still be run in a
// goroutine like all suite methods, but this method will not return
// while the fixture goroutine is not done, because the fixture must be
// run in a desired order.
func (runner *suiteRunner) runFixture(method *methodType, testName string, logb *logger) *C ***REMOVED***
	if method != nil ***REMOVED***
		c := runner.runFunc(method, fixtureKd, testName, logb, func(c *C) ***REMOVED***
			c.ResetTimer()
			c.StartTimer()
			defer c.StopTimer()
			c.method.Call([]reflect.Value***REMOVED***reflect.ValueOf(c)***REMOVED***)
		***REMOVED***)
		return c
	***REMOVED***
	return nil
***REMOVED***

// Run the fixture method with runFixture(), but panic with a fixturePanic***REMOVED******REMOVED***
// in case the fixture method panics.  This makes it easier to track the
// fixture panic together with other call panics within forkTest().
func (runner *suiteRunner) runFixtureWithPanic(method *methodType, testName string, logb *logger, skipped *bool) *C ***REMOVED***
	if skipped != nil && *skipped ***REMOVED***
		return nil
	***REMOVED***
	c := runner.runFixture(method, testName, logb)
	if c != nil && c.status() != succeededSt ***REMOVED***
		if skipped != nil ***REMOVED***
			*skipped = c.status() == skippedSt
		***REMOVED***
		panic(&fixturePanic***REMOVED***c.status(), method***REMOVED***)
	***REMOVED***
	return c
***REMOVED***

type fixturePanic struct ***REMOVED***
	status funcStatus
	method *methodType
***REMOVED***

// Run the suite test method, together with the test-specific fixture,
// asynchronously.
func (runner *suiteRunner) forkTest(method *methodType) *C ***REMOVED***
	testName := method.String()
	return runner.forkCall(method, testKd, testName, nil, func(c *C) ***REMOVED***
		var skipped bool
		defer runner.runFixtureWithPanic(runner.tearDownTest, testName, nil, &skipped)
		defer c.StopTimer()
		benchN := 1
		for ***REMOVED***
			runner.runFixtureWithPanic(runner.setUpTest, testName, c.logb, &skipped)
			mt := c.method.Type()
			if mt.NumIn() != 1 || mt.In(0) != reflect.TypeOf(c) ***REMOVED***
				// Rather than a plain panic, provide a more helpful message when
				// the argument type is incorrect.
				c.setStatus(panickedSt)
				c.logArgPanic(c.method, "*check.C")
				return
			***REMOVED***

			if strings.HasPrefix(c.method.Info.Name, "Test") ***REMOVED***
				c.ResetTimer()
				c.StartTimer()
				c.method.Call([]reflect.Value***REMOVED***reflect.ValueOf(c)***REMOVED***)
				return
			***REMOVED***

			if !strings.HasPrefix(c.method.Info.Name, "Benchmark") ***REMOVED***
				panic("unexpected method prefix: " + c.method.Info.Name)
			***REMOVED***

			runtime.GC()
			c.N = benchN
			c.ResetTimer()
			c.StartTimer()

			c.method.Call([]reflect.Value***REMOVED***reflect.ValueOf(c)***REMOVED***)
			c.StopTimer()
			if c.status() != succeededSt || c.duration >= c.benchTime || benchN >= 1e9 ***REMOVED***
				return
			***REMOVED***
			perOpN := int(1e9)
			if c.nsPerOp() != 0 ***REMOVED***
				perOpN = int(c.benchTime.Nanoseconds() / c.nsPerOp())
			***REMOVED***

			// Logic taken from the stock testing package:
			// - Run more iterations than we think we'll need for a second (1.5x).
			// - Don't grow too fast in case we had timing errors previously.
			// - Be sure to run at least one more than last time.
			benchN = max(min(perOpN+perOpN/2, 100*benchN), benchN+1)
			benchN = roundUp(benchN)

			skipped = true // Don't run the deferred one if this panics.
			runner.runFixtureWithPanic(runner.tearDownTest, testName, nil, nil)
			skipped = false
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Same as forkTest(), but wait for the test to finish before returning.
func (runner *suiteRunner) runTest(method *methodType) *C ***REMOVED***
	var timeout <-chan time.Time
	if runner.checkTimeout != 0 ***REMOVED***
		timeout = time.After(runner.checkTimeout)
	***REMOVED***
	c := runner.forkTest(method)
	select ***REMOVED***
	case <-c.done:
	case <-timeout:
		if runner.onTimeout != nil ***REMOVED***
			// run the OnTimeout callback, allowing the suite to collect any sort of debug information it can
			// `runFixture` is syncronous, so run this in a separate goroutine with a timeout
			cChan := make(chan *C)
			go func() ***REMOVED***
				cChan <- runner.runFixture(runner.onTimeout, c.testName, c.logb)
			***REMOVED***()
			select ***REMOVED***
			case <-cChan:
			case <-time.After(runner.checkTimeout):
			***REMOVED***
		***REMOVED***
		panic(timeoutErr***REMOVED***method, runner.checkTimeout***REMOVED***)
	***REMOVED***
	return c
***REMOVED***

// Helper to mark tests as skipped or missed.  A bit heavy for what
// it does, but it enables homogeneous handling of tracking, including
// nice verbose output.
func (runner *suiteRunner) skipTests(status funcStatus, methods []*methodType) ***REMOVED***
	for _, method := range methods ***REMOVED***
		runner.runFunc(method, testKd, "", nil, func(c *C) ***REMOVED***
			c.setStatus(status)
		***REMOVED***)
	***REMOVED***
***REMOVED***

// Verify if the fixture arguments are *check.C.  In case of errors,
// log the error as a panic in the fixture method call, and return false.
func (runner *suiteRunner) checkFixtureArgs() bool ***REMOVED***
	succeeded := true
	argType := reflect.TypeOf(&C***REMOVED******REMOVED***)
	for _, method := range []*methodType***REMOVED***runner.setUpSuite, runner.tearDownSuite, runner.setUpTest, runner.tearDownTest, runner.onTimeout***REMOVED*** ***REMOVED***
		if method != nil ***REMOVED***
			mt := method.Type()
			if mt.NumIn() != 1 || mt.In(0) != argType ***REMOVED***
				succeeded = false
				runner.runFunc(method, fixtureKd, "", nil, func(c *C) ***REMOVED***
					c.logArgPanic(method, "*check.C")
					c.setStatus(panickedSt)
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return succeeded
***REMOVED***

func (runner *suiteRunner) reportCallStarted(c *C) ***REMOVED***
	runner.output.WriteCallStarted("START", c)
***REMOVED***

func (runner *suiteRunner) reportCallDone(c *C) ***REMOVED***
	runner.tracker.callDone(c)
	switch c.status() ***REMOVED***
	case succeededSt:
		if c.mustFail ***REMOVED***
			runner.output.WriteCallSuccess("FAIL EXPECTED", c)
		***REMOVED*** else ***REMOVED***
			runner.output.WriteCallSuccess("PASS", c)
		***REMOVED***
	case skippedSt:
		runner.output.WriteCallSuccess("SKIP", c)
	case failedSt:
		runner.output.WriteCallProblem("FAIL", c)
	case panickedSt:
		runner.output.WriteCallProblem("PANIC", c)
	case fixturePanickedSt:
		// That's a testKd call reporting that its fixture
		// has panicked. The fixture call which caused the
		// panic itself was tracked above. We'll report to
		// aid debugging.
		runner.output.WriteCallProblem("PANIC", c)
	case missedSt:
		runner.output.WriteCallSuccess("MISS", c)
	***REMOVED***
***REMOVED***
