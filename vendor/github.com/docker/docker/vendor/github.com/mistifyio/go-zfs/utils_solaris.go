// +build solaris

package zfs

import (
	"strings"
)

// List of ZFS properties to retrieve from zfs list command on a Solaris platform
var dsPropList = []string***REMOVED***"name", "origin", "used", "available", "mountpoint", "compression", "type", "volsize", "quota"***REMOVED***

var dsPropListOptions = strings.Join(dsPropList, ",")

// List of Zpool properties to retrieve from zpool list command on a non-Solaris platform
var zpoolPropList = []string***REMOVED***"name", "health", "allocated", "size", "free", "readonly", "dedupratio"***REMOVED***
var zpoolPropListOptions = strings.Join(zpoolPropList, ",")
var zpoolArgs = []string***REMOVED***"get", "-p", zpoolPropListOptions***REMOVED***
