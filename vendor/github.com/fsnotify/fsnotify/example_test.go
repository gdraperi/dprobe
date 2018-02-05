// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package fsnotify_test

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func ExampleNewWatcher() ***REMOVED***
	watcher, err := fsnotify.NewWatcher()
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	defer watcher.Close()

	done := make(chan bool)
	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write ***REMOVED***
					log.Println("modified file:", event.Name)
				***REMOVED***
			case err := <-watcher.Errors:
				log.Println("error:", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	err = watcher.Add("/tmp/foo")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	<-done
***REMOVED***
