// Package ast declares the types used to represent syntax trees for HCL
// (HashiCorp Configuration Language)
package ast

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/hcl/token"
)

// Node is an element in the abstract syntax tree.
type Node interface ***REMOVED***
	node()
	Pos() token.Pos
***REMOVED***

func (File) node()         ***REMOVED******REMOVED***
func (ObjectList) node()   ***REMOVED******REMOVED***
func (ObjectKey) node()    ***REMOVED******REMOVED***
func (ObjectItem) node()   ***REMOVED******REMOVED***
func (Comment) node()      ***REMOVED******REMOVED***
func (CommentGroup) node() ***REMOVED******REMOVED***
func (ObjectType) node()   ***REMOVED******REMOVED***
func (LiteralType) node()  ***REMOVED******REMOVED***
func (ListType) node()     ***REMOVED******REMOVED***

// File represents a single HCL file
type File struct ***REMOVED***
	Node     Node            // usually a *ObjectList
	Comments []*CommentGroup // list of all comments in the source
***REMOVED***

func (f *File) Pos() token.Pos ***REMOVED***
	return f.Node.Pos()
***REMOVED***

// ObjectList represents a list of ObjectItems. An HCL file itself is an
// ObjectList.
type ObjectList struct ***REMOVED***
	Items []*ObjectItem
***REMOVED***

func (o *ObjectList) Add(item *ObjectItem) ***REMOVED***
	o.Items = append(o.Items, item)
***REMOVED***

// Filter filters out the objects with the given key list as a prefix.
//
// The returned list of objects contain ObjectItems where the keys have
// this prefix already stripped off. This might result in objects with
// zero-length key lists if they have no children.
//
// If no matches are found, an empty ObjectList (non-nil) is returned.
func (o *ObjectList) Filter(keys ...string) *ObjectList ***REMOVED***
	var result ObjectList
	for _, item := range o.Items ***REMOVED***
		// If there aren't enough keys, then ignore this
		if len(item.Keys) < len(keys) ***REMOVED***
			continue
		***REMOVED***

		match := true
		for i, key := range item.Keys[:len(keys)] ***REMOVED***
			key := key.Token.Value().(string)
			if key != keys[i] && !strings.EqualFold(key, keys[i]) ***REMOVED***
				match = false
				break
			***REMOVED***
		***REMOVED***
		if !match ***REMOVED***
			continue
		***REMOVED***

		// Strip off the prefix from the children
		newItem := *item
		newItem.Keys = newItem.Keys[len(keys):]
		result.Add(&newItem)
	***REMOVED***

	return &result
***REMOVED***

// Children returns further nested objects (key length > 0) within this
// ObjectList. This should be used with Filter to get at child items.
func (o *ObjectList) Children() *ObjectList ***REMOVED***
	var result ObjectList
	for _, item := range o.Items ***REMOVED***
		if len(item.Keys) > 0 ***REMOVED***
			result.Add(item)
		***REMOVED***
	***REMOVED***

	return &result
***REMOVED***

// Elem returns items in the list that are direct element assignments
// (key length == 0). This should be used with Filter to get at elements.
func (o *ObjectList) Elem() *ObjectList ***REMOVED***
	var result ObjectList
	for _, item := range o.Items ***REMOVED***
		if len(item.Keys) == 0 ***REMOVED***
			result.Add(item)
		***REMOVED***
	***REMOVED***

	return &result
***REMOVED***

func (o *ObjectList) Pos() token.Pos ***REMOVED***
	// always returns the uninitiliazed position
	return o.Items[0].Pos()
***REMOVED***

// ObjectItem represents a HCL Object Item. An item is represented with a key
// (or keys). It can be an assignment or an object (both normal and nested)
type ObjectItem struct ***REMOVED***
	// keys is only one length long if it's of type assignment. If it's a
	// nested object it can be larger than one. In that case "assign" is
	// invalid as there is no assignments for a nested object.
	Keys []*ObjectKey

	// assign contains the position of "=", if any
	Assign token.Pos

	// val is the item itself. It can be an object,list, number, bool or a
	// string. If key length is larger than one, val can be only of type
	// Object.
	Val Node

	LeadComment *CommentGroup // associated lead comment
	LineComment *CommentGroup // associated line comment
***REMOVED***

func (o *ObjectItem) Pos() token.Pos ***REMOVED***
	// I'm not entirely sure what causes this, but removing this causes
	// a test failure. We should investigate at some point.
	if len(o.Keys) == 0 ***REMOVED***
		return token.Pos***REMOVED******REMOVED***
	***REMOVED***

	return o.Keys[0].Pos()
***REMOVED***

// ObjectKeys are either an identifier or of type string.
type ObjectKey struct ***REMOVED***
	Token token.Token
***REMOVED***

func (o *ObjectKey) Pos() token.Pos ***REMOVED***
	return o.Token.Pos
***REMOVED***

// LiteralType represents a literal of basic type. Valid types are:
// token.NUMBER, token.FLOAT, token.BOOL and token.STRING
type LiteralType struct ***REMOVED***
	Token token.Token

	// comment types, only used when in a list
	LeadComment *CommentGroup
	LineComment *CommentGroup
***REMOVED***

func (l *LiteralType) Pos() token.Pos ***REMOVED***
	return l.Token.Pos
***REMOVED***

// ListStatement represents a HCL List type
type ListType struct ***REMOVED***
	Lbrack token.Pos // position of "["
	Rbrack token.Pos // position of "]"
	List   []Node    // the elements in lexical order
***REMOVED***

func (l *ListType) Pos() token.Pos ***REMOVED***
	return l.Lbrack
***REMOVED***

func (l *ListType) Add(node Node) ***REMOVED***
	l.List = append(l.List, node)
***REMOVED***

// ObjectType represents a HCL Object Type
type ObjectType struct ***REMOVED***
	Lbrace token.Pos   // position of "***REMOVED***"
	Rbrace token.Pos   // position of "***REMOVED***"
	List   *ObjectList // the nodes in lexical order
***REMOVED***

func (o *ObjectType) Pos() token.Pos ***REMOVED***
	return o.Lbrace
***REMOVED***

// Comment node represents a single //, # style or /*- style commment
type Comment struct ***REMOVED***
	Start token.Pos // position of / or #
	Text  string
***REMOVED***

func (c *Comment) Pos() token.Pos ***REMOVED***
	return c.Start
***REMOVED***

// CommentGroup node represents a sequence of comments with no other tokens and
// no empty lines between.
type CommentGroup struct ***REMOVED***
	List []*Comment // len(List) > 0
***REMOVED***

func (c *CommentGroup) Pos() token.Pos ***REMOVED***
	return c.List[0].Pos()
***REMOVED***

//-------------------------------------------------------------------
// GoStringer
//-------------------------------------------------------------------

func (o *ObjectKey) GoString() string  ***REMOVED*** return fmt.Sprintf("*%#v", *o) ***REMOVED***
func (o *ObjectList) GoString() string ***REMOVED*** return fmt.Sprintf("*%#v", *o) ***REMOVED***
