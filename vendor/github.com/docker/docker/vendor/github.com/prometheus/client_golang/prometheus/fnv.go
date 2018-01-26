package prometheus

// Inline and byte-free variant of hash/fnv's fnv64a.

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

// hashNew initializies a new fnv64a hash value.
func hashNew() uint64 ***REMOVED***
	return offset64
***REMOVED***

// hashAdd adds a string to a fnv64a hash value, returning the updated hash.
func hashAdd(h uint64, s string) uint64 ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		h ^= uint64(s[i])
		h *= prime64
	***REMOVED***
	return h
***REMOVED***

// hashAddByte adds a byte to a fnv64a hash value, returning the updated hash.
func hashAddByte(h uint64, b byte) uint64 ***REMOVED***
	h ^= uint64(b)
	h *= prime64
	return h
***REMOVED***
