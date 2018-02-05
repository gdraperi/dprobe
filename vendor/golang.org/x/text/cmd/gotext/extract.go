// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"golang.org/x/text/message/pipeline"
)

// TODO:
// - merge information into existing files
// - handle different file formats (PO, XLIFF)
// - handle features (gender, plural)
// - message rewriting

func init() ***REMOVED***
	lang = cmdExtract.Flag.String("lang", "en-US", "comma-separated list of languages to process")
***REMOVED***

var cmdExtract = &Command***REMOVED***
	Run:       runExtract,
	UsageLine: "extract <package>*",
	Short:     "extracts strings to be translated from code",
***REMOVED***

func runExtract(cmd *Command, config *pipeline.Config, args []string) error ***REMOVED***
	config.Packages = args
	state, err := pipeline.Extract(config)
	if err != nil ***REMOVED***
		return wrap(err, "extract failed")
	***REMOVED***
	if err := state.Import(); err != nil ***REMOVED***
		return wrap(err, "import failed")
	***REMOVED***
	if err := state.Merge(); err != nil ***REMOVED***
		return wrap(err, "merge failed")
	***REMOVED***
	return wrap(state.Export(), "export failed")
***REMOVED***
