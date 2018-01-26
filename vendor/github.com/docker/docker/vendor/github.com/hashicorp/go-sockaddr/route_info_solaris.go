package sockaddr

import (
	"errors"
	"os/exec"
)

var cmds map[string][]string = map[string][]string***REMOVED***
	"route": ***REMOVED***"/usr/sbin/route", "-n", "get", "default"***REMOVED***,
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
	out, err := exec.Command(cmds["route"][0], cmds["route"][1:]...).Output()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var ifName string
	if ifName, err = parseDefaultIfNameFromRoute(string(out)); err != nil ***REMOVED***
		return "", errors.New("No default interface found")
	***REMOVED***
	return ifName, nil
***REMOVED***
