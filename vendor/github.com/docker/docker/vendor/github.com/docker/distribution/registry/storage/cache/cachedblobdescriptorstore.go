package cache

import (
	"github.com/docker/distribution/context"
	"github.com/opencontainers/go-digest"

	"github.com/docker/distribution"
)

// Metrics is used to hold metric counters
// related to the number of times a cache was
// hit or missed.
type Metrics struct ***REMOVED***
	Requests uint64
	Hits     uint64
	Misses   uint64
***REMOVED***

// MetricsTracker represents a metric tracker
// which simply counts the number of hits and misses.
type MetricsTracker interface ***REMOVED***
	Hit()
	Miss()
	Metrics() Metrics
***REMOVED***

type cachedBlobStatter struct ***REMOVED***
	cache   distribution.BlobDescriptorService
	backend distribution.BlobDescriptorService
	tracker MetricsTracker
***REMOVED***

// NewCachedBlobStatter creates a new statter which prefers a cache and
// falls back to a backend.
func NewCachedBlobStatter(cache distribution.BlobDescriptorService, backend distribution.BlobDescriptorService) distribution.BlobDescriptorService ***REMOVED***
	return &cachedBlobStatter***REMOVED***
		cache:   cache,
		backend: backend,
	***REMOVED***
***REMOVED***

// NewCachedBlobStatterWithMetrics creates a new statter which prefers a cache and
// falls back to a backend. Hits and misses will send to the tracker.
func NewCachedBlobStatterWithMetrics(cache distribution.BlobDescriptorService, backend distribution.BlobDescriptorService, tracker MetricsTracker) distribution.BlobStatter ***REMOVED***
	return &cachedBlobStatter***REMOVED***
		cache:   cache,
		backend: backend,
		tracker: tracker,
	***REMOVED***
***REMOVED***

func (cbds *cachedBlobStatter) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) ***REMOVED***
	desc, err := cbds.cache.Stat(ctx, dgst)
	if err != nil ***REMOVED***
		if err != distribution.ErrBlobUnknown ***REMOVED***
			context.GetLogger(ctx).Errorf("error retrieving descriptor from cache: %v", err)
		***REMOVED***

		goto fallback
	***REMOVED***

	if cbds.tracker != nil ***REMOVED***
		cbds.tracker.Hit()
	***REMOVED***
	return desc, nil
fallback:
	if cbds.tracker != nil ***REMOVED***
		cbds.tracker.Miss()
	***REMOVED***
	desc, err = cbds.backend.Stat(ctx, dgst)
	if err != nil ***REMOVED***
		return desc, err
	***REMOVED***

	if err := cbds.cache.SetDescriptor(ctx, dgst, desc); err != nil ***REMOVED***
		context.GetLogger(ctx).Errorf("error adding descriptor %v to cache: %v", desc.Digest, err)
	***REMOVED***

	return desc, err

***REMOVED***

func (cbds *cachedBlobStatter) Clear(ctx context.Context, dgst digest.Digest) error ***REMOVED***
	err := cbds.cache.Clear(ctx, dgst)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = cbds.backend.Clear(ctx, dgst)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (cbds *cachedBlobStatter) SetDescriptor(ctx context.Context, dgst digest.Digest, desc distribution.Descriptor) error ***REMOVED***
	if err := cbds.cache.SetDescriptor(ctx, dgst, desc); err != nil ***REMOVED***
		context.GetLogger(ctx).Errorf("error adding descriptor %v to cache: %v", desc.Digest, err)
	***REMOVED***
	return nil
***REMOVED***
