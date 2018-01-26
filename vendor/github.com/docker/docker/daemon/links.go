package daemon

import (
	"sync"

	"github.com/docker/docker/container"
)

// linkIndex stores link relationships between containers, including their specified alias
// The alias is the name the parent uses to reference the child
type linkIndex struct ***REMOVED***
	// idx maps a parent->alias->child relationship
	idx map[*container.Container]map[string]*container.Container
	// childIdx maps  child->parent->aliases
	childIdx map[*container.Container]map[*container.Container]map[string]struct***REMOVED******REMOVED***
	mu       sync.Mutex
***REMOVED***

func newLinkIndex() *linkIndex ***REMOVED***
	return &linkIndex***REMOVED***
		idx:      make(map[*container.Container]map[string]*container.Container),
		childIdx: make(map[*container.Container]map[*container.Container]map[string]struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// link adds indexes for the passed in parent/child/alias relationships
func (l *linkIndex) link(parent, child *container.Container, alias string) ***REMOVED***
	l.mu.Lock()

	if l.idx[parent] == nil ***REMOVED***
		l.idx[parent] = make(map[string]*container.Container)
	***REMOVED***
	l.idx[parent][alias] = child
	if l.childIdx[child] == nil ***REMOVED***
		l.childIdx[child] = make(map[*container.Container]map[string]struct***REMOVED******REMOVED***)
	***REMOVED***
	if l.childIdx[child][parent] == nil ***REMOVED***
		l.childIdx[child][parent] = make(map[string]struct***REMOVED******REMOVED***)
	***REMOVED***
	l.childIdx[child][parent][alias] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	l.mu.Unlock()
***REMOVED***

// unlink removes the requested alias for the given parent/child
func (l *linkIndex) unlink(alias string, child, parent *container.Container) ***REMOVED***
	l.mu.Lock()
	delete(l.idx[parent], alias)
	delete(l.childIdx[child], parent)
	l.mu.Unlock()
***REMOVED***

// children maps all the aliases-> children for the passed in parent
// aliases here are the aliases the parent uses to refer to the child
func (l *linkIndex) children(parent *container.Container) map[string]*container.Container ***REMOVED***
	l.mu.Lock()
	children := l.idx[parent]
	l.mu.Unlock()
	return children
***REMOVED***

// parents maps all the aliases->parent for the passed in child
// aliases here are the aliases the parents use to refer to the child
func (l *linkIndex) parents(child *container.Container) map[string]*container.Container ***REMOVED***
	l.mu.Lock()

	parents := make(map[string]*container.Container)
	for parent, aliases := range l.childIdx[child] ***REMOVED***
		for alias := range aliases ***REMOVED***
			parents[alias] = parent
		***REMOVED***
	***REMOVED***

	l.mu.Unlock()
	return parents
***REMOVED***

// delete deletes all link relationships referencing this container
func (l *linkIndex) delete(container *container.Container) []string ***REMOVED***
	l.mu.Lock()

	var aliases []string
	for alias, child := range l.idx[container] ***REMOVED***
		aliases = append(aliases, alias)
		delete(l.childIdx[child], container)
	***REMOVED***
	delete(l.idx, container)
	delete(l.childIdx, container)
	l.mu.Unlock()
	return aliases
***REMOVED***
