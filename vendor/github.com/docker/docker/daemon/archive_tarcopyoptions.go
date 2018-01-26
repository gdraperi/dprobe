package daemon

import (
	"github.com/docker/docker/pkg/archive"
)

// defaultTarCopyOptions is the setting that is used when unpacking an archive
// for a copy API event.
func (daemon *Daemon) defaultTarCopyOptions(noOverwriteDirNonDir bool) *archive.TarOptions ***REMOVED***
	return &archive.TarOptions***REMOVED***
		NoOverwriteDirNonDir: noOverwriteDirNonDir,
		UIDMaps:              daemon.idMappings.UIDs(),
		GIDMaps:              daemon.idMappings.GIDs(),
	***REMOVED***
***REMOVED***
