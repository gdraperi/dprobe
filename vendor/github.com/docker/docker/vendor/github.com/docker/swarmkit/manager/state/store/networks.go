package store

import (
	"strings"

	"github.com/docker/swarmkit/api"
	memdb "github.com/hashicorp/go-memdb"
)

const tableNetwork = "network"

func init() ***REMOVED***
	register(ObjectStoreConfig***REMOVED***
		Table: &memdb.TableSchema***REMOVED***
			Name: tableNetwork,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				indexID: ***REMOVED***
					Name:    indexID,
					Unique:  true,
					Indexer: api.NetworkIndexerByID***REMOVED******REMOVED***,
				***REMOVED***,
				indexName: ***REMOVED***
					Name:    indexName,
					Unique:  true,
					Indexer: api.NetworkIndexerByName***REMOVED******REMOVED***,
				***REMOVED***,
				indexCustom: ***REMOVED***
					Name:         indexCustom,
					Indexer:      api.NetworkCustomIndexer***REMOVED******REMOVED***,
					AllowMissing: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Save: func(tx ReadTx, snapshot *api.StoreSnapshot) error ***REMOVED***
			var err error
			snapshot.Networks, err = FindNetworks(tx, All)
			return err
		***REMOVED***,
		Restore: func(tx Tx, snapshot *api.StoreSnapshot) error ***REMOVED***
			toStoreObj := make([]api.StoreObject, len(snapshot.Networks))
			for i, x := range snapshot.Networks ***REMOVED***
				toStoreObj[i] = x
			***REMOVED***
			return RestoreTable(tx, tableNetwork, toStoreObj)
		***REMOVED***,
		ApplyStoreAction: func(tx Tx, sa api.StoreAction) error ***REMOVED***
			switch v := sa.Target.(type) ***REMOVED***
			case *api.StoreAction_Network:
				obj := v.Network
				switch sa.Action ***REMOVED***
				case api.StoreActionKindCreate:
					return CreateNetwork(tx, obj)
				case api.StoreActionKindUpdate:
					return UpdateNetwork(tx, obj)
				case api.StoreActionKindRemove:
					return DeleteNetwork(tx, obj.ID)
				***REMOVED***
			***REMOVED***
			return errUnknownStoreAction
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// CreateNetwork adds a new network to the store.
// Returns ErrExist if the ID is already taken.
func CreateNetwork(tx Tx, n *api.Network) error ***REMOVED***
	// Ensure the name is not already in use.
	if tx.lookup(tableNetwork, indexName, strings.ToLower(n.Spec.Annotations.Name)) != nil ***REMOVED***
		return ErrNameConflict
	***REMOVED***

	return tx.create(tableNetwork, n)
***REMOVED***

// UpdateNetwork updates an existing network in the store.
// Returns ErrNotExist if the network doesn't exist.
func UpdateNetwork(tx Tx, n *api.Network) error ***REMOVED***
	// Ensure the name is either not in use or already used by this same Network.
	if existing := tx.lookup(tableNetwork, indexName, strings.ToLower(n.Spec.Annotations.Name)); existing != nil ***REMOVED***
		if existing.GetID() != n.ID ***REMOVED***
			return ErrNameConflict
		***REMOVED***
	***REMOVED***

	return tx.update(tableNetwork, n)
***REMOVED***

// DeleteNetwork removes a network from the store.
// Returns ErrNotExist if the network doesn't exist.
func DeleteNetwork(tx Tx, id string) error ***REMOVED***
	return tx.delete(tableNetwork, id)
***REMOVED***

// GetNetwork looks up a network by ID.
// Returns nil if the network doesn't exist.
func GetNetwork(tx ReadTx, id string) *api.Network ***REMOVED***
	n := tx.get(tableNetwork, id)
	if n == nil ***REMOVED***
		return nil
	***REMOVED***
	return n.(*api.Network)
***REMOVED***

// FindNetworks selects a set of networks and returns them.
func FindNetworks(tx ReadTx, by By) ([]*api.Network, error) ***REMOVED***
	checkType := func(by By) error ***REMOVED***
		switch by.(type) ***REMOVED***
		case byName, byNamePrefix, byIDPrefix, byCustom, byCustomPrefix:
			return nil
		default:
			return ErrInvalidFindBy
		***REMOVED***
	***REMOVED***

	networkList := []*api.Network***REMOVED******REMOVED***
	appendResult := func(o api.StoreObject) ***REMOVED***
		networkList = append(networkList, o.(*api.Network))
	***REMOVED***

	err := tx.find(tableNetwork, by, checkType, appendResult)
	return networkList, err
***REMOVED***
