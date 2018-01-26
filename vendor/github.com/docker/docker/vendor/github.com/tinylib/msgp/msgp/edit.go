package msgp

import (
	"math"
)

// Locate returns a []byte pointing to the field
// in a messagepack map with the provided key. (The returned []byte
// points to a sub-slice of 'raw'; Locate does no allocations.) If the
// key doesn't exist in the map, a zero-length []byte will be returned.
func Locate(key string, raw []byte) []byte ***REMOVED***
	s, n := locate(raw, key)
	return raw[s:n]
***REMOVED***

// Replace takes a key ("key") in a messagepack map ("raw")
// and replaces its value with the one provided and returns
// the new []byte. The returned []byte may point to the same
// memory as "raw". Replace makes no effort to evaluate the validity
// of the contents of 'val'. It may use up to the full capacity of 'raw.'
// Replace returns 'nil' if the field doesn't exist or if the object in 'raw'
// is not a map.
func Replace(key string, raw []byte, val []byte) []byte ***REMOVED***
	start, end := locate(raw, key)
	if start == end ***REMOVED***
		return nil
	***REMOVED***
	return replace(raw, start, end, val, true)
***REMOVED***

// CopyReplace works similarly to Replace except that the returned
// byte slice does not point to the same memory as 'raw'. CopyReplace
// returns 'nil' if the field doesn't exist or 'raw' isn't a map.
func CopyReplace(key string, raw []byte, val []byte) []byte ***REMOVED***
	start, end := locate(raw, key)
	if start == end ***REMOVED***
		return nil
	***REMOVED***
	return replace(raw, start, end, val, false)
***REMOVED***

// Remove removes a key-value pair from 'raw'. It returns
// 'raw' unchanged if the key didn't exist.
func Remove(key string, raw []byte) []byte ***REMOVED***
	start, end := locateKV(raw, key)
	if start == end ***REMOVED***
		return raw
	***REMOVED***
	raw = raw[:start+copy(raw[start:], raw[end:])]
	return resizeMap(raw, -1)
***REMOVED***

// HasKey returns whether the map in 'raw' has
// a field with key 'key'
func HasKey(key string, raw []byte) bool ***REMOVED***
	sz, bts, err := ReadMapHeaderBytes(raw)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	var field []byte
	for i := uint32(0); i < sz; i++ ***REMOVED***
		field, bts, err = ReadStringZC(bts)
		if err != nil ***REMOVED***
			return false
		***REMOVED***
		if UnsafeString(field) == key ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func replace(raw []byte, start int, end int, val []byte, inplace bool) []byte ***REMOVED***
	ll := end - start // length of segment to replace
	lv := len(val)

	if inplace ***REMOVED***
		extra := lv - ll

		// fastest case: we're doing
		// a 1:1 replacement
		if extra == 0 ***REMOVED***
			copy(raw[start:], val)
			return raw

		***REMOVED*** else if extra < 0 ***REMOVED***
			// 'val' smaller than replaced value
			// copy in place and shift back

			x := copy(raw[start:], val)
			y := copy(raw[start+x:], raw[end:])
			return raw[:start+x+y]

		***REMOVED*** else if extra < cap(raw)-len(raw) ***REMOVED***
			// 'val' less than (cap-len) extra bytes
			// copy in place and shift forward
			raw = raw[0 : len(raw)+extra]
			// shift end forward
			copy(raw[end+extra:], raw[end:])
			copy(raw[start:], val)
			return raw
		***REMOVED***
	***REMOVED***

	// we have to allocate new space
	out := make([]byte, len(raw)+len(val)-ll)
	x := copy(out, raw[:start])
	y := copy(out[x:], val)
	copy(out[x+y:], raw[end:])
	return out
***REMOVED***

// locate does a naive O(n) search for the map key; returns start, end
// (returns 0,0 on error)
func locate(raw []byte, key string) (start int, end int) ***REMOVED***
	var (
		sz    uint32
		bts   []byte
		field []byte
		err   error
	)
	sz, bts, err = ReadMapHeaderBytes(raw)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	// loop and locate field
	for i := uint32(0); i < sz; i++ ***REMOVED***
		field, bts, err = ReadStringZC(bts)
		if err != nil ***REMOVED***
			return 0, 0
		***REMOVED***
		if UnsafeString(field) == key ***REMOVED***
			// start location
			l := len(raw)
			start = l - len(bts)
			bts, err = Skip(bts)
			if err != nil ***REMOVED***
				return 0, 0
			***REMOVED***
			end = l - len(bts)
			return
		***REMOVED***
		bts, err = Skip(bts)
		if err != nil ***REMOVED***
			return 0, 0
		***REMOVED***
	***REMOVED***
	return 0, 0
***REMOVED***

// locate key AND value
func locateKV(raw []byte, key string) (start int, end int) ***REMOVED***
	var (
		sz    uint32
		bts   []byte
		field []byte
		err   error
	)
	sz, bts, err = ReadMapHeaderBytes(raw)
	if err != nil ***REMOVED***
		return 0, 0
	***REMOVED***

	for i := uint32(0); i < sz; i++ ***REMOVED***
		tmp := len(bts)
		field, bts, err = ReadStringZC(bts)
		if err != nil ***REMOVED***
			return 0, 0
		***REMOVED***
		if UnsafeString(field) == key ***REMOVED***
			start = len(raw) - tmp
			bts, err = Skip(bts)
			if err != nil ***REMOVED***
				return 0, 0
			***REMOVED***
			end = len(raw) - len(bts)
			return
		***REMOVED***
		bts, err = Skip(bts)
		if err != nil ***REMOVED***
			return 0, 0
		***REMOVED***
	***REMOVED***
	return 0, 0
***REMOVED***

// delta is delta on map size
func resizeMap(raw []byte, delta int64) []byte ***REMOVED***
	var sz int64
	switch raw[0] ***REMOVED***
	case mmap16:
		sz = int64(big.Uint16(raw[1:]))
		if sz+delta <= math.MaxUint16 ***REMOVED***
			big.PutUint16(raw[1:], uint16(sz+delta))
			return raw
		***REMOVED***
		if cap(raw)-len(raw) >= 2 ***REMOVED***
			raw = raw[0 : len(raw)+2]
			copy(raw[5:], raw[3:])
			big.PutUint32(raw[1:], uint32(sz+delta))
			return raw
		***REMOVED***
		n := make([]byte, 0, len(raw)+5)
		n = AppendMapHeader(n, uint32(sz+delta))
		return append(n, raw[3:]...)

	case mmap32:
		sz = int64(big.Uint32(raw[1:]))
		big.PutUint32(raw[1:], uint32(sz+delta))
		return raw

	default:
		sz = int64(rfixmap(raw[0]))
		if sz+delta < 16 ***REMOVED***
			raw[0] = wfixmap(uint8(sz + delta))
			return raw
		***REMOVED*** else if sz+delta <= math.MaxUint16 ***REMOVED***
			if cap(raw)-len(raw) >= 2 ***REMOVED***
				raw = raw[0 : len(raw)+2]
				copy(raw[3:], raw[1:])
				raw[0] = mmap16
				big.PutUint16(raw[1:], uint16(sz+delta))
				return raw
			***REMOVED***
			n := make([]byte, 0, len(raw)+5)
			n = AppendMapHeader(n, uint32(sz+delta))
			return append(n, raw[1:]...)
		***REMOVED***
		if cap(raw)-len(raw) >= 4 ***REMOVED***
			raw = raw[0 : len(raw)+4]
			copy(raw[5:], raw[1:])
			raw[0] = mmap32
			big.PutUint32(raw[1:], uint32(sz+delta))
			return raw
		***REMOVED***
		n := make([]byte, 0, len(raw)+5)
		n = AppendMapHeader(n, uint32(sz+delta))
		return append(n, raw[1:]...)
	***REMOVED***
***REMOVED***
