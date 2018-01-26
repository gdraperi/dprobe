package main

import (
	"strings"
	"text/template"
)

func printArgs(args []arg) string ***REMOVED***
	var argStr []string
	for _, arg := range args ***REMOVED***
		argStr = append(argStr, arg.String())
	***REMOVED***
	return strings.Join(argStr, ", ")
***REMOVED***

func buildImports(specs []importSpec) string ***REMOVED***
	if len(specs) == 0 ***REMOVED***
		return `import "errors"`
	***REMOVED***
	imports := "import(\n"
	imports += "\t\"errors\"\n"
	for _, i := range specs ***REMOVED***
		imports += "\t" + i.String() + "\n"
	***REMOVED***
	imports += ")"
	return imports
***REMOVED***

func marshalType(t string) string ***REMOVED***
	switch t ***REMOVED***
	case "error":
		// convert error types to plain strings to ensure the values are encoded/decoded properly
		return "string"
	default:
		return t
	***REMOVED***
***REMOVED***

func isErr(t string) bool ***REMOVED***
	switch t ***REMOVED***
	case "error":
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// Need to use this helper due to issues with go-vet
func buildTag(s string) string ***REMOVED***
	return "+build " + s
***REMOVED***

var templFuncs = template.FuncMap***REMOVED***
	"printArgs":   printArgs,
	"marshalType": marshalType,
	"isErr":       isErr,
	"lower":       strings.ToLower,
	"title":       title,
	"tag":         buildTag,
	"imports":     buildImports,
***REMOVED***

func title(s string) string ***REMOVED***
	if strings.ToLower(s) == "id" ***REMOVED***
		return "ID"
	***REMOVED***
	return strings.Title(s)
***REMOVED***

var generatedTempl = template.Must(template.New("rpc_cient").Funcs(templFuncs).Parse(`
// generated code - DO NOT EDIT
***REMOVED******REMOVED*** range $k, $v := .BuildTags ***REMOVED******REMOVED***
	// ***REMOVED******REMOVED*** tag $k ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** end ***REMOVED******REMOVED***

package ***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***

***REMOVED******REMOVED*** imports .Imports ***REMOVED******REMOVED***

type client interface***REMOVED***
	Call(string, interface***REMOVED******REMOVED***, interface***REMOVED******REMOVED***) error
***REMOVED***

type ***REMOVED******REMOVED*** .InterfaceType ***REMOVED******REMOVED***Proxy struct ***REMOVED***
	client
***REMOVED***

***REMOVED******REMOVED*** range .Functions ***REMOVED******REMOVED***
	type ***REMOVED******REMOVED*** $.InterfaceType ***REMOVED******REMOVED***Proxy***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***Request struct***REMOVED***
		***REMOVED******REMOVED*** range .Args ***REMOVED******REMOVED***
			***REMOVED******REMOVED*** title .Name ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .ArgType ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
	***REMOVED***

	type ***REMOVED******REMOVED*** $.InterfaceType ***REMOVED******REMOVED***Proxy***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***Response struct***REMOVED***
		***REMOVED******REMOVED*** range .Returns ***REMOVED******REMOVED***
			***REMOVED******REMOVED*** title .Name ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** marshalType .ArgType ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
	***REMOVED***

	func (pp ****REMOVED******REMOVED*** $.InterfaceType ***REMOVED******REMOVED***Proxy) ***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***(***REMOVED******REMOVED*** printArgs .Args ***REMOVED******REMOVED***) (***REMOVED******REMOVED*** printArgs .Returns ***REMOVED******REMOVED***) ***REMOVED***
		var(
			req ***REMOVED******REMOVED*** $.InterfaceType ***REMOVED******REMOVED***Proxy***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***Request
			ret ***REMOVED******REMOVED*** $.InterfaceType ***REMOVED******REMOVED***Proxy***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***Response
		)
		***REMOVED******REMOVED*** range .Args ***REMOVED******REMOVED***
			req.***REMOVED******REMOVED*** title .Name ***REMOVED******REMOVED*** = ***REMOVED******REMOVED*** lower .Name ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
		if err = pp.Call("***REMOVED******REMOVED*** $.RPCName ***REMOVED******REMOVED***.***REMOVED******REMOVED*** .Name ***REMOVED******REMOVED***", req, &ret); err != nil ***REMOVED***
			return
		***REMOVED***
		***REMOVED******REMOVED*** range $r := .Returns ***REMOVED******REMOVED***
			***REMOVED******REMOVED*** if isErr .ArgType ***REMOVED******REMOVED***
				if ret.***REMOVED******REMOVED*** title .Name ***REMOVED******REMOVED*** != "" ***REMOVED***
					***REMOVED******REMOVED*** lower .Name ***REMOVED******REMOVED*** = errors.New(ret.***REMOVED******REMOVED*** title .Name ***REMOVED******REMOVED***)
				***REMOVED*** ***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
			***REMOVED******REMOVED*** if isErr .ArgType | not ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** lower .Name ***REMOVED******REMOVED*** = ret.***REMOVED******REMOVED*** title .Name ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** end ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** end ***REMOVED******REMOVED***

		return
	***REMOVED***
***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
`))
