package iradix

import (
	"bytes"
	"sort"
)

// WalkFn is used when walking the tree. Takes a
// key and value, returning if iteration should
// be terminated.
type WalkFn func(k []byte, v interface***REMOVED******REMOVED***) bool

// leafNode is used to represent a value
type leafNode struct ***REMOVED***
	key []byte
	val interface***REMOVED******REMOVED***
***REMOVED***

// edge is used to represent an edge node
type edge struct ***REMOVED***
	label byte
	node  *Node
***REMOVED***

// Node is an immutable node in the radix tree
type Node struct ***REMOVED***
	// leaf is used to store possible leaf
	leaf *leafNode

	// prefix is the common prefix we ignore
	prefix []byte

	// Edges should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	edges edges
***REMOVED***

func (n *Node) isLeaf() bool ***REMOVED***
	return n.leaf != nil
***REMOVED***

func (n *Node) addEdge(e edge) ***REMOVED***
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool ***REMOVED***
		return n.edges[i].label >= e.label
	***REMOVED***)
	n.edges = append(n.edges, e)
	if idx != num ***REMOVED***
		copy(n.edges[idx+1:], n.edges[idx:num])
		n.edges[idx] = e
	***REMOVED***
***REMOVED***

func (n *Node) replaceEdge(e edge) ***REMOVED***
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool ***REMOVED***
		return n.edges[i].label >= e.label
	***REMOVED***)
	if idx < num && n.edges[idx].label == e.label ***REMOVED***
		n.edges[idx].node = e.node
		return
	***REMOVED***
	panic("replacing missing edge")
***REMOVED***

func (n *Node) getEdge(label byte) (int, *Node) ***REMOVED***
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool ***REMOVED***
		return n.edges[i].label >= label
	***REMOVED***)
	if idx < num && n.edges[idx].label == label ***REMOVED***
		return idx, n.edges[idx].node
	***REMOVED***
	return -1, nil
***REMOVED***

func (n *Node) delEdge(label byte) ***REMOVED***
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool ***REMOVED***
		return n.edges[i].label >= label
	***REMOVED***)
	if idx < num && n.edges[idx].label == label ***REMOVED***
		copy(n.edges[idx:], n.edges[idx+1:])
		n.edges[len(n.edges)-1] = edge***REMOVED******REMOVED***
		n.edges = n.edges[:len(n.edges)-1]
	***REMOVED***
***REMOVED***

func (n *Node) mergeChild() ***REMOVED***
	e := n.edges[0]
	child := e.node
	n.prefix = concat(n.prefix, child.prefix)
	if child.leaf != nil ***REMOVED***
		n.leaf = new(leafNode)
		*n.leaf = *child.leaf
	***REMOVED*** else ***REMOVED***
		n.leaf = nil
	***REMOVED***
	if len(child.edges) != 0 ***REMOVED***
		n.edges = make([]edge, len(child.edges))
		copy(n.edges, child.edges)
	***REMOVED*** else ***REMOVED***
		n.edges = nil
	***REMOVED***
***REMOVED***

func (n *Node) Get(k []byte) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	search := k
	for ***REMOVED***
		// Check for key exhaution
		if len(search) == 0 ***REMOVED***
			if n.isLeaf() ***REMOVED***
				return n.leaf.val, true
			***REMOVED***
			break
		***REMOVED***

		// Look for an edge
		_, n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			break
		***REMOVED***

		// Consume the search prefix
		if bytes.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

// LongestPrefix is like Get, but instead of an
// exact match, it will return the longest prefix match.
func (n *Node) LongestPrefix(k []byte) ([]byte, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	var last *leafNode
	search := k
	for ***REMOVED***
		// Look for a leaf node
		if n.isLeaf() ***REMOVED***
			last = n.leaf
		***REMOVED***

		// Check for key exhaution
		if len(search) == 0 ***REMOVED***
			break
		***REMOVED***

		// Look for an edge
		_, n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			break
		***REMOVED***

		// Consume the search prefix
		if bytes.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if last != nil ***REMOVED***
		return last.key, last.val, true
	***REMOVED***
	return nil, nil, false
***REMOVED***

// Minimum is used to return the minimum value in the tree
func (n *Node) Minimum() ([]byte, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	for ***REMOVED***
		if n.isLeaf() ***REMOVED***
			return n.leaf.key, n.leaf.val, true
		***REMOVED***
		if len(n.edges) > 0 ***REMOVED***
			n = n.edges[0].node
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return nil, nil, false
***REMOVED***

// Maximum is used to return the maximum value in the tree
func (n *Node) Maximum() ([]byte, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	for ***REMOVED***
		if num := len(n.edges); num > 0 ***REMOVED***
			n = n.edges[num-1].node
			continue
		***REMOVED***
		if n.isLeaf() ***REMOVED***
			return n.leaf.key, n.leaf.val, true
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return nil, nil, false
***REMOVED***

// Iterator is used to return an iterator at
// the given node to walk the tree
func (n *Node) Iterator() *Iterator ***REMOVED***
	return &Iterator***REMOVED***node: n***REMOVED***
***REMOVED***

// Walk is used to walk the tree
func (n *Node) Walk(fn WalkFn) ***REMOVED***
	recursiveWalk(n, fn)
***REMOVED***

// WalkPrefix is used to walk the tree under a prefix
func (n *Node) WalkPrefix(prefix []byte, fn WalkFn) ***REMOVED***
	search := prefix
	for ***REMOVED***
		// Check for key exhaution
		if len(search) == 0 ***REMOVED***
			recursiveWalk(n, fn)
			return
		***REMOVED***

		// Look for an edge
		_, n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			break
		***REMOVED***

		// Consume the search prefix
		if bytes.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]

		***REMOVED*** else if bytes.HasPrefix(n.prefix, search) ***REMOVED***
			// Child may be under our search prefix
			recursiveWalk(n, fn)
			return
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// WalkPath is used to walk the tree, but only visiting nodes
// from the root down to a given leaf. Where WalkPrefix walks
// all the entries *under* the given prefix, this walks the
// entries *above* the given prefix.
func (n *Node) WalkPath(path []byte, fn WalkFn) ***REMOVED***
	search := path
	for ***REMOVED***
		// Visit the leaf values if any
		if n.leaf != nil && fn(n.leaf.key, n.leaf.val) ***REMOVED***
			return
		***REMOVED***

		// Check for key exhaution
		if len(search) == 0 ***REMOVED***
			return
		***REMOVED***

		// Look for an edge
		_, n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			return
		***REMOVED***

		// Consume the search prefix
		if bytes.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// recursiveWalk is used to do a pre-order walk of a node
// recursively. Returns true if the walk should be aborted
func recursiveWalk(n *Node, fn WalkFn) bool ***REMOVED***
	// Visit the leaf values if any
	if n.leaf != nil && fn(n.leaf.key, n.leaf.val) ***REMOVED***
		return true
	***REMOVED***

	// Recurse on the children
	for _, e := range n.edges ***REMOVED***
		if recursiveWalk(e.node, fn) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
