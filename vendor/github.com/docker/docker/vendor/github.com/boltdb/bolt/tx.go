package bolt

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
	"unsafe"
)

// txid represents the internal transaction identifier.
type txid uint64

// Tx represents a read-only or read/write transaction on the database.
// Read-only transactions can be used for retrieving values for keys and creating cursors.
// Read/write transactions can create and remove buckets and create and remove keys.
//
// IMPORTANT: You must commit or rollback transactions when you are done with
// them. Pages can not be reclaimed by the writer until no more transactions
// are using them. A long running read transaction can cause the database to
// quickly grow.
type Tx struct ***REMOVED***
	writable       bool
	managed        bool
	db             *DB
	meta           *meta
	root           Bucket
	pages          map[pgid]*page
	stats          TxStats
	commitHandlers []func()

	// WriteFlag specifies the flag for write-related methods like WriteTo().
	// Tx opens the database file with the specified flag to copy the data.
	//
	// By default, the flag is unset, which works well for mostly in-memory
	// workloads. For databases that are much larger than available RAM,
	// set the flag to syscall.O_DIRECT to avoid trashing the page cache.
	WriteFlag int
***REMOVED***

// init initializes the transaction.
func (tx *Tx) init(db *DB) ***REMOVED***
	tx.db = db
	tx.pages = nil

	// Copy the meta page since it can be changed by the writer.
	tx.meta = &meta***REMOVED******REMOVED***
	db.meta().copy(tx.meta)

	// Copy over the root bucket.
	tx.root = newBucket(tx)
	tx.root.bucket = &bucket***REMOVED******REMOVED***
	*tx.root.bucket = tx.meta.root

	// Increment the transaction id and add a page cache for writable transactions.
	if tx.writable ***REMOVED***
		tx.pages = make(map[pgid]*page)
		tx.meta.txid += txid(1)
	***REMOVED***
***REMOVED***

// ID returns the transaction id.
func (tx *Tx) ID() int ***REMOVED***
	return int(tx.meta.txid)
***REMOVED***

// DB returns a reference to the database that created the transaction.
func (tx *Tx) DB() *DB ***REMOVED***
	return tx.db
***REMOVED***

// Size returns current database size in bytes as seen by this transaction.
func (tx *Tx) Size() int64 ***REMOVED***
	return int64(tx.meta.pgid) * int64(tx.db.pageSize)
***REMOVED***

// Writable returns whether the transaction can perform write operations.
func (tx *Tx) Writable() bool ***REMOVED***
	return tx.writable
***REMOVED***

// Cursor creates a cursor associated with the root bucket.
// All items in the cursor will return a nil value because all root bucket keys point to buckets.
// The cursor is only valid as long as the transaction is open.
// Do not use a cursor after the transaction is closed.
func (tx *Tx) Cursor() *Cursor ***REMOVED***
	return tx.root.Cursor()
***REMOVED***

// Stats retrieves a copy of the current transaction statistics.
func (tx *Tx) Stats() TxStats ***REMOVED***
	return tx.stats
***REMOVED***

// Bucket retrieves a bucket by name.
// Returns nil if the bucket does not exist.
// The bucket instance is only valid for the lifetime of the transaction.
func (tx *Tx) Bucket(name []byte) *Bucket ***REMOVED***
	return tx.root.Bucket(name)
***REMOVED***

// CreateBucket creates a new bucket.
// Returns an error if the bucket already exists, if the bucket name is blank, or if the bucket name is too long.
// The bucket instance is only valid for the lifetime of the transaction.
func (tx *Tx) CreateBucket(name []byte) (*Bucket, error) ***REMOVED***
	return tx.root.CreateBucket(name)
***REMOVED***

// CreateBucketIfNotExists creates a new bucket if it doesn't already exist.
// Returns an error if the bucket name is blank, or if the bucket name is too long.
// The bucket instance is only valid for the lifetime of the transaction.
func (tx *Tx) CreateBucketIfNotExists(name []byte) (*Bucket, error) ***REMOVED***
	return tx.root.CreateBucketIfNotExists(name)
***REMOVED***

// DeleteBucket deletes a bucket.
// Returns an error if the bucket cannot be found or if the key represents a non-bucket value.
func (tx *Tx) DeleteBucket(name []byte) error ***REMOVED***
	return tx.root.DeleteBucket(name)
***REMOVED***

// ForEach executes a function for each bucket in the root.
// If the provided function returns an error then the iteration is stopped and
// the error is returned to the caller.
func (tx *Tx) ForEach(fn func(name []byte, b *Bucket) error) error ***REMOVED***
	return tx.root.ForEach(func(k, v []byte) error ***REMOVED***
		if err := fn(k, tx.root.Bucket(k)); err != nil ***REMOVED***
			return err
		***REMOVED***
		return nil
	***REMOVED***)
***REMOVED***

// OnCommit adds a handler function to be executed after the transaction successfully commits.
func (tx *Tx) OnCommit(fn func()) ***REMOVED***
	tx.commitHandlers = append(tx.commitHandlers, fn)
***REMOVED***

// Commit writes all changes to disk and updates the meta page.
// Returns an error if a disk write error occurs, or if Commit is
// called on a read-only transaction.
func (tx *Tx) Commit() error ***REMOVED***
	_assert(!tx.managed, "managed tx commit not allowed")
	if tx.db == nil ***REMOVED***
		return ErrTxClosed
	***REMOVED*** else if !tx.writable ***REMOVED***
		return ErrTxNotWritable
	***REMOVED***

	// TODO(benbjohnson): Use vectorized I/O to write out dirty pages.

	// Rebalance nodes which have had deletions.
	var startTime = time.Now()
	tx.root.rebalance()
	if tx.stats.Rebalance > 0 ***REMOVED***
		tx.stats.RebalanceTime += time.Since(startTime)
	***REMOVED***

	// spill data onto dirty pages.
	startTime = time.Now()
	if err := tx.root.spill(); err != nil ***REMOVED***
		tx.rollback()
		return err
	***REMOVED***
	tx.stats.SpillTime += time.Since(startTime)

	// Free the old root bucket.
	tx.meta.root.root = tx.root.root

	opgid := tx.meta.pgid

	// Free the freelist and allocate new pages for it. This will overestimate
	// the size of the freelist but not underestimate the size (which would be bad).
	tx.db.freelist.free(tx.meta.txid, tx.db.page(tx.meta.freelist))
	p, err := tx.allocate((tx.db.freelist.size() / tx.db.pageSize) + 1)
	if err != nil ***REMOVED***
		tx.rollback()
		return err
	***REMOVED***
	if err := tx.db.freelist.write(p); err != nil ***REMOVED***
		tx.rollback()
		return err
	***REMOVED***
	tx.meta.freelist = p.id

	// If the high water mark has moved up then attempt to grow the database.
	if tx.meta.pgid > opgid ***REMOVED***
		if err := tx.db.grow(int(tx.meta.pgid+1) * tx.db.pageSize); err != nil ***REMOVED***
			tx.rollback()
			return err
		***REMOVED***
	***REMOVED***

	// Write dirty pages to disk.
	startTime = time.Now()
	if err := tx.write(); err != nil ***REMOVED***
		tx.rollback()
		return err
	***REMOVED***

	// If strict mode is enabled then perform a consistency check.
	// Only the first consistency error is reported in the panic.
	if tx.db.StrictMode ***REMOVED***
		ch := tx.Check()
		var errs []string
		for ***REMOVED***
			err, ok := <-ch
			if !ok ***REMOVED***
				break
			***REMOVED***
			errs = append(errs, err.Error())
		***REMOVED***
		if len(errs) > 0 ***REMOVED***
			panic("check fail: " + strings.Join(errs, "\n"))
		***REMOVED***
	***REMOVED***

	// Write meta to disk.
	if err := tx.writeMeta(); err != nil ***REMOVED***
		tx.rollback()
		return err
	***REMOVED***
	tx.stats.WriteTime += time.Since(startTime)

	// Finalize the transaction.
	tx.close()

	// Execute commit handlers now that the locks have been removed.
	for _, fn := range tx.commitHandlers ***REMOVED***
		fn()
	***REMOVED***

	return nil
***REMOVED***

// Rollback closes the transaction and ignores all previous updates. Read-only
// transactions must be rolled back and not committed.
func (tx *Tx) Rollback() error ***REMOVED***
	_assert(!tx.managed, "managed tx rollback not allowed")
	if tx.db == nil ***REMOVED***
		return ErrTxClosed
	***REMOVED***
	tx.rollback()
	return nil
***REMOVED***

func (tx *Tx) rollback() ***REMOVED***
	if tx.db == nil ***REMOVED***
		return
	***REMOVED***
	if tx.writable ***REMOVED***
		tx.db.freelist.rollback(tx.meta.txid)
		tx.db.freelist.reload(tx.db.page(tx.db.meta().freelist))
	***REMOVED***
	tx.close()
***REMOVED***

func (tx *Tx) close() ***REMOVED***
	if tx.db == nil ***REMOVED***
		return
	***REMOVED***
	if tx.writable ***REMOVED***
		// Grab freelist stats.
		var freelistFreeN = tx.db.freelist.free_count()
		var freelistPendingN = tx.db.freelist.pending_count()
		var freelistAlloc = tx.db.freelist.size()

		// Remove transaction ref & writer lock.
		tx.db.rwtx = nil
		tx.db.rwlock.Unlock()

		// Merge statistics.
		tx.db.statlock.Lock()
		tx.db.stats.FreePageN = freelistFreeN
		tx.db.stats.PendingPageN = freelistPendingN
		tx.db.stats.FreeAlloc = (freelistFreeN + freelistPendingN) * tx.db.pageSize
		tx.db.stats.FreelistInuse = freelistAlloc
		tx.db.stats.TxStats.add(&tx.stats)
		tx.db.statlock.Unlock()
	***REMOVED*** else ***REMOVED***
		tx.db.removeTx(tx)
	***REMOVED***

	// Clear all references.
	tx.db = nil
	tx.meta = nil
	tx.root = Bucket***REMOVED***tx: tx***REMOVED***
	tx.pages = nil
***REMOVED***

// Copy writes the entire database to a writer.
// This function exists for backwards compatibility. Use WriteTo() instead.
func (tx *Tx) Copy(w io.Writer) error ***REMOVED***
	_, err := tx.WriteTo(w)
	return err
***REMOVED***

// WriteTo writes the entire database to a writer.
// If err == nil then exactly tx.Size() bytes will be written into the writer.
func (tx *Tx) WriteTo(w io.Writer) (n int64, err error) ***REMOVED***
	// Attempt to open reader with WriteFlag
	f, err := os.OpenFile(tx.db.path, os.O_RDONLY|tx.WriteFlag, 0)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer func() ***REMOVED*** _ = f.Close() ***REMOVED***()

	// Generate a meta page. We use the same page data for both meta pages.
	buf := make([]byte, tx.db.pageSize)
	page := (*page)(unsafe.Pointer(&buf[0]))
	page.flags = metaPageFlag
	*page.meta() = *tx.meta

	// Write meta 0.
	page.id = 0
	page.meta().checksum = page.meta().sum64()
	nn, err := w.Write(buf)
	n += int64(nn)
	if err != nil ***REMOVED***
		return n, fmt.Errorf("meta 0 copy: %s", err)
	***REMOVED***

	// Write meta 1 with a lower transaction id.
	page.id = 1
	page.meta().txid -= 1
	page.meta().checksum = page.meta().sum64()
	nn, err = w.Write(buf)
	n += int64(nn)
	if err != nil ***REMOVED***
		return n, fmt.Errorf("meta 1 copy: %s", err)
	***REMOVED***

	// Move past the meta pages in the file.
	if _, err := f.Seek(int64(tx.db.pageSize*2), os.SEEK_SET); err != nil ***REMOVED***
		return n, fmt.Errorf("seek: %s", err)
	***REMOVED***

	// Copy data pages.
	wn, err := io.CopyN(w, f, tx.Size()-int64(tx.db.pageSize*2))
	n += wn
	if err != nil ***REMOVED***
		return n, err
	***REMOVED***

	return n, f.Close()
***REMOVED***

// CopyFile copies the entire database to file at the given path.
// A reader transaction is maintained during the copy so it is safe to continue
// using the database while a copy is in progress.
func (tx *Tx) CopyFile(path string, mode os.FileMode) error ***REMOVED***
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = tx.Copy(f)
	if err != nil ***REMOVED***
		_ = f.Close()
		return err
	***REMOVED***
	return f.Close()
***REMOVED***

// Check performs several consistency checks on the database for this transaction.
// An error is returned if any inconsistency is found.
//
// It can be safely run concurrently on a writable transaction. However, this
// incurs a high cost for large databases and databases with a lot of subbuckets
// because of caching. This overhead can be removed if running on a read-only
// transaction, however, it is not safe to execute other writer transactions at
// the same time.
func (tx *Tx) Check() <-chan error ***REMOVED***
	ch := make(chan error)
	go tx.check(ch)
	return ch
***REMOVED***

func (tx *Tx) check(ch chan error) ***REMOVED***
	// Check if any pages are double freed.
	freed := make(map[pgid]bool)
	for _, id := range tx.db.freelist.all() ***REMOVED***
		if freed[id] ***REMOVED***
			ch <- fmt.Errorf("page %d: already freed", id)
		***REMOVED***
		freed[id] = true
	***REMOVED***

	// Track every reachable page.
	reachable := make(map[pgid]*page)
	reachable[0] = tx.page(0) // meta0
	reachable[1] = tx.page(1) // meta1
	for i := uint32(0); i <= tx.page(tx.meta.freelist).overflow; i++ ***REMOVED***
		reachable[tx.meta.freelist+pgid(i)] = tx.page(tx.meta.freelist)
	***REMOVED***

	// Recursively check buckets.
	tx.checkBucket(&tx.root, reachable, freed, ch)

	// Ensure all pages below high water mark are either reachable or freed.
	for i := pgid(0); i < tx.meta.pgid; i++ ***REMOVED***
		_, isReachable := reachable[i]
		if !isReachable && !freed[i] ***REMOVED***
			ch <- fmt.Errorf("page %d: unreachable unfreed", int(i))
		***REMOVED***
	***REMOVED***

	// Close the channel to signal completion.
	close(ch)
***REMOVED***

func (tx *Tx) checkBucket(b *Bucket, reachable map[pgid]*page, freed map[pgid]bool, ch chan error) ***REMOVED***
	// Ignore inline buckets.
	if b.root == 0 ***REMOVED***
		return
	***REMOVED***

	// Check every page used by this bucket.
	b.tx.forEachPage(b.root, 0, func(p *page, _ int) ***REMOVED***
		if p.id > tx.meta.pgid ***REMOVED***
			ch <- fmt.Errorf("page %d: out of bounds: %d", int(p.id), int(b.tx.meta.pgid))
		***REMOVED***

		// Ensure each page is only referenced once.
		for i := pgid(0); i <= pgid(p.overflow); i++ ***REMOVED***
			var id = p.id + i
			if _, ok := reachable[id]; ok ***REMOVED***
				ch <- fmt.Errorf("page %d: multiple references", int(id))
			***REMOVED***
			reachable[id] = p
		***REMOVED***

		// We should only encounter un-freed leaf and branch pages.
		if freed[p.id] ***REMOVED***
			ch <- fmt.Errorf("page %d: reachable freed", int(p.id))
		***REMOVED*** else if (p.flags&branchPageFlag) == 0 && (p.flags&leafPageFlag) == 0 ***REMOVED***
			ch <- fmt.Errorf("page %d: invalid type: %s", int(p.id), p.typ())
		***REMOVED***
	***REMOVED***)

	// Check each bucket within this bucket.
	_ = b.ForEach(func(k, v []byte) error ***REMOVED***
		if child := b.Bucket(k); child != nil ***REMOVED***
			tx.checkBucket(child, reachable, freed, ch)
		***REMOVED***
		return nil
	***REMOVED***)
***REMOVED***

// allocate returns a contiguous block of memory starting at a given page.
func (tx *Tx) allocate(count int) (*page, error) ***REMOVED***
	p, err := tx.db.allocate(count)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Save to our page cache.
	tx.pages[p.id] = p

	// Update statistics.
	tx.stats.PageCount++
	tx.stats.PageAlloc += count * tx.db.pageSize

	return p, nil
***REMOVED***

// write writes any dirty pages to disk.
func (tx *Tx) write() error ***REMOVED***
	// Sort pages by id.
	pages := make(pages, 0, len(tx.pages))
	for _, p := range tx.pages ***REMOVED***
		pages = append(pages, p)
	***REMOVED***
	// Clear out page cache early.
	tx.pages = make(map[pgid]*page)
	sort.Sort(pages)

	// Write pages to disk in order.
	for _, p := range pages ***REMOVED***
		size := (int(p.overflow) + 1) * tx.db.pageSize
		offset := int64(p.id) * int64(tx.db.pageSize)

		// Write out page in "max allocation" sized chunks.
		ptr := (*[maxAllocSize]byte)(unsafe.Pointer(p))
		for ***REMOVED***
			// Limit our write to our max allocation size.
			sz := size
			if sz > maxAllocSize-1 ***REMOVED***
				sz = maxAllocSize - 1
			***REMOVED***

			// Write chunk to disk.
			buf := ptr[:sz]
			if _, err := tx.db.ops.writeAt(buf, offset); err != nil ***REMOVED***
				return err
			***REMOVED***

			// Update statistics.
			tx.stats.Write++

			// Exit inner for loop if we've written all the chunks.
			size -= sz
			if size == 0 ***REMOVED***
				break
			***REMOVED***

			// Otherwise move offset forward and move pointer to next chunk.
			offset += int64(sz)
			ptr = (*[maxAllocSize]byte)(unsafe.Pointer(&ptr[sz]))
		***REMOVED***
	***REMOVED***

	// Ignore file sync if flag is set on DB.
	if !tx.db.NoSync || IgnoreNoSync ***REMOVED***
		if err := fdatasync(tx.db); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Put small pages back to page pool.
	for _, p := range pages ***REMOVED***
		// Ignore page sizes over 1 page.
		// These are allocated using make() instead of the page pool.
		if int(p.overflow) != 0 ***REMOVED***
			continue
		***REMOVED***

		buf := (*[maxAllocSize]byte)(unsafe.Pointer(p))[:tx.db.pageSize]

		// See https://go.googlesource.com/go/+/f03c9202c43e0abb130669852082117ca50aa9b1
		for i := range buf ***REMOVED***
			buf[i] = 0
		***REMOVED***
		tx.db.pagePool.Put(buf)
	***REMOVED***

	return nil
***REMOVED***

// writeMeta writes the meta to the disk.
func (tx *Tx) writeMeta() error ***REMOVED***
	// Create a temporary buffer for the meta page.
	buf := make([]byte, tx.db.pageSize)
	p := tx.db.pageInBuffer(buf, 0)
	tx.meta.write(p)

	// Write the meta page to file.
	if _, err := tx.db.ops.writeAt(buf, int64(p.id)*int64(tx.db.pageSize)); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !tx.db.NoSync || IgnoreNoSync ***REMOVED***
		if err := fdatasync(tx.db); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Update statistics.
	tx.stats.Write++

	return nil
***REMOVED***

// page returns a reference to the page with a given id.
// If page has been written to then a temporary buffered page is returned.
func (tx *Tx) page(id pgid) *page ***REMOVED***
	// Check the dirty pages first.
	if tx.pages != nil ***REMOVED***
		if p, ok := tx.pages[id]; ok ***REMOVED***
			return p
		***REMOVED***
	***REMOVED***

	// Otherwise return directly from the mmap.
	return tx.db.page(id)
***REMOVED***

// forEachPage iterates over every page within a given page and executes a function.
func (tx *Tx) forEachPage(pgid pgid, depth int, fn func(*page, int)) ***REMOVED***
	p := tx.page(pgid)

	// Execute function.
	fn(p, depth)

	// Recursively loop over children.
	if (p.flags & branchPageFlag) != 0 ***REMOVED***
		for i := 0; i < int(p.count); i++ ***REMOVED***
			elem := p.branchPageElement(uint16(i))
			tx.forEachPage(elem.pgid, depth+1, fn)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Page returns page information for a given page number.
// This is only safe for concurrent use when used by a writable transaction.
func (tx *Tx) Page(id int) (*PageInfo, error) ***REMOVED***
	if tx.db == nil ***REMOVED***
		return nil, ErrTxClosed
	***REMOVED*** else if pgid(id) >= tx.meta.pgid ***REMOVED***
		return nil, nil
	***REMOVED***

	// Build the page info.
	p := tx.db.page(pgid(id))
	info := &PageInfo***REMOVED***
		ID:            id,
		Count:         int(p.count),
		OverflowCount: int(p.overflow),
	***REMOVED***

	// Determine the type (or if it's free).
	if tx.db.freelist.freed(pgid(id)) ***REMOVED***
		info.Type = "free"
	***REMOVED*** else ***REMOVED***
		info.Type = p.typ()
	***REMOVED***

	return info, nil
***REMOVED***

// TxStats represents statistics about the actions performed by the transaction.
type TxStats struct ***REMOVED***
	// Page statistics.
	PageCount int // number of page allocations
	PageAlloc int // total bytes allocated

	// Cursor statistics.
	CursorCount int // number of cursors created

	// Node statistics
	NodeCount int // number of node allocations
	NodeDeref int // number of node dereferences

	// Rebalance statistics.
	Rebalance     int           // number of node rebalances
	RebalanceTime time.Duration // total time spent rebalancing

	// Split/Spill statistics.
	Split     int           // number of nodes split
	Spill     int           // number of nodes spilled
	SpillTime time.Duration // total time spent spilling

	// Write statistics.
	Write     int           // number of writes performed
	WriteTime time.Duration // total time spent writing to disk
***REMOVED***

func (s *TxStats) add(other *TxStats) ***REMOVED***
	s.PageCount += other.PageCount
	s.PageAlloc += other.PageAlloc
	s.CursorCount += other.CursorCount
	s.NodeCount += other.NodeCount
	s.NodeDeref += other.NodeDeref
	s.Rebalance += other.Rebalance
	s.RebalanceTime += other.RebalanceTime
	s.Split += other.Split
	s.Spill += other.Spill
	s.SpillTime += other.SpillTime
	s.Write += other.Write
	s.WriteTime += other.WriteTime
***REMOVED***

// Sub calculates and returns the difference between two sets of transaction stats.
// This is useful when obtaining stats at two different points and time and
// you need the performance counters that occurred within that time span.
func (s *TxStats) Sub(other *TxStats) TxStats ***REMOVED***
	var diff TxStats
	diff.PageCount = s.PageCount - other.PageCount
	diff.PageAlloc = s.PageAlloc - other.PageAlloc
	diff.CursorCount = s.CursorCount - other.CursorCount
	diff.NodeCount = s.NodeCount - other.NodeCount
	diff.NodeDeref = s.NodeDeref - other.NodeDeref
	diff.Rebalance = s.Rebalance - other.Rebalance
	diff.RebalanceTime = s.RebalanceTime - other.RebalanceTime
	diff.Split = s.Split - other.Split
	diff.Spill = s.Spill - other.Spill
	diff.SpillTime = s.SpillTime - other.SpillTime
	diff.Write = s.Write - other.Write
	diff.WriteTime = s.WriteTime - other.WriteTime
	return diff
***REMOVED***
