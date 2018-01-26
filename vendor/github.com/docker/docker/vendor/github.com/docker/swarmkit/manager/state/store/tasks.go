package store

import (
	"strconv"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/naming"
	memdb "github.com/hashicorp/go-memdb"
)

const tableTask = "task"

func init() ***REMOVED***
	register(ObjectStoreConfig***REMOVED***
		Table: &memdb.TableSchema***REMOVED***
			Name: tableTask,
			Indexes: map[string]*memdb.IndexSchema***REMOVED***
				indexID: ***REMOVED***
					Name:    indexID,
					Unique:  true,
					Indexer: api.TaskIndexerByID***REMOVED******REMOVED***,
				***REMOVED***,
				indexName: ***REMOVED***
					Name:         indexName,
					AllowMissing: true,
					Indexer:      taskIndexerByName***REMOVED******REMOVED***,
				***REMOVED***,
				indexRuntime: ***REMOVED***
					Name:         indexRuntime,
					AllowMissing: true,
					Indexer:      taskIndexerByRuntime***REMOVED******REMOVED***,
				***REMOVED***,
				indexServiceID: ***REMOVED***
					Name:         indexServiceID,
					AllowMissing: true,
					Indexer:      taskIndexerByServiceID***REMOVED******REMOVED***,
				***REMOVED***,
				indexNodeID: ***REMOVED***
					Name:         indexNodeID,
					AllowMissing: true,
					Indexer:      taskIndexerByNodeID***REMOVED******REMOVED***,
				***REMOVED***,
				indexSlot: ***REMOVED***
					Name:         indexSlot,
					AllowMissing: true,
					Indexer:      taskIndexerBySlot***REMOVED******REMOVED***,
				***REMOVED***,
				indexDesiredState: ***REMOVED***
					Name:    indexDesiredState,
					Indexer: taskIndexerByDesiredState***REMOVED******REMOVED***,
				***REMOVED***,
				indexTaskState: ***REMOVED***
					Name:    indexTaskState,
					Indexer: taskIndexerByTaskState***REMOVED******REMOVED***,
				***REMOVED***,
				indexNetwork: ***REMOVED***
					Name:         indexNetwork,
					AllowMissing: true,
					Indexer:      taskIndexerByNetwork***REMOVED******REMOVED***,
				***REMOVED***,
				indexSecret: ***REMOVED***
					Name:         indexSecret,
					AllowMissing: true,
					Indexer:      taskIndexerBySecret***REMOVED******REMOVED***,
				***REMOVED***,
				indexConfig: ***REMOVED***
					Name:         indexConfig,
					AllowMissing: true,
					Indexer:      taskIndexerByConfig***REMOVED******REMOVED***,
				***REMOVED***,
				indexCustom: ***REMOVED***
					Name:         indexCustom,
					Indexer:      api.TaskCustomIndexer***REMOVED******REMOVED***,
					AllowMissing: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Save: func(tx ReadTx, snapshot *api.StoreSnapshot) error ***REMOVED***
			var err error
			snapshot.Tasks, err = FindTasks(tx, All)
			return err
		***REMOVED***,
		Restore: func(tx Tx, snapshot *api.StoreSnapshot) error ***REMOVED***
			toStoreObj := make([]api.StoreObject, len(snapshot.Tasks))
			for i, x := range snapshot.Tasks ***REMOVED***
				toStoreObj[i] = x
			***REMOVED***
			return RestoreTable(tx, tableTask, toStoreObj)
		***REMOVED***,
		ApplyStoreAction: func(tx Tx, sa api.StoreAction) error ***REMOVED***
			switch v := sa.Target.(type) ***REMOVED***
			case *api.StoreAction_Task:
				obj := v.Task
				switch sa.Action ***REMOVED***
				case api.StoreActionKindCreate:
					return CreateTask(tx, obj)
				case api.StoreActionKindUpdate:
					return UpdateTask(tx, obj)
				case api.StoreActionKindRemove:
					return DeleteTask(tx, obj.ID)
				***REMOVED***
			***REMOVED***
			return errUnknownStoreAction
		***REMOVED***,
	***REMOVED***)
***REMOVED***

// CreateTask adds a new task to the store.
// Returns ErrExist if the ID is already taken.
func CreateTask(tx Tx, t *api.Task) error ***REMOVED***
	return tx.create(tableTask, t)
***REMOVED***

// UpdateTask updates an existing task in the store.
// Returns ErrNotExist if the node doesn't exist.
func UpdateTask(tx Tx, t *api.Task) error ***REMOVED***
	return tx.update(tableTask, t)
***REMOVED***

// DeleteTask removes a task from the store.
// Returns ErrNotExist if the task doesn't exist.
func DeleteTask(tx Tx, id string) error ***REMOVED***
	return tx.delete(tableTask, id)
***REMOVED***

// GetTask looks up a task by ID.
// Returns nil if the task doesn't exist.
func GetTask(tx ReadTx, id string) *api.Task ***REMOVED***
	t := tx.get(tableTask, id)
	if t == nil ***REMOVED***
		return nil
	***REMOVED***
	return t.(*api.Task)
***REMOVED***

// FindTasks selects a set of tasks and returns them.
func FindTasks(tx ReadTx, by By) ([]*api.Task, error) ***REMOVED***
	checkType := func(by By) error ***REMOVED***
		switch by.(type) ***REMOVED***
		case byName, byNamePrefix, byIDPrefix, byRuntime, byDesiredState, byTaskState, byNode, byService, bySlot, byReferencedNetworkID, byReferencedSecretID, byReferencedConfigID, byCustom, byCustomPrefix:
			return nil
		default:
			return ErrInvalidFindBy
		***REMOVED***
	***REMOVED***

	taskList := []*api.Task***REMOVED******REMOVED***
	appendResult := func(o api.StoreObject) ***REMOVED***
		taskList = append(taskList, o.(*api.Task))
	***REMOVED***

	err := tx.find(tableTask, by, checkType, appendResult)
	return taskList, err
***REMOVED***

type taskIndexerByName struct***REMOVED******REMOVED***

func (ti taskIndexerByName) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ti taskIndexerByName) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	t := obj.(*api.Task)

	name := naming.Task(t)

	// Add the null character as a terminator
	return true, []byte(strings.ToLower(name) + "\x00"), nil
***REMOVED***

func (ti taskIndexerByName) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return prefixFromArgs(args...)
***REMOVED***

type taskIndexerByRuntime struct***REMOVED******REMOVED***

func (ti taskIndexerByRuntime) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ti taskIndexerByRuntime) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	t := obj.(*api.Task)
	r, err := naming.Runtime(t.Spec)
	if err != nil ***REMOVED***
		return false, nil, nil
	***REMOVED***
	return true, []byte(r + "\x00"), nil
***REMOVED***

func (ti taskIndexerByRuntime) PrefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return prefixFromArgs(args...)
***REMOVED***

type taskIndexerByServiceID struct***REMOVED******REMOVED***

func (ti taskIndexerByServiceID) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ti taskIndexerByServiceID) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	t := obj.(*api.Task)

	// Add the null character as a terminator
	val := t.ServiceID + "\x00"
	return true, []byte(val), nil
***REMOVED***

type taskIndexerByNodeID struct***REMOVED******REMOVED***

func (ti taskIndexerByNodeID) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ti taskIndexerByNodeID) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	t := obj.(*api.Task)

	// Add the null character as a terminator
	val := t.NodeID + "\x00"
	return true, []byte(val), nil
***REMOVED***

type taskIndexerBySlot struct***REMOVED******REMOVED***

func (ti taskIndexerBySlot) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ti taskIndexerBySlot) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	t := obj.(*api.Task)

	// Add the null character as a terminator
	val := t.ServiceID + "\x00" + strconv.FormatUint(t.Slot, 10) + "\x00"
	return true, []byte(val), nil
***REMOVED***

type taskIndexerByDesiredState struct***REMOVED******REMOVED***

func (ti taskIndexerByDesiredState) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ti taskIndexerByDesiredState) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	t := obj.(*api.Task)

	// Add the null character as a terminator
	return true, []byte(strconv.FormatInt(int64(t.DesiredState), 10) + "\x00"), nil
***REMOVED***

type taskIndexerByNetwork struct***REMOVED******REMOVED***

func (ti taskIndexerByNetwork) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ti taskIndexerByNetwork) FromObject(obj interface***REMOVED******REMOVED***) (bool, [][]byte, error) ***REMOVED***
	t := obj.(*api.Task)

	var networkIDs [][]byte

	for _, na := range t.Spec.Networks ***REMOVED***
		// Add the null character as a terminator
		networkIDs = append(networkIDs, []byte(na.Target+"\x00"))
	***REMOVED***

	return len(networkIDs) != 0, networkIDs, nil
***REMOVED***

type taskIndexerBySecret struct***REMOVED******REMOVED***

func (ti taskIndexerBySecret) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ti taskIndexerBySecret) FromObject(obj interface***REMOVED******REMOVED***) (bool, [][]byte, error) ***REMOVED***
	t := obj.(*api.Task)

	container := t.Spec.GetContainer()
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

type taskIndexerByConfig struct***REMOVED******REMOVED***

func (ti taskIndexerByConfig) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ti taskIndexerByConfig) FromObject(obj interface***REMOVED******REMOVED***) (bool, [][]byte, error) ***REMOVED***
	t, ok := obj.(*api.Task)
	if !ok ***REMOVED***
		panic("unexpected type passed to FromObject")
	***REMOVED***

	container := t.Spec.GetContainer()
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

type taskIndexerByTaskState struct***REMOVED******REMOVED***

func (ts taskIndexerByTaskState) FromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return fromArgs(args...)
***REMOVED***

func (ts taskIndexerByTaskState) FromObject(obj interface***REMOVED******REMOVED***) (bool, []byte, error) ***REMOVED***
	t := obj.(*api.Task)

	// Add the null character as a terminator
	return true, []byte(strconv.FormatInt(int64(t.Status.State), 10) + "\x00"), nil
***REMOVED***
