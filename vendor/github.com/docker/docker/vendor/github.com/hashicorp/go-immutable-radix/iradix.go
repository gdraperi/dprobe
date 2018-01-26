package iradix

import (
	"bytes"

	"github.com/hashicorp/golang-lru/simplelru"
)

const (
	// defaultModifiedCache is the default size of the modified node
	// cache used per transaction. This is used to cache the updates
	// to the nodes near the root, while the leaves do not need to be
	// cached. This is important for very large transactions to prevent
	// the modified cache from growing to be enormous.
	defaultModifiedCache = 8192
)

// Tree implements an immutable radix tree. This can be treated as a
// Dictionary abstract data type. The main advantage over a standard
// hash map is prefix-based lookups and ordered iteration. The immutability
// means that it is safe to concurrently read from a Tree without any
// coordination.
type Tree struct ***REMOVED***
	root *Node
	size int
***REMOVED***

// New returns an empty Tree
func New() *Tree ***REMOVED***
	t := &Tree***REMOVED***root: &Node***REMOVED******REMOVED******REMOVED***
	return t
***REMOVED***

// Len is used to return the number of elements in the tree
func (t *Tree) Len() int ***REMOVED***
	return t.size
***REMOVED***

// Txn is a transaction on the tree. This transaction is applied
// atomically and returns a new tree when committed. A transaction
// is not thread safe, and should only be used by a single goroutine.
type Txn struct ***REMOVED***
	root     *Node
	size     int
	modified *simplelru.LRU
***REMOVED***

// Txn starts a new transaction that can be used to mutate the tree
func (t *Tree) Txn() *Txn ***REMOVED***
	txn := &Txn***REMOVED***
		root: t.root,
		size: t.size,
	***REMOVED***
	return txn
***REMOVED***

// writeNode returns a node to be modified, if the current
// node as already been modified during the course of
// the transaction, it is used in-place.
func (t *Txn) writeNode(n *Node) *Node ***REMOVED***
	// Ensure the modified set exists
	if t.modified == nil ***REMOVED***
		lru, err := simplelru.NewLRU(defaultModifiedCache, nil)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		t.modified = lru
	***REMOVED***

	// If this node has already been modified, we can
	// continue to use it during this transaction.
	if _, ok := t.modified.Get(n); ok ***REMOVED***
		return n
	***REMOVED***

	// Copy the existing node
	nc := new(Node)
	if n.prefix != nil ***REMOVED***
		nc.prefix = make([]byte, len(n.prefix))
		copy(nc.prefix, n.prefix)
	***REMOVED***
	if n.leaf != nil ***REMOVED***
		nc.leaf = new(leafNode)
		*nc.leaf = *n.leaf
	***REMOVED***
	if len(n.edges) != 0 ***REMOVED***
		nc.edges = make([]edge, len(n.edges))
		copy(nc.edges, n.edges)
	***REMOVED***

	// Mark this node as modified
	t.modified.Add(n, nil)
	return nc
***REMOVED***

// insert does a recursive insertion
func (t *Txn) insert(n *Node, k, search []byte, v interface***REMOVED******REMOVED***) (*Node, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	// Handle key exhaution
	if len(search) == 0 ***REMOVED***
		nc := t.writeNode(n)
		if n.isLeaf() ***REMOVED***
			old := nc.leaf.val
			nc.leaf.val = v
			return nc, old, true
		***REMOVED*** else ***REMOVED***
			nc.leaf = &leafNode***REMOVED***
				key: k,
				val: v,
			***REMOVED***
			return nc, nil, false
		***REMOVED***
	***REMOVED***

	// Look for the edge
	idx, child := n.getEdge(search[0])

	// No edge, create one
	if child == nil ***REMOVED***
		e := edge***REMOVED***
			label: search[0],
			node: &Node***REMOVED***
				leaf: &leafNode***REMOVED***
					key: k,
					val: v,
				***REMOVED***,
				prefix: search,
			***REMOVED***,
		***REMOVED***
		nc := t.writeNode(n)
		nc.addEdge(e)
		return nc, nil, false
	***REMOVED***

	// Determine longest prefix of the search key on match
	commonPrefix := longestPrefix(search, child.prefix)
	if commonPrefix == len(child.prefix) ***REMOVED***
		search = search[commonPrefix:]
		newChild, oldVal, didUpdate := t.insert(child, k, search, v)
		if newChild != nil ***REMOVED***
			nc := t.writeNode(n)
			nc.edges[idx].node = newChild
			return nc, oldVal, didUpdate
		***REMOVED***
		return nil, oldVal, didUpdate
	***REMOVED***

	// Split the node
	nc := t.writeNode(n)
	splitNode := &Node***REMOVED***
		prefix: search[:commonPrefix],
	***REMOVED***
	nc.replaceEdge(edge***REMOVED***
		label: search[0],
		node:  splitNode,
	***REMOVED***)

	// Restore the existing child node
	modChild := t.writeNode(child)
	splitNode.addEdge(edge***REMOVED***
		label: modChild.prefix[commonPrefix],
		node:  modChild,
	***REMOVED***)
	modChild.prefix = modChild.prefix[commonPrefix:]

	// Create a new leaf node
	leaf := &leafNode***REMOVED***
		key: k,
		val: v,
	***REMOVED***

	// If the new key is a subset, add to to this node
	search = search[commonPrefix:]
	if len(search) == 0 ***REMOVED***
		splitNode.leaf = leaf
		return nc, nil, false
	***REMOVED***

	// Create a new edge for the node
	splitNode.addEdge(edge***REMOVED***
		label: search[0],
		node: &Node***REMOVED***
			leaf:   leaf,
			prefix: search,
		***REMOVED***,
	***REMOVED***)
	return nc, nil, false
***REMOVED***

// delete does a recursive deletion
func (t *Txn) delete(parent, n *Node, search []byte) (*Node, *leafNode) ***REMOVED***
	// Check for key exhaution
	if len(search) == 0 ***REMOVED***
		if !n.isLeaf() ***REMOVED***
			return nil, nil
		***REMOVED***

		// Remove the leaf node
		nc := t.writeNode(n)
		nc.leaf = nil

		// Check if this node should be merged
		if n != t.root && len(nc.edges) == 1 ***REMOVED***
			nc.mergeChild()
		***REMOVED***
		return nc, n.leaf
	***REMOVED***

	// Look for an edge
	label := search[0]
	idx, child := n.getEdge(label)
	if child == nil || !bytes.HasPrefix(search, child.prefix) ***REMOVED***
		return nil, nil
	***REMOVED***

	// Consume the search prefix
	search = search[len(child.prefix):]
	newChild, leaf := t.delete(n, child, search)
	if newChild == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	// Copy this node
	nc := t.writeNode(n)

	// Delete the edge if the node has no edges
	if newChild.leaf == nil && len(newChild.edges) == 0 ***REMOVED***
		nc.delEdge(label)
		if n != t.root && len(nc.edges) == 1 && !nc.isLeaf() ***REMOVED***
			nc.mergeChild()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		nc.edges[idx].node = newChild
	***REMOVED***
	return nc, leaf
***REMOVED***

// Insert is used to add or update a given key. The return provides
// the previous value and a bool indicating if any was set.
func (t *Txn) Insert(k []byte, v interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	newRoot, oldVal, didUpdate := t.insert(t.root, k, k, v)
	if newRoot != nil ***REMOVED***
		t.root = newRoot
	***REMOVED***
	if !didUpdate ***REMOVED***
		t.size++
	***REMOVED***
	return oldVal, didUpdate
***REMOVED***

// Delete is used to delete a given key. Returns the old value if any,
// and a bool indicating if the key was set.
func (t *Txn) Delete(k []byte) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	newRoot, leaf := t.delete(nil, t.root, k)
	if newRoot != nil ***REMOVED***
		t.root = newRoot
	***REMOVED***
	if leaf != nil ***REMOVED***
		t.size--
		return leaf.val, true
	***REMOVED***
	return nil, false
***REMOVED***

// Root returns the current root of the radix tree within this
// transaction. The root is not safe across insert and delete operations,
// but can be used to read the current state during a transaction.
func (t *Txn) Root() *Node ***REMOVED***
	return t.root
***REMOVED***

// Get is used to lookup a specific key, returning
// the value and if it was found
func (t *Txn) Get(k []byte) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	return t.root.Get(k)
***REMOVED***

// Commit is used to finalize the transaction and return a new tree
func (t *Txn) Commit() *Tree ***REMOVED***
	t.modified = nil
	return &Tree***REMOVED***t.root, t.size***REMOVED***
***REMOVED***

// Insert is used to add or update a given key. The return provides
// the new tree, previous value and a bool indicating if any was set.
func (t *Tree) Insert(k []byte, v interface***REMOVED******REMOVED***) (*Tree, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	txn := t.Txn()
	old, ok := txn.Insert(k, v)
	return txn.Commit(), old, ok
***REMOVED***

// Delete is used to delete a given key. Returns the new tree,
// old value if any, and a bool indicating if the key was set.
func (t *Tree) Delete(k []byte) (*Tree, interface***REMOVED******REMOVED***, bool) ***REMOVED***
	txn := t.Txn()
	old, ok := txn.Delete(k)
	return txn.Commit(), old, ok
***REMOVED***

// Root returns the root node of the tree which can be used for richer
// query operations.
func (t *Tree) Root() *Node ***REMOVED***
	return t.root
***REMOVED***

// Get is used to lookup a specific key, returning
// the value and if it was found
func (t *Tree) Get(k []byte) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	return t.root.Get(k)
***REMOVED***

// longestPrefix finds the length of the shared prefix
// of two strings
func longestPrefix(k1, k2 []byte) int ***REMOVED***
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

// concat two byte slices, returning a third new copy
func concat(a, b []byte) []byte ***REMOVED***
	c := make([]byte, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
***REMOVED***
