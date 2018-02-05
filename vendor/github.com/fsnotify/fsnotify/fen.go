// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package fsnotify

import (
	"errors"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct ***REMOVED***
	Events chan Event
	Errors chan error
***REMOVED***

// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
func NewWatcher() (*Watcher, error) ***REMOVED***
	return nil, errors.New("FEN based watcher not yet supported for fsnotify\n")
***REMOVED***

// Close removes all watches and closes the events channel.
func (w *Watcher) Close() error ***REMOVED***
	return nil
***REMOVED***

// Add starts watching the named file or directory (non-recursively).
func (w *Watcher) Add(name string) error ***REMOVED***
	return nil
***REMOVED***

// Remove stops watching the the named file or directory (non-recursively).
func (w *Watcher) Remove(name string) error ***REMOVED***
	return nil
***REMOVED***
