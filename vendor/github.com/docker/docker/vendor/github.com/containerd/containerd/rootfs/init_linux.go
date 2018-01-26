package rootfs

import (
	"os"
	"path/filepath"
	"syscall"
)

const (
	defaultInitializer = "linux-init"
)

func init() ***REMOVED***
	initializers[defaultInitializer] = initFS
***REMOVED***

func createDirectory(name string, uid, gid int) initializerFunc ***REMOVED***
	return func(root string) error ***REMOVED***
		dname := filepath.Join(root, name)
		st, err := os.Stat(dname)
		if err != nil && !os.IsNotExist(err) ***REMOVED***
			return err
		***REMOVED*** else if err == nil ***REMOVED***
			if st.IsDir() ***REMOVED***
				stat := st.Sys().(*syscall.Stat_t)
				if int(stat.Gid) == gid && int(stat.Uid) == uid ***REMOVED***
					return nil
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if err := os.Remove(dname); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := os.Mkdir(dname, 0755); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err := os.Mkdir(dname, 0755); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		return os.Chown(dname, uid, gid)
	***REMOVED***
***REMOVED***

func touchFile(name string, uid, gid int) initializerFunc ***REMOVED***
	return func(root string) error ***REMOVED***
		fname := filepath.Join(root, name)

		st, err := os.Stat(fname)
		if err != nil && !os.IsNotExist(err) ***REMOVED***
			return err
		***REMOVED*** else if err == nil ***REMOVED***
			stat := st.Sys().(*syscall.Stat_t)
			if int(stat.Gid) == gid && int(stat.Uid) == uid ***REMOVED***
				return nil
			***REMOVED***
			return os.Chown(fname, uid, gid)
		***REMOVED***

		f, err := os.OpenFile(fname, os.O_CREATE, 0644)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer f.Close()

		return f.Chown(uid, gid)
	***REMOVED***
***REMOVED***

func symlink(oldname, newname string) initializerFunc ***REMOVED***
	return func(root string) error ***REMOVED***
		linkName := filepath.Join(root, newname)
		if _, err := os.Stat(linkName); err != nil && !os.IsNotExist(err) ***REMOVED***
			return err
		***REMOVED*** else if err == nil ***REMOVED***
			return nil
		***REMOVED***
		return os.Symlink(oldname, linkName)
	***REMOVED***
***REMOVED***

func initFS(root string) error ***REMOVED***
	st, err := os.Stat(root)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	stat := st.Sys().(*syscall.Stat_t)
	uid := int(stat.Uid)
	gid := int(stat.Gid)

	initFuncs := []initializerFunc***REMOVED***
		createDirectory("/dev", uid, gid),
		createDirectory("/dev/pts", uid, gid),
		createDirectory("/dev/shm", uid, gid),
		touchFile("/dev/console", uid, gid),
		createDirectory("/proc", uid, gid),
		createDirectory("/sys", uid, gid),
		createDirectory("/etc", uid, gid),
		touchFile("/etc/resolv.conf", uid, gid),
		touchFile("/etc/hosts", uid, gid),
		touchFile("/etc/hostname", uid, gid),
		symlink("/proc/mounts", "/etc/mtab"),
	***REMOVED***

	for _, fn := range initFuncs ***REMOVED***
		if err := fn(root); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
