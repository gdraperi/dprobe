package store

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/docker/go-events"
	"github.com/docker/go-metrics"
	"github.com/docker/swarmkit/api"
	pb "github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/watch"
	gogotypes "github.com/gogo/protobuf/types"
	memdb "github.com/hashicorp/go-memdb"
	"golang.org/x/net/context"
)

const (
	indexID           = "id"
	indexName         = "name"
	indexRuntime      = "runtime"
	indexServiceID    = "serviceid"
	indexNodeID       = "nodeid"
	indexSlot         = "slot"
	indexDesiredState = "desiredstate"
	indexTaskState    = "taskstate"
	indexRole         = "role"
	indexMembership   = "membership"
	indexNetwork      = "network"
	indexSecret       = "secret"
	indexConfig       = "config"
	indexKind         = "kind"
	indexCustom       = "custom"

	prefix = "_prefix"

	// MaxChangesPerTransaction is the number of changes after which a new
	// transaction should be started within Batch.
	MaxChangesPerTransaction = 200

	// MaxTransactionBytes is the maximum serialized transaction size.
	MaxTransactionBytes = 1.5 * 1024 * 1024
)

var (
	// ErrExist is returned by create operations if the provided ID is already
	// taken.
	ErrExist = errors.New("object already exists")

	// ErrNotExist is returned by altering operations (update, delete) if the
	// provided ID is not found.
	ErrNotExist = errors.New("object does not exist")

	// ErrNameConflict is returned by create/update if the object name is
	// already in use by another object.
	ErrNameConflict = errors.New("name conflicts with an existing object")

	// ErrInvalidFindBy is returned if an unrecognized type is passed to Find.
	ErrInvalidFindBy = errors.New("invalid find argument type")

	// ErrSequenceConflict is returned when trying to update an object
	// whose sequence information does not match the object in the store's.
	ErrSequenceConflict = errors.New("update out of sequence")

	objectStorers []ObjectStoreConfig
	schema        = &memdb.DBSchema***REMOVED***
		Tables: map[string]*memdb.TableSchema***REMOVED******REMOVED***,
	***REMOVED***
	errUnknownStoreAction = errors.New("unknown store action")

	// WedgeTimeout is the maximum amount of time the store lock may be
	// held before declaring a suspected deadlock.
	WedgeTimeout = 30 * time.Second

	// update()/write tx latency timer.
	updateLatencyTimer metrics.Timer

	// view()/read tx latency timer.
	viewLatencyTimer metrics.Timer

	// lookup() latency timer.
	lookupLatencyTimer metrics.Timer

	// Batch() latency timer.
	batchLatencyTimer metrics.Timer

	// timer to capture the duration for which the memory store mutex is locked.
	storeLockDurationTimer metrics.Timer
)

func init() ***REMOVED***
	ns := metrics.NewNamespace("swarm", "store", nil)
	updateLatencyTimer = ns.NewTimer("write_tx_latency",
		"Raft store write tx latency.")
	viewLatencyTimer = ns.NewTimer("read_tx_latency",
		"Raft store read tx latency.")
	lookupLatencyTimer = ns.NewTimer("lookup_latency",
		"Raft store read latency.")
	batchLatencyTimer = ns.NewTimer("batch_latency",
		"Raft store batch latency.")
	storeLockDurationTimer = ns.NewTimer("memory_store_lock_duration",
		"Duration for which the raft memory store lock was held.")
	metrics.Register(ns)
***REMOVED***

func register(os ObjectStoreConfig) ***REMOVED***
	objectStorers = append(objectStorers, os)
	schema.Tables[os.Table.Name] = os.Table
***REMOVED***

// timedMutex wraps a sync.Mutex, and keeps track of when it was locked.
type timedMutex struct ***REMOVED***
	sync.Mutex
	lockedAt atomic.Value
***REMOVED***

func (m *timedMutex) Lock() ***REMOVED***
	m.Mutex.Lock()
	m.lockedAt.Store(time.Now())
***REMOVED***

// Unlocks the timedMutex and captures the duration
// for which it was locked in a metric.
func (m *timedMutex) Unlock() ***REMOVED***
	unlockedTimestamp := m.lockedAt.Load()
	m.Mutex.Unlock()
	lockedFor := time.Since(unlockedTimestamp.(time.Time))
	storeLockDurationTimer.Update(lockedFor)
	m.lockedAt.Store(time.Time***REMOVED******REMOVED***)
***REMOVED***

func (m *timedMutex) LockedAt() time.Time ***REMOVED***
	lockedTimestamp := m.lockedAt.Load()
	if lockedTimestamp == nil ***REMOVED***
		return time.Time***REMOVED******REMOVED***
	***REMOVED***
	return lockedTimestamp.(time.Time)
***REMOVED***

// MemoryStore is a concurrency-safe, in-memory implementation of the Store
// interface.
type MemoryStore struct ***REMOVED***
	// updateLock must be held during an update transaction.
	updateLock timedMutex

	memDB *memdb.MemDB
	queue *watch.Queue

	proposer state.Proposer
***REMOVED***

// NewMemoryStore returns an in-memory store. The argument is an optional
// Proposer which will be used to propagate changes to other members in a
// cluster.
func NewMemoryStore(proposer state.Proposer) *MemoryStore ***REMOVED***
	memDB, err := memdb.NewMemDB(schema)
	if err != nil ***REMOVED***
		// This shouldn't fail
		panic(err)
	***REMOVED***

	return &MemoryStore***REMOVED***
		memDB:    memDB,
		queue:    watch.NewQueue(),
		proposer: proposer,
	***REMOVED***
***REMOVED***

// Close closes the memory store and frees its associated resources.
func (s *MemoryStore) Close() error ***REMOVED***
	return s.queue.Close()
***REMOVED***

func fromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***
	arg, ok := args[0].(string)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("argument must be a string: %#v", args[0])
	***REMOVED***
	// Add the null character as a terminator
	arg += "\x00"
	return []byte(arg), nil
***REMOVED***

func prefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	val, err := fromArgs(args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Strip the null terminator, the rest is a prefix
	n := len(val)
	if n > 0 ***REMOVED***
		return val[:n-1], nil
	***REMOVED***
	return val, nil
***REMOVED***

// ReadTx is a read transaction. Note that transaction does not imply
// any internal batching. It only means that the transaction presents a
// consistent view of the data that cannot be affected by other
// transactions.
type ReadTx interface ***REMOVED***
	lookup(table, index, id string) api.StoreObject
	get(table, id string) api.StoreObject
	find(table string, by By, checkType func(By) error, appendResult func(api.StoreObject)) error
***REMOVED***

type readTx struct ***REMOVED***
	memDBTx *memdb.Txn
***REMOVED***

// View executes a read transaction.
func (s *MemoryStore) View(cb func(ReadTx)) ***REMOVED***
	defer metrics.StartTimer(viewLatencyTimer)()
	memDBTx := s.memDB.Txn(false)

	readTx := readTx***REMOVED***
		memDBTx: memDBTx,
	***REMOVED***
	cb(readTx)
	memDBTx.Commit()
***REMOVED***

// Tx is a read/write transaction. Note that transaction does not imply
// any internal batching. The purpose of this transaction is to give the
// user a guarantee that its changes won't be visible to other transactions
// until the transaction is over.
type Tx interface ***REMOVED***
	ReadTx
	create(table string, o api.StoreObject) error
	update(table string, o api.StoreObject) error
	delete(table, id string) error
***REMOVED***

type tx struct ***REMOVED***
	readTx
	curVersion *api.Version
	changelist []api.Event
***REMOVED***

// changelistBetweenVersions returns the changes after "from" up to and
// including "to".
func (s *MemoryStore) changelistBetweenVersions(from, to api.Version) ([]api.Event, error) ***REMOVED***
	if s.proposer == nil ***REMOVED***
		return nil, errors.New("store does not support versioning")
	***REMOVED***
	changes, err := s.proposer.ChangesBetween(from, to)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var changelist []api.Event

	for _, change := range changes ***REMOVED***
		for _, sa := range change.StoreActions ***REMOVED***
			event, err := api.EventFromStoreAction(sa, nil)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			changelist = append(changelist, event)
		***REMOVED***
		changelist = append(changelist, state.EventCommit***REMOVED***Version: change.Version.Copy()***REMOVED***)
	***REMOVED***

	return changelist, nil
***REMOVED***

// ApplyStoreActions updates a store based on StoreAction messages.
func (s *MemoryStore) ApplyStoreActions(actions []api.StoreAction) error ***REMOVED***
	s.updateLock.Lock()
	memDBTx := s.memDB.Txn(true)

	tx := tx***REMOVED***
		readTx: readTx***REMOVED***
			memDBTx: memDBTx,
		***REMOVED***,
	***REMOVED***

	for _, sa := range actions ***REMOVED***
		if err := applyStoreAction(&tx, sa); err != nil ***REMOVED***
			memDBTx.Abort()
			s.updateLock.Unlock()
			return err
		***REMOVED***
	***REMOVED***

	memDBTx.Commit()

	for _, c := range tx.changelist ***REMOVED***
		s.queue.Publish(c)
	***REMOVED***
	if len(tx.changelist) != 0 ***REMOVED***
		s.queue.Publish(state.EventCommit***REMOVED******REMOVED***)
	***REMOVED***
	s.updateLock.Unlock()
	return nil
***REMOVED***

func applyStoreAction(tx Tx, sa api.StoreAction) error ***REMOVED***
	for _, os := range objectStorers ***REMOVED***
		err := os.ApplyStoreAction(tx, sa)
		if err != errUnknownStoreAction ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return errors.New("unrecognized action type")
***REMOVED***

func (s *MemoryStore) update(proposer state.Proposer, cb func(Tx) error) error ***REMOVED***
	defer metrics.StartTimer(updateLatencyTimer)()
	s.updateLock.Lock()
	memDBTx := s.memDB.Txn(true)

	var curVersion *api.Version

	if proposer != nil ***REMOVED***
		curVersion = proposer.GetVersion()
	***REMOVED***

	var tx tx
	tx.init(memDBTx, curVersion)

	err := cb(&tx)

	if err == nil ***REMOVED***
		if proposer == nil ***REMOVED***
			memDBTx.Commit()
		***REMOVED*** else ***REMOVED***
			var sa []api.StoreAction
			sa, err = tx.changelistStoreActions()

			if err == nil ***REMOVED***
				if len(sa) != 0 ***REMOVED***
					err = proposer.ProposeValue(context.Background(), sa, func() ***REMOVED***
						memDBTx.Commit()
					***REMOVED***)
				***REMOVED*** else ***REMOVED***
					memDBTx.Commit()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err == nil ***REMOVED***
		for _, c := range tx.changelist ***REMOVED***
			s.queue.Publish(c)
		***REMOVED***
		if len(tx.changelist) != 0 ***REMOVED***
			if proposer != nil ***REMOVED***
				curVersion = proposer.GetVersion()
			***REMOVED***

			s.queue.Publish(state.EventCommit***REMOVED***Version: curVersion***REMOVED***)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		memDBTx.Abort()
	***REMOVED***
	s.updateLock.Unlock()
	return err
***REMOVED***

func (s *MemoryStore) updateLocal(cb func(Tx) error) error ***REMOVED***
	return s.update(nil, cb)
***REMOVED***

// Update executes a read/write transaction.
func (s *MemoryStore) Update(cb func(Tx) error) error ***REMOVED***
	return s.update(s.proposer, cb)
***REMOVED***

// Batch provides a mechanism to batch updates to a store.
type Batch struct ***REMOVED***
	tx    tx
	store *MemoryStore
	// applied counts the times Update has run successfully
	applied int
	// transactionSizeEstimate is the running count of the size of the
	// current transaction.
	transactionSizeEstimate int
	// changelistLen is the last known length of the transaction's
	// changelist.
	changelistLen int
	err           error
***REMOVED***

// Update adds a single change to a batch. Each call to Update is atomic, but
// different calls to Update may be spread across multiple transactions to
// circumvent transaction size limits.
func (batch *Batch) Update(cb func(Tx) error) error ***REMOVED***
	if batch.err != nil ***REMOVED***
		return batch.err
	***REMOVED***

	if err := cb(&batch.tx); err != nil ***REMOVED***
		return err
	***REMOVED***

	batch.applied++

	for batch.changelistLen < len(batch.tx.changelist) ***REMOVED***
		sa, err := api.NewStoreAction(batch.tx.changelist[batch.changelistLen])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		batch.transactionSizeEstimate += sa.Size()
		batch.changelistLen++
	***REMOVED***

	if batch.changelistLen >= MaxChangesPerTransaction || batch.transactionSizeEstimate >= (MaxTransactionBytes*3)/4 ***REMOVED***
		if err := batch.commit(); err != nil ***REMOVED***
			return err
		***REMOVED***

		// Yield the update lock
		batch.store.updateLock.Unlock()
		runtime.Gosched()
		batch.store.updateLock.Lock()

		batch.newTx()
	***REMOVED***

	return nil
***REMOVED***

func (batch *Batch) newTx() ***REMOVED***
	var curVersion *api.Version

	if batch.store.proposer != nil ***REMOVED***
		curVersion = batch.store.proposer.GetVersion()
	***REMOVED***

	batch.tx.init(batch.store.memDB.Txn(true), curVersion)
	batch.transactionSizeEstimate = 0
	batch.changelistLen = 0
***REMOVED***

func (batch *Batch) commit() error ***REMOVED***
	if batch.store.proposer != nil ***REMOVED***
		var sa []api.StoreAction
		sa, batch.err = batch.tx.changelistStoreActions()

		if batch.err == nil ***REMOVED***
			if len(sa) != 0 ***REMOVED***
				batch.err = batch.store.proposer.ProposeValue(context.Background(), sa, func() ***REMOVED***
					batch.tx.memDBTx.Commit()
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				batch.tx.memDBTx.Commit()
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		batch.tx.memDBTx.Commit()
	***REMOVED***

	if batch.err != nil ***REMOVED***
		batch.tx.memDBTx.Abort()
		return batch.err
	***REMOVED***

	for _, c := range batch.tx.changelist ***REMOVED***
		batch.store.queue.Publish(c)
	***REMOVED***
	if len(batch.tx.changelist) != 0 ***REMOVED***
		batch.store.queue.Publish(state.EventCommit***REMOVED******REMOVED***)
	***REMOVED***

	return nil
***REMOVED***

// Batch performs one or more transactions that allow reads and writes
// It invokes a callback that is passed a Batch object. The callback may
// call batch.Update for each change it wants to make as part of the
// batch. The changes in the batch may be split over multiple
// transactions if necessary to keep transactions below the size limit.
// Batch holds a lock over the state, but will yield this lock every
// it creates a new transaction to allow other writers to proceed.
// Thus, unrelated changes to the state may occur between calls to
// batch.Update.
//
// This method allows the caller to iterate over a data set and apply
// changes in sequence without holding the store write lock for an
// excessive time, or producing a transaction that exceeds the maximum
// size.
//
// If Batch returns an error, no guarantees are made about how many updates
// were committed successfully.
func (s *MemoryStore) Batch(cb func(*Batch) error) error ***REMOVED***
	defer metrics.StartTimer(batchLatencyTimer)()
	s.updateLock.Lock()

	batch := Batch***REMOVED***
		store: s,
	***REMOVED***
	batch.newTx()

	if err := cb(&batch); err != nil ***REMOVED***
		batch.tx.memDBTx.Abort()
		s.updateLock.Unlock()
		return err
	***REMOVED***

	err := batch.commit()
	s.updateLock.Unlock()
	return err
***REMOVED***

func (tx *tx) init(memDBTx *memdb.Txn, curVersion *api.Version) ***REMOVED***
	tx.memDBTx = memDBTx
	tx.curVersion = curVersion
	tx.changelist = nil
***REMOVED***

func (tx tx) changelistStoreActions() ([]api.StoreAction, error) ***REMOVED***
	var actions []api.StoreAction

	for _, c := range tx.changelist ***REMOVED***
		sa, err := api.NewStoreAction(c)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		actions = append(actions, sa)
	***REMOVED***

	return actions, nil
***REMOVED***

// lookup is an internal typed wrapper around memdb.
func (tx readTx) lookup(table, index, id string) api.StoreObject ***REMOVED***
	defer metrics.StartTimer(lookupLatencyTimer)()
	j, err := tx.memDBTx.First(table, index, id)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	if j != nil ***REMOVED***
		return j.(api.StoreObject)
	***REMOVED***
	return nil
***REMOVED***

// create adds a new object to the store.
// Returns ErrExist if the ID is already taken.
func (tx *tx) create(table string, o api.StoreObject) error ***REMOVED***
	if tx.lookup(table, indexID, o.GetID()) != nil ***REMOVED***
		return ErrExist
	***REMOVED***

	copy := o.CopyStoreObject()
	meta := copy.GetMeta()
	if err := touchMeta(&meta, tx.curVersion); err != nil ***REMOVED***
		return err
	***REMOVED***
	copy.SetMeta(meta)

	err := tx.memDBTx.Insert(table, copy)
	if err == nil ***REMOVED***
		tx.changelist = append(tx.changelist, copy.EventCreate())
		o.SetMeta(meta)
	***REMOVED***
	return err
***REMOVED***

// Update updates an existing object in the store.
// Returns ErrNotExist if the object doesn't exist.
func (tx *tx) update(table string, o api.StoreObject) error ***REMOVED***
	oldN := tx.lookup(table, indexID, o.GetID())
	if oldN == nil ***REMOVED***
		return ErrNotExist
	***REMOVED***

	meta := o.GetMeta()

	if tx.curVersion != nil ***REMOVED***
		if oldN.GetMeta().Version != meta.Version ***REMOVED***
			return ErrSequenceConflict
		***REMOVED***
	***REMOVED***

	copy := o.CopyStoreObject()
	if err := touchMeta(&meta, tx.curVersion); err != nil ***REMOVED***
		return err
	***REMOVED***
	copy.SetMeta(meta)

	err := tx.memDBTx.Insert(table, copy)
	if err == nil ***REMOVED***
		tx.changelist = append(tx.changelist, copy.EventUpdate(oldN))
		o.SetMeta(meta)
	***REMOVED***
	return err
***REMOVED***

// Delete removes an object from the store.
// Returns ErrNotExist if the object doesn't exist.
func (tx *tx) delete(table, id string) error ***REMOVED***
	n := tx.lookup(table, indexID, id)
	if n == nil ***REMOVED***
		return ErrNotExist
	***REMOVED***

	err := tx.memDBTx.Delete(table, n)
	if err == nil ***REMOVED***
		tx.changelist = append(tx.changelist, n.EventDelete())
	***REMOVED***
	return err
***REMOVED***

// Get looks up an object by ID.
// Returns nil if the object doesn't exist.
func (tx readTx) get(table, id string) api.StoreObject ***REMOVED***
	o := tx.lookup(table, indexID, id)
	if o == nil ***REMOVED***
		return nil
	***REMOVED***
	return o.CopyStoreObject()
***REMOVED***

// findIterators returns a slice of iterators. The union of items from these
// iterators provides the result of the query.
func (tx readTx) findIterators(table string, by By, checkType func(By) error) ([]memdb.ResultIterator, error) ***REMOVED***
	switch by.(type) ***REMOVED***
	case byAll, orCombinator: // generic types
	default: // all other types
		if err := checkType(by); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	switch v := by.(type) ***REMOVED***
	case byAll:
		it, err := tx.memDBTx.Get(table, indexID)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case orCombinator:
		var iters []memdb.ResultIterator
		for _, subBy := range v.bys ***REMOVED***
			it, err := tx.findIterators(table, subBy, checkType)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			iters = append(iters, it...)
		***REMOVED***
		return iters, nil
	case byName:
		it, err := tx.memDBTx.Get(table, indexName, strings.ToLower(string(v)))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byIDPrefix:
		it, err := tx.memDBTx.Get(table, indexID+prefix, string(v))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byNamePrefix:
		it, err := tx.memDBTx.Get(table, indexName+prefix, strings.ToLower(string(v)))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byRuntime:
		it, err := tx.memDBTx.Get(table, indexRuntime, string(v))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byNode:
		it, err := tx.memDBTx.Get(table, indexNodeID, string(v))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byService:
		it, err := tx.memDBTx.Get(table, indexServiceID, string(v))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case bySlot:
		it, err := tx.memDBTx.Get(table, indexSlot, v.serviceID+"\x00"+strconv.FormatUint(uint64(v.slot), 10))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byDesiredState:
		it, err := tx.memDBTx.Get(table, indexDesiredState, strconv.FormatInt(int64(v), 10))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byTaskState:
		it, err := tx.memDBTx.Get(table, indexTaskState, strconv.FormatInt(int64(v), 10))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byRole:
		it, err := tx.memDBTx.Get(table, indexRole, strconv.FormatInt(int64(v), 10))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byMembership:
		it, err := tx.memDBTx.Get(table, indexMembership, strconv.FormatInt(int64(v), 10))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byReferencedNetworkID:
		it, err := tx.memDBTx.Get(table, indexNetwork, string(v))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byReferencedSecretID:
		it, err := tx.memDBTx.Get(table, indexSecret, string(v))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byReferencedConfigID:
		it, err := tx.memDBTx.Get(table, indexConfig, string(v))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byKind:
		it, err := tx.memDBTx.Get(table, indexKind, string(v))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byCustom:
		var key string
		if v.objType != "" ***REMOVED***
			key = v.objType + "|" + v.index + "|" + v.value
		***REMOVED*** else ***REMOVED***
			key = v.index + "|" + v.value
		***REMOVED***
		it, err := tx.memDBTx.Get(table, indexCustom, key)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	case byCustomPrefix:
		var key string
		if v.objType != "" ***REMOVED***
			key = v.objType + "|" + v.index + "|" + v.value
		***REMOVED*** else ***REMOVED***
			key = v.index + "|" + v.value
		***REMOVED***
		it, err := tx.memDBTx.Get(table, indexCustom+prefix, key)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return []memdb.ResultIterator***REMOVED***it***REMOVED***, nil
	default:
		return nil, ErrInvalidFindBy
	***REMOVED***
***REMOVED***

// find selects a set of objects calls a callback for each matching object.
func (tx readTx) find(table string, by By, checkType func(By) error, appendResult func(api.StoreObject)) error ***REMOVED***
	fromResultIterators := func(its ...memdb.ResultIterator) ***REMOVED***
		ids := make(map[string]struct***REMOVED******REMOVED***)
		for _, it := range its ***REMOVED***
			for ***REMOVED***
				obj := it.Next()
				if obj == nil ***REMOVED***
					break
				***REMOVED***
				o := obj.(api.StoreObject)
				id := o.GetID()
				if _, exists := ids[id]; !exists ***REMOVED***
					appendResult(o.CopyStoreObject())
					ids[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	iters, err := tx.findIterators(table, by, checkType)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	fromResultIterators(iters...)

	return nil
***REMOVED***

// Save serializes the data in the store.
func (s *MemoryStore) Save(tx ReadTx) (*pb.StoreSnapshot, error) ***REMOVED***
	var snapshot pb.StoreSnapshot
	for _, os := range objectStorers ***REMOVED***
		if err := os.Save(tx, &snapshot); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return &snapshot, nil
***REMOVED***

// Restore sets the contents of the store to the serialized data in the
// argument.
func (s *MemoryStore) Restore(snapshot *pb.StoreSnapshot) error ***REMOVED***
	return s.updateLocal(func(tx Tx) error ***REMOVED***
		for _, os := range objectStorers ***REMOVED***
			if err := os.Restore(tx, snapshot); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
***REMOVED***

// WatchQueue returns the publish/subscribe queue.
func (s *MemoryStore) WatchQueue() *watch.Queue ***REMOVED***
	return s.queue
***REMOVED***

// ViewAndWatch calls a callback which can observe the state of this
// MemoryStore. It also returns a channel that will return further events from
// this point so the snapshot can be kept up to date. The watch channel must be
// released with watch.StopWatch when it is no longer needed. The channel is
// guaranteed to get all events after the moment of the snapshot, and only
// those events.
func ViewAndWatch(store *MemoryStore, cb func(ReadTx) error, specifiers ...api.Event) (watch chan events.Event, cancel func(), err error) ***REMOVED***
	// Using Update to lock the store and guarantee consistency between
	// the watcher and the the state seen by the callback. snapshotReadTx
	// exposes this Tx as a ReadTx so the callback can't modify it.
	err = store.Update(func(tx Tx) error ***REMOVED***
		if err := cb(tx); err != nil ***REMOVED***
			return err
		***REMOVED***
		watch, cancel = state.Watch(store.WatchQueue(), specifiers...)
		return nil
	***REMOVED***)
	if watch != nil && err != nil ***REMOVED***
		cancel()
		cancel = nil
		watch = nil
	***REMOVED***
	return
***REMOVED***

// WatchFrom returns a channel that will return past events from starting
// from "version", and new events until the channel is closed. If "version"
// is nil, this function is equivalent to
//
//     state.Watch(store.WatchQueue(), specifiers...).
//
// If the log has been compacted and it's not possible to produce the exact
// set of events leading from "version" to the current state, this function
// will return an error, and the caller should re-sync.
//
// The watch channel must be released with watch.StopWatch when it is no
// longer needed.
func WatchFrom(store *MemoryStore, version *api.Version, specifiers ...api.Event) (chan events.Event, func(), error) ***REMOVED***
	if version == nil ***REMOVED***
		ch, cancel := state.Watch(store.WatchQueue(), specifiers...)
		return ch, cancel, nil
	***REMOVED***

	if store.proposer == nil ***REMOVED***
		return nil, nil, errors.New("store does not support versioning")
	***REMOVED***

	var (
		curVersion  *api.Version
		watch       chan events.Event
		cancelWatch func()
	)
	// Using Update to lock the store
	err := store.Update(func(tx Tx) error ***REMOVED***
		// Get current version
		curVersion = store.proposer.GetVersion()
		// Start the watch with the store locked so events cannot be
		// missed
		watch, cancelWatch = state.Watch(store.WatchQueue(), specifiers...)
		return nil
	***REMOVED***)
	if watch != nil && err != nil ***REMOVED***
		cancelWatch()
		return nil, nil, err
	***REMOVED***

	if curVersion == nil ***REMOVED***
		cancelWatch()
		return nil, nil, errors.New("could not get current version from store")
	***REMOVED***

	changelist, err := store.changelistBetweenVersions(*version, *curVersion)
	if err != nil ***REMOVED***
		cancelWatch()
		return nil, nil, err
	***REMOVED***

	ch := make(chan events.Event)
	stop := make(chan struct***REMOVED******REMOVED***)
	cancel := func() ***REMOVED***
		close(stop)
	***REMOVED***

	go func() ***REMOVED***
		defer cancelWatch()

		matcher := state.Matcher(specifiers...)
		for _, change := range changelist ***REMOVED***
			if matcher(change) ***REMOVED***
				select ***REMOVED***
				case ch <- change:
				case <-stop:
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***

		for ***REMOVED***
			select ***REMOVED***
			case <-stop:
				return
			case e := <-watch:
				ch <- e
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch, cancel, nil
***REMOVED***

// touchMeta updates an object's timestamps when necessary and bumps the version
// if provided.
func touchMeta(meta *api.Meta, version *api.Version) error ***REMOVED***
	// Skip meta update if version is not defined as it means we're applying
	// from raft or restoring from a snapshot.
	if version == nil ***REMOVED***
		return nil
	***REMOVED***

	now, err := gogotypes.TimestampProto(time.Now())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	meta.Version = *version

	// Updated CreatedAt if not defined
	if meta.CreatedAt == nil ***REMOVED***
		meta.CreatedAt = now
	***REMOVED***

	meta.UpdatedAt = now

	return nil
***REMOVED***

// Wedged returns true if the store lock has been held for a long time,
// possibly indicating a deadlock.
func (s *MemoryStore) Wedged() bool ***REMOVED***
	lockedAt := s.updateLock.LockedAt()
	if lockedAt.IsZero() ***REMOVED***
		return false
	***REMOVED***

	return time.Since(lockedAt) > WedgeTimeout
***REMOVED***
