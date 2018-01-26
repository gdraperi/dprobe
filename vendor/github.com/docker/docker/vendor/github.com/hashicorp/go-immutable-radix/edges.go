package iradix

import "sort"

type edges []edge

func (e edges) Len() int ***REMOVED***
	return len(e)
***REMOVED***

func (e edges) Less(i, j int) bool ***REMOVED***
	return e[i].label < e[j].label
***REMOVED***

func (e edges) Swap(i, j int) ***REMOVED***
	e[i], e[j] = e[j], e[i]
***REMOVED***

func (e edges) Sort() ***REMOVED***
	sort.Sort(e)
***REMOVED***
