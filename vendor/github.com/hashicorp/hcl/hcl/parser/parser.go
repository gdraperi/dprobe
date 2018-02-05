// Package parser implements a parser for HCL (HashiCorp Configuration
// Language)
package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/scanner"
	"github.com/hashicorp/hcl/hcl/token"
)

type Parser struct ***REMOVED***
	sc *scanner.Scanner

	// Last read token
	tok       token.Token
	commaPrev token.Token

	comments    []*ast.CommentGroup
	leadComment *ast.CommentGroup // last lead comment
	lineComment *ast.CommentGroup // last line comment

	enableTrace bool
	indent      int
	n           int // buffer size (max = 1)
***REMOVED***

func newParser(src []byte) *Parser ***REMOVED***
	return &Parser***REMOVED***
		sc: scanner.New(src),
	***REMOVED***
***REMOVED***

// Parse returns the fully parsed source and returns the abstract syntax tree.
func Parse(src []byte) (*ast.File, error) ***REMOVED***
	// normalize all line endings
	// since the scanner and output only work with "\n" line endings, we may
	// end up with dangling "\r" characters in the parsed data.
	src = bytes.Replace(src, []byte("\r\n"), []byte("\n"), -1)

	p := newParser(src)
	return p.Parse()
***REMOVED***

var errEofToken = errors.New("EOF token found")

// Parse returns the fully parsed source and returns the abstract syntax tree.
func (p *Parser) Parse() (*ast.File, error) ***REMOVED***
	f := &ast.File***REMOVED******REMOVED***
	var err, scerr error
	p.sc.Error = func(pos token.Pos, msg string) ***REMOVED***
		scerr = &PosError***REMOVED***Pos: pos, Err: errors.New(msg)***REMOVED***
	***REMOVED***

	f.Node, err = p.objectList(false)
	if scerr != nil ***REMOVED***
		return nil, scerr
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	f.Comments = p.comments
	return f, nil
***REMOVED***

// objectList parses a list of items within an object (generally k/v pairs).
// The parameter" obj" tells this whether to we are within an object (braces:
// '***REMOVED***', '***REMOVED***') or just at the top level. If we're within an object, we end
// at an RBRACE.
func (p *Parser) objectList(obj bool) (*ast.ObjectList, error) ***REMOVED***
	defer un(trace(p, "ParseObjectList"))
	node := &ast.ObjectList***REMOVED******REMOVED***

	for ***REMOVED***
		if obj ***REMOVED***
			tok := p.scan()
			p.unscan()
			if tok.Type == token.RBRACE ***REMOVED***
				break
			***REMOVED***
		***REMOVED***

		n, err := p.objectItem()
		if err == errEofToken ***REMOVED***
			break // we are finished
		***REMOVED***

		// we don't return a nil node, because might want to use already
		// collected items.
		if err != nil ***REMOVED***
			return node, err
		***REMOVED***

		node.Add(n)

		// object lists can be optionally comma-delimited e.g. when a list of maps
		// is being expressed, so a comma is allowed here - it's simply consumed
		tok := p.scan()
		if tok.Type != token.COMMA ***REMOVED***
			p.unscan()
		***REMOVED***
	***REMOVED***
	return node, nil
***REMOVED***

func (p *Parser) consumeComment() (comment *ast.Comment, endline int) ***REMOVED***
	endline = p.tok.Pos.Line

	// count the endline if it's multiline comment, ie starting with /*
	if len(p.tok.Text) > 1 && p.tok.Text[1] == '*' ***REMOVED***
		// don't use range here - no need to decode Unicode code points
		for i := 0; i < len(p.tok.Text); i++ ***REMOVED***
			if p.tok.Text[i] == '\n' ***REMOVED***
				endline++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	comment = &ast.Comment***REMOVED***Start: p.tok.Pos, Text: p.tok.Text***REMOVED***
	p.tok = p.sc.Scan()
	return
***REMOVED***

func (p *Parser) consumeCommentGroup(n int) (comments *ast.CommentGroup, endline int) ***REMOVED***
	var list []*ast.Comment
	endline = p.tok.Pos.Line

	for p.tok.Type == token.COMMENT && p.tok.Pos.Line <= endline+n ***REMOVED***
		var comment *ast.Comment
		comment, endline = p.consumeComment()
		list = append(list, comment)
	***REMOVED***

	// add comment group to the comments list
	comments = &ast.CommentGroup***REMOVED***List: list***REMOVED***
	p.comments = append(p.comments, comments)

	return
***REMOVED***

// objectItem parses a single object item
func (p *Parser) objectItem() (*ast.ObjectItem, error) ***REMOVED***
	defer un(trace(p, "ParseObjectItem"))

	keys, err := p.objectKey()
	if len(keys) > 0 && err == errEofToken ***REMOVED***
		// We ignore eof token here since it is an error if we didn't
		// receive a value (but we did receive a key) for the item.
		err = nil
	***REMOVED***
	if len(keys) > 0 && err != nil && p.tok.Type == token.RBRACE ***REMOVED***
		// This is a strange boolean statement, but what it means is:
		// We have keys with no value, and we're likely in an object
		// (since RBrace ends an object). For this, we set err to nil so
		// we continue and get the error below of having the wrong value
		// type.
		err = nil

		// Reset the token type so we don't think it completed fine. See
		// objectType which uses p.tok.Type to check if we're done with
		// the object.
		p.tok.Type = token.EOF
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	o := &ast.ObjectItem***REMOVED***
		Keys: keys,
	***REMOVED***

	if p.leadComment != nil ***REMOVED***
		o.LeadComment = p.leadComment
		p.leadComment = nil
	***REMOVED***

	switch p.tok.Type ***REMOVED***
	case token.ASSIGN:
		o.Assign = p.tok.Pos
		o.Val, err = p.object()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	case token.LBRACE:
		o.Val, err = p.objectType()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	default:
		keyStr := make([]string, 0, len(keys))
		for _, k := range keys ***REMOVED***
			keyStr = append(keyStr, k.Token.Text)
		***REMOVED***

		return nil, &PosError***REMOVED***
			Pos: p.tok.Pos,
			Err: fmt.Errorf(
				"key '%s' expected start of object ('***REMOVED***') or assignment ('=')",
				strings.Join(keyStr, " ")),
		***REMOVED***
	***REMOVED***

	// do a look-ahead for line comment
	p.scan()
	if len(keys) > 0 && o.Val.Pos().Line == keys[0].Pos().Line && p.lineComment != nil ***REMOVED***
		o.LineComment = p.lineComment
		p.lineComment = nil
	***REMOVED***
	p.unscan()
	return o, nil
***REMOVED***

// objectKey parses an object key and returns a ObjectKey AST
func (p *Parser) objectKey() ([]*ast.ObjectKey, error) ***REMOVED***
	keyCount := 0
	keys := make([]*ast.ObjectKey, 0)

	for ***REMOVED***
		tok := p.scan()
		switch tok.Type ***REMOVED***
		case token.EOF:
			// It is very important to also return the keys here as well as
			// the error. This is because we need to be able to tell if we
			// did parse keys prior to finding the EOF, or if we just found
			// a bare EOF.
			return keys, errEofToken
		case token.ASSIGN:
			// assignment or object only, but not nested objects. this is not
			// allowed: `foo bar = ***REMOVED******REMOVED***`
			if keyCount > 1 ***REMOVED***
				return nil, &PosError***REMOVED***
					Pos: p.tok.Pos,
					Err: fmt.Errorf("nested object expected: LBRACE got: %s", p.tok.Type),
				***REMOVED***
			***REMOVED***

			if keyCount == 0 ***REMOVED***
				return nil, &PosError***REMOVED***
					Pos: p.tok.Pos,
					Err: errors.New("no object keys found!"),
				***REMOVED***
			***REMOVED***

			return keys, nil
		case token.LBRACE:
			var err error

			// If we have no keys, then it is a syntax error. i.e. ***REMOVED******REMOVED******REMOVED******REMOVED*** is not
			// allowed.
			if len(keys) == 0 ***REMOVED***
				err = &PosError***REMOVED***
					Pos: p.tok.Pos,
					Err: fmt.Errorf("expected: IDENT | STRING got: %s", p.tok.Type),
				***REMOVED***
			***REMOVED***

			// object
			return keys, err
		case token.IDENT, token.STRING:
			keyCount++
			keys = append(keys, &ast.ObjectKey***REMOVED***Token: p.tok***REMOVED***)
		case token.ILLEGAL:
			return keys, &PosError***REMOVED***
				Pos: p.tok.Pos,
				Err: fmt.Errorf("illegal character"),
			***REMOVED***
		default:
			return keys, &PosError***REMOVED***
				Pos: p.tok.Pos,
				Err: fmt.Errorf("expected: IDENT | STRING | ASSIGN | LBRACE got: %s", p.tok.Type),
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// object parses any type of object, such as number, bool, string, object or
// list.
func (p *Parser) object() (ast.Node, error) ***REMOVED***
	defer un(trace(p, "ParseType"))
	tok := p.scan()

	switch tok.Type ***REMOVED***
	case token.NUMBER, token.FLOAT, token.BOOL, token.STRING, token.HEREDOC:
		return p.literalType()
	case token.LBRACE:
		return p.objectType()
	case token.LBRACK:
		return p.listType()
	case token.COMMENT:
		// implement comment
	case token.EOF:
		return nil, errEofToken
	***REMOVED***

	return nil, &PosError***REMOVED***
		Pos: tok.Pos,
		Err: fmt.Errorf("Unknown token: %+v", tok),
	***REMOVED***
***REMOVED***

// objectType parses an object type and returns a ObjectType AST
func (p *Parser) objectType() (*ast.ObjectType, error) ***REMOVED***
	defer un(trace(p, "ParseObjectType"))

	// we assume that the currently scanned token is a LBRACE
	o := &ast.ObjectType***REMOVED***
		Lbrace: p.tok.Pos,
	***REMOVED***

	l, err := p.objectList(true)

	// if we hit RBRACE, we are good to go (means we parsed all Items), if it's
	// not a RBRACE, it's an syntax error and we just return it.
	if err != nil && p.tok.Type != token.RBRACE ***REMOVED***
		return nil, err
	***REMOVED***

	// No error, scan and expect the ending to be a brace
	if tok := p.scan(); tok.Type != token.RBRACE ***REMOVED***
		return nil, &PosError***REMOVED***
			Pos: tok.Pos,
			Err: fmt.Errorf("object expected closing RBRACE got: %s", tok.Type),
		***REMOVED***
	***REMOVED***

	o.List = l
	o.Rbrace = p.tok.Pos // advanced via parseObjectList
	return o, nil
***REMOVED***

// listType parses a list type and returns a ListType AST
func (p *Parser) listType() (*ast.ListType, error) ***REMOVED***
	defer un(trace(p, "ParseListType"))

	// we assume that the currently scanned token is a LBRACK
	l := &ast.ListType***REMOVED***
		Lbrack: p.tok.Pos,
	***REMOVED***

	needComma := false
	for ***REMOVED***
		tok := p.scan()
		if needComma ***REMOVED***
			switch tok.Type ***REMOVED***
			case token.COMMA, token.RBRACK:
			default:
				return nil, &PosError***REMOVED***
					Pos: tok.Pos,
					Err: fmt.Errorf(
						"error parsing list, expected comma or list end, got: %s",
						tok.Type),
				***REMOVED***
			***REMOVED***
		***REMOVED***
		switch tok.Type ***REMOVED***
		case token.BOOL, token.NUMBER, token.FLOAT, token.STRING, token.HEREDOC:
			node, err := p.literalType()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			// If there is a lead comment, apply it
			if p.leadComment != nil ***REMOVED***
				node.LeadComment = p.leadComment
				p.leadComment = nil
			***REMOVED***

			l.Add(node)
			needComma = true
		case token.COMMA:
			// get next list item or we are at the end
			// do a look-ahead for line comment
			p.scan()
			if p.lineComment != nil && len(l.List) > 0 ***REMOVED***
				lit, ok := l.List[len(l.List)-1].(*ast.LiteralType)
				if ok ***REMOVED***
					lit.LineComment = p.lineComment
					l.List[len(l.List)-1] = lit
					p.lineComment = nil
				***REMOVED***
			***REMOVED***
			p.unscan()

			needComma = false
			continue
		case token.LBRACE:
			// Looks like a nested object, so parse it out
			node, err := p.objectType()
			if err != nil ***REMOVED***
				return nil, &PosError***REMOVED***
					Pos: tok.Pos,
					Err: fmt.Errorf(
						"error while trying to parse object within list: %s", err),
				***REMOVED***
			***REMOVED***
			l.Add(node)
			needComma = true
		case token.LBRACK:
			node, err := p.listType()
			if err != nil ***REMOVED***
				return nil, &PosError***REMOVED***
					Pos: tok.Pos,
					Err: fmt.Errorf(
						"error while trying to parse list within list: %s", err),
				***REMOVED***
			***REMOVED***
			l.Add(node)
		case token.RBRACK:
			// finished
			l.Rbrack = p.tok.Pos
			return l, nil
		default:
			return nil, &PosError***REMOVED***
				Pos: tok.Pos,
				Err: fmt.Errorf("unexpected token while parsing list: %s", tok.Type),
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// literalType parses a literal type and returns a LiteralType AST
func (p *Parser) literalType() (*ast.LiteralType, error) ***REMOVED***
	defer un(trace(p, "ParseLiteral"))

	return &ast.LiteralType***REMOVED***
		Token: p.tok,
	***REMOVED***, nil
***REMOVED***

// scan returns the next token from the underlying scanner. If a token has
// been unscanned then read that instead. In the process, it collects any
// comment groups encountered, and remembers the last lead and line comments.
func (p *Parser) scan() token.Token ***REMOVED***
	// If we have a token on the buffer, then return it.
	if p.n != 0 ***REMOVED***
		p.n = 0
		return p.tok
	***REMOVED***

	// Otherwise read the next token from the scanner and Save it to the buffer
	// in case we unscan later.
	prev := p.tok
	p.tok = p.sc.Scan()

	if p.tok.Type == token.COMMENT ***REMOVED***
		var comment *ast.CommentGroup
		var endline int

		// fmt.Printf("p.tok.Pos.Line = %+v prev: %d endline %d \n",
		// p.tok.Pos.Line, prev.Pos.Line, endline)
		if p.tok.Pos.Line == prev.Pos.Line ***REMOVED***
			// The comment is on same line as the previous token; it
			// cannot be a lead comment but may be a line comment.
			comment, endline = p.consumeCommentGroup(0)
			if p.tok.Pos.Line != endline ***REMOVED***
				// The next token is on a different line, thus
				// the last comment group is a line comment.
				p.lineComment = comment
			***REMOVED***
		***REMOVED***

		// consume successor comments, if any
		endline = -1
		for p.tok.Type == token.COMMENT ***REMOVED***
			comment, endline = p.consumeCommentGroup(1)
		***REMOVED***

		if endline+1 == p.tok.Pos.Line && p.tok.Type != token.RBRACE ***REMOVED***
			switch p.tok.Type ***REMOVED***
			case token.RBRACE, token.RBRACK:
				// Do not count for these cases
			default:
				// The next token is following on the line immediately after the
				// comment group, thus the last comment group is a lead comment.
				p.leadComment = comment
			***REMOVED***
		***REMOVED***

	***REMOVED***

	return p.tok
***REMOVED***

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() ***REMOVED***
	p.n = 1
***REMOVED***

// ----------------------------------------------------------------------------
// Parsing support

func (p *Parser) printTrace(a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if !p.enableTrace ***REMOVED***
		return
	***REMOVED***

	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	fmt.Printf("%5d:%3d: ", p.tok.Pos.Line, p.tok.Pos.Column)

	i := 2 * p.indent
	for i > n ***REMOVED***
		fmt.Print(dots)
		i -= n
	***REMOVED***
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
***REMOVED***

func trace(p *Parser, msg string) *Parser ***REMOVED***
	p.printTrace(msg, "(")
	p.indent++
	return p
***REMOVED***

// Usage pattern: defer un(trace(p, "..."))
func un(p *Parser) ***REMOVED***
	p.indent--
	p.printTrace(")")
***REMOVED***
