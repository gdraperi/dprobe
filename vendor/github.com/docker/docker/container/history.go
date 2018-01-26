package container

import "sort"

// History is a convenience type for storing a list of containers,
// sorted by creation date in descendant order.
type History []*Container

// Len returns the number of containers in the history.
func (history *History) Len() int ***REMOVED***
	return len(*history)
***REMOVED***

// Less compares two containers and returns true if the second one
// was created before the first one.
func (history *History) Less(i, j int) bool ***REMOVED***
	containers := *history
	return containers[j].Created.Before(containers[i].Created)
***REMOVED***

// Swap switches containers i and j positions in the history.
func (history *History) Swap(i, j int) ***REMOVED***
	containers := *history
	containers[i], containers[j] = containers[j], containers[i]
***REMOVED***

// sort orders the history by creation date in descendant order.
func (history *History) sort() ***REMOVED***
	sort.Sort(history)
***REMOVED***
