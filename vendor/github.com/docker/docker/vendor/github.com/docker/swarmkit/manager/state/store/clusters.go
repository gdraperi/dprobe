package store

import (
	"strings"

	"github.com/docker/swarmkit/api"
	memdb "github.com/hashicorp/go-memdb"
)

const (
	tableCluster = "cluster"

	// DefaultClusterName is the default name to use for the cluster
	// object.
	DefaultClusterName = "default"
)

func init() ***REMOVED***
	register(ObjectStoreConfig***REMOVED***
		Table: &memdb.TableSchema***REMOVED***
			Name: tableCluster,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				indexID: ***REMOVED***
					Name:    indexID,
					Unique:  true,
					Indexer: api.ClusterIndexerByID***REMOVED******REMOVED***,
				***REMOVED***,
				indexName: ***REMOVED***
					Name:    indexName,
					Unique:  true,
					Indexer: api.ClusterIndexerByName***REMOVED******REMOVED***,
				***REMOVED***,
				indexCustom: ***REMOVED***
					Name:         indexCustom,
					Indexer:      api.ClusterCustomIndexer***REMOVED******REMOVED***,
					AllowMissing: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Save: func(tx ReadTx, snapshot *api.StoreSnapshot) error ***REMOVED***
			var err error
			snapshot.Clusters, err = FindClusters(tx, All)
			return err
		***REMOVED***,
		Restore: func(tx Tx, snapshot *api.StoreSnapshot) error ***REMOVED***
			toStoreObj := make([]api.StoreObject, len(snapshot.Clusters))
			for i, x := range snapshot.Clusters ***REMOVED***
				toStoreObj[i] = x
			***REMOVED***
			return RestoreTable(tx, tableCluster, toStoreObj)
		***REMOVED***,
		ApplyStoreAction: func(tx Tx, sa api.StoreAction) error ***REMOVED***
			switch v := sa.Target.(type) ***REMOVED***
			case *api.StoreAction_Cluster:
				obj := v.Cluster
				switch sa.Action ***REMOVED***
				case api.StoreActionKindCreate:
					return CreateCluster(tx, obj)
				case api.StoreActionKindUpdate:
					return UpdateCluster(tx, obj)
				case api.StoreActionKindRemove:
					return DeleteCluster(tx, obj.ID)
				***REMOVED***
			***REMOVED***
			return errUnknownStoreAction
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// CreateCluster adds a new cluster to the store.
// Returns ErrExist if the ID is already taken.
func CreateCluster(tx Tx, c *api.Cluster) error ***REMOVED***
	// Ensure the name is not already in use.
	if tx.lookup(tableCluster, indexName, strings.ToLower(c.Spec.Annotations.Name)) != nil ***REMOVED***
		return ErrNameConflict
	***REMOVED***

	return tx.create(tableCluster, c)
***REMOVED***

// UpdateCluster updates an existing cluster in the store.
// Returns ErrNotExist if the cluster doesn't exist.
func UpdateCluster(tx Tx, c *api.Cluster) error ***REMOVED***
	// Ensure the name is either not in use or already used by this same Cluster.
	if existing := tx.lookup(tableCluster, indexName, strings.ToLower(c.Spec.Annotations.Name)); existing != nil ***REMOVED***
		if existing.GetID() != c.ID ***REMOVED***
			return ErrNameConflict
		***REMOVED***
	***REMOVED***

	return tx.update(tableCluster, c)
***REMOVED***

// DeleteCluster removes a cluster from the store.
// Returns ErrNotExist if the cluster doesn't exist.
func DeleteCluster(tx Tx, id string) error ***REMOVED***
	return tx.delete(tableCluster, id)
***REMOVED***

// GetCluster looks up a cluster by ID.
// Returns nil if the cluster doesn't exist.
func GetCluster(tx ReadTx, id string) *api.Cluster ***REMOVED***
	n := tx.get(tableCluster, id)
	if n == nil ***REMOVED***
		return nil
	***REMOVED***
	return n.(*api.Cluster)
***REMOVED***

// FindClusters selects a set of clusters and returns them.
func FindClusters(tx ReadTx, by By) ([]*api.Cluster, error) ***REMOVED***
	checkType := func(by By) error ***REMOVED***
		switch by.(type) ***REMOVED***
		case byName, byNamePrefix, byIDPrefix, byCustom, byCustomPrefix:
			return nil
		default:
			return ErrInvalidFindBy
		***REMOVED***
	***REMOVED***

	clusterList := []*api.Cluster***REMOVED******REMOVED***
	appendResult := func(o api.StoreObject) ***REMOVED***
		clusterList = append(clusterList, o.(*api.Cluster))
	***REMOVED***

	err := tx.find(tableCluster, by, checkType, appendResult)
	return clusterList, err
***REMOVED***
