// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"golang.org/x/text/message/pipeline"
)

const printerType = "golang.org/x/text/message.Printer"

// TODO:
// - merge information into existing files
// - handle different file formats (PO, XLIFF)
// - handle features (gender, plural)
// - message rewriting

func init() ***REMOVED***
	overwrite = cmdRewrite.Flag.Bool("w", false, "write files in place")
***REMOVED***

var (
	overwrite *bool
)

var cmdRewrite = &Command***REMOVED***
	Run:       runRewrite,
	UsageLine: "rewrite <package>",
	Short:     "rewrites fmt functions to use a message Printer",
	Long: `
rewrite is typically done once for a project. It rewrites all usages of
fmt to use x/text's message package whenever a message.Printer is in scope.
It rewrites Print and Println calls with constant strings to the equivalent
using Printf to allow translators to reorder arguments.
`,
***REMOVED***

func runRewrite(cmd *Command, _ *pipeline.Config, args []string) error ***REMOVED***
	w := os.Stdout
	if *overwrite ***REMOVED***
		w = nil
	***REMOVED***
	pkg := "."
	switch len(args) ***REMOVED***
	case 0:
	case 1:
		pkg = args[0]
	default:
		return errorf("can only specify at most one package")
	***REMOVED***
	return pipeline.Rewrite(w, pkg)
***REMOVED***
