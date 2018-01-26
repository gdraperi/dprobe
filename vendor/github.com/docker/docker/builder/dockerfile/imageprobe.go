package dockerfile

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder"
	"github.com/sirupsen/logrus"
)

// ImageProber exposes an Image cache to the Builder. It supports resetting a
// cache.
type ImageProber interface ***REMOVED***
	Reset()
	Probe(parentID string, runConfig *container.Config) (string, error)
***REMOVED***

type imageProber struct ***REMOVED***
	cache       builder.ImageCache
	reset       func() builder.ImageCache
	cacheBusted bool
***REMOVED***

func newImageProber(cacheBuilder builder.ImageCacheBuilder, cacheFrom []string, noCache bool) ImageProber ***REMOVED***
	if noCache ***REMOVED***
		return &nopProber***REMOVED******REMOVED***
	***REMOVED***

	reset := func() builder.ImageCache ***REMOVED***
		return cacheBuilder.MakeImageCache(cacheFrom)
	***REMOVED***
	return &imageProber***REMOVED***cache: reset(), reset: reset***REMOVED***
***REMOVED***

func (c *imageProber) Reset() ***REMOVED***
	c.cache = c.reset()
	c.cacheBusted = false
***REMOVED***

// Probe checks if cache match can be found for current build instruction.
// It returns the cachedID if there is a hit, and the empty string on miss
func (c *imageProber) Probe(parentID string, runConfig *container.Config) (string, error) ***REMOVED***
	if c.cacheBusted ***REMOVED***
		return "", nil
	***REMOVED***
	cacheID, err := c.cache.GetCache(parentID, runConfig)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if len(cacheID) == 0 ***REMOVED***
		logrus.Debugf("[BUILDER] Cache miss: %s", runConfig.Cmd)
		c.cacheBusted = true
		return "", nil
	***REMOVED***
	logrus.Debugf("[BUILDER] Use cached version: %s", runConfig.Cmd)
	return cacheID, nil
***REMOVED***

type nopProber struct***REMOVED******REMOVED***

func (c *nopProber) Reset() ***REMOVED******REMOVED***

func (c *nopProber) Probe(_ string, _ *container.Config) (string, error) ***REMOVED***
	return "", nil
***REMOVED***
