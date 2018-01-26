package cluster

import (
	"fmt"

	apitypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/runconfig"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// GetNetworks returns all current cluster managed networks.
func (c *Cluster) GetNetworks() ([]apitypes.NetworkResource, error) ***REMOVED***
	list, err := c.getNetworks(nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	removePredefinedNetworks(&list)
	return list, nil
***REMOVED***

func removePredefinedNetworks(networks *[]apitypes.NetworkResource) ***REMOVED***
	if networks == nil ***REMOVED***
		return
	***REMOVED***
	var idxs []int
	for i, n := range *networks ***REMOVED***
		if v, ok := n.Labels["com.docker.swarm.predefined"]; ok && v == "true" ***REMOVED***
			idxs = append(idxs, i)
		***REMOVED***
	***REMOVED***
	for i, idx := range idxs ***REMOVED***
		idx -= i
		*networks = append((*networks)[:idx], (*networks)[idx+1:]...)
	***REMOVED***
***REMOVED***

func (c *Cluster) getNetworks(filters *swarmapi.ListNetworksRequest_Filters) ([]apitypes.NetworkResource, error) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() ***REMOVED***
		return nil, c.errNoManager(state)
	***REMOVED***

	ctx, cancel := c.getRequestContext()
	defer cancel()

	r, err := state.controlClient.ListNetworks(ctx, &swarmapi.ListNetworksRequest***REMOVED***Filters: filters***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	networks := make([]apitypes.NetworkResource, 0, len(r.Networks))

	for _, network := range r.Networks ***REMOVED***
		networks = append(networks, convert.BasicNetworkFromGRPC(*network))
	***REMOVED***

	return networks, nil
***REMOVED***

// GetNetwork returns a cluster network by an ID.
func (c *Cluster) GetNetwork(input string) (apitypes.NetworkResource, error) ***REMOVED***
	var network *swarmapi.Network

	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		n, err := getNetwork(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		network = n
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return apitypes.NetworkResource***REMOVED******REMOVED***, err
	***REMOVED***
	return convert.BasicNetworkFromGRPC(*network), nil
***REMOVED***

// GetNetworksByName returns cluster managed networks by name.
// It is ok to have multiple networks here. #18864
func (c *Cluster) GetNetworksByName(name string) ([]apitypes.NetworkResource, error) ***REMOVED***
	// Note that swarmapi.GetNetworkRequest.Name is not functional.
	// So we cannot just use that with c.GetNetwork.
	return c.getNetworks(&swarmapi.ListNetworksRequest_Filters***REMOVED***
		Names: []string***REMOVED***name***REMOVED***,
	***REMOVED***)
***REMOVED***

func attacherKey(target, containerID string) string ***REMOVED***
	return containerID + ":" + target
***REMOVED***

// UpdateAttachment signals the attachment config to the attachment
// waiter who is trying to start or attach the container to the
// network.
func (c *Cluster) UpdateAttachment(target, containerID string, config *network.NetworkingConfig) error ***REMOVED***
	c.mu.Lock()
	attacher, ok := c.attachers[attacherKey(target, containerID)]
	if !ok || attacher == nil ***REMOVED***
		c.mu.Unlock()
		return fmt.Errorf("could not find attacher for container %s to network %s", containerID, target)
	***REMOVED***
	if attacher.inProgress ***REMOVED***
		logrus.Debugf("Discarding redundant notice of resource allocation on network %s for task id %s", target, attacher.taskID)
		c.mu.Unlock()
		return nil
	***REMOVED***
	attacher.inProgress = true
	c.mu.Unlock()

	attacher.attachWaitCh <- config

	return nil
***REMOVED***

// WaitForDetachment waits for the container to stop or detach from
// the network.
func (c *Cluster) WaitForDetachment(ctx context.Context, networkName, networkID, taskID, containerID string) error ***REMOVED***
	c.mu.RLock()
	attacher, ok := c.attachers[attacherKey(networkName, containerID)]
	if !ok ***REMOVED***
		attacher, ok = c.attachers[attacherKey(networkID, containerID)]
	***REMOVED***
	state := c.currentNodeState()
	if state.swarmNode == nil || state.swarmNode.Agent() == nil ***REMOVED***
		c.mu.RUnlock()
		return errors.New("invalid cluster node while waiting for detachment")
	***REMOVED***

	c.mu.RUnlock()
	agent := state.swarmNode.Agent()
	if ok && attacher != nil &&
		attacher.detachWaitCh != nil &&
		attacher.attachCompleteCh != nil ***REMOVED***
		// Attachment may be in progress still so wait for
		// attachment to complete.
		select ***REMOVED***
		case <-attacher.attachCompleteCh:
		case <-ctx.Done():
			return ctx.Err()
		***REMOVED***

		if attacher.taskID == taskID ***REMOVED***
			select ***REMOVED***
			case <-attacher.detachWaitCh:
			case <-ctx.Done():
				return ctx.Err()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return agent.ResourceAllocator().DetachNetwork(ctx, taskID)
***REMOVED***

// AttachNetwork generates an attachment request towards the manager.
func (c *Cluster) AttachNetwork(target string, containerID string, addresses []string) (*network.NetworkingConfig, error) ***REMOVED***
	aKey := attacherKey(target, containerID)
	c.mu.Lock()
	state := c.currentNodeState()
	if state.swarmNode == nil || state.swarmNode.Agent() == nil ***REMOVED***
		c.mu.Unlock()
		return nil, errors.New("invalid cluster node while attaching to network")
	***REMOVED***
	if attacher, ok := c.attachers[aKey]; ok ***REMOVED***
		c.mu.Unlock()
		return attacher.config, nil
	***REMOVED***

	agent := state.swarmNode.Agent()
	attachWaitCh := make(chan *network.NetworkingConfig)
	detachWaitCh := make(chan struct***REMOVED******REMOVED***)
	attachCompleteCh := make(chan struct***REMOVED******REMOVED***)
	c.attachers[aKey] = &attacher***REMOVED***
		attachWaitCh:     attachWaitCh,
		attachCompleteCh: attachCompleteCh,
		detachWaitCh:     detachWaitCh,
	***REMOVED***
	c.mu.Unlock()

	ctx, cancel := c.getRequestContext()
	defer cancel()

	taskID, err := agent.ResourceAllocator().AttachNetwork(ctx, containerID, target, addresses)
	if err != nil ***REMOVED***
		c.mu.Lock()
		delete(c.attachers, aKey)
		c.mu.Unlock()
		return nil, fmt.Errorf("Could not attach to network %s: %v", target, err)
	***REMOVED***

	c.mu.Lock()
	c.attachers[aKey].taskID = taskID
	close(attachCompleteCh)
	c.mu.Unlock()

	logrus.Debugf("Successfully attached to network %s with task id %s", target, taskID)

	release := func() ***REMOVED***
		ctx, cancel := c.getRequestContext()
		defer cancel()
		if err := agent.ResourceAllocator().DetachNetwork(ctx, taskID); err != nil ***REMOVED***
			logrus.Errorf("Failed remove network attachment %s to network %s on allocation failure: %v",
				taskID, target, err)
		***REMOVED***
	***REMOVED***

	var config *network.NetworkingConfig
	select ***REMOVED***
	case config = <-attachWaitCh:
	case <-ctx.Done():
		release()
		return nil, fmt.Errorf("attaching to network failed, make sure your network options are correct and check manager logs: %v", ctx.Err())
	***REMOVED***

	c.mu.Lock()
	c.attachers[aKey].config = config
	c.mu.Unlock()

	logrus.Debugf("Successfully allocated resources on network %s for task id %s", target, taskID)

	return config, nil
***REMOVED***

// DetachNetwork unblocks the waiters waiting on WaitForDetachment so
// that a request to detach can be generated towards the manager.
func (c *Cluster) DetachNetwork(target string, containerID string) error ***REMOVED***
	aKey := attacherKey(target, containerID)

	c.mu.Lock()
	attacher, ok := c.attachers[aKey]
	delete(c.attachers, aKey)
	c.mu.Unlock()

	if !ok ***REMOVED***
		return fmt.Errorf("could not find network attachment for container %s to network %s", containerID, target)
	***REMOVED***

	close(attacher.detachWaitCh)
	return nil
***REMOVED***

// CreateNetwork creates a new cluster managed network.
func (c *Cluster) CreateNetwork(s apitypes.NetworkCreateRequest) (string, error) ***REMOVED***
	if runconfig.IsPreDefinedNetwork(s.Name) ***REMOVED***
		err := notAllowedError(fmt.Sprintf("%s is a pre-defined network and cannot be created", s.Name))
		return "", errors.WithStack(err)
	***REMOVED***

	var resp *swarmapi.CreateNetworkResponse
	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		networkSpec := convert.BasicNetworkCreateToGRPC(s)
		r, err := state.controlClient.CreateNetwork(ctx, &swarmapi.CreateNetworkRequest***REMOVED***Spec: &networkSpec***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		resp = r
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return resp.Network.ID, nil
***REMOVED***

// RemoveNetwork removes a cluster network.
func (c *Cluster) RemoveNetwork(input string) error ***REMOVED***
	return c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		network, err := getNetwork(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		_, err = state.controlClient.RemoveNetwork(ctx, &swarmapi.RemoveNetworkRequest***REMOVED***NetworkID: network.ID***REMOVED***)
		return err
	***REMOVED***)
***REMOVED***

func (c *Cluster) populateNetworkID(ctx context.Context, client swarmapi.ControlClient, s *types.ServiceSpec) error ***REMOVED***
	// Always prefer NetworkAttachmentConfigs from TaskTemplate
	// but fallback to service spec for backward compatibility
	networks := s.TaskTemplate.Networks
	if len(networks) == 0 ***REMOVED***
		networks = s.Networks
	***REMOVED***
	for i, n := range networks ***REMOVED***
		apiNetwork, err := getNetwork(ctx, client, n.Target)
		if err != nil ***REMOVED***
			ln, _ := c.config.Backend.FindNetwork(n.Target)
			if ln != nil && runconfig.IsPreDefinedNetwork(ln.Name()) ***REMOVED***
				// Need to retrieve the corresponding predefined swarm network
				// and use its id for the request.
				apiNetwork, err = getNetwork(ctx, client, ln.Name())
				if err != nil ***REMOVED***
					return errors.Wrap(errdefs.NotFound(err), "could not find the corresponding predefined swarm network")
				***REMOVED***
				goto setid
			***REMOVED***
			if ln != nil && !ln.Info().Dynamic() ***REMOVED***
				errMsg := fmt.Sprintf("The network %s cannot be used with services. Only networks scoped to the swarm can be used, such as those created with the overlay driver.", ln.Name())
				return errors.WithStack(notAllowedError(errMsg))
			***REMOVED***
			return err
		***REMOVED***
	setid:
		networks[i].Target = apiNetwork.ID
	***REMOVED***
	return nil
***REMOVED***
