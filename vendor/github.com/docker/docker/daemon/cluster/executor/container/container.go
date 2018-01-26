package container

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	enginecontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	enginemount "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/daemon/cluster/convert"
	executorpkg "github.com/docker/docker/daemon/cluster/executor"
	clustertypes "github.com/docker/docker/daemon/cluster/provider"
	"github.com/docker/go-connections/nat"
	netconst "github.com/docker/libnetwork/datastore"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/genericresource"
	"github.com/docker/swarmkit/template"
	gogotypes "github.com/gogo/protobuf/types"
)

const (
	// Explicitly use the kernel's default setting for CPU quota of 100ms.
	// https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt
	cpuQuotaPeriod = 100 * time.Millisecond

	// systemLabelPrefix represents the reserved namespace for system labels.
	systemLabelPrefix = "com.docker.swarm"
)

// containerConfig converts task properties into docker container compatible
// components.
type containerConfig struct ***REMOVED***
	task                *api.Task
	networksAttachments map[string]*api.NetworkAttachment
***REMOVED***

// newContainerConfig returns a validated container config. No methods should
// return an error if this function returns without error.
func newContainerConfig(t *api.Task, node *api.NodeDescription) (*containerConfig, error) ***REMOVED***
	var c containerConfig
	return &c, c.setTask(t, node)
***REMOVED***

func (c *containerConfig) setTask(t *api.Task, node *api.NodeDescription) error ***REMOVED***
	if t.Spec.GetContainer() == nil && t.Spec.GetAttachment() == nil ***REMOVED***
		return exec.ErrRuntimeUnsupported
	***REMOVED***

	container := t.Spec.GetContainer()
	if container != nil ***REMOVED***
		if container.Image == "" ***REMOVED***
			return ErrImageRequired
		***REMOVED***

		if err := validateMounts(container.Mounts); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// index the networks by name
	c.networksAttachments = make(map[string]*api.NetworkAttachment, len(t.Networks))
	for _, attachment := range t.Networks ***REMOVED***
		c.networksAttachments[attachment.Network.Spec.Annotations.Name] = attachment
	***REMOVED***

	c.task = t

	if t.Spec.GetContainer() != nil ***REMOVED***
		preparedSpec, err := template.ExpandContainerSpec(node, t)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.task.Spec.Runtime = &api.TaskSpec_Container***REMOVED***
			Container: preparedSpec,
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *containerConfig) networkAttachmentContainerID() string ***REMOVED***
	attachment := c.task.Spec.GetAttachment()
	if attachment == nil ***REMOVED***
		return ""
	***REMOVED***

	return attachment.ContainerID
***REMOVED***

func (c *containerConfig) taskID() string ***REMOVED***
	return c.task.ID
***REMOVED***

func (c *containerConfig) endpoint() *api.Endpoint ***REMOVED***
	return c.task.Endpoint
***REMOVED***

func (c *containerConfig) spec() *api.ContainerSpec ***REMOVED***
	return c.task.Spec.GetContainer()
***REMOVED***

func (c *containerConfig) nameOrID() string ***REMOVED***
	if c.task.Spec.GetContainer() != nil ***REMOVED***
		return c.name()
	***REMOVED***

	return c.networkAttachmentContainerID()
***REMOVED***

func (c *containerConfig) name() string ***REMOVED***
	if c.task.Annotations.Name != "" ***REMOVED***
		// if set, use the container Annotations.Name field, set in the orchestrator.
		return c.task.Annotations.Name
	***REMOVED***

	slot := fmt.Sprint(c.task.Slot)
	if slot == "" || c.task.Slot == 0 ***REMOVED***
		slot = c.task.NodeID
	***REMOVED***

	// fallback to service.slot.id.
	return fmt.Sprintf("%s.%s.%s", c.task.ServiceAnnotations.Name, slot, c.task.ID)
***REMOVED***

func (c *containerConfig) image() string ***REMOVED***
	raw := c.spec().Image
	ref, err := reference.ParseNormalizedNamed(raw)
	if err != nil ***REMOVED***
		return raw
	***REMOVED***
	return reference.FamiliarString(reference.TagNameOnly(ref))
***REMOVED***

func (c *containerConfig) portBindings() nat.PortMap ***REMOVED***
	portBindings := nat.PortMap***REMOVED******REMOVED***
	if c.task.Endpoint == nil ***REMOVED***
		return portBindings
	***REMOVED***

	for _, portConfig := range c.task.Endpoint.Ports ***REMOVED***
		if portConfig.PublishMode != api.PublishModeHost ***REMOVED***
			continue
		***REMOVED***

		port := nat.Port(fmt.Sprintf("%d/%s", portConfig.TargetPort, strings.ToLower(portConfig.Protocol.String())))
		binding := []nat.PortBinding***REMOVED***
			***REMOVED******REMOVED***,
		***REMOVED***

		if portConfig.PublishedPort != 0 ***REMOVED***
			binding[0].HostPort = strconv.Itoa(int(portConfig.PublishedPort))
		***REMOVED***
		portBindings[port] = binding
	***REMOVED***

	return portBindings
***REMOVED***

func (c *containerConfig) isolation() enginecontainer.Isolation ***REMOVED***
	return convert.IsolationFromGRPC(c.spec().Isolation)
***REMOVED***

func (c *containerConfig) exposedPorts() map[nat.Port]struct***REMOVED******REMOVED*** ***REMOVED***
	exposedPorts := make(map[nat.Port]struct***REMOVED******REMOVED***)
	if c.task.Endpoint == nil ***REMOVED***
		return exposedPorts
	***REMOVED***

	for _, portConfig := range c.task.Endpoint.Ports ***REMOVED***
		if portConfig.PublishMode != api.PublishModeHost ***REMOVED***
			continue
		***REMOVED***

		port := nat.Port(fmt.Sprintf("%d/%s", portConfig.TargetPort, strings.ToLower(portConfig.Protocol.String())))
		exposedPorts[port] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	return exposedPorts
***REMOVED***

func (c *containerConfig) config() *enginecontainer.Config ***REMOVED***
	genericEnvs := genericresource.EnvFormat(c.task.AssignedGenericResources, "DOCKER_RESOURCE")
	env := append(c.spec().Env, genericEnvs...)

	config := &enginecontainer.Config***REMOVED***
		Labels:       c.labels(),
		StopSignal:   c.spec().StopSignal,
		Tty:          c.spec().TTY,
		OpenStdin:    c.spec().OpenStdin,
		User:         c.spec().User,
		Env:          env,
		Hostname:     c.spec().Hostname,
		WorkingDir:   c.spec().Dir,
		Image:        c.image(),
		ExposedPorts: c.exposedPorts(),
		Healthcheck:  c.healthcheck(),
	***REMOVED***

	if len(c.spec().Command) > 0 ***REMOVED***
		// If Command is provided, we replace the whole invocation with Command
		// by replacing Entrypoint and specifying Cmd. Args is ignored in this
		// case.
		config.Entrypoint = append(config.Entrypoint, c.spec().Command...)
		config.Cmd = append(config.Cmd, c.spec().Args...)
	***REMOVED*** else if len(c.spec().Args) > 0 ***REMOVED***
		// In this case, we assume the image has an Entrypoint and Args
		// specifies the arguments for that entrypoint.
		config.Cmd = c.spec().Args
	***REMOVED***

	return config
***REMOVED***

func (c *containerConfig) labels() map[string]string ***REMOVED***
	var (
		system = map[string]string***REMOVED***
			"task":         "", // mark as cluster task
			"task.id":      c.task.ID,
			"task.name":    c.name(),
			"node.id":      c.task.NodeID,
			"service.id":   c.task.ServiceID,
			"service.name": c.task.ServiceAnnotations.Name,
		***REMOVED***
		labels = make(map[string]string)
	)

	// base labels are those defined in the spec.
	for k, v := range c.spec().Labels ***REMOVED***
		labels[k] = v
	***REMOVED***

	// we then apply the overrides from the task, which may be set via the
	// orchestrator.
	for k, v := range c.task.Annotations.Labels ***REMOVED***
		labels[k] = v
	***REMOVED***

	// finally, we apply the system labels, which override all labels.
	for k, v := range system ***REMOVED***
		labels[strings.Join([]string***REMOVED***systemLabelPrefix, k***REMOVED***, ".")] = v
	***REMOVED***

	return labels
***REMOVED***

func (c *containerConfig) mounts() []enginemount.Mount ***REMOVED***
	var r []enginemount.Mount
	for _, mount := range c.spec().Mounts ***REMOVED***
		r = append(r, convertMount(mount))
	***REMOVED***
	return r
***REMOVED***

func convertMount(m api.Mount) enginemount.Mount ***REMOVED***
	mount := enginemount.Mount***REMOVED***
		Source:   m.Source,
		Target:   m.Target,
		ReadOnly: m.ReadOnly,
	***REMOVED***

	switch m.Type ***REMOVED***
	case api.MountTypeBind:
		mount.Type = enginemount.TypeBind
	case api.MountTypeVolume:
		mount.Type = enginemount.TypeVolume
	case api.MountTypeTmpfs:
		mount.Type = enginemount.TypeTmpfs
	***REMOVED***

	if m.BindOptions != nil ***REMOVED***
		mount.BindOptions = &enginemount.BindOptions***REMOVED******REMOVED***
		switch m.BindOptions.Propagation ***REMOVED***
		case api.MountPropagationRPrivate:
			mount.BindOptions.Propagation = enginemount.PropagationRPrivate
		case api.MountPropagationPrivate:
			mount.BindOptions.Propagation = enginemount.PropagationPrivate
		case api.MountPropagationRSlave:
			mount.BindOptions.Propagation = enginemount.PropagationRSlave
		case api.MountPropagationSlave:
			mount.BindOptions.Propagation = enginemount.PropagationSlave
		case api.MountPropagationRShared:
			mount.BindOptions.Propagation = enginemount.PropagationRShared
		case api.MountPropagationShared:
			mount.BindOptions.Propagation = enginemount.PropagationShared
		***REMOVED***
	***REMOVED***

	if m.VolumeOptions != nil ***REMOVED***
		mount.VolumeOptions = &enginemount.VolumeOptions***REMOVED***
			NoCopy: m.VolumeOptions.NoCopy,
		***REMOVED***
		if m.VolumeOptions.Labels != nil ***REMOVED***
			mount.VolumeOptions.Labels = make(map[string]string, len(m.VolumeOptions.Labels))
			for k, v := range m.VolumeOptions.Labels ***REMOVED***
				mount.VolumeOptions.Labels[k] = v
			***REMOVED***
		***REMOVED***
		if m.VolumeOptions.DriverConfig != nil ***REMOVED***
			mount.VolumeOptions.DriverConfig = &enginemount.Driver***REMOVED***
				Name: m.VolumeOptions.DriverConfig.Name,
			***REMOVED***
			if m.VolumeOptions.DriverConfig.Options != nil ***REMOVED***
				mount.VolumeOptions.DriverConfig.Options = make(map[string]string, len(m.VolumeOptions.DriverConfig.Options))
				for k, v := range m.VolumeOptions.DriverConfig.Options ***REMOVED***
					mount.VolumeOptions.DriverConfig.Options[k] = v
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if m.TmpfsOptions != nil ***REMOVED***
		mount.TmpfsOptions = &enginemount.TmpfsOptions***REMOVED***
			SizeBytes: m.TmpfsOptions.SizeBytes,
			Mode:      m.TmpfsOptions.Mode,
		***REMOVED***
	***REMOVED***

	return mount
***REMOVED***

func (c *containerConfig) healthcheck() *enginecontainer.HealthConfig ***REMOVED***
	hcSpec := c.spec().Healthcheck
	if hcSpec == nil ***REMOVED***
		return nil
	***REMOVED***
	interval, _ := gogotypes.DurationFromProto(hcSpec.Interval)
	timeout, _ := gogotypes.DurationFromProto(hcSpec.Timeout)
	startPeriod, _ := gogotypes.DurationFromProto(hcSpec.StartPeriod)
	return &enginecontainer.HealthConfig***REMOVED***
		Test:        hcSpec.Test,
		Interval:    interval,
		Timeout:     timeout,
		Retries:     int(hcSpec.Retries),
		StartPeriod: startPeriod,
	***REMOVED***
***REMOVED***

func (c *containerConfig) hostConfig() *enginecontainer.HostConfig ***REMOVED***
	hc := &enginecontainer.HostConfig***REMOVED***
		Resources:      c.resources(),
		GroupAdd:       c.spec().Groups,
		PortBindings:   c.portBindings(),
		Mounts:         c.mounts(),
		ReadonlyRootfs: c.spec().ReadOnly,
		Isolation:      c.isolation(),
	***REMOVED***

	if c.spec().DNSConfig != nil ***REMOVED***
		hc.DNS = c.spec().DNSConfig.Nameservers
		hc.DNSSearch = c.spec().DNSConfig.Search
		hc.DNSOptions = c.spec().DNSConfig.Options
	***REMOVED***

	c.applyPrivileges(hc)

	// The format of extra hosts on swarmkit is specified in:
	// http://man7.org/linux/man-pages/man5/hosts.5.html
	//    IP_address canonical_hostname [aliases...]
	// However, the format of ExtraHosts in HostConfig is
	//    <host>:<ip>
	// We need to do the conversion here
	// (Alias is ignored for now)
	for _, entry := range c.spec().Hosts ***REMOVED***
		parts := strings.Fields(entry)
		if len(parts) > 1 ***REMOVED***
			hc.ExtraHosts = append(hc.ExtraHosts, fmt.Sprintf("%s:%s", parts[1], parts[0]))
		***REMOVED***
	***REMOVED***

	if c.task.LogDriver != nil ***REMOVED***
		hc.LogConfig = enginecontainer.LogConfig***REMOVED***
			Type:   c.task.LogDriver.Name,
			Config: c.task.LogDriver.Options,
		***REMOVED***
	***REMOVED***

	if len(c.task.Networks) > 0 ***REMOVED***
		labels := c.task.Networks[0].Network.Spec.Annotations.Labels
		name := c.task.Networks[0].Network.Spec.Annotations.Name
		if v, ok := labels["com.docker.swarm.predefined"]; ok && v == "true" ***REMOVED***
			hc.NetworkMode = enginecontainer.NetworkMode(name)
		***REMOVED***
	***REMOVED***

	return hc
***REMOVED***

// This handles the case of volumes that are defined inside a service Mount
func (c *containerConfig) volumeCreateRequest(mount *api.Mount) *volumetypes.VolumesCreateBody ***REMOVED***
	var (
		driverName string
		driverOpts map[string]string
		labels     map[string]string
	)

	if mount.VolumeOptions != nil && mount.VolumeOptions.DriverConfig != nil ***REMOVED***
		driverName = mount.VolumeOptions.DriverConfig.Name
		driverOpts = mount.VolumeOptions.DriverConfig.Options
		labels = mount.VolumeOptions.Labels
	***REMOVED***

	if mount.VolumeOptions != nil ***REMOVED***
		return &volumetypes.VolumesCreateBody***REMOVED***
			Name:       mount.Source,
			Driver:     driverName,
			DriverOpts: driverOpts,
			Labels:     labels,
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *containerConfig) resources() enginecontainer.Resources ***REMOVED***
	resources := enginecontainer.Resources***REMOVED******REMOVED***

	// If no limits are specified let the engine use its defaults.
	//
	// TODO(aluzzardi): We might want to set some limits anyway otherwise
	// "unlimited" tasks will step over the reservation of other tasks.
	r := c.task.Spec.Resources
	if r == nil || r.Limits == nil ***REMOVED***
		return resources
	***REMOVED***

	if r.Limits.MemoryBytes > 0 ***REMOVED***
		resources.Memory = r.Limits.MemoryBytes
	***REMOVED***

	if r.Limits.NanoCPUs > 0 ***REMOVED***
		// CPU Period must be set in microseconds.
		resources.CPUPeriod = int64(cpuQuotaPeriod / time.Microsecond)
		resources.CPUQuota = r.Limits.NanoCPUs * resources.CPUPeriod / 1e9
	***REMOVED***

	return resources
***REMOVED***

// Docker daemon supports just 1 network during container create.
func (c *containerConfig) createNetworkingConfig(b executorpkg.Backend) *network.NetworkingConfig ***REMOVED***
	var networks []*api.NetworkAttachment
	if c.task.Spec.GetContainer() != nil || c.task.Spec.GetAttachment() != nil ***REMOVED***
		networks = c.task.Networks
	***REMOVED***

	epConfig := make(map[string]*network.EndpointSettings)
	if len(networks) > 0 ***REMOVED***
		epConfig[networks[0].Network.Spec.Annotations.Name] = getEndpointConfig(networks[0], b)
	***REMOVED***

	return &network.NetworkingConfig***REMOVED***EndpointsConfig: epConfig***REMOVED***
***REMOVED***

// TODO: Merge this function with createNetworkingConfig after daemon supports multiple networks in container create
func (c *containerConfig) connectNetworkingConfig(b executorpkg.Backend) *network.NetworkingConfig ***REMOVED***
	var networks []*api.NetworkAttachment
	if c.task.Spec.GetContainer() != nil ***REMOVED***
		networks = c.task.Networks
	***REMOVED***
	// First network is used during container create. Other networks are used in "docker network connect"
	if len(networks) < 2 ***REMOVED***
		return nil
	***REMOVED***

	epConfig := make(map[string]*network.EndpointSettings)
	for _, na := range networks[1:] ***REMOVED***
		epConfig[na.Network.Spec.Annotations.Name] = getEndpointConfig(na, b)
	***REMOVED***
	return &network.NetworkingConfig***REMOVED***EndpointsConfig: epConfig***REMOVED***
***REMOVED***

func getEndpointConfig(na *api.NetworkAttachment, b executorpkg.Backend) *network.EndpointSettings ***REMOVED***
	var ipv4, ipv6 string
	for _, addr := range na.Addresses ***REMOVED***
		ip, _, err := net.ParseCIDR(addr)
		if err != nil ***REMOVED***
			continue
		***REMOVED***

		if ip.To4() != nil ***REMOVED***
			ipv4 = ip.String()
			continue
		***REMOVED***

		if ip.To16() != nil ***REMOVED***
			ipv6 = ip.String()
		***REMOVED***
	***REMOVED***

	n := &network.EndpointSettings***REMOVED***
		NetworkID: na.Network.ID,
		IPAMConfig: &network.EndpointIPAMConfig***REMOVED***
			IPv4Address: ipv4,
			IPv6Address: ipv6,
		***REMOVED***,
		DriverOpts: na.DriverAttachmentOpts,
	***REMOVED***
	if v, ok := na.Network.Spec.Annotations.Labels["com.docker.swarm.predefined"]; ok && v == "true" ***REMOVED***
		if ln, err := b.FindNetwork(na.Network.Spec.Annotations.Name); err == nil ***REMOVED***
			n.NetworkID = ln.ID()
		***REMOVED***
	***REMOVED***
	return n
***REMOVED***

func (c *containerConfig) virtualIP(networkID string) string ***REMOVED***
	if c.task.Endpoint == nil ***REMOVED***
		return ""
	***REMOVED***

	for _, eVip := range c.task.Endpoint.VirtualIPs ***REMOVED***
		// We only support IPv4 VIPs for now.
		if eVip.NetworkID == networkID ***REMOVED***
			vip, _, err := net.ParseCIDR(eVip.Addr)
			if err != nil ***REMOVED***
				return ""
			***REMOVED***

			return vip.String()
		***REMOVED***
	***REMOVED***

	return ""
***REMOVED***

func (c *containerConfig) serviceConfig() *clustertypes.ServiceConfig ***REMOVED***
	if len(c.task.Networks) == 0 ***REMOVED***
		return nil
	***REMOVED***

	logrus.Debugf("Creating service config in agent for t = %+v", c.task)
	svcCfg := &clustertypes.ServiceConfig***REMOVED***
		Name:             c.task.ServiceAnnotations.Name,
		Aliases:          make(map[string][]string),
		ID:               c.task.ServiceID,
		VirtualAddresses: make(map[string]*clustertypes.VirtualAddress),
	***REMOVED***

	for _, na := range c.task.Networks ***REMOVED***
		svcCfg.VirtualAddresses[na.Network.ID] = &clustertypes.VirtualAddress***REMOVED***
			// We support only IPv4 virtual IP for now.
			IPv4: c.virtualIP(na.Network.ID),
		***REMOVED***
		if len(na.Aliases) > 0 ***REMOVED***
			svcCfg.Aliases[na.Network.ID] = na.Aliases
		***REMOVED***
	***REMOVED***

	if c.task.Endpoint != nil ***REMOVED***
		for _, ePort := range c.task.Endpoint.Ports ***REMOVED***
			if ePort.PublishMode != api.PublishModeIngress ***REMOVED***
				continue
			***REMOVED***

			svcCfg.ExposedPorts = append(svcCfg.ExposedPorts, &clustertypes.PortConfig***REMOVED***
				Name:          ePort.Name,
				Protocol:      int32(ePort.Protocol),
				TargetPort:    ePort.TargetPort,
				PublishedPort: ePort.PublishedPort,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	return svcCfg
***REMOVED***

func (c *containerConfig) networkCreateRequest(name string) (clustertypes.NetworkCreateRequest, error) ***REMOVED***
	na, ok := c.networksAttachments[name]
	if !ok ***REMOVED***
		return clustertypes.NetworkCreateRequest***REMOVED******REMOVED***, errors.New("container: unknown network referenced")
	***REMOVED***

	options := types.NetworkCreate***REMOVED***
		// ID:     na.Network.ID,
		Labels:         na.Network.Spec.Annotations.Labels,
		Internal:       na.Network.Spec.Internal,
		Attachable:     na.Network.Spec.Attachable,
		Ingress:        convert.IsIngressNetwork(na.Network),
		EnableIPv6:     na.Network.Spec.Ipv6Enabled,
		CheckDuplicate: true,
		Scope:          netconst.SwarmScope,
	***REMOVED***

	if na.Network.Spec.GetNetwork() != "" ***REMOVED***
		options.ConfigFrom = &network.ConfigReference***REMOVED***
			Network: na.Network.Spec.GetNetwork(),
		***REMOVED***
	***REMOVED***

	if na.Network.DriverState != nil ***REMOVED***
		options.Driver = na.Network.DriverState.Name
		options.Options = na.Network.DriverState.Options
	***REMOVED***
	if na.Network.IPAM != nil ***REMOVED***
		options.IPAM = &network.IPAM***REMOVED***
			Driver:  na.Network.IPAM.Driver.Name,
			Options: na.Network.IPAM.Driver.Options,
		***REMOVED***
		for _, ic := range na.Network.IPAM.Configs ***REMOVED***
			c := network.IPAMConfig***REMOVED***
				Subnet:  ic.Subnet,
				IPRange: ic.Range,
				Gateway: ic.Gateway,
			***REMOVED***
			options.IPAM.Config = append(options.IPAM.Config, c)
		***REMOVED***
	***REMOVED***

	return clustertypes.NetworkCreateRequest***REMOVED***
		ID: na.Network.ID,
		NetworkCreateRequest: types.NetworkCreateRequest***REMOVED***
			Name:          name,
			NetworkCreate: options,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (c *containerConfig) applyPrivileges(hc *enginecontainer.HostConfig) ***REMOVED***
	privileges := c.spec().Privileges
	if privileges == nil ***REMOVED***
		return
	***REMOVED***

	credentials := privileges.CredentialSpec
	if credentials != nil ***REMOVED***
		switch credentials.Source.(type) ***REMOVED***
		case *api.Privileges_CredentialSpec_File:
			hc.SecurityOpt = append(hc.SecurityOpt, "credentialspec=file://"+credentials.GetFile())
		case *api.Privileges_CredentialSpec_Registry:
			hc.SecurityOpt = append(hc.SecurityOpt, "credentialspec=registry://"+credentials.GetRegistry())
		***REMOVED***
	***REMOVED***

	selinux := privileges.SELinuxContext
	if selinux != nil ***REMOVED***
		if selinux.Disable ***REMOVED***
			hc.SecurityOpt = append(hc.SecurityOpt, "label=disable")
		***REMOVED***
		if selinux.User != "" ***REMOVED***
			hc.SecurityOpt = append(hc.SecurityOpt, "label=user:"+selinux.User)
		***REMOVED***
		if selinux.Role != "" ***REMOVED***
			hc.SecurityOpt = append(hc.SecurityOpt, "label=role:"+selinux.Role)
		***REMOVED***
		if selinux.Level != "" ***REMOVED***
			hc.SecurityOpt = append(hc.SecurityOpt, "label=level:"+selinux.Level)
		***REMOVED***
		if selinux.Type != "" ***REMOVED***
			hc.SecurityOpt = append(hc.SecurityOpt, "label=type:"+selinux.Type)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c containerConfig) eventFilter() filters.Args ***REMOVED***
	filter := filters.NewArgs()
	filter.Add("type", events.ContainerEventType)
	filter.Add("name", c.name())
	filter.Add("label", fmt.Sprintf("%v.task.id=%v", systemLabelPrefix, c.task.ID))
	return filter
***REMOVED***
