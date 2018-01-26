package raft

import (
	"fmt"
	"sync"
)

type waitItem struct ***REMOVED***
	// channel to wait up the waiter
	ch chan interface***REMOVED******REMOVED***
	// callback which is called synchronously when the wait is triggered
	cb func()
	// callback which is called to cancel a waiter
	cancel func()
***REMOVED***

type wait struct ***REMOVED***
	l sync.Mutex
	m map[uint64]waitItem
***REMOVED***

func newWait() *wait ***REMOVED***
	return &wait***REMOVED***m: make(map[uint64]waitItem)***REMOVED***
***REMOVED***

func (w *wait) register(id uint64, cb func(), cancel func()) <-chan interface***REMOVED******REMOVED*** ***REMOVED***
	w.l.Lock()
	defer w.l.Unlock()
	_, ok := w.m[id]
	if !ok ***REMOVED***
		ch := make(chan interface***REMOVED******REMOVED***, 1)
		w.m[id] = waitItem***REMOVED***ch: ch, cb: cb, cancel: cancel***REMOVED***
		return ch
	***REMOVED***
	panic(fmt.Sprintf("duplicate id %x", id))
***REMOVED***

func (w *wait) trigger(id uint64, x interface***REMOVED******REMOVED***) bool ***REMOVED***
	w.l.Lock()
	waitItem, ok := w.m[id]
	delete(w.m, id)
	w.l.Unlock()
	if ok ***REMOVED***
		if waitItem.cb != nil ***REMOVED***
			waitItem.cb()
		***REMOVED***
		waitItem.ch <- x
		return true
	***REMOVED***
	return false
***REMOVED***

func (w *wait) cancel(id uint64) ***REMOVED***
	w.l.Lock()
	waitItem, ok := w.m[id]
	delete(w.m, id)
	w.l.Unlock()
	if ok ***REMOVED***
		if waitItem.cancel != nil ***REMOVED***
			waitItem.cancel()
		***REMOVED***
		close(waitItem.ch)
	***REMOVED***
***REMOVED***

func (w *wait) cancelAll() ***REMOVED***
	w.l.Lock()
	defer w.l.Unlock()

	for id, waitItem := range w.m ***REMOVED***
		delete(w.m, id)
		if waitItem.cancel != nil ***REMOVED***
			waitItem.cancel()
		***REMOVED***
		close(waitItem.ch)
	***REMOVED***
***REMOVED***
