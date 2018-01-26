package cluster

import (
	"fmt"

	"github.com/docker/docker/errdefs"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func getSwarm(ctx context.Context, c swarmapi.ControlClient) (*swarmapi.Cluster, error) ***REMOVED***
	rl, err := c.ListClusters(ctx, &swarmapi.ListClustersRequest***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(rl.Clusters) == 0 ***REMOVED***
		return nil, errors.WithStack(errNoSwarm)
	***REMOVED***

	// TODO: assume one cluster only
	return rl.Clusters[0], nil
***REMOVED***

func getNode(ctx context.Context, c swarmapi.ControlClient, input string) (*swarmapi.Node, error) ***REMOVED***
	// GetNode to match via full ID.
	if rg, err := c.GetNode(ctx, &swarmapi.GetNodeRequest***REMOVED***NodeID: input***REMOVED***); err == nil ***REMOVED***
		return rg.Node, nil
	***REMOVED***

	// If any error (including NotFound), ListNodes to match via full name.
	rl, err := c.ListNodes(ctx, &swarmapi.ListNodesRequest***REMOVED***
		Filters: &swarmapi.ListNodesRequest_Filters***REMOVED***
			Names: []string***REMOVED***input***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	if err != nil || len(rl.Nodes) == 0 ***REMOVED***
		// If any error or 0 result, ListNodes to match via ID prefix.
		rl, err = c.ListNodes(ctx, &swarmapi.ListNodesRequest***REMOVED***
			Filters: &swarmapi.ListNodesRequest_Filters***REMOVED***
				IDPrefixes: []string***REMOVED***input***REMOVED***,
			***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(rl.Nodes) == 0 ***REMOVED***
		err := fmt.Errorf("node %s not found", input)
		return nil, errdefs.NotFound(err)
	***REMOVED***

	if l := len(rl.Nodes); l > 1 ***REMOVED***
		return nil, errdefs.InvalidParameter(fmt.Errorf("node %s is ambiguous (%d matches found)", input, l))
	***REMOVED***

	return rl.Nodes[0], nil
***REMOVED***

func getService(ctx context.Context, c swarmapi.ControlClient, input string, insertDefaults bool) (*swarmapi.Service, error) ***REMOVED***
	// GetService to match via full ID.
	if rg, err := c.GetService(ctx, &swarmapi.GetServiceRequest***REMOVED***ServiceID: input, InsertDefaults: insertDefaults***REMOVED***); err == nil ***REMOVED***
		return rg.Service, nil
	***REMOVED***

	// If any error (including NotFound), ListServices to match via full name.
	rl, err := c.ListServices(ctx, &swarmapi.ListServicesRequest***REMOVED***
		Filters: &swarmapi.ListServicesRequest_Filters***REMOVED***
			Names: []string***REMOVED***input***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	if err != nil || len(rl.Services) == 0 ***REMOVED***
		// If any error or 0 result, ListServices to match via ID prefix.
		rl, err = c.ListServices(ctx, &swarmapi.ListServicesRequest***REMOVED***
			Filters: &swarmapi.ListServicesRequest_Filters***REMOVED***
				IDPrefixes: []string***REMOVED***input***REMOVED***,
			***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(rl.Services) == 0 ***REMOVED***
		err := fmt.Errorf("service %s not found", input)
		return nil, errdefs.NotFound(err)
	***REMOVED***

	if l := len(rl.Services); l > 1 ***REMOVED***
		return nil, errdefs.InvalidParameter(fmt.Errorf("service %s is ambiguous (%d matches found)", input, l))
	***REMOVED***

	if !insertDefaults ***REMOVED***
		return rl.Services[0], nil
	***REMOVED***

	rg, err := c.GetService(ctx, &swarmapi.GetServiceRequest***REMOVED***ServiceID: rl.Services[0].ID, InsertDefaults: true***REMOVED***)
	if err == nil ***REMOVED***
		return rg.Service, nil
	***REMOVED***
	return nil, err
***REMOVED***

func getTask(ctx context.Context, c swarmapi.ControlClient, input string) (*swarmapi.Task, error) ***REMOVED***
	// GetTask to match via full ID.
	if rg, err := c.GetTask(ctx, &swarmapi.GetTaskRequest***REMOVED***TaskID: input***REMOVED***); err == nil ***REMOVED***
		return rg.Task, nil
	***REMOVED***

	// If any error (including NotFound), ListTasks to match via full name.
	rl, err := c.ListTasks(ctx, &swarmapi.ListTasksRequest***REMOVED***
		Filters: &swarmapi.ListTasksRequest_Filters***REMOVED***
			Names: []string***REMOVED***input***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	if err != nil || len(rl.Tasks) == 0 ***REMOVED***
		// If any error or 0 result, ListTasks to match via ID prefix.
		rl, err = c.ListTasks(ctx, &swarmapi.ListTasksRequest***REMOVED***
			Filters: &swarmapi.ListTasksRequest_Filters***REMOVED***
				IDPrefixes: []string***REMOVED***input***REMOVED***,
			***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(rl.Tasks) == 0 ***REMOVED***
		err := fmt.Errorf("task %s not found", input)
		return nil, errdefs.NotFound(err)
	***REMOVED***

	if l := len(rl.Tasks); l > 1 ***REMOVED***
		return nil, errdefs.InvalidParameter(fmt.Errorf("task %s is ambiguous (%d matches found)", input, l))
	***REMOVED***

	return rl.Tasks[0], nil
***REMOVED***

func getSecret(ctx context.Context, c swarmapi.ControlClient, input string) (*swarmapi.Secret, error) ***REMOVED***
	// attempt to lookup secret by full ID
	if rg, err := c.GetSecret(ctx, &swarmapi.GetSecretRequest***REMOVED***SecretID: input***REMOVED***); err == nil ***REMOVED***
		return rg.Secret, nil
	***REMOVED***

	// If any error (including NotFound), ListSecrets to match via full name.
	rl, err := c.ListSecrets(ctx, &swarmapi.ListSecretsRequest***REMOVED***
		Filters: &swarmapi.ListSecretsRequest_Filters***REMOVED***
			Names: []string***REMOVED***input***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	if err != nil || len(rl.Secrets) == 0 ***REMOVED***
		// If any error or 0 result, ListSecrets to match via ID prefix.
		rl, err = c.ListSecrets(ctx, &swarmapi.ListSecretsRequest***REMOVED***
			Filters: &swarmapi.ListSecretsRequest_Filters***REMOVED***
				IDPrefixes: []string***REMOVED***input***REMOVED***,
			***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(rl.Secrets) == 0 ***REMOVED***
		err := fmt.Errorf("secret %s not found", input)
		return nil, errdefs.NotFound(err)
	***REMOVED***

	if l := len(rl.Secrets); l > 1 ***REMOVED***
		return nil, errdefs.InvalidParameter(fmt.Errorf("secret %s is ambiguous (%d matches found)", input, l))
	***REMOVED***

	return rl.Secrets[0], nil
***REMOVED***

func getConfig(ctx context.Context, c swarmapi.ControlClient, input string) (*swarmapi.Config, error) ***REMOVED***
	// attempt to lookup config by full ID
	if rg, err := c.GetConfig(ctx, &swarmapi.GetConfigRequest***REMOVED***ConfigID: input***REMOVED***); err == nil ***REMOVED***
		return rg.Config, nil
	***REMOVED***

	// If any error (including NotFound), ListConfigs to match via full name.
	rl, err := c.ListConfigs(ctx, &swarmapi.ListConfigsRequest***REMOVED***
		Filters: &swarmapi.ListConfigsRequest_Filters***REMOVED***
			Names: []string***REMOVED***input***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	if err != nil || len(rl.Configs) == 0 ***REMOVED***
		// If any error or 0 result, ListConfigs to match via ID prefix.
		rl, err = c.ListConfigs(ctx, &swarmapi.ListConfigsRequest***REMOVED***
			Filters: &swarmapi.ListConfigsRequest_Filters***REMOVED***
				IDPrefixes: []string***REMOVED***input***REMOVED***,
			***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(rl.Configs) == 0 ***REMOVED***
		err := fmt.Errorf("config %s not found", input)
		return nil, errdefs.NotFound(err)
	***REMOVED***

	if l := len(rl.Configs); l > 1 ***REMOVED***
		return nil, errdefs.InvalidParameter(fmt.Errorf("config %s is ambiguous (%d matches found)", input, l))
	***REMOVED***

	return rl.Configs[0], nil
***REMOVED***

func getNetwork(ctx context.Context, c swarmapi.ControlClient, input string) (*swarmapi.Network, error) ***REMOVED***
	// GetNetwork to match via full ID.
	if rg, err := c.GetNetwork(ctx, &swarmapi.GetNetworkRequest***REMOVED***NetworkID: input***REMOVED***); err == nil ***REMOVED***
		return rg.Network, nil
	***REMOVED***

	// If any error (including NotFound), ListNetworks to match via ID prefix and full name.
	rl, err := c.ListNetworks(ctx, &swarmapi.ListNetworksRequest***REMOVED***
		Filters: &swarmapi.ListNetworksRequest_Filters***REMOVED***
			Names: []string***REMOVED***input***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	if err != nil || len(rl.Networks) == 0 ***REMOVED***
		rl, err = c.ListNetworks(ctx, &swarmapi.ListNetworksRequest***REMOVED***
			Filters: &swarmapi.ListNetworksRequest_Filters***REMOVED***
				IDPrefixes: []string***REMOVED***input***REMOVED***,
			***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(rl.Networks) == 0 ***REMOVED***
		return nil, fmt.Errorf("network %s not found", input)
	***REMOVED***

	if l := len(rl.Networks); l > 1 ***REMOVED***
		return nil, errdefs.InvalidParameter(fmt.Errorf("network %s is ambiguous (%d matches found)", input, l))
	***REMOVED***

	return rl.Networks[0], nil
***REMOVED***
