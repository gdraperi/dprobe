// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipeline

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/go/loader"
)

const printerType = "golang.org/x/text/message.Printer"

// Rewrite rewrites the Go files in a single package to use the localization
// machinery and rewrites strings to adopt best practices when possible.
// If w is not nil the generated files are written to it, each files with a
// "--- <filename>" header. Otherwise the files are overwritten.
func Rewrite(w io.Writer, args ...string) error ***REMOVED***
	conf := &loader.Config***REMOVED***
		AllowErrors: true, // Allow unused instances of message.Printer.
	***REMOVED***
	prog, err := loadPackages(conf, args)
	if err != nil ***REMOVED***
		return wrap(err, "")
	***REMOVED***

	for _, info := range prog.InitialPackages() ***REMOVED***
		for _, f := range info.Files ***REMOVED***
			// Associate comments with nodes.

			// Pick up initialized Printers at the package level.
			r := rewriter***REMOVED***info: info, conf: conf***REMOVED***
			for _, n := range info.InitOrder ***REMOVED***
				if t := r.info.Types[n.Rhs].Type.String(); strings.HasSuffix(t, printerType) ***REMOVED***
					r.printerVar = n.Lhs[0].Name()
				***REMOVED***
			***REMOVED***

			ast.Walk(&r, f)

			w := w
			if w == nil ***REMOVED***
				var err error
				if w, err = os.Create(conf.Fset.File(f.Pos()).Name()); err != nil ***REMOVED***
					return wrap(err, "open failed")
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				fmt.Fprintln(w, "---", conf.Fset.File(f.Pos()).Name())
			***REMOVED***

			if err := format.Node(w, conf.Fset, f); err != nil ***REMOVED***
				return wrap(err, "go format failed")
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type rewriter struct ***REMOVED***
	info       *loader.PackageInfo
	conf       *loader.Config
	printerVar string
***REMOVED***

// print returns Go syntax for the specified node.
func (r *rewriter) print(n ast.Node) string ***REMOVED***
	var buf bytes.Buffer
	format.Node(&buf, r.conf.Fset, n)
	return buf.String()
***REMOVED***

func (r *rewriter) Visit(n ast.Node) ast.Visitor ***REMOVED***
	// Save the state by scope.
	if _, ok := n.(*ast.BlockStmt); ok ***REMOVED***
		r := *r
		return &r
	***REMOVED***
	// Find Printers created by assignment.
	stmt, ok := n.(*ast.AssignStmt)
	if ok ***REMOVED***
		for _, v := range stmt.Lhs ***REMOVED***
			if r.printerVar == r.print(v) ***REMOVED***
				r.printerVar = ""
			***REMOVED***
		***REMOVED***
		for i, v := range stmt.Rhs ***REMOVED***
			if t := r.info.Types[v].Type.String(); strings.HasSuffix(t, printerType) ***REMOVED***
				r.printerVar = r.print(stmt.Lhs[i])
				return r
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Find Printers created by variable declaration.
	spec, ok := n.(*ast.ValueSpec)
	if ok ***REMOVED***
		for _, v := range spec.Names ***REMOVED***
			if r.printerVar == r.print(v) ***REMOVED***
				r.printerVar = ""
			***REMOVED***
		***REMOVED***
		for i, v := range spec.Values ***REMOVED***
			if t := r.info.Types[v].Type.String(); strings.HasSuffix(t, printerType) ***REMOVED***
				r.printerVar = r.print(spec.Names[i])
				return r
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if r.printerVar == "" ***REMOVED***
		return r
	***REMOVED***
	call, ok := n.(*ast.CallExpr)
	if !ok ***REMOVED***
		return r
	***REMOVED***

	// TODO: Handle literal values?
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok ***REMOVED***
		return r
	***REMOVED***
	meth := r.info.Selections[sel]

	source := r.print(sel.X)
	fun := r.print(sel.Sel)
	if meth != nil ***REMOVED***
		source = meth.Recv().String()
		fun = meth.Obj().Name()
	***REMOVED***

	// TODO: remove cheap hack and check if the type either
	// implements some interface or is specifically of type
	// "golang.org/x/text/message".Printer.
	m, ok := rewriteFuncs[source]
	if !ok ***REMOVED***
		return r
	***REMOVED***

	rewriteType, ok := m[fun]
	if !ok ***REMOVED***
		return r
	***REMOVED***
	ident := ast.NewIdent(r.printerVar)
	ident.NamePos = sel.X.Pos()
	sel.X = ident
	if rewriteType.method != "" ***REMOVED***
		sel.Sel.Name = rewriteType.method
	***REMOVED***

	// Analyze arguments.
	argn := rewriteType.arg
	if rewriteType.format || argn >= len(call.Args) ***REMOVED***
		return r
	***REMOVED***
	hasConst := false
	for _, a := range call.Args[argn:] ***REMOVED***
		if v := r.info.Types[a].Value; v != nil && v.Kind() == constant.String ***REMOVED***
			hasConst = true
			break
		***REMOVED***
	***REMOVED***
	if !hasConst ***REMOVED***
		return r
	***REMOVED***
	sel.Sel.Name = rewriteType.methodf

	// We are done if there is only a single string that does not need to be
	// escaped.
	if len(call.Args) == 1 ***REMOVED***
		s, ok := constStr(r.info, call.Args[0])
		if ok && !strings.Contains(s, "%") && !rewriteType.newLine ***REMOVED***
			return r
		***REMOVED***
	***REMOVED***

	// Rewrite arguments as format string.
	expr := &ast.BasicLit***REMOVED***
		ValuePos: call.Lparen,
		Kind:     token.STRING,
	***REMOVED***
	newArgs := append(call.Args[:argn:argn], expr)
	newStr := []string***REMOVED******REMOVED***
	for i, a := range call.Args[argn:] ***REMOVED***
		if s, ok := constStr(r.info, a); ok ***REMOVED***
			newStr = append(newStr, strings.Replace(s, "%", "%%", -1))
		***REMOVED*** else ***REMOVED***
			newStr = append(newStr, "%v")
			newArgs = append(newArgs, call.Args[argn+i])
		***REMOVED***
	***REMOVED***
	s := strings.Join(newStr, rewriteType.sep)
	if rewriteType.newLine ***REMOVED***
		s += "\n"
	***REMOVED***
	expr.Value = fmt.Sprintf("%q", s)

	call.Args = newArgs

	// TODO: consider creating an expression instead of a constant string and
	// then wrapping it in an escape function or so:
	// call.Args[argn+i] = &ast.CallExpr***REMOVED***
	// 		Fun: &ast.SelectorExpr***REMOVED***
	// 			X:   ast.NewIdent("message"),
	// 			Sel: ast.NewIdent("Lookup"),
	// 		***REMOVED***,
	// 		Args: []ast.Expr***REMOVED***a***REMOVED***,
	// 	***REMOVED***
	// ***REMOVED***

	return r
***REMOVED***

type rewriteType struct ***REMOVED***
	// method is the name of the equivalent method on a printer, or "" if it is
	// the same.
	method string

	// methodf is the method to use if the arguments can be rewritten as a
	// arguments to a printf-style call.
	methodf string

	// format is true if the method takes a formatting string followed by
	// substitution arguments.
	format bool

	// arg indicates the position of the argument to extract. If all is
	// positive, all arguments from this argument onwards needs to be extracted.
	arg int

	sep     string
	newLine bool
***REMOVED***

// rewriteFuncs list functions that can be directly mapped to the printer
// functions of the message package.
var rewriteFuncs = map[string]map[string]rewriteType***REMOVED***
	// TODO: Printer -> *golang.org/x/text/message.Printer
	"fmt": ***REMOVED***
		"Print":  rewriteType***REMOVED***methodf: "Printf"***REMOVED***,
		"Sprint": rewriteType***REMOVED***methodf: "Sprintf"***REMOVED***,
		"Fprint": rewriteType***REMOVED***methodf: "Fprintf"***REMOVED***,

		"Println":  rewriteType***REMOVED***methodf: "Printf", sep: " ", newLine: true***REMOVED***,
		"Sprintln": rewriteType***REMOVED***methodf: "Sprintf", sep: " ", newLine: true***REMOVED***,
		"Fprintln": rewriteType***REMOVED***methodf: "Fprintf", sep: " ", newLine: true***REMOVED***,

		"Printf":  rewriteType***REMOVED***method: "Printf", format: true***REMOVED***,
		"Sprintf": rewriteType***REMOVED***method: "Sprintf", format: true***REMOVED***,
		"Fprintf": rewriteType***REMOVED***method: "Fprintf", format: true***REMOVED***,
	***REMOVED***,
***REMOVED***

func constStr(info *loader.PackageInfo, e ast.Expr) (s string, ok bool) ***REMOVED***
	v := info.Types[e].Value
	if v == nil || v.Kind() != constant.String ***REMOVED***
		return "", false
	***REMOVED***
	return constant.StringVal(v), true
***REMOVED***
