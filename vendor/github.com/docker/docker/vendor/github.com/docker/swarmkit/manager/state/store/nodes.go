package store

import (
	"strconv"
	"strings"

	"github.com/docker/swarmkit/api"
	memdb "github.com/hashicorp/go-memdb"
)

const tableNode = "node"

func init() ***REMOVED***
	register(ObjectStoreConfig***REMOVED***
		Table: &memdb.TableSchema***REMOVED***
			Name: tableNode,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				indexID: ***REMOVED***
					Name:    indexID,
					Unique:  true,
					Indexer: api.NodeIndexerByID***REMOVED******REMOVED***,
				***REMOVED***,
				// TODO(aluzzardi): Use `indexHostname` instead.
				indexName: ***REMOVED***
					Name:         indexName,
					AllowMissing: true,
					Indexer:      nodeIndexerByHostname***REMOVED******REMOVED***,
				***REMOVED***,
				indexRole: ***REMOVED***
					Name:    indexRole,
					Indexer: nodeIndexerByRole***REMOVED******REMOVED***,
				***REMOVED***,
				indexMembership: ***REMOVED***
					Name:    indexMembership,
					Indexer: nodeIndexerByMembership***REMOVED******REMOVED***,
				***REMOVED***,
				indexCustom: ***REMOVED***
					Name:         indexCustom,
					Indexer:      api.NodeCustomIndexer***REMOVED******REMOVED***,
					AllowMissing: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Save: func(tx ReadTx, snapshot *api.StoreSnapshot) error ***REMOVED***
			var err error
			snapshot.Nodes, err = FindNodes(tx, All)
			return err
		***REMOVED***,
		Restore: func(tx Tx, snapshot *api.StoreSnapshot) error ***REMOVED***
			toStoreObj := make([]api.StoreObject, len(snapshot.Nodes))
			for i, x := range snapshot.Nodes ***REMOVED***
				toStoreObj[i] = x
			***REMOVED***
			return RestoreTable(tx, tableNode, toStoreObj)
		***REMOVED***,
		ApplyStoreAction: func(tx Tx, sa api.StoreAction) error ***REMOVED***
			switch v := sa.Target.(type) ***REMOVED***
			case *api.StoreAction_Node:
				obj := v.Node
				switch sa.Action ***REMOVED***
				case api.StoreActionKindCreate:
					return CreateNode(tx, obj)
				case api.StoreActionKindUpdate:
					return UpdateNode(tx, obj)
				case api.StoreActionKindRemove:
					return DeleteNode(tx, obj.ID)
				***REMOVED***
			***REMOVED***
			return errUnknownStoreAction
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// CreateNode adds a new node to the store.
// Returns ErrExist if the ID is already taken.
func CreateNode(tx Tx, n *api.Node) error ***REMOVED***
	return tx.create(tableNode, n)
***REMOVED***

// UpdateNode updates an existing node in the store.
// Returns ErrNotExist if the node doesn't exist.
func UpdateNode(tx Tx, n *api.Node) error ***REMOVED***
	return tx.update(tableNode, n)
***REMOVED***

// DeleteNode removes a node from the store.
// Returns ErrNotExist if the node doesn't exist.
func DeleteNode(tx Tx, id string) error ***REMOVED***
	return tx.delete(tableNode, id)
***REMOVED***

// GetNode looks up a node by ID.
// Returns nil if the node doesn't exist.
func GetNode(tx ReadTx, id string) *api.Node ***REMOVED***
	n := tx.get(tableNode, id)
	if n == nil ***REMOVED***
		return nil
	***REMOVED***
	return n.(*api.Node)
***REMOVED***

// FindNodes selects a set of nodes and returns them.
func FindNodes(tx ReadTx, by By) ([]*api.Node, error) ***REMOVED***
	checkType := func(by By) error ***REMOVED***
		switch by.(type) ***REMOVED***
		case byName, byNamePrefix, byIDPrefix, byRole, byMembership, byCustom, byCustomPrefix:
			return nil
		default:
			return ErrInvalidFindBy
		***REMOVED***
	***REMOVED***

	nodeList := []*api.Node***REMOVED******REMOVED***
	appendResult := func(o api.StoreObject) ***REMOVED***
		nodeList = append(nodeList, o.(*api.Node))
	***REMOVED***

	err := tx.find(tableNode, by, checkType, appendResult)
	return nodeList, err
***REMOVED***

type nodeIndexerByHostname struct***REMOVED******REMOVED***

func (ni nodeIndexerByHostname) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ni nodeIndexerByHostname) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	n := obj.(*api.Node)

	if n.Description == nil ***REMOVED***
		return false, nil, nil
	***REMOVED***
	// Add the null character as a terminator
	return true, []byte(strings.ToLower(n.Description.Hostname) + "\x00"), nil
***REMOVED***

func (ni nodeIndexerByHostname) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return prefixFromArgs(args...)
***REMOVED***

type nodeIndexerByRole struct***REMOVED******REMOVED***

func (ni nodeIndexerByRole) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ni nodeIndexerByRole) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	n := obj.(*api.Node)

	// Add the null character as a terminator
	return true, []byte(strconv.FormatInt(int64(n.Role), 10) + "\x00"), nil
***REMOVED***

type nodeIndexerByMembership struct***REMOVED******REMOVED***

func (ni nodeIndexerByMembership) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ni nodeIndexerByMembership) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	n := obj.(*api.Node)

	// Add the null character as a terminator
	return true, []byte(strconv.FormatInt(int64(n.Spec.Membership), 10) + "\x00"), nil
***REMOVED***
