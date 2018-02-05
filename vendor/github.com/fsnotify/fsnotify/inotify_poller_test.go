// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package fsnotify

import (
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

type testFd [2]int

func makeTestFd(t *testing.T) testFd ***REMOVED***
	var tfd testFd
	errno := unix.Pipe(tfd[:])
	if errno != nil ***REMOVED***
		t.Fatalf("Failed to create pipe: %v", errno)
	***REMOVED***
	return tfd
***REMOVED***

func (tfd testFd) fd() int ***REMOVED***
	return tfd[0]
***REMOVED***

func (tfd testFd) closeWrite(t *testing.T) ***REMOVED***
	errno := unix.Close(tfd[1])
	if errno != nil ***REMOVED***
		t.Fatalf("Failed to close write end of pipe: %v", errno)
	***REMOVED***
***REMOVED***

func (tfd testFd) put(t *testing.T) ***REMOVED***
	buf := make([]byte, 10)
	_, errno := unix.Write(tfd[1], buf)
	if errno != nil ***REMOVED***
		t.Fatalf("Failed to write to pipe: %v", errno)
	***REMOVED***
***REMOVED***

func (tfd testFd) get(t *testing.T) ***REMOVED***
	buf := make([]byte, 10)
	_, errno := unix.Read(tfd[0], buf)
	if errno != nil ***REMOVED***
		t.Fatalf("Failed to read from pipe: %v", errno)
	***REMOVED***
***REMOVED***

func (tfd testFd) close() ***REMOVED***
	unix.Close(tfd[1])
	unix.Close(tfd[0])
***REMOVED***

func makePoller(t *testing.T) (testFd, *fdPoller) ***REMOVED***
	tfd := makeTestFd(t)
	poller, err := newFdPoller(tfd.fd())
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create poller: %v", err)
	***REMOVED***
	return tfd, poller
***REMOVED***

func TestPollerWithBadFd(t *testing.T) ***REMOVED***
	_, err := newFdPoller(-1)
	if err != unix.EBADF ***REMOVED***
		t.Fatalf("Expected EBADF, got: %v", err)
	***REMOVED***
***REMOVED***

func TestPollerWithData(t *testing.T) ***REMOVED***
	tfd, poller := makePoller(t)
	defer tfd.close()
	defer poller.close()

	tfd.put(t)
	ok, err := poller.wait()
	if err != nil ***REMOVED***
		t.Fatalf("poller failed: %v", err)
	***REMOVED***
	if !ok ***REMOVED***
		t.Fatalf("expected poller to return true")
	***REMOVED***
	tfd.get(t)
***REMOVED***

func TestPollerWithWakeup(t *testing.T) ***REMOVED***
	tfd, poller := makePoller(t)
	defer tfd.close()
	defer poller.close()

	err := poller.wake()
	if err != nil ***REMOVED***
		t.Fatalf("wake failed: %v", err)
	***REMOVED***
	ok, err := poller.wait()
	if err != nil ***REMOVED***
		t.Fatalf("poller failed: %v", err)
	***REMOVED***
	if ok ***REMOVED***
		t.Fatalf("expected poller to return false")
	***REMOVED***
***REMOVED***

func TestPollerWithClose(t *testing.T) ***REMOVED***
	tfd, poller := makePoller(t)
	defer tfd.close()
	defer poller.close()

	tfd.closeWrite(t)
	ok, err := poller.wait()
	if err != nil ***REMOVED***
		t.Fatalf("poller failed: %v", err)
	***REMOVED***
	if !ok ***REMOVED***
		t.Fatalf("expected poller to return true")
	***REMOVED***
***REMOVED***

func TestPollerWithWakeupAndData(t *testing.T) ***REMOVED***
	tfd, poller := makePoller(t)
	defer tfd.close()
	defer poller.close()

	tfd.put(t)
	err := poller.wake()
	if err != nil ***REMOVED***
		t.Fatalf("wake failed: %v", err)
	***REMOVED***

	// both data and wakeup
	ok, err := poller.wait()
	if err != nil ***REMOVED***
		t.Fatalf("poller failed: %v", err)
	***REMOVED***
	if !ok ***REMOVED***
		t.Fatalf("expected poller to return true")
	***REMOVED***

	// data is still in the buffer, wakeup is cleared
	ok, err = poller.wait()
	if err != nil ***REMOVED***
		t.Fatalf("poller failed: %v", err)
	***REMOVED***
	if !ok ***REMOVED***
		t.Fatalf("expected poller to return true")
	***REMOVED***

	tfd.get(t)
	// data is gone, only wakeup now
	err = poller.wake()
	if err != nil ***REMOVED***
		t.Fatalf("wake failed: %v", err)
	***REMOVED***
	ok, err = poller.wait()
	if err != nil ***REMOVED***
		t.Fatalf("poller failed: %v", err)
	***REMOVED***
	if ok ***REMOVED***
		t.Fatalf("expected poller to return false")
	***REMOVED***
***REMOVED***

func TestPollerConcurrent(t *testing.T) ***REMOVED***
	tfd, poller := makePoller(t)
	defer tfd.close()
	defer poller.close()

	oks := make(chan bool)
	live := make(chan bool)
	defer close(live)
	go func() ***REMOVED***
		defer close(oks)
		for ***REMOVED***
			ok, err := poller.wait()
			if err != nil ***REMOVED***
				t.Fatalf("poller failed: %v", err)
			***REMOVED***
			oks <- ok
			if !<-live ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Try a write
	select ***REMOVED***
	case <-time.After(50 * time.Millisecond):
	case <-oks:
		t.Fatalf("poller did not wait")
	***REMOVED***
	tfd.put(t)
	if !<-oks ***REMOVED***
		t.Fatalf("expected true")
	***REMOVED***
	tfd.get(t)
	live <- true

	// Try a wakeup
	select ***REMOVED***
	case <-time.After(50 * time.Millisecond):
	case <-oks:
		t.Fatalf("poller did not wait")
	***REMOVED***
	err := poller.wake()
	if err != nil ***REMOVED***
		t.Fatalf("wake failed: %v", err)
	***REMOVED***
	if <-oks ***REMOVED***
		t.Fatalf("expected false")
	***REMOVED***
	live <- true

	// Try a close
	select ***REMOVED***
	case <-time.After(50 * time.Millisecond):
	case <-oks:
		t.Fatalf("poller did not wait")
	***REMOVED***
	tfd.closeWrite(t)
	if !<-oks ***REMOVED***
		t.Fatalf("expected true")
	***REMOVED***
	tfd.get(t)
***REMOVED***
