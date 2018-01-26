// +build windows

package client

import (
	"fmt"

	"github.com/Microsoft/hcsshim"
	"github.com/sirupsen/logrus"
)

// HotRemoveVhd hot-removes a VHD from a utility VM. This is used in the global one-utility-VM-
// service-VM per host scenario.
func (config *Config) HotRemoveVhd(hostPath string) error ***REMOVED***
	logrus.Debugf("opengcs: HotRemoveVhd: %s", hostPath)

	if config.Uvm == nil ***REMOVED***
		return fmt.Errorf("cannot hot-add VHD as no utility VM is in configuration")
	***REMOVED***

	defer config.DebugGCS()

	modification := &hcsshim.ResourceModificationRequestResponse***REMOVED***
		Resource: "MappedVirtualDisk",
		Data: hcsshim.MappedVirtualDisk***REMOVED***
			HostPath:          hostPath,
			CreateInUtilityVM: true,
		***REMOVED***,
		Request: "Remove",
	***REMOVED***
	if err := config.Uvm.Modify(modification); err != nil ***REMOVED***
		return fmt.Errorf("failed modifying utility VM for hot-remove %s: %s", hostPath, err)
	***REMOVED***
	logrus.Debugf("opengcs: HotRemoveVhd: %s removed successfully", hostPath)
	return nil
***REMOVED***
