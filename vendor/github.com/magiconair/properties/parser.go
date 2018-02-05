// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"fmt"
	"runtime"
)

type parser struct ***REMOVED***
	lex *lexer
***REMOVED***

func parse(input string) (properties *Properties, err error) ***REMOVED***
	p := &parser***REMOVED***lex: lex(input)***REMOVED***
	defer p.recover(&err)

	properties = NewProperties()
	key := ""
	comments := []string***REMOVED******REMOVED***

	for ***REMOVED***
		token := p.expectOneOf(itemComment, itemKey, itemEOF)
		switch token.typ ***REMOVED***
		case itemEOF:
			goto done
		case itemComment:
			comments = append(comments, token.val)
			continue
		case itemKey:
			key = token.val
			if _, ok := properties.m[key]; !ok ***REMOVED***
				properties.k = append(properties.k, key)
			***REMOVED***
		***REMOVED***

		token = p.expectOneOf(itemValue, itemEOF)
		if len(comments) > 0 ***REMOVED***
			properties.c[key] = comments
			comments = []string***REMOVED******REMOVED***
		***REMOVED***
		switch token.typ ***REMOVED***
		case itemEOF:
			properties.m[key] = ""
			goto done
		case itemValue:
			properties.m[key] = token.val
		***REMOVED***
	***REMOVED***

done:
	return properties, nil
***REMOVED***

func (p *parser) errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	format = fmt.Sprintf("properties: Line %d: %s", p.lex.lineNumber(), format)
	panic(fmt.Errorf(format, args...))
***REMOVED***

func (p *parser) expect(expected itemType) (token item) ***REMOVED***
	token = p.lex.nextItem()
	if token.typ != expected ***REMOVED***
		p.unexpected(token)
	***REMOVED***
	return token
***REMOVED***

func (p *parser) expectOneOf(expected ...itemType) (token item) ***REMOVED***
	token = p.lex.nextItem()
	for _, v := range expected ***REMOVED***
		if token.typ == v ***REMOVED***
			return token
		***REMOVED***
	***REMOVED***
	p.unexpected(token)
	panic("unexpected token")
***REMOVED***

func (p *parser) unexpected(token item) ***REMOVED***
	p.errorf(token.String())
***REMOVED***

// recover is the handler that turns panics into returns from the top level of Parse.
func (p *parser) recover(errp *error) ***REMOVED***
	e := recover()
	if e != nil ***REMOVED***
		if _, ok := e.(runtime.Error); ok ***REMOVED***
			panic(e)
		***REMOVED***
		*errp = e.(error)
	***REMOVED***
	return
***REMOVED***
