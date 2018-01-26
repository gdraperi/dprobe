package xfer

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/docker/docker/pkg/progress"
)

func TestTransfer(t *testing.T) ***REMOVED***
	makeXferFunc := func(id string) DoFunc ***REMOVED***
		return func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer ***REMOVED***
			select ***REMOVED***
			case <-start:
			default:
				t.Fatalf("transfer function not started even though concurrency limit not reached")
			***REMOVED***

			xfer := NewTransfer()
			go func() ***REMOVED***
				for i := 0; i <= 10; i++ ***REMOVED***
					progressChan <- progress.Progress***REMOVED***ID: id, Action: "testing", Current: int64(i), Total: 10***REMOVED***
					time.Sleep(10 * time.Millisecond)
				***REMOVED***
				close(progressChan)
			***REMOVED***()
			return xfer
		***REMOVED***
	***REMOVED***

	tm := NewTransferManager(5)
	progressChan := make(chan progress.Progress)
	progressDone := make(chan struct***REMOVED******REMOVED***)
	receivedProgress := make(map[string]int64)

	go func() ***REMOVED***
		for p := range progressChan ***REMOVED***
			val, present := receivedProgress[p.ID]
			if present && p.Current <= val ***REMOVED***
				t.Fatalf("got unexpected progress value: %d (expected %d)", p.Current, val+1)
			***REMOVED***
			receivedProgress[p.ID] = p.Current
		***REMOVED***
		close(progressDone)
	***REMOVED***()

	// Start a few transfers
	ids := []string***REMOVED***"id1", "id2", "id3"***REMOVED***
	xfers := make([]Transfer, len(ids))
	watchers := make([]*Watcher, len(ids))
	for i, id := range ids ***REMOVED***
		xfers[i], watchers[i] = tm.Transfer(id, makeXferFunc(id), progress.ChanOutput(progressChan))
	***REMOVED***

	for i, xfer := range xfers ***REMOVED***
		<-xfer.Done()
		xfer.Release(watchers[i])
	***REMOVED***
	close(progressChan)
	<-progressDone

	for _, id := range ids ***REMOVED***
		if receivedProgress[id] != 10 ***REMOVED***
			t.Fatalf("final progress value %d instead of 10", receivedProgress[id])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestConcurrencyLimit(t *testing.T) ***REMOVED***
	concurrencyLimit := 3
	var runningJobs int32

	makeXferFunc := func(id string) DoFunc ***REMOVED***
		return func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer ***REMOVED***
			xfer := NewTransfer()
			go func() ***REMOVED***
				<-start
				totalJobs := atomic.AddInt32(&runningJobs, 1)
				if int(totalJobs) > concurrencyLimit ***REMOVED***
					t.Fatalf("too many jobs running")
				***REMOVED***
				for i := 0; i <= 10; i++ ***REMOVED***
					progressChan <- progress.Progress***REMOVED***ID: id, Action: "testing", Current: int64(i), Total: 10***REMOVED***
					time.Sleep(10 * time.Millisecond)
				***REMOVED***
				atomic.AddInt32(&runningJobs, -1)
				close(progressChan)
			***REMOVED***()
			return xfer
		***REMOVED***
	***REMOVED***

	tm := NewTransferManager(concurrencyLimit)
	progressChan := make(chan progress.Progress)
	progressDone := make(chan struct***REMOVED******REMOVED***)
	receivedProgress := make(map[string]int64)

	go func() ***REMOVED***
		for p := range progressChan ***REMOVED***
			receivedProgress[p.ID] = p.Current
		***REMOVED***
		close(progressDone)
	***REMOVED***()

	// Start more transfers than the concurrency limit
	ids := []string***REMOVED***"id1", "id2", "id3", "id4", "id5", "id6", "id7", "id8"***REMOVED***
	xfers := make([]Transfer, len(ids))
	watchers := make([]*Watcher, len(ids))
	for i, id := range ids ***REMOVED***
		xfers[i], watchers[i] = tm.Transfer(id, makeXferFunc(id), progress.ChanOutput(progressChan))
	***REMOVED***

	for i, xfer := range xfers ***REMOVED***
		<-xfer.Done()
		xfer.Release(watchers[i])
	***REMOVED***
	close(progressChan)
	<-progressDone

	for _, id := range ids ***REMOVED***
		if receivedProgress[id] != 10 ***REMOVED***
			t.Fatalf("final progress value %d instead of 10", receivedProgress[id])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestInactiveJobs(t *testing.T) ***REMOVED***
	concurrencyLimit := 3
	var runningJobs int32
	testDone := make(chan struct***REMOVED******REMOVED***)

	makeXferFunc := func(id string) DoFunc ***REMOVED***
		return func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer ***REMOVED***
			xfer := NewTransfer()
			go func() ***REMOVED***
				<-start
				totalJobs := atomic.AddInt32(&runningJobs, 1)
				if int(totalJobs) > concurrencyLimit ***REMOVED***
					t.Fatalf("too many jobs running")
				***REMOVED***
				for i := 0; i <= 10; i++ ***REMOVED***
					progressChan <- progress.Progress***REMOVED***ID: id, Action: "testing", Current: int64(i), Total: 10***REMOVED***
					time.Sleep(10 * time.Millisecond)
				***REMOVED***
				atomic.AddInt32(&runningJobs, -1)
				close(inactive)
				<-testDone
				close(progressChan)
			***REMOVED***()
			return xfer
		***REMOVED***
	***REMOVED***

	tm := NewTransferManager(concurrencyLimit)
	progressChan := make(chan progress.Progress)
	progressDone := make(chan struct***REMOVED******REMOVED***)
	receivedProgress := make(map[string]int64)

	go func() ***REMOVED***
		for p := range progressChan ***REMOVED***
			receivedProgress[p.ID] = p.Current
		***REMOVED***
		close(progressDone)
	***REMOVED***()

	// Start more transfers than the concurrency limit
	ids := []string***REMOVED***"id1", "id2", "id3", "id4", "id5", "id6", "id7", "id8"***REMOVED***
	xfers := make([]Transfer, len(ids))
	watchers := make([]*Watcher, len(ids))
	for i, id := range ids ***REMOVED***
		xfers[i], watchers[i] = tm.Transfer(id, makeXferFunc(id), progress.ChanOutput(progressChan))
	***REMOVED***

	close(testDone)
	for i, xfer := range xfers ***REMOVED***
		<-xfer.Done()
		xfer.Release(watchers[i])
	***REMOVED***
	close(progressChan)
	<-progressDone

	for _, id := range ids ***REMOVED***
		if receivedProgress[id] != 10 ***REMOVED***
			t.Fatalf("final progress value %d instead of 10", receivedProgress[id])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWatchRelease(t *testing.T) ***REMOVED***
	ready := make(chan struct***REMOVED******REMOVED***)

	makeXferFunc := func(id string) DoFunc ***REMOVED***
		return func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer ***REMOVED***
			xfer := NewTransfer()
			go func() ***REMOVED***
				defer func() ***REMOVED***
					close(progressChan)
				***REMOVED***()
				<-ready
				for i := int64(0); ; i++ ***REMOVED***
					select ***REMOVED***
					case <-time.After(10 * time.Millisecond):
					case <-xfer.Context().Done():
						return
					***REMOVED***
					progressChan <- progress.Progress***REMOVED***ID: id, Action: "testing", Current: i, Total: 10***REMOVED***
				***REMOVED***
			***REMOVED***()
			return xfer
		***REMOVED***
	***REMOVED***

	tm := NewTransferManager(5)

	type watcherInfo struct ***REMOVED***
		watcher               *Watcher
		progressChan          chan progress.Progress
		progressDone          chan struct***REMOVED******REMOVED***
		receivedFirstProgress chan struct***REMOVED******REMOVED***
	***REMOVED***

	progressConsumer := func(w watcherInfo) ***REMOVED***
		first := true
		for range w.progressChan ***REMOVED***
			if first ***REMOVED***
				close(w.receivedFirstProgress)
			***REMOVED***
			first = false
		***REMOVED***
		close(w.progressDone)
	***REMOVED***

	// Start a transfer
	watchers := make([]watcherInfo, 5)
	var xfer Transfer
	watchers[0].progressChan = make(chan progress.Progress)
	watchers[0].progressDone = make(chan struct***REMOVED******REMOVED***)
	watchers[0].receivedFirstProgress = make(chan struct***REMOVED******REMOVED***)
	xfer, watchers[0].watcher = tm.Transfer("id1", makeXferFunc("id1"), progress.ChanOutput(watchers[0].progressChan))
	go progressConsumer(watchers[0])

	// Give it multiple watchers
	for i := 1; i != len(watchers); i++ ***REMOVED***
		watchers[i].progressChan = make(chan progress.Progress)
		watchers[i].progressDone = make(chan struct***REMOVED******REMOVED***)
		watchers[i].receivedFirstProgress = make(chan struct***REMOVED******REMOVED***)
		watchers[i].watcher = xfer.Watch(progress.ChanOutput(watchers[i].progressChan))
		go progressConsumer(watchers[i])
	***REMOVED***

	// Now that the watchers are set up, allow the transfer goroutine to
	// proceed.
	close(ready)

	// Confirm that each watcher gets progress output.
	for _, w := range watchers ***REMOVED***
		<-w.receivedFirstProgress
	***REMOVED***

	// Release one watcher every 5ms
	for _, w := range watchers ***REMOVED***
		xfer.Release(w.watcher)
		<-time.After(5 * time.Millisecond)
	***REMOVED***

	// Now that all watchers have been released, Released() should
	// return a closed channel.
	<-xfer.Released()

	// Done() should return a closed channel because the xfer func returned
	// due to cancellation.
	<-xfer.Done()

	for _, w := range watchers ***REMOVED***
		close(w.progressChan)
		<-w.progressDone
	***REMOVED***
***REMOVED***

func TestWatchFinishedTransfer(t *testing.T) ***REMOVED***
	makeXferFunc := func(id string) DoFunc ***REMOVED***
		return func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer ***REMOVED***
			xfer := NewTransfer()
			go func() ***REMOVED***
				// Finish immediately
				close(progressChan)
			***REMOVED***()
			return xfer
		***REMOVED***
	***REMOVED***

	tm := NewTransferManager(5)

	// Start a transfer
	watchers := make([]*Watcher, 3)
	var xfer Transfer
	xfer, watchers[0] = tm.Transfer("id1", makeXferFunc("id1"), progress.ChanOutput(make(chan progress.Progress)))

	// Give it a watcher immediately
	watchers[1] = xfer.Watch(progress.ChanOutput(make(chan progress.Progress)))

	// Wait for the transfer to complete
	<-xfer.Done()

	// Set up another watcher
	watchers[2] = xfer.Watch(progress.ChanOutput(make(chan progress.Progress)))

	// Release the watchers
	for _, w := range watchers ***REMOVED***
		xfer.Release(w)
	***REMOVED***

	// Now that all watchers have been released, Released() should
	// return a closed channel.
	<-xfer.Released()
***REMOVED***

func TestDuplicateTransfer(t *testing.T) ***REMOVED***
	ready := make(chan struct***REMOVED******REMOVED***)

	var xferFuncCalls int32

	makeXferFunc := func(id string) DoFunc ***REMOVED***
		return func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer ***REMOVED***
			atomic.AddInt32(&xferFuncCalls, 1)
			xfer := NewTransfer()
			go func() ***REMOVED***
				defer func() ***REMOVED***
					close(progressChan)
				***REMOVED***()
				<-ready
				for i := int64(0); ; i++ ***REMOVED***
					select ***REMOVED***
					case <-time.After(10 * time.Millisecond):
					case <-xfer.Context().Done():
						return
					***REMOVED***
					progressChan <- progress.Progress***REMOVED***ID: id, Action: "testing", Current: i, Total: 10***REMOVED***
				***REMOVED***
			***REMOVED***()
			return xfer
		***REMOVED***
	***REMOVED***

	tm := NewTransferManager(5)

	type transferInfo struct ***REMOVED***
		xfer                  Transfer
		watcher               *Watcher
		progressChan          chan progress.Progress
		progressDone          chan struct***REMOVED******REMOVED***
		receivedFirstProgress chan struct***REMOVED******REMOVED***
	***REMOVED***

	progressConsumer := func(t transferInfo) ***REMOVED***
		first := true
		for range t.progressChan ***REMOVED***
			if first ***REMOVED***
				close(t.receivedFirstProgress)
			***REMOVED***
			first = false
		***REMOVED***
		close(t.progressDone)
	***REMOVED***

	// Try to start multiple transfers with the same ID
	transfers := make([]transferInfo, 5)
	for i := range transfers ***REMOVED***
		t := &transfers[i]
		t.progressChan = make(chan progress.Progress)
		t.progressDone = make(chan struct***REMOVED******REMOVED***)
		t.receivedFirstProgress = make(chan struct***REMOVED******REMOVED***)
		t.xfer, t.watcher = tm.Transfer("id1", makeXferFunc("id1"), progress.ChanOutput(t.progressChan))
		go progressConsumer(*t)
	***REMOVED***

	// Allow the transfer goroutine to proceed.
	close(ready)

	// Confirm that each watcher gets progress output.
	for _, t := range transfers ***REMOVED***
		<-t.receivedFirstProgress
	***REMOVED***

	// Confirm that the transfer function was called exactly once.
	if xferFuncCalls != 1 ***REMOVED***
		t.Fatal("transfer function wasn't called exactly once")
	***REMOVED***

	// Release one watcher every 5ms
	for _, t := range transfers ***REMOVED***
		t.xfer.Release(t.watcher)
		<-time.After(5 * time.Millisecond)
	***REMOVED***

	for _, t := range transfers ***REMOVED***
		// Now that all watchers have been released, Released() should
		// return a closed channel.
		<-t.xfer.Released()
		// Done() should return a closed channel because the xfer func returned
		// due to cancellation.
		<-t.xfer.Done()
	***REMOVED***

	for _, t := range transfers ***REMOVED***
		close(t.progressChan)
		<-t.progressDone
	***REMOVED***
***REMOVED***
