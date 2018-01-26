package v2

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/oci"
	"github.com/docker/docker/pkg/system"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

// InitSpec creates an OCI spec from the plugin's config.
func (p *Plugin) InitSpec(execRoot string) (*specs.Spec, error) ***REMOVED***
	s := oci.DefaultSpec()
	s.Root = &specs.Root***REMOVED***
		Path:     p.Rootfs,
		Readonly: false, // TODO: all plugins should be readonly? settable in config?
	***REMOVED***

	userMounts := make(map[string]struct***REMOVED******REMOVED***, len(p.PluginObj.Settings.Mounts))
	for _, m := range p.PluginObj.Settings.Mounts ***REMOVED***
		userMounts[m.Destination] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	execRoot = filepath.Join(execRoot, p.PluginObj.ID)
	if err := os.MkdirAll(execRoot, 0700); err != nil ***REMOVED***
		return nil, errors.WithStack(err)
	***REMOVED***

	mounts := append(p.PluginObj.Config.Mounts, types.PluginMount***REMOVED***
		Source:      &execRoot,
		Destination: defaultPluginRuntimeDestination,
		Type:        "bind",
		Options:     []string***REMOVED***"rbind", "rshared"***REMOVED***,
	***REMOVED***)

	if p.PluginObj.Config.Network.Type != "" ***REMOVED***
		// TODO: if net == bridge, use libnetwork controller to create a new plugin-specific bridge, bind mount /etc/hosts and /etc/resolv.conf look at the docker code (allocateNetwork, initialize)
		if p.PluginObj.Config.Network.Type == "host" ***REMOVED***
			oci.RemoveNamespace(&s, specs.LinuxNamespaceType("network"))
		***REMOVED***
		etcHosts := "/etc/hosts"
		resolvConf := "/etc/resolv.conf"
		mounts = append(mounts,
			types.PluginMount***REMOVED***
				Source:      &etcHosts,
				Destination: etcHosts,
				Type:        "bind",
				Options:     []string***REMOVED***"rbind", "ro"***REMOVED***,
			***REMOVED***,
			types.PluginMount***REMOVED***
				Source:      &resolvConf,
				Destination: resolvConf,
				Type:        "bind",
				Options:     []string***REMOVED***"rbind", "ro"***REMOVED***,
			***REMOVED***)
	***REMOVED***
	if p.PluginObj.Config.PidHost ***REMOVED***
		oci.RemoveNamespace(&s, specs.LinuxNamespaceType("pid"))
	***REMOVED***

	if p.PluginObj.Config.IpcHost ***REMOVED***
		oci.RemoveNamespace(&s, specs.LinuxNamespaceType("ipc"))
	***REMOVED***

	for _, mnt := range mounts ***REMOVED***
		m := specs.Mount***REMOVED***
			Destination: mnt.Destination,
			Type:        mnt.Type,
			Options:     mnt.Options,
		***REMOVED***
		if mnt.Source == nil ***REMOVED***
			return nil, errors.New("mount source is not specified")
		***REMOVED***
		m.Source = *mnt.Source
		s.Mounts = append(s.Mounts, m)
	***REMOVED***

	for i, m := range s.Mounts ***REMOVED***
		if strings.HasPrefix(m.Destination, "/dev/") ***REMOVED***
			if _, ok := userMounts[m.Destination]; ok ***REMOVED***
				s.Mounts = append(s.Mounts[:i], s.Mounts[i+1:]...)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if p.PluginObj.Config.PropagatedMount != "" ***REMOVED***
		p.PropagatedMount = filepath.Join(p.Rootfs, p.PluginObj.Config.PropagatedMount)
		s.Linux.RootfsPropagation = "rshared"
	***REMOVED***

	if p.PluginObj.Config.Linux.AllowAllDevices ***REMOVED***
		s.Linux.Resources.Devices = []specs.LinuxDeviceCgroup***REMOVED******REMOVED***Allow: true, Access: "rwm"***REMOVED******REMOVED***
	***REMOVED***
	for _, dev := range p.PluginObj.Settings.Devices ***REMOVED***
		path := *dev.Path
		d, dPermissions, err := oci.DevicesFromPath(path, path, "rwm")
		if err != nil ***REMOVED***
			return nil, errors.WithStack(err)
		***REMOVED***
		s.Linux.Devices = append(s.Linux.Devices, d...)
		s.Linux.Resources.Devices = append(s.Linux.Resources.Devices, dPermissions...)
	***REMOVED***

	envs := make([]string, 1, len(p.PluginObj.Settings.Env)+1)
	envs[0] = "PATH=" + system.DefaultPathEnv(runtime.GOOS)
	envs = append(envs, p.PluginObj.Settings.Env...)

	args := append(p.PluginObj.Config.Entrypoint, p.PluginObj.Settings.Args...)
	cwd := p.PluginObj.Config.WorkDir
	if len(cwd) == 0 ***REMOVED***
		cwd = "/"
	***REMOVED***
	s.Process.Terminal = false
	s.Process.Args = args
	s.Process.Cwd = cwd
	s.Process.Env = envs

	caps := s.Process.Capabilities
	caps.Bounding = append(caps.Bounding, p.PluginObj.Config.Linux.Capabilities...)
	caps.Permitted = append(caps.Permitted, p.PluginObj.Config.Linux.Capabilities...)
	caps.Inheritable = append(caps.Inheritable, p.PluginObj.Config.Linux.Capabilities...)
	caps.Effective = append(caps.Effective, p.PluginObj.Config.Linux.Capabilities...)

	return &s, nil
***REMOVED***
