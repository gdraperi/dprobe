package memdb

import (
	"bytes"
	"fmt"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/hashicorp/go-immutable-radix"
)

const (
	id = "id"
)

// tableIndex is a tuple of (Table, Index) used for lookups
type tableIndex struct ***REMOVED***
	Table string
	Index string
***REMOVED***

// Txn is a transaction against a MemDB.
// This can be a read or write transaction.
type Txn struct ***REMOVED***
	db      *MemDB
	write   bool
	rootTxn *iradix.Txn
	after   []func()

	modified map[tableIndex]*iradix.Txn
***REMOVED***

// readableIndex returns a transaction usable for reading the given
// index in a table. If a write transaction is in progress, we may need
// to use an existing modified txn.
func (txn *Txn) readableIndex(table, index string) *iradix.Txn ***REMOVED***
	// Look for existing transaction
	if txn.write && txn.modified != nil ***REMOVED***
		key := tableIndex***REMOVED***table, index***REMOVED***
		exist, ok := txn.modified[key]
		if ok ***REMOVED***
			return exist
		***REMOVED***
	***REMOVED***

	// Create a read transaction
	path := indexPath(table, index)
	raw, _ := txn.rootTxn.Get(path)
	indexTxn := raw.(*iradix.Tree).Txn()
	return indexTxn
***REMOVED***

// writableIndex returns a transaction usable for modifying the
// given index in a table.
func (txn *Txn) writableIndex(table, index string) *iradix.Txn ***REMOVED***
	if txn.modified == nil ***REMOVED***
		txn.modified = make(map[tableIndex]*iradix.Txn)
	***REMOVED***

	// Look for existing transaction
	key := tableIndex***REMOVED***table, index***REMOVED***
	exist, ok := txn.modified[key]
	if ok ***REMOVED***
		return exist
	***REMOVED***

	// Start a new transaction
	path := indexPath(table, index)
	raw, _ := txn.rootTxn.Get(path)
	indexTxn := raw.(*iradix.Tree).Txn()

	// Keep this open for the duration of the txn
	txn.modified[key] = indexTxn
	return indexTxn
***REMOVED***

// Abort is used to cancel this transaction.
// This is a noop for read transactions.
func (txn *Txn) Abort() ***REMOVED***
	// Noop for a read transaction
	if !txn.write ***REMOVED***
		return
	***REMOVED***

	// Check if already aborted or committed
	if txn.rootTxn == nil ***REMOVED***
		return
	***REMOVED***

	// Clear the txn
	txn.rootTxn = nil
	txn.modified = nil

	// Release the writer lock since this is invalid
	txn.db.writer.Unlock()
***REMOVED***

// Commit is used to finalize this transaction.
// This is a noop for read transactions.
func (txn *Txn) Commit() ***REMOVED***
	// Noop for a read transaction
	if !txn.write ***REMOVED***
		return
	***REMOVED***

	// Check if already aborted or committed
	if txn.rootTxn == nil ***REMOVED***
		return
	***REMOVED***

	// Commit each sub-transaction scoped to (table, index)
	for key, subTxn := range txn.modified ***REMOVED***
		path := indexPath(key.Table, key.Index)
		final := subTxn.Commit()
		txn.rootTxn.Insert(path, final)
	***REMOVED***

	// Update the root of the DB
	newRoot := txn.rootTxn.Commit()
	atomic.StorePointer(&txn.db.root, unsafe.Pointer(newRoot))

	// Clear the txn
	txn.rootTxn = nil
	txn.modified = nil

	// Release the writer lock since this is invalid
	txn.db.writer.Unlock()

	// Run the deferred functions, if any
	for i := len(txn.after); i > 0; i-- ***REMOVED***
		fn := txn.after[i-1]
		fn()
	***REMOVED***
***REMOVED***

// Insert is used to add or update an object into the given table
func (txn *Txn) Insert(table string, obj interface***REMOVED******REMOVED***) error ***REMOVED***
	if !txn.write ***REMOVED***
		return fmt.Errorf("cannot insert in read-only transaction")
	***REMOVED***

	// Get the table schema
	tableSchema, ok := txn.db.schema.Tables[table]
	if !ok ***REMOVED***
		return fmt.Errorf("invalid table '%s'", table)
	***REMOVED***

	// Get the primary ID of the object
	idSchema := tableSchema.Indexes[id]
	idIndexer := idSchema.Indexer.(SingleIndexer)
	ok, idVal, err := idIndexer.FromObject(obj)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to build primary index: %v", err)
	***REMOVED***
	if !ok ***REMOVED***
		return fmt.Errorf("object missing primary index")
	***REMOVED***

	// Lookup the object by ID first, to see if this is an update
	idTxn := txn.writableIndex(table, id)
	existing, update := idTxn.Get(idVal)

	// On an update, there is an existing object with the given
	// primary ID. We do the update by deleting the current object
	// and inserting the new object.
	for name, indexSchema := range tableSchema.Indexes ***REMOVED***
		indexTxn := txn.writableIndex(table, name)

		// Determine the new index value
		var (
			ok   bool
			vals [][]byte
			err  error
		)
		switch indexer := indexSchema.Indexer.(type) ***REMOVED***
		case SingleIndexer:
			var val []byte
			ok, val, err = indexer.FromObject(obj)
			vals = [][]byte***REMOVED***val***REMOVED***
		case MultiIndexer:
			ok, vals, err = indexer.FromObject(obj)
		***REMOVED***
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to build index '%s': %v", name, err)
		***REMOVED***

		// Handle non-unique index by computing a unique index.
		// This is done by appending the primary key which must
		// be unique anyways.
		if ok && !indexSchema.Unique ***REMOVED***
			for i := range vals ***REMOVED***
				vals[i] = append(vals[i], idVal...)
			***REMOVED***
		***REMOVED***

		// Handle the update by deleting from the index first
		if update ***REMOVED***
			var (
				okExist   bool
				valsExist [][]byte
				err       error
			)
			switch indexer := indexSchema.Indexer.(type) ***REMOVED***
			case SingleIndexer:
				var valExist []byte
				okExist, valExist, err = indexer.FromObject(existing)
				valsExist = [][]byte***REMOVED***valExist***REMOVED***
			case MultiIndexer:
				okExist, valsExist, err = indexer.FromObject(existing)
			***REMOVED***
			if err != nil ***REMOVED***
				return fmt.Errorf("failed to build index '%s': %v", name, err)
			***REMOVED***
			if okExist ***REMOVED***
				for i, valExist := range valsExist ***REMOVED***
					// Handle non-unique index by computing a unique index.
					// This is done by appending the primary key which must
					// be unique anyways.
					if !indexSchema.Unique ***REMOVED***
						valExist = append(valExist, idVal...)
					***REMOVED***

					// If we are writing to the same index with the same value,
					// we can avoid the delete as the insert will overwrite the
					// value anyways.
					if i >= len(vals) || !bytes.Equal(valExist, vals[i]) ***REMOVED***
						indexTxn.Delete(valExist)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// If there is no index value, either this is an error or an expected
		// case and we can skip updating
		if !ok ***REMOVED***
			if indexSchema.AllowMissing ***REMOVED***
				continue
			***REMOVED*** else ***REMOVED***
				return fmt.Errorf("missing value for index '%s'", name)
			***REMOVED***
		***REMOVED***

		// Update the value of the index
		for _, val := range vals ***REMOVED***
			indexTxn.Insert(val, obj)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Delete is used to delete a single object from the given table
// This object must already exist in the table
func (txn *Txn) Delete(table string, obj interface***REMOVED******REMOVED***) error ***REMOVED***
	if !txn.write ***REMOVED***
		return fmt.Errorf("cannot delete in read-only transaction")
	***REMOVED***

	// Get the table schema
	tableSchema, ok := txn.db.schema.Tables[table]
	if !ok ***REMOVED***
		return fmt.Errorf("invalid table '%s'", table)
	***REMOVED***

	// Get the primary ID of the object
	idSchema := tableSchema.Indexes[id]
	idIndexer := idSchema.Indexer.(SingleIndexer)
	ok, idVal, err := idIndexer.FromObject(obj)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to build primary index: %v", err)
	***REMOVED***
	if !ok ***REMOVED***
		return fmt.Errorf("object missing primary index")
	***REMOVED***

	// Lookup the object by ID first, check fi we should continue
	idTxn := txn.writableIndex(table, id)
	existing, ok := idTxn.Get(idVal)
	if !ok ***REMOVED***
		return fmt.Errorf("not found")
	***REMOVED***

	// Remove the object from all the indexes
	for name, indexSchema := range tableSchema.Indexes ***REMOVED***
		indexTxn := txn.writableIndex(table, name)

		// Handle the update by deleting from the index first
		var (
			ok   bool
			vals [][]byte
			err  error
		)
		switch indexer := indexSchema.Indexer.(type) ***REMOVED***
		case SingleIndexer:
			var val []byte
			ok, val, err = indexer.FromObject(existing)
			vals = [][]byte***REMOVED***val***REMOVED***
		case MultiIndexer:
			ok, vals, err = indexer.FromObject(existing)
		***REMOVED***
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to build index '%s': %v", name, err)
		***REMOVED***
		if ok ***REMOVED***
			// Handle non-unique index by computing a unique index.
			// This is done by appending the primary key which must
			// be unique anyways.
			for _, val := range vals ***REMOVED***
				if !indexSchema.Unique ***REMOVED***
					val = append(val, idVal...)
				***REMOVED***
				indexTxn.Delete(val)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// DeleteAll is used to delete all the objects in a given table
// matching the constraints on the index
func (txn *Txn) DeleteAll(table, index string, args ...interface***REMOVED******REMOVED***) (int, error) ***REMOVED***
	if !txn.write ***REMOVED***
		return 0, fmt.Errorf("cannot delete in read-only transaction")
	***REMOVED***

	// Get all the objects
	iter, err := txn.Get(table, index, args...)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	// Put them into a slice so there are no safety concerns while actually
	// performing the deletes
	var objs []interface***REMOVED******REMOVED***
	for ***REMOVED***
		obj := iter.Next()
		if obj == nil ***REMOVED***
			break
		***REMOVED***

		objs = append(objs, obj)
	***REMOVED***

	// Do the deletes
	num := 0
	for _, obj := range objs ***REMOVED***
		if err := txn.Delete(table, obj); err != nil ***REMOVED***
			return num, err
		***REMOVED***
		num++
	***REMOVED***
	return num, nil
***REMOVED***

// First is used to return the first matching object for
// the given constraints on the index
func (txn *Txn) First(table, index string, args ...interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// Get the index value
	indexSchema, val, err := txn.getIndexValue(table, index, args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Get the index itself
	indexTxn := txn.readableIndex(table, indexSchema.Name)

	// Do an exact lookup
	if indexSchema.Unique && val != nil && indexSchema.Name == index ***REMOVED***
		obj, ok := indexTxn.Get(val)
		if !ok ***REMOVED***
			return nil, nil
		***REMOVED***
		return obj, nil
	***REMOVED***

	// Handle non-unique index by using an iterator and getting the first value
	iter := indexTxn.Root().Iterator()
	iter.SeekPrefix(val)
	_, value, _ := iter.Next()
	return value, nil
***REMOVED***

// LongestPrefix is used to fetch the longest prefix match for the given
// constraints on the index. Note that this will not work with the memdb
// StringFieldIndex because it adds null terminators which prevent the
// algorithm from correctly finding a match (it will get to right before the
// null and fail to find a leaf node). This should only be used where the prefix
// given is capable of matching indexed entries directly, which typically only
// applies to a custom indexer. See the unit test for an example.
func (txn *Txn) LongestPrefix(table, index string, args ...interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// Enforce that this only works on prefix indexes.
	if !strings.HasSuffix(index, "_prefix") ***REMOVED***
		return nil, fmt.Errorf("must use '%s_prefix' on index", index)
	***REMOVED***

	// Get the index value.
	indexSchema, val, err := txn.getIndexValue(table, index, args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// This algorithm only makes sense against a unique index, otherwise the
	// index keys will have the IDs appended to them.
	if !indexSchema.Unique ***REMOVED***
		return nil, fmt.Errorf("index '%s' is not unique", index)
	***REMOVED***

	// Find the longest prefix match with the given index.
	indexTxn := txn.readableIndex(table, indexSchema.Name)
	if _, value, ok := indexTxn.Root().LongestPrefix(val); ok ***REMOVED***
		return value, nil
	***REMOVED***
	return nil, nil
***REMOVED***

// getIndexValue is used to get the IndexSchema and the value
// used to scan the index given the parameters. This handles prefix based
// scans when the index has the "_prefix" suffix. The index must support
// prefix iteration.
func (txn *Txn) getIndexValue(table, index string, args ...interface***REMOVED******REMOVED***) (*IndexSchema, []byte, error) ***REMOVED***
	// Get the table schema
	tableSchema, ok := txn.db.schema.Tables[table]
	if !ok ***REMOVED***
		return nil, nil, fmt.Errorf("invalid table '%s'", table)
	***REMOVED***

	// Check for a prefix scan
	prefixScan := false
	if strings.HasSuffix(index, "_prefix") ***REMOVED***
		index = strings.TrimSuffix(index, "_prefix")
		prefixScan = true
	***REMOVED***

	// Get the index schema
	indexSchema, ok := tableSchema.Indexes[index]
	if !ok ***REMOVED***
		return nil, nil, fmt.Errorf("invalid index '%s'", index)
	***REMOVED***

	// Hot-path for when there are no arguments
	if len(args) == 0 ***REMOVED***
		return indexSchema, nil, nil
	***REMOVED***

	// Special case the prefix scanning
	if prefixScan ***REMOVED***
		prefixIndexer, ok := indexSchema.Indexer.(PrefixIndexer)
		if !ok ***REMOVED***
			return indexSchema, nil,
				fmt.Errorf("index '%s' does not support prefix scanning", index)
		***REMOVED***

		val, err := prefixIndexer.PrefixFromArgs(args...)
		if err != nil ***REMOVED***
			return indexSchema, nil, fmt.Errorf("index error: %v", err)
		***REMOVED***
		return indexSchema, val, err
	***REMOVED***

	// Get the exact match index
	val, err := indexSchema.Indexer.FromArgs(args...)
	if err != nil ***REMOVED***
		return indexSchema, nil, fmt.Errorf("index error: %v", err)
	***REMOVED***
	return indexSchema, val, err
***REMOVED***

// ResultIterator is used to iterate over a list of results
// from a Get query on a table.
type ResultIterator interface ***REMOVED***
	Next() interface***REMOVED******REMOVED***
***REMOVED***

// Get is used to construct a ResultIterator over all the
// rows that match the given constraints of an index.
func (txn *Txn) Get(table, index string, args ...interface***REMOVED******REMOVED***) (ResultIterator, error) ***REMOVED***
	// Get the index value to scan
	indexSchema, val, err := txn.getIndexValue(table, index, args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Get the index itself
	indexTxn := txn.readableIndex(table, indexSchema.Name)
	indexRoot := indexTxn.Root()

	// Get an interator over the index
	indexIter := indexRoot.Iterator()

	// Seek the iterator to the appropriate sub-set
	indexIter.SeekPrefix(val)

	// Create an iterator
	iter := &radixIterator***REMOVED***
		iter: indexIter,
	***REMOVED***
	return iter, nil
***REMOVED***

// Defer is used to push a new arbitrary function onto a stack which
// gets called when a transaction is committed and finished. Deferred
// functions are called in LIFO order, and only invoked at the end of
// write transactions.
func (txn *Txn) Defer(fn func()) ***REMOVED***
	txn.after = append(txn.after, fn)
***REMOVED***

// radixIterator is used to wrap an underlying iradix iterator.
// This is much mroe efficient than a sliceIterator as we are not
// materializing the entire view.
type radixIterator struct ***REMOVED***
	iter *iradix.Iterator
***REMOVED***

func (r *radixIterator) Next() interface***REMOVED******REMOVED*** ***REMOVED***
	_, value, ok := r.iter.Next()
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return value
***REMOVED***
