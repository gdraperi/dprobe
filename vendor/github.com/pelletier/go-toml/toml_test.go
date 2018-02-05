// Testing support for go-toml

package toml

import (
	"testing"
)

func TestTomlHas(t *testing.T) ***REMOVED***
	tree, _ := Load(`
		[test]
		key = "value"
	`)

	if !tree.Has("test.key") ***REMOVED***
		t.Errorf("Has - expected test.key to exists")
	***REMOVED***

	if tree.Has("") ***REMOVED***
		t.Errorf("Should return false if the key is not provided")
	***REMOVED***
***REMOVED***

func TestTomlGet(t *testing.T) ***REMOVED***
	tree, _ := Load(`
		[test]
		key = "value"
	`)

	if tree.Get("") != tree ***REMOVED***
		t.Errorf("Get should return the tree itself when given an empty path")
	***REMOVED***

	if tree.Get("test.key") != "value" ***REMOVED***
		t.Errorf("Get should return the value")
	***REMOVED***
	if tree.Get(`\`) != nil ***REMOVED***
		t.Errorf("should return nil when the key is malformed")
	***REMOVED***
***REMOVED***

func TestTomlGetDefault(t *testing.T) ***REMOVED***
	tree, _ := Load(`
		[test]
		key = "value"
	`)

	if tree.GetDefault("", "hello") != tree ***REMOVED***
		t.Error("GetDefault should return the tree itself when given an empty path")
	***REMOVED***

	if tree.GetDefault("test.key", "hello") != "value" ***REMOVED***
		t.Error("Get should return the value")
	***REMOVED***

	if tree.GetDefault("whatever", "hello") != "hello" ***REMOVED***
		t.Error("GetDefault should return the default value if the key does not exist")
	***REMOVED***
***REMOVED***

func TestTomlHasPath(t *testing.T) ***REMOVED***
	tree, _ := Load(`
		[test]
		key = "value"
	`)

	if !tree.HasPath([]string***REMOVED***"test", "key"***REMOVED***) ***REMOVED***
		t.Errorf("HasPath - expected test.key to exists")
	***REMOVED***
***REMOVED***

func TestTomlGetPath(t *testing.T) ***REMOVED***
	node := newTree()
	//TODO: set other node data

	for idx, item := range []struct ***REMOVED***
		Path     []string
		Expected *Tree
	***REMOVED******REMOVED***
		***REMOVED*** // empty path test
			[]string***REMOVED******REMOVED***,
			node,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		result := node.GetPath(item.Path)
		if result != item.Expected ***REMOVED***
			t.Errorf("GetPath[%d] %v - expected %v, got %v instead.", idx, item.Path, item.Expected, result)
		***REMOVED***
	***REMOVED***

	tree, _ := Load("[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6")
	if tree.GetPath([]string***REMOVED***"whatever"***REMOVED***) != nil ***REMOVED***
		t.Error("GetPath should return nil when the key does not exist")
	***REMOVED***
***REMOVED***

func TestTomlFromMap(t *testing.T) ***REMOVED***
	simpleMap := map[string]interface***REMOVED******REMOVED******REMOVED***"hello": 42***REMOVED***
	tree, err := TreeFromMap(simpleMap)
	if err != nil ***REMOVED***
		t.Fatal("unexpected error:", err)
	***REMOVED***
	if tree.Get("hello") != int64(42) ***REMOVED***
		t.Fatal("hello should be 42, not", tree.Get("hello"))
	***REMOVED***
***REMOVED***
