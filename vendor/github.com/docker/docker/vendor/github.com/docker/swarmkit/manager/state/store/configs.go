package store

import (
	"strings"

	"github.com/docker/swarmkit/api"
	memdb "github.com/hashicorp/go-memdb"
)

const tableConfig = "config"

func init() ***REMOVED***
	register(ObjectStoreConfig***REMOVED***
		Table: &memdb.TableSchema***REMOVED***
			Name: tableConfig,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				indexID: ***REMOVED***
					Name:    indexID,
					Unique:  true,
					Indexer: api.ConfigIndexerByID***REMOVED******REMOVED***,
				***REMOVED***,
				indexName: ***REMOVED***
					Name:    indexName,
					Unique:  true,
					Indexer: api.ConfigIndexerByName***REMOVED******REMOVED***,
				***REMOVED***,
				indexCustom: ***REMOVED***
					Name:         indexCustom,
					Indexer:      api.ConfigCustomIndexer***REMOVED******REMOVED***,
					AllowMissing: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Save: func(tx ReadTx, snapshot *api.StoreSnapshot) error ***REMOVED***
			var err error
			snapshot.Configs, err = FindConfigs(tx, All)
			return err
		***REMOVED***,
		Restore: func(tx Tx, snapshot *api.StoreSnapshot) error ***REMOVED***
			toStoreObj := make([]api.StoreObject, len(snapshot.Configs))
			for i, x := range snapshot.Configs ***REMOVED***
				toStoreObj[i] = x
			***REMOVED***
			return RestoreTable(tx, tableConfig, toStoreObj)
		***REMOVED***,
		ApplyStoreAction: func(tx Tx, sa api.StoreAction) error ***REMOVED***
			switch v := sa.Target.(type) ***REMOVED***
			case *api.StoreAction_Config:
				obj := v.Config
				switch sa.Action ***REMOVED***
				case api.StoreActionKindCreate:
					return CreateConfig(tx, obj)
				case api.StoreActionKindUpdate:
					return UpdateConfig(tx, obj)
				case api.StoreActionKindRemove:
					return DeleteConfig(tx, obj.ID)
				***REMOVED***
			***REMOVED***
			return errUnknownStoreAction
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// CreateConfig adds a new config to the store.
// Returns ErrExist if the ID is already taken.
func CreateConfig(tx Tx, c *api.Config) error ***REMOVED***
	// Ensure the name is not already in use.
	if tx.lookup(tableConfig, indexName, strings.ToLower(c.Spec.Annotations.Name)) != nil ***REMOVED***
		return ErrNameConflict
	***REMOVED***

	return tx.create(tableConfig, c)
***REMOVED***

// UpdateConfig updates an existing config in the store.
// Returns ErrNotExist if the config doesn't exist.
func UpdateConfig(tx Tx, c *api.Config) error ***REMOVED***
	// Ensure the name is either not in use or already used by this same Config.
	if existing := tx.lookup(tableConfig, indexName, strings.ToLower(c.Spec.Annotations.Name)); existing != nil ***REMOVED***
		if existing.GetID() != c.ID ***REMOVED***
			return ErrNameConflict
		***REMOVED***
	***REMOVED***

	return tx.update(tableConfig, c)
***REMOVED***

// DeleteConfig removes a config from the store.
// Returns ErrNotExist if the config doesn't exist.
func DeleteConfig(tx Tx, id string) error ***REMOVED***
	return tx.delete(tableConfig, id)
***REMOVED***

// GetConfig looks up a config by ID.
// Returns nil if the config doesn't exist.
func GetConfig(tx ReadTx, id string) *api.Config ***REMOVED***
	c := tx.get(tableConfig, id)
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	return c.(*api.Config)
***REMOVED***

// FindConfigs selects a set of configs and returns them.
func FindConfigs(tx ReadTx, by By) ([]*api.Config, error) ***REMOVED***
	checkType := func(by By) error ***REMOVED***
		switch by.(type) ***REMOVED***
		case byName, byNamePrefix, byIDPrefix, byCustom, byCustomPrefix:
			return nil
		default:
			return ErrInvalidFindBy
		***REMOVED***
	***REMOVED***

	configList := []*api.Config***REMOVED******REMOVED***
	appendResult := func(o api.StoreObject) ***REMOVED***
		configList = append(configList, o.(*api.Config))
	***REMOVED***

	err := tx.find(tableConfig, by, checkType, appendResult)
	return configList, err
***REMOVED***
