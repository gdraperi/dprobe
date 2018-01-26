package radix

import (
	"sort"
	"strings"
)

// WalkFn is used when walking the tree. Takes a
// key and value, returning if iteration should
// be terminated.
type WalkFn func(s string, v interface***REMOVED******REMOVED***) bool

// leafNode is used to represent a value
type leafNode struct ***REMOVED***
	key string
	val interface***REMOVED******REMOVED***
***REMOVED***

// edge is used to represent an edge node
type edge struct ***REMOVED***
	label byte
	node  *node
***REMOVED***

type node struct ***REMOVED***
	// leaf is used to store possible leaf
	leaf *leafNode

	// prefix is the common prefix we ignore
	prefix string

	// Edges should be stored in-order for iteration.
	// We avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse
	edges edges
***REMOVED***

func (n *node) isLeaf() bool ***REMOVED***
	return n.leaf != nil
***REMOVED***

func (n *node) addEdge(e edge) ***REMOVED***
	n.edges = append(n.edges, e)
	n.edges.Sort()
***REMOVED***

func (n *node) replaceEdge(e edge) ***REMOVED***
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

func (n *node) getEdge(label byte) *node ***REMOVED***
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool ***REMOVED***
		return n.edges[i].label >= label
	***REMOVED***)
	if idx < num && n.edges[idx].label == label ***REMOVED***
		return n.edges[idx].node
	***REMOVED***
	return nil
***REMOVED***

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

// Tree implements a radix tree. This can be treated as a
// Dictionary abstract data type. The main advantage over
// a standard hash map is prefix-based lookups and
// ordered iteration,
type Tree struct ***REMOVED***
	root *node
	size int
***REMOVED***

// New returns an empty Tree
func New() *Tree ***REMOVED***
	return NewFromMap(nil)
***REMOVED***

// NewFromMap returns a new tree containing the keys
// from an existing map
func NewFromMap(m map[string]interface***REMOVED******REMOVED***) *Tree ***REMOVED***
	t := &Tree***REMOVED***root: &node***REMOVED******REMOVED******REMOVED***
	for k, v := range m ***REMOVED***
		t.Insert(k, v)
	***REMOVED***
	return t
***REMOVED***

// Len is used to return the number of elements in the tree
func (t *Tree) Len() int ***REMOVED***
	return t.size
***REMOVED***

// longestPrefix finds the length of the shared prefix
// of two strings
func longestPrefix(k1, k2 string) int ***REMOVED***
	max := len(k1)
	if l := len(k2); l < max ***REMOVED***
		max = l
	***REMOVED***
	var i int
	for i = 0; i < max; i++ ***REMOVED***
		if k1[i] != k2[i] ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return i
***REMOVED***

// Insert is used to add a newentry or update
// an existing entry. Returns if updated.
func (t *Tree) Insert(s string, v interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	var parent *node
	n := t.root
	search := s
	for ***REMOVED***
		// Handle key exhaution
		if len(search) == 0 ***REMOVED***
			if n.isLeaf() ***REMOVED***
				old := n.leaf.val
				n.leaf.val = v
				return old, true
			***REMOVED*** else ***REMOVED***
				n.leaf = &leafNode***REMOVED***
					key: s,
					val: v,
				***REMOVED***
				t.size++
				return nil, false
			***REMOVED***
		***REMOVED***

		// Look for the edge
		parent = n
		n = n.getEdge(search[0])

		// No edge, create one
		if n == nil ***REMOVED***
			e := edge***REMOVED***
				label: search[0],
				node: &node***REMOVED***
					leaf: &leafNode***REMOVED***
						key: s,
						val: v,
					***REMOVED***,
					prefix: search,
				***REMOVED***,
			***REMOVED***
			parent.addEdge(e)
			t.size++
			return nil, false
		***REMOVED***

		// Determine longest prefix of the search key on match
		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) ***REMOVED***
			search = search[commonPrefix:]
			continue
		***REMOVED***

		// Split the node
		t.size++
		child := &node***REMOVED***
			prefix: search[:commonPrefix],
		***REMOVED***
		parent.replaceEdge(edge***REMOVED***
			label: search[0],
			node:  child,
		***REMOVED***)

		// Restore the existing node
		child.addEdge(edge***REMOVED***
			label: n.prefix[commonPrefix],
			node:  n,
		***REMOVED***)
		n.prefix = n.prefix[commonPrefix:]

		// Create a new leaf node
		leaf := &leafNode***REMOVED***
			key: s,
			val: v,
		***REMOVED***

		// If the new key is a subset, add to to this node
		search = search[commonPrefix:]
		if len(search) == 0 ***REMOVED***
			child.leaf = leaf
			return nil, false
		***REMOVED***

		// Create a new edge for the node
		child.addEdge(edge***REMOVED***
			label: search[0],
			node: &node***REMOVED***
				leaf:   leaf,
				prefix: search,
			***REMOVED***,
		***REMOVED***)
		return nil, false
	***REMOVED***
	return nil, false
***REMOVED***

// Delete is used to delete a key, returning the previous
// value and if it was deleted
func (t *Tree) Delete(s string) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	n := t.root
	search := s
	for ***REMOVED***
		// Check for key exhaution
		if len(search) == 0 ***REMOVED***
			if !n.isLeaf() ***REMOVED***
				break
			***REMOVED***
			goto DELETE
		***REMOVED***

		// Look for an edge
		n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			break
		***REMOVED***

		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return nil, false

DELETE:
	// Delete the leaf
	leaf := n.leaf
	n.leaf = nil
	t.size--

	// Check if we should merge this node
	if len(n.edges) == 1 ***REMOVED***
		e := n.edges[0]
		child := e.node
		n.prefix = n.prefix + child.prefix
		n.leaf = child.leaf
		n.edges = child.edges
	***REMOVED***
	return leaf.val, true
***REMOVED***

// Get is used to lookup a specific key, returning
// the value and if it was found
func (t *Tree) Get(s string) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	n := t.root
	search := s
	for ***REMOVED***
		// Check for key exhaution
		if len(search) == 0 ***REMOVED***
			if n.isLeaf() ***REMOVED***
				return n.leaf.val, true
			***REMOVED***
			break
		***REMOVED***

		// Look for an edge
		n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			break
		***REMOVED***

		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

// LongestPrefix is like Get, but instead of an
// exact match, it will return the longest prefix match.
func (t *Tree) LongestPrefix(s string) (string, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	var last *leafNode
	n := t.root
	search := s
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
		n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			break
		***REMOVED***

		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if last != nil ***REMOVED***
		return last.key, last.val, true
	***REMOVED***
	return "", nil, false
***REMOVED***

// Minimum is used to return the minimum value in the tree
func (t *Tree) Minimum() (string, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	n := t.root
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
	return "", nil, false
***REMOVED***

// Maximum is used to return the maximum value in the tree
func (t *Tree) Maximum() (string, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	n := t.root
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
	return "", nil, false
***REMOVED***

// Walk is used to walk the tree
func (t *Tree) Walk(fn WalkFn) ***REMOVED***
	recursiveWalk(t.root, fn)
***REMOVED***

// WalkPrefix is used to walk the tree under a prefix
func (t *Tree) WalkPrefix(prefix string, fn WalkFn) ***REMOVED***
	n := t.root
	search := prefix
	for ***REMOVED***
		// Check for key exhaution
		if len(search) == 0 ***REMOVED***
			recursiveWalk(n, fn)
			return
		***REMOVED***

		// Look for an edge
		n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			break
		***REMOVED***

		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]

		***REMOVED*** else if strings.HasPrefix(n.prefix, search) ***REMOVED***
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
func (t *Tree) WalkPath(path string, fn WalkFn) ***REMOVED***
	n := t.root
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
		n = n.getEdge(search[0])
		if n == nil ***REMOVED***
			return
		***REMOVED***

		// Consume the search prefix
		if strings.HasPrefix(search, n.prefix) ***REMOVED***
			search = search[len(n.prefix):]
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// recursiveWalk is used to do a pre-order walk of a node
// recursively. Returns true if the walk should be aborted
func recursiveWalk(n *node, fn WalkFn) bool ***REMOVED***
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

// ToMap is used to walk the tree and convert it into a map
func (t *Tree) ToMap() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	out := make(map[string]interface***REMOVED******REMOVED***, t.size)
	t.Walk(func(k string, v interface***REMOVED******REMOVED***) bool ***REMOVED***
		out[k] = v
		return false
	***REMOVED***)
	return out
***REMOVED***
