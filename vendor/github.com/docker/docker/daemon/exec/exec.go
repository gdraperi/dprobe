package exec

import (
	"runtime"
	"sync"

	"github.com/containerd/containerd/cio"
	"github.com/docker/docker/container/stream"
	"github.com/docker/docker/pkg/stringid"
	"github.com/sirupsen/logrus"
)

// Config holds the configurations for execs. The Daemon keeps
// track of both running and finished execs so that they can be
// examined both during and after completion.
type Config struct ***REMOVED***
	sync.Mutex
	StreamConfig *stream.Config
	ID           string
	Running      bool
	ExitCode     *int
	OpenStdin    bool
	OpenStderr   bool
	OpenStdout   bool
	CanRemove    bool
	ContainerID  string
	DetachKeys   []byte
	Entrypoint   string
	Args         []string
	Tty          bool
	Privileged   bool
	User         string
	WorkingDir   string
	Env          []string
	Pid          int
***REMOVED***

// NewConfig initializes the a new exec configuration
func NewConfig() *Config ***REMOVED***
	return &Config***REMOVED***
		ID:           stringid.GenerateNonCryptoID(),
		StreamConfig: stream.NewConfig(),
	***REMOVED***
***REMOVED***

type rio struct ***REMOVED***
	cio.IO

	sc *stream.Config
***REMOVED***

func (i *rio) Close() error ***REMOVED***
	i.IO.Close()

	return i.sc.CloseStreams()
***REMOVED***

func (i *rio) Wait() ***REMOVED***
	i.sc.Wait()

	i.IO.Wait()
***REMOVED***

// InitializeStdio is called by libcontainerd to connect the stdio.
func (c *Config) InitializeStdio(iop *cio.DirectIO) (cio.IO, error) ***REMOVED***
	c.StreamConfig.CopyToPipe(iop)

	if c.StreamConfig.Stdin() == nil && !c.Tty && runtime.GOOS == "windows" ***REMOVED***
		if iop.Stdin != nil ***REMOVED***
			if err := iop.Stdin.Close(); err != nil ***REMOVED***
				logrus.Errorf("error closing exec stdin: %+v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return &rio***REMOVED***IO: iop, sc: c.StreamConfig***REMOVED***, nil
***REMOVED***

// CloseStreams closes the stdio streams for the exec
func (c *Config) CloseStreams() error ***REMOVED***
	return c.StreamConfig.CloseStreams()
***REMOVED***

// SetExitCode sets the exec config's exit code
func (c *Config) SetExitCode(code int) ***REMOVED***
	c.ExitCode = &code
***REMOVED***

// Store keeps track of the exec configurations.
type Store struct ***REMOVED***
	byID map[string]*Config
	sync.RWMutex
***REMOVED***

// NewStore initializes a new exec store.
func NewStore() *Store ***REMOVED***
	return &Store***REMOVED***
		byID: make(map[string]*Config),
	***REMOVED***
***REMOVED***

// Commands returns the exec configurations in the store.
func (e *Store) Commands() map[string]*Config ***REMOVED***
	e.RLock()
	byID := make(map[string]*Config, len(e.byID))
	for id, config := range e.byID ***REMOVED***
		byID[id] = config
	***REMOVED***
	e.RUnlock()
	return byID
***REMOVED***

// Add adds a new exec configuration to the store.
func (e *Store) Add(id string, Config *Config) ***REMOVED***
	e.Lock()
	e.byID[id] = Config
	e.Unlock()
***REMOVED***

// Get returns an exec configuration by its id.
func (e *Store) Get(id string) *Config ***REMOVED***
	e.RLock()
	res := e.byID[id]
	e.RUnlock()
	return res
***REMOVED***

// Delete removes an exec configuration from the store.
func (e *Store) Delete(id string, pid int) ***REMOVED***
	e.Lock()
	delete(e.byID, id)
	e.Unlock()
***REMOVED***

// List returns the list of exec ids in the store.
func (e *Store) List() []string ***REMOVED***
	var IDs []string
	e.RLock()
	for id := range e.byID ***REMOVED***
		IDs = append(IDs, id)
	***REMOVED***
	e.RUnlock()
	return IDs
***REMOVED***
