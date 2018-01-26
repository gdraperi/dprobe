package serf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/serf/coordinate"
)

/*
Serf supports using a "snapshot" file that contains various
transactional data that is used to help Serf recover quickly
and gracefully from a failure. We append member events, as well
as the latest clock values to the file during normal operation,
and periodically checkpoint and roll over the file. During a restore,
we can replay the various member events to recall a list of known
nodes to re-join, as well as restore our clock values to avoid replaying
old events.
*/

const flushInterval = 500 * time.Millisecond
const clockUpdateInterval = 500 * time.Millisecond
const coordinateUpdateInterval = 60 * time.Second
const tmpExt = ".compact"

// Snapshotter is responsible for ingesting events and persisting
// them to disk, and providing a recovery mechanism at start time.
type Snapshotter struct ***REMOVED***
	aliveNodes       map[string]string
	clock            *LamportClock
	coordClient      *coordinate.Client
	fh               *os.File
	buffered         *bufio.Writer
	inCh             <-chan Event
	lastFlush        time.Time
	lastClock        LamportTime
	lastEventClock   LamportTime
	lastQueryClock   LamportTime
	leaveCh          chan struct***REMOVED******REMOVED***
	leaving          bool
	logger           *log.Logger
	maxSize          int64
	path             string
	offset           int64
	outCh            chan<- Event
	rejoinAfterLeave bool
	shutdownCh       <-chan struct***REMOVED******REMOVED***
	waitCh           chan struct***REMOVED******REMOVED***
***REMOVED***

// PreviousNode is used to represent the previously known alive nodes
type PreviousNode struct ***REMOVED***
	Name string
	Addr string
***REMOVED***

func (p PreviousNode) String() string ***REMOVED***
	return fmt.Sprintf("%s: %s", p.Name, p.Addr)
***REMOVED***

// NewSnapshotter creates a new Snapshotter that records events up to a
// max byte size before rotating the file. It can also be used to
// recover old state. Snapshotter works by reading an event channel it returns,
// passing through to an output channel, and persisting relevant events to disk.
// Setting rejoinAfterLeave makes leave not clear the state, and can be used
// if you intend to rejoin the same cluster after a leave.
func NewSnapshotter(path string,
	maxSize int,
	rejoinAfterLeave bool,
	logger *log.Logger,
	clock *LamportClock,
	coordClient *coordinate.Client,
	outCh chan<- Event,
	shutdownCh <-chan struct***REMOVED******REMOVED***) (chan<- Event, *Snapshotter, error) ***REMOVED***
	inCh := make(chan Event, 1024)

	// Try to open the file
	fh, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("failed to open snapshot: %v", err)
	***REMOVED***

	// Determine the offset
	info, err := fh.Stat()
	if err != nil ***REMOVED***
		fh.Close()
		return nil, nil, fmt.Errorf("failed to stat snapshot: %v", err)
	***REMOVED***
	offset := info.Size()

	// Create the snapshotter
	snap := &Snapshotter***REMOVED***
		aliveNodes:       make(map[string]string),
		clock:            clock,
		coordClient:      coordClient,
		fh:               fh,
		buffered:         bufio.NewWriter(fh),
		inCh:             inCh,
		lastClock:        0,
		lastEventClock:   0,
		lastQueryClock:   0,
		leaveCh:          make(chan struct***REMOVED******REMOVED***),
		logger:           logger,
		maxSize:          int64(maxSize),
		path:             path,
		offset:           offset,
		outCh:            outCh,
		rejoinAfterLeave: rejoinAfterLeave,
		shutdownCh:       shutdownCh,
		waitCh:           make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	// Recover the last known state
	if err := snap.replay(); err != nil ***REMOVED***
		fh.Close()
		return nil, nil, err
	***REMOVED***

	// Start handling new commands
	go snap.stream()
	return inCh, snap, nil
***REMOVED***

// LastClock returns the last known clock time
func (s *Snapshotter) LastClock() LamportTime ***REMOVED***
	return s.lastClock
***REMOVED***

// LastEventClock returns the last known event clock time
func (s *Snapshotter) LastEventClock() LamportTime ***REMOVED***
	return s.lastEventClock
***REMOVED***

// LastQueryClock returns the last known query clock time
func (s *Snapshotter) LastQueryClock() LamportTime ***REMOVED***
	return s.lastQueryClock
***REMOVED***

// AliveNodes returns the last known alive nodes
func (s *Snapshotter) AliveNodes() []*PreviousNode ***REMOVED***
	// Copy the previously known
	previous := make([]*PreviousNode, 0, len(s.aliveNodes))
	for name, addr := range s.aliveNodes ***REMOVED***
		previous = append(previous, &PreviousNode***REMOVED***name, addr***REMOVED***)
	***REMOVED***

	// Randomize the order, prevents hot shards
	for i := range previous ***REMOVED***
		j := rand.Intn(i + 1)
		previous[i], previous[j] = previous[j], previous[i]
	***REMOVED***
	return previous
***REMOVED***

// Wait is used to wait until the snapshotter finishes shut down
func (s *Snapshotter) Wait() ***REMOVED***
	<-s.waitCh
***REMOVED***

// Leave is used to remove known nodes to prevent a restart from
// causing a join. Otherwise nodes will re-join after leaving!
func (s *Snapshotter) Leave() ***REMOVED***
	select ***REMOVED***
	case s.leaveCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	case <-s.shutdownCh:
	***REMOVED***
***REMOVED***

// stream is a long running routine that is used to handle events
func (s *Snapshotter) stream() ***REMOVED***
	clockTicker := time.NewTicker(clockUpdateInterval)
	defer clockTicker.Stop()

	coordinateTicker := time.NewTicker(coordinateUpdateInterval)
	defer coordinateTicker.Stop()

	for ***REMOVED***
		select ***REMOVED***
		case <-s.leaveCh:
			s.leaving = true

			// If we plan to re-join, keep our state
			if !s.rejoinAfterLeave ***REMOVED***
				s.aliveNodes = make(map[string]string)
			***REMOVED***
			s.tryAppend("leave\n")
			if err := s.buffered.Flush(); err != nil ***REMOVED***
				s.logger.Printf("[ERR] serf: failed to flush leave to snapshot: %v", err)
			***REMOVED***
			if err := s.fh.Sync(); err != nil ***REMOVED***
				s.logger.Printf("[ERR] serf: failed to sync leave to snapshot: %v", err)
			***REMOVED***

		case e := <-s.inCh:
			// Forward the event immediately
			if s.outCh != nil ***REMOVED***
				s.outCh <- e
			***REMOVED***

			// Stop recording events after a leave is issued
			if s.leaving ***REMOVED***
				continue
			***REMOVED***
			switch typed := e.(type) ***REMOVED***
			case MemberEvent:
				s.processMemberEvent(typed)
			case UserEvent:
				s.processUserEvent(typed)
			case *Query:
				s.processQuery(typed)
			default:
				s.logger.Printf("[ERR] serf: Unknown event to snapshot: %#v", e)
			***REMOVED***

		case <-clockTicker.C:
			s.updateClock()

		case <-coordinateTicker.C:
			s.updateCoordinate()

		case <-s.shutdownCh:
			if err := s.buffered.Flush(); err != nil ***REMOVED***
				s.logger.Printf("[ERR] serf: failed to flush snapshot: %v", err)
			***REMOVED***
			if err := s.fh.Sync(); err != nil ***REMOVED***
				s.logger.Printf("[ERR] serf: failed to sync snapshot: %v", err)
			***REMOVED***
			s.fh.Close()
			close(s.waitCh)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// processMemberEvent is used to handle a single member event
func (s *Snapshotter) processMemberEvent(e MemberEvent) ***REMOVED***
	switch e.Type ***REMOVED***
	case EventMemberJoin:
		for _, mem := range e.Members ***REMOVED***
			addr := net.TCPAddr***REMOVED***IP: mem.Addr, Port: int(mem.Port)***REMOVED***
			s.aliveNodes[mem.Name] = addr.String()
			s.tryAppend(fmt.Sprintf("alive: %s %s\n", mem.Name, addr.String()))
		***REMOVED***

	case EventMemberLeave:
		fallthrough
	case EventMemberFailed:
		for _, mem := range e.Members ***REMOVED***
			delete(s.aliveNodes, mem.Name)
			s.tryAppend(fmt.Sprintf("not-alive: %s\n", mem.Name))
		***REMOVED***
	***REMOVED***
	s.updateClock()
***REMOVED***

// updateClock is called periodically to check if we should udpate our
// clock value. This is done after member events but should also be done
// periodically due to race conditions with join and leave intents
func (s *Snapshotter) updateClock() ***REMOVED***
	lastSeen := s.clock.Time() - 1
	if lastSeen > s.lastClock ***REMOVED***
		s.lastClock = lastSeen
		s.tryAppend(fmt.Sprintf("clock: %d\n", s.lastClock))
	***REMOVED***
***REMOVED***

// updateCoordinate is called periodically to write out the current local
// coordinate. It's safe to call this if coordinates aren't enabled (nil
// client) and it will be a no-op.
func (s *Snapshotter) updateCoordinate() ***REMOVED***
	if s.coordClient != nil ***REMOVED***
		encoded, err := json.Marshal(s.coordClient.GetCoordinate())
		if err != nil ***REMOVED***
			s.logger.Printf("[ERR] serf: Failed to encode coordinate: %v", err)
		***REMOVED*** else ***REMOVED***
			s.tryAppend(fmt.Sprintf("coordinate: %s\n", encoded))
		***REMOVED***
	***REMOVED***
***REMOVED***

// processUserEvent is used to handle a single user event
func (s *Snapshotter) processUserEvent(e UserEvent) ***REMOVED***
	// Ignore old clocks
	if e.LTime <= s.lastEventClock ***REMOVED***
		return
	***REMOVED***
	s.lastEventClock = e.LTime
	s.tryAppend(fmt.Sprintf("event-clock: %d\n", e.LTime))
***REMOVED***

// processQuery is used to handle a single query event
func (s *Snapshotter) processQuery(q *Query) ***REMOVED***
	// Ignore old clocks
	if q.LTime <= s.lastQueryClock ***REMOVED***
		return
	***REMOVED***
	s.lastQueryClock = q.LTime
	s.tryAppend(fmt.Sprintf("query-clock: %d\n", q.LTime))
***REMOVED***

// tryAppend will invoke append line but will not return an error
func (s *Snapshotter) tryAppend(l string) ***REMOVED***
	if err := s.appendLine(l); err != nil ***REMOVED***
		s.logger.Printf("[ERR] serf: Failed to update snapshot: %v", err)
	***REMOVED***
***REMOVED***

// appendLine is used to append a line to the existing log
func (s *Snapshotter) appendLine(l string) error ***REMOVED***
	defer metrics.MeasureSince([]string***REMOVED***"serf", "snapshot", "appendLine"***REMOVED***, time.Now())

	n, err := s.buffered.WriteString(l)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Check if we should flush
	now := time.Now()
	if now.Sub(s.lastFlush) > flushInterval ***REMOVED***
		s.lastFlush = now
		if err := s.buffered.Flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Check if a compaction is necessary
	s.offset += int64(n)
	if s.offset > s.maxSize ***REMOVED***
		return s.compact()
	***REMOVED***
	return nil
***REMOVED***

// Compact is used to compact the snapshot once it is too large
func (s *Snapshotter) compact() error ***REMOVED***
	defer metrics.MeasureSince([]string***REMOVED***"serf", "snapshot", "compact"***REMOVED***, time.Now())

	// Try to open the file to new fiel
	newPath := s.path + tmpExt
	fh, err := os.OpenFile(newPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to open new snapshot: %v", err)
	***REMOVED***

	// Create a buffered writer
	buf := bufio.NewWriter(fh)

	// Write out the live nodes
	var offset int64
	for name, addr := range s.aliveNodes ***REMOVED***
		line := fmt.Sprintf("alive: %s %s\n", name, addr)
		n, err := buf.WriteString(line)
		if err != nil ***REMOVED***
			fh.Close()
			return err
		***REMOVED***
		offset += int64(n)
	***REMOVED***

	// Write out the clocks
	line := fmt.Sprintf("clock: %d\n", s.lastClock)
	n, err := buf.WriteString(line)
	if err != nil ***REMOVED***
		fh.Close()
		return err
	***REMOVED***
	offset += int64(n)

	line = fmt.Sprintf("event-clock: %d\n", s.lastEventClock)
	n, err = buf.WriteString(line)
	if err != nil ***REMOVED***
		fh.Close()
		return err
	***REMOVED***
	offset += int64(n)

	line = fmt.Sprintf("query-clock: %d\n", s.lastQueryClock)
	n, err = buf.WriteString(line)
	if err != nil ***REMOVED***
		fh.Close()
		return err
	***REMOVED***
	offset += int64(n)

	// Write out the coordinate.
	if s.coordClient != nil ***REMOVED***
		encoded, err := json.Marshal(s.coordClient.GetCoordinate())
		if err != nil ***REMOVED***
			fh.Close()
			return err
		***REMOVED***

		line = fmt.Sprintf("coordinate: %s\n", encoded)
		n, err = buf.WriteString(line)
		if err != nil ***REMOVED***
			fh.Close()
			return err
		***REMOVED***
		offset += int64(n)
	***REMOVED***

	// Flush the new snapshot
	err = buf.Flush()
	fh.Close()
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to flush new snapshot: %v", err)
	***REMOVED***

	// We now need to swap the old snapshot file with the new snapshot.
	// Turns out, Windows won't let us rename the files if we have
	// open handles to them or if the destination already exists. This
	// means we are forced to close the existing handles, delete the
	// old file, move the new one in place, and then re-open the file
	// handles.

	// Flush the existing snapshot, ignoring errors since we will
	// delete it momentarily.
	s.buffered.Flush()
	s.buffered = nil

	// Close the file handle to the old snapshot
	s.fh.Close()
	s.fh = nil

	// Delete the old file
	if err := os.Remove(s.path); err != nil ***REMOVED***
		return fmt.Errorf("failed to remove old snapshot: %v", err)
	***REMOVED***

	// Move the new file into place
	if err := os.Rename(newPath, s.path); err != nil ***REMOVED***
		return fmt.Errorf("failed to install new snapshot: %v", err)
	***REMOVED***

	// Open the new snapshot
	fh, err = os.OpenFile(s.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to open snapshot: %v", err)
	***REMOVED***
	buf = bufio.NewWriter(fh)

	// Rotate our handles
	s.fh = fh
	s.buffered = buf
	s.offset = offset
	s.lastFlush = time.Now()
	return nil
***REMOVED***

// replay is used to seek to reset our internal state by replaying
// the snapshot file. It is used at initialization time to read old
// state
func (s *Snapshotter) replay() error ***REMOVED***
	// Seek to the beginning
	if _, err := s.fh.Seek(0, os.SEEK_SET); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Read each line
	reader := bufio.NewReader(s.fh)
	for ***REMOVED***
		line, err := reader.ReadString('\n')
		if err != nil ***REMOVED***
			break
		***REMOVED***

		// Skip the newline
		line = line[:len(line)-1]

		// Switch on the prefix
		if strings.HasPrefix(line, "alive: ") ***REMOVED***
			info := strings.TrimPrefix(line, "alive: ")
			addrIdx := strings.LastIndex(info, " ")
			if addrIdx == -1 ***REMOVED***
				s.logger.Printf("[WARN] serf: Failed to parse address: %v", line)
				continue
			***REMOVED***
			addr := info[addrIdx+1:]
			name := info[:addrIdx]
			s.aliveNodes[name] = addr

		***REMOVED*** else if strings.HasPrefix(line, "not-alive: ") ***REMOVED***
			name := strings.TrimPrefix(line, "not-alive: ")
			delete(s.aliveNodes, name)

		***REMOVED*** else if strings.HasPrefix(line, "clock: ") ***REMOVED***
			timeStr := strings.TrimPrefix(line, "clock: ")
			timeInt, err := strconv.ParseUint(timeStr, 10, 64)
			if err != nil ***REMOVED***
				s.logger.Printf("[WARN] serf: Failed to convert clock time: %v", err)
				continue
			***REMOVED***
			s.lastClock = LamportTime(timeInt)

		***REMOVED*** else if strings.HasPrefix(line, "event-clock: ") ***REMOVED***
			timeStr := strings.TrimPrefix(line, "event-clock: ")
			timeInt, err := strconv.ParseUint(timeStr, 10, 64)
			if err != nil ***REMOVED***
				s.logger.Printf("[WARN] serf: Failed to convert event clock time: %v", err)
				continue
			***REMOVED***
			s.lastEventClock = LamportTime(timeInt)

		***REMOVED*** else if strings.HasPrefix(line, "query-clock: ") ***REMOVED***
			timeStr := strings.TrimPrefix(line, "query-clock: ")
			timeInt, err := strconv.ParseUint(timeStr, 10, 64)
			if err != nil ***REMOVED***
				s.logger.Printf("[WARN] serf: Failed to convert query clock time: %v", err)
				continue
			***REMOVED***
			s.lastQueryClock = LamportTime(timeInt)

		***REMOVED*** else if strings.HasPrefix(line, "coordinate: ") ***REMOVED***
			if s.coordClient == nil ***REMOVED***
				s.logger.Printf("[WARN] serf: Ignoring snapshot coordinates since they are disabled")
				continue
			***REMOVED***

			coordStr := strings.TrimPrefix(line, "coordinate: ")
			var coord coordinate.Coordinate
			err := json.Unmarshal([]byte(coordStr), &coord)
			if err != nil ***REMOVED***
				s.logger.Printf("[WARN] serf: Failed to decode coordinate: %v", err)
				continue
			***REMOVED***
			s.coordClient.SetCoordinate(&coord)
		***REMOVED*** else if line == "leave" ***REMOVED***
			// Ignore a leave if we plan on re-joining
			if s.rejoinAfterLeave ***REMOVED***
				s.logger.Printf("[INFO] serf: Ignoring previous leave in snapshot")
				continue
			***REMOVED***
			s.aliveNodes = make(map[string]string)
			s.lastClock = 0
			s.lastEventClock = 0
			s.lastQueryClock = 0

		***REMOVED*** else if strings.HasPrefix(line, "#") ***REMOVED***
			// Skip comment lines

		***REMOVED*** else ***REMOVED***
			s.logger.Printf("[WARN] serf: Unrecognized snapshot line: %v", line)
		***REMOVED***
	***REMOVED***

	// Seek to the end
	if _, err := s.fh.Seek(0, os.SEEK_END); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
