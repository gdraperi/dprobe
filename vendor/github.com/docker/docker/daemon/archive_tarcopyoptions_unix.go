// +build !windows

package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
)

func (daemon *Daemon) tarCopyOptions(container *container.Container, noOverwriteDirNonDir bool) (*archive.TarOptions, error) ***REMOVED***
	if container.Config.User == "" ***REMOVED***
		return daemon.defaultTarCopyOptions(noOverwriteDirNonDir), nil
	***REMOVED***

	user, err := idtools.LookupUser(container.Config.User)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &archive.TarOptions***REMOVED***
		NoOverwriteDirNonDir: noOverwriteDirNonDir,
		ChownOpts:            &idtools.IDPair***REMOVED***UID: user.Uid, GID: user.Gid***REMOVED***,
	***REMOVED***, nil
***REMOVED***
