package daemon

import (
	"github.com/docker/docker/builder"
	"github.com/docker/docker/image/cache"
	"github.com/sirupsen/logrus"
)

// MakeImageCache creates a stateful image cache.
func (daemon *Daemon) MakeImageCache(sourceRefs []string) builder.ImageCache ***REMOVED***
	if len(sourceRefs) == 0 ***REMOVED***
		return cache.NewLocal(daemon.imageStore)
	***REMOVED***

	cache := cache.New(daemon.imageStore)

	for _, ref := range sourceRefs ***REMOVED***
		img, err := daemon.GetImage(ref)
		if err != nil ***REMOVED***
			logrus.Warnf("Could not look up %s for cache resolution, skipping: %+v", ref, err)
			continue
		***REMOVED***
		cache.Populate(img)
	***REMOVED***

	return cache
***REMOVED***
