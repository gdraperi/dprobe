package template

import (
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
	"github.com/pkg/errors"
)

type templatedSecretGetter struct ***REMOVED***
	dependencies exec.DependencyGetter
	t            *api.Task
	node         *api.NodeDescription
***REMOVED***

// NewTemplatedSecretGetter returns a SecretGetter that evaluates templates.
func NewTemplatedSecretGetter(dependencies exec.DependencyGetter, t *api.Task, node *api.NodeDescription) exec.SecretGetter ***REMOVED***
	return templatedSecretGetter***REMOVED***dependencies: dependencies, t: t, node: node***REMOVED***
***REMOVED***

func (t templatedSecretGetter) Get(secretID string) (*api.Secret, error) ***REMOVED***
	if t.dependencies == nil ***REMOVED***
		return nil, errors.New("no secret provider available")
	***REMOVED***

	secrets := t.dependencies.Secrets()
	if secrets == nil ***REMOVED***
		return nil, errors.New("no secret provider available")
	***REMOVED***

	secret, err := secrets.Get(secretID)
	if err != nil ***REMOVED***
		return secret, err
	***REMOVED***

	newSpec, err := ExpandSecretSpec(secret, t.node, t.t, t.dependencies)
	if err != nil ***REMOVED***
		return secret, errors.Wrapf(err, "failed to expand templated secret %s", secretID)
	***REMOVED***

	secretCopy := *secret
	secretCopy.Spec = *newSpec
	return &secretCopy, nil
***REMOVED***

// TemplatedConfigGetter is a ConfigGetter with an additional method to expose
// whether a config contains sensitive data.
type TemplatedConfigGetter interface ***REMOVED***
	exec.ConfigGetter

	// GetAndFlagSecretData returns the interpolated config, and also
	// returns true if the config has been interpolated with data from a
	// secret. In this case, the config should be handled specially and
	// should not be written to disk.
	GetAndFlagSecretData(configID string) (*api.Config, bool, error)
***REMOVED***

type templatedConfigGetter struct ***REMOVED***
	dependencies exec.DependencyGetter
	t            *api.Task
	node         *api.NodeDescription
***REMOVED***

// NewTemplatedConfigGetter returns a ConfigGetter that evaluates templates.
func NewTemplatedConfigGetter(dependencies exec.DependencyGetter, t *api.Task, node *api.NodeDescription) TemplatedConfigGetter ***REMOVED***
	return templatedConfigGetter***REMOVED***dependencies: dependencies, t: t, node: node***REMOVED***
***REMOVED***

func (t templatedConfigGetter) Get(configID string) (*api.Config, error) ***REMOVED***
	config, _, err := t.GetAndFlagSecretData(configID)
	return config, err
***REMOVED***

func (t templatedConfigGetter) GetAndFlagSecretData(configID string) (*api.Config, bool, error) ***REMOVED***
	if t.dependencies == nil ***REMOVED***
		return nil, false, errors.New("no config provider available")
	***REMOVED***

	configs := t.dependencies.Configs()
	if configs == nil ***REMOVED***
		return nil, false, errors.New("no config provider available")
	***REMOVED***

	config, err := configs.Get(configID)
	if err != nil ***REMOVED***
		return config, false, err
	***REMOVED***

	newSpec, sensitive, err := ExpandConfigSpec(config, t.node, t.t, t.dependencies)
	if err != nil ***REMOVED***
		return config, false, errors.Wrapf(err, "failed to expand templated config %s", configID)
	***REMOVED***

	configCopy := *config
	configCopy.Spec = *newSpec
	return &configCopy, sensitive, nil
***REMOVED***

type templatedDependencyGetter struct ***REMOVED***
	secrets exec.SecretGetter
	configs TemplatedConfigGetter
***REMOVED***

// NewTemplatedDependencyGetter returns a DependencyGetter that evaluates templates.
func NewTemplatedDependencyGetter(dependencies exec.DependencyGetter, t *api.Task, node *api.NodeDescription) exec.DependencyGetter ***REMOVED***
	return templatedDependencyGetter***REMOVED***
		secrets: NewTemplatedSecretGetter(dependencies, t, node),
		configs: NewTemplatedConfigGetter(dependencies, t, node),
	***REMOVED***
***REMOVED***

func (t templatedDependencyGetter) Secrets() exec.SecretGetter ***REMOVED***
	return t.secrets
***REMOVED***

func (t templatedDependencyGetter) Configs() exec.ConfigGetter ***REMOVED***
	return t.configs
***REMOVED***
