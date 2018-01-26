// Package bitseq provides a structure and utilities for representing long bitmask
// as sequence of run-length encoded blocks. It operates directly on the encoded
// representation, it does not decode/encode.
package bitseq

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

// block sequence constants
// If needed we can think of making these configurable
const (
	blockLen      = uint32(32)
	blockBytes    = uint64(blockLen / 8)
	blockMAX      = uint32(1<<blockLen - 1)
	blockFirstBit = uint32(1) << (blockLen - 1)
	invalidPos    = uint64(0xFFFFFFFFFFFFFFFF)
)

var (
	// ErrNoBitAvailable is returned when no more bits are available to set
	ErrNoBitAvailable = errors.New("no bit available")
	// ErrBitAllocated is returned when the specific bit requested is already set
	ErrBitAllocated = errors.New("requested bit is already allocated")
)

// Handle contains the sequece representing the bitmask and its identifier
type Handle struct ***REMOVED***
	bits       uint64
	unselected uint64
	head       *sequence
	app        string
	id         string
	dbIndex    uint64
	dbExists   bool
	curr       uint64
	store      datastore.DataStore
	sync.Mutex
***REMOVED***

// NewHandle returns a thread-safe instance of the bitmask handler
func NewHandle(app string, ds datastore.DataStore, id string, numElements uint64) (*Handle, error) ***REMOVED***
	h := &Handle***REMOVED***
		app:        app,
		id:         id,
		store:      ds,
		bits:       numElements,
		unselected: numElements,
		head: &sequence***REMOVED***
			block: 0x0,
			count: getNumBlocks(numElements),
		***REMOVED***,
	***REMOVED***

	if h.store == nil ***REMOVED***
		return h, nil
	***REMOVED***

	// Get the initial status from the ds if present.
	if err := h.store.GetObject(datastore.Key(h.Key()...), h); err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		return nil, err
	***REMOVED***

	// If the handle is not in store, write it.
	if !h.Exists() ***REMOVED***
		if err := h.writeToStore(); err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to write bitsequence to store: %v", err)
		***REMOVED***
	***REMOVED***

	return h, nil
***REMOVED***

// sequence represents a recurring sequence of 32 bits long bitmasks
type sequence struct ***REMOVED***
	block uint32    // block is a symbol representing 4 byte long allocation bitmask
	count uint64    // number of consecutive blocks (symbols)
	next  *sequence // next sequence
***REMOVED***

// String returns a string representation of the block sequence starting from this block
func (s *sequence) toString() string ***REMOVED***
	var nextBlock string
	if s.next == nil ***REMOVED***
		nextBlock = "end"
	***REMOVED*** else ***REMOVED***
		nextBlock = s.next.toString()
	***REMOVED***
	return fmt.Sprintf("(0x%x, %d)->%s", s.block, s.count, nextBlock)
***REMOVED***

// GetAvailableBit returns the position of the first unset bit in the bitmask represented by this sequence
func (s *sequence) getAvailableBit(from uint64) (uint64, uint64, error) ***REMOVED***
	if s.block == blockMAX || s.count == 0 ***REMOVED***
		return invalidPos, invalidPos, ErrNoBitAvailable
	***REMOVED***
	bits := from
	bitSel := blockFirstBit >> from
	for bitSel > 0 && s.block&bitSel != 0 ***REMOVED***
		bitSel >>= 1
		bits++
	***REMOVED***
	return bits / 8, bits % 8, nil
***REMOVED***

// GetCopy returns a copy of the linked list rooted at this node
func (s *sequence) getCopy() *sequence ***REMOVED***
	n := &sequence***REMOVED***block: s.block, count: s.count***REMOVED***
	pn := n
	ps := s.next
	for ps != nil ***REMOVED***
		pn.next = &sequence***REMOVED***block: ps.block, count: ps.count***REMOVED***
		pn = pn.next
		ps = ps.next
	***REMOVED***
	return n
***REMOVED***

// Equal checks if this sequence is equal to the passed one
func (s *sequence) equal(o *sequence) bool ***REMOVED***
	this := s
	other := o
	for this != nil ***REMOVED***
		if other == nil ***REMOVED***
			return false
		***REMOVED***
		if this.block != other.block || this.count != other.count ***REMOVED***
			return false
		***REMOVED***
		this = this.next
		other = other.next
	***REMOVED***
	// Check if other is longer than this
	if other != nil ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// ToByteArray converts the sequence into a byte array
func (s *sequence) toByteArray() ([]byte, error) ***REMOVED***
	var bb []byte

	p := s
	for p != nil ***REMOVED***
		b := make([]byte, 12)
		binary.BigEndian.PutUint32(b[0:], p.block)
		binary.BigEndian.PutUint64(b[4:], p.count)
		bb = append(bb, b...)
		p = p.next
	***REMOVED***

	return bb, nil
***REMOVED***

// fromByteArray construct the sequence from the byte array
func (s *sequence) fromByteArray(data []byte) error ***REMOVED***
	l := len(data)
	if l%12 != 0 ***REMOVED***
		return fmt.Errorf("cannot deserialize byte sequence of length %d (%v)", l, data)
	***REMOVED***

	p := s
	i := 0
	for ***REMOVED***
		p.block = binary.BigEndian.Uint32(data[i : i+4])
		p.count = binary.BigEndian.Uint64(data[i+4 : i+12])
		i += 12
		if i == l ***REMOVED***
			break
		***REMOVED***
		p.next = &sequence***REMOVED******REMOVED***
		p = p.next
	***REMOVED***

	return nil
***REMOVED***

func (h *Handle) getCopy() *Handle ***REMOVED***
	return &Handle***REMOVED***
		bits:       h.bits,
		unselected: h.unselected,
		head:       h.head.getCopy(),
		app:        h.app,
		id:         h.id,
		dbIndex:    h.dbIndex,
		dbExists:   h.dbExists,
		store:      h.store,
		curr:       h.curr,
	***REMOVED***
***REMOVED***

// SetAnyInRange atomically sets the first unset bit in the specified range in the sequence and returns the corresponding ordinal
func (h *Handle) SetAnyInRange(start, end uint64, serial bool) (uint64, error) ***REMOVED***
	if end < start || end >= h.bits ***REMOVED***
		return invalidPos, fmt.Errorf("invalid bit range [%d, %d]", start, end)
	***REMOVED***
	if h.Unselected() == 0 ***REMOVED***
		return invalidPos, ErrNoBitAvailable
	***REMOVED***
	return h.set(0, start, end, true, false, serial)
***REMOVED***

// SetAny atomically sets the first unset bit in the sequence and returns the corresponding ordinal
func (h *Handle) SetAny(serial bool) (uint64, error) ***REMOVED***
	if h.Unselected() == 0 ***REMOVED***
		return invalidPos, ErrNoBitAvailable
	***REMOVED***
	return h.set(0, 0, h.bits-1, true, false, serial)
***REMOVED***

// Set atomically sets the corresponding bit in the sequence
func (h *Handle) Set(ordinal uint64) error ***REMOVED***
	if err := h.validateOrdinal(ordinal); err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err := h.set(ordinal, 0, 0, false, false, false)
	return err
***REMOVED***

// Unset atomically unsets the corresponding bit in the sequence
func (h *Handle) Unset(ordinal uint64) error ***REMOVED***
	if err := h.validateOrdinal(ordinal); err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err := h.set(ordinal, 0, 0, false, true, false)
	return err
***REMOVED***

// IsSet atomically checks if the ordinal bit is set. In case ordinal
// is outside of the bit sequence limits, false is returned.
func (h *Handle) IsSet(ordinal uint64) bool ***REMOVED***
	if err := h.validateOrdinal(ordinal); err != nil ***REMOVED***
		return false
	***REMOVED***
	h.Lock()
	_, _, err := checkIfAvailable(h.head, ordinal)
	h.Unlock()
	return err != nil
***REMOVED***

func (h *Handle) runConsistencyCheck() bool ***REMOVED***
	corrupted := false
	for p, c := h.head, h.head.next; c != nil; c = c.next ***REMOVED***
		if c.count == 0 ***REMOVED***
			corrupted = true
			p.next = c.next
			continue // keep same p
		***REMOVED***
		p = c
	***REMOVED***
	return corrupted
***REMOVED***

// CheckConsistency checks if the bit sequence is in an inconsistent state and attempts to fix it.
// It looks for a corruption signature that may happen in docker 1.9.0 and 1.9.1.
func (h *Handle) CheckConsistency() error ***REMOVED***
	for ***REMOVED***
		h.Lock()
		store := h.store
		h.Unlock()

		if store != nil ***REMOVED***
			if err := store.GetObject(datastore.Key(h.Key()...), h); err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		h.Lock()
		nh := h.getCopy()
		h.Unlock()

		if !nh.runConsistencyCheck() ***REMOVED***
			return nil
		***REMOVED***

		if err := nh.writeToStore(); err != nil ***REMOVED***
			if _, ok := err.(types.RetryError); !ok ***REMOVED***
				return fmt.Errorf("internal failure while fixing inconsistent bitsequence: %v", err)
			***REMOVED***
			continue
		***REMOVED***

		logrus.Infof("Fixed inconsistent bit sequence in datastore:\n%s\n%s", h, nh)

		h.Lock()
		h.head = nh.head
		h.Unlock()

		return nil
	***REMOVED***
***REMOVED***

// set/reset the bit
func (h *Handle) set(ordinal, start, end uint64, any bool, release bool, serial bool) (uint64, error) ***REMOVED***
	var (
		bitPos  uint64
		bytePos uint64
		ret     uint64
		err     error
	)

	for ***REMOVED***
		var store datastore.DataStore
		curr := uint64(0)
		h.Lock()
		store = h.store
		h.Unlock()
		if store != nil ***REMOVED***
			if err := store.GetObject(datastore.Key(h.Key()...), h); err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
				return ret, err
			***REMOVED***
		***REMOVED***

		h.Lock()
		if serial ***REMOVED***
			curr = h.curr
		***REMOVED***
		// Get position if available
		if release ***REMOVED***
			bytePos, bitPos = ordinalToPos(ordinal)
		***REMOVED*** else ***REMOVED***
			if any ***REMOVED***
				bytePos, bitPos, err = getAvailableFromCurrent(h.head, start, curr, end)
				ret = posToOrdinal(bytePos, bitPos)
				if err == nil ***REMOVED***
					h.curr = ret + 1
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				bytePos, bitPos, err = checkIfAvailable(h.head, ordinal)
				ret = ordinal
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			h.Unlock()
			return ret, err
		***REMOVED***

		// Create a private copy of h and work on it
		nh := h.getCopy()
		h.Unlock()

		nh.head = pushReservation(bytePos, bitPos, nh.head, release)
		if release ***REMOVED***
			nh.unselected++
		***REMOVED*** else ***REMOVED***
			nh.unselected--
		***REMOVED***

		// Attempt to write private copy to store
		if err := nh.writeToStore(); err != nil ***REMOVED***
			if _, ok := err.(types.RetryError); !ok ***REMOVED***
				return ret, fmt.Errorf("internal failure while setting the bit: %v", err)
			***REMOVED***
			// Retry
			continue
		***REMOVED***

		// Previous atomic push was succesfull. Save private copy to local copy
		h.Lock()
		defer h.Unlock()
		h.unselected = nh.unselected
		h.head = nh.head
		h.dbExists = nh.dbExists
		h.dbIndex = nh.dbIndex
		return ret, nil
	***REMOVED***
***REMOVED***

// checks is needed because to cover the case where the number of bits is not a multiple of blockLen
func (h *Handle) validateOrdinal(ordinal uint64) error ***REMOVED***
	h.Lock()
	defer h.Unlock()
	if ordinal >= h.bits ***REMOVED***
		return errors.New("bit does not belong to the sequence")
	***REMOVED***
	return nil
***REMOVED***

// Destroy removes from the datastore the data belonging to this handle
func (h *Handle) Destroy() error ***REMOVED***
	for ***REMOVED***
		if err := h.deleteFromStore(); err != nil ***REMOVED***
			if _, ok := err.(types.RetryError); !ok ***REMOVED***
				return fmt.Errorf("internal failure while destroying the sequence: %v", err)
			***REMOVED***
			// Fetch latest
			if err := h.store.GetObject(datastore.Key(h.Key()...), h); err != nil ***REMOVED***
				if err == datastore.ErrKeyNotFound ***REMOVED*** // already removed
					return nil
				***REMOVED***
				return fmt.Errorf("failed to fetch from store when destroying the sequence: %v", err)
			***REMOVED***
			continue
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// ToByteArray converts this handle's data into a byte array
func (h *Handle) ToByteArray() ([]byte, error) ***REMOVED***

	h.Lock()
	defer h.Unlock()
	ba := make([]byte, 16)
	binary.BigEndian.PutUint64(ba[0:], h.bits)
	binary.BigEndian.PutUint64(ba[8:], h.unselected)
	bm, err := h.head.toByteArray()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to serialize head: %s", err.Error())
	***REMOVED***
	ba = append(ba, bm...)

	return ba, nil
***REMOVED***

// FromByteArray reads his handle's data from a byte array
func (h *Handle) FromByteArray(ba []byte) error ***REMOVED***
	if ba == nil ***REMOVED***
		return errors.New("nil byte array")
	***REMOVED***

	nh := &sequence***REMOVED******REMOVED***
	err := nh.fromByteArray(ba[16:])
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to deserialize head: %s", err.Error())
	***REMOVED***

	h.Lock()
	h.head = nh
	h.bits = binary.BigEndian.Uint64(ba[0:8])
	h.unselected = binary.BigEndian.Uint64(ba[8:16])
	h.Unlock()

	return nil
***REMOVED***

// Bits returns the length of the bit sequence
func (h *Handle) Bits() uint64 ***REMOVED***
	return h.bits
***REMOVED***

// Unselected returns the number of bits which are not selected
func (h *Handle) Unselected() uint64 ***REMOVED***
	h.Lock()
	defer h.Unlock()
	return h.unselected
***REMOVED***

func (h *Handle) String() string ***REMOVED***
	h.Lock()
	defer h.Unlock()
	return fmt.Sprintf("App: %s, ID: %s, DBIndex: 0x%x, bits: %d, unselected: %d, sequence: %s",
		h.app, h.id, h.dbIndex, h.bits, h.unselected, h.head.toString())
***REMOVED***

// MarshalJSON encodes Handle into json message
func (h *Handle) MarshalJSON() ([]byte, error) ***REMOVED***
	m := map[string]interface***REMOVED******REMOVED******REMOVED***
		"id": h.id,
	***REMOVED***

	b, err := h.ToByteArray()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m["sequence"] = b
	return json.Marshal(m)
***REMOVED***

// UnmarshalJSON decodes json message into Handle
func (h *Handle) UnmarshalJSON(data []byte) error ***REMOVED***
	var (
		m   map[string]interface***REMOVED******REMOVED***
		b   []byte
		err error
	)
	if err = json.Unmarshal(data, &m); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.id = m["id"].(string)
	bi, _ := json.Marshal(m["sequence"])
	if err := json.Unmarshal(bi, &b); err != nil ***REMOVED***
		return err
	***REMOVED***
	return h.FromByteArray(b)
***REMOVED***

// getFirstAvailable looks for the first unset bit in passed mask starting from start
func getFirstAvailable(head *sequence, start uint64) (uint64, uint64, error) ***REMOVED***
	// Find sequence which contains the start bit
	byteStart, bitStart := ordinalToPos(start)
	current, _, _, inBlockBytePos := findSequence(head, byteStart)

	// Derive the this sequence offsets
	byteOffset := byteStart - inBlockBytePos
	bitOffset := inBlockBytePos*8 + bitStart
	var firstOffset uint64
	if current == head ***REMOVED***
		firstOffset = byteOffset
	***REMOVED***
	for current != nil ***REMOVED***
		if current.block != blockMAX ***REMOVED***
			bytePos, bitPos, err := current.getAvailableBit(bitOffset)
			return byteOffset + bytePos, bitPos, err
		***REMOVED***
		// Moving to next block: Reset bit offset.
		bitOffset = 0
		byteOffset += (current.count * blockBytes) - firstOffset
		firstOffset = 0
		current = current.next
	***REMOVED***
	return invalidPos, invalidPos, ErrNoBitAvailable
***REMOVED***

// getAvailableFromCurrent will look for available ordinal from the current ordinal.
// If none found then it will loop back to the start to check of the available bit.
// This can be further optimized to check from start till curr in case of a rollover
func getAvailableFromCurrent(head *sequence, start, curr, end uint64) (uint64, uint64, error) ***REMOVED***
	var bytePos, bitPos uint64
	if curr != 0 && curr > start ***REMOVED***
		bytePos, bitPos, _ = getFirstAvailable(head, curr)
		ret := posToOrdinal(bytePos, bitPos)
		if end < ret ***REMOVED***
			goto begin
		***REMOVED***
		return bytePos, bitPos, nil
	***REMOVED***

begin:
	bytePos, bitPos, _ = getFirstAvailable(head, start)
	ret := posToOrdinal(bytePos, bitPos)
	if end < ret ***REMOVED***
		return invalidPos, invalidPos, ErrNoBitAvailable
	***REMOVED***
	return bytePos, bitPos, nil
***REMOVED***

// checkIfAvailable checks if the bit correspondent to the specified ordinal is unset
// If the ordinal is beyond the sequence limits, a negative response is returned
func checkIfAvailable(head *sequence, ordinal uint64) (uint64, uint64, error) ***REMOVED***
	bytePos, bitPos := ordinalToPos(ordinal)

	// Find the sequence containing this byte
	current, _, _, inBlockBytePos := findSequence(head, bytePos)
	if current != nil ***REMOVED***
		// Check whether the bit corresponding to the ordinal address is unset
		bitSel := blockFirstBit >> (inBlockBytePos*8 + bitPos)
		if current.block&bitSel == 0 ***REMOVED***
			return bytePos, bitPos, nil
		***REMOVED***
	***REMOVED***

	return invalidPos, invalidPos, ErrBitAllocated
***REMOVED***

// Given the byte position and the sequences list head, return the pointer to the
// sequence containing the byte (current), the pointer to the previous sequence,
// the number of blocks preceding the block containing the byte inside the current sequence.
// If bytePos is outside of the list, function will return (nil, nil, 0, invalidPos)
func findSequence(head *sequence, bytePos uint64) (*sequence, *sequence, uint64, uint64) ***REMOVED***
	// Find the sequence containing this byte
	previous := head
	current := head
	n := bytePos
	for current.next != nil && n >= (current.count*blockBytes) ***REMOVED*** // Nil check for less than 32 addresses masks
		n -= (current.count * blockBytes)
		previous = current
		current = current.next
	***REMOVED***

	// If byte is outside of the list, let caller know
	if n >= (current.count * blockBytes) ***REMOVED***
		return nil, nil, 0, invalidPos
	***REMOVED***

	// Find the byte position inside the block and the number of blocks
	// preceding the block containing the byte inside this sequence
	precBlocks := n / blockBytes
	inBlockBytePos := bytePos % blockBytes

	return current, previous, precBlocks, inBlockBytePos
***REMOVED***

// PushReservation pushes the bit reservation inside the bitmask.
// Given byte and bit positions, identify the sequence (current) which holds the block containing the affected bit.
// Create a new block with the modified bit according to the operation (allocate/release).
// Create a new sequence containing the new block and insert it in the proper position.
// Remove current sequence if empty.
// Check if new sequence can be merged with neighbour (previous/next) sequences.
//
//
// Identify "current" sequence containing block:
//                                      [prev seq] [current seq] [next seq]
//
// Based on block position, resulting list of sequences can be any of three forms:
//
//        block position                        Resulting list of sequences
// A) block is first in current:         [prev seq] [new] [modified current seq] [next seq]
// B) block is last in current:          [prev seq] [modified current seq] [new] [next seq]
// C) block is in the middle of current: [prev seq] [curr pre] [new] [curr post] [next seq]
func pushReservation(bytePos, bitPos uint64, head *sequence, release bool) *sequence ***REMOVED***
	// Store list's head
	newHead := head

	// Find the sequence containing this byte
	current, previous, precBlocks, inBlockBytePos := findSequence(head, bytePos)
	if current == nil ***REMOVED***
		return newHead
	***REMOVED***

	// Construct updated block
	bitSel := blockFirstBit >> (inBlockBytePos*8 + bitPos)
	newBlock := current.block
	if release ***REMOVED***
		newBlock &^= bitSel
	***REMOVED*** else ***REMOVED***
		newBlock |= bitSel
	***REMOVED***

	// Quit if it was a redundant request
	if current.block == newBlock ***REMOVED***
		return newHead
	***REMOVED***

	// Current sequence inevitably looses one block, upadate count
	current.count--

	// Create new sequence
	newSequence := &sequence***REMOVED***block: newBlock, count: 1***REMOVED***

	// Insert the new sequence in the list based on block position
	if precBlocks == 0 ***REMOVED*** // First in sequence (A)
		newSequence.next = current
		if current == head ***REMOVED***
			newHead = newSequence
			previous = newHead
		***REMOVED*** else ***REMOVED***
			previous.next = newSequence
		***REMOVED***
		removeCurrentIfEmpty(&newHead, newSequence, current)
		mergeSequences(previous)
	***REMOVED*** else if precBlocks == current.count ***REMOVED*** // Last in sequence (B)
		newSequence.next = current.next
		current.next = newSequence
		mergeSequences(current)
	***REMOVED*** else ***REMOVED*** // In between the sequence (C)
		currPre := &sequence***REMOVED***block: current.block, count: precBlocks, next: newSequence***REMOVED***
		currPost := current
		currPost.count -= precBlocks
		newSequence.next = currPost
		if currPost == head ***REMOVED***
			newHead = currPre
		***REMOVED*** else ***REMOVED***
			previous.next = currPre
		***REMOVED***
		// No merging or empty current possible here
	***REMOVED***

	return newHead
***REMOVED***

// Removes the current sequence from the list if empty, adjusting the head pointer if needed
func removeCurrentIfEmpty(head **sequence, previous, current *sequence) ***REMOVED***
	if current.count == 0 ***REMOVED***
		if current == *head ***REMOVED***
			*head = current.next
		***REMOVED*** else ***REMOVED***
			previous.next = current.next
			current = current.next
		***REMOVED***
	***REMOVED***
***REMOVED***

// Given a pointer to a sequence, it checks if it can be merged with any following sequences
// It stops when no more merging is possible.
// TODO: Optimization: only attempt merge from start to end sequence, no need to scan till the end of the list
func mergeSequences(seq *sequence) ***REMOVED***
	if seq != nil ***REMOVED***
		// Merge all what possible from seq
		for seq.next != nil && seq.block == seq.next.block ***REMOVED***
			seq.count += seq.next.count
			seq.next = seq.next.next
		***REMOVED***
		// Move to next
		mergeSequences(seq.next)
	***REMOVED***
***REMOVED***

func getNumBlocks(numBits uint64) uint64 ***REMOVED***
	numBlocks := numBits / uint64(blockLen)
	if numBits%uint64(blockLen) != 0 ***REMOVED***
		numBlocks++
	***REMOVED***
	return numBlocks
***REMOVED***

func ordinalToPos(ordinal uint64) (uint64, uint64) ***REMOVED***
	return ordinal / 8, ordinal % 8
***REMOVED***

func posToOrdinal(bytePos, bitPos uint64) uint64 ***REMOVED***
	return bytePos*8 + bitPos
***REMOVED***
