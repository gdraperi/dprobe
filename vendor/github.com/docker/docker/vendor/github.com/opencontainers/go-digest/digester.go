package digest

import "hash"

// Digester calculates the digest of written data. Writes should go directly
// to the return value of Hash, while calling Digest will return the current
// value of the digest.
type Digester interface ***REMOVED***
	Hash() hash.Hash // provides direct access to underlying hash instance.
	Digest() Digest
***REMOVED***

// digester provides a simple digester definition that embeds a hasher.
type digester struct ***REMOVED***
	alg  Algorithm
	hash hash.Hash
***REMOVED***

func (d *digester) Hash() hash.Hash ***REMOVED***
	return d.hash
***REMOVED***

func (d *digester) Digest() Digest ***REMOVED***
	return NewDigest(d.alg, d.hash)
***REMOVED***
