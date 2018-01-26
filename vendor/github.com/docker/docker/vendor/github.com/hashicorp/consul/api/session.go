package api

import (
	"fmt"
	"time"
)

const (
	// SessionBehaviorRelease is the default behavior and causes
	// all associated locks to be released on session invalidation.
	SessionBehaviorRelease = "release"

	// SessionBehaviorDelete is new in Consul 0.5 and changes the
	// behavior to delete all associated locks on session invalidation.
	// It can be used in a way similar to Ephemeral Nodes in ZooKeeper.
	SessionBehaviorDelete = "delete"
)

// SessionEntry represents a session in consul
type SessionEntry struct ***REMOVED***
	CreateIndex uint64
	ID          string
	Name        string
	Node        string
	Checks      []string
	LockDelay   time.Duration
	Behavior    string
	TTL         string
***REMOVED***

// Session can be used to query the Session endpoints
type Session struct ***REMOVED***
	c *Client
***REMOVED***

// Session returns a handle to the session endpoints
func (c *Client) Session() *Session ***REMOVED***
	return &Session***REMOVED***c***REMOVED***
***REMOVED***

// CreateNoChecks is like Create but is used specifically to create
// a session with no associated health checks.
func (s *Session) CreateNoChecks(se *SessionEntry, q *WriteOptions) (string, *WriteMeta, error) ***REMOVED***
	body := make(map[string]interface***REMOVED******REMOVED***)
	body["Checks"] = []string***REMOVED******REMOVED***
	if se != nil ***REMOVED***
		if se.Name != "" ***REMOVED***
			body["Name"] = se.Name
		***REMOVED***
		if se.Node != "" ***REMOVED***
			body["Node"] = se.Node
		***REMOVED***
		if se.LockDelay != 0 ***REMOVED***
			body["LockDelay"] = durToMsec(se.LockDelay)
		***REMOVED***
		if se.Behavior != "" ***REMOVED***
			body["Behavior"] = se.Behavior
		***REMOVED***
		if se.TTL != "" ***REMOVED***
			body["TTL"] = se.TTL
		***REMOVED***
	***REMOVED***
	return s.create(body, q)

***REMOVED***

// Create makes a new session. Providing a session entry can
// customize the session. It can also be nil to use defaults.
func (s *Session) Create(se *SessionEntry, q *WriteOptions) (string, *WriteMeta, error) ***REMOVED***
	var obj interface***REMOVED******REMOVED***
	if se != nil ***REMOVED***
		body := make(map[string]interface***REMOVED******REMOVED***)
		obj = body
		if se.Name != "" ***REMOVED***
			body["Name"] = se.Name
		***REMOVED***
		if se.Node != "" ***REMOVED***
			body["Node"] = se.Node
		***REMOVED***
		if se.LockDelay != 0 ***REMOVED***
			body["LockDelay"] = durToMsec(se.LockDelay)
		***REMOVED***
		if len(se.Checks) > 0 ***REMOVED***
			body["Checks"] = se.Checks
		***REMOVED***
		if se.Behavior != "" ***REMOVED***
			body["Behavior"] = se.Behavior
		***REMOVED***
		if se.TTL != "" ***REMOVED***
			body["TTL"] = se.TTL
		***REMOVED***
	***REMOVED***
	return s.create(obj, q)
***REMOVED***

func (s *Session) create(obj interface***REMOVED******REMOVED***, q *WriteOptions) (string, *WriteMeta, error) ***REMOVED***
	var out struct***REMOVED*** ID string ***REMOVED***
	wm, err := s.c.write("/v1/session/create", obj, &out, q)
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	return out.ID, wm, nil
***REMOVED***

// Destroy invalides a given session
func (s *Session) Destroy(id string, q *WriteOptions) (*WriteMeta, error) ***REMOVED***
	wm, err := s.c.write("/v1/session/destroy/"+id, nil, nil, q)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return wm, nil
***REMOVED***

// Renew renews the TTL on a given session
func (s *Session) Renew(id string, q *WriteOptions) (*SessionEntry, *WriteMeta, error) ***REMOVED***
	var entries []*SessionEntry
	wm, err := s.c.write("/v1/session/renew/"+id, nil, &entries, q)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if len(entries) > 0 ***REMOVED***
		return entries[0], wm, nil
	***REMOVED***
	return nil, wm, nil
***REMOVED***

// RenewPeriodic is used to periodically invoke Session.Renew on a
// session until a doneCh is closed. This is meant to be used in a long running
// goroutine to ensure a session stays valid.
func (s *Session) RenewPeriodic(initialTTL string, id string, q *WriteOptions, doneCh chan struct***REMOVED******REMOVED***) error ***REMOVED***
	ttl, err := time.ParseDuration(initialTTL)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	waitDur := ttl / 2
	lastRenewTime := time.Now()
	var lastErr error
	for ***REMOVED***
		if time.Since(lastRenewTime) > ttl ***REMOVED***
			return lastErr
		***REMOVED***
		select ***REMOVED***
		case <-time.After(waitDur):
			entry, _, err := s.Renew(id, q)
			if err != nil ***REMOVED***
				waitDur = time.Second
				lastErr = err
				continue
			***REMOVED***
			if entry == nil ***REMOVED***
				waitDur = time.Second
				lastErr = fmt.Errorf("No SessionEntry returned")
				continue
			***REMOVED***

			// Handle the server updating the TTL
			ttl, _ = time.ParseDuration(entry.TTL)
			waitDur = ttl / 2
			lastRenewTime = time.Now()

		case <-doneCh:
			// Attempt a session destroy
			s.Destroy(id, q)
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// Info looks up a single session
func (s *Session) Info(id string, q *QueryOptions) (*SessionEntry, *QueryMeta, error) ***REMOVED***
	var entries []*SessionEntry
	qm, err := s.c.query("/v1/session/info/"+id, &entries, q)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if len(entries) > 0 ***REMOVED***
		return entries[0], qm, nil
	***REMOVED***
	return nil, qm, nil
***REMOVED***

// List gets sessions for a node
func (s *Session) Node(node string, q *QueryOptions) ([]*SessionEntry, *QueryMeta, error) ***REMOVED***
	var entries []*SessionEntry
	qm, err := s.c.query("/v1/session/node/"+node, &entries, q)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return entries, qm, nil
***REMOVED***

// List gets all active sessions
func (s *Session) List(q *QueryOptions) ([]*SessionEntry, *QueryMeta, error) ***REMOVED***
	var entries []*SessionEntry
	qm, err := s.c.query("/v1/session/list", &entries, q)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return entries, qm, nil
***REMOVED***
