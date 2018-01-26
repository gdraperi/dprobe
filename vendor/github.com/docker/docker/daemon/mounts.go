package daemon

import (
	"fmt"
	"strings"

	mounttypes "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/container"
	volumestore "github.com/docker/docker/volume/store"
)

func (daemon *Daemon) prepareMountPoints(container *container.Container) error ***REMOVED***
	for _, config := range container.MountPoints ***REMOVED***
		if err := daemon.lazyInitializeVolume(container.ID, config); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) removeMountPoints(container *container.Container, rm bool) error ***REMOVED***
	var rmErrors []string
	for _, m := range container.MountPoints ***REMOVED***
		if m.Type != mounttypes.TypeVolume || m.Volume == nil ***REMOVED***
			continue
		***REMOVED***
		daemon.volumes.Dereference(m.Volume, container.ID)
		if !rm ***REMOVED***
			continue
		***REMOVED***

		// Do not remove named mountpoints
		// these are mountpoints specified like `docker run -v <name>:/foo`
		if m.Spec.Source != "" ***REMOVED***
			continue
		***REMOVED***

		err := daemon.volumes.Remove(m.Volume)
		// Ignore volume in use errors because having this
		// volume being referenced by other container is
		// not an error, but an implementation detail.
		// This prevents docker from logging "ERROR: Volume in use"
		// where there is another container using the volume.
		if err != nil && !volumestore.IsInUse(err) ***REMOVED***
			rmErrors = append(rmErrors, err.Error())
		***REMOVED***
	***REMOVED***

	if len(rmErrors) > 0 ***REMOVED***
		return fmt.Errorf("Error removing volumes:\n%v", strings.Join(rmErrors, "\n"))
	***REMOVED***
	return nil
***REMOVED***
