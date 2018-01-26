package iradix

import "bytes"

// Iterator is used to iterate over a set of nodes
// in pre-order
type Iterator struct ***REMOVED***
	node  *Node
	stack []edges
***REMOVED***

// SeekPrefix is used to seek the iterator to a given prefix
func (i *Iterator) SeekPrefix(prefix []byte) ***REMOVED***
	// Wipe the stack
	i.stack = nil
	n := i.node
	search := prefix
	for ***REMOVED***
		// Check for key exhaution
		if len(search) == 0 ***REMOVED***
			i.node = n
			return
		***REMOVED***

		// Look for an edge
		_, n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			i.node = nil
			return
		***REMOVED***

		// Consume the search prefix
		if bytes.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]

		***REMOVED*** else if bytes.HasPrefix(n.prefix, search) ***REMOVED***
			i.node = n
			return
		***REMOVED*** else ***REMOVED***
			i.node = nil
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Next returns the next node in order
func (i *Iterator) Next() ([]byte, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	// Initialize our stack if needed
	if i.stack == nil && i.node != nil ***REMOVED***
		i.stack = []edges***REMOVED***
			edges***REMOVED***
				edge***REMOVED***node: i.node***REMOVED***,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	for len(i.stack) > 0 ***REMOVED***
		// Inspect the last element of the stack
		n := len(i.stack)
		last := i.stack[n-1]
		elem := last[0].node

		// Update the stack
		if len(last) > 1 ***REMOVED***
			i.stack[n-1] = last[1:]
		***REMOVED*** else ***REMOVED***
			i.stack = i.stack[:n-1]
		***REMOVED***

		// Push the edges onto the frontier
		if len(elem.edges) > 0 ***REMOVED***
			i.stack = append(i.stack, elem.edges)
		***REMOVED***

		// Return the leaf values if any
		if elem.leaf != nil ***REMOVED***
			return elem.leaf.key, elem.leaf.val, true
		***REMOVED***
	***REMOVED***
	return nil, nil, false
***REMOVED***
