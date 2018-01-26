package store

import (
	"strings"

	"github.com/docker/swarmkit/api"
	memdb "github.com/hashicorp/go-memdb"
)

const tableSecret = "secret"

func init() ***REMOVED***
	register(ObjectStoreConfig***REMOVED***
		Table: &memdb.TableSchema***REMOVED***
			Name: tableSecret,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				indexID: ***REMOVED***
					Name:    indexID,
					Unique:  true,
					Indexer: api.SecretIndexerByID***REMOVED******REMOVED***,
				***REMOVED***,
				indexName: ***REMOVED***
					Name:    indexName,
					Unique:  true,
					Indexer: api.SecretIndexerByName***REMOVED******REMOVED***,
				***REMOVED***,
				indexCustom: ***REMOVED***
					Name:         indexCustom,
					Indexer:      api.SecretCustomIndexer***REMOVED******REMOVED***,
					AllowMissing: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Save: func(tx ReadTx, snapshot *api.StoreSnapshot) error ***REMOVED***
			var err error
			snapshot.Secrets, err = FindSecrets(tx, All)
			return err
		***REMOVED***,
		Restore: func(tx Tx, snapshot *api.StoreSnapshot) error ***REMOVED***
			toStoreObj := make([]api.StoreObject, len(snapshot.Secrets))
			for i, x := range snapshot.Secrets ***REMOVED***
				toStoreObj[i] = x
			***REMOVED***
			return RestoreTable(tx, tableSecret, toStoreObj)
		***REMOVED***,
		ApplyStoreAction: func(tx Tx, sa api.StoreAction) error ***REMOVED***
			switch v := sa.Target.(type) ***REMOVED***
			case *api.StoreAction_Secret:
				obj := v.Secret
				switch sa.Action ***REMOVED***
				case api.StoreActionKindCreate:
					return CreateSecret(tx, obj)
				case api.StoreActionKindUpdate:
					return UpdateSecret(tx, obj)
				case api.StoreActionKindRemove:
					return DeleteSecret(tx, obj.ID)
				***REMOVED***
			***REMOVED***
			return errUnknownStoreAction
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// CreateSecret adds a new secret to the store.
// Returns ErrExist if the ID is already taken.
func CreateSecret(tx Tx, s *api.Secret) error ***REMOVED***
	// Ensure the name is not already in use.
	if tx.lookup(tableSecret, indexName, strings.ToLower(s.Spec.Annotations.Name)) != nil ***REMOVED***
		return ErrNameConflict
	***REMOVED***

	return tx.create(tableSecret, s)
***REMOVED***

// UpdateSecret updates an existing secret in the store.
// Returns ErrNotExist if the secret doesn't exist.
func UpdateSecret(tx Tx, s *api.Secret) error ***REMOVED***
	// Ensure the name is either not in use or already used by this same Secret.
	if existing := tx.lookup(tableSecret, indexName, strings.ToLower(s.Spec.Annotations.Name)); existing != nil ***REMOVED***
		if existing.GetID() != s.ID ***REMOVED***
			return ErrNameConflict
		***REMOVED***
	***REMOVED***

	return tx.update(tableSecret, s)
***REMOVED***

// DeleteSecret removes a secret from the store.
// Returns ErrNotExist if the secret doesn't exist.
func DeleteSecret(tx Tx, id string) error ***REMOVED***
	return tx.delete(tableSecret, id)
***REMOVED***

// GetSecret looks up a secret by ID.
// Returns nil if the secret doesn't exist.
func GetSecret(tx ReadTx, id string) *api.Secret ***REMOVED***
	n := tx.get(tableSecret, id)
	if n == nil ***REMOVED***
		return nil
	***REMOVED***
	return n.(*api.Secret)
***REMOVED***

// FindSecrets selects a set of secrets and returns them.
func FindSecrets(tx ReadTx, by By) ([]*api.Secret, error) ***REMOVED***
	checkType := func(by By) error ***REMOVED***
		switch by.(type) ***REMOVED***
		case byName, byNamePrefix, byIDPrefix, byCustom, byCustomPrefix:
			return nil
		default:
			return ErrInvalidFindBy
		***REMOVED***
	***REMOVED***

	secretList := []*api.Secret***REMOVED******REMOVED***
	appendResult := func(o api.StoreObject) ***REMOVED***
		secretList = append(secretList, o.(*api.Secret))
	***REMOVED***

	err := tx.find(tableSecret, by, checkType, appendResult)
	return secretList, err
***REMOVED***
