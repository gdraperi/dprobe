// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"golang.org/x/text/message/pipeline"
)

func init() ***REMOVED***
	out = cmdGenerate.Flag.String("out", "", "output file to write to")
***REMOVED***

var cmdGenerate = &Command***REMOVED***
	Run:       runGenerate,
	UsageLine: "generate <package>",
	Short:     "generates code to insert translated messages",
***REMOVED***

func runGenerate(cmd *Command, config *pipeline.Config, args []string) error ***REMOVED***
	config.Packages = args
	s, err := pipeline.Extract(config)
	if err != nil ***REMOVED***
		return wrap(err, "extraction failed")
	***REMOVED***
	if err := s.Import(); err != nil ***REMOVED***
		return wrap(err, "import failed")
	***REMOVED***
	return wrap(s.Generate(), "generation failed")
***REMOVED***
