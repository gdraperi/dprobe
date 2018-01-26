package file

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/docker/pkg/discovery"
)

// Discovery is exported
type Discovery struct ***REMOVED***
	heartbeat time.Duration
	path      string
***REMOVED***

func init() ***REMOVED***
	Init()
***REMOVED***

// Init is exported
func Init() ***REMOVED***
	discovery.Register("file", &Discovery***REMOVED******REMOVED***)
***REMOVED***

// Initialize is exported
func (s *Discovery) Initialize(path string, heartbeat time.Duration, ttl time.Duration, _ map[string]string) error ***REMOVED***
	s.path = path
	s.heartbeat = heartbeat
	return nil
***REMOVED***

func parseFileContent(content []byte) []string ***REMOVED***
	var result []string
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") ***REMOVED***
		line = strings.TrimSpace(line)
		// Ignoring line starts with #
		if strings.HasPrefix(line, "#") ***REMOVED***
			continue
		***REMOVED***
		// Inlined # comment also ignored.
		if strings.Contains(line, "#") ***REMOVED***
			line = line[0:strings.Index(line, "#")]
			// Trim additional spaces caused by above stripping.
			line = strings.TrimSpace(line)
		***REMOVED***
		result = append(result, discovery.Generate(line)...)
	***REMOVED***
	return result
***REMOVED***

func (s *Discovery) fetch() (discovery.Entries, error) ***REMOVED***
	fileContent, err := ioutil.ReadFile(s.path)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to read '%s': %v", s.path, err)
	***REMOVED***
	return discovery.CreateEntries(parseFileContent(fileContent))
***REMOVED***

// Watch is exported
func (s *Discovery) Watch(stopCh <-chan struct***REMOVED******REMOVED***) (<-chan discovery.Entries, <-chan error) ***REMOVED***
	ch := make(chan discovery.Entries)
	errCh := make(chan error)
	ticker := time.NewTicker(s.heartbeat)

	go func() ***REMOVED***
		defer close(errCh)
		defer close(ch)

		// Send the initial entries if available.
		currentEntries, err := s.fetch()
		if err != nil ***REMOVED***
			errCh <- err
		***REMOVED*** else ***REMOVED***
			ch <- currentEntries
		***REMOVED***

		// Periodically send updates.
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				newEntries, err := s.fetch()
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

// Register is exported
func (s *Discovery) Register(addr string) error ***REMOVED***
	return discovery.ErrNotImplemented
***REMOVED***
