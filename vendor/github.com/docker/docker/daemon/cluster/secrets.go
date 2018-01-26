package cluster

import (
	apitypes "github.com/docker/docker/api/types"
	types "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	swarmapi "github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

// GetSecret returns a secret from a managed swarm cluster
func (c *Cluster) GetSecret(input string) (types.Secret, error) ***REMOVED***
	var secret *swarmapi.Secret

	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		s, err := getSecret(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		secret = s
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return types.Secret***REMOVED******REMOVED***, err
	***REMOVED***
	return convert.SecretFromGRPC(secret), nil
***REMOVED***

// GetSecrets returns all secrets of a managed swarm cluster.
func (c *Cluster) GetSecrets(options apitypes.SecretListOptions) ([]types.Secret, error) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := c.currentNodeState()
	if !state.IsActiveManager() ***REMOVED***
		return nil, c.errNoManager(state)
	***REMOVED***

	filters, err := newListSecretsFilters(options.Filters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ctx, cancel := c.getRequestContext()
	defer cancel()

	r, err := state.controlClient.ListSecrets(ctx,
		&swarmapi.ListSecretsRequest***REMOVED***Filters: filters***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	secrets := make([]types.Secret, 0, len(r.Secrets))

	for _, secret := range r.Secrets ***REMOVED***
		secrets = append(secrets, convert.SecretFromGRPC(secret))
	***REMOVED***

	return secrets, nil
***REMOVED***

// CreateSecret creates a new secret in a managed swarm cluster.
func (c *Cluster) CreateSecret(s types.SecretSpec) (string, error) ***REMOVED***
	var resp *swarmapi.CreateSecretResponse
	if err := c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		secretSpec := convert.SecretSpecToGRPC(s)

		r, err := state.controlClient.CreateSecret(ctx,
			&swarmapi.CreateSecretRequest***REMOVED***Spec: &secretSpec***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		resp = r
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return resp.Secret.ID, nil
***REMOVED***

// RemoveSecret removes a secret from a managed swarm cluster.
func (c *Cluster) RemoveSecret(input string) error ***REMOVED***
	return c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		secret, err := getSecret(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		req := &swarmapi.RemoveSecretRequest***REMOVED***
			SecretID: secret.ID,
		***REMOVED***

		_, err = state.controlClient.RemoveSecret(ctx, req)
		return err
	***REMOVED***)
***REMOVED***

// UpdateSecret updates a secret in a managed swarm cluster.
// Note: this is not exposed to the CLI but is available from the API only
func (c *Cluster) UpdateSecret(input string, version uint64, spec types.SecretSpec) error ***REMOVED***
	return c.lockedManagerAction(func(ctx context.Context, state nodeState) error ***REMOVED***
		secret, err := getSecret(ctx, state.controlClient, input)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		secretSpec := convert.SecretSpecToGRPC(spec)

		_, err = state.controlClient.UpdateSecret(ctx,
			&swarmapi.UpdateSecretRequest***REMOVED***
				SecretID: secret.ID,
				SecretVersion: &swarmapi.Version***REMOVED***
					Index: version,
				***REMOVED***,
				Spec: &secretSpec,
			***REMOVED***)
		return err
	***REMOVED***)
***REMOVED***
