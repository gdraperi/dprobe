package query

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"testing"
)

// dump path tree to a string
func pathString(root pathFn) string ***REMOVED***
	result := fmt.Sprintf("%T:", root)
	switch fn := root.(type) ***REMOVED***
	case *terminatingFn:
		result += "***REMOVED******REMOVED***"
	case *matchKeyFn:
		result += fmt.Sprintf("***REMOVED***%s***REMOVED***", fn.Name)
		result += pathString(fn.next)
	case *matchIndexFn:
		result += fmt.Sprintf("***REMOVED***%d***REMOVED***", fn.Idx)
		result += pathString(fn.next)
	case *matchSliceFn:
		result += fmt.Sprintf("***REMOVED***%d:%d:%d***REMOVED***",
			fn.Start, fn.End, fn.Step)
		result += pathString(fn.next)
	case *matchAnyFn:
		result += "***REMOVED******REMOVED***"
		result += pathString(fn.next)
	case *matchUnionFn:
		result += "***REMOVED***["
		for _, v := range fn.Union ***REMOVED***
			result += pathString(v) + ", "
		***REMOVED***
		result += "]***REMOVED***"
	case *matchRecursiveFn:
		result += "***REMOVED******REMOVED***"
		result += pathString(fn.next)
	case *matchFilterFn:
		result += fmt.Sprintf("***REMOVED***%s***REMOVED***", fn.Name)
		result += pathString(fn.next)
	***REMOVED***
	return result
***REMOVED***

func assertPathMatch(t *testing.T, path, ref *Query) bool ***REMOVED***
	pathStr := pathString(path.root)
	refStr := pathString(ref.root)
	if pathStr != refStr ***REMOVED***
		t.Errorf("paths do not match")
		t.Log("test:", pathStr)
		t.Log("ref: ", refStr)
		return false
	***REMOVED***
	return true
***REMOVED***

func assertPath(t *testing.T, query string, ref *Query) ***REMOVED***
	path, _ := parseQuery(lexQuery(query))
	assertPathMatch(t, path, ref)
***REMOVED***

func buildPath(parts ...pathFn) *Query ***REMOVED***
	query := newQuery()
	for _, v := range parts ***REMOVED***
		query.appendPath(v)
	***REMOVED***
	return query
***REMOVED***

func TestPathRoot(t *testing.T) ***REMOVED***
	assertPath(t,
		"$",
		buildPath(
		// empty
		))
***REMOVED***

func TestPathKey(t *testing.T) ***REMOVED***
	assertPath(t,
		"$.foo",
		buildPath(
			newMatchKeyFn("foo"),
		))
***REMOVED***

func TestPathBracketKey(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[foo]",
		buildPath(
			newMatchKeyFn("foo"),
		))
***REMOVED***

func TestPathBracketStringKey(t *testing.T) ***REMOVED***
	assertPath(t,
		"$['foo']",
		buildPath(
			newMatchKeyFn("foo"),
		))
***REMOVED***

func TestPathIndex(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[123]",
		buildPath(
			newMatchIndexFn(123),
		))
***REMOVED***

func TestPathSliceStart(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[123:]",
		buildPath(
			newMatchSliceFn(123, maxInt, 1),
		))
***REMOVED***

func TestPathSliceStartEnd(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[123:456]",
		buildPath(
			newMatchSliceFn(123, 456, 1),
		))
***REMOVED***

func TestPathSliceStartEndColon(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[123:456:]",
		buildPath(
			newMatchSliceFn(123, 456, 1),
		))
***REMOVED***

func TestPathSliceStartStep(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[123::7]",
		buildPath(
			newMatchSliceFn(123, maxInt, 7),
		))
***REMOVED***

func TestPathSliceEndStep(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[:456:7]",
		buildPath(
			newMatchSliceFn(0, 456, 7),
		))
***REMOVED***

func TestPathSliceStep(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[::7]",
		buildPath(
			newMatchSliceFn(0, maxInt, 7),
		))
***REMOVED***

func TestPathSliceAll(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[123:456:7]",
		buildPath(
			newMatchSliceFn(123, 456, 7),
		))
***REMOVED***

func TestPathAny(t *testing.T) ***REMOVED***
	assertPath(t,
		"$.*",
		buildPath(
			newMatchAnyFn(),
		))
***REMOVED***

func TestPathUnion(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[foo, bar, baz]",
		buildPath(
			&matchUnionFn***REMOVED***[]pathFn***REMOVED***
				newMatchKeyFn("foo"),
				newMatchKeyFn("bar"),
				newMatchKeyFn("baz"),
			***REMOVED******REMOVED***,
		))
***REMOVED***

func TestPathRecurse(t *testing.T) ***REMOVED***
	assertPath(t,
		"$..*",
		buildPath(
			newMatchRecursiveFn(),
		))
***REMOVED***

func TestPathFilterExpr(t *testing.T) ***REMOVED***
	assertPath(t,
		"$[?('foo'),?(bar)]",
		buildPath(
			&matchUnionFn***REMOVED***[]pathFn***REMOVED***
				newMatchFilterFn("foo", toml.Position***REMOVED******REMOVED***),
				newMatchFilterFn("bar", toml.Position***REMOVED******REMOVED***),
			***REMOVED******REMOVED***,
		))
***REMOVED***
