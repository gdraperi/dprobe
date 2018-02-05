package toml

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type tomlValue struct ***REMOVED***
	value     interface***REMOVED******REMOVED*** // string, int64, uint64, float64, bool, time.Time, [] of any of this list
	comment   string
	commented bool
	position  Position
***REMOVED***

// Tree is the result of the parsing of a TOML file.
type Tree struct ***REMOVED***
	values    map[string]interface***REMOVED******REMOVED*** // string -> *tomlValue, *Tree, []*Tree
	comment   string
	commented bool
	position  Position
***REMOVED***

func newTree() *Tree ***REMOVED***
	return &Tree***REMOVED***
		values:   make(map[string]interface***REMOVED******REMOVED***),
		position: Position***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// TreeFromMap initializes a new Tree object using the given map.
func TreeFromMap(m map[string]interface***REMOVED******REMOVED***) (*Tree, error) ***REMOVED***
	result, err := toTree(m)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return result.(*Tree), nil
***REMOVED***

// Position returns the position of the tree.
func (t *Tree) Position() Position ***REMOVED***
	return t.position
***REMOVED***

// Has returns a boolean indicating if the given key exists.
func (t *Tree) Has(key string) bool ***REMOVED***
	if key == "" ***REMOVED***
		return false
	***REMOVED***
	return t.HasPath(strings.Split(key, "."))
***REMOVED***

// HasPath returns true if the given path of keys exists, false otherwise.
func (t *Tree) HasPath(keys []string) bool ***REMOVED***
	return t.GetPath(keys) != nil
***REMOVED***

// Keys returns the keys of the toplevel tree (does not recurse).
func (t *Tree) Keys() []string ***REMOVED***
	keys := make([]string, len(t.values))
	i := 0
	for k := range t.values ***REMOVED***
		keys[i] = k
		i++
	***REMOVED***
	return keys
***REMOVED***

// Get the value at key in the Tree.
// Key is a dot-separated path (e.g. a.b.c) without single/double quoted strings.
// If you need to retrieve non-bare keys, use GetPath.
// Returns nil if the path does not exist in the tree.
// If keys is of length zero, the current tree is returned.
func (t *Tree) Get(key string) interface***REMOVED******REMOVED*** ***REMOVED***
	if key == "" ***REMOVED***
		return t
	***REMOVED***
	return t.GetPath(strings.Split(key, "."))
***REMOVED***

// GetPath returns the element in the tree indicated by 'keys'.
// If keys is of length zero, the current tree is returned.
func (t *Tree) GetPath(keys []string) interface***REMOVED******REMOVED*** ***REMOVED***
	if len(keys) == 0 ***REMOVED***
		return t
	***REMOVED***
	subtree := t
	for _, intermediateKey := range keys[:len(keys)-1] ***REMOVED***
		value, exists := subtree.values[intermediateKey]
		if !exists ***REMOVED***
			return nil
		***REMOVED***
		switch node := value.(type) ***REMOVED***
		case *Tree:
			subtree = node
		case []*Tree:
			// go to most recent element
			if len(node) == 0 ***REMOVED***
				return nil
			***REMOVED***
			subtree = node[len(node)-1]
		default:
			return nil // cannot navigate through other node types
		***REMOVED***
	***REMOVED***
	// branch based on final node type
	switch node := subtree.values[keys[len(keys)-1]].(type) ***REMOVED***
	case *tomlValue:
		return node.value
	default:
		return node
	***REMOVED***
***REMOVED***

// GetPosition returns the position of the given key.
func (t *Tree) GetPosition(key string) Position ***REMOVED***
	if key == "" ***REMOVED***
		return t.position
	***REMOVED***
	return t.GetPositionPath(strings.Split(key, "."))
***REMOVED***

// GetPositionPath returns the element in the tree indicated by 'keys'.
// If keys is of length zero, the current tree is returned.
func (t *Tree) GetPositionPath(keys []string) Position ***REMOVED***
	if len(keys) == 0 ***REMOVED***
		return t.position
	***REMOVED***
	subtree := t
	for _, intermediateKey := range keys[:len(keys)-1] ***REMOVED***
		value, exists := subtree.values[intermediateKey]
		if !exists ***REMOVED***
			return Position***REMOVED***0, 0***REMOVED***
		***REMOVED***
		switch node := value.(type) ***REMOVED***
		case *Tree:
			subtree = node
		case []*Tree:
			// go to most recent element
			if len(node) == 0 ***REMOVED***
				return Position***REMOVED***0, 0***REMOVED***
			***REMOVED***
			subtree = node[len(node)-1]
		default:
			return Position***REMOVED***0, 0***REMOVED***
		***REMOVED***
	***REMOVED***
	// branch based on final node type
	switch node := subtree.values[keys[len(keys)-1]].(type) ***REMOVED***
	case *tomlValue:
		return node.position
	case *Tree:
		return node.position
	case []*Tree:
		// go to most recent element
		if len(node) == 0 ***REMOVED***
			return Position***REMOVED***0, 0***REMOVED***
		***REMOVED***
		return node[len(node)-1].position
	default:
		return Position***REMOVED***0, 0***REMOVED***
	***REMOVED***
***REMOVED***

// GetDefault works like Get but with a default value
func (t *Tree) GetDefault(key string, def interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	val := t.Get(key)
	if val == nil ***REMOVED***
		return def
	***REMOVED***
	return val
***REMOVED***

// Set an element in the tree.
// Key is a dot-separated path (e.g. a.b.c).
// Creates all necessary intermediate trees, if needed.
func (t *Tree) Set(key string, value interface***REMOVED******REMOVED***) ***REMOVED***
	t.SetWithComment(key, "", false, value)
***REMOVED***

// SetWithComment is the same as Set, but allows you to provide comment
// information to the key, that will be reused by Marshal().
func (t *Tree) SetWithComment(key string, comment string, commented bool, value interface***REMOVED******REMOVED***) ***REMOVED***
	t.SetPathWithComment(strings.Split(key, "."), comment, commented, value)
***REMOVED***

// SetPath sets an element in the tree.
// Keys is an array of path elements (e.g. ***REMOVED***"a","b","c"***REMOVED***).
// Creates all necessary intermediate trees, if needed.
func (t *Tree) SetPath(keys []string, value interface***REMOVED******REMOVED***) ***REMOVED***
	t.SetPathWithComment(keys, "", false, value)
***REMOVED***

// SetPathWithComment is the same as SetPath, but allows you to provide comment
// information to the key, that will be reused by Marshal().
func (t *Tree) SetPathWithComment(keys []string, comment string, commented bool, value interface***REMOVED******REMOVED***) ***REMOVED***
	subtree := t
	for _, intermediateKey := range keys[:len(keys)-1] ***REMOVED***
		nextTree, exists := subtree.values[intermediateKey]
		if !exists ***REMOVED***
			nextTree = newTree()
			subtree.values[intermediateKey] = nextTree // add new element here
		***REMOVED***
		switch node := nextTree.(type) ***REMOVED***
		case *Tree:
			subtree = node
		case []*Tree:
			// go to most recent element
			if len(node) == 0 ***REMOVED***
				// create element if it does not exist
				subtree.values[intermediateKey] = append(node, newTree())
			***REMOVED***
			subtree = node[len(node)-1]
		***REMOVED***
	***REMOVED***

	var toInsert interface***REMOVED******REMOVED***

	switch value.(type) ***REMOVED***
	case *Tree:
		tt := value.(*Tree)
		tt.comment = comment
		toInsert = value
	case []*Tree:
		toInsert = value
	case *tomlValue:
		tt := value.(*tomlValue)
		tt.comment = comment
		toInsert = tt
	default:
		toInsert = &tomlValue***REMOVED***value: value, comment: comment, commented: commented***REMOVED***
	***REMOVED***

	subtree.values[keys[len(keys)-1]] = toInsert
***REMOVED***

// createSubTree takes a tree and a key and create the necessary intermediate
// subtrees to create a subtree at that point. In-place.
//
// e.g. passing a.b.c will create (assuming tree is empty) tree[a], tree[a][b]
// and tree[a][b][c]
//
// Returns nil on success, error object on failure
func (t *Tree) createSubTree(keys []string, pos Position) error ***REMOVED***
	subtree := t
	for _, intermediateKey := range keys ***REMOVED***
		nextTree, exists := subtree.values[intermediateKey]
		if !exists ***REMOVED***
			tree := newTree()
			tree.position = pos
			subtree.values[intermediateKey] = tree
			nextTree = tree
		***REMOVED***

		switch node := nextTree.(type) ***REMOVED***
		case []*Tree:
			subtree = node[len(node)-1]
		case *Tree:
			subtree = node
		default:
			return fmt.Errorf("unknown type for path %s (%s): %T (%#v)",
				strings.Join(keys, "."), intermediateKey, nextTree, nextTree)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// LoadBytes creates a Tree from a []byte.
func LoadBytes(b []byte) (tree *Tree, err error) ***REMOVED***
	defer func() ***REMOVED***
		if r := recover(); r != nil ***REMOVED***
			if _, ok := r.(runtime.Error); ok ***REMOVED***
				panic(r)
			***REMOVED***
			err = errors.New(r.(string))
		***REMOVED***
	***REMOVED***()
	tree = parseToml(lexToml(b))
	return
***REMOVED***

// LoadReader creates a Tree from any io.Reader.
func LoadReader(reader io.Reader) (tree *Tree, err error) ***REMOVED***
	inputBytes, err := ioutil.ReadAll(reader)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	tree, err = LoadBytes(inputBytes)
	return
***REMOVED***

// Load creates a Tree from a string.
func Load(content string) (tree *Tree, err error) ***REMOVED***
	return LoadBytes([]byte(content))
***REMOVED***

// LoadFile creates a Tree from a file.
func LoadFile(path string) (tree *Tree, err error) ***REMOVED***
	file, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer file.Close()
	return LoadReader(file)
***REMOVED***
