package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/libcontainerd"
)

// ContainerResize changes the size of the TTY of the process running
// in the container with the given name to the given height and width.
func (daemon *Daemon) ContainerResize(name string, height, width int) error ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !container.IsRunning() ***REMOVED***
		return errNotRunning(container.ID)
	***REMOVED***

	if err = daemon.containerd.ResizeTerminal(context.Background(), container.ID, libcontainerd.InitProcessName, width, height); err == nil ***REMOVED***
		attributes := map[string]string***REMOVED***
			"height": fmt.Sprintf("%d", height),
			"width":  fmt.Sprintf("%d", width),
		***REMOVED***
		daemon.LogContainerEventWithAttributes(container, "resize", attributes)
	***REMOVED***
	return err
***REMOVED***

// ContainerExecResize changes the size of the TTY of the process
// running in the exec with the given name to the given height and
// width.
func (daemon *Daemon) ContainerExecResize(name string, height, width int) error ***REMOVED***
	ec, err := daemon.getExecConfig(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return daemon.containerd.ResizeTerminal(context.Background(), ec.ContainerID, ec.ID, width, height)
***REMOVED***
