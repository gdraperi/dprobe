package jmespath

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type astNodeType int

//go:generate stringer -type astNodeType
const (
	ASTEmpty astNodeType = iota
	ASTComparator
	ASTCurrentNode
	ASTExpRef
	ASTFunctionExpression
	ASTField
	ASTFilterProjection
	ASTFlatten
	ASTIdentity
	ASTIndex
	ASTIndexExpression
	ASTKeyValPair
	ASTLiteral
	ASTMultiSelectHash
	ASTMultiSelectList
	ASTOrExpression
	ASTAndExpression
	ASTNotExpression
	ASTPipe
	ASTProjection
	ASTSubexpression
	ASTSlice
	ASTValueProjection
)

// ASTNode represents the abstract syntax tree of a JMESPath expression.
type ASTNode struct ***REMOVED***
	nodeType astNodeType
	value    interface***REMOVED******REMOVED***
	children []ASTNode
***REMOVED***

func (node ASTNode) String() string ***REMOVED***
	return node.PrettyPrint(0)
***REMOVED***

// PrettyPrint will pretty print the parsed AST.
// The AST is an implementation detail and this pretty print
// function is provided as a convenience method to help with
// debugging.  You should not rely on its output as the internal
// structure of the AST may change at any time.
func (node ASTNode) PrettyPrint(indent int) string ***REMOVED***
	spaces := strings.Repeat(" ", indent)
	output := fmt.Sprintf("%s%s ***REMOVED***\n", spaces, node.nodeType)
	nextIndent := indent + 2
	if node.value != nil ***REMOVED***
		if converted, ok := node.value.(fmt.Stringer); ok ***REMOVED***
			// Account for things like comparator nodes
			// that are enums with a String() method.
			output += fmt.Sprintf("%svalue: %s\n", strings.Repeat(" ", nextIndent), converted.String())
		***REMOVED*** else ***REMOVED***
			output += fmt.Sprintf("%svalue: %#v\n", strings.Repeat(" ", nextIndent), node.value)
		***REMOVED***
	***REMOVED***
	lastIndex := len(node.children)
	if lastIndex > 0 ***REMOVED***
		output += fmt.Sprintf("%schildren: ***REMOVED***\n", strings.Repeat(" ", nextIndent))
		childIndent := nextIndent + 2
		for _, elem := range node.children ***REMOVED***
			output += elem.PrettyPrint(childIndent)
		***REMOVED***
	***REMOVED***
	output += fmt.Sprintf("%s***REMOVED***\n", spaces)
	return output
***REMOVED***

var bindingPowers = map[tokType]int***REMOVED***
	tEOF:                0,
	tUnquotedIdentifier: 0,
	tQuotedIdentifier:   0,
	tRbracket:           0,
	tRparen:             0,
	tComma:              0,
	tRbrace:             0,
	tNumber:             0,
	tCurrent:            0,
	tExpref:             0,
	tColon:              0,
	tPipe:               1,
	tOr:                 2,
	tAnd:                3,
	tEQ:                 5,
	tLT:                 5,
	tLTE:                5,
	tGT:                 5,
	tGTE:                5,
	tNE:                 5,
	tFlatten:            9,
	tStar:               20,
	tFilter:             21,
	tDot:                40,
	tNot:                45,
	tLbrace:             50,
	tLbracket:           55,
	tLparen:             60,
***REMOVED***

// Parser holds state about the current expression being parsed.
type Parser struct ***REMOVED***
	expression string
	tokens     []token
	index      int
***REMOVED***

// NewParser creates a new JMESPath parser.
func NewParser() *Parser ***REMOVED***
	p := Parser***REMOVED******REMOVED***
	return &p
***REMOVED***

// Parse will compile a JMESPath expression.
func (p *Parser) Parse(expression string) (ASTNode, error) ***REMOVED***
	lexer := NewLexer()
	p.expression = expression
	p.index = 0
	tokens, err := lexer.tokenize(expression)
	if err != nil ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, err
	***REMOVED***
	p.tokens = tokens
	parsed, err := p.parseExpression(0)
	if err != nil ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, err
	***REMOVED***
	if p.current() != tEOF ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, p.syntaxError(fmt.Sprintf(
			"Unexpected token at the end of the expresssion: %s", p.current()))
	***REMOVED***
	return parsed, nil
***REMOVED***

func (p *Parser) parseExpression(bindingPower int) (ASTNode, error) ***REMOVED***
	var err error
	leftToken := p.lookaheadToken(0)
	p.advance()
	leftNode, err := p.nud(leftToken)
	if err != nil ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, err
	***REMOVED***
	currentToken := p.current()
	for bindingPower < bindingPowers[currentToken] ***REMOVED***
		p.advance()
		leftNode, err = p.led(currentToken, leftNode)
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		currentToken = p.current()
	***REMOVED***
	return leftNode, nil
***REMOVED***

func (p *Parser) parseIndexExpression() (ASTNode, error) ***REMOVED***
	if p.lookahead(0) == tColon || p.lookahead(1) == tColon ***REMOVED***
		return p.parseSliceExpression()
	***REMOVED***
	indexStr := p.lookaheadToken(0).value
	parsedInt, err := strconv.Atoi(indexStr)
	if err != nil ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, err
	***REMOVED***
	indexNode := ASTNode***REMOVED***nodeType: ASTIndex, value: parsedInt***REMOVED***
	p.advance()
	if err := p.match(tRbracket); err != nil ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, err
	***REMOVED***
	return indexNode, nil
***REMOVED***

func (p *Parser) parseSliceExpression() (ASTNode, error) ***REMOVED***
	parts := []*int***REMOVED***nil, nil, nil***REMOVED***
	index := 0
	current := p.current()
	for current != tRbracket && index < 3 ***REMOVED***
		if current == tColon ***REMOVED***
			index++
			p.advance()
		***REMOVED*** else if current == tNumber ***REMOVED***
			parsedInt, err := strconv.Atoi(p.lookaheadToken(0).value)
			if err != nil ***REMOVED***
				return ASTNode***REMOVED******REMOVED***, err
			***REMOVED***
			parts[index] = &parsedInt
			p.advance()
		***REMOVED*** else ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, p.syntaxError(
				"Expected tColon or tNumber" + ", received: " + p.current().String())
		***REMOVED***
		current = p.current()
	***REMOVED***
	if err := p.match(tRbracket); err != nil ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, err
	***REMOVED***
	return ASTNode***REMOVED***
		nodeType: ASTSlice,
		value:    parts,
	***REMOVED***, nil
***REMOVED***

func (p *Parser) match(tokenType tokType) error ***REMOVED***
	if p.current() == tokenType ***REMOVED***
		p.advance()
		return nil
	***REMOVED***
	return p.syntaxError("Expected " + tokenType.String() + ", received: " + p.current().String())
***REMOVED***

func (p *Parser) led(tokenType tokType, node ASTNode) (ASTNode, error) ***REMOVED***
	switch tokenType ***REMOVED***
	case tDot:
		if p.current() != tStar ***REMOVED***
			right, err := p.parseDotRHS(bindingPowers[tDot])
			return ASTNode***REMOVED***
				nodeType: ASTSubexpression,
				children: []ASTNode***REMOVED***node, right***REMOVED***,
			***REMOVED***, err
		***REMOVED***
		p.advance()
		right, err := p.parseProjectionRHS(bindingPowers[tDot])
		return ASTNode***REMOVED***
			nodeType: ASTValueProjection,
			children: []ASTNode***REMOVED***node, right***REMOVED***,
		***REMOVED***, err
	case tPipe:
		right, err := p.parseExpression(bindingPowers[tPipe])
		return ASTNode***REMOVED***nodeType: ASTPipe, children: []ASTNode***REMOVED***node, right***REMOVED******REMOVED***, err
	case tOr:
		right, err := p.parseExpression(bindingPowers[tOr])
		return ASTNode***REMOVED***nodeType: ASTOrExpression, children: []ASTNode***REMOVED***node, right***REMOVED******REMOVED***, err
	case tAnd:
		right, err := p.parseExpression(bindingPowers[tAnd])
		return ASTNode***REMOVED***nodeType: ASTAndExpression, children: []ASTNode***REMOVED***node, right***REMOVED******REMOVED***, err
	case tLparen:
		name := node.value
		var args []ASTNode
		for p.current() != tRparen ***REMOVED***
			expression, err := p.parseExpression(0)
			if err != nil ***REMOVED***
				return ASTNode***REMOVED******REMOVED***, err
			***REMOVED***
			if p.current() == tComma ***REMOVED***
				if err := p.match(tComma); err != nil ***REMOVED***
					return ASTNode***REMOVED******REMOVED***, err
				***REMOVED***
			***REMOVED***
			args = append(args, expression)
		***REMOVED***
		if err := p.match(tRparen); err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return ASTNode***REMOVED***
			nodeType: ASTFunctionExpression,
			value:    name,
			children: args,
		***REMOVED***, nil
	case tFilter:
		return p.parseFilter(node)
	case tFlatten:
		left := ASTNode***REMOVED***nodeType: ASTFlatten, children: []ASTNode***REMOVED***node***REMOVED******REMOVED***
		right, err := p.parseProjectionRHS(bindingPowers[tFlatten])
		return ASTNode***REMOVED***
			nodeType: ASTProjection,
			children: []ASTNode***REMOVED***left, right***REMOVED***,
		***REMOVED***, err
	case tEQ, tNE, tGT, tGTE, tLT, tLTE:
		right, err := p.parseExpression(bindingPowers[tokenType])
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return ASTNode***REMOVED***
			nodeType: ASTComparator,
			value:    tokenType,
			children: []ASTNode***REMOVED***node, right***REMOVED***,
		***REMOVED***, nil
	case tLbracket:
		tokenType := p.current()
		var right ASTNode
		var err error
		if tokenType == tNumber || tokenType == tColon ***REMOVED***
			right, err = p.parseIndexExpression()
			if err != nil ***REMOVED***
				return ASTNode***REMOVED******REMOVED***, err
			***REMOVED***
			return p.projectIfSlice(node, right)
		***REMOVED***
		// Otherwise this is a projection.
		if err := p.match(tStar); err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		if err := p.match(tRbracket); err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		right, err = p.parseProjectionRHS(bindingPowers[tStar])
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return ASTNode***REMOVED***
			nodeType: ASTProjection,
			children: []ASTNode***REMOVED***node, right***REMOVED***,
		***REMOVED***, nil
	***REMOVED***
	return ASTNode***REMOVED******REMOVED***, p.syntaxError("Unexpected token: " + tokenType.String())
***REMOVED***

func (p *Parser) nud(token token) (ASTNode, error) ***REMOVED***
	switch token.tokenType ***REMOVED***
	case tJSONLiteral:
		var parsed interface***REMOVED******REMOVED***
		err := json.Unmarshal([]byte(token.value), &parsed)
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return ASTNode***REMOVED***nodeType: ASTLiteral, value: parsed***REMOVED***, nil
	case tStringLiteral:
		return ASTNode***REMOVED***nodeType: ASTLiteral, value: token.value***REMOVED***, nil
	case tUnquotedIdentifier:
		return ASTNode***REMOVED***
			nodeType: ASTField,
			value:    token.value,
		***REMOVED***, nil
	case tQuotedIdentifier:
		node := ASTNode***REMOVED***nodeType: ASTField, value: token.value***REMOVED***
		if p.current() == tLparen ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, p.syntaxErrorToken("Can't have quoted identifier as function name.", token)
		***REMOVED***
		return node, nil
	case tStar:
		left := ASTNode***REMOVED***nodeType: ASTIdentity***REMOVED***
		var right ASTNode
		var err error
		if p.current() == tRbracket ***REMOVED***
			right = ASTNode***REMOVED***nodeType: ASTIdentity***REMOVED***
		***REMOVED*** else ***REMOVED***
			right, err = p.parseProjectionRHS(bindingPowers[tStar])
		***REMOVED***
		return ASTNode***REMOVED***nodeType: ASTValueProjection, children: []ASTNode***REMOVED***left, right***REMOVED******REMOVED***, err
	case tFilter:
		return p.parseFilter(ASTNode***REMOVED***nodeType: ASTIdentity***REMOVED***)
	case tLbrace:
		return p.parseMultiSelectHash()
	case tFlatten:
		left := ASTNode***REMOVED***
			nodeType: ASTFlatten,
			children: []ASTNode***REMOVED******REMOVED***nodeType: ASTIdentity***REMOVED******REMOVED***,
		***REMOVED***
		right, err := p.parseProjectionRHS(bindingPowers[tFlatten])
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return ASTNode***REMOVED***nodeType: ASTProjection, children: []ASTNode***REMOVED***left, right***REMOVED******REMOVED***, nil
	case tLbracket:
		tokenType := p.current()
		//var right ASTNode
		if tokenType == tNumber || tokenType == tColon ***REMOVED***
			right, err := p.parseIndexExpression()
			if err != nil ***REMOVED***
				return ASTNode***REMOVED******REMOVED***, nil
			***REMOVED***
			return p.projectIfSlice(ASTNode***REMOVED***nodeType: ASTIdentity***REMOVED***, right)
		***REMOVED*** else if tokenType == tStar && p.lookahead(1) == tRbracket ***REMOVED***
			p.advance()
			p.advance()
			right, err := p.parseProjectionRHS(bindingPowers[tStar])
			if err != nil ***REMOVED***
				return ASTNode***REMOVED******REMOVED***, err
			***REMOVED***
			return ASTNode***REMOVED***
				nodeType: ASTProjection,
				children: []ASTNode***REMOVED******REMOVED***nodeType: ASTIdentity***REMOVED***, right***REMOVED***,
			***REMOVED***, nil
		***REMOVED*** else ***REMOVED***
			return p.parseMultiSelectList()
		***REMOVED***
	case tCurrent:
		return ASTNode***REMOVED***nodeType: ASTCurrentNode***REMOVED***, nil
	case tExpref:
		expression, err := p.parseExpression(bindingPowers[tExpref])
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return ASTNode***REMOVED***nodeType: ASTExpRef, children: []ASTNode***REMOVED***expression***REMOVED******REMOVED***, nil
	case tNot:
		expression, err := p.parseExpression(bindingPowers[tNot])
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return ASTNode***REMOVED***nodeType: ASTNotExpression, children: []ASTNode***REMOVED***expression***REMOVED******REMOVED***, nil
	case tLparen:
		expression, err := p.parseExpression(0)
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		if err := p.match(tRparen); err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return expression, nil
	case tEOF:
		return ASTNode***REMOVED******REMOVED***, p.syntaxErrorToken("Incomplete expression", token)
	***REMOVED***

	return ASTNode***REMOVED******REMOVED***, p.syntaxErrorToken("Invalid token: "+token.tokenType.String(), token)
***REMOVED***

func (p *Parser) parseMultiSelectList() (ASTNode, error) ***REMOVED***
	var expressions []ASTNode
	for ***REMOVED***
		expression, err := p.parseExpression(0)
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		expressions = append(expressions, expression)
		if p.current() == tRbracket ***REMOVED***
			break
		***REMOVED***
		err = p.match(tComma)
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***
	err := p.match(tRbracket)
	if err != nil ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, err
	***REMOVED***
	return ASTNode***REMOVED***
		nodeType: ASTMultiSelectList,
		children: expressions,
	***REMOVED***, nil
***REMOVED***

func (p *Parser) parseMultiSelectHash() (ASTNode, error) ***REMOVED***
	var children []ASTNode
	for ***REMOVED***
		keyToken := p.lookaheadToken(0)
		if err := p.match(tUnquotedIdentifier); err != nil ***REMOVED***
			if err := p.match(tQuotedIdentifier); err != nil ***REMOVED***
				return ASTNode***REMOVED******REMOVED***, p.syntaxError("Expected tQuotedIdentifier or tUnquotedIdentifier")
			***REMOVED***
		***REMOVED***
		keyName := keyToken.value
		err := p.match(tColon)
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		value, err := p.parseExpression(0)
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		node := ASTNode***REMOVED***
			nodeType: ASTKeyValPair,
			value:    keyName,
			children: []ASTNode***REMOVED***value***REMOVED***,
		***REMOVED***
		children = append(children, node)
		if p.current() == tComma ***REMOVED***
			err := p.match(tComma)
			if err != nil ***REMOVED***
				return ASTNode***REMOVED******REMOVED***, nil
			***REMOVED***
		***REMOVED*** else if p.current() == tRbrace ***REMOVED***
			err := p.match(tRbrace)
			if err != nil ***REMOVED***
				return ASTNode***REMOVED******REMOVED***, nil
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return ASTNode***REMOVED***
		nodeType: ASTMultiSelectHash,
		children: children,
	***REMOVED***, nil
***REMOVED***

func (p *Parser) projectIfSlice(left ASTNode, right ASTNode) (ASTNode, error) ***REMOVED***
	indexExpr := ASTNode***REMOVED***
		nodeType: ASTIndexExpression,
		children: []ASTNode***REMOVED***left, right***REMOVED***,
	***REMOVED***
	if right.nodeType == ASTSlice ***REMOVED***
		right, err := p.parseProjectionRHS(bindingPowers[tStar])
		return ASTNode***REMOVED***
			nodeType: ASTProjection,
			children: []ASTNode***REMOVED***indexExpr, right***REMOVED***,
		***REMOVED***, err
	***REMOVED***
	return indexExpr, nil
***REMOVED***
func (p *Parser) parseFilter(node ASTNode) (ASTNode, error) ***REMOVED***
	var right, condition ASTNode
	var err error
	condition, err = p.parseExpression(0)
	if err != nil ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, err
	***REMOVED***
	if err := p.match(tRbracket); err != nil ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, err
	***REMOVED***
	if p.current() == tFlatten ***REMOVED***
		right = ASTNode***REMOVED***nodeType: ASTIdentity***REMOVED***
	***REMOVED*** else ***REMOVED***
		right, err = p.parseProjectionRHS(bindingPowers[tFilter])
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***

	return ASTNode***REMOVED***
		nodeType: ASTFilterProjection,
		children: []ASTNode***REMOVED***node, right, condition***REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (p *Parser) parseDotRHS(bindingPower int) (ASTNode, error) ***REMOVED***
	lookahead := p.current()
	if tokensOneOf([]tokType***REMOVED***tQuotedIdentifier, tUnquotedIdentifier, tStar***REMOVED***, lookahead) ***REMOVED***
		return p.parseExpression(bindingPower)
	***REMOVED*** else if lookahead == tLbracket ***REMOVED***
		if err := p.match(tLbracket); err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return p.parseMultiSelectList()
	***REMOVED*** else if lookahead == tLbrace ***REMOVED***
		if err := p.match(tLbrace); err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return p.parseMultiSelectHash()
	***REMOVED***
	return ASTNode***REMOVED******REMOVED***, p.syntaxError("Expected identifier, lbracket, or lbrace")
***REMOVED***

func (p *Parser) parseProjectionRHS(bindingPower int) (ASTNode, error) ***REMOVED***
	current := p.current()
	if bindingPowers[current] < 10 ***REMOVED***
		return ASTNode***REMOVED***nodeType: ASTIdentity***REMOVED***, nil
	***REMOVED*** else if current == tLbracket ***REMOVED***
		return p.parseExpression(bindingPower)
	***REMOVED*** else if current == tFilter ***REMOVED***
		return p.parseExpression(bindingPower)
	***REMOVED*** else if current == tDot ***REMOVED***
		err := p.match(tDot)
		if err != nil ***REMOVED***
			return ASTNode***REMOVED******REMOVED***, err
		***REMOVED***
		return p.parseDotRHS(bindingPower)
	***REMOVED*** else ***REMOVED***
		return ASTNode***REMOVED******REMOVED***, p.syntaxError("Error")
	***REMOVED***
***REMOVED***

func (p *Parser) lookahead(number int) tokType ***REMOVED***
	return p.lookaheadToken(number).tokenType
***REMOVED***

func (p *Parser) current() tokType ***REMOVED***
	return p.lookahead(0)
***REMOVED***

func (p *Parser) lookaheadToken(number int) token ***REMOVED***
	return p.tokens[p.index+number]
***REMOVED***

func (p *Parser) advance() ***REMOVED***
	p.index++
***REMOVED***

func tokensOneOf(elements []tokType, token tokType) bool ***REMOVED***
	for _, elem := range elements ***REMOVED***
		if elem == token ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (p *Parser) syntaxError(msg string) SyntaxError ***REMOVED***
	return SyntaxError***REMOVED***
		msg:        msg,
		Expression: p.expression,
		Offset:     p.lookaheadToken(0).position,
	***REMOVED***
***REMOVED***

// Create a SyntaxError based on the provided token.
// This differs from syntaxError() which creates a SyntaxError
// based on the current lookahead token.
func (p *Parser) syntaxErrorToken(msg string, t token) SyntaxError ***REMOVED***
	return SyntaxError***REMOVED***
		msg:        msg,
		Expression: p.expression,
		Offset:     t.position,
	***REMOVED***
***REMOVED***
