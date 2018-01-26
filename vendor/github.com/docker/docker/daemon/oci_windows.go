package daemon

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/oci"
	"github.com/docker/docker/pkg/sysinfo"
	"github.com/docker/docker/pkg/system"
	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	credentialSpecRegistryLocation = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Virtualization\Containers\CredentialSpecs`
	credentialSpecFileLocation     = "CredentialSpecs"
)

func (daemon *Daemon) createSpec(c *container.Container) (*specs.Spec, error) ***REMOVED***
	img, err := daemon.GetImage(string(c.ImageID))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	s := oci.DefaultOSSpec(img.OS)

	linkedEnv, err := daemon.setupLinkedContainers(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Note, unlike Unix, we do NOT call into SetupWorkingDirectory as
	// this is done in VMCompute. Further, we couldn't do it for Hyper-V
	// containers anyway.

	// In base spec
	s.Hostname = c.FullHostname()

	if err := daemon.setupSecretDir(c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := daemon.setupConfigDir(c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// In s.Mounts
	mounts, err := daemon.setupMounts(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var isHyperV bool
	if c.HostConfig.Isolation.IsDefault() ***REMOVED***
		// Container using default isolation, so take the default from the daemon configuration
		isHyperV = daemon.defaultIsolation.IsHyperV()
	***REMOVED*** else ***REMOVED***
		// Container may be requesting an explicit isolation mode.
		isHyperV = c.HostConfig.Isolation.IsHyperV()
	***REMOVED***

	if isHyperV ***REMOVED***
		s.Windows.HyperV = &specs.WindowsHyperV***REMOVED******REMOVED***
	***REMOVED***

	// If the container has not been started, and has configs or secrets
	// secrets, create symlinks to each config and secret. If it has been
	// started before, the symlinks should have already been created. Also, it
	// is important to not mount a Hyper-V  container that has been started
	// before, to protect the host from the container; for example, from
	// malicious mutation of NTFS data structures.
	if !c.HasBeenStartedBefore && (len(c.SecretReferences) > 0 || len(c.ConfigReferences) > 0) ***REMOVED***
		// The container file system is mounted before this function is called,
		// except for Hyper-V containers, so mount it here in that case.
		if isHyperV ***REMOVED***
			if err := daemon.Mount(c); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			defer daemon.Unmount(c)
		***REMOVED***
		if err := c.CreateSecretSymlinks(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if err := c.CreateConfigSymlinks(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	secretMounts, err := c.SecretMounts()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if secretMounts != nil ***REMOVED***
		mounts = append(mounts, secretMounts...)
	***REMOVED***

	configMounts, err := c.ConfigMounts()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if configMounts != nil ***REMOVED***
		mounts = append(mounts, configMounts...)
	***REMOVED***

	for _, mount := range mounts ***REMOVED***
		m := specs.Mount***REMOVED***
			Source:      mount.Source,
			Destination: mount.Destination,
		***REMOVED***
		if !mount.Writable ***REMOVED***
			m.Options = append(m.Options, "ro")
		***REMOVED***
		if img.OS != runtime.GOOS ***REMOVED***
			m.Type = "bind"
			m.Options = append(m.Options, "rbind")
			m.Options = append(m.Options, fmt.Sprintf("uvmpath=/tmp/gcs/%s/binds", c.ID))
		***REMOVED***
		s.Mounts = append(s.Mounts, m)
	***REMOVED***

	// In s.Process
	s.Process.Args = append([]string***REMOVED***c.Path***REMOVED***, c.Args...)
	if !c.Config.ArgsEscaped && img.OS == "windows" ***REMOVED***
		s.Process.Args = escapeArgs(s.Process.Args)
	***REMOVED***

	s.Process.Cwd = c.Config.WorkingDir
	s.Process.Env = c.CreateDaemonEnvironment(c.Config.Tty, linkedEnv)
	if c.Config.Tty ***REMOVED***
		s.Process.Terminal = c.Config.Tty
		s.Process.ConsoleSize = &specs.Box***REMOVED***
			Height: c.HostConfig.ConsoleSize[0],
			Width:  c.HostConfig.ConsoleSize[1],
		***REMOVED***
	***REMOVED***
	s.Process.User.Username = c.Config.User

	// Get the layer path for each layer.
	max := len(img.RootFS.DiffIDs)
	for i := 1; i <= max; i++ ***REMOVED***
		img.RootFS.DiffIDs = img.RootFS.DiffIDs[:i]
		if !system.IsOSSupported(img.OperatingSystem()) ***REMOVED***
			return nil, fmt.Errorf("cannot get layerpath for ImageID %s: %s ", img.RootFS.ChainID(), system.ErrNotSupportedOperatingSystem)
		***REMOVED***
		layerPath, err := layer.GetLayerPath(daemon.layerStores[img.OperatingSystem()], img.RootFS.ChainID())
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to get layer path from graphdriver %s for ImageID %s - %s", daemon.layerStores[img.OperatingSystem()], img.RootFS.ChainID(), err)
		***REMOVED***
		// Reverse order, expecting parent most first
		s.Windows.LayerFolders = append([]string***REMOVED***layerPath***REMOVED***, s.Windows.LayerFolders...)
	***REMOVED***
	m, err := c.RWLayer.Metadata()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get layer metadata - %s", err)
	***REMOVED***
	s.Windows.LayerFolders = append(s.Windows.LayerFolders, m["dir"])

	dnsSearch := daemon.getDNSSearchSettings(c)

	// Get endpoints for the libnetwork allocated networks to the container
	var epList []string
	AllowUnqualifiedDNSQuery := false
	gwHNSID := ""
	if c.NetworkSettings != nil ***REMOVED***
		for n := range c.NetworkSettings.Networks ***REMOVED***
			sn, err := daemon.FindNetwork(n)
			if err != nil ***REMOVED***
				continue
			***REMOVED***

			ep, err := c.GetEndpointInNetwork(sn)
			if err != nil ***REMOVED***
				continue
			***REMOVED***

			data, err := ep.DriverInfo()
			if err != nil ***REMOVED***
				continue
			***REMOVED***

			if data["GW_INFO"] != nil ***REMOVED***
				gwInfo := data["GW_INFO"].(map[string]interface***REMOVED******REMOVED***)
				if gwInfo["hnsid"] != nil ***REMOVED***
					gwHNSID = gwInfo["hnsid"].(string)
				***REMOVED***
			***REMOVED***

			if data["hnsid"] != nil ***REMOVED***
				epList = append(epList, data["hnsid"].(string))
			***REMOVED***

			if data["AllowUnqualifiedDNSQuery"] != nil ***REMOVED***
				AllowUnqualifiedDNSQuery = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var networkSharedContainerID string
	if c.HostConfig.NetworkMode.IsContainer() ***REMOVED***
		networkSharedContainerID = c.NetworkSharedContainerID
		for _, ep := range c.SharedEndpointList ***REMOVED***
			epList = append(epList, ep)
		***REMOVED***
	***REMOVED***

	if gwHNSID != "" ***REMOVED***
		epList = append(epList, gwHNSID)
	***REMOVED***

	s.Windows.Network = &specs.WindowsNetwork***REMOVED***
		AllowUnqualifiedDNSQuery:   AllowUnqualifiedDNSQuery,
		DNSSearchList:              dnsSearch,
		EndpointList:               epList,
		NetworkSharedContainerName: networkSharedContainerID,
	***REMOVED***

	switch img.OS ***REMOVED***
	case "windows":
		if err := daemon.createSpecWindowsFields(c, &s, isHyperV); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	case "linux":
		if !system.LCOWSupported() ***REMOVED***
			return nil, fmt.Errorf("Linux containers on Windows are not supported")
		***REMOVED***
		daemon.createSpecLinuxFields(c, &s)
	default:
		return nil, fmt.Errorf("Unsupported platform %q", img.OS)
	***REMOVED***

	return (*specs.Spec)(&s), nil
***REMOVED***

// Sets the Windows-specific fields of the OCI spec
func (daemon *Daemon) createSpecWindowsFields(c *container.Container, s *specs.Spec, isHyperV bool) error ***REMOVED***
	if len(s.Process.Cwd) == 0 ***REMOVED***
		// We default to C:\ to workaround the oddity of the case that the
		// default directory for cmd running as LocalSystem (or
		// ContainerAdministrator) is c:\windows\system32. Hence docker run
		// <image> cmd will by default end in c:\windows\system32, rather
		// than 'root' (/) on Linux. The oddity is that if you have a dockerfile
		// which has no WORKDIR and has a COPY file ., . will be interpreted
		// as c:\. Hence, setting it to default of c:\ makes for consistency.
		s.Process.Cwd = `C:\`
	***REMOVED***

	s.Root.Readonly = false // Windows does not support a read-only root filesystem
	if !isHyperV ***REMOVED***
		s.Root.Path = c.BaseFS.Path() // This is not set for Hyper-V containers
		if !strings.HasSuffix(s.Root.Path, `\`) ***REMOVED***
			s.Root.Path = s.Root.Path + `\` // Ensure a correctly formatted volume GUID path \\?\Volume***REMOVED***GUID***REMOVED***\
		***REMOVED***
	***REMOVED***

	// First boot optimization
	s.Windows.IgnoreFlushesDuringBoot = !c.HasBeenStartedBefore

	// In s.Windows.Resources
	cpuShares := uint16(c.HostConfig.CPUShares)
	cpuMaximum := uint16(c.HostConfig.CPUPercent) * 100
	cpuCount := uint64(c.HostConfig.CPUCount)
	if c.HostConfig.NanoCPUs > 0 ***REMOVED***
		if isHyperV ***REMOVED***
			cpuCount = uint64(c.HostConfig.NanoCPUs / 1e9)
			leftoverNanoCPUs := c.HostConfig.NanoCPUs % 1e9
			if leftoverNanoCPUs != 0 ***REMOVED***
				cpuCount++
				cpuMaximum = uint16(c.HostConfig.NanoCPUs / int64(cpuCount) / (1e9 / 10000))
				if cpuMaximum < 1 ***REMOVED***
					// The requested NanoCPUs is so small that we rounded to 0, use 1 instead
					cpuMaximum = 1
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			cpuMaximum = uint16(c.HostConfig.NanoCPUs / int64(sysinfo.NumCPU()) / (1e9 / 10000))
			if cpuMaximum < 1 ***REMOVED***
				// The requested NanoCPUs is so small that we rounded to 0, use 1 instead
				cpuMaximum = 1
			***REMOVED***
		***REMOVED***
	***REMOVED***
	memoryLimit := uint64(c.HostConfig.Memory)
	s.Windows.Resources = &specs.WindowsResources***REMOVED***
		CPU: &specs.WindowsCPUResources***REMOVED***
			Maximum: &cpuMaximum,
			Shares:  &cpuShares,
			Count:   &cpuCount,
		***REMOVED***,
		Memory: &specs.WindowsMemoryResources***REMOVED***
			Limit: &memoryLimit,
		***REMOVED***,
		Storage: &specs.WindowsStorageResources***REMOVED***
			Bps:  &c.HostConfig.IOMaximumBandwidth,
			Iops: &c.HostConfig.IOMaximumIOps,
		***REMOVED***,
	***REMOVED***

	// Read and add credentials from the security options if a credential spec has been provided.
	if c.HostConfig.SecurityOpt != nil ***REMOVED***
		cs := ""
		for _, sOpt := range c.HostConfig.SecurityOpt ***REMOVED***
			sOpt = strings.ToLower(sOpt)
			if !strings.Contains(sOpt, "=") ***REMOVED***
				return fmt.Errorf("invalid security option: no equals sign in supplied value %s", sOpt)
			***REMOVED***
			var splitsOpt []string
			splitsOpt = strings.SplitN(sOpt, "=", 2)
			if len(splitsOpt) != 2 ***REMOVED***
				return fmt.Errorf("invalid security option: %s", sOpt)
			***REMOVED***
			if splitsOpt[0] != "credentialspec" ***REMOVED***
				return fmt.Errorf("security option not supported: %s", splitsOpt[0])
			***REMOVED***

			var (
				match   bool
				csValue string
				err     error
			)
			if match, csValue = getCredentialSpec("file://", splitsOpt[1]); match ***REMOVED***
				if csValue == "" ***REMOVED***
					return fmt.Errorf("no value supplied for file:// credential spec security option")
				***REMOVED***
				if cs, err = readCredentialSpecFile(c.ID, daemon.root, filepath.Clean(csValue)); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else if match, csValue = getCredentialSpec("registry://", splitsOpt[1]); match ***REMOVED***
				if csValue == "" ***REMOVED***
					return fmt.Errorf("no value supplied for registry:// credential spec security option")
				***REMOVED***
				if cs, err = readCredentialSpecRegistry(c.ID, csValue); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				return fmt.Errorf("invalid credential spec security option - value must be prefixed file:// or registry:// followed by a value")
			***REMOVED***
		***REMOVED***
		s.Windows.CredentialSpec = cs
	***REMOVED***

	// Assume we are not starting a container for a servicing operation
	s.Windows.Servicing = false

	return nil
***REMOVED***

// Sets the Linux-specific fields of the OCI spec
// TODO: @jhowardmsft LCOW Support. We need to do a lot more pulling in what can
// be pulled in from oci_linux.go.
func (daemon *Daemon) createSpecLinuxFields(c *container.Container, s *specs.Spec) ***REMOVED***
	if len(s.Process.Cwd) == 0 ***REMOVED***
		s.Process.Cwd = `/`
	***REMOVED***
	s.Root.Path = "rootfs"
	s.Root.Readonly = c.HostConfig.ReadonlyRootfs
***REMOVED***

func escapeArgs(args []string) []string ***REMOVED***
	escapedArgs := make([]string, len(args))
	for i, a := range args ***REMOVED***
		escapedArgs[i] = windows.EscapeArg(a)
	***REMOVED***
	return escapedArgs
***REMOVED***

// mergeUlimits merge the Ulimits from HostConfig with daemon defaults, and update HostConfig
// It will do nothing on non-Linux platform
func (daemon *Daemon) mergeUlimits(c *containertypes.HostConfig) ***REMOVED***
	return
***REMOVED***

// getCredentialSpec is a helper function to get the value of a credential spec supplied
// on the CLI, stripping the prefix
func getCredentialSpec(prefix, value string) (bool, string) ***REMOVED***
	if strings.HasPrefix(value, prefix) ***REMOVED***
		return true, strings.TrimPrefix(value, prefix)
	***REMOVED***
	return false, ""
***REMOVED***

// readCredentialSpecRegistry is a helper function to read a credential spec from
// the registry. If not found, we return an empty string and warn in the log.
// This allows for staging on machines which do not have the necessary components.
func readCredentialSpecRegistry(id, name string) (string, error) ***REMOVED***
	var (
		k   registry.Key
		err error
		val string
	)
	if k, err = registry.OpenKey(registry.LOCAL_MACHINE, credentialSpecRegistryLocation, registry.QUERY_VALUE); err != nil ***REMOVED***
		return "", fmt.Errorf("failed handling spec %q for container %s - %s could not be opened", name, id, credentialSpecRegistryLocation)
	***REMOVED***
	if val, _, err = k.GetStringValue(name); err != nil ***REMOVED***
		if err == registry.ErrNotExist ***REMOVED***
			return "", fmt.Errorf("credential spec %q for container %s as it was not found", name, id)
		***REMOVED***
		return "", fmt.Errorf("error %v reading credential spec %q from registry for container %s", err, name, id)
	***REMOVED***
	return val, nil
***REMOVED***

// readCredentialSpecFile is a helper function to read a credential spec from
// a file. If not found, we return an empty string and warn in the log.
// This allows for staging on machines which do not have the necessary components.
func readCredentialSpecFile(id, root, location string) (string, error) ***REMOVED***
	if filepath.IsAbs(location) ***REMOVED***
		return "", fmt.Errorf("invalid credential spec - file:// path cannot be absolute")
	***REMOVED***
	base := filepath.Join(root, credentialSpecFileLocation)
	full := filepath.Join(base, location)
	if !strings.HasPrefix(full, base) ***REMOVED***
		return "", fmt.Errorf("invalid credential spec - file:// path must be under %s", base)
	***REMOVED***
	bcontents, err := ioutil.ReadFile(full)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("credential spec '%s' for container %s as the file could not be read: %q", full, id, err)
	***REMOVED***
	return string(bcontents[:]), nil
***REMOVED***
