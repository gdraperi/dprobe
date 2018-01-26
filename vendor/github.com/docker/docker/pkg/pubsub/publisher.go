package pubsub

import (
	"sync"
	"time"
)

var wgPool = sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(sync.WaitGroup) ***REMOVED******REMOVED***

// NewPublisher creates a new pub/sub publisher to broadcast messages.
// The duration is used as the send timeout as to not block the publisher publishing
// messages to other clients if one client is slow or unresponsive.
// The buffer is used when creating new channels for subscribers.
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher ***REMOVED***
	return &Publisher***REMOVED***
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	***REMOVED***
***REMOVED***

type subscriber chan interface***REMOVED******REMOVED***
type topicFunc func(v interface***REMOVED******REMOVED***) bool

// Publisher is basic pub/sub structure. Allows to send events and subscribe
// to them. Can be safely used from multiple goroutines.
type Publisher struct ***REMOVED***
	m           sync.RWMutex
	buffer      int
	timeout     time.Duration
	subscribers map[subscriber]topicFunc
***REMOVED***

// Len returns the number of subscribers for the publisher
func (p *Publisher) Len() int ***REMOVED***
	p.m.RLock()
	i := len(p.subscribers)
	p.m.RUnlock()
	return i
***REMOVED***

// Subscribe adds a new subscriber to the publisher returning the channel.
func (p *Publisher) Subscribe() chan interface***REMOVED******REMOVED*** ***REMOVED***
	return p.SubscribeTopic(nil)
***REMOVED***

// SubscribeTopic adds a new subscriber that filters messages sent by a topic.
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface***REMOVED******REMOVED*** ***REMOVED***
	ch := make(chan interface***REMOVED******REMOVED***, p.buffer)
	p.m.Lock()
	p.subscribers[ch] = topic
	p.m.Unlock()
	return ch
***REMOVED***

// SubscribeTopicWithBuffer adds a new subscriber that filters messages sent by a topic.
// The returned channel has a buffer of the specified size.
func (p *Publisher) SubscribeTopicWithBuffer(topic topicFunc, buffer int) chan interface***REMOVED******REMOVED*** ***REMOVED***
	ch := make(chan interface***REMOVED******REMOVED***, buffer)
	p.m.Lock()
	p.subscribers[ch] = topic
	p.m.Unlock()
	return ch
***REMOVED***

// Evict removes the specified subscriber from receiving any more messages.
func (p *Publisher) Evict(sub chan interface***REMOVED******REMOVED***) ***REMOVED***
	p.m.Lock()
	delete(p.subscribers, sub)
	close(sub)
	p.m.Unlock()
***REMOVED***

// Publish sends the data in v to all subscribers currently registered with the publisher.
func (p *Publisher) Publish(v interface***REMOVED******REMOVED***) ***REMOVED***
	p.m.RLock()
	if len(p.subscribers) == 0 ***REMOVED***
		p.m.RUnlock()
		return
	***REMOVED***

	wg := wgPool.Get().(*sync.WaitGroup)
	for sub, topic := range p.subscribers ***REMOVED***
		wg.Add(1)
		go p.sendTopic(sub, topic, v, wg)
	***REMOVED***
	wg.Wait()
	wgPool.Put(wg)
	p.m.RUnlock()
***REMOVED***

// Close closes the channels to all subscribers registered with the publisher.
func (p *Publisher) Close() ***REMOVED***
	p.m.Lock()
	for sub := range p.subscribers ***REMOVED***
		delete(p.subscribers, sub)
		close(sub)
	***REMOVED***
	p.m.Unlock()
***REMOVED***

func (p *Publisher) sendTopic(sub subscriber, topic topicFunc, v interface***REMOVED******REMOVED***, wg *sync.WaitGroup) ***REMOVED***
	defer wg.Done()
	if topic != nil && !topic(v) ***REMOVED***
		return
	***REMOVED***

	// send under a select as to not block if the receiver is unavailable
	if p.timeout > 0 ***REMOVED***
		select ***REMOVED***
		case sub <- v:
		case <-time.After(p.timeout):
		***REMOVED***
		return
	***REMOVED***

	select ***REMOVED***
	case sub <- v:
	default:
	***REMOVED***
***REMOVED***
