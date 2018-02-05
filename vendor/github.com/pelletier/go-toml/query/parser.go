/*
  Based on the "jsonpath" spec/concept.

  http://goessner.net/articles/JsonPath/
  https://code.google.com/p/json-path/
*/

package query

import (
	"fmt"
)

const maxInt = int(^uint(0) >> 1)

type queryParser struct ***REMOVED***
	flow         chan token
	tokensBuffer []token
	query        *Query
	union        []pathFn
	err          error
***REMOVED***

type queryParserStateFn func() queryParserStateFn

// Formats and panics an error message based on a token
func (p *queryParser) parseError(tok *token, msg string, args ...interface***REMOVED******REMOVED***) queryParserStateFn ***REMOVED***
	p.err = fmt.Errorf(tok.Position.String()+": "+msg, args...)
	return nil // trigger parse to end
***REMOVED***

func (p *queryParser) run() ***REMOVED***
	for state := p.parseStart; state != nil; ***REMOVED***
		state = state()
	***REMOVED***
***REMOVED***

func (p *queryParser) backup(tok *token) ***REMOVED***
	p.tokensBuffer = append(p.tokensBuffer, *tok)
***REMOVED***

func (p *queryParser) peek() *token ***REMOVED***
	if len(p.tokensBuffer) != 0 ***REMOVED***
		return &(p.tokensBuffer[0])
	***REMOVED***

	tok, ok := <-p.flow
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	p.backup(&tok)
	return &tok
***REMOVED***

func (p *queryParser) lookahead(types ...tokenType) bool ***REMOVED***
	result := true
	buffer := []token***REMOVED******REMOVED***

	for _, typ := range types ***REMOVED***
		tok := p.getToken()
		if tok == nil ***REMOVED***
			result = false
			break
		***REMOVED***
		buffer = append(buffer, *tok)
		if tok.typ != typ ***REMOVED***
			result = false
			break
		***REMOVED***
	***REMOVED***
	// add the tokens back to the buffer, and return
	p.tokensBuffer = append(p.tokensBuffer, buffer...)
	return result
***REMOVED***

func (p *queryParser) getToken() *token ***REMOVED***
	if len(p.tokensBuffer) != 0 ***REMOVED***
		tok := p.tokensBuffer[0]
		p.tokensBuffer = p.tokensBuffer[1:]
		return &tok
	***REMOVED***
	tok, ok := <-p.flow
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return &tok
***REMOVED***

func (p *queryParser) parseStart() queryParserStateFn ***REMOVED***
	tok := p.getToken()

	if tok == nil || tok.typ == tokenEOF ***REMOVED***
		return nil
	***REMOVED***

	if tok.typ != tokenDollar ***REMOVED***
		return p.parseError(tok, "Expected '$' at start of expression")
	***REMOVED***

	return p.parseMatchExpr
***REMOVED***

// handle '.' prefix, '[]', and '..'
func (p *queryParser) parseMatchExpr() queryParserStateFn ***REMOVED***
	tok := p.getToken()
	switch tok.typ ***REMOVED***
	case tokenDotDot:
		p.query.appendPath(&matchRecursiveFn***REMOVED******REMOVED***)
		// nested parse for '..'
		tok := p.getToken()
		switch tok.typ ***REMOVED***
		case tokenKey:
			p.query.appendPath(newMatchKeyFn(tok.val))
			return p.parseMatchExpr
		case tokenLeftBracket:
			return p.parseBracketExpr
		case tokenStar:
			// do nothing - the recursive predicate is enough
			return p.parseMatchExpr
		***REMOVED***

	case tokenDot:
		// nested parse for '.'
		tok := p.getToken()
		switch tok.typ ***REMOVED***
		case tokenKey:
			p.query.appendPath(newMatchKeyFn(tok.val))
			return p.parseMatchExpr
		case tokenStar:
			p.query.appendPath(&matchAnyFn***REMOVED******REMOVED***)
			return p.parseMatchExpr
		***REMOVED***

	case tokenLeftBracket:
		return p.parseBracketExpr

	case tokenEOF:
		return nil // allow EOF at this stage
	***REMOVED***
	return p.parseError(tok, "expected match expression")
***REMOVED***

func (p *queryParser) parseBracketExpr() queryParserStateFn ***REMOVED***
	if p.lookahead(tokenInteger, tokenColon) ***REMOVED***
		return p.parseSliceExpr
	***REMOVED***
	if p.peek().typ == tokenColon ***REMOVED***
		return p.parseSliceExpr
	***REMOVED***
	return p.parseUnionExpr
***REMOVED***

func (p *queryParser) parseUnionExpr() queryParserStateFn ***REMOVED***
	var tok *token

	// this state can be traversed after some sub-expressions
	// so be careful when setting up state in the parser
	if p.union == nil ***REMOVED***
		p.union = []pathFn***REMOVED******REMOVED***
	***REMOVED***

loop: // labeled loop for easy breaking
	for ***REMOVED***
		if len(p.union) > 0 ***REMOVED***
			// parse delimiter or terminator
			tok = p.getToken()
			switch tok.typ ***REMOVED***
			case tokenComma:
				// do nothing
			case tokenRightBracket:
				break loop
			default:
				return p.parseError(tok, "expected ',' or ']', not '%s'", tok.val)
			***REMOVED***
		***REMOVED***

		// parse sub expression
		tok = p.getToken()
		switch tok.typ ***REMOVED***
		case tokenInteger:
			p.union = append(p.union, newMatchIndexFn(tok.Int()))
		case tokenKey:
			p.union = append(p.union, newMatchKeyFn(tok.val))
		case tokenString:
			p.union = append(p.union, newMatchKeyFn(tok.val))
		case tokenQuestion:
			return p.parseFilterExpr
		default:
			return p.parseError(tok, "expected union sub expression, not '%s', %d", tok.val, len(p.union))
		***REMOVED***
	***REMOVED***

	// if there is only one sub-expression, use that instead
	if len(p.union) == 1 ***REMOVED***
		p.query.appendPath(p.union[0])
	***REMOVED*** else ***REMOVED***
		p.query.appendPath(&matchUnionFn***REMOVED***p.union***REMOVED***)
	***REMOVED***

	p.union = nil // clear out state
	return p.parseMatchExpr
***REMOVED***

func (p *queryParser) parseSliceExpr() queryParserStateFn ***REMOVED***
	// init slice to grab all elements
	start, end, step := 0, maxInt, 1

	// parse optional start
	tok := p.getToken()
	if tok.typ == tokenInteger ***REMOVED***
		start = tok.Int()
		tok = p.getToken()
	***REMOVED***
	if tok.typ != tokenColon ***REMOVED***
		return p.parseError(tok, "expected ':'")
	***REMOVED***

	// parse optional end
	tok = p.getToken()
	if tok.typ == tokenInteger ***REMOVED***
		end = tok.Int()
		tok = p.getToken()
	***REMOVED***
	if tok.typ == tokenRightBracket ***REMOVED***
		p.query.appendPath(newMatchSliceFn(start, end, step))
		return p.parseMatchExpr
	***REMOVED***
	if tok.typ != tokenColon ***REMOVED***
		return p.parseError(tok, "expected ']' or ':'")
	***REMOVED***

	// parse optional step
	tok = p.getToken()
	if tok.typ == tokenInteger ***REMOVED***
		step = tok.Int()
		if step < 0 ***REMOVED***
			return p.parseError(tok, "step must be a positive value")
		***REMOVED***
		tok = p.getToken()
	***REMOVED***
	if tok.typ != tokenRightBracket ***REMOVED***
		return p.parseError(tok, "expected ']'")
	***REMOVED***

	p.query.appendPath(newMatchSliceFn(start, end, step))
	return p.parseMatchExpr
***REMOVED***

func (p *queryParser) parseFilterExpr() queryParserStateFn ***REMOVED***
	tok := p.getToken()
	if tok.typ != tokenLeftParen ***REMOVED***
		return p.parseError(tok, "expected left-parenthesis for filter expression")
	***REMOVED***
	tok = p.getToken()
	if tok.typ != tokenKey && tok.typ != tokenString ***REMOVED***
		return p.parseError(tok, "expected key or string for filter function name")
	***REMOVED***
	name := tok.val
	tok = p.getToken()
	if tok.typ != tokenRightParen ***REMOVED***
		return p.parseError(tok, "expected right-parenthesis for filter expression")
	***REMOVED***
	p.union = append(p.union, newMatchFilterFn(name, tok.Position))
	return p.parseUnionExpr
***REMOVED***

func parseQuery(flow chan token) (*Query, error) ***REMOVED***
	parser := &queryParser***REMOVED***
		flow:         flow,
		tokensBuffer: []token***REMOVED******REMOVED***,
		query:        newQuery(),
	***REMOVED***
	parser.run()
	return parser.query, parser.err
***REMOVED***
