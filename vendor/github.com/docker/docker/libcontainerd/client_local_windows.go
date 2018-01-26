package libcontainerd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Microsoft/hcsshim"
	opengcs "github.com/Microsoft/opengcs/client"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/docker/docker/pkg/sysinfo"
	"github.com/docker/docker/pkg/system"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

const InitProcessName = "init"

type process struct ***REMOVED***
	id         string
	pid        int
	hcsProcess hcsshim.Process
***REMOVED***

type container struct ***REMOVED***
	sync.Mutex

	// The ociSpec is required, as client.Create() needs a spec, but can
	// be called from the RestartManager context which does not otherwise
	// have access to the Spec
	ociSpec *specs.Spec

	isWindows           bool
	manualStopRequested bool
	hcsContainer        hcsshim.Container

	id            string
	status        Status
	exitedAt      time.Time
	exitCode      uint32
	waitCh        chan struct***REMOVED******REMOVED***
	init          *process
	execs         map[string]*process
	updatePending bool
***REMOVED***

// Win32 error codes that are used for various workarounds
// These really should be ALL_CAPS to match golangs syscall library and standard
// Win32 error conventions, but golint insists on CamelCase.
const (
	CoEClassstring     = syscall.Errno(0x800401F3) // Invalid class string
	ErrorNoNetwork     = syscall.Errno(1222)       // The network is not present or not started
	ErrorBadPathname   = syscall.Errno(161)        // The specified path is invalid
	ErrorInvalidObject = syscall.Errno(0x800710D8) // The object identifier does not represent a valid object
)

// defaultOwner is a tag passed to HCS to allow it to differentiate between
// container creator management stacks. We hard code "docker" in the case
// of docker.
const defaultOwner = "docker"

func (c *client) Version(ctx context.Context) (containerd.Version, error) ***REMOVED***
	return containerd.Version***REMOVED******REMOVED***, errors.New("not implemented on Windows")
***REMOVED***

// Create is the entrypoint to create a container from a spec.
// Table below shows the fields required for HCS JSON calling parameters,
// where if not populated, is omitted.
// +-----------------+--------------------------------------------+---------------------------------------------------+
// |                 | Isolation=Process                          | Isolation=Hyper-V                                 |
// +-----------------+--------------------------------------------+---------------------------------------------------+
// | VolumePath      | \\?\\Volume***REMOVED***GUIDa***REMOVED***                         |                                                   |
// | LayerFolderPath | %root%\windowsfilter\containerID           | %root%\windowsfilter\containerID (servicing only) |
// | Layers[]        | ID=GUIDb;Path=%root%\windowsfilter\layerID | ID=GUIDb;Path=%root%\windowsfilter\layerID        |
// | HvRuntime       |                                            | ImagePath=%root%\BaseLayerID\UtilityVM            |
// +-----------------+--------------------------------------------+---------------------------------------------------+
//
// Isolation=Process example:
//
// ***REMOVED***
//	"SystemType": "Container",
//	"Name": "5e0055c814a6005b8e57ac59f9a522066e0af12b48b3c26a9416e23907698776",
//	"Owner": "docker",
//	"VolumePath": "\\\\\\\\?\\\\Volume***REMOVED***66d1ef4c-7a00-11e6-8948-00155ddbef9d***REMOVED***",
//	"IgnoreFlushesDuringBoot": true,
//	"LayerFolderPath": "C:\\\\control\\\\windowsfilter\\\\5e0055c814a6005b8e57ac59f9a522066e0af12b48b3c26a9416e23907698776",
//	"Layers": [***REMOVED***
//		"ID": "18955d65-d45a-557b-bf1c-49d6dfefc526",
//		"Path": "C:\\\\control\\\\windowsfilter\\\\65bf96e5760a09edf1790cb229e2dfb2dbd0fcdc0bf7451bae099106bfbfea0c"
//	***REMOVED***],
//	"HostName": "5e0055c814a6",
//	"MappedDirectories": [],
//	"HvPartition": false,
//	"EndpointList": ["eef2649d-bb17-4d53-9937-295a8efe6f2c"],
//	"Servicing": false
//***REMOVED***
//
// Isolation=Hyper-V example:
//
//***REMOVED***
//	"SystemType": "Container",
//	"Name": "475c2c58933b72687a88a441e7e0ca4bd72d76413c5f9d5031fee83b98f6045d",
//	"Owner": "docker",
//	"IgnoreFlushesDuringBoot": true,
//	"Layers": [***REMOVED***
//		"ID": "18955d65-d45a-557b-bf1c-49d6dfefc526",
//		"Path": "C:\\\\control\\\\windowsfilter\\\\65bf96e5760a09edf1790cb229e2dfb2dbd0fcdc0bf7451bae099106bfbfea0c"
//	***REMOVED***],
//	"HostName": "475c2c58933b",
//	"MappedDirectories": [],
//	"HvPartition": true,
//	"EndpointList": ["e1bb1e61-d56f-405e-b75d-fd520cefa0cb"],
//	"DNSSearchList": "a.com,b.com,c.com",
//	"HvRuntime": ***REMOVED***
//		"ImagePath": "C:\\\\control\\\\windowsfilter\\\\65bf96e5760a09edf1790cb229e2dfb2dbd0fcdc0bf7451bae099106bfbfea0c\\\\UtilityVM"
//	***REMOVED***,
//	"Servicing": false
//***REMOVED***
func (c *client) Create(_ context.Context, id string, spec *specs.Spec, runtimeOptions interface***REMOVED******REMOVED***) error ***REMOVED***
	if ctr := c.getContainer(id); ctr != nil ***REMOVED***
		return errors.WithStack(newConflictError("id already in use"))
	***REMOVED***

	// spec.Linux must be nil for Windows containers, but spec.Windows
	// will be filled in regardless of container platform.  This is a
	// temporary workaround due to LCOW requiring layer folder paths,
	// which are stored under spec.Windows.
	//
	// TODO: @darrenstahlmsft fix this once the OCI spec is updated to
	// support layer folder paths for LCOW
	if spec.Linux == nil ***REMOVED***
		return c.createWindows(id, spec, runtimeOptions)
	***REMOVED***
	return c.createLinux(id, spec, runtimeOptions)
***REMOVED***

func (c *client) createWindows(id string, spec *specs.Spec, runtimeOptions interface***REMOVED******REMOVED***) error ***REMOVED***
	logger := c.logger.WithField("container", id)
	configuration := &hcsshim.ContainerConfig***REMOVED***
		SystemType: "Container",
		Name:       id,
		Owner:      defaultOwner,
		IgnoreFlushesDuringBoot: spec.Windows.IgnoreFlushesDuringBoot,
		HostName:                spec.Hostname,
		HvPartition:             false,
		Servicing:               spec.Windows.Servicing,
	***REMOVED***

	if spec.Windows.Resources != nil ***REMOVED***
		if spec.Windows.Resources.CPU != nil ***REMOVED***
			if spec.Windows.Resources.CPU.Count != nil ***REMOVED***
				// This check is being done here rather than in adaptContainerSettings
				// because we don't want to update the HostConfig in case this container
				// is moved to a host with more CPUs than this one.
				cpuCount := *spec.Windows.Resources.CPU.Count
				hostCPUCount := uint64(sysinfo.NumCPU())
				if cpuCount > hostCPUCount ***REMOVED***
					c.logger.Warnf("Changing requested CPUCount of %d to current number of processors, %d", cpuCount, hostCPUCount)
					cpuCount = hostCPUCount
				***REMOVED***
				configuration.ProcessorCount = uint32(cpuCount)
			***REMOVED***
			if spec.Windows.Resources.CPU.Shares != nil ***REMOVED***
				configuration.ProcessorWeight = uint64(*spec.Windows.Resources.CPU.Shares)
			***REMOVED***
			if spec.Windows.Resources.CPU.Maximum != nil ***REMOVED***
				configuration.ProcessorMaximum = int64(*spec.Windows.Resources.CPU.Maximum)
			***REMOVED***
		***REMOVED***
		if spec.Windows.Resources.Memory != nil ***REMOVED***
			if spec.Windows.Resources.Memory.Limit != nil ***REMOVED***
				configuration.MemoryMaximumInMB = int64(*spec.Windows.Resources.Memory.Limit) / 1024 / 1024
			***REMOVED***
		***REMOVED***
		if spec.Windows.Resources.Storage != nil ***REMOVED***
			if spec.Windows.Resources.Storage.Bps != nil ***REMOVED***
				configuration.StorageBandwidthMaximum = *spec.Windows.Resources.Storage.Bps
			***REMOVED***
			if spec.Windows.Resources.Storage.Iops != nil ***REMOVED***
				configuration.StorageIOPSMaximum = *spec.Windows.Resources.Storage.Iops
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if spec.Windows.HyperV != nil ***REMOVED***
		configuration.HvPartition = true
	***REMOVED***

	if spec.Windows.Network != nil ***REMOVED***
		configuration.EndpointList = spec.Windows.Network.EndpointList
		configuration.AllowUnqualifiedDNSQuery = spec.Windows.Network.AllowUnqualifiedDNSQuery
		if spec.Windows.Network.DNSSearchList != nil ***REMOVED***
			configuration.DNSSearchList = strings.Join(spec.Windows.Network.DNSSearchList, ",")
		***REMOVED***
		configuration.NetworkSharedContainerName = spec.Windows.Network.NetworkSharedContainerName
	***REMOVED***

	if cs, ok := spec.Windows.CredentialSpec.(string); ok ***REMOVED***
		configuration.Credentials = cs
	***REMOVED***

	// We must have least two layers in the spec, the bottom one being a
	// base image, the top one being the RW layer.
	if spec.Windows.LayerFolders == nil || len(spec.Windows.LayerFolders) < 2 ***REMOVED***
		return fmt.Errorf("OCI spec is invalid - at least two LayerFolders must be supplied to the runtime")
	***REMOVED***

	// Strip off the top-most layer as that's passed in separately to HCS
	configuration.LayerFolderPath = spec.Windows.LayerFolders[len(spec.Windows.LayerFolders)-1]
	layerFolders := spec.Windows.LayerFolders[:len(spec.Windows.LayerFolders)-1]

	if configuration.HvPartition ***REMOVED***
		// We don't currently support setting the utility VM image explicitly.
		// TODO @swernli/jhowardmsft circa RS3/4, this may be re-locatable.
		if spec.Windows.HyperV.UtilityVMPath != "" ***REMOVED***
			return errors.New("runtime does not support an explicit utility VM path for Hyper-V containers")
		***REMOVED***

		// Find the upper-most utility VM image.
		var uvmImagePath string
		for _, path := range layerFolders ***REMOVED***
			fullPath := filepath.Join(path, "UtilityVM")
			_, err := os.Stat(fullPath)
			if err == nil ***REMOVED***
				uvmImagePath = fullPath
				break
			***REMOVED***
			if !os.IsNotExist(err) ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if uvmImagePath == "" ***REMOVED***
			return errors.New("utility VM image could not be found")
		***REMOVED***
		configuration.HvRuntime = &hcsshim.HvRuntime***REMOVED***ImagePath: uvmImagePath***REMOVED***

		if spec.Root.Path != "" ***REMOVED***
			return errors.New("OCI spec is invalid - Root.Path must be omitted for a Hyper-V container")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		const volumeGUIDRegex = `^\\\\\?\\(Volume)\***REMOVED******REMOVED***0,1***REMOVED***[0-9a-fA-F]***REMOVED***8***REMOVED***\-[0-9a-fA-F]***REMOVED***4***REMOVED***\-[0-9a-fA-F]***REMOVED***4***REMOVED***\-[0-9a-fA-F]***REMOVED***4***REMOVED***\-[0-9a-fA-F]***REMOVED***12***REMOVED***(\***REMOVED***)***REMOVED***0,1***REMOVED***\***REMOVED***\\$`
		if _, err := regexp.MatchString(volumeGUIDRegex, spec.Root.Path); err != nil ***REMOVED***
			return fmt.Errorf(`OCI spec is invalid - Root.Path '%s' must be a volume GUID path in the format '\\?\Volume***REMOVED***GUID***REMOVED***\'`, spec.Root.Path)
		***REMOVED***
		// HCS API requires the trailing backslash to be removed
		configuration.VolumePath = spec.Root.Path[:len(spec.Root.Path)-1]
	***REMOVED***

	if spec.Root.Readonly ***REMOVED***
		return errors.New(`OCI spec is invalid - Root.Readonly must not be set on Windows`)
	***REMOVED***

	for _, layerPath := range layerFolders ***REMOVED***
		_, filename := filepath.Split(layerPath)
		g, err := hcsshim.NameToGuid(filename)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		configuration.Layers = append(configuration.Layers, hcsshim.Layer***REMOVED***
			ID:   g.ToString(),
			Path: layerPath,
		***REMOVED***)
	***REMOVED***

	// Add the mounts (volumes, bind mounts etc) to the structure
	var mds []hcsshim.MappedDir
	var mps []hcsshim.MappedPipe
	for _, mount := range spec.Mounts ***REMOVED***
		const pipePrefix = `\\.\pipe\`
		if mount.Type != "" ***REMOVED***
			return fmt.Errorf("OCI spec is invalid - Mount.Type '%s' must not be set", mount.Type)
		***REMOVED***
		if strings.HasPrefix(mount.Destination, pipePrefix) ***REMOVED***
			mp := hcsshim.MappedPipe***REMOVED***
				HostPath:          mount.Source,
				ContainerPipeName: mount.Destination[len(pipePrefix):],
			***REMOVED***
			mps = append(mps, mp)
		***REMOVED*** else ***REMOVED***
			md := hcsshim.MappedDir***REMOVED***
				HostPath:      mount.Source,
				ContainerPath: mount.Destination,
				ReadOnly:      false,
			***REMOVED***
			for _, o := range mount.Options ***REMOVED***
				if strings.ToLower(o) == "ro" ***REMOVED***
					md.ReadOnly = true
				***REMOVED***
			***REMOVED***
			mds = append(mds, md)
		***REMOVED***
	***REMOVED***
	configuration.MappedDirectories = mds
	if len(mps) > 0 && system.GetOSVersion().Build < 16210 ***REMOVED*** // replace with Win10 RS3 build number at RTM
		return errors.New("named pipe mounts are not supported on this version of Windows")
	***REMOVED***
	configuration.MappedPipes = mps

	hcsContainer, err := hcsshim.CreateContainer(id, configuration)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Construct a container object for calling start on it.
	ctr := &container***REMOVED***
		id:           id,
		execs:        make(map[string]*process),
		isWindows:    true,
		ociSpec:      spec,
		hcsContainer: hcsContainer,
		status:       StatusCreated,
		waitCh:       make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	// Start the container. If this is a servicing container, this call
	// will block until the container is done with the servicing
	// execution.
	logger.Debug("starting container")
	if err = hcsContainer.Start(); err != nil ***REMOVED***
		c.logger.WithError(err).Error("failed to start container")
		ctr.debugGCS()
		if err := c.terminateContainer(ctr); err != nil ***REMOVED***
			c.logger.WithError(err).Error("failed to cleanup after a failed Start")
		***REMOVED*** else ***REMOVED***
			c.logger.Debug("cleaned up after failed Start by calling Terminate")
		***REMOVED***
		return err
	***REMOVED***
	ctr.debugGCS()

	c.Lock()
	c.containers[id] = ctr
	c.Unlock()

	logger.Debug("createWindows() completed successfully")
	return nil

***REMOVED***

func (c *client) createLinux(id string, spec *specs.Spec, runtimeOptions interface***REMOVED******REMOVED***) error ***REMOVED***
	logrus.Debugf("libcontainerd: createLinux(): containerId %s ", id)
	logger := c.logger.WithField("container", id)

	if runtimeOptions == nil ***REMOVED***
		return fmt.Errorf("lcow option must be supplied to the runtime")
	***REMOVED***
	lcowConfig, ok := runtimeOptions.(*opengcs.Config)
	if !ok ***REMOVED***
		return fmt.Errorf("lcow option must be supplied to the runtime")
	***REMOVED***

	configuration := &hcsshim.ContainerConfig***REMOVED***
		HvPartition:   true,
		Name:          id,
		SystemType:    "container",
		ContainerType: "linux",
		Owner:         defaultOwner,
		TerminateOnLastHandleClosed: true,
	***REMOVED***

	if lcowConfig.ActualMode == opengcs.ModeActualVhdx ***REMOVED***
		configuration.HvRuntime = &hcsshim.HvRuntime***REMOVED***
			ImagePath:          lcowConfig.Vhdx,
			BootSource:         "Vhd",
			WritableBootSource: false,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		configuration.HvRuntime = &hcsshim.HvRuntime***REMOVED***
			ImagePath:           lcowConfig.KirdPath,
			LinuxKernelFile:     lcowConfig.KernelFile,
			LinuxInitrdFile:     lcowConfig.InitrdFile,
			LinuxBootParameters: lcowConfig.BootParameters,
		***REMOVED***
	***REMOVED***

	if spec.Windows == nil ***REMOVED***
		return fmt.Errorf("spec.Windows must not be nil for LCOW containers")
	***REMOVED***

	// We must have least one layer in the spec
	if spec.Windows.LayerFolders == nil || len(spec.Windows.LayerFolders) == 0 ***REMOVED***
		return fmt.Errorf("OCI spec is invalid - at least one LayerFolders must be supplied to the runtime")
	***REMOVED***

	// Strip off the top-most layer as that's passed in separately to HCS
	configuration.LayerFolderPath = spec.Windows.LayerFolders[len(spec.Windows.LayerFolders)-1]
	layerFolders := spec.Windows.LayerFolders[:len(spec.Windows.LayerFolders)-1]

	for _, layerPath := range layerFolders ***REMOVED***
		_, filename := filepath.Split(layerPath)
		g, err := hcsshim.NameToGuid(filename)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		configuration.Layers = append(configuration.Layers, hcsshim.Layer***REMOVED***
			ID:   g.ToString(),
			Path: filepath.Join(layerPath, "layer.vhd"),
		***REMOVED***)
	***REMOVED***

	if spec.Windows.Network != nil ***REMOVED***
		configuration.EndpointList = spec.Windows.Network.EndpointList
		configuration.AllowUnqualifiedDNSQuery = spec.Windows.Network.AllowUnqualifiedDNSQuery
		if spec.Windows.Network.DNSSearchList != nil ***REMOVED***
			configuration.DNSSearchList = strings.Join(spec.Windows.Network.DNSSearchList, ",")
		***REMOVED***
		configuration.NetworkSharedContainerName = spec.Windows.Network.NetworkSharedContainerName
	***REMOVED***

	// Add the mounts (volumes, bind mounts etc) to the structure. We have to do
	// some translation for both the mapped directories passed into HCS and in
	// the spec.
	//
	// For HCS, we only pass in the mounts from the spec which are type "bind".
	// Further, the "ContainerPath" field (which is a little mis-leadingly
	// named when it applies to the utility VM rather than the container in the
	// utility VM) is moved to under /tmp/gcs/<ID>/binds, where this is passed
	// by the caller through a 'uvmpath' option.
	//
	// We do similar translation for the mounts in the spec by stripping out
	// the uvmpath option, and translating the Source path to the location in the
	// utility VM calculated above.
	//
	// From inside the utility VM, you would see a 9p mount such as in the following
	// where a host folder has been mapped to /target. The line with /tmp/gcs/<ID>/binds
	// specifically:
	//
	//	/ # mount
	//	rootfs on / type rootfs (rw,size=463736k,nr_inodes=115934)
	//	proc on /proc type proc (rw,relatime)
	//	sysfs on /sys type sysfs (rw,relatime)
	//	udev on /dev type devtmpfs (rw,relatime,size=498100k,nr_inodes=124525,mode=755)
	//	tmpfs on /run type tmpfs (rw,relatime)
	//	cgroup on /sys/fs/cgroup type cgroup (rw,relatime,cpuset,cpu,cpuacct,blkio,memory,devices,freezer,net_cls,perf_event,net_prio,hugetlb,pids,rdma)
	//	mqueue on /dev/mqueue type mqueue (rw,relatime)
	//	devpts on /dev/pts type devpts (rw,relatime,mode=600,ptmxmode=000)
	//	/binds/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc/target on /binds/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc/target type 9p (rw,sync,dirsync,relatime,trans=fd,rfdno=6,wfdno=6)
	//	/dev/pmem0 on /tmp/gcs/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc/layer0 type ext4 (ro,relatime,block_validity,delalloc,norecovery,barrier,dax,user_xattr,acl)
	//	/dev/sda on /tmp/gcs/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc/scratch type ext4 (rw,relatime,block_validity,delalloc,barrier,user_xattr,acl)
	//	overlay on /tmp/gcs/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc/rootfs type overlay (rw,relatime,lowerdir=/tmp/base/:/tmp/gcs/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc/layer0,upperdir=/tmp/gcs/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc/scratch/upper,workdir=/tmp/gcs/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc/scratch/work)
	//
	//  /tmp/gcs/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc # ls -l
	//	total 16
	//	drwx------    3 0        0               60 Sep  7 18:54 binds
	//	-rw-r--r--    1 0        0             3345 Sep  7 18:54 config.json
	//	drwxr-xr-x   10 0        0             4096 Sep  6 17:26 layer0
	//	drwxr-xr-x    1 0        0             4096 Sep  7 18:54 rootfs
	//	drwxr-xr-x    5 0        0             4096 Sep  7 18:54 scratch
	//
	//	/tmp/gcs/b3ea9126d67702173647ece2744f7c11181c0150e9890fc9a431849838033edc # ls -l binds
	//	total 0
	//	drwxrwxrwt    2 0        0             4096 Sep  7 16:51 target

	mds := []hcsshim.MappedDir***REMOVED******REMOVED***
	specMounts := []specs.Mount***REMOVED******REMOVED***
	for _, mount := range spec.Mounts ***REMOVED***
		specMount := mount
		if mount.Type == "bind" ***REMOVED***
			// Strip out the uvmpath from the options
			updatedOptions := []string***REMOVED******REMOVED***
			uvmPath := ""
			readonly := false
			for _, opt := range mount.Options ***REMOVED***
				dropOption := false
				elements := strings.SplitN(opt, "=", 2)
				switch elements[0] ***REMOVED***
				case "uvmpath":
					uvmPath = elements[1]
					dropOption = true
				case "rw":
				case "ro":
					readonly = true
				case "rbind":
				default:
					return fmt.Errorf("unsupported option %q", opt)
				***REMOVED***
				if !dropOption ***REMOVED***
					updatedOptions = append(updatedOptions, opt)
				***REMOVED***
			***REMOVED***
			mount.Options = updatedOptions
			if uvmPath == "" ***REMOVED***
				return fmt.Errorf("no uvmpath for bind mount %+v", mount)
			***REMOVED***
			md := hcsshim.MappedDir***REMOVED***
				HostPath:          mount.Source,
				ContainerPath:     path.Join(uvmPath, mount.Destination),
				CreateInUtilityVM: true,
				ReadOnly:          readonly,
			***REMOVED***
			mds = append(mds, md)
			specMount.Source = path.Join(uvmPath, mount.Destination)
		***REMOVED***
		specMounts = append(specMounts, specMount)
	***REMOVED***
	configuration.MappedDirectories = mds

	hcsContainer, err := hcsshim.CreateContainer(id, configuration)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	spec.Mounts = specMounts

	// Construct a container object for calling start on it.
	ctr := &container***REMOVED***
		id:           id,
		execs:        make(map[string]*process),
		isWindows:    false,
		ociSpec:      spec,
		hcsContainer: hcsContainer,
		status:       StatusCreated,
		waitCh:       make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	// Start the container. If this is a servicing container, this call
	// will block until the container is done with the servicing
	// execution.
	logger.Debug("starting container")
	if err = hcsContainer.Start(); err != nil ***REMOVED***
		c.logger.WithError(err).Error("failed to start container")
		ctr.debugGCS()
		if err := c.terminateContainer(ctr); err != nil ***REMOVED***
			c.logger.WithError(err).Error("failed to cleanup after a failed Start")
		***REMOVED*** else ***REMOVED***
			c.logger.Debug("cleaned up after failed Start by calling Terminate")
		***REMOVED***
		return err
	***REMOVED***
	ctr.debugGCS()

	c.Lock()
	c.containers[id] = ctr
	c.Unlock()

	c.eventQ.append(id, func() ***REMOVED***
		ei := EventInfo***REMOVED***
			ContainerID: id,
		***REMOVED***
		c.logger.WithFields(logrus.Fields***REMOVED***
			"container": ctr.id,
			"event":     EventCreate,
		***REMOVED***).Info("sending event")
		err := c.backend.ProcessEvent(id, EventCreate, ei)
		if err != nil ***REMOVED***
			c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
				"container": id,
				"event":     EventCreate,
			***REMOVED***).Error("failed to process event")
		***REMOVED***
	***REMOVED***)

	logger.Debug("createLinux() completed successfully")
	return nil
***REMOVED***

func (c *client) Start(_ context.Context, id, _ string, withStdin bool, attachStdio StdioCallback) (int, error) ***REMOVED***
	ctr := c.getContainer(id)
	switch ***REMOVED***
	case ctr == nil:
		return -1, errors.WithStack(newNotFoundError("no such container"))
	case ctr.init != nil:
		return -1, errors.WithStack(newConflictError("container already started"))
	***REMOVED***

	logger := c.logger.WithField("container", id)

	// Note we always tell HCS to create stdout as it's required
	// regardless of '-i' or '-t' options, so that docker can always grab
	// the output through logs. We also tell HCS to always create stdin,
	// even if it's not used - it will be closed shortly. Stderr is only
	// created if it we're not -t.
	var (
		emulateConsole   bool
		createStdErrPipe bool
	)
	if ctr.ociSpec.Process != nil ***REMOVED***
		emulateConsole = ctr.ociSpec.Process.Terminal
		createStdErrPipe = !ctr.ociSpec.Process.Terminal && !ctr.ociSpec.Windows.Servicing
	***REMOVED***

	createProcessParms := &hcsshim.ProcessConfig***REMOVED***
		EmulateConsole:   emulateConsole,
		WorkingDirectory: ctr.ociSpec.Process.Cwd,
		CreateStdInPipe:  !ctr.ociSpec.Windows.Servicing,
		CreateStdOutPipe: !ctr.ociSpec.Windows.Servicing,
		CreateStdErrPipe: createStdErrPipe,
	***REMOVED***

	if ctr.ociSpec.Process != nil && ctr.ociSpec.Process.ConsoleSize != nil ***REMOVED***
		createProcessParms.ConsoleSize[0] = uint(ctr.ociSpec.Process.ConsoleSize.Height)
		createProcessParms.ConsoleSize[1] = uint(ctr.ociSpec.Process.ConsoleSize.Width)
	***REMOVED***

	// Configure the environment for the process
	createProcessParms.Environment = setupEnvironmentVariables(ctr.ociSpec.Process.Env)
	if ctr.isWindows ***REMOVED***
		createProcessParms.CommandLine = strings.Join(ctr.ociSpec.Process.Args, " ")
	***REMOVED*** else ***REMOVED***
		createProcessParms.CommandArgs = ctr.ociSpec.Process.Args
	***REMOVED***
	createProcessParms.User = ctr.ociSpec.Process.User.Username

	// LCOW requires the raw OCI spec passed through HCS and onwards to
	// GCS for the utility VM.
	if !ctr.isWindows ***REMOVED***
		ociBuf, err := json.Marshal(ctr.ociSpec)
		if err != nil ***REMOVED***
			return -1, err
		***REMOVED***
		ociRaw := json.RawMessage(ociBuf)
		createProcessParms.OCISpecification = &ociRaw
	***REMOVED***

	ctr.Lock()
	defer ctr.Unlock()

	// Start the command running in the container.
	newProcess, err := ctr.hcsContainer.CreateProcess(createProcessParms)
	if err != nil ***REMOVED***
		logger.WithError(err).Error("CreateProcess() failed")
		return -1, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if err := newProcess.Kill(); err != nil ***REMOVED***
				logger.WithError(err).Error("failed to kill process")
			***REMOVED***
			go func() ***REMOVED***
				if err := newProcess.Wait(); err != nil ***REMOVED***
					logger.WithError(err).Error("failed to wait for process")
				***REMOVED***
				if err := newProcess.Close(); err != nil ***REMOVED***
					logger.WithError(err).Error("failed to clean process resources")
				***REMOVED***
			***REMOVED***()
		***REMOVED***
	***REMOVED***()
	p := &process***REMOVED***
		hcsProcess: newProcess,
		id:         InitProcessName,
		pid:        newProcess.Pid(),
	***REMOVED***
	logger.WithField("pid", p.pid).Debug("init process started")

	// If this is a servicing container, wait on the process synchronously here and
	// if it succeeds, wait for it cleanly shutdown and merge into the parent container.
	if ctr.ociSpec.Windows.Servicing ***REMOVED***
		// reapProcess takes the lock
		ctr.Unlock()
		defer ctr.Lock()
		exitCode := c.reapProcess(ctr, p)

		if exitCode != 0 ***REMOVED***
			return -1, errors.Errorf("libcontainerd: servicing container %s returned non-zero exit code %d", ctr.id, exitCode)
		***REMOVED***

		return p.pid, nil
	***REMOVED***

	dio, err := newIOFromProcess(newProcess, ctr.ociSpec.Process.Terminal)
	if err != nil ***REMOVED***
		logger.WithError(err).Error("failed to get stdio pipes")
		return -1, err
	***REMOVED***
	_, err = attachStdio(dio)
	if err != nil ***REMOVED***
		logger.WithError(err).Error("failed to attache stdio")
		return -1, err
	***REMOVED***
	ctr.status = StatusRunning
	ctr.init = p

	// Spin up a go routine waiting for exit to handle cleanup
	go c.reapProcess(ctr, p)

	// Generate the associated event
	c.eventQ.append(id, func() ***REMOVED***
		ei := EventInfo***REMOVED***
			ContainerID: id,
			ProcessID:   InitProcessName,
			Pid:         uint32(p.pid),
		***REMOVED***
		c.logger.WithFields(logrus.Fields***REMOVED***
			"container":  ctr.id,
			"event":      EventStart,
			"event-info": ei,
		***REMOVED***).Info("sending event")
		err := c.backend.ProcessEvent(ei.ContainerID, EventStart, ei)
		if err != nil ***REMOVED***
			c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
				"container":  id,
				"event":      EventStart,
				"event-info": ei,
			***REMOVED***).Error("failed to process event")
		***REMOVED***
	***REMOVED***)
	logger.Debug("start() completed")
	return p.pid, nil
***REMOVED***

func newIOFromProcess(newProcess hcsshim.Process, terminal bool) (*cio.DirectIO, error) ***REMOVED***
	stdin, stdout, stderr, err := newProcess.Stdio()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	dio := cio.NewDirectIO(createStdInCloser(stdin, newProcess), nil, nil, terminal)

	// Convert io.ReadClosers to io.Readers
	if stdout != nil ***REMOVED***
		dio.Stdout = ioutil.NopCloser(&autoClosingReader***REMOVED***ReadCloser: stdout***REMOVED***)
	***REMOVED***
	if stderr != nil ***REMOVED***
		dio.Stderr = ioutil.NopCloser(&autoClosingReader***REMOVED***ReadCloser: stderr***REMOVED***)
	***REMOVED***
	return dio, nil
***REMOVED***

// Exec adds a process in an running container
func (c *client) Exec(ctx context.Context, containerID, processID string, spec *specs.Process, withStdin bool, attachStdio StdioCallback) (int, error) ***REMOVED***
	ctr := c.getContainer(containerID)
	switch ***REMOVED***
	case ctr == nil:
		return -1, errors.WithStack(newNotFoundError("no such container"))
	case ctr.hcsContainer == nil:
		return -1, errors.WithStack(newInvalidParameterError("container is not running"))
	case ctr.execs != nil && ctr.execs[processID] != nil:
		return -1, errors.WithStack(newConflictError("id already in use"))
	***REMOVED***
	logger := c.logger.WithFields(logrus.Fields***REMOVED***
		"container": containerID,
		"exec":      processID,
	***REMOVED***)

	// Note we always tell HCS to
	// create stdout as it's required regardless of '-i' or '-t' options, so that
	// docker can always grab the output through logs. We also tell HCS to always
	// create stdin, even if it's not used - it will be closed shortly. Stderr
	// is only created if it we're not -t.
	createProcessParms := hcsshim.ProcessConfig***REMOVED***
		CreateStdInPipe:  true,
		CreateStdOutPipe: true,
		CreateStdErrPipe: !spec.Terminal,
	***REMOVED***
	if spec.Terminal ***REMOVED***
		createProcessParms.EmulateConsole = true
		if spec.ConsoleSize != nil ***REMOVED***
			createProcessParms.ConsoleSize[0] = uint(spec.ConsoleSize.Height)
			createProcessParms.ConsoleSize[1] = uint(spec.ConsoleSize.Width)
		***REMOVED***
	***REMOVED***

	// Take working directory from the process to add if it is defined,
	// otherwise take from the first process.
	if spec.Cwd != "" ***REMOVED***
		createProcessParms.WorkingDirectory = spec.Cwd
	***REMOVED*** else ***REMOVED***
		createProcessParms.WorkingDirectory = ctr.ociSpec.Process.Cwd
	***REMOVED***

	// Configure the environment for the process
	createProcessParms.Environment = setupEnvironmentVariables(spec.Env)
	if ctr.isWindows ***REMOVED***
		createProcessParms.CommandLine = strings.Join(spec.Args, " ")
	***REMOVED*** else ***REMOVED***
		createProcessParms.CommandArgs = spec.Args
	***REMOVED***
	createProcessParms.User = spec.User.Username

	logger.Debugf("exec commandLine: %s", createProcessParms.CommandLine)

	// Start the command running in the container.
	newProcess, err := ctr.hcsContainer.CreateProcess(&createProcessParms)
	if err != nil ***REMOVED***
		logger.WithError(err).Errorf("exec's CreateProcess() failed")
		return -1, err
	***REMOVED***
	pid := newProcess.Pid()
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if err := newProcess.Kill(); err != nil ***REMOVED***
				logger.WithError(err).Error("failed to kill process")
			***REMOVED***
			go func() ***REMOVED***
				if err := newProcess.Wait(); err != nil ***REMOVED***
					logger.WithError(err).Error("failed to wait for process")
				***REMOVED***
				if err := newProcess.Close(); err != nil ***REMOVED***
					logger.WithError(err).Error("failed to clean process resources")
				***REMOVED***
			***REMOVED***()
		***REMOVED***
	***REMOVED***()

	dio, err := newIOFromProcess(newProcess, spec.Terminal)
	if err != nil ***REMOVED***
		logger.WithError(err).Error("failed to get stdio pipes")
		return -1, err
	***REMOVED***
	// Tell the engine to attach streams back to the client
	_, err = attachStdio(dio)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	p := &process***REMOVED***
		id:         processID,
		pid:        pid,
		hcsProcess: newProcess,
	***REMOVED***

	// Add the process to the container's list of processes
	ctr.Lock()
	ctr.execs[processID] = p
	ctr.Unlock()

	// Spin up a go routine waiting for exit to handle cleanup
	go c.reapProcess(ctr, p)

	c.eventQ.append(ctr.id, func() ***REMOVED***
		ei := EventInfo***REMOVED***
			ContainerID: ctr.id,
			ProcessID:   p.id,
			Pid:         uint32(p.pid),
		***REMOVED***
		c.logger.WithFields(logrus.Fields***REMOVED***
			"container":  ctr.id,
			"event":      EventExecAdded,
			"event-info": ei,
		***REMOVED***).Info("sending event")
		err := c.backend.ProcessEvent(ctr.id, EventExecAdded, ei)
		if err != nil ***REMOVED***
			c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
				"container":  ctr.id,
				"event":      EventExecAdded,
				"event-info": ei,
			***REMOVED***).Error("failed to process event")
		***REMOVED***
		err = c.backend.ProcessEvent(ctr.id, EventExecStarted, ei)
		if err != nil ***REMOVED***
			c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
				"container":  ctr.id,
				"event":      EventExecStarted,
				"event-info": ei,
			***REMOVED***).Error("failed to process event")
		***REMOVED***
	***REMOVED***)

	return pid, nil
***REMOVED***

// Signal handles `docker stop` on Windows. While Linux has support for
// the full range of signals, signals aren't really implemented on Windows.
// We fake supporting regular stop and -9 to force kill.
func (c *client) SignalProcess(_ context.Context, containerID, processID string, signal int) error ***REMOVED***
	ctr, p, err := c.getProcess(containerID, processID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	ctr.manualStopRequested = true

	logger := c.logger.WithFields(logrus.Fields***REMOVED***
		"container": containerID,
		"process":   processID,
		"pid":       p.pid,
		"signal":    signal,
	***REMOVED***)
	logger.Debug("Signal()")

	if processID == InitProcessName ***REMOVED***
		if syscall.Signal(signal) == syscall.SIGKILL ***REMOVED***
			// Terminate the compute system
			if err := ctr.hcsContainer.Terminate(); err != nil ***REMOVED***
				if !hcsshim.IsPending(err) ***REMOVED***
					logger.WithError(err).Error("failed to terminate hccshim container")
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Shut down the container
			if err := ctr.hcsContainer.Shutdown(); err != nil ***REMOVED***
				if !hcsshim.IsPending(err) && !hcsshim.IsAlreadyStopped(err) ***REMOVED***
					// ignore errors
					logger.WithError(err).Error("failed to shutdown hccshim container")
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return p.hcsProcess.Kill()
	***REMOVED***

	return nil
***REMOVED***

// Resize handles a CLI event to resize an interactive docker run or docker
// exec window.
func (c *client) ResizeTerminal(_ context.Context, containerID, processID string, width, height int) error ***REMOVED***
	_, p, err := c.getProcess(containerID, processID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.logger.WithFields(logrus.Fields***REMOVED***
		"container": containerID,
		"process":   processID,
		"height":    height,
		"width":     width,
		"pid":       p.pid,
	***REMOVED***).Debug("resizing")
	return p.hcsProcess.ResizeConsole(uint16(width), uint16(height))
***REMOVED***

func (c *client) CloseStdin(_ context.Context, containerID, processID string) error ***REMOVED***
	_, p, err := c.getProcess(containerID, processID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return p.hcsProcess.CloseStdin()
***REMOVED***

// Pause handles pause requests for containers
func (c *client) Pause(_ context.Context, containerID string) error ***REMOVED***
	ctr, _, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if ctr.ociSpec.Windows.HyperV == nil ***REMOVED***
		return errors.New("cannot pause Windows Server Containers")
	***REMOVED***

	ctr.Lock()
	defer ctr.Unlock()

	if err = ctr.hcsContainer.Pause(); err != nil ***REMOVED***
		return err
	***REMOVED***

	ctr.status = StatusPaused

	c.eventQ.append(containerID, func() ***REMOVED***
		err := c.backend.ProcessEvent(containerID, EventPaused, EventInfo***REMOVED***
			ContainerID: containerID,
			ProcessID:   InitProcessName,
		***REMOVED***)
		c.logger.WithFields(logrus.Fields***REMOVED***
			"container": ctr.id,
			"event":     EventPaused,
		***REMOVED***).Info("sending event")
		if err != nil ***REMOVED***
			c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
				"container": containerID,
				"event":     EventPaused,
			***REMOVED***).Error("failed to process event")
		***REMOVED***
	***REMOVED***)

	return nil
***REMOVED***

// Resume handles resume requests for containers
func (c *client) Resume(_ context.Context, containerID string) error ***REMOVED***
	ctr, _, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if ctr.ociSpec.Windows.HyperV == nil ***REMOVED***
		return errors.New("cannot resume Windows Server Containers")
	***REMOVED***

	ctr.Lock()
	defer ctr.Unlock()

	if err = ctr.hcsContainer.Resume(); err != nil ***REMOVED***
		return err
	***REMOVED***

	ctr.status = StatusRunning

	c.eventQ.append(containerID, func() ***REMOVED***
		err := c.backend.ProcessEvent(containerID, EventResumed, EventInfo***REMOVED***
			ContainerID: containerID,
			ProcessID:   InitProcessName,
		***REMOVED***)
		c.logger.WithFields(logrus.Fields***REMOVED***
			"container": ctr.id,
			"event":     EventResumed,
		***REMOVED***).Info("sending event")
		if err != nil ***REMOVED***
			c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
				"container": containerID,
				"event":     EventResumed,
			***REMOVED***).Error("failed to process event")
		***REMOVED***
	***REMOVED***)

	return nil
***REMOVED***

// Stats handles stats requests for containers
func (c *client) Stats(_ context.Context, containerID string) (*Stats, error) ***REMOVED***
	ctr, _, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	readAt := time.Now()
	s, err := ctr.hcsContainer.Statistics()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Stats***REMOVED***
		Read:     readAt,
		HCSStats: &s,
	***REMOVED***, nil
***REMOVED***

// Restore is the handler for restoring a container
func (c *client) Restore(ctx context.Context, id string, attachStdio StdioCallback) (bool, int, error) ***REMOVED***
	c.logger.WithField("container", id).Debug("restore()")

	// TODO Windows: On RS1, a re-attach isn't possible.
	// However, there is a scenario in which there is an issue.
	// Consider a background container. The daemon dies unexpectedly.
	// HCS will still have the compute service alive and running.
	// For consistence, we call in to shoot it regardless if HCS knows about it
	// We explicitly just log a warning if the terminate fails.
	// Then we tell the backend the container exited.
	if hc, err := hcsshim.OpenContainer(id); err == nil ***REMOVED***
		const terminateTimeout = time.Minute * 2
		err := hc.Terminate()

		if hcsshim.IsPending(err) ***REMOVED***
			err = hc.WaitTimeout(terminateTimeout)
		***REMOVED*** else if hcsshim.IsAlreadyStopped(err) ***REMOVED***
			err = nil
		***REMOVED***

		if err != nil ***REMOVED***
			c.logger.WithField("container", id).WithError(err).Debug("terminate failed on restore")
			return false, -1, err
		***REMOVED***
	***REMOVED***
	return false, -1, nil
***REMOVED***

// GetPidsForContainer returns a list of process IDs running in a container.
// Not used on Windows.
func (c *client) ListPids(_ context.Context, _ string) ([]uint32, error) ***REMOVED***
	return nil, errors.New("not implemented on Windows")
***REMOVED***

// Summary returns a summary of the processes running in a container.
// This is present in Windows to support docker top. In linux, the
// engine shells out to ps to get process information. On Windows, as
// the containers could be Hyper-V containers, they would not be
// visible on the container host. However, libcontainerd does have
// that information.
func (c *client) Summary(_ context.Context, containerID string) ([]Summary, error) ***REMOVED***
	ctr, _, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p, err := ctr.hcsContainer.ProcessList()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pl := make([]Summary, len(p))
	for i := range p ***REMOVED***
		pl[i] = Summary(p[i])
	***REMOVED***
	return pl, nil
***REMOVED***

func (c *client) DeleteTask(ctx context.Context, containerID string) (uint32, time.Time, error) ***REMOVED***
	ec := -1
	ctr := c.getContainer(containerID)
	if ctr == nil ***REMOVED***
		return uint32(ec), time.Now(), errors.WithStack(newNotFoundError("no such container"))
	***REMOVED***

	select ***REMOVED***
	case <-ctx.Done():
		return uint32(ec), time.Now(), errors.WithStack(ctx.Err())
	case <-ctr.waitCh:
	default:
		return uint32(ec), time.Now(), errors.New("container is not stopped")
	***REMOVED***

	ctr.Lock()
	defer ctr.Unlock()
	return ctr.exitCode, ctr.exitedAt, nil
***REMOVED***

func (c *client) Delete(_ context.Context, containerID string) error ***REMOVED***
	c.Lock()
	defer c.Unlock()
	ctr := c.containers[containerID]
	if ctr == nil ***REMOVED***
		return errors.WithStack(newNotFoundError("no such container"))
	***REMOVED***

	ctr.Lock()
	defer ctr.Unlock()

	switch ctr.status ***REMOVED***
	case StatusCreated:
		if err := c.shutdownContainer(ctr); err != nil ***REMOVED***
			return err
		***REMOVED***
		fallthrough
	case StatusStopped:
		delete(c.containers, containerID)
		return nil
	***REMOVED***

	return errors.WithStack(newInvalidParameterError("container is not stopped"))
***REMOVED***

func (c *client) Status(ctx context.Context, containerID string) (Status, error) ***REMOVED***
	c.Lock()
	defer c.Unlock()
	ctr := c.containers[containerID]
	if ctr == nil ***REMOVED***
		return StatusUnknown, errors.WithStack(newNotFoundError("no such container"))
	***REMOVED***

	ctr.Lock()
	defer ctr.Unlock()
	return ctr.status, nil
***REMOVED***

func (c *client) UpdateResources(ctx context.Context, containerID string, resources *Resources) error ***REMOVED***
	// Updating resource isn't supported on Windows
	// but we should return nil for enabling updating container
	return nil
***REMOVED***

func (c *client) CreateCheckpoint(ctx context.Context, containerID, checkpointDir string, exit bool) error ***REMOVED***
	return errors.New("Windows: Containers do not support checkpoints")
***REMOVED***

func (c *client) getContainer(id string) *container ***REMOVED***
	c.Lock()
	ctr := c.containers[id]
	c.Unlock()

	return ctr
***REMOVED***

func (c *client) getProcess(containerID, processID string) (*container, *process, error) ***REMOVED***
	ctr := c.getContainer(containerID)
	switch ***REMOVED***
	case ctr == nil:
		return nil, nil, errors.WithStack(newNotFoundError("no such container"))
	case ctr.init == nil:
		return nil, nil, errors.WithStack(newNotFoundError("container is not running"))
	case processID == InitProcessName:
		return ctr, ctr.init, nil
	default:
		ctr.Lock()
		defer ctr.Unlock()
		if ctr.execs == nil ***REMOVED***
			return nil, nil, errors.WithStack(newNotFoundError("no execs"))
		***REMOVED***
	***REMOVED***

	p := ctr.execs[processID]
	if p == nil ***REMOVED***
		return nil, nil, errors.WithStack(newNotFoundError("no such exec"))
	***REMOVED***

	return ctr, p, nil
***REMOVED***

func (c *client) shutdownContainer(ctr *container) error ***REMOVED***
	const shutdownTimeout = time.Minute * 5
	err := ctr.hcsContainer.Shutdown()

	if hcsshim.IsPending(err) ***REMOVED***
		err = ctr.hcsContainer.WaitTimeout(shutdownTimeout)
	***REMOVED*** else if hcsshim.IsAlreadyStopped(err) ***REMOVED***
		err = nil
	***REMOVED***

	if err != nil ***REMOVED***
		c.logger.WithError(err).WithField("container", ctr.id).
			Debug("failed to shutdown container, terminating it")
		return c.terminateContainer(ctr)
	***REMOVED***

	return nil
***REMOVED***

func (c *client) terminateContainer(ctr *container) error ***REMOVED***
	const terminateTimeout = time.Minute * 5
	err := ctr.hcsContainer.Terminate()

	if hcsshim.IsPending(err) ***REMOVED***
		err = ctr.hcsContainer.WaitTimeout(terminateTimeout)
	***REMOVED*** else if hcsshim.IsAlreadyStopped(err) ***REMOVED***
		err = nil
	***REMOVED***

	if err != nil ***REMOVED***
		c.logger.WithError(err).WithField("container", ctr.id).
			Debug("failed to terminate container")
		return err
	***REMOVED***

	return nil
***REMOVED***

func (c *client) reapProcess(ctr *container, p *process) int ***REMOVED***
	logger := c.logger.WithFields(logrus.Fields***REMOVED***
		"container": ctr.id,
		"process":   p.id,
	***REMOVED***)

	// Block indefinitely for the process to exit.
	if err := p.hcsProcess.Wait(); err != nil ***REMOVED***
		if herr, ok := err.(*hcsshim.ProcessError); ok && herr.Err != windows.ERROR_BROKEN_PIPE ***REMOVED***
			logger.WithError(err).Warnf("Wait() failed (container may have been killed)")
		***REMOVED***
		// Fall through here, do not return. This ensures we attempt to
		// continue the shutdown in HCS and tell the docker engine that the
		// process/container has exited to avoid a container being dropped on
		// the floor.
	***REMOVED***
	exitedAt := time.Now()

	exitCode, err := p.hcsProcess.ExitCode()
	if err != nil ***REMOVED***
		if herr, ok := err.(*hcsshim.ProcessError); ok && herr.Err != windows.ERROR_BROKEN_PIPE ***REMOVED***
			logger.WithError(err).Warnf("unable to get exit code for process")
		***REMOVED***
		// Since we got an error retrieving the exit code, make sure that the
		// code we return doesn't incorrectly indicate success.
		exitCode = -1

		// Fall through here, do not return. This ensures we attempt to
		// continue the shutdown in HCS and tell the docker engine that the
		// process/container has exited to avoid a container being dropped on
		// the floor.
	***REMOVED***

	if err := p.hcsProcess.Close(); err != nil ***REMOVED***
		logger.WithError(err).Warnf("failed to cleanup hcs process resources")
	***REMOVED***

	var pendingUpdates bool
	if p.id == InitProcessName ***REMOVED***
		// Update container status
		ctr.Lock()
		ctr.status = StatusStopped
		ctr.exitedAt = exitedAt
		ctr.exitCode = uint32(exitCode)
		close(ctr.waitCh)
		ctr.Unlock()

		// Handle any servicing
		if exitCode == 0 && ctr.isWindows && !ctr.ociSpec.Windows.Servicing ***REMOVED***
			pendingUpdates, err = ctr.hcsContainer.HasPendingUpdates()
			logger.Infof("Pending updates: %v", pendingUpdates)
			if err != nil ***REMOVED***
				logger.WithError(err).
					Warnf("failed to check for pending updates (container may have been killed)")
			***REMOVED***
		***REMOVED***

		if err := c.shutdownContainer(ctr); err != nil ***REMOVED***
			logger.WithError(err).Warn("failed to shutdown container")
		***REMOVED*** else ***REMOVED***
			logger.Debug("completed container shutdown")
		***REMOVED***

		if err := ctr.hcsContainer.Close(); err != nil ***REMOVED***
			logger.WithError(err).Error("failed to clean hcs container resources")
		***REMOVED***
	***REMOVED***

	if !(ctr.isWindows && ctr.ociSpec.Windows.Servicing) ***REMOVED***
		c.eventQ.append(ctr.id, func() ***REMOVED***
			ei := EventInfo***REMOVED***
				ContainerID:   ctr.id,
				ProcessID:     p.id,
				Pid:           uint32(p.pid),
				ExitCode:      uint32(exitCode),
				ExitedAt:      exitedAt,
				UpdatePending: pendingUpdates,
			***REMOVED***
			c.logger.WithFields(logrus.Fields***REMOVED***
				"container":  ctr.id,
				"event":      EventExit,
				"event-info": ei,
			***REMOVED***).Info("sending event")
			err := c.backend.ProcessEvent(ctr.id, EventExit, ei)
			if err != nil ***REMOVED***
				c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
					"container":  ctr.id,
					"event":      EventExit,
					"event-info": ei,
				***REMOVED***).Error("failed to process event")
			***REMOVED***
			if p.id != InitProcessName ***REMOVED***
				ctr.Lock()
				delete(ctr.execs, p.id)
				ctr.Unlock()
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	return exitCode
***REMOVED***
