package ast

import "fmt"

// WalkFunc describes a function to be called for each node during a Walk. The
// returned node can be used to rewrite the AST. Walking stops the returned
// bool is false.
type WalkFunc func(Node) (Node, bool)

// Walk traverses an AST in depth-first order: It starts by calling fn(node);
// node must not be nil. If fn returns true, Walk invokes fn recursively for
// each of the non-nil children of node, followed by a call of fn(nil). The
// returned node of fn can be used to rewrite the passed node to fn.
func Walk(node Node, fn WalkFunc) Node ***REMOVED***
	rewritten, ok := fn(node)
	if !ok ***REMOVED***
		return rewritten
	***REMOVED***

	switch n := node.(type) ***REMOVED***
	case *File:
		n.Node = Walk(n.Node, fn)
	case *ObjectList:
		for i, item := range n.Items ***REMOVED***
			n.Items[i] = Walk(item, fn).(*ObjectItem)
		***REMOVED***
	case *ObjectKey:
		// nothing to do
	case *ObjectItem:
		for i, k := range n.Keys ***REMOVED***
			n.Keys[i] = Walk(k, fn).(*ObjectKey)
		***REMOVED***

		if n.Val != nil ***REMOVED***
			n.Val = Walk(n.Val, fn)
		***REMOVED***
	case *LiteralType:
		// nothing to do
	case *ListType:
		for i, l := range n.List ***REMOVED***
			n.List[i] = Walk(l, fn)
		***REMOVED***
	case *ObjectType:
		n.List = Walk(n.List, fn).(*ObjectList)
	default:
		// should we panic here?
		fmt.Printf("unknown type: %T\n", n)
	***REMOVED***

	fn(nil)
	return rewritten
***REMOVED***
