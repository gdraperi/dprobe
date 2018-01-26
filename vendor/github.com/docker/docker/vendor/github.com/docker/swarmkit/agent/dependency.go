package agent

import (
	"github.com/docker/swarmkit/agent/configs"
	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/agent/secrets"
	"github.com/docker/swarmkit/api"
)

type dependencyManager struct ***REMOVED***
	secrets exec.SecretsManager
	configs exec.ConfigsManager
***REMOVED***

// NewDependencyManager creates a dependency manager object that wraps
// objects which provide access to various dependency types.
func NewDependencyManager() exec.DependencyManager ***REMOVED***
	return &dependencyManager***REMOVED***
		secrets: secrets.NewManager(),
		configs: configs.NewManager(),
	***REMOVED***
***REMOVED***

func (d *dependencyManager) Secrets() exec.SecretsManager ***REMOVED***
	return d.secrets
***REMOVED***

func (d *dependencyManager) Configs() exec.ConfigsManager ***REMOVED***
	return d.configs
***REMOVED***

type dependencyGetter struct ***REMOVED***
	secrets exec.SecretGetter
	configs exec.ConfigGetter
***REMOVED***

func (d *dependencyGetter) Secrets() exec.SecretGetter ***REMOVED***
	return d.secrets
***REMOVED***

func (d *dependencyGetter) Configs() exec.ConfigGetter ***REMOVED***
	return d.configs
***REMOVED***

// Restrict provides getters that only allows access to the dependencies
// referenced by the task.
func Restrict(dependencies exec.DependencyManager, t *api.Task) exec.DependencyGetter ***REMOVED***
	return &dependencyGetter***REMOVED***
		secrets: secrets.Restrict(dependencies.Secrets(), t),
		configs: configs.Restrict(dependencies.Configs(), t),
	***REMOVED***
***REMOVED***
