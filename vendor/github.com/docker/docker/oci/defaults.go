package oci

import (
	"os"
	"runtime"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func iPtr(i int64) *int64        ***REMOVED*** return &i ***REMOVED***
func u32Ptr(i int64) *uint32     ***REMOVED*** u := uint32(i); return &u ***REMOVED***
func fmPtr(i int64) *os.FileMode ***REMOVED*** fm := os.FileMode(i); return &fm ***REMOVED***

func defaultCapabilities() []string ***REMOVED***
	return []string***REMOVED***
		"CAP_CHOWN",
		"CAP_DAC_OVERRIDE",
		"CAP_FSETID",
		"CAP_FOWNER",
		"CAP_MKNOD",
		"CAP_NET_RAW",
		"CAP_SETGID",
		"CAP_SETUID",
		"CAP_SETFCAP",
		"CAP_SETPCAP",
		"CAP_NET_BIND_SERVICE",
		"CAP_SYS_CHROOT",
		"CAP_KILL",
		"CAP_AUDIT_WRITE",
	***REMOVED***
***REMOVED***

// DefaultSpec returns the default spec used by docker for the current Platform
func DefaultSpec() specs.Spec ***REMOVED***
	return DefaultOSSpec(runtime.GOOS)
***REMOVED***

// DefaultOSSpec returns the spec for a given OS
func DefaultOSSpec(osName string) specs.Spec ***REMOVED***
	if osName == "windows" ***REMOVED***
		return DefaultWindowsSpec()
	***REMOVED***
	return DefaultLinuxSpec()
***REMOVED***

// DefaultWindowsSpec create a default spec for running Windows containers
func DefaultWindowsSpec() specs.Spec ***REMOVED***
	return specs.Spec***REMOVED***
		Version: specs.Version,
		Windows: &specs.Windows***REMOVED******REMOVED***,
		Process: &specs.Process***REMOVED******REMOVED***,
		Root:    &specs.Root***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// DefaultLinuxSpec create a default spec for running Linux containers
func DefaultLinuxSpec() specs.Spec ***REMOVED***
	s := specs.Spec***REMOVED***
		Version: specs.Version,
		Process: &specs.Process***REMOVED***
			Capabilities: &specs.LinuxCapabilities***REMOVED***
				Bounding:    defaultCapabilities(),
				Permitted:   defaultCapabilities(),
				Inheritable: defaultCapabilities(),
				Effective:   defaultCapabilities(),
			***REMOVED***,
		***REMOVED***,
		Root: &specs.Root***REMOVED******REMOVED***,
	***REMOVED***
	s.Mounts = []specs.Mount***REMOVED***
		***REMOVED***
			Destination: "/proc",
			Type:        "proc",
			Source:      "proc",
			Options:     []string***REMOVED***"nosuid", "noexec", "nodev"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Destination: "/dev",
			Type:        "tmpfs",
			Source:      "tmpfs",
			Options:     []string***REMOVED***"nosuid", "strictatime", "mode=755", "size=65536k"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Destination: "/dev/pts",
			Type:        "devpts",
			Source:      "devpts",
			Options:     []string***REMOVED***"nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620", "gid=5"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Destination: "/sys",
			Type:        "sysfs",
			Source:      "sysfs",
			Options:     []string***REMOVED***"nosuid", "noexec", "nodev", "ro"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Destination: "/sys/fs/cgroup",
			Type:        "cgroup",
			Source:      "cgroup",
			Options:     []string***REMOVED***"ro", "nosuid", "noexec", "nodev"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Destination: "/dev/mqueue",
			Type:        "mqueue",
			Source:      "mqueue",
			Options:     []string***REMOVED***"nosuid", "noexec", "nodev"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Destination: "/dev/shm",
			Type:        "tmpfs",
			Source:      "shm",
			Options:     []string***REMOVED***"nosuid", "noexec", "nodev", "mode=1777"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	s.Linux = &specs.Linux***REMOVED***
		MaskedPaths: []string***REMOVED***
			"/proc/kcore",
			"/proc/latency_stats",
			"/proc/timer_list",
			"/proc/timer_stats",
			"/proc/sched_debug",
			"/proc/scsi",
			"/sys/firmware",
		***REMOVED***,
		ReadonlyPaths: []string***REMOVED***
			"/proc/asound",
			"/proc/bus",
			"/proc/fs",
			"/proc/irq",
			"/proc/sys",
			"/proc/sysrq-trigger",
		***REMOVED***,
		Namespaces: []specs.LinuxNamespace***REMOVED***
			***REMOVED***Type: "mount"***REMOVED***,
			***REMOVED***Type: "network"***REMOVED***,
			***REMOVED***Type: "uts"***REMOVED***,
			***REMOVED***Type: "pid"***REMOVED***,
			***REMOVED***Type: "ipc"***REMOVED***,
		***REMOVED***,
		// Devices implicitly contains the following devices:
		// null, zero, full, random, urandom, tty, console, and ptmx.
		// ptmx is a bind mount or symlink of the container's ptmx.
		// See also: https://github.com/opencontainers/runtime-spec/blob/master/config-linux.md#default-devices
		Devices: []specs.LinuxDevice***REMOVED******REMOVED***,
		Resources: &specs.LinuxResources***REMOVED***
			Devices: []specs.LinuxDeviceCgroup***REMOVED***
				***REMOVED***
					Allow:  false,
					Access: "rwm",
				***REMOVED***,
				***REMOVED***
					Allow:  true,
					Type:   "c",
					Major:  iPtr(1),
					Minor:  iPtr(5),
					Access: "rwm",
				***REMOVED***,
				***REMOVED***
					Allow:  true,
					Type:   "c",
					Major:  iPtr(1),
					Minor:  iPtr(3),
					Access: "rwm",
				***REMOVED***,
				***REMOVED***
					Allow:  true,
					Type:   "c",
					Major:  iPtr(1),
					Minor:  iPtr(9),
					Access: "rwm",
				***REMOVED***,
				***REMOVED***
					Allow:  true,
					Type:   "c",
					Major:  iPtr(1),
					Minor:  iPtr(8),
					Access: "rwm",
				***REMOVED***,
				***REMOVED***
					Allow:  true,
					Type:   "c",
					Major:  iPtr(5),
					Minor:  iPtr(0),
					Access: "rwm",
				***REMOVED***,
				***REMOVED***
					Allow:  true,
					Type:   "c",
					Major:  iPtr(5),
					Minor:  iPtr(1),
					Access: "rwm",
				***REMOVED***,
				***REMOVED***
					Allow:  false,
					Type:   "c",
					Major:  iPtr(10),
					Minor:  iPtr(229),
					Access: "rwm",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	// For LCOW support, populate a blank Windows spec
	if runtime.GOOS == "windows" ***REMOVED***
		s.Windows = &specs.Windows***REMOVED******REMOVED***
	***REMOVED***

	return s
***REMOVED***
