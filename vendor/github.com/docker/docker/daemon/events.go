package daemon

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/container"
	daemonevents "github.com/docker/docker/daemon/events"
	"github.com/docker/libnetwork"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/sirupsen/logrus"
)

var (
	clusterEventAction = map[swarmapi.WatchActionKind]string***REMOVED***
		swarmapi.WatchActionKindCreate: "create",
		swarmapi.WatchActionKindUpdate: "update",
		swarmapi.WatchActionKindRemove: "remove",
	***REMOVED***
)

// LogContainerEvent generates an event related to a container with only the default attributes.
func (daemon *Daemon) LogContainerEvent(container *container.Container, action string) ***REMOVED***
	daemon.LogContainerEventWithAttributes(container, action, map[string]string***REMOVED******REMOVED***)
***REMOVED***

// LogContainerEventWithAttributes generates an event related to a container with specific given attributes.
func (daemon *Daemon) LogContainerEventWithAttributes(container *container.Container, action string, attributes map[string]string) ***REMOVED***
	copyAttributes(attributes, container.Config.Labels)
	if container.Config.Image != "" ***REMOVED***
		attributes["image"] = container.Config.Image
	***REMOVED***
	attributes["name"] = strings.TrimLeft(container.Name, "/")

	actor := events.Actor***REMOVED***
		ID:         container.ID,
		Attributes: attributes,
	***REMOVED***
	daemon.EventsService.Log(action, events.ContainerEventType, actor)
***REMOVED***

// LogImageEvent generates an event related to an image with only the default attributes.
func (daemon *Daemon) LogImageEvent(imageID, refName, action string) ***REMOVED***
	daemon.LogImageEventWithAttributes(imageID, refName, action, map[string]string***REMOVED******REMOVED***)
***REMOVED***

// LogImageEventWithAttributes generates an event related to an image with specific given attributes.
func (daemon *Daemon) LogImageEventWithAttributes(imageID, refName, action string, attributes map[string]string) ***REMOVED***
	img, err := daemon.GetImage(imageID)
	if err == nil && img.Config != nil ***REMOVED***
		// image has not been removed yet.
		// it could be missing if the event is `delete`.
		copyAttributes(attributes, img.Config.Labels)
	***REMOVED***
	if refName != "" ***REMOVED***
		attributes["name"] = refName
	***REMOVED***
	actor := events.Actor***REMOVED***
		ID:         imageID,
		Attributes: attributes,
	***REMOVED***

	daemon.EventsService.Log(action, events.ImageEventType, actor)
***REMOVED***

// LogPluginEvent generates an event related to a plugin with only the default attributes.
func (daemon *Daemon) LogPluginEvent(pluginID, refName, action string) ***REMOVED***
	daemon.LogPluginEventWithAttributes(pluginID, refName, action, map[string]string***REMOVED******REMOVED***)
***REMOVED***

// LogPluginEventWithAttributes generates an event related to a plugin with specific given attributes.
func (daemon *Daemon) LogPluginEventWithAttributes(pluginID, refName, action string, attributes map[string]string) ***REMOVED***
	attributes["name"] = refName
	actor := events.Actor***REMOVED***
		ID:         pluginID,
		Attributes: attributes,
	***REMOVED***
	daemon.EventsService.Log(action, events.PluginEventType, actor)
***REMOVED***

// LogVolumeEvent generates an event related to a volume.
func (daemon *Daemon) LogVolumeEvent(volumeID, action string, attributes map[string]string) ***REMOVED***
	actor := events.Actor***REMOVED***
		ID:         volumeID,
		Attributes: attributes,
	***REMOVED***
	daemon.EventsService.Log(action, events.VolumeEventType, actor)
***REMOVED***

// LogNetworkEvent generates an event related to a network with only the default attributes.
func (daemon *Daemon) LogNetworkEvent(nw libnetwork.Network, action string) ***REMOVED***
	daemon.LogNetworkEventWithAttributes(nw, action, map[string]string***REMOVED******REMOVED***)
***REMOVED***

// LogNetworkEventWithAttributes generates an event related to a network with specific given attributes.
func (daemon *Daemon) LogNetworkEventWithAttributes(nw libnetwork.Network, action string, attributes map[string]string) ***REMOVED***
	attributes["name"] = nw.Name()
	attributes["type"] = nw.Type()
	actor := events.Actor***REMOVED***
		ID:         nw.ID(),
		Attributes: attributes,
	***REMOVED***
	daemon.EventsService.Log(action, events.NetworkEventType, actor)
***REMOVED***

// LogDaemonEventWithAttributes generates an event related to the daemon itself with specific given attributes.
func (daemon *Daemon) LogDaemonEventWithAttributes(action string, attributes map[string]string) ***REMOVED***
	if daemon.EventsService != nil ***REMOVED***
		if info, err := daemon.SystemInfo(); err == nil && info.Name != "" ***REMOVED***
			attributes["name"] = info.Name
		***REMOVED***
		actor := events.Actor***REMOVED***
			ID:         daemon.ID,
			Attributes: attributes,
		***REMOVED***
		daemon.EventsService.Log(action, events.DaemonEventType, actor)
	***REMOVED***
***REMOVED***

// SubscribeToEvents returns the currently record of events, a channel to stream new events from, and a function to cancel the stream of events.
func (daemon *Daemon) SubscribeToEvents(since, until time.Time, filter filters.Args) ([]events.Message, chan interface***REMOVED******REMOVED***) ***REMOVED***
	ef := daemonevents.NewFilter(filter)
	return daemon.EventsService.SubscribeTopic(since, until, ef)
***REMOVED***

// UnsubscribeFromEvents stops the event subscription for a client by closing the
// channel where the daemon sends events to.
func (daemon *Daemon) UnsubscribeFromEvents(listener chan interface***REMOVED******REMOVED***) ***REMOVED***
	daemon.EventsService.Evict(listener)
***REMOVED***

// copyAttributes guarantees that labels are not mutated by event triggers.
func copyAttributes(attributes, labels map[string]string) ***REMOVED***
	if labels == nil ***REMOVED***
		return
	***REMOVED***
	for k, v := range labels ***REMOVED***
		attributes[k] = v
	***REMOVED***
***REMOVED***

// ProcessClusterNotifications gets changes from store and add them to event list
func (daemon *Daemon) ProcessClusterNotifications(ctx context.Context, watchStream chan *swarmapi.WatchMessage) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return
		case message, ok := <-watchStream:
			if !ok ***REMOVED***
				logrus.Debug("cluster event channel has stopped")
				return
			***REMOVED***
			daemon.generateClusterEvent(message)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (daemon *Daemon) generateClusterEvent(msg *swarmapi.WatchMessage) ***REMOVED***
	for _, event := range msg.Events ***REMOVED***
		if event.Object == nil ***REMOVED***
			logrus.Errorf("event without object: %v", event)
			continue
		***REMOVED***
		switch v := event.Object.GetObject().(type) ***REMOVED***
		case *swarmapi.Object_Node:
			daemon.logNodeEvent(event.Action, v.Node, event.OldObject.GetNode())
		case *swarmapi.Object_Service:
			daemon.logServiceEvent(event.Action, v.Service, event.OldObject.GetService())
		case *swarmapi.Object_Network:
			daemon.logNetworkEvent(event.Action, v.Network, event.OldObject.GetNetwork())
		case *swarmapi.Object_Secret:
			daemon.logSecretEvent(event.Action, v.Secret, event.OldObject.GetSecret())
		case *swarmapi.Object_Config:
			daemon.logConfigEvent(event.Action, v.Config, event.OldObject.GetConfig())
		default:
			logrus.Warnf("unrecognized event: %v", event)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (daemon *Daemon) logNetworkEvent(action swarmapi.WatchActionKind, net *swarmapi.Network, oldNet *swarmapi.Network) ***REMOVED***
	attributes := map[string]string***REMOVED***
		"name": net.Spec.Annotations.Name,
	***REMOVED***
	eventTime := eventTimestamp(net.Meta, action)
	daemon.logClusterEvent(action, net.ID, "network", attributes, eventTime)
***REMOVED***

func (daemon *Daemon) logSecretEvent(action swarmapi.WatchActionKind, secret *swarmapi.Secret, oldSecret *swarmapi.Secret) ***REMOVED***
	attributes := map[string]string***REMOVED***
		"name": secret.Spec.Annotations.Name,
	***REMOVED***
	eventTime := eventTimestamp(secret.Meta, action)
	daemon.logClusterEvent(action, secret.ID, "secret", attributes, eventTime)
***REMOVED***

func (daemon *Daemon) logConfigEvent(action swarmapi.WatchActionKind, config *swarmapi.Config, oldConfig *swarmapi.Config) ***REMOVED***
	attributes := map[string]string***REMOVED***
		"name": config.Spec.Annotations.Name,
	***REMOVED***
	eventTime := eventTimestamp(config.Meta, action)
	daemon.logClusterEvent(action, config.ID, "config", attributes, eventTime)
***REMOVED***

func (daemon *Daemon) logNodeEvent(action swarmapi.WatchActionKind, node *swarmapi.Node, oldNode *swarmapi.Node) ***REMOVED***
	name := node.Spec.Annotations.Name
	if name == "" && node.Description != nil ***REMOVED***
		name = node.Description.Hostname
	***REMOVED***
	attributes := map[string]string***REMOVED***
		"name": name,
	***REMOVED***
	eventTime := eventTimestamp(node.Meta, action)
	// In an update event, display the changes in attributes
	if action == swarmapi.WatchActionKindUpdate && oldNode != nil ***REMOVED***
		if node.Spec.Availability != oldNode.Spec.Availability ***REMOVED***
			attributes["availability.old"] = strings.ToLower(oldNode.Spec.Availability.String())
			attributes["availability.new"] = strings.ToLower(node.Spec.Availability.String())
		***REMOVED***
		if node.Role != oldNode.Role ***REMOVED***
			attributes["role.old"] = strings.ToLower(oldNode.Role.String())
			attributes["role.new"] = strings.ToLower(node.Role.String())
		***REMOVED***
		if node.Status.State != oldNode.Status.State ***REMOVED***
			attributes["state.old"] = strings.ToLower(oldNode.Status.State.String())
			attributes["state.new"] = strings.ToLower(node.Status.State.String())
		***REMOVED***
		// This handles change within manager role
		if node.ManagerStatus != nil && oldNode.ManagerStatus != nil ***REMOVED***
			// leader change
			if node.ManagerStatus.Leader != oldNode.ManagerStatus.Leader ***REMOVED***
				if node.ManagerStatus.Leader ***REMOVED***
					attributes["leader.old"] = "false"
					attributes["leader.new"] = "true"
				***REMOVED*** else ***REMOVED***
					attributes["leader.old"] = "true"
					attributes["leader.new"] = "false"
				***REMOVED***
			***REMOVED***
			if node.ManagerStatus.Reachability != oldNode.ManagerStatus.Reachability ***REMOVED***
				attributes["reachability.old"] = strings.ToLower(oldNode.ManagerStatus.Reachability.String())
				attributes["reachability.new"] = strings.ToLower(node.ManagerStatus.Reachability.String())
			***REMOVED***
		***REMOVED***
	***REMOVED***

	daemon.logClusterEvent(action, node.ID, "node", attributes, eventTime)
***REMOVED***

func (daemon *Daemon) logServiceEvent(action swarmapi.WatchActionKind, service *swarmapi.Service, oldService *swarmapi.Service) ***REMOVED***
	attributes := map[string]string***REMOVED***
		"name": service.Spec.Annotations.Name,
	***REMOVED***
	eventTime := eventTimestamp(service.Meta, action)

	if action == swarmapi.WatchActionKindUpdate && oldService != nil ***REMOVED***
		// check image
		if x, ok := service.Spec.Task.GetRuntime().(*swarmapi.TaskSpec_Container); ok ***REMOVED***
			containerSpec := x.Container
			if y, ok := oldService.Spec.Task.GetRuntime().(*swarmapi.TaskSpec_Container); ok ***REMOVED***
				oldContainerSpec := y.Container
				if containerSpec.Image != oldContainerSpec.Image ***REMOVED***
					attributes["image.old"] = oldContainerSpec.Image
					attributes["image.new"] = containerSpec.Image
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// This should not happen.
				logrus.Errorf("service %s runtime changed from %T to %T", service.Spec.Annotations.Name, oldService.Spec.Task.GetRuntime(), service.Spec.Task.GetRuntime())
			***REMOVED***
		***REMOVED***
		// check replicated count change
		if x, ok := service.Spec.GetMode().(*swarmapi.ServiceSpec_Replicated); ok ***REMOVED***
			replicas := x.Replicated.Replicas
			if y, ok := oldService.Spec.GetMode().(*swarmapi.ServiceSpec_Replicated); ok ***REMOVED***
				oldReplicas := y.Replicated.Replicas
				if replicas != oldReplicas ***REMOVED***
					attributes["replicas.old"] = strconv.FormatUint(oldReplicas, 10)
					attributes["replicas.new"] = strconv.FormatUint(replicas, 10)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// This should not happen.
				logrus.Errorf("service %s mode changed from %T to %T", service.Spec.Annotations.Name, oldService.Spec.GetMode(), service.Spec.GetMode())
			***REMOVED***
		***REMOVED***
		if service.UpdateStatus != nil ***REMOVED***
			if oldService.UpdateStatus == nil ***REMOVED***
				attributes["updatestate.new"] = strings.ToLower(service.UpdateStatus.State.String())
			***REMOVED*** else if service.UpdateStatus.State != oldService.UpdateStatus.State ***REMOVED***
				attributes["updatestate.old"] = strings.ToLower(oldService.UpdateStatus.State.String())
				attributes["updatestate.new"] = strings.ToLower(service.UpdateStatus.State.String())
			***REMOVED***
		***REMOVED***
	***REMOVED***
	daemon.logClusterEvent(action, service.ID, "service", attributes, eventTime)
***REMOVED***

func (daemon *Daemon) logClusterEvent(action swarmapi.WatchActionKind, id, eventType string, attributes map[string]string, eventTime time.Time) ***REMOVED***
	actor := events.Actor***REMOVED***
		ID:         id,
		Attributes: attributes,
	***REMOVED***

	jm := events.Message***REMOVED***
		Action:   clusterEventAction[action],
		Type:     eventType,
		Actor:    actor,
		Scope:    "swarm",
		Time:     eventTime.UTC().Unix(),
		TimeNano: eventTime.UTC().UnixNano(),
	***REMOVED***
	daemon.EventsService.PublishMessage(jm)
***REMOVED***

func eventTimestamp(meta swarmapi.Meta, action swarmapi.WatchActionKind) time.Time ***REMOVED***
	var eventTime time.Time
	switch action ***REMOVED***
	case swarmapi.WatchActionKindCreate:
		eventTime, _ = gogotypes.TimestampFromProto(meta.CreatedAt)
	case swarmapi.WatchActionKindUpdate:
		eventTime, _ = gogotypes.TimestampFromProto(meta.UpdatedAt)
	case swarmapi.WatchActionKindRemove:
		// There is no timestamp from store message for remove operations.
		// Use current time.
		eventTime = time.Now()
	***REMOVED***
	return eventTime
***REMOVED***
