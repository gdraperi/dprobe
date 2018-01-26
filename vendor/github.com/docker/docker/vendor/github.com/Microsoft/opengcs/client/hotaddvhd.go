// +build windows

package client

import (
	"fmt"

	"github.com/Microsoft/hcsshim"
	"github.com/sirupsen/logrus"
)

// HotAddVhd hot-adds a VHD to a utility VM. This is used in the global one-utility-VM-
// service-VM per host scenario. In order to do a graphdriver `Diff`, we hot-add the
// sandbox to /mnt/<id> so that we can run `exportSandbox` inside the utility VM to
// get a tar-stream of the sandboxes contents back to the daemon.
func (config *Config) HotAddVhd(hostPath string, containerPath string, readOnly bool, mount bool) error ***REMOVED***
	logrus.Debugf("opengcs: HotAddVhd: %s: %s", hostPath, containerPath)

	if config.Uvm == nil ***REMOVED***
		return fmt.Errorf("cannot hot-add VHD as no utility VM is in configuration")
	***REMOVED***

	defer config.DebugGCS()

	modification := &hcsshim.ResourceModificationRequestResponse***REMOVED***
		Resource: "MappedVirtualDisk",
		Data: hcsshim.MappedVirtualDisk***REMOVED***
			HostPath:          hostPath,
			ContainerPath:     containerPath,
			CreateInUtilityVM: true,
			ReadOnly:          readOnly,
			AttachOnly:        !mount,
		***REMOVED***,
		Request: "Add",
	***REMOVED***

	if err := config.Uvm.Modify(modification); err != nil ***REMOVED***
		return fmt.Errorf("failed to modify utility VM configuration for hot-add: %s", err)
	***REMOVED***
	logrus.Debugf("opengcs: HotAddVhd: %s added successfully", hostPath)
	return nil
***REMOVED***
