package jmespath

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type jpFunction func(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error)

type jpType string

const (
	jpUnknown     jpType = "unknown"
	jpNumber      jpType = "number"
	jpString      jpType = "string"
	jpArray       jpType = "array"
	jpObject      jpType = "object"
	jpArrayNumber jpType = "array[number]"
	jpArrayString jpType = "array[string]"
	jpExpref      jpType = "expref"
	jpAny         jpType = "any"
)

type functionEntry struct ***REMOVED***
	name      string
	arguments []argSpec
	handler   jpFunction
	hasExpRef bool
***REMOVED***

type argSpec struct ***REMOVED***
	types    []jpType
	variadic bool
***REMOVED***

type byExprString struct ***REMOVED***
	intr     *treeInterpreter
	node     ASTNode
	items    []interface***REMOVED******REMOVED***
	hasError bool
***REMOVED***

func (a *byExprString) Len() int ***REMOVED***
	return len(a.items)
***REMOVED***
func (a *byExprString) Swap(i, j int) ***REMOVED***
	a.items[i], a.items[j] = a.items[j], a.items[i]
***REMOVED***
func (a *byExprString) Less(i, j int) bool ***REMOVED***
	first, err := a.intr.Execute(a.node, a.items[i])
	if err != nil ***REMOVED***
		a.hasError = true
		// Return a dummy value.
		return true
	***REMOVED***
	ith, ok := first.(string)
	if !ok ***REMOVED***
		a.hasError = true
		return true
	***REMOVED***
	second, err := a.intr.Execute(a.node, a.items[j])
	if err != nil ***REMOVED***
		a.hasError = true
		// Return a dummy value.
		return true
	***REMOVED***
	jth, ok := second.(string)
	if !ok ***REMOVED***
		a.hasError = true
		return true
	***REMOVED***
	return ith < jth
***REMOVED***

type byExprFloat struct ***REMOVED***
	intr     *treeInterpreter
	node     ASTNode
	items    []interface***REMOVED******REMOVED***
	hasError bool
***REMOVED***

func (a *byExprFloat) Len() int ***REMOVED***
	return len(a.items)
***REMOVED***
func (a *byExprFloat) Swap(i, j int) ***REMOVED***
	a.items[i], a.items[j] = a.items[j], a.items[i]
***REMOVED***
func (a *byExprFloat) Less(i, j int) bool ***REMOVED***
	first, err := a.intr.Execute(a.node, a.items[i])
	if err != nil ***REMOVED***
		a.hasError = true
		// Return a dummy value.
		return true
	***REMOVED***
	ith, ok := first.(float64)
	if !ok ***REMOVED***
		a.hasError = true
		return true
	***REMOVED***
	second, err := a.intr.Execute(a.node, a.items[j])
	if err != nil ***REMOVED***
		a.hasError = true
		// Return a dummy value.
		return true
	***REMOVED***
	jth, ok := second.(float64)
	if !ok ***REMOVED***
		a.hasError = true
		return true
	***REMOVED***
	return ith < jth
***REMOVED***

type functionCaller struct ***REMOVED***
	functionTable map[string]functionEntry
***REMOVED***

func newFunctionCaller() *functionCaller ***REMOVED***
	caller := &functionCaller***REMOVED******REMOVED***
	caller.functionTable = map[string]functionEntry***REMOVED***
		"length": ***REMOVED***
			name: "length",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpString, jpArray, jpObject***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfLength,
		***REMOVED***,
		"starts_with": ***REMOVED***
			name: "starts_with",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpString***REMOVED******REMOVED***,
				***REMOVED***types: []jpType***REMOVED***jpString***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfStartsWith,
		***REMOVED***,
		"abs": ***REMOVED***
			name: "abs",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpNumber***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfAbs,
		***REMOVED***,
		"avg": ***REMOVED***
			name: "avg",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArrayNumber***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfAvg,
		***REMOVED***,
		"ceil": ***REMOVED***
			name: "ceil",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpNumber***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfCeil,
		***REMOVED***,
		"contains": ***REMOVED***
			name: "contains",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArray, jpString***REMOVED******REMOVED***,
				***REMOVED***types: []jpType***REMOVED***jpAny***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfContains,
		***REMOVED***,
		"ends_with": ***REMOVED***
			name: "ends_with",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpString***REMOVED******REMOVED***,
				***REMOVED***types: []jpType***REMOVED***jpString***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfEndsWith,
		***REMOVED***,
		"floor": ***REMOVED***
			name: "floor",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpNumber***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfFloor,
		***REMOVED***,
		"map": ***REMOVED***
			name: "amp",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpExpref***REMOVED******REMOVED***,
				***REMOVED***types: []jpType***REMOVED***jpArray***REMOVED******REMOVED***,
			***REMOVED***,
			handler:   jpfMap,
			hasExpRef: true,
		***REMOVED***,
		"max": ***REMOVED***
			name: "max",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArrayNumber, jpArrayString***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfMax,
		***REMOVED***,
		"merge": ***REMOVED***
			name: "merge",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpObject***REMOVED***, variadic: true***REMOVED***,
			***REMOVED***,
			handler: jpfMerge,
		***REMOVED***,
		"max_by": ***REMOVED***
			name: "max_by",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArray***REMOVED******REMOVED***,
				***REMOVED***types: []jpType***REMOVED***jpExpref***REMOVED******REMOVED***,
			***REMOVED***,
			handler:   jpfMaxBy,
			hasExpRef: true,
		***REMOVED***,
		"sum": ***REMOVED***
			name: "sum",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArrayNumber***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfSum,
		***REMOVED***,
		"min": ***REMOVED***
			name: "min",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArrayNumber, jpArrayString***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfMin,
		***REMOVED***,
		"min_by": ***REMOVED***
			name: "min_by",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArray***REMOVED******REMOVED***,
				***REMOVED***types: []jpType***REMOVED***jpExpref***REMOVED******REMOVED***,
			***REMOVED***,
			handler:   jpfMinBy,
			hasExpRef: true,
		***REMOVED***,
		"type": ***REMOVED***
			name: "type",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpAny***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfType,
		***REMOVED***,
		"keys": ***REMOVED***
			name: "keys",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpObject***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfKeys,
		***REMOVED***,
		"values": ***REMOVED***
			name: "values",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpObject***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfValues,
		***REMOVED***,
		"sort": ***REMOVED***
			name: "sort",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArrayString, jpArrayNumber***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfSort,
		***REMOVED***,
		"sort_by": ***REMOVED***
			name: "sort_by",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArray***REMOVED******REMOVED***,
				***REMOVED***types: []jpType***REMOVED***jpExpref***REMOVED******REMOVED***,
			***REMOVED***,
			handler:   jpfSortBy,
			hasExpRef: true,
		***REMOVED***,
		"join": ***REMOVED***
			name: "join",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpString***REMOVED******REMOVED***,
				***REMOVED***types: []jpType***REMOVED***jpArrayString***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfJoin,
		***REMOVED***,
		"reverse": ***REMOVED***
			name: "reverse",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpArray, jpString***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfReverse,
		***REMOVED***,
		"to_array": ***REMOVED***
			name: "to_array",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpAny***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfToArray,
		***REMOVED***,
		"to_string": ***REMOVED***
			name: "to_string",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpAny***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfToString,
		***REMOVED***,
		"to_number": ***REMOVED***
			name: "to_number",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpAny***REMOVED******REMOVED***,
			***REMOVED***,
			handler: jpfToNumber,
		***REMOVED***,
		"not_null": ***REMOVED***
			name: "not_null",
			arguments: []argSpec***REMOVED***
				***REMOVED***types: []jpType***REMOVED***jpAny***REMOVED***, variadic: true***REMOVED***,
			***REMOVED***,
			handler: jpfNotNull,
		***REMOVED***,
	***REMOVED***
	return caller
***REMOVED***

func (e *functionEntry) resolveArgs(arguments []interface***REMOVED******REMOVED***) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	if len(e.arguments) == 0 ***REMOVED***
		return arguments, nil
	***REMOVED***
	if !e.arguments[len(e.arguments)-1].variadic ***REMOVED***
		if len(e.arguments) != len(arguments) ***REMOVED***
			return nil, errors.New("incorrect number of args")
		***REMOVED***
		for i, spec := range e.arguments ***REMOVED***
			userArg := arguments[i]
			err := spec.typeCheck(userArg)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		return arguments, nil
	***REMOVED***
	if len(arguments) < len(e.arguments) ***REMOVED***
		return nil, errors.New("Invalid arity.")
	***REMOVED***
	return arguments, nil
***REMOVED***

func (a *argSpec) typeCheck(arg interface***REMOVED******REMOVED***) error ***REMOVED***
	for _, t := range a.types ***REMOVED***
		switch t ***REMOVED***
		case jpNumber:
			if _, ok := arg.(float64); ok ***REMOVED***
				return nil
			***REMOVED***
		case jpString:
			if _, ok := arg.(string); ok ***REMOVED***
				return nil
			***REMOVED***
		case jpArray:
			if isSliceType(arg) ***REMOVED***
				return nil
			***REMOVED***
		case jpObject:
			if _, ok := arg.(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
				return nil
			***REMOVED***
		case jpArrayNumber:
			if _, ok := toArrayNum(arg); ok ***REMOVED***
				return nil
			***REMOVED***
		case jpArrayString:
			if _, ok := toArrayStr(arg); ok ***REMOVED***
				return nil
			***REMOVED***
		case jpAny:
			return nil
		case jpExpref:
			if _, ok := arg.(expRef); ok ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return fmt.Errorf("Invalid type for: %v, expected: %#v", arg, a.types)
***REMOVED***

func (f *functionCaller) CallFunction(name string, arguments []interface***REMOVED******REMOVED***, intr *treeInterpreter) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	entry, ok := f.functionTable[name]
	if !ok ***REMOVED***
		return nil, errors.New("unknown function: " + name)
	***REMOVED***
	resolvedArgs, err := entry.resolveArgs(arguments)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if entry.hasExpRef ***REMOVED***
		var extra []interface***REMOVED******REMOVED***
		extra = append(extra, intr)
		resolvedArgs = append(extra, resolvedArgs...)
	***REMOVED***
	return entry.handler(resolvedArgs)
***REMOVED***

func jpfAbs(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	num := arguments[0].(float64)
	return math.Abs(num), nil
***REMOVED***

func jpfLength(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	arg := arguments[0]
	if c, ok := arg.(string); ok ***REMOVED***
		return float64(utf8.RuneCountInString(c)), nil
	***REMOVED*** else if isSliceType(arg) ***REMOVED***
		v := reflect.ValueOf(arg)
		return float64(v.Len()), nil
	***REMOVED*** else if c, ok := arg.(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return float64(len(c)), nil
	***REMOVED***
	return nil, errors.New("could not compute length()")
***REMOVED***

func jpfStartsWith(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	search := arguments[0].(string)
	prefix := arguments[1].(string)
	return strings.HasPrefix(search, prefix), nil
***REMOVED***

func jpfAvg(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// We've already type checked the value so we can safely use
	// type assertions.
	args := arguments[0].([]interface***REMOVED******REMOVED***)
	length := float64(len(args))
	numerator := 0.0
	for _, n := range args ***REMOVED***
		numerator += n.(float64)
	***REMOVED***
	return numerator / length, nil
***REMOVED***
func jpfCeil(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val := arguments[0].(float64)
	return math.Ceil(val), nil
***REMOVED***
func jpfContains(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	search := arguments[0]
	el := arguments[1]
	if searchStr, ok := search.(string); ok ***REMOVED***
		if elStr, ok := el.(string); ok ***REMOVED***
			return strings.Index(searchStr, elStr) != -1, nil
		***REMOVED***
		return false, nil
	***REMOVED***
	// Otherwise this is a generic contains for []interface***REMOVED******REMOVED***
	general := search.([]interface***REMOVED******REMOVED***)
	for _, item := range general ***REMOVED***
		if item == el ***REMOVED***
			return true, nil
		***REMOVED***
	***REMOVED***
	return false, nil
***REMOVED***
func jpfEndsWith(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	search := arguments[0].(string)
	suffix := arguments[1].(string)
	return strings.HasSuffix(search, suffix), nil
***REMOVED***
func jpfFloor(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val := arguments[0].(float64)
	return math.Floor(val), nil
***REMOVED***
func jpfMap(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	intr := arguments[0].(*treeInterpreter)
	exp := arguments[1].(expRef)
	node := exp.ref
	arr := arguments[2].([]interface***REMOVED******REMOVED***)
	mapped := make([]interface***REMOVED******REMOVED***, 0, len(arr))
	for _, value := range arr ***REMOVED***
		current, err := intr.Execute(node, value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		mapped = append(mapped, current)
	***REMOVED***
	return mapped, nil
***REMOVED***
func jpfMax(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if items, ok := toArrayNum(arguments[0]); ok ***REMOVED***
		if len(items) == 0 ***REMOVED***
			return nil, nil
		***REMOVED***
		if len(items) == 1 ***REMOVED***
			return items[0], nil
		***REMOVED***
		best := items[0]
		for _, item := range items[1:] ***REMOVED***
			if item > best ***REMOVED***
				best = item
			***REMOVED***
		***REMOVED***
		return best, nil
	***REMOVED***
	// Otherwise we're dealing with a max() of strings.
	items, _ := toArrayStr(arguments[0])
	if len(items) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	if len(items) == 1 ***REMOVED***
		return items[0], nil
	***REMOVED***
	best := items[0]
	for _, item := range items[1:] ***REMOVED***
		if item > best ***REMOVED***
			best = item
		***REMOVED***
	***REMOVED***
	return best, nil
***REMOVED***
func jpfMerge(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	final := make(map[string]interface***REMOVED******REMOVED***)
	for _, m := range arguments ***REMOVED***
		mapped := m.(map[string]interface***REMOVED******REMOVED***)
		for key, value := range mapped ***REMOVED***
			final[key] = value
		***REMOVED***
	***REMOVED***
	return final, nil
***REMOVED***
func jpfMaxBy(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	intr := arguments[0].(*treeInterpreter)
	arr := arguments[1].([]interface***REMOVED******REMOVED***)
	exp := arguments[2].(expRef)
	node := exp.ref
	if len(arr) == 0 ***REMOVED***
		return nil, nil
	***REMOVED*** else if len(arr) == 1 ***REMOVED***
		return arr[0], nil
	***REMOVED***
	start, err := intr.Execute(node, arr[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch t := start.(type) ***REMOVED***
	case float64:
		bestVal := t
		bestItem := arr[0]
		for _, item := range arr[1:] ***REMOVED***
			result, err := intr.Execute(node, item)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			current, ok := result.(float64)
			if !ok ***REMOVED***
				return nil, errors.New("invalid type, must be number")
			***REMOVED***
			if current > bestVal ***REMOVED***
				bestVal = current
				bestItem = item
			***REMOVED***
		***REMOVED***
		return bestItem, nil
	case string:
		bestVal := t
		bestItem := arr[0]
		for _, item := range arr[1:] ***REMOVED***
			result, err := intr.Execute(node, item)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			current, ok := result.(string)
			if !ok ***REMOVED***
				return nil, errors.New("invalid type, must be string")
			***REMOVED***
			if current > bestVal ***REMOVED***
				bestVal = current
				bestItem = item
			***REMOVED***
		***REMOVED***
		return bestItem, nil
	default:
		return nil, errors.New("invalid type, must be number of string")
	***REMOVED***
***REMOVED***
func jpfSum(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	items, _ := toArrayNum(arguments[0])
	sum := 0.0
	for _, item := range items ***REMOVED***
		sum += item
	***REMOVED***
	return sum, nil
***REMOVED***

func jpfMin(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if items, ok := toArrayNum(arguments[0]); ok ***REMOVED***
		if len(items) == 0 ***REMOVED***
			return nil, nil
		***REMOVED***
		if len(items) == 1 ***REMOVED***
			return items[0], nil
		***REMOVED***
		best := items[0]
		for _, item := range items[1:] ***REMOVED***
			if item < best ***REMOVED***
				best = item
			***REMOVED***
		***REMOVED***
		return best, nil
	***REMOVED***
	items, _ := toArrayStr(arguments[0])
	if len(items) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	if len(items) == 1 ***REMOVED***
		return items[0], nil
	***REMOVED***
	best := items[0]
	for _, item := range items[1:] ***REMOVED***
		if item < best ***REMOVED***
			best = item
		***REMOVED***
	***REMOVED***
	return best, nil
***REMOVED***

func jpfMinBy(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	intr := arguments[0].(*treeInterpreter)
	arr := arguments[1].([]interface***REMOVED******REMOVED***)
	exp := arguments[2].(expRef)
	node := exp.ref
	if len(arr) == 0 ***REMOVED***
		return nil, nil
	***REMOVED*** else if len(arr) == 1 ***REMOVED***
		return arr[0], nil
	***REMOVED***
	start, err := intr.Execute(node, arr[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if t, ok := start.(float64); ok ***REMOVED***
		bestVal := t
		bestItem := arr[0]
		for _, item := range arr[1:] ***REMOVED***
			result, err := intr.Execute(node, item)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			current, ok := result.(float64)
			if !ok ***REMOVED***
				return nil, errors.New("invalid type, must be number")
			***REMOVED***
			if current < bestVal ***REMOVED***
				bestVal = current
				bestItem = item
			***REMOVED***
		***REMOVED***
		return bestItem, nil
	***REMOVED*** else if t, ok := start.(string); ok ***REMOVED***
		bestVal := t
		bestItem := arr[0]
		for _, item := range arr[1:] ***REMOVED***
			result, err := intr.Execute(node, item)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			current, ok := result.(string)
			if !ok ***REMOVED***
				return nil, errors.New("invalid type, must be string")
			***REMOVED***
			if current < bestVal ***REMOVED***
				bestVal = current
				bestItem = item
			***REMOVED***
		***REMOVED***
		return bestItem, nil
	***REMOVED*** else ***REMOVED***
		return nil, errors.New("invalid type, must be number of string")
	***REMOVED***
***REMOVED***
func jpfType(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	arg := arguments[0]
	if _, ok := arg.(float64); ok ***REMOVED***
		return "number", nil
	***REMOVED***
	if _, ok := arg.(string); ok ***REMOVED***
		return "string", nil
	***REMOVED***
	if _, ok := arg.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return "array", nil
	***REMOVED***
	if _, ok := arg.(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return "object", nil
	***REMOVED***
	if arg == nil ***REMOVED***
		return "null", nil
	***REMOVED***
	if arg == true || arg == false ***REMOVED***
		return "boolean", nil
	***REMOVED***
	return nil, errors.New("unknown type")
***REMOVED***
func jpfKeys(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	arg := arguments[0].(map[string]interface***REMOVED******REMOVED***)
	collected := make([]interface***REMOVED******REMOVED***, 0, len(arg))
	for key := range arg ***REMOVED***
		collected = append(collected, key)
	***REMOVED***
	return collected, nil
***REMOVED***
func jpfValues(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	arg := arguments[0].(map[string]interface***REMOVED******REMOVED***)
	collected := make([]interface***REMOVED******REMOVED***, 0, len(arg))
	for _, value := range arg ***REMOVED***
		collected = append(collected, value)
	***REMOVED***
	return collected, nil
***REMOVED***
func jpfSort(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if items, ok := toArrayNum(arguments[0]); ok ***REMOVED***
		d := sort.Float64Slice(items)
		sort.Stable(d)
		final := make([]interface***REMOVED******REMOVED***, len(d))
		for i, val := range d ***REMOVED***
			final[i] = val
		***REMOVED***
		return final, nil
	***REMOVED***
	// Otherwise we're dealing with sort()'ing strings.
	items, _ := toArrayStr(arguments[0])
	d := sort.StringSlice(items)
	sort.Stable(d)
	final := make([]interface***REMOVED******REMOVED***, len(d))
	for i, val := range d ***REMOVED***
		final[i] = val
	***REMOVED***
	return final, nil
***REMOVED***
func jpfSortBy(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	intr := arguments[0].(*treeInterpreter)
	arr := arguments[1].([]interface***REMOVED******REMOVED***)
	exp := arguments[2].(expRef)
	node := exp.ref
	if len(arr) == 0 ***REMOVED***
		return arr, nil
	***REMOVED*** else if len(arr) == 1 ***REMOVED***
		return arr, nil
	***REMOVED***
	start, err := intr.Execute(node, arr[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, ok := start.(float64); ok ***REMOVED***
		sortable := &byExprFloat***REMOVED***intr, node, arr, false***REMOVED***
		sort.Stable(sortable)
		if sortable.hasError ***REMOVED***
			return nil, errors.New("error in sort_by comparison")
		***REMOVED***
		return arr, nil
	***REMOVED*** else if _, ok := start.(string); ok ***REMOVED***
		sortable := &byExprString***REMOVED***intr, node, arr, false***REMOVED***
		sort.Stable(sortable)
		if sortable.hasError ***REMOVED***
			return nil, errors.New("error in sort_by comparison")
		***REMOVED***
		return arr, nil
	***REMOVED*** else ***REMOVED***
		return nil, errors.New("invalid type, must be number of string")
	***REMOVED***
***REMOVED***
func jpfJoin(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	sep := arguments[0].(string)
	// We can't just do arguments[1].([]string), we have to
	// manually convert each item to a string.
	arrayStr := []string***REMOVED******REMOVED***
	for _, item := range arguments[1].([]interface***REMOVED******REMOVED***) ***REMOVED***
		arrayStr = append(arrayStr, item.(string))
	***REMOVED***
	return strings.Join(arrayStr, sep), nil
***REMOVED***
func jpfReverse(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if s, ok := arguments[0].(string); ok ***REMOVED***
		r := []rune(s)
		for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 ***REMOVED***
			r[i], r[j] = r[j], r[i]
		***REMOVED***
		return string(r), nil
	***REMOVED***
	items := arguments[0].([]interface***REMOVED******REMOVED***)
	length := len(items)
	reversed := make([]interface***REMOVED******REMOVED***, length)
	for i, item := range items ***REMOVED***
		reversed[length-(i+1)] = item
	***REMOVED***
	return reversed, nil
***REMOVED***
func jpfToArray(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if _, ok := arguments[0].([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return arguments[0], nil
	***REMOVED***
	return arguments[:1:1], nil
***REMOVED***
func jpfToString(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if v, ok := arguments[0].(string); ok ***REMOVED***
		return v, nil
	***REMOVED***
	result, err := json.Marshal(arguments[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return string(result), nil
***REMOVED***
func jpfToNumber(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	arg := arguments[0]
	if v, ok := arg.(float64); ok ***REMOVED***
		return v, nil
	***REMOVED***
	if v, ok := arg.(string); ok ***REMOVED***
		conv, err := strconv.ParseFloat(v, 64)
		if err != nil ***REMOVED***
			return nil, nil
		***REMOVED***
		return conv, nil
	***REMOVED***
	if _, ok := arg.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return nil, nil
	***REMOVED***
	if _, ok := arg.(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return nil, nil
	***REMOVED***
	if arg == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	if arg == true || arg == false ***REMOVED***
		return nil, nil
	***REMOVED***
	return nil, errors.New("unknown type")
***REMOVED***
func jpfNotNull(arguments []interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	for _, arg := range arguments ***REMOVED***
		if arg != nil ***REMOVED***
			return arg, nil
		***REMOVED***
	***REMOVED***
	return nil, nil
***REMOVED***
