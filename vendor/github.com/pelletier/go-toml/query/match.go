package query

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

// base match
type matchBase struct ***REMOVED***
	next pathFn
***REMOVED***

func (f *matchBase) setNext(next pathFn) ***REMOVED***
	f.next = next
***REMOVED***

// terminating functor - gathers results
type terminatingFn struct ***REMOVED***
	// empty
***REMOVED***

func newTerminatingFn() *terminatingFn ***REMOVED***
	return &terminatingFn***REMOVED******REMOVED***
***REMOVED***

func (f *terminatingFn) setNext(next pathFn) ***REMOVED***
	// do nothing
***REMOVED***

func (f *terminatingFn) call(node interface***REMOVED******REMOVED***, ctx *queryContext) ***REMOVED***
	ctx.result.appendResult(node, ctx.lastPosition)
***REMOVED***

// match single key
type matchKeyFn struct ***REMOVED***
	matchBase
	Name string
***REMOVED***

func newMatchKeyFn(name string) *matchKeyFn ***REMOVED***
	return &matchKeyFn***REMOVED***Name: name***REMOVED***
***REMOVED***

func (f *matchKeyFn) call(node interface***REMOVED******REMOVED***, ctx *queryContext) ***REMOVED***
	if array, ok := node.([]*toml.Tree); ok ***REMOVED***
		for _, tree := range array ***REMOVED***
			item := tree.Get(f.Name)
			if item != nil ***REMOVED***
				ctx.lastPosition = tree.GetPosition(f.Name)
				f.next.call(item, ctx)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if tree, ok := node.(*toml.Tree); ok ***REMOVED***
		item := tree.Get(f.Name)
		if item != nil ***REMOVED***
			ctx.lastPosition = tree.GetPosition(f.Name)
			f.next.call(item, ctx)
		***REMOVED***
	***REMOVED***
***REMOVED***

// match single index
type matchIndexFn struct ***REMOVED***
	matchBase
	Idx int
***REMOVED***

func newMatchIndexFn(idx int) *matchIndexFn ***REMOVED***
	return &matchIndexFn***REMOVED***Idx: idx***REMOVED***
***REMOVED***

func (f *matchIndexFn) call(node interface***REMOVED******REMOVED***, ctx *queryContext) ***REMOVED***
	if arr, ok := node.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		if f.Idx < len(arr) && f.Idx >= 0 ***REMOVED***
			if treesArray, ok := node.([]*toml.Tree); ok ***REMOVED***
				if len(treesArray) > 0 ***REMOVED***
					ctx.lastPosition = treesArray[0].Position()
				***REMOVED***
			***REMOVED***
			f.next.call(arr[f.Idx], ctx)
		***REMOVED***
	***REMOVED***
***REMOVED***

// filter by slicing
type matchSliceFn struct ***REMOVED***
	matchBase
	Start, End, Step int
***REMOVED***

func newMatchSliceFn(start, end, step int) *matchSliceFn ***REMOVED***
	return &matchSliceFn***REMOVED***Start: start, End: end, Step: step***REMOVED***
***REMOVED***

func (f *matchSliceFn) call(node interface***REMOVED******REMOVED***, ctx *queryContext) ***REMOVED***
	if arr, ok := node.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		// adjust indexes for negative values, reverse ordering
		realStart, realEnd := f.Start, f.End
		if realStart < 0 ***REMOVED***
			realStart = len(arr) + realStart
		***REMOVED***
		if realEnd < 0 ***REMOVED***
			realEnd = len(arr) + realEnd
		***REMOVED***
		if realEnd < realStart ***REMOVED***
			realEnd, realStart = realStart, realEnd // swap
		***REMOVED***
		// loop and gather
		for idx := realStart; idx < realEnd; idx += f.Step ***REMOVED***
			if treesArray, ok := node.([]*toml.Tree); ok ***REMOVED***
				if len(treesArray) > 0 ***REMOVED***
					ctx.lastPosition = treesArray[0].Position()
				***REMOVED***
			***REMOVED***
			f.next.call(arr[idx], ctx)
		***REMOVED***
	***REMOVED***
***REMOVED***

// match anything
type matchAnyFn struct ***REMOVED***
	matchBase
***REMOVED***

func newMatchAnyFn() *matchAnyFn ***REMOVED***
	return &matchAnyFn***REMOVED******REMOVED***
***REMOVED***

func (f *matchAnyFn) call(node interface***REMOVED******REMOVED***, ctx *queryContext) ***REMOVED***
	if tree, ok := node.(*toml.Tree); ok ***REMOVED***
		for _, k := range tree.Keys() ***REMOVED***
			v := tree.Get(k)
			ctx.lastPosition = tree.GetPosition(k)
			f.next.call(v, ctx)
		***REMOVED***
	***REMOVED***
***REMOVED***

// filter through union
type matchUnionFn struct ***REMOVED***
	Union []pathFn
***REMOVED***

func (f *matchUnionFn) setNext(next pathFn) ***REMOVED***
	for _, fn := range f.Union ***REMOVED***
		fn.setNext(next)
	***REMOVED***
***REMOVED***

func (f *matchUnionFn) call(node interface***REMOVED******REMOVED***, ctx *queryContext) ***REMOVED***
	for _, fn := range f.Union ***REMOVED***
		fn.call(node, ctx)
	***REMOVED***
***REMOVED***

// match every single last node in the tree
type matchRecursiveFn struct ***REMOVED***
	matchBase
***REMOVED***

func newMatchRecursiveFn() *matchRecursiveFn ***REMOVED***
	return &matchRecursiveFn***REMOVED******REMOVED***
***REMOVED***

func (f *matchRecursiveFn) call(node interface***REMOVED******REMOVED***, ctx *queryContext) ***REMOVED***
	originalPosition := ctx.lastPosition
	if tree, ok := node.(*toml.Tree); ok ***REMOVED***
		var visit func(tree *toml.Tree)
		visit = func(tree *toml.Tree) ***REMOVED***
			for _, k := range tree.Keys() ***REMOVED***
				v := tree.Get(k)
				ctx.lastPosition = tree.GetPosition(k)
				f.next.call(v, ctx)
				switch node := v.(type) ***REMOVED***
				case *toml.Tree:
					visit(node)
				case []*toml.Tree:
					for _, subtree := range node ***REMOVED***
						visit(subtree)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		ctx.lastPosition = originalPosition
		f.next.call(tree, ctx)
		visit(tree)
	***REMOVED***
***REMOVED***

// match based on an externally provided functional filter
type matchFilterFn struct ***REMOVED***
	matchBase
	Pos  toml.Position
	Name string
***REMOVED***

func newMatchFilterFn(name string, pos toml.Position) *matchFilterFn ***REMOVED***
	return &matchFilterFn***REMOVED***Name: name, Pos: pos***REMOVED***
***REMOVED***

func (f *matchFilterFn) call(node interface***REMOVED******REMOVED***, ctx *queryContext) ***REMOVED***
	fn, ok := (*ctx.filters)[f.Name]
	if !ok ***REMOVED***
		panic(fmt.Sprintf("%s: query context does not have filter '%s'",
			f.Pos.String(), f.Name))
	***REMOVED***
	switch castNode := node.(type) ***REMOVED***
	case *toml.Tree:
		for _, k := range castNode.Keys() ***REMOVED***
			v := castNode.Get(k)
			if fn(v) ***REMOVED***
				ctx.lastPosition = castNode.GetPosition(k)
				f.next.call(v, ctx)
			***REMOVED***
		***REMOVED***
	case []*toml.Tree:
		for _, v := range castNode ***REMOVED***
			if fn(v) ***REMOVED***
				if len(castNode) > 0 ***REMOVED***
					ctx.lastPosition = castNode[0].Position()
				***REMOVED***
				f.next.call(v, ctx)
			***REMOVED***
		***REMOVED***
	case []interface***REMOVED******REMOVED***:
		for _, v := range castNode ***REMOVED***
			if fn(v) ***REMOVED***
				f.next.call(v, ctx)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
