package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"log"
)

type cursor_context struct ***REMOVED***
	decl         *decl
	partial      string
	struct_field bool
	decl_import  bool

	// store expression that was supposed to be deduced to "decl", however
	// if decl is nil, then deduction failed, we could try to resolve it to
	// unimported package instead
	expr ast.Expr
***REMOVED***

type token_iterator struct ***REMOVED***
	tokens      []token_item
	token_index int
***REMOVED***

type token_item struct ***REMOVED***
	off int
	tok token.Token
	lit string
***REMOVED***

func (i token_item) literal() string ***REMOVED***
	if i.tok.IsLiteral() ***REMOVED***
		return i.lit
	***REMOVED***
	return i.tok.String()
***REMOVED***

func new_token_iterator(src []byte, cursor int) token_iterator ***REMOVED***
	tokens := make([]token_item, 0, 1000)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, 0)
	for ***REMOVED***
		pos, tok, lit := s.Scan()
		off := fset.Position(pos).Offset
		if tok == token.EOF || cursor <= off ***REMOVED***
			break
		***REMOVED***
		tokens = append(tokens, token_item***REMOVED***
			off: off,
			tok: tok,
			lit: lit,
		***REMOVED***)
	***REMOVED***
	return token_iterator***REMOVED***
		tokens:      tokens,
		token_index: len(tokens) - 1,
	***REMOVED***
***REMOVED***

func (this *token_iterator) token() token_item ***REMOVED***
	return this.tokens[this.token_index]
***REMOVED***

func (this *token_iterator) go_back() bool ***REMOVED***
	if this.token_index <= 0 ***REMOVED***
		return false
	***REMOVED***
	this.token_index--
	return true
***REMOVED***

var bracket_pairs_map = map[token.Token]token.Token***REMOVED***
	token.RPAREN: token.LPAREN,
	token.RBRACK: token.LBRACK,
	token.RBRACE: token.LBRACE,
***REMOVED***

func (ti *token_iterator) skip_to_left(left, right token.Token) bool ***REMOVED***
	if ti.token().tok == left ***REMOVED***
		return true
	***REMOVED***
	balance := 1
	for balance != 0 ***REMOVED***
		if !ti.go_back() ***REMOVED***
			return false
		***REMOVED***
		switch ti.token().tok ***REMOVED***
		case right:
			balance++
		case left:
			balance--
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// when the cursor is at the ')' or ']' or '***REMOVED***', move the cursor to an opposite
// bracket pair, this functions takes nested bracket pairs into account
func (this *token_iterator) skip_to_balanced_pair() bool ***REMOVED***
	right := this.token().tok
	left := bracket_pairs_map[right]
	return this.skip_to_left(left, right)
***REMOVED***

// Move the cursor to the open brace of the current block, taking nested blocks
// into account.
func (this *token_iterator) skip_to_left_curly() bool ***REMOVED***
	return this.skip_to_left(token.LBRACE, token.RBRACE)
***REMOVED***

func (ti *token_iterator) extract_type_alike() string ***REMOVED***
	if ti.token().tok != token.IDENT ***REMOVED*** // not Foo, return nothing
		return ""
	***REMOVED***
	b := ti.token().literal()
	if !ti.go_back() ***REMOVED*** // just Foo
		return b
	***REMOVED***
	if ti.token().tok != token.PERIOD ***REMOVED*** // not .Foo, return Foo
		return b
	***REMOVED***
	if !ti.go_back() ***REMOVED*** // just .Foo, return Foo (best choice recovery)
		return b
	***REMOVED***
	if ti.token().tok != token.IDENT ***REMOVED*** // not lib.Foo, return Foo
		return b
	***REMOVED***
	out := ti.token().literal() + "." + b // lib.Foo
	ti.go_back()
	return out
***REMOVED***

// Extract the type expression right before the enclosing curly bracket block.
// Examples (# - the cursor):
//   &lib.Struct***REMOVED***Whatever: 1, Hel#***REMOVED*** // returns "lib.Struct"
//   X***REMOVED***#***REMOVED***                           // returns X
// The idea is that we check if this type expression is a type and it is, we
// can apply special filtering for autocompletion results.
// Sadly, this doesn't cover anonymous structs.
func (ti *token_iterator) extract_struct_type() string ***REMOVED***
	if !ti.skip_to_left_curly() ***REMOVED***
		return ""
	***REMOVED***
	if !ti.go_back() ***REMOVED***
		return ""
	***REMOVED***
	if ti.token().tok == token.LBRACE ***REMOVED*** // Foo***REMOVED***#***REMOVED******REMOVED******REMOVED***
		if !ti.go_back() ***REMOVED***
			return ""
		***REMOVED***
	***REMOVED*** else if ti.token().tok == token.COMMA ***REMOVED*** // Foo***REMOVED***abc,#***REMOVED******REMOVED******REMOVED***
		return ti.extract_struct_type()
	***REMOVED***
	typ := ti.extract_type_alike()
	if typ == "" ***REMOVED***
		return ""
	***REMOVED***
	if ti.token().tok == token.RPAREN || ti.token().tok == token.MUL ***REMOVED***
		return ""
	***REMOVED***
	return typ
***REMOVED***

// Starting from the token under the cursor move back and extract something
// that resembles a valid Go primary expression. Examples of primary expressions
// from Go spec:
//   x
//   2
//   (s + ".txt")
//   f(3.1415, true)
//   Point***REMOVED***1, 2***REMOVED***
//   m["foo"]
//   s[i : j + 1]
//   obj.color
//   f.p[i].x()
//
// As you can see we can move through all of them using balanced bracket
// matching and applying simple rules
// E.g.
//   Point***REMOVED***1, 2***REMOVED***.m["foo"].s[i : j + 1].MethodCall(a, func(a, b int) int ***REMOVED*** return a + b ***REMOVED***).
// Can be seen as:
//   Point***REMOVED******REMOVED***.m[     ].s[         ].MethodCall(                                      ).
// Which boils the rules down to these connected via dots:
//   ident
//   ident[]
//   ident***REMOVED******REMOVED***
//   ident()
// Of course there are also slightly more complicated rules for brackets:
//   ident***REMOVED******REMOVED***.ident()[5][4](), etc.
func (this *token_iterator) extract_go_expr() string ***REMOVED***
	orig := this.token_index

	// Contains the type of the previously scanned token (initialized with
	// the token right under the cursor). This is the token to the *right* of
	// the current one.
	prev := this.token().tok
loop:
	for ***REMOVED***
		if !this.go_back() ***REMOVED***
			return token_items_to_string(this.tokens[:orig])
		***REMOVED***
		switch this.token().tok ***REMOVED***
		case token.PERIOD:
			// If the '.' is not followed by IDENT, it's invalid.
			if prev != token.IDENT ***REMOVED***
				break loop
			***REMOVED***
		case token.IDENT:
			// Valid tokens after IDENT are '.', '[', '***REMOVED***' and '('.
			switch prev ***REMOVED***
			case token.PERIOD, token.LBRACK, token.LBRACE, token.LPAREN:
				// all ok
			default:
				break loop
			***REMOVED***
		case token.RBRACE:
			// This one can only be a part of type initialization, like:
			//   Dummy***REMOVED******REMOVED***.Hello()
			// It is valid Go if Hello method is defined on a non-pointer receiver.
			if prev != token.PERIOD ***REMOVED***
				break loop
			***REMOVED***
			this.skip_to_balanced_pair()
		case token.RPAREN, token.RBRACK:
			// After ']' and ')' their opening counterparts are valid '[', '(',
			// as well as the dot.
			switch prev ***REMOVED***
			case token.PERIOD, token.LBRACK, token.LPAREN:
				// all ok
			default:
				break loop
			***REMOVED***
			this.skip_to_balanced_pair()
		default:
			break loop
		***REMOVED***
		prev = this.token().tok
	***REMOVED***
	expr := token_items_to_string(this.tokens[this.token_index+1 : orig])
	if *g_debug ***REMOVED***
		log.Printf("extracted expression tokens: %s", expr)
	***REMOVED***
	return expr
***REMOVED***

// Given a slice of token_item, reassembles them into the original literal
// expression.
func token_items_to_string(tokens []token_item) string ***REMOVED***
	var buf bytes.Buffer
	for _, t := range tokens ***REMOVED***
		buf.WriteString(t.literal())
	***REMOVED***
	return buf.String()
***REMOVED***

// this function is called when the cursor is at the '.' and you need to get the
// declaration before that dot
func (c *auto_complete_context) deduce_cursor_decl(iter *token_iterator) (*decl, ast.Expr) ***REMOVED***
	expr, err := parser.ParseExpr(iter.extract_go_expr())
	if err != nil ***REMOVED***
		return nil, nil
	***REMOVED***
	return expr_to_decl(expr, c.current.scope), expr
***REMOVED***

// try to find and extract the surrounding struct literal type
func (c *auto_complete_context) deduce_struct_type_decl(iter *token_iterator) *decl ***REMOVED***
	typ := iter.extract_struct_type()
	if typ == "" ***REMOVED***
		return nil
	***REMOVED***

	expr, err := parser.ParseExpr(typ)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	decl := type_to_decl(expr, c.current.scope)
	if decl == nil ***REMOVED***
		return nil
	***REMOVED***

	// we allow only struct types here, but also support type aliases
	if decl.is_alias() ***REMOVED***
		dd := decl.type_dealias()
		if _, ok := dd.typ.(*ast.StructType); !ok ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED*** else if _, ok := decl.typ.(*ast.StructType); !ok ***REMOVED***
		return nil
	***REMOVED***
	return decl
***REMOVED***

// Entry point from autocompletion, the function looks at text before the cursor
// and figures out the declaration the cursor is on. This declaration is
// used in filtering the resulting set of autocompletion suggestions.
func (c *auto_complete_context) deduce_cursor_context(file []byte, cursor int) (cursor_context, bool) ***REMOVED***
	if cursor <= 0 ***REMOVED***
		return cursor_context***REMOVED******REMOVED***, true
	***REMOVED***

	iter := new_token_iterator(file, cursor)
	if len(iter.tokens) == 0 ***REMOVED***
		return cursor_context***REMOVED******REMOVED***, false
	***REMOVED***

	// figure out what is just before the cursor
	switch tok := iter.token(); tok.tok ***REMOVED***
	case token.STRING:
		// make sure cursor is inside the string
		s := tok.literal()
		if len(s) > 1 && s[len(s)-1] == '"' && tok.off+len(s) <= cursor ***REMOVED***
			return cursor_context***REMOVED******REMOVED***, true
		***REMOVED***
		// now figure out if inside an import declaration
		var ptok = token.STRING
		for iter.go_back() ***REMOVED***
			itok := iter.token().tok
			switch itok ***REMOVED***
			case token.STRING:
				switch ptok ***REMOVED***
				case token.SEMICOLON, token.IDENT, token.PERIOD:
				default:
					return cursor_context***REMOVED******REMOVED***, true
				***REMOVED***
			case token.LPAREN, token.SEMICOLON:
				switch ptok ***REMOVED***
				case token.STRING, token.IDENT, token.PERIOD:
				default:
					return cursor_context***REMOVED******REMOVED***, true
				***REMOVED***
			case token.IDENT, token.PERIOD:
				switch ptok ***REMOVED***
				case token.STRING:
				default:
					return cursor_context***REMOVED******REMOVED***, true
				***REMOVED***
			case token.IMPORT:
				switch ptok ***REMOVED***
				case token.STRING, token.IDENT, token.PERIOD, token.LPAREN:
					path_len := cursor - tok.off
					path := s[1:path_len]
					return cursor_context***REMOVED***decl_import: true, partial: path***REMOVED***, true
				default:
					return cursor_context***REMOVED******REMOVED***, true
				***REMOVED***
			default:
				return cursor_context***REMOVED******REMOVED***, true
			***REMOVED***
			ptok = itok
		***REMOVED***
	case token.PERIOD:
		// we're '<whatever>.'
		// figure out decl, Partial is ""
		decl, expr := c.deduce_cursor_decl(&iter)
		return cursor_context***REMOVED***decl: decl, expr: expr***REMOVED***, decl != nil
	case token.IDENT, token.TYPE, token.CONST, token.VAR, token.FUNC, token.PACKAGE:
		// we're '<whatever>.<ident>'
		// parse <ident> as Partial and figure out decl
		var partial string
		if tok.tok == token.IDENT ***REMOVED***
			// Calculate the offset of the cursor position within the identifier.
			// For instance, if we are 'ab#c', we want partial_len = 2 and partial = ab.
			partial_len := cursor - tok.off

			// If it happens that the cursor is past the end of the literal,
			// means there is a space between the literal and the cursor, think
			// of it as no context, because that's what it really is.
			if partial_len > len(tok.literal()) ***REMOVED***
				return cursor_context***REMOVED******REMOVED***, true
			***REMOVED***
			partial = tok.literal()[0:partial_len]
		***REMOVED*** else ***REMOVED***
			// Do not try to truncate if it is not an identifier.
			partial = tok.literal()
		***REMOVED***

		iter.go_back()
		switch iter.token().tok ***REMOVED***
		case token.PERIOD:
			decl, expr := c.deduce_cursor_decl(&iter)
			return cursor_context***REMOVED***decl: decl, partial: partial, expr: expr***REMOVED***, decl != nil
		case token.COMMA, token.LBRACE:
			// This can happen for struct fields:
			// &Struct***REMOVED***Hello: 1, Wor#***REMOVED*** // (# - the cursor)
			// Let's try to find the struct type
			decl := c.deduce_struct_type_decl(&iter)
			return cursor_context***REMOVED***
				decl:         decl,
				partial:      partial,
				struct_field: decl != nil,
			***REMOVED***, true
		default:
			return cursor_context***REMOVED***partial: partial***REMOVED***, true
		***REMOVED***
	case token.COMMA, token.LBRACE:
		// Try to parse the current expression as a structure initialization.
		decl := c.deduce_struct_type_decl(&iter)
		return cursor_context***REMOVED***
			decl:         decl,
			partial:      "",
			struct_field: decl != nil,
		***REMOVED***, true
	***REMOVED***

	return cursor_context***REMOVED******REMOVED***, true
***REMOVED***

// Decl deduction failed, but we're on "<ident>.", this ident can be an
// unexported package, let's try to match the ident against a set of known
// packages and if it matches try to import it.
// TODO: Right now I've made a static list of built-in packages, but in theory
// we could scan all GOPATH packages as well. Now, don't forget that default
// package name has nothing to do with package file name, that's why we need to
// scan the packages. And many of them will have conflicts. Can we make a smart
// prediction algorithm which will prefer certain packages over another ones?
func resolveKnownPackageIdent(ident string, filename string, context *package_lookup_context) *package_file_cache ***REMOVED***
	importPath, ok := knownPackageIdents[ident]
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	path, ok := abs_path_for_package(filename, importPath, context)
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	p := new_package_file_cache(path, importPath)
	p.update_cache()
	return p
***REMOVED***

var knownPackageIdents = map[string]string***REMOVED***
	"adler32":         "hash/adler32",
	"aes":             "crypto/aes",
	"ascii85":         "encoding/ascii85",
	"asn1":            "encoding/asn1",
	"ast":             "go/ast",
	"atomic":          "sync/atomic",
	"base32":          "encoding/base32",
	"base64":          "encoding/base64",
	"big":             "math/big",
	"binary":          "encoding/binary",
	"bufio":           "bufio",
	"build":           "go/build",
	"bytes":           "bytes",
	"bzip2":           "compress/bzip2",
	"cgi":             "net/http/cgi",
	"cgo":             "runtime/cgo",
	"cipher":          "crypto/cipher",
	"cmplx":           "math/cmplx",
	"color":           "image/color",
	"constant":        "go/constant",
	"context":         "context",
	"cookiejar":       "net/http/cookiejar",
	"crc32":           "hash/crc32",
	"crc64":           "hash/crc64",
	"crypto":          "crypto",
	"csv":             "encoding/csv",
	"debug":           "runtime/debug",
	"des":             "crypto/des",
	"doc":             "go/doc",
	"draw":            "image/draw",
	"driver":          "database/sql/driver",
	"dsa":             "crypto/dsa",
	"dwarf":           "debug/dwarf",
	"ecdsa":           "crypto/ecdsa",
	"elf":             "debug/elf",
	"elliptic":        "crypto/elliptic",
	"encoding":        "encoding",
	"errors":          "errors",
	"exec":            "os/exec",
	"expvar":          "expvar",
	"fcgi":            "net/http/fcgi",
	"filepath":        "path/filepath",
	"flag":            "flag",
	"flate":           "compress/flate",
	"fmt":             "fmt",
	"fnv":             "hash/fnv",
	"format":          "go/format",
	"gif":             "image/gif",
	"gob":             "encoding/gob",
	"gosym":           "debug/gosym",
	"gzip":            "compress/gzip",
	"hash":            "hash",
	"heap":            "container/heap",
	"hex":             "encoding/hex",
	"hmac":            "crypto/hmac",
	"hpack":           "vendor/golang_org/x/net/http2/hpack",
	"html":            "html",
	"http":            "net/http",
	"httplex":         "vendor/golang_org/x/net/lex/httplex",
	"httptest":        "net/http/httptest",
	"httptrace":       "net/http/httptrace",
	"httputil":        "net/http/httputil",
	"image":           "image",
	"importer":        "go/importer",
	"io":              "io",
	"iotest":          "testing/iotest",
	"ioutil":          "io/ioutil",
	"jpeg":            "image/jpeg",
	"json":            "encoding/json",
	"jsonrpc":         "net/rpc/jsonrpc",
	"list":            "container/list",
	"log":             "log",
	"lzw":             "compress/lzw",
	"macho":           "debug/macho",
	"mail":            "net/mail",
	"math":            "math",
	"md5":             "crypto/md5",
	"mime":            "mime",
	"multipart":       "mime/multipart",
	"net":             "net",
	"os":              "os",
	"palette":         "image/color/palette",
	"parse":           "text/template/parse",
	"parser":          "go/parser",
	"path":            "path",
	"pe":              "debug/pe",
	"pem":             "encoding/pem",
	"pkix":            "crypto/x509/pkix",
	"plan9obj":        "debug/plan9obj",
	"png":             "image/png",
	"pprof":           "net/http/pprof",
	"printer":         "go/printer",
	"quick":           "testing/quick",
	"quotedprintable": "mime/quotedprintable",
	"race":            "runtime/race",
	"rand":            "math/rand",
	"rc4":             "crypto/rc4",
	"reflect":         "reflect",
	"regexp":          "regexp",
	"ring":            "container/ring",
	"rpc":             "net/rpc",
	"rsa":             "crypto/rsa",
	"runtime":         "runtime",
	"scanner":         "text/scanner",
	"sha1":            "crypto/sha1",
	"sha256":          "crypto/sha256",
	"sha512":          "crypto/sha512",
	"signal":          "os/signal",
	"smtp":            "net/smtp",
	"sort":            "sort",
	"sql":             "database/sql",
	"strconv":         "strconv",
	"strings":         "strings",
	"subtle":          "crypto/subtle",
	"suffixarray":     "index/suffixarray",
	"sync":            "sync",
	"syntax":          "regexp/syntax",
	"syscall":         "syscall",
	"syslog":          "log/syslog",
	"tabwriter":       "text/tabwriter",
	"tar":             "archive/tar",
	"template":        "html/template",
	"testing":         "testing",
	"textproto":       "net/textproto",
	"time":            "time",
	"tls":             "crypto/tls",
	"token":           "go/token",
	"trace":           "runtime/trace",
	"types":           "go/types",
	"unicode":         "unicode",
	"url":             "net/url",
	"user":            "os/user",
	"utf16":           "unicode/utf16",
	"utf8":            "unicode/utf8",
	"x509":            "crypto/x509",
	"xml":             "encoding/xml",
	"zip":             "archive/zip",
	"zlib":            "compress/zlib",
	//"scanner": "go/scanner", // DUP: prefer text/scanner
	//"template": "text/template", // DUP: prefer html/template
	//"pprof": "runtime/pprof", // DUP: prefer net/http/pprof
	//"rand": "crypto/rand", // DUP: prefer math/rand
***REMOVED***
