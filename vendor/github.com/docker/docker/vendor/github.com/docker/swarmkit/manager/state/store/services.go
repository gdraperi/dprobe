package store

import (
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/naming"
	memdb "github.com/hashicorp/go-memdb"
)

const tableService = "service"

func init() ***REMOVED***
	register(ObjectStoreConfig***REMOVED***
		Table: &memdb.TableSchema***REMOVED***
			Name: tableService,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				indexID: ***REMOVED***
					Name:    indexID,
					Unique:  true,
					Indexer: api.ServiceIndexerByID***REMOVED******REMOVED***,
				***REMOVED***,
				indexName: ***REMOVED***
					Name:    indexName,
					Unique:  true,
					Indexer: api.ServiceIndexerByName***REMOVED******REMOVED***,
				***REMOVED***,
				indexRuntime: ***REMOVED***
					Name:         indexRuntime,
					AllowMissing: true,
					Indexer:      serviceIndexerByRuntime***REMOVED******REMOVED***,
				***REMOVED***,
				indexNetwork: ***REMOVED***
					Name:         indexNetwork,
					AllowMissing: true,
					Indexer:      serviceIndexerByNetwork***REMOVED******REMOVED***,
				***REMOVED***,
				indexSecret: ***REMOVED***
					Name:         indexSecret,
					AllowMissing: true,
					Indexer:      serviceIndexerBySecret***REMOVED******REMOVED***,
				***REMOVED***,
				indexConfig: ***REMOVED***
					Name:         indexConfig,
					AllowMissing: true,
					Indexer:      serviceIndexerByConfig***REMOVED******REMOVED***,
				***REMOVED***,
				indexCustom: ***REMOVED***
					Name:         indexCustom,
					Indexer:      api.ServiceCustomIndexer***REMOVED******REMOVED***,
					AllowMissing: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Save: func(tx ReadTx, snapshot *api.StoreSnapshot) error ***REMOVED***
			var err error
			snapshot.Services, err = FindServices(tx, All)
			return err
		***REMOVED***,
		Restore: func(tx Tx, snapshot *api.StoreSnapshot) error ***REMOVED***
			toStoreObj := make([]api.StoreObject, len(snapshot.Services))
			for i, x := range snapshot.Services ***REMOVED***
				toStoreObj[i] = x
			***REMOVED***
			return RestoreTable(tx, tableService, toStoreObj)
		***REMOVED***,
		ApplyStoreAction: func(tx Tx, sa api.StoreAction) error ***REMOVED***
			switch v := sa.Target.(type) ***REMOVED***
			case *api.StoreAction_Service:
				obj := v.Service
				switch sa.Action ***REMOVED***
				case api.StoreActionKindCreate:
					return CreateService(tx, obj)
				case api.StoreActionKindUpdate:
					return UpdateService(tx, obj)
				case api.StoreActionKindRemove:
					return DeleteService(tx, obj.ID)
				***REMOVED***
			***REMOVED***
			return errUnknownStoreAction
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// CreateService adds a new service to the store.
// Returns ErrExist if the ID is already taken.
func CreateService(tx Tx, s *api.Service) error ***REMOVED***
	// Ensure the name is not already in use.
	if tx.lookup(tableService, indexName, strings.ToLower(s.Spec.Annotations.Name)) != nil ***REMOVED***
		return ErrNameConflict
	***REMOVED***

	return tx.create(tableService, s)
***REMOVED***

// UpdateService updates an existing service in the store.
// Returns ErrNotExist if the service doesn't exist.
func UpdateService(tx Tx, s *api.Service) error ***REMOVED***
	// Ensure the name is either not in use or already used by this same Service.
	if existing := tx.lookup(tableService, indexName, strings.ToLower(s.Spec.Annotations.Name)); existing != nil ***REMOVED***
		if existing.GetID() != s.ID ***REMOVED***
			return ErrNameConflict
		***REMOVED***
	***REMOVED***

	return tx.update(tableService, s)
***REMOVED***

// DeleteService removes a service from the store.
// Returns ErrNotExist if the service doesn't exist.
func DeleteService(tx Tx, id string) error ***REMOVED***
	return tx.delete(tableService, id)
***REMOVED***

// GetService looks up a service by ID.
// Returns nil if the service doesn't exist.
func GetService(tx ReadTx, id string) *api.Service ***REMOVED***
	s := tx.get(tableService, id)
	if s == nil ***REMOVED***
		return nil
	***REMOVED***
	return s.(*api.Service)
***REMOVED***

// FindServices selects a set of services and returns them.
func FindServices(tx ReadTx, by By) ([]*api.Service, error) ***REMOVED***
	checkType := func(by By) error ***REMOVED***
		switch by.(type) ***REMOVED***
		case byName, byNamePrefix, byIDPrefix, byRuntime, byReferencedNetworkID, byReferencedSecretID, byReferencedConfigID, byCustom, byCustomPrefix:
			return nil
		default:
			return ErrInvalidFindBy
		***REMOVED***
	***REMOVED***

	serviceList := []*api.Service***REMOVED******REMOVED***
	appendResult := func(o api.StoreObject) ***REMOVED***
		serviceList = append(serviceList, o.(*api.Service))
	***REMOVED***

	err := tx.find(tableService, by, checkType, appendResult)
	return serviceList, err
***REMOVED***

type serviceIndexerByRuntime struct***REMOVED******REMOVED***

func (si serviceIndexerByRuntime) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (si serviceIndexerByRuntime) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	s := obj.(*api.Service)
	r, err := naming.Runtime(s.Spec.Task)
	if err != nil ***REMOVED***
		return false, nil, nil
	***REMOVED***
	return true, []byte(r + "\x00"), nil
***REMOVED***

func (si serviceIndexerByRuntime) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return prefixFromArgs(args...)
***REMOVED***

type serviceIndexerByNetwork struct***REMOVED******REMOVED***

func (si serviceIndexerByNetwork) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (si serviceIndexerByNetwork) FromObject(obj interface***REMOVED******REMOVED***) (bool, [][]byte, error) ***REMOVED***
	s := obj.(*api.Service)

	var networkIDs [][]byte

	specNetworks := s.Spec.Task.Networks

	if len(specNetworks) == 0 ***REMOVED***
		specNetworks = s.Spec.Networks
	***REMOVED***

	for _, na := range specNetworks ***REMOVED***
		// Add the null character as a terminator
		networkIDs = append(networkIDs, []byte(na.Target+"\x00"))
	***REMOVED***

	return len(networkIDs) != 0, networkIDs, nil
***REMOVED***

type serviceIndexerBySecret struct***REMOVED******REMOVED***

func (si serviceIndexerBySecret) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (si serviceIndexerBySecret) FromObject(obj interface***REMOVED******REMOVED***) (bool, [][]byte, error) ***REMOVED***
	s := obj.(*api.Service)

	container := s.Spec.Task.GetContainer()
	if container == nil ***REMOVED***
		return false, nil, nil
	***REMOVED***

	var secretIDs [][]byte

	for _, secretRef := range container.Secrets ***REMOVED***
		// Add the null character as a terminator
		secretIDs = append(secretIDs, []byte(secretRef.SecretID+"\x00"))
	***REMOVED***

	return len(secretIDs) != 0, secretIDs, nil
***REMOVED***

type serviceIndexerByConfig struct***REMOVED******REMOVED***

func (si serviceIndexerByConfig) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (si serviceIndexerByConfig) FromObject(obj interface***REMOVED******REMOVED***) (bool, [][]byte, error) ***REMOVED***
	s, ok := obj.(*api.Service)
	if !ok ***REMOVED***
		panic("unexpected type passed to FromObject")
	***REMOVED***

	container := s.Spec.Task.GetContainer()
	if container == nil ***REMOVED***
		return false, nil, nil
	***REMOVED***

	var configIDs [][]byte

	for _, configRef := range container.Configs ***REMOVED***
		// Add the null character as a terminator
		configIDs = append(configIDs, []byte(configRef.ConfigID+"\x00"))
	***REMOVED***

	return len(configIDs) != 0, configIDs, nil
***REMOVED***
