package xfer

import (
	"runtime"
	"sync"

	"github.com/docker/docker/pkg/progress"
	"golang.org/x/net/context"
)

// DoNotRetry is an error wrapper indicating that the error cannot be resolved
// with a retry.
type DoNotRetry struct ***REMOVED***
	Err error
***REMOVED***

// Error returns the stringified representation of the encapsulated error.
func (e DoNotRetry) Error() string ***REMOVED***
	return e.Err.Error()
***REMOVED***

// Watcher is returned by Watch and can be passed to Release to stop watching.
type Watcher struct ***REMOVED***
	// signalChan is used to signal to the watcher goroutine that
	// new progress information is available, or that the transfer
	// has finished.
	signalChan chan struct***REMOVED******REMOVED***
	// releaseChan signals to the watcher goroutine that the watcher
	// should be detached.
	releaseChan chan struct***REMOVED******REMOVED***
	// running remains open as long as the watcher is watching the
	// transfer. It gets closed if the transfer finishes or the
	// watcher is detached.
	running chan struct***REMOVED******REMOVED***
***REMOVED***

// Transfer represents an in-progress transfer.
type Transfer interface ***REMOVED***
	Watch(progressOutput progress.Output) *Watcher
	Release(*Watcher)
	Context() context.Context
	Close()
	Done() <-chan struct***REMOVED******REMOVED***
	Released() <-chan struct***REMOVED******REMOVED***
	Broadcast(masterProgressChan <-chan progress.Progress)
***REMOVED***

type transfer struct ***REMOVED***
	mu sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	// watchers keeps track of the goroutines monitoring progress output,
	// indexed by the channels that release them.
	watchers map[chan struct***REMOVED******REMOVED***]*Watcher

	// lastProgress is the most recently received progress event.
	lastProgress progress.Progress
	// hasLastProgress is true when lastProgress has been set.
	hasLastProgress bool

	// running remains open as long as the transfer is in progress.
	running chan struct***REMOVED******REMOVED***
	// released stays open until all watchers release the transfer and
	// the transfer is no longer tracked by the transfer manager.
	released chan struct***REMOVED******REMOVED***

	// broadcastDone is true if the master progress channel has closed.
	broadcastDone bool
	// closed is true if Close has been called
	closed bool
	// broadcastSyncChan allows watchers to "ping" the broadcasting
	// goroutine to wait for it for deplete its input channel. This ensures
	// a detaching watcher won't miss an event that was sent before it
	// started detaching.
	broadcastSyncChan chan struct***REMOVED******REMOVED***
***REMOVED***

// NewTransfer creates a new transfer.
func NewTransfer() Transfer ***REMOVED***
	t := &transfer***REMOVED***
		watchers:          make(map[chan struct***REMOVED******REMOVED***]*Watcher),
		running:           make(chan struct***REMOVED******REMOVED***),
		released:          make(chan struct***REMOVED******REMOVED***),
		broadcastSyncChan: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	// This uses context.Background instead of a caller-supplied context
	// so that a transfer won't be cancelled automatically if the client
	// which requested it is ^C'd (there could be other viewers).
	t.ctx, t.cancel = context.WithCancel(context.Background())

	return t
***REMOVED***

// Broadcast copies the progress and error output to all viewers.
func (t *transfer) Broadcast(masterProgressChan <-chan progress.Progress) ***REMOVED***
	for ***REMOVED***
		var (
			p  progress.Progress
			ok bool
		)
		select ***REMOVED***
		case p, ok = <-masterProgressChan:
		default:
			// We've depleted the channel, so now we can handle
			// reads on broadcastSyncChan to let detaching watchers
			// know we're caught up.
			select ***REMOVED***
			case <-t.broadcastSyncChan:
				continue
			case p, ok = <-masterProgressChan:
			***REMOVED***
		***REMOVED***

		t.mu.Lock()
		if ok ***REMOVED***
			t.lastProgress = p
			t.hasLastProgress = true
			for _, w := range t.watchers ***REMOVED***
				select ***REMOVED***
				case w.signalChan <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
				default:
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			t.broadcastDone = true
		***REMOVED***
		t.mu.Unlock()
		if !ok ***REMOVED***
			close(t.running)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Watch adds a watcher to the transfer. The supplied channel gets progress
// updates and is closed when the transfer finishes.
func (t *transfer) Watch(progressOutput progress.Output) *Watcher ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()

	w := &Watcher***REMOVED***
		releaseChan: make(chan struct***REMOVED******REMOVED***),
		signalChan:  make(chan struct***REMOVED******REMOVED***),
		running:     make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	t.watchers[w.releaseChan] = w

	if t.broadcastDone ***REMOVED***
		close(w.running)
		return w
	***REMOVED***

	go func() ***REMOVED***
		defer func() ***REMOVED***
			close(w.running)
		***REMOVED***()
		var (
			done           bool
			lastWritten    progress.Progress
			hasLastWritten bool
		)
		for ***REMOVED***
			t.mu.Lock()
			hasLastProgress := t.hasLastProgress
			lastProgress := t.lastProgress
			t.mu.Unlock()

			// Make sure we don't write the last progress item
			// twice.
			if hasLastProgress && (!done || !hasLastWritten || lastProgress != lastWritten) ***REMOVED***
				progressOutput.WriteProgress(lastProgress)
				lastWritten = lastProgress
				hasLastWritten = true
			***REMOVED***

			if done ***REMOVED***
				return
			***REMOVED***

			select ***REMOVED***
			case <-w.signalChan:
			case <-w.releaseChan:
				done = true
				// Since the watcher is going to detach, make
				// sure the broadcaster is caught up so we
				// don't miss anything.
				select ***REMOVED***
				case t.broadcastSyncChan <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
				case <-t.running:
				***REMOVED***
			case <-t.running:
				done = true
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return w
***REMOVED***

// Release is the inverse of Watch; indicating that the watcher no longer wants
// to be notified about the progress of the transfer. All calls to Watch must
// be paired with later calls to Release so that the lifecycle of the transfer
// is properly managed.
func (t *transfer) Release(watcher *Watcher) ***REMOVED***
	t.mu.Lock()
	delete(t.watchers, watcher.releaseChan)

	if len(t.watchers) == 0 ***REMOVED***
		if t.closed ***REMOVED***
			// released may have been closed already if all
			// watchers were released, then another one was added
			// while waiting for a previous watcher goroutine to
			// finish.
			select ***REMOVED***
			case <-t.released:
			default:
				close(t.released)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			t.cancel()
		***REMOVED***
	***REMOVED***
	t.mu.Unlock()

	close(watcher.releaseChan)
	// Block until the watcher goroutine completes
	<-watcher.running
***REMOVED***

// Done returns a channel which is closed if the transfer completes or is
// cancelled. Note that having 0 watchers causes a transfer to be cancelled.
func (t *transfer) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	// Note that this doesn't return t.ctx.Done() because that channel will
	// be closed the moment Cancel is called, and we need to return a
	// channel that blocks until a cancellation is actually acknowledged by
	// the transfer function.
	return t.running
***REMOVED***

// Released returns a channel which is closed once all watchers release the
// transfer AND the transfer is no longer tracked by the transfer manager.
func (t *transfer) Released() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return t.released
***REMOVED***

// Context returns the context associated with the transfer.
func (t *transfer) Context() context.Context ***REMOVED***
	return t.ctx
***REMOVED***

// Close is called by the transfer manager when the transfer is no longer
// being tracked.
func (t *transfer) Close() ***REMOVED***
	t.mu.Lock()
	t.closed = true
	if len(t.watchers) == 0 ***REMOVED***
		close(t.released)
	***REMOVED***
	t.mu.Unlock()
***REMOVED***

// DoFunc is a function called by the transfer manager to actually perform
// a transfer. It should be non-blocking. It should wait until the start channel
// is closed before transferring any data. If the function closes inactive, that
// signals to the transfer manager that the job is no longer actively moving
// data - for example, it may be waiting for a dependent transfer to finish.
// This prevents it from taking up a slot.
type DoFunc func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer

// TransferManager is used by LayerDownloadManager and LayerUploadManager to
// schedule and deduplicate transfers. It is up to the TransferManager
// implementation to make the scheduling and concurrency decisions.
type TransferManager interface ***REMOVED***
	// Transfer checks if a transfer with the given key is in progress. If
	// so, it returns progress and error output from that transfer.
	// Otherwise, it will call xferFunc to initiate the transfer.
	Transfer(key string, xferFunc DoFunc, progressOutput progress.Output) (Transfer, *Watcher)
	// SetConcurrency set the concurrencyLimit so that it is adjustable daemon reload
	SetConcurrency(concurrency int)
***REMOVED***

type transferManager struct ***REMOVED***
	mu sync.Mutex

	concurrencyLimit int
	activeTransfers  int
	transfers        map[string]Transfer
	waitingTransfers []chan struct***REMOVED******REMOVED***
***REMOVED***

// NewTransferManager returns a new TransferManager.
func NewTransferManager(concurrencyLimit int) TransferManager ***REMOVED***
	return &transferManager***REMOVED***
		concurrencyLimit: concurrencyLimit,
		transfers:        make(map[string]Transfer),
	***REMOVED***
***REMOVED***

// SetConcurrency sets the concurrencyLimit
func (tm *transferManager) SetConcurrency(concurrency int) ***REMOVED***
	tm.mu.Lock()
	tm.concurrencyLimit = concurrency
	tm.mu.Unlock()
***REMOVED***

// Transfer checks if a transfer matching the given key is in progress. If not,
// it starts one by calling xferFunc. The caller supplies a channel which
// receives progress output from the transfer.
func (tm *transferManager) Transfer(key string, xferFunc DoFunc, progressOutput progress.Output) (Transfer, *Watcher) ***REMOVED***
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for ***REMOVED***
		xfer, present := tm.transfers[key]
		if !present ***REMOVED***
			break
		***REMOVED***
		// Transfer is already in progress.
		watcher := xfer.Watch(progressOutput)

		select ***REMOVED***
		case <-xfer.Context().Done():
			// We don't want to watch a transfer that has been cancelled.
			// Wait for it to be removed from the map and try again.
			xfer.Release(watcher)
			tm.mu.Unlock()
			// The goroutine that removes this transfer from the
			// map is also waiting for xfer.Done(), so yield to it.
			// This could be avoided by adding a Closed method
			// to Transfer to allow explicitly waiting for it to be
			// removed the map, but forcing a scheduling round in
			// this very rare case seems better than bloating the
			// interface definition.
			runtime.Gosched()
			<-xfer.Done()
			tm.mu.Lock()
		default:
			return xfer, watcher
		***REMOVED***
	***REMOVED***

	start := make(chan struct***REMOVED******REMOVED***)
	inactive := make(chan struct***REMOVED******REMOVED***)

	if tm.concurrencyLimit == 0 || tm.activeTransfers < tm.concurrencyLimit ***REMOVED***
		close(start)
		tm.activeTransfers++
	***REMOVED*** else ***REMOVED***
		tm.waitingTransfers = append(tm.waitingTransfers, start)
	***REMOVED***

	masterProgressChan := make(chan progress.Progress)
	xfer := xferFunc(masterProgressChan, start, inactive)
	watcher := xfer.Watch(progressOutput)
	go xfer.Broadcast(masterProgressChan)
	tm.transfers[key] = xfer

	// When the transfer is finished, remove from the map.
	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-inactive:
				tm.mu.Lock()
				tm.inactivate(start)
				tm.mu.Unlock()
				inactive = nil
			case <-xfer.Done():
				tm.mu.Lock()
				if inactive != nil ***REMOVED***
					tm.inactivate(start)
				***REMOVED***
				delete(tm.transfers, key)
				tm.mu.Unlock()
				xfer.Close()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return xfer, watcher
***REMOVED***

func (tm *transferManager) inactivate(start chan struct***REMOVED******REMOVED***) ***REMOVED***
	// If the transfer was started, remove it from the activeTransfers
	// count.
	select ***REMOVED***
	case <-start:
		// Start next transfer if any are waiting
		if len(tm.waitingTransfers) != 0 ***REMOVED***
			close(tm.waitingTransfers[0])
			tm.waitingTransfers = tm.waitingTransfers[1:]
		***REMOVED*** else ***REMOVED***
			tm.activeTransfers--
		***REMOVED***
	default:
	***REMOVED***
***REMOVED***
