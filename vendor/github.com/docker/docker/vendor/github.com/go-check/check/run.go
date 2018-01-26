package check

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
)

// -----------------------------------------------------------------------
// Test suite registry.

var allSuites []interface***REMOVED******REMOVED***

// Suite registers the given value as a test suite to be run. Any methods
// starting with the Test prefix in the given value will be considered as
// a test method.
func Suite(suite interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	allSuites = append(allSuites, suite)
	return suite
***REMOVED***

// -----------------------------------------------------------------------
// Public running interface.

var (
	oldFilterFlag  = flag.String("gocheck.f", "", "Regular expression selecting which tests and/or suites to run")
	oldVerboseFlag = flag.Bool("gocheck.v", false, "Verbose mode")
	oldStreamFlag  = flag.Bool("gocheck.vv", false, "Super verbose mode (disables output caching)")
	oldBenchFlag   = flag.Bool("gocheck.b", false, "Run benchmarks")
	oldBenchTime   = flag.Duration("gocheck.btime", 1*time.Second, "approximate run time for each benchmark")
	oldListFlag    = flag.Bool("gocheck.list", false, "List the names of all tests that will be run")
	oldWorkFlag    = flag.Bool("gocheck.work", false, "Display and do not remove the test working directory")

	newFilterFlag  = flag.String("check.f", "", "Regular expression selecting which tests and/or suites to run")
	newVerboseFlag = flag.Bool("check.v", false, "Verbose mode")
	newStreamFlag  = flag.Bool("check.vv", false, "Super verbose mode (disables output caching)")
	newBenchFlag   = flag.Bool("check.b", false, "Run benchmarks")
	newBenchTime   = flag.Duration("check.btime", 1*time.Second, "approximate run time for each benchmark")
	newBenchMem    = flag.Bool("check.bmem", false, "Report memory benchmarks")
	newListFlag    = flag.Bool("check.list", false, "List the names of all tests that will be run")
	newWorkFlag    = flag.Bool("check.work", false, "Display and do not remove the test working directory")
	checkTimeout   = flag.String("check.timeout", "", "Panic if test runs longer than specified duration")
)

// TestingT runs all test suites registered with the Suite function,
// printing results to stdout, and reporting any failures back to
// the "testing" package.
func TestingT(testingT *testing.T) ***REMOVED***
	benchTime := *newBenchTime
	if benchTime == 1*time.Second ***REMOVED***
		benchTime = *oldBenchTime
	***REMOVED***
	conf := &RunConf***REMOVED***
		Filter:        *oldFilterFlag + *newFilterFlag,
		Verbose:       *oldVerboseFlag || *newVerboseFlag,
		Stream:        *oldStreamFlag || *newStreamFlag,
		Benchmark:     *oldBenchFlag || *newBenchFlag,
		BenchmarkTime: benchTime,
		BenchmarkMem:  *newBenchMem,
		KeepWorkDir:   *oldWorkFlag || *newWorkFlag,
	***REMOVED***
	if *checkTimeout != "" ***REMOVED***
		timeout, err := time.ParseDuration(*checkTimeout)
		if err != nil ***REMOVED***
			testingT.Fatalf("error parsing specified timeout flag: %v", err)
		***REMOVED***
		conf.CheckTimeout = timeout
	***REMOVED***
	if *oldListFlag || *newListFlag ***REMOVED***
		w := bufio.NewWriter(os.Stdout)
		for _, name := range ListAll(conf) ***REMOVED***
			fmt.Fprintln(w, name)
		***REMOVED***
		w.Flush()
		return
	***REMOVED***
	result := RunAll(conf)
	println(result.String())
	if !result.Passed() ***REMOVED***
		testingT.Fail()
	***REMOVED***
***REMOVED***

// RunAll runs all test suites registered with the Suite function, using the
// provided run configuration.
func RunAll(runConf *RunConf) *Result ***REMOVED***
	result := Result***REMOVED******REMOVED***
	for _, suite := range allSuites ***REMOVED***
		result.Add(Run(suite, runConf))
	***REMOVED***
	return &result
***REMOVED***

// Run runs the provided test suite using the provided run configuration.
func Run(suite interface***REMOVED******REMOVED***, runConf *RunConf) *Result ***REMOVED***
	runner := newSuiteRunner(suite, runConf)
	return runner.run()
***REMOVED***

// ListAll returns the names of all the test functions registered with the
// Suite function that will be run with the provided run configuration.
func ListAll(runConf *RunConf) []string ***REMOVED***
	var names []string
	for _, suite := range allSuites ***REMOVED***
		names = append(names, List(suite, runConf)...)
	***REMOVED***
	return names
***REMOVED***

// List returns the names of the test functions in the given
// suite that will be run with the provided run configuration.
func List(suite interface***REMOVED******REMOVED***, runConf *RunConf) []string ***REMOVED***
	var names []string
	runner := newSuiteRunner(suite, runConf)
	for _, t := range runner.tests ***REMOVED***
		names = append(names, t.String())
	***REMOVED***
	return names
***REMOVED***

// -----------------------------------------------------------------------
// Result methods.

func (r *Result) Add(other *Result) ***REMOVED***
	r.Succeeded += other.Succeeded
	r.Skipped += other.Skipped
	r.Failed += other.Failed
	r.Panicked += other.Panicked
	r.FixturePanicked += other.FixturePanicked
	r.ExpectedFailures += other.ExpectedFailures
	r.Missed += other.Missed
	if r.WorkDir != "" && other.WorkDir != "" ***REMOVED***
		r.WorkDir += ":" + other.WorkDir
	***REMOVED*** else if other.WorkDir != "" ***REMOVED***
		r.WorkDir = other.WorkDir
	***REMOVED***
***REMOVED***

func (r *Result) Passed() bool ***REMOVED***
	return (r.Failed == 0 && r.Panicked == 0 &&
		r.FixturePanicked == 0 && r.Missed == 0 &&
		r.RunError == nil)
***REMOVED***

func (r *Result) String() string ***REMOVED***
	if r.RunError != nil ***REMOVED***
		return "ERROR: " + r.RunError.Error()
	***REMOVED***

	var value string
	if r.Failed == 0 && r.Panicked == 0 && r.FixturePanicked == 0 &&
		r.Missed == 0 ***REMOVED***
		value = "OK: "
	***REMOVED*** else ***REMOVED***
		value = "OOPS: "
	***REMOVED***
	value += fmt.Sprintf("%d passed", r.Succeeded)
	if r.Skipped != 0 ***REMOVED***
		value += fmt.Sprintf(", %d skipped", r.Skipped)
	***REMOVED***
	if r.ExpectedFailures != 0 ***REMOVED***
		value += fmt.Sprintf(", %d expected failures", r.ExpectedFailures)
	***REMOVED***
	if r.Failed != 0 ***REMOVED***
		value += fmt.Sprintf(", %d FAILED", r.Failed)
	***REMOVED***
	if r.Panicked != 0 ***REMOVED***
		value += fmt.Sprintf(", %d PANICKED", r.Panicked)
	***REMOVED***
	if r.FixturePanicked != 0 ***REMOVED***
		value += fmt.Sprintf(", %d FIXTURE-PANICKED", r.FixturePanicked)
	***REMOVED***
	if r.Missed != 0 ***REMOVED***
		value += fmt.Sprintf(", %d MISSED", r.Missed)
	***REMOVED***
	if r.WorkDir != "" ***REMOVED***
		value += "\nWORK=" + r.WorkDir
	***REMOVED***
	return value
***REMOVED***
