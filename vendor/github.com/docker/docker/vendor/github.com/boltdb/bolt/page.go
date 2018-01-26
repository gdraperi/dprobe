package bolt

import (
	"fmt"
	"os"
	"sort"
	"unsafe"
)

const pageHeaderSize = int(unsafe.Offsetof(((*page)(nil)).ptr))

const minKeysPerPage = 2

const branchPageElementSize = int(unsafe.Sizeof(branchPageElement***REMOVED******REMOVED***))
const leafPageElementSize = int(unsafe.Sizeof(leafPageElement***REMOVED******REMOVED***))

const (
	branchPageFlag   = 0x01
	leafPageFlag     = 0x02
	metaPageFlag     = 0x04
	freelistPageFlag = 0x10
)

const (
	bucketLeafFlag = 0x01
)

type pgid uint64

type page struct ***REMOVED***
	id       pgid
	flags    uint16
	count    uint16
	overflow uint32
	ptr      uintptr
***REMOVED***

// typ returns a human readable page type string used for debugging.
func (p *page) typ() string ***REMOVED***
	if (p.flags & branchPageFlag) != 0 ***REMOVED***
		return "branch"
	***REMOVED*** else if (p.flags & leafPageFlag) != 0 ***REMOVED***
		return "leaf"
	***REMOVED*** else if (p.flags & metaPageFlag) != 0 ***REMOVED***
		return "meta"
	***REMOVED*** else if (p.flags & freelistPageFlag) != 0 ***REMOVED***
		return "freelist"
	***REMOVED***
	return fmt.Sprintf("unknown<%02x>", p.flags)
***REMOVED***

// meta returns a pointer to the metadata section of the page.
func (p *page) meta() *meta ***REMOVED***
	return (*meta)(unsafe.Pointer(&p.ptr))
***REMOVED***

// leafPageElement retrieves the leaf node by index
func (p *page) leafPageElement(index uint16) *leafPageElement ***REMOVED***
	n := &((*[0x7FFFFFF]leafPageElement)(unsafe.Pointer(&p.ptr)))[index]
	return n
***REMOVED***

// leafPageElements retrieves a list of leaf nodes.
func (p *page) leafPageElements() []leafPageElement ***REMOVED***
	if p.count == 0 ***REMOVED***
		return nil
	***REMOVED***
	return ((*[0x7FFFFFF]leafPageElement)(unsafe.Pointer(&p.ptr)))[:]
***REMOVED***

// branchPageElement retrieves the branch node by index
func (p *page) branchPageElement(index uint16) *branchPageElement ***REMOVED***
	return &((*[0x7FFFFFF]branchPageElement)(unsafe.Pointer(&p.ptr)))[index]
***REMOVED***

// branchPageElements retrieves a list of branch nodes.
func (p *page) branchPageElements() []branchPageElement ***REMOVED***
	if p.count == 0 ***REMOVED***
		return nil
	***REMOVED***
	return ((*[0x7FFFFFF]branchPageElement)(unsafe.Pointer(&p.ptr)))[:]
***REMOVED***

// dump writes n bytes of the page to STDERR as hex output.
func (p *page) hexdump(n int) ***REMOVED***
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(p))[:n]
	fmt.Fprintf(os.Stderr, "%x\n", buf)
***REMOVED***

type pages []*page

func (s pages) Len() int           ***REMOVED*** return len(s) ***REMOVED***
func (s pages) Swap(i, j int)      ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***
func (s pages) Less(i, j int) bool ***REMOVED*** return s[i].id < s[j].id ***REMOVED***

// branchPageElement represents a node on a branch page.
type branchPageElement struct ***REMOVED***
	pos   uint32
	ksize uint32
	pgid  pgid
***REMOVED***

// key returns a byte slice of the node key.
func (n *branchPageElement) key() []byte ***REMOVED***
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	return (*[maxAllocSize]byte)(unsafe.Pointer(&buf[n.pos]))[:n.ksize]
***REMOVED***

// leafPageElement represents a node on a leaf page.
type leafPageElement struct ***REMOVED***
	flags uint32
	pos   uint32
	ksize uint32
	vsize uint32
***REMOVED***

// key returns a byte slice of the node key.
func (n *leafPageElement) key() []byte ***REMOVED***
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	return (*[maxAllocSize]byte)(unsafe.Pointer(&buf[n.pos]))[:n.ksize:n.ksize]
***REMOVED***

// value returns a byte slice of the node value.
func (n *leafPageElement) value() []byte ***REMOVED***
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	return (*[maxAllocSize]byte)(unsafe.Pointer(&buf[n.pos+n.ksize]))[:n.vsize:n.vsize]
***REMOVED***

// PageInfo represents human readable information about a page.
type PageInfo struct ***REMOVED***
	ID            int
	Type          string
	Count         int
	OverflowCount int
***REMOVED***

type pgids []pgid

func (s pgids) Len() int           ***REMOVED*** return len(s) ***REMOVED***
func (s pgids) Swap(i, j int)      ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***
func (s pgids) Less(i, j int) bool ***REMOVED*** return s[i] < s[j] ***REMOVED***

// merge returns the sorted union of a and b.
func (a pgids) merge(b pgids) pgids ***REMOVED***
	// Return the opposite slice if one is nil.
	if len(a) == 0 ***REMOVED***
		return b
	***REMOVED*** else if len(b) == 0 ***REMOVED***
		return a
	***REMOVED***

	// Create a list to hold all elements from both lists.
	merged := make(pgids, 0, len(a)+len(b))

	// Assign lead to the slice with a lower starting value, follow to the higher value.
	lead, follow := a, b
	if b[0] < a[0] ***REMOVED***
		lead, follow = b, a
	***REMOVED***

	// Continue while there are elements in the lead.
	for len(lead) > 0 ***REMOVED***
		// Merge largest prefix of lead that is ahead of follow[0].
		n := sort.Search(len(lead), func(i int) bool ***REMOVED*** return lead[i] > follow[0] ***REMOVED***)
		merged = append(merged, lead[:n]...)
		if n >= len(lead) ***REMOVED***
			break
		***REMOVED***

		// Swap lead and follow.
		lead, follow = follow, lead[n:]
	***REMOVED***

	// Append what's left in follow.
	merged = append(merged, follow...)

	return merged
***REMOVED***
