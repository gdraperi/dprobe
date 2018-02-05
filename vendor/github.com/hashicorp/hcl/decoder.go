package hcl

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/token"
)

// This is the tag to use with structures to have settings for HCL
const tagName = "hcl"

var (
	// nodeType holds a reference to the type of ast.Node
	nodeType reflect.Type = findNodeType()
)

// Unmarshal accepts a byte slice as input and writes the
// data to the value pointed to by v.
func Unmarshal(bs []byte, v interface***REMOVED******REMOVED***) error ***REMOVED***
	root, err := parse(bs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return DecodeObject(v, root)
***REMOVED***

// Decode reads the given input and decodes it into the structure
// given by `out`.
func Decode(out interface***REMOVED******REMOVED***, in string) error ***REMOVED***
	obj, err := Parse(in)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return DecodeObject(out, obj)
***REMOVED***

// DecodeObject is a lower-level version of Decode. It decodes a
// raw Object into the given output.
func DecodeObject(out interface***REMOVED******REMOVED***, n ast.Node) error ***REMOVED***
	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr ***REMOVED***
		return errors.New("result must be a pointer")
	***REMOVED***

	// If we have the file, we really decode the root node
	if f, ok := n.(*ast.File); ok ***REMOVED***
		n = f.Node
	***REMOVED***

	var d decoder
	return d.decode("root", n, val.Elem())
***REMOVED***

type decoder struct ***REMOVED***
	stack []reflect.Kind
***REMOVED***

func (d *decoder) decode(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	k := result

	// If we have an interface with a valid value, we use that
	// for the check.
	if result.Kind() == reflect.Interface ***REMOVED***
		elem := result.Elem()
		if elem.IsValid() ***REMOVED***
			k = elem
		***REMOVED***
	***REMOVED***

	// Push current onto stack unless it is an interface.
	if k.Kind() != reflect.Interface ***REMOVED***
		d.stack = append(d.stack, k.Kind())

		// Schedule a pop
		defer func() ***REMOVED***
			d.stack = d.stack[:len(d.stack)-1]
		***REMOVED***()
	***REMOVED***

	switch k.Kind() ***REMOVED***
	case reflect.Bool:
		return d.decodeBool(name, node, result)
	case reflect.Float32, reflect.Float64:
		return d.decodeFloat(name, node, result)
	case reflect.Int, reflect.Int32, reflect.Int64:
		return d.decodeInt(name, node, result)
	case reflect.Interface:
		// When we see an interface, we make our own thing
		return d.decodeInterface(name, node, result)
	case reflect.Map:
		return d.decodeMap(name, node, result)
	case reflect.Ptr:
		return d.decodePtr(name, node, result)
	case reflect.Slice:
		return d.decodeSlice(name, node, result)
	case reflect.String:
		return d.decodeString(name, node, result)
	case reflect.Struct:
		return d.decodeStruct(name, node, result)
	default:
		return &parser.PosError***REMOVED***
			Pos: node.Pos(),
			Err: fmt.Errorf("%s: unknown kind to decode into: %s", name, k.Kind()),
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *decoder) decodeBool(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	switch n := node.(type) ***REMOVED***
	case *ast.LiteralType:
		if n.Token.Type == token.BOOL ***REMOVED***
			v, err := strconv.ParseBool(n.Token.Text)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			result.Set(reflect.ValueOf(v))
			return nil
		***REMOVED***
	***REMOVED***

	return &parser.PosError***REMOVED***
		Pos: node.Pos(),
		Err: fmt.Errorf("%s: unknown type %T", name, node),
	***REMOVED***
***REMOVED***

func (d *decoder) decodeFloat(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	switch n := node.(type) ***REMOVED***
	case *ast.LiteralType:
		if n.Token.Type == token.FLOAT || n.Token.Type == token.NUMBER ***REMOVED***
			v, err := strconv.ParseFloat(n.Token.Text, 64)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			result.Set(reflect.ValueOf(v).Convert(result.Type()))
			return nil
		***REMOVED***
	***REMOVED***

	return &parser.PosError***REMOVED***
		Pos: node.Pos(),
		Err: fmt.Errorf("%s: unknown type %T", name, node),
	***REMOVED***
***REMOVED***

func (d *decoder) decodeInt(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	switch n := node.(type) ***REMOVED***
	case *ast.LiteralType:
		switch n.Token.Type ***REMOVED***
		case token.NUMBER:
			v, err := strconv.ParseInt(n.Token.Text, 0, 0)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if result.Kind() == reflect.Interface ***REMOVED***
				result.Set(reflect.ValueOf(int(v)))
			***REMOVED*** else ***REMOVED***
				result.SetInt(v)
			***REMOVED***
			return nil
		case token.STRING:
			v, err := strconv.ParseInt(n.Token.Value().(string), 0, 0)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if result.Kind() == reflect.Interface ***REMOVED***
				result.Set(reflect.ValueOf(int(v)))
			***REMOVED*** else ***REMOVED***
				result.SetInt(v)
			***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	return &parser.PosError***REMOVED***
		Pos: node.Pos(),
		Err: fmt.Errorf("%s: unknown type %T", name, node),
	***REMOVED***
***REMOVED***

func (d *decoder) decodeInterface(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	// When we see an ast.Node, we retain the value to enable deferred decoding.
	// Very useful in situations where we want to preserve ast.Node information
	// like Pos
	if result.Type() == nodeType && result.CanSet() ***REMOVED***
		result.Set(reflect.ValueOf(node))
		return nil
	***REMOVED***

	var set reflect.Value
	redecode := true

	// For testing types, ObjectType should just be treated as a list. We
	// set this to a temporary var because we want to pass in the real node.
	testNode := node
	if ot, ok := node.(*ast.ObjectType); ok ***REMOVED***
		testNode = ot.List
	***REMOVED***

	switch n := testNode.(type) ***REMOVED***
	case *ast.ObjectList:
		// If we're at the root or we're directly within a slice, then we
		// decode objects into map[string]interface***REMOVED******REMOVED***, otherwise we decode
		// them into lists.
		if len(d.stack) == 0 || d.stack[len(d.stack)-1] == reflect.Slice ***REMOVED***
			var temp map[string]interface***REMOVED******REMOVED***
			tempVal := reflect.ValueOf(temp)
			result := reflect.MakeMap(
				reflect.MapOf(
					reflect.TypeOf(""),
					tempVal.Type().Elem()))

			set = result
		***REMOVED*** else ***REMOVED***
			var temp []map[string]interface***REMOVED******REMOVED***
			tempVal := reflect.ValueOf(temp)
			result := reflect.MakeSlice(
				reflect.SliceOf(tempVal.Type().Elem()), 0, len(n.Items))
			set = result
		***REMOVED***
	case *ast.ObjectType:
		// If we're at the root or we're directly within a slice, then we
		// decode objects into map[string]interface***REMOVED******REMOVED***, otherwise we decode
		// them into lists.
		if len(d.stack) == 0 || d.stack[len(d.stack)-1] == reflect.Slice ***REMOVED***
			var temp map[string]interface***REMOVED******REMOVED***
			tempVal := reflect.ValueOf(temp)
			result := reflect.MakeMap(
				reflect.MapOf(
					reflect.TypeOf(""),
					tempVal.Type().Elem()))

			set = result
		***REMOVED*** else ***REMOVED***
			var temp []map[string]interface***REMOVED******REMOVED***
			tempVal := reflect.ValueOf(temp)
			result := reflect.MakeSlice(
				reflect.SliceOf(tempVal.Type().Elem()), 0, 1)
			set = result
		***REMOVED***
	case *ast.ListType:
		var temp []interface***REMOVED******REMOVED***
		tempVal := reflect.ValueOf(temp)
		result := reflect.MakeSlice(
			reflect.SliceOf(tempVal.Type().Elem()), 0, 0)
		set = result
	case *ast.LiteralType:
		switch n.Token.Type ***REMOVED***
		case token.BOOL:
			var result bool
			set = reflect.Indirect(reflect.New(reflect.TypeOf(result)))
		case token.FLOAT:
			var result float64
			set = reflect.Indirect(reflect.New(reflect.TypeOf(result)))
		case token.NUMBER:
			var result int
			set = reflect.Indirect(reflect.New(reflect.TypeOf(result)))
		case token.STRING, token.HEREDOC:
			set = reflect.Indirect(reflect.New(reflect.TypeOf("")))
		default:
			return &parser.PosError***REMOVED***
				Pos: node.Pos(),
				Err: fmt.Errorf("%s: cannot decode into interface: %T", name, node),
			***REMOVED***
		***REMOVED***
	default:
		return fmt.Errorf(
			"%s: cannot decode into interface: %T",
			name, node)
	***REMOVED***

	// Set the result to what its supposed to be, then reset
	// result so we don't reflect into this method anymore.
	result.Set(set)

	if redecode ***REMOVED***
		// Revisit the node so that we can use the newly instantiated
		// thing and populate it.
		if err := d.decode(name, node, result); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *decoder) decodeMap(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	if item, ok := node.(*ast.ObjectItem); ok ***REMOVED***
		node = &ast.ObjectList***REMOVED***Items: []*ast.ObjectItem***REMOVED***item***REMOVED******REMOVED***
	***REMOVED***

	if ot, ok := node.(*ast.ObjectType); ok ***REMOVED***
		node = ot.List
	***REMOVED***

	n, ok := node.(*ast.ObjectList)
	if !ok ***REMOVED***
		return &parser.PosError***REMOVED***
			Pos: node.Pos(),
			Err: fmt.Errorf("%s: not an object type for map (%T)", name, node),
		***REMOVED***
	***REMOVED***

	// If we have an interface, then we can address the interface,
	// but not the slice itself, so get the element but set the interface
	set := result
	if result.Kind() == reflect.Interface ***REMOVED***
		result = result.Elem()
	***REMOVED***

	resultType := result.Type()
	resultElemType := resultType.Elem()
	resultKeyType := resultType.Key()
	if resultKeyType.Kind() != reflect.String ***REMOVED***
		return &parser.PosError***REMOVED***
			Pos: node.Pos(),
			Err: fmt.Errorf("%s: map must have string keys", name),
		***REMOVED***
	***REMOVED***

	// Make a map if it is nil
	resultMap := result
	if result.IsNil() ***REMOVED***
		resultMap = reflect.MakeMap(
			reflect.MapOf(resultKeyType, resultElemType))
	***REMOVED***

	// Go through each element and decode it.
	done := make(map[string]struct***REMOVED******REMOVED***)
	for _, item := range n.Items ***REMOVED***
		if item.Val == nil ***REMOVED***
			continue
		***REMOVED***

		// github.com/hashicorp/terraform/issue/5740
		if len(item.Keys) == 0 ***REMOVED***
			return &parser.PosError***REMOVED***
				Pos: node.Pos(),
				Err: fmt.Errorf("%s: map must have string keys", name),
			***REMOVED***
		***REMOVED***

		// Get the key we're dealing with, which is the first item
		keyStr := item.Keys[0].Token.Value().(string)

		// If we've already processed this key, then ignore it
		if _, ok := done[keyStr]; ok ***REMOVED***
			continue
		***REMOVED***

		// Determine the value. If we have more than one key, then we
		// get the objectlist of only these keys.
		itemVal := item.Val
		if len(item.Keys) > 1 ***REMOVED***
			itemVal = n.Filter(keyStr)
			done[keyStr] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***

		// Make the field name
		fieldName := fmt.Sprintf("%s.%s", name, keyStr)

		// Get the key/value as reflection values
		key := reflect.ValueOf(keyStr)
		val := reflect.Indirect(reflect.New(resultElemType))

		// If we have a pre-existing value in the map, use that
		oldVal := resultMap.MapIndex(key)
		if oldVal.IsValid() ***REMOVED***
			val.Set(oldVal)
		***REMOVED***

		// Decode!
		if err := d.decode(fieldName, itemVal, val); err != nil ***REMOVED***
			return err
		***REMOVED***

		// Set the value on the map
		resultMap.SetMapIndex(key, val)
	***REMOVED***

	// Set the final map if we can
	set.Set(resultMap)
	return nil
***REMOVED***

func (d *decoder) decodePtr(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	resultType := result.Type()
	resultElemType := resultType.Elem()
	val := reflect.New(resultElemType)
	if err := d.decode(name, node, reflect.Indirect(val)); err != nil ***REMOVED***
		return err
	***REMOVED***

	result.Set(val)
	return nil
***REMOVED***

func (d *decoder) decodeSlice(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	// If we have an interface, then we can address the interface,
	// but not the slice itself, so get the element but set the interface
	set := result
	if result.Kind() == reflect.Interface ***REMOVED***
		result = result.Elem()
	***REMOVED***
	// Create the slice if it isn't nil
	resultType := result.Type()
	resultElemType := resultType.Elem()
	if result.IsNil() ***REMOVED***
		resultSliceType := reflect.SliceOf(resultElemType)
		result = reflect.MakeSlice(
			resultSliceType, 0, 0)
	***REMOVED***

	// Figure out the items we'll be copying into the slice
	var items []ast.Node
	switch n := node.(type) ***REMOVED***
	case *ast.ObjectList:
		items = make([]ast.Node, len(n.Items))
		for i, item := range n.Items ***REMOVED***
			items[i] = item
		***REMOVED***
	case *ast.ObjectType:
		items = []ast.Node***REMOVED***n***REMOVED***
	case *ast.ListType:
		items = n.List
	default:
		return &parser.PosError***REMOVED***
			Pos: node.Pos(),
			Err: fmt.Errorf("unknown slice type: %T", node),
		***REMOVED***
	***REMOVED***

	for i, item := range items ***REMOVED***
		fieldName := fmt.Sprintf("%s[%d]", name, i)

		// Decode
		val := reflect.Indirect(reflect.New(resultElemType))

		// if item is an object that was decoded from ambiguous JSON and
		// flattened, make sure it's expanded if it needs to decode into a
		// defined structure.
		item := expandObject(item, val)

		if err := d.decode(fieldName, item, val); err != nil ***REMOVED***
			return err
		***REMOVED***

		// Append it onto the slice
		result = reflect.Append(result, val)
	***REMOVED***

	set.Set(result)
	return nil
***REMOVED***

// expandObject detects if an ambiguous JSON object was flattened to a List which
// should be decoded into a struct, and expands the ast to properly deocode.
func expandObject(node ast.Node, result reflect.Value) ast.Node ***REMOVED***
	item, ok := node.(*ast.ObjectItem)
	if !ok ***REMOVED***
		return node
	***REMOVED***

	elemType := result.Type()

	// our target type must be a struct
	switch elemType.Kind() ***REMOVED***
	case reflect.Ptr:
		switch elemType.Elem().Kind() ***REMOVED***
		case reflect.Struct:
			//OK
		default:
			return node
		***REMOVED***
	case reflect.Struct:
		//OK
	default:
		return node
	***REMOVED***

	// A list value will have a key and field name. If it had more fields,
	// it wouldn't have been flattened.
	if len(item.Keys) != 2 ***REMOVED***
		return node
	***REMOVED***

	keyToken := item.Keys[0].Token
	item.Keys = item.Keys[1:]

	// we need to un-flatten the ast enough to decode
	newNode := &ast.ObjectItem***REMOVED***
		Keys: []*ast.ObjectKey***REMOVED***
			&ast.ObjectKey***REMOVED***
				Token: keyToken,
			***REMOVED***,
		***REMOVED***,
		Val: &ast.ObjectType***REMOVED***
			List: &ast.ObjectList***REMOVED***
				Items: []*ast.ObjectItem***REMOVED***item***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	return newNode
***REMOVED***

func (d *decoder) decodeString(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	switch n := node.(type) ***REMOVED***
	case *ast.LiteralType:
		switch n.Token.Type ***REMOVED***
		case token.NUMBER:
			result.Set(reflect.ValueOf(n.Token.Text).Convert(result.Type()))
			return nil
		case token.STRING, token.HEREDOC:
			result.Set(reflect.ValueOf(n.Token.Value()).Convert(result.Type()))
			return nil
		***REMOVED***
	***REMOVED***

	return &parser.PosError***REMOVED***
		Pos: node.Pos(),
		Err: fmt.Errorf("%s: unknown type for string %T", name, node),
	***REMOVED***
***REMOVED***

func (d *decoder) decodeStruct(name string, node ast.Node, result reflect.Value) error ***REMOVED***
	var item *ast.ObjectItem
	if it, ok := node.(*ast.ObjectItem); ok ***REMOVED***
		item = it
		node = it.Val
	***REMOVED***

	if ot, ok := node.(*ast.ObjectType); ok ***REMOVED***
		node = ot.List
	***REMOVED***

	// Handle the special case where the object itself is a literal. Previously
	// the yacc parser would always ensure top-level elements were arrays. The new
	// parser does not make the same guarantees, thus we need to convert any
	// top-level literal elements into a list.
	if _, ok := node.(*ast.LiteralType); ok && item != nil ***REMOVED***
		node = &ast.ObjectList***REMOVED***Items: []*ast.ObjectItem***REMOVED***item***REMOVED******REMOVED***
	***REMOVED***

	list, ok := node.(*ast.ObjectList)
	if !ok ***REMOVED***
		return &parser.PosError***REMOVED***
			Pos: node.Pos(),
			Err: fmt.Errorf("%s: not an object type for struct (%T)", name, node),
		***REMOVED***
	***REMOVED***

	// This slice will keep track of all the structs we'll be decoding.
	// There can be more than one struct if there are embedded structs
	// that are squashed.
	structs := make([]reflect.Value, 1, 5)
	structs[0] = result

	// Compile the list of all the fields that we're going to be decoding
	// from all the structs.
	type field struct ***REMOVED***
		field reflect.StructField
		val   reflect.Value
	***REMOVED***
	fields := []field***REMOVED******REMOVED***
	for len(structs) > 0 ***REMOVED***
		structVal := structs[0]
		structs = structs[1:]

		structType := structVal.Type()
		for i := 0; i < structType.NumField(); i++ ***REMOVED***
			fieldType := structType.Field(i)
			tagParts := strings.Split(fieldType.Tag.Get(tagName), ",")

			// Ignore fields with tag name "-"
			if tagParts[0] == "-" ***REMOVED***
				continue
			***REMOVED***

			if fieldType.Anonymous ***REMOVED***
				fieldKind := fieldType.Type.Kind()
				if fieldKind != reflect.Struct ***REMOVED***
					return &parser.PosError***REMOVED***
						Pos: node.Pos(),
						Err: fmt.Errorf("%s: unsupported type to struct: %s",
							fieldType.Name, fieldKind),
					***REMOVED***
				***REMOVED***

				// We have an embedded field. We "squash" the fields down
				// if specified in the tag.
				squash := false
				for _, tag := range tagParts[1:] ***REMOVED***
					if tag == "squash" ***REMOVED***
						squash = true
						break
					***REMOVED***
				***REMOVED***

				if squash ***REMOVED***
					structs = append(
						structs, result.FieldByName(fieldType.Name))
					continue
				***REMOVED***
			***REMOVED***

			// Normal struct field, store it away
			fields = append(fields, field***REMOVED***fieldType, structVal.Field(i)***REMOVED***)
		***REMOVED***
	***REMOVED***

	usedKeys := make(map[string]struct***REMOVED******REMOVED***)
	decodedFields := make([]string, 0, len(fields))
	decodedFieldsVal := make([]reflect.Value, 0)
	unusedKeysVal := make([]reflect.Value, 0)
	for _, f := range fields ***REMOVED***
		field, fieldValue := f.field, f.val
		if !fieldValue.IsValid() ***REMOVED***
			// This should never happen
			panic("field is not valid")
		***REMOVED***

		// If we can't set the field, then it is unexported or something,
		// and we just continue onwards.
		if !fieldValue.CanSet() ***REMOVED***
			continue
		***REMOVED***

		fieldName := field.Name

		tagValue := field.Tag.Get(tagName)
		tagParts := strings.SplitN(tagValue, ",", 2)
		if len(tagParts) >= 2 ***REMOVED***
			switch tagParts[1] ***REMOVED***
			case "decodedFields":
				decodedFieldsVal = append(decodedFieldsVal, fieldValue)
				continue
			case "key":
				if item == nil ***REMOVED***
					return &parser.PosError***REMOVED***
						Pos: node.Pos(),
						Err: fmt.Errorf("%s: %s asked for 'key', impossible",
							name, fieldName),
					***REMOVED***
				***REMOVED***

				fieldValue.SetString(item.Keys[0].Token.Value().(string))
				continue
			case "unusedKeys":
				unusedKeysVal = append(unusedKeysVal, fieldValue)
				continue
			***REMOVED***
		***REMOVED***

		if tagParts[0] != "" ***REMOVED***
			fieldName = tagParts[0]
		***REMOVED***

		// Determine the element we'll use to decode. If it is a single
		// match (only object with the field), then we decode it exactly.
		// If it is a prefix match, then we decode the matches.
		filter := list.Filter(fieldName)

		prefixMatches := filter.Children()
		matches := filter.Elem()
		if len(matches.Items) == 0 && len(prefixMatches.Items) == 0 ***REMOVED***
			continue
		***REMOVED***

		// Track the used key
		usedKeys[fieldName] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

		// Create the field name and decode. We range over the elements
		// because we actually want the value.
		fieldName = fmt.Sprintf("%s.%s", name, fieldName)
		if len(prefixMatches.Items) > 0 ***REMOVED***
			if err := d.decode(fieldName, prefixMatches, fieldValue); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		for _, match := range matches.Items ***REMOVED***
			var decodeNode ast.Node = match.Val
			if ot, ok := decodeNode.(*ast.ObjectType); ok ***REMOVED***
				decodeNode = &ast.ObjectList***REMOVED***Items: ot.List.Items***REMOVED***
			***REMOVED***

			if err := d.decode(fieldName, decodeNode, fieldValue); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		decodedFields = append(decodedFields, field.Name)
	***REMOVED***

	if len(decodedFieldsVal) > 0 ***REMOVED***
		// Sort it so that it is deterministic
		sort.Strings(decodedFields)

		for _, v := range decodedFieldsVal ***REMOVED***
			v.Set(reflect.ValueOf(decodedFields))
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// findNodeType returns the type of ast.Node
func findNodeType() reflect.Type ***REMOVED***
	var nodeContainer struct ***REMOVED***
		Node ast.Node
	***REMOVED***
	value := reflect.ValueOf(nodeContainer).FieldByName("Node")
	return value.Type()
***REMOVED***
