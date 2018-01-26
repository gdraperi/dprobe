package mount

// MakeShared ensures a mounted filesystem has the SHARED mount option enabled.
// See the supported options in flags.go for further reference.
func MakeShared(mountPoint string) error ***REMOVED***
	return ensureMountedAs(mountPoint, "shared")
***REMOVED***

// MakeRShared ensures a mounted filesystem has the RSHARED mount option enabled.
// See the supported options in flags.go for further reference.
func MakeRShared(mountPoint string) error ***REMOVED***
	return ensureMountedAs(mountPoint, "rshared")
***REMOVED***

// MakePrivate ensures a mounted filesystem has the PRIVATE mount option enabled.
// See the supported options in flags.go for further reference.
func MakePrivate(mountPoint string) error ***REMOVED***
	return ensureMountedAs(mountPoint, "private")
***REMOVED***

// MakeRPrivate ensures a mounted filesystem has the RPRIVATE mount option
// enabled. See the supported options in flags.go for further reference.
func MakeRPrivate(mountPoint string) error ***REMOVED***
	return ensureMountedAs(mountPoint, "rprivate")
***REMOVED***

// MakeSlave ensures a mounted filesystem has the SLAVE mount option enabled.
// See the supported options in flags.go for further reference.
func MakeSlave(mountPoint string) error ***REMOVED***
	return ensureMountedAs(mountPoint, "slave")
***REMOVED***

// MakeRSlave ensures a mounted filesystem has the RSLAVE mount option enabled.
// See the supported options in flags.go for further reference.
func MakeRSlave(mountPoint string) error ***REMOVED***
	return ensureMountedAs(mountPoint, "rslave")
***REMOVED***

// MakeUnbindable ensures a mounted filesystem has the UNBINDABLE mount option
// enabled. See the supported options in flags.go for further reference.
func MakeUnbindable(mountPoint string) error ***REMOVED***
	return ensureMountedAs(mountPoint, "unbindable")
***REMOVED***

// MakeRUnbindable ensures a mounted filesystem has the RUNBINDABLE mount
// option enabled. See the supported options in flags.go for further reference.
func MakeRUnbindable(mountPoint string) error ***REMOVED***
	return ensureMountedAs(mountPoint, "runbindable")
***REMOVED***

func ensureMountedAs(mountPoint, options string) error ***REMOVED***
	mounted, err := Mounted(mountPoint)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !mounted ***REMOVED***
		if err := Mount(mountPoint, mountPoint, "none", "bind,rw"); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if _, err = Mounted(mountPoint); err != nil ***REMOVED***
		return err
	***REMOVED***

	return ForceMount("", mountPoint, "none", options)
***REMOVED***
