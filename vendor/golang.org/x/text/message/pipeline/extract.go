// Copyright 2016 The Go Authors. All rights reserved.
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
	"go/types"
	"path"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	fmtparser "golang.org/x/text/internal/format"
	"golang.org/x/tools/go/loader"
)

// TODO:
// - merge information into existing files
// - handle different file formats (PO, XLIFF)
// - handle features (gender, plural)
// - message rewriting

// - %m substitutions
// - `msg:"etc"` tags
// - msg/Msg top-level vars and strings.

// Extract extracts all strings form the package defined in Config.
func Extract(c *Config) (*State, error) ***REMOVED***
	conf := loader.Config***REMOVED******REMOVED***
	prog, err := loadPackages(&conf, c.Packages)
	if err != nil ***REMOVED***
		return nil, wrap(err, "")
	***REMOVED***

	// print returns Go syntax for the specified node.
	print := func(n ast.Node) string ***REMOVED***
		var buf bytes.Buffer
		format.Node(&buf, conf.Fset, n)
		return buf.String()
	***REMOVED***

	var messages []Message

	for _, info := range prog.AllPackages ***REMOVED***
		for _, f := range info.Files ***REMOVED***
			// Associate comments with nodes.
			cmap := ast.NewCommentMap(prog.Fset, f, f.Comments)
			getComment := func(n ast.Node) string ***REMOVED***
				cs := cmap.Filter(n).Comments()
				if len(cs) > 0 ***REMOVED***
					return strings.TrimSpace(cs[0].Text())
				***REMOVED***
				return ""
			***REMOVED***

			// Find function calls.
			ast.Inspect(f, func(n ast.Node) bool ***REMOVED***
				call, ok := n.(*ast.CallExpr)
				if !ok ***REMOVED***
					return true
				***REMOVED***

				// Skip calls of functions other than
				// (*message.Printer).***REMOVED***Sp,Fp,P***REMOVED***rintf.
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok ***REMOVED***
					return true
				***REMOVED***
				meth := info.Selections[sel]
				if meth == nil || meth.Kind() != types.MethodVal ***REMOVED***
					return true
				***REMOVED***
				// TODO: remove cheap hack and check if the type either
				// implements some interface or is specifically of type
				// "golang.org/x/text/message".Printer.
				m, ok := extractFuncs[path.Base(meth.Recv().String())]
				if !ok ***REMOVED***
					return true
				***REMOVED***

				fmtType, ok := m[meth.Obj().Name()]
				if !ok ***REMOVED***
					return true
				***REMOVED***
				// argn is the index of the format string.
				argn := fmtType.arg
				if argn >= len(call.Args) ***REMOVED***
					return true
				***REMOVED***

				args := call.Args[fmtType.arg:]

				fmtMsg, ok := msgStr(info, args[0])
				if !ok ***REMOVED***
					// TODO: identify the type of the format argument. If it
					// is not a string, multiple keys may be defined.
					return true
				***REMOVED***
				comment := ""
				key := []string***REMOVED******REMOVED***
				if ident, ok := args[0].(*ast.Ident); ok ***REMOVED***
					key = append(key, ident.Name)
					if v, ok := ident.Obj.Decl.(*ast.ValueSpec); ok && v.Comment != nil ***REMOVED***
						// TODO: get comment above ValueSpec as well
						comment = v.Comment.Text()
					***REMOVED***
				***REMOVED***

				arguments := []argument***REMOVED******REMOVED***
				args = args[1:]
				simArgs := make([]interface***REMOVED******REMOVED***, len(args))
				for i, arg := range args ***REMOVED***
					expr := print(arg)
					val := ""
					if v := info.Types[arg].Value; v != nil ***REMOVED***
						val = v.ExactString()
						simArgs[i] = val
						switch arg.(type) ***REMOVED***
						case *ast.BinaryExpr, *ast.UnaryExpr:
							expr = val
						***REMOVED***
					***REMOVED***
					arguments = append(arguments, argument***REMOVED***
						ArgNum:         i + 1,
						Type:           info.Types[arg].Type.String(),
						UnderlyingType: info.Types[arg].Type.Underlying().String(),
						Expr:           expr,
						Value:          val,
						Comment:        getComment(arg),
						Position:       posString(conf, info, arg.Pos()),
						// TODO report whether it implements
						// interfaces plural.Interface,
						// gender.Interface.
					***REMOVED***)
				***REMOVED***
				msg := ""

				ph := placeholders***REMOVED***index: map[string]string***REMOVED******REMOVED******REMOVED***

				trimmed, _, _ := trimWS(fmtMsg)

				p := fmtparser.Parser***REMOVED******REMOVED***
				p.Reset(simArgs)
				for p.SetFormat(trimmed); p.Scan(); ***REMOVED***
					switch p.Status ***REMOVED***
					case fmtparser.StatusText:
						msg += p.Text()
					case fmtparser.StatusSubstitution,
						fmtparser.StatusBadWidthSubstitution,
						fmtparser.StatusBadPrecSubstitution:
						arguments[p.ArgNum-1].used = true
						arg := arguments[p.ArgNum-1]
						sub := p.Text()
						if !p.HasIndex ***REMOVED***
							r, sz := utf8.DecodeLastRuneInString(sub)
							sub = fmt.Sprintf("%s[%d]%c", sub[:len(sub)-sz], p.ArgNum, r)
						***REMOVED***
						msg += fmt.Sprintf("***REMOVED***%s***REMOVED***", ph.addArg(&arg, sub))
					***REMOVED***
				***REMOVED***
				key = append(key, msg)

				// Add additional Placeholders that can be used in translations
				// that are not present in the string.
				for _, arg := range arguments ***REMOVED***
					if arg.used ***REMOVED***
						continue
					***REMOVED***
					ph.addArg(&arg, fmt.Sprintf("%%[%d]v", arg.ArgNum))
				***REMOVED***

				if c := getComment(call.Args[0]); c != "" ***REMOVED***
					comment = c
				***REMOVED***

				messages = append(messages, Message***REMOVED***
					ID:      key,
					Key:     fmtMsg,
					Message: Text***REMOVED***Msg: msg***REMOVED***,
					// TODO(fix): this doesn't get the before comment.
					Comment:      comment,
					Placeholders: ph.slice,
					Position:     posString(conf, info, call.Lparen),
				***REMOVED***)
				return true
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	return &State***REMOVED***
		Config:  *c,
		program: prog,
		Extracted: Messages***REMOVED***
			Language: c.SourceLanguage,
			Messages: messages,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

func posString(conf loader.Config, info *loader.PackageInfo, pos token.Pos) string ***REMOVED***
	p := conf.Fset.Position(pos)
	file := fmt.Sprintf("%s:%d:%d", filepath.Base(p.Filename), p.Line, p.Column)
	return filepath.Join(info.Pkg.Path(), file)
***REMOVED***

// extractFuncs indicates the types and methods for which to extract strings,
// and which argument to extract.
// TODO: use the types in conf.Import("golang.org/x/text/message") to extract
// the correct instances.
var extractFuncs = map[string]map[string]extractType***REMOVED***
	// TODO: Printer -> *golang.org/x/text/message.Printer
	"message.Printer": ***REMOVED***
		"Printf":  extractType***REMOVED***arg: 0, format: true***REMOVED***,
		"Sprintf": extractType***REMOVED***arg: 0, format: true***REMOVED***,
		"Fprintf": extractType***REMOVED***arg: 1, format: true***REMOVED***,

		"Lookup": extractType***REMOVED***arg: 0***REMOVED***,
	***REMOVED***,
***REMOVED***

type extractType struct ***REMOVED***
	// format indicates if the next arg is a formatted string or whether to
	// concatenate all arguments
	format bool
	// arg indicates the position of the argument to extract.
	arg int
***REMOVED***

func getID(arg *argument) string ***REMOVED***
	s := getLastComponent(arg.Expr)
	s = strip(s)
	s = strings.Replace(s, " ", "", -1)
	// For small variable names, use user-defined types for more info.
	if len(s) <= 2 && arg.UnderlyingType != arg.Type ***REMOVED***
		s = getLastComponent(arg.Type)
	***REMOVED***
	return strings.Title(s)
***REMOVED***

// strip is a dirty hack to convert function calls to placeholder IDs.
func strip(s string) string ***REMOVED***
	s = strings.Map(func(r rune) rune ***REMOVED***
		if unicode.IsSpace(r) || r == '-' ***REMOVED***
			return '_'
		***REMOVED***
		if !unicode.In(r, unicode.Letter, unicode.Mark, unicode.Number) ***REMOVED***
			return -1
		***REMOVED***
		return r
	***REMOVED***, s)
	// Strip "Get" from getter functions.
	if strings.HasPrefix(s, "Get") || strings.HasPrefix(s, "get") ***REMOVED***
		if len(s) > len("get") ***REMOVED***
			r, _ := utf8.DecodeRuneInString(s)
			if !unicode.In(r, unicode.Ll, unicode.M) ***REMOVED*** // not lower or mark
				s = s[len("get"):]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***

type placeholders struct ***REMOVED***
	index map[string]string
	slice []Placeholder
***REMOVED***

func (p *placeholders) addArg(arg *argument, sub string) (id string) ***REMOVED***
	id = getID(arg)
	id1 := id
	alt, ok := p.index[id1]
	for i := 1; ok && alt != sub; i++ ***REMOVED***
		id1 = fmt.Sprintf("%s_%d", id, i)
		alt, ok = p.index[id1]
	***REMOVED***
	p.index[id1] = sub
	p.slice = append(p.slice, Placeholder***REMOVED***
		ID:             id1,
		String:         sub,
		Type:           arg.Type,
		UnderlyingType: arg.UnderlyingType,
		ArgNum:         arg.ArgNum,
		Expr:           arg.Expr,
		Comment:        arg.Comment,
	***REMOVED***)
	return id1
***REMOVED***

func getLastComponent(s string) string ***REMOVED***
	return s[1+strings.LastIndexByte(s, '.'):]
***REMOVED***

func msgStr(info *loader.PackageInfo, e ast.Expr) (s string, ok bool) ***REMOVED***
	v := info.Types[e].Value
	if v == nil || v.Kind() != constant.String ***REMOVED***
		return "", false
	***REMOVED***
	s = constant.StringVal(v)
	// Only record strings with letters.
	for _, r := range s ***REMOVED***
		if unicode.In(r, unicode.L) ***REMOVED***
			return s, true
		***REMOVED***
	***REMOVED***
	return "", false
***REMOVED***
