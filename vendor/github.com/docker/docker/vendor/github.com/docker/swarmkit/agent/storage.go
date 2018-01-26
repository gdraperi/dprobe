package agent

import (
	"github.com/boltdb/bolt"
	"github.com/docker/swarmkit/api"
	"github.com/gogo/protobuf/proto"
)

// Layout:
//
//  bucket(v1.tasks.<id>) ->
//			data (task protobuf)
//			status (task status protobuf)
//			assigned (key present)
var (
	bucketKeyStorageVersion = []byte("v1")
	bucketKeyTasks          = []byte("tasks")
	bucketKeyAssigned       = []byte("assigned")
	bucketKeyData           = []byte("data")
	bucketKeyStatus         = []byte("status")
)

// InitDB prepares a database for writing task data.
//
// Proper buckets will be created if they don't already exist.
func InitDB(db *bolt.DB) error ***REMOVED***
	return db.Update(func(tx *bolt.Tx) error ***REMOVED***
		_, err := createBucketIfNotExists(tx, bucketKeyStorageVersion, bucketKeyTasks)
		return err
	***REMOVED***)
***REMOVED***

// GetTask retrieves the task with id from the datastore.
func GetTask(tx *bolt.Tx, id string) (*api.Task, error) ***REMOVED***
	var t api.Task

	if err := withTaskBucket(tx, id, func(bkt *bolt.Bucket) error ***REMOVED***
		p := bkt.Get(bucketKeyData)
		if p == nil ***REMOVED***
			return errTaskUnknown
		***REMOVED***

		return proto.Unmarshal(p, &t)
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &t, nil
***REMOVED***

// WalkTasks walks all tasks in the datastore.
func WalkTasks(tx *bolt.Tx, fn func(task *api.Task) error) error ***REMOVED***
	bkt := getTasksBucket(tx)
	if bkt == nil ***REMOVED***
		return nil
	***REMOVED***

	return bkt.ForEach(func(k, v []byte) error ***REMOVED***
		tbkt := bkt.Bucket(k)

		p := tbkt.Get(bucketKeyData)
		var t api.Task
		if err := proto.Unmarshal(p, &t); err != nil ***REMOVED***
			return err
		***REMOVED***

		return fn(&t)
	***REMOVED***)
***REMOVED***

// TaskAssigned returns true if the task is assigned to the node.
func TaskAssigned(tx *bolt.Tx, id string) bool ***REMOVED***
	bkt := getTaskBucket(tx, id)
	if bkt == nil ***REMOVED***
		return false
	***REMOVED***

	return len(bkt.Get(bucketKeyAssigned)) > 0
***REMOVED***

// GetTaskStatus returns the current status for the task.
func GetTaskStatus(tx *bolt.Tx, id string) (*api.TaskStatus, error) ***REMOVED***
	var ts api.TaskStatus
	if err := withTaskBucket(tx, id, func(bkt *bolt.Bucket) error ***REMOVED***
		p := bkt.Get(bucketKeyStatus)
		if p == nil ***REMOVED***
			return errTaskUnknown
		***REMOVED***

		return proto.Unmarshal(p, &ts)
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &ts, nil
***REMOVED***

// WalkTaskStatus calls fn for the status of each task.
func WalkTaskStatus(tx *bolt.Tx, fn func(id string, status *api.TaskStatus) error) error ***REMOVED***
	bkt := getTasksBucket(tx)
	if bkt == nil ***REMOVED***
		return nil
	***REMOVED***

	return bkt.ForEach(func(k, v []byte) error ***REMOVED***
		tbkt := bkt.Bucket(k)

		p := tbkt.Get(bucketKeyStatus)
		var ts api.TaskStatus
		if err := proto.Unmarshal(p, &ts); err != nil ***REMOVED***
			return err
		***REMOVED***

		return fn(string(k), &ts)
	***REMOVED***)
***REMOVED***

// PutTask places the task into the database.
func PutTask(tx *bolt.Tx, task *api.Task) error ***REMOVED***
	return withCreateTaskBucketIfNotExists(tx, task.ID, func(bkt *bolt.Bucket) error ***REMOVED***
		taskCopy := *task
		taskCopy.Status = api.TaskStatus***REMOVED******REMOVED*** // blank out the status.

		p, err := proto.Marshal(&taskCopy)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return bkt.Put(bucketKeyData, p)
	***REMOVED***)
***REMOVED***

// PutTaskStatus updates the status for the task with id.
func PutTaskStatus(tx *bolt.Tx, id string, status *api.TaskStatus) error ***REMOVED***
	return withCreateTaskBucketIfNotExists(tx, id, func(bkt *bolt.Bucket) error ***REMOVED***
		p, err := proto.Marshal(status)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return bkt.Put(bucketKeyStatus, p)
	***REMOVED***)
***REMOVED***

// DeleteTask completely removes the task from the database.
func DeleteTask(tx *bolt.Tx, id string) error ***REMOVED***
	bkt := getTasksBucket(tx)
	if bkt == nil ***REMOVED***
		return nil
	***REMOVED***

	return bkt.DeleteBucket([]byte(id))
***REMOVED***

// SetTaskAssignment sets the current assignment state.
func SetTaskAssignment(tx *bolt.Tx, id string, assigned bool) error ***REMOVED***
	return withTaskBucket(tx, id, func(bkt *bolt.Bucket) error ***REMOVED***
		if assigned ***REMOVED***
			return bkt.Put(bucketKeyAssigned, []byte***REMOVED***0xFF***REMOVED***)
		***REMOVED***
		return bkt.Delete(bucketKeyAssigned)
	***REMOVED***)
***REMOVED***

func createBucketIfNotExists(tx *bolt.Tx, keys ...[]byte) (*bolt.Bucket, error) ***REMOVED***
	bkt, err := tx.CreateBucketIfNotExists(keys[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, key := range keys[1:] ***REMOVED***
		bkt, err = bkt.CreateBucketIfNotExists(key)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return bkt, nil
***REMOVED***

func withCreateTaskBucketIfNotExists(tx *bolt.Tx, id string, fn func(bkt *bolt.Bucket) error) error ***REMOVED***
	bkt, err := createBucketIfNotExists(tx, bucketKeyStorageVersion, bucketKeyTasks, []byte(id))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return fn(bkt)
***REMOVED***

func withTaskBucket(tx *bolt.Tx, id string, fn func(bkt *bolt.Bucket) error) error ***REMOVED***
	bkt := getTaskBucket(tx, id)
	if bkt == nil ***REMOVED***
		return errTaskUnknown
	***REMOVED***

	return fn(bkt)
***REMOVED***

func getTaskBucket(tx *bolt.Tx, id string) *bolt.Bucket ***REMOVED***
	return getBucket(tx, bucketKeyStorageVersion, bucketKeyTasks, []byte(id))
***REMOVED***

func getTasksBucket(tx *bolt.Tx) *bolt.Bucket ***REMOVED***
	return getBucket(tx, bucketKeyStorageVersion, bucketKeyTasks)
***REMOVED***

func getBucket(tx *bolt.Tx, keys ...[]byte) *bolt.Bucket ***REMOVED***
	bkt := tx.Bucket(keys[0])

	for _, key := range keys[1:] ***REMOVED***
		if bkt == nil ***REMOVED***
			break
		***REMOVED***
		bkt = bkt.Bucket(key)
	***REMOVED***

	return bkt
***REMOVED***
