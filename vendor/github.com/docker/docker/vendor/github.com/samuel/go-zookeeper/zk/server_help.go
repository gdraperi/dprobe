package zk

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type TestServer struct ***REMOVED***
	Port int
	Path string
	Srv  *Server
***REMOVED***

type TestCluster struct ***REMOVED***
	Path    string
	Servers []TestServer
***REMOVED***

func StartTestCluster(size int, stdout, stderr io.Writer) (*TestCluster, error) ***REMOVED***
	tmpPath, err := ioutil.TempDir("", "gozk")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	success := false
	startPort := int(rand.Int31n(6000) + 10000)
	cluster := &TestCluster***REMOVED***Path: tmpPath***REMOVED***
	defer func() ***REMOVED***
		if !success ***REMOVED***
			cluster.Stop()
		***REMOVED***
	***REMOVED***()
	for serverN := 0; serverN < size; serverN++ ***REMOVED***
		srvPath := filepath.Join(tmpPath, fmt.Sprintf("srv%d", serverN))
		if err := os.Mkdir(srvPath, 0700); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		port := startPort + serverN*3
		cfg := ServerConfig***REMOVED***
			ClientPort: port,
			DataDir:    srvPath,
		***REMOVED***
		for i := 0; i < size; i++ ***REMOVED***
			cfg.Servers = append(cfg.Servers, ServerConfigServer***REMOVED***
				ID:                 i + 1,
				Host:               "127.0.0.1",
				PeerPort:           startPort + i*3 + 1,
				LeaderElectionPort: startPort + i*3 + 2,
			***REMOVED***)
		***REMOVED***
		cfgPath := filepath.Join(srvPath, "zoo.cfg")
		fi, err := os.Create(cfgPath)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		err = cfg.Marshall(fi)
		fi.Close()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		fi, err = os.Create(filepath.Join(srvPath, "myid"))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		_, err = fmt.Fprintf(fi, "%d\n", serverN+1)
		fi.Close()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		srv := &Server***REMOVED***
			ConfigPath: cfgPath,
			Stdout:     stdout,
			Stderr:     stderr,
		***REMOVED***
		if err := srv.Start(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cluster.Servers = append(cluster.Servers, TestServer***REMOVED***
			Path: srvPath,
			Port: cfg.ClientPort,
			Srv:  srv,
		***REMOVED***)
	***REMOVED***
	success = true
	time.Sleep(time.Second) // Give the server time to become active. Should probably actually attempt to connect to verify.
	return cluster, nil
***REMOVED***

func (ts *TestCluster) Connect(idx int) (*Conn, error) ***REMOVED***
	zk, _, err := Connect([]string***REMOVED***fmt.Sprintf("127.0.0.1:%d", ts.Servers[idx].Port)***REMOVED***, time.Second*15)
	return zk, err
***REMOVED***

func (ts *TestCluster) ConnectAll() (*Conn, <-chan Event, error) ***REMOVED***
	return ts.ConnectAllTimeout(time.Second * 15)
***REMOVED***

func (ts *TestCluster) ConnectAllTimeout(sessionTimeout time.Duration) (*Conn, <-chan Event, error) ***REMOVED***
	hosts := make([]string, len(ts.Servers))
	for i, srv := range ts.Servers ***REMOVED***
		hosts[i] = fmt.Sprintf("127.0.0.1:%d", srv.Port)
	***REMOVED***
	zk, ch, err := Connect(hosts, sessionTimeout)
	return zk, ch, err
***REMOVED***

func (ts *TestCluster) Stop() error ***REMOVED***
	for _, srv := range ts.Servers ***REMOVED***
		srv.Srv.Stop()
	***REMOVED***
	defer os.RemoveAll(ts.Path)
	return nil
***REMOVED***
