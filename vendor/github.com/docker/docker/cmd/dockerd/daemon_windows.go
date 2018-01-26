package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/docker/docker/libcontainerd"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

var defaultDaemonConfigFile = ""

// setDefaultUmask doesn't do anything on windows
func setDefaultUmask() error ***REMOVED***
	return nil
***REMOVED***

func getDaemonConfDir(root string) string ***REMOVED***
	return filepath.Join(root, `\config`)
***REMOVED***

// preNotifySystem sends a message to the host when the API is active, but before the daemon is
func preNotifySystem() ***REMOVED***
	// start the service now to prevent timeouts waiting for daemon to start
	// but still (eventually) complete all requests that are sent after this
	if service != nil ***REMOVED***
		err := service.started()
		if err != nil ***REMOVED***
			logrus.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// notifySystem sends a message to the host when the server is ready to be used
func notifySystem() ***REMOVED***
***REMOVED***

// notifyShutdown is called after the daemon shuts down but before the process exits.
func notifyShutdown(err error) ***REMOVED***
	if service != nil ***REMOVED***
		if err != nil ***REMOVED***
			logrus.Fatal(err)
		***REMOVED***
		service.stopped(err)
	***REMOVED***
***REMOVED***

func (cli *DaemonCli) getPlatformRemoteOptions() ([]libcontainerd.RemoteOption, error) ***REMOVED***
	return nil, nil
***REMOVED***

// setupConfigReloadTrap configures a Win32 event to reload the configuration.
func (cli *DaemonCli) setupConfigReloadTrap() ***REMOVED***
	go func() ***REMOVED***
		sa := windows.SecurityAttributes***REMOVED***
			Length: 0,
		***REMOVED***
		event := "Global\\docker-daemon-config-" + fmt.Sprint(os.Getpid())
		ev, _ := windows.UTF16PtrFromString(event)
		if h, _ := windows.CreateEvent(&sa, 0, 0, ev); h != 0 ***REMOVED***
			logrus.Debugf("Config reload - waiting signal at %s", event)
			for ***REMOVED***
				windows.WaitForSingleObject(h, windows.INFINITE)
				cli.reloadConfig()
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***

// getSwarmRunRoot gets the root directory for swarm to store runtime state
// For example, the control socket
func (cli *DaemonCli) getSwarmRunRoot() string ***REMOVED***
	return ""
***REMOVED***

func allocateDaemonPort(addr string) error ***REMOVED***
	return nil
***REMOVED***

func wrapListeners(proto string, ls []net.Listener) []net.Listener ***REMOVED***
	return ls
***REMOVED***
