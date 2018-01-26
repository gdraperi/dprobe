package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func TestSendToOneSub(t *testing.T) ***REMOVED***
	p := NewPublisher(100*time.Millisecond, 10)
	c := p.Subscribe()

	p.Publish("hi")

	msg := <-c
	if msg.(string) != "hi" ***REMOVED***
		t.Fatalf("expected message hi but received %v", msg)
	***REMOVED***
***REMOVED***

func TestSendToMultipleSubs(t *testing.T) ***REMOVED***
	p := NewPublisher(100*time.Millisecond, 10)
	subs := []chan interface***REMOVED******REMOVED******REMOVED******REMOVED***
	subs = append(subs, p.Subscribe(), p.Subscribe(), p.Subscribe())

	p.Publish("hi")

	for _, c := range subs ***REMOVED***
		msg := <-c
		if msg.(string) != "hi" ***REMOVED***
			t.Fatalf("expected message hi but received %v", msg)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEvictOneSub(t *testing.T) ***REMOVED***
	p := NewPublisher(100*time.Millisecond, 10)
	s1 := p.Subscribe()
	s2 := p.Subscribe()

	p.Evict(s1)
	p.Publish("hi")
	if _, ok := <-s1; ok ***REMOVED***
		t.Fatal("expected s1 to not receive the published message")
	***REMOVED***

	msg := <-s2
	if msg.(string) != "hi" ***REMOVED***
		t.Fatalf("expected message hi but received %v", msg)
	***REMOVED***
***REMOVED***

func TestClosePublisher(t *testing.T) ***REMOVED***
	p := NewPublisher(100*time.Millisecond, 10)
	subs := []chan interface***REMOVED******REMOVED******REMOVED******REMOVED***
	subs = append(subs, p.Subscribe(), p.Subscribe(), p.Subscribe())
	p.Close()

	for _, c := range subs ***REMOVED***
		if _, ok := <-c; ok ***REMOVED***
			t.Fatal("expected all subscriber channels to be closed")
		***REMOVED***
	***REMOVED***
***REMOVED***

const sampleText = "test"

type testSubscriber struct ***REMOVED***
	dataCh chan interface***REMOVED******REMOVED***
	ch     chan error
***REMOVED***

func (s *testSubscriber) Wait() error ***REMOVED***
	return <-s.ch
***REMOVED***

func newTestSubscriber(p *Publisher) *testSubscriber ***REMOVED***
	ts := &testSubscriber***REMOVED***
		dataCh: p.Subscribe(),
		ch:     make(chan error),
	***REMOVED***
	go func() ***REMOVED***
		for data := range ts.dataCh ***REMOVED***
			s, ok := data.(string)
			if !ok ***REMOVED***
				ts.ch <- fmt.Errorf("Unexpected type %T", data)
				break
			***REMOVED***
			if s != sampleText ***REMOVED***
				ts.ch <- fmt.Errorf("Unexpected text %s", s)
				break
			***REMOVED***
		***REMOVED***
		close(ts.ch)
	***REMOVED***()
	return ts
***REMOVED***

// for testing with -race
func TestPubSubRace(t *testing.T) ***REMOVED***
	p := NewPublisher(0, 1024)
	var subs [](*testSubscriber)
	for j := 0; j < 50; j++ ***REMOVED***
		subs = append(subs, newTestSubscriber(p))
	***REMOVED***
	for j := 0; j < 1000; j++ ***REMOVED***
		p.Publish(sampleText)
	***REMOVED***
	time.AfterFunc(1*time.Second, func() ***REMOVED***
		for _, s := range subs ***REMOVED***
			p.Evict(s.dataCh)
		***REMOVED***
	***REMOVED***)
	for _, s := range subs ***REMOVED***
		s.Wait()
	***REMOVED***
***REMOVED***

func BenchmarkPubSub(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		b.StopTimer()
		p := NewPublisher(0, 1024)
		var subs [](*testSubscriber)
		for j := 0; j < 50; j++ ***REMOVED***
			subs = append(subs, newTestSubscriber(p))
		***REMOVED***
		b.StartTimer()
		for j := 0; j < 1000; j++ ***REMOVED***
			p.Publish(sampleText)
		***REMOVED***
		time.AfterFunc(1*time.Second, func() ***REMOVED***
			for _, s := range subs ***REMOVED***
				p.Evict(s.dataCh)
			***REMOVED***
		***REMOVED***)
		for _, s := range subs ***REMOVED***
			if err := s.Wait(); err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
