package cluster

import (
	apitypes "github.com/docker/docker/api/types"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	swarmapi "github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

// GetConfig returns a config from a managed swarm cluster
func (c *Cluster) GetConfig(input string) (types.Config, error) ***REMOVED***
	var config *swarmapi.Config

	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		s, err := getConfig(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		config = s
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return types.Config***REMOVED******REMOVED***, err
	***REMOVED***
	return convert.ConfigFromGRPC(config), nil
***REMOVED***

// GetConfigs returns all configs of a managed swarm cluster.
func (c *Cluster) GetConfigs(options apitypes.ConfigListOptions) ([]types.Config, error) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() ***REMOVED***
		return nil, c.errNoManager(state)
	***REMOVED***

	filters, err := newListConfigsFilters(options.Filters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ctx, cancel := c.getRequestContext()
	defer cancel()

	r, err := state.controlClient.ListConfigs(ctx,
		&swarmapi.ListConfigsRequest***REMOVED***Filters: filters***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	configs := []types.Config***REMOVED******REMOVED***

	for _, config := range r.Configs ***REMOVED***
		configs = append(configs, convert.ConfigFromGRPC(config))
	***REMOVED***

	return configs, nil
***REMOVED***

// CreateConfig creates a new config in a managed swarm cluster.
func (c *Cluster) CreateConfig(s types.ConfigSpec) (string, error) ***REMOVED***
	var resp *swarmapi.CreateConfigResponse
	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		configSpec := convert.ConfigSpecToGRPC(s)

		r, err := state.controlClient.CreateConfig(ctx,
			&swarmapi.CreateConfigRequest***REMOVED***Spec: &configSpec***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		resp = r
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return resp.Config.ID, nil
***REMOVED***

// RemoveConfig removes a config from a managed swarm cluster.
func (c *Cluster) RemoveConfig(input string) error ***REMOVED***
	return c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		config, err := getConfig(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		req := &swarmapi.RemoveConfigRequest***REMOVED***
			ConfigID: config.ID,
		***REMOVED***

		_, err = state.controlClient.RemoveConfig(ctx, req)
		return err
	***REMOVED***)
***REMOVED***

// UpdateConfig updates a config in a managed swarm cluster.
// Note: this is not exposed to the CLI but is available from the API only
func (c *Cluster) UpdateConfig(input string, version uint64, spec types.ConfigSpec) error ***REMOVED***
	return c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		config, err := getConfig(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		configSpec := convert.ConfigSpecToGRPC(spec)

		_, err = state.controlClient.UpdateConfig(ctx,
			&swarmapi.UpdateConfigRequest***REMOVED***
				ConfigID: config.ID,
				ConfigVersion: &swarmapi.Version***REMOVED***
					Index: version,
				***REMOVED***,
				Spec: &configSpec,
			***REMOVED***)
		return err
	***REMOVED***)
***REMOVED***
