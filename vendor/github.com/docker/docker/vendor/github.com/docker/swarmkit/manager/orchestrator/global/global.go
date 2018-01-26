package global

import (
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/constraint"
	"github.com/docker/swarmkit/manager/orchestrator"
	"github.com/docker/swarmkit/manager/orchestrator/restart"
	"github.com/docker/swarmkit/manager/orchestrator/taskinit"
	"github.com/docker/swarmkit/manager/orchestrator/update"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

type globalService struct ***REMOVED***
	*api.Service

	// Compiled constraints
	constraints []constraint.Constraint
***REMOVED***

// Orchestrator runs a reconciliation loop to create and destroy tasks as
// necessary for global services.
type Orchestrator struct ***REMOVED***
	store *store.MemoryStore
	// nodes is the set of non-drained nodes in the cluster, indexed by node ID
	nodes map[string]*api.Node
	// globalServices has all the global services in the cluster, indexed by ServiceID
	globalServices map[string]globalService
	restartTasks   map[string]struct***REMOVED******REMOVED***

	// stopChan signals to the state machine to stop running.
	stopChan chan struct***REMOVED******REMOVED***
	// doneChan is closed when the state machine terminates.
	doneChan chan struct***REMOVED******REMOVED***

	updater  *update.Supervisor
	restarts *restart.Supervisor

	cluster *api.Cluster // local instance of the cluster
***REMOVED***

// NewGlobalOrchestrator creates a new global Orchestrator
func NewGlobalOrchestrator(store *store.MemoryStore) *Orchestrator ***REMOVED***
	restartSupervisor := restart.NewSupervisor(store)
	updater := update.NewSupervisor(store, restartSupervisor)
	return &Orchestrator***REMOVED***
		store:          store,
		nodes:          make(map[string]*api.Node),
		globalServices: make(map[string]globalService),
		stopChan:       make(chan struct***REMOVED******REMOVED***),
		doneChan:       make(chan struct***REMOVED******REMOVED***),
		updater:        updater,
		restarts:       restartSupervisor,
		restartTasks:   make(map[string]struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

func (g *Orchestrator) initTasks(ctx context.Context, readTx store.ReadTx) error ***REMOVED***
	return taskinit.CheckTasks(ctx, g.store, readTx, g, g.restarts)
***REMOVED***

// Run contains the global orchestrator event loop
func (g *Orchestrator) Run(ctx context.Context) error ***REMOVED***
	defer close(g.doneChan)

	// Watch changes to services and tasks
	queue := g.store.WatchQueue()
	watcher, cancel := queue.Watch()
	defer cancel()

	// lookup the cluster
	var err error
	g.store.View(func(readTx store.ReadTx) ***REMOVED***
		var clusters []*api.Cluster
		clusters, err = store.FindClusters(readTx, store.ByName(store.DefaultClusterName))

		if len(clusters) != 1 ***REMOVED***
			return // just pick up the cluster when it is created.
		***REMOVED***
		g.cluster = clusters[0]
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Get list of nodes
	var nodes []*api.Node
	g.store.View(func(readTx store.ReadTx) ***REMOVED***
		nodes, err = store.FindNodes(readTx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, n := range nodes ***REMOVED***
		g.updateNode(n)
	***REMOVED***

	// Lookup global services
	var existingServices []*api.Service
	g.store.View(func(readTx store.ReadTx) ***REMOVED***
		existingServices, err = store.FindServices(readTx, store.All)
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var reconcileServiceIDs []string
	for _, s := range existingServices ***REMOVED***
		if orchestrator.IsGlobalService(s) ***REMOVED***
			g.updateService(s)
			reconcileServiceIDs = append(reconcileServiceIDs, s.ID)
		***REMOVED***
	***REMOVED***

	// fix tasks in store before reconciliation loop
	g.store.View(func(readTx store.ReadTx) ***REMOVED***
		err = g.initTasks(ctx, readTx)
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	g.tickTasks(ctx)
	g.reconcileServices(ctx, reconcileServiceIDs)

	for ***REMOVED***
		select ***REMOVED***
		case event := <-watcher:
			// TODO(stevvooe): Use ctx to limit running time of operation.
			switch v := event.(type) ***REMOVED***
			case api.EventUpdateCluster:
				g.cluster = v.Cluster
			case api.EventCreateService:
				if !orchestrator.IsGlobalService(v.Service) ***REMOVED***
					continue
				***REMOVED***
				g.updateService(v.Service)
				g.reconcileServices(ctx, []string***REMOVED***v.Service.ID***REMOVED***)
			case api.EventUpdateService:
				if !orchestrator.IsGlobalService(v.Service) ***REMOVED***
					continue
				***REMOVED***
				g.updateService(v.Service)
				g.reconcileServices(ctx, []string***REMOVED***v.Service.ID***REMOVED***)
			case api.EventDeleteService:
				if !orchestrator.IsGlobalService(v.Service) ***REMOVED***
					continue
				***REMOVED***
				orchestrator.SetServiceTasksRemove(ctx, g.store, v.Service)
				// delete the service from service map
				delete(g.globalServices, v.Service.ID)
				g.restarts.ClearServiceHistory(v.Service.ID)
			case api.EventCreateNode:
				g.updateNode(v.Node)
				g.reconcileOneNode(ctx, v.Node)
			case api.EventUpdateNode:
				g.updateNode(v.Node)
				g.reconcileOneNode(ctx, v.Node)
			case api.EventDeleteNode:
				g.foreachTaskFromNode(ctx, v.Node, g.deleteTask)
				delete(g.nodes, v.Node.ID)
			case api.EventUpdateTask:
				g.handleTaskChange(ctx, v.Task)
			***REMOVED***
		case <-g.stopChan:
			return nil
		***REMOVED***
		g.tickTasks(ctx)
	***REMOVED***
***REMOVED***

// FixTask validates a task with the current cluster settings, and takes
// action to make it conformant to node state and service constraint
// it's called at orchestrator initialization
func (g *Orchestrator) FixTask(ctx context.Context, batch *store.Batch, t *api.Task) ***REMOVED***
	if _, exists := g.globalServices[t.ServiceID]; !exists ***REMOVED***
		return
	***REMOVED***
	// if a task's DesiredState has past running, the task has been processed
	if t.DesiredState > api.TaskStateRunning ***REMOVED***
		return
	***REMOVED***

	var node *api.Node
	if t.NodeID != "" ***REMOVED***
		node = g.nodes[t.NodeID]
	***REMOVED***
	// if the node no longer valid, remove the task
	if t.NodeID == "" || orchestrator.InvalidNode(node) ***REMOVED***
		g.shutdownTask(ctx, batch, t)
		return
	***REMOVED***

	// restart a task if it fails
	if t.Status.State > api.TaskStateRunning ***REMOVED***
		g.restartTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// handleTaskChange defines what orchestrator does when a task is updated by agent
func (g *Orchestrator) handleTaskChange(ctx context.Context, t *api.Task) ***REMOVED***
	if _, exists := g.globalServices[t.ServiceID]; !exists ***REMOVED***
		return
	***REMOVED***
	// if a task's DesiredState has passed running, it
	// means the task has been processed
	if t.DesiredState > api.TaskStateRunning ***REMOVED***
		return
	***REMOVED***

	// if a task has passed running, restart it
	if t.Status.State > api.TaskStateRunning ***REMOVED***
		g.restartTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// Stop stops the orchestrator.
func (g *Orchestrator) Stop() ***REMOVED***
	close(g.stopChan)
	<-g.doneChan
	g.updater.CancelAll()
	g.restarts.CancelAll()
***REMOVED***

func (g *Orchestrator) foreachTaskFromNode(ctx context.Context, node *api.Node, cb func(context.Context, *store.Batch, *api.Task)) ***REMOVED***
	var (
		tasks []*api.Task
		err   error
	)
	g.store.View(func(tx store.ReadTx) ***REMOVED***
		tasks, err = store.FindTasks(tx, store.ByNodeID(node.ID))
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("global orchestrator: foreachTaskFromNode failed finding tasks")
		return
	***REMOVED***

	err = g.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, t := range tasks ***REMOVED***
			// Global orchestrator only removes tasks from globalServices
			if _, exists := g.globalServices[t.ServiceID]; exists ***REMOVED***
				cb(ctx, batch, t)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("global orchestrator: foreachTaskFromNode failed batching tasks")
	***REMOVED***
***REMOVED***

func (g *Orchestrator) reconcileServices(ctx context.Context, serviceIDs []string) ***REMOVED***
	nodeTasks := make(map[string]map[string][]*api.Task)

	g.store.View(func(tx store.ReadTx) ***REMOVED***
		for _, serviceID := range serviceIDs ***REMOVED***
			service := g.globalServices[serviceID].Service
			if service == nil ***REMOVED***
				continue
			***REMOVED***

			tasks, err := store.FindTasks(tx, store.ByServiceID(serviceID))
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("global orchestrator: reconcileServices failed finding tasks for service %s", serviceID)
				continue
			***REMOVED***

			// nodeID -> task list
			nodeTasks[serviceID] = make(map[string][]*api.Task)

			for _, t := range tasks ***REMOVED***
				nodeTasks[serviceID][t.NodeID] = append(nodeTasks[serviceID][t.NodeID], t)
			***REMOVED***

			// Keep all runnable instances of this service,
			// and instances that were not be restarted due
			// to restart policy but may be updated if the
			// service spec changed.
			for nodeID, slot := range nodeTasks[serviceID] ***REMOVED***
				updatable := g.restarts.UpdatableTasksInSlot(ctx, slot, g.globalServices[serviceID].Service)
				if len(updatable) != 0 ***REMOVED***
					nodeTasks[serviceID][nodeID] = updatable
				***REMOVED*** else ***REMOVED***
					delete(nodeTasks[serviceID], nodeID)
				***REMOVED***
			***REMOVED***

		***REMOVED***
	***REMOVED***)

	updates := make(map[*api.Service][]orchestrator.Slot)

	err := g.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for _, serviceID := range serviceIDs ***REMOVED***
			var updateTasks []orchestrator.Slot

			if _, exists := nodeTasks[serviceID]; !exists ***REMOVED***
				continue
			***REMOVED***

			service := g.globalServices[serviceID]

			for nodeID, node := range g.nodes ***REMOVED***
				meetsConstraints := constraint.NodeMatches(service.constraints, node)
				ntasks := nodeTasks[serviceID][nodeID]
				delete(nodeTasks[serviceID], nodeID)

				if !meetsConstraints ***REMOVED***
					g.shutdownTasks(ctx, batch, ntasks)
					continue
				***REMOVED***

				if node.Spec.Availability == api.NodeAvailabilityPause ***REMOVED***
					// the node is paused, so we won't add or update
					// any tasks
					continue
				***REMOVED***

				// this node needs to run 1 copy of the task
				if len(ntasks) == 0 ***REMOVED***
					g.addTask(ctx, batch, service.Service, nodeID)
				***REMOVED*** else ***REMOVED***
					updateTasks = append(updateTasks, ntasks)
				***REMOVED***
			***REMOVED***

			if len(updateTasks) > 0 ***REMOVED***
				updates[service.Service] = updateTasks
			***REMOVED***

			// Remove any tasks assigned to nodes not found in g.nodes.
			// These must be associated with nodes that are drained, or
			// nodes that no longer exist.
			for _, ntasks := range nodeTasks[serviceID] ***REMOVED***
				g.shutdownTasks(ctx, batch, ntasks)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("global orchestrator: reconcileServices transaction failed")
	***REMOVED***

	for service, updateTasks := range updates ***REMOVED***
		g.updater.Update(ctx, g.cluster, service, updateTasks)
	***REMOVED***
***REMOVED***

// updateNode updates g.nodes based on the current node value
func (g *Orchestrator) updateNode(node *api.Node) ***REMOVED***
	if node.Spec.Availability == api.NodeAvailabilityDrain || node.Status.State == api.NodeStatus_DOWN ***REMOVED***
		delete(g.nodes, node.ID)
	***REMOVED*** else ***REMOVED***
		g.nodes[node.ID] = node
	***REMOVED***
***REMOVED***

// updateService updates g.globalServices based on the current service value
func (g *Orchestrator) updateService(service *api.Service) ***REMOVED***
	var constraints []constraint.Constraint

	if service.Spec.Task.Placement != nil && len(service.Spec.Task.Placement.Constraints) != 0 ***REMOVED***
		constraints, _ = constraint.Parse(service.Spec.Task.Placement.Constraints)
	***REMOVED***

	g.globalServices[service.ID] = globalService***REMOVED***
		Service:     service,
		constraints: constraints,
	***REMOVED***
***REMOVED***

// reconcileOneNode checks all global services on one node
func (g *Orchestrator) reconcileOneNode(ctx context.Context, node *api.Node) ***REMOVED***
	if node.Spec.Availability == api.NodeAvailabilityDrain ***REMOVED***
		log.G(ctx).Debugf("global orchestrator: node %s in drain state, shutting down its tasks", node.ID)
		g.foreachTaskFromNode(ctx, node, g.shutdownTask)
		return
	***REMOVED***

	if node.Status.State == api.NodeStatus_DOWN ***REMOVED***
		log.G(ctx).Debugf("global orchestrator: node %s is down, shutting down its tasks", node.ID)
		g.foreachTaskFromNode(ctx, node, g.shutdownTask)
		return
	***REMOVED***

	if node.Spec.Availability == api.NodeAvailabilityPause ***REMOVED***
		// the node is paused, so we won't add or update tasks
		return
	***REMOVED***

	node, exists := g.nodes[node.ID]
	if !exists ***REMOVED***
		return
	***REMOVED***

	// tasks by service
	tasks := make(map[string][]*api.Task)

	var (
		tasksOnNode []*api.Task
		err         error
	)

	g.store.View(func(tx store.ReadTx) ***REMOVED***
		tasksOnNode, err = store.FindTasks(tx, store.ByNodeID(node.ID))
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("global orchestrator: reconcile failed finding tasks on node %s", node.ID)
		return
	***REMOVED***

	for serviceID, service := range g.globalServices ***REMOVED***
		for _, t := range tasksOnNode ***REMOVED***
			if t.ServiceID != serviceID ***REMOVED***
				continue
			***REMOVED***
			tasks[serviceID] = append(tasks[serviceID], t)
		***REMOVED***

		// Keep all runnable instances of this service,
		// and instances that were not be restarted due
		// to restart policy but may be updated if the
		// service spec changed.
		for serviceID, slot := range tasks ***REMOVED***
			updatable := g.restarts.UpdatableTasksInSlot(ctx, slot, service.Service)

			if len(updatable) != 0 ***REMOVED***
				tasks[serviceID] = updatable
			***REMOVED*** else ***REMOVED***
				delete(tasks, serviceID)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	err = g.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for serviceID, service := range g.globalServices ***REMOVED***
			if !constraint.NodeMatches(service.constraints, node) ***REMOVED***
				continue
			***REMOVED***

			if len(tasks) == 0 ***REMOVED***
				g.addTask(ctx, batch, service.Service, node.ID)
			***REMOVED*** else ***REMOVED***
				// If task is out of date, update it. This can happen
				// on node reconciliation if, for example, we pause a
				// node, update the service, and then activate the node
				// later.

				// We don't use g.updater here for two reasons:
				// - This is not a rolling update. Since it was not
				//   triggered directly by updating the service, it
				//   should not observe the rolling update parameters
				//   or show status in UpdateStatus.
				// - Calling Update cancels any current rolling updates
				//   for the service, such as one triggered by service
				//   reconciliation.

				var (
					dirtyTasks []*api.Task
					cleanTasks []*api.Task
				)

				for _, t := range tasks[serviceID] ***REMOVED***
					if orchestrator.IsTaskDirty(service.Service, t) ***REMOVED***
						dirtyTasks = append(dirtyTasks, t)
					***REMOVED*** else ***REMOVED***
						cleanTasks = append(cleanTasks, t)
					***REMOVED***
				***REMOVED***

				if len(cleanTasks) == 0 ***REMOVED***
					g.addTask(ctx, batch, service.Service, node.ID)
				***REMOVED*** else ***REMOVED***
					dirtyTasks = append(dirtyTasks, cleanTasks[1:]...)
				***REMOVED***
				g.shutdownTasks(ctx, batch, dirtyTasks)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("global orchestrator: reconcileServiceOneNode batch failed")
	***REMOVED***
***REMOVED***

func (g *Orchestrator) tickTasks(ctx context.Context) ***REMOVED***
	if len(g.restartTasks) == 0 ***REMOVED***
		return
	***REMOVED***
	err := g.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for taskID := range g.restartTasks ***REMOVED***
			err := batch.Update(func(tx store.Tx) error ***REMOVED***
				t := store.GetTask(tx, taskID)
				if t == nil || t.DesiredState > api.TaskStateRunning ***REMOVED***
					return nil
				***REMOVED***

				service := store.GetService(tx, t.ServiceID)
				if service == nil ***REMOVED***
					return nil
				***REMOVED***

				node, nodeExists := g.nodes[t.NodeID]
				serviceEntry, serviceExists := g.globalServices[t.ServiceID]
				if !nodeExists || !serviceExists ***REMOVED***
					return nil
				***REMOVED***

				if node.Spec.Availability == api.NodeAvailabilityPause ||
					!constraint.NodeMatches(serviceEntry.constraints, node) ***REMOVED***
					t.DesiredState = api.TaskStateShutdown
					return store.UpdateTask(tx, t)
				***REMOVED***

				return g.restarts.Restart(ctx, tx, g.cluster, service, *t)
			***REMOVED***)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).Errorf("orchestrator restartTask transaction failed")
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("global orchestrator: restartTask transaction failed")
	***REMOVED***
	g.restartTasks = make(map[string]struct***REMOVED******REMOVED***)
***REMOVED***

func (g *Orchestrator) shutdownTask(ctx context.Context, batch *store.Batch, t *api.Task) ***REMOVED***
	// set existing task DesiredState to TaskStateShutdown
	// TODO(aaronl): optimistic update?
	err := batch.Update(func(tx store.Tx) error ***REMOVED***
		t = store.GetTask(tx, t.ID)
		if t != nil && t.DesiredState < api.TaskStateShutdown ***REMOVED***
			t.DesiredState = api.TaskStateShutdown
			return store.UpdateTask(tx, t)
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("global orchestrator: shutdownTask failed to shut down %s", t.ID)
	***REMOVED***
***REMOVED***

func (g *Orchestrator) addTask(ctx context.Context, batch *store.Batch, service *api.Service, nodeID string) ***REMOVED***
	task := orchestrator.NewTask(g.cluster, service, 0, nodeID)

	err := batch.Update(func(tx store.Tx) error ***REMOVED***
		if store.GetService(tx, service.ID) == nil ***REMOVED***
			return nil
		***REMOVED***
		return store.CreateTask(tx, task)
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("global orchestrator: failed to create task")
	***REMOVED***
***REMOVED***

func (g *Orchestrator) shutdownTasks(ctx context.Context, batch *store.Batch, tasks []*api.Task) ***REMOVED***
	for _, t := range tasks ***REMOVED***
		g.shutdownTask(ctx, batch, t)
	***REMOVED***
***REMOVED***

func (g *Orchestrator) deleteTask(ctx context.Context, batch *store.Batch, t *api.Task) ***REMOVED***
	err := batch.Update(func(tx store.Tx) error ***REMOVED***
		return store.DeleteTask(tx, t.ID)
	***REMOVED***)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("global orchestrator: deleteTask failed to delete %s", t.ID)
	***REMOVED***
***REMOVED***

// IsRelatedService returns true if the service should be governed by this orchestrator
func (g *Orchestrator) IsRelatedService(service *api.Service) bool ***REMOVED***
	return orchestrator.IsGlobalService(service)
***REMOVED***

// SlotTuple returns a slot tuple for the global service task.
func (g *Orchestrator) SlotTuple(t *api.Task) orchestrator.SlotTuple ***REMOVED***
	return orchestrator.SlotTuple***REMOVED***
		ServiceID: t.ServiceID,
		NodeID:    t.NodeID,
	***REMOVED***
***REMOVED***

func isTaskCompleted(t *api.Task, restartPolicy api.RestartPolicy_RestartCondition) bool ***REMOVED***
	if t == nil || t.DesiredState <= api.TaskStateRunning ***REMOVED***
		return false
	***REMOVED***
	return restartPolicy == api.RestartOnNone ||
		(restartPolicy == api.RestartOnFailure && t.Status.State == api.TaskStateCompleted)
***REMOVED***
