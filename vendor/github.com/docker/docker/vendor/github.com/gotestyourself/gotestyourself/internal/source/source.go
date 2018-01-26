package source

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"runtime"

	"github.com/pkg/errors"
)

const baseStackIndex = 1

// GetCondition returns the condition string by reading it from the file
// identified in the callstack. In golang 1.9 the line number changed from
// being the line where the statement ended to the line where the statement began.
func GetCondition(stackIndex int, argPos int) (string, error) ***REMOVED***
	_, filename, lineNum, ok := runtime.Caller(baseStackIndex + stackIndex)
	if !ok ***REMOVED***
		return "", errors.New("failed to get caller info")
	***REMOVED***

	node, err := getNodeAtLine(filename, lineNum)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return getArgSourceFromAST(node, argPos)
***REMOVED***

func getNodeAtLine(filename string, lineNum int) (ast.Node, error) ***REMOVED***
	fileset := token.NewFileSet()
	astFile, err := parser.ParseFile(fileset, filename, nil, parser.AllErrors)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to parse source file: %s", filename)
	***REMOVED***

	node := scanToLine(fileset, astFile, lineNum)
	if node == nil ***REMOVED***
		return nil, errors.Wrapf(err,
			"failed to find an expression on line %d in %s", lineNum, filename)
	***REMOVED***
	return node, nil
***REMOVED***

func scanToLine(fileset *token.FileSet, node ast.Node, lineNum int) ast.Node ***REMOVED***
	v := &scanToLineVisitor***REMOVED***lineNum: lineNum, fileset: fileset***REMOVED***
	ast.Walk(v, node)
	return v.matchedNode
***REMOVED***

type scanToLineVisitor struct ***REMOVED***
	lineNum     int
	matchedNode ast.Node
	fileset     *token.FileSet
***REMOVED***

func (v *scanToLineVisitor) Visit(node ast.Node) ast.Visitor ***REMOVED***
	if node == nil || v.matchedNode != nil ***REMOVED***
		return nil
	***REMOVED***

	var position token.Position
	switch ***REMOVED***
	case runtime.Version() < "go1.9":
		position = v.fileset.Position(node.End())
	default:
		position = v.fileset.Position(node.Pos())
	***REMOVED***

	if position.Line == v.lineNum ***REMOVED***
		v.matchedNode = node
		return nil
	***REMOVED***
	return v
***REMOVED***

func getArgSourceFromAST(node ast.Node, argPos int) (string, error) ***REMOVED***
	visitor := &callExprVisitor***REMOVED******REMOVED***
	ast.Walk(visitor, node)
	if visitor.expr == nil ***REMOVED***
		return "", errors.Errorf("unexpected ast")
	***REMOVED***

	buf := new(bytes.Buffer)
	err := format.Node(buf, token.NewFileSet(), visitor.expr.Args[argPos])
	return buf.String(), err
***REMOVED***

type callExprVisitor struct ***REMOVED***
	expr *ast.CallExpr
***REMOVED***

func (v *callExprVisitor) Visit(node ast.Node) ast.Visitor ***REMOVED***
	switch typed := node.(type) ***REMOVED***
	case nil:
		return nil
	case *ast.IfStmt:
		ast.Walk(v, typed.Cond)
	case *ast.CallExpr:
		v.expr = typed
	***REMOVED***

	if v.expr != nil ***REMOVED***
		return nil
	***REMOVED***
	return v
***REMOVED***
