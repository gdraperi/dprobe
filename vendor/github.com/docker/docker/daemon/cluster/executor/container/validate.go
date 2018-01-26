package container

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/docker/swarmkit/api"
)

func validateMounts(mounts []api.Mount) error ***REMOVED***
	for _, mount := range mounts ***REMOVED***
		// Target must always be absolute
		if !filepath.IsAbs(mount.Target) ***REMOVED***
			return fmt.Errorf("invalid mount target, must be an absolute path: %s", mount.Target)
		***REMOVED***

		switch mount.Type ***REMOVED***
		// The checks on abs paths are required due to the container API confusing
		// volume mounts as bind mounts when the source is absolute (and vice-versa)
		// See #25253
		// TODO: This is probably not necessary once #22373 is merged
		case api.MountTypeBind:
			if !filepath.IsAbs(mount.Source) ***REMOVED***
				return fmt.Errorf("invalid bind mount source, must be an absolute path: %s", mount.Source)
			***REMOVED***
		case api.MountTypeVolume:
			if filepath.IsAbs(mount.Source) ***REMOVED***
				return fmt.Errorf("invalid volume mount source, must not be an absolute path: %s", mount.Source)
			***REMOVED***
		case api.MountTypeTmpfs:
			if mount.Source != "" ***REMOVED***
				return errors.New("invalid tmpfs source, source must be empty")
			***REMOVED***
		default:
			return fmt.Errorf("invalid mount type: %s", mount.Type)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
