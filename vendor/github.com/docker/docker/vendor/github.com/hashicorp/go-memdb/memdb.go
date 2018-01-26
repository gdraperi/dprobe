package memdb

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/hashicorp/go-immutable-radix"
)

// MemDB is an in-memory database. It provides a table abstraction,
// which is used to store objects (rows) with multiple indexes based
// on values. The database makes use of immutable radix trees to provide
// transactions and MVCC.
type MemDB struct ***REMOVED***
	schema *DBSchema
	root   unsafe.Pointer // *iradix.Tree underneath

	// There can only be a single writter at once
	writer sync.Mutex
***REMOVED***

// NewMemDB creates a new MemDB with the given schema
func NewMemDB(schema *DBSchema) (*MemDB, error) ***REMOVED***
	// Validate the schema
	if err := schema.Validate(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Create the MemDB
	db := &MemDB***REMOVED***
		schema: schema,
		root:   unsafe.Pointer(iradix.New()),
	***REMOVED***
	if err := db.initialize(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return db, nil
***REMOVED***

// getRoot is used to do an atomic load of the root pointer
func (db *MemDB) getRoot() *iradix.Tree ***REMOVED***
	root := (*iradix.Tree)(atomic.LoadPointer(&db.root))
	return root
***REMOVED***

// Txn is used to start a new transaction, in either read or write mode.
// There can only be a single concurrent writer, but any number of readers.
func (db *MemDB) Txn(write bool) *Txn ***REMOVED***
	if write ***REMOVED***
		db.writer.Lock()
	***REMOVED***
	txn := &Txn***REMOVED***
		db:      db,
		write:   write,
		rootTxn: db.getRoot().Txn(),
	***REMOVED***
	return txn
***REMOVED***

// Snapshot is used to capture a point-in-time snapshot
// of the database that will not be affected by any write
// operations to the existing DB.
func (db *MemDB) Snapshot() *MemDB ***REMOVED***
	clone := &MemDB***REMOVED***
		schema: db.schema,
		root:   unsafe.Pointer(db.getRoot()),
	***REMOVED***
	return clone
***REMOVED***

// initialize is used to setup the DB for use after creation
func (db *MemDB) initialize() error ***REMOVED***
	root := db.getRoot()
	for tName, tableSchema := range db.schema.Tables ***REMOVED***
		for iName, _ := range tableSchema.Indexes ***REMOVED***
			index := iradix.New()
			path := indexPath(tName, iName)
			root, _, _ = root.Insert(path, index)
		***REMOVED***
	***REMOVED***
	db.root = unsafe.Pointer(root)
	return nil
***REMOVED***

// indexPath returns the path from the root to the given table index
func indexPath(table, index string) []byte ***REMOVED***
	return []byte(table + "." + index)
***REMOVED***
