// Package gc experiments with providing central gc tooling to ensure
// deterministic resource removal within containerd.
//
// For now, we just have a single exported implementation that can be used
// under certain use cases.
package gc

import (
	"context"
	"sync"
)

// ResourceType represents type of resource at a node
type ResourceType uint8

// Node presents a resource which has a type and key,
// this node can be used to lookup other nodes.
type Node struct ***REMOVED***
	Type      ResourceType
	Namespace string
	Key       string
***REMOVED***

// Tricolor implements basic, single-thread tri-color GC. Given the roots, the
// complete set and a refs function, this function returns a map of all
// reachable objects.
//
// Correct usage requires that the caller not allow the arguments to change
// until the result is used to delete objects in the system.
//
// It will allocate memory proportional to the size of the reachable set.
//
// We can probably use this to inform a design for incremental GC by injecting
// callbacks to the set modification algorithms.
func Tricolor(roots []Node, refs func(ref Node) ([]Node, error)) (map[Node]struct***REMOVED******REMOVED***, error) ***REMOVED***
	var (
		grays     []Node                // maintain a gray "stack"
		seen      = map[Node]struct***REMOVED******REMOVED******REMOVED******REMOVED*** // or not "white", basically "seen"
		reachable = map[Node]struct***REMOVED******REMOVED******REMOVED******REMOVED*** // or "black", in tri-color parlance
	)

	grays = append(grays, roots...)

	for len(grays) > 0 ***REMOVED***
		// Pick any gray object
		id := grays[len(grays)-1] // effectively "depth first" because first element
		grays = grays[:len(grays)-1]
		seen[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED*** // post-mark this as not-white
		rs, err := refs(id)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// mark all the referenced objects as gray
		for _, target := range rs ***REMOVED***
			if _, ok := seen[target]; !ok ***REMOVED***
				grays = append(grays, target)
			***REMOVED***
		***REMOVED***

		// mark as black when done
		reachable[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	return reachable, nil
***REMOVED***

// ConcurrentMark implements simple, concurrent GC. All the roots are scanned
// and the complete set of references is formed by calling the refs function
// for each seen object. This function returns a map of all object reachable
// from a root.
//
// Correct usage requires that the caller not allow the arguments to change
// until the result is used to delete objects in the system.
//
// It will allocate memory proportional to the size of the reachable set.
func ConcurrentMark(ctx context.Context, root <-chan Node, refs func(context.Context, Node, func(Node)) error) (map[Node]struct***REMOVED******REMOVED***, error) ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		grays = make(chan Node)
		seen  = map[Node]struct***REMOVED******REMOVED******REMOVED******REMOVED*** // or not "white", basically "seen"
		wg    sync.WaitGroup

		errOnce sync.Once
		refErr  error
	)

	go func() ***REMOVED***
		for gray := range grays ***REMOVED***
			if _, ok := seen[gray]; ok ***REMOVED***
				wg.Done()
				continue
			***REMOVED***
			seen[gray] = struct***REMOVED******REMOVED******REMOVED******REMOVED*** // post-mark this as non-white

			go func(gray Node) ***REMOVED***
				defer wg.Done()

				send := func(n Node) ***REMOVED***
					wg.Add(1)
					select ***REMOVED***
					case grays <- n:
					case <-ctx.Done():
						wg.Done()
					***REMOVED***
				***REMOVED***

				if err := refs(ctx, gray, send); err != nil ***REMOVED***
					errOnce.Do(func() ***REMOVED***
						refErr = err
						cancel()
					***REMOVED***)
				***REMOVED***

			***REMOVED***(gray)
		***REMOVED***
	***REMOVED***()

	for r := range root ***REMOVED***
		wg.Add(1)
		select ***REMOVED***
		case grays <- r:
		case <-ctx.Done():
			wg.Done()
		***REMOVED***

	***REMOVED***

	// Wait for outstanding grays to be processed
	wg.Wait()

	close(grays)

	if refErr != nil ***REMOVED***
		return nil, refErr
	***REMOVED***
	if cErr := ctx.Err(); cErr != nil ***REMOVED***
		return nil, cErr
	***REMOVED***

	return seen, nil
***REMOVED***

// Sweep removes all nodes returned through the channel which are not in
// the reachable set by calling the provided remove function.
func Sweep(reachable map[Node]struct***REMOVED******REMOVED***, all []Node, remove func(Node) error) error ***REMOVED***
	// All black objects are now reachable, and all white objects are
	// unreachable. Free those that are white!
	for _, node := range all ***REMOVED***
		if _, ok := reachable[node]; !ok ***REMOVED***
			if err := remove(node); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
