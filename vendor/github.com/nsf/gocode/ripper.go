package main

import (
	"go/scanner"
	"go/token"
)

// All the code in this file serves single purpose:
// It separates a function with the cursor inside and the rest of the code. I'm
// doing that, because sometimes parser is not able to recover itself from an
// error and the autocompletion results become less complete.

type tok_pos_pair struct ***REMOVED***
	tok token.Token
	pos token.Pos
***REMOVED***

type tok_collection struct ***REMOVED***
	tokens []tok_pos_pair
	fset   *token.FileSet
***REMOVED***

func (this *tok_collection) next(s *scanner.Scanner) bool ***REMOVED***
	pos, tok, _ := s.Scan()
	if tok == token.EOF ***REMOVED***
		return false
	***REMOVED***

	this.tokens = append(this.tokens, tok_pos_pair***REMOVED***tok, pos***REMOVED***)
	return true
***REMOVED***

func (this *tok_collection) find_decl_beg(pos int) int ***REMOVED***
	lowest := 0
	lowpos := -1
	lowi := -1
	cur := 0
	for i := pos; i >= 0; i-- ***REMOVED***
		t := this.tokens[i]
		switch t.tok ***REMOVED***
		case token.RBRACE:
			cur++
		case token.LBRACE:
			cur--
		***REMOVED***

		if cur < lowest ***REMOVED***
			lowest = cur
			lowpos = this.fset.Position(t.pos).Offset
			lowi = i
		***REMOVED***
	***REMOVED***

	cur = lowest
	for i := lowi - 1; i >= 0; i-- ***REMOVED***
		t := this.tokens[i]
		switch t.tok ***REMOVED***
		case token.RBRACE:
			cur++
		case token.LBRACE:
			cur--
		***REMOVED***
		if t.tok == token.SEMICOLON && cur == lowest ***REMOVED***
			lowpos = this.fset.Position(t.pos).Offset
			break
		***REMOVED***
	***REMOVED***

	return lowpos
***REMOVED***

func (this *tok_collection) find_decl_end(pos int) int ***REMOVED***
	highest := 0
	highpos := -1
	cur := 0

	if this.tokens[pos].tok == token.LBRACE ***REMOVED***
		pos++
	***REMOVED***

	for i := pos; i < len(this.tokens); i++ ***REMOVED***
		t := this.tokens[i]
		switch t.tok ***REMOVED***
		case token.RBRACE:
			cur++
		case token.LBRACE:
			cur--
		***REMOVED***

		if cur > highest ***REMOVED***
			highest = cur
			highpos = this.fset.Position(t.pos).Offset
		***REMOVED***
	***REMOVED***

	return highpos
***REMOVED***

func (this *tok_collection) find_outermost_scope(cursor int) (int, int) ***REMOVED***
	pos := 0

	for i, t := range this.tokens ***REMOVED***
		if cursor <= this.fset.Position(t.pos).Offset ***REMOVED***
			break
		***REMOVED***
		pos = i
	***REMOVED***

	return this.find_decl_beg(pos), this.find_decl_end(pos)
***REMOVED***

// return new cursor position, file without ripped part and the ripped part itself
// variants:
//   new-cursor, file-without-ripped-part, ripped-part
//   old-cursor, file, nil
func (this *tok_collection) rip_off_decl(file []byte, cursor int) (int, []byte, []byte) ***REMOVED***
	this.fset = token.NewFileSet()
	var s scanner.Scanner
	s.Init(this.fset.AddFile("", this.fset.Base(), len(file)), file, nil, scanner.ScanComments)
	for this.next(&s) ***REMOVED***
	***REMOVED***

	beg, end := this.find_outermost_scope(cursor)
	if beg == -1 || end == -1 ***REMOVED***
		return cursor, file, nil
	***REMOVED***

	ripped := make([]byte, end+1-beg)
	copy(ripped, file[beg:end+1])

	newfile := make([]byte, len(file)-len(ripped))
	copy(newfile, file[:beg])
	copy(newfile[beg:], file[end+1:])

	return cursor - beg, newfile, ripped
***REMOVED***

func rip_off_decl(file []byte, cursor int) (int, []byte, []byte) ***REMOVED***
	var tc tok_collection
	return tc.rip_off_decl(file, cursor)
***REMOVED***
