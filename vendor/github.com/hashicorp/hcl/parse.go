package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/ast"
	hclParser "github.com/hashicorp/hcl/hcl/parser"
	jsonParser "github.com/hashicorp/hcl/json/parser"
)

// ParseBytes accepts as input byte slice and returns ast tree.
//
// Input can be either JSON or HCL
func ParseBytes(in []byte) (*ast.File, error) ***REMOVED***
	return parse(in)
***REMOVED***

// ParseString accepts input as a string and returns ast tree.
func ParseString(input string) (*ast.File, error) ***REMOVED***
	return parse([]byte(input))
***REMOVED***

func parse(in []byte) (*ast.File, error) ***REMOVED***
	switch lexMode(in) ***REMOVED***
	case lexModeHcl:
		return hclParser.Parse(in)
	case lexModeJson:
		return jsonParser.Parse(in)
	***REMOVED***

	return nil, fmt.Errorf("unknown config format")
***REMOVED***

// Parse parses the given input and returns the root object.
//
// The input format can be either HCL or JSON.
func Parse(input string) (*ast.File, error) ***REMOVED***
	return parse([]byte(input))
***REMOVED***
