package libcontainerd

import "sync"

type queue struct ***REMOVED***
	sync.Mutex
	fns map[string]chan struct***REMOVED******REMOVED***
***REMOVED***

func (q *queue) append(id string, f func()) ***REMOVED***
	q.Lock()
	defer q.Unlock()

	if q.fns == nil ***REMOVED***
		q.fns = make(map[string]chan struct***REMOVED******REMOVED***)
	***REMOVED***

	done := make(chan struct***REMOVED******REMOVED***)

	fn, ok := q.fns[id]
	q.fns[id] = done
	go func() ***REMOVED***
		if ok ***REMOVED***
			<-fn
		***REMOVED***
		f()
		close(done)

		q.Lock()
		if q.fns[id] == done ***REMOVED***
			delete(q.fns, id)
		***REMOVED***
		q.Unlock()
	***REMOVED***()
***REMOVED***
