// Copyright (c) 2014 The go-patricia AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package patricia

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

//------------------------------------------------------------------------------
// Trie
//------------------------------------------------------------------------------

const (
	DefaultMaxPrefixPerNode         = 10
	DefaultMaxChildrenPerSparseNode = 8
)

type (
	Prefix      []byte
	Item        interface***REMOVED******REMOVED***
	VisitorFunc func(prefix Prefix, item Item) error
)

// Trie is a generic patricia trie that allows fast retrieval of items by prefix.
// and other funky stuff.
//
// Trie is not thread-safe.
type Trie struct ***REMOVED***
	prefix Prefix
	item   Item

	maxPrefixPerNode         int
	maxChildrenPerSparseNode int

	children childList
***REMOVED***

// Public API ------------------------------------------------------------------

type Option func(*Trie)

// Trie constructor.
func NewTrie(options ...Option) *Trie ***REMOVED***
	trie := &Trie***REMOVED******REMOVED***

	for _, opt := range options ***REMOVED***
		opt(trie)
	***REMOVED***

	if trie.maxPrefixPerNode <= 0 ***REMOVED***
		trie.maxPrefixPerNode = DefaultMaxPrefixPerNode
	***REMOVED***
	if trie.maxChildrenPerSparseNode <= 0 ***REMOVED***
		trie.maxChildrenPerSparseNode = DefaultMaxChildrenPerSparseNode
	***REMOVED***

	trie.children = newSparseChildList(trie.maxChildrenPerSparseNode)
	return trie
***REMOVED***

func MaxPrefixPerNode(value int) Option ***REMOVED***
	return func(trie *Trie) ***REMOVED***
		trie.maxPrefixPerNode = value
	***REMOVED***
***REMOVED***

func MaxChildrenPerSparseNode(value int) Option ***REMOVED***
	return func(trie *Trie) ***REMOVED***
		trie.maxChildrenPerSparseNode = value
	***REMOVED***
***REMOVED***

// Item returns the item stored in the root of this trie.
func (trie *Trie) Item() Item ***REMOVED***
	return trie.item
***REMOVED***

// Insert inserts a new item into the trie using the given prefix. Insert does
// not replace existing items. It returns false if an item was already in place.
func (trie *Trie) Insert(key Prefix, item Item) (inserted bool) ***REMOVED***
	return trie.put(key, item, false)
***REMOVED***

// Set works much like Insert, but it always sets the item, possibly replacing
// the item previously inserted.
func (trie *Trie) Set(key Prefix, item Item) ***REMOVED***
	trie.put(key, item, true)
***REMOVED***

// Get returns the item located at key.
//
// This method is a bit dangerous, because Get can as well end up in an internal
// node that is not really representing any user-defined value. So when nil is
// a valid value being used, it is not possible to tell if the value was inserted
// into the tree by the user or not. A possible workaround for this is not to use
// nil interface as a valid value, even using zero value of any type is enough
// to prevent this bad behaviour.
func (trie *Trie) Get(key Prefix) (item Item) ***REMOVED***
	_, node, found, leftover := trie.findSubtree(key)
	if !found || len(leftover) != 0 ***REMOVED***
		return nil
	***REMOVED***
	return node.item
***REMOVED***

// Match returns what Get(prefix) != nil would return. The same warning as for
// Get applies here as well.
func (trie *Trie) Match(prefix Prefix) (matchedExactly bool) ***REMOVED***
	return trie.Get(prefix) != nil
***REMOVED***

// MatchSubtree returns true when there is a subtree representing extensions
// to key, that is if there are any keys in the tree which have key as prefix.
func (trie *Trie) MatchSubtree(key Prefix) (matched bool) ***REMOVED***
	_, _, matched, _ = trie.findSubtree(key)
	return
***REMOVED***

// Visit calls visitor on every node containing a non-nil item
// in alphabetical order.
//
// If an error is returned from visitor, the function stops visiting the tree
// and returns that error, unless it is a special error - SkipSubtree. In that
// case Visit skips the subtree represented by the current node and continues
// elsewhere.
func (trie *Trie) Visit(visitor VisitorFunc) error ***REMOVED***
	return trie.walk(nil, visitor)
***REMOVED***

func (trie *Trie) size() int ***REMOVED***
	n := 0

	trie.walk(nil, func(prefix Prefix, item Item) error ***REMOVED***
		n++
		return nil
	***REMOVED***)

	return n
***REMOVED***

func (trie *Trie) total() int ***REMOVED***
	return 1 + trie.children.total()
***REMOVED***

// VisitSubtree works much like Visit, but it only visits nodes matching prefix.
func (trie *Trie) VisitSubtree(prefix Prefix, visitor VisitorFunc) error ***REMOVED***
	// Nil prefix not allowed.
	if prefix == nil ***REMOVED***
		panic(ErrNilPrefix)
	***REMOVED***

	// Empty trie must be handled explicitly.
	if trie.prefix == nil ***REMOVED***
		return nil
	***REMOVED***

	// Locate the relevant subtree.
	_, root, found, leftover := trie.findSubtree(prefix)
	if !found ***REMOVED***
		return nil
	***REMOVED***
	prefix = append(prefix, leftover...)

	// Visit it.
	return root.walk(prefix, visitor)
***REMOVED***

// VisitPrefixes visits only nodes that represent prefixes of key.
// To say the obvious, returning SkipSubtree from visitor makes no sense here.
func (trie *Trie) VisitPrefixes(key Prefix, visitor VisitorFunc) error ***REMOVED***
	// Nil key not allowed.
	if key == nil ***REMOVED***
		panic(ErrNilPrefix)
	***REMOVED***

	// Empty trie must be handled explicitly.
	if trie.prefix == nil ***REMOVED***
		return nil
	***REMOVED***

	// Walk the path matching key prefixes.
	node := trie
	prefix := key
	offset := 0
	for ***REMOVED***
		// Compute what part of prefix matches.
		common := node.longestCommonPrefixLength(key)
		key = key[common:]
		offset += common

		// Partial match means that there is no subtree matching prefix.
		if common < len(node.prefix) ***REMOVED***
			return nil
		***REMOVED***

		// Call the visitor.
		if item := node.item; item != nil ***REMOVED***
			if err := visitor(prefix[:offset], item); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if len(key) == 0 ***REMOVED***
			// This node represents key, we are finished.
			return nil
		***REMOVED***

		// There is some key suffix left, move to the children.
		child := node.children.next(key[0])
		if child == nil ***REMOVED***
			// There is nowhere to continue, return.
			return nil
		***REMOVED***

		node = child
	***REMOVED***
***REMOVED***

// Delete deletes the item represented by the given prefix.
//
// True is returned if the matching node was found and deleted.
func (trie *Trie) Delete(key Prefix) (deleted bool) ***REMOVED***
	// Nil prefix not allowed.
	if key == nil ***REMOVED***
		panic(ErrNilPrefix)
	***REMOVED***

	// Empty trie must be handled explicitly.
	if trie.prefix == nil ***REMOVED***
		return false
	***REMOVED***

	// Find the relevant node.
	path, found, _ := trie.findSubtreePath(key)
	if !found ***REMOVED***
		return false
	***REMOVED***

	node := path[len(path)-1]
	var parent *Trie
	if len(path) != 1 ***REMOVED***
		parent = path[len(path)-2]
	***REMOVED***

	// If the item is already set to nil, there is nothing to do.
	if node.item == nil ***REMOVED***
		return false
	***REMOVED***

	// Delete the item.
	node.item = nil

	// Initialise i before goto.
	// Will be used later in a loop.
	i := len(path) - 1

	// In case there are some child nodes, we cannot drop the whole subtree.
	// We can try to compact nodes, though.
	if node.children.length() != 0 ***REMOVED***
		goto Compact
	***REMOVED***

	// In case we are at the root, just reset it and we are done.
	if parent == nil ***REMOVED***
		node.reset()
		return true
	***REMOVED***

	// We can drop a subtree.
	// Find the first ancestor that has its value set or it has 2 or more child nodes.
	// That will be the node where to drop the subtree at.
	for ; i >= 0; i-- ***REMOVED***
		if current := path[i]; current.item != nil || current.children.length() >= 2 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	// Handle the case when there is no such node.
	// In other words, we can reset the whole tree.
	if i == -1 ***REMOVED***
		path[0].reset()
		return true
	***REMOVED***

	// We can just remove the subtree here.
	node = path[i]
	if i == 0 ***REMOVED***
		parent = nil
	***REMOVED*** else ***REMOVED***
		parent = path[i-1]
	***REMOVED***
	// i+1 is always a valid index since i is never pointing to the last node.
	// The loop above skips at least the last node since we are sure that the item
	// is set to nil and it has no children, othewise we would be compacting instead.
	node.children.remove(path[i+1].prefix[0])

Compact:
	// The node is set to the first non-empty ancestor,
	// so try to compact since that might be possible now.
	if compacted := node.compact(); compacted != node ***REMOVED***
		if parent == nil ***REMOVED***
			*node = *compacted
		***REMOVED*** else ***REMOVED***
			parent.children.replace(node.prefix[0], compacted)
			*parent = *parent.compact()
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// DeleteSubtree finds the subtree exactly matching prefix and deletes it.
//
// True is returned if the subtree was found and deleted.
func (trie *Trie) DeleteSubtree(prefix Prefix) (deleted bool) ***REMOVED***
	// Nil prefix not allowed.
	if prefix == nil ***REMOVED***
		panic(ErrNilPrefix)
	***REMOVED***

	// Empty trie must be handled explicitly.
	if trie.prefix == nil ***REMOVED***
		return false
	***REMOVED***

	// Locate the relevant subtree.
	parent, root, found, _ := trie.findSubtree(prefix)
	if !found ***REMOVED***
		return false
	***REMOVED***

	// If we are in the root of the trie, reset the trie.
	if parent == nil ***REMOVED***
		root.reset()
		return true
	***REMOVED***

	// Otherwise remove the root node from its parent.
	parent.children.remove(root.prefix[0])
	return true
***REMOVED***

// Internal helper methods -----------------------------------------------------

func (trie *Trie) empty() bool ***REMOVED***
	return trie.item == nil && trie.children.length() == 0
***REMOVED***

func (trie *Trie) reset() ***REMOVED***
	trie.prefix = nil
	trie.children = newSparseChildList(trie.maxPrefixPerNode)
***REMOVED***

func (trie *Trie) put(key Prefix, item Item, replace bool) (inserted bool) ***REMOVED***
	// Nil prefix not allowed.
	if key == nil ***REMOVED***
		panic(ErrNilPrefix)
	***REMOVED***

	var (
		common int
		node   *Trie = trie
		child  *Trie
	)

	if node.prefix == nil ***REMOVED***
		if len(key) <= trie.maxPrefixPerNode ***REMOVED***
			node.prefix = key
			goto InsertItem
		***REMOVED***
		node.prefix = key[:trie.maxPrefixPerNode]
		key = key[trie.maxPrefixPerNode:]
		goto AppendChild
	***REMOVED***

	for ***REMOVED***
		// Compute the longest common prefix length.
		common = node.longestCommonPrefixLength(key)
		key = key[common:]

		// Only a part matches, split.
		if common < len(node.prefix) ***REMOVED***
			goto SplitPrefix
		***REMOVED***

		// common == len(node.prefix) since never (common > len(node.prefix))
		// common == len(former key) <-> 0 == len(key)
		// -> former key == node.prefix
		if len(key) == 0 ***REMOVED***
			goto InsertItem
		***REMOVED***

		// Check children for matching prefix.
		child = node.children.next(key[0])
		if child == nil ***REMOVED***
			goto AppendChild
		***REMOVED***
		node = child
	***REMOVED***

SplitPrefix:
	// Split the prefix if necessary.
	child = new(Trie)
	*child = *node
	*node = *NewTrie()
	node.prefix = child.prefix[:common]
	child.prefix = child.prefix[common:]
	child = child.compact()
	node.children = node.children.add(child)

AppendChild:
	// Keep appending children until whole prefix is inserted.
	// This loop starts with empty node.prefix that needs to be filled.
	for len(key) != 0 ***REMOVED***
		child := NewTrie()
		if len(key) <= trie.maxPrefixPerNode ***REMOVED***
			child.prefix = key
			node.children = node.children.add(child)
			node = child
			goto InsertItem
		***REMOVED*** else ***REMOVED***
			child.prefix = key[:trie.maxPrefixPerNode]
			key = key[trie.maxPrefixPerNode:]
			node.children = node.children.add(child)
			node = child
		***REMOVED***
	***REMOVED***

InsertItem:
	// Try to insert the item if possible.
	if replace || node.item == nil ***REMOVED***
		node.item = item
		return true
	***REMOVED***
	return false
***REMOVED***

func (trie *Trie) compact() *Trie ***REMOVED***
	// Only a node with a single child can be compacted.
	if trie.children.length() != 1 ***REMOVED***
		return trie
	***REMOVED***

	child := trie.children.head()

	// If any item is set, we cannot compact since we want to retain
	// the ability to do searching by key. This makes compaction less usable,
	// but that simply cannot be avoided.
	if trie.item != nil || child.item != nil ***REMOVED***
		return trie
	***REMOVED***

	// Make sure the combined prefixes fit into a single node.
	if len(trie.prefix)+len(child.prefix) > trie.maxPrefixPerNode ***REMOVED***
		return trie
	***REMOVED***

	// Concatenate the prefixes, move the items.
	child.prefix = append(trie.prefix, child.prefix...)
	if trie.item != nil ***REMOVED***
		child.item = trie.item
	***REMOVED***

	return child
***REMOVED***

func (trie *Trie) findSubtree(prefix Prefix) (parent *Trie, root *Trie, found bool, leftover Prefix) ***REMOVED***
	// Find the subtree matching prefix.
	root = trie
	for ***REMOVED***
		// Compute what part of prefix matches.
		common := root.longestCommonPrefixLength(prefix)
		prefix = prefix[common:]

		// We used up the whole prefix, subtree found.
		if len(prefix) == 0 ***REMOVED***
			found = true
			leftover = root.prefix[common:]
			return
		***REMOVED***

		// Partial match means that there is no subtree matching prefix.
		if common < len(root.prefix) ***REMOVED***
			leftover = root.prefix[common:]
			return
		***REMOVED***

		// There is some prefix left, move to the children.
		child := root.children.next(prefix[0])
		if child == nil ***REMOVED***
			// There is nowhere to continue, there is no subtree matching prefix.
			return
		***REMOVED***

		parent = root
		root = child
	***REMOVED***
***REMOVED***

func (trie *Trie) findSubtreePath(prefix Prefix) (path []*Trie, found bool, leftover Prefix) ***REMOVED***
	// Find the subtree matching prefix.
	root := trie
	var subtreePath []*Trie
	for ***REMOVED***
		// Append the current root to the path.
		subtreePath = append(subtreePath, root)

		// Compute what part of prefix matches.
		common := root.longestCommonPrefixLength(prefix)
		prefix = prefix[common:]

		// We used up the whole prefix, subtree found.
		if len(prefix) == 0 ***REMOVED***
			path = subtreePath
			found = true
			leftover = root.prefix[common:]
			return
		***REMOVED***

		// Partial match means that there is no subtree matching prefix.
		if common < len(root.prefix) ***REMOVED***
			leftover = root.prefix[common:]
			return
		***REMOVED***

		// There is some prefix left, move to the children.
		child := root.children.next(prefix[0])
		if child == nil ***REMOVED***
			// There is nowhere to continue, there is no subtree matching prefix.
			return
		***REMOVED***

		root = child
	***REMOVED***
***REMOVED***

func (trie *Trie) walk(actualRootPrefix Prefix, visitor VisitorFunc) error ***REMOVED***
	var prefix Prefix
	// Allocate a bit more space for prefix at the beginning.
	if actualRootPrefix == nil ***REMOVED***
		prefix = make(Prefix, 32+len(trie.prefix))
		copy(prefix, trie.prefix)
		prefix = prefix[:len(trie.prefix)]
	***REMOVED*** else ***REMOVED***
		prefix = make(Prefix, 32+len(actualRootPrefix))
		copy(prefix, actualRootPrefix)
		prefix = prefix[:len(actualRootPrefix)]
	***REMOVED***

	// Visit the root first. Not that this works for empty trie as well since
	// in that case item == nil && len(children) == 0.
	if trie.item != nil ***REMOVED***
		if err := visitor(prefix, trie.item); err != nil ***REMOVED***
			if err == SkipSubtree ***REMOVED***
				return nil
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Then continue to the children.
	return trie.children.walk(&prefix, visitor)
***REMOVED***

func (trie *Trie) longestCommonPrefixLength(prefix Prefix) (i int) ***REMOVED***
	for ; i < len(prefix) && i < len(trie.prefix) && prefix[i] == trie.prefix[i]; i++ ***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (trie *Trie) dump() string ***REMOVED***
	writer := &bytes.Buffer***REMOVED******REMOVED***
	trie.print(writer, 0)
	return writer.String()
***REMOVED***

func (trie *Trie) print(writer io.Writer, indent int) ***REMOVED***
	fmt.Fprintf(writer, "%s%s %v\n", strings.Repeat(" ", indent), string(trie.prefix), trie.item)
	trie.children.print(writer, indent+2)
***REMOVED***

// Errors ----------------------------------------------------------------------

var (
	SkipSubtree  = errors.New("Skip this subtree")
	ErrNilPrefix = errors.New("Nil prefix passed into a method call")
)
