package logger

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

const (
	defaultRingMaxSize = 1e6 // 1MB
)

// RingLogger is a ring buffer that implements the Logger interface.
// This is used when lossy logging is OK.
type RingLogger struct ***REMOVED***
	buffer    *messageRing
	l         Logger
	logInfo   Info
	closeFlag int32
***REMOVED***

type ringWithReader struct ***REMOVED***
	*RingLogger
***REMOVED***

func (r *ringWithReader) ReadLogs(cfg ReadConfig) *LogWatcher ***REMOVED***
	reader, ok := r.l.(LogReader)
	if !ok ***REMOVED***
		// something is wrong if we get here
		panic("expected log reader")
	***REMOVED***
	return reader.ReadLogs(cfg)
***REMOVED***

func newRingLogger(driver Logger, logInfo Info, maxSize int64) *RingLogger ***REMOVED***
	l := &RingLogger***REMOVED***
		buffer:  newRing(maxSize),
		l:       driver,
		logInfo: logInfo,
	***REMOVED***
	go l.run()
	return l
***REMOVED***

// NewRingLogger creates a new Logger that is implemented as a RingBuffer wrapping
// the passed in logger.
func NewRingLogger(driver Logger, logInfo Info, maxSize int64) Logger ***REMOVED***
	if maxSize < 0 ***REMOVED***
		maxSize = defaultRingMaxSize
	***REMOVED***
	l := newRingLogger(driver, logInfo, maxSize)
	if _, ok := driver.(LogReader); ok ***REMOVED***
		return &ringWithReader***REMOVED***l***REMOVED***
	***REMOVED***
	return l
***REMOVED***

// Log queues messages into the ring buffer
func (r *RingLogger) Log(msg *Message) error ***REMOVED***
	if r.closed() ***REMOVED***
		return errClosed
	***REMOVED***
	return r.buffer.Enqueue(msg)
***REMOVED***

// Name returns the name of the underlying logger
func (r *RingLogger) Name() string ***REMOVED***
	return r.l.Name()
***REMOVED***

func (r *RingLogger) closed() bool ***REMOVED***
	return atomic.LoadInt32(&r.closeFlag) == 1
***REMOVED***

func (r *RingLogger) setClosed() ***REMOVED***
	atomic.StoreInt32(&r.closeFlag, 1)
***REMOVED***

// Close closes the logger
func (r *RingLogger) Close() error ***REMOVED***
	r.setClosed()
	r.buffer.Close()
	// empty out the queue
	var logErr bool
	for _, msg := range r.buffer.Drain() ***REMOVED***
		if logErr ***REMOVED***
			// some error logging a previous message, so re-insert to message pool
			// and assume log driver is hosed
			PutMessage(msg)
			continue
		***REMOVED***

		if err := r.l.Log(msg); err != nil ***REMOVED***
			logrus.WithField("driver", r.l.Name()).WithField("container", r.logInfo.ContainerID).Errorf("Error writing log message: %v", r.l)
			logErr = true
		***REMOVED***
	***REMOVED***
	return r.l.Close()
***REMOVED***

// run consumes messages from the ring buffer and forwards them to the underling
// logger.
// This is run in a goroutine when the RingLogger is created
func (r *RingLogger) run() ***REMOVED***
	for ***REMOVED***
		if r.closed() ***REMOVED***
			return
		***REMOVED***
		msg, err := r.buffer.Dequeue()
		if err != nil ***REMOVED***
			// buffer is closed
			return
		***REMOVED***
		if err := r.l.Log(msg); err != nil ***REMOVED***
			logrus.WithField("driver", r.l.Name()).WithField("container", r.logInfo.ContainerID).Errorf("Error writing log message: %v", r.l)
		***REMOVED***
	***REMOVED***
***REMOVED***

type messageRing struct ***REMOVED***
	mu sync.Mutex
	// signals callers of `Dequeue` to wake up either on `Close` or when a new `Message` is added
	wait *sync.Cond

	sizeBytes int64 // current buffer size
	maxBytes  int64 // max buffer size size
	queue     []*Message
	closed    bool
***REMOVED***

func newRing(maxBytes int64) *messageRing ***REMOVED***
	queueSize := 1000
	if maxBytes == 0 || maxBytes == 1 ***REMOVED***
		// With 0 or 1 max byte size, the maximum size of the queue would only ever be 1
		// message long.
		queueSize = 1
	***REMOVED***

	r := &messageRing***REMOVED***queue: make([]*Message, 0, queueSize), maxBytes: maxBytes***REMOVED***
	r.wait = sync.NewCond(&r.mu)
	return r
***REMOVED***

// Enqueue adds a message to the buffer queue
// If the message is too big for the buffer it drops the oldest messages to make room
// If there are no messages in the queue and the message is still too big, it adds the message anyway.
func (r *messageRing) Enqueue(m *Message) error ***REMOVED***
	mSize := int64(len(m.Line))

	r.mu.Lock()
	if r.closed ***REMOVED***
		r.mu.Unlock()
		return errClosed
	***REMOVED***
	if mSize+r.sizeBytes > r.maxBytes && len(r.queue) > 0 ***REMOVED***
		r.wait.Signal()
		r.mu.Unlock()
		return nil
	***REMOVED***

	r.queue = append(r.queue, m)
	r.sizeBytes += mSize
	r.wait.Signal()
	r.mu.Unlock()
	return nil
***REMOVED***

// Dequeue pulls a message off the queue
// If there are no messages, it waits for one.
// If the buffer is closed, it will return immediately.
func (r *messageRing) Dequeue() (*Message, error) ***REMOVED***
	r.mu.Lock()
	for len(r.queue) == 0 && !r.closed ***REMOVED***
		r.wait.Wait()
	***REMOVED***

	if r.closed ***REMOVED***
		r.mu.Unlock()
		return nil, errClosed
	***REMOVED***

	msg := r.queue[0]
	r.queue = r.queue[1:]
	r.sizeBytes -= int64(len(msg.Line))
	r.mu.Unlock()
	return msg, nil
***REMOVED***

var errClosed = errors.New("closed")

// Close closes the buffer ensuring no new messages can be added.
// Any callers waiting to dequeue a message will be woken up.
func (r *messageRing) Close() ***REMOVED***
	r.mu.Lock()
	if r.closed ***REMOVED***
		r.mu.Unlock()
		return
	***REMOVED***

	r.closed = true
	r.wait.Broadcast()
	r.mu.Unlock()
***REMOVED***

// Drain drains all messages from the queue.
// This can be used after `Close()` to get any remaining messages that were in queue.
func (r *messageRing) Drain() []*Message ***REMOVED***
	r.mu.Lock()
	ls := make([]*Message, 0, len(r.queue))
	ls = append(ls, r.queue...)
	r.sizeBytes = 0
	r.queue = r.queue[:0]
	r.mu.Unlock()
	return ls
***REMOVED***
