package query

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml"
)

func assertArrayContainsInAnyOrder(t *testing.T, array []interface***REMOVED******REMOVED***, objects ...interface***REMOVED******REMOVED***) ***REMOVED***
	if len(array) != len(objects) ***REMOVED***
		t.Fatalf("array contains %d objects but %d are expected", len(array), len(objects))
	***REMOVED***

	for _, o := range objects ***REMOVED***
		found := false
		for _, a := range array ***REMOVED***
			if a == o ***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			t.Fatal(o, "not found in array", array)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestQueryExample(t *testing.T) ***REMOVED***
	config, _ := toml.Load(`
      [[book]]
      title = "The Stand"
      author = "Stephen King"
      [[book]]
      title = "For Whom the Bell Tolls"
      author = "Ernest Hemmingway"
      [[book]]
      title = "Neuromancer"
      author = "William Gibson"
    `)
	authors, err := CompileAndExecute("$.book.author", config)
	if err != nil ***REMOVED***
		t.Fatal("unexpected error:", err)
	***REMOVED***
	names := authors.Values()
	if len(names) != 3 ***REMOVED***
		t.Fatalf("query should return 3 names but returned %d", len(names))
	***REMOVED***
	assertArrayContainsInAnyOrder(t, names, "Stephen King", "Ernest Hemmingway", "William Gibson")
***REMOVED***

func TestQueryReadmeExample(t *testing.T) ***REMOVED***
	config, _ := toml.Load(`
[postgres]
user = "pelletier"
password = "mypassword"
`)

	query, err := Compile("$..[user,password]")
	if err != nil ***REMOVED***
		t.Fatal("unexpected error:", err)
	***REMOVED***
	results := query.Execute(config)
	values := results.Values()
	if len(values) != 2 ***REMOVED***
		t.Fatalf("query should return 2 values but returned %d", len(values))
	***REMOVED***
	assertArrayContainsInAnyOrder(t, values, "pelletier", "mypassword")
***REMOVED***

func TestQueryPathNotPresent(t *testing.T) ***REMOVED***
	config, _ := toml.Load(`a = "hello"`)
	query, err := Compile("$.foo.bar")
	if err != nil ***REMOVED***
		t.Fatal("unexpected error:", err)
	***REMOVED***
	results := query.Execute(config)
	if err != nil ***REMOVED***
		t.Fatalf("err should be nil. got %s instead", err)
	***REMOVED***
	if len(results.items) != 0 ***REMOVED***
		t.Fatalf("no items should be matched. %d matched instead", len(results.items))
	***REMOVED***
***REMOVED***

func ExampleNodeFilterFn_filterExample() ***REMOVED***
	tree, _ := toml.Load(`
      [struct_one]
      foo = "foo"
      bar = "bar"

      [struct_two]
      baz = "baz"
      gorf = "gorf"
    `)

	// create a query that references a user-defined-filter
	query, _ := Compile("$[?(bazOnly)]")

	// define the filter, and assign it to the query
	query.SetFilter("bazOnly", func(node interface***REMOVED******REMOVED***) bool ***REMOVED***
		if tree, ok := node.(*toml.Tree); ok ***REMOVED***
			return tree.Has("baz")
		***REMOVED***
		return false // reject all other node types
	***REMOVED***)

	// results contain only the 'struct_two' Tree
	query.Execute(tree)
***REMOVED***

func ExampleQuery_queryExample() ***REMOVED***
	config, _ := toml.Load(`
      [[book]]
      title = "The Stand"
      author = "Stephen King"
      [[book]]
      title = "For Whom the Bell Tolls"
      author = "Ernest Hemmingway"
      [[book]]
      title = "Neuromancer"
      author = "William Gibson"
    `)

	// find and print all the authors in the document
	query, _ := Compile("$.book.author")
	authors := query.Execute(config)
	for _, name := range authors.Values() ***REMOVED***
		fmt.Println(name)
	***REMOVED***
***REMOVED***

func TestTomlQuery(t *testing.T) ***REMOVED***
	tree, err := toml.Load("[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6")
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	query, err := Compile("$.foo.bar")
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	result := query.Execute(tree)
	values := result.Values()
	if len(values) != 1 ***REMOVED***
		t.Errorf("Expected resultset of 1, got %d instead: %v", len(values), values)
	***REMOVED***

	if tt, ok := values[0].(*toml.Tree); !ok ***REMOVED***
		t.Errorf("Expected type of Tree: %T", values[0])
	***REMOVED*** else if tt.Get("a") != int64(1) ***REMOVED***
		t.Errorf("Expected 'a' with a value 1: %v", tt.Get("a"))
	***REMOVED*** else if tt.Get("b") != int64(2) ***REMOVED***
		t.Errorf("Expected 'b' with a value 2: %v", tt.Get("b"))
	***REMOVED***
***REMOVED***
