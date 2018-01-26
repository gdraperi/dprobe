package store

import (
	"errors"
	"strings"

	"github.com/docker/swarmkit/api"
	memdb "github.com/hashicorp/go-memdb"
)

const tableExtension = "extension"

func init() ***REMOVED***
	register(ObjectStoreConfig***REMOVED***
		Table: &memdb.TableSchema***REMOVED***
			Name: tableExtension,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				indexID: ***REMOVED***
					Name:    indexID,
					Unique:  true,
					Indexer: extensionIndexerByID***REMOVED******REMOVED***,
				***REMOVED***,
				indexName: ***REMOVED***
					Name:    indexName,
					Unique:  true,
					Indexer: extensionIndexerByName***REMOVED******REMOVED***,
				***REMOVED***,
				indexCustom: ***REMOVED***
					Name:         indexCustom,
					Indexer:      extensionCustomIndexer***REMOVED******REMOVED***,
					AllowMissing: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Save: func(tx ReadTx, snapshot *api.StoreSnapshot) error ***REMOVED***
			var err error
			snapshot.Extensions, err = FindExtensions(tx, All)
			return err
		***REMOVED***,
		Restore: func(tx Tx, snapshot *api.StoreSnapshot) error ***REMOVED***
			toStoreObj := make([]api.StoreObject, len(snapshot.Extensions))
			for i, x := range snapshot.Extensions ***REMOVED***
				toStoreObj[i] = extensionEntry***REMOVED***x***REMOVED***
			***REMOVED***
			return RestoreTable(tx, tableExtension, toStoreObj)
		***REMOVED***,
		ApplyStoreAction: func(tx Tx, sa api.StoreAction) error ***REMOVED***
			switch v := sa.Target.(type) ***REMOVED***
			case *api.StoreAction_Extension:
				obj := v.Extension
				switch sa.Action ***REMOVED***
				case api.StoreActionKindCreate:
					return CreateExtension(tx, obj)
				case api.StoreActionKindUpdate:
					return UpdateExtension(tx, obj)
				case api.StoreActionKindRemove:
					return DeleteExtension(tx, obj.ID)
				***REMOVED***
			***REMOVED***
			return errUnknownStoreAction
		***REMOVED***,
	***REMOVED***)
***REMOVED***

type extensionEntry struct ***REMOVED***
	*api.Extension
***REMOVED***

func (e extensionEntry) CopyStoreObject() api.StoreObject ***REMOVED***
	return extensionEntry***REMOVED***Extension: e.Extension.Copy()***REMOVED***
***REMOVED***

// ensure that when update events are emitted, we unwrap extensionEntry
func (e extensionEntry) EventUpdate(oldObject api.StoreObject) api.Event ***REMOVED***
	if oldObject != nil ***REMOVED***
		return api.EventUpdateExtension***REMOVED***Extension: e.Extension, OldExtension: oldObject.(extensionEntry).Extension***REMOVED***
	***REMOVED***
	return api.EventUpdateExtension***REMOVED***Extension: e.Extension***REMOVED***
***REMOVED***

// CreateExtension adds a new extension to the store.
// Returns ErrExist if the ID is already taken.
func CreateExtension(tx Tx, e *api.Extension) error ***REMOVED***
	// Ensure the name is not already in use.
	if tx.lookup(tableExtension, indexName, strings.ToLower(e.Annotations.Name)) != nil ***REMOVED***
		return ErrNameConflict
	***REMOVED***

	// It can't conflict with built-in kinds either.
	if _, ok := schema.Tables[e.Annotations.Name]; ok ***REMOVED***
		return ErrNameConflict
	***REMOVED***

	return tx.create(tableExtension, extensionEntry***REMOVED***e***REMOVED***)
***REMOVED***

// UpdateExtension updates an existing extension in the store.
// Returns ErrNotExist if the object doesn't exist.
func UpdateExtension(tx Tx, e *api.Extension) error ***REMOVED***
	// TODO(aaronl): For the moment, extensions are immutable
	return errors.New("extensions are immutable")
***REMOVED***

// DeleteExtension removes an extension from the store.
// Returns ErrNotExist if the object doesn't exist.
func DeleteExtension(tx Tx, id string) error ***REMOVED***
	e := tx.get(tableExtension, id)
	if e == nil ***REMOVED***
		return ErrNotExist
	***REMOVED***

	resources, err := FindResources(tx, ByKind(e.(extensionEntry).Annotations.Name))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(resources) != 0 ***REMOVED***
		return errors.New("cannot delete extension because objects of this type exist in the data store")
	***REMOVED***

	return tx.delete(tableExtension, id)
***REMOVED***

// GetExtension looks up an extension by ID.
// Returns nil if the object doesn't exist.
func GetExtension(tx ReadTx, id string) *api.Extension ***REMOVED***
	e := tx.get(tableExtension, id)
	if e == nil ***REMOVED***
		return nil
	***REMOVED***
	return e.(extensionEntry).Extension
***REMOVED***

// FindExtensions selects a set of extensions and returns them.
func FindExtensions(tx ReadTx, by By) ([]*api.Extension, error) ***REMOVED***
	checkType := func(by By) error ***REMOVED***
		switch by.(type) ***REMOVED***
		case byIDPrefix, byName, byCustom, byCustomPrefix:
			return nil
		default:
			return ErrInvalidFindBy
		***REMOVED***
	***REMOVED***

	extensionList := []*api.Extension***REMOVED******REMOVED***
	appendResult := func(o api.StoreObject) ***REMOVED***
		extensionList = append(extensionList, o.(extensionEntry).Extension)
	***REMOVED***

	err := tx.find(tableExtension, by, checkType, appendResult)
	return extensionList, err
***REMOVED***

type extensionIndexerByID struct***REMOVED******REMOVED***

func (indexer extensionIndexerByID) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ExtensionIndexerByID***REMOVED******REMOVED***.FromArgs(args...)
***REMOVED***
func (indexer extensionIndexerByID) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ExtensionIndexerByID***REMOVED******REMOVED***.PrefixFromArgs(args...)
***REMOVED***
func (indexer extensionIndexerByID) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	return api.ExtensionIndexerByID***REMOVED******REMOVED***.FromObject(obj.(extensionEntry).Extension)
***REMOVED***

type extensionIndexerByName struct***REMOVED******REMOVED***

func (indexer extensionIndexerByName) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ExtensionIndexerByName***REMOVED******REMOVED***.FromArgs(args...)
***REMOVED***
func (indexer extensionIndexerByName) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ExtensionIndexerByName***REMOVED******REMOVED***.PrefixFromArgs(args...)
***REMOVED***
func (indexer extensionIndexerByName) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	return api.ExtensionIndexerByName***REMOVED******REMOVED***.FromObject(obj.(extensionEntry).Extension)
***REMOVED***

type extensionCustomIndexer struct***REMOVED******REMOVED***

func (indexer extensionCustomIndexer) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ExtensionCustomIndexer***REMOVED******REMOVED***.FromArgs(args...)
***REMOVED***
func (indexer extensionCustomIndexer) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return api.ExtensionCustomIndexer***REMOVED******REMOVED***.PrefixFromArgs(args...)
***REMOVED***
func (indexer extensionCustomIndexer) FromObject(obj interface***REMOVED******REMOVED***) (bool, [][]byte, error) ***REMOVED***
	return api.ExtensionCustomIndexer***REMOVED******REMOVED***.FromObject(obj.(extensionEntry).Extension)
***REMOVED***
