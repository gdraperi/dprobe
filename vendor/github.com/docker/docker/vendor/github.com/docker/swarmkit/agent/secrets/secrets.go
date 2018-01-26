package secrets

import (
	"fmt"
	"sync"

	"github.com/docker/swarmkit/agent/exec"
	"github.com/docker/swarmkit/api"
)

// secrets is a map that keeps all the currently available secrets to the agent
// mapped by secret ID.
type secrets struct ***REMOVED***
	mu sync.RWMutex
	m  map[string]*api.Secret
***REMOVED***

// NewManager returns a place to store secrets.
func NewManager() exec.SecretsManager ***REMOVED***
	return &secrets***REMOVED***
		m: make(map[string]*api.Secret),
	***REMOVED***
***REMOVED***

// Get returns a secret by ID.  If the secret doesn't exist, returns nil.
func (s *secrets) Get(secretID string) (*api.Secret, error) ***REMOVED***
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s, ok := s.m[secretID]; ok ***REMOVED***
		return s, nil
	***REMOVED***
	return nil, fmt.Errorf("secret %s not found", secretID)
***REMOVED***

// Add adds one or more secrets to the secret map.
func (s *secrets) Add(secrets ...api.Secret) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, secret := range secrets ***REMOVED***
		s.m[secret.ID] = secret.Copy()
	***REMOVED***
***REMOVED***

// Remove removes one or more secrets by ID from the secret map.  Succeeds
// whether or not the given IDs are in the map.
func (s *secrets) Remove(secrets []string) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, secret := range secrets ***REMOVED***
		delete(s.m, secret)
	***REMOVED***
***REMOVED***

// Reset removes all the secrets.
func (s *secrets) Reset() ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = make(map[string]*api.Secret)
***REMOVED***

// taskRestrictedSecretsProvider restricts the ids to the task.
type taskRestrictedSecretsProvider struct ***REMOVED***
	secrets   exec.SecretGetter
	secretIDs map[string]struct***REMOVED******REMOVED*** // allow list of secret ids
***REMOVED***

func (sp *taskRestrictedSecretsProvider) Get(secretID string) (*api.Secret, error) ***REMOVED***
	if _, ok := sp.secretIDs[secretID]; !ok ***REMOVED***
		return nil, fmt.Errorf("task not authorized to access secret %s", secretID)
	***REMOVED***

	return sp.secrets.Get(secretID)
***REMOVED***

// Restrict provides a getter that only allows access to the secrets
// referenced by the task.
func Restrict(secrets exec.SecretGetter, t *api.Task) exec.SecretGetter ***REMOVED***
	sids := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***

	container := t.Spec.GetContainer()
	if container != nil ***REMOVED***
		for _, ref := range container.Secrets ***REMOVED***
			sids[ref.SecretID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	return &taskRestrictedSecretsProvider***REMOVED***secrets: secrets, secretIDs: sids***REMOVED***
***REMOVED***
