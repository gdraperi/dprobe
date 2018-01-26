package metrics

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	// We force flush the statsite metrics after this period of
	// inactivity. Prevents stats from getting stuck in a buffer
	// forever.
	flushInterval = 100 * time.Millisecond
)

// StatsiteSink provides a MetricSink that can be used with a
// statsite metrics server
type StatsiteSink struct ***REMOVED***
	addr        string
	metricQueue chan string
***REMOVED***

// NewStatsiteSink is used to create a new StatsiteSink
func NewStatsiteSink(addr string) (*StatsiteSink, error) ***REMOVED***
	s := &StatsiteSink***REMOVED***
		addr:        addr,
		metricQueue: make(chan string, 4096),
	***REMOVED***
	go s.flushMetrics()
	return s, nil
***REMOVED***

// Close is used to stop flushing to statsite
func (s *StatsiteSink) Shutdown() ***REMOVED***
	close(s.metricQueue)
***REMOVED***

func (s *StatsiteSink) SetGauge(key []string, val float32) ***REMOVED***
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|g\n", flatKey, val))
***REMOVED***

func (s *StatsiteSink) EmitKey(key []string, val float32) ***REMOVED***
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|kv\n", flatKey, val))
***REMOVED***

func (s *StatsiteSink) IncrCounter(key []string, val float32) ***REMOVED***
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|c\n", flatKey, val))
***REMOVED***

func (s *StatsiteSink) AddSample(key []string, val float32) ***REMOVED***
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|ms\n", flatKey, val))
***REMOVED***

// Flattens the key for formatting, removes spaces
func (s *StatsiteSink) flattenKey(parts []string) string ***REMOVED***
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
func (s *StatsiteSink) pushMetric(m string) ***REMOVED***
	select ***REMOVED***
	case s.metricQueue <- m:
	default:
	***REMOVED***
***REMOVED***

// Flushes metrics
func (s *StatsiteSink) flushMetrics() ***REMOVED***
	var sock net.Conn
	var err error
	var wait <-chan time.Time
	var buffered *bufio.Writer
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

CONNECT:
	// Attempt to connect
	sock, err = net.Dial("tcp", s.addr)
	if err != nil ***REMOVED***
		log.Printf("[ERR] Error connecting to statsite! Err: %s", err)
		goto WAIT
	***REMOVED***

	// Create a buffered writer
	buffered = bufio.NewWriter(sock)

	for ***REMOVED***
		select ***REMOVED***
		case metric, ok := <-s.metricQueue:
			// Get a metric from the queue
			if !ok ***REMOVED***
				goto QUIT
			***REMOVED***

			// Try to send to statsite
			_, err := buffered.Write([]byte(metric))
			if err != nil ***REMOVED***
				log.Printf("[ERR] Error writing to statsite! Err: %s", err)
				goto WAIT
			***REMOVED***
		case <-ticker.C:
			if err := buffered.Flush(); err != nil ***REMOVED***
				log.Printf("[ERR] Error flushing to statsite! Err: %s", err)
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
