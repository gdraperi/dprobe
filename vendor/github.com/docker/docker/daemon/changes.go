package daemon

import (
	"errors"
	"runtime"
	"time"

	"github.com/docker/docker/pkg/archive"
)

// ContainerChanges returns a list of container fs changes
func (daemon *Daemon) ContainerChanges(name string) ([]archive.Change, error) ***REMOVED***
	start := time.Now()
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if runtime.GOOS == "windows" && container.IsRunning() ***REMOVED***
		return nil, errors.New("Windows does not support diff of a running container")
	***REMOVED***

	container.Lock()
	defer container.Unlock()
	c, err := container.RWLayer.Changes()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	containerActions.WithValues("changes").UpdateSince(start)
	return c, nil
***REMOVED***
