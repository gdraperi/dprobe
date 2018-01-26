package daemon

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/caps"
	daemonconfig "github.com/docker/docker/daemon/config"
	"github.com/docker/docker/oci"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/volume"
	"github.com/opencontainers/runc/libcontainer/apparmor"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/devices"
	"github.com/opencontainers/runc/libcontainer/user"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// nolint: gosimple
var (
	deviceCgroupRuleRegex = regexp.MustCompile("^([acb]) ([0-9]+|\\*):([0-9]+|\\*) ([rwm]***REMOVED***1,3***REMOVED***)$")
)

func setResources(s *specs.Spec, r containertypes.Resources) error ***REMOVED***
	weightDevices, err := getBlkioWeightDevices(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	readBpsDevice, err := getBlkioThrottleDevices(r.BlkioDeviceReadBps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	writeBpsDevice, err := getBlkioThrottleDevices(r.BlkioDeviceWriteBps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	readIOpsDevice, err := getBlkioThrottleDevices(r.BlkioDeviceReadIOps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	writeIOpsDevice, err := getBlkioThrottleDevices(r.BlkioDeviceWriteIOps)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	memoryRes := getMemoryResources(r)
	cpuRes, err := getCPUResources(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	blkioWeight := r.BlkioWeight

	specResources := &specs.LinuxResources***REMOVED***
		Memory: memoryRes,
		CPU:    cpuRes,
		BlockIO: &specs.LinuxBlockIO***REMOVED***
			Weight:                  &blkioWeight,
			WeightDevice:            weightDevices,
			ThrottleReadBpsDevice:   readBpsDevice,
			ThrottleWriteBpsDevice:  writeBpsDevice,
			ThrottleReadIOPSDevice:  readIOpsDevice,
			ThrottleWriteIOPSDevice: writeIOpsDevice,
		***REMOVED***,
		Pids: &specs.LinuxPids***REMOVED***
			Limit: r.PidsLimit,
		***REMOVED***,
	***REMOVED***

	if s.Linux.Resources != nil && len(s.Linux.Resources.Devices) > 0 ***REMOVED***
		specResources.Devices = s.Linux.Resources.Devices
	***REMOVED***

	s.Linux.Resources = specResources
	return nil
***REMOVED***

func setDevices(s *specs.Spec, c *container.Container) error ***REMOVED***
	// Build lists of devices allowed and created within the container.
	var devs []specs.LinuxDevice
	devPermissions := s.Linux.Resources.Devices
	if c.HostConfig.Privileged ***REMOVED***
		hostDevices, err := devices.HostDevices()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, d := range hostDevices ***REMOVED***
			devs = append(devs, oci.Device(d))
		***REMOVED***
		devPermissions = []specs.LinuxDeviceCgroup***REMOVED***
			***REMOVED***
				Allow:  true,
				Access: "rwm",
			***REMOVED***,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, deviceMapping := range c.HostConfig.Devices ***REMOVED***
			d, dPermissions, err := oci.DevicesFromPath(deviceMapping.PathOnHost, deviceMapping.PathInContainer, deviceMapping.CgroupPermissions)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			devs = append(devs, d...)
			devPermissions = append(devPermissions, dPermissions...)
		***REMOVED***

		for _, deviceCgroupRule := range c.HostConfig.DeviceCgroupRules ***REMOVED***
			ss := deviceCgroupRuleRegex.FindAllStringSubmatch(deviceCgroupRule, -1)
			if len(ss[0]) != 5 ***REMOVED***
				return fmt.Errorf("invalid device cgroup rule format: '%s'", deviceCgroupRule)
			***REMOVED***
			matches := ss[0]

			dPermissions := specs.LinuxDeviceCgroup***REMOVED***
				Allow:  true,
				Type:   matches[1],
				Access: matches[4],
			***REMOVED***
			if matches[2] == "*" ***REMOVED***
				major := int64(-1)
				dPermissions.Major = &major
			***REMOVED*** else ***REMOVED***
				major, err := strconv.ParseInt(matches[2], 10, 64)
				if err != nil ***REMOVED***
					return fmt.Errorf("invalid major value in device cgroup rule format: '%s'", deviceCgroupRule)
				***REMOVED***
				dPermissions.Major = &major
			***REMOVED***
			if matches[3] == "*" ***REMOVED***
				minor := int64(-1)
				dPermissions.Minor = &minor
			***REMOVED*** else ***REMOVED***
				minor, err := strconv.ParseInt(matches[3], 10, 64)
				if err != nil ***REMOVED***
					return fmt.Errorf("invalid minor value in device cgroup rule format: '%s'", deviceCgroupRule)
				***REMOVED***
				dPermissions.Minor = &minor
			***REMOVED***
			devPermissions = append(devPermissions, dPermissions)
		***REMOVED***
	***REMOVED***

	s.Linux.Devices = append(s.Linux.Devices, devs...)
	s.Linux.Resources.Devices = devPermissions
	return nil
***REMOVED***

func (daemon *Daemon) setRlimits(s *specs.Spec, c *container.Container) error ***REMOVED***
	var rlimits []specs.POSIXRlimit

	// We want to leave the original HostConfig alone so make a copy here
	hostConfig := *c.HostConfig
	// Merge with the daemon defaults
	daemon.mergeUlimits(&hostConfig)
	for _, ul := range hostConfig.Ulimits ***REMOVED***
		rlimits = append(rlimits, specs.POSIXRlimit***REMOVED***
			Type: "RLIMIT_" + strings.ToUpper(ul.Name),
			Soft: uint64(ul.Soft),
			Hard: uint64(ul.Hard),
		***REMOVED***)
	***REMOVED***

	s.Process.Rlimits = rlimits
	return nil
***REMOVED***

func setUser(s *specs.Spec, c *container.Container) error ***REMOVED***
	uid, gid, additionalGids, err := getUser(c, c.Config.User)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	s.Process.User.UID = uid
	s.Process.User.GID = gid
	s.Process.User.AdditionalGids = additionalGids
	return nil
***REMOVED***

func readUserFile(c *container.Container, p string) (io.ReadCloser, error) ***REMOVED***
	fp, err := c.GetResourcePath(p)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return os.Open(fp)
***REMOVED***

func getUser(c *container.Container, username string) (uint32, uint32, []uint32, error) ***REMOVED***
	passwdPath, err := user.GetPasswdPath()
	if err != nil ***REMOVED***
		return 0, 0, nil, err
	***REMOVED***
	groupPath, err := user.GetGroupPath()
	if err != nil ***REMOVED***
		return 0, 0, nil, err
	***REMOVED***
	passwdFile, err := readUserFile(c, passwdPath)
	if err == nil ***REMOVED***
		defer passwdFile.Close()
	***REMOVED***
	groupFile, err := readUserFile(c, groupPath)
	if err == nil ***REMOVED***
		defer groupFile.Close()
	***REMOVED***

	execUser, err := user.GetExecUser(username, nil, passwdFile, groupFile)
	if err != nil ***REMOVED***
		return 0, 0, nil, err
	***REMOVED***

	// todo: fix this double read by a change to libcontainer/user pkg
	groupFile, err = readUserFile(c, groupPath)
	if err == nil ***REMOVED***
		defer groupFile.Close()
	***REMOVED***
	var addGroups []int
	if len(c.HostConfig.GroupAdd) > 0 ***REMOVED***
		addGroups, err = user.GetAdditionalGroups(c.HostConfig.GroupAdd, groupFile)
		if err != nil ***REMOVED***
			return 0, 0, nil, err
		***REMOVED***
	***REMOVED***
	uid := uint32(execUser.Uid)
	gid := uint32(execUser.Gid)
	sgids := append(execUser.Sgids, addGroups...)
	var additionalGids []uint32
	for _, g := range sgids ***REMOVED***
		additionalGids = append(additionalGids, uint32(g))
	***REMOVED***
	return uid, gid, additionalGids, nil
***REMOVED***

func setNamespace(s *specs.Spec, ns specs.LinuxNamespace) ***REMOVED***
	for i, n := range s.Linux.Namespaces ***REMOVED***
		if n.Type == ns.Type ***REMOVED***
			s.Linux.Namespaces[i] = ns
			return
		***REMOVED***
	***REMOVED***
	s.Linux.Namespaces = append(s.Linux.Namespaces, ns)
***REMOVED***

func setCapabilities(s *specs.Spec, c *container.Container) error ***REMOVED***
	var caplist []string
	var err error
	if c.HostConfig.Privileged ***REMOVED***
		caplist = caps.GetAllCapabilities()
	***REMOVED*** else ***REMOVED***
		caplist, err = caps.TweakCapabilities(s.Process.Capabilities.Effective, c.HostConfig.CapAdd, c.HostConfig.CapDrop)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	s.Process.Capabilities.Effective = caplist
	s.Process.Capabilities.Bounding = caplist
	s.Process.Capabilities.Permitted = caplist
	s.Process.Capabilities.Inheritable = caplist
	return nil
***REMOVED***

func setNamespaces(daemon *Daemon, s *specs.Spec, c *container.Container) error ***REMOVED***
	userNS := false
	// user
	if c.HostConfig.UsernsMode.IsPrivate() ***REMOVED***
		uidMap := daemon.idMappings.UIDs()
		if uidMap != nil ***REMOVED***
			userNS = true
			ns := specs.LinuxNamespace***REMOVED***Type: "user"***REMOVED***
			setNamespace(s, ns)
			s.Linux.UIDMappings = specMapping(uidMap)
			s.Linux.GIDMappings = specMapping(daemon.idMappings.GIDs())
		***REMOVED***
	***REMOVED***
	// network
	if !c.Config.NetworkDisabled ***REMOVED***
		ns := specs.LinuxNamespace***REMOVED***Type: "network"***REMOVED***
		parts := strings.SplitN(string(c.HostConfig.NetworkMode), ":", 2)
		if parts[0] == "container" ***REMOVED***
			nc, err := daemon.getNetworkedContainer(c.ID, c.HostConfig.NetworkMode.ConnectedContainer())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			ns.Path = fmt.Sprintf("/proc/%d/ns/net", nc.State.GetPID())
			if userNS ***REMOVED***
				// to share a net namespace, they must also share a user namespace
				nsUser := specs.LinuxNamespace***REMOVED***Type: "user"***REMOVED***
				nsUser.Path = fmt.Sprintf("/proc/%d/ns/user", nc.State.GetPID())
				setNamespace(s, nsUser)
			***REMOVED***
		***REMOVED*** else if c.HostConfig.NetworkMode.IsHost() ***REMOVED***
			ns.Path = c.NetworkSettings.SandboxKey
		***REMOVED***
		setNamespace(s, ns)
	***REMOVED***

	// ipc
	ipcMode := c.HostConfig.IpcMode
	switch ***REMOVED***
	case ipcMode.IsContainer():
		ns := specs.LinuxNamespace***REMOVED***Type: "ipc"***REMOVED***
		ic, err := daemon.getIpcContainer(ipcMode.Container())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		ns.Path = fmt.Sprintf("/proc/%d/ns/ipc", ic.State.GetPID())
		setNamespace(s, ns)
		if userNS ***REMOVED***
			// to share an IPC namespace, they must also share a user namespace
			nsUser := specs.LinuxNamespace***REMOVED***Type: "user"***REMOVED***
			nsUser.Path = fmt.Sprintf("/proc/%d/ns/user", ic.State.GetPID())
			setNamespace(s, nsUser)
		***REMOVED***
	case ipcMode.IsHost():
		oci.RemoveNamespace(s, specs.LinuxNamespaceType("ipc"))
	case ipcMode.IsEmpty():
		// A container was created by an older version of the daemon.
		// The default behavior used to be what is now called "shareable".
		fallthrough
	case ipcMode.IsPrivate(), ipcMode.IsShareable(), ipcMode.IsNone():
		ns := specs.LinuxNamespace***REMOVED***Type: "ipc"***REMOVED***
		setNamespace(s, ns)
	default:
		return fmt.Errorf("Invalid IPC mode: %v", ipcMode)
	***REMOVED***

	// pid
	if c.HostConfig.PidMode.IsContainer() ***REMOVED***
		ns := specs.LinuxNamespace***REMOVED***Type: "pid"***REMOVED***
		pc, err := daemon.getPidContainer(c)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		ns.Path = fmt.Sprintf("/proc/%d/ns/pid", pc.State.GetPID())
		setNamespace(s, ns)
		if userNS ***REMOVED***
			// to share a PID namespace, they must also share a user namespace
			nsUser := specs.LinuxNamespace***REMOVED***Type: "user"***REMOVED***
			nsUser.Path = fmt.Sprintf("/proc/%d/ns/user", pc.State.GetPID())
			setNamespace(s, nsUser)
		***REMOVED***
	***REMOVED*** else if c.HostConfig.PidMode.IsHost() ***REMOVED***
		oci.RemoveNamespace(s, specs.LinuxNamespaceType("pid"))
	***REMOVED*** else ***REMOVED***
		ns := specs.LinuxNamespace***REMOVED***Type: "pid"***REMOVED***
		setNamespace(s, ns)
	***REMOVED***
	// uts
	if c.HostConfig.UTSMode.IsHost() ***REMOVED***
		oci.RemoveNamespace(s, specs.LinuxNamespaceType("uts"))
		s.Hostname = ""
	***REMOVED***

	return nil
***REMOVED***

func specMapping(s []idtools.IDMap) []specs.LinuxIDMapping ***REMOVED***
	var ids []specs.LinuxIDMapping
	for _, item := range s ***REMOVED***
		ids = append(ids, specs.LinuxIDMapping***REMOVED***
			HostID:      uint32(item.HostID),
			ContainerID: uint32(item.ContainerID),
			Size:        uint32(item.Size),
		***REMOVED***)
	***REMOVED***
	return ids
***REMOVED***

func getMountInfo(mountinfo []*mount.Info, dir string) *mount.Info ***REMOVED***
	for _, m := range mountinfo ***REMOVED***
		if m.Mountpoint == dir ***REMOVED***
			return m
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Get the source mount point of directory passed in as argument. Also return
// optional fields.
func getSourceMount(source string) (string, string, error) ***REMOVED***
	// Ensure any symlinks are resolved.
	sourcePath, err := filepath.EvalSymlinks(source)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	mountinfos, err := mount.GetMounts()
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	mountinfo := getMountInfo(mountinfos, sourcePath)
	if mountinfo != nil ***REMOVED***
		return sourcePath, mountinfo.Optional, nil
	***REMOVED***

	path := sourcePath
	for ***REMOVED***
		path = filepath.Dir(path)

		mountinfo = getMountInfo(mountinfos, path)
		if mountinfo != nil ***REMOVED***
			return path, mountinfo.Optional, nil
		***REMOVED***

		if path == "/" ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	// If we are here, we did not find parent mount. Something is wrong.
	return "", "", fmt.Errorf("Could not find source mount of %s", source)
***REMOVED***

// Ensure mount point on which path is mounted, is shared.
func ensureShared(path string) error ***REMOVED***
	sharedMount := false

	sourceMount, optionalOpts, err := getSourceMount(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Make sure source mount point is shared.
	optsSplit := strings.Split(optionalOpts, " ")
	for _, opt := range optsSplit ***REMOVED***
		if strings.HasPrefix(opt, "shared:") ***REMOVED***
			sharedMount = true
			break
		***REMOVED***
	***REMOVED***

	if !sharedMount ***REMOVED***
		return fmt.Errorf("path %s is mounted on %s but it is not a shared mount", path, sourceMount)
	***REMOVED***
	return nil
***REMOVED***

// Ensure mount point on which path is mounted, is either shared or slave.
func ensureSharedOrSlave(path string) error ***REMOVED***
	sharedMount := false
	slaveMount := false

	sourceMount, optionalOpts, err := getSourceMount(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Make sure source mount point is shared.
	optsSplit := strings.Split(optionalOpts, " ")
	for _, opt := range optsSplit ***REMOVED***
		if strings.HasPrefix(opt, "shared:") ***REMOVED***
			sharedMount = true
			break
		***REMOVED*** else if strings.HasPrefix(opt, "master:") ***REMOVED***
			slaveMount = true
			break
		***REMOVED***
	***REMOVED***

	if !sharedMount && !slaveMount ***REMOVED***
		return fmt.Errorf("path %s is mounted on %s but it is not a shared or slave mount", path, sourceMount)
	***REMOVED***
	return nil
***REMOVED***

// Get the set of mount flags that are set on the mount that contains the given
// path and are locked by CL_UNPRIVILEGED. This is necessary to ensure that
// bind-mounting "with options" will not fail with user namespaces, due to
// kernel restrictions that require user namespace mounts to preserve
// CL_UNPRIVILEGED locked flags.
func getUnprivilegedMountFlags(path string) ([]string, error) ***REMOVED***
	var statfs unix.Statfs_t
	if err := unix.Statfs(path, &statfs); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// The set of keys come from https://github.com/torvalds/linux/blob/v4.13/fs/namespace.c#L1034-L1048.
	unprivilegedFlags := map[uint64]string***REMOVED***
		unix.MS_RDONLY:     "ro",
		unix.MS_NODEV:      "nodev",
		unix.MS_NOEXEC:     "noexec",
		unix.MS_NOSUID:     "nosuid",
		unix.MS_NOATIME:    "noatime",
		unix.MS_RELATIME:   "relatime",
		unix.MS_NODIRATIME: "nodiratime",
	***REMOVED***

	var flags []string
	for mask, flag := range unprivilegedFlags ***REMOVED***
		if uint64(statfs.Flags)&mask == mask ***REMOVED***
			flags = append(flags, flag)
		***REMOVED***
	***REMOVED***

	return flags, nil
***REMOVED***

var (
	mountPropagationMap = map[string]int***REMOVED***
		"private":  mount.PRIVATE,
		"rprivate": mount.RPRIVATE,
		"shared":   mount.SHARED,
		"rshared":  mount.RSHARED,
		"slave":    mount.SLAVE,
		"rslave":   mount.RSLAVE,
	***REMOVED***

	mountPropagationReverseMap = map[int]string***REMOVED***
		mount.PRIVATE:  "private",
		mount.RPRIVATE: "rprivate",
		mount.SHARED:   "shared",
		mount.RSHARED:  "rshared",
		mount.SLAVE:    "slave",
		mount.RSLAVE:   "rslave",
	***REMOVED***
)

// inSlice tests whether a string is contained in a slice of strings or not.
// Comparison is case sensitive
func inSlice(slice []string, s string) bool ***REMOVED***
	for _, ss := range slice ***REMOVED***
		if s == ss ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func setMounts(daemon *Daemon, s *specs.Spec, c *container.Container, mounts []container.Mount) error ***REMOVED***
	userMounts := make(map[string]struct***REMOVED******REMOVED***)
	for _, m := range mounts ***REMOVED***
		userMounts[m.Destination] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	// Copy all mounts from spec to defaultMounts, except for
	//  - mounts overriden by a user supplied mount;
	//  - all mounts under /dev if a user supplied /dev is present;
	//  - /dev/shm, in case IpcMode is none.
	// While at it, also
	//  - set size for /dev/shm from shmsize.
	var defaultMounts []specs.Mount
	_, mountDev := userMounts["/dev"]
	for _, m := range s.Mounts ***REMOVED***
		if _, ok := userMounts[m.Destination]; ok ***REMOVED***
			// filter out mount overridden by a user supplied mount
			continue
		***REMOVED***
		if mountDev && strings.HasPrefix(m.Destination, "/dev/") ***REMOVED***
			// filter out everything under /dev if /dev is user-mounted
			continue
		***REMOVED***

		if m.Destination == "/dev/shm" ***REMOVED***
			if c.HostConfig.IpcMode.IsNone() ***REMOVED***
				// filter out /dev/shm for "none" IpcMode
				continue
			***REMOVED***
			// set size for /dev/shm mount from spec
			sizeOpt := "size=" + strconv.FormatInt(c.HostConfig.ShmSize, 10)
			m.Options = append(m.Options, sizeOpt)
		***REMOVED***

		defaultMounts = append(defaultMounts, m)
	***REMOVED***

	s.Mounts = defaultMounts
	for _, m := range mounts ***REMOVED***
		for _, cm := range s.Mounts ***REMOVED***
			if cm.Destination == m.Destination ***REMOVED***
				return duplicateMountPointError(m.Destination)
			***REMOVED***
		***REMOVED***

		if m.Source == "tmpfs" ***REMOVED***
			data := m.Data
			parser := volume.NewParser("linux")
			options := []string***REMOVED***"noexec", "nosuid", "nodev", string(parser.DefaultPropagationMode())***REMOVED***
			if data != "" ***REMOVED***
				options = append(options, strings.Split(data, ",")...)
			***REMOVED***

			merged, err := mount.MergeTmpfsOptions(options)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			s.Mounts = append(s.Mounts, specs.Mount***REMOVED***Destination: m.Destination, Source: m.Source, Type: "tmpfs", Options: merged***REMOVED***)
			continue
		***REMOVED***

		mt := specs.Mount***REMOVED***Destination: m.Destination, Source: m.Source, Type: "bind"***REMOVED***

		// Determine property of RootPropagation based on volume
		// properties. If a volume is shared, then keep root propagation
		// shared. This should work for slave and private volumes too.
		//
		// For slave volumes, it can be either [r]shared/[r]slave.
		//
		// For private volumes any root propagation value should work.
		pFlag := mountPropagationMap[m.Propagation]
		if pFlag == mount.SHARED || pFlag == mount.RSHARED ***REMOVED***
			if err := ensureShared(m.Source); err != nil ***REMOVED***
				return err
			***REMOVED***
			rootpg := mountPropagationMap[s.Linux.RootfsPropagation]
			if rootpg != mount.SHARED && rootpg != mount.RSHARED ***REMOVED***
				s.Linux.RootfsPropagation = mountPropagationReverseMap[mount.SHARED]
			***REMOVED***
		***REMOVED*** else if pFlag == mount.SLAVE || pFlag == mount.RSLAVE ***REMOVED***
			if err := ensureSharedOrSlave(m.Source); err != nil ***REMOVED***
				return err
			***REMOVED***
			rootpg := mountPropagationMap[s.Linux.RootfsPropagation]
			if rootpg != mount.SHARED && rootpg != mount.RSHARED && rootpg != mount.SLAVE && rootpg != mount.RSLAVE ***REMOVED***
				s.Linux.RootfsPropagation = mountPropagationReverseMap[mount.RSLAVE]
			***REMOVED***
		***REMOVED***

		opts := []string***REMOVED***"rbind"***REMOVED***
		if !m.Writable ***REMOVED***
			opts = append(opts, "ro")
		***REMOVED***
		if pFlag != 0 ***REMOVED***
			opts = append(opts, mountPropagationReverseMap[pFlag])
		***REMOVED***

		// If we are using user namespaces, then we must make sure that we
		// don't drop any of the CL_UNPRIVILEGED "locked" flags of the source
		// "mount" when we bind-mount. The reason for this is that at the point
		// when runc sets up the root filesystem, it is already inside a user
		// namespace, and thus cannot change any flags that are locked.
		if daemon.configStore.RemappedRoot != "" ***REMOVED***
			unprivOpts, err := getUnprivilegedMountFlags(m.Source)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			opts = append(opts, unprivOpts...)
		***REMOVED***

		mt.Options = opts
		s.Mounts = append(s.Mounts, mt)
	***REMOVED***

	if s.Root.Readonly ***REMOVED***
		for i, m := range s.Mounts ***REMOVED***
			switch m.Destination ***REMOVED***
			case "/proc", "/dev/pts", "/dev/mqueue", "/dev":
				continue
			***REMOVED***
			if _, ok := userMounts[m.Destination]; !ok ***REMOVED***
				if !inSlice(m.Options, "ro") ***REMOVED***
					s.Mounts[i].Options = append(s.Mounts[i].Options, "ro")
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if c.HostConfig.Privileged ***REMOVED***
		if !s.Root.Readonly ***REMOVED***
			// clear readonly for /sys
			for i := range s.Mounts ***REMOVED***
				if s.Mounts[i].Destination == "/sys" ***REMOVED***
					clearReadOnly(&s.Mounts[i])
				***REMOVED***
			***REMOVED***
		***REMOVED***
		s.Linux.ReadonlyPaths = nil
		s.Linux.MaskedPaths = nil
	***REMOVED***

	// TODO: until a kernel/mount solution exists for handling remount in a user namespace,
	// we must clear the readonly flag for the cgroups mount (@mrunalp concurs)
	if uidMap := daemon.idMappings.UIDs(); uidMap != nil || c.HostConfig.Privileged ***REMOVED***
		for i, m := range s.Mounts ***REMOVED***
			if m.Type == "cgroup" ***REMOVED***
				clearReadOnly(&s.Mounts[i])
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (daemon *Daemon) populateCommonSpec(s *specs.Spec, c *container.Container) error ***REMOVED***
	linkedEnv, err := daemon.setupLinkedContainers(c)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	s.Root = &specs.Root***REMOVED***
		Path:     c.BaseFS.Path(),
		Readonly: c.HostConfig.ReadonlyRootfs,
	***REMOVED***
	if err := c.SetupWorkingDirectory(daemon.idMappings.RootPair()); err != nil ***REMOVED***
		return err
	***REMOVED***
	cwd := c.Config.WorkingDir
	if len(cwd) == 0 ***REMOVED***
		cwd = "/"
	***REMOVED***
	s.Process.Args = append([]string***REMOVED***c.Path***REMOVED***, c.Args...)

	// only add the custom init if it is specified and the container is running in its
	// own private pid namespace.  It does not make sense to add if it is running in the
	// host namespace or another container's pid namespace where we already have an init
	if c.HostConfig.PidMode.IsPrivate() ***REMOVED***
		if (c.HostConfig.Init != nil && *c.HostConfig.Init) ||
			(c.HostConfig.Init == nil && daemon.configStore.Init) ***REMOVED***
			s.Process.Args = append([]string***REMOVED***"/dev/init", "--", c.Path***REMOVED***, c.Args...)
			var path string
			if daemon.configStore.InitPath == "" ***REMOVED***
				path, err = exec.LookPath(daemonconfig.DefaultInitBinary)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			if daemon.configStore.InitPath != "" ***REMOVED***
				path = daemon.configStore.InitPath
			***REMOVED***
			s.Mounts = append(s.Mounts, specs.Mount***REMOVED***
				Destination: "/dev/init",
				Type:        "bind",
				Source:      path,
				Options:     []string***REMOVED***"bind", "ro"***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	s.Process.Cwd = cwd
	s.Process.Env = c.CreateDaemonEnvironment(c.Config.Tty, linkedEnv)
	s.Process.Terminal = c.Config.Tty
	s.Hostname = c.FullHostname()

	return nil
***REMOVED***

func (daemon *Daemon) createSpec(c *container.Container) (*specs.Spec, error) ***REMOVED***
	s := oci.DefaultSpec()
	if err := daemon.populateCommonSpec(&s, c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var cgroupsPath string
	scopePrefix := "docker"
	parent := "/docker"
	useSystemd := UsingSystemd(daemon.configStore)
	if useSystemd ***REMOVED***
		parent = "system.slice"
	***REMOVED***

	if c.HostConfig.CgroupParent != "" ***REMOVED***
		parent = c.HostConfig.CgroupParent
	***REMOVED*** else if daemon.configStore.CgroupParent != "" ***REMOVED***
		parent = daemon.configStore.CgroupParent
	***REMOVED***

	if useSystemd ***REMOVED***
		cgroupsPath = parent + ":" + scopePrefix + ":" + c.ID
		logrus.Debugf("createSpec: cgroupsPath: %s", cgroupsPath)
	***REMOVED*** else ***REMOVED***
		cgroupsPath = filepath.Join(parent, c.ID)
	***REMOVED***
	s.Linux.CgroupsPath = cgroupsPath

	if err := setResources(&s, c.HostConfig.Resources); err != nil ***REMOVED***
		return nil, fmt.Errorf("linux runtime spec resources: %v", err)
	***REMOVED***
	s.Linux.Sysctl = c.HostConfig.Sysctls

	p := s.Linux.CgroupsPath
	if useSystemd ***REMOVED***
		initPath, err := cgroups.GetInitCgroup("cpu")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		_, err = cgroups.GetOwnCgroup("cpu")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		p = filepath.Join(initPath, s.Linux.CgroupsPath)
	***REMOVED***

	// Clean path to guard against things like ../../../BAD
	parentPath := filepath.Dir(p)
	if !filepath.IsAbs(parentPath) ***REMOVED***
		parentPath = filepath.Clean("/" + parentPath)
	***REMOVED***

	if err := daemon.initCgroupsPath(parentPath); err != nil ***REMOVED***
		return nil, fmt.Errorf("linux init cgroups path: %v", err)
	***REMOVED***
	if err := setDevices(&s, c); err != nil ***REMOVED***
		return nil, fmt.Errorf("linux runtime spec devices: %v", err)
	***REMOVED***
	if err := daemon.setRlimits(&s, c); err != nil ***REMOVED***
		return nil, fmt.Errorf("linux runtime spec rlimits: %v", err)
	***REMOVED***
	if err := setUser(&s, c); err != nil ***REMOVED***
		return nil, fmt.Errorf("linux spec user: %v", err)
	***REMOVED***
	if err := setNamespaces(daemon, &s, c); err != nil ***REMOVED***
		return nil, fmt.Errorf("linux spec namespaces: %v", err)
	***REMOVED***
	if err := setCapabilities(&s, c); err != nil ***REMOVED***
		return nil, fmt.Errorf("linux spec capabilities: %v", err)
	***REMOVED***
	if err := setSeccomp(daemon, &s, c); err != nil ***REMOVED***
		return nil, fmt.Errorf("linux seccomp: %v", err)
	***REMOVED***

	if err := daemon.setupContainerMountsRoot(c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := daemon.setupIpcDirs(c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := daemon.setupSecretDir(c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := daemon.setupConfigDir(c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ms, err := daemon.setupMounts(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !c.HostConfig.IpcMode.IsPrivate() && !c.HostConfig.IpcMode.IsEmpty() ***REMOVED***
		ms = append(ms, c.IpcMounts()...)
	***REMOVED***

	tmpfsMounts, err := c.TmpfsMounts()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ms = append(ms, tmpfsMounts...)

	secretMounts, err := c.SecretMounts()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ms = append(ms, secretMounts...)

	configMounts, err := c.ConfigMounts()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ms = append(ms, configMounts...)

	sort.Sort(mounts(ms))
	if err := setMounts(daemon, &s, c, ms); err != nil ***REMOVED***
		return nil, fmt.Errorf("linux mounts: %v", err)
	***REMOVED***

	for _, ns := range s.Linux.Namespaces ***REMOVED***
		if ns.Type == "network" && ns.Path == "" && !c.Config.NetworkDisabled ***REMOVED***
			target, err := os.Readlink(filepath.Join("/proc", strconv.Itoa(os.Getpid()), "exe"))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			s.Hooks = &specs.Hooks***REMOVED***
				Prestart: []specs.Hook***REMOVED******REMOVED***
					Path: target, // FIXME: cross-platform
					Args: []string***REMOVED***"libnetwork-setkey", c.ID, daemon.netController.ID()***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if apparmor.IsEnabled() ***REMOVED***
		var appArmorProfile string
		if c.AppArmorProfile != "" ***REMOVED***
			appArmorProfile = c.AppArmorProfile
		***REMOVED*** else if c.HostConfig.Privileged ***REMOVED***
			appArmorProfile = "unconfined"
		***REMOVED*** else ***REMOVED***
			appArmorProfile = "docker-default"
		***REMOVED***

		if appArmorProfile == "docker-default" ***REMOVED***
			// Unattended upgrades and other fun services can unload AppArmor
			// profiles inadvertently. Since we cannot store our profile in
			// /etc/apparmor.d, nor can we practically add other ways of
			// telling the system to keep our profile loaded, in order to make
			// sure that we keep the default profile enabled we dynamically
			// reload it if necessary.
			if err := ensureDefaultAppArmorProfile(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***

		s.Process.ApparmorProfile = appArmorProfile
	***REMOVED***
	s.Process.SelinuxLabel = c.GetProcessLabel()
	s.Process.NoNewPrivileges = c.NoNewPrivileges
	s.Process.OOMScoreAdj = &c.HostConfig.OomScoreAdj
	s.Linux.MountLabel = c.MountLabel

	return &s, nil
***REMOVED***

func clearReadOnly(m *specs.Mount) ***REMOVED***
	var opt []string
	for _, o := range m.Options ***REMOVED***
		if o != "ro" ***REMOVED***
			opt = append(opt, o)
		***REMOVED***
	***REMOVED***
	m.Options = opt
***REMOVED***

// mergeUlimits merge the Ulimits from HostConfig with daemon defaults, and update HostConfig
func (daemon *Daemon) mergeUlimits(c *containertypes.HostConfig) ***REMOVED***
	ulimits := c.Ulimits
	// Merge ulimits with daemon defaults
	ulIdx := make(map[string]struct***REMOVED******REMOVED***)
	for _, ul := range ulimits ***REMOVED***
		ulIdx[ul.Name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	for name, ul := range daemon.configStore.Ulimits ***REMOVED***
		if _, exists := ulIdx[name]; !exists ***REMOVED***
			ulimits = append(ulimits, ul)
		***REMOVED***
	***REMOVED***
	c.Ulimits = ulimits
***REMOVED***
