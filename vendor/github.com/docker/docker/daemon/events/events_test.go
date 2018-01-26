package events

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/events"
	timetypes "github.com/docker/docker/api/types/time"
	eventstestutils "github.com/docker/docker/daemon/events/testutils"
)

func TestEventsLog(t *testing.T) ***REMOVED***
	e := New()
	_, l1, _ := e.Subscribe()
	_, l2, _ := e.Subscribe()
	defer e.Evict(l1)
	defer e.Evict(l2)
	count := e.SubscribersCount()
	if count != 2 ***REMOVED***
		t.Fatalf("Must be 2 subscribers, got %d", count)
	***REMOVED***
	actor := events.Actor***REMOVED***
		ID:         "cont",
		Attributes: map[string]string***REMOVED***"image": "image"***REMOVED***,
	***REMOVED***
	e.Log("test", events.ContainerEventType, actor)
	select ***REMOVED***
	case msg := <-l1:
		jmsg, ok := msg.(events.Message)
		if !ok ***REMOVED***
			t.Fatalf("Unexpected type %T", msg)
		***REMOVED***
		if len(e.events) != 1 ***REMOVED***
			t.Fatalf("Must be only one event, got %d", len(e.events))
		***REMOVED***
		if jmsg.Status != "test" ***REMOVED***
			t.Fatalf("Status should be test, got %s", jmsg.Status)
		***REMOVED***
		if jmsg.ID != "cont" ***REMOVED***
			t.Fatalf("ID should be cont, got %s", jmsg.ID)
		***REMOVED***
		if jmsg.From != "image" ***REMOVED***
			t.Fatalf("From should be image, got %s", jmsg.From)
		***REMOVED***
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for broadcasted message")
	***REMOVED***
	select ***REMOVED***
	case msg := <-l2:
		jmsg, ok := msg.(events.Message)
		if !ok ***REMOVED***
			t.Fatalf("Unexpected type %T", msg)
		***REMOVED***
		if len(e.events) != 1 ***REMOVED***
			t.Fatalf("Must be only one event, got %d", len(e.events))
		***REMOVED***
		if jmsg.Status != "test" ***REMOVED***
			t.Fatalf("Status should be test, got %s", jmsg.Status)
		***REMOVED***
		if jmsg.ID != "cont" ***REMOVED***
			t.Fatalf("ID should be cont, got %s", jmsg.ID)
		***REMOVED***
		if jmsg.From != "image" ***REMOVED***
			t.Fatalf("From should be image, got %s", jmsg.From)
		***REMOVED***
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for broadcasted message")
	***REMOVED***
***REMOVED***

func TestEventsLogTimeout(t *testing.T) ***REMOVED***
	e := New()
	_, l, _ := e.Subscribe()
	defer e.Evict(l)

	c := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		actor := events.Actor***REMOVED***
			ID: "image",
		***REMOVED***
		e.Log("test", events.ImageEventType, actor)
		close(c)
	***REMOVED***()

	select ***REMOVED***
	case <-c:
	case <-time.After(time.Second):
		t.Fatal("Timeout publishing message")
	***REMOVED***
***REMOVED***

func TestLogEvents(t *testing.T) ***REMOVED***
	e := New()

	for i := 0; i < eventsLimit+16; i++ ***REMOVED***
		action := fmt.Sprintf("action_%d", i)
		id := fmt.Sprintf("cont_%d", i)
		from := fmt.Sprintf("image_%d", i)

		actor := events.Actor***REMOVED***
			ID:         id,
			Attributes: map[string]string***REMOVED***"image": from***REMOVED***,
		***REMOVED***
		e.Log(action, events.ContainerEventType, actor)
	***REMOVED***
	time.Sleep(50 * time.Millisecond)
	current, l, _ := e.Subscribe()
	for i := 0; i < 10; i++ ***REMOVED***
		num := i + eventsLimit + 16
		action := fmt.Sprintf("action_%d", num)
		id := fmt.Sprintf("cont_%d", num)
		from := fmt.Sprintf("image_%d", num)

		actor := events.Actor***REMOVED***
			ID:         id,
			Attributes: map[string]string***REMOVED***"image": from***REMOVED***,
		***REMOVED***
		e.Log(action, events.ContainerEventType, actor)
	***REMOVED***
	if len(e.events) != eventsLimit ***REMOVED***
		t.Fatalf("Must be %d events, got %d", eventsLimit, len(e.events))
	***REMOVED***

	var msgs []events.Message
	for len(msgs) < 10 ***REMOVED***
		m := <-l
		jm, ok := (m).(events.Message)
		if !ok ***REMOVED***
			t.Fatalf("Unexpected type %T", m)
		***REMOVED***
		msgs = append(msgs, jm)
	***REMOVED***
	if len(current) != eventsLimit ***REMOVED***
		t.Fatalf("Must be %d events, got %d", eventsLimit, len(current))
	***REMOVED***
	first := current[0]

	// TODO remove this once we removed the deprecated `ID`, `Status`, and `From` fields
	if first.Action != first.Status ***REMOVED***
		// Verify that the (deprecated) Status is set to the expected value
		t.Fatalf("Action (%s) does not match Status (%s)", first.Action, first.Status)
	***REMOVED***

	if first.Action != "action_16" ***REMOVED***
		t.Fatalf("First action is %s, must be action_16", first.Action)
	***REMOVED***
	last := current[len(current)-1]
	if last.Action != "action_271" ***REMOVED***
		t.Fatalf("Last action is %s, must be action_271", last.Action)
	***REMOVED***

	firstC := msgs[0]
	if firstC.Action != "action_272" ***REMOVED***
		t.Fatalf("First action is %s, must be action_272", firstC.Action)
	***REMOVED***
	lastC := msgs[len(msgs)-1]
	if lastC.Action != "action_281" ***REMOVED***
		t.Fatalf("Last action is %s, must be action_281", lastC.Action)
	***REMOVED***
***REMOVED***

// https://github.com/docker/docker/issues/20999
// Fixtures:
//
//2016-03-07T17:28:03.022433271+02:00 container die 0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079 (image=ubuntu, name=small_hoover)
//2016-03-07T17:28:03.091719377+02:00 network disconnect 19c5ed41acb798f26b751e0035cd7821741ab79e2bbd59a66b5fd8abf954eaa0 (type=bridge, container=0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079, name=bridge)
//2016-03-07T17:28:03.129014751+02:00 container destroy 0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079 (image=ubuntu, name=small_hoover)
func TestLoadBufferedEvents(t *testing.T) ***REMOVED***
	now := time.Now()
	f, err := timetypes.GetTimestamp("2016-03-07T17:28:03.100000000+02:00", now)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	s, sNano, err := timetypes.ParseTimestamps(f, -1)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	m1, err := eventstestutils.Scan("2016-03-07T17:28:03.022433271+02:00 container die 0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079 (image=ubuntu, name=small_hoover)")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	m2, err := eventstestutils.Scan("2016-03-07T17:28:03.091719377+02:00 network disconnect 19c5ed41acb798f26b751e0035cd7821741ab79e2bbd59a66b5fd8abf954eaa0 (type=bridge, container=0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079, name=bridge)")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	m3, err := eventstestutils.Scan("2016-03-07T17:28:03.129014751+02:00 container destroy 0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079 (image=ubuntu, name=small_hoover)")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	events := &Events***REMOVED***
		events: []events.Message***REMOVED****m1, *m2, *m3***REMOVED***,
	***REMOVED***

	since := time.Unix(s, sNano)
	until := time.Time***REMOVED******REMOVED***

	out := events.loadBufferedEvents(since, until, nil)
	if len(out) != 1 ***REMOVED***
		t.Fatalf("expected 1 message, got %d: %v", len(out), out)
	***REMOVED***
***REMOVED***

func TestLoadBufferedEventsOnlyFromPast(t *testing.T) ***REMOVED***
	now := time.Now()
	f, err := timetypes.GetTimestamp("2016-03-07T17:28:03.090000000+02:00", now)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	s, sNano, err := timetypes.ParseTimestamps(f, 0)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	f, err = timetypes.GetTimestamp("2016-03-07T17:28:03.100000000+02:00", now)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	u, uNano, err := timetypes.ParseTimestamps(f, 0)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	m1, err := eventstestutils.Scan("2016-03-07T17:28:03.022433271+02:00 container die 0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079 (image=ubuntu, name=small_hoover)")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	m2, err := eventstestutils.Scan("2016-03-07T17:28:03.091719377+02:00 network disconnect 19c5ed41acb798f26b751e0035cd7821741ab79e2bbd59a66b5fd8abf954eaa0 (type=bridge, container=0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079, name=bridge)")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	m3, err := eventstestutils.Scan("2016-03-07T17:28:03.129014751+02:00 container destroy 0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079 (image=ubuntu, name=small_hoover)")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	events := &Events***REMOVED***
		events: []events.Message***REMOVED****m1, *m2, *m3***REMOVED***,
	***REMOVED***

	since := time.Unix(s, sNano)
	until := time.Unix(u, uNano)

	out := events.loadBufferedEvents(since, until, nil)
	if len(out) != 1 ***REMOVED***
		t.Fatalf("expected 1 message, got %d: %v", len(out), out)
	***REMOVED***

	if out[0].Type != "network" ***REMOVED***
		t.Fatalf("expected network event, got %s", out[0].Type)
	***REMOVED***
***REMOVED***

// #13753
func TestIgnoreBufferedWhenNoTimes(t *testing.T) ***REMOVED***
	m1, err := eventstestutils.Scan("2016-03-07T17:28:03.022433271+02:00 container die 0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079 (image=ubuntu, name=small_hoover)")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	m2, err := eventstestutils.Scan("2016-03-07T17:28:03.091719377+02:00 network disconnect 19c5ed41acb798f26b751e0035cd7821741ab79e2bbd59a66b5fd8abf954eaa0 (type=bridge, container=0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079, name=bridge)")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	m3, err := eventstestutils.Scan("2016-03-07T17:28:03.129014751+02:00 container destroy 0b863f2a26c18557fc6cdadda007c459f9ec81b874780808138aea78a3595079 (image=ubuntu, name=small_hoover)")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	events := &Events***REMOVED***
		events: []events.Message***REMOVED****m1, *m2, *m3***REMOVED***,
	***REMOVED***

	since := time.Time***REMOVED******REMOVED***
	until := time.Time***REMOVED******REMOVED***

	out := events.loadBufferedEvents(since, until, nil)
	if len(out) != 0 ***REMOVED***
		t.Fatalf("expected 0 buffered events, got %q", out)
	***REMOVED***
***REMOVED***
