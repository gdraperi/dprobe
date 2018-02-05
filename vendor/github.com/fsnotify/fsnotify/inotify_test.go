// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package fsnotify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInotifyCloseRightAway(t *testing.T) ***REMOVED***
	w, err := NewWatcher()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create watcher")
	***REMOVED***

	// Close immediately; it won't even reach the first unix.Read.
	w.Close()

	// Wait for the close to complete.
	<-time.After(50 * time.Millisecond)
	isWatcherReallyClosed(t, w)
***REMOVED***

func TestInotifyCloseSlightlyLater(t *testing.T) ***REMOVED***
	w, err := NewWatcher()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create watcher")
	***REMOVED***

	// Wait until readEvents has reached unix.Read, and Close.
	<-time.After(50 * time.Millisecond)
	w.Close()

	// Wait for the close to complete.
	<-time.After(50 * time.Millisecond)
	isWatcherReallyClosed(t, w)
***REMOVED***

func TestInotifyCloseSlightlyLaterWithWatch(t *testing.T) ***REMOVED***
	testDir := tempMkdir(t)
	defer os.RemoveAll(testDir)

	w, err := NewWatcher()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create watcher")
	***REMOVED***
	w.Add(testDir)

	// Wait until readEvents has reached unix.Read, and Close.
	<-time.After(50 * time.Millisecond)
	w.Close()

	// Wait for the close to complete.
	<-time.After(50 * time.Millisecond)
	isWatcherReallyClosed(t, w)
***REMOVED***

func TestInotifyCloseAfterRead(t *testing.T) ***REMOVED***
	testDir := tempMkdir(t)
	defer os.RemoveAll(testDir)

	w, err := NewWatcher()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create watcher")
	***REMOVED***

	err = w.Add(testDir)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to add .")
	***REMOVED***

	// Generate an event.
	os.Create(filepath.Join(testDir, "somethingSOMETHINGsomethingSOMETHING"))

	// Wait for readEvents to read the event, then close the watcher.
	<-time.After(50 * time.Millisecond)
	w.Close()

	// Wait for the close to complete.
	<-time.After(50 * time.Millisecond)
	isWatcherReallyClosed(t, w)
***REMOVED***

func isWatcherReallyClosed(t *testing.T, w *Watcher) ***REMOVED***
	select ***REMOVED***
	case err, ok := <-w.Errors:
		if ok ***REMOVED***
			t.Fatalf("w.Errors is not closed; readEvents is still alive after closing (error: %v)", err)
		***REMOVED***
	default:
		t.Fatalf("w.Errors would have blocked; readEvents is still alive!")
	***REMOVED***

	select ***REMOVED***
	case _, ok := <-w.Events:
		if ok ***REMOVED***
			t.Fatalf("w.Events is not closed; readEvents is still alive after closing")
		***REMOVED***
	default:
		t.Fatalf("w.Events would have blocked; readEvents is still alive!")
	***REMOVED***
***REMOVED***

func TestInotifyCloseCreate(t *testing.T) ***REMOVED***
	testDir := tempMkdir(t)
	defer os.RemoveAll(testDir)

	w, err := NewWatcher()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create watcher: %v", err)
	***REMOVED***
	defer w.Close()

	err = w.Add(testDir)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to add testDir: %v", err)
	***REMOVED***
	h, err := os.Create(filepath.Join(testDir, "testfile"))
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create file in testdir: %v", err)
	***REMOVED***
	h.Close()
	select ***REMOVED***
	case _ = <-w.Events:
	case err := <-w.Errors:
		t.Fatalf("Error from watcher: %v", err)
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("Took too long to wait for event")
	***REMOVED***

	// At this point, we've received one event, so the goroutine is ready.
	// It's also blocking on unix.Read.
	// Now we try to swap the file descriptor under its nose.
	w.Close()
	w, err = NewWatcher()
	defer w.Close()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create second watcher: %v", err)
	***REMOVED***

	<-time.After(50 * time.Millisecond)
	err = w.Add(testDir)
	if err != nil ***REMOVED***
		t.Fatalf("Error adding testDir again: %v", err)
	***REMOVED***
***REMOVED***

// This test verifies the watcher can keep up with file creations/deletions
// when under load.
func TestInotifyStress(t *testing.T) ***REMOVED***
	maxNumToCreate := 1000

	testDir := tempMkdir(t)
	defer os.RemoveAll(testDir)
	testFilePrefix := filepath.Join(testDir, "testfile")

	w, err := NewWatcher()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create watcher: %v", err)
	***REMOVED***
	defer w.Close()

	err = w.Add(testDir)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to add testDir: %v", err)
	***REMOVED***

	doneChan := make(chan struct***REMOVED******REMOVED***)
	// The buffer ensures that the file generation goroutine is never blocked.
	errChan := make(chan error, 2*maxNumToCreate)

	go func() ***REMOVED***
		for i := 0; i < maxNumToCreate; i++ ***REMOVED***
			testFile := fmt.Sprintf("%s%d", testFilePrefix, i)

			handle, err := os.Create(testFile)
			if err != nil ***REMOVED***
				errChan <- fmt.Errorf("Create failed: %v", err)
				continue
			***REMOVED***

			err = handle.Close()
			if err != nil ***REMOVED***
				errChan <- fmt.Errorf("Close failed: %v", err)
				continue
			***REMOVED***
		***REMOVED***

		// If we delete a newly created file too quickly, inotify will skip the
		// create event and only send the delete event.
		time.Sleep(100 * time.Millisecond)

		for i := 0; i < maxNumToCreate; i++ ***REMOVED***
			testFile := fmt.Sprintf("%s%d", testFilePrefix, i)
			err = os.Remove(testFile)
			if err != nil ***REMOVED***
				errChan <- fmt.Errorf("Remove failed: %v", err)
			***REMOVED***
		***REMOVED***

		close(doneChan)
	***REMOVED***()

	creates := 0
	removes := 0

	finished := false
	after := time.After(10 * time.Second)
	for !finished ***REMOVED***
		select ***REMOVED***
		case <-after:
			t.Fatalf("Not done")
		case <-doneChan:
			finished = true
		case err := <-errChan:
			t.Fatalf("Got an error from file creator goroutine: %v", err)
		case err := <-w.Errors:
			t.Fatalf("Got an error from watcher: %v", err)
		case evt := <-w.Events:
			if !strings.HasPrefix(evt.Name, testFilePrefix) ***REMOVED***
				t.Fatalf("Got an event for an unknown file: %s", evt.Name)
			***REMOVED***
			if evt.Op == Create ***REMOVED***
				creates++
			***REMOVED***
			if evt.Op == Remove ***REMOVED***
				removes++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Drain remaining events from channels
	count := 0
	for count < 10 ***REMOVED***
		select ***REMOVED***
		case err := <-errChan:
			t.Fatalf("Got an error from file creator goroutine: %v", err)
		case err := <-w.Errors:
			t.Fatalf("Got an error from watcher: %v", err)
		case evt := <-w.Events:
			if !strings.HasPrefix(evt.Name, testFilePrefix) ***REMOVED***
				t.Fatalf("Got an event for an unknown file: %s", evt.Name)
			***REMOVED***
			if evt.Op == Create ***REMOVED***
				creates++
			***REMOVED***
			if evt.Op == Remove ***REMOVED***
				removes++
			***REMOVED***
			count = 0
		default:
			count++
			// Give the watcher chances to fill the channels.
			time.Sleep(time.Millisecond)
		***REMOVED***
	***REMOVED***

	if creates-removes > 1 || creates-removes < -1 ***REMOVED***
		t.Fatalf("Creates and removes should not be off by more than one: %d creates, %d removes", creates, removes)
	***REMOVED***
	if creates < 50 ***REMOVED***
		t.Fatalf("Expected at least 50 creates, got %d", creates)
	***REMOVED***
***REMOVED***

func TestInotifyRemoveTwice(t *testing.T) ***REMOVED***
	testDir := tempMkdir(t)
	defer os.RemoveAll(testDir)
	testFile := filepath.Join(testDir, "testfile")

	handle, err := os.Create(testFile)
	if err != nil ***REMOVED***
		t.Fatalf("Create failed: %v", err)
	***REMOVED***
	handle.Close()

	w, err := NewWatcher()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create watcher: %v", err)
	***REMOVED***
	defer w.Close()

	err = w.Add(testFile)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to add testFile: %v", err)
	***REMOVED***

	err = w.Remove(testFile)
	if err != nil ***REMOVED***
		t.Fatalf("wanted successful remove but got: %v", err)
	***REMOVED***

	err = w.Remove(testFile)
	if err == nil ***REMOVED***
		t.Fatalf("no error on removing invalid file")
	***REMOVED***

	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.watches) != 0 ***REMOVED***
		t.Fatalf("Expected watches len is 0, but got: %d, %v", len(w.watches), w.watches)
	***REMOVED***
	if len(w.paths) != 0 ***REMOVED***
		t.Fatalf("Expected paths len is 0, but got: %d, %v", len(w.paths), w.paths)
	***REMOVED***
***REMOVED***

func TestInotifyInnerMapLength(t *testing.T) ***REMOVED***
	testDir := tempMkdir(t)
	defer os.RemoveAll(testDir)
	testFile := filepath.Join(testDir, "testfile")

	handle, err := os.Create(testFile)
	if err != nil ***REMOVED***
		t.Fatalf("Create failed: %v", err)
	***REMOVED***
	handle.Close()

	w, err := NewWatcher()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create watcher: %v", err)
	***REMOVED***
	defer w.Close()

	err = w.Add(testFile)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to add testFile: %v", err)
	***REMOVED***
	go func() ***REMOVED***
		for err := range w.Errors ***REMOVED***
			t.Fatalf("error received: %s", err)
		***REMOVED***
	***REMOVED***()

	err = os.Remove(testFile)
	if err != nil ***REMOVED***
		t.Fatalf("Failed to remove testFile: %v", err)
	***REMOVED***
	_ = <-w.Events                      // consume Remove event
	<-time.After(50 * time.Millisecond) // wait IN_IGNORE propagated

	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.watches) != 0 ***REMOVED***
		t.Fatalf("Expected watches len is 0, but got: %d, %v", len(w.watches), w.watches)
	***REMOVED***
	if len(w.paths) != 0 ***REMOVED***
		t.Fatalf("Expected paths len is 0, but got: %d, %v", len(w.paths), w.paths)
	***REMOVED***
***REMOVED***

func TestInotifyOverflow(t *testing.T) ***REMOVED***
	// We need to generate many more events than the
	// fs.inotify.max_queued_events sysctl setting.
	// We use multiple goroutines (one per directory)
	// to speed up file creation.
	numDirs := 128
	numFiles := 1024

	testDir := tempMkdir(t)
	defer os.RemoveAll(testDir)

	w, err := NewWatcher()
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create watcher: %v", err)
	***REMOVED***
	defer w.Close()

	for dn := 0; dn < numDirs; dn++ ***REMOVED***
		testSubdir := fmt.Sprintf("%s/%d", testDir, dn)

		err := os.Mkdir(testSubdir, 0777)
		if err != nil ***REMOVED***
			t.Fatalf("Cannot create subdir: %v", err)
		***REMOVED***

		err = w.Add(testSubdir)
		if err != nil ***REMOVED***
			t.Fatalf("Failed to add subdir: %v", err)
		***REMOVED***
	***REMOVED***

	errChan := make(chan error, numDirs*numFiles)

	for dn := 0; dn < numDirs; dn++ ***REMOVED***
		testSubdir := fmt.Sprintf("%s/%d", testDir, dn)

		go func() ***REMOVED***
			for fn := 0; fn < numFiles; fn++ ***REMOVED***
				testFile := fmt.Sprintf("%s/%d", testSubdir, fn)

				handle, err := os.Create(testFile)
				if err != nil ***REMOVED***
					errChan <- fmt.Errorf("Create failed: %v", err)
					continue
				***REMOVED***

				err = handle.Close()
				if err != nil ***REMOVED***
					errChan <- fmt.Errorf("Close failed: %v", err)
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	creates := 0
	overflows := 0

	after := time.After(10 * time.Second)
	for overflows == 0 && creates < numDirs*numFiles ***REMOVED***
		select ***REMOVED***
		case <-after:
			t.Fatalf("Not done")
		case err := <-errChan:
			t.Fatalf("Got an error from file creator goroutine: %v", err)
		case err := <-w.Errors:
			if err == ErrEventOverflow ***REMOVED***
				overflows++
			***REMOVED*** else ***REMOVED***
				t.Fatalf("Got an error from watcher: %v", err)
			***REMOVED***
		case evt := <-w.Events:
			if !strings.HasPrefix(evt.Name, testDir) ***REMOVED***
				t.Fatalf("Got an event for an unknown file: %s", evt.Name)
			***REMOVED***
			if evt.Op == Create ***REMOVED***
				creates++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if creates == numDirs*numFiles ***REMOVED***
		t.Fatalf("Could not trigger overflow")
	***REMOVED***

	if overflows == 0 ***REMOVED***
		t.Fatalf("No overflow and not enough creates (expected %d, got %d)",
			numDirs*numFiles, creates)
	***REMOVED***
***REMOVED***
