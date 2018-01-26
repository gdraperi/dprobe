package manager

import (
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/raft"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
)

const roleReconcileInterval = 5 * time.Second

// roleManager reconciles the raft member list with desired role changes.
type roleManager struct ***REMOVED***
	ctx    context.Context
	cancel func()

	store    *store.MemoryStore
	raft     *raft.Node
	doneChan chan struct***REMOVED******REMOVED***

	// pending contains changed nodes that have not yet been reconciled in
	// the raft member list.
	pending map[string]*api.Node
***REMOVED***

// newRoleManager creates a new roleManager.
func newRoleManager(store *store.MemoryStore, raftNode *raft.Node) *roleManager ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	return &roleManager***REMOVED***
		ctx:      ctx,
		cancel:   cancel,
		store:    store,
		raft:     raftNode,
		doneChan: make(chan struct***REMOVED******REMOVED***),
		pending:  make(map[string]*api.Node),
	***REMOVED***
***REMOVED***

// Run is roleManager's main loop.
// ctx is only used for logging.
func (rm *roleManager) Run(ctx context.Context) ***REMOVED***
	defer close(rm.doneChan)

	var (
		nodes    []*api.Node
		ticker   *time.Ticker
		tickerCh <-chan time.Time
	)

	watcher, cancelWatch, err := store.ViewAndWatch(rm.store,
		func(readTx store.ReadTx) error ***REMOVED***
			var err error
			nodes, err = store.FindNodes(readTx, store.All)
			return err
		***REMOVED***,
		api.EventUpdateNode***REMOVED******REMOVED***)
	defer cancelWatch()

	if err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed to check nodes for role changes")
	***REMOVED*** else ***REMOVED***
		for _, node := range nodes ***REMOVED***
			rm.pending[node.ID] = node
			rm.reconcileRole(ctx, node)
		***REMOVED***
		if len(rm.pending) != 0 ***REMOVED***
			ticker = time.NewTicker(roleReconcileInterval)
			tickerCh = ticker.C
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case event := <-watcher:
			node := event.(api.EventUpdateNode).Node
			rm.pending[node.ID] = node
			rm.reconcileRole(ctx, node)
			if len(rm.pending) != 0 && ticker == nil ***REMOVED***
				ticker = time.NewTicker(roleReconcileInterval)
				tickerCh = ticker.C
			***REMOVED***
		case <-tickerCh:
			for _, node := range rm.pending ***REMOVED***
				rm.reconcileRole(ctx, node)
			***REMOVED***
			if len(rm.pending) == 0 ***REMOVED***
				ticker.Stop()
				ticker = nil
				tickerCh = nil
			***REMOVED***
		case <-rm.ctx.Done():
			if ticker != nil ***REMOVED***
				ticker.Stop()
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (rm *roleManager) reconcileRole(ctx context.Context, node *api.Node) ***REMOVED***
	if node.Role == node.Spec.DesiredRole ***REMOVED***
		// Nothing to do.
		delete(rm.pending, node.ID)
		return
	***REMOVED***

	// Promotion can proceed right away.
	if node.Spec.DesiredRole == api.NodeRoleManager && node.Role == api.NodeRoleWorker ***REMOVED***
		err := rm.store.Update(func(tx store.Tx) error ***REMOVED***
			updatedNode := store.GetNode(tx, node.ID)
			if updatedNode == nil || updatedNode.Spec.DesiredRole != node.Spec.DesiredRole || updatedNode.Role != node.Role ***REMOVED***
				return nil
			***REMOVED***
			updatedNode.Role = api.NodeRoleManager
			return store.UpdateNode(tx, updatedNode)
		***REMOVED***)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed to promote node %s", node.ID)
		***REMOVED*** else ***REMOVED***
			delete(rm.pending, node.ID)
		***REMOVED***
	***REMOVED*** else if node.Spec.DesiredRole == api.NodeRoleWorker && node.Role == api.NodeRoleManager ***REMOVED***
		// Check for node in memberlist
		member := rm.raft.GetMemberByNodeID(node.ID)
		if member != nil ***REMOVED***
			// Quorum safeguard
			if !rm.raft.CanRemoveMember(member.RaftID) ***REMOVED***
				// TODO(aaronl): Retry later
				log.G(ctx).Debugf("can't demote node %s at this time: removing member from raft would result in a loss of quorum", node.ID)
				return
			***REMOVED***

			rmCtx, rmCancel := context.WithTimeout(rm.ctx, 5*time.Second)
			defer rmCancel()

			if member.RaftID == rm.raft.Config.ID ***REMOVED***
				// Don't use rmCtx, because we expect to lose
				// leadership, which will cancel this context.
				log.G(ctx).Info("demoted; transferring leadership")
				err := rm.raft.TransferLeadership(context.Background())
				if err == nil ***REMOVED***
					return
				***REMOVED***
				log.G(ctx).WithError(err).Info("failed to transfer leadership")
			***REMOVED***
			if err := rm.raft.RemoveMember(rmCtx, member.RaftID); err != nil ***REMOVED***
				// TODO(aaronl): Retry later
				log.G(ctx).WithError(err).Debugf("can't demote node %s at this time", node.ID)
			***REMOVED***
			return
		***REMOVED***

		err := rm.store.Update(func(tx store.Tx) error ***REMOVED***
			updatedNode := store.GetNode(tx, node.ID)
			if updatedNode == nil || updatedNode.Spec.DesiredRole != node.Spec.DesiredRole || updatedNode.Role != node.Role ***REMOVED***
				return nil
			***REMOVED***
			updatedNode.Role = api.NodeRoleWorker

			return store.UpdateNode(tx, updatedNode)
		***REMOVED***)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Errorf("failed to demote node %s", node.ID)
		***REMOVED*** else ***REMOVED***
			delete(rm.pending, node.ID)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop stops the roleManager and waits for the main loop to exit.
func (rm *roleManager) Stop() ***REMOVED***
	rm.cancel()
	<-rm.doneChan
***REMOVED***
