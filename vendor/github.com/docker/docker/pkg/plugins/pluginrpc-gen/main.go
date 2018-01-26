package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"unicode"
	"unicode/utf8"
)

type stringSet struct ***REMOVED***
	values map[string]struct***REMOVED******REMOVED***
***REMOVED***

func (s stringSet) String() string ***REMOVED***
	return ""
***REMOVED***

func (s stringSet) Set(value string) error ***REMOVED***
	s.values[value] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return nil
***REMOVED***
func (s stringSet) GetValues() map[string]struct***REMOVED******REMOVED*** ***REMOVED***
	return s.values
***REMOVED***

var (
	typeName   = flag.String("type", "", "interface type to generate plugin rpc proxy for")
	rpcName    = flag.String("name", *typeName, "RPC name, set if different from type")
	inputFile  = flag.String("i", "", "input file path")
	outputFile = flag.String("o", *inputFile+"_proxy.go", "output file path")

	skipFuncs   map[string]struct***REMOVED******REMOVED***
	flSkipFuncs = stringSet***REMOVED***make(map[string]struct***REMOVED******REMOVED***)***REMOVED***

	flBuildTags = stringSet***REMOVED***make(map[string]struct***REMOVED******REMOVED***)***REMOVED***
)

func errorOut(msg string, err error) ***REMOVED***
	if err == nil ***REMOVED***
		return
	***REMOVED***
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(1)
***REMOVED***

func checkFlags() error ***REMOVED***
	if *outputFile == "" ***REMOVED***
		return fmt.Errorf("missing required flag `-o`")
	***REMOVED***
	if *inputFile == "" ***REMOVED***
		return fmt.Errorf("missing required flag `-i`")
	***REMOVED***
	return nil
***REMOVED***

func main() ***REMOVED***
	flag.Var(flSkipFuncs, "skip", "skip parsing for function")
	flag.Var(flBuildTags, "tag", "build tags to add to generated files")
	flag.Parse()
	skipFuncs = flSkipFuncs.GetValues()

	errorOut("error", checkFlags())

	pkg, err := Parse(*inputFile, *typeName)
	errorOut(fmt.Sprintf("error parsing requested type %s", *typeName), err)

	var analysis = struct ***REMOVED***
		InterfaceType string
		RPCName       string
		BuildTags     map[string]struct***REMOVED******REMOVED***
		*ParsedPkg
	***REMOVED******REMOVED***toLower(*typeName), *rpcName, flBuildTags.GetValues(), pkg***REMOVED***
	var buf bytes.Buffer

	errorOut("parser error", generatedTempl.Execute(&buf, analysis))
	src, err := format.Source(buf.Bytes())
	errorOut("error formatting generated source:\n"+buf.String(), err)
	errorOut("error writing file", ioutil.WriteFile(*outputFile, src, 0644))
***REMOVED***

func toLower(s string) string ***REMOVED***
	if s == "" ***REMOVED***
		return ""
	***REMOVED***
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
***REMOVED***
