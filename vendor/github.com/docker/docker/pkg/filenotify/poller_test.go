package filenotify

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestPollerAddRemove(t *testing.T) ***REMOVED***
	w := NewPollingWatcher()

	if err := w.Add("no-such-file"); err == nil ***REMOVED***
		t.Fatal("should have gotten error when adding a non-existent file")
	***REMOVED***
	if err := w.Remove("no-such-file"); err == nil ***REMOVED***
		t.Fatal("should have gotten error when removing non-existent watch")
	***REMOVED***

	f, err := ioutil.TempFile("", "asdf")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(f.Name())

	if err := w.Add(f.Name()); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := w.Remove(f.Name()); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestPollerEvent(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("No chmod on Windows")
	***REMOVED***
	w := NewPollingWatcher()

	f, err := ioutil.TempFile("", "test-poller")
	if err != nil ***REMOVED***
		t.Fatal("error creating temp file")
	***REMOVED***
	defer os.RemoveAll(f.Name())
	f.Close()

	if err := w.Add(f.Name()); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	select ***REMOVED***
	case <-w.Events():
		t.Fatal("got event before anything happened")
	case <-w.Errors():
		t.Fatal("got error before anything happened")
	default:
	***REMOVED***

	if err := ioutil.WriteFile(f.Name(), []byte("hello"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := assertEvent(w, fsnotify.Write); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Chmod(f.Name(), 600); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := assertEvent(w, fsnotify.Chmod); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Remove(f.Name()); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := assertEvent(w, fsnotify.Remove); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestPollerClose(t *testing.T) ***REMOVED***
	w := NewPollingWatcher()
	if err := w.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// test double-close
	if err := w.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	f, err := ioutil.TempFile("", "asdf")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(f.Name())
	if err := w.Add(f.Name()); err == nil ***REMOVED***
		t.Fatal("should have gotten error adding watch for closed watcher")
	***REMOVED***
***REMOVED***

func assertEvent(w FileWatcher, eType fsnotify.Op) error ***REMOVED***
	var err error
	select ***REMOVED***
	case e := <-w.Events():
		if e.Op != eType ***REMOVED***
			err = fmt.Errorf("got wrong event type, expected %q: %v", eType, e.Op)
		***REMOVED***
	case e := <-w.Errors():
		err = fmt.Errorf("got unexpected error waiting for events %v: %v", eType, e)
	case <-time.After(watchWaitTime * 3):
		err = fmt.Errorf("timeout waiting for event %v", eType)
	***REMOVED***
	return err
***REMOVED***
