// Copyright 2012 Neal van Veen. All rights reserved.
// Usage of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package gotty

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var exp = [...]string***REMOVED***
	"%%",
	"%c",
	"%s",
	"%p(\\d)",
	"%P([A-z])",
	"%g([A-z])",
	"%'(.)'",
	"%***REMOVED***([0-9]+)***REMOVED***",
	"%l",
	"%\\+|%-|%\\*|%/|%m",
	"%&|%\\||%\\^",
	"%=|%>|%<",
	"%A|%O",
	"%!|%~",
	"%i",
	"%(:[\\ #\\-\\+]***REMOVED***0,4***REMOVED***)?(\\d+\\.\\d+|\\d+)?[doxXs]",
	"%\\?(.*?);",
***REMOVED***

var regex *regexp.Regexp
var staticVar map[byte]stacker

// Parses the attribute that is received with name attr and parameters params.
func (term *TermInfo) Parse(attr string, params ...interface***REMOVED******REMOVED***) (string, error) ***REMOVED***
	// Get the attribute name first.
	iface, err := term.GetAttribute(attr)
	str, ok := iface.(string)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if !ok ***REMOVED***
		return str, errors.New("Only string capabilities can be parsed.")
	***REMOVED***
	// Construct the hidden parser struct so we can use a recursive stack based
	// parser.
	ps := &parser***REMOVED******REMOVED***
	// Dynamic variables only exist in this context.
	ps.dynamicVar = make(map[byte]stacker, 26)
	ps.parameters = make([]stacker, len(params))
	// Convert the parameters to insert them into the parser struct.
	for i, x := range params ***REMOVED***
		ps.parameters[i] = x
	***REMOVED***
	// Recursively walk and return.
	result, err := ps.walk(str)
	return result, err
***REMOVED***

// Parses the attribute that is received with name attr and parameters params.
// Only works on full name of a capability that is given, which it uses to
// search for the termcap name.
func (term *TermInfo) ParseName(attr string, params ...interface***REMOVED******REMOVED***) (string, error) ***REMOVED***
	tc := GetTermcapName(attr)
	return term.Parse(tc, params)
***REMOVED***

// Identify each token in a stack based manner and do the actual parsing.
func (ps *parser) walk(attr string) (string, error) ***REMOVED***
	// We use a buffer to get the modified string.
	var buf bytes.Buffer
	// Next, find and identify all tokens by their indices and strings.
	tokens := regex.FindAllStringSubmatch(attr, -1)
	if len(tokens) == 0 ***REMOVED***
		return attr, nil
	***REMOVED***
	indices := regex.FindAllStringIndex(attr, -1)
	q := 0 // q counts the matches of one token
	// Iterate through the string per character.
	for i := 0; i < len(attr); i++ ***REMOVED***
		// If the current position is an identified token, execute the following
		// steps.
		if q < len(indices) && i >= indices[q][0] && i < indices[q][1] ***REMOVED***
			// Switch on token.
			switch ***REMOVED***
			case tokens[q][0][:2] == "%%":
				// Literal percentage character.
				buf.WriteByte('%')
			case tokens[q][0][:2] == "%c":
				// Pop a character.
				c, err := ps.st.pop()
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				buf.WriteByte(c.(byte))
			case tokens[q][0][:2] == "%s":
				// Pop a string.
				str, err := ps.st.pop()
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				if _, ok := str.(string); !ok ***REMOVED***
					return buf.String(), errors.New("Stack head is not a string")
				***REMOVED***
				buf.WriteString(str.(string))
			case tokens[q][0][:2] == "%p":
				// Push a parameter on the stack.
				index, err := strconv.ParseInt(tokens[q][1], 10, 8)
				index--
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				if int(index) >= len(ps.parameters) ***REMOVED***
					return buf.String(), errors.New("Parameters index out of bound")
				***REMOVED***
				ps.st.push(ps.parameters[index])
			case tokens[q][0][:2] == "%P":
				// Pop a variable from the stack as a dynamic or static variable.
				val, err := ps.st.pop()
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				index := tokens[q][2]
				if len(index) > 1 ***REMOVED***
					errorStr := fmt.Sprintf("%s is not a valid dynamic variables index",
						index)
					return buf.String(), errors.New(errorStr)
				***REMOVED***
				// Specify either dynamic or static.
				if index[0] >= 'a' && index[0] <= 'z' ***REMOVED***
					ps.dynamicVar[index[0]] = val
				***REMOVED*** else if index[0] >= 'A' && index[0] <= 'Z' ***REMOVED***
					staticVar[index[0]] = val
				***REMOVED***
			case tokens[q][0][:2] == "%g":
				// Push a variable from the stack as a dynamic or static variable.
				index := tokens[q][3]
				if len(index) > 1 ***REMOVED***
					errorStr := fmt.Sprintf("%s is not a valid static variables index",
						index)
					return buf.String(), errors.New(errorStr)
				***REMOVED***
				var val stacker
				if index[0] >= 'a' && index[0] <= 'z' ***REMOVED***
					val = ps.dynamicVar[index[0]]
				***REMOVED*** else if index[0] >= 'A' && index[0] <= 'Z' ***REMOVED***
					val = staticVar[index[0]]
				***REMOVED***
				ps.st.push(val)
			case tokens[q][0][:2] == "%'":
				// Push a character constant.
				con := tokens[q][4]
				if len(con) > 1 ***REMOVED***
					errorStr := fmt.Sprintf("%s is not a valid character constant", con)
					return buf.String(), errors.New(errorStr)
				***REMOVED***
				ps.st.push(con[0])
			case tokens[q][0][:2] == "%***REMOVED***":
				// Push an integer constant.
				con, err := strconv.ParseInt(tokens[q][5], 10, 32)
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				ps.st.push(con)
			case tokens[q][0][:2] == "%l":
				// Push the length of the string that is popped from the stack.
				popStr, err := ps.st.pop()
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				if _, ok := popStr.(string); !ok ***REMOVED***
					errStr := fmt.Sprintf("Stack head is not a string")
					return buf.String(), errors.New(errStr)
				***REMOVED***
				ps.st.push(len(popStr.(string)))
			case tokens[q][0][:2] == "%?":
				// If-then-else construct. First, the whole string is identified and
				// then inside this substring, we can specify which parts to switch on.
				ifReg, _ := regexp.Compile("%\\?(.*)%t(.*)%e(.*);|%\\?(.*)%t(.*);")
				ifTokens := ifReg.FindStringSubmatch(tokens[q][0])
				var (
					ifStr string
					err   error
				)
				// Parse the if-part to determine if-else.
				if len(ifTokens[1]) > 0 ***REMOVED***
					ifStr, err = ps.walk(ifTokens[1])
				***REMOVED*** else ***REMOVED*** // else
					ifStr, err = ps.walk(ifTokens[4])
				***REMOVED***
				// Return any errors
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED*** else if len(ifStr) > 0 ***REMOVED***
					// Self-defined limitation, not sure if this is correct, but didn't
					// seem like it.
					return buf.String(), errors.New("If-clause cannot print statements")
				***REMOVED***
				var thenStr string
				// Pop the first value that is set by parsing the if-clause.
				choose, err := ps.st.pop()
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				// Switch to if or else.
				if choose.(int) == 0 && len(ifTokens[1]) > 0 ***REMOVED***
					thenStr, err = ps.walk(ifTokens[3])
				***REMOVED*** else if choose.(int) != 0 ***REMOVED***
					if len(ifTokens[1]) > 0 ***REMOVED***
						thenStr, err = ps.walk(ifTokens[2])
					***REMOVED*** else ***REMOVED***
						thenStr, err = ps.walk(ifTokens[5])
					***REMOVED***
				***REMOVED***
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				buf.WriteString(thenStr)
			case tokens[q][0][len(tokens[q][0])-1] == 'd': // Fallthrough for printing
				fallthrough
			case tokens[q][0][len(tokens[q][0])-1] == 'o': // digits.
				fallthrough
			case tokens[q][0][len(tokens[q][0])-1] == 'x':
				fallthrough
			case tokens[q][0][len(tokens[q][0])-1] == 'X':
				fallthrough
			case tokens[q][0][len(tokens[q][0])-1] == 's':
				token := tokens[q][0]
				// Remove the : that comes before a flag.
				if token[1] == ':' ***REMOVED***
					token = token[:1] + token[2:]
				***REMOVED***
				digit, err := ps.st.pop()
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				// The rest is determined like the normal formatted prints.
				digitStr := fmt.Sprintf(token, digit.(int))
				buf.WriteString(digitStr)
			case tokens[q][0][:2] == "%i":
				// Increment the parameters by one.
				if len(ps.parameters) < 2 ***REMOVED***
					return buf.String(), errors.New("Not enough parameters to increment.")
				***REMOVED***
				val1, val2 := ps.parameters[0].(int), ps.parameters[1].(int)
				val1++
				val2++
				ps.parameters[0], ps.parameters[1] = val1, val2
			default:
				// The rest of the tokens is a special case, where two values are
				// popped and then operated on by the token that comes after them.
				op1, err := ps.st.pop()
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				op2, err := ps.st.pop()
				if err != nil ***REMOVED***
					return buf.String(), err
				***REMOVED***
				var result stacker
				switch tokens[q][0][:2] ***REMOVED***
				case "%+":
					// Addition
					result = op2.(int) + op1.(int)
				case "%-":
					// Subtraction
					result = op2.(int) - op1.(int)
				case "%*":
					// Multiplication
					result = op2.(int) * op1.(int)
				case "%/":
					// Division
					result = op2.(int) / op1.(int)
				case "%m":
					// Modulo
					result = op2.(int) % op1.(int)
				case "%&":
					// Bitwise AND
					result = op2.(int) & op1.(int)
				case "%|":
					// Bitwise OR
					result = op2.(int) | op1.(int)
				case "%^":
					// Bitwise XOR
					result = op2.(int) ^ op1.(int)
				case "%=":
					// Equals
					result = op2 == op1
				case "%>":
					// Greater-than
					result = op2.(int) > op1.(int)
				case "%<":
					// Lesser-than
					result = op2.(int) < op1.(int)
				case "%A":
					// Logical AND
					result = op2.(bool) && op1.(bool)
				case "%O":
					// Logical OR
					result = op2.(bool) || op1.(bool)
				case "%!":
					// Logical complement
					result = !op1.(bool)
				case "%~":
					// Bitwise complement
					result = ^(op1.(int))
				***REMOVED***
				ps.st.push(result)
			***REMOVED***

			i = indices[q][1] - 1
			q++
		***REMOVED*** else ***REMOVED***
			// We are not "inside" a token, so just skip until the end or the next
			// token, and add all characters to the buffer.
			j := i
			if q != len(indices) ***REMOVED***
				for !(j >= indices[q][0] && j < indices[q][1]) ***REMOVED***
					j++
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				j = len(attr)
			***REMOVED***
			buf.WriteString(string(attr[i:j]))
			i = j
		***REMOVED***
	***REMOVED***
	// Return the buffer as a string.
	return buf.String(), nil
***REMOVED***

// Push a stacker-value onto the stack.
func (st *stack) push(s stacker) ***REMOVED***
	*st = append(*st, s)
***REMOVED***

// Pop a stacker-value from the stack.
func (st *stack) pop() (stacker, error) ***REMOVED***
	if len(*st) == 0 ***REMOVED***
		return nil, errors.New("Stack is empty.")
	***REMOVED***
	newStack := make(stack, len(*st)-1)
	val := (*st)[len(*st)-1]
	copy(newStack, (*st)[:len(*st)-1])
	*st = newStack
	return val, nil
***REMOVED***

// Initialize regexes and the static vars (that don't get changed between
// calls.
func init() ***REMOVED***
	// Initialize the main regex.
	expStr := strings.Join(exp[:], "|")
	regex, _ = regexp.Compile(expStr)
	// Initialize the static variables.
	staticVar = make(map[byte]stacker, 26)
***REMOVED***
