package scheduler

import (
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/genericresource"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/protobuf/ptypes"
	"golang.org/x/net/context"
)

const (
	// monitorFailures is the lookback period for counting failures of
	// a task to determine if a node is faulty for a particular service.
	monitorFailures = 5 * time.Minute

	// maxFailures is the number of failures within monitorFailures that
	// triggers downweighting of a node in the sorting function.
	maxFailures = 5
)

type schedulingDecision struct ***REMOVED***
	old *api.Task
	new *api.Task
***REMOVED***

// Scheduler assigns tasks to nodes.
type Scheduler struct ***REMOVED***
	store           *store.MemoryStore
	unassignedTasks map[string]*api.Task
	// pendingPreassignedTasks already have NodeID, need resource validation
	pendingPreassignedTasks map[string]*api.Task
	// preassignedTasks tracks tasks that were preassigned, including those
	// past the pending state.
	preassignedTasks map[string]struct***REMOVED******REMOVED***
	nodeSet          nodeSet
	allTasks         map[string]*api.Task
	pipeline         *Pipeline

	// stopChan signals to the state machine to stop running
	stopChan chan struct***REMOVED******REMOVED***
	// doneChan is closed when the state machine terminates
	doneChan chan struct***REMOVED******REMOVED***
***REMOVED***

// New creates a new scheduler.
func New(store *store.MemoryStore) *Scheduler ***REMOVED***
	return &Scheduler***REMOVED***
		store:                   store,
		unassignedTasks:         make(map[string]*api.Task),
		pendingPreassignedTasks: make(map[string]*api.Task),
		preassignedTasks:        make(map[string]struct***REMOVED******REMOVED***),
		allTasks:                make(map[string]*api.Task),
		stopChan:                make(chan struct***REMOVED******REMOVED***),
		doneChan:                make(chan struct***REMOVED******REMOVED***),
		pipeline:                NewPipeline(),
	***REMOVED***
***REMOVED***

func (s *Scheduler) setupTasksList(tx store.ReadTx) error ***REMOVED***
	tasks, err := store.FindTasks(tx, store.All)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	tasksByNode := make(map[string]map[string]*api.Task)
	for _, t := range tasks ***REMOVED***
		// Ignore all tasks that have not reached PENDING
		// state and tasks that no longer consume resources.
		if t.Status.State < api.TaskStatePending || t.Status.State > api.TaskStateRunning ***REMOVED***
			continue
		***REMOVED***

		s.allTasks[t.ID] = t
		if t.NodeID == "" ***REMOVED***
			s.enqueue(t)
			continue
		***REMOVED***
		// preassigned tasks need to validate resource requirement on corresponding node
		if t.Status.State == api.TaskStatePending ***REMOVED***
			s.preassignedTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			s.pendingPreassignedTasks[t.ID] = t
			continue
		***REMOVED***

		if tasksByNode[t.NodeID] == nil ***REMOVED***
			tasksByNode[t.NodeID] = make(map[string]*api.Task)
		***REMOVED***
		tasksByNode[t.NodeID][t.ID] = t
	***REMOVED***

	return s.buildNodeSet(tx, tasksByNode)
***REMOVED***

// Run is the scheduler event loop.
func (s *Scheduler) Run(ctx context.Context) error ***REMOVED***
	defer close(s.doneChan)

	updates, cancel, err := store.ViewAndWatch(s.store, s.setupTasksList)
	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Errorf("snapshot store update failed")
		return err
	***REMOVED***
	defer cancel()

	// Validate resource for tasks from preassigned tasks
	// do this before other tasks because preassigned tasks like
	// global service should start before other tasks
	s.processPreassignedTasks(ctx)

	// Queue all unassigned tasks before processing changes.
	s.tick(ctx)

	const (
		// commitDebounceGap is the amount of time to wait between
		// commit events to debounce them.
		commitDebounceGap = 50 * time.Millisecond
		// maxLatency is a time limit on the debouncing.
		maxLatency = time.Second
	)
	var (
		debouncingStarted     time.Time
		commitDebounceTimer   *time.Timer
		commitDebounceTimeout <-chan time.Time
	)

	tickRequired := false

	schedule := func() ***REMOVED***
		if len(s.pendingPreassignedTasks) > 0 ***REMOVED***
			s.processPreassignedTasks(ctx)
		***REMOVED***
		if tickRequired ***REMOVED***
			s.tick(ctx)
			tickRequired = false
		***REMOVED***
	***REMOVED***

	// Watch for changes.
	for ***REMOVED***
		select ***REMOVED***
		case event := <-updates:
			switch v := event.(type) ***REMOVED***
			case api.EventCreateTask:
				if s.createTask(ctx, v.Task) ***REMOVED***
					tickRequired = true
				***REMOVED***
			case api.EventUpdateTask:
				if s.updateTask(ctx, v.Task) ***REMOVED***
					tickRequired = true
				***REMOVED***
			case api.EventDeleteTask:
				if s.deleteTask(v.Task) ***REMOVED***
					// deleting tasks may free up node resource, pending tasks should be re-evaluated.
					tickRequired = true
				***REMOVED***
			case api.EventCreateNode:
				s.createOrUpdateNode(v.Node)
				tickRequired = true
			case api.EventUpdateNode:
				s.createOrUpdateNode(v.Node)
				tickRequired = true
			case api.EventDeleteNode:
				s.nodeSet.remove(v.Node.ID)
			case state.EventCommit:
				if commitDebounceTimer != nil ***REMOVED***
					if time.Since(debouncingStarted) > maxLatency ***REMOVED***
						commitDebounceTimer.Stop()
						commitDebounceTimer = nil
						commitDebounceTimeout = nil
						schedule()
					***REMOVED*** else ***REMOVED***
						commitDebounceTimer.Reset(commitDebounceGap)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					commitDebounceTimer = time.NewTimer(commitDebounceGap)
					commitDebounceTimeout = commitDebounceTimer.C
					debouncingStarted = time.Now()
				***REMOVED***
			***REMOVED***
		case <-commitDebounceTimeout:
			schedule()
			commitDebounceTimer = nil
			commitDebounceTimeout = nil
		case <-s.stopChan:
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop causes the scheduler event loop to stop running.
func (s *Scheduler) Stop() ***REMOVED***
	close(s.stopChan)
	<-s.doneChan
***REMOVED***

// enqueue queues a task for scheduling.
func (s *Scheduler) enqueue(t *api.Task) ***REMOVED***
	s.unassignedTasks[t.ID] = t
***REMOVED***

func (s *Scheduler) createTask(ctx context.Context, t *api.Task) bool ***REMOVED***
	// Ignore all tasks that have not reached PENDING
	// state, and tasks that no longer consume resources.
	if t.Status.State < api.TaskStatePending || t.Status.State > api.TaskStateRunning ***REMOVED***
		return false
	***REMOVED***

	s.allTasks[t.ID] = t
	if t.NodeID == "" ***REMOVED***
		// unassigned task
		s.enqueue(t)
		return true
	***REMOVED***

	if t.Status.State == api.TaskStatePending ***REMOVED***
		s.preassignedTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		s.pendingPreassignedTasks[t.ID] = t
		// preassigned tasks do not contribute to running tasks count
		return false
	***REMOVED***

	nodeInfo, err := s.nodeSet.nodeInfo(t.NodeID)
	if err == nil && nodeInfo.addTask(t) ***REMOVED***
		s.nodeSet.updateNode(nodeInfo)
	***REMOVED***

	return false
***REMOVED***

func (s *Scheduler) updateTask(ctx context.Context, t *api.Task) bool ***REMOVED***
	// Ignore all tasks that have not reached PENDING
	// state.
	if t.Status.State < api.TaskStatePending ***REMOVED***
		return false
	***REMOVED***

	oldTask := s.allTasks[t.ID]

	// Ignore all tasks that have not reached Pending
	// state, and tasks that no longer consume resources.
	if t.Status.State > api.TaskStateRunning ***REMOVED***
		if oldTask == nil ***REMOVED***
			return false
		***REMOVED***

		if t.Status.State != oldTask.Status.State &&
			(t.Status.State == api.TaskStateFailed || t.Status.State == api.TaskStateRejected) ***REMOVED***
			// Keep track of task failures, so other nodes can be preferred
			// for scheduling this service if it looks like the service is
			// failing in a loop on this node. However, skip this for
			// preassigned tasks, because the scheduler does not choose
			// which nodes those run on.
			if _, wasPreassigned := s.preassignedTasks[t.ID]; !wasPreassigned ***REMOVED***
				nodeInfo, err := s.nodeSet.nodeInfo(t.NodeID)
				if err == nil ***REMOVED***
					nodeInfo.taskFailed(ctx, t)
					s.nodeSet.updateNode(nodeInfo)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		s.deleteTask(oldTask)

		return true
	***REMOVED***

	if t.NodeID == "" ***REMOVED***
		// unassigned task
		if oldTask != nil ***REMOVED***
			s.deleteTask(oldTask)
		***REMOVED***
		s.allTasks[t.ID] = t
		s.enqueue(t)
		return true
	***REMOVED***

	if t.Status.State == api.TaskStatePending ***REMOVED***
		if oldTask != nil ***REMOVED***
			s.deleteTask(oldTask)
		***REMOVED***
		s.preassignedTasks[t.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		s.allTasks[t.ID] = t
		s.pendingPreassignedTasks[t.ID] = t
		// preassigned tasks do not contribute to running tasks count
		return false
	***REMOVED***

	s.allTasks[t.ID] = t
	nodeInfo, err := s.nodeSet.nodeInfo(t.NodeID)
	if err == nil && nodeInfo.addTask(t) ***REMOVED***
		s.nodeSet.updateNode(nodeInfo)
	***REMOVED***

	return false
***REMOVED***

func (s *Scheduler) deleteTask(t *api.Task) bool ***REMOVED***
	delete(s.allTasks, t.ID)
	delete(s.preassignedTasks, t.ID)
	delete(s.pendingPreassignedTasks, t.ID)
	nodeInfo, err := s.nodeSet.nodeInfo(t.NodeID)
	if err == nil && nodeInfo.removeTask(t) ***REMOVED***
		s.nodeSet.updateNode(nodeInfo)
		return true
	***REMOVED***
	return false
***REMOVED***

func (s *Scheduler) createOrUpdateNode(n *api.Node) ***REMOVED***
	nodeInfo, nodeInfoErr := s.nodeSet.nodeInfo(n.ID)
	var resources *api.Resources
	if n.Description != nil && n.Description.Resources != nil ***REMOVED***
		resources = n.Description.Resources.Copy()
		// reconcile resources by looping over all tasks in this node
		if nodeInfoErr == nil ***REMOVED***
			for _, task := range nodeInfo.Tasks ***REMOVED***
				reservations := taskReservations(task.Spec)

				resources.MemoryBytes -= reservations.MemoryBytes
				resources.NanoCPUs -= reservations.NanoCPUs

				genericresource.ConsumeNodeResources(&resources.Generic,
					task.AssignedGenericResources)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		resources = &api.Resources***REMOVED******REMOVED***
	***REMOVED***

	if nodeInfoErr != nil ***REMOVED***
		nodeInfo = newNodeInfo(n, nil, *resources)
	***REMOVED*** else ***REMOVED***
		nodeInfo.Node = n
		nodeInfo.AvailableResources = resources
	***REMOVED***
	s.nodeSet.addOrUpdateNode(nodeInfo)
***REMOVED***

func (s *Scheduler) processPreassignedTasks(ctx context.Context) ***REMOVED***
	schedulingDecisions := make(map[string]schedulingDecision, len(s.pendingPreassignedTasks))
	for _, t := range s.pendingPreassignedTasks ***REMOVED***
		newT := s.taskFitNode(ctx, t, t.NodeID)
		if newT == nil ***REMOVED***
			continue
		***REMOVED***
		schedulingDecisions[t.ID] = schedulingDecision***REMOVED***old: t, new: newT***REMOVED***
	***REMOVED***

	successful, failed := s.applySchedulingDecisions(ctx, schedulingDecisions)

	for _, decision := range successful ***REMOVED***
		if decision.new.Status.State == api.TaskStateAssigned ***REMOVED***
			delete(s.pendingPreassignedTasks, decision.old.ID)
		***REMOVED***
	***REMOVED***
	for _, decision := range failed ***REMOVED***
		s.allTasks[decision.old.ID] = decision.old
		nodeInfo, err := s.nodeSet.nodeInfo(decision.new.NodeID)
		if err == nil && nodeInfo.removeTask(decision.new) ***REMOVED***
			s.nodeSet.updateNode(nodeInfo)
		***REMOVED***
	***REMOVED***
***REMOVED***

// tick attempts to schedule the queue.
func (s *Scheduler) tick(ctx context.Context) ***REMOVED***
	type commonSpecKey struct ***REMOVED***
		serviceID   string
		specVersion api.Version
	***REMOVED***
	tasksByCommonSpec := make(map[commonSpecKey]map[string]*api.Task)
	var oneOffTasks []*api.Task
	schedulingDecisions := make(map[string]schedulingDecision, len(s.unassignedTasks))

	for taskID, t := range s.unassignedTasks ***REMOVED***
		if t == nil || t.NodeID != "" ***REMOVED***
			// task deleted or already assigned
			delete(s.unassignedTasks, taskID)
			continue
		***REMOVED***

		// Group tasks with common specs
		if t.SpecVersion != nil ***REMOVED***
			taskGroupKey := commonSpecKey***REMOVED***
				serviceID:   t.ServiceID,
				specVersion: *t.SpecVersion,
			***REMOVED***

			if tasksByCommonSpec[taskGroupKey] == nil ***REMOVED***
				tasksByCommonSpec[taskGroupKey] = make(map[string]*api.Task)
			***REMOVED***
			tasksByCommonSpec[taskGroupKey][taskID] = t
		***REMOVED*** else ***REMOVED***
			// This task doesn't have a spec version. We have to
			// schedule it as a one-off.
			oneOffTasks = append(oneOffTasks, t)
		***REMOVED***
		delete(s.unassignedTasks, taskID)
	***REMOVED***

	for _, taskGroup := range tasksByCommonSpec ***REMOVED***
		s.scheduleTaskGroup(ctx, taskGroup, schedulingDecisions)
	***REMOVED***
	for _, t := range oneOffTasks ***REMOVED***
		s.scheduleTaskGroup(ctx, map[string]*api.Task***REMOVED***t.ID: t***REMOVED***, schedulingDecisions)
	***REMOVED***

	_, failed := s.applySchedulingDecisions(ctx, schedulingDecisions)
	for _, decision := range failed ***REMOVED***
		s.allTasks[decision.old.ID] = decision.old

		nodeInfo, err := s.nodeSet.nodeInfo(decision.new.NodeID)
		if err == nil && nodeInfo.removeTask(decision.new) ***REMOVED***
			s.nodeSet.updateNode(nodeInfo)
		***REMOVED***

		// enqueue task for next scheduling attempt
		s.enqueue(decision.old)
	***REMOVED***
***REMOVED***

func (s *Scheduler) applySchedulingDecisions(ctx context.Context, schedulingDecisions map[string]schedulingDecision) (successful, failed []schedulingDecision) ***REMOVED***
	if len(schedulingDecisions) == 0 ***REMOVED***
		return
	***REMOVED***

	successful = make([]schedulingDecision, 0, len(schedulingDecisions))

	// Apply changes to master store
	err := s.store.Batch(func(batch *store.Batch) error ***REMOVED***
		for len(schedulingDecisions) > 0 ***REMOVED***
			err := batch.Update(func(tx store.Tx) error ***REMOVED***
				// Update exactly one task inside this Update
				// callback.
				for taskID, decision := range schedulingDecisions ***REMOVED***
					delete(schedulingDecisions, taskID)

					t := store.GetTask(tx, taskID)
					if t == nil ***REMOVED***
						// Task no longer exists
						s.deleteTask(decision.new)
						continue
					***REMOVED***

					if t.Status.State == decision.new.Status.State &&
						t.Status.Message == decision.new.Status.Message &&
						t.Status.Err == decision.new.Status.Err ***REMOVED***
						// No changes, ignore
						continue
					***REMOVED***

					if t.Status.State >= api.TaskStateAssigned ***REMOVED***
						nodeInfo, err := s.nodeSet.nodeInfo(decision.new.NodeID)
						if err != nil ***REMOVED***
							failed = append(failed, decision)
							continue
						***REMOVED***
						node := store.GetNode(tx, decision.new.NodeID)
						if node == nil || node.Meta.Version != nodeInfo.Meta.Version ***REMOVED***
							// node is out of date
							failed = append(failed, decision)
							continue
						***REMOVED***
					***REMOVED***

					if err := store.UpdateTask(tx, decision.new); err != nil ***REMOVED***
						log.G(ctx).Debugf("scheduler failed to update task %s; will retry", taskID)
						failed = append(failed, decision)
						continue
					***REMOVED***
					successful = append(successful, decision)
					return nil
				***REMOVED***
				return nil
			***REMOVED***)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("scheduler tick transaction failed")
		failed = append(failed, successful...)
		successful = nil
	***REMOVED***
	return
***REMOVED***

// taskFitNode checks if a node has enough resources to accommodate a task.
func (s *Scheduler) taskFitNode(ctx context.Context, t *api.Task, nodeID string) *api.Task ***REMOVED***
	nodeInfo, err := s.nodeSet.nodeInfo(nodeID)
	if err != nil ***REMOVED***
		// node does not exist in set (it may have been deleted)
		return nil
	***REMOVED***
	newT := *t
	s.pipeline.SetTask(t)
	if !s.pipeline.Process(&nodeInfo) ***REMOVED***
		// this node cannot accommodate this task
		newT.Status.Timestamp = ptypes.MustTimestampProto(time.Now())
		newT.Status.Err = s.pipeline.Explain()
		s.allTasks[t.ID] = &newT

		return &newT
	***REMOVED***
	newT.Status = api.TaskStatus***REMOVED***
		State:     api.TaskStateAssigned,
		Timestamp: ptypes.MustTimestampProto(time.Now()),
		Message:   "scheduler confirmed task can run on preassigned node",
	***REMOVED***
	s.allTasks[t.ID] = &newT

	if nodeInfo.addTask(&newT) ***REMOVED***
		s.nodeSet.updateNode(nodeInfo)
	***REMOVED***
	return &newT
***REMOVED***

// scheduleTaskGroup schedules a batch of tasks that are part of the same
// service and share the same version of the spec.
func (s *Scheduler) scheduleTaskGroup(ctx context.Context, taskGroup map[string]*api.Task, schedulingDecisions map[string]schedulingDecision) ***REMOVED***
	// Pick at task at random from taskGroup to use for constraint
	// evaluation. It doesn't matter which one we pick because all the
	// tasks in the group are equal in terms of the fields the constraint
	// filters consider.
	var t *api.Task
	for _, t = range taskGroup ***REMOVED***
		break
	***REMOVED***

	s.pipeline.SetTask(t)

	now := time.Now()

	nodeLess := func(a *NodeInfo, b *NodeInfo) bool ***REMOVED***
		// If either node has at least maxFailures recent failures,
		// that's the deciding factor.
		recentFailuresA := a.countRecentFailures(now, t)
		recentFailuresB := b.countRecentFailures(now, t)

		if recentFailuresA >= maxFailures || recentFailuresB >= maxFailures ***REMOVED***
			if recentFailuresA > recentFailuresB ***REMOVED***
				return false
			***REMOVED***
			if recentFailuresB > recentFailuresA ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***

		tasksByServiceA := a.ActiveTasksCountByService[t.ServiceID]
		tasksByServiceB := b.ActiveTasksCountByService[t.ServiceID]

		if tasksByServiceA < tasksByServiceB ***REMOVED***
			return true
		***REMOVED***
		if tasksByServiceA > tasksByServiceB ***REMOVED***
			return false
		***REMOVED***

		// Total number of tasks breaks ties.
		return a.ActiveTasksCount < b.ActiveTasksCount
	***REMOVED***

	var prefs []*api.PlacementPreference
	if t.Spec.Placement != nil ***REMOVED***
		prefs = t.Spec.Placement.Preferences
	***REMOVED***

	tree := s.nodeSet.tree(t.ServiceID, prefs, len(taskGroup), s.pipeline.Process, nodeLess)

	s.scheduleNTasksOnSubtree(ctx, len(taskGroup), taskGroup, &tree, schedulingDecisions, nodeLess)
	if len(taskGroup) != 0 ***REMOVED***
		s.noSuitableNode(ctx, taskGroup, schedulingDecisions)
	***REMOVED***
***REMOVED***

func (s *Scheduler) scheduleNTasksOnSubtree(ctx context.Context, n int, taskGroup map[string]*api.Task, tree *decisionTree, schedulingDecisions map[string]schedulingDecision, nodeLess func(a *NodeInfo, b *NodeInfo) bool) int ***REMOVED***
	if tree.next == nil ***REMOVED***
		nodes := tree.orderedNodes(s.pipeline.Process, nodeLess)
		if len(nodes) == 0 ***REMOVED***
			return 0
		***REMOVED***

		return s.scheduleNTasksOnNodes(ctx, n, taskGroup, nodes, schedulingDecisions, nodeLess)
	***REMOVED***

	// Walk the tree and figure out how the tasks should be split at each
	// level.
	tasksScheduled := 0
	tasksInUsableBranches := tree.tasks
	var noRoom map[*decisionTree]struct***REMOVED******REMOVED***

	// Try to make branches even until either all branches are
	// full, or all tasks have been scheduled.
	for tasksScheduled != n && len(noRoom) != len(tree.next) ***REMOVED***
		desiredTasksPerBranch := (tasksInUsableBranches + n - tasksScheduled) / (len(tree.next) - len(noRoom))
		remainder := (tasksInUsableBranches + n - tasksScheduled) % (len(tree.next) - len(noRoom))

		for _, subtree := range tree.next ***REMOVED***
			if noRoom != nil ***REMOVED***
				if _, ok := noRoom[subtree]; ok ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			subtreeTasks := subtree.tasks
			if subtreeTasks < desiredTasksPerBranch || (subtreeTasks == desiredTasksPerBranch && remainder > 0) ***REMOVED***
				tasksToAssign := desiredTasksPerBranch - subtreeTasks
				if remainder > 0 ***REMOVED***
					tasksToAssign++
				***REMOVED***
				res := s.scheduleNTasksOnSubtree(ctx, tasksToAssign, taskGroup, subtree, schedulingDecisions, nodeLess)
				if res < tasksToAssign ***REMOVED***
					if noRoom == nil ***REMOVED***
						noRoom = make(map[*decisionTree]struct***REMOVED******REMOVED***)
					***REMOVED***
					noRoom[subtree] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					tasksInUsableBranches -= subtreeTasks
				***REMOVED*** else if remainder > 0 ***REMOVED***
					remainder--
				***REMOVED***
				tasksScheduled += res
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return tasksScheduled
***REMOVED***

func (s *Scheduler) scheduleNTasksOnNodes(ctx context.Context, n int, taskGroup map[string]*api.Task, nodes []NodeInfo, schedulingDecisions map[string]schedulingDecision, nodeLess func(a *NodeInfo, b *NodeInfo) bool) int ***REMOVED***
	tasksScheduled := 0
	failedConstraints := make(map[int]bool) // key is index in nodes slice
	nodeIter := 0
	nodeCount := len(nodes)
	for taskID, t := range taskGroup ***REMOVED***
		// Skip tasks which were already scheduled because they ended
		// up in two groups at once.
		if _, exists := schedulingDecisions[taskID]; exists ***REMOVED***
			continue
		***REMOVED***

		node := &nodes[nodeIter%nodeCount]

		log.G(ctx).WithField("task.id", t.ID).Debugf("assigning to node %s", node.ID)
		newT := *t
		newT.NodeID = node.ID
		newT.Status = api.TaskStatus***REMOVED***
			State:     api.TaskStateAssigned,
			Timestamp: ptypes.MustTimestampProto(time.Now()),
			Message:   "scheduler assigned task to node",
		***REMOVED***
		s.allTasks[t.ID] = &newT

		nodeInfo, err := s.nodeSet.nodeInfo(node.ID)
		if err == nil && nodeInfo.addTask(&newT) ***REMOVED***
			s.nodeSet.updateNode(nodeInfo)
			nodes[nodeIter%nodeCount] = nodeInfo
		***REMOVED***

		schedulingDecisions[taskID] = schedulingDecision***REMOVED***old: t, new: &newT***REMOVED***
		delete(taskGroup, taskID)
		tasksScheduled++
		if tasksScheduled == n ***REMOVED***
			return tasksScheduled
		***REMOVED***

		if nodeIter+1 < nodeCount ***REMOVED***
			// First pass fills the nodes until they have the same
			// number of tasks from this service.
			nextNode := nodes[(nodeIter+1)%nodeCount]
			if nodeLess(&nextNode, &nodeInfo) ***REMOVED***
				nodeIter++
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// In later passes, we just assign one task at a time
			// to each node that still meets the constraints.
			nodeIter++
		***REMOVED***

		origNodeIter := nodeIter
		for failedConstraints[nodeIter%nodeCount] || !s.pipeline.Process(&nodes[nodeIter%nodeCount]) ***REMOVED***
			failedConstraints[nodeIter%nodeCount] = true
			nodeIter++
			if nodeIter-origNodeIter == nodeCount ***REMOVED***
				// None of the nodes meet the constraints anymore.
				return tasksScheduled
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return tasksScheduled
***REMOVED***

func (s *Scheduler) noSuitableNode(ctx context.Context, taskGroup map[string]*api.Task, schedulingDecisions map[string]schedulingDecision) ***REMOVED***
	explanation := s.pipeline.Explain()
	for _, t := range taskGroup ***REMOVED***
		log.G(ctx).WithField("task.id", t.ID).Debug("no suitable node available for task")

		newT := *t
		newT.Status.Timestamp = ptypes.MustTimestampProto(time.Now())
		if explanation != "" ***REMOVED***
			newT.Status.Err = "no suitable node (" + explanation + ")"
		***REMOVED*** else ***REMOVED***
			newT.Status.Err = "no suitable node"
		***REMOVED***
		s.allTasks[t.ID] = &newT
		schedulingDecisions[t.ID] = schedulingDecision***REMOVED***old: t, new: &newT***REMOVED***

		s.enqueue(&newT)
	***REMOVED***
***REMOVED***

func (s *Scheduler) buildNodeSet(tx store.ReadTx, tasksByNode map[string]map[string]*api.Task) error ***REMOVED***
	nodes, err := store.FindNodes(tx, store.All)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	s.nodeSet.alloc(len(nodes))

	for _, n := range nodes ***REMOVED***
		var resources api.Resources
		if n.Description != nil && n.Description.Resources != nil ***REMOVED***
			resources = *n.Description.Resources
		***REMOVED***
		s.nodeSet.addOrUpdateNode(newNodeInfo(n, tasksByNode[n.ID], resources))
	***REMOVED***

	return nil
***REMOVED***
