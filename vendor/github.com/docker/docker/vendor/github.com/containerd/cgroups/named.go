package cgroups

import "path/filepath"

func NewNamed(root string, name Name) *namedController ***REMOVED***
	return &namedController***REMOVED***
		root: root,
		name: name,
	***REMOVED***
***REMOVED***

type namedController struct ***REMOVED***
	root string
	name Name
***REMOVED***

func (n *namedController) Name() Name ***REMOVED***
	return n.name
***REMOVED***

func (n *namedController) Path(path string) string ***REMOVED***
	return filepath.Join(n.root, string(n.name), path)
***REMOVED***
