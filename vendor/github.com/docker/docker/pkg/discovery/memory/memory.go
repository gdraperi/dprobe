package memory

import (
	"sync"
	"time"

	"github.com/docker/docker/pkg/discovery"
)

// Discovery implements a discovery backend that keeps
// data in memory.
type Discovery struct ***REMOVED***
	heartbeat time.Duration
	values    []string
	mu        sync.Mutex
***REMOVED***

func init() ***REMOVED***
	Init()
***REMOVED***

// Init registers the memory backend on demand.
func Init() ***REMOVED***
	discovery.Register("memory", &Discovery***REMOVED******REMOVED***)
***REMOVED***

// Initialize sets the heartbeat for the memory backend.
func (s *Discovery) Initialize(_ string, heartbeat time.Duration, _ time.Duration, _ map[string]string) error ***REMOVED***
	s.heartbeat = heartbeat
	s.values = make([]string, 0)
	return nil
***REMOVED***

// Watch sends periodic discovery updates to a channel.
func (s *Discovery) Watch(stopCh <-chan struct***REMOVED******REMOVED***) (<-chan discovery.Entries, <-chan error) ***REMOVED***
	ch := make(chan discovery.Entries)
	errCh := make(chan error)
	ticker := time.NewTicker(s.heartbeat)

	go func() ***REMOVED***
		defer close(errCh)
		defer close(ch)

		// Send the initial entries if available.
		var currentEntries discovery.Entries
		var err error

		s.mu.Lock()
		if len(s.values) > 0 ***REMOVED***
			currentEntries, err = discovery.CreateEntries(s.values)
		***REMOVED***
		s.mu.Unlock()

		if err != nil ***REMOVED***
			errCh <- err
		***REMOVED*** else if currentEntries != nil ***REMOVED***
			ch <- currentEntries
		***REMOVED***

		// Periodically send updates.
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				s.mu.Lock()
				newEntries, err := discovery.CreateEntries(s.values)
				s.mu.Unlock()
				if err != nil ***REMOVED***
					errCh <- err
					continue
				***REMOVED***

				// Check if the file has really changed.
				if !newEntries.Equals(currentEntries) ***REMOVED***
					ch <- newEntries
				***REMOVED***
				currentEntries = newEntries
			case <-stopCh:
				ticker.Stop()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch, errCh
***REMOVED***

// Register adds a new address to the discovery.
func (s *Discovery) Register(addr string) error ***REMOVED***
	s.mu.Lock()
	s.values = append(s.values, addr)
	s.mu.Unlock()
	return nil
***REMOVED***
