package agent

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"golang.org/x/net/context"
)

const (
	initialSessionFailureBackoff = 100 * time.Millisecond
	maxSessionFailureBackoff     = 8 * time.Second
	nodeUpdatePeriod             = 20 * time.Second
)

// Agent implements the primary node functionality for a member of a swarm
// cluster. The primary functionality is to run and report on the status of
// tasks assigned to the node.
type Agent struct ***REMOVED***
	config *Config

	// The latest node object state from manager
	// for this node known to the agent.
	node *api.Node

	keys []*api.EncryptionKey

	sessionq chan sessionOperation
	worker   Worker

	started   chan struct***REMOVED******REMOVED***
	startOnce sync.Once // start only once
	ready     chan struct***REMOVED******REMOVED***
	leaving   chan struct***REMOVED******REMOVED***
	leaveOnce sync.Once
	left      chan struct***REMOVED******REMOVED*** // closed after "run" processes "leaving" and will no longer accept new assignments
	stopped   chan struct***REMOVED******REMOVED*** // requests shutdown
	stopOnce  sync.Once     // only allow stop to be called once
	closed    chan struct***REMOVED******REMOVED*** // only closed in run
	err       error         // read only after closed is closed

	nodeUpdatePeriod time.Duration
***REMOVED***

// New returns a new agent, ready for task dispatch.
func New(config *Config) (*Agent, error) ***REMOVED***
	if err := config.validate(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	a := &Agent***REMOVED***
		config:           config,
		sessionq:         make(chan sessionOperation),
		started:          make(chan struct***REMOVED******REMOVED***),
		leaving:          make(chan struct***REMOVED******REMOVED***),
		left:             make(chan struct***REMOVED******REMOVED***),
		stopped:          make(chan struct***REMOVED******REMOVED***),
		closed:           make(chan struct***REMOVED******REMOVED***),
		ready:            make(chan struct***REMOVED******REMOVED***),
		nodeUpdatePeriod: nodeUpdatePeriod,
	***REMOVED***

	a.worker = newWorker(config.DB, config.Executor, a)
	return a, nil
***REMOVED***

// Start begins execution of the agent in the provided context, if not already
// started.
//
// Start returns an error if the agent has already started.
func (a *Agent) Start(ctx context.Context) error ***REMOVED***
	err := errAgentStarted

	a.startOnce.Do(func() ***REMOVED***
		close(a.started)
		go a.run(ctx)
		err = nil // clear error above, only once.
	***REMOVED***)

	return err
***REMOVED***

// Leave instructs the agent to leave the cluster. This method will shutdown
// assignment processing and remove all assignments from the node.
// Leave blocks until worker has finished closing all task managers or agent
// is closed.
func (a *Agent) Leave(ctx context.Context) error ***REMOVED***
	select ***REMOVED***
	case <-a.started:
	default:
		return errAgentNotStarted
	***REMOVED***

	a.leaveOnce.Do(func() ***REMOVED***
		close(a.leaving)
	***REMOVED***)

	// Do not call Wait until we have confirmed that the agent is no longer
	// accepting assignments. Starting a worker might race with Wait.
	select ***REMOVED***
	case <-a.left:
	case <-a.closed:
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***

	// agent could be closed while Leave is in progress
	var err error
	ch := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		err = a.worker.Wait(ctx)
		close(ch)
	***REMOVED***()

	select ***REMOVED***
	case <-ch:
		return err
	case <-a.closed:
		return ErrClosed
	***REMOVED***
***REMOVED***

// Stop shuts down the agent, blocking until full shutdown. If the agent is not
// started, Stop will block until the agent has fully shutdown.
func (a *Agent) Stop(ctx context.Context) error ***REMOVED***
	select ***REMOVED***
	case <-a.started:
	default:
		return errAgentNotStarted
	***REMOVED***

	a.stop()

	// wait till closed or context cancelled
	select ***REMOVED***
	case <-a.closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***
***REMOVED***

// stop signals the agent shutdown process, returning true if this call was the
// first to actually shutdown the agent.
func (a *Agent) stop() bool ***REMOVED***
	var stopped bool
	a.stopOnce.Do(func() ***REMOVED***
		close(a.stopped)
		stopped = true
	***REMOVED***)

	return stopped
***REMOVED***

// Err returns the error that caused the agent to shutdown or nil. Err blocks
// until the agent is fully shutdown.
func (a *Agent) Err(ctx context.Context) error ***REMOVED***
	select ***REMOVED***
	case <-a.closed:
		return a.err
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***
***REMOVED***

// Ready returns a channel that will be closed when agent first becomes ready.
func (a *Agent) Ready() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return a.ready
***REMOVED***

func (a *Agent) run(ctx context.Context) ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer close(a.closed) // full shutdown.

	ctx = log.WithModule(ctx, "agent")

	log.G(ctx).Debug("(*Agent).run")
	defer log.G(ctx).Debug("(*Agent).run exited")

	nodeTLSInfo := a.config.NodeTLSInfo

	// get the node description
	nodeDescription, err := a.nodeDescriptionWithHostname(ctx, nodeTLSInfo)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).WithField("agent", a.config.Executor).Error("agent: node description unavailable")
	***REMOVED***
	// nodeUpdateTicker is used to periodically check for updates to node description
	nodeUpdateTicker := time.NewTicker(a.nodeUpdatePeriod)
	defer nodeUpdateTicker.Stop()

	var (
		backoff       time.Duration
		session       = newSession(ctx, a, backoff, "", nodeDescription) // start the initial session
		registered    = session.registered
		ready         = a.ready // first session ready
		sessionq      chan sessionOperation
		leaving       = a.leaving
		subscriptions = map[string]context.CancelFunc***REMOVED******REMOVED***
	)
	defer func() ***REMOVED***
		session.close()
	***REMOVED***()

	if err := a.worker.Init(ctx); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("worker initialization failed")
		a.err = err
		return // fatal?
	***REMOVED***
	defer a.worker.Close()

	// setup a reliable reporter to call back to us.
	reporter := newStatusReporter(ctx, a)
	defer reporter.Close()

	a.worker.Listen(ctx, reporter)

	updateNode := func() ***REMOVED***
		// skip updating if the registration isn't finished
		if registered != nil ***REMOVED***
			return
		***REMOVED***
		// get the current node description
		newNodeDescription, err := a.nodeDescriptionWithHostname(ctx, nodeTLSInfo)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).WithField("agent", a.config.Executor).Error("agent: updated node description unavailable")
		***REMOVED***

		// if newNodeDescription is nil, it will cause a panic when
		// trying to create a session. Typically this can happen
		// if the engine goes down
		if newNodeDescription == nil ***REMOVED***
			return
		***REMOVED***

		// if the node description has changed, update it to the new one
		// and close the session. The old session will be stopped and a
		// new one will be created with the updated description
		if !reflect.DeepEqual(nodeDescription, newNodeDescription) ***REMOVED***
			nodeDescription = newNodeDescription
			// close the session
			log.G(ctx).Info("agent: found node update")

			if err := session.close(); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("agent: closing session failed")
			***REMOVED***
			sessionq = nil
			registered = nil
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case operation := <-sessionq:
			operation.response <- operation.fn(session)
		case <-leaving:
			leaving = nil

			// TODO(stevvooe): Signal to the manager that the node is leaving.

			// when leaving we remove all assignments.
			if err := a.worker.Assign(ctx, nil); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("failed removing all assignments")
			***REMOVED***

			close(a.left)
		case msg := <-session.assignments:
			// if we have left, accept no more assignments
			if leaving == nil ***REMOVED***
				continue
			***REMOVED***

			switch msg.Type ***REMOVED***
			case api.AssignmentsMessage_COMPLETE:
				// Need to assign secrets and configs before tasks,
				// because tasks might depend on new secrets or configs
				if err := a.worker.Assign(ctx, msg.Changes); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to synchronize worker assignments")
				***REMOVED***
			case api.AssignmentsMessage_INCREMENTAL:
				if err := a.worker.Update(ctx, msg.Changes); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed to update worker assignments")
				***REMOVED***
			***REMOVED***
		case msg := <-session.messages:
			if err := a.handleSessionMessage(ctx, msg, nodeTLSInfo); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("session message handler failed")
			***REMOVED***
		case sub := <-session.subscriptions:
			if sub.Close ***REMOVED***
				if cancel, ok := subscriptions[sub.ID]; ok ***REMOVED***
					cancel()
				***REMOVED***
				delete(subscriptions, sub.ID)
				continue
			***REMOVED***

			if _, ok := subscriptions[sub.ID]; ok ***REMOVED***
				// Duplicate subscription
				continue
			***REMOVED***

			subCtx, subCancel := context.WithCancel(ctx)
			subscriptions[sub.ID] = subCancel
			// TODO(dperny) we're tossing the error here, that seems wrong
			go a.worker.Subscribe(subCtx, sub)
		case <-registered:
			log.G(ctx).Debugln("agent: registered")
			if ready != nil ***REMOVED***
				close(ready)
			***REMOVED***
			if a.config.SessionTracker != nil ***REMOVED***
				a.config.SessionTracker.SessionEstablished()
			***REMOVED***
			ready = nil
			registered = nil // we only care about this once per session
			backoff = 0      // reset backoff
			sessionq = a.sessionq
		case err := <-session.errs:
			// TODO(stevvooe): This may actually block if a session is closed
			// but no error was sent. This must be the only place
			// session.close is called in response to errors, for this to work.
			if err != nil ***REMOVED***
				if a.config.SessionTracker != nil ***REMOVED***
					a.config.SessionTracker.SessionError(err)
				***REMOVED***

				log.G(ctx).WithError(err).Error("agent: session failed")
				backoff = initialSessionFailureBackoff + 2*backoff
				if backoff > maxSessionFailureBackoff ***REMOVED***
					backoff = maxSessionFailureBackoff
				***REMOVED***
			***REMOVED***

			if err := session.close(); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("agent: closing session failed")
			***REMOVED***
			sessionq = nil
			// if we're here before <-registered, do nothing for that event
			registered = nil
		case <-session.closed:
			if a.config.SessionTracker != nil ***REMOVED***
				if err := a.config.SessionTracker.SessionClosed(); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("agent: exiting")
					a.err = err
					return
				***REMOVED***
			***REMOVED***

			log.G(ctx).Debugf("agent: rebuild session")

			// select a session registration delay from backoff range.
			delay := time.Duration(0)
			if backoff > 0 ***REMOVED***
				delay = time.Duration(rand.Int63n(int64(backoff)))
			***REMOVED***
			session = newSession(ctx, a, delay, session.sessionID, nodeDescription)
			registered = session.registered
		case ev := <-a.config.NotifyTLSChange:
			// the TLS info has changed, so force a check to see if we need to restart the session
			if tlsInfo, ok := ev.(*api.NodeTLSInfo); ok ***REMOVED***
				nodeTLSInfo = tlsInfo
				updateNode()
				nodeUpdateTicker.Stop()
				nodeUpdateTicker = time.NewTicker(a.nodeUpdatePeriod)
			***REMOVED***
		case <-nodeUpdateTicker.C:
			// periodically check to see whether the node information has changed, and if so, restart the session
			updateNode()
		case <-a.stopped:
			// TODO(stevvooe): Wait on shutdown and cleanup. May need to pump
			// this loop a few times.
			return
		case <-ctx.Done():
			if a.err == nil ***REMOVED***
				a.err = ctx.Err()
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (a *Agent) handleSessionMessage(ctx context.Context, message *api.SessionMessage, nti *api.NodeTLSInfo) error ***REMOVED***
	seen := map[api.Peer]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, manager := range message.Managers ***REMOVED***
		if manager.Peer.Addr == "" ***REMOVED***
			continue
		***REMOVED***

		a.config.ConnBroker.Remotes().Observe(*manager.Peer, int(manager.Weight))
		seen[*manager.Peer] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	var changes *NodeChanges
	if message.Node != nil && (a.node == nil || !nodesEqual(a.node, message.Node)) ***REMOVED***
		if a.config.NotifyNodeChange != nil ***REMOVED***
			changes = &NodeChanges***REMOVED***Node: message.Node.Copy()***REMOVED***
		***REMOVED***
		a.node = message.Node.Copy()
		if err := a.config.Executor.Configure(ctx, a.node); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("node configure failed")
		***REMOVED***
	***REMOVED***
	if len(message.RootCA) > 0 && !bytes.Equal(message.RootCA, nti.TrustRoot) ***REMOVED***
		if changes == nil ***REMOVED***
			changes = &NodeChanges***REMOVED***RootCert: message.RootCA***REMOVED***
		***REMOVED*** else ***REMOVED***
			changes.RootCert = message.RootCA
		***REMOVED***
	***REMOVED***

	if changes != nil ***REMOVED***
		a.config.NotifyNodeChange <- changes
	***REMOVED***

	// prune managers not in list.
	for peer := range a.config.ConnBroker.Remotes().Weights() ***REMOVED***
		if _, ok := seen[peer]; !ok ***REMOVED***
			a.config.ConnBroker.Remotes().Remove(peer)
		***REMOVED***
	***REMOVED***

	if message.NetworkBootstrapKeys == nil ***REMOVED***
		return nil
	***REMOVED***

	for _, key := range message.NetworkBootstrapKeys ***REMOVED***
		same := false
		for _, agentKey := range a.keys ***REMOVED***
			if agentKey.LamportTime == key.LamportTime ***REMOVED***
				same = true
			***REMOVED***
		***REMOVED***
		if !same ***REMOVED***
			a.keys = message.NetworkBootstrapKeys
			if err := a.config.Executor.SetNetworkBootstrapKeys(a.keys); err != nil ***REMOVED***
				panic(fmt.Errorf("configuring network key failed"))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type sessionOperation struct ***REMOVED***
	fn       func(session *session) error
	response chan error
***REMOVED***

// withSession runs fn with the current session.
func (a *Agent) withSession(ctx context.Context, fn func(session *session) error) error ***REMOVED***
	response := make(chan error, 1)
	select ***REMOVED***
	case a.sessionq <- sessionOperation***REMOVED***
		fn:       fn,
		response: response,
	***REMOVED***:
		select ***REMOVED***
		case err := <-response:
			return err
		case <-a.closed:
			return ErrClosed
		case <-ctx.Done():
			return ctx.Err()
		***REMOVED***
	case <-a.closed:
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***
***REMOVED***

// UpdateTaskStatus attempts to send a task status update over the current session,
// blocking until the operation is completed.
//
// If an error is returned, the operation should be retried.
func (a *Agent) UpdateTaskStatus(ctx context.Context, taskID string, status *api.TaskStatus) error ***REMOVED***
	log.G(ctx).WithField("task.id", taskID).Debug("(*Agent).UpdateTaskStatus")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, 1)
	if err := a.withSession(ctx, func(session *session) error ***REMOVED***
		go func() ***REMOVED***
			err := session.sendTaskStatus(ctx, taskID, status)
			if err != nil ***REMOVED***
				if err == errTaskUnknown ***REMOVED***
					err = nil // dispatcher no longer cares about this task.
				***REMOVED*** else ***REMOVED***
					log.G(ctx).WithError(err).Error("closing session after fatal error")
					session.sendError(err)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				log.G(ctx).Debug("task status reported")
			***REMOVED***

			errs <- err
		***REMOVED***()

		return nil
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	select ***REMOVED***
	case err := <-errs:
		return err
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***
***REMOVED***

// Publisher returns a LogPublisher for the given subscription
// as well as a cancel function that should be called when the log stream
// is completed.
func (a *Agent) Publisher(ctx context.Context, subscriptionID string) (exec.LogPublisher, func(), error) ***REMOVED***
	// TODO(stevvooe): The level of coordination here is WAY too much for logs.
	// These should only be best effort and really just buffer until a session is
	// ready. Ideally, they would use a separate connection completely.

	var (
		err       error
		publisher api.LogBroker_PublishLogsClient
	)

	err = a.withSession(ctx, func(session *session) error ***REMOVED***
		publisher, err = api.NewLogBrokerClient(session.conn.ClientConn).PublishLogs(ctx)
		return err
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// make little closure for ending the log stream
	sendCloseMsg := func() ***REMOVED***
		// send a close message, to tell the manager our logs are done
		publisher.Send(&api.PublishLogsMessage***REMOVED***
			SubscriptionID: subscriptionID,
			Close:          true,
		***REMOVED***)
		// close the stream forreal
		publisher.CloseSend()
	***REMOVED***

	return exec.LogPublisherFunc(func(ctx context.Context, message api.LogMessage) error ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				sendCloseMsg()
				return ctx.Err()
			default:
			***REMOVED***

			return publisher.Send(&api.PublishLogsMessage***REMOVED***
				SubscriptionID: subscriptionID,
				Messages:       []api.LogMessage***REMOVED***message***REMOVED***,
			***REMOVED***)
		***REMOVED***), func() ***REMOVED***
			sendCloseMsg()
		***REMOVED***, nil
***REMOVED***

// nodeDescriptionWithHostname retrieves node description, and overrides hostname if available
func (a *Agent) nodeDescriptionWithHostname(ctx context.Context, tlsInfo *api.NodeTLSInfo) (*api.NodeDescription, error) ***REMOVED***
	desc, err := a.config.Executor.Describe(ctx)

	// Override hostname and TLS info
	if desc != nil ***REMOVED***
		if a.config.Hostname != "" && desc != nil ***REMOVED***
			desc.Hostname = a.config.Hostname
		***REMOVED***
		desc.TLSInfo = tlsInfo
	***REMOVED***
	return desc, err
***REMOVED***

// nodesEqual returns true if the node states are functionally equal, ignoring status,
// version and other superfluous fields.
//
// This used to decide whether or not to propagate a node update to executor.
func nodesEqual(a, b *api.Node) bool ***REMOVED***
	a, b = a.Copy(), b.Copy()

	a.Status, b.Status = api.NodeStatus***REMOVED******REMOVED***, api.NodeStatus***REMOVED******REMOVED***
	a.Meta, b.Meta = api.Meta***REMOVED******REMOVED***, api.Meta***REMOVED******REMOVED***

	return reflect.DeepEqual(a, b)
***REMOVED***
