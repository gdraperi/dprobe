package bolt

import (
	"bytes"
	"fmt"
	"sort"
	"unsafe"
)

// node represents an in-memory, deserialized page.
type node struct ***REMOVED***
	bucket     *Bucket
	isLeaf     bool
	unbalanced bool
	spilled    bool
	key        []byte
	pgid       pgid
	parent     *node
	children   nodes
	inodes     inodes
***REMOVED***

// root returns the top-level node this node is attached to.
func (n *node) root() *node ***REMOVED***
	if n.parent == nil ***REMOVED***
		return n
	***REMOVED***
	return n.parent.root()
***REMOVED***

// minKeys returns the minimum number of inodes this node should have.
func (n *node) minKeys() int ***REMOVED***
	if n.isLeaf ***REMOVED***
		return 1
	***REMOVED***
	return 2
***REMOVED***

// size returns the size of the node after serialization.
func (n *node) size() int ***REMOVED***
	sz, elsz := pageHeaderSize, n.pageElementSize()
	for i := 0; i < len(n.inodes); i++ ***REMOVED***
		item := &n.inodes[i]
		sz += elsz + len(item.key) + len(item.value)
	***REMOVED***
	return sz
***REMOVED***

// sizeLessThan returns true if the node is less than a given size.
// This is an optimization to avoid calculating a large node when we only need
// to know if it fits inside a certain page size.
func (n *node) sizeLessThan(v int) bool ***REMOVED***
	sz, elsz := pageHeaderSize, n.pageElementSize()
	for i := 0; i < len(n.inodes); i++ ***REMOVED***
		item := &n.inodes[i]
		sz += elsz + len(item.key) + len(item.value)
		if sz >= v ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// pageElementSize returns the size of each page element based on the type of node.
func (n *node) pageElementSize() int ***REMOVED***
	if n.isLeaf ***REMOVED***
		return leafPageElementSize
	***REMOVED***
	return branchPageElementSize
***REMOVED***

// childAt returns the child node at a given index.
func (n *node) childAt(index int) *node ***REMOVED***
	if n.isLeaf ***REMOVED***
		panic(fmt.Sprintf("invalid childAt(%d) on a leaf node", index))
	***REMOVED***
	return n.bucket.node(n.inodes[index].pgid, n)
***REMOVED***

// childIndex returns the index of a given child node.
func (n *node) childIndex(child *node) int ***REMOVED***
	index := sort.Search(len(n.inodes), func(i int) bool ***REMOVED*** return bytes.Compare(n.inodes[i].key, child.key) != -1 ***REMOVED***)
	return index
***REMOVED***

// numChildren returns the number of children.
func (n *node) numChildren() int ***REMOVED***
	return len(n.inodes)
***REMOVED***

// nextSibling returns the next node with the same parent.
func (n *node) nextSibling() *node ***REMOVED***
	if n.parent == nil ***REMOVED***
		return nil
	***REMOVED***
	index := n.parent.childIndex(n)
	if index >= n.parent.numChildren()-1 ***REMOVED***
		return nil
	***REMOVED***
	return n.parent.childAt(index + 1)
***REMOVED***

// prevSibling returns the previous node with the same parent.
func (n *node) prevSibling() *node ***REMOVED***
	if n.parent == nil ***REMOVED***
		return nil
	***REMOVED***
	index := n.parent.childIndex(n)
	if index == 0 ***REMOVED***
		return nil
	***REMOVED***
	return n.parent.childAt(index - 1)
***REMOVED***

// put inserts a key/value.
func (n *node) put(oldKey, newKey, value []byte, pgid pgid, flags uint32) ***REMOVED***
	if pgid >= n.bucket.tx.meta.pgid ***REMOVED***
		panic(fmt.Sprintf("pgid (%d) above high water mark (%d)", pgid, n.bucket.tx.meta.pgid))
	***REMOVED*** else if len(oldKey) <= 0 ***REMOVED***
		panic("put: zero-length old key")
	***REMOVED*** else if len(newKey) <= 0 ***REMOVED***
		panic("put: zero-length new key")
	***REMOVED***

	// Find insertion index.
	index := sort.Search(len(n.inodes), func(i int) bool ***REMOVED*** return bytes.Compare(n.inodes[i].key, oldKey) != -1 ***REMOVED***)

	// Add capacity and shift nodes if we don't have an exact match and need to insert.
	exact := (len(n.inodes) > 0 && index < len(n.inodes) && bytes.Equal(n.inodes[index].key, oldKey))
	if !exact ***REMOVED***
		n.inodes = append(n.inodes, inode***REMOVED******REMOVED***)
		copy(n.inodes[index+1:], n.inodes[index:])
	***REMOVED***

	inode := &n.inodes[index]
	inode.flags = flags
	inode.key = newKey
	inode.value = value
	inode.pgid = pgid
	_assert(len(inode.key) > 0, "put: zero-length inode key")
***REMOVED***

// del removes a key from the node.
func (n *node) del(key []byte) ***REMOVED***
	// Find index of key.
	index := sort.Search(len(n.inodes), func(i int) bool ***REMOVED*** return bytes.Compare(n.inodes[i].key, key) != -1 ***REMOVED***)

	// Exit if the key isn't found.
	if index >= len(n.inodes) || !bytes.Equal(n.inodes[index].key, key) ***REMOVED***
		return
	***REMOVED***

	// Delete inode from the node.
	n.inodes = append(n.inodes[:index], n.inodes[index+1:]...)

	// Mark the node as needing rebalancing.
	n.unbalanced = true
***REMOVED***

// read initializes the node from a page.
func (n *node) read(p *page) ***REMOVED***
	n.pgid = p.id
	n.isLeaf = ((p.flags & leafPageFlag) != 0)
	n.inodes = make(inodes, int(p.count))

	for i := 0; i < int(p.count); i++ ***REMOVED***
		inode := &n.inodes[i]
		if n.isLeaf ***REMOVED***
			elem := p.leafPageElement(uint16(i))
			inode.flags = elem.flags
			inode.key = elem.key()
			inode.value = elem.value()
		***REMOVED*** else ***REMOVED***
			elem := p.branchPageElement(uint16(i))
			inode.pgid = elem.pgid
			inode.key = elem.key()
		***REMOVED***
		_assert(len(inode.key) > 0, "read: zero-length inode key")
	***REMOVED***

	// Save first key so we can find the node in the parent when we spill.
	if len(n.inodes) > 0 ***REMOVED***
		n.key = n.inodes[0].key
		_assert(len(n.key) > 0, "read: zero-length node key")
	***REMOVED*** else ***REMOVED***
		n.key = nil
	***REMOVED***
***REMOVED***

// write writes the items onto one or more pages.
func (n *node) write(p *page) ***REMOVED***
	// Initialize page.
	if n.isLeaf ***REMOVED***
		p.flags |= leafPageFlag
	***REMOVED*** else ***REMOVED***
		p.flags |= branchPageFlag
	***REMOVED***

	if len(n.inodes) >= 0xFFFF ***REMOVED***
		panic(fmt.Sprintf("inode overflow: %d (pgid=%d)", len(n.inodes), p.id))
	***REMOVED***
	p.count = uint16(len(n.inodes))

	// Stop here if there are no items to write.
	if p.count == 0 ***REMOVED***
		return
	***REMOVED***

	// Loop over each item and write it to the page.
	b := (*[maxAllocSize]byte)(unsafe.Pointer(&p.ptr))[n.pageElementSize()*len(n.inodes):]
	for i, item := range n.inodes ***REMOVED***
		_assert(len(item.key) > 0, "write: zero-length inode key")

		// Write the page element.
		if n.isLeaf ***REMOVED***
			elem := p.leafPageElement(uint16(i))
			elem.pos = uint32(uintptr(unsafe.Pointer(&b[0])) - uintptr(unsafe.Pointer(elem)))
			elem.flags = item.flags
			elem.ksize = uint32(len(item.key))
			elem.vsize = uint32(len(item.value))
		***REMOVED*** else ***REMOVED***
			elem := p.branchPageElement(uint16(i))
			elem.pos = uint32(uintptr(unsafe.Pointer(&b[0])) - uintptr(unsafe.Pointer(elem)))
			elem.ksize = uint32(len(item.key))
			elem.pgid = item.pgid
			_assert(elem.pgid != p.id, "write: circular dependency occurred")
		***REMOVED***

		// If the length of key+value is larger than the max allocation size
		// then we need to reallocate the byte array pointer.
		//
		// See: https://github.com/boltdb/bolt/pull/335
		klen, vlen := len(item.key), len(item.value)
		if len(b) < klen+vlen ***REMOVED***
			b = (*[maxAllocSize]byte)(unsafe.Pointer(&b[0]))[:]
		***REMOVED***

		// Write data for the element to the end of the page.
		copy(b[0:], item.key)
		b = b[klen:]
		copy(b[0:], item.value)
		b = b[vlen:]
	***REMOVED***

	// DEBUG ONLY: n.dump()
***REMOVED***

// split breaks up a node into multiple smaller nodes, if appropriate.
// This should only be called from the spill() function.
func (n *node) split(pageSize int) []*node ***REMOVED***
	var nodes []*node

	node := n
	for ***REMOVED***
		// Split node into two.
		a, b := node.splitTwo(pageSize)
		nodes = append(nodes, a)

		// If we can't split then exit the loop.
		if b == nil ***REMOVED***
			break
		***REMOVED***

		// Set node to b so it gets split on the next iteration.
		node = b
	***REMOVED***

	return nodes
***REMOVED***

// splitTwo breaks up a node into two smaller nodes, if appropriate.
// This should only be called from the split() function.
func (n *node) splitTwo(pageSize int) (*node, *node) ***REMOVED***
	// Ignore the split if the page doesn't have at least enough nodes for
	// two pages or if the nodes can fit in a single page.
	if len(n.inodes) <= (minKeysPerPage*2) || n.sizeLessThan(pageSize) ***REMOVED***
		return n, nil
	***REMOVED***

	// Determine the threshold before starting a new node.
	var fillPercent = n.bucket.FillPercent
	if fillPercent < minFillPercent ***REMOVED***
		fillPercent = minFillPercent
	***REMOVED*** else if fillPercent > maxFillPercent ***REMOVED***
		fillPercent = maxFillPercent
	***REMOVED***
	threshold := int(float64(pageSize) * fillPercent)

	// Determine split position and sizes of the two pages.
	splitIndex, _ := n.splitIndex(threshold)

	// Split node into two separate nodes.
	// If there's no parent then we'll need to create one.
	if n.parent == nil ***REMOVED***
		n.parent = &node***REMOVED***bucket: n.bucket, children: []*node***REMOVED***n***REMOVED******REMOVED***
	***REMOVED***

	// Create a new node and add it to the parent.
	next := &node***REMOVED***bucket: n.bucket, isLeaf: n.isLeaf, parent: n.parent***REMOVED***
	n.parent.children = append(n.parent.children, next)

	// Split inodes across two nodes.
	next.inodes = n.inodes[splitIndex:]
	n.inodes = n.inodes[:splitIndex]

	// Update the statistics.
	n.bucket.tx.stats.Split++

	return n, next
***REMOVED***

// splitIndex finds the position where a page will fill a given threshold.
// It returns the index as well as the size of the first page.
// This is only be called from split().
func (n *node) splitIndex(threshold int) (index, sz int) ***REMOVED***
	sz = pageHeaderSize

	// Loop until we only have the minimum number of keys required for the second page.
	for i := 0; i < len(n.inodes)-minKeysPerPage; i++ ***REMOVED***
		index = i
		inode := n.inodes[i]
		elsize := n.pageElementSize() + len(inode.key) + len(inode.value)

		// If we have at least the minimum number of keys and adding another
		// node would put us over the threshold then exit and return.
		if i >= minKeysPerPage && sz+elsize > threshold ***REMOVED***
			break
		***REMOVED***

		// Add the element size to the total size.
		sz += elsize
	***REMOVED***

	return
***REMOVED***

// spill writes the nodes to dirty pages and splits nodes as it goes.
// Returns an error if dirty pages cannot be allocated.
func (n *node) spill() error ***REMOVED***
	var tx = n.bucket.tx
	if n.spilled ***REMOVED***
		return nil
	***REMOVED***

	// Spill child nodes first. Child nodes can materialize sibling nodes in
	// the case of split-merge so we cannot use a range loop. We have to check
	// the children size on every loop iteration.
	sort.Sort(n.children)
	for i := 0; i < len(n.children); i++ ***REMOVED***
		if err := n.children[i].spill(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// We no longer need the child list because it's only used for spill tracking.
	n.children = nil

	// Split nodes into appropriate sizes. The first node will always be n.
	var nodes = n.split(tx.db.pageSize)
	for _, node := range nodes ***REMOVED***
		// Add node's page to the freelist if it's not new.
		if node.pgid > 0 ***REMOVED***
			tx.db.freelist.free(tx.meta.txid, tx.page(node.pgid))
			node.pgid = 0
		***REMOVED***

		// Allocate contiguous space for the node.
		p, err := tx.allocate((node.size() / tx.db.pageSize) + 1)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Write the node.
		if p.id >= tx.meta.pgid ***REMOVED***
			panic(fmt.Sprintf("pgid (%d) above high water mark (%d)", p.id, tx.meta.pgid))
		***REMOVED***
		node.pgid = p.id
		node.write(p)
		node.spilled = true

		// Insert into parent inodes.
		if node.parent != nil ***REMOVED***
			var key = node.key
			if key == nil ***REMOVED***
				key = node.inodes[0].key
			***REMOVED***

			node.parent.put(key, node.inodes[0].key, nil, node.pgid, 0)
			node.key = node.inodes[0].key
			_assert(len(node.key) > 0, "spill: zero-length node key")
		***REMOVED***

		// Update the statistics.
		tx.stats.Spill++
	***REMOVED***

	// If the root node split and created a new root then we need to spill that
	// as well. We'll clear out the children to make sure it doesn't try to respill.
	if n.parent != nil && n.parent.pgid == 0 ***REMOVED***
		n.children = nil
		return n.parent.spill()
	***REMOVED***

	return nil
***REMOVED***

// rebalance attempts to combine the node with sibling nodes if the node fill
// size is below a threshold or if there are not enough keys.
func (n *node) rebalance() ***REMOVED***
	if !n.unbalanced ***REMOVED***
		return
	***REMOVED***
	n.unbalanced = false

	// Update statistics.
	n.bucket.tx.stats.Rebalance++

	// Ignore if node is above threshold (25%) and has enough keys.
	var threshold = n.bucket.tx.db.pageSize / 4
	if n.size() > threshold && len(n.inodes) > n.minKeys() ***REMOVED***
		return
	***REMOVED***

	// Root node has special handling.
	if n.parent == nil ***REMOVED***
		// If root node is a branch and only has one node then collapse it.
		if !n.isLeaf && len(n.inodes) == 1 ***REMOVED***
			// Move root's child up.
			child := n.bucket.node(n.inodes[0].pgid, n)
			n.isLeaf = child.isLeaf
			n.inodes = child.inodes[:]
			n.children = child.children

			// Reparent all child nodes being moved.
			for _, inode := range n.inodes ***REMOVED***
				if child, ok := n.bucket.nodes[inode.pgid]; ok ***REMOVED***
					child.parent = n
				***REMOVED***
			***REMOVED***

			// Remove old child.
			child.parent = nil
			delete(n.bucket.nodes, child.pgid)
			child.free()
		***REMOVED***

		return
	***REMOVED***

	// If node has no keys then just remove it.
	if n.numChildren() == 0 ***REMOVED***
		n.parent.del(n.key)
		n.parent.removeChild(n)
		delete(n.bucket.nodes, n.pgid)
		n.free()
		n.parent.rebalance()
		return
	***REMOVED***

	_assert(n.parent.numChildren() > 1, "parent must have at least 2 children")

	// Destination node is right sibling if idx == 0, otherwise left sibling.
	var target *node
	var useNextSibling = (n.parent.childIndex(n) == 0)
	if useNextSibling ***REMOVED***
		target = n.nextSibling()
	***REMOVED*** else ***REMOVED***
		target = n.prevSibling()
	***REMOVED***

	// If both this node and the target node are too small then merge them.
	if useNextSibling ***REMOVED***
		// Reparent all child nodes being moved.
		for _, inode := range target.inodes ***REMOVED***
			if child, ok := n.bucket.nodes[inode.pgid]; ok ***REMOVED***
				child.parent.removeChild(child)
				child.parent = n
				child.parent.children = append(child.parent.children, child)
			***REMOVED***
		***REMOVED***

		// Copy over inodes from target and remove target.
		n.inodes = append(n.inodes, target.inodes...)
		n.parent.del(target.key)
		n.parent.removeChild(target)
		delete(n.bucket.nodes, target.pgid)
		target.free()
	***REMOVED*** else ***REMOVED***
		// Reparent all child nodes being moved.
		for _, inode := range n.inodes ***REMOVED***
			if child, ok := n.bucket.nodes[inode.pgid]; ok ***REMOVED***
				child.parent.removeChild(child)
				child.parent = target
				child.parent.children = append(child.parent.children, child)
			***REMOVED***
		***REMOVED***

		// Copy over inodes to target and remove node.
		target.inodes = append(target.inodes, n.inodes...)
		n.parent.del(n.key)
		n.parent.removeChild(n)
		delete(n.bucket.nodes, n.pgid)
		n.free()
	***REMOVED***

	// Either this node or the target node was deleted from the parent so rebalance it.
	n.parent.rebalance()
***REMOVED***

// removes a node from the list of in-memory children.
// This does not affect the inodes.
func (n *node) removeChild(target *node) ***REMOVED***
	for i, child := range n.children ***REMOVED***
		if child == target ***REMOVED***
			n.children = append(n.children[:i], n.children[i+1:]...)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// dereference causes the node to copy all its inode key/value references to heap memory.
// This is required when the mmap is reallocated so inodes are not pointing to stale data.
func (n *node) dereference() ***REMOVED***
	if n.key != nil ***REMOVED***
		key := make([]byte, len(n.key))
		copy(key, n.key)
		n.key = key
		_assert(n.pgid == 0 || len(n.key) > 0, "dereference: zero-length node key on existing node")
	***REMOVED***

	for i := range n.inodes ***REMOVED***
		inode := &n.inodes[i]

		key := make([]byte, len(inode.key))
		copy(key, inode.key)
		inode.key = key
		_assert(len(inode.key) > 0, "dereference: zero-length inode key")

		value := make([]byte, len(inode.value))
		copy(value, inode.value)
		inode.value = value
	***REMOVED***

	// Recursively dereference children.
	for _, child := range n.children ***REMOVED***
		child.dereference()
	***REMOVED***

	// Update statistics.
	n.bucket.tx.stats.NodeDeref++
***REMOVED***

// free adds the node's underlying page to the freelist.
func (n *node) free() ***REMOVED***
	if n.pgid != 0 ***REMOVED***
		n.bucket.tx.db.freelist.free(n.bucket.tx.meta.txid, n.bucket.tx.page(n.pgid))
		n.pgid = 0
	***REMOVED***
***REMOVED***

// dump writes the contents of the node to STDERR for debugging purposes.
/*
func (n *node) dump() ***REMOVED***
	// Write node header.
	var typ = "branch"
	if n.isLeaf ***REMOVED***
		typ = "leaf"
	***REMOVED***
	warnf("[NODE %d ***REMOVED***type=%s count=%d***REMOVED***]", n.pgid, typ, len(n.inodes))

	// Write out abbreviated version of each item.
	for _, item := range n.inodes ***REMOVED***
		if n.isLeaf ***REMOVED***
			if item.flags&bucketLeafFlag != 0 ***REMOVED***
				bucket := (*bucket)(unsafe.Pointer(&item.value[0]))
				warnf("+L %08x -> (bucket root=%d)", trunc(item.key, 4), bucket.root)
			***REMOVED*** else ***REMOVED***
				warnf("+L %08x -> %08x", trunc(item.key, 4), trunc(item.value, 4))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			warnf("+B %08x -> pgid=%d", trunc(item.key, 4), item.pgid)
		***REMOVED***
	***REMOVED***
	warn("")
***REMOVED***
*/

type nodes []*node

func (s nodes) Len() int           ***REMOVED*** return len(s) ***REMOVED***
func (s nodes) Swap(i, j int)      ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***
func (s nodes) Less(i, j int) bool ***REMOVED*** return bytes.Compare(s[i].inodes[0].key, s[j].inodes[0].key) == -1 ***REMOVED***

// inode represents an internal node inside of a node.
// It can be used to point to elements in a page or point
// to an element which hasn't been added to a page yet.
type inode struct ***REMOVED***
	flags uint32
	pgid  pgid
	key   []byte
	value []byte
***REMOVED***

type inodes []inode
