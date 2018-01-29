// +build go1.9 go1.8.typealias

package main

import (
	"go/ast"
)

func typeAliasSpec(name string, typ ast.Expr) *ast.TypeSpec ***REMOVED***
	return &ast.TypeSpec***REMOVED***
		Name:   ast.NewIdent(name),
		Assign: 1,
		Type:   typ,
	***REMOVED***
***REMOVED***

func isAliasTypeSpec(t *ast.TypeSpec) bool ***REMOVED***
	return t.Assign != 0
***REMOVED***
