package daemon

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/libcontainerd"
)

func toContainerdResources(resources container.Resources) *libcontainerd.Resources ***REMOVED***
	// We don't support update, so do nothing
	return nil
***REMOVED***
