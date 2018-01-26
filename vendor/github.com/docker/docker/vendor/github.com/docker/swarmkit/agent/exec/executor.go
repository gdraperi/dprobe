package exec

import (
	"github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

// Executor provides controllers for tasks.
type Executor interface ***REMOVED***
	// Describe returns the underlying node description.
	Describe(ctx context.Context) (*api.NodeDescription, error)

	// Configure uses the node object state to propagate node
	// state to the underlying executor.
	Configure(ctx context.Context, node *api.Node) error

	// Controller provides a controller for the given task.
	Controller(t *api.Task) (Controller, error)

	// SetNetworkBootstrapKeys passes the symmetric keys from the
	// manager to the executor.
	SetNetworkBootstrapKeys([]*api.EncryptionKey) error
***REMOVED***

// SecretsProvider is implemented by objects that can store secrets, typically
// an executor.
type SecretsProvider interface ***REMOVED***
	Secrets() SecretsManager
***REMOVED***

// ConfigsProvider is implemented by objects that can store configs,
// typically an executor.
type ConfigsProvider interface ***REMOVED***
	Configs() ConfigsManager
***REMOVED***

// DependencyManager is a meta-object that can keep track of typed objects
// such as secrets and configs.
type DependencyManager interface ***REMOVED***
	SecretsProvider
	ConfigsProvider
***REMOVED***

// DependencyGetter is a meta-object that can provide access to typed objects
// such as secrets and configs.
type DependencyGetter interface ***REMOVED***
	Secrets() SecretGetter
	Configs() ConfigGetter
***REMOVED***

// SecretGetter contains secret data necessary for the Controller.
type SecretGetter interface ***REMOVED***
	// Get returns the the secret with a specific secret ID, if available.
	// When the secret is not available, the return will be nil.
	Get(secretID string) (*api.Secret, error)
***REMOVED***

// SecretsManager is the interface for secret storage and updates.
type SecretsManager interface ***REMOVED***
	SecretGetter

	Add(secrets ...api.Secret) // add one or more secrets
	Remove(secrets []string)   // remove the secrets by ID
	Reset()                    // remove all secrets
***REMOVED***

// ConfigGetter contains config data necessary for the Controller.
type ConfigGetter interface ***REMOVED***
	// Get returns the the config with a specific config ID, if available.
	// When the config is not available, the return will be nil.
	Get(configID string) (*api.Config, error)
***REMOVED***

// ConfigsManager is the interface for config storage and updates.
type ConfigsManager interface ***REMOVED***
	ConfigGetter

	Add(configs ...api.Config) // add one or more configs
	Remove(configs []string)   // remove the configs by ID
	Reset()                    // remove all configs
***REMOVED***
