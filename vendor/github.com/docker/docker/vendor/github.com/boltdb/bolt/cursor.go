package bolt

import (
	"bytes"
	"fmt"
	"sort"
)

// Cursor represents an iterator that can traverse over all key/value pairs in a bucket in sorted order.
// Cursors see nested buckets with value == nil.
// Cursors can be obtained from a transaction and are valid as long as the transaction is open.
//
// Keys and values returned from the cursor are only valid for the life of the transaction.
//
// Changing data while traversing with a cursor may cause it to be invalidated
// and return unexpected keys and/or values. You must reposition your cursor
// after mutating data.
type Cursor struct ***REMOVED***
	bucket *Bucket
	stack  []elemRef
***REMOVED***

// Bucket returns the bucket that this cursor was created from.
func (c *Cursor) Bucket() *Bucket ***REMOVED***
	return c.bucket
***REMOVED***

// First moves the cursor to the first item in the bucket and returns its key and value.
// If the bucket is empty then a nil key and value are returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) First() (key []byte, value []byte) ***REMOVED***
	_assert(c.bucket.tx.db != nil, "tx closed")
	c.stack = c.stack[:0]
	p, n := c.bucket.pageNode(c.bucket.root)
	c.stack = append(c.stack, elemRef***REMOVED***page: p, node: n, index: 0***REMOVED***)
	c.first()

	// If we land on an empty page then move to the next value.
	// https://github.com/boltdb/bolt/issues/450
	if c.stack[len(c.stack)-1].count() == 0 ***REMOVED***
		c.next()
	***REMOVED***

	k, v, flags := c.keyValue()
	if (flags & uint32(bucketLeafFlag)) != 0 ***REMOVED***
		return k, nil
	***REMOVED***
	return k, v

***REMOVED***

// Last moves the cursor to the last item in the bucket and returns its key and value.
// If the bucket is empty then a nil key and value are returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) Last() (key []byte, value []byte) ***REMOVED***
	_assert(c.bucket.tx.db != nil, "tx closed")
	c.stack = c.stack[:0]
	p, n := c.bucket.pageNode(c.bucket.root)
	ref := elemRef***REMOVED***page: p, node: n***REMOVED***
	ref.index = ref.count() - 1
	c.stack = append(c.stack, ref)
	c.last()
	k, v, flags := c.keyValue()
	if (flags & uint32(bucketLeafFlag)) != 0 ***REMOVED***
		return k, nil
	***REMOVED***
	return k, v
***REMOVED***

// Next moves the cursor to the next item in the bucket and returns its key and value.
// If the cursor is at the end of the bucket then a nil key and value are returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) Next() (key []byte, value []byte) ***REMOVED***
	_assert(c.bucket.tx.db != nil, "tx closed")
	k, v, flags := c.next()
	if (flags & uint32(bucketLeafFlag)) != 0 ***REMOVED***
		return k, nil
	***REMOVED***
	return k, v
***REMOVED***

// Prev moves the cursor to the previous item in the bucket and returns its key and value.
// If the cursor is at the beginning of the bucket then a nil key and value are returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) Prev() (key []byte, value []byte) ***REMOVED***
	_assert(c.bucket.tx.db != nil, "tx closed")

	// Attempt to move back one element until we're successful.
	// Move up the stack as we hit the beginning of each page in our stack.
	for i := len(c.stack) - 1; i >= 0; i-- ***REMOVED***
		elem := &c.stack[i]
		if elem.index > 0 ***REMOVED***
			elem.index--
			break
		***REMOVED***
		c.stack = c.stack[:i]
	***REMOVED***

	// If we've hit the end then return nil.
	if len(c.stack) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	// Move down the stack to find the last element of the last leaf under this branch.
	c.last()
	k, v, flags := c.keyValue()
	if (flags & uint32(bucketLeafFlag)) != 0 ***REMOVED***
		return k, nil
	***REMOVED***
	return k, v
***REMOVED***

// Seek moves the cursor to a given key and returns it.
// If the key does not exist then the next key is used. If no keys
// follow, a nil key is returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) Seek(seek []byte) (key []byte, value []byte) ***REMOVED***
	k, v, flags := c.seek(seek)

	// If we ended up after the last element of a page then move to the next one.
	if ref := &c.stack[len(c.stack)-1]; ref.index >= ref.count() ***REMOVED***
		k, v, flags = c.next()
	***REMOVED***

	if k == nil ***REMOVED***
		return nil, nil
	***REMOVED*** else if (flags & uint32(bucketLeafFlag)) != 0 ***REMOVED***
		return k, nil
	***REMOVED***
	return k, v
***REMOVED***

// Delete removes the current key/value under the cursor from the bucket.
// Delete fails if current key/value is a bucket or if the transaction is not writable.
func (c *Cursor) Delete() error ***REMOVED***
	if c.bucket.tx.db == nil ***REMOVED***
		return ErrTxClosed
	***REMOVED*** else if !c.bucket.Writable() ***REMOVED***
		return ErrTxNotWritable
	***REMOVED***

	key, _, flags := c.keyValue()
	// Return an error if current value is a bucket.
	if (flags & bucketLeafFlag) != 0 ***REMOVED***
		return ErrIncompatibleValue
	***REMOVED***
	c.node().del(key)

	return nil
***REMOVED***

// seek moves the cursor to a given key and returns it.
// If the key does not exist then the next key is used.
func (c *Cursor) seek(seek []byte) (key []byte, value []byte, flags uint32) ***REMOVED***
	_assert(c.bucket.tx.db != nil, "tx closed")

	// Start from root page/node and traverse to correct page.
	c.stack = c.stack[:0]
	c.search(seek, c.bucket.root)
	ref := &c.stack[len(c.stack)-1]

	// If the cursor is pointing to the end of page/node then return nil.
	if ref.index >= ref.count() ***REMOVED***
		return nil, nil, 0
	***REMOVED***

	// If this is a bucket then return a nil value.
	return c.keyValue()
***REMOVED***

// first moves the cursor to the first leaf element under the last page in the stack.
func (c *Cursor) first() ***REMOVED***
	for ***REMOVED***
		// Exit when we hit a leaf page.
		var ref = &c.stack[len(c.stack)-1]
		if ref.isLeaf() ***REMOVED***
			break
		***REMOVED***

		// Keep adding pages pointing to the first element to the stack.
		var pgid pgid
		if ref.node != nil ***REMOVED***
			pgid = ref.node.inodes[ref.index].pgid
		***REMOVED*** else ***REMOVED***
			pgid = ref.page.branchPageElement(uint16(ref.index)).pgid
		***REMOVED***
		p, n := c.bucket.pageNode(pgid)
		c.stack = append(c.stack, elemRef***REMOVED***page: p, node: n, index: 0***REMOVED***)
	***REMOVED***
***REMOVED***

// last moves the cursor to the last leaf element under the last page in the stack.
func (c *Cursor) last() ***REMOVED***
	for ***REMOVED***
		// Exit when we hit a leaf page.
		ref := &c.stack[len(c.stack)-1]
		if ref.isLeaf() ***REMOVED***
			break
		***REMOVED***

		// Keep adding pages pointing to the last element in the stack.
		var pgid pgid
		if ref.node != nil ***REMOVED***
			pgid = ref.node.inodes[ref.index].pgid
		***REMOVED*** else ***REMOVED***
			pgid = ref.page.branchPageElement(uint16(ref.index)).pgid
		***REMOVED***
		p, n := c.bucket.pageNode(pgid)

		var nextRef = elemRef***REMOVED***page: p, node: n***REMOVED***
		nextRef.index = nextRef.count() - 1
		c.stack = append(c.stack, nextRef)
	***REMOVED***
***REMOVED***

// next moves to the next leaf element and returns the key and value.
// If the cursor is at the last leaf element then it stays there and returns nil.
func (c *Cursor) next() (key []byte, value []byte, flags uint32) ***REMOVED***
	for ***REMOVED***
		// Attempt to move over one element until we're successful.
		// Move up the stack as we hit the end of each page in our stack.
		var i int
		for i = len(c.stack) - 1; i >= 0; i-- ***REMOVED***
			elem := &c.stack[i]
			if elem.index < elem.count()-1 ***REMOVED***
				elem.index++
				break
			***REMOVED***
		***REMOVED***

		// If we've hit the root page then stop and return. This will leave the
		// cursor on the last element of the last page.
		if i == -1 ***REMOVED***
			return nil, nil, 0
		***REMOVED***

		// Otherwise start from where we left off in the stack and find the
		// first element of the first leaf page.
		c.stack = c.stack[:i+1]
		c.first()

		// If this is an empty page then restart and move back up the stack.
		// https://github.com/boltdb/bolt/issues/450
		if c.stack[len(c.stack)-1].count() == 0 ***REMOVED***
			continue
		***REMOVED***

		return c.keyValue()
	***REMOVED***
***REMOVED***

// search recursively performs a binary search against a given page/node until it finds a given key.
func (c *Cursor) search(key []byte, pgid pgid) ***REMOVED***
	p, n := c.bucket.pageNode(pgid)
	if p != nil && (p.flags&(branchPageFlag|leafPageFlag)) == 0 ***REMOVED***
		panic(fmt.Sprintf("invalid page type: %d: %x", p.id, p.flags))
	***REMOVED***
	e := elemRef***REMOVED***page: p, node: n***REMOVED***
	c.stack = append(c.stack, e)

	// If we're on a leaf page/node then find the specific node.
	if e.isLeaf() ***REMOVED***
		c.nsearch(key)
		return
	***REMOVED***

	if n != nil ***REMOVED***
		c.searchNode(key, n)
		return
	***REMOVED***
	c.searchPage(key, p)
***REMOVED***

func (c *Cursor) searchNode(key []byte, n *node) ***REMOVED***
	var exact bool
	index := sort.Search(len(n.inodes), func(i int) bool ***REMOVED***
		// TODO(benbjohnson): Optimize this range search. It's a bit hacky right now.
		// sort.Search() finds the lowest index where f() != -1 but we need the highest index.
		ret := bytes.Compare(n.inodes[i].key, key)
		if ret == 0 ***REMOVED***
			exact = true
		***REMOVED***
		return ret != -1
	***REMOVED***)
	if !exact && index > 0 ***REMOVED***
		index--
	***REMOVED***
	c.stack[len(c.stack)-1].index = index

	// Recursively search to the next page.
	c.search(key, n.inodes[index].pgid)
***REMOVED***

func (c *Cursor) searchPage(key []byte, p *page) ***REMOVED***
	// Binary search for the correct range.
	inodes := p.branchPageElements()

	var exact bool
	index := sort.Search(int(p.count), func(i int) bool ***REMOVED***
		// TODO(benbjohnson): Optimize this range search. It's a bit hacky right now.
		// sort.Search() finds the lowest index where f() != -1 but we need the highest index.
		ret := bytes.Compare(inodes[i].key(), key)
		if ret == 0 ***REMOVED***
			exact = true
		***REMOVED***
		return ret != -1
	***REMOVED***)
	if !exact && index > 0 ***REMOVED***
		index--
	***REMOVED***
	c.stack[len(c.stack)-1].index = index

	// Recursively search to the next page.
	c.search(key, inodes[index].pgid)
***REMOVED***

// nsearch searches the leaf node on the top of the stack for a key.
func (c *Cursor) nsearch(key []byte) ***REMOVED***
	e := &c.stack[len(c.stack)-1]
	p, n := e.page, e.node

	// If we have a node then search its inodes.
	if n != nil ***REMOVED***
		index := sort.Search(len(n.inodes), func(i int) bool ***REMOVED***
			return bytes.Compare(n.inodes[i].key, key) != -1
		***REMOVED***)
		e.index = index
		return
	***REMOVED***

	// If we have a page then search its leaf elements.
	inodes := p.leafPageElements()
	index := sort.Search(int(p.count), func(i int) bool ***REMOVED***
		return bytes.Compare(inodes[i].key(), key) != -1
	***REMOVED***)
	e.index = index
***REMOVED***

// keyValue returns the key and value of the current leaf element.
func (c *Cursor) keyValue() ([]byte, []byte, uint32) ***REMOVED***
	ref := &c.stack[len(c.stack)-1]
	if ref.count() == 0 || ref.index >= ref.count() ***REMOVED***
		return nil, nil, 0
	***REMOVED***

	// Retrieve value from node.
	if ref.node != nil ***REMOVED***
		inode := &ref.node.inodes[ref.index]
		return inode.key, inode.value, inode.flags
	***REMOVED***

	// Or retrieve value from page.
	elem := ref.page.leafPageElement(uint16(ref.index))
	return elem.key(), elem.value(), elem.flags
***REMOVED***

// node returns the node that the cursor is currently positioned on.
func (c *Cursor) node() *node ***REMOVED***
	_assert(len(c.stack) > 0, "accessing a node with a zero-length cursor stack")

	// If the top of the stack is a leaf node then just return it.
	if ref := &c.stack[len(c.stack)-1]; ref.node != nil && ref.isLeaf() ***REMOVED***
		return ref.node
	***REMOVED***

	// Start from root and traverse down the hierarchy.
	var n = c.stack[0].node
	if n == nil ***REMOVED***
		n = c.bucket.node(c.stack[0].page.id, nil)
	***REMOVED***
	for _, ref := range c.stack[:len(c.stack)-1] ***REMOVED***
		_assert(!n.isLeaf, "expected branch node")
		n = n.childAt(int(ref.index))
	***REMOVED***
	_assert(n.isLeaf, "expected leaf node")
	return n
***REMOVED***

// elemRef represents a reference to an element on a given page/node.
type elemRef struct ***REMOVED***
	page  *page
	node  *node
	index int
***REMOVED***

// isLeaf returns whether the ref is pointing at a leaf page/node.
func (r *elemRef) isLeaf() bool ***REMOVED***
	if r.node != nil ***REMOVED***
		return r.node.isLeaf
	***REMOVED***
	return (r.page.flags & leafPageFlag) != 0
***REMOVED***

// count returns the number of inodes or page elements.
func (r *elemRef) count() int ***REMOVED***
	if r.node != nil ***REMOVED***
		return len(r.node.inodes)
	***REMOVED***
	return int(r.page.count)
***REMOVED***
