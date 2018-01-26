package memory

import (
	"sync"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/storage/cache"
	"github.com/opencontainers/go-digest"
)

type inMemoryBlobDescriptorCacheProvider struct ***REMOVED***
	global       *mapBlobDescriptorCache
	repositories map[string]*mapBlobDescriptorCache
	mu           sync.RWMutex
***REMOVED***

// NewInMemoryBlobDescriptorCacheProvider returns a new mapped-based cache for
// storing blob descriptor data.
func NewInMemoryBlobDescriptorCacheProvider() cache.BlobDescriptorCacheProvider ***REMOVED***
	return &inMemoryBlobDescriptorCacheProvider***REMOVED***
		global:       newMapBlobDescriptorCache(),
		repositories: make(map[string]*mapBlobDescriptorCache),
	***REMOVED***
***REMOVED***

func (imbdcp *inMemoryBlobDescriptorCacheProvider) RepositoryScoped(repo string) (distribution.BlobDescriptorService, error) ***REMOVED***
	if _, err := reference.ParseNormalizedNamed(repo); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	imbdcp.mu.RLock()
	defer imbdcp.mu.RUnlock()

	return &repositoryScopedInMemoryBlobDescriptorCache***REMOVED***
		repo:       repo,
		parent:     imbdcp,
		repository: imbdcp.repositories[repo],
	***REMOVED***, nil
***REMOVED***

func (imbdcp *inMemoryBlobDescriptorCacheProvider) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) ***REMOVED***
	return imbdcp.global.Stat(ctx, dgst)
***REMOVED***

func (imbdcp *inMemoryBlobDescriptorCacheProvider) Clear(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	return imbdcp.global.Clear(ctx, dgst)
***REMOVED***

func (imbdcp *inMemoryBlobDescriptorCacheProvider) SetDescriptor(ctx context.Context, dgst digest.Digest, desc distribution.Descriptor) error ***REMOVED***
	_, err := imbdcp.Stat(ctx, dgst)
	if err == distribution.ErrBlobUnknown ***REMOVED***

		if dgst.Algorithm() != desc.Digest.Algorithm() && dgst != desc.Digest ***REMOVED***
			// if the digests differ, set the other canonical mapping
			if err := imbdcp.global.SetDescriptor(ctx, desc.Digest, desc); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		// unknown, just set it
		return imbdcp.global.SetDescriptor(ctx, dgst, desc)
	***REMOVED***

	// we already know it, do nothing
	return err
***REMOVED***

// repositoryScopedInMemoryBlobDescriptorCache provides the request scoped
// repository cache. Instances are not thread-safe but the delegated
// operations are.
type repositoryScopedInMemoryBlobDescriptorCache struct ***REMOVED***
	repo       string
	parent     *inMemoryBlobDescriptorCacheProvider // allows lazy allocation of repo's map
	repository *mapBlobDescriptorCache
***REMOVED***

func (rsimbdcp *repositoryScopedInMemoryBlobDescriptorCache) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) ***REMOVED***
	rsimbdcp.parent.mu.Lock()
	repo := rsimbdcp.repository
	rsimbdcp.parent.mu.Unlock()

	if repo == nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, distribution.ErrBlobUnknown
	***REMOVED***

	return repo.Stat(ctx, dgst)
***REMOVED***

func (rsimbdcp *repositoryScopedInMemoryBlobDescriptorCache) Clear(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	rsimbdcp.parent.mu.Lock()
	repo := rsimbdcp.repository
	rsimbdcp.parent.mu.Unlock()

	if repo == nil ***REMOVED***
		return distribution.ErrBlobUnknown
	***REMOVED***

	return repo.Clear(ctx, dgst)
***REMOVED***

func (rsimbdcp *repositoryScopedInMemoryBlobDescriptorCache) SetDescriptor(ctx context.Context, dgst digest.Digest, desc distribution.Descriptor) error ***REMOVED***
	rsimbdcp.parent.mu.Lock()
	repo := rsimbdcp.repository
	if repo == nil ***REMOVED***
		// allocate map since we are setting it now.
		var ok bool
		// have to read back value since we may have allocated elsewhere.
		repo, ok = rsimbdcp.parent.repositories[rsimbdcp.repo]
		if !ok ***REMOVED***
			repo = newMapBlobDescriptorCache()
			rsimbdcp.parent.repositories[rsimbdcp.repo] = repo
		***REMOVED***
		rsimbdcp.repository = repo
	***REMOVED***
	rsimbdcp.parent.mu.Unlock()

	if err := repo.SetDescriptor(ctx, dgst, desc); err != nil ***REMOVED***
		return err
	***REMOVED***

	return rsimbdcp.parent.SetDescriptor(ctx, dgst, desc)
***REMOVED***

// mapBlobDescriptorCache provides a simple map-based implementation of the
// descriptor cache.
type mapBlobDescriptorCache struct ***REMOVED***
	descriptors map[digest.Digest]distribution.Descriptor
	mu          sync.RWMutex
***REMOVED***

var _ distribution.BlobDescriptorService = &mapBlobDescriptorCache***REMOVED******REMOVED***

func newMapBlobDescriptorCache() *mapBlobDescriptorCache ***REMOVED***
	return &mapBlobDescriptorCache***REMOVED***
		descriptors: make(map[digest.Digest]distribution.Descriptor),
	***REMOVED***
***REMOVED***

func (mbdc *mapBlobDescriptorCache) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) ***REMOVED***
	if err := dgst.Validate(); err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	mbdc.mu.RLock()
	defer mbdc.mu.RUnlock()

	desc, ok := mbdc.descriptors[dgst]
	if !ok ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, distribution.ErrBlobUnknown
	***REMOVED***

	return desc, nil
***REMOVED***

func (mbdc *mapBlobDescriptorCache) Clear(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	mbdc.mu.Lock()
	defer mbdc.mu.Unlock()

	delete(mbdc.descriptors, dgst)
	return nil
***REMOVED***

func (mbdc *mapBlobDescriptorCache) SetDescriptor(ctx context.Context, dgst digest.Digest, desc distribution.Descriptor) error ***REMOVED***
	if err := dgst.Validate(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := cache.ValidateDescriptor(desc); err != nil ***REMOVED***
		return err
	***REMOVED***

	mbdc.mu.Lock()
	defer mbdc.mu.Unlock()

	mbdc.descriptors[dgst] = desc
	return nil
***REMOVED***
