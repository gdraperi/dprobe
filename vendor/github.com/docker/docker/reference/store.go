package reference

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

var (
	// ErrDoesNotExist is returned if a reference is not found in the
	// store.
	ErrDoesNotExist notFoundError = "reference does not exist"
)

// An Association is a tuple associating a reference with an image ID.
type Association struct ***REMOVED***
	Ref reference.Named
	ID  digest.Digest
***REMOVED***

// Store provides the set of methods which can operate on a reference store.
type Store interface ***REMOVED***
	References(id digest.Digest) []reference.Named
	ReferencesByName(ref reference.Named) []Association
	AddTag(ref reference.Named, id digest.Digest, force bool) error
	AddDigest(ref reference.Canonical, id digest.Digest, force bool) error
	Delete(ref reference.Named) (bool, error)
	Get(ref reference.Named) (digest.Digest, error)
***REMOVED***

type store struct ***REMOVED***
	mu sync.RWMutex
	// jsonPath is the path to the file where the serialized tag data is
	// stored.
	jsonPath string
	// Repositories is a map of repositories, indexed by name.
	Repositories map[string]repository
	// referencesByIDCache is a cache of references indexed by ID, to speed
	// up References.
	referencesByIDCache map[digest.Digest]map[string]reference.Named
***REMOVED***

// Repository maps tags to digests. The key is a stringified Reference,
// including the repository name.
type repository map[string]digest.Digest

type lexicalRefs []reference.Named

func (a lexicalRefs) Len() int      ***REMOVED*** return len(a) ***REMOVED***
func (a lexicalRefs) Swap(i, j int) ***REMOVED*** a[i], a[j] = a[j], a[i] ***REMOVED***
func (a lexicalRefs) Less(i, j int) bool ***REMOVED***
	return a[i].String() < a[j].String()
***REMOVED***

type lexicalAssociations []Association

func (a lexicalAssociations) Len() int      ***REMOVED*** return len(a) ***REMOVED***
func (a lexicalAssociations) Swap(i, j int) ***REMOVED*** a[i], a[j] = a[j], a[i] ***REMOVED***
func (a lexicalAssociations) Less(i, j int) bool ***REMOVED***
	return a[i].Ref.String() < a[j].Ref.String()
***REMOVED***

// NewReferenceStore creates a new reference store, tied to a file path where
// the set of references are serialized in JSON format.
func NewReferenceStore(jsonPath string) (Store, error) ***REMOVED***
	abspath, err := filepath.Abs(jsonPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	store := &store***REMOVED***
		jsonPath:            abspath,
		Repositories:        make(map[string]repository),
		referencesByIDCache: make(map[digest.Digest]map[string]reference.Named),
	***REMOVED***
	// Load the json file if it exists, otherwise create it.
	if err := store.reload(); os.IsNotExist(err) ***REMOVED***
		if err := store.save(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return store, nil
***REMOVED***

// AddTag adds a tag reference to the store. If force is set to true, existing
// references can be overwritten. This only works for tags, not digests.
func (store *store) AddTag(ref reference.Named, id digest.Digest, force bool) error ***REMOVED***
	if _, isCanonical := ref.(reference.Canonical); isCanonical ***REMOVED***
		return errors.WithStack(invalidTagError("refusing to create a tag with a digest reference"))
	***REMOVED***
	return store.addReference(reference.TagNameOnly(ref), id, force)
***REMOVED***

// AddDigest adds a digest reference to the store.
func (store *store) AddDigest(ref reference.Canonical, id digest.Digest, force bool) error ***REMOVED***
	return store.addReference(ref, id, force)
***REMOVED***

func favorDigest(originalRef reference.Named) (reference.Named, error) ***REMOVED***
	ref := originalRef
	// If the reference includes a digest and a tag, we must store only the
	// digest.
	canonical, isCanonical := originalRef.(reference.Canonical)
	_, isNamedTagged := originalRef.(reference.NamedTagged)

	if isCanonical && isNamedTagged ***REMOVED***
		trimmed, err := reference.WithDigest(reference.TrimNamed(canonical), canonical.Digest())
		if err != nil ***REMOVED***
			// should never happen
			return originalRef, err
		***REMOVED***
		ref = trimmed
	***REMOVED***
	return ref, nil
***REMOVED***

func (store *store) addReference(ref reference.Named, id digest.Digest, force bool) error ***REMOVED***
	ref, err := favorDigest(ref)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	refName := reference.FamiliarName(ref)
	refStr := reference.FamiliarString(ref)

	if refName == string(digest.Canonical) ***REMOVED***
		return errors.WithStack(invalidTagError("refusing to create an ambiguous tag using digest algorithm as name"))
	***REMOVED***

	store.mu.Lock()
	defer store.mu.Unlock()

	repository, exists := store.Repositories[refName]
	if !exists || repository == nil ***REMOVED***
		repository = make(map[string]digest.Digest)
		store.Repositories[refName] = repository
	***REMOVED***

	oldID, exists := repository[refStr]

	if exists ***REMOVED***
		// force only works for tags
		if digested, isDigest := ref.(reference.Canonical); isDigest ***REMOVED***
			return errors.WithStack(conflictingTagError("Cannot overwrite digest " + digested.Digest().String()))
		***REMOVED***

		if !force ***REMOVED***
			return errors.WithStack(
				conflictingTagError(
					fmt.Sprintf("Conflict: Tag %s is already set to image %s, if you want to replace it, please use the force option", refStr, oldID.String()),
				),
			)
		***REMOVED***

		if store.referencesByIDCache[oldID] != nil ***REMOVED***
			delete(store.referencesByIDCache[oldID], refStr)
			if len(store.referencesByIDCache[oldID]) == 0 ***REMOVED***
				delete(store.referencesByIDCache, oldID)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	repository[refStr] = id
	if store.referencesByIDCache[id] == nil ***REMOVED***
		store.referencesByIDCache[id] = make(map[string]reference.Named)
	***REMOVED***
	store.referencesByIDCache[id][refStr] = ref

	return store.save()
***REMOVED***

// Delete deletes a reference from the store. It returns true if a deletion
// happened, or false otherwise.
func (store *store) Delete(ref reference.Named) (bool, error) ***REMOVED***
	ref, err := favorDigest(ref)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	ref = reference.TagNameOnly(ref)

	refName := reference.FamiliarName(ref)
	refStr := reference.FamiliarString(ref)

	store.mu.Lock()
	defer store.mu.Unlock()

	repository, exists := store.Repositories[refName]
	if !exists ***REMOVED***
		return false, ErrDoesNotExist
	***REMOVED***

	if id, exists := repository[refStr]; exists ***REMOVED***
		delete(repository, refStr)
		if len(repository) == 0 ***REMOVED***
			delete(store.Repositories, refName)
		***REMOVED***
		if store.referencesByIDCache[id] != nil ***REMOVED***
			delete(store.referencesByIDCache[id], refStr)
			if len(store.referencesByIDCache[id]) == 0 ***REMOVED***
				delete(store.referencesByIDCache, id)
			***REMOVED***
		***REMOVED***
		return true, store.save()
	***REMOVED***

	return false, ErrDoesNotExist
***REMOVED***

// Get retrieves an item from the store by reference
func (store *store) Get(ref reference.Named) (digest.Digest, error) ***REMOVED***
	if canonical, ok := ref.(reference.Canonical); ok ***REMOVED***
		// If reference contains both tag and digest, only
		// lookup by digest as it takes precedence over
		// tag, until tag/digest combos are stored.
		if _, ok := ref.(reference.Tagged); ok ***REMOVED***
			var err error
			ref, err = reference.WithDigest(reference.TrimNamed(canonical), canonical.Digest())
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		ref = reference.TagNameOnly(ref)
	***REMOVED***

	refName := reference.FamiliarName(ref)
	refStr := reference.FamiliarString(ref)

	store.mu.RLock()
	defer store.mu.RUnlock()

	repository, exists := store.Repositories[refName]
	if !exists || repository == nil ***REMOVED***
		return "", ErrDoesNotExist
	***REMOVED***

	id, exists := repository[refStr]
	if !exists ***REMOVED***
		return "", ErrDoesNotExist
	***REMOVED***

	return id, nil
***REMOVED***

// References returns a slice of references to the given ID. The slice
// will be nil if there are no references to this ID.
func (store *store) References(id digest.Digest) []reference.Named ***REMOVED***
	store.mu.RLock()
	defer store.mu.RUnlock()

	// Convert the internal map to an array for two reasons:
	// 1) We must not return a mutable
	// 2) It would be ugly to expose the extraneous map keys to callers.

	var references []reference.Named
	for _, ref := range store.referencesByIDCache[id] ***REMOVED***
		references = append(references, ref)
	***REMOVED***

	sort.Sort(lexicalRefs(references))

	return references
***REMOVED***

// ReferencesByName returns the references for a given repository name.
// If there are no references known for this repository name,
// ReferencesByName returns nil.
func (store *store) ReferencesByName(ref reference.Named) []Association ***REMOVED***
	refName := reference.FamiliarName(ref)

	store.mu.RLock()
	defer store.mu.RUnlock()

	repository, exists := store.Repositories[refName]
	if !exists ***REMOVED***
		return nil
	***REMOVED***

	var associations []Association
	for refStr, refID := range repository ***REMOVED***
		ref, err := reference.ParseNormalizedNamed(refStr)
		if err != nil ***REMOVED***
			// Should never happen
			return nil
		***REMOVED***
		associations = append(associations,
			Association***REMOVED***
				Ref: ref,
				ID:  refID,
			***REMOVED***)
	***REMOVED***

	sort.Sort(lexicalAssociations(associations))

	return associations
***REMOVED***

func (store *store) save() error ***REMOVED***
	// Store the json
	jsonData, err := json.Marshal(store)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutils.AtomicWriteFile(store.jsonPath, jsonData, 0600)
***REMOVED***

func (store *store) reload() error ***REMOVED***
	f, err := os.Open(store.jsonPath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&store); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, repository := range store.Repositories ***REMOVED***
		for refStr, refID := range repository ***REMOVED***
			ref, err := reference.ParseNormalizedNamed(refStr)
			if err != nil ***REMOVED***
				// Should never happen
				continue
			***REMOVED***
			if store.referencesByIDCache[refID] == nil ***REMOVED***
				store.referencesByIDCache[refID] = make(map[string]reference.Named)
			***REMOVED***
			store.referencesByIDCache[refID][refStr] = ref
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
