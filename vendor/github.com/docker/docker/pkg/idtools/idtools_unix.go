// +build !windows

package idtools

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/docker/docker/pkg/system"
	"github.com/opencontainers/runc/libcontainer/user"
)

var (
	entOnce   sync.Once
	getentCmd string
)

func mkdirAs(path string, mode os.FileMode, ownerUID, ownerGID int, mkAll, chownExisting bool) error ***REMOVED***
	// make an array containing the original path asked for, plus (for mkAll == true)
	// all path components leading up to the complete path that don't exist before we MkdirAll
	// so that we can chown all of them properly at the end.  If chownExisting is false, we won't
	// chown the full directory path if it exists
	var paths []string

	stat, err := system.Stat(path)
	if err == nil ***REMOVED***
		if !stat.IsDir() ***REMOVED***
			return &os.PathError***REMOVED***Op: "mkdir", Path: path, Err: syscall.ENOTDIR***REMOVED***
		***REMOVED***
		if !chownExisting ***REMOVED***
			return nil
		***REMOVED***

		// short-circuit--we were called with an existing directory and chown was requested
		return lazyChown(path, ownerUID, ownerGID, stat)
	***REMOVED***

	if os.IsNotExist(err) ***REMOVED***
		paths = []string***REMOVED***path***REMOVED***
	***REMOVED***

	if mkAll ***REMOVED***
		// walk back to "/" looking for directories which do not exist
		// and add them to the paths array for chown after creation
		dirPath := path
		for ***REMOVED***
			dirPath = filepath.Dir(dirPath)
			if dirPath == "/" ***REMOVED***
				break
			***REMOVED***
			if _, err := os.Stat(dirPath); err != nil && os.IsNotExist(err) ***REMOVED***
				paths = append(paths, dirPath)
			***REMOVED***
		***REMOVED***
		if err := system.MkdirAll(path, mode, ""); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := os.Mkdir(path, mode); err != nil && !os.IsExist(err) ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// even if it existed, we will chown the requested path + any subpaths that
	// didn't exist when we called MkdirAll
	for _, pathComponent := range paths ***REMOVED***
		if err := lazyChown(pathComponent, ownerUID, ownerGID, nil); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// CanAccess takes a valid (existing) directory and a uid, gid pair and determines
// if that uid, gid pair has access (execute bit) to the directory
func CanAccess(path string, pair IDPair) bool ***REMOVED***
	statInfo, err := system.Stat(path)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	fileMode := os.FileMode(statInfo.Mode())
	permBits := fileMode.Perm()
	return accessible(statInfo.UID() == uint32(pair.UID),
		statInfo.GID() == uint32(pair.GID), permBits)
***REMOVED***

func accessible(isOwner, isGroup bool, perms os.FileMode) bool ***REMOVED***
	if isOwner && (perms&0100 == 0100) ***REMOVED***
		return true
	***REMOVED***
	if isGroup && (perms&0010 == 0010) ***REMOVED***
		return true
	***REMOVED***
	if perms&0001 == 0001 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// LookupUser uses traditional local system files lookup (from libcontainer/user) on a username,
// followed by a call to `getent` for supporting host configured non-files passwd and group dbs
func LookupUser(username string) (user.User, error) ***REMOVED***
	// first try a local system files lookup using existing capabilities
	usr, err := user.LookupUser(username)
	if err == nil ***REMOVED***
		return usr, nil
	***REMOVED***
	// local files lookup failed; attempt to call `getent` to query configured passwd dbs
	usr, err = getentUser(fmt.Sprintf("%s %s", "passwd", username))
	if err != nil ***REMOVED***
		return user.User***REMOVED******REMOVED***, err
	***REMOVED***
	return usr, nil
***REMOVED***

// LookupUID uses traditional local system files lookup (from libcontainer/user) on a uid,
// followed by a call to `getent` for supporting host configured non-files passwd and group dbs
func LookupUID(uid int) (user.User, error) ***REMOVED***
	// first try a local system files lookup using existing capabilities
	usr, err := user.LookupUid(uid)
	if err == nil ***REMOVED***
		return usr, nil
	***REMOVED***
	// local files lookup failed; attempt to call `getent` to query configured passwd dbs
	return getentUser(fmt.Sprintf("%s %d", "passwd", uid))
***REMOVED***

func getentUser(args string) (user.User, error) ***REMOVED***
	reader, err := callGetent(args)
	if err != nil ***REMOVED***
		return user.User***REMOVED******REMOVED***, err
	***REMOVED***
	users, err := user.ParsePasswd(reader)
	if err != nil ***REMOVED***
		return user.User***REMOVED******REMOVED***, err
	***REMOVED***
	if len(users) == 0 ***REMOVED***
		return user.User***REMOVED******REMOVED***, fmt.Errorf("getent failed to find passwd entry for %q", strings.Split(args, " ")[1])
	***REMOVED***
	return users[0], nil
***REMOVED***

// LookupGroup uses traditional local system files lookup (from libcontainer/user) on a group name,
// followed by a call to `getent` for supporting host configured non-files passwd and group dbs
func LookupGroup(groupname string) (user.Group, error) ***REMOVED***
	// first try a local system files lookup using existing capabilities
	group, err := user.LookupGroup(groupname)
	if err == nil ***REMOVED***
		return group, nil
	***REMOVED***
	// local files lookup failed; attempt to call `getent` to query configured group dbs
	return getentGroup(fmt.Sprintf("%s %s", "group", groupname))
***REMOVED***

// LookupGID uses traditional local system files lookup (from libcontainer/user) on a group ID,
// followed by a call to `getent` for supporting host configured non-files passwd and group dbs
func LookupGID(gid int) (user.Group, error) ***REMOVED***
	// first try a local system files lookup using existing capabilities
	group, err := user.LookupGid(gid)
	if err == nil ***REMOVED***
		return group, nil
	***REMOVED***
	// local files lookup failed; attempt to call `getent` to query configured group dbs
	return getentGroup(fmt.Sprintf("%s %d", "group", gid))
***REMOVED***

func getentGroup(args string) (user.Group, error) ***REMOVED***
	reader, err := callGetent(args)
	if err != nil ***REMOVED***
		return user.Group***REMOVED******REMOVED***, err
	***REMOVED***
	groups, err := user.ParseGroup(reader)
	if err != nil ***REMOVED***
		return user.Group***REMOVED******REMOVED***, err
	***REMOVED***
	if len(groups) == 0 ***REMOVED***
		return user.Group***REMOVED******REMOVED***, fmt.Errorf("getent failed to find groups entry for %q", strings.Split(args, " ")[1])
	***REMOVED***
	return groups[0], nil
***REMOVED***

func callGetent(args string) (io.Reader, error) ***REMOVED***
	entOnce.Do(func() ***REMOVED*** getentCmd, _ = resolveBinary("getent") ***REMOVED***)
	// if no `getent` command on host, can't do anything else
	if getentCmd == "" ***REMOVED***
		return nil, fmt.Errorf("")
	***REMOVED***
	out, err := execCmd(getentCmd, args)
	if err != nil ***REMOVED***
		exitCode, errC := system.GetExitCode(err)
		if errC != nil ***REMOVED***
			return nil, err
		***REMOVED***
		switch exitCode ***REMOVED***
		case 1:
			return nil, fmt.Errorf("getent reported invalid parameters/database unknown")
		case 2:
			terms := strings.Split(args, " ")
			return nil, fmt.Errorf("getent unable to find entry %q in %s database", terms[1], terms[0])
		case 3:
			return nil, fmt.Errorf("getent database doesn't support enumeration")
		default:
			return nil, err
		***REMOVED***

	***REMOVED***
	return bytes.NewReader(out), nil
***REMOVED***

// lazyChown performs a chown only if the uid/gid don't match what's requested
// Normally a Chown is a no-op if uid/gid match, but in some cases this can still cause an error, e.g. if the
// dir is on an NFS share, so don't call chown unless we absolutely must.
func lazyChown(p string, uid, gid int, stat *system.StatT) error ***REMOVED***
	if stat == nil ***REMOVED***
		var err error
		stat, err = system.Stat(p)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if stat.UID() == uint32(uid) && stat.GID() == uint32(gid) ***REMOVED***
		return nil
	***REMOVED***
	return os.Chown(p, uid, gid)
***REMOVED***
