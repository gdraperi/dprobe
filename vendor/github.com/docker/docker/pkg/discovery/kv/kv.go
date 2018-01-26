package kv

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/docker/docker/pkg/discovery"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/sirupsen/logrus"
)

const (
	defaultDiscoveryPath = "docker/nodes"
)

// Discovery is exported
type Discovery struct ***REMOVED***
	backend   store.Backend
	store     store.Store
	heartbeat time.Duration
	ttl       time.Duration
	prefix    string
	path      string
***REMOVED***

func init() ***REMOVED***
	Init()
***REMOVED***

// Init is exported
func Init() ***REMOVED***
	// Register to libkv
	zookeeper.Register()
	consul.Register()
	etcd.Register()

	// Register to internal discovery service
	discovery.Register("zk", &Discovery***REMOVED***backend: store.ZK***REMOVED***)
	discovery.Register("consul", &Discovery***REMOVED***backend: store.CONSUL***REMOVED***)
	discovery.Register("etcd", &Discovery***REMOVED***backend: store.ETCD***REMOVED***)
***REMOVED***

// Initialize is exported
func (s *Discovery) Initialize(uris string, heartbeat time.Duration, ttl time.Duration, clusterOpts map[string]string) error ***REMOVED***
	var (
		parts = strings.SplitN(uris, "/", 2)
		addrs = strings.Split(parts[0], ",")
		err   error
	)

	// A custom prefix to the path can be optionally used.
	if len(parts) == 2 ***REMOVED***
		s.prefix = parts[1]
	***REMOVED***

	s.heartbeat = heartbeat
	s.ttl = ttl

	// Use a custom path if specified in discovery options
	dpath := defaultDiscoveryPath
	if clusterOpts["kv.path"] != "" ***REMOVED***
		dpath = clusterOpts["kv.path"]
	***REMOVED***

	s.path = path.Join(s.prefix, dpath)

	var config *store.Config
	if clusterOpts["kv.cacertfile"] != "" && clusterOpts["kv.certfile"] != "" && clusterOpts["kv.keyfile"] != "" ***REMOVED***
		logrus.Info("Initializing discovery with TLS")
		tlsConfig, err := tlsconfig.Client(tlsconfig.Options***REMOVED***
			CAFile:   clusterOpts["kv.cacertfile"],
			CertFile: clusterOpts["kv.certfile"],
			KeyFile:  clusterOpts["kv.keyfile"],
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		config = &store.Config***REMOVED***
			// Set ClientTLS to trigger https (bug in libkv/etcd)
			ClientTLS: &store.ClientTLSConfig***REMOVED***
				CACertFile: clusterOpts["kv.cacertfile"],
				CertFile:   clusterOpts["kv.certfile"],
				KeyFile:    clusterOpts["kv.keyfile"],
			***REMOVED***,
			// The actual TLS config that will be used
			TLS: tlsConfig,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		logrus.Info("Initializing discovery without TLS")
	***REMOVED***

	// Creates a new store, will ignore options given
	// if not supported by the chosen store
	s.store, err = libkv.NewStore(s.backend, addrs, config)
	return err
***REMOVED***

// Watch the store until either there's a store error or we receive a stop request.
// Returns false if we shouldn't attempt watching the store anymore (stop request received).
func (s *Discovery) watchOnce(stopCh <-chan struct***REMOVED******REMOVED***, watchCh <-chan []*store.KVPair, discoveryCh chan discovery.Entries, errCh chan error) bool ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case pairs := <-watchCh:
			if pairs == nil ***REMOVED***
				return true
			***REMOVED***

			logrus.WithField("discovery", s.backend).Debugf("Watch triggered with %d nodes", len(pairs))

			// Convert `KVPair` into `discovery.Entry`.
			addrs := make([]string, len(pairs))
			for _, pair := range pairs ***REMOVED***
				addrs = append(addrs, string(pair.Value))
			***REMOVED***

			entries, err := discovery.CreateEntries(addrs)
			if err != nil ***REMOVED***
				errCh <- err
			***REMOVED*** else ***REMOVED***
				discoveryCh <- entries
			***REMOVED***
		case <-stopCh:
			// We were requested to stop watching.
			return false
		***REMOVED***
	***REMOVED***
***REMOVED***

// Watch is exported
func (s *Discovery) Watch(stopCh <-chan struct***REMOVED******REMOVED***) (<-chan discovery.Entries, <-chan error) ***REMOVED***
	ch := make(chan discovery.Entries)
	errCh := make(chan error)

	go func() ***REMOVED***
		defer close(ch)
		defer close(errCh)

		// Forever: Create a store watch, watch until we get an error and then try again.
		// Will only stop if we receive a stopCh request.
		for ***REMOVED***
			// Create the path to watch if it does not exist yet
			exists, err := s.store.Exists(s.path)
			if err != nil ***REMOVED***
				errCh <- err
			***REMOVED***
			if !exists ***REMOVED***
				if err := s.store.Put(s.path, []byte(""), &store.WriteOptions***REMOVED***IsDir: true***REMOVED***); err != nil ***REMOVED***
					errCh <- err
				***REMOVED***
			***REMOVED***

			// Set up a watch.
			watchCh, err := s.store.WatchTree(s.path, stopCh)
			if err != nil ***REMOVED***
				errCh <- err
			***REMOVED*** else ***REMOVED***
				if !s.watchOnce(stopCh, watchCh, ch, errCh) ***REMOVED***
					return
				***REMOVED***
			***REMOVED***

			// If we get here it means the store watch channel was closed. This
			// is unexpected so let's retry later.
			errCh <- fmt.Errorf("Unexpected watch error")
			time.Sleep(s.heartbeat)
		***REMOVED***
	***REMOVED***()
	return ch, errCh
***REMOVED***

// Register is exported
func (s *Discovery) Register(addr string) error ***REMOVED***
	opts := &store.WriteOptions***REMOVED***TTL: s.ttl***REMOVED***
	return s.store.Put(path.Join(s.path, addr), []byte(addr), opts)
***REMOVED***

// Store returns the underlying store used by KV discovery.
func (s *Discovery) Store() store.Store ***REMOVED***
	return s.store
***REMOVED***

// Prefix returns the store prefix
func (s *Discovery) Prefix() string ***REMOVED***
	return s.prefix
***REMOVED***
