// +build !windows

package oci

import (
	"context"
	"path/filepath"

	"github.com/containerd/containerd/namespaces"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const (
	rwm               = "rwm"
	defaultRootfsPath = "rootfs"
)

var (
	defaultEnv = []string***REMOVED***
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	***REMOVED***
)

func defaultCaps() []string ***REMOVED***
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

func defaultNamespaces() []specs.LinuxNamespace ***REMOVED***
	return []specs.LinuxNamespace***REMOVED***
		***REMOVED***
			Type: specs.PIDNamespace,
		***REMOVED***,
		***REMOVED***
			Type: specs.IPCNamespace,
		***REMOVED***,
		***REMOVED***
			Type: specs.UTSNamespace,
		***REMOVED***,
		***REMOVED***
			Type: specs.MountNamespace,
		***REMOVED***,
		***REMOVED***
			Type: specs.NetworkNamespace,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func createDefaultSpec(ctx context.Context, id string) (*specs.Spec, error) ***REMOVED***
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	s := &specs.Spec***REMOVED***
		Version: specs.Version,
		Root: &specs.Root***REMOVED***
			Path: defaultRootfsPath,
		***REMOVED***,
		Process: &specs.Process***REMOVED***
			Env:             defaultEnv,
			Cwd:             "/",
			NoNewPrivileges: true,
			User: specs.User***REMOVED***
				UID: 0,
				GID: 0,
			***REMOVED***,
			Capabilities: &specs.LinuxCapabilities***REMOVED***
				Bounding:    defaultCaps(),
				Permitted:   defaultCaps(),
				Inheritable: defaultCaps(),
				Effective:   defaultCaps(),
			***REMOVED***,
			Rlimits: []specs.POSIXRlimit***REMOVED***
				***REMOVED***
					Type: "RLIMIT_NOFILE",
					Hard: uint64(1024),
					Soft: uint64(1024),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Mounts: []specs.Mount***REMOVED***
			***REMOVED***
				Destination: "/proc",
				Type:        "proc",
				Source:      "proc",
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
				Destination: "/dev/shm",
				Type:        "tmpfs",
				Source:      "shm",
				Options:     []string***REMOVED***"nosuid", "noexec", "nodev", "mode=1777", "size=65536k"***REMOVED***,
			***REMOVED***,
			***REMOVED***
				Destination: "/dev/mqueue",
				Type:        "mqueue",
				Source:      "mqueue",
				Options:     []string***REMOVED***"nosuid", "noexec", "nodev"***REMOVED***,
			***REMOVED***,
			***REMOVED***
				Destination: "/sys",
				Type:        "sysfs",
				Source:      "sysfs",
				Options:     []string***REMOVED***"nosuid", "noexec", "nodev", "ro"***REMOVED***,
			***REMOVED***,
			***REMOVED***
				Destination: "/run",
				Type:        "tmpfs",
				Source:      "tmpfs",
				Options:     []string***REMOVED***"nosuid", "strictatime", "mode=755", "size=65536k"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Linux: &specs.Linux***REMOVED***
			MaskedPaths: []string***REMOVED***
				"/proc/kcore",
				"/proc/latency_stats",
				"/proc/timer_list",
				"/proc/timer_stats",
				"/proc/sched_debug",
				"/sys/firmware",
				"/proc/scsi",
			***REMOVED***,
			ReadonlyPaths: []string***REMOVED***
				"/proc/asound",
				"/proc/bus",
				"/proc/fs",
				"/proc/irq",
				"/proc/sys",
				"/proc/sysrq-trigger",
			***REMOVED***,
			CgroupsPath: filepath.Join("/", ns, id),
			Resources: &specs.LinuxResources***REMOVED***
				Devices: []specs.LinuxDeviceCgroup***REMOVED***
					***REMOVED***
						Allow:  false,
						Access: rwm,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			Namespaces: defaultNamespaces(),
		***REMOVED***,
	***REMOVED***
	return s, nil
***REMOVED***
