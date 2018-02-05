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

var (
	lang *string
	out  *string
)

func init() ***REMOVED***
	lang = cmdUpdate.Flag.String("lang", "en-US", "comma-separated list of languages to process")
	out = cmdUpdate.Flag.String("out", "", "output file to write to")
***REMOVED***

var cmdUpdate = &Command***REMOVED***
	Run:       runUpdate,
	UsageLine: "update <package>* [-out <gofile>]",
	Short:     "merge translations and generate catalog",
***REMOVED***

func runUpdate(cmd *Command, config *pipeline.Config, args []string) error ***REMOVED***
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
	if err := state.Export(); err != nil ***REMOVED***
		return wrap(err, "export failed")
	***REMOVED***
	if *out != "" ***REMOVED***
		return wrap(state.Generate(), "generation failed")
	***REMOVED***
	return nil
***REMOVED***
