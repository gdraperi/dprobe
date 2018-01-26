package manager

import (
	"reflect"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/state/store"
)

// IsStateDirty returns true if any objects have been added to raft which make
// the state "dirty". Currently, the existence of any object other than the
// default cluster or the local node implies a dirty state.
func (m *Manager) IsStateDirty() (bool, error) ***REMOVED***
	var (
		storeSnapshot *api.StoreSnapshot
		err           error
	)
	m.raftNode.MemoryStore().View(func(readTx store.ReadTx) ***REMOVED***
		storeSnapshot, err = m.raftNode.MemoryStore().Save(readTx)
	***REMOVED***)

	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	// Check Nodes and Clusters fields.
	nodeID := m.config.SecurityConfig.ClientTLSCreds.NodeID()
	if len(storeSnapshot.Nodes) > 1 || (len(storeSnapshot.Nodes) == 1 && storeSnapshot.Nodes[0].ID != nodeID) ***REMOVED***
		return true, nil
	***REMOVED***

	clusterID := m.config.SecurityConfig.ClientTLSCreds.Organization()
	if len(storeSnapshot.Clusters) > 1 || (len(storeSnapshot.Clusters) == 1 && storeSnapshot.Clusters[0].ID != clusterID) ***REMOVED***
		return true, nil
	***REMOVED***

	// Use reflection to check that other fields don't have values. This
	// lets us implement a whitelist-type approach, where we don't need to
	// remember to add individual types here.

	val := reflect.ValueOf(*storeSnapshot)
	numFields := val.NumField()

	for i := 0; i != numFields; i++ ***REMOVED***
		field := val.Field(i)
		structField := val.Type().Field(i)
		if structField.Type.Kind() != reflect.Slice ***REMOVED***
			panic("unexpected field type in StoreSnapshot")
		***REMOVED***
		if structField.Name != "Nodes" && structField.Name != "Clusters" && structField.Name != "Networks" && field.Len() != 0 ***REMOVED***
			// One of the other data types has an entry
			return true, nil
		***REMOVED***
	***REMOVED***

	return false, nil
***REMOVED***
