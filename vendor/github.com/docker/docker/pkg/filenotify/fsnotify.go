package filenotify

import "github.com/fsnotify/fsnotify"

// fsNotifyWatcher wraps the fsnotify package to satisfy the FileNotifier interface
type fsNotifyWatcher struct ***REMOVED***
	*fsnotify.Watcher
***REMOVED***

// Events returns the fsnotify event channel receiver
func (w *fsNotifyWatcher) Events() <-chan fsnotify.Event ***REMOVED***
	return w.Watcher.Events
***REMOVED***

// Errors returns the fsnotify error channel receiver
func (w *fsNotifyWatcher) Errors() <-chan error ***REMOVED***
	return w.Watcher.Errors
***REMOVED***
