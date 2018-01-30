// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
)

// checkTreeConsistency checks that a node and its descendants are all
// consistent in their parent/child/sibling relationships.
func checkTreeConsistency(n *Node) error ***REMOVED***
	return checkTreeConsistency1(n, 0)
***REMOVED***

func checkTreeConsistency1(n *Node, depth int) error ***REMOVED***
	if depth == 1e4 ***REMOVED***
		return fmt.Errorf("html: tree looks like it contains a cycle")
	***REMOVED***
	if err := checkNodeConsistency(n); err != nil ***REMOVED***
		return err
	***REMOVED***
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if err := checkTreeConsistency1(c, depth+1); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// checkNodeConsistency checks that a node's parent/child/sibling relationships
// are consistent.
func checkNodeConsistency(n *Node) error ***REMOVED***
	if n == nil ***REMOVED***
		return nil
	***REMOVED***

	nParent := 0
	for p := n.Parent; p != nil; p = p.Parent ***REMOVED***
		nParent++
		if nParent == 1e4 ***REMOVED***
			return fmt.Errorf("html: parent list looks like an infinite loop")
		***REMOVED***
	***REMOVED***

	nForward := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		nForward++
		if nForward == 1e6 ***REMOVED***
			return fmt.Errorf("html: forward list of children looks like an infinite loop")
		***REMOVED***
		if c.Parent != n ***REMOVED***
			return fmt.Errorf("html: inconsistent child/parent relationship")
		***REMOVED***
	***REMOVED***

	nBackward := 0
	for c := n.LastChild; c != nil; c = c.PrevSibling ***REMOVED***
		nBackward++
		if nBackward == 1e6 ***REMOVED***
			return fmt.Errorf("html: backward list of children looks like an infinite loop")
		***REMOVED***
		if c.Parent != n ***REMOVED***
			return fmt.Errorf("html: inconsistent child/parent relationship")
		***REMOVED***
	***REMOVED***

	if n.Parent != nil ***REMOVED***
		if n.Parent == n ***REMOVED***
			return fmt.Errorf("html: inconsistent parent relationship")
		***REMOVED***
		if n.Parent == n.FirstChild ***REMOVED***
			return fmt.Errorf("html: inconsistent parent/first relationship")
		***REMOVED***
		if n.Parent == n.LastChild ***REMOVED***
			return fmt.Errorf("html: inconsistent parent/last relationship")
		***REMOVED***
		if n.Parent == n.PrevSibling ***REMOVED***
			return fmt.Errorf("html: inconsistent parent/prev relationship")
		***REMOVED***
		if n.Parent == n.NextSibling ***REMOVED***
			return fmt.Errorf("html: inconsistent parent/next relationship")
		***REMOVED***

		parentHasNAsAChild := false
		for c := n.Parent.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
			if c == n ***REMOVED***
				parentHasNAsAChild = true
				break
			***REMOVED***
		***REMOVED***
		if !parentHasNAsAChild ***REMOVED***
			return fmt.Errorf("html: inconsistent parent/child relationship")
		***REMOVED***
	***REMOVED***

	if n.PrevSibling != nil && n.PrevSibling.NextSibling != n ***REMOVED***
		return fmt.Errorf("html: inconsistent prev/next relationship")
	***REMOVED***
	if n.NextSibling != nil && n.NextSibling.PrevSibling != n ***REMOVED***
		return fmt.Errorf("html: inconsistent next/prev relationship")
	***REMOVED***

	if (n.FirstChild == nil) != (n.LastChild == nil) ***REMOVED***
		return fmt.Errorf("html: inconsistent first/last relationship")
	***REMOVED***
	if n.FirstChild != nil && n.FirstChild == n.LastChild ***REMOVED***
		// We have a sole child.
		if n.FirstChild.PrevSibling != nil || n.FirstChild.NextSibling != nil ***REMOVED***
			return fmt.Errorf("html: inconsistent sole child's sibling relationship")
		***REMOVED***
	***REMOVED***

	seen := map[*Node]bool***REMOVED******REMOVED***

	var last *Node
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if seen[c] ***REMOVED***
			return fmt.Errorf("html: inconsistent repeated child")
		***REMOVED***
		seen[c] = true
		last = c
	***REMOVED***
	if last != n.LastChild ***REMOVED***
		return fmt.Errorf("html: inconsistent last relationship")
	***REMOVED***

	var first *Node
	for c := n.LastChild; c != nil; c = c.PrevSibling ***REMOVED***
		if !seen[c] ***REMOVED***
			return fmt.Errorf("html: inconsistent missing child")
		***REMOVED***
		delete(seen, c)
		first = c
	***REMOVED***
	if first != n.FirstChild ***REMOVED***
		return fmt.Errorf("html: inconsistent first relationship")
	***REMOVED***

	if len(seen) != 0 ***REMOVED***
		return fmt.Errorf("html: inconsistent forwards/backwards child list")
	***REMOVED***

	return nil
***REMOVED***
