package configs

import (
	"fmt"
	"sync"

	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
)

// configs is a map that keeps all the currently available configs to the agent
// mapped by config ID.
type configs struct ***REMOVED***
	mu sync.RWMutex
	m  map[string]*api.Config
***REMOVED***

// NewManager returns a place to store configs.
func NewManager() exec.ConfigsManager ***REMOVED***
	return &configs***REMOVED***
		m: make(map[string]*api.Config),
	***REMOVED***
***REMOVED***

// Get returns a config by ID.  If the config doesn't exist, returns nil.
func (r *configs) Get(configID string) (*api.Config, error) ***REMOVED***
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r, ok := r.m[configID]; ok ***REMOVED***
		return r, nil
	***REMOVED***
	return nil, fmt.Errorf("config %s not found", configID)
***REMOVED***

// Add adds one or more configs to the config map.
func (r *configs) Add(configs ...api.Config) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, config := range configs ***REMOVED***
		r.m[config.ID] = config.Copy()
	***REMOVED***
***REMOVED***

// Remove removes one or more configs by ID from the config map. Succeeds
// whether or not the given IDs are in the map.
func (r *configs) Remove(configs []string) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, config := range configs ***REMOVED***
		delete(r.m, config)
	***REMOVED***
***REMOVED***

// Reset removes all the configs.
func (r *configs) Reset() ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m = make(map[string]*api.Config)
***REMOVED***

// taskRestrictedConfigsProvider restricts the ids to the task.
type taskRestrictedConfigsProvider struct ***REMOVED***
	configs   exec.ConfigGetter
	configIDs map[string]struct***REMOVED******REMOVED*** // allow list of config ids
***REMOVED***

func (sp *taskRestrictedConfigsProvider) Get(configID string) (*api.Config, error) ***REMOVED***
	if _, ok := sp.configIDs[configID]; !ok ***REMOVED***
		return nil, fmt.Errorf("task not authorized to access config %s", configID)
	***REMOVED***

	return sp.configs.Get(configID)
***REMOVED***

// Restrict provides a getter that only allows access to the configs
// referenced by the task.
func Restrict(configs exec.ConfigGetter, t *api.Task) exec.ConfigGetter ***REMOVED***
	cids := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***

	container := t.Spec.GetContainer()
	if container != nil ***REMOVED***
		for _, configRef := range container.Configs ***REMOVED***
			cids[configRef.ConfigID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	return &taskRestrictedConfigsProvider***REMOVED***configs: configs, configIDs: cids***REMOVED***
***REMOVED***
