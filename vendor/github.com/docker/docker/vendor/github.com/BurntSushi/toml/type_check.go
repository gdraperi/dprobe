package toml

// tomlType represents any Go type that corresponds to a TOML type.
// While the first draft of the TOML spec has a simplistic type system that
// probably doesn't need this level of sophistication, we seem to be militating
// toward adding real composite types.
type tomlType interface ***REMOVED***
	typeString() string
***REMOVED***

// typeEqual accepts any two types and returns true if they are equal.
func typeEqual(t1, t2 tomlType) bool ***REMOVED***
	if t1 == nil || t2 == nil ***REMOVED***
		return false
	***REMOVED***
	return t1.typeString() == t2.typeString()
***REMOVED***

func typeIsHash(t tomlType) bool ***REMOVED***
	return typeEqual(t, tomlHash) || typeEqual(t, tomlArrayHash)
***REMOVED***

type tomlBaseType string

func (btype tomlBaseType) typeString() string ***REMOVED***
	return string(btype)
***REMOVED***

func (btype tomlBaseType) String() string ***REMOVED***
	return btype.typeString()
***REMOVED***

var (
	tomlInteger   tomlBaseType = "Integer"
	tomlFloat     tomlBaseType = "Float"
	tomlDatetime  tomlBaseType = "Datetime"
	tomlString    tomlBaseType = "String"
	tomlBool      tomlBaseType = "Bool"
	tomlArray     tomlBaseType = "Array"
	tomlHash      tomlBaseType = "Hash"
	tomlArrayHash tomlBaseType = "ArrayHash"
)

// typeOfPrimitive returns a tomlType of any primitive value in TOML.
// Primitive values are: Integer, Float, Datetime, String and Bool.
//
// Passing a lexer item other than the following will cause a BUG message
// to occur: itemString, itemBool, itemInteger, itemFloat, itemDatetime.
func (p *parser) typeOfPrimitive(lexItem item) tomlType ***REMOVED***
	switch lexItem.typ ***REMOVED***
	case itemInteger:
		return tomlInteger
	case itemFloat:
		return tomlFloat
	case itemDatetime:
		return tomlDatetime
	case itemString:
		return tomlString
	case itemMultilineString:
		return tomlString
	case itemRawString:
		return tomlString
	case itemRawMultilineString:
		return tomlString
	case itemBool:
		return tomlBool
	***REMOVED***
	p.bug("Cannot infer primitive type of lex item '%s'.", lexItem)
	panic("unreachable")
***REMOVED***

// typeOfArray returns a tomlType for an array given a list of types of its
// values.
//
// In the current spec, if an array is homogeneous, then its type is always
// "Array". If the array is not homogeneous, an error is generated.
func (p *parser) typeOfArray(types []tomlType) tomlType ***REMOVED***
	// Empty arrays are cool.
	if len(types) == 0 ***REMOVED***
		return tomlArray
	***REMOVED***

	theType := types[0]
	for _, t := range types[1:] ***REMOVED***
		if !typeEqual(theType, t) ***REMOVED***
			p.panicf("Array contains values of type '%s' and '%s', but "+
				"arrays must be homogeneous.", theType, t)
		***REMOVED***
	***REMOVED***
	return tomlArray
***REMOVED***
