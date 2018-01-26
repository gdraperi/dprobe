package events

import (
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

// Filter can filter out docker events from a stream
type Filter struct ***REMOVED***
	filter filters.Args
***REMOVED***

// NewFilter creates a new Filter
func NewFilter(filter filters.Args) *Filter ***REMOVED***
	return &Filter***REMOVED***filter: filter***REMOVED***
***REMOVED***

// Include returns true when the event ev is included by the filters
func (ef *Filter) Include(ev events.Message) bool ***REMOVED***
	return ef.matchEvent(ev) &&
		ef.filter.ExactMatch("type", ev.Type) &&
		ef.matchScope(ev.Scope) &&
		ef.matchDaemon(ev) &&
		ef.matchContainer(ev) &&
		ef.matchPlugin(ev) &&
		ef.matchVolume(ev) &&
		ef.matchNetwork(ev) &&
		ef.matchImage(ev) &&
		ef.matchNode(ev) &&
		ef.matchService(ev) &&
		ef.matchSecret(ev) &&
		ef.matchConfig(ev) &&
		ef.matchLabels(ev.Actor.Attributes)
***REMOVED***

func (ef *Filter) matchEvent(ev events.Message) bool ***REMOVED***
	// #25798 if an event filter contains either health_status, exec_create or exec_start without a colon
	// Let's to a FuzzyMatch instead of an ExactMatch.
	if ef.filterContains("event", map[string]struct***REMOVED******REMOVED******REMOVED***"health_status": ***REMOVED******REMOVED***, "exec_create": ***REMOVED******REMOVED***, "exec_start": ***REMOVED******REMOVED******REMOVED***) ***REMOVED***
		return ef.filter.FuzzyMatch("event", ev.Action)
	***REMOVED***
	return ef.filter.ExactMatch("event", ev.Action)
***REMOVED***

func (ef *Filter) filterContains(field string, values map[string]struct***REMOVED******REMOVED***) bool ***REMOVED***
	for _, v := range ef.filter.Get(field) ***REMOVED***
		if _, ok := values[v]; ok ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (ef *Filter) matchScope(scope string) bool ***REMOVED***
	if !ef.filter.Contains("scope") ***REMOVED***
		return true
	***REMOVED***
	return ef.filter.ExactMatch("scope", scope)
***REMOVED***

func (ef *Filter) matchLabels(attributes map[string]string) bool ***REMOVED***
	if !ef.filter.Contains("label") ***REMOVED***
		return true
	***REMOVED***
	return ef.filter.MatchKVList("label", attributes)
***REMOVED***

func (ef *Filter) matchDaemon(ev events.Message) bool ***REMOVED***
	return ef.fuzzyMatchName(ev, events.DaemonEventType)
***REMOVED***

func (ef *Filter) matchContainer(ev events.Message) bool ***REMOVED***
	return ef.fuzzyMatchName(ev, events.ContainerEventType)
***REMOVED***

func (ef *Filter) matchPlugin(ev events.Message) bool ***REMOVED***
	return ef.fuzzyMatchName(ev, events.PluginEventType)
***REMOVED***

func (ef *Filter) matchVolume(ev events.Message) bool ***REMOVED***
	return ef.fuzzyMatchName(ev, events.VolumeEventType)
***REMOVED***

func (ef *Filter) matchNetwork(ev events.Message) bool ***REMOVED***
	return ef.fuzzyMatchName(ev, events.NetworkEventType)
***REMOVED***

func (ef *Filter) matchService(ev events.Message) bool ***REMOVED***
	return ef.fuzzyMatchName(ev, events.ServiceEventType)
***REMOVED***

func (ef *Filter) matchNode(ev events.Message) bool ***REMOVED***
	return ef.fuzzyMatchName(ev, events.NodeEventType)
***REMOVED***

func (ef *Filter) matchSecret(ev events.Message) bool ***REMOVED***
	return ef.fuzzyMatchName(ev, events.SecretEventType)
***REMOVED***

func (ef *Filter) matchConfig(ev events.Message) bool ***REMOVED***
	return ef.fuzzyMatchName(ev, events.ConfigEventType)
***REMOVED***

func (ef *Filter) fuzzyMatchName(ev events.Message, eventType string) bool ***REMOVED***
	return ef.filter.FuzzyMatch(eventType, ev.Actor.ID) ||
		ef.filter.FuzzyMatch(eventType, ev.Actor.Attributes["name"])
***REMOVED***

// matchImage matches against both event.Actor.ID (for image events)
// and event.Actor.Attributes["image"] (for container events), so that any container that was created
// from an image will be included in the image events. Also compare both
// against the stripped repo name without any tags.
func (ef *Filter) matchImage(ev events.Message) bool ***REMOVED***
	id := ev.Actor.ID
	nameAttr := "image"
	var imageName string

	if ev.Type == events.ImageEventType ***REMOVED***
		nameAttr = "name"
	***REMOVED***

	if n, ok := ev.Actor.Attributes[nameAttr]; ok ***REMOVED***
		imageName = n
	***REMOVED***
	return ef.filter.ExactMatch("image", id) ||
		ef.filter.ExactMatch("image", imageName) ||
		ef.filter.ExactMatch("image", stripTag(id)) ||
		ef.filter.ExactMatch("image", stripTag(imageName))
***REMOVED***

func stripTag(image string) string ***REMOVED***
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil ***REMOVED***
		return image
	***REMOVED***
	return reference.FamiliarName(ref)
***REMOVED***
