package daemon

import (
	"fmt"
	"os"
	"unsafe"

	winio "github.com/Microsoft/go-winio"
	"github.com/docker/docker/pkg/signal"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

func (d *Daemon) setupDumpStackTrap(root string) ***REMOVED***
	// Windows does not support signals like *nix systems. So instead of
	// trapping on SIGUSR1 to dump stacks, we wait on a Win32 event to be
	// signaled. ACL'd to builtin administrators and local system
	event := "Global\\docker-daemon-" + fmt.Sprint(os.Getpid())
	ev, _ := windows.UTF16PtrFromString(event)
	sd, err := winio.SddlToSecurityDescriptor("D:P(A;;GA;;;BA)(A;;GA;;;SY)")
	if err != nil ***REMOVED***
		logrus.Errorf("failed to get security descriptor for debug stackdump event %s: %s", event, err.Error())
		return
	***REMOVED***
	var sa windows.SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	sa.SecurityDescriptor = uintptr(unsafe.Pointer(&sd[0]))
	h, err := windows.CreateEvent(&sa, 0, 0, ev)
	if h == 0 || err != nil ***REMOVED***
		logrus.Errorf("failed to create debug stackdump event %s: %s", event, err.Error())
		return
	***REMOVED***
	go func() ***REMOVED***
		logrus.Debugf("Stackdump - waiting signal at %s", event)
		for ***REMOVED***
			windows.WaitForSingleObject(h, windows.INFINITE)
			path, err := signal.DumpStacks(root)
			if err != nil ***REMOVED***
				logrus.WithError(err).Error("failed to write goroutines dump")
			***REMOVED*** else ***REMOVED***
				logrus.Infof("goroutine stacks written to %s", path)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***
