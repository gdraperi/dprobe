package container

import (
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
)

// Health holds the current container health-check state
type Health struct ***REMOVED***
	types.Health
	stop chan struct***REMOVED******REMOVED*** // Write struct***REMOVED******REMOVED*** to stop the monitor
	mu   sync.Mutex
***REMOVED***

// String returns a human-readable description of the health-check state
func (s *Health) String() string ***REMOVED***
	status := s.Status()

	switch status ***REMOVED***
	case types.Starting:
		return "health: starting"
	default: // Healthy and Unhealthy are clear on their own
		return s.Health.Status
	***REMOVED***
***REMOVED***

// Status returns the current health status.
//
// Note that this takes a lock and the value may change after being read.
func (s *Health) Status() string ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	// This happens when the monitor has yet to be setup.
	if s.Health.Status == "" ***REMOVED***
		return types.Unhealthy
	***REMOVED***

	return s.Health.Status
***REMOVED***

// SetStatus writes the current status to the underlying health structure,
// obeying the locking semantics.
//
// Status may be set directly if another lock is used.
func (s *Health) SetStatus(new string) ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Health.Status = new
***REMOVED***

// OpenMonitorChannel creates and returns a new monitor channel. If there
// already is one, it returns nil.
func (s *Health) OpenMonitorChannel() chan struct***REMOVED******REMOVED*** ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stop == nil ***REMOVED***
		logrus.Debug("OpenMonitorChannel")
		s.stop = make(chan struct***REMOVED******REMOVED***)
		return s.stop
	***REMOVED***
	return nil
***REMOVED***

// CloseMonitorChannel closes any existing monitor channel.
func (s *Health) CloseMonitorChannel() ***REMOVED***
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stop != nil ***REMOVED***
		logrus.Debug("CloseMonitorChannel: waiting for probe to stop")
		close(s.stop)
		s.stop = nil
		// unhealthy when the monitor has stopped for compatibility reasons
		s.Health.Status = types.Unhealthy
		logrus.Debug("CloseMonitorChannel done")
	***REMOVED***
***REMOVED***
