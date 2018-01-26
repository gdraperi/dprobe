package exec

import (
	"fmt"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/equality"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// Controller controls execution of a task.
type Controller interface ***REMOVED***
	// Update the task definition seen by the controller. Will return
	// ErrTaskUpdateFailed if the provided task definition changes fields that
	// cannot be changed.
	//
	// Will be ignored if the task has exited.
	Update(ctx context.Context, t *api.Task) error

	// Prepare the task for execution. This should ensure that all resources
	// are created such that a call to start should execute immediately.
	Prepare(ctx context.Context) error

	// Start the target and return when it has started successfully.
	Start(ctx context.Context) error

	// Wait blocks until the target has exited.
	Wait(ctx context.Context) error

	// Shutdown requests to exit the target gracefully.
	Shutdown(ctx context.Context) error

	// Terminate the target.
	Terminate(ctx context.Context) error

	// Remove all resources allocated by the controller.
	Remove(ctx context.Context) error

	// Close closes any ephemeral resources associated with controller instance.
	Close() error
***REMOVED***

// ControllerLogs defines a component that makes logs accessible.
//
// Can usually be accessed on a controller instance via type assertion.
type ControllerLogs interface ***REMOVED***
	// Logs will write publisher until the context is cancelled or an error
	// occurs.
	Logs(ctx context.Context, publisher LogPublisher, options api.LogSubscriptionOptions) error
***REMOVED***

// LogPublisher defines the protocol for receiving a log message.
type LogPublisher interface ***REMOVED***
	Publish(ctx context.Context, message api.LogMessage) error
***REMOVED***

// LogPublisherFunc implements publisher with just a function.
type LogPublisherFunc func(ctx context.Context, message api.LogMessage) error

// Publish calls the wrapped function.
func (fn LogPublisherFunc) Publish(ctx context.Context, message api.LogMessage) error ***REMOVED***
	return fn(ctx, message)
***REMOVED***

// LogPublisherProvider defines the protocol for receiving a log publisher
type LogPublisherProvider interface ***REMOVED***
	Publisher(ctx context.Context, subscriptionID string) (LogPublisher, func(), error)
***REMOVED***

// ContainerStatuser reports status of a container.
//
// This can be implemented by controllers or error types.
type ContainerStatuser interface ***REMOVED***
	// ContainerStatus returns the status of the target container, if
	// available. When the container is not available, the status will be nil.
	ContainerStatus(ctx context.Context) (*api.ContainerStatus, error)
***REMOVED***

// PortStatuser reports status of ports which are allocated by the executor
type PortStatuser interface ***REMOVED***
	// PortStatus returns the status on a list of PortConfigs
	// which are managed at the host level by the controller.
	PortStatus(ctx context.Context) (*api.PortStatus, error)
***REMOVED***

// Resolve attempts to get a controller from the executor and reports the
// correct status depending on the tasks current state according to the result.
//
// Unlike Do, if an error is returned, the status should still be reported. The
// error merely reports the failure at getting the controller.
func Resolve(ctx context.Context, task *api.Task, executor Executor) (Controller, *api.TaskStatus, error) ***REMOVED***
	status := task.Status.Copy()

	defer func() ***REMOVED***
		logStateChange(ctx, task.DesiredState, task.Status.State, status.State)
	***REMOVED***()

	ctlr, err := executor.Controller(task)

	// depending on the tasks state, a failed controller resolution has varying
	// impact. The following expresses that impact.
	if err != nil ***REMOVED***
		status.Message = "resolving controller failed"
		status.Err = err.Error()
		// before the task has been started, we consider it a rejection.
		// if task is running, consider the task has failed
		// otherwise keep the existing state
		if task.Status.State < api.TaskStateStarting ***REMOVED***
			status.State = api.TaskStateRejected
		***REMOVED*** else if task.Status.State <= api.TaskStateRunning ***REMOVED***
			status.State = api.TaskStateFailed
		***REMOVED***
	***REMOVED*** else if task.Status.State < api.TaskStateAccepted ***REMOVED***
		// we always want to proceed to accepted when we resolve the controller
		status.Message = "accepted"
		status.State = api.TaskStateAccepted
		status.Err = ""
	***REMOVED***

	return ctlr, status, err
***REMOVED***

// Do progresses the task state using the controller performing a single
// operation on the controller. The return TaskStatus should be marked as the
// new state of the task.
//
// The returned status should be reported and placed back on to task
// before the next call. The operation can be cancelled by creating a
// cancelling context.
//
// Errors from the task controller will reported on the returned status. Any
// errors coming from this function should not be reported as related to the
// individual task.
//
// If ErrTaskNoop is returned, it means a second call to Do will result in no
// change. If ErrTaskDead is returned, calls to Do will no longer result in any
// action.
func Do(ctx context.Context, task *api.Task, ctlr Controller) (*api.TaskStatus, error) ***REMOVED***
	status := task.Status.Copy()

	// stay in the current state.
	noop := func(errs ...error) (*api.TaskStatus, error) ***REMOVED***
		return status, ErrTaskNoop
	***REMOVED***

	retry := func() (*api.TaskStatus, error) ***REMOVED***
		// while we retry on all errors, this allows us to explicitly declare
		// retry cases.
		return status, ErrTaskRetry
	***REMOVED***

	// transition moves the task to the next state.
	transition := func(state api.TaskState, msg string) (*api.TaskStatus, error) ***REMOVED***
		current := status.State
		status.State = state
		status.Message = msg
		status.Err = ""

		if current > state ***REMOVED***
			panic("invalid state transition")
		***REMOVED***
		return status, nil
	***REMOVED***

	// containerStatus exitCode keeps track of whether or not we've set it in
	// this particular method. Eventually, we assemble this as part of a defer.
	var (
		containerStatus *api.ContainerStatus
		portStatus      *api.PortStatus
		exitCode        int
	)

	// returned when a fatal execution of the task is fatal. In this case, we
	// proceed to a terminal error state and set the appropriate fields.
	//
	// Common checks for the nature of an error should be included here. If the
	// error is determined not to be fatal for the task,
	fatal := func(err error) (*api.TaskStatus, error) ***REMOVED***
		if err == nil ***REMOVED***
			panic("err must not be nil when fatal")
		***REMOVED***

		if cs, ok := err.(ContainerStatuser); ok ***REMOVED***
			var err error
			containerStatus, err = cs.ContainerStatus(ctx)
			if err != nil && !contextDoneError(err) ***REMOVED***
				log.G(ctx).WithError(err).Error("error resolving container status on fatal")
			***REMOVED***
		***REMOVED***

		// make sure we've set the *correct* exit code
		if ec, ok := err.(ExitCoder); ok ***REMOVED***
			exitCode = ec.ExitCode()
		***REMOVED***

		if cause := errors.Cause(err); cause == context.DeadlineExceeded || cause == context.Canceled ***REMOVED***
			return retry()
		***REMOVED***

		status.Err = err.Error() // still reported on temporary
		if IsTemporary(err) ***REMOVED***
			return retry()
		***REMOVED***

		// only at this point do we consider the error fatal to the task.
		log.G(ctx).WithError(err).Error("fatal task error")

		// NOTE(stevvooe): The following switch dictates the terminal failure
		// state based on the state in which the failure was encountered.
		switch ***REMOVED***
		case status.State < api.TaskStateStarting:
			status.State = api.TaskStateRejected
		case status.State >= api.TaskStateStarting:
			status.State = api.TaskStateFailed
		***REMOVED***

		return status, nil
	***REMOVED***

	// below, we have several callbacks that are run after the state transition
	// is completed.
	defer func() ***REMOVED***
		logStateChange(ctx, task.DesiredState, task.Status.State, status.State)

		if !equality.TaskStatusesEqualStable(status, &task.Status) ***REMOVED***
			status.Timestamp = ptypes.MustTimestampProto(time.Now())
		***REMOVED***
	***REMOVED***()

	// extract the container status from the container, if supported.
	defer func() ***REMOVED***
		// only do this if in an active state
		if status.State < api.TaskStateStarting ***REMOVED***
			return
		***REMOVED***

		if containerStatus == nil ***REMOVED***
			// collect this, if we haven't
			cctlr, ok := ctlr.(ContainerStatuser)
			if !ok ***REMOVED***
				return
			***REMOVED***

			var err error
			containerStatus, err = cctlr.ContainerStatus(ctx)
			if err != nil && !contextDoneError(err) ***REMOVED***
				log.G(ctx).WithError(err).Error("container status unavailable")
			***REMOVED***

			// at this point, things have gone fairly wrong. Remain positive
			// and let's get something out the door.
			if containerStatus == nil ***REMOVED***
				containerStatus = new(api.ContainerStatus)
				containerStatusTask := task.Status.GetContainer()
				if containerStatusTask != nil ***REMOVED***
					*containerStatus = *containerStatusTask // copy it over.
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// at this point, we *must* have a containerStatus.
		if exitCode != 0 ***REMOVED***
			containerStatus.ExitCode = int32(exitCode)
		***REMOVED***

		status.RuntimeStatus = &api.TaskStatus_Container***REMOVED***
			Container: containerStatus,
		***REMOVED***

		if portStatus == nil ***REMOVED***
			pctlr, ok := ctlr.(PortStatuser)
			if !ok ***REMOVED***
				return
			***REMOVED***

			var err error
			portStatus, err = pctlr.PortStatus(ctx)
			if err != nil && !contextDoneError(err) ***REMOVED***
				log.G(ctx).WithError(err).Error("container port status unavailable")
			***REMOVED***
		***REMOVED***

		status.PortStatus = portStatus
	***REMOVED***()

	// this branch bounds the largest state achievable in the agent as SHUTDOWN, which
	// is exactly the correct behavior for the agent.
	if task.DesiredState >= api.TaskStateShutdown ***REMOVED***
		if status.State >= api.TaskStateCompleted ***REMOVED***
			return noop()
		***REMOVED***

		if err := ctlr.Shutdown(ctx); err != nil ***REMOVED***
			return fatal(err)
		***REMOVED***

		return transition(api.TaskStateShutdown, "shutdown")
	***REMOVED***

	if status.State > task.DesiredState ***REMOVED***
		return noop() // way beyond desired state, pause
	***REMOVED***

	// the following states may proceed past desired state.
	switch status.State ***REMOVED***
	case api.TaskStatePreparing:
		if err := ctlr.Prepare(ctx); err != nil && err != ErrTaskPrepared ***REMOVED***
			return fatal(err)
		***REMOVED***

		return transition(api.TaskStateReady, "prepared")
	case api.TaskStateStarting:
		if err := ctlr.Start(ctx); err != nil && err != ErrTaskStarted ***REMOVED***
			return fatal(err)
		***REMOVED***

		return transition(api.TaskStateRunning, "started")
	case api.TaskStateRunning:
		if err := ctlr.Wait(ctx); err != nil ***REMOVED***
			return fatal(err)
		***REMOVED***

		return transition(api.TaskStateCompleted, "finished")
	***REMOVED***

	// The following represent "pause" states. We can only proceed when the
	// desired state is beyond our current state.
	if status.State >= task.DesiredState ***REMOVED***
		return noop()
	***REMOVED***

	switch status.State ***REMOVED***
	case api.TaskStateNew, api.TaskStatePending, api.TaskStateAssigned:
		return transition(api.TaskStateAccepted, "accepted")
	case api.TaskStateAccepted:
		return transition(api.TaskStatePreparing, "preparing")
	case api.TaskStateReady:
		return transition(api.TaskStateStarting, "starting")
	default: // terminal states
		return noop()
	***REMOVED***
***REMOVED***

func logStateChange(ctx context.Context, desired, previous, next api.TaskState) ***REMOVED***
	if previous != next ***REMOVED***
		fields := logrus.Fields***REMOVED***
			"state.transition": fmt.Sprintf("%v->%v", previous, next),
			"state.desired":    desired,
		***REMOVED***
		log.G(ctx).WithFields(fields).Debug("state changed")
	***REMOVED***
***REMOVED***

func contextDoneError(err error) bool ***REMOVED***
	cause := errors.Cause(err)
	return cause == context.Canceled || cause == context.DeadlineExceeded
***REMOVED***
