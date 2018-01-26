package container

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	containerpkg "github.com/docker/docker/container"
	"github.com/docker/docker/daemon/cluster/convert"
	executorpkg "github.com/docker/docker/daemon/cluster/executor"
	"github.com/docker/libnetwork"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
)

// containerAdapter conducts remote operations for a container. All calls
// are mostly naked calls to the client API, seeded with information from
// containerConfig.
type containerAdapter struct ***REMOVED***
	backend      executorpkg.Backend
	container    *containerConfig
	dependencies exec.DependencyGetter
***REMOVED***

func newContainerAdapter(b executorpkg.Backend, task *api.Task, node *api.NodeDescription, dependencies exec.DependencyGetter) (*containerAdapter, error) ***REMOVED***
	ctnr, err := newContainerConfig(task, node)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &containerAdapter***REMOVED***
		container:    ctnr,
		backend:      b,
		dependencies: dependencies,
	***REMOVED***, nil
***REMOVED***

func (c *containerAdapter) pullImage(ctx context.Context) error ***REMOVED***
	spec := c.container.spec()

	// Skip pulling if the image is referenced by image ID.
	if _, err := digest.Parse(spec.Image); err == nil ***REMOVED***
		return nil
	***REMOVED***

	// Skip pulling if the image is referenced by digest and already
	// exists locally.
	named, err := reference.ParseNormalizedNamed(spec.Image)
	if err == nil ***REMOVED***
		if _, ok := named.(reference.Canonical); ok ***REMOVED***
			_, err := c.backend.LookupImage(spec.Image)
			if err == nil ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// if the image needs to be pulled, the auth config will be retrieved and updated
	var encodedAuthConfig string
	if spec.PullOptions != nil ***REMOVED***
		encodedAuthConfig = spec.PullOptions.RegistryAuth
	***REMOVED***

	authConfig := &types.AuthConfig***REMOVED******REMOVED***
	if encodedAuthConfig != "" ***REMOVED***
		if err := json.NewDecoder(base64.NewDecoder(base64.URLEncoding, strings.NewReader(encodedAuthConfig))).Decode(authConfig); err != nil ***REMOVED***
			logrus.Warnf("invalid authconfig: %v", err)
		***REMOVED***
	***REMOVED***

	pr, pw := io.Pipe()
	metaHeaders := map[string][]string***REMOVED******REMOVED***
	go func() ***REMOVED***
		// TODO @jhowardmsft LCOW Support: This will need revisiting as
		// the stack is built up to include LCOW support for swarm.
		platform := runtime.GOOS
		err := c.backend.PullImage(ctx, c.container.image(), "", platform, metaHeaders, authConfig, pw)
		pw.CloseWithError(err)
	***REMOVED***()

	dec := json.NewDecoder(pr)
	dec.UseNumber()
	m := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	spamLimiter := rate.NewLimiter(rate.Every(time.Second), 1)

	lastStatus := ""
	for ***REMOVED***
		if err := dec.Decode(&m); err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			return err
		***REMOVED***
		l := log.G(ctx)
		// limit pull progress logs unless the status changes
		if spamLimiter.Allow() || lastStatus != m["status"] ***REMOVED***
			// if we have progress details, we have everything we need
			if progress, ok := m["progressDetail"].(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
				// first, log the image and status
				l = l.WithFields(logrus.Fields***REMOVED***
					"image":  c.container.image(),
					"status": m["status"],
				***REMOVED***)
				// then, if we have progress, log the progress
				if progress["current"] != nil && progress["total"] != nil ***REMOVED***
					l = l.WithFields(logrus.Fields***REMOVED***
						"current": progress["current"],
						"total":   progress["total"],
					***REMOVED***)
				***REMOVED***
			***REMOVED***
			l.Debug("pull in progress")
		***REMOVED***
		// sometimes, we get no useful information at all, and add no fields
		if status, ok := m["status"].(string); ok ***REMOVED***
			lastStatus = status
		***REMOVED***
	***REMOVED***

	// if the final stream object contained an error, return it
	if errMsg, ok := m["error"]; ok ***REMOVED***
		return fmt.Errorf("%v", errMsg)
	***REMOVED***
	return nil
***REMOVED***

func (c *containerAdapter) createNetworks(ctx context.Context) error ***REMOVED***
	for name := range c.container.networksAttachments ***REMOVED***
		ncr, err := c.container.networkCreateRequest(name)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := c.backend.CreateManagedNetwork(ncr); err != nil ***REMOVED*** // todo name missing
			if _, ok := err.(libnetwork.NetworkNameError); ok ***REMOVED***
				continue
			***REMOVED***

			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *containerAdapter) removeNetworks(ctx context.Context) error ***REMOVED***
	for name, v := range c.container.networksAttachments ***REMOVED***
		if err := c.backend.DeleteManagedNetwork(v.Network.ID); err != nil ***REMOVED***
			switch err.(type) ***REMOVED***
			case *libnetwork.ActiveEndpointsError:
				continue
			case libnetwork.ErrNoSuchNetwork:
				continue
			default:
				log.G(ctx).Errorf("network %s remove failed: %v", name, err)
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *containerAdapter) networkAttach(ctx context.Context) error ***REMOVED***
	config := c.container.createNetworkingConfig(c.backend)

	var (
		networkName string
		networkID   string
	)

	if config != nil ***REMOVED***
		for n, epConfig := range config.EndpointsConfig ***REMOVED***
			networkName = n
			networkID = epConfig.NetworkID
			break
		***REMOVED***
	***REMOVED***

	return c.backend.UpdateAttachment(networkName, networkID, c.container.networkAttachmentContainerID(), config)
***REMOVED***

func (c *containerAdapter) waitForDetach(ctx context.Context) error ***REMOVED***
	config := c.container.createNetworkingConfig(c.backend)

	var (
		networkName string
		networkID   string
	)

	if config != nil ***REMOVED***
		for n, epConfig := range config.EndpointsConfig ***REMOVED***
			networkName = n
			networkID = epConfig.NetworkID
			break
		***REMOVED***
	***REMOVED***

	return c.backend.WaitForDetachment(ctx, networkName, networkID, c.container.taskID(), c.container.networkAttachmentContainerID())
***REMOVED***

func (c *containerAdapter) create(ctx context.Context) error ***REMOVED***
	var cr containertypes.ContainerCreateCreatedBody
	var err error
	if cr, err = c.backend.CreateManagedContainer(types.ContainerCreateConfig***REMOVED***
		Name:       c.container.name(),
		Config:     c.container.config(),
		HostConfig: c.container.hostConfig(),
		// Use the first network in container create
		NetworkingConfig: c.container.createNetworkingConfig(c.backend),
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Docker daemon currently doesn't support multiple networks in container create
	// Connect to all other networks
	nc := c.container.connectNetworkingConfig(c.backend)

	if nc != nil ***REMOVED***
		for n, ep := range nc.EndpointsConfig ***REMOVED***
			if err := c.backend.ConnectContainerToNetwork(cr.ID, n, ep); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	container := c.container.task.Spec.GetContainer()
	if container == nil ***REMOVED***
		return errors.New("unable to get container from task spec")
	***REMOVED***

	if err := c.backend.SetContainerDependencyStore(cr.ID, c.dependencies); err != nil ***REMOVED***
		return err
	***REMOVED***

	// configure secrets
	secretRefs := convert.SecretReferencesFromGRPC(container.Secrets)
	if err := c.backend.SetContainerSecretReferences(cr.ID, secretRefs); err != nil ***REMOVED***
		return err
	***REMOVED***

	configRefs := convert.ConfigReferencesFromGRPC(container.Configs)
	if err := c.backend.SetContainerConfigReferences(cr.ID, configRefs); err != nil ***REMOVED***
		return err
	***REMOVED***

	return c.backend.UpdateContainerServiceConfig(cr.ID, c.container.serviceConfig())
***REMOVED***

// checkMounts ensures that the provided mounts won't have any host-specific
// problems at start up. For example, we disallow bind mounts without an
// existing path, which slightly different from the container API.
func (c *containerAdapter) checkMounts() error ***REMOVED***
	spec := c.container.spec()
	for _, mount := range spec.Mounts ***REMOVED***
		switch mount.Type ***REMOVED***
		case api.MountTypeBind:
			if _, err := os.Stat(mount.Source); os.IsNotExist(err) ***REMOVED***
				return fmt.Errorf("invalid bind mount source, source path not found: %s", mount.Source)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *containerAdapter) start(ctx context.Context) error ***REMOVED***
	if err := c.checkMounts(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return c.backend.ContainerStart(c.container.name(), nil, "", "")
***REMOVED***

func (c *containerAdapter) inspect(ctx context.Context) (types.ContainerJSON, error) ***REMOVED***
	cs, err := c.backend.ContainerInspectCurrent(c.container.name(), false)
	if ctx.Err() != nil ***REMOVED***
		return types.ContainerJSON***REMOVED******REMOVED***, ctx.Err()
	***REMOVED***
	if err != nil ***REMOVED***
		return types.ContainerJSON***REMOVED******REMOVED***, err
	***REMOVED***
	return *cs, nil
***REMOVED***

// events issues a call to the events API and returns a channel with all
// events. The stream of events can be shutdown by cancelling the context.
func (c *containerAdapter) events(ctx context.Context) <-chan events.Message ***REMOVED***
	log.G(ctx).Debugf("waiting on events")
	buffer, l := c.backend.SubscribeToEvents(time.Time***REMOVED******REMOVED***, time.Time***REMOVED******REMOVED***, c.container.eventFilter())
	eventsq := make(chan events.Message, len(buffer))

	for _, event := range buffer ***REMOVED***
		eventsq <- event
	***REMOVED***

	go func() ***REMOVED***
		defer c.backend.UnsubscribeFromEvents(l)

		for ***REMOVED***
			select ***REMOVED***
			case ev := <-l:
				jev, ok := ev.(events.Message)
				if !ok ***REMOVED***
					log.G(ctx).Warnf("unexpected event message: %q", ev)
					continue
				***REMOVED***
				select ***REMOVED***
				case eventsq <- jev:
				case <-ctx.Done():
					return
				***REMOVED***
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return eventsq
***REMOVED***

func (c *containerAdapter) wait(ctx context.Context) (<-chan containerpkg.StateStatus, error) ***REMOVED***
	return c.backend.ContainerWait(ctx, c.container.nameOrID(), containerpkg.WaitConditionNotRunning)
***REMOVED***

func (c *containerAdapter) shutdown(ctx context.Context) error ***REMOVED***
	// Default stop grace period to nil (daemon will use the stopTimeout of the container)
	var stopgrace *int
	spec := c.container.spec()
	if spec.StopGracePeriod != nil ***REMOVED***
		stopgraceValue := int(spec.StopGracePeriod.Seconds)
		stopgrace = &stopgraceValue
	***REMOVED***
	return c.backend.ContainerStop(c.container.name(), stopgrace)
***REMOVED***

func (c *containerAdapter) terminate(ctx context.Context) error ***REMOVED***
	return c.backend.ContainerKill(c.container.name(), uint64(syscall.SIGKILL))
***REMOVED***

func (c *containerAdapter) remove(ctx context.Context) error ***REMOVED***
	return c.backend.ContainerRm(c.container.name(), &types.ContainerRmConfig***REMOVED***
		RemoveVolume: true,
		ForceRemove:  true,
	***REMOVED***)
***REMOVED***

func (c *containerAdapter) createVolumes(ctx context.Context) error ***REMOVED***
	// Create plugin volumes that are embedded inside a Mount
	for _, mount := range c.container.task.Spec.GetContainer().Mounts ***REMOVED***
		if mount.Type != api.MountTypeVolume ***REMOVED***
			continue
		***REMOVED***

		if mount.VolumeOptions == nil ***REMOVED***
			continue
		***REMOVED***

		if mount.VolumeOptions.DriverConfig == nil ***REMOVED***
			continue
		***REMOVED***

		req := c.container.volumeCreateRequest(&mount)

		// Check if this volume exists on the engine
		if _, err := c.backend.VolumeCreate(req.Name, req.Driver, req.DriverOpts, req.Labels); err != nil ***REMOVED***
			// TODO(amitshukla): Today, volume create through the engine api does not return an error
			// when the named volume with the same parameters already exists.
			// It returns an error if the driver name is different - that is a valid error
			return err
		***REMOVED***

	***REMOVED***

	return nil
***REMOVED***

func (c *containerAdapter) activateServiceBinding() error ***REMOVED***
	return c.backend.ActivateContainerServiceBinding(c.container.name())
***REMOVED***

func (c *containerAdapter) deactivateServiceBinding() error ***REMOVED***
	return c.backend.DeactivateContainerServiceBinding(c.container.name())
***REMOVED***

func (c *containerAdapter) logs(ctx context.Context, options api.LogSubscriptionOptions) (<-chan *backend.LogMessage, error) ***REMOVED***
	apiOptions := &types.ContainerLogsOptions***REMOVED***
		Follow: options.Follow,

		// Always say yes to Timestamps and Details. we make the decision
		// of whether to return these to the user or not way higher up the
		// stack.
		Timestamps: true,
		Details:    true,
	***REMOVED***

	if options.Since != nil ***REMOVED***
		since, err := gogotypes.TimestampFromProto(options.Since)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// print since as this formatted string because the docker container
		// logs interface expects it like this.
		// see github.com/docker/docker/api/types/time.ParseTimestamps
		apiOptions.Since = fmt.Sprintf("%d.%09d", since.Unix(), int64(since.Nanosecond()))
	***REMOVED***

	if options.Tail < 0 ***REMOVED***
		// See protobuf documentation for details of how this works.
		apiOptions.Tail = fmt.Sprint(-options.Tail - 1)
	***REMOVED*** else if options.Tail > 0 ***REMOVED***
		return nil, errors.New("tail relative to start of logs not supported via docker API")
	***REMOVED***

	if len(options.Streams) == 0 ***REMOVED***
		// empty == all
		apiOptions.ShowStdout, apiOptions.ShowStderr = true, true
	***REMOVED*** else ***REMOVED***
		for _, stream := range options.Streams ***REMOVED***
			switch stream ***REMOVED***
			case api.LogStreamStdout:
				apiOptions.ShowStdout = true
			case api.LogStreamStderr:
				apiOptions.ShowStderr = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	msgs, _, err := c.backend.ContainerLogs(ctx, c.container.name(), apiOptions)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return msgs, nil
***REMOVED***

// todo: typed/wrapped errors
func isContainerCreateNameConflict(err error) bool ***REMOVED***
	return strings.Contains(err.Error(), "Conflict. The name")
***REMOVED***

func isUnknownContainer(err error) bool ***REMOVED***
	return strings.Contains(err.Error(), "No such container:")
***REMOVED***

func isStoppedContainer(err error) bool ***REMOVED***
	return strings.Contains(err.Error(), "is already stopped")
***REMOVED***
