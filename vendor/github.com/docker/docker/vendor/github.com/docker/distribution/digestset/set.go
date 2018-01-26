package digestset

import (
	"errors"
	"sort"
	"strings"
	"sync"

	digest "github.com/opencontainers/go-digest"
)

var (
	// ErrDigestNotFound is used when a matching digest
	// could not be found in a set.
	ErrDigestNotFound = errors.New("digest not found")

	// ErrDigestAmbiguous is used when multiple digests
	// are found in a set. None of the matching digests
	// should be considered valid matches.
	ErrDigestAmbiguous = errors.New("ambiguous digest string")
)

// Set is used to hold a unique set of digests which
// may be easily referenced by easily  referenced by a string
// representation of the digest as well as short representation.
// The uniqueness of the short representation is based on other
// digests in the set. If digests are omitted from this set,
// collisions in a larger set may not be detected, therefore it
// is important to always do short representation lookups on
// the complete set of digests. To mitigate collisions, an
// appropriately long short code should be used.
type Set struct ***REMOVED***
	mutex   sync.RWMutex
	entries digestEntries
***REMOVED***

// NewSet creates an empty set of digests
// which may have digests added.
func NewSet() *Set ***REMOVED***
	return &Set***REMOVED***
		entries: digestEntries***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// checkShortMatch checks whether two digests match as either whole
// values or short values. This function does not test equality,
// rather whether the second value could match against the first
// value.
func checkShortMatch(alg digest.Algorithm, hex, shortAlg, shortHex string) bool ***REMOVED***
	if len(hex) == len(shortHex) ***REMOVED***
		if hex != shortHex ***REMOVED***
			return false
		***REMOVED***
		if len(shortAlg) > 0 && string(alg) != shortAlg ***REMOVED***
			return false
		***REMOVED***
	***REMOVED*** else if !strings.HasPrefix(hex, shortHex) ***REMOVED***
		return false
	***REMOVED*** else if len(shortAlg) > 0 && string(alg) != shortAlg ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// Lookup looks for a digest matching the given string representation.
// If no digests could be found ErrDigestNotFound will be returned
// with an empty digest value. If multiple matches are found
// ErrDigestAmbiguous will be returned with an empty digest value.
func (dst *Set) Lookup(d string) (digest.Digest, error) ***REMOVED***
	dst.mutex.RLock()
	defer dst.mutex.RUnlock()
	if len(dst.entries) == 0 ***REMOVED***
		return "", ErrDigestNotFound
	***REMOVED***
	var (
		searchFunc func(int) bool
		alg        digest.Algorithm
		hex        string
	)
	dgst, err := digest.Parse(d)
	if err == digest.ErrDigestInvalidFormat ***REMOVED***
		hex = d
		searchFunc = func(i int) bool ***REMOVED***
			return dst.entries[i].val >= d
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		hex = dgst.Hex()
		alg = dgst.Algorithm()
		searchFunc = func(i int) bool ***REMOVED***
			if dst.entries[i].val == hex ***REMOVED***
				return dst.entries[i].alg >= alg
			***REMOVED***
			return dst.entries[i].val >= hex
		***REMOVED***
	***REMOVED***
	idx := sort.Search(len(dst.entries), searchFunc)
	if idx == len(dst.entries) || !checkShortMatch(dst.entries[idx].alg, dst.entries[idx].val, string(alg), hex) ***REMOVED***
		return "", ErrDigestNotFound
	***REMOVED***
	if dst.entries[idx].alg == alg && dst.entries[idx].val == hex ***REMOVED***
		return dst.entries[idx].digest, nil
	***REMOVED***
	if idx+1 < len(dst.entries) && checkShortMatch(dst.entries[idx+1].alg, dst.entries[idx+1].val, string(alg), hex) ***REMOVED***
		return "", ErrDigestAmbiguous
	***REMOVED***

	return dst.entries[idx].digest, nil
***REMOVED***

// Add adds the given digest to the set. An error will be returned
// if the given digest is invalid. If the digest already exists in the
// set, this operation will be a no-op.
func (dst *Set) Add(d digest.Digest) error ***REMOVED***
	if err := d.Validate(); err != nil ***REMOVED***
		return err
	***REMOVED***
	dst.mutex.Lock()
	defer dst.mutex.Unlock()
	entry := &digestEntry***REMOVED***alg: d.Algorithm(), val: d.Hex(), digest: d***REMOVED***
	searchFunc := func(i int) bool ***REMOVED***
		if dst.entries[i].val == entry.val ***REMOVED***
			return dst.entries[i].alg >= entry.alg
		***REMOVED***
		return dst.entries[i].val >= entry.val
	***REMOVED***
	idx := sort.Search(len(dst.entries), searchFunc)
	if idx == len(dst.entries) ***REMOVED***
		dst.entries = append(dst.entries, entry)
		return nil
	***REMOVED*** else if dst.entries[idx].digest == d ***REMOVED***
		return nil
	***REMOVED***

	entries := append(dst.entries, nil)
	copy(entries[idx+1:], entries[idx:len(entries)-1])
	entries[idx] = entry
	dst.entries = entries
	return nil
***REMOVED***

// Remove removes the given digest from the set. An err will be
// returned if the given digest is invalid. If the digest does
// not exist in the set, this operation will be a no-op.
func (dst *Set) Remove(d digest.Digest) error ***REMOVED***
	if err := d.Validate(); err != nil ***REMOVED***
		return err
	***REMOVED***
	dst.mutex.Lock()
	defer dst.mutex.Unlock()
	entry := &digestEntry***REMOVED***alg: d.Algorithm(), val: d.Hex(), digest: d***REMOVED***
	searchFunc := func(i int) bool ***REMOVED***
		if dst.entries[i].val == entry.val ***REMOVED***
			return dst.entries[i].alg >= entry.alg
		***REMOVED***
		return dst.entries[i].val >= entry.val
	***REMOVED***
	idx := sort.Search(len(dst.entries), searchFunc)
	// Not found if idx is after or value at idx is not digest
	if idx == len(dst.entries) || dst.entries[idx].digest != d ***REMOVED***
		return nil
	***REMOVED***

	entries := dst.entries
	copy(entries[idx:], entries[idx+1:])
	entries = entries[:len(entries)-1]
	dst.entries = entries

	return nil
***REMOVED***

// All returns all the digests in the set
func (dst *Set) All() []digest.Digest ***REMOVED***
	dst.mutex.RLock()
	defer dst.mutex.RUnlock()
	retValues := make([]digest.Digest, len(dst.entries))
	for i := range dst.entries ***REMOVED***
		retValues[i] = dst.entries[i].digest
	***REMOVED***

	return retValues
***REMOVED***

// ShortCodeTable returns a map of Digest to unique short codes. The
// length represents the minimum value, the maximum length may be the
// entire value of digest if uniqueness cannot be achieved without the
// full value. This function will attempt to make short codes as short
// as possible to be unique.
func ShortCodeTable(dst *Set, length int) map[digest.Digest]string ***REMOVED***
	dst.mutex.RLock()
	defer dst.mutex.RUnlock()
	m := make(map[digest.Digest]string, len(dst.entries))
	l := length
	resetIdx := 0
	for i := 0; i < len(dst.entries); i++ ***REMOVED***
		var short string
		extended := true
		for extended ***REMOVED***
			extended = false
			if len(dst.entries[i].val) <= l ***REMOVED***
				short = dst.entries[i].digest.String()
			***REMOVED*** else ***REMOVED***
				short = dst.entries[i].val[:l]
				for j := i + 1; j < len(dst.entries); j++ ***REMOVED***
					if checkShortMatch(dst.entries[j].alg, dst.entries[j].val, "", short) ***REMOVED***
						if j > resetIdx ***REMOVED***
							resetIdx = j
						***REMOVED***
						extended = true
					***REMOVED*** else ***REMOVED***
						break
					***REMOVED***
				***REMOVED***
				if extended ***REMOVED***
					l++
				***REMOVED***
			***REMOVED***
		***REMOVED***
		m[dst.entries[i].digest] = short
		if i >= resetIdx ***REMOVED***
			l = length
		***REMOVED***
	***REMOVED***
	return m
***REMOVED***

type digestEntry struct ***REMOVED***
	alg    digest.Algorithm
	val    string
	digest digest.Digest
***REMOVED***

type digestEntries []*digestEntry

func (d digestEntries) Len() int ***REMOVED***
	return len(d)
***REMOVED***

func (d digestEntries) Less(i, j int) bool ***REMOVED***
	if d[i].val != d[j].val ***REMOVED***
		return d[i].val < d[j].val
	***REMOVED***
	return d[i].alg < d[j].alg
***REMOVED***

func (d digestEntries) Swap(i, j int) ***REMOVED***
	d[i], d[j] = d[j], d[i]
***REMOVED***
