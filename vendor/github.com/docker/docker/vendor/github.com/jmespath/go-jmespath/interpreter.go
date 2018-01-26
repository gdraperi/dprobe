package jmespath

import (
	"errors"
	"reflect"
	"unicode"
	"unicode/utf8"
)

/* This is a tree based interpreter.  It walks the AST and directly
   interprets the AST to search through a JSON document.
*/

type treeInterpreter struct ***REMOVED***
	fCall *functionCaller
***REMOVED***

func newInterpreter() *treeInterpreter ***REMOVED***
	interpreter := treeInterpreter***REMOVED******REMOVED***
	interpreter.fCall = newFunctionCaller()
	return &interpreter
***REMOVED***

type expRef struct ***REMOVED***
	ref ASTNode
***REMOVED***

// Execute takes an ASTNode and input data and interprets the AST directly.
// It will produce the result of applying the JMESPath expression associated
// with the ASTNode to the input data "value".
func (intr *treeInterpreter) Execute(node ASTNode, value interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch node.nodeType ***REMOVED***
	case ASTComparator:
		left, err := intr.Execute(node.children[0], value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		right, err := intr.Execute(node.children[1], value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		switch node.value ***REMOVED***
		case tEQ:
			return objsEqual(left, right), nil
		case tNE:
			return !objsEqual(left, right), nil
		***REMOVED***
		leftNum, ok := left.(float64)
		if !ok ***REMOVED***
			return nil, nil
		***REMOVED***
		rightNum, ok := right.(float64)
		if !ok ***REMOVED***
			return nil, nil
		***REMOVED***
		switch node.value ***REMOVED***
		case tGT:
			return leftNum > rightNum, nil
		case tGTE:
			return leftNum >= rightNum, nil
		case tLT:
			return leftNum < rightNum, nil
		case tLTE:
			return leftNum <= rightNum, nil
		***REMOVED***
	case ASTExpRef:
		return expRef***REMOVED***ref: node.children[0]***REMOVED***, nil
	case ASTFunctionExpression:
		resolvedArgs := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
		for _, arg := range node.children ***REMOVED***
			current, err := intr.Execute(arg, value)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			resolvedArgs = append(resolvedArgs, current)
		***REMOVED***
		return intr.fCall.CallFunction(node.value.(string), resolvedArgs, intr)
	case ASTField:
		if m, ok := value.(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
			key := node.value.(string)
			return m[key], nil
		***REMOVED***
		return intr.fieldFromStruct(node.value.(string), value)
	case ASTFilterProjection:
		left, err := intr.Execute(node.children[0], value)
		if err != nil ***REMOVED***
			return nil, nil
		***REMOVED***
		sliceType, ok := left.([]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			if isSliceType(left) ***REMOVED***
				return intr.filterProjectionWithReflection(node, left)
			***REMOVED***
			return nil, nil
		***REMOVED***
		compareNode := node.children[2]
		collected := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
		for _, element := range sliceType ***REMOVED***
			result, err := intr.Execute(compareNode, element)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if !isFalse(result) ***REMOVED***
				current, err := intr.Execute(node.children[1], element)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				if current != nil ***REMOVED***
					collected = append(collected, current)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return collected, nil
	case ASTFlatten:
		left, err := intr.Execute(node.children[0], value)
		if err != nil ***REMOVED***
			return nil, nil
		***REMOVED***
		sliceType, ok := left.([]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			// If we can't type convert to []interface***REMOVED******REMOVED***, there's
			// a chance this could still work via reflection if we're
			// dealing with user provided types.
			if isSliceType(left) ***REMOVED***
				return intr.flattenWithReflection(left)
			***REMOVED***
			return nil, nil
		***REMOVED***
		flattened := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
		for _, element := range sliceType ***REMOVED***
			if elementSlice, ok := element.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
				flattened = append(flattened, elementSlice...)
			***REMOVED*** else if isSliceType(element) ***REMOVED***
				reflectFlat := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
				v := reflect.ValueOf(element)
				for i := 0; i < v.Len(); i++ ***REMOVED***
					reflectFlat = append(reflectFlat, v.Index(i).Interface())
				***REMOVED***
				flattened = append(flattened, reflectFlat...)
			***REMOVED*** else ***REMOVED***
				flattened = append(flattened, element)
			***REMOVED***
		***REMOVED***
		return flattened, nil
	case ASTIdentity, ASTCurrentNode:
		return value, nil
	case ASTIndex:
		if sliceType, ok := value.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
			index := node.value.(int)
			if index < 0 ***REMOVED***
				index += len(sliceType)
			***REMOVED***
			if index < len(sliceType) && index >= 0 ***REMOVED***
				return sliceType[index], nil
			***REMOVED***
			return nil, nil
		***REMOVED***
		// Otherwise try via reflection.
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Slice ***REMOVED***
			index := node.value.(int)
			if index < 0 ***REMOVED***
				index += rv.Len()
			***REMOVED***
			if index < rv.Len() && index >= 0 ***REMOVED***
				v := rv.Index(index)
				return v.Interface(), nil
			***REMOVED***
		***REMOVED***
		return nil, nil
	case ASTKeyValPair:
		return intr.Execute(node.children[0], value)
	case ASTLiteral:
		return node.value, nil
	case ASTMultiSelectHash:
		if value == nil ***REMOVED***
			return nil, nil
		***REMOVED***
		collected := make(map[string]interface***REMOVED******REMOVED***)
		for _, child := range node.children ***REMOVED***
			current, err := intr.Execute(child, value)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			key := child.value.(string)
			collected[key] = current
		***REMOVED***
		return collected, nil
	case ASTMultiSelectList:
		if value == nil ***REMOVED***
			return nil, nil
		***REMOVED***
		collected := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
		for _, child := range node.children ***REMOVED***
			current, err := intr.Execute(child, value)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			collected = append(collected, current)
		***REMOVED***
		return collected, nil
	case ASTOrExpression:
		matched, err := intr.Execute(node.children[0], value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if isFalse(matched) ***REMOVED***
			matched, err = intr.Execute(node.children[1], value)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		return matched, nil
	case ASTAndExpression:
		matched, err := intr.Execute(node.children[0], value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if isFalse(matched) ***REMOVED***
			return matched, nil
		***REMOVED***
		return intr.Execute(node.children[1], value)
	case ASTNotExpression:
		matched, err := intr.Execute(node.children[0], value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if isFalse(matched) ***REMOVED***
			return true, nil
		***REMOVED***
		return false, nil
	case ASTPipe:
		result := value
		var err error
		for _, child := range node.children ***REMOVED***
			result, err = intr.Execute(child, result)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		return result, nil
	case ASTProjection:
		left, err := intr.Execute(node.children[0], value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		sliceType, ok := left.([]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			if isSliceType(left) ***REMOVED***
				return intr.projectWithReflection(node, left)
			***REMOVED***
			return nil, nil
		***REMOVED***
		collected := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
		var current interface***REMOVED******REMOVED***
		for _, element := range sliceType ***REMOVED***
			current, err = intr.Execute(node.children[1], element)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if current != nil ***REMOVED***
				collected = append(collected, current)
			***REMOVED***
		***REMOVED***
		return collected, nil
	case ASTSubexpression, ASTIndexExpression:
		left, err := intr.Execute(node.children[0], value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return intr.Execute(node.children[1], left)
	case ASTSlice:
		sliceType, ok := value.([]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			if isSliceType(value) ***REMOVED***
				return intr.sliceWithReflection(node, value)
			***REMOVED***
			return nil, nil
		***REMOVED***
		parts := node.value.([]*int)
		sliceParams := make([]sliceParam, 3)
		for i, part := range parts ***REMOVED***
			if part != nil ***REMOVED***
				sliceParams[i].Specified = true
				sliceParams[i].N = *part
			***REMOVED***
		***REMOVED***
		return slice(sliceType, sliceParams)
	case ASTValueProjection:
		left, err := intr.Execute(node.children[0], value)
		if err != nil ***REMOVED***
			return nil, nil
		***REMOVED***
		mapType, ok := left.(map[string]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			return nil, nil
		***REMOVED***
		values := make([]interface***REMOVED******REMOVED***, len(mapType))
		for _, value := range mapType ***REMOVED***
			values = append(values, value)
		***REMOVED***
		collected := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
		for _, element := range values ***REMOVED***
			current, err := intr.Execute(node.children[1], element)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if current != nil ***REMOVED***
				collected = append(collected, current)
			***REMOVED***
		***REMOVED***
		return collected, nil
	***REMOVED***
	return nil, errors.New("Unknown AST node: " + node.nodeType.String())
***REMOVED***

func (intr *treeInterpreter) fieldFromStruct(key string, value interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	rv := reflect.ValueOf(value)
	first, n := utf8.DecodeRuneInString(key)
	fieldName := string(unicode.ToUpper(first)) + key[n:]
	if rv.Kind() == reflect.Struct ***REMOVED***
		v := rv.FieldByName(fieldName)
		if !v.IsValid() ***REMOVED***
			return nil, nil
		***REMOVED***
		return v.Interface(), nil
	***REMOVED*** else if rv.Kind() == reflect.Ptr ***REMOVED***
		// Handle multiple levels of indirection?
		if rv.IsNil() ***REMOVED***
			return nil, nil
		***REMOVED***
		rv = rv.Elem()
		v := rv.FieldByName(fieldName)
		if !v.IsValid() ***REMOVED***
			return nil, nil
		***REMOVED***
		return v.Interface(), nil
	***REMOVED***
	return nil, nil
***REMOVED***

func (intr *treeInterpreter) flattenWithReflection(value interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	v := reflect.ValueOf(value)
	flattened := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
	for i := 0; i < v.Len(); i++ ***REMOVED***
		element := v.Index(i).Interface()
		if reflect.TypeOf(element).Kind() == reflect.Slice ***REMOVED***
			// Then insert the contents of the element
			// slice into the flattened slice,
			// i.e flattened = append(flattened, mySlice...)
			elementV := reflect.ValueOf(element)
			for j := 0; j < elementV.Len(); j++ ***REMOVED***
				flattened = append(
					flattened, elementV.Index(j).Interface())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			flattened = append(flattened, element)
		***REMOVED***
	***REMOVED***
	return flattened, nil
***REMOVED***

func (intr *treeInterpreter) sliceWithReflection(node ASTNode, value interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	v := reflect.ValueOf(value)
	parts := node.value.([]*int)
	sliceParams := make([]sliceParam, 3)
	for i, part := range parts ***REMOVED***
		if part != nil ***REMOVED***
			sliceParams[i].Specified = true
			sliceParams[i].N = *part
		***REMOVED***
	***REMOVED***
	final := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
	for i := 0; i < v.Len(); i++ ***REMOVED***
		element := v.Index(i).Interface()
		final = append(final, element)
	***REMOVED***
	return slice(final, sliceParams)
***REMOVED***

func (intr *treeInterpreter) filterProjectionWithReflection(node ASTNode, value interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	compareNode := node.children[2]
	collected := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
	v := reflect.ValueOf(value)
	for i := 0; i < v.Len(); i++ ***REMOVED***
		element := v.Index(i).Interface()
		result, err := intr.Execute(compareNode, element)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if !isFalse(result) ***REMOVED***
			current, err := intr.Execute(node.children[1], element)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if current != nil ***REMOVED***
				collected = append(collected, current)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return collected, nil
***REMOVED***

func (intr *treeInterpreter) projectWithReflection(node ASTNode, value interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	collected := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
	v := reflect.ValueOf(value)
	for i := 0; i < v.Len(); i++ ***REMOVED***
		element := v.Index(i).Interface()
		result, err := intr.Execute(node.children[1], element)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if result != nil ***REMOVED***
			collected = append(collected, result)
		***REMOVED***
	***REMOVED***
	return collected, nil
***REMOVED***
