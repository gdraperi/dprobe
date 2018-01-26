// +build !windows

package daemon

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/containerd/containerd/linux/runctypes"
	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/pkg/errors"
)

func (daemon *Daemon) getRuntimeScript(container *container.Container) (string, error) ***REMOVED***
	name := container.HostConfig.Runtime
	rt := daemon.configStore.GetRuntime(name)
	if rt == nil ***REMOVED***
		return "", errdefs.InvalidParameter(errors.Errorf("no such runtime '%s'", name))
	***REMOVED***

	if len(rt.Args) > 0 ***REMOVED***
		// First check that the target exist, as using it in a script won't
		// give us the right error
		if _, err := exec.LookPath(rt.Path); err != nil ***REMOVED***
			return "", translateContainerdStartErr(container.Path, container.SetExitCode, err)
		***REMOVED***
		return filepath.Join(daemon.configStore.Root, "runtimes", name), nil
	***REMOVED***
	return rt.Path, nil
***REMOVED***

// getLibcontainerdCreateOptions callers must hold a lock on the container
func (daemon *Daemon) getLibcontainerdCreateOptions(container *container.Container) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// Ensure a runtime has been assigned to this container
	if container.HostConfig.Runtime == "" ***REMOVED***
		container.HostConfig.Runtime = daemon.configStore.GetDefaultRuntimeName()
		container.CheckpointTo(daemon.containersReplica)
	***REMOVED***

	path, err := daemon.getRuntimeScript(container)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	opts := &runctypes.RuncOptions***REMOVED***
		Runtime: path,
		RuntimeRoot: filepath.Join(daemon.configStore.ExecRoot,
			fmt.Sprintf("runtime-%s", container.HostConfig.Runtime)),
	***REMOVED***

	if UsingSystemd(daemon.configStore) ***REMOVED***
		opts.SystemdCgroup = true
	***REMOVED***

	return opts, nil
***REMOVED***
