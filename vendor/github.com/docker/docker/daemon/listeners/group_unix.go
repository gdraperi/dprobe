// +build !windows

package listeners

import (
	"fmt"
	"strconv"

	"github.com/opencontainers/runc/libcontainer/user"
	"github.com/pkg/errors"
)

const defaultSocketGroup = "docker"

func lookupGID(name string) (int, error) ***REMOVED***
	groupFile, err := user.GetGroupPath()
	if err != nil ***REMOVED***
		return -1, errors.Wrap(err, "error looking up groups")
	***REMOVED***
	groups, err := user.ParseGroupFileFilter(groupFile, func(g user.Group) bool ***REMOVED***
		return g.Name == name || strconv.Itoa(g.Gid) == name
	***REMOVED***)
	if err != nil ***REMOVED***
		return -1, errors.Wrapf(err, "error parsing groups for %s", name)
	***REMOVED***
	if len(groups) > 0 ***REMOVED***
		return groups[0].Gid, nil
	***REMOVED***
	gid, err := strconv.Atoi(name)
	if err == nil ***REMOVED***
		return gid, nil
	***REMOVED***
	return -1, fmt.Errorf("group %s not found", name)
***REMOVED***
