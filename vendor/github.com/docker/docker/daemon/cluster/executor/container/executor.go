package container

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	swarmtypes "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/controllers/plugin"
	"github.com/docker/docker/daemon/cluster/convert"
	executorpkg "github.com/docker/docker/daemon/cluster/executor"
	clustertypes "github.com/docker/docker/daemon/cluster/provider"
	networktypes "github.com/docker/libnetwork/types"
	"github.com/docker/swarmkit/agent"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/naming"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type executor struct ***REMOVED***
	backend       executorpkg.Backend
	pluginBackend plugin.Backend
	dependencies  exec.DependencyManager
	mutex         sync.Mutex // This mutex protects the following node field
	node          *api.NodeDescription
***REMOVED***

// NewExecutor returns an executor from the docker client.
func NewExecutor(b executorpkg.Backend, p plugin.Backend) exec.Executor ***REMOVED***
	return &executor***REMOVED***
		backend:       b,
		pluginBackend: p,
		dependencies:  agent.NewDependencyManager(),
	***REMOVED***
***REMOVED***

// Describe returns the underlying node description from the docker client.
func (e *executor) Describe(ctx context.Context) (*api.NodeDescription, error) ***REMOVED***
	info, err := e.backend.SystemInfo()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	plugins := map[api.PluginDescription]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	addPlugins := func(typ string, names []string) ***REMOVED***
		for _, name := range names ***REMOVED***
			plugins[api.PluginDescription***REMOVED***
				Type: typ,
				Name: name,
			***REMOVED***] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	// add v1 plugins
	addPlugins("Volume", info.Plugins.Volume)
	// Add builtin driver "overlay" (the only builtin multi-host driver) to
	// the plugin list by default.
	addPlugins("Network", append([]string***REMOVED***"overlay"***REMOVED***, info.Plugins.Network...))
	addPlugins("Authorization", info.Plugins.Authorization)
	addPlugins("Log", info.Plugins.Log)

	// add v2 plugins
	v2Plugins, err := e.backend.PluginManager().List(filters.NewArgs())
	if err == nil ***REMOVED***
		for _, plgn := range v2Plugins ***REMOVED***
			for _, typ := range plgn.Config.Interface.Types ***REMOVED***
				if typ.Prefix != "docker" || !plgn.Enabled ***REMOVED***
					continue
				***REMOVED***
				plgnTyp := typ.Capability
				switch typ.Capability ***REMOVED***
				case "volumedriver":
					plgnTyp = "Volume"
				case "networkdriver":
					plgnTyp = "Network"
				case "logdriver":
					plgnTyp = "Log"
				***REMOVED***

				plugins[api.PluginDescription***REMOVED***
					Type: plgnTyp,
					Name: plgn.Name,
				***REMOVED***] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	pluginFields := make([]api.PluginDescription, 0, len(plugins))
	for k := range plugins ***REMOVED***
		pluginFields = append(pluginFields, k)
	***REMOVED***

	sort.Sort(sortedPlugins(pluginFields))

	// parse []string labels into a map[string]string
	labels := map[string]string***REMOVED******REMOVED***
	for _, l := range info.Labels ***REMOVED***
		stringSlice := strings.SplitN(l, "=", 2)
		// this will take the last value in the list for a given key
		// ideally, one shouldn't assign multiple values to the same key
		if len(stringSlice) > 1 ***REMOVED***
			labels[stringSlice[0]] = stringSlice[1]
		***REMOVED***
	***REMOVED***

	description := &api.NodeDescription***REMOVED***
		Hostname: info.Name,
		Platform: &api.Platform***REMOVED***
			Architecture: info.Architecture,
			OS:           info.OSType,
		***REMOVED***,
		Engine: &api.EngineDescription***REMOVED***
			EngineVersion: info.ServerVersion,
			Labels:        labels,
			Plugins:       pluginFields,
		***REMOVED***,
		Resources: &api.Resources***REMOVED***
			NanoCPUs:    int64(info.NCPU) * 1e9,
			MemoryBytes: info.MemTotal,
			Generic:     convert.GenericResourcesToGRPC(info.GenericResources),
		***REMOVED***,
	***REMOVED***

	// Save the node information in the executor field
	e.mutex.Lock()
	e.node = description
	e.mutex.Unlock()

	return description, nil
***REMOVED***

func (e *executor) Configure(ctx context.Context, node *api.Node) error ***REMOVED***
	var ingressNA *api.NetworkAttachment
	attachments := make(map[string]string)

	for _, na := range node.Attachments ***REMOVED***
		if na.Network.Spec.Ingress ***REMOVED***
			ingressNA = na
		***REMOVED***
		attachments[na.Network.ID] = na.Addresses[0]
	***REMOVED***

	if (ingressNA == nil) && (node.Attachment != nil) ***REMOVED***
		ingressNA = node.Attachment
		attachments[ingressNA.Network.ID] = ingressNA.Addresses[0]
	***REMOVED***

	if ingressNA == nil ***REMOVED***
		e.backend.ReleaseIngress()
		return e.backend.GetAttachmentStore().ResetAttachments(attachments)
	***REMOVED***

	options := types.NetworkCreate***REMOVED***
		Driver: ingressNA.Network.DriverState.Name,
		IPAM: &network.IPAM***REMOVED***
			Driver: ingressNA.Network.IPAM.Driver.Name,
		***REMOVED***,
		Options:        ingressNA.Network.DriverState.Options,
		Ingress:        true,
		CheckDuplicate: true,
	***REMOVED***

	for _, ic := range ingressNA.Network.IPAM.Configs ***REMOVED***
		c := network.IPAMConfig***REMOVED***
			Subnet:  ic.Subnet,
			IPRange: ic.Range,
			Gateway: ic.Gateway,
		***REMOVED***
		options.IPAM.Config = append(options.IPAM.Config, c)
	***REMOVED***

	_, err := e.backend.SetupIngress(clustertypes.NetworkCreateRequest***REMOVED***
		ID: ingressNA.Network.ID,
		NetworkCreateRequest: types.NetworkCreateRequest***REMOVED***
			Name:          ingressNA.Network.Spec.Annotations.Name,
			NetworkCreate: options,
		***REMOVED***,
	***REMOVED***, ingressNA.Addresses[0])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return e.backend.GetAttachmentStore().ResetAttachments(attachments)
***REMOVED***

// Controller returns a docker container runner.
func (e *executor) Controller(t *api.Task) (exec.Controller, error) ***REMOVED***
	dependencyGetter := agent.Restrict(e.dependencies, t)

	// Get the node description from the executor field
	e.mutex.Lock()
	nodeDescription := e.node
	e.mutex.Unlock()

	if t.Spec.GetAttachment() != nil ***REMOVED***
		return newNetworkAttacherController(e.backend, t, nodeDescription, dependencyGetter)
	***REMOVED***

	var ctlr exec.Controller
	switch r := t.Spec.GetRuntime().(type) ***REMOVED***
	case *api.TaskSpec_Generic:
		logrus.WithFields(logrus.Fields***REMOVED***
			"kind":     r.Generic.Kind,
			"type_url": r.Generic.Payload.TypeUrl,
		***REMOVED***).Debug("custom runtime requested")
		runtimeKind, err := naming.Runtime(t.Spec)
		if err != nil ***REMOVED***
			return ctlr, err
		***REMOVED***
		switch runtimeKind ***REMOVED***
		case string(swarmtypes.RuntimePlugin):
			info, _ := e.backend.SystemInfo()
			if !info.ExperimentalBuild ***REMOVED***
				return ctlr, fmt.Errorf("runtime type %q only supported in experimental", swarmtypes.RuntimePlugin)
			***REMOVED***
			c, err := plugin.NewController(e.pluginBackend, t)
			if err != nil ***REMOVED***
				return ctlr, err
			***REMOVED***
			ctlr = c
		default:
			return ctlr, fmt.Errorf("unsupported runtime type: %q", runtimeKind)
		***REMOVED***
	case *api.TaskSpec_Container:
		c, err := newController(e.backend, t, nodeDescription, dependencyGetter)
		if err != nil ***REMOVED***
			return ctlr, err
		***REMOVED***
		ctlr = c
	default:
		return ctlr, fmt.Errorf("unsupported runtime: %q", r)
	***REMOVED***

	return ctlr, nil
***REMOVED***

func (e *executor) SetNetworkBootstrapKeys(keys []*api.EncryptionKey) error ***REMOVED***
	nwKeys := []*networktypes.EncryptionKey***REMOVED******REMOVED***
	for _, key := range keys ***REMOVED***
		nwKey := &networktypes.EncryptionKey***REMOVED***
			Subsystem:   key.Subsystem,
			Algorithm:   int32(key.Algorithm),
			Key:         make([]byte, len(key.Key)),
			LamportTime: key.LamportTime,
		***REMOVED***
		copy(nwKey.Key, key.Key)
		nwKeys = append(nwKeys, nwKey)
	***REMOVED***
	e.backend.SetNetworkBootstrapKeys(nwKeys)

	return nil
***REMOVED***

func (e *executor) Secrets() exec.SecretsManager ***REMOVED***
	return e.dependencies.Secrets()
***REMOVED***

func (e *executor) Configs() exec.ConfigsManager ***REMOVED***
	return e.dependencies.Configs()
***REMOVED***

type sortedPlugins []api.PluginDescription

func (sp sortedPlugins) Len() int ***REMOVED*** return len(sp) ***REMOVED***

func (sp sortedPlugins) Swap(i, j int) ***REMOVED*** sp[i], sp[j] = sp[j], sp[i] ***REMOVED***

func (sp sortedPlugins) Less(i, j int) bool ***REMOVED***
	if sp[i].Type != sp[j].Type ***REMOVED***
		return sp[i].Type < sp[j].Type
	***REMOVED***
	return sp[i].Name < sp[j].Name
***REMOVED***
