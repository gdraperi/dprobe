package daemon

import (
	"fmt"
	"io"
	"runtime"

	"github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/system"
)

// ContainerExport writes the contents of the container to the given
// writer. An error is returned if the container cannot be found.
func (daemon *Daemon) ContainerExport(name string, out io.Writer) error ***REMOVED***
	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if runtime.GOOS == "windows" && container.OS == "windows" ***REMOVED***
		return fmt.Errorf("the daemon on this operating system does not support exporting Windows containers")
	***REMOVED***

	if container.IsDead() ***REMOVED***
		err := fmt.Errorf("You cannot export container %s which is Dead", container.ID)
		return errdefs.Conflict(err)
	***REMOVED***

	if container.IsRemovalInProgress() ***REMOVED***
		err := fmt.Errorf("You cannot export container %s which is being removed", container.ID)
		return errdefs.Conflict(err)
	***REMOVED***

	data, err := daemon.containerExport(container)
	if err != nil ***REMOVED***
		return fmt.Errorf("Error exporting container %s: %v", name, err)
	***REMOVED***
	defer data.Close()

	// Stream the entire contents of the container (basically a volatile snapshot)
	if _, err := io.Copy(out, data); err != nil ***REMOVED***
		return fmt.Errorf("Error exporting container %s: %v", name, err)
	***REMOVED***
	return nil
***REMOVED***

func (daemon *Daemon) containerExport(container *container.Container) (arch io.ReadCloser, err error) ***REMOVED***
	if !system.IsOSSupported(container.OS) ***REMOVED***
		return nil, fmt.Errorf("cannot export %s: %s ", container.ID, system.ErrNotSupportedOperatingSystem)
	***REMOVED***
	rwlayer, err := daemon.layerStores[container.OS].GetRWLayer(container.ID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			daemon.layerStores[container.OS].ReleaseRWLayer(rwlayer)
		***REMOVED***
	***REMOVED***()

	_, err = rwlayer.Mount(container.GetMountLabel())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	archive, err := archivePath(container.BaseFS, container.BaseFS.Path(), &archive.TarOptions***REMOVED***
		Compression: archive.Uncompressed,
		UIDMaps:     daemon.idMappings.UIDs(),
		GIDMaps:     daemon.idMappings.GIDs(),
	***REMOVED***)
	if err != nil ***REMOVED***
		rwlayer.Unmount()
		return nil, err
	***REMOVED***
	arch = ioutils.NewReadCloserWrapper(archive, func() error ***REMOVED***
		err := archive.Close()
		rwlayer.Unmount()
		daemon.layerStores[container.OS].ReleaseRWLayer(rwlayer)
		return err
	***REMOVED***)
	daemon.LogContainerEvent(container, "export")
	return arch, err
***REMOVED***
