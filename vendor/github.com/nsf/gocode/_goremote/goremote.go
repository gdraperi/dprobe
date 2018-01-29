package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"reflect"
	"strings"
)

const prefix = "server_"

func pretty_print_type_expr(out io.Writer, e ast.Expr) ***REMOVED***
	ty := reflect.TypeOf(e)
	switch t := e.(type) ***REMOVED***
	case *ast.StarExpr:
		fmt.Fprintf(out, "*")
		pretty_print_type_expr(out, t.X)
	case *ast.Ident:
		fmt.Fprintf(out, t.Name)
	case *ast.ArrayType:
		fmt.Fprintf(out, "[]")
		pretty_print_type_expr(out, t.Elt)
	case *ast.SelectorExpr:
		pretty_print_type_expr(out, t.X)
		fmt.Fprintf(out, ".%s", t.Sel.Name)
	case *ast.FuncType:
		fmt.Fprintf(out, "func(")
		pretty_print_func_field_list(out, t.Params)
		fmt.Fprintf(out, ")")

		buf := bytes.NewBuffer(make([]byte, 0, 256))
		nresults := pretty_print_func_field_list(buf, t.Results)
		if nresults > 0 ***REMOVED***
			results := buf.String()
			if strings.Index(results, " ") != -1 ***REMOVED***
				results = "(" + results + ")"
			***REMOVED***
			fmt.Fprintf(out, " %s", results)
		***REMOVED***
	case *ast.MapType:
		fmt.Fprintf(out, "map[")
		pretty_print_type_expr(out, t.Key)
		fmt.Fprintf(out, "]")
		pretty_print_type_expr(out, t.Value)
	case *ast.InterfaceType:
		fmt.Fprintf(out, "interface***REMOVED******REMOVED***")
	case *ast.Ellipsis:
		fmt.Fprintf(out, "...")
		pretty_print_type_expr(out, t.Elt)
	default:
		fmt.Fprintf(out, "\n[!!] unknown type: %s\n", ty.String())
	***REMOVED***
***REMOVED***

func pretty_print_func_field_list(out io.Writer, f *ast.FieldList) int ***REMOVED***
	count := 0
	if f == nil ***REMOVED***
		return count
	***REMOVED***
	for i, field := range f.List ***REMOVED***
		// names
		if field.Names != nil ***REMOVED***
			for j, name := range field.Names ***REMOVED***
				fmt.Fprintf(out, "%s", name.Name)
				if j != len(field.Names)-1 ***REMOVED***
					fmt.Fprintf(out, ", ")
				***REMOVED***
				count++
			***REMOVED***
			fmt.Fprintf(out, " ")
		***REMOVED*** else ***REMOVED***
			count++
		***REMOVED***

		// type
		pretty_print_type_expr(out, field.Type)

		// ,
		if i != len(f.List)-1 ***REMOVED***
			fmt.Fprintf(out, ", ")
		***REMOVED***
	***REMOVED***
	return count
***REMOVED***

func pretty_print_func_field_list_using_args(out io.Writer, f *ast.FieldList) int ***REMOVED***
	count := 0
	if f == nil ***REMOVED***
		return count
	***REMOVED***
	for i, field := range f.List ***REMOVED***
		// names
		if field.Names != nil ***REMOVED***
			for j := range field.Names ***REMOVED***
				fmt.Fprintf(out, "Arg%d", count)
				if j != len(field.Names)-1 ***REMOVED***
					fmt.Fprintf(out, ", ")
				***REMOVED***
				count++
			***REMOVED***
			fmt.Fprintf(out, " ")
		***REMOVED*** else ***REMOVED***
			count++
		***REMOVED***

		// type
		pretty_print_type_expr(out, field.Type)

		// ,
		if i != len(f.List)-1 ***REMOVED***
			fmt.Fprintf(out, ", ")
		***REMOVED***
	***REMOVED***
	return count
***REMOVED***

func generate_struct_wrapper(out io.Writer, fun *ast.FieldList, structname, name string) int ***REMOVED***
	fmt.Fprintf(out, "type %s_%s struct ***REMOVED***\n", structname, name)
	argn := 0
	for _, field := range fun.List ***REMOVED***
		fmt.Fprintf(out, "\t")
		// names
		if field.Names != nil ***REMOVED***
			for j := range field.Names ***REMOVED***
				fmt.Fprintf(out, "Arg%d", argn)
				if j != len(field.Names)-1 ***REMOVED***
					fmt.Fprintf(out, ", ")
				***REMOVED***
				argn++
			***REMOVED***
			fmt.Fprintf(out, " ")
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(out, "Arg%d ", argn)
			argn++
		***REMOVED***

		// type
		pretty_print_type_expr(out, field.Type)

		// \n
		fmt.Fprintf(out, "\n")
	***REMOVED***
	fmt.Fprintf(out, "***REMOVED***\n")
	return argn
***REMOVED***

// function that is being exposed to an RPC API, but calls simple "Server_" one
func generate_server_rpc_wrapper(out io.Writer, fun *ast.FuncDecl, name string, argcnt, replycnt int) ***REMOVED***
	fmt.Fprintf(out, "func (r *RPC) RPC_%s(args *Args_%s, reply *Reply_%s) error ***REMOVED***\n",
		name, name, name)

	fmt.Fprintf(out, "\t")
	for i := 0; i < replycnt; i++ ***REMOVED***
		fmt.Fprintf(out, "reply.Arg%d", i)
		if i != replycnt-1 ***REMOVED***
			fmt.Fprintf(out, ", ")
		***REMOVED***
	***REMOVED***
	fmt.Fprintf(out, " = %s(", fun.Name.Name)
	for i := 0; i < argcnt; i++ ***REMOVED***
		fmt.Fprintf(out, "args.Arg%d", i)
		if i != argcnt-1 ***REMOVED***
			fmt.Fprintf(out, ", ")
		***REMOVED***
	***REMOVED***
	fmt.Fprintf(out, ")\n")
	fmt.Fprintf(out, "\treturn nil\n***REMOVED***\n")
***REMOVED***

func generate_client_rpc_wrapper(out io.Writer, fun *ast.FuncDecl, name string, argcnt, replycnt int) ***REMOVED***
	fmt.Fprintf(out, "func client_%s(cli *rpc.Client, ", name)
	pretty_print_func_field_list_using_args(out, fun.Type.Params)
	fmt.Fprintf(out, ")")

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	nresults := pretty_print_func_field_list(buf, fun.Type.Results)
	if nresults > 0 ***REMOVED***
		results := buf.String()
		if strings.Index(results, " ") != -1 ***REMOVED***
			results = "(" + results + ")"
		***REMOVED***
		fmt.Fprintf(out, " %s", results)
	***REMOVED***
	fmt.Fprintf(out, " ***REMOVED***\n")
	fmt.Fprintf(out, "\tvar args Args_%s\n", name)
	fmt.Fprintf(out, "\tvar reply Reply_%s\n", name)
	for i := 0; i < argcnt; i++ ***REMOVED***
		fmt.Fprintf(out, "\targs.Arg%d = Arg%d\n", i, i)
	***REMOVED***
	fmt.Fprintf(out, "\terr := cli.Call(\"RPC.RPC_%s\", &args, &reply)\n", name)
	fmt.Fprintf(out, "\tif err != nil ***REMOVED***\n")
	fmt.Fprintf(out, "\t\tpanic(err)\n\t***REMOVED***\n")

	fmt.Fprintf(out, "\treturn ")
	for i := 0; i < replycnt; i++ ***REMOVED***
		fmt.Fprintf(out, "reply.Arg%d", i)
		if i != replycnt-1 ***REMOVED***
			fmt.Fprintf(out, ", ")
		***REMOVED***
	***REMOVED***
	fmt.Fprintf(out, "\n***REMOVED***\n")
***REMOVED***

func wrap_function(out io.Writer, fun *ast.FuncDecl) ***REMOVED***
	name := fun.Name.Name[len(prefix):]
	fmt.Fprintf(out, "// wrapper for: %s\n\n", fun.Name.Name)
	argcnt := generate_struct_wrapper(out, fun.Type.Params, "Args", name)
	replycnt := generate_struct_wrapper(out, fun.Type.Results, "Reply", name)
	generate_server_rpc_wrapper(out, fun, name, argcnt, replycnt)
	generate_client_rpc_wrapper(out, fun, name, argcnt, replycnt)
	fmt.Fprintf(out, "\n")
***REMOVED***

func process_file(out io.Writer, filename string) ***REMOVED***
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	for _, decl := range file.Decls ***REMOVED***
		if fdecl, ok := decl.(*ast.FuncDecl); ok ***REMOVED***
			namelen := len(fdecl.Name.Name)
			if namelen >= len(prefix) && fdecl.Name.Name[0:len(prefix)] == prefix ***REMOVED***
				wrap_function(out, fdecl)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

const head = `// WARNING! Autogenerated by goremote, don't touch.

package main

import (
	"net/rpc"
)

type RPC struct ***REMOVED***
***REMOVED***

`

func main() ***REMOVED***
	flag.Parse()
	fmt.Fprintf(os.Stdout, head)
	for _, file := range flag.Args() ***REMOVED***
		process_file(os.Stdout, file)
	***REMOVED***
***REMOVED***
