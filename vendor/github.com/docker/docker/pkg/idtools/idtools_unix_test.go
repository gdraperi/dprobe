// +build !windows

package idtools

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/gotestyourself/gotestyourself/skip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

const (
	tempUser = "tempuser"
)

type node struct ***REMOVED***
	uid int
	gid int
***REMOVED***

func TestMkdirAllAndChown(t *testing.T) ***REMOVED***
	RequiresRoot(t)
	dirName, err := ioutil.TempDir("", "mkdirall")
	if err != nil ***REMOVED***
		t.Fatalf("Couldn't create temp dir: %v", err)
	***REMOVED***
	defer os.RemoveAll(dirName)

	testTree := map[string]node***REMOVED***
		"usr":              ***REMOVED***0, 0***REMOVED***,
		"usr/bin":          ***REMOVED***0, 0***REMOVED***,
		"lib":              ***REMOVED***33, 33***REMOVED***,
		"lib/x86_64":       ***REMOVED***45, 45***REMOVED***,
		"lib/x86_64/share": ***REMOVED***1, 1***REMOVED***,
	***REMOVED***

	if err := buildTree(dirName, testTree); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// test adding a directory to a pre-existing dir; only the new dir is owned by the uid/gid
	if err := MkdirAllAndChown(filepath.Join(dirName, "usr", "share"), 0755, IDPair***REMOVED***UID: 99, GID: 99***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	testTree["usr/share"] = node***REMOVED***99, 99***REMOVED***
	verifyTree, err := readTree(dirName, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareTrees(testTree, verifyTree); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// test 2-deep new directories--both should be owned by the uid/gid pair
	if err := MkdirAllAndChown(filepath.Join(dirName, "lib", "some", "other"), 0755, IDPair***REMOVED***UID: 101, GID: 101***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	testTree["lib/some"] = node***REMOVED***101, 101***REMOVED***
	testTree["lib/some/other"] = node***REMOVED***101, 101***REMOVED***
	verifyTree, err = readTree(dirName, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareTrees(testTree, verifyTree); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// test a directory that already exists; should be chowned, but nothing else
	if err := MkdirAllAndChown(filepath.Join(dirName, "usr"), 0755, IDPair***REMOVED***UID: 102, GID: 102***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	testTree["usr"] = node***REMOVED***102, 102***REMOVED***
	verifyTree, err = readTree(dirName, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareTrees(testTree, verifyTree); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestMkdirAllAndChownNew(t *testing.T) ***REMOVED***
	RequiresRoot(t)
	dirName, err := ioutil.TempDir("", "mkdirnew")
	require.NoError(t, err)
	defer os.RemoveAll(dirName)

	testTree := map[string]node***REMOVED***
		"usr":              ***REMOVED***0, 0***REMOVED***,
		"usr/bin":          ***REMOVED***0, 0***REMOVED***,
		"lib":              ***REMOVED***33, 33***REMOVED***,
		"lib/x86_64":       ***REMOVED***45, 45***REMOVED***,
		"lib/x86_64/share": ***REMOVED***1, 1***REMOVED***,
	***REMOVED***
	require.NoError(t, buildTree(dirName, testTree))

	// test adding a directory to a pre-existing dir; only the new dir is owned by the uid/gid
	err = MkdirAllAndChownNew(filepath.Join(dirName, "usr", "share"), 0755, IDPair***REMOVED***UID: 99, GID: 99***REMOVED***)
	require.NoError(t, err)

	testTree["usr/share"] = node***REMOVED***99, 99***REMOVED***
	verifyTree, err := readTree(dirName, "")
	require.NoError(t, err)
	require.NoError(t, compareTrees(testTree, verifyTree))

	// test 2-deep new directories--both should be owned by the uid/gid pair
	err = MkdirAllAndChownNew(filepath.Join(dirName, "lib", "some", "other"), 0755, IDPair***REMOVED***UID: 101, GID: 101***REMOVED***)
	require.NoError(t, err)
	testTree["lib/some"] = node***REMOVED***101, 101***REMOVED***
	testTree["lib/some/other"] = node***REMOVED***101, 101***REMOVED***
	verifyTree, err = readTree(dirName, "")
	require.NoError(t, err)
	require.NoError(t, compareTrees(testTree, verifyTree))

	// test a directory that already exists; should NOT be chowned
	err = MkdirAllAndChownNew(filepath.Join(dirName, "usr"), 0755, IDPair***REMOVED***UID: 102, GID: 102***REMOVED***)
	require.NoError(t, err)
	verifyTree, err = readTree(dirName, "")
	require.NoError(t, err)
	require.NoError(t, compareTrees(testTree, verifyTree))
***REMOVED***

func TestMkdirAndChown(t *testing.T) ***REMOVED***
	RequiresRoot(t)
	dirName, err := ioutil.TempDir("", "mkdir")
	if err != nil ***REMOVED***
		t.Fatalf("Couldn't create temp dir: %v", err)
	***REMOVED***
	defer os.RemoveAll(dirName)

	testTree := map[string]node***REMOVED***
		"usr": ***REMOVED***0, 0***REMOVED***,
	***REMOVED***
	if err := buildTree(dirName, testTree); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// test a directory that already exists; should just chown to the requested uid/gid
	if err := MkdirAndChown(filepath.Join(dirName, "usr"), 0755, IDPair***REMOVED***UID: 99, GID: 99***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	testTree["usr"] = node***REMOVED***99, 99***REMOVED***
	verifyTree, err := readTree(dirName, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareTrees(testTree, verifyTree); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// create a subdir under a dir which doesn't exist--should fail
	if err := MkdirAndChown(filepath.Join(dirName, "usr", "bin", "subdir"), 0755, IDPair***REMOVED***UID: 102, GID: 102***REMOVED***); err == nil ***REMOVED***
		t.Fatalf("Trying to create a directory with Mkdir where the parent doesn't exist should have failed")
	***REMOVED***

	// create a subdir under an existing dir; should only change the ownership of the new subdir
	if err := MkdirAndChown(filepath.Join(dirName, "usr", "bin"), 0755, IDPair***REMOVED***UID: 102, GID: 102***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	testTree["usr/bin"] = node***REMOVED***102, 102***REMOVED***
	verifyTree, err = readTree(dirName, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := compareTrees(testTree, verifyTree); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func buildTree(base string, tree map[string]node) error ***REMOVED***
	for path, node := range tree ***REMOVED***
		fullPath := filepath.Join(base, path)
		if err := os.MkdirAll(fullPath, 0755); err != nil ***REMOVED***
			return fmt.Errorf("Couldn't create path: %s; error: %v", fullPath, err)
		***REMOVED***
		if err := os.Chown(fullPath, node.uid, node.gid); err != nil ***REMOVED***
			return fmt.Errorf("Couldn't chown path: %s; error: %v", fullPath, err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func readTree(base, root string) (map[string]node, error) ***REMOVED***
	tree := make(map[string]node)

	dirInfos, err := ioutil.ReadDir(base)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Couldn't read directory entries for %q: %v", base, err)
	***REMOVED***

	for _, info := range dirInfos ***REMOVED***
		s := &unix.Stat_t***REMOVED******REMOVED***
		if err := unix.Stat(filepath.Join(base, info.Name()), s); err != nil ***REMOVED***
			return nil, fmt.Errorf("Can't stat file %q: %v", filepath.Join(base, info.Name()), err)
		***REMOVED***
		tree[filepath.Join(root, info.Name())] = node***REMOVED***int(s.Uid), int(s.Gid)***REMOVED***
		if info.IsDir() ***REMOVED***
			// read the subdirectory
			subtree, err := readTree(filepath.Join(base, info.Name()), filepath.Join(root, info.Name()))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			for path, nodeinfo := range subtree ***REMOVED***
				tree[path] = nodeinfo
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return tree, nil
***REMOVED***

func compareTrees(left, right map[string]node) error ***REMOVED***
	if len(left) != len(right) ***REMOVED***
		return fmt.Errorf("Trees aren't the same size")
	***REMOVED***
	for path, nodeLeft := range left ***REMOVED***
		if nodeRight, ok := right[path]; ok ***REMOVED***
			if nodeRight.uid != nodeLeft.uid || nodeRight.gid != nodeLeft.gid ***REMOVED***
				// mismatch
				return fmt.Errorf("mismatched ownership for %q: expected: %d:%d, got: %d:%d", path,
					nodeLeft.uid, nodeLeft.gid, nodeRight.uid, nodeRight.gid)
			***REMOVED***
			continue
		***REMOVED***
		return fmt.Errorf("right tree didn't contain path %q", path)
	***REMOVED***
	return nil
***REMOVED***

func delUser(t *testing.T, name string) ***REMOVED***
	_, err := execCmd("userdel", name)
	assert.NoError(t, err)
***REMOVED***

func TestParseSubidFileWithNewlinesAndComments(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "parsesubid")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	fnamePath := filepath.Join(tmpDir, "testsubuid")
	fcontent := `tss:100000:65536
# empty default subuid/subgid file

dockremap:231072:65536`
	if err := ioutil.WriteFile(fnamePath, []byte(fcontent), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ranges, err := parseSubidFile(fnamePath, "dockremap")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(ranges) != 1 ***REMOVED***
		t.Fatalf("wanted 1 element in ranges, got %d instead", len(ranges))
	***REMOVED***
	if ranges[0].Start != 231072 ***REMOVED***
		t.Fatalf("wanted 231072, got %d instead", ranges[0].Start)
	***REMOVED***
	if ranges[0].Length != 65536 ***REMOVED***
		t.Fatalf("wanted 65536, got %d instead", ranges[0].Length)
	***REMOVED***
***REMOVED***

func TestGetRootUIDGID(t *testing.T) ***REMOVED***
	uidMap := []IDMap***REMOVED***
		***REMOVED***
			ContainerID: 0,
			HostID:      os.Getuid(),
			Size:        1,
		***REMOVED***,
	***REMOVED***
	gidMap := []IDMap***REMOVED***
		***REMOVED***
			ContainerID: 0,
			HostID:      os.Getgid(),
			Size:        1,
		***REMOVED***,
	***REMOVED***

	uid, gid, err := GetRootUIDGID(uidMap, gidMap)
	assert.NoError(t, err)
	assert.Equal(t, os.Getegid(), uid)
	assert.Equal(t, os.Getegid(), gid)

	uidMapError := []IDMap***REMOVED***
		***REMOVED***
			ContainerID: 1,
			HostID:      os.Getuid(),
			Size:        1,
		***REMOVED***,
	***REMOVED***
	_, _, err = GetRootUIDGID(uidMapError, gidMap)
	assert.EqualError(t, err, "Container ID 0 cannot be mapped to a host ID")
***REMOVED***

func TestToContainer(t *testing.T) ***REMOVED***
	uidMap := []IDMap***REMOVED***
		***REMOVED***
			ContainerID: 2,
			HostID:      2,
			Size:        1,
		***REMOVED***,
	***REMOVED***

	containerID, err := toContainer(2, uidMap)
	assert.NoError(t, err)
	assert.Equal(t, uidMap[0].ContainerID, containerID)
***REMOVED***

func TestNewIDMappings(t *testing.T) ***REMOVED***
	RequiresRoot(t)
	_, _, err := AddNamespaceRangesUser(tempUser)
	assert.NoError(t, err)
	defer delUser(t, tempUser)

	tempUser, err := user.Lookup(tempUser)
	assert.NoError(t, err)

	gids, err := tempUser.GroupIds()
	assert.NoError(t, err)
	group, err := user.LookupGroupId(string(gids[0]))
	assert.NoError(t, err)

	idMappings, err := NewIDMappings(tempUser.Username, group.Name)
	assert.NoError(t, err)

	rootUID, rootGID, err := GetRootUIDGID(idMappings.UIDs(), idMappings.GIDs())
	assert.NoError(t, err)

	dirName, err := ioutil.TempDir("", "mkdirall")
	assert.NoError(t, err, "Couldn't create temp directory")
	defer os.RemoveAll(dirName)

	err = MkdirAllAndChown(dirName, 0700, IDPair***REMOVED***UID: rootUID, GID: rootGID***REMOVED***)
	assert.NoError(t, err, "Couldn't change ownership of file path. Got error")
	assert.True(t, CanAccess(dirName, idMappings.RootPair()), fmt.Sprintf("Unable to access %s directory with user UID:%d and GID:%d", dirName, rootUID, rootGID))
***REMOVED***

func TestLookupUserAndGroup(t *testing.T) ***REMOVED***
	RequiresRoot(t)
	uid, gid, err := AddNamespaceRangesUser(tempUser)
	assert.NoError(t, err)
	defer delUser(t, tempUser)

	fetchedUser, err := LookupUser(tempUser)
	assert.NoError(t, err)

	fetchedUserByID, err := LookupUID(uid)
	assert.NoError(t, err)
	assert.Equal(t, fetchedUserByID, fetchedUser)

	fetchedGroup, err := LookupGroup(tempUser)
	assert.NoError(t, err)

	fetchedGroupByID, err := LookupGID(gid)
	assert.NoError(t, err)
	assert.Equal(t, fetchedGroupByID, fetchedGroup)
***REMOVED***

func TestLookupUserAndGroupThatDoesNotExist(t *testing.T) ***REMOVED***
	fakeUser := "fakeuser"
	_, err := LookupUser(fakeUser)
	assert.EqualError(t, err, "getent unable to find entry \""+fakeUser+"\" in passwd database")

	_, err = LookupUID(-1)
	assert.Error(t, err)

	fakeGroup := "fakegroup"
	_, err = LookupGroup(fakeGroup)
	assert.EqualError(t, err, "getent unable to find entry \""+fakeGroup+"\" in group database")

	_, err = LookupGID(-1)
	assert.Error(t, err)
***REMOVED***

// TestMkdirIsNotDir checks that mkdirAs() function (used by MkdirAll...)
// returns a correct error in case a directory which it is about to create
// already exists but is a file (rather than a directory).
func TestMkdirIsNotDir(t *testing.T) ***REMOVED***
	file, err := ioutil.TempFile("", t.Name())
	if err != nil ***REMOVED***
		t.Fatalf("Couldn't create temp dir: %v", err)
	***REMOVED***
	defer os.Remove(file.Name())

	err = mkdirAs(file.Name(), 0755, 0, 0, false, false)
	assert.EqualError(t, err, "mkdir "+file.Name()+": not a directory")
***REMOVED***

func RequiresRoot(t *testing.T) ***REMOVED***
	skip.IfCondition(t, os.Getuid() != 0, "skipping test that requires root")
***REMOVED***
