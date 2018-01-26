package dockerfile

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/symlink"
	lcUser "github.com/opencontainers/runc/libcontainer/user"
	"github.com/pkg/errors"
)

func parseChownFlag(chown, ctrRootPath string, idMappings *idtools.IDMappings) (idtools.IDPair, error) ***REMOVED***
	var userStr, grpStr string
	parts := strings.Split(chown, ":")
	if len(parts) > 2 ***REMOVED***
		return idtools.IDPair***REMOVED******REMOVED***, errors.New("invalid chown string format: " + chown)
	***REMOVED***
	if len(parts) == 1 ***REMOVED***
		// if no group specified, use the user spec as group as well
		userStr, grpStr = parts[0], parts[0]
	***REMOVED*** else ***REMOVED***
		userStr, grpStr = parts[0], parts[1]
	***REMOVED***

	passwdPath, err := symlink.FollowSymlinkInScope(filepath.Join(ctrRootPath, "etc", "passwd"), ctrRootPath)
	if err != nil ***REMOVED***
		return idtools.IDPair***REMOVED******REMOVED***, errors.Wrapf(err, "can't resolve /etc/passwd path in container rootfs")
	***REMOVED***
	groupPath, err := symlink.FollowSymlinkInScope(filepath.Join(ctrRootPath, "etc", "group"), ctrRootPath)
	if err != nil ***REMOVED***
		return idtools.IDPair***REMOVED******REMOVED***, errors.Wrapf(err, "can't resolve /etc/group path in container rootfs")
	***REMOVED***
	uid, err := lookupUser(userStr, passwdPath)
	if err != nil ***REMOVED***
		return idtools.IDPair***REMOVED******REMOVED***, errors.Wrapf(err, "can't find uid for user "+userStr)
	***REMOVED***
	gid, err := lookupGroup(grpStr, groupPath)
	if err != nil ***REMOVED***
		return idtools.IDPair***REMOVED******REMOVED***, errors.Wrapf(err, "can't find gid for group "+grpStr)
	***REMOVED***

	// convert as necessary because of user namespaces
	chownPair, err := idMappings.ToHost(idtools.IDPair***REMOVED***UID: uid, GID: gid***REMOVED***)
	if err != nil ***REMOVED***
		return idtools.IDPair***REMOVED******REMOVED***, errors.Wrapf(err, "unable to convert uid/gid to host mapping")
	***REMOVED***
	return chownPair, nil
***REMOVED***

func lookupUser(userStr, filepath string) (int, error) ***REMOVED***
	// if the string is actually a uid integer, parse to int and return
	// as we don't need to translate with the help of files
	uid, err := strconv.Atoi(userStr)
	if err == nil ***REMOVED***
		return uid, nil
	***REMOVED***
	users, err := lcUser.ParsePasswdFileFilter(filepath, func(u lcUser.User) bool ***REMOVED***
		return u.Name == userStr
	***REMOVED***)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if len(users) == 0 ***REMOVED***
		return 0, errors.New("no such user: " + userStr)
	***REMOVED***
	return users[0].Uid, nil
***REMOVED***

func lookupGroup(groupStr, filepath string) (int, error) ***REMOVED***
	// if the string is actually a gid integer, parse to int and return
	// as we don't need to translate with the help of files
	gid, err := strconv.Atoi(groupStr)
	if err == nil ***REMOVED***
		return gid, nil
	***REMOVED***
	groups, err := lcUser.ParseGroupFileFilter(filepath, func(g lcUser.Group) bool ***REMOVED***
		return g.Name == groupStr
	***REMOVED***)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if len(groups) == 0 ***REMOVED***
		return 0, errors.New("no such group: " + groupStr)
	***REMOVED***
	return groups[0].Gid, nil
***REMOVED***
