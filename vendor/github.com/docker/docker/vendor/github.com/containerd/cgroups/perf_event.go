package cgroups

import "path/filepath"

func NewPerfEvent(root string) *PerfEventController ***REMOVED***
	return &PerfEventController***REMOVED***
		root: filepath.Join(root, string(PerfEvent)),
	***REMOVED***
***REMOVED***

type PerfEventController struct ***REMOVED***
	root string
***REMOVED***

func (p *PerfEventController) Name() Name ***REMOVED***
	return PerfEvent
***REMOVED***

func (p *PerfEventController) Path(path string) string ***REMOVED***
	return filepath.Join(p.root, path)
***REMOVED***
