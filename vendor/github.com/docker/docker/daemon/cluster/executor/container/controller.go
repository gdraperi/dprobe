package container

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	executorpkg "github.com/docker/docker/daemon/cluster/executor"
	"github.com/docker/go-connections/nat"
	"github.com/docker/libnetwork"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
)

const defaultGossipConvergeDelay = 2 * time.Second

// controller implements agent.Controller against docker's API.
//
// Most operations against docker's API are done through the container name,
// which is unique to the task.
type controller struct ***REMOVED***
	task       *api.Task
	adapter    *containerAdapter
	closed     chan struct***REMOVED******REMOVED***
	err        error
	pulled     chan struct***REMOVED******REMOVED*** // closed after pull
	cancelPull func()        // cancels pull context if not nil
	pullErr    error         // pull error, only read after pulled closed
***REMOVED***

var _ exec.Controller = &controller***REMOVED******REMOVED***

// NewController returns a docker exec runner for the provided task.
func newController(b executorpkg.Backend, task *api.Task, node *api.NodeDescription, dependencies exec.DependencyGetter) (*controller, error) ***REMOVED***
	adapter, err := newContainerAdapter(b, task, node, dependencies)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &controller***REMOVED***
		task:    task,
		adapter: adapter,
		closed:  make(chan struct***REMOVED******REMOVED***),
	***REMOVED***, nil
***REMOVED***

func (r *controller) Task() (*api.Task, error) ***REMOVED***
	return r.task, nil
***REMOVED***

// ContainerStatus returns the container-specific status for the task.
func (r *controller) ContainerStatus(ctx context.Context) (*api.ContainerStatus, error) ***REMOVED***
	ctnr, err := r.adapter.inspect(ctx)
	if err != nil ***REMOVED***
		if isUnknownContainer(err) ***REMOVED***
			return nil, nil
		***REMOVED***
		return nil, err
	***REMOVED***
	return parseContainerStatus(ctnr)
***REMOVED***

func (r *controller) PortStatus(ctx context.Context) (*api.PortStatus, error) ***REMOVED***
	ctnr, err := r.adapter.inspect(ctx)
	if err != nil ***REMOVED***
		if isUnknownContainer(err) ***REMOVED***
			return nil, nil
		***REMOVED***

		return nil, err
	***REMOVED***

	return parsePortStatus(ctnr)
***REMOVED***

// Update tasks a recent task update and applies it to the container.
func (r *controller) Update(ctx context.Context, t *api.Task) error ***REMOVED***
	// TODO(stevvooe): While assignment of tasks is idempotent, we do allow
	// updates of metadata, such as labelling, as well as any other properties
	// that make sense.
	return nil
***REMOVED***

// Prepare creates a container and ensures the image is pulled.
//
// If the container has already be created, exec.ErrTaskPrepared is returned.
func (r *controller) Prepare(ctx context.Context) error ***REMOVED***
	if err := r.checkClosed(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Make sure all the networks that the task needs are created.
	if err := r.adapter.createNetworks(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Make sure all the volumes that the task needs are created.
	if err := r.adapter.createVolumes(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	if os.Getenv("DOCKER_SERVICE_PREFER_OFFLINE_IMAGE") != "1" ***REMOVED***
		if r.pulled == nil ***REMOVED***
			// Fork the pull to a different context to allow pull to continue
			// on re-entrant calls to Prepare. This ensures that Prepare can be
			// idempotent and not incur the extra cost of pulling when
			// cancelled on updates.
			var pctx context.Context

			r.pulled = make(chan struct***REMOVED******REMOVED***)
			pctx, r.cancelPull = context.WithCancel(context.Background()) // TODO(stevvooe): Bind a context to the entire controller.

			go func() ***REMOVED***
				defer close(r.pulled)
				r.pullErr = r.adapter.pullImage(pctx) // protected by closing r.pulled
			***REMOVED***()
		***REMOVED***

		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-r.pulled:
			if r.pullErr != nil ***REMOVED***
				// NOTE(stevvooe): We always try to pull the image to make sure we have
				// the most up to date version. This will return an error, but we only
				// log it. If the image truly doesn't exist, the create below will
				// error out.
				//
				// This gives us some nice behavior where we use up to date versions of
				// mutable tags, but will still run if the old image is available but a
				// registry is down.
				//
				// If you don't want this behavior, lock down your image to an
				// immutable tag or digest.
				log.G(ctx).WithError(r.pullErr).Error("pulling image failed")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if err := r.adapter.create(ctx); err != nil ***REMOVED***
		if isContainerCreateNameConflict(err) ***REMOVED***
			if _, err := r.adapter.inspect(ctx); err != nil ***REMOVED***
				return err
			***REMOVED***

			// container is already created. success!
			return exec.ErrTaskPrepared
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

// Start the container. An error will be returned if the container is already started.
func (r *controller) Start(ctx context.Context) error ***REMOVED***
	if err := r.checkClosed(); err != nil ***REMOVED***
		return err
	***REMOVED***

	ctnr, err := r.adapter.inspect(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Detect whether the container has *ever* been started. If so, we don't
	// issue the start.
	//
	// TODO(stevvooe): This is very racy. While reading inspect, another could
	// start the process and we could end up starting it twice.
	if ctnr.State.Status != "created" ***REMOVED***
		return exec.ErrTaskStarted
	***REMOVED***

	for ***REMOVED***
		if err := r.adapter.start(ctx); err != nil ***REMOVED***
			if _, ok := errors.Cause(err).(libnetwork.ErrNoSuchNetwork); ok ***REMOVED***
				// Retry network creation again if we
				// failed because some of the networks
				// were not found.
				if err := r.adapter.createNetworks(ctx); err != nil ***REMOVED***
					return err
				***REMOVED***

				continue
			***REMOVED***

			return errors.Wrap(err, "starting container failed")
		***REMOVED***

		break
	***REMOVED***

	// no health check
	if ctnr.Config == nil || ctnr.Config.Healthcheck == nil || len(ctnr.Config.Healthcheck.Test) == 0 || ctnr.Config.Healthcheck.Test[0] == "NONE" ***REMOVED***
		if err := r.adapter.activateServiceBinding(); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed to activate service binding for container %s which has no healthcheck config", r.adapter.container.name())
			return err
		***REMOVED***
		return nil
	***REMOVED***

	// wait for container to be healthy
	eventq := r.adapter.events(ctx)

	var healthErr error
	for ***REMOVED***
		select ***REMOVED***
		case event := <-eventq:
			if !r.matchevent(event) ***REMOVED***
				continue
			***REMOVED***

			switch event.Action ***REMOVED***
			case "die": // exit on terminal events
				ctnr, err := r.adapter.inspect(ctx)
				if err != nil ***REMOVED***
					return errors.Wrap(err, "die event received")
				***REMOVED*** else if ctnr.State.ExitCode != 0 ***REMOVED***
					return &exitError***REMOVED***code: ctnr.State.ExitCode, cause: healthErr***REMOVED***
				***REMOVED***

				return nil
			case "destroy":
				// If we get here, something has gone wrong but we want to exit
				// and report anyways.
				return ErrContainerDestroyed
			case "health_status: unhealthy":
				// in this case, we stop the container and report unhealthy status
				if err := r.Shutdown(ctx); err != nil ***REMOVED***
					return errors.Wrap(err, "unhealthy container shutdown failed")
				***REMOVED***
				// set health check error, and wait for container to fully exit ("die" event)
				healthErr = ErrContainerUnhealthy
			case "health_status: healthy":
				if err := r.adapter.activateServiceBinding(); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Errorf("failed to activate service binding for container %s after healthy event", r.adapter.container.name())
					return err
				***REMOVED***
				return nil
			***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-r.closed:
			return r.err
		***REMOVED***
	***REMOVED***
***REMOVED***

// Wait on the container to exit.
func (r *controller) Wait(pctx context.Context) error ***REMOVED***
	if err := r.checkClosed(); err != nil ***REMOVED***
		return err
	***REMOVED***

	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	healthErr := make(chan error, 1)
	go func() ***REMOVED***
		ectx, cancel := context.WithCancel(ctx) // cancel event context on first event
		defer cancel()
		if err := r.checkHealth(ectx); err == ErrContainerUnhealthy ***REMOVED***
			healthErr <- ErrContainerUnhealthy
			if err := r.Shutdown(ectx); err != nil ***REMOVED***
				log.G(ectx).WithError(err).Debug("shutdown failed on unhealthy")
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	waitC, err := r.adapter.wait(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if status := <-waitC; status.ExitCode() != 0 ***REMOVED***
		exitErr := &exitError***REMOVED***
			code: status.ExitCode(),
		***REMOVED***

		// Set the cause if it is knowable.
		select ***REMOVED***
		case e := <-healthErr:
			exitErr.cause = e
		default:
			if status.Err() != nil ***REMOVED***
				exitErr.cause = status.Err()
			***REMOVED***
		***REMOVED***

		return exitErr
	***REMOVED***

	return nil
***REMOVED***

func (r *controller) hasServiceBinding() bool ***REMOVED***
	if r.task == nil ***REMOVED***
		return false
	***REMOVED***

	// service is attached to a network besides the default bridge
	for _, na := range r.task.Networks ***REMOVED***
		if na.Network == nil ||
			na.Network.DriverState == nil ||
			na.Network.DriverState.Name == "bridge" && na.Network.Spec.Annotations.Name == "bridge" ***REMOVED***
			continue
		***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// Shutdown the container cleanly.
func (r *controller) Shutdown(ctx context.Context) error ***REMOVED***
	if err := r.checkClosed(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if r.cancelPull != nil ***REMOVED***
		r.cancelPull()
	***REMOVED***

	if r.hasServiceBinding() ***REMOVED***
		// remove container from service binding
		if err := r.adapter.deactivateServiceBinding(); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Warningf("failed to deactivate service binding for container %s", r.adapter.container.name())
			// Don't return an error here, because failure to deactivate
			// the service binding is expected if the container was never
			// started.
		***REMOVED***

		// add a delay for gossip converge
		// TODO(dongluochen): this delay should be configurable to fit different cluster size and network delay.
		time.Sleep(defaultGossipConvergeDelay)
	***REMOVED***

	if err := r.adapter.shutdown(ctx); err != nil ***REMOVED***
		if isUnknownContainer(err) || isStoppedContainer(err) ***REMOVED***
			return nil
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

// Terminate the container, with force.
func (r *controller) Terminate(ctx context.Context) error ***REMOVED***
	if err := r.checkClosed(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if r.cancelPull != nil ***REMOVED***
		r.cancelPull()
	***REMOVED***

	if err := r.adapter.terminate(ctx); err != nil ***REMOVED***
		if isUnknownContainer(err) ***REMOVED***
			return nil
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

// Remove the container and its resources.
func (r *controller) Remove(ctx context.Context) error ***REMOVED***
	if err := r.checkClosed(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if r.cancelPull != nil ***REMOVED***
		r.cancelPull()
	***REMOVED***

	// It may be necessary to shut down the task before removing it.
	if err := r.Shutdown(ctx); err != nil ***REMOVED***
		if isUnknownContainer(err) ***REMOVED***
			return nil
		***REMOVED***
		// This may fail if the task was already shut down.
		log.G(ctx).WithError(err).Debug("shutdown failed on removal")
	***REMOVED***

	// Try removing networks referenced in this task in case this
	// task is the last one referencing it
	if err := r.adapter.removeNetworks(ctx); err != nil ***REMOVED***
		if isUnknownContainer(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	if err := r.adapter.remove(ctx); err != nil ***REMOVED***
		if isUnknownContainer(err) ***REMOVED***
			return nil
		***REMOVED***

		return err
	***REMOVED***
	return nil
***REMOVED***

// waitReady waits for a container to be "ready".
// Ready means it's past the started state.
func (r *controller) waitReady(pctx context.Context) error ***REMOVED***
	if err := r.checkClosed(); err != nil ***REMOVED***
		return err
	***REMOVED***

	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	eventq := r.adapter.events(ctx)

	ctnr, err := r.adapter.inspect(ctx)
	if err != nil ***REMOVED***
		if !isUnknownContainer(err) ***REMOVED***
			return errors.Wrap(err, "inspect container failed")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		switch ctnr.State.Status ***REMOVED***
		case "running", "exited", "dead":
			return nil
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case event := <-eventq:
			if !r.matchevent(event) ***REMOVED***
				continue
			***REMOVED***

			switch event.Action ***REMOVED***
			case "start":
				return nil
			***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-r.closed:
			return r.err
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *controller) Logs(ctx context.Context, publisher exec.LogPublisher, options api.LogSubscriptionOptions) error ***REMOVED***
	if err := r.checkClosed(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// if we're following, wait for this container to be ready. there is a
	// problem here: if the container will never be ready (for example, it has
	// been totally deleted) then this will wait forever. however, this doesn't
	// actually cause any UI issues, and shouldn't be a problem. the stuck wait
	// will go away when the follow (context) is canceled.
	if options.Follow ***REMOVED***
		if err := r.waitReady(ctx); err != nil ***REMOVED***
			return errors.Wrap(err, "container not ready for logs")
		***REMOVED***
	***REMOVED***
	// if we're not following, we're not gonna wait for the container to be
	// ready. just call logs. if the container isn't ready, the call will fail
	// and return an error. no big deal, we don't care, we only want the logs
	// we can get RIGHT NOW with no follow

	logsContext, cancel := context.WithCancel(ctx)
	msgs, err := r.adapter.logs(logsContext, options)
	defer cancel()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed getting container logs")
	***REMOVED***

	var (
		// use a rate limiter to keep things under control but also provides some
		// ability coalesce messages.
		limiter = rate.NewLimiter(rate.Every(time.Second), 10<<20) // 10 MB/s
		msgctx  = api.LogContext***REMOVED***
			NodeID:    r.task.NodeID,
			ServiceID: r.task.ServiceID,
			TaskID:    r.task.ID,
		***REMOVED***
	)

	for ***REMOVED***
		msg, ok := <-msgs
		if !ok ***REMOVED***
			// we're done here, no more messages
			return nil
		***REMOVED***

		if msg.Err != nil ***REMOVED***
			// the defered cancel closes the adapter's log stream
			return msg.Err
		***REMOVED***

		// wait here for the limiter to catch up
		if err := limiter.WaitN(ctx, len(msg.Line)); err != nil ***REMOVED***
			return errors.Wrap(err, "failed rate limiter")
		***REMOVED***
		tsp, err := gogotypes.TimestampProto(msg.Timestamp)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "failed to convert timestamp")
		***REMOVED***
		var stream api.LogStream
		if msg.Source == "stdout" ***REMOVED***
			stream = api.LogStreamStdout
		***REMOVED*** else if msg.Source == "stderr" ***REMOVED***
			stream = api.LogStreamStderr
		***REMOVED***

		// parse the details out of the Attrs map
		var attrs []api.LogAttr
		if len(msg.Attrs) != 0 ***REMOVED***
			attrs = make([]api.LogAttr, 0, len(msg.Attrs))
			for _, attr := range msg.Attrs ***REMOVED***
				attrs = append(attrs, api.LogAttr***REMOVED***Key: attr.Key, Value: attr.Value***REMOVED***)
			***REMOVED***
		***REMOVED***

		if err := publisher.Publish(ctx, api.LogMessage***REMOVED***
			Context:   msgctx,
			Timestamp: tsp,
			Stream:    stream,
			Attrs:     attrs,
			Data:      msg.Line,
		***REMOVED***); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to publish log message")
		***REMOVED***
	***REMOVED***
***REMOVED***

// Close the runner and clean up any ephemeral resources.
func (r *controller) Close() error ***REMOVED***
	select ***REMOVED***
	case <-r.closed:
		return r.err
	default:
		if r.cancelPull != nil ***REMOVED***
			r.cancelPull()
		***REMOVED***

		r.err = exec.ErrControllerClosed
		close(r.closed)
	***REMOVED***
	return nil
***REMOVED***

func (r *controller) matchevent(event events.Message) bool ***REMOVED***
	if event.Type != events.ContainerEventType ***REMOVED***
		return false
	***REMOVED***
	// we can't filter using id since it will have huge chances to introduce a deadlock. see #33377.
	return event.Actor.Attributes["name"] == r.adapter.container.name()
***REMOVED***

func (r *controller) checkClosed() error ***REMOVED***
	select ***REMOVED***
	case <-r.closed:
		return r.err
	default:
		return nil
	***REMOVED***
***REMOVED***

func parseContainerStatus(ctnr types.ContainerJSON) (*api.ContainerStatus, error) ***REMOVED***
	status := &api.ContainerStatus***REMOVED***
		ContainerID: ctnr.ID,
		PID:         int32(ctnr.State.Pid),
		ExitCode:    int32(ctnr.State.ExitCode),
	***REMOVED***

	return status, nil
***REMOVED***

func parsePortStatus(ctnr types.ContainerJSON) (*api.PortStatus, error) ***REMOVED***
	status := &api.PortStatus***REMOVED******REMOVED***

	if ctnr.NetworkSettings != nil && len(ctnr.NetworkSettings.Ports) > 0 ***REMOVED***
		exposedPorts, err := parsePortMap(ctnr.NetworkSettings.Ports)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		status.Ports = exposedPorts
	***REMOVED***

	return status, nil
***REMOVED***

func parsePortMap(portMap nat.PortMap) ([]*api.PortConfig, error) ***REMOVED***
	exposedPorts := make([]*api.PortConfig, 0, len(portMap))

	for portProtocol, mapping := range portMap ***REMOVED***
		parts := strings.SplitN(string(portProtocol), "/", 2)
		if len(parts) != 2 ***REMOVED***
			return nil, fmt.Errorf("invalid port mapping: %s", portProtocol)
		***REMOVED***

		port, err := strconv.ParseUint(parts[0], 10, 16)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		protocol := api.ProtocolTCP
		switch strings.ToLower(parts[1]) ***REMOVED***
		case "tcp":
			protocol = api.ProtocolTCP
		case "udp":
			protocol = api.ProtocolUDP
		default:
			return nil, fmt.Errorf("invalid protocol: %s", parts[1])
		***REMOVED***

		for _, binding := range mapping ***REMOVED***
			hostPort, err := strconv.ParseUint(binding.HostPort, 10, 16)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			// TODO(aluzzardi): We're losing the port `name` here since
			// there's no way to retrieve it back from the Engine.
			exposedPorts = append(exposedPorts, &api.PortConfig***REMOVED***
				PublishMode:   api.PublishModeHost,
				Protocol:      protocol,
				TargetPort:    uint32(port),
				PublishedPort: uint32(hostPort),
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	return exposedPorts, nil
***REMOVED***

type exitError struct ***REMOVED***
	code  int
	cause error
***REMOVED***

func (e *exitError) Error() string ***REMOVED***
	if e.cause != nil ***REMOVED***
		return fmt.Sprintf("task: non-zero exit (%v): %v", e.code, e.cause)
	***REMOVED***

	return fmt.Sprintf("task: non-zero exit (%v)", e.code)
***REMOVED***

func (e *exitError) ExitCode() int ***REMOVED***
	return e.code
***REMOVED***

func (e *exitError) Cause() error ***REMOVED***
	return e.cause
***REMOVED***

// checkHealth blocks until unhealthy container is detected or ctx exits
func (r *controller) checkHealth(ctx context.Context) error ***REMOVED***
	eventq := r.adapter.events(ctx)

	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil
		case <-r.closed:
			return nil
		case event := <-eventq:
			if !r.matchevent(event) ***REMOVED***
				continue
			***REMOVED***

			switch event.Action ***REMOVED***
			case "health_status: unhealthy":
				return ErrContainerUnhealthy
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
