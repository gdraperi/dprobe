package mount

// Mount is the lingua franca of containerd. A mount represents a
// serialized mount syscall. Components either emit or consume mounts.
type Mount struct ***REMOVED***
	// Type specifies the host-specific of the mount.
	Type string
	// Source specifies where to mount from. Depending on the host system, this
	// can be a source path or device.
	Source string
	// Options contains zero or more fstab-style mount options. Typically,
	// these are platform specific.
	Options []string
***REMOVED***

// All mounts all the provided mounts to the provided target
func All(mounts []Mount, target string) error ***REMOVED***
	for _, m := range mounts ***REMOVED***
		if err := m.Mount(target); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
