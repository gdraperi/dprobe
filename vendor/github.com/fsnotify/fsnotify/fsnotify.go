// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

// Package fsnotify provides a platform-independent interface for file system notifications.
package fsnotify

import (
	"bytes"
	"errors"
	"fmt"
)

// Event represents a single file system notification.
type Event struct ***REMOVED***
	Name string // Relative path to the file or directory.
	Op   Op     // File operation that triggered the event.
***REMOVED***

// Op describes a set of file operations.
type Op uint32

// These are the generalized file operations that can trigger a notification.
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

func (op Op) String() string ***REMOVED***
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if op&Create == Create ***REMOVED***
		buffer.WriteString("|CREATE")
	***REMOVED***
	if op&Remove == Remove ***REMOVED***
		buffer.WriteString("|REMOVE")
	***REMOVED***
	if op&Write == Write ***REMOVED***
		buffer.WriteString("|WRITE")
	***REMOVED***
	if op&Rename == Rename ***REMOVED***
		buffer.WriteString("|RENAME")
	***REMOVED***
	if op&Chmod == Chmod ***REMOVED***
		buffer.WriteString("|CHMOD")
	***REMOVED***
	if buffer.Len() == 0 ***REMOVED***
		return ""
	***REMOVED***
	return buffer.String()[1:] // Strip leading pipe
***REMOVED***

// String returns a string representation of the event in the form
// "file: REMOVE|WRITE|..."
func (e Event) String() string ***REMOVED***
	return fmt.Sprintf("%q: %s", e.Name, e.Op.String())
***REMOVED***

// Common errors that can be reported by a watcher
var ErrEventOverflow = errors.New("fsnotify queue overflow")
