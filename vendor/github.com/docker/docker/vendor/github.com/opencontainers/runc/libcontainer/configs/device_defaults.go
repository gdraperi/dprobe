// +build linux

package configs

var (
	// DefaultSimpleDevices are devices that are to be both allowed and created.
	DefaultSimpleDevices = []*Device***REMOVED***
		// /dev/null and zero
		***REMOVED***
			Path:        "/dev/null",
			Type:        'c',
			Major:       1,
			Minor:       3,
			Permissions: "rwm",
			FileMode:    0666,
		***REMOVED***,
		***REMOVED***
			Path:        "/dev/zero",
			Type:        'c',
			Major:       1,
			Minor:       5,
			Permissions: "rwm",
			FileMode:    0666,
		***REMOVED***,

		***REMOVED***
			Path:        "/dev/full",
			Type:        'c',
			Major:       1,
			Minor:       7,
			Permissions: "rwm",
			FileMode:    0666,
		***REMOVED***,

		// consoles and ttys
		***REMOVED***
			Path:        "/dev/tty",
			Type:        'c',
			Major:       5,
			Minor:       0,
			Permissions: "rwm",
			FileMode:    0666,
		***REMOVED***,

		// /dev/urandom,/dev/random
		***REMOVED***
			Path:        "/dev/urandom",
			Type:        'c',
			Major:       1,
			Minor:       9,
			Permissions: "rwm",
			FileMode:    0666,
		***REMOVED***,
		***REMOVED***
			Path:        "/dev/random",
			Type:        'c',
			Major:       1,
			Minor:       8,
			Permissions: "rwm",
			FileMode:    0666,
		***REMOVED***,
	***REMOVED***
	DefaultAllowedDevices = append([]*Device***REMOVED***
		// allow mknod for any device
		***REMOVED***
			Type:        'c',
			Major:       Wildcard,
			Minor:       Wildcard,
			Permissions: "m",
		***REMOVED***,
		***REMOVED***
			Type:        'b',
			Major:       Wildcard,
			Minor:       Wildcard,
			Permissions: "m",
		***REMOVED***,

		***REMOVED***
			Path:        "/dev/console",
			Type:        'c',
			Major:       5,
			Minor:       1,
			Permissions: "rwm",
		***REMOVED***,
		// /dev/pts/ - pts namespaces are "coming soon"
		***REMOVED***
			Path:        "",
			Type:        'c',
			Major:       136,
			Minor:       Wildcard,
			Permissions: "rwm",
		***REMOVED***,
		***REMOVED***
			Path:        "",
			Type:        'c',
			Major:       5,
			Minor:       2,
			Permissions: "rwm",
		***REMOVED***,

		// tuntap
		***REMOVED***
			Path:        "",
			Type:        'c',
			Major:       10,
			Minor:       200,
			Permissions: "rwm",
		***REMOVED***,
	***REMOVED***, DefaultSimpleDevices...)
	DefaultAutoCreatedDevices = append([]*Device***REMOVED******REMOVED***, DefaultSimpleDevices...)
)
