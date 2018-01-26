// Package truncindex provides a general 'index tree', used by Docker
// in order to be able to reference containers by only a few unambiguous
// characters of their id.
package truncindex

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/tchap/go-patricia/patricia"
)

var (
	// ErrEmptyPrefix is an error returned if the prefix was empty.
	ErrEmptyPrefix = errors.New("Prefix can't be empty")

	// ErrIllegalChar is returned when a space is in the ID
	ErrIllegalChar = errors.New("illegal character: ' '")

	// ErrNotExist is returned when ID or its prefix not found in index.
	ErrNotExist = errors.New("ID does not exist")
)

// ErrAmbiguousPrefix is returned if the prefix was ambiguous
// (multiple ids for the prefix).
type ErrAmbiguousPrefix struct ***REMOVED***
	prefix string
***REMOVED***

func (e ErrAmbiguousPrefix) Error() string ***REMOVED***
	return fmt.Sprintf("Multiple IDs found with provided prefix: %s", e.prefix)
***REMOVED***

// TruncIndex allows the retrieval of string identifiers by any of their unique prefixes.
// This is used to retrieve image and container IDs by more convenient shorthand prefixes.
type TruncIndex struct ***REMOVED***
	sync.RWMutex
	trie *patricia.Trie
	ids  map[string]struct***REMOVED******REMOVED***
***REMOVED***

// NewTruncIndex creates a new TruncIndex and initializes with a list of IDs.
func NewTruncIndex(ids []string) (idx *TruncIndex) ***REMOVED***
	idx = &TruncIndex***REMOVED***
		ids: make(map[string]struct***REMOVED******REMOVED***),

		// Change patricia max prefix per node length,
		// because our len(ID) always 64
		trie: patricia.NewTrie(patricia.MaxPrefixPerNode(64)),
	***REMOVED***
	for _, id := range ids ***REMOVED***
		idx.addID(id)
	***REMOVED***
	return
***REMOVED***

func (idx *TruncIndex) addID(id string) error ***REMOVED***
	if strings.Contains(id, " ") ***REMOVED***
		return ErrIllegalChar
	***REMOVED***
	if id == "" ***REMOVED***
		return ErrEmptyPrefix
	***REMOVED***
	if _, exists := idx.ids[id]; exists ***REMOVED***
		return fmt.Errorf("id already exists: '%s'", id)
	***REMOVED***
	idx.ids[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	if inserted := idx.trie.Insert(patricia.Prefix(id), struct***REMOVED******REMOVED******REMOVED******REMOVED***); !inserted ***REMOVED***
		return fmt.Errorf("failed to insert id: %s", id)
	***REMOVED***
	return nil
***REMOVED***

// Add adds a new ID to the TruncIndex.
func (idx *TruncIndex) Add(id string) error ***REMOVED***
	idx.Lock()
	defer idx.Unlock()
	return idx.addID(id)
***REMOVED***

// Delete removes an ID from the TruncIndex. If there are multiple IDs
// with the given prefix, an error is thrown.
func (idx *TruncIndex) Delete(id string) error ***REMOVED***
	idx.Lock()
	defer idx.Unlock()
	if _, exists := idx.ids[id]; !exists || id == "" ***REMOVED***
		return fmt.Errorf("no such id: '%s'", id)
	***REMOVED***
	delete(idx.ids, id)
	if deleted := idx.trie.Delete(patricia.Prefix(id)); !deleted ***REMOVED***
		return fmt.Errorf("no such id: '%s'", id)
	***REMOVED***
	return nil
***REMOVED***

// Get retrieves an ID from the TruncIndex. If there are multiple IDs
// with the given prefix, an error is thrown.
func (idx *TruncIndex) Get(s string) (string, error) ***REMOVED***
	if s == "" ***REMOVED***
		return "", ErrEmptyPrefix
	***REMOVED***
	var (
		id string
	)
	subTreeVisitFunc := func(prefix patricia.Prefix, item patricia.Item) error ***REMOVED***
		if id != "" ***REMOVED***
			// we haven't found the ID if there are two or more IDs
			id = ""
			return ErrAmbiguousPrefix***REMOVED***prefix: string(prefix)***REMOVED***
		***REMOVED***
		id = string(prefix)
		return nil
	***REMOVED***

	idx.RLock()
	defer idx.RUnlock()
	if err := idx.trie.VisitSubtree(patricia.Prefix(s), subTreeVisitFunc); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if id != "" ***REMOVED***
		return id, nil
	***REMOVED***
	return "", ErrNotExist
***REMOVED***

// Iterate iterates over all stored IDs and passes each of them to the given
// handler. Take care that the handler method does not call any public
// method on truncindex as the internal locking is not reentrant/recursive
// and will result in deadlock.
func (idx *TruncIndex) Iterate(handler func(id string)) ***REMOVED***
	idx.Lock()
	defer idx.Unlock()
	idx.trie.Visit(func(prefix patricia.Prefix, item patricia.Item) error ***REMOVED***
		handler(string(prefix))
		return nil
	***REMOVED***)
***REMOVED***
