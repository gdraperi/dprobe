package container

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

var root string

func TestMain(m *testing.M) ***REMOVED***
	var err error
	root, err = ioutil.TempDir("", "docker-container-test-")
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	defer os.RemoveAll(root)

	os.Exit(m.Run())
***REMOVED***

func newContainer(t *testing.T) *Container ***REMOVED***
	var (
		id    = uuid.New()
		cRoot = filepath.Join(root, id)
	)
	if err := os.MkdirAll(cRoot, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	c := NewBaseContainer(id, cRoot)
	c.HostConfig = &containertypes.HostConfig***REMOVED******REMOVED***
	return c
***REMOVED***

func TestViewSaveDelete(t *testing.T) ***REMOVED***
	db, err := NewViewDB()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	c := newContainer(t)
	if err := c.CheckpointTo(db); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := db.Delete(c); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestViewAll(t *testing.T) ***REMOVED***
	var (
		db, _ = NewViewDB()
		one   = newContainer(t)
		two   = newContainer(t)
	)
	one.Pid = 10
	if err := one.CheckpointTo(db); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	two.Pid = 20
	if err := two.CheckpointTo(db); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	all, err := db.Snapshot().All()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if l := len(all); l != 2 ***REMOVED***
		t.Fatalf("expected 2 items, got %d", l)
	***REMOVED***
	byID := make(map[string]Snapshot)
	for i := range all ***REMOVED***
		byID[all[i].ID] = all[i]
	***REMOVED***
	if s, ok := byID[one.ID]; !ok || s.Pid != 10 ***REMOVED***
		t.Fatalf("expected something different with for id=%s: %v", one.ID, s)
	***REMOVED***
	if s, ok := byID[two.ID]; !ok || s.Pid != 20 ***REMOVED***
		t.Fatalf("expected something different with for id=%s: %v", two.ID, s)
	***REMOVED***
***REMOVED***

func TestViewGet(t *testing.T) ***REMOVED***
	var (
		db, _ = NewViewDB()
		one   = newContainer(t)
	)
	one.ImageID = "some-image-123"
	if err := one.CheckpointTo(db); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	s, err := db.Snapshot().Get(one.ID)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if s == nil || s.ImageID != "some-image-123" ***REMOVED***
		t.Fatalf("expected ImageID=some-image-123. Got: %v", s)
	***REMOVED***
***REMOVED***

func TestNames(t *testing.T) ***REMOVED***
	db, err := NewViewDB()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assert.NoError(t, db.ReserveName("name1", "containerid1"))
	assert.NoError(t, db.ReserveName("name1", "containerid1")) // idempotent
	assert.NoError(t, db.ReserveName("name2", "containerid2"))
	assert.EqualError(t, db.ReserveName("name2", "containerid3"), ErrNameReserved.Error())

	// Releasing a name allows the name to point to something else later.
	assert.NoError(t, db.ReleaseName("name2"))
	assert.NoError(t, db.ReserveName("name2", "containerid3"))

	view := db.Snapshot()

	id, err := view.GetID("name1")
	assert.NoError(t, err)
	assert.Equal(t, "containerid1", id)

	id, err = view.GetID("name2")
	assert.NoError(t, err)
	assert.Equal(t, "containerid3", id)

	_, err = view.GetID("notreserved")
	assert.EqualError(t, err, ErrNameNotReserved.Error())

	// Releasing and re-reserving a name doesn't affect the snapshot.
	assert.NoError(t, db.ReleaseName("name2"))
	assert.NoError(t, db.ReserveName("name2", "containerid4"))

	id, err = view.GetID("name1")
	assert.NoError(t, err)
	assert.Equal(t, "containerid1", id)

	id, err = view.GetID("name2")
	assert.NoError(t, err)
	assert.Equal(t, "containerid3", id)

	// GetAllNames
	assert.Equal(t, map[string][]string***REMOVED***"containerid1": ***REMOVED***"name1"***REMOVED***, "containerid3": ***REMOVED***"name2"***REMOVED******REMOVED***, view.GetAllNames())

	assert.NoError(t, db.ReserveName("name3", "containerid1"))
	assert.NoError(t, db.ReserveName("name4", "containerid1"))

	view = db.Snapshot()
	assert.Equal(t, map[string][]string***REMOVED***"containerid1": ***REMOVED***"name1", "name3", "name4"***REMOVED***, "containerid4": ***REMOVED***"name2"***REMOVED******REMOVED***, view.GetAllNames())

	// Release containerid1's names with Delete even though no container exists
	assert.NoError(t, db.Delete(&Container***REMOVED***ID: "containerid1"***REMOVED***))

	// Reusing one of those names should work
	assert.NoError(t, db.ReserveName("name1", "containerid4"))
	view = db.Snapshot()
	assert.Equal(t, map[string][]string***REMOVED***"containerid4": ***REMOVED***"name1", "name2"***REMOVED******REMOVED***, view.GetAllNames())
***REMOVED***

// Test case for GitHub issue 35920
func TestViewWithHealthCheck(t *testing.T) ***REMOVED***
	var (
		db, _ = NewViewDB()
		one   = newContainer(t)
	)
	one.Health = &Health***REMOVED***
		Health: types.Health***REMOVED***
			Status: "starting",
		***REMOVED***,
	***REMOVED***
	if err := one.CheckpointTo(db); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	s, err := db.Snapshot().Get(one.ID)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if s == nil || s.Health != "starting" ***REMOVED***
		t.Fatalf("expected Health=starting. Got: %+v", s)
	***REMOVED***
***REMOVED***
