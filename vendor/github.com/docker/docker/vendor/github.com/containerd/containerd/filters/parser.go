package filters

import (
	"fmt"
	"io"

	"github.com/containerd/containerd/errdefs"
	"github.com/pkg/errors"
)

/*
Parse the strings into a filter that may be used with an adaptor.

The filter is made up of zero or more selectors.

The format is a comma separated list of expressions, in the form of
`<fieldpath><op><value>`, known as selectors. All selectors must match the
target object for the filter to be true.

We define the operators "==" for equality, "!=" for not equal and "~=" for a
regular expression. If the operator and value are not present, the matcher will
test for the presence of a value, as defined by the target object.

The formal grammar is as follows:

selectors := selector ("," selector)*
selector  := fieldpath (operator value)
fieldpath := field ('.' field)*
field     := quoted | [A-Za-z] [A-Za-z0-9_]+
operator  := "==" | "!=" | "~="
value     := quoted | [^\s,]+
quoted    := <go string syntax>

*/
func Parse(s string) (Filter, error) ***REMOVED***
	// special case empty to match all
	if s == "" ***REMOVED***
		return Always, nil
	***REMOVED***

	p := parser***REMOVED***input: s***REMOVED***
	return p.parse()
***REMOVED***

// ParseAll parses each filter in ss and returns a filter that will return true
// if any filter matches the expression.
//
// If no filters are provided, the filter will match anything.
func ParseAll(ss ...string) (Filter, error) ***REMOVED***
	if len(ss) == 0 ***REMOVED***
		return Always, nil
	***REMOVED***

	var fs []Filter
	for _, s := range ss ***REMOVED***
		f, err := Parse(s)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(errdefs.ErrInvalidArgument, err.Error())
		***REMOVED***

		fs = append(fs, f)
	***REMOVED***

	return Any(fs), nil
***REMOVED***

type parser struct ***REMOVED***
	input   string
	scanner scanner
***REMOVED***

func (p *parser) parse() (Filter, error) ***REMOVED***
	p.scanner.init(p.input)

	ss, err := p.selectors()
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "filters")
	***REMOVED***

	return ss, nil
***REMOVED***

func (p *parser) selectors() (Filter, error) ***REMOVED***
	s, err := p.selector()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ss := All***REMOVED***s***REMOVED***

loop:
	for ***REMOVED***
		tok := p.scanner.peek()
		switch tok ***REMOVED***
		case ',':
			pos, tok, _ := p.scanner.scan()
			if tok != tokenSeparator ***REMOVED***
				return nil, p.mkerr(pos, "expected a separator")
			***REMOVED***

			s, err := p.selector()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			ss = append(ss, s)
		case tokenEOF:
			break loop
		default:
			return nil, p.mkerr(p.scanner.ppos, "unexpected input: %v", string(tok))
		***REMOVED***
	***REMOVED***

	return ss, nil
***REMOVED***

func (p *parser) selector() (selector, error) ***REMOVED***
	fieldpath, err := p.fieldpath()
	if err != nil ***REMOVED***
		return selector***REMOVED******REMOVED***, err
	***REMOVED***

	switch p.scanner.peek() ***REMOVED***
	case ',', tokenSeparator, tokenEOF:
		return selector***REMOVED***
			fieldpath: fieldpath,
			operator:  operatorPresent,
		***REMOVED***, nil
	***REMOVED***

	op, err := p.operator()
	if err != nil ***REMOVED***
		return selector***REMOVED******REMOVED***, err
	***REMOVED***

	var allowAltQuotes bool
	if op == operatorMatches ***REMOVED***
		allowAltQuotes = true
	***REMOVED***

	value, err := p.value(allowAltQuotes)
	if err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			return selector***REMOVED******REMOVED***, io.ErrUnexpectedEOF
		***REMOVED***
		return selector***REMOVED******REMOVED***, err
	***REMOVED***

	return selector***REMOVED***
		fieldpath: fieldpath,
		value:     value,
		operator:  op,
	***REMOVED***, nil
***REMOVED***

func (p *parser) fieldpath() ([]string, error) ***REMOVED***
	f, err := p.field()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	fs := []string***REMOVED***f***REMOVED***
loop:
	for ***REMOVED***
		tok := p.scanner.peek() // lookahead to consume field separator

		switch tok ***REMOVED***
		case '.':
			pos, tok, _ := p.scanner.scan() // consume separator
			if tok != tokenSeparator ***REMOVED***
				return nil, p.mkerr(pos, "expected a field separator (`.`)")
			***REMOVED***

			f, err := p.field()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			fs = append(fs, f)
		default:
			// let the layer above handle the other bad cases.
			break loop
		***REMOVED***
	***REMOVED***

	return fs, nil
***REMOVED***

func (p *parser) field() (string, error) ***REMOVED***
	pos, tok, s := p.scanner.scan()
	switch tok ***REMOVED***
	case tokenField:
		return s, nil
	case tokenQuoted:
		return p.unquote(pos, s, false)
	***REMOVED***

	return "", p.mkerr(pos, "expected field or quoted")
***REMOVED***

func (p *parser) operator() (operator, error) ***REMOVED***
	pos, tok, s := p.scanner.scan()
	switch tok ***REMOVED***
	case tokenOperator:
		switch s ***REMOVED***
		case "==":
			return operatorEqual, nil
		case "!=":
			return operatorNotEqual, nil
		case "~=":
			return operatorMatches, nil
		default:
			return 0, p.mkerr(pos, "unsupported operator %q", s)
		***REMOVED***
	***REMOVED***

	return 0, p.mkerr(pos, `expected an operator ("=="|"!="|"~=")`)
***REMOVED***

func (p *parser) value(allowAltQuotes bool) (string, error) ***REMOVED***
	pos, tok, s := p.scanner.scan()

	switch tok ***REMOVED***
	case tokenValue, tokenField:
		return s, nil
	case tokenQuoted:
		return p.unquote(pos, s, allowAltQuotes)
	***REMOVED***

	return "", p.mkerr(pos, "expected value or quoted")
***REMOVED***

func (p *parser) unquote(pos int, s string, allowAlts bool) (string, error) ***REMOVED***
	if !allowAlts && s[0] != '\'' && s[0] != '"' ***REMOVED***
		return "", p.mkerr(pos, "invalid quote encountered")
	***REMOVED***

	uq, err := unquote(s)
	if err != nil ***REMOVED***
		return "", p.mkerr(pos, "unquoting failed: %v", err)
	***REMOVED***

	return uq, nil
***REMOVED***

type parseError struct ***REMOVED***
	input string
	pos   int
	msg   string
***REMOVED***

func (pe parseError) Error() string ***REMOVED***
	if pe.pos < len(pe.input) ***REMOVED***
		before := pe.input[:pe.pos]
		location := pe.input[pe.pos : pe.pos+1] // need to handle end
		after := pe.input[pe.pos+1:]

		return fmt.Sprintf("[%s >|%s|< %s]: %v", before, location, after, pe.msg)
	***REMOVED***

	return fmt.Sprintf("[%s]: %v", pe.input, pe.msg)
***REMOVED***

func (p *parser) mkerr(pos int, format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return errors.Wrap(parseError***REMOVED***
		input: p.input,
		pos:   pos,
		msg:   fmt.Sprintf(format, args...),
	***REMOVED***, "parse error")
***REMOVED***
