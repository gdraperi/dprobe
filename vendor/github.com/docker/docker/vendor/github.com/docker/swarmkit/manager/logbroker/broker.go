package logbroker

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/watch"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errAlreadyRunning = errors.New("broker is already running")
	errNotRunning     = errors.New("broker is not running")
)

type logMessage struct ***REMOVED***
	*api.PublishLogsMessage
	completed bool
	err       error
***REMOVED***

// LogBroker coordinates log subscriptions to services and tasks. Clients can
// publish and subscribe to logs channels.
//
// Log subscriptions are pushed to the work nodes by creating log subscsription
// tasks. As such, the LogBroker also acts as an orchestrator of these tasks.
type LogBroker struct ***REMOVED***
	mu                sync.RWMutex
	logQueue          *watch.Queue
	subscriptionQueue *watch.Queue

	registeredSubscriptions map[string]*subscription
	subscriptionsByNode     map[string]map[*subscription]struct***REMOVED******REMOVED***

	pctx      context.Context
	cancelAll context.CancelFunc

	store *store.MemoryStore
***REMOVED***

// New initializes and returns a new LogBroker
func New(store *store.MemoryStore) *LogBroker ***REMOVED***
	return &LogBroker***REMOVED***
		store: store,
	***REMOVED***
***REMOVED***

// Start starts the log broker
func (lb *LogBroker) Start(ctx context.Context) error ***REMOVED***
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.cancelAll != nil ***REMOVED***
		return errAlreadyRunning
	***REMOVED***

	lb.pctx, lb.cancelAll = context.WithCancel(ctx)
	lb.logQueue = watch.NewQueue()
	lb.subscriptionQueue = watch.NewQueue()
	lb.registeredSubscriptions = make(map[string]*subscription)
	lb.subscriptionsByNode = make(map[string]map[*subscription]struct***REMOVED******REMOVED***)
	return nil
***REMOVED***

// Stop stops the log broker
func (lb *LogBroker) Stop() error ***REMOVED***
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.cancelAll == nil ***REMOVED***
		return errNotRunning
	***REMOVED***
	lb.cancelAll()
	lb.cancelAll = nil

	lb.logQueue.Close()
	lb.subscriptionQueue.Close()

	return nil
***REMOVED***

func validateSelector(selector *api.LogSelector) error ***REMOVED***
	if selector == nil ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "log selector must be provided")
	***REMOVED***

	if len(selector.ServiceIDs) == 0 && len(selector.TaskIDs) == 0 && len(selector.NodeIDs) == 0 ***REMOVED***
		return status.Errorf(codes.InvalidArgument, "log selector must not be empty")
	***REMOVED***

	return nil
***REMOVED***

func (lb *LogBroker) newSubscription(selector *api.LogSelector, options *api.LogSubscriptionOptions) *subscription ***REMOVED***
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	subscription := newSubscription(lb.store, &api.SubscriptionMessage***REMOVED***
		ID:       identity.NewID(),
		Selector: selector,
		Options:  options,
	***REMOVED***, lb.subscriptionQueue)

	return subscription
***REMOVED***

func (lb *LogBroker) getSubscription(id string) *subscription ***REMOVED***
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	subscription, ok := lb.registeredSubscriptions[id]
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return subscription
***REMOVED***

func (lb *LogBroker) registerSubscription(subscription *subscription) ***REMOVED***
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.registeredSubscriptions[subscription.message.ID] = subscription
	lb.subscriptionQueue.Publish(subscription)

	for _, node := range subscription.Nodes() ***REMOVED***
		if _, ok := lb.subscriptionsByNode[node]; !ok ***REMOVED***
			// Mark nodes that won't receive the message as done.
			subscription.Done(node, fmt.Errorf("node %s is not available", node))
		***REMOVED*** else ***REMOVED***
			// otherwise, add the subscription to the node's subscriptions list
			lb.subscriptionsByNode[node][subscription] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (lb *LogBroker) unregisterSubscription(subscription *subscription) ***REMOVED***
	lb.mu.Lock()
	defer lb.mu.Unlock()

	delete(lb.registeredSubscriptions, subscription.message.ID)

	// remove the subscription from all of the nodes
	for _, node := range subscription.Nodes() ***REMOVED***
		// but only if a node exists
		if _, ok := lb.subscriptionsByNode[node]; ok ***REMOVED***
			delete(lb.subscriptionsByNode[node], subscription)
		***REMOVED***
	***REMOVED***

	subscription.Close()
	lb.subscriptionQueue.Publish(subscription)
***REMOVED***

// watchSubscriptions grabs all current subscriptions and notifies of any
// subscription change for this node.
//
// Subscriptions may fire multiple times and the caller has to protect against
// dupes.
func (lb *LogBroker) watchSubscriptions(nodeID string) ([]*subscription, chan events.Event, func()) ***REMOVED***
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// Watch for subscription changes for this node.
	ch, cancel := lb.subscriptionQueue.CallbackWatch(events.MatcherFunc(func(event events.Event) bool ***REMOVED***
		s := event.(*subscription)
		return s.Contains(nodeID)
	***REMOVED***))

	// Grab current subscriptions.
	var subscriptions []*subscription
	for _, s := range lb.registeredSubscriptions ***REMOVED***
		if s.Contains(nodeID) ***REMOVED***
			subscriptions = append(subscriptions, s)
		***REMOVED***
	***REMOVED***

	return subscriptions, ch, cancel
***REMOVED***

func (lb *LogBroker) subscribe(id string) (chan events.Event, func()) ***REMOVED***
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	return lb.logQueue.CallbackWatch(events.MatcherFunc(func(event events.Event) bool ***REMOVED***
		publish := event.(*logMessage)
		return publish.SubscriptionID == id
	***REMOVED***))
***REMOVED***

func (lb *LogBroker) publish(log *api.PublishLogsMessage) ***REMOVED***
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	lb.logQueue.Publish(&logMessage***REMOVED***PublishLogsMessage: log***REMOVED***)
***REMOVED***

// markDone wraps (*Subscription).Done() so that the removal of the sub from
// the node's subscription list is possible
func (lb *LogBroker) markDone(sub *subscription, nodeID string, err error) ***REMOVED***
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// remove the subscription from the node's subscription list, if it exists
	if _, ok := lb.subscriptionsByNode[nodeID]; ok ***REMOVED***
		delete(lb.subscriptionsByNode[nodeID], sub)
	***REMOVED***

	// mark the sub as done
	sub.Done(nodeID, err)
***REMOVED***

// SubscribeLogs creates a log subscription and streams back logs
func (lb *LogBroker) SubscribeLogs(request *api.SubscribeLogsRequest, stream api.Logs_SubscribeLogsServer) error ***REMOVED***
	ctx := stream.Context()

	if err := validateSelector(request.Selector); err != nil ***REMOVED***
		return err
	***REMOVED***

	lb.mu.Lock()
	pctx := lb.pctx
	lb.mu.Unlock()
	if pctx == nil ***REMOVED***
		return errNotRunning
	***REMOVED***

	subscription := lb.newSubscription(request.Selector, request.Options)
	subscription.Run(pctx)
	defer subscription.Stop()

	log := log.G(ctx).WithFields(
		logrus.Fields***REMOVED***
			"method":          "(*LogBroker).SubscribeLogs",
			"subscription.id": subscription.message.ID,
		***REMOVED***,
	)
	log.Debug("subscribed")

	publishCh, publishCancel := lb.subscribe(subscription.message.ID)
	defer publishCancel()

	lb.registerSubscription(subscription)
	defer lb.unregisterSubscription(subscription)

	completed := subscription.Wait(ctx)
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-pctx.Done():
			return pctx.Err()
		case event := <-publishCh:
			publish := event.(*logMessage)
			if publish.completed ***REMOVED***
				return publish.err
			***REMOVED***
			if err := stream.Send(&api.SubscribeLogsMessage***REMOVED***
				Messages: publish.Messages,
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		case <-completed:
			completed = nil
			lb.logQueue.Publish(&logMessage***REMOVED***
				PublishLogsMessage: &api.PublishLogsMessage***REMOVED***
					SubscriptionID: subscription.message.ID,
				***REMOVED***,
				completed: true,
				err:       subscription.Err(),
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (lb *LogBroker) nodeConnected(nodeID string) ***REMOVED***
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if _, ok := lb.subscriptionsByNode[nodeID]; !ok ***REMOVED***
		lb.subscriptionsByNode[nodeID] = make(map[*subscription]struct***REMOVED******REMOVED***)
	***REMOVED***
***REMOVED***

func (lb *LogBroker) nodeDisconnected(nodeID string) ***REMOVED***
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for sub := range lb.subscriptionsByNode[nodeID] ***REMOVED***
		sub.Done(nodeID, fmt.Errorf("node %s disconnected unexpectedly", nodeID))
	***REMOVED***
	delete(lb.subscriptionsByNode, nodeID)
***REMOVED***

// ListenSubscriptions returns a stream of matching subscriptions for the current node
func (lb *LogBroker) ListenSubscriptions(request *api.ListenSubscriptionsRequest, stream api.LogBroker_ListenSubscriptionsServer) error ***REMOVED***
	remote, err := ca.RemoteNode(stream.Context())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	lb.mu.Lock()
	pctx := lb.pctx
	lb.mu.Unlock()
	if pctx == nil ***REMOVED***
		return errNotRunning
	***REMOVED***

	lb.nodeConnected(remote.NodeID)
	defer lb.nodeDisconnected(remote.NodeID)

	log := log.G(stream.Context()).WithFields(
		logrus.Fields***REMOVED***
			"method": "(*LogBroker).ListenSubscriptions",
			"node":   remote.NodeID,
		***REMOVED***,
	)
	subscriptions, subscriptionCh, subscriptionCancel := lb.watchSubscriptions(remote.NodeID)
	defer subscriptionCancel()

	log.Debug("node registered")

	activeSubscriptions := make(map[string]*subscription)

	// Start by sending down all active subscriptions.
	for _, subscription := range subscriptions ***REMOVED***
		select ***REMOVED***
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-pctx.Done():
			return nil
		default:
		***REMOVED***

		if err := stream.Send(subscription.message); err != nil ***REMOVED***
			log.Error(err)
			return err
		***REMOVED***
		activeSubscriptions[subscription.message.ID] = subscription
	***REMOVED***

	// Send down new subscriptions.
	for ***REMOVED***
		select ***REMOVED***
		case v := <-subscriptionCh:
			subscription := v.(*subscription)

			if subscription.Closed() ***REMOVED***
				delete(activeSubscriptions, subscription.message.ID)
			***REMOVED*** else ***REMOVED***
				// Avoid sending down the same subscription multiple times
				if _, ok := activeSubscriptions[subscription.message.ID]; ok ***REMOVED***
					continue
				***REMOVED***
				activeSubscriptions[subscription.message.ID] = subscription
			***REMOVED***
			if err := stream.Send(subscription.message); err != nil ***REMOVED***
				log.Error(err)
				return err
			***REMOVED***
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-pctx.Done():
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// PublishLogs publishes log messages for a given subscription
func (lb *LogBroker) PublishLogs(stream api.LogBroker_PublishLogsServer) (err error) ***REMOVED***
	remote, err := ca.RemoteNode(stream.Context())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var currentSubscription *subscription
	defer func() ***REMOVED***
		if currentSubscription != nil ***REMOVED***
			lb.markDone(currentSubscription, remote.NodeID, err)
		***REMOVED***
	***REMOVED***()

	for ***REMOVED***
		logMsg, err := stream.Recv()
		if err == io.EOF ***REMOVED***
			return stream.SendAndClose(&api.PublishLogsResponse***REMOVED******REMOVED***)
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if logMsg.SubscriptionID == "" ***REMOVED***
			return status.Errorf(codes.InvalidArgument, "missing subscription ID")
		***REMOVED***

		if currentSubscription == nil ***REMOVED***
			currentSubscription = lb.getSubscription(logMsg.SubscriptionID)
			if currentSubscription == nil ***REMOVED***
				return status.Errorf(codes.NotFound, "unknown subscription ID")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if logMsg.SubscriptionID != currentSubscription.message.ID ***REMOVED***
				return status.Errorf(codes.InvalidArgument, "different subscription IDs in the same session")
			***REMOVED***
		***REMOVED***

		// if we have a close message, close out the subscription
		if logMsg.Close ***REMOVED***
			// Mark done and then set to nil so if we error after this point,
			// we don't try to close again in the defer
			lb.markDone(currentSubscription, remote.NodeID, err)
			currentSubscription = nil
			return nil
		***REMOVED***

		// Make sure logs are emitted using the right Node ID to avoid impersonation.
		for _, msg := range logMsg.Messages ***REMOVED***
			if msg.Context.NodeID != remote.NodeID ***REMOVED***
				return status.Errorf(codes.PermissionDenied, "invalid NodeID: expected=%s;received=%s", remote.NodeID, msg.Context.NodeID)
			***REMOVED***
		***REMOVED***

		lb.publish(logMsg)
	***REMOVED***
***REMOVED***
