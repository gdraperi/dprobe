// +build windows

package libcontainerd

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type remote struct ***REMOVED***
	sync.RWMutex

	logger  *logrus.Entry
	clients []*client

	// Options
	rootDir  string
	stateDir string
***REMOVED***

// New creates a fresh instance of libcontainerd remote.
func New(rootDir, stateDir string, options ...RemoteOption) (Remote, error) ***REMOVED***
	return &remote***REMOVED***
		logger:   logrus.WithField("module", "libcontainerd"),
		rootDir:  rootDir,
		stateDir: stateDir,
	***REMOVED***, nil
***REMOVED***

type client struct ***REMOVED***
	sync.Mutex

	rootDir    string
	stateDir   string
	backend    Backend
	logger     *logrus.Entry
	eventQ     queue
	containers map[string]*container
***REMOVED***

func (r *remote) NewClient(ns string, b Backend) (Client, error) ***REMOVED***
	c := &client***REMOVED***
		rootDir:    r.rootDir,
		stateDir:   r.stateDir,
		backend:    b,
		logger:     r.logger.WithField("namespace", ns),
		containers: make(map[string]*container),
	***REMOVED***
	r.Lock()
	r.clients = append(r.clients, c)
	r.Unlock()

	return c, nil
***REMOVED***

func (r *remote) Cleanup() ***REMOVED***
	// Nothing to do
***REMOVED***
