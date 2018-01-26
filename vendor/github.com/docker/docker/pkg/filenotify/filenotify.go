// Package filenotify provides a mechanism for watching file(s) for changes.
// Generally leans on fsnotify, but provides a poll-based notifier which fsnotify does not support.
// These are wrapped up in a common interface so that either can be used interchangeably in your code.
package filenotify

import "github.com/fsnotify/fsnotify"

// FileWatcher is an interface for implementing file notification watchers
type FileWatcher interface ***REMOVED***
	Events() <-chan fsnotify.Event
	Errors() <-chan error
	Add(name string) error
	Remove(name string) error
	Close() error
***REMOVED***

// New tries to use an fs-event watcher, and falls back to the poller if there is an error
func New() (FileWatcher, error) ***REMOVED***
	if watcher, err := NewEventWatcher(); err == nil ***REMOVED***
		return watcher, nil
	***REMOVED***
	return NewPollingWatcher(), nil
***REMOVED***

// NewPollingWatcher returns a poll-based file watcher
func NewPollingWatcher() FileWatcher ***REMOVED***
	return &filePoller***REMOVED***
		events: make(chan fsnotify.Event),
		errors: make(chan error),
	***REMOVED***
***REMOVED***

// NewEventWatcher returns an fs-event based file watcher
func NewEventWatcher() (FileWatcher, error) ***REMOVED***
	watcher, err := fsnotify.NewWatcher()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &fsNotifyWatcher***REMOVED***watcher***REMOVED***, nil
***REMOVED***
