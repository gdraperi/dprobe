// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package fsnotify

import (
	"os"
	"testing"
	"time"
)

func TestEventStringWithValue(t *testing.T) ***REMOVED***
	for opMask, expectedString := range map[Op]string***REMOVED***
		Chmod | Create: `"/usr/someFile": CREATE|CHMOD`,
		Rename:         `"/usr/someFile": RENAME`,
		Remove:         `"/usr/someFile": REMOVE`,
		Write | Chmod:  `"/usr/someFile": WRITE|CHMOD`,
	***REMOVED*** ***REMOVED***
		event := Event***REMOVED***Name: "/usr/someFile", Op: opMask***REMOVED***
		if event.String() != expectedString ***REMOVED***
			t.Fatalf("Expected %s, got: %v", expectedString, event.String())
		***REMOVED***

	***REMOVED***
***REMOVED***

func TestEventOpStringWithValue(t *testing.T) ***REMOVED***
	expectedOpString := "WRITE|CHMOD"
	event := Event***REMOVED***Name: "someFile", Op: Write | Chmod***REMOVED***
	if event.Op.String() != expectedOpString ***REMOVED***
		t.Fatalf("Expected %s, got: %v", expectedOpString, event.Op.String())
	***REMOVED***
***REMOVED***

func TestEventOpStringWithNoValue(t *testing.T) ***REMOVED***
	expectedOpString := ""
	event := Event***REMOVED***Name: "testFile", Op: 0***REMOVED***
	if event.Op.String() != expectedOpString ***REMOVED***
		t.Fatalf("Expected %s, got: %v", expectedOpString, event.Op.String())
	***REMOVED***
***REMOVED***

// TestWatcherClose tests that the goroutine started by creating the watcher can be
// signalled to return at any time, even if there is no goroutine listening on the events
// or errors channels.
func TestWatcherClose(t *testing.T) ***REMOVED***
	t.Parallel()

	name := tempMkFile(t, "")
	w := newWatcher(t)
	err := w.Add(name)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = os.Remove(name)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// Allow the watcher to receive the event.
	time.Sleep(time.Millisecond * 100)

	err = w.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
