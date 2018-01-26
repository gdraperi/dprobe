package sockaddr

import "os/exec"

var cmds map[string][]string = map[string][]string***REMOVED***
	"netstat":  ***REMOVED***"netstat", "-rn"***REMOVED***,
	"ipconfig": ***REMOVED***"ipconfig"***REMOVED***,
***REMOVED***

type routeInfo struct ***REMOVED***
	cmds map[string][]string
***REMOVED***

// NewRouteInfo returns a BSD-specific implementation of the RouteInfo
// interface.
func NewRouteInfo() (routeInfo, error) ***REMOVED***
	return routeInfo***REMOVED***
		cmds: cmds,
	***REMOVED***, nil
***REMOVED***

// GetDefaultInterfaceName returns the interface name attached to the default
// route on the default interface.
func (ri routeInfo) GetDefaultInterfaceName() (string, error) ***REMOVED***
	ifNameOut, err := exec.Command(cmds["netstat"][0], cmds["netstat"][1:]...).Output()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	ipconfigOut, err := exec.Command(cmds["ipconfig"][0], cmds["ipconfig"][1:]...).Output()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	ifName, err := parseDefaultIfNameWindows(string(ifNameOut), string(ipconfigOut))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return ifName, nil
***REMOVED***
