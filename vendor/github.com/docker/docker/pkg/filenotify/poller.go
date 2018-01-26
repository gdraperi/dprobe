package filenotify

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

var (
	// errPollerClosed is returned when the poller is closed
	errPollerClosed = errors.New("poller is closed")
	// errNoSuchWatch is returned when trying to remove a watch that doesn't exist
	errNoSuchWatch = errors.New("watch does not exist")
)

// watchWaitTime is the time to wait between file poll loops
const watchWaitTime = 200 * time.Millisecond

// filePoller is used to poll files for changes, especially in cases where fsnotify
// can't be run (e.g. when inotify handles are exhausted)
// filePoller satisfies the FileWatcher interface
type filePoller struct ***REMOVED***
	// watches is the list of files currently being polled, close the associated channel to stop the watch
	watches map[string]chan struct***REMOVED******REMOVED***
	// events is the channel to listen to for watch events
	events chan fsnotify.Event
	// errors is the channel to listen to for watch errors
	errors chan error
	// mu locks the poller for modification
	mu sync.Mutex
	// closed is used to specify when the poller has already closed
	closed bool
***REMOVED***

// Add adds a filename to the list of watches
// once added the file is polled for changes in a separate goroutine
func (w *filePoller) Add(name string) error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed ***REMOVED***
		return errPollerClosed
	***REMOVED***

	f, err := os.Open(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	fi, err := os.Stat(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if w.watches == nil ***REMOVED***
		w.watches = make(map[string]chan struct***REMOVED******REMOVED***)
	***REMOVED***
	if _, exists := w.watches[name]; exists ***REMOVED***
		return fmt.Errorf("watch exists")
	***REMOVED***
	chClose := make(chan struct***REMOVED******REMOVED***)
	w.watches[name] = chClose

	go w.watch(f, fi, chClose)
	return nil
***REMOVED***

// Remove stops and removes watch with the specified name
func (w *filePoller) Remove(name string) error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.remove(name)
***REMOVED***

func (w *filePoller) remove(name string) error ***REMOVED***
	if w.closed ***REMOVED***
		return errPollerClosed
	***REMOVED***

	chClose, exists := w.watches[name]
	if !exists ***REMOVED***
		return errNoSuchWatch
	***REMOVED***
	close(chClose)
	delete(w.watches, name)
	return nil
***REMOVED***

// Events returns the event channel
// This is used for notifications on events about watched files
func (w *filePoller) Events() <-chan fsnotify.Event ***REMOVED***
	return w.events
***REMOVED***

// Errors returns the errors channel
// This is used for notifications about errors on watched files
func (w *filePoller) Errors() <-chan error ***REMOVED***
	return w.errors
***REMOVED***

// Close closes the poller
// All watches are stopped, removed, and the poller cannot be added to
func (w *filePoller) Close() error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed ***REMOVED***
		return nil
	***REMOVED***

	w.closed = true
	for name := range w.watches ***REMOVED***
		w.remove(name)
		delete(w.watches, name)
	***REMOVED***
	return nil
***REMOVED***

// sendEvent publishes the specified event to the events channel
func (w *filePoller) sendEvent(e fsnotify.Event, chClose <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	select ***REMOVED***
	case w.events <- e:
	case <-chClose:
		return fmt.Errorf("closed")
	***REMOVED***
	return nil
***REMOVED***

// sendErr publishes the specified error to the errors channel
func (w *filePoller) sendErr(e error, chClose <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	select ***REMOVED***
	case w.errors <- e:
	case <-chClose:
		return fmt.Errorf("closed")
	***REMOVED***
	return nil
***REMOVED***

// watch is responsible for polling the specified file for changes
// upon finding changes to a file or errors, sendEvent/sendErr is called
func (w *filePoller) watch(f *os.File, lastFi os.FileInfo, chClose chan struct***REMOVED******REMOVED***) ***REMOVED***
	defer f.Close()
	for ***REMOVED***
		time.Sleep(watchWaitTime)
		select ***REMOVED***
		case <-chClose:
			logrus.Debugf("watch for %s closed", f.Name())
			return
		default:
		***REMOVED***

		fi, err := os.Stat(f.Name())
		if err != nil ***REMOVED***
			// if we got an error here and lastFi is not set, we can presume that nothing has changed
			// This should be safe since before `watch()` is called, a stat is performed, there is any error `watch` is not called
			if lastFi == nil ***REMOVED***
				continue
			***REMOVED***
			// If it doesn't exist at this point, it must have been removed
			// no need to send the error here since this is a valid operation
			if os.IsNotExist(err) ***REMOVED***
				if err := w.sendEvent(fsnotify.Event***REMOVED***Op: fsnotify.Remove, Name: f.Name()***REMOVED***, chClose); err != nil ***REMOVED***
					return
				***REMOVED***
				lastFi = nil
				continue
			***REMOVED***
			// at this point, send the error
			if err := w.sendErr(err, chClose); err != nil ***REMOVED***
				return
			***REMOVED***
			continue
		***REMOVED***

		if lastFi == nil ***REMOVED***
			if err := w.sendEvent(fsnotify.Event***REMOVED***Op: fsnotify.Create, Name: fi.Name()***REMOVED***, chClose); err != nil ***REMOVED***
				return
			***REMOVED***
			lastFi = fi
			continue
		***REMOVED***

		if fi.Mode() != lastFi.Mode() ***REMOVED***
			if err := w.sendEvent(fsnotify.Event***REMOVED***Op: fsnotify.Chmod, Name: fi.Name()***REMOVED***, chClose); err != nil ***REMOVED***
				return
			***REMOVED***
			lastFi = fi
			continue
		***REMOVED***

		if fi.ModTime() != lastFi.ModTime() || fi.Size() != lastFi.Size() ***REMOVED***
			if err := w.sendEvent(fsnotify.Event***REMOVED***Op: fsnotify.Write, Name: fi.Name()***REMOVED***, chClose); err != nil ***REMOVED***
				return
			***REMOVED***
			lastFi = fi
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***
