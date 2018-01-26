// Package cache provides facilities to speed up access to the storage
// backend.
package cache

import (
	"fmt"

	"github.com/docker/distribution"
)

// BlobDescriptorCacheProvider provides repository scoped
// BlobDescriptorService cache instances and a global descriptor cache.
type BlobDescriptorCacheProvider interface ***REMOVED***
	distribution.BlobDescriptorService

	RepositoryScoped(repo string) (distribution.BlobDescriptorService, error)
***REMOVED***

// ValidateDescriptor provides a helper function to ensure that caches have
// common criteria for admitting descriptors.
func ValidateDescriptor(desc distribution.Descriptor) error ***REMOVED***
	if err := desc.Digest.Validate(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if desc.Size < 0 ***REMOVED***
		return fmt.Errorf("cache: invalid length in descriptor: %v < 0", desc.Size)
	***REMOVED***

	if desc.MediaType == "" ***REMOVED***
		return fmt.Errorf("cache: empty mediatype on descriptor: %v", desc)
	***REMOVED***

	return nil
***REMOVED***
