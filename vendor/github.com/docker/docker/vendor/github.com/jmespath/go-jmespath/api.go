package jmespath

import "strconv"

// JmesPath is the epresentation of a compiled JMES path query. A JmesPath is
// safe for concurrent use by multiple goroutines.
type JMESPath struct ***REMOVED***
	ast  ASTNode
	intr *treeInterpreter
***REMOVED***

// Compile parses a JMESPath expression and returns, if successful, a JMESPath
// object that can be used to match against data.
func Compile(expression string) (*JMESPath, error) ***REMOVED***
	parser := NewParser()
	ast, err := parser.Parse(expression)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	jmespath := &JMESPath***REMOVED***ast: ast, intr: newInterpreter()***REMOVED***
	return jmespath, nil
***REMOVED***

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled
// JMESPaths.
func MustCompile(expression string) *JMESPath ***REMOVED***
	jmespath, err := Compile(expression)
	if err != nil ***REMOVED***
		panic(`jmespath: Compile(` + strconv.Quote(expression) + `): ` + err.Error())
	***REMOVED***
	return jmespath
***REMOVED***

// Search evaluates a JMESPath expression against input data and returns the result.
func (jp *JMESPath) Search(data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return jp.intr.Execute(jp.ast, data)
***REMOVED***

// Search evaluates a JMESPath expression against input data and returns the result.
func Search(expression string, data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	intr := newInterpreter()
	parser := NewParser()
	ast, err := parser.Parse(expression)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return intr.Execute(ast, data)
***REMOVED***
