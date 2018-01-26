package metrics

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	// statsdMaxLen is the maximum size of a packet
	// to send to statsd
	statsdMaxLen = 1400
)

// StatsdSink provides a MetricSink that can be used
// with a statsite or statsd metrics server. It uses
// only UDP packets, while StatsiteSink uses TCP.
type StatsdSink struct ***REMOVED***
	addr        string
	metricQueue chan string
***REMOVED***

// NewStatsdSink is used to create a new StatsdSink
func NewStatsdSink(addr string) (*StatsdSink, error) ***REMOVED***
	s := &StatsdSink***REMOVED***
		addr:        addr,
		metricQueue: make(chan string, 4096),
	***REMOVED***
	go s.flushMetrics()
	return s, nil
***REMOVED***

// Close is used to stop flushing to statsd
func (s *StatsdSink) Shutdown() ***REMOVED***
	close(s.metricQueue)
***REMOVED***

func (s *StatsdSink) SetGauge(key []string, val float32) ***REMOVED***
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|g\n", flatKey, val))
***REMOVED***

func (s *StatsdSink) EmitKey(key []string, val float32) ***REMOVED***
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|kv\n", flatKey, val))
***REMOVED***

func (s *StatsdSink) IncrCounter(key []string, val float32) ***REMOVED***
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|c\n", flatKey, val))
***REMOVED***

func (s *StatsdSink) AddSample(key []string, val float32) ***REMOVED***
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|ms\n", flatKey, val))
***REMOVED***

// Flattens the key for formatting, removes spaces
func (s *StatsdSink) flattenKey(parts []string) string ***REMOVED***
	joined := strings.Join(parts, ".")
	return strings.Map(func(r rune) rune ***REMOVED***
		switch r ***REMOVED***
		case ':':
			fallthrough
		case ' ':
			return '_'
		default:
			return r
		***REMOVED***
	***REMOVED***, joined)
***REMOVED***

// Does a non-blocking push to the metrics queue
func (s *StatsdSink) pushMetric(m string) ***REMOVED***
	select ***REMOVED***
	case s.metricQueue <- m:
	default:
	***REMOVED***
***REMOVED***

// Flushes metrics
func (s *StatsdSink) flushMetrics() ***REMOVED***
	var sock net.Conn
	var err error
	var wait <-chan time.Time
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

CONNECT:
	// Create a buffer
	buf := bytes.NewBuffer(nil)

	// Attempt to connect
	sock, err = net.Dial("udp", s.addr)
	if err != nil ***REMOVED***
		log.Printf("[ERR] Error connecting to statsd! Err: %s", err)
		goto WAIT
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case metric, ok := <-s.metricQueue:
			// Get a metric from the queue
			if !ok ***REMOVED***
				goto QUIT
			***REMOVED***

			// Check if this would overflow the packet size
			if len(metric)+buf.Len() > statsdMaxLen ***REMOVED***
				_, err := sock.Write(buf.Bytes())
				buf.Reset()
				if err != nil ***REMOVED***
					log.Printf("[ERR] Error writing to statsd! Err: %s", err)
					goto WAIT
				***REMOVED***
			***REMOVED***

			// Append to the buffer
			buf.WriteString(metric)

		case <-ticker.C:
			if buf.Len() == 0 ***REMOVED***
				continue
			***REMOVED***

			_, err := sock.Write(buf.Bytes())
			buf.Reset()
			if err != nil ***REMOVED***
				log.Printf("[ERR] Error flushing to statsd! Err: %s", err)
				goto WAIT
			***REMOVED***
		***REMOVED***
	***REMOVED***

WAIT:
	// Wait for a while
	wait = time.After(time.Duration(5) * time.Second)
	for ***REMOVED***
		select ***REMOVED***
		// Dequeue the messages to avoid backlog
		case _, ok := <-s.metricQueue:
			if !ok ***REMOVED***
				goto QUIT
			***REMOVED***
		case <-wait:
			goto CONNECT
		***REMOVED***
	***REMOVED***
QUIT:
	s.metricQueue = nil
***REMOVED***
