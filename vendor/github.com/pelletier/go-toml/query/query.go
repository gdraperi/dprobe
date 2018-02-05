package query

import (
	"time"

	"github.com/pelletier/go-toml"
)

// NodeFilterFn represents a user-defined filter function, for use with
// Query.SetFilter().
//
// The return value of the function must indicate if 'node' is to be included
// at this stage of the TOML path.  Returning true will include the node, and
// returning false will exclude it.
//
// NOTE: Care should be taken to write script callbacks such that they are safe
// to use from multiple goroutines.
type NodeFilterFn func(node interface***REMOVED******REMOVED***) bool

// Result is the result of Executing a Query.
type Result struct ***REMOVED***
	items     []interface***REMOVED******REMOVED***
	positions []toml.Position
***REMOVED***

// appends a value/position pair to the result set.
func (r *Result) appendResult(node interface***REMOVED******REMOVED***, pos toml.Position) ***REMOVED***
	r.items = append(r.items, node)
	r.positions = append(r.positions, pos)
***REMOVED***

// Values is a set of values within a Result.  The order of values is not
// guaranteed to be in document order, and may be different each time a query is
// executed.
func (r Result) Values() []interface***REMOVED******REMOVED*** ***REMOVED***
	return r.items
***REMOVED***

// Positions is a set of positions for values within a Result.  Each index
// in Positions() corresponds to the entry in Value() of the same index.
func (r Result) Positions() []toml.Position ***REMOVED***
	return r.positions
***REMOVED***

// runtime context for executing query paths
type queryContext struct ***REMOVED***
	result       *Result
	filters      *map[string]NodeFilterFn
	lastPosition toml.Position
***REMOVED***

// generic path functor interface
type pathFn interface ***REMOVED***
	setNext(next pathFn)
	// it is the caller's responsibility to set the ctx.lastPosition before invoking call()
	// node can be one of: *toml.Tree, []*toml.Tree, or a scalar
	call(node interface***REMOVED******REMOVED***, ctx *queryContext)
***REMOVED***

// A Query is the representation of a compiled TOML path.  A Query is safe
// for concurrent use by multiple goroutines.
type Query struct ***REMOVED***
	root    pathFn
	tail    pathFn
	filters *map[string]NodeFilterFn
***REMOVED***

func newQuery() *Query ***REMOVED***
	return &Query***REMOVED***
		root:    nil,
		tail:    nil,
		filters: &defaultFilterFunctions,
	***REMOVED***
***REMOVED***

func (q *Query) appendPath(next pathFn) ***REMOVED***
	if q.root == nil ***REMOVED***
		q.root = next
	***REMOVED*** else ***REMOVED***
		q.tail.setNext(next)
	***REMOVED***
	q.tail = next
	next.setNext(newTerminatingFn()) // init the next functor
***REMOVED***

// Compile compiles a TOML path expression. The returned Query can be used
// to match elements within a Tree and its descendants. See Execute.
func Compile(path string) (*Query, error) ***REMOVED***
	return parseQuery(lexQuery(path))
***REMOVED***

// Execute executes a query against a Tree, and returns the result of the query.
func (q *Query) Execute(tree *toml.Tree) *Result ***REMOVED***
	result := &Result***REMOVED***
		items:     []interface***REMOVED******REMOVED******REMOVED******REMOVED***,
		positions: []toml.Position***REMOVED******REMOVED***,
	***REMOVED***
	if q.root == nil ***REMOVED***
		result.appendResult(tree, tree.GetPosition(""))
	***REMOVED*** else ***REMOVED***
		ctx := &queryContext***REMOVED***
			result:  result,
			filters: q.filters,
		***REMOVED***
		ctx.lastPosition = tree.Position()
		q.root.call(tree, ctx)
	***REMOVED***
	return result
***REMOVED***

// CompileAndExecute is a shorthand for Compile(path) followed by Execute(tree).
func CompileAndExecute(path string, tree *toml.Tree) (*Result, error) ***REMOVED***
	query, err := Compile(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return query.Execute(tree), nil
***REMOVED***

// SetFilter sets a user-defined filter function.  These may be used inside
// "?(..)" query expressions to filter TOML document elements within a query.
func (q *Query) SetFilter(name string, fn NodeFilterFn) ***REMOVED***
	if q.filters == &defaultFilterFunctions ***REMOVED***
		// clone the static table
		q.filters = &map[string]NodeFilterFn***REMOVED******REMOVED***
		for k, v := range defaultFilterFunctions ***REMOVED***
			(*q.filters)[k] = v
		***REMOVED***
	***REMOVED***
	(*q.filters)[name] = fn
***REMOVED***

var defaultFilterFunctions = map[string]NodeFilterFn***REMOVED***
	"tree": func(node interface***REMOVED******REMOVED***) bool ***REMOVED***
		_, ok := node.(*toml.Tree)
		return ok
	***REMOVED***,
	"int": func(node interface***REMOVED******REMOVED***) bool ***REMOVED***
		_, ok := node.(int64)
		return ok
	***REMOVED***,
	"float": func(node interface***REMOVED******REMOVED***) bool ***REMOVED***
		_, ok := node.(float64)
		return ok
	***REMOVED***,
	"string": func(node interface***REMOVED******REMOVED***) bool ***REMOVED***
		_, ok := node.(string)
		return ok
	***REMOVED***,
	"time": func(node interface***REMOVED******REMOVED***) bool ***REMOVED***
		_, ok := node.(time.Time)
		return ok
	***REMOVED***,
	"bool": func(node interface***REMOVED******REMOVED***) bool ***REMOVED***
		_, ok := node.(bool)
		return ok
	***REMOVED***,
***REMOVED***
