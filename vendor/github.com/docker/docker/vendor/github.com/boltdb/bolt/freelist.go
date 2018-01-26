package bolt

import (
	"fmt"
	"sort"
	"unsafe"
)

// freelist represents a list of all pages that are available for allocation.
// It also tracks pages that have been freed but are still in use by open transactions.
type freelist struct ***REMOVED***
	ids     []pgid          // all free and available free page ids.
	pending map[txid][]pgid // mapping of soon-to-be free page ids by tx.
	cache   map[pgid]bool   // fast lookup of all free and pending page ids.
***REMOVED***

// newFreelist returns an empty, initialized freelist.
func newFreelist() *freelist ***REMOVED***
	return &freelist***REMOVED***
		pending: make(map[txid][]pgid),
		cache:   make(map[pgid]bool),
	***REMOVED***
***REMOVED***

// size returns the size of the page after serialization.
func (f *freelist) size() int ***REMOVED***
	return pageHeaderSize + (int(unsafe.Sizeof(pgid(0))) * f.count())
***REMOVED***

// count returns count of pages on the freelist
func (f *freelist) count() int ***REMOVED***
	return f.free_count() + f.pending_count()
***REMOVED***

// free_count returns count of free pages
func (f *freelist) free_count() int ***REMOVED***
	return len(f.ids)
***REMOVED***

// pending_count returns count of pending pages
func (f *freelist) pending_count() int ***REMOVED***
	var count int
	for _, list := range f.pending ***REMOVED***
		count += len(list)
	***REMOVED***
	return count
***REMOVED***

// all returns a list of all free ids and all pending ids in one sorted list.
func (f *freelist) all() []pgid ***REMOVED***
	m := make(pgids, 0)

	for _, list := range f.pending ***REMOVED***
		m = append(m, list...)
	***REMOVED***

	sort.Sort(m)
	return pgids(f.ids).merge(m)
***REMOVED***

// allocate returns the starting page id of a contiguous list of pages of a given size.
// If a contiguous block cannot be found then 0 is returned.
func (f *freelist) allocate(n int) pgid ***REMOVED***
	if len(f.ids) == 0 ***REMOVED***
		return 0
	***REMOVED***

	var initial, previd pgid
	for i, id := range f.ids ***REMOVED***
		if id <= 1 ***REMOVED***
			panic(fmt.Sprintf("invalid page allocation: %d", id))
		***REMOVED***

		// Reset initial page if this is not contiguous.
		if previd == 0 || id-previd != 1 ***REMOVED***
			initial = id
		***REMOVED***

		// If we found a contiguous block then remove it and return it.
		if (id-initial)+1 == pgid(n) ***REMOVED***
			// If we're allocating off the beginning then take the fast path
			// and just adjust the existing slice. This will use extra memory
			// temporarily but the append() in free() will realloc the slice
			// as is necessary.
			if (i + 1) == n ***REMOVED***
				f.ids = f.ids[i+1:]
			***REMOVED*** else ***REMOVED***
				copy(f.ids[i-n+1:], f.ids[i+1:])
				f.ids = f.ids[:len(f.ids)-n]
			***REMOVED***

			// Remove from the free cache.
			for i := pgid(0); i < pgid(n); i++ ***REMOVED***
				delete(f.cache, initial+i)
			***REMOVED***

			return initial
		***REMOVED***

		previd = id
	***REMOVED***
	return 0
***REMOVED***

// free releases a page and its overflow for a given transaction id.
// If the page is already free then a panic will occur.
func (f *freelist) free(txid txid, p *page) ***REMOVED***
	if p.id <= 1 ***REMOVED***
		panic(fmt.Sprintf("cannot free page 0 or 1: %d", p.id))
	***REMOVED***

	// Free page and all its overflow pages.
	var ids = f.pending[txid]
	for id := p.id; id <= p.id+pgid(p.overflow); id++ ***REMOVED***
		// Verify that page is not already free.
		if f.cache[id] ***REMOVED***
			panic(fmt.Sprintf("page %d already freed", id))
		***REMOVED***

		// Add to the freelist and cache.
		ids = append(ids, id)
		f.cache[id] = true
	***REMOVED***
	f.pending[txid] = ids
***REMOVED***

// release moves all page ids for a transaction id (or older) to the freelist.
func (f *freelist) release(txid txid) ***REMOVED***
	m := make(pgids, 0)
	for tid, ids := range f.pending ***REMOVED***
		if tid <= txid ***REMOVED***
			// Move transaction's pending pages to the available freelist.
			// Don't remove from the cache since the page is still free.
			m = append(m, ids...)
			delete(f.pending, tid)
		***REMOVED***
	***REMOVED***
	sort.Sort(m)
	f.ids = pgids(f.ids).merge(m)
***REMOVED***

// rollback removes the pages from a given pending tx.
func (f *freelist) rollback(txid txid) ***REMOVED***
	// Remove page ids from cache.
	for _, id := range f.pending[txid] ***REMOVED***
		delete(f.cache, id)
	***REMOVED***

	// Remove pages from pending list.
	delete(f.pending, txid)
***REMOVED***

// freed returns whether a given page is in the free list.
func (f *freelist) freed(pgid pgid) bool ***REMOVED***
	return f.cache[pgid]
***REMOVED***

// read initializes the freelist from a freelist page.
func (f *freelist) read(p *page) ***REMOVED***
	// If the page.count is at the max uint16 value (64k) then it's considered
	// an overflow and the size of the freelist is stored as the first element.
	idx, count := 0, int(p.count)
	if count == 0xFFFF ***REMOVED***
		idx = 1
		count = int(((*[maxAllocSize]pgid)(unsafe.Pointer(&p.ptr)))[0])
	***REMOVED***

	// Copy the list of page ids from the freelist.
	if count == 0 ***REMOVED***
		f.ids = nil
	***REMOVED*** else ***REMOVED***
		ids := ((*[maxAllocSize]pgid)(unsafe.Pointer(&p.ptr)))[idx:count]
		f.ids = make([]pgid, len(ids))
		copy(f.ids, ids)

		// Make sure they're sorted.
		sort.Sort(pgids(f.ids))
	***REMOVED***

	// Rebuild the page cache.
	f.reindex()
***REMOVED***

// write writes the page ids onto a freelist page. All free and pending ids are
// saved to disk since in the event of a program crash, all pending ids will
// become free.
func (f *freelist) write(p *page) error ***REMOVED***
	// Combine the old free pgids and pgids waiting on an open transaction.
	ids := f.all()

	// Update the header flag.
	p.flags |= freelistPageFlag

	// The page.count can only hold up to 64k elements so if we overflow that
	// number then we handle it by putting the size in the first element.
	if len(ids) == 0 ***REMOVED***
		p.count = uint16(len(ids))
	***REMOVED*** else if len(ids) < 0xFFFF ***REMOVED***
		p.count = uint16(len(ids))
		copy(((*[maxAllocSize]pgid)(unsafe.Pointer(&p.ptr)))[:], ids)
	***REMOVED*** else ***REMOVED***
		p.count = 0xFFFF
		((*[maxAllocSize]pgid)(unsafe.Pointer(&p.ptr)))[0] = pgid(len(ids))
		copy(((*[maxAllocSize]pgid)(unsafe.Pointer(&p.ptr)))[1:], ids)
	***REMOVED***

	return nil
***REMOVED***

// reload reads the freelist from a page and filters out pending items.
func (f *freelist) reload(p *page) ***REMOVED***
	f.read(p)

	// Build a cache of only pending pages.
	pcache := make(map[pgid]bool)
	for _, pendingIDs := range f.pending ***REMOVED***
		for _, pendingID := range pendingIDs ***REMOVED***
			pcache[pendingID] = true
		***REMOVED***
	***REMOVED***

	// Check each page in the freelist and build a new available freelist
	// with any pages not in the pending lists.
	var a []pgid
	for _, id := range f.ids ***REMOVED***
		if !pcache[id] ***REMOVED***
			a = append(a, id)
		***REMOVED***
	***REMOVED***
	f.ids = a

	// Once the available list is rebuilt then rebuild the free cache so that
	// it includes the available and pending free pages.
	f.reindex()
***REMOVED***

// reindex rebuilds the free cache based on available and pending free lists.
func (f *freelist) reindex() ***REMOVED***
	f.cache = make(map[pgid]bool, len(f.ids))
	for _, id := range f.ids ***REMOVED***
		f.cache[id] = true
	***REMOVED***
	for _, pendingIDs := range f.pending ***REMOVED***
		for _, pendingID := range pendingIDs ***REMOVED***
			f.cache[pendingID] = true
		***REMOVED***
	***REMOVED***
***REMOVED***
