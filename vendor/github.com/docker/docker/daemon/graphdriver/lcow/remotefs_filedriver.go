// +build windows

package lcow

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"

	"github.com/Microsoft/opengcs/service/gcsutils/remotefs"

	"github.com/containerd/continuity/driver"
	"github.com/sirupsen/logrus"
)

var _ driver.Driver = &lcowfs***REMOVED******REMOVED***

func (l *lcowfs) Readlink(p string) (string, error) ***REMOVED***
	logrus.Debugf("removefs.readlink args: %s", p)

	result := &bytes.Buffer***REMOVED******REMOVED***
	if err := l.runRemoteFSProcess(nil, result, remotefs.ReadlinkCmd, p); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return result.String(), nil
***REMOVED***

func (l *lcowfs) Mkdir(path string, mode os.FileMode) error ***REMOVED***
	return l.mkdir(path, mode, remotefs.MkdirCmd)
***REMOVED***

func (l *lcowfs) MkdirAll(path string, mode os.FileMode) error ***REMOVED***
	return l.mkdir(path, mode, remotefs.MkdirAllCmd)
***REMOVED***

func (l *lcowfs) mkdir(path string, mode os.FileMode, cmd string) error ***REMOVED***
	modeStr := strconv.FormatUint(uint64(mode), 8)
	logrus.Debugf("remotefs.%s args: %s %s", cmd, path, modeStr)
	return l.runRemoteFSProcess(nil, nil, cmd, path, modeStr)
***REMOVED***

func (l *lcowfs) Remove(path string) error ***REMOVED***
	return l.remove(path, remotefs.RemoveCmd)
***REMOVED***

func (l *lcowfs) RemoveAll(path string) error ***REMOVED***
	return l.remove(path, remotefs.RemoveAllCmd)
***REMOVED***

func (l *lcowfs) remove(path string, cmd string) error ***REMOVED***
	logrus.Debugf("remotefs.%s args: %s", cmd, path)
	return l.runRemoteFSProcess(nil, nil, cmd, path)
***REMOVED***

func (l *lcowfs) Link(oldname, newname string) error ***REMOVED***
	return l.link(oldname, newname, remotefs.LinkCmd)
***REMOVED***

func (l *lcowfs) Symlink(oldname, newname string) error ***REMOVED***
	return l.link(oldname, newname, remotefs.SymlinkCmd)
***REMOVED***

func (l *lcowfs) link(oldname, newname, cmd string) error ***REMOVED***
	logrus.Debugf("remotefs.%s args: %s %s", cmd, oldname, newname)
	return l.runRemoteFSProcess(nil, nil, cmd, oldname, newname)
***REMOVED***

func (l *lcowfs) Lchown(name string, uid, gid int64) error ***REMOVED***
	uidStr := strconv.FormatInt(uid, 10)
	gidStr := strconv.FormatInt(gid, 10)

	logrus.Debugf("remotefs.lchown args: %s %s %s", name, uidStr, gidStr)
	return l.runRemoteFSProcess(nil, nil, remotefs.LchownCmd, name, uidStr, gidStr)
***REMOVED***

// Lchmod changes the mode of an file not following symlinks.
func (l *lcowfs) Lchmod(path string, mode os.FileMode) error ***REMOVED***
	modeStr := strconv.FormatUint(uint64(mode), 8)
	logrus.Debugf("remotefs.lchmod args: %s %s", path, modeStr)
	return l.runRemoteFSProcess(nil, nil, remotefs.LchmodCmd, path, modeStr)
***REMOVED***

func (l *lcowfs) Mknod(path string, mode os.FileMode, major, minor int) error ***REMOVED***
	modeStr := strconv.FormatUint(uint64(mode), 8)
	majorStr := strconv.FormatUint(uint64(major), 10)
	minorStr := strconv.FormatUint(uint64(minor), 10)

	logrus.Debugf("remotefs.mknod args: %s %s %s %s", path, modeStr, majorStr, minorStr)
	return l.runRemoteFSProcess(nil, nil, remotefs.MknodCmd, path, modeStr, majorStr, minorStr)
***REMOVED***

func (l *lcowfs) Mkfifo(path string, mode os.FileMode) error ***REMOVED***
	modeStr := strconv.FormatUint(uint64(mode), 8)
	logrus.Debugf("remotefs.mkfifo args: %s %s", path, modeStr)
	return l.runRemoteFSProcess(nil, nil, remotefs.MkfifoCmd, path, modeStr)
***REMOVED***

func (l *lcowfs) Stat(p string) (os.FileInfo, error) ***REMOVED***
	return l.stat(p, remotefs.StatCmd)
***REMOVED***

func (l *lcowfs) Lstat(p string) (os.FileInfo, error) ***REMOVED***
	return l.stat(p, remotefs.LstatCmd)
***REMOVED***

func (l *lcowfs) stat(path string, cmd string) (os.FileInfo, error) ***REMOVED***
	logrus.Debugf("remotefs.stat inputs: %s %s", cmd, path)

	output := &bytes.Buffer***REMOVED******REMOVED***
	err := l.runRemoteFSProcess(nil, output, cmd, path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var fi remotefs.FileInfo
	if err := json.Unmarshal(output.Bytes(), &fi); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	logrus.Debugf("remotefs.stat success. got: %v\n", fi)
	return &fi, nil
***REMOVED***
