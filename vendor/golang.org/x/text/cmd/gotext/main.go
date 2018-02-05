// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go build -o gotext.latest
//go:generate ./gotext.latest help gendocumentation
//go:generate rm gotext.latest

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/message/pipeline"

	"golang.org/x/text/language"
	"golang.org/x/tools/go/buildutil"
)

func init() ***REMOVED***
	flag.Var((*buildutil.TagsFlag)(&build.Default.BuildTags), "tags", buildutil.TagsFlagDoc)
***REMOVED***

var (
	srcLang = flag.String("srclang", "en-US", "the source-code language")
	dir     = flag.String("dir", "locales", "default subdirectory to store translation files")
)

func config() (*pipeline.Config, error) ***REMOVED***
	tag, err := language.Parse(*srcLang)
	if err != nil ***REMOVED***
		return nil, wrap(err, "invalid srclang")
	***REMOVED***
	return &pipeline.Config***REMOVED***
		SourceLanguage:      tag,
		Supported:           getLangs(),
		TranslationsPattern: `messages\.(.*)\.json`,
		GenFile:             *out,
	***REMOVED***, nil
***REMOVED***

// NOTE: the Command struct is copied from the go tool in core.

// A Command is an implementation of a go command
// like go build or go fix.
type Command struct ***REMOVED***
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, c *pipeline.Config, args []string) error

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
***REMOVED***

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string ***REMOVED***
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 ***REMOVED***
		name = name[:i]
	***REMOVED***
	return name
***REMOVED***

func (c *Command) Usage() ***REMOVED***
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
***REMOVED***

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool ***REMOVED***
	return c.Run != nil
***REMOVED***

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'go help'.
var commands = []*Command***REMOVED***
	cmdUpdate,
	cmdExtract,
	cmdRewrite,
	cmdGenerate,
	// TODO:
	// - update: full-cycle update of extraction, sending, and integration
	// - report: report of freshness of translations
***REMOVED***

var exitStatus = 0
var exitMu sync.Mutex

func setExitStatus(n int) ***REMOVED***
	exitMu.Lock()
	if exitStatus < n ***REMOVED***
		exitStatus = n
	***REMOVED***
	exitMu.Unlock()
***REMOVED***

var origEnv []string

func main() ***REMOVED***
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 ***REMOVED***
		usage()
	***REMOVED***

	if args[0] == "help" ***REMOVED***
		help(args[1:])
		return
	***REMOVED***

	for _, cmd := range commands ***REMOVED***
		if cmd.Name() == args[0] && cmd.Runnable() ***REMOVED***
			cmd.Flag.Usage = func() ***REMOVED*** cmd.Usage() ***REMOVED***
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			config, err := config()
			if err != nil ***REMOVED***
				fatalf("gotext: %+v", err)
			***REMOVED***
			if err := cmd.Run(cmd, config, args); err != nil ***REMOVED***
				fatalf("gotext: %+v", err)
			***REMOVED***
			exit()
			return
		***REMOVED***
	***REMOVED***

	fmt.Fprintf(os.Stderr, "gotext: unknown subcommand %q\nRun 'go help' for usage.\n", args[0])
	setExitStatus(2)
	exit()
***REMOVED***

var usageTemplate = `gotext is a tool for managing text in Go source code.

Usage:

	gotext command [arguments]

The commands are:
***REMOVED******REMOVED***range .***REMOVED******REMOVED******REMOVED******REMOVED***if .Runnable***REMOVED******REMOVED***
	***REMOVED******REMOVED***.Name | printf "%-11s"***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Short***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***

Use "go help [command]" for more information about a command.

Additional help topics:
***REMOVED******REMOVED***range .***REMOVED******REMOVED******REMOVED******REMOVED***if not .Runnable***REMOVED******REMOVED***
	***REMOVED******REMOVED***.Name | printf "%-11s"***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Short***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***

Use "gotext help [topic]" for more information about that topic.

`

var helpTemplate = `***REMOVED******REMOVED***if .Runnable***REMOVED******REMOVED***usage: go ***REMOVED******REMOVED***.UsageLine***REMOVED******REMOVED***

***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***.Long | trim***REMOVED******REMOVED***
`

var documentationTemplate = `***REMOVED******REMOVED***range .***REMOVED******REMOVED******REMOVED******REMOVED***if .Short***REMOVED******REMOVED******REMOVED******REMOVED***.Short | capitalize***REMOVED******REMOVED***

***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .Runnable***REMOVED******REMOVED***Usage:

	go ***REMOVED******REMOVED***.UsageLine***REMOVED******REMOVED***

***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***.Long | trim***REMOVED******REMOVED***


***REMOVED******REMOVED***end***REMOVED******REMOVED***`

// commentWriter writes a Go comment to the underlying io.Writer,
// using line comment form (//).
type commentWriter struct ***REMOVED***
	W            io.Writer
	wroteSlashes bool // Wrote "//" at the beginning of the current line.
***REMOVED***

func (c *commentWriter) Write(p []byte) (int, error) ***REMOVED***
	var n int
	for i, b := range p ***REMOVED***
		if !c.wroteSlashes ***REMOVED***
			s := "//"
			if b != '\n' ***REMOVED***
				s = "// "
			***REMOVED***
			if _, err := io.WriteString(c.W, s); err != nil ***REMOVED***
				return n, err
			***REMOVED***
			c.wroteSlashes = true
		***REMOVED***
		n0, err := c.W.Write(p[i : i+1])
		n += n0
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		if b == '\n' ***REMOVED***
			c.wroteSlashes = false
		***REMOVED***
	***REMOVED***
	return len(p), nil
***REMOVED***

// An errWriter wraps a writer, recording whether a write error occurred.
type errWriter struct ***REMOVED***
	w   io.Writer
	err error
***REMOVED***

func (w *errWriter) Write(b []byte) (int, error) ***REMOVED***
	n, err := w.w.Write(b)
	if err != nil ***REMOVED***
		w.err = err
	***REMOVED***
	return n, err
***REMOVED***

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface***REMOVED******REMOVED***) ***REMOVED***
	t := template.New("top")
	t.Funcs(template.FuncMap***REMOVED***"trim": strings.TrimSpace, "capitalize": capitalize***REMOVED***)
	template.Must(t.Parse(text))
	ew := &errWriter***REMOVED***w: w***REMOVED***
	err := t.Execute(ew, data)
	if ew.err != nil ***REMOVED***
		// I/O error writing. Ignore write on closed pipe.
		if strings.Contains(ew.err.Error(), "pipe") ***REMOVED***
			os.Exit(1)
		***REMOVED***
		fatalf("writing output: %v", ew.err)
	***REMOVED***
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func capitalize(s string) string ***REMOVED***
	if s == "" ***REMOVED***
		return s
	***REMOVED***
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
***REMOVED***

func printUsage(w io.Writer) ***REMOVED***
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, commands)
	bw.Flush()
***REMOVED***

func usage() ***REMOVED***
	printUsage(os.Stderr)
	os.Exit(2)
***REMOVED***

// help implements the 'help' command.
func help(args []string) ***REMOVED***
	if len(args) == 0 ***REMOVED***
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'go help'.
		return
	***REMOVED***
	if len(args) != 1 ***REMOVED***
		fmt.Fprintf(os.Stderr, "usage: go help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'go help'
	***REMOVED***

	arg := args[0]

	// 'go help documentation' generates doc.go.
	if strings.HasSuffix(arg, "documentation") ***REMOVED***
		w := &bytes.Buffer***REMOVED******REMOVED***

		fmt.Fprintln(w, "// Code generated by go generate. DO NOT EDIT.")
		fmt.Fprintln(w)
		buf := new(bytes.Buffer)
		printUsage(buf)
		usage := &Command***REMOVED***Long: buf.String()***REMOVED***
		tmpl(&commentWriter***REMOVED***W: w***REMOVED***, documentationTemplate, append([]*Command***REMOVED***usage***REMOVED***, commands...))
		fmt.Fprintln(w, "package main")
		if arg == "gendocumentation" ***REMOVED***
			b, err := format.Source(w.Bytes())
			if err != nil ***REMOVED***
				logf("Could not format generated docs: %v\n", err)
			***REMOVED***
			if err := ioutil.WriteFile("doc.go", b, 0666); err != nil ***REMOVED***
				logf("Could not create file alldocs.go: %v\n", err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			fmt.Println(w.String())
		***REMOVED***
		return
	***REMOVED***

	for _, cmd := range commands ***REMOVED***
		if cmd.Name() == arg ***REMOVED***
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'go help cmd'.
			return
		***REMOVED***
	***REMOVED***

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'go help'.\n", arg)
	os.Exit(2) // failed at 'go help cmd'
***REMOVED***

func getLangs() (tags []language.Tag) ***REMOVED***
	for _, t := range strings.Split(*lang, ",") ***REMOVED***
		if t == "" ***REMOVED***
			continue
		***REMOVED***
		tag, err := language.Parse(t)
		if err != nil ***REMOVED***
			fatalf("gotext: could not parse language %q: %v", t, err)
		***REMOVED***
		tags = append(tags, tag)
	***REMOVED***
	return tags
***REMOVED***

var atexitFuncs []func()

func atexit(f func()) ***REMOVED***
	atexitFuncs = append(atexitFuncs, f)
***REMOVED***

func exit() ***REMOVED***
	for _, f := range atexitFuncs ***REMOVED***
		f()
	***REMOVED***
	os.Exit(exitStatus)
***REMOVED***

func fatalf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	logf(format, args...)
	exit()
***REMOVED***

func logf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	log.Printf(format, args...)
	setExitStatus(1)
***REMOVED***

func exitIfErrors() ***REMOVED***
	if exitStatus != 0 ***REMOVED***
		exit()
	***REMOVED***
***REMOVED***
