package parser

import "github.com/hashicorp/hcl/hcl/ast"

// flattenObjects takes an AST node, walks it, and flattens
func flattenObjects(node ast.Node) ***REMOVED***
	ast.Walk(node, func(n ast.Node) (ast.Node, bool) ***REMOVED***
		// We only care about lists, because this is what we modify
		list, ok := n.(*ast.ObjectList)
		if !ok ***REMOVED***
			return n, true
		***REMOVED***

		// Rebuild the item list
		items := make([]*ast.ObjectItem, 0, len(list.Items))
		frontier := make([]*ast.ObjectItem, len(list.Items))
		copy(frontier, list.Items)
		for len(frontier) > 0 ***REMOVED***
			// Pop the current item
			n := len(frontier)
			item := frontier[n-1]
			frontier = frontier[:n-1]

			switch v := item.Val.(type) ***REMOVED***
			case *ast.ObjectType:
				items, frontier = flattenObjectType(v, item, items, frontier)
			case *ast.ListType:
				items, frontier = flattenListType(v, item, items, frontier)
			default:
				items = append(items, item)
			***REMOVED***
		***REMOVED***

		// Reverse the list since the frontier model runs things backwards
		for i := len(items)/2 - 1; i >= 0; i-- ***REMOVED***
			opp := len(items) - 1 - i
			items[i], items[opp] = items[opp], items[i]
		***REMOVED***

		// Done! Set the original items
		list.Items = items
		return n, true
	***REMOVED***)
***REMOVED***

func flattenListType(
	ot *ast.ListType,
	item *ast.ObjectItem,
	items []*ast.ObjectItem,
	frontier []*ast.ObjectItem) ([]*ast.ObjectItem, []*ast.ObjectItem) ***REMOVED***
	// If the list is empty, keep the original list
	if len(ot.List) == 0 ***REMOVED***
		items = append(items, item)
		return items, frontier
	***REMOVED***

	// All the elements of this object must also be objects!
	for _, subitem := range ot.List ***REMOVED***
		if _, ok := subitem.(*ast.ObjectType); !ok ***REMOVED***
			items = append(items, item)
			return items, frontier
		***REMOVED***
	***REMOVED***

	// Great! We have a match go through all the items and flatten
	for _, elem := range ot.List ***REMOVED***
		// Add it to the frontier so that we can recurse
		frontier = append(frontier, &ast.ObjectItem***REMOVED***
			Keys:        item.Keys,
			Assign:      item.Assign,
			Val:         elem,
			LeadComment: item.LeadComment,
			LineComment: item.LineComment,
		***REMOVED***)
	***REMOVED***

	return items, frontier
***REMOVED***

func flattenObjectType(
	ot *ast.ObjectType,
	item *ast.ObjectItem,
	items []*ast.ObjectItem,
	frontier []*ast.ObjectItem) ([]*ast.ObjectItem, []*ast.ObjectItem) ***REMOVED***
	// If the list has no items we do not have to flatten anything
	if ot.List.Items == nil ***REMOVED***
		items = append(items, item)
		return items, frontier
	***REMOVED***

	// All the elements of this object must also be objects!
	for _, subitem := range ot.List.Items ***REMOVED***
		if _, ok := subitem.Val.(*ast.ObjectType); !ok ***REMOVED***
			items = append(items, item)
			return items, frontier
		***REMOVED***
	***REMOVED***

	// Great! We have a match go through all the items and flatten
	for _, subitem := range ot.List.Items ***REMOVED***
		// Copy the new key
		keys := make([]*ast.ObjectKey, len(item.Keys)+len(subitem.Keys))
		copy(keys, item.Keys)
		copy(keys[len(item.Keys):], subitem.Keys)

		// Add it to the frontier so that we can recurse
		frontier = append(frontier, &ast.ObjectItem***REMOVED***
			Keys:        keys,
			Assign:      item.Assign,
			Val:         subitem.Val,
			LeadComment: item.LeadComment,
			LineComment: item.LineComment,
		***REMOVED***)
	***REMOVED***

	return items, frontier
***REMOVED***
