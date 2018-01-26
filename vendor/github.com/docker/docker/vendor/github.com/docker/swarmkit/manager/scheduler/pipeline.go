package scheduler

import (
	"sort"

	"github.com/docker/swarmkit/api"
)

var (
	defaultFilters = []Filter***REMOVED***
		// Always check for readiness first.
		&ReadyFilter***REMOVED******REMOVED***,
		&ResourceFilter***REMOVED******REMOVED***,
		&PluginFilter***REMOVED******REMOVED***,
		&ConstraintFilter***REMOVED******REMOVED***,
		&PlatformFilter***REMOVED******REMOVED***,
		&HostPortFilter***REMOVED******REMOVED***,
	***REMOVED***
)

type checklistEntry struct ***REMOVED***
	f       Filter
	enabled bool

	// failureCount counts the number of nodes that this filter failed
	// against.
	failureCount int
***REMOVED***

type checklistByFailures []checklistEntry

func (c checklistByFailures) Len() int           ***REMOVED*** return len(c) ***REMOVED***
func (c checklistByFailures) Swap(i, j int)      ***REMOVED*** c[i], c[j] = c[j], c[i] ***REMOVED***
func (c checklistByFailures) Less(i, j int) bool ***REMOVED*** return c[i].failureCount < c[j].failureCount ***REMOVED***

// Pipeline runs a set of filters against nodes.
type Pipeline struct ***REMOVED***
	// checklist is a slice of filters to run
	checklist []checklistEntry
***REMOVED***

// NewPipeline returns a pipeline with the default set of filters.
func NewPipeline() *Pipeline ***REMOVED***
	p := &Pipeline***REMOVED******REMOVED***

	for _, f := range defaultFilters ***REMOVED***
		p.checklist = append(p.checklist, checklistEntry***REMOVED***f: f***REMOVED***)
	***REMOVED***

	return p
***REMOVED***

// Process a node through the filter pipeline.
// Returns true if all filters pass, false otherwise.
func (p *Pipeline) Process(n *NodeInfo) bool ***REMOVED***
	for i, entry := range p.checklist ***REMOVED***
		if entry.enabled && !entry.f.Check(n) ***REMOVED***
			// Immediately stop on first failure.
			p.checklist[i].failureCount++
			return false
		***REMOVED***
	***REMOVED***
	for i := range p.checklist ***REMOVED***
		p.checklist[i].failureCount = 0
	***REMOVED***
	return true
***REMOVED***

// SetTask sets up the filters to process a new task. Once this is called,
// Process can be called repeatedly to try to assign the task various nodes.
func (p *Pipeline) SetTask(t *api.Task) ***REMOVED***
	for i := range p.checklist ***REMOVED***
		p.checklist[i].enabled = p.checklist[i].f.SetTask(t)
		p.checklist[i].failureCount = 0
	***REMOVED***
***REMOVED***

// Explain returns a string explaining why a task could not be scheduled.
func (p *Pipeline) Explain() string ***REMOVED***
	var explanation string

	// Sort from most failures to least

	sortedByFailures := make([]checklistEntry, len(p.checklist))
	copy(sortedByFailures, p.checklist)
	sort.Sort(sort.Reverse(checklistByFailures(sortedByFailures)))

	for _, entry := range sortedByFailures ***REMOVED***
		if entry.failureCount > 0 ***REMOVED***
			if len(explanation) > 0 ***REMOVED***
				explanation += "; "
			***REMOVED***
			explanation += entry.f.Explain(entry.failureCount)
		***REMOVED***
	***REMOVED***

	return explanation
***REMOVED***
