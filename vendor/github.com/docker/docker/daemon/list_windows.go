package daemon

import (
	"strings"

	"github.com/docker/docker/container"
)

// excludeByIsolation is a platform specific helper function to support PS
// filtering by Isolation. This is a Windows-only concept, so is a no-op on Unix.
func excludeByIsolation(container *container.Snapshot, ctx *listContext) iterationAction ***REMOVED***
	i := strings.ToLower(string(container.HostConfig.Isolation))
	if i == "" ***REMOVED***
		i = "default"
	***REMOVED***
	if !ctx.filters.Match("isolation", i) ***REMOVED***
		return excludeContainer
	***REMOVED***
	return includeContainer
***REMOVED***
