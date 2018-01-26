package discovery

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/pkg/discovery"
	"github.com/sirupsen/logrus"

	// Register the libkv backends for discovery.
	_ "github.com/docker/docker/pkg/discovery/kv"
)

const (
	// defaultDiscoveryHeartbeat is the default value for discovery heartbeat interval.
	defaultDiscoveryHeartbeat = 20 * time.Second
	// defaultDiscoveryTTLFactor is the default TTL factor for discovery
	defaultDiscoveryTTLFactor = 3
)

// ErrDiscoveryDisabled is an error returned if the discovery is disabled
var ErrDiscoveryDisabled = errors.New("discovery is disabled")

// Reloader is the discovery reloader of the daemon
type Reloader interface ***REMOVED***
	discovery.Watcher
	Stop()
	Reload(backend, address string, clusterOpts map[string]string) error
	ReadyCh() <-chan struct***REMOVED******REMOVED***
***REMOVED***

type daemonDiscoveryReloader struct ***REMOVED***
	backend discovery.Backend
	ticker  *time.Ticker
	term    chan bool
	readyCh chan struct***REMOVED******REMOVED***
***REMOVED***

func (d *daemonDiscoveryReloader) Watch(stopCh <-chan struct***REMOVED******REMOVED***) (<-chan discovery.Entries, <-chan error) ***REMOVED***
	return d.backend.Watch(stopCh)
***REMOVED***

func (d *daemonDiscoveryReloader) ReadyCh() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return d.readyCh
***REMOVED***

func discoveryOpts(clusterOpts map[string]string) (time.Duration, time.Duration, error) ***REMOVED***
	var (
		heartbeat = defaultDiscoveryHeartbeat
		ttl       = defaultDiscoveryTTLFactor * defaultDiscoveryHeartbeat
	)

	if hb, ok := clusterOpts["discovery.heartbeat"]; ok ***REMOVED***
		h, err := strconv.Atoi(hb)
		if err != nil ***REMOVED***
			return time.Duration(0), time.Duration(0), err
		***REMOVED***

		if h <= 0 ***REMOVED***
			return time.Duration(0), time.Duration(0),
				fmt.Errorf("discovery.heartbeat must be positive")
		***REMOVED***

		heartbeat = time.Duration(h) * time.Second
		ttl = defaultDiscoveryTTLFactor * heartbeat
	***REMOVED***

	if tstr, ok := clusterOpts["discovery.ttl"]; ok ***REMOVED***
		t, err := strconv.Atoi(tstr)
		if err != nil ***REMOVED***
			return time.Duration(0), time.Duration(0), err
		***REMOVED***

		if t <= 0 ***REMOVED***
			return time.Duration(0), time.Duration(0),
				fmt.Errorf("discovery.ttl must be positive")
		***REMOVED***

		ttl = time.Duration(t) * time.Second

		if _, ok := clusterOpts["discovery.heartbeat"]; !ok ***REMOVED***
			heartbeat = time.Duration(t) * time.Second / time.Duration(defaultDiscoveryTTLFactor)
		***REMOVED***

		if ttl <= heartbeat ***REMOVED***
			return time.Duration(0), time.Duration(0),
				fmt.Errorf("discovery.ttl timer must be greater than discovery.heartbeat")
		***REMOVED***
	***REMOVED***

	return heartbeat, ttl, nil
***REMOVED***

// Init initializes the nodes discovery subsystem by connecting to the specified backend
// and starts a registration loop to advertise the current node under the specified address.
func Init(backendAddress, advertiseAddress string, clusterOpts map[string]string) (Reloader, error) ***REMOVED***
	heartbeat, backend, err := parseDiscoveryOptions(backendAddress, clusterOpts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	reloader := &daemonDiscoveryReloader***REMOVED***
		backend: backend,
		ticker:  time.NewTicker(heartbeat),
		term:    make(chan bool),
		readyCh: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	// We call Register() on the discovery backend in a loop for the whole lifetime of the daemon,
	// but we never actually Watch() for nodes appearing and disappearing for the moment.
	go reloader.advertiseHeartbeat(advertiseAddress)
	return reloader, nil
***REMOVED***

// advertiseHeartbeat registers the current node against the discovery backend using the specified
// address. The function never returns, as registration against the backend comes with a TTL and
// requires regular heartbeats.
func (d *daemonDiscoveryReloader) advertiseHeartbeat(address string) ***REMOVED***
	var ready bool
	if err := d.initHeartbeat(address); err == nil ***REMOVED***
		ready = true
		close(d.readyCh)
	***REMOVED*** else ***REMOVED***
		logrus.WithError(err).Debug("First discovery heartbeat failed")
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case <-d.ticker.C:
			if err := d.backend.Register(address); err != nil ***REMOVED***
				logrus.Warnf("Registering as %q in discovery failed: %v", address, err)
			***REMOVED*** else ***REMOVED***
				if !ready ***REMOVED***
					close(d.readyCh)
					ready = true
				***REMOVED***
			***REMOVED***
		case <-d.term:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// initHeartbeat is used to do the first heartbeat. It uses a tight loop until
// either the timeout period is reached or the heartbeat is successful and returns.
func (d *daemonDiscoveryReloader) initHeartbeat(address string) error ***REMOVED***
	// Setup a short ticker until the first heartbeat has succeeded
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	// timeout makes sure that after a period of time we stop being so aggressive trying to reach the discovery service
	timeout := time.After(60 * time.Second)

	for ***REMOVED***
		select ***REMOVED***
		case <-timeout:
			return errors.New("timeout waiting for initial discovery")
		case <-d.term:
			return errors.New("terminated")
		case <-t.C:
			if err := d.backend.Register(address); err == nil ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Reload makes the watcher to stop advertising and reconfigures it to advertise in a new address.
func (d *daemonDiscoveryReloader) Reload(backendAddress, advertiseAddress string, clusterOpts map[string]string) error ***REMOVED***
	d.Stop()

	heartbeat, backend, err := parseDiscoveryOptions(backendAddress, clusterOpts)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	d.backend = backend
	d.ticker = time.NewTicker(heartbeat)
	d.readyCh = make(chan struct***REMOVED******REMOVED***)

	go d.advertiseHeartbeat(advertiseAddress)
	return nil
***REMOVED***

// Stop terminates the discovery advertising.
func (d *daemonDiscoveryReloader) Stop() ***REMOVED***
	d.ticker.Stop()
	d.term <- true
***REMOVED***

func parseDiscoveryOptions(backendAddress string, clusterOpts map[string]string) (time.Duration, discovery.Backend, error) ***REMOVED***
	heartbeat, ttl, err := discoveryOpts(clusterOpts)
	if err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***

	backend, err := discovery.New(backendAddress, heartbeat, ttl, clusterOpts)
	if err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***
	return heartbeat, backend, nil
***REMOVED***
