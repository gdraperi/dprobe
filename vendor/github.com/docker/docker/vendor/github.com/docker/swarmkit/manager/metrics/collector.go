package metrics

import (
	"context"

	"strings"

	metrics "github.com/docker/go-metrics"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/state/store"
)

var (
	ns          = metrics.NewNamespace("swarm", "manager", nil)
	nodesMetric metrics.LabeledGauge
)

func init() ***REMOVED***
	nodesMetric = ns.NewLabeledGauge("nodes", "The number of nodes", "", "state")
	for _, state := range api.NodeStatus_State_name ***REMOVED***
		nodesMetric.WithValues(strings.ToLower(state)).Set(0)
	***REMOVED***
	metrics.Register(ns)
***REMOVED***

// Collector collects swarmkit metrics
type Collector struct ***REMOVED***
	store *store.MemoryStore

	// stopChan signals to the state machine to stop running.
	stopChan chan struct***REMOVED******REMOVED***
	// doneChan is closed when the state machine terminates.
	doneChan chan struct***REMOVED******REMOVED***
***REMOVED***

// NewCollector creates a new metrics collector
func NewCollector(store *store.MemoryStore) *Collector ***REMOVED***
	return &Collector***REMOVED***
		store:    store,
		stopChan: make(chan struct***REMOVED******REMOVED***),
		doneChan: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

func (c *Collector) updateNodeState(prevNode, newNode *api.Node) ***REMOVED***
	// Skip updates if nothing changed.
	if prevNode != nil && newNode != nil && prevNode.Status.State == newNode.Status.State ***REMOVED***
		return
	***REMOVED***

	if prevNode != nil ***REMOVED***
		nodesMetric.WithValues(strings.ToLower(prevNode.Status.State.String())).Dec(1)
	***REMOVED***
	if newNode != nil ***REMOVED***
		nodesMetric.WithValues(strings.ToLower(newNode.Status.State.String())).Inc(1)
	***REMOVED***
***REMOVED***

// Run contains the collector event loop
func (c *Collector) Run(ctx context.Context) error ***REMOVED***
	defer close(c.doneChan)

	watcher, cancel, err := store.ViewAndWatch(c.store, func(readTx store.ReadTx) error ***REMOVED***
		nodes, err := store.FindNodes(readTx, store.All)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, node := range nodes ***REMOVED***
			c.updateNodeState(nil, node)
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer cancel()

	for ***REMOVED***
		select ***REMOVED***
		case event := <-watcher:
			switch v := event.(type) ***REMOVED***
			case api.EventCreateNode:
				c.updateNodeState(nil, v.Node)
			case api.EventUpdateNode:
				c.updateNodeState(v.OldNode, v.Node)
			case api.EventDeleteNode:
				c.updateNodeState(v.Node, nil)
			***REMOVED***
		case <-c.stopChan:
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop stops the collector.
func (c *Collector) Stop() ***REMOVED***
	close(c.stopChan)
	<-c.doneChan

	// Clean the metrics on exit.
	for _, state := range api.NodeStatus_State_name ***REMOVED***
		nodesMetric.WithValues(strings.ToLower(state)).Set(0)
	***REMOVED***
***REMOVED***
