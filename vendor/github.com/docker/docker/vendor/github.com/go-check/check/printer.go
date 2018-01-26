package check

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func indent(s, with string) (r string) ***REMOVED***
	eol := true
	for i := 0; i != len(s); i++ ***REMOVED***
		c := s[i]
		switch ***REMOVED***
		case eol && c == '\n' || c == '\r':
		case c == '\n' || c == '\r':
			eol = true
		case eol:
			eol = false
			s = s[:i] + with + s[i:]
			i += len(with)
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***

func printLine(filename string, line int) (string, error) ***REMOVED***
	fset := token.NewFileSet()
	file, err := os.Open(filename)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	fnode, err := parser.ParseFile(fset, filename, file, parser.ParseComments)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	config := &printer.Config***REMOVED***Mode: printer.UseSpaces, Tabwidth: 4***REMOVED***
	lp := &linePrinter***REMOVED***fset: fset, fnode: fnode, line: line, config: config***REMOVED***
	ast.Walk(lp, fnode)
	result := lp.output.Bytes()
	// Comments leave \n at the end.
	n := len(result)
	for n > 0 && result[n-1] == '\n' ***REMOVED***
		n--
	***REMOVED***
	return string(result[:n]), nil
***REMOVED***

type linePrinter struct ***REMOVED***
	config *printer.Config
	fset   *token.FileSet
	fnode  *ast.File
	line   int
	output bytes.Buffer
	stmt   ast.Stmt
***REMOVED***

func (lp *linePrinter) emit() bool ***REMOVED***
	if lp.stmt != nil ***REMOVED***
		lp.trim(lp.stmt)
		lp.printWithComments(lp.stmt)
		lp.stmt = nil
		return true
	***REMOVED***
	return false
***REMOVED***

func (lp *linePrinter) printWithComments(n ast.Node) ***REMOVED***
	nfirst := lp.fset.Position(n.Pos()).Line
	nlast := lp.fset.Position(n.End()).Line
	for _, g := range lp.fnode.Comments ***REMOVED***
		cfirst := lp.fset.Position(g.Pos()).Line
		clast := lp.fset.Position(g.End()).Line
		if clast == nfirst-1 && lp.fset.Position(n.Pos()).Column == lp.fset.Position(g.Pos()).Column ***REMOVED***
			for _, c := range g.List ***REMOVED***
				lp.output.WriteString(c.Text)
				lp.output.WriteByte('\n')
			***REMOVED***
		***REMOVED***
		if cfirst >= nfirst && cfirst <= nlast && n.End() <= g.List[0].Slash ***REMOVED***
			// The printer will not include the comment if it starts past
			// the node itself. Trick it into printing by overlapping the
			// slash with the end of the statement.
			g.List[0].Slash = n.End() - 1
		***REMOVED***
	***REMOVED***
	node := &printer.CommentedNode***REMOVED***n, lp.fnode.Comments***REMOVED***
	lp.config.Fprint(&lp.output, lp.fset, node)
***REMOVED***

func (lp *linePrinter) Visit(n ast.Node) (w ast.Visitor) ***REMOVED***
	if n == nil ***REMOVED***
		if lp.output.Len() == 0 ***REMOVED***
			lp.emit()
		***REMOVED***
		return nil
	***REMOVED***
	first := lp.fset.Position(n.Pos()).Line
	last := lp.fset.Position(n.End()).Line
	if first <= lp.line && last >= lp.line ***REMOVED***
		// Print the innermost statement containing the line.
		if stmt, ok := n.(ast.Stmt); ok ***REMOVED***
			if _, ok := n.(*ast.BlockStmt); !ok ***REMOVED***
				lp.stmt = stmt
			***REMOVED***
		***REMOVED***
		if first == lp.line && lp.emit() ***REMOVED***
			return nil
		***REMOVED***
		return lp
	***REMOVED***
	return nil
***REMOVED***

func (lp *linePrinter) trim(n ast.Node) bool ***REMOVED***
	stmt, ok := n.(ast.Stmt)
	if !ok ***REMOVED***
		return true
	***REMOVED***
	line := lp.fset.Position(n.Pos()).Line
	if line != lp.line ***REMOVED***
		return false
	***REMOVED***
	switch stmt := stmt.(type) ***REMOVED***
	case *ast.IfStmt:
		stmt.Body = lp.trimBlock(stmt.Body)
	case *ast.SwitchStmt:
		stmt.Body = lp.trimBlock(stmt.Body)
	case *ast.TypeSwitchStmt:
		stmt.Body = lp.trimBlock(stmt.Body)
	case *ast.CaseClause:
		stmt.Body = lp.trimList(stmt.Body)
	case *ast.CommClause:
		stmt.Body = lp.trimList(stmt.Body)
	case *ast.BlockStmt:
		stmt.List = lp.trimList(stmt.List)
	***REMOVED***
	return true
***REMOVED***

func (lp *linePrinter) trimBlock(stmt *ast.BlockStmt) *ast.BlockStmt ***REMOVED***
	if !lp.trim(stmt) ***REMOVED***
		return lp.emptyBlock(stmt)
	***REMOVED***
	stmt.Rbrace = stmt.Lbrace
	return stmt
***REMOVED***

func (lp *linePrinter) trimList(stmts []ast.Stmt) []ast.Stmt ***REMOVED***
	for i := 0; i != len(stmts); i++ ***REMOVED***
		if !lp.trim(stmts[i]) ***REMOVED***
			stmts[i] = lp.emptyStmt(stmts[i])
			break
		***REMOVED***
	***REMOVED***
	return stmts
***REMOVED***

func (lp *linePrinter) emptyStmt(n ast.Node) *ast.ExprStmt ***REMOVED***
	return &ast.ExprStmt***REMOVED***&ast.Ellipsis***REMOVED***n.Pos(), nil***REMOVED******REMOVED***
***REMOVED***

func (lp *linePrinter) emptyBlock(n ast.Node) *ast.BlockStmt ***REMOVED***
	p := n.Pos()
	return &ast.BlockStmt***REMOVED***p, []ast.Stmt***REMOVED***lp.emptyStmt(n)***REMOVED***, p***REMOVED***
***REMOVED***
