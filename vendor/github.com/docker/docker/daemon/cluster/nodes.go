package cluster

import (
	apitypes "github.com/docker/docker/api/types"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	"github.com/docker/docker/errdefs"
	swarmapi "github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

// GetNodes returns a list of all nodes known to a cluster.
func (c *Cluster) GetNodes(options apitypes.NodeListOptions) ([]types.Node, error) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() ***REMOVED***
		return nil, c.errNoManager(state)
	***REMOVED***

	filters, err := newListNodesFilters(options.Filters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ctx, cancel := c.getRequestContext()
	defer cancel()

	r, err := state.controlClient.ListNodes(
		ctx,
		&swarmapi.ListNodesRequest***REMOVED***Filters: filters***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	nodes := make([]types.Node, 0, len(r.Nodes))

	for _, node := range r.Nodes ***REMOVED***
		nodes = append(nodes, convert.NodeFromGRPC(*node))
	***REMOVED***
	return nodes, nil
***REMOVED***

// GetNode returns a node based on an ID.
func (c *Cluster) GetNode(input string) (types.Node, error) ***REMOVED***
	var node *swarmapi.Node

	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		n, err := getNode(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		node = n
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return types.Node***REMOVED******REMOVED***, err
	***REMOVED***

	return convert.NodeFromGRPC(*node), nil
***REMOVED***

// UpdateNode updates existing nodes properties.
func (c *Cluster) UpdateNode(input string, version uint64, spec types.NodeSpec) error ***REMOVED***
	return c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		nodeSpec, err := convert.NodeSpecToGRPC(spec)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***

		ctx, cancel := c.getRequestContext()
		defer cancel()

		currentNode, err := getNode(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		_, err = state.controlClient.UpdateNode(
			ctx,
			&swarmapi.UpdateNodeRequest***REMOVED***
				NodeID: currentNode.ID,
				Spec:   &nodeSpec,
				NodeVersion: &swarmapi.Version***REMOVED***
					Index: version,
				***REMOVED***,
			***REMOVED***,
		)
		return err
	***REMOVED***)
***REMOVED***

// RemoveNode removes a node from a cluster
func (c *Cluster) RemoveNode(input string, force bool) error ***REMOVED***
	return c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		node, err := getNode(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		_, err = state.controlClient.RemoveNode(ctx, &swarmapi.RemoveNodeRequest***REMOVED***NodeID: node.ID, Force: force***REMOVED***)
		return err
	***REMOVED***)
***REMOVED***
