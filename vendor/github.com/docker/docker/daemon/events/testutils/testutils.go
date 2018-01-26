package testutils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types/events"
	timetypes "github.com/docker/docker/api/types/time"
)

var (
	reTimestamp  = `(?P<timestamp>\d***REMOVED***4***REMOVED***-\d***REMOVED***2***REMOVED***-\d***REMOVED***2***REMOVED***T\d***REMOVED***2***REMOVED***:\d***REMOVED***2***REMOVED***:\d***REMOVED***2***REMOVED***.\d***REMOVED***9***REMOVED***(:?(:?(:?-|\+)\d***REMOVED***2***REMOVED***:\d***REMOVED***2***REMOVED***)|Z))`
	reEventType  = `(?P<eventType>\w+)`
	reAction     = `(?P<action>\w+)`
	reID         = `(?P<id>[^\s]+)`
	reAttributes = `(\s\((?P<attributes>[^\)]+)\))?`
	reString     = fmt.Sprintf(`\A%s\s%s\s%s\s%s%s\z`, reTimestamp, reEventType, reAction, reID, reAttributes)

	// eventCliRegexp is a regular expression that matches all possible event outputs in the cli
	eventCliRegexp = regexp.MustCompile(reString)
)

// ScanMap turns an event string like the default ones formatted in the cli output
// and turns it into map.
func ScanMap(text string) map[string]string ***REMOVED***
	matches := eventCliRegexp.FindAllStringSubmatch(text, -1)
	md := map[string]string***REMOVED******REMOVED***
	if len(matches) == 0 ***REMOVED***
		return md
	***REMOVED***

	names := eventCliRegexp.SubexpNames()
	for i, n := range matches[0] ***REMOVED***
		md[names[i]] = n
	***REMOVED***
	return md
***REMOVED***

// Scan turns an event string like the default ones formatted in the cli output
// and turns it into an event message.
func Scan(text string) (*events.Message, error) ***REMOVED***
	md := ScanMap(text)
	if len(md) == 0 ***REMOVED***
		return nil, fmt.Errorf("text is not an event: %s", text)
	***REMOVED***

	f, err := timetypes.GetTimestamp(md["timestamp"], time.Now())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	t, tn, err := timetypes.ParseTimestamps(f, -1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	attrs := make(map[string]string)
	for _, a := range strings.SplitN(md["attributes"], ", ", -1) ***REMOVED***
		kv := strings.SplitN(a, "=", 2)
		attrs[kv[0]] = kv[1]
	***REMOVED***

	tu := time.Unix(t, tn)
	return &events.Message***REMOVED***
		Time:     t,
		TimeNano: tu.UnixNano(),
		Type:     md["eventType"],
		Action:   md["action"],
		Actor: events.Actor***REMOVED***
			ID:         md["id"],
			Attributes: attrs,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***
