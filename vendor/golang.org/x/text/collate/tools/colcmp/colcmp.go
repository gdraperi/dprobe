// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "golang.org/x/text/collate/tools/colcmp"

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/unicode/norm"
)

var (
	doNorm  = flag.Bool("norm", false, "normalize input strings")
	cases   = flag.Bool("case", false, "generate case variants")
	verbose = flag.Bool("verbose", false, "print results")
	debug   = flag.Bool("debug", false, "output debug information")
	locales = flag.String("locale", "en_US", "the locale to use. May be a comma-separated list for some commands.")
	col     = flag.String("col", "go", "collator to test")
	gold    = flag.String("gold", "go", "collator used as the gold standard")
	usecmp  = flag.Bool("usecmp", false,
		`use comparison instead of sort keys when sorting.  Must be "test", "gold" or "both"`)
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	exclude    = flag.String("exclude", "", "exclude errors that contain any of the characters")
	limit      = flag.Int("limit", 5000000, "maximum number of samples to generate for one run")
)

func failOnError(err error) ***REMOVED***
	if err != nil ***REMOVED***
		log.Panic(err)
	***REMOVED***
***REMOVED***

// Test holds test data for testing a locale-collator pair.
// Test also provides functionality that is commonly used by the various commands.
type Test struct ***REMOVED***
	ctxt    *Context
	Name    string
	Locale  string
	ColName string

	Col        Collator
	UseCompare bool

	Input    []Input
	Duration time.Duration

	start time.Time
	msg   string
	count int
***REMOVED***

func (t *Test) clear() ***REMOVED***
	t.Col = nil
	t.Input = nil
***REMOVED***

const (
	msgGeneratingInput = "generating input"
	msgGeneratingKeys  = "generating keys"
	msgSorting         = "sorting"
)

var lastLen = 0

func (t *Test) SetStatus(msg string) ***REMOVED***
	if *debug || *verbose ***REMOVED***
		fmt.Printf("%s: %s...\n", t.Name, msg)
	***REMOVED*** else if t.ctxt.out != nil ***REMOVED***
		fmt.Fprint(t.ctxt.out, strings.Repeat(" ", lastLen))
		fmt.Fprint(t.ctxt.out, strings.Repeat("\b", lastLen))
		fmt.Fprint(t.ctxt.out, msg, "...")
		lastLen = len(msg) + 3
		fmt.Fprint(t.ctxt.out, strings.Repeat("\b", lastLen))
	***REMOVED***
***REMOVED***

// Start is used by commands to signal the start of an operation.
func (t *Test) Start(msg string) ***REMOVED***
	t.SetStatus(msg)
	t.count = 0
	t.msg = msg
	t.start = time.Now()
***REMOVED***

// Stop is used by commands to signal the end of an operation.
func (t *Test) Stop() (time.Duration, int) ***REMOVED***
	d := time.Now().Sub(t.start)
	t.Duration += d
	if *debug || *verbose ***REMOVED***
		fmt.Printf("%s: %s done. (%.3fs /%dK ops)\n", t.Name, t.msg, d.Seconds(), t.count/1000)
	***REMOVED***
	return d, t.count
***REMOVED***

// generateKeys generates sort keys for all the inputs.
func (t *Test) generateKeys() ***REMOVED***
	for i, s := range t.Input ***REMOVED***
		b := t.Col.Key(s)
		t.Input[i].key = b
		if *debug ***REMOVED***
			fmt.Printf("%s (%X): %X\n", string(s.UTF8), s.UTF16, b)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Sort sorts the inputs. It generates sort keys if this is required by the
// chosen sort method.
func (t *Test) Sort() (tkey, tsort time.Duration, nkey, nsort int) ***REMOVED***
	if *cpuprofile != "" ***REMOVED***
		f, err := os.Create(*cpuprofile)
		failOnError(err)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	***REMOVED***
	if t.UseCompare || t.Col.Key(t.Input[0]) == nil ***REMOVED***
		t.Start(msgSorting)
		sort.Sort(&testCompare***REMOVED****t***REMOVED***)
		tsort, nsort = t.Stop()
	***REMOVED*** else ***REMOVED***
		t.Start(msgGeneratingKeys)
		t.generateKeys()
		t.count = len(t.Input)
		tkey, nkey = t.Stop()
		t.Start(msgSorting)
		sort.Sort(t)
		tsort, nsort = t.Stop()
	***REMOVED***
	return
***REMOVED***

func (t *Test) Swap(a, b int) ***REMOVED***
	t.Input[a], t.Input[b] = t.Input[b], t.Input[a]
***REMOVED***

func (t *Test) Less(a, b int) bool ***REMOVED***
	t.count++
	return bytes.Compare(t.Input[a].key, t.Input[b].key) == -1
***REMOVED***

func (t Test) Len() int ***REMOVED***
	return len(t.Input)
***REMOVED***

type testCompare struct ***REMOVED***
	Test
***REMOVED***

func (t *testCompare) Less(a, b int) bool ***REMOVED***
	t.count++
	return t.Col.Compare(t.Input[a], t.Input[b]) == -1
***REMOVED***

type testRestore struct ***REMOVED***
	Test
***REMOVED***

func (t *testRestore) Less(a, b int) bool ***REMOVED***
	return t.Input[a].index < t.Input[b].index
***REMOVED***

// GenerateInput generates input phrases for the locale tested by t.
func (t *Test) GenerateInput() ***REMOVED***
	t.Input = nil
	if t.ctxt.lastLocale != t.Locale ***REMOVED***
		gen := phraseGenerator***REMOVED******REMOVED***
		gen.init(t.Locale)
		t.SetStatus(msgGeneratingInput)
		t.ctxt.lastInput = nil // allow the previous value to be garbage collected.
		t.Input = gen.generate(*doNorm)
		t.ctxt.lastInput = t.Input
		t.ctxt.lastLocale = t.Locale
	***REMOVED*** else ***REMOVED***
		t.Input = t.ctxt.lastInput
		for i := range t.Input ***REMOVED***
			t.Input[i].key = nil
		***REMOVED***
		sort.Sort(&testRestore***REMOVED****t***REMOVED***)
	***REMOVED***
***REMOVED***

// Context holds all tests and settings translated from command line options.
type Context struct ***REMOVED***
	test []*Test
	last *Test

	lastLocale string
	lastInput  []Input

	out io.Writer
***REMOVED***

func (ts *Context) Printf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	ts.assertBuf()
	fmt.Fprintf(ts.out, format, a...)
***REMOVED***

func (ts *Context) Print(a ...interface***REMOVED******REMOVED***) ***REMOVED***
	ts.assertBuf()
	fmt.Fprint(ts.out, a...)
***REMOVED***

// assertBuf sets up an io.Writer for output, if it doesn't already exist.
// In debug and verbose mode, output is buffered so that the regular output
// will not interfere with the additional output.  Otherwise, output is
// written directly to stdout for a more responsive feel.
func (ts *Context) assertBuf() ***REMOVED***
	if ts.out != nil ***REMOVED***
		return
	***REMOVED***
	if *debug || *verbose ***REMOVED***
		ts.out = &bytes.Buffer***REMOVED******REMOVED***
	***REMOVED*** else ***REMOVED***
		ts.out = os.Stdout
	***REMOVED***
***REMOVED***

// flush flushes the contents of ts.out to stdout, if it is not stdout already.
func (ts *Context) flush() ***REMOVED***
	if ts.out != nil ***REMOVED***
		if _, ok := ts.out.(io.ReadCloser); !ok ***REMOVED***
			io.Copy(os.Stdout, ts.out.(io.Reader))
		***REMOVED***
	***REMOVED***
***REMOVED***

// parseTests creates all tests from command lines and returns
// a Context to hold them.
func parseTests() *Context ***REMOVED***
	ctxt := &Context***REMOVED******REMOVED***
	colls := strings.Split(*col, ",")
	for _, loc := range strings.Split(*locales, ",") ***REMOVED***
		loc = strings.TrimSpace(loc)
		for _, name := range colls ***REMOVED***
			name = strings.TrimSpace(name)
			col := getCollator(name, loc)
			ctxt.test = append(ctxt.test, &Test***REMOVED***
				ctxt:       ctxt,
				Locale:     loc,
				ColName:    name,
				UseCompare: *usecmp,
				Col:        col,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return ctxt
***REMOVED***

func (c *Context) Len() int ***REMOVED***
	return len(c.test)
***REMOVED***

func (c *Context) Test(i int) *Test ***REMOVED***
	if c.last != nil ***REMOVED***
		c.last.clear()
	***REMOVED***
	c.last = c.test[i]
	return c.last
***REMOVED***

func parseInput(args []string) []Input ***REMOVED***
	input := []Input***REMOVED******REMOVED***
	for _, s := range args ***REMOVED***
		rs := []rune***REMOVED******REMOVED***
		for len(s) > 0 ***REMOVED***
			var r rune
			r, _, s, _ = strconv.UnquoteChar(s, '\'')
			rs = append(rs, r)
		***REMOVED***
		s = string(rs)
		if *doNorm ***REMOVED***
			s = norm.NFD.String(s)
		***REMOVED***
		input = append(input, makeInputString(s))
	***REMOVED***
	return input
***REMOVED***

// A Command is an implementation of a colcmp command.
type Command struct ***REMOVED***
	Run   func(cmd *Context, args []string)
	Usage string
	Short string
	Long  string
***REMOVED***

func (cmd Command) Name() string ***REMOVED***
	return strings.SplitN(cmd.Usage, " ", 2)[0]
***REMOVED***

var commands = []*Command***REMOVED***
	cmdSort,
	cmdBench,
	cmdRegress,
***REMOVED***

const sortHelp = `
Sort sorts a given list of strings.  Strings are separated by whitespace.
`

var cmdSort = &Command***REMOVED***
	Run:   runSort,
	Usage: "sort <string>*",
	Short: "sort a given list of strings",
	Long:  sortHelp,
***REMOVED***

func runSort(ctxt *Context, args []string) ***REMOVED***
	input := parseInput(args)
	if len(input) == 0 ***REMOVED***
		log.Fatalf("Nothing to sort.")
	***REMOVED***
	if ctxt.Len() > 1 ***REMOVED***
		ctxt.Print("COLL  LOCALE RESULT\n")
	***REMOVED***
	for i := 0; i < ctxt.Len(); i++ ***REMOVED***
		t := ctxt.Test(i)
		t.Input = append(t.Input, input...)
		t.Sort()
		if ctxt.Len() > 1 ***REMOVED***
			ctxt.Printf("%-5s %-5s  ", t.ColName, t.Locale)
		***REMOVED***
		for _, s := range t.Input ***REMOVED***
			ctxt.Print(string(s.UTF8), " ")
		***REMOVED***
		ctxt.Print("\n")
	***REMOVED***
***REMOVED***

const benchHelp = `
Bench runs a benchmark for the given list of collator implementations.
If no collator implementations are given, the go collator will be used.
`

var cmdBench = &Command***REMOVED***
	Run:   runBench,
	Usage: "bench",
	Short: "benchmark a given list of collator implementations",
	Long:  benchHelp,
***REMOVED***

func runBench(ctxt *Context, args []string) ***REMOVED***
	ctxt.Printf("%-7s %-5s %-6s %-24s %-24s %-5s %s\n", "LOCALE", "COLL", "N", "KEYS", "SORT", "AVGLN", "TOTAL")
	for i := 0; i < ctxt.Len(); i++ ***REMOVED***
		t := ctxt.Test(i)
		ctxt.Printf("%-7s %-5s ", t.Locale, t.ColName)
		t.GenerateInput()
		ctxt.Printf("%-6s ", fmt.Sprintf("%dK", t.Len()/1000))
		tkey, tsort, nkey, nsort := t.Sort()
		p := func(dur time.Duration, n int) ***REMOVED***
			s := ""
			if dur > 0 ***REMOVED***
				s = fmt.Sprintf("%6.3fs ", dur.Seconds())
				if n > 0 ***REMOVED***
					s += fmt.Sprintf("%15s", fmt.Sprintf("(%4.2f ns/op)", float64(dur)/float64(n)))
				***REMOVED***
			***REMOVED***
			ctxt.Printf("%-24s ", s)
		***REMOVED***
		p(tkey, nkey)
		p(tsort, nsort)

		total := 0
		for _, s := range t.Input ***REMOVED***
			total += len(s.key)
		***REMOVED***
		ctxt.Printf("%-5d ", total/t.Len())
		ctxt.Printf("%6.3fs\n", t.Duration.Seconds())
		if *debug ***REMOVED***
			for _, s := range t.Input ***REMOVED***
				fmt.Print(string(s.UTF8), " ")
			***REMOVED***
			fmt.Println()
		***REMOVED***
	***REMOVED***
***REMOVED***

const regressHelp = `
Regress runs a monkey test by comparing the results of randomly generated tests
between two implementations of a collator. The user may optionally pass a list
of strings to regress against instead of the default test set.
`

var cmdRegress = &Command***REMOVED***
	Run:   runRegress,
	Usage: "regress -gold=<col> -test=<col> [string]*",
	Short: "run a monkey test between two collators",
	Long:  regressHelp,
***REMOVED***

const failedKeyCompare = `
%s:%d: incorrect comparison result for input:
    a:   %q (%.4X)
    key: %s
    b:   %q (%.4X)
    key: %s
    Compare(a, b) = %d; want %d.

  gold keys:
	a:   %s
	b:   %s
`

const failedCompare = `
%s:%d: incorrect comparison result for input:
    a:   %q (%.4X)
    b:   %q (%.4X)
    Compare(a, b) = %d; want %d.
`

func keyStr(b []byte) string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***
	for _, v := range b ***REMOVED***
		fmt.Fprintf(buf, "%.2X ", v)
	***REMOVED***
	return buf.String()
***REMOVED***

func runRegress(ctxt *Context, args []string) ***REMOVED***
	input := parseInput(args)
	for i := 0; i < ctxt.Len(); i++ ***REMOVED***
		t := ctxt.Test(i)
		if len(input) > 0 ***REMOVED***
			t.Input = append(t.Input, input...)
		***REMOVED*** else ***REMOVED***
			t.GenerateInput()
		***REMOVED***
		t.Sort()
		count := 0
		gold := getCollator(*gold, t.Locale)
		for i := 1; i < len(t.Input); i++ ***REMOVED***
			ia := t.Input[i-1]
			ib := t.Input[i]
			if bytes.IndexAny(ib.UTF8, *exclude) != -1 ***REMOVED***
				i++
				continue
			***REMOVED***
			if bytes.IndexAny(ia.UTF8, *exclude) != -1 ***REMOVED***
				continue
			***REMOVED***
			goldCmp := gold.Compare(ia, ib)
			if cmp := bytes.Compare(ia.key, ib.key); cmp != goldCmp ***REMOVED***
				count++
				a := string(ia.UTF8)
				b := string(ib.UTF8)
				fmt.Printf(failedKeyCompare, t.Locale, i-1, a, []rune(a), keyStr(ia.key), b, []rune(b), keyStr(ib.key), cmp, goldCmp, keyStr(gold.Key(ia)), keyStr(gold.Key(ib)))
			***REMOVED*** else if cmp := t.Col.Compare(ia, ib); cmp != goldCmp ***REMOVED***
				count++
				a := string(ia.UTF8)
				b := string(ib.UTF8)
				fmt.Printf(failedCompare, t.Locale, i-1, a, []rune(a), b, []rune(b), cmp, goldCmp)
			***REMOVED***
		***REMOVED***
		if count > 0 ***REMOVED***
			ctxt.Printf("Found %d inconsistencies in %d entries.\n", count, t.Len()-1)
		***REMOVED***
	***REMOVED***
***REMOVED***

const helpTemplate = `
colcmp is a tool for testing and benchmarking collation

Usage: colcmp command [arguments]

The commands are:
***REMOVED******REMOVED***range .***REMOVED******REMOVED***
    ***REMOVED******REMOVED***.Name | printf "%-11s"***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Short***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***

Use "col help [topic]" for more information about that topic.
`

const detailedHelpTemplate = `
Usage: colcmp ***REMOVED******REMOVED***.Usage***REMOVED******REMOVED***

***REMOVED******REMOVED***.Long | trim***REMOVED******REMOVED***
`

func runHelp(args []string) ***REMOVED***
	t := template.New("help")
	t.Funcs(template.FuncMap***REMOVED***"trim": strings.TrimSpace***REMOVED***)
	if len(args) < 1 ***REMOVED***
		template.Must(t.Parse(helpTemplate))
		failOnError(t.Execute(os.Stderr, &commands))
	***REMOVED*** else ***REMOVED***
		for _, cmd := range commands ***REMOVED***
			if cmd.Name() == args[0] ***REMOVED***
				template.Must(t.Parse(detailedHelpTemplate))
				failOnError(t.Execute(os.Stderr, cmd))
				os.Exit(0)
			***REMOVED***
		***REMOVED***
		log.Fatalf("Unknown command %q. Run 'colcmp help'.", args[0])
	***REMOVED***
	os.Exit(0)
***REMOVED***

func main() ***REMOVED***
	flag.Parse()
	log.SetFlags(0)

	ctxt := parseTests()

	if flag.NArg() < 1 ***REMOVED***
		runHelp(nil)
	***REMOVED***
	args := flag.Args()[1:]
	if flag.Arg(0) == "help" ***REMOVED***
		runHelp(args)
	***REMOVED***
	for _, cmd := range commands ***REMOVED***
		if cmd.Name() == flag.Arg(0) ***REMOVED***
			cmd.Run(ctxt, args)
			ctxt.flush()
			return
		***REMOVED***
	***REMOVED***
	runHelp(flag.Args())
***REMOVED***
