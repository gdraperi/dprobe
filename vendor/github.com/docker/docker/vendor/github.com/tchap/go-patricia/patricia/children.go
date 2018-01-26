// Copyright (c) 2014 The go-patricia AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package patricia

import (
	"fmt"
	"io"
	"sort"
)

type childList interface ***REMOVED***
	length() int
	head() *Trie
	add(child *Trie) childList
	remove(b byte)
	replace(b byte, child *Trie)
	next(b byte) *Trie
	walk(prefix *Prefix, visitor VisitorFunc) error
	print(w io.Writer, indent int)
	total() int
***REMOVED***

type tries []*Trie

func (t tries) Len() int ***REMOVED***
	return len(t)
***REMOVED***

func (t tries) Less(i, j int) bool ***REMOVED***
	strings := sort.StringSlice***REMOVED***string(t[i].prefix), string(t[j].prefix)***REMOVED***
	return strings.Less(0, 1)
***REMOVED***

func (t tries) Swap(i, j int) ***REMOVED***
	t[i], t[j] = t[j], t[i]
***REMOVED***

type sparseChildList struct ***REMOVED***
	children tries
***REMOVED***

func newSparseChildList(maxChildrenPerSparseNode int) childList ***REMOVED***
	return &sparseChildList***REMOVED***
		children: make(tries, 0, maxChildrenPerSparseNode),
	***REMOVED***
***REMOVED***

func (list *sparseChildList) length() int ***REMOVED***
	return len(list.children)
***REMOVED***

func (list *sparseChildList) head() *Trie ***REMOVED***
	return list.children[0]
***REMOVED***

func (list *sparseChildList) add(child *Trie) childList ***REMOVED***
	// Search for an empty spot and insert the child if possible.
	if len(list.children) != cap(list.children) ***REMOVED***
		list.children = append(list.children, child)
		return list
	***REMOVED***

	// Otherwise we have to transform to the dense list type.
	return newDenseChildList(list, child)
***REMOVED***

func (list *sparseChildList) remove(b byte) ***REMOVED***
	for i, node := range list.children ***REMOVED***
		if node.prefix[0] == b ***REMOVED***
			list.children[i] = list.children[len(list.children)-1]
			list.children[len(list.children)-1] = nil
			list.children = list.children[:len(list.children)-1]
			return
		***REMOVED***
	***REMOVED***

	// This is not supposed to be reached.
	panic("removing non-existent child")
***REMOVED***

func (list *sparseChildList) replace(b byte, child *Trie) ***REMOVED***
	// Make a consistency check.
	if p0 := child.prefix[0]; p0 != b ***REMOVED***
		panic(fmt.Errorf("child prefix mismatch: %v != %v", p0, b))
	***REMOVED***

	// Seek the child and replace it.
	for i, node := range list.children ***REMOVED***
		if node.prefix[0] == b ***REMOVED***
			list.children[i] = child
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (list *sparseChildList) next(b byte) *Trie ***REMOVED***
	for _, child := range list.children ***REMOVED***
		if child.prefix[0] == b ***REMOVED***
			return child
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (list *sparseChildList) walk(prefix *Prefix, visitor VisitorFunc) error ***REMOVED***

	sort.Sort(list.children)

	for _, child := range list.children ***REMOVED***
		*prefix = append(*prefix, child.prefix...)
		if child.item != nil ***REMOVED***
			err := visitor(*prefix, child.item)
			if err != nil ***REMOVED***
				if err == SkipSubtree ***REMOVED***
					*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
					continue
				***REMOVED***
				*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
				return err
			***REMOVED***
		***REMOVED***

		err := child.children.walk(prefix, visitor)
		*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (list *sparseChildList) total() int ***REMOVED***
	tot := 0
	for _, child := range list.children ***REMOVED***
		if child != nil ***REMOVED***
			tot = tot + child.total()
		***REMOVED***
	***REMOVED***
	return tot
***REMOVED***

func (list *sparseChildList) print(w io.Writer, indent int) ***REMOVED***
	for _, child := range list.children ***REMOVED***
		if child != nil ***REMOVED***
			child.print(w, indent)
		***REMOVED***
	***REMOVED***
***REMOVED***

type denseChildList struct ***REMOVED***
	min         int
	max         int
	numChildren int
	headIndex   int
	children    []*Trie
***REMOVED***

func newDenseChildList(list *sparseChildList, child *Trie) childList ***REMOVED***
	var (
		min int = 255
		max int = 0
	)
	for _, child := range list.children ***REMOVED***
		b := int(child.prefix[0])
		if b < min ***REMOVED***
			min = b
		***REMOVED***
		if b > max ***REMOVED***
			max = b
		***REMOVED***
	***REMOVED***

	b := int(child.prefix[0])
	if b < min ***REMOVED***
		min = b
	***REMOVED***
	if b > max ***REMOVED***
		max = b
	***REMOVED***

	children := make([]*Trie, max-min+1)
	for _, child := range list.children ***REMOVED***
		children[int(child.prefix[0])-min] = child
	***REMOVED***
	children[int(child.prefix[0])-min] = child

	return &denseChildList***REMOVED***
		min:         min,
		max:         max,
		numChildren: list.length() + 1,
		headIndex:   0,
		children:    children,
	***REMOVED***
***REMOVED***

func (list *denseChildList) length() int ***REMOVED***
	return list.numChildren
***REMOVED***

func (list *denseChildList) head() *Trie ***REMOVED***
	return list.children[list.headIndex]
***REMOVED***

func (list *denseChildList) add(child *Trie) childList ***REMOVED***
	b := int(child.prefix[0])
	var i int

	switch ***REMOVED***
	case list.min <= b && b <= list.max:
		if list.children[b-list.min] != nil ***REMOVED***
			panic("dense child list collision detected")
		***REMOVED***
		i = b - list.min
		list.children[i] = child

	case b < list.min:
		children := make([]*Trie, list.max-b+1)
		i = 0
		children[i] = child
		copy(children[list.min-b:], list.children)
		list.children = children
		list.min = b

	default: // b > list.max
		children := make([]*Trie, b-list.min+1)
		i = b - list.min
		children[i] = child
		copy(children, list.children)
		list.children = children
		list.max = b
	***REMOVED***

	list.numChildren++
	if i < list.headIndex ***REMOVED***
		list.headIndex = i
	***REMOVED***
	return list
***REMOVED***

func (list *denseChildList) remove(b byte) ***REMOVED***
	i := int(b) - list.min
	if list.children[i] == nil ***REMOVED***
		// This is not supposed to be reached.
		panic("removing non-existent child")
	***REMOVED***
	list.numChildren--
	list.children[i] = nil

	// Update head index.
	if i == list.headIndex ***REMOVED***
		for ; i < len(list.children); i++ ***REMOVED***
			if list.children[i] != nil ***REMOVED***
				list.headIndex = i
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (list *denseChildList) replace(b byte, child *Trie) ***REMOVED***
	// Make a consistency check.
	if p0 := child.prefix[0]; p0 != b ***REMOVED***
		panic(fmt.Errorf("child prefix mismatch: %v != %v", p0, b))
	***REMOVED***

	// Replace the child.
	list.children[int(b)-list.min] = child
***REMOVED***

func (list *denseChildList) next(b byte) *Trie ***REMOVED***
	i := int(b)
	if i < list.min || list.max < i ***REMOVED***
		return nil
	***REMOVED***
	return list.children[i-list.min]
***REMOVED***

func (list *denseChildList) walk(prefix *Prefix, visitor VisitorFunc) error ***REMOVED***
	for _, child := range list.children ***REMOVED***
		if child == nil ***REMOVED***
			continue
		***REMOVED***
		*prefix = append(*prefix, child.prefix...)
		if child.item != nil ***REMOVED***
			if err := visitor(*prefix, child.item); err != nil ***REMOVED***
				if err == SkipSubtree ***REMOVED***
					*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
					continue
				***REMOVED***
				*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
				return err
			***REMOVED***
		***REMOVED***

		err := child.children.walk(prefix, visitor)
		*prefix = (*prefix)[:len(*prefix)-len(child.prefix)]
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (list *denseChildList) print(w io.Writer, indent int) ***REMOVED***
	for _, child := range list.children ***REMOVED***
		if child != nil ***REMOVED***
			child.print(w, indent)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (list *denseChildList) total() int ***REMOVED***
	tot := 0
	for _, child := range list.children ***REMOVED***
		if child != nil ***REMOVED***
			tot = tot + child.total()
		***REMOVED***
	***REMOVED***
	return tot
***REMOVED***
