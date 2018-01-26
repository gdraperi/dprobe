package store

import (
	"github.com/docker/swarmkit/api"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/pkg/errors"
)

const tableResource = "resource"

func init() ***REMOVED***
	register(ObjectStoreConfig***REMOVED***
		Table: &memdb.TableSchema***REMOVED***
			Name: tableResource,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				indexID: ***REMOVED***
					Name:    indexID,
					Unique:  true,
					Indexer: resourceIndexerByID***REMOVED******REMOVED***,
				***REMOVED***,
				indexName: ***REMOVED***
					Name:    indexName,
					Unique:  true,
					Indexer: resourceIndexerByName***REMOVED******REMOVED***,
				***REMOVED***,
				indexKind: ***REMOVED***
					Name:    indexKind,
					Indexer: resourceIndexerByKind***REMOVED******REMOVED***,
				***REMOVED***,
				indexCustom: ***REMOVED***
					Name:         indexCustom,
					Indexer:      resourceCustomIndexer***REMOVED******REMOVED***,
					AllowMissing: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Save: func(tx ReadTx, snapshot *api.StoreSnapshot) error ***REMOVED***
			var err error
			snapshot.Resources, err = FindResources(tx, All)
			return err
		***REMOVED***,
		Restore: func(tx Tx, snapshot *api.StoreSnapshot) error ***REMOVED***
			toStoreObj := make([]api.StoreObject, len(snapshot.Resources))
			for i, x := range snapshot.Resources ***REMOVED***
				toStoreObj[i] = resourceEntry***REMOVED***x***REMOVED***
			***REMOVED***
			return RestoreTable(tx, tableResource, toStoreObj)
		***REMOVED***,
		ApplyStoreAction: func(tx Tx, sa api.StoreAction) error ***REMOVED***
			switch v := sa.Target.(type) ***REMOVED***
			case *api.StoreAction_Resource:
				obj := v.Resource
				switch sa.Action ***REMOVED***
				case api.StoreActionKindCreate:
					return CreateResource(tx, obj)
				case api.StoreActionKindUpdate:
					return UpdateResource(tx, obj)
				case api.StoreActionKindRemove:
					return DeleteResource(tx, obj.ID)
				***REMOVED***
			***REMOVED***
			return errUnknownStoreAction
		***REMOVED***,
	***REMOVED***)
***REMOVED***

type resourceEntry struct ***REMOVED***
	*api.Resource
***REMOVED***

func (r resourceEntry) CopyStoreObject() api.StoreObject ***REMOVED***
	return resourceEntry***REMOVED***Resource: r.Resource.Copy()***REMOVED***
***REMOVED***

// ensure that when update events are emitted, we unwrap resourceEntry
func (r resourceEntry) EventUpdate(oldObject api.StoreObject) api.Event ***REMOVED***
	if oldObject != nil ***REMOVED***
		return api.EventUpdateResource***REMOVED***Resource: r.Resource, OldResource: oldObject.(resourceEntry).Resource***REMOVED***
	***REMOVED***
	return api.EventUpdateResource***REMOVED***Resource: r.Resource***REMOVED***
***REMOVED***

func confirmExtension(tx Tx, r *api.Resource) error ***REMOVED***
	// There must be an extension corresponding to the Kind field.
	extensions, err := FindExtensions(tx, ByName(r.Kind))
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to query extensions")
	***REMOVED***
	if len(extensions) == 0 ***REMOVED***
		return errors.Errorf("object kind %s is unregistered", r.Kind)
	***REMOVED***
	return nil
***REMOVED***

// CreateResource adds a new resource object to the store.
// Returns ErrExist if the ID is already taken.
func CreateResource(tx Tx, r *api.Resource) error ***REMOVED***
	if err := confirmExtension(tx, r); err != nil ***REMOVED***
		return err
	***REMOVED***
	return tx.create(tableResource, resourceEntry***REMOVED***r***REMOVED***)
***REMOVED***

// UpdateResource updates an existing resource object in the store.
// Returns ErrNotExist if the object doesn't exist.
func UpdateResource(tx Tx, r *api.Resource) error ***REMOVED***
	if err := confirmExtension(tx, r); err != nil ***REMOVED***
		return err
	***REMOVED***
	return tx.update(tableResource, resourceEntry***REMOVED***r***REMOVED***)
***REMOVED***

// DeleteResource removes a resource object from the store.
// Returns ErrNotExist if the object doesn't exist.
func DeleteResource(tx Tx, id string) error ***REMOVED***
	return tx.delete(tableResource, id)
***REMOVED***

// GetResource looks up a resource object by ID.
// Returns nil if the object doesn't exist.
func GetResource(tx ReadTx, id string) *api.Resource ***REMOVED***
	r := tx.get(tableResource, id)
	if r == nil ***REMOVED***
		return nil
	***REMOVED***
	return r.(resourceEntry).Resource
***REMOVED***

// FindResources selects a set of resource objects and returns them.
func FindResources(tx ReadTx, by By) ([]*api.Resource, error) ***REMOVED***
	checkType := func(by By) error ***REMOVED***
		switch by.(type) ***REMOVED***
		case byIDPrefix, byName, byKind, byCustom, byCustomPrefix:
			return nil
		default:
			return ErrInvalidFindBy
		***REMOVED***
	***REMOVED***

	resourceList := []*api.Resource***REMOVED******REMOVED***
	appendResult := func(o api.StoreObject) ***REMOVED***
		resourceList = append(resourceList, o.(resourceEntry).Resource)
	***REMOVED***

	err := tx.find(tableResource, by, checkType, appendResult)
	return resourceList, err
***REMOVED***

type resourceIndexerByKind struct***REMOVED******REMOVED***

func (ri resourceIndexerByKind) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ri resourceIndexerByKind) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	r := obj.(resourceEntry)

	// Add the null character as a terminator
	val := r.Resource.Kind + "\x00"
	return true, []byte(val), nil
***REMOVED***

type resourceIndexerByID struct***REMOVED******REMOVED***

func (indexer resourceIndexerByID) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ResourceIndexerByID***REMOVED******REMOVED***.FromArgs(args...)
***REMOVED***
func (indexer resourceIndexerByID) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ResourceIndexerByID***REMOVED******REMOVED***.PrefixFromArgs(args...)
***REMOVED***
func (indexer resourceIndexerByID) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	return api.ResourceIndexerByID***REMOVED******REMOVED***.FromObject(obj.(resourceEntry).Resource)
***REMOVED***

type resourceIndexerByName struct***REMOVED******REMOVED***

func (indexer resourceIndexerByName) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ResourceIndexerByName***REMOVED******REMOVED***.FromArgs(args...)
***REMOVED***
func (indexer resourceIndexerByName) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ResourceIndexerByName***REMOVED******REMOVED***.PrefixFromArgs(args...)
***REMOVED***
func (indexer resourceIndexerByName) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	return api.ResourceIndexerByName***REMOVED******REMOVED***.FromObject(obj.(resourceEntry).Resource)
***REMOVED***

type resourceCustomIndexer struct***REMOVED******REMOVED***

func (indexer resourceCustomIndexer) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ResourceCustomIndexer***REMOVED******REMOVED***.FromArgs(args...)
***REMOVED***
func (indexer resourceCustomIndexer) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ResourceCustomIndexer***REMOVED******REMOVED***.PrefixFromArgs(args...)
***REMOVED***
func (indexer resourceCustomIndexer) FromObject(obj interface***REMOVED******REMOVED***) (bool, [][]byte, error) ***REMOVED***
	return api.ResourceCustomIndexer***REMOVED******REMOVED***.FromObject(obj.(resourceEntry).Resource)
***REMOVED***
