# libcontainer

[![GoDoc](https://godoc.org/github.com/opencontainers/runc/libcontainer?status.svg)](https://godoc.org/github.com/opencontainers/runc/libcontainer)

Libcontainer provides a native Go implementation for creating containers
with namespaces, cgroups, capabilities, and filesystem access controls.
It allows you to manage the lifecycle of the container performing additional operations
after the container is created.


#### Container
A container is a self contained execution environment that shares the kernel of the
host system and which is (optionally) isolated from other containers in the system.

#### Using libcontainer

Because containers are spawned in a two step process you will need a binary that
will be executed as the init process for the container. In libcontainer, we use
the current binary (/proc/self/exe) to be executed as the init process, and use
arg "init", we call the first step process "bootstrap", so you always need a "init"
function as the entry of "bootstrap".

In addition to the go init function the early stage bootstrap is handled by importing
[nsenter](https://github.com/opencontainers/runc/blob/master/libcontainer/nsenter/README.md).

```go
import (
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
)

func init() ***REMOVED***
	if len(os.Args) > 1 && os.Args[1] == "init" ***REMOVED***
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		factory, _ := libcontainer.New("")
		if err := factory.StartInitialization(); err != nil ***REMOVED***
			logrus.Fatal(err)
		***REMOVED***
		panic("--this line should have never been executed, congratulations--")
	***REMOVED***
***REMOVED***
```

Then to create a container you first have to initialize an instance of a factory
that will handle the creation and initialization for a container.

```go
factory, err := libcontainer.New("/var/lib/container", libcontainer.Cgroupfs, libcontainer.InitArgs(os.Args[0], "init"))
if err != nil ***REMOVED***
	logrus.Fatal(err)
	return
***REMOVED***
```

Once you have an instance of the factory created we can create a configuration
struct describing how the container is to be created. A sample would look similar to this:

```go
defaultMountFlags := unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV
config := &configs.Config***REMOVED***
	Rootfs: "/your/path/to/rootfs",
	Capabilities: &configs.Capabilities***REMOVED***
                Bounding: []string***REMOVED***
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
            ***REMOVED***,
                Effective: []string***REMOVED***
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
            ***REMOVED***,
                Inheritable: []string***REMOVED***
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
            ***REMOVED***,
                Permitted: []string***REMOVED***
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
            ***REMOVED***,
                Ambient: []string***REMOVED***
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
            ***REMOVED***,
    ***REMOVED***,
	Namespaces: configs.Namespaces([]configs.Namespace***REMOVED***
		***REMOVED***Type: configs.NEWNS***REMOVED***,
		***REMOVED***Type: configs.NEWUTS***REMOVED***,
		***REMOVED***Type: configs.NEWIPC***REMOVED***,
		***REMOVED***Type: configs.NEWPID***REMOVED***,
		***REMOVED***Type: configs.NEWUSER***REMOVED***,
		***REMOVED***Type: configs.NEWNET***REMOVED***,
	***REMOVED***),
	Cgroups: &configs.Cgroup***REMOVED***
		Name:   "test-container",
		Parent: "system",
		Resources: &configs.Resources***REMOVED***
			MemorySwappiness: nil,
			AllowAllDevices:  nil,
			AllowedDevices:   configs.DefaultAllowedDevices,
		***REMOVED***,
	***REMOVED***,
	MaskPaths: []string***REMOVED***
		"/proc/kcore",
		"/sys/firmware",
	***REMOVED***,
	ReadonlyPaths: []string***REMOVED***
		"/proc/sys", "/proc/sysrq-trigger", "/proc/irq", "/proc/bus",
	***REMOVED***,
	Devices:  configs.DefaultAutoCreatedDevices,
	Hostname: "testing",
	Mounts: []*configs.Mount***REMOVED***
		***REMOVED***
			Source:      "proc",
			Destination: "/proc",
			Device:      "proc",
			Flags:       defaultMountFlags,
		***REMOVED***,
		***REMOVED***
			Source:      "tmpfs",
			Destination: "/dev",
			Device:      "tmpfs",
			Flags:       unix.MS_NOSUID | unix.MS_STRICTATIME,
			Data:        "mode=755",
		***REMOVED***,
		***REMOVED***
			Source:      "devpts",
			Destination: "/dev/pts",
			Device:      "devpts",
			Flags:       unix.MS_NOSUID | unix.MS_NOEXEC,
			Data:        "newinstance,ptmxmode=0666,mode=0620,gid=5",
		***REMOVED***,
		***REMOVED***
			Device:      "tmpfs",
			Source:      "shm",
			Destination: "/dev/shm",
			Data:        "mode=1777,size=65536k",
			Flags:       defaultMountFlags,
		***REMOVED***,
		***REMOVED***
			Source:      "mqueue",
			Destination: "/dev/mqueue",
			Device:      "mqueue",
			Flags:       defaultMountFlags,
		***REMOVED***,
		***REMOVED***
			Source:      "sysfs",
			Destination: "/sys",
			Device:      "sysfs",
			Flags:       defaultMountFlags | unix.MS_RDONLY,
		***REMOVED***,
	***REMOVED***,
	UidMappings: []configs.IDMap***REMOVED***
		***REMOVED***
			ContainerID: 0,
			HostID: 1000,
			Size: 65536,
		***REMOVED***,
	***REMOVED***,
	GidMappings: []configs.IDMap***REMOVED***
		***REMOVED***
			ContainerID: 0,
			HostID: 1000,
			Size: 65536,
		***REMOVED***,
	***REMOVED***,
	Networks: []*configs.Network***REMOVED***
		***REMOVED***
			Type:    "loopback",
			Address: "127.0.0.1/0",
			Gateway: "localhost",
		***REMOVED***,
	***REMOVED***,
	Rlimits: []configs.Rlimit***REMOVED***
		***REMOVED***
			Type: unix.RLIMIT_NOFILE,
			Hard: uint64(1025),
			Soft: uint64(1025),
		***REMOVED***,
	***REMOVED***,
***REMOVED***
```

Once you have the configuration populated you can create a container:

```go
container, err := factory.Create("container-id", config)
if err != nil ***REMOVED***
	logrus.Fatal(err)
	return
***REMOVED***
```

To spawn bash as the initial process inside the container and have the
processes pid returned in order to wait, signal, or kill the process:

```go
process := &libcontainer.Process***REMOVED***
	Args:   []string***REMOVED***"/bin/bash"***REMOVED***,
	Env:    []string***REMOVED***"PATH=/bin"***REMOVED***,
	User:   "daemon",
	Stdin:  os.Stdin,
	Stdout: os.Stdout,
	Stderr: os.Stderr,
***REMOVED***

err := container.Run(process)
if err != nil ***REMOVED***
	container.Destroy()
	logrus.Fatal(err)
	return
***REMOVED***

// wait for the process to finish.
_, err := process.Wait()
if err != nil ***REMOVED***
	logrus.Fatal(err)
***REMOVED***

// destroy the container.
container.Destroy()
```

Additional ways to interact with a running container are:

```go
// return all the pids for all processes running inside the container.
processes, err := container.Processes()

// get detailed cpu, memory, io, and network statistics for the container and
// it's processes.
stats, err := container.Stats()

// pause all processes inside the container.
container.Pause()

// resume all paused processes.
container.Resume()

// send signal to container's init process.
container.Signal(signal)

// update container resource constraints.
container.Set(config)

// get current status of the container.
status, err := container.Status()

// get current container's state information.
state, err := container.State()
```


#### Checkpoint & Restore

libcontainer now integrates [CRIU](http://criu.org/) for checkpointing and restoring containers.
This let's you save the state of a process running inside a container to disk, and then restore
that state into a new process, on the same machine or on another machine.

`criu` version 1.5.2 or higher is required to use checkpoint and restore.
If you don't already  have `criu` installed, you can build it from source, following the
[online instructions](http://criu.org/Installation). `criu` is also installed in the docker image
generated when building libcontainer with docker.


## Copyright and license

Code and documentation copyright 2014 Docker, inc. Code released under the Apache 2.0 license.
Docs released under Creative commons.

