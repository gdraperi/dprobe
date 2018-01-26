package bolt

import (
	"bytes"
	"fmt"
	"unsafe"
)

const (
	// MaxKeySize is the maximum length of a key, in bytes.
	MaxKeySize = 32768

	// MaxValueSize is the maximum length of a value, in bytes.
	MaxValueSize = (1 << 31) - 2
)

const (
	maxUint = ^uint(0)
	minUint = 0
	maxInt  = int(^uint(0) >> 1)
	minInt  = -maxInt - 1
)

const bucketHeaderSize = int(unsafe.Sizeof(bucket***REMOVED******REMOVED***))

const (
	minFillPercent = 0.1
	maxFillPercent = 1.0
)

// DefaultFillPercent is the percentage that split pages are filled.
// This value can be changed by setting Bucket.FillPercent.
const DefaultFillPercent = 0.5

// Bucket represents a collection of key/value pairs inside the database.
type Bucket struct ***REMOVED***
	*bucket
	tx       *Tx                // the associated transaction
	buckets  map[string]*Bucket // subbucket cache
	page     *page              // inline page reference
	rootNode *node              // materialized node for the root page.
	nodes    map[pgid]*node     // node cache

	// Sets the threshold for filling nodes when they split. By default,
	// the bucket will fill to 50% but it can be useful to increase this
	// amount if you know that your write workloads are mostly append-only.
	//
	// This is non-persisted across transactions so it must be set in every Tx.
	FillPercent float64
***REMOVED***

// bucket represents the on-file representation of a bucket.
// This is stored as the "value" of a bucket key. If the bucket is small enough,
// then its root page can be stored inline in the "value", after the bucket
// header. In the case of inline buckets, the "root" will be 0.
type bucket struct ***REMOVED***
	root     pgid   // page id of the bucket's root-level page
	sequence uint64 // monotonically incrementing, used by NextSequence()
***REMOVED***

// newBucket returns a new bucket associated with a transaction.
func newBucket(tx *Tx) Bucket ***REMOVED***
	var b = Bucket***REMOVED***tx: tx, FillPercent: DefaultFillPercent***REMOVED***
	if tx.writable ***REMOVED***
		b.buckets = make(map[string]*Bucket)
		b.nodes = make(map[pgid]*node)
	***REMOVED***
	return b
***REMOVED***

// Tx returns the tx of the bucket.
func (b *Bucket) Tx() *Tx ***REMOVED***
	return b.tx
***REMOVED***

// Root returns the root of the bucket.
func (b *Bucket) Root() pgid ***REMOVED***
	return b.root
***REMOVED***

// Writable returns whether the bucket is writable.
func (b *Bucket) Writable() bool ***REMOVED***
	return b.tx.writable
***REMOVED***

// Cursor creates a cursor associated with the bucket.
// The cursor is only valid as long as the transaction is open.
// Do not use a cursor after the transaction is closed.
func (b *Bucket) Cursor() *Cursor ***REMOVED***
	// Update transaction statistics.
	b.tx.stats.CursorCount++

	// Allocate and return a cursor.
	return &Cursor***REMOVED***
		bucket: b,
		stack:  make([]elemRef, 0),
	***REMOVED***
***REMOVED***

// Bucket retrieves a nested bucket by name.
// Returns nil if the bucket does not exist.
// The bucket instance is only valid for the lifetime of the transaction.
func (b *Bucket) Bucket(name []byte) *Bucket ***REMOVED***
	if b.buckets != nil ***REMOVED***
		if child := b.buckets[string(name)]; child != nil ***REMOVED***
			return child
		***REMOVED***
	***REMOVED***

	// Move cursor to key.
	c := b.Cursor()
	k, v, flags := c.seek(name)

	// Return nil if the key doesn't exist or it is not a bucket.
	if !bytes.Equal(name, k) || (flags&bucketLeafFlag) == 0 ***REMOVED***
		return nil
	***REMOVED***

	// Otherwise create a bucket and cache it.
	var child = b.openBucket(v)
	if b.buckets != nil ***REMOVED***
		b.buckets[string(name)] = child
	***REMOVED***

	return child
***REMOVED***

// Helper method that re-interprets a sub-bucket value
// from a parent into a Bucket
func (b *Bucket) openBucket(value []byte) *Bucket ***REMOVED***
	var child = newBucket(b.tx)

	// If unaligned load/stores are broken on this arch and value is
	// unaligned simply clone to an aligned byte array.
	unaligned := brokenUnaligned && uintptr(unsafe.Pointer(&value[0]))&3 != 0

	if unaligned ***REMOVED***
		value = cloneBytes(value)
	***REMOVED***

	// If this is a writable transaction then we need to copy the bucket entry.
	// Read-only transactions can point directly at the mmap entry.
	if b.tx.writable && !unaligned ***REMOVED***
		child.bucket = &bucket***REMOVED******REMOVED***
		*child.bucket = *(*bucket)(unsafe.Pointer(&value[0]))
	***REMOVED*** else ***REMOVED***
		child.bucket = (*bucket)(unsafe.Pointer(&value[0]))
	***REMOVED***

	// Save a reference to the inline page if the bucket is inline.
	if child.root == 0 ***REMOVED***
		child.page = (*page)(unsafe.Pointer(&value[bucketHeaderSize]))
	***REMOVED***

	return &child
***REMOVED***

// CreateBucket creates a new bucket at the given key and returns the new bucket.
// Returns an error if the key already exists, if the bucket name is blank, or if the bucket name is too long.
// The bucket instance is only valid for the lifetime of the transaction.
func (b *Bucket) CreateBucket(key []byte) (*Bucket, error) ***REMOVED***
	if b.tx.db == nil ***REMOVED***
		return nil, ErrTxClosed
	***REMOVED*** else if !b.tx.writable ***REMOVED***
		return nil, ErrTxNotWritable
	***REMOVED*** else if len(key) == 0 ***REMOVED***
		return nil, ErrBucketNameRequired
	***REMOVED***

	// Move cursor to correct position.
	c := b.Cursor()
	k, _, flags := c.seek(key)

	// Return an error if there is an existing key.
	if bytes.Equal(key, k) ***REMOVED***
		if (flags & bucketLeafFlag) != 0 ***REMOVED***
			return nil, ErrBucketExists
		***REMOVED*** else ***REMOVED***
			return nil, ErrIncompatibleValue
		***REMOVED***
	***REMOVED***

	// Create empty, inline bucket.
	var bucket = Bucket***REMOVED***
		bucket:      &bucket***REMOVED******REMOVED***,
		rootNode:    &node***REMOVED***isLeaf: true***REMOVED***,
		FillPercent: DefaultFillPercent,
	***REMOVED***
	var value = bucket.write()

	// Insert into node.
	key = cloneBytes(key)
	c.node().put(key, key, value, 0, bucketLeafFlag)

	// Since subbuckets are not allowed on inline buckets, we need to
	// dereference the inline page, if it exists. This will cause the bucket
	// to be treated as a regular, non-inline bucket for the rest of the tx.
	b.page = nil

	return b.Bucket(key), nil
***REMOVED***

// CreateBucketIfNotExists creates a new bucket if it doesn't already exist and returns a reference to it.
// Returns an error if the bucket name is blank, or if the bucket name is too long.
// The bucket instance is only valid for the lifetime of the transaction.
func (b *Bucket) CreateBucketIfNotExists(key []byte) (*Bucket, error) ***REMOVED***
	child, err := b.CreateBucket(key)
	if err == ErrBucketExists ***REMOVED***
		return b.Bucket(key), nil
	***REMOVED*** else if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return child, nil
***REMOVED***

// DeleteBucket deletes a bucket at the given key.
// Returns an error if the bucket does not exists, or if the key represents a non-bucket value.
func (b *Bucket) DeleteBucket(key []byte) error ***REMOVED***
	if b.tx.db == nil ***REMOVED***
		return ErrTxClosed
	***REMOVED*** else if !b.Writable() ***REMOVED***
		return ErrTxNotWritable
	***REMOVED***

	// Move cursor to correct position.
	c := b.Cursor()
	k, _, flags := c.seek(key)

	// Return an error if bucket doesn't exist or is not a bucket.
	if !bytes.Equal(key, k) ***REMOVED***
		return ErrBucketNotFound
	***REMOVED*** else if (flags & bucketLeafFlag) == 0 ***REMOVED***
		return ErrIncompatibleValue
	***REMOVED***

	// Recursively delete all child buckets.
	child := b.Bucket(key)
	err := child.ForEach(func(k, v []byte) error ***REMOVED***
		if v == nil ***REMOVED***
			if err := child.DeleteBucket(k); err != nil ***REMOVED***
				return fmt.Errorf("delete bucket: %s", err)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Remove cached copy.
	delete(b.buckets, string(key))

	// Release all bucket pages to freelist.
	child.nodes = nil
	child.rootNode = nil
	child.free()

	// Delete the node if we have a matching key.
	c.node().del(key)

	return nil
***REMOVED***

// Get retrieves the value for a key in the bucket.
// Returns a nil value if the key does not exist or if the key is a nested bucket.
// The returned value is only valid for the life of the transaction.
func (b *Bucket) Get(key []byte) []byte ***REMOVED***
	k, v, flags := b.Cursor().seek(key)

	// Return nil if this is a bucket.
	if (flags & bucketLeafFlag) != 0 ***REMOVED***
		return nil
	***REMOVED***

	// If our target node isn't the same key as what's passed in then return nil.
	if !bytes.Equal(key, k) ***REMOVED***
		return nil
	***REMOVED***
	return v
***REMOVED***

// Put sets the value for a key in the bucket.
// If the key exist then its previous value will be overwritten.
// Supplied value must remain valid for the life of the transaction.
// Returns an error if the bucket was created from a read-only transaction, if the key is blank, if the key is too large, or if the value is too large.
func (b *Bucket) Put(key []byte, value []byte) error ***REMOVED***
	if b.tx.db == nil ***REMOVED***
		return ErrTxClosed
	***REMOVED*** else if !b.Writable() ***REMOVED***
		return ErrTxNotWritable
	***REMOVED*** else if len(key) == 0 ***REMOVED***
		return ErrKeyRequired
	***REMOVED*** else if len(key) > MaxKeySize ***REMOVED***
		return ErrKeyTooLarge
	***REMOVED*** else if int64(len(value)) > MaxValueSize ***REMOVED***
		return ErrValueTooLarge
	***REMOVED***

	// Move cursor to correct position.
	c := b.Cursor()
	k, _, flags := c.seek(key)

	// Return an error if there is an existing key with a bucket value.
	if bytes.Equal(key, k) && (flags&bucketLeafFlag) != 0 ***REMOVED***
		return ErrIncompatibleValue
	***REMOVED***

	// Insert into node.
	key = cloneBytes(key)
	c.node().put(key, key, value, 0, 0)

	return nil
***REMOVED***

// Delete removes a key from the bucket.
// If the key does not exist then nothing is done and a nil error is returned.
// Returns an error if the bucket was created from a read-only transaction.
func (b *Bucket) Delete(key []byte) error ***REMOVED***
	if b.tx.db == nil ***REMOVED***
		return ErrTxClosed
	***REMOVED*** else if !b.Writable() ***REMOVED***
		return ErrTxNotWritable
	***REMOVED***

	// Move cursor to correct position.
	c := b.Cursor()
	_, _, flags := c.seek(key)

	// Return an error if there is already existing bucket value.
	if (flags & bucketLeafFlag) != 0 ***REMOVED***
		return ErrIncompatibleValue
	***REMOVED***

	// Delete the node if we have a matching key.
	c.node().del(key)

	return nil
***REMOVED***

// Sequence returns the current integer for the bucket without incrementing it.
func (b *Bucket) Sequence() uint64 ***REMOVED*** return b.bucket.sequence ***REMOVED***

// SetSequence updates the sequence number for the bucket.
func (b *Bucket) SetSequence(v uint64) error ***REMOVED***
	if b.tx.db == nil ***REMOVED***
		return ErrTxClosed
	***REMOVED*** else if !b.Writable() ***REMOVED***
		return ErrTxNotWritable
	***REMOVED***

	// Materialize the root node if it hasn't been already so that the
	// bucket will be saved during commit.
	if b.rootNode == nil ***REMOVED***
		_ = b.node(b.root, nil)
	***REMOVED***

	// Increment and return the sequence.
	b.bucket.sequence = v
	return nil
***REMOVED***

// NextSequence returns an autoincrementing integer for the bucket.
func (b *Bucket) NextSequence() (uint64, error) ***REMOVED***
	if b.tx.db == nil ***REMOVED***
		return 0, ErrTxClosed
	***REMOVED*** else if !b.Writable() ***REMOVED***
		return 0, ErrTxNotWritable
	***REMOVED***

	// Materialize the root node if it hasn't been already so that the
	// bucket will be saved during commit.
	if b.rootNode == nil ***REMOVED***
		_ = b.node(b.root, nil)
	***REMOVED***

	// Increment and return the sequence.
	b.bucket.sequence++
	return b.bucket.sequence, nil
***REMOVED***

// ForEach executes a function for each key/value pair in a bucket.
// If the provided function returns an error then the iteration is stopped and
// the error is returned to the caller. The provided function must not modify
// the bucket; this will result in undefined behavior.
func (b *Bucket) ForEach(fn func(k, v []byte) error) error ***REMOVED***
	if b.tx.db == nil ***REMOVED***
		return ErrTxClosed
	***REMOVED***
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() ***REMOVED***
		if err := fn(k, v); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Stat returns stats on a bucket.
func (b *Bucket) Stats() BucketStats ***REMOVED***
	var s, subStats BucketStats
	pageSize := b.tx.db.pageSize
	s.BucketN += 1
	if b.root == 0 ***REMOVED***
		s.InlineBucketN += 1
	***REMOVED***
	b.forEachPage(func(p *page, depth int) ***REMOVED***
		if (p.flags & leafPageFlag) != 0 ***REMOVED***
			s.KeyN += int(p.count)

			// used totals the used bytes for the page
			used := pageHeaderSize

			if p.count != 0 ***REMOVED***
				// If page has any elements, add all element headers.
				used += leafPageElementSize * int(p.count-1)

				// Add all element key, value sizes.
				// The computation takes advantage of the fact that the position
				// of the last element's key/value equals to the total of the sizes
				// of all previous elements' keys and values.
				// It also includes the last element's header.
				lastElement := p.leafPageElement(p.count - 1)
				used += int(lastElement.pos + lastElement.ksize + lastElement.vsize)
			***REMOVED***

			if b.root == 0 ***REMOVED***
				// For inlined bucket just update the inline stats
				s.InlineBucketInuse += used
			***REMOVED*** else ***REMOVED***
				// For non-inlined bucket update all the leaf stats
				s.LeafPageN++
				s.LeafInuse += used
				s.LeafOverflowN += int(p.overflow)

				// Collect stats from sub-buckets.
				// Do that by iterating over all element headers
				// looking for the ones with the bucketLeafFlag.
				for i := uint16(0); i < p.count; i++ ***REMOVED***
					e := p.leafPageElement(i)
					if (e.flags & bucketLeafFlag) != 0 ***REMOVED***
						// For any bucket element, open the element value
						// and recursively call Stats on the contained bucket.
						subStats.Add(b.openBucket(e.value()).Stats())
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if (p.flags & branchPageFlag) != 0 ***REMOVED***
			s.BranchPageN++
			lastElement := p.branchPageElement(p.count - 1)

			// used totals the used bytes for the page
			// Add header and all element headers.
			used := pageHeaderSize + (branchPageElementSize * int(p.count-1))

			// Add size of all keys and values.
			// Again, use the fact that last element's position equals to
			// the total of key, value sizes of all previous elements.
			used += int(lastElement.pos + lastElement.ksize)
			s.BranchInuse += used
			s.BranchOverflowN += int(p.overflow)
		***REMOVED***

		// Keep track of maximum page depth.
		if depth+1 > s.Depth ***REMOVED***
			s.Depth = (depth + 1)
		***REMOVED***
	***REMOVED***)

	// Alloc stats can be computed from page counts and pageSize.
	s.BranchAlloc = (s.BranchPageN + s.BranchOverflowN) * pageSize
	s.LeafAlloc = (s.LeafPageN + s.LeafOverflowN) * pageSize

	// Add the max depth of sub-buckets to get total nested depth.
	s.Depth += subStats.Depth
	// Add the stats for all sub-buckets
	s.Add(subStats)
	return s
***REMOVED***

// forEachPage iterates over every page in a bucket, including inline pages.
func (b *Bucket) forEachPage(fn func(*page, int)) ***REMOVED***
	// If we have an inline page then just use that.
	if b.page != nil ***REMOVED***
		fn(b.page, 0)
		return
	***REMOVED***

	// Otherwise traverse the page hierarchy.
	b.tx.forEachPage(b.root, 0, fn)
***REMOVED***

// forEachPageNode iterates over every page (or node) in a bucket.
// This also includes inline pages.
func (b *Bucket) forEachPageNode(fn func(*page, *node, int)) ***REMOVED***
	// If we have an inline page or root node then just use that.
	if b.page != nil ***REMOVED***
		fn(b.page, nil, 0)
		return
	***REMOVED***
	b._forEachPageNode(b.root, 0, fn)
***REMOVED***

func (b *Bucket) _forEachPageNode(pgid pgid, depth int, fn func(*page, *node, int)) ***REMOVED***
	var p, n = b.pageNode(pgid)

	// Execute function.
	fn(p, n, depth)

	// Recursively loop over children.
	if p != nil ***REMOVED***
		if (p.flags & branchPageFlag) != 0 ***REMOVED***
			for i := 0; i < int(p.count); i++ ***REMOVED***
				elem := p.branchPageElement(uint16(i))
				b._forEachPageNode(elem.pgid, depth+1, fn)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !n.isLeaf ***REMOVED***
			for _, inode := range n.inodes ***REMOVED***
				b._forEachPageNode(inode.pgid, depth+1, fn)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// spill writes all the nodes for this bucket to dirty pages.
func (b *Bucket) spill() error ***REMOVED***
	// Spill all child buckets first.
	for name, child := range b.buckets ***REMOVED***
		// If the child bucket is small enough and it has no child buckets then
		// write it inline into the parent bucket's page. Otherwise spill it
		// like a normal bucket and make the parent value a pointer to the page.
		var value []byte
		if child.inlineable() ***REMOVED***
			child.free()
			value = child.write()
		***REMOVED*** else ***REMOVED***
			if err := child.spill(); err != nil ***REMOVED***
				return err
			***REMOVED***

			// Update the child bucket header in this bucket.
			value = make([]byte, unsafe.Sizeof(bucket***REMOVED******REMOVED***))
			var bucket = (*bucket)(unsafe.Pointer(&value[0]))
			*bucket = *child.bucket
		***REMOVED***

		// Skip writing the bucket if there are no materialized nodes.
		if child.rootNode == nil ***REMOVED***
			continue
		***REMOVED***

		// Update parent node.
		var c = b.Cursor()
		k, _, flags := c.seek([]byte(name))
		if !bytes.Equal([]byte(name), k) ***REMOVED***
			panic(fmt.Sprintf("misplaced bucket header: %x -> %x", []byte(name), k))
		***REMOVED***
		if flags&bucketLeafFlag == 0 ***REMOVED***
			panic(fmt.Sprintf("unexpected bucket header flag: %x", flags))
		***REMOVED***
		c.node().put([]byte(name), []byte(name), value, 0, bucketLeafFlag)
	***REMOVED***

	// Ignore if there's not a materialized root node.
	if b.rootNode == nil ***REMOVED***
		return nil
	***REMOVED***

	// Spill nodes.
	if err := b.rootNode.spill(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.rootNode = b.rootNode.root()

	// Update the root node for this bucket.
	if b.rootNode.pgid >= b.tx.meta.pgid ***REMOVED***
		panic(fmt.Sprintf("pgid (%d) above high water mark (%d)", b.rootNode.pgid, b.tx.meta.pgid))
	***REMOVED***
	b.root = b.rootNode.pgid

	return nil
***REMOVED***

// inlineable returns true if a bucket is small enough to be written inline
// and if it contains no subbuckets. Otherwise returns false.
func (b *Bucket) inlineable() bool ***REMOVED***
	var n = b.rootNode

	// Bucket must only contain a single leaf node.
	if n == nil || !n.isLeaf ***REMOVED***
		return false
	***REMOVED***

	// Bucket is not inlineable if it contains subbuckets or if it goes beyond
	// our threshold for inline bucket size.
	var size = pageHeaderSize
	for _, inode := range n.inodes ***REMOVED***
		size += leafPageElementSize + len(inode.key) + len(inode.value)

		if inode.flags&bucketLeafFlag != 0 ***REMOVED***
			return false
		***REMOVED*** else if size > b.maxInlineBucketSize() ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// Returns the maximum total size of a bucket to make it a candidate for inlining.
func (b *Bucket) maxInlineBucketSize() int ***REMOVED***
	return b.tx.db.pageSize / 4
***REMOVED***

// write allocates and writes a bucket to a byte slice.
func (b *Bucket) write() []byte ***REMOVED***
	// Allocate the appropriate size.
	var n = b.rootNode
	var value = make([]byte, bucketHeaderSize+n.size())

	// Write a bucket header.
	var bucket = (*bucket)(unsafe.Pointer(&value[0]))
	*bucket = *b.bucket

	// Convert byte slice to a fake page and write the root node.
	var p = (*page)(unsafe.Pointer(&value[bucketHeaderSize]))
	n.write(p)

	return value
***REMOVED***

// rebalance attempts to balance all nodes.
func (b *Bucket) rebalance() ***REMOVED***
	for _, n := range b.nodes ***REMOVED***
		n.rebalance()
	***REMOVED***
	for _, child := range b.buckets ***REMOVED***
		child.rebalance()
	***REMOVED***
***REMOVED***

// node creates a node from a page and associates it with a given parent.
func (b *Bucket) node(pgid pgid, parent *node) *node ***REMOVED***
	_assert(b.nodes != nil, "nodes map expected")

	// Retrieve node if it's already been created.
	if n := b.nodes[pgid]; n != nil ***REMOVED***
		return n
	***REMOVED***

	// Otherwise create a node and cache it.
	n := &node***REMOVED***bucket: b, parent: parent***REMOVED***
	if parent == nil ***REMOVED***
		b.rootNode = n
	***REMOVED*** else ***REMOVED***
		parent.children = append(parent.children, n)
	***REMOVED***

	// Use the inline page if this is an inline bucket.
	var p = b.page
	if p == nil ***REMOVED***
		p = b.tx.page(pgid)
	***REMOVED***

	// Read the page into the node and cache it.
	n.read(p)
	b.nodes[pgid] = n

	// Update statistics.
	b.tx.stats.NodeCount++

	return n
***REMOVED***

// free recursively frees all pages in the bucket.
func (b *Bucket) free() ***REMOVED***
	if b.root == 0 ***REMOVED***
		return
	***REMOVED***

	var tx = b.tx
	b.forEachPageNode(func(p *page, n *node, _ int) ***REMOVED***
		if p != nil ***REMOVED***
			tx.db.freelist.free(tx.meta.txid, p)
		***REMOVED*** else ***REMOVED***
			n.free()
		***REMOVED***
	***REMOVED***)
	b.root = 0
***REMOVED***

// dereference removes all references to the old mmap.
func (b *Bucket) dereference() ***REMOVED***
	if b.rootNode != nil ***REMOVED***
		b.rootNode.root().dereference()
	***REMOVED***

	for _, child := range b.buckets ***REMOVED***
		child.dereference()
	***REMOVED***
***REMOVED***

// pageNode returns the in-memory node, if it exists.
// Otherwise returns the underlying page.
func (b *Bucket) pageNode(id pgid) (*page, *node) ***REMOVED***
	// Inline buckets have a fake page embedded in their value so treat them
	// differently. We'll return the rootNode (if available) or the fake page.
	if b.root == 0 ***REMOVED***
		if id != 0 ***REMOVED***
			panic(fmt.Sprintf("inline bucket non-zero page access(2): %d != 0", id))
		***REMOVED***
		if b.rootNode != nil ***REMOVED***
			return nil, b.rootNode
		***REMOVED***
		return b.page, nil
	***REMOVED***

	// Check the node cache for non-inline buckets.
	if b.nodes != nil ***REMOVED***
		if n := b.nodes[id]; n != nil ***REMOVED***
			return nil, n
		***REMOVED***
	***REMOVED***

	// Finally lookup the page from the transaction if no node is materialized.
	return b.tx.page(id), nil
***REMOVED***

// BucketStats records statistics about resources used by a bucket.
type BucketStats struct ***REMOVED***
	// Page count statistics.
	BranchPageN     int // number of logical branch pages
	BranchOverflowN int // number of physical branch overflow pages
	LeafPageN       int // number of logical leaf pages
	LeafOverflowN   int // number of physical leaf overflow pages

	// Tree statistics.
	KeyN  int // number of keys/value pairs
	Depth int // number of levels in B+tree

	// Page size utilization.
	BranchAlloc int // bytes allocated for physical branch pages
	BranchInuse int // bytes actually used for branch data
	LeafAlloc   int // bytes allocated for physical leaf pages
	LeafInuse   int // bytes actually used for leaf data

	// Bucket statistics
	BucketN           int // total number of buckets including the top bucket
	InlineBucketN     int // total number on inlined buckets
	InlineBucketInuse int // bytes used for inlined buckets (also accounted for in LeafInuse)
***REMOVED***

func (s *BucketStats) Add(other BucketStats) ***REMOVED***
	s.BranchPageN += other.BranchPageN
	s.BranchOverflowN += other.BranchOverflowN
	s.LeafPageN += other.LeafPageN
	s.LeafOverflowN += other.LeafOverflowN
	s.KeyN += other.KeyN
	if s.Depth < other.Depth ***REMOVED***
		s.Depth = other.Depth
	***REMOVED***
	s.BranchAlloc += other.BranchAlloc
	s.BranchInuse += other.BranchInuse
	s.LeafAlloc += other.LeafAlloc
	s.LeafInuse += other.LeafInuse

	s.BucketN += other.BucketN
	s.InlineBucketN += other.InlineBucketN
	s.InlineBucketInuse += other.InlineBucketInuse
***REMOVED***

// cloneBytes returns a copy of a given slice.
func cloneBytes(v []byte) []byte ***REMOVED***
	var clone = make([]byte, len(v))
	copy(clone, v)
	return clone
***REMOVED***
