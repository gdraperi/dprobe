package allocator

import (
	"fmt"
	"time"

	"github.com/docker/go-events"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/allocator/cnmallocator"
	"github.com/docker/swarmkit/manager/allocator/networkallocator"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/protobuf/ptypes"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	// Network allocator Voter ID for task allocation vote.
	networkVoter           = "network"
	allocatedStatusMessage = "pending task scheduling"
)

var (
	// ErrNoIngress is returned when no ingress network is found in store
	ErrNoIngress = errors.New("no ingress network found")
	errNoChanges = errors.New("task unchanged")

	retryInterval = 5 * time.Minute
)

// Network context information which is used throughout the network allocation code.
type networkContext struct ***REMOVED***
	ingressNetwork *api.Network
	// Instance of the low-level network allocator which performs
	// the actual network allocation.
	nwkAllocator networkallocator.NetworkAllocator

	// A set of tasks which are ready to be allocated as a batch. This is
	// distinct from "unallocatedTasks" which are tasks that failed to
	// allocate on the first try, being held for a future retry.
	pendingTasks map[string]*api.Task

	// A set of unallocated tasks which will be revisited if any thing
	// changes in system state that might help task allocation.
	unallocatedTasks map[string]*api.Task

	// A set of unallocated services which will be revisited if
	// any thing changes in system state that might help service
	// allocation.
	unallocatedServices map[string]*api.Service

	// A set of unallocated networks which will be revisited if
	// any thing changes in system state that might help network
	// allocation.
	unallocatedNetworks map[string]*api.Network

	// lastRetry is the last timestamp when unallocated
	// tasks/services/networks were retried.
	lastRetry time.Time

	// somethingWasDeallocated indicates that we just deallocated at
	// least one service/task/network, so we should retry failed
	// allocations (in we are experiencing IP exhaustion and an IP was
	// released).
	somethingWasDeallocated bool
***REMOVED***

func (a *Allocator) doNetworkInit(ctx context.Context) (err error) ***REMOVED***
	na, err := cnmallocator.New(a.pluginGetter)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	nc := &networkContext***REMOVED***
		nwkAllocator:        na,
		pendingTasks:        make(map[string]*api.Task),
		unallocatedTasks:    make(map[string]*api.Task),
		unallocatedServices: make(map[string]*api.Service),
		unallocatedNetworks: make(map[string]*api.Network),
		lastRetry:           time.Now(),
	***REMOVED***
	a.netCtx = nc
	defer func() ***REMOVED***
		// Clear a.netCtx if initialization was unsuccessful.
		if err != nil ***REMOVED***
			a.netCtx = nil
		***REMOVED***
	***REMOVED***()

	// Ingress network is now created at cluster's first time creation.
	// Check if we have the ingress network. If found, make sure it is
	// allocated, before reading all network objects for allocation.
	// If not found, it means it was removed by user, nothing to do here.
	ingressNetwork, err := GetIngressNetwork(a.store)
	switch err ***REMOVED***
	case nil:
		// Try to complete ingress network allocation before anything else so
		// that the we can get the preferred subnet for ingress network.
		nc.ingressNetwork = ingressNetwork
		if !na.IsAllocated(nc.ingressNetwork) ***REMOVED***
			if err := a.allocateNetwork(ctx, nc.ingressNetwork); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("failed allocating ingress network during init")
			***REMOVED*** else if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
				if err := a.commitAllocatedNetwork(ctx, batch, nc.ingressNetwork); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Error("failed committing allocation of ingress network during init")
				***REMOVED***
				return nil
			***REMOVED***); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("failed committing allocation of ingress network during init")
			***REMOVED***
		***REMOVED***
	case ErrNoIngress:
		// Ingress network is not present in store, It means user removed it
		// and did not create a new one.
	default:
		return errors.Wrap(err, "failure while looking for ingress network during init")
	***REMOVED***

	// Allocate networks in the store so far before we started
	// watching.
	var networks []*api.Network
	a.store.View(func(tx store.ReadTx) ***REMOVED***
		networks, err = store.FindNetworks(tx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error listing all networks in store while trying to allocate during init")
	***REMOVED***

	var allocatedNetworks []*api.Network
	for _, n := range networks ***REMOVED***
		if na.IsAllocated(n) ***REMOVED***
			continue
		***REMOVED***

		if err := a.allocateNetwork(ctx, n); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed allocating network %s during init", n.ID)
			continue
		***REMOVED***
		allocatedNetworks = append(allocatedNetworks, n)
	***REMOVED***

	if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, n := range allocatedNetworks ***REMOVED***
			if err := a.commitAllocatedNetwork(ctx, batch, n); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("failed committing allocation of network %s during init", n.ID)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed committing allocation of networks during init")
	***REMOVED***

	// First, allocate objects that already have addresses associated with
	// them, to reserve these IP addresses in internal state.
	if err := a.allocateNodes(ctx, true); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := a.allocateServices(ctx, true); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := a.allocateTasks(ctx, true); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := a.allocateNodes(ctx, false); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := a.allocateServices(ctx, false); err != nil ***REMOVED***
		return err
	***REMOVED***
	return a.allocateTasks(ctx, false)
***REMOVED***

func (a *Allocator) doNetworkAlloc(ctx context.Context, ev events.Event) ***REMOVED***
	nc := a.netCtx

	switch v := ev.(type) ***REMOVED***
	case api.EventCreateNetwork:
		n := v.Network.Copy()
		if nc.nwkAllocator.IsAllocated(n) ***REMOVED***
			break
		***REMOVED***

		if IsIngressNetwork(n) && nc.ingressNetwork != nil ***REMOVED***
			log.G(ctx).Errorf("Cannot allocate ingress network %s (%s) because another ingress network is already present: %s (%s)",
				n.ID, n.Spec.Annotations.Name, nc.ingressNetwork.ID, nc.ingressNetwork.Spec.Annotations)
			break
		***REMOVED***

		if err := a.allocateNetwork(ctx, n); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed allocation for network %s", n.ID)
			break
		***REMOVED***

		if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
			return a.commitAllocatedNetwork(ctx, batch, n)
		***REMOVED***); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed to commit allocation for network %s", n.ID)
		***REMOVED***
		if IsIngressNetwork(n) ***REMOVED***
			nc.ingressNetwork = n
		***REMOVED***
		err := a.allocateNodes(ctx, false)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error(err)
		***REMOVED***
	case api.EventDeleteNetwork:
		n := v.Network.Copy()

		if IsIngressNetwork(n) && nc.ingressNetwork != nil && nc.ingressNetwork.ID == n.ID ***REMOVED***
			nc.ingressNetwork = nil
		***REMOVED***

		if err := a.deallocateNodeAttachments(ctx, n.ID); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error(err)
		***REMOVED***

		// The assumption here is that all dependent objects
		// have been cleaned up when we are here so the only
		// thing that needs to happen is free the network
		// resources.
		if err := nc.nwkAllocator.Deallocate(n); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed during network free for network %s", n.ID)
		***REMOVED*** else ***REMOVED***
			nc.somethingWasDeallocated = true
		***REMOVED***

		delete(nc.unallocatedNetworks, n.ID)
	case api.EventCreateService:
		var s *api.Service
		a.store.View(func(tx store.ReadTx) ***REMOVED***
			s = store.GetService(tx, v.Service.ID)
		***REMOVED***)

		if s == nil ***REMOVED***
			break
		***REMOVED***

		if nc.nwkAllocator.IsServiceAllocated(s) ***REMOVED***
			break
		***REMOVED***

		if err := a.allocateService(ctx, s); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed allocation for service %s", s.ID)
			break
		***REMOVED***

		if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
			return a.commitAllocatedService(ctx, batch, s)
		***REMOVED***); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed to commit allocation for service %s", s.ID)
		***REMOVED***
	case api.EventUpdateService:
		// We may have already allocated this service. If a create or
		// update event is older than the current version in the store,
		// we run the risk of allocating the service a second time.
		// Only operate on the latest version of the service.
		var s *api.Service
		a.store.View(func(tx store.ReadTx) ***REMOVED***
			s = store.GetService(tx, v.Service.ID)
		***REMOVED***)

		if s == nil ***REMOVED***
			break
		***REMOVED***

		if nc.nwkAllocator.IsServiceAllocated(s) ***REMOVED***
			if !nc.nwkAllocator.HostPublishPortsNeedUpdate(s) ***REMOVED***
				break
			***REMOVED***
			updatePortsInHostPublishMode(s)
		***REMOVED*** else ***REMOVED***
			if err := a.allocateService(ctx, s); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("Failed allocation during update of service %s", s.ID)
				break
			***REMOVED***
		***REMOVED***

		if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
			return a.commitAllocatedService(ctx, batch, s)
		***REMOVED***); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed to commit allocation during update for service %s", s.ID)
			nc.unallocatedServices[s.ID] = s
		***REMOVED*** else ***REMOVED***
			delete(nc.unallocatedServices, s.ID)
		***REMOVED***
	case api.EventDeleteService:
		s := v.Service.Copy()

		if err := nc.nwkAllocator.DeallocateService(s); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed deallocation during delete of service %s", s.ID)
		***REMOVED*** else ***REMOVED***
			nc.somethingWasDeallocated = true
		***REMOVED***

		// Remove it from unallocatedServices just in case
		// it's still there.
		delete(nc.unallocatedServices, s.ID)
	case api.EventCreateNode, api.EventUpdateNode, api.EventDeleteNode:
		a.doNodeAlloc(ctx, ev)
	case api.EventCreateTask, api.EventUpdateTask, api.EventDeleteTask:
		a.doTaskAlloc(ctx, ev)
	case state.EventCommit:
		a.procTasksNetwork(ctx, false)

		if time.Since(nc.lastRetry) > retryInterval || nc.somethingWasDeallocated ***REMOVED***
			a.procUnallocatedNetworks(ctx)
			a.procUnallocatedServices(ctx)
			a.procTasksNetwork(ctx, true)
			nc.lastRetry = time.Now()
			nc.somethingWasDeallocated = false
		***REMOVED***

		// Any left over tasks are moved to the unallocated set
		for _, t := range nc.pendingTasks ***REMOVED***
			nc.unallocatedTasks[t.ID] = t
		***REMOVED***
		nc.pendingTasks = make(map[string]*api.Task)
	***REMOVED***
***REMOVED***

func (a *Allocator) doNodeAlloc(ctx context.Context, ev events.Event) ***REMOVED***
	var (
		isDelete bool
		node     *api.Node
	)

	// We may have already allocated this node. If a create or update
	// event is older than the current version in the store, we run the
	// risk of allocating the node a second time. Only operate on the
	// latest version of the node.
	switch v := ev.(type) ***REMOVED***
	case api.EventCreateNode:
		a.store.View(func(tx store.ReadTx) ***REMOVED***
			node = store.GetNode(tx, v.Node.ID)
		***REMOVED***)
	case api.EventUpdateNode:
		a.store.View(func(tx store.ReadTx) ***REMOVED***
			node = store.GetNode(tx, v.Node.ID)
		***REMOVED***)
	case api.EventDeleteNode:
		isDelete = true
		node = v.Node.Copy()
	***REMOVED***

	if node == nil ***REMOVED***
		return
	***REMOVED***

	nc := a.netCtx

	if isDelete ***REMOVED***
		if err := a.deallocateNode(node); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed freeing network resources for node %s", node.ID)
		***REMOVED*** else ***REMOVED***
			nc.somethingWasDeallocated = true
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		allocatedNetworks, err := a.getAllocatedNetworks()
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Error listing allocated networks in network %s", node.ID)
		***REMOVED***

		isAllocated := a.allocateNode(ctx, node, false, allocatedNetworks)

		if isAllocated ***REMOVED***
			if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
				return a.commitAllocatedNode(ctx, batch, node)
			***REMOVED***); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("Failed to commit allocation of network resources for node %s", node.ID)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func isOverlayNetwork(n *api.Network) bool ***REMOVED***
	if n.DriverState != nil && n.DriverState.Name == "overlay" ***REMOVED***
		return true
	***REMOVED***

	if n.Spec.DriverConfig != nil && n.Spec.DriverConfig.Name == "overlay" ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

func (a *Allocator) getAllocatedNetworks() ([]*api.Network, error) ***REMOVED***
	var (
		err               error
		nc                = a.netCtx
		na                = nc.nwkAllocator
		allocatedNetworks []*api.Network
	)

	// Find allocated networks
	var networks []*api.Network
	a.store.View(func(tx store.ReadTx) ***REMOVED***
		networks, err = store.FindNetworks(tx, store.All)
	***REMOVED***)

	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "error listing all networks in store while trying to allocate during init")
	***REMOVED***

	for _, n := range networks ***REMOVED***

		if isOverlayNetwork(n) && na.IsAllocated(n) ***REMOVED***
			allocatedNetworks = append(allocatedNetworks, n)
		***REMOVED***
	***REMOVED***

	return allocatedNetworks, nil
***REMOVED***

func (a *Allocator) allocateNodes(ctx context.Context, existingAddressesOnly bool) error ***REMOVED***
	// Allocate nodes in the store so far before we process watched events.
	var (
		allocatedNodes []*api.Node
		nodes          []*api.Node
		err            error
	)

	a.store.View(func(tx store.ReadTx) ***REMOVED***
		nodes, err = store.FindNodes(tx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error listing all nodes in store while trying to allocate network resources")
	***REMOVED***

	allocatedNetworks, err := a.getAllocatedNetworks()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error listing all nodes in store while trying to allocate network resources")
	***REMOVED***

	for _, node := range nodes ***REMOVED***
		isAllocated := a.allocateNode(ctx, node, existingAddressesOnly, allocatedNetworks)
		if isAllocated ***REMOVED***
			allocatedNodes = append(allocatedNodes, node)
		***REMOVED***
	***REMOVED***

	if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, node := range allocatedNodes ***REMOVED***
			if err := a.commitAllocatedNode(ctx, batch, node); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("Failed to commit allocation of network resources for node %s", node.ID)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("Failed to commit allocation of network resources for nodes")
	***REMOVED***

	return nil
***REMOVED***

func (a *Allocator) deallocateNodes(ctx context.Context) error ***REMOVED***
	var (
		nodes []*api.Node
		nc    = a.netCtx
		err   error
	)

	a.store.View(func(tx store.ReadTx) ***REMOVED***
		nodes, err = store.FindNodes(tx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return fmt.Errorf("error listing all nodes in store while trying to free network resources")
	***REMOVED***

	for _, node := range nodes ***REMOVED***
		if err := a.deallocateNode(node); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed freeing network resources for node %s", node.ID)
		***REMOVED*** else ***REMOVED***
			nc.somethingWasDeallocated = true
		***REMOVED***
		if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
			return a.commitAllocatedNode(ctx, batch, node)
		***REMOVED***); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed to commit deallocation of network resources for node %s", node.ID)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (a *Allocator) deallocateNodeAttachments(ctx context.Context, nid string) error ***REMOVED***
	var (
		nodes []*api.Node
		nc    = a.netCtx
		err   error
	)

	a.store.View(func(tx store.ReadTx) ***REMOVED***
		nodes, err = store.FindNodes(tx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return fmt.Errorf("error listing all nodes in store while trying to free network resources")
	***REMOVED***

	for _, node := range nodes ***REMOVED***

		var networkAttachment *api.NetworkAttachment
		var naIndex int
		for index, na := range node.Attachments ***REMOVED***
			if na.Network.ID == nid ***REMOVED***
				networkAttachment = na
				naIndex = index
				break
			***REMOVED***
		***REMOVED***

		if networkAttachment == nil ***REMOVED***
			log.G(ctx).Errorf("Failed to find network %s on node %s", nid, node.ID)
			continue
		***REMOVED***

		if nc.nwkAllocator.IsAttachmentAllocated(node, networkAttachment) ***REMOVED***
			if err := nc.nwkAllocator.DeallocateAttachment(node, networkAttachment); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("Failed to commit deallocation of network resources for node %s", node.ID)
			***REMOVED*** else ***REMOVED***

				// Delete the lbattachment
				node.Attachments[naIndex] = node.Attachments[len(node.Attachments)-1]
				node.Attachments[len(node.Attachments)-1] = nil
				node.Attachments = node.Attachments[:len(node.Attachments)-1]

				if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
					return a.commitAllocatedNode(ctx, batch, node)
				***REMOVED***); err != nil ***REMOVED***
					log.G(ctx).WithError(err).Errorf("Failed to commit deallocation of network resources for node %s", node.ID)
				***REMOVED***

			***REMOVED***
		***REMOVED***

	***REMOVED***
	return nil
***REMOVED***

func (a *Allocator) deallocateNode(node *api.Node) error ***REMOVED***
	var (
		nc = a.netCtx
	)

	for _, na := range node.Attachments ***REMOVED***
		if nc.nwkAllocator.IsAttachmentAllocated(node, na) ***REMOVED***
			if err := nc.nwkAllocator.DeallocateAttachment(node, na); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	node.Attachments = nil

	return nil
***REMOVED***

// allocateServices allocates services in the store so far before we process
// watched events.
func (a *Allocator) allocateServices(ctx context.Context, existingAddressesOnly bool) error ***REMOVED***
	var (
		nc       = a.netCtx
		services []*api.Service
		err      error
	)
	a.store.View(func(tx store.ReadTx) ***REMOVED***
		services, err = store.FindServices(tx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error listing all services in store while trying to allocate during init")
	***REMOVED***

	var allocatedServices []*api.Service
	for _, s := range services ***REMOVED***
		if nc.nwkAllocator.IsServiceAllocated(s, networkallocator.OnInit) ***REMOVED***
			continue
		***REMOVED***

		if existingAddressesOnly &&
			(s.Endpoint == nil ||
				len(s.Endpoint.VirtualIPs) == 0) ***REMOVED***
			continue
		***REMOVED***

		if err := a.allocateService(ctx, s); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed allocating service %s during init", s.ID)
			continue
		***REMOVED***
		allocatedServices = append(allocatedServices, s)
	***REMOVED***

	if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, s := range allocatedServices ***REMOVED***
			if err := a.commitAllocatedService(ctx, batch, s); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("failed committing allocation of service %s during init", s.ID)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed committing allocation of services during init")
	***REMOVED***

	return nil
***REMOVED***

// allocateTasks allocates tasks in the store so far before we started watching.
func (a *Allocator) allocateTasks(ctx context.Context, existingAddressesOnly bool) error ***REMOVED***
	var (
		nc             = a.netCtx
		tasks          []*api.Task
		allocatedTasks []*api.Task
		err            error
	)
	a.store.View(func(tx store.ReadTx) ***REMOVED***
		tasks, err = store.FindTasks(tx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error listing all tasks in store while trying to allocate during init")
	***REMOVED***

	for _, t := range tasks ***REMOVED***
		if t.Status.State > api.TaskStateRunning ***REMOVED***
			continue
		***REMOVED***

		if existingAddressesOnly ***REMOVED***
			hasAddresses := false
			for _, nAttach := range t.Networks ***REMOVED***
				if len(nAttach.Addresses) != 0 ***REMOVED***
					hasAddresses = true
					break
				***REMOVED***
			***REMOVED***
			if !hasAddresses ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		var s *api.Service
		if t.ServiceID != "" ***REMOVED***
			a.store.View(func(tx store.ReadTx) ***REMOVED***
				s = store.GetService(tx, t.ServiceID)
			***REMOVED***)
		***REMOVED***

		// Populate network attachments in the task
		// based on service spec.
		a.taskCreateNetworkAttachments(t, s)

		if taskReadyForNetworkVote(t, s, nc) ***REMOVED***
			if t.Status.State >= api.TaskStatePending ***REMOVED***
				continue
			***REMOVED***

			if a.taskAllocateVote(networkVoter, t.ID) ***REMOVED***
				// If the task is not attached to any network, network
				// allocators job is done. Immediately cast a vote so
				// that the task can be moved to the PENDING state as
				// soon as possible.
				updateTaskStatus(t, api.TaskStatePending, allocatedStatusMessage)
				allocatedTasks = append(allocatedTasks, t)
			***REMOVED***
			continue
		***REMOVED***

		err := a.allocateTask(ctx, t)
		if err == nil ***REMOVED***
			allocatedTasks = append(allocatedTasks, t)
		***REMOVED*** else if err != errNoChanges ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed allocating task %s during init", t.ID)
			nc.unallocatedTasks[t.ID] = t
		***REMOVED***
	***REMOVED***

	if err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, t := range allocatedTasks ***REMOVED***
			if err := a.commitAllocatedTask(ctx, batch, t); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("failed committing allocation of task %s during init", t.ID)
			***REMOVED***
		***REMOVED***

		return nil
	***REMOVED***); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed committing allocation of tasks during init")
	***REMOVED***

	return nil
***REMOVED***

// taskReadyForNetworkVote checks if the task is ready for a network
// vote to move it to PENDING state.
func taskReadyForNetworkVote(t *api.Task, s *api.Service, nc *networkContext) bool ***REMOVED***
	// Task is ready for vote if the following is true:
	//
	// Task has no network attached or networks attached but all
	// of them allocated AND Task's service has no endpoint or
	// network configured or service endpoints have been
	// allocated.
	return (len(t.Networks) == 0 || nc.nwkAllocator.IsTaskAllocated(t)) &&
		(s == nil || nc.nwkAllocator.IsServiceAllocated(s))
***REMOVED***

func taskUpdateNetworks(t *api.Task, networks []*api.NetworkAttachment) ***REMOVED***
	networksCopy := make([]*api.NetworkAttachment, 0, len(networks))
	for _, n := range networks ***REMOVED***
		networksCopy = append(networksCopy, n.Copy())
	***REMOVED***

	t.Networks = networksCopy
***REMOVED***

func taskUpdateEndpoint(t *api.Task, endpoint *api.Endpoint) ***REMOVED***
	t.Endpoint = endpoint.Copy()
***REMOVED***

// IsIngressNetworkNeeded checks whether the service requires the routing-mesh
func IsIngressNetworkNeeded(s *api.Service) bool ***REMOVED***
	return networkallocator.IsIngressNetworkNeeded(s)
***REMOVED***

func (a *Allocator) taskCreateNetworkAttachments(t *api.Task, s *api.Service) ***REMOVED***
	// If task network attachments have already been filled in no
	// need to do anything else.
	if len(t.Networks) != 0 ***REMOVED***
		return
	***REMOVED***

	var networks []*api.NetworkAttachment
	if IsIngressNetworkNeeded(s) && a.netCtx.ingressNetwork != nil ***REMOVED***
		networks = append(networks, &api.NetworkAttachment***REMOVED***Network: a.netCtx.ingressNetwork***REMOVED***)
	***REMOVED***

	a.store.View(func(tx store.ReadTx) ***REMOVED***
		// Always prefer NetworkAttachmentConfig in the TaskSpec
		specNetworks := t.Spec.Networks
		if len(specNetworks) == 0 && s != nil && len(s.Spec.Networks) != 0 ***REMOVED***
			specNetworks = s.Spec.Networks
		***REMOVED***

		for _, na := range specNetworks ***REMOVED***
			n := store.GetNetwork(tx, na.Target)
			if n == nil ***REMOVED***
				continue
			***REMOVED***

			attachment := api.NetworkAttachment***REMOVED***Network: n***REMOVED***
			attachment.Aliases = append(attachment.Aliases, na.Aliases...)
			attachment.Addresses = append(attachment.Addresses, na.Addresses...)
			attachment.DriverAttachmentOpts = na.DriverAttachmentOpts
			networks = append(networks, &attachment)
		***REMOVED***
	***REMOVED***)

	taskUpdateNetworks(t, networks)
***REMOVED***

func (a *Allocator) doTaskAlloc(ctx context.Context, ev events.Event) ***REMOVED***
	var (
		isDelete bool
		t        *api.Task
	)

	// We may have already allocated this task. If a create or update
	// event is older than the current version in the store, we run the
	// risk of allocating the task a second time. Only operate on the
	// latest version of the task.
	switch v := ev.(type) ***REMOVED***
	case api.EventCreateTask:
		a.store.View(func(tx store.ReadTx) ***REMOVED***
			t = store.GetTask(tx, v.Task.ID)
		***REMOVED***)
	case api.EventUpdateTask:
		a.store.View(func(tx store.ReadTx) ***REMOVED***
			t = store.GetTask(tx, v.Task.ID)
		***REMOVED***)
	case api.EventDeleteTask:
		isDelete = true
		t = v.Task.Copy()
	***REMOVED***

	if t == nil ***REMOVED***
		return
	***REMOVED***

	nc := a.netCtx

	// If the task has stopped running then we should free the network
	// resources associated with the task right away.
	if t.Status.State > api.TaskStateRunning || isDelete ***REMOVED***
		if nc.nwkAllocator.IsTaskAllocated(t) ***REMOVED***
			if err := nc.nwkAllocator.DeallocateTask(t); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("Failed freeing network resources for task %s", t.ID)
			***REMOVED*** else ***REMOVED***
				nc.somethingWasDeallocated = true
			***REMOVED***
		***REMOVED***

		// Cleanup any task references that might exist
		delete(nc.pendingTasks, t.ID)
		delete(nc.unallocatedTasks, t.ID)

		return
	***REMOVED***

	// If we are already in allocated state, there is
	// absolutely nothing else to do.
	if t.Status.State >= api.TaskStatePending ***REMOVED***
		delete(nc.pendingTasks, t.ID)
		delete(nc.unallocatedTasks, t.ID)
		return
	***REMOVED***

	var s *api.Service
	if t.ServiceID != "" ***REMOVED***
		a.store.View(func(tx store.ReadTx) ***REMOVED***
			s = store.GetService(tx, t.ServiceID)
		***REMOVED***)
		if s == nil ***REMOVED***
			// If the task is running it is not normal to
			// not be able to find the associated
			// service. If the task is not running (task
			// is either dead or the desired state is set
			// to dead) then the service may not be
			// available in store. But we still need to
			// cleanup network resources associated with
			// the task.
			if t.Status.State <= api.TaskStateRunning && !isDelete ***REMOVED***
				log.G(ctx).Errorf("Event %T: Failed to get service %s for task %s state %s: could not find service %s", ev, t.ServiceID, t.ID, t.Status.State, t.ServiceID)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Populate network attachments in the task
	// based on service spec.
	a.taskCreateNetworkAttachments(t, s)

	nc.pendingTasks[t.ID] = t
***REMOVED***

func (a *Allocator) allocateNode(ctx context.Context, node *api.Node, existingAddressesOnly bool, networks []*api.Network) bool ***REMOVED***
	var allocated bool

	nc := a.netCtx

	for _, network := range networks ***REMOVED***

		var lbAttachment *api.NetworkAttachment
		for _, na := range node.Attachments ***REMOVED***
			if na.Network != nil && na.Network.ID == network.ID ***REMOVED***
				lbAttachment = na
				break
			***REMOVED***
		***REMOVED***

		if lbAttachment != nil ***REMOVED***
			if nc.nwkAllocator.IsAttachmentAllocated(node, lbAttachment) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if lbAttachment == nil ***REMOVED***
			lbAttachment = &api.NetworkAttachment***REMOVED******REMOVED***
			node.Attachments = append(node.Attachments, lbAttachment)
		***REMOVED***

		if existingAddressesOnly && len(lbAttachment.Addresses) == 0 ***REMOVED***
			continue
		***REMOVED***

		lbAttachment.Network = network.Copy()
		if err := a.netCtx.nwkAllocator.AllocateAttachment(node, lbAttachment); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("Failed to allocate network resources for node %s", node.ID)
			// TODO: Should we add a unallocatedNode and retry allocating resources like we do for network, tasks, services?
			// right now, we will only retry allocating network resources for the node when the node is updated.
			continue
		***REMOVED***

		allocated = true
	***REMOVED***
	return allocated

***REMOVED***

func (a *Allocator) commitAllocatedNode(ctx context.Context, batch *store.Batch, node *api.Node) error ***REMOVED***
	if err := batch.Update(func(tx store.Tx) error ***REMOVED***
		err := store.UpdateNode(tx, node)

		if err == store.ErrSequenceConflict ***REMOVED***
			storeNode := store.GetNode(tx, node.ID)
			storeNode.Attachments = node.Attachments
			err = store.UpdateNode(tx, storeNode)
		***REMOVED***

		return errors.Wrapf(err, "failed updating state in store transaction for node %s", node.ID)
	***REMOVED***); err != nil ***REMOVED***
		if err := a.deallocateNode(node); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed rolling back allocation of node %s", node.ID)
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

// This function prepares the service object for being updated when the change regards
// the published ports in host mode: It resets the runtime state ports (s.Endpoint.Ports)
// to the current ingress mode runtime state ports plus the newly configured publish mode ports,
// so that the service allocation invoked on this new service object will trigger the deallocation
// of any old publish mode port and allocation of any new one.
func updatePortsInHostPublishMode(s *api.Service) ***REMOVED***
	// First, remove all host-mode ports from s.Endpoint.Ports
	if s.Endpoint != nil ***REMOVED***
		var portConfigs []*api.PortConfig
		for _, portConfig := range s.Endpoint.Ports ***REMOVED***
			if portConfig.PublishMode != api.PublishModeHost ***REMOVED***
				portConfigs = append(portConfigs, portConfig)
			***REMOVED***
		***REMOVED***
		s.Endpoint.Ports = portConfigs
	***REMOVED***

	// Add back all host-mode ports
	if s.Spec.Endpoint != nil ***REMOVED***
		if s.Endpoint == nil ***REMOVED***
			s.Endpoint = &api.Endpoint***REMOVED******REMOVED***
		***REMOVED***
		for _, portConfig := range s.Spec.Endpoint.Ports ***REMOVED***
			if portConfig.PublishMode == api.PublishModeHost ***REMOVED***
				s.Endpoint.Ports = append(s.Endpoint.Ports, portConfig.Copy())
			***REMOVED***
		***REMOVED***
	***REMOVED***
	s.Endpoint.Spec = s.Spec.Endpoint.Copy()
***REMOVED***

func (a *Allocator) allocateService(ctx context.Context, s *api.Service) error ***REMOVED***
	nc := a.netCtx

	if s.Spec.Endpoint != nil ***REMOVED***
		// service has user-defined endpoint
		if s.Endpoint == nil ***REMOVED***
			// service currently has no allocated endpoint, need allocated.
			s.Endpoint = &api.Endpoint***REMOVED***
				Spec: s.Spec.Endpoint.Copy(),
			***REMOVED***
		***REMOVED***

		// The service is trying to expose ports to the external
		// world. Automatically attach the service to the ingress
		// network only if it is not already done.
		if IsIngressNetworkNeeded(s) ***REMOVED***
			if nc.ingressNetwork == nil ***REMOVED***
				return fmt.Errorf("ingress network is missing")
			***REMOVED***
			var found bool
			for _, vip := range s.Endpoint.VirtualIPs ***REMOVED***
				if vip.NetworkID == nc.ingressNetwork.ID ***REMOVED***
					found = true
					break
				***REMOVED***
			***REMOVED***

			if !found ***REMOVED***
				s.Endpoint.VirtualIPs = append(s.Endpoint.VirtualIPs,
					&api.Endpoint_VirtualIP***REMOVED***NetworkID: nc.ingressNetwork.ID***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if s.Endpoint != nil ***REMOVED***
		// service has no user-defined endpoints while has already allocated network resources,
		// need deallocated.
		if err := nc.nwkAllocator.DeallocateService(s); err != nil ***REMOVED***
			return err
		***REMOVED***
		nc.somethingWasDeallocated = true
	***REMOVED***

	if err := nc.nwkAllocator.AllocateService(s); err != nil ***REMOVED***
		nc.unallocatedServices[s.ID] = s
		return err
	***REMOVED***

	// If the service doesn't expose ports any more and if we have
	// any lingering virtual IP references for ingress network
	// clean them up here.
	if !IsIngressNetworkNeeded(s) && nc.ingressNetwork != nil ***REMOVED***
		if s.Endpoint != nil ***REMOVED***
			for i, vip := range s.Endpoint.VirtualIPs ***REMOVED***
				if vip.NetworkID == nc.ingressNetwork.ID ***REMOVED***
					n := len(s.Endpoint.VirtualIPs)
					s.Endpoint.VirtualIPs[i], s.Endpoint.VirtualIPs[n-1] = s.Endpoint.VirtualIPs[n-1], nil
					s.Endpoint.VirtualIPs = s.Endpoint.VirtualIPs[:n-1]
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (a *Allocator) commitAllocatedService(ctx context.Context, batch *store.Batch, s *api.Service) error ***REMOVED***
	if err := batch.Update(func(tx store.Tx) error ***REMOVED***
		err := store.UpdateService(tx, s)

		if err == store.ErrSequenceConflict ***REMOVED***
			storeService := store.GetService(tx, s.ID)
			storeService.Endpoint = s.Endpoint
			err = store.UpdateService(tx, storeService)
		***REMOVED***

		return errors.Wrapf(err, "failed updating state in store transaction for service %s", s.ID)
	***REMOVED***); err != nil ***REMOVED***
		if err := a.netCtx.nwkAllocator.DeallocateService(s); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed rolling back allocation of service %s", s.ID)
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

func (a *Allocator) allocateNetwork(ctx context.Context, n *api.Network) error ***REMOVED***
	nc := a.netCtx

	if err := nc.nwkAllocator.Allocate(n); err != nil ***REMOVED***
		nc.unallocatedNetworks[n.ID] = n
		return err
	***REMOVED***

	return nil
***REMOVED***

func (a *Allocator) commitAllocatedNetwork(ctx context.Context, batch *store.Batch, n *api.Network) error ***REMOVED***
	if err := batch.Update(func(tx store.Tx) error ***REMOVED***
		if err := store.UpdateNetwork(tx, n); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed updating state in store transaction for network %s", n.ID)
		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		if err := a.netCtx.nwkAllocator.Deallocate(n); err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed rolling back allocation of network %s", n.ID)
		***REMOVED***

		return err
	***REMOVED***

	return nil
***REMOVED***

func (a *Allocator) allocateTask(ctx context.Context, t *api.Task) (err error) ***REMOVED***
	taskUpdated := false
	nc := a.netCtx

	// We might be here even if a task allocation has already
	// happened but wasn't successfully committed to store. In such
	// cases skip allocation and go straight ahead to updating the
	// store.
	if !nc.nwkAllocator.IsTaskAllocated(t) ***REMOVED***
		a.store.View(func(tx store.ReadTx) ***REMOVED***
			if t.ServiceID != "" ***REMOVED***
				s := store.GetService(tx, t.ServiceID)
				if s == nil ***REMOVED***
					err = fmt.Errorf("could not find service %s", t.ServiceID)
					return
				***REMOVED***

				if !nc.nwkAllocator.IsServiceAllocated(s) ***REMOVED***
					err = fmt.Errorf("service %s to which this task %s belongs has pending allocations", s.ID, t.ID)
					return
				***REMOVED***

				if s.Endpoint != nil ***REMOVED***
					taskUpdateEndpoint(t, s.Endpoint)
					taskUpdated = true
				***REMOVED***
			***REMOVED***

			for _, na := range t.Networks ***REMOVED***
				n := store.GetNetwork(tx, na.Network.ID)
				if n == nil ***REMOVED***
					err = fmt.Errorf("failed to retrieve network %s while allocating task %s", na.Network.ID, t.ID)
					return
				***REMOVED***

				if !nc.nwkAllocator.IsAllocated(n) ***REMOVED***
					err = fmt.Errorf("network %s attached to task %s not allocated yet", n.ID, t.ID)
					return
				***REMOVED***

				na.Network = n
			***REMOVED***

			if err = nc.nwkAllocator.AllocateTask(t); err != nil ***REMOVED***
				return
			***REMOVED***
			if nc.nwkAllocator.IsTaskAllocated(t) ***REMOVED***
				taskUpdated = true
			***REMOVED***
		***REMOVED***)

		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Update the network allocations and moving to
	// PENDING state on top of the latest store state.
	if a.taskAllocateVote(networkVoter, t.ID) ***REMOVED***
		if t.Status.State < api.TaskStatePending ***REMOVED***
			updateTaskStatus(t, api.TaskStatePending, allocatedStatusMessage)
			taskUpdated = true
		***REMOVED***
	***REMOVED***

	if !taskUpdated ***REMOVED***
		return errNoChanges
	***REMOVED***

	return nil
***REMOVED***

func (a *Allocator) commitAllocatedTask(ctx context.Context, batch *store.Batch, t *api.Task) error ***REMOVED***
	return batch.Update(func(tx store.Tx) error ***REMOVED***
		err := store.UpdateTask(tx, t)

		if err == store.ErrSequenceConflict ***REMOVED***
			storeTask := store.GetTask(tx, t.ID)
			taskUpdateNetworks(storeTask, t.Networks)
			taskUpdateEndpoint(storeTask, t.Endpoint)
			if storeTask.Status.State < api.TaskStatePending ***REMOVED***
				storeTask.Status = t.Status
			***REMOVED***
			err = store.UpdateTask(tx, storeTask)
		***REMOVED***

		return errors.Wrapf(err, "failed updating state in store transaction for task %s", t.ID)
	***REMOVED***)
***REMOVED***

func (a *Allocator) procUnallocatedNetworks(ctx context.Context) ***REMOVED***
	nc := a.netCtx
	var allocatedNetworks []*api.Network
	for _, n := range nc.unallocatedNetworks ***REMOVED***
		if !nc.nwkAllocator.IsAllocated(n) ***REMOVED***
			if err := a.allocateNetwork(ctx, n); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Debugf("Failed allocation of unallocated network %s", n.ID)
				continue
			***REMOVED***
			allocatedNetworks = append(allocatedNetworks, n)
		***REMOVED***
	***REMOVED***

	if len(allocatedNetworks) == 0 ***REMOVED***
		return
	***REMOVED***

	err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, n := range allocatedNetworks ***REMOVED***
			if err := a.commitAllocatedNetwork(ctx, batch, n); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Debugf("Failed to commit allocation of unallocated network %s", n.ID)
				continue
			***REMOVED***
			delete(nc.unallocatedNetworks, n.ID)
		***REMOVED***
		return nil
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("Failed to commit allocation of unallocated networks")
		// We optimistically removed these from nc.unallocatedNetworks
		// above in anticipation of successfully committing the batch,
		// but since the transaction has failed, we requeue them here.
		for _, n := range allocatedNetworks ***REMOVED***
			nc.unallocatedNetworks[n.ID] = n
		***REMOVED***
	***REMOVED***
***REMOVED***

func (a *Allocator) procUnallocatedServices(ctx context.Context) ***REMOVED***
	nc := a.netCtx
	var allocatedServices []*api.Service
	for _, s := range nc.unallocatedServices ***REMOVED***
		if !nc.nwkAllocator.IsServiceAllocated(s) ***REMOVED***
			if err := a.allocateService(ctx, s); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Debugf("Failed allocation of unallocated service %s", s.ID)
				continue
			***REMOVED***
			allocatedServices = append(allocatedServices, s)
		***REMOVED***
	***REMOVED***

	if len(allocatedServices) == 0 ***REMOVED***
		return
	***REMOVED***

	err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, s := range allocatedServices ***REMOVED***
			if err := a.commitAllocatedService(ctx, batch, s); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Debugf("Failed to commit allocation of unallocated service %s", s.ID)
				continue
			***REMOVED***
			delete(nc.unallocatedServices, s.ID)
		***REMOVED***
		return nil
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("Failed to commit allocation of unallocated services")
		// We optimistically removed these from nc.unallocatedServices
		// above in anticipation of successfully committing the batch,
		// but since the transaction has failed, we requeue them here.
		for _, s := range allocatedServices ***REMOVED***
			nc.unallocatedServices[s.ID] = s
		***REMOVED***
	***REMOVED***
***REMOVED***

func (a *Allocator) procTasksNetwork(ctx context.Context, onRetry bool) ***REMOVED***
	nc := a.netCtx
	quiet := false
	toAllocate := nc.pendingTasks
	if onRetry ***REMOVED***
		toAllocate = nc.unallocatedTasks
		quiet = true
	***REMOVED***
	allocatedTasks := make([]*api.Task, 0, len(toAllocate))

	for _, t := range toAllocate ***REMOVED***
		if err := a.allocateTask(ctx, t); err == nil ***REMOVED***
			allocatedTasks = append(allocatedTasks, t)
		***REMOVED*** else if err != errNoChanges ***REMOVED***
			if quiet ***REMOVED***
				log.G(ctx).WithError(err).Debug("task allocation failure")
			***REMOVED*** else ***REMOVED***
				log.G(ctx).WithError(err).Error("task allocation failure")
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if len(allocatedTasks) == 0 ***REMOVED***
		return
	***REMOVED***

	err := a.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, t := range allocatedTasks ***REMOVED***
			err := a.commitAllocatedTask(ctx, batch, t)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("task allocation commit failure")
				continue
			***REMOVED***
			delete(toAllocate, t.ID)
		***REMOVED***

		return nil
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed a store batch operation while processing tasks")
		// We optimistically removed these from toAllocate above in
		// anticipation of successfully committing the batch, but since
		// the transaction has failed, we requeue them here.
		for _, t := range allocatedTasks ***REMOVED***
			toAllocate[t.ID] = t
		***REMOVED***
	***REMOVED***
***REMOVED***

// IsBuiltInNetworkDriver returns whether the passed driver is an internal network driver
func IsBuiltInNetworkDriver(name string) bool ***REMOVED***
	return cnmallocator.IsBuiltInDriver(name)
***REMOVED***

// PredefinedNetworks returns the list of predefined network structures for a given network model
func PredefinedNetworks() []networkallocator.PredefinedNetworkData ***REMOVED***
	return cnmallocator.PredefinedNetworks()
***REMOVED***

// updateTaskStatus sets TaskStatus and updates timestamp.
func updateTaskStatus(t *api.Task, newStatus api.TaskState, message string) ***REMOVED***
	t.Status = api.TaskStatus***REMOVED***
		State:     newStatus,
		Message:   message,
		Timestamp: ptypes.MustTimestampProto(time.Now()),
	***REMOVED***
***REMOVED***

// IsIngressNetwork returns whether the passed network is an ingress network.
func IsIngressNetwork(nw *api.Network) bool ***REMOVED***
	return networkallocator.IsIngressNetwork(nw)
***REMOVED***

// GetIngressNetwork fetches the ingress network from store.
// ErrNoIngress will be returned if the ingress network is not present,
// nil otherwise. In case of any other failure in accessing the store,
// the respective error will be reported as is.
func GetIngressNetwork(s *store.MemoryStore) (*api.Network, error) ***REMOVED***
	var (
		networks []*api.Network
		err      error
	)
	s.View(func(tx store.ReadTx) ***REMOVED***
		networks, err = store.FindNetworks(tx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, n := range networks ***REMOVED***
		if IsIngressNetwork(n) ***REMOVED***
			return n, nil
		***REMOVED***
	***REMOVED***
	return nil, ErrNoIngress
***REMOVED***
