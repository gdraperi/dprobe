// +build linux

package cgroups

import (
	"fmt"

	"github.com/opencontainers/runc/libcontainer/configs"
)

type Manager interface ***REMOVED***
	// Applies cgroup configuration to the process with the specified pid
	Apply(pid int) error

	// Returns the PIDs inside the cgroup set
	GetPids() ([]int, error)

	// Returns the PIDs inside the cgroup set & all sub-cgroups
	GetAllPids() ([]int, error)

	// Returns statistics for the cgroup set
	GetStats() (*Stats, error)

	// Toggles the freezer cgroup according with specified state
	Freeze(state configs.FreezerState) error

	// Destroys the cgroup set
	Destroy() error

	// The option func SystemdCgroups() and Cgroupfs() require following attributes:
	// 	Paths   map[string]string
	// 	Cgroups *configs.Cgroup
	// Paths maps cgroup subsystem to path at which it is mounted.
	// Cgroups specifies specific cgroup settings for the various subsystems

	// Returns cgroup paths to save in a state file and to be able to
	// restore the object later.
	GetPaths() map[string]string

	// Sets the cgroup as configured.
	Set(container *configs.Config) error
***REMOVED***

type NotFoundError struct ***REMOVED***
	Subsystem string
***REMOVED***

func (e *NotFoundError) Error() string ***REMOVED***
	return fmt.Sprintf("mountpoint for %s not found", e.Subsystem)
***REMOVED***

func NewNotFoundError(sub string) error ***REMOVED***
	return &NotFoundError***REMOVED***
		Subsystem: sub,
	***REMOVED***
***REMOVED***

func IsNotFound(err error) bool ***REMOVED***
	if err == nil ***REMOVED***
		return false
	***REMOVED***
	_, ok := err.(*NotFoundError)
	return ok
***REMOVED***
