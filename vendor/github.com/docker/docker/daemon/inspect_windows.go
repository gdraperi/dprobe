package daemon

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/exec"
)

// This sets platform-specific fields
func setPlatformSpecificContainerFields(container *container.Container, contJSONBase *types.ContainerJSONBase) *types.ContainerJSONBase ***REMOVED***
	return contJSONBase
***REMOVED***

// containerInspectPre120 get containers for pre 1.20 APIs.
func (daemon *Daemon) containerInspectPre120(name string) (*types.ContainerJSON, error) ***REMOVED***
	return daemon.ContainerInspectCurrent(name, false)
***REMOVED***

func inspectExecProcessConfig(e *exec.Config) *backend.ExecProcessConfig ***REMOVED***
	return &backend.ExecProcessConfig***REMOVED***
		Tty:        e.Tty,
		Entrypoint: e.Entrypoint,
		Arguments:  e.Args,
	***REMOVED***
***REMOVED***
