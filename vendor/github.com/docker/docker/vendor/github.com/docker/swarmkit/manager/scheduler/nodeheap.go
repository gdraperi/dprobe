package scheduler

type nodeMaxHeap struct ***REMOVED***
	nodes    []NodeInfo
	lessFunc func(*NodeInfo, *NodeInfo) bool
	length   int
***REMOVED***

func (h nodeMaxHeap) Len() int ***REMOVED***
	return h.length
***REMOVED***

func (h nodeMaxHeap) Swap(i, j int) ***REMOVED***
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
***REMOVED***

func (h nodeMaxHeap) Less(i, j int) bool ***REMOVED***
	// reversed to make a max-heap
	return h.lessFunc(&h.nodes[j], &h.nodes[i])
***REMOVED***

func (h *nodeMaxHeap) Push(x interface***REMOVED******REMOVED***) ***REMOVED***
	h.nodes = append(h.nodes, x.(NodeInfo))
	h.length++
***REMOVED***

func (h *nodeMaxHeap) Pop() interface***REMOVED******REMOVED*** ***REMOVED***
	h.length--
	// return value is never used
	return nil
***REMOVED***
