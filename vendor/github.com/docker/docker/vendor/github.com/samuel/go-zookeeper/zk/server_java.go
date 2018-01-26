package zk

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type ErrMissingServerConfigField string

func (e ErrMissingServerConfigField) Error() string ***REMOVED***
	return fmt.Sprintf("zk: missing server config field '%s'", string(e))
***REMOVED***

const (
	DefaultServerTickTime                 = 2000
	DefaultServerInitLimit                = 10
	DefaultServerSyncLimit                = 5
	DefaultServerAutoPurgeSnapRetainCount = 3
	DefaultPeerPort                       = 2888
	DefaultLeaderElectionPort             = 3888
)

type ServerConfigServer struct ***REMOVED***
	ID                 int
	Host               string
	PeerPort           int
	LeaderElectionPort int
***REMOVED***

type ServerConfig struct ***REMOVED***
	TickTime                 int    // Number of milliseconds of each tick
	InitLimit                int    // Number of ticks that the initial synchronization phase can take
	SyncLimit                int    // Number of ticks that can pass between sending a request and getting an acknowledgement
	DataDir                  string // Direcrory where the snapshot is stored
	ClientPort               int    // Port at which clients will connect
	AutoPurgeSnapRetainCount int    // Number of snapshots to retain in dataDir
	AutoPurgePurgeInterval   int    // Purge task internal in hours (0 to disable auto purge)
	Servers                  []ServerConfigServer
***REMOVED***

func (sc ServerConfig) Marshall(w io.Writer) error ***REMOVED***
	if sc.DataDir == "" ***REMOVED***
		return ErrMissingServerConfigField("dataDir")
	***REMOVED***
	fmt.Fprintf(w, "dataDir=%s\n", sc.DataDir)
	if sc.TickTime <= 0 ***REMOVED***
		sc.TickTime = DefaultServerTickTime
	***REMOVED***
	fmt.Fprintf(w, "tickTime=%d\n", sc.TickTime)
	if sc.InitLimit <= 0 ***REMOVED***
		sc.InitLimit = DefaultServerInitLimit
	***REMOVED***
	fmt.Fprintf(w, "initLimit=%d\n", sc.InitLimit)
	if sc.SyncLimit <= 0 ***REMOVED***
		sc.SyncLimit = DefaultServerSyncLimit
	***REMOVED***
	fmt.Fprintf(w, "syncLimit=%d\n", sc.SyncLimit)
	if sc.ClientPort <= 0 ***REMOVED***
		sc.ClientPort = DefaultPort
	***REMOVED***
	fmt.Fprintf(w, "clientPort=%d\n", sc.ClientPort)
	if sc.AutoPurgePurgeInterval > 0 ***REMOVED***
		if sc.AutoPurgeSnapRetainCount <= 0 ***REMOVED***
			sc.AutoPurgeSnapRetainCount = DefaultServerAutoPurgeSnapRetainCount
		***REMOVED***
		fmt.Fprintf(w, "autopurge.snapRetainCount=%d\n", sc.AutoPurgeSnapRetainCount)
		fmt.Fprintf(w, "autopurge.purgeInterval=%d\n", sc.AutoPurgePurgeInterval)
	***REMOVED***
	if len(sc.Servers) > 0 ***REMOVED***
		for _, srv := range sc.Servers ***REMOVED***
			if srv.PeerPort <= 0 ***REMOVED***
				srv.PeerPort = DefaultPeerPort
			***REMOVED***
			if srv.LeaderElectionPort <= 0 ***REMOVED***
				srv.LeaderElectionPort = DefaultLeaderElectionPort
			***REMOVED***
			fmt.Fprintf(w, "server.%d=%s:%d:%d\n", srv.ID, srv.Host, srv.PeerPort, srv.LeaderElectionPort)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var jarSearchPaths = []string***REMOVED***
	"zookeeper-*/contrib/fatjar/zookeeper-*-fatjar.jar",
	"../zookeeper-*/contrib/fatjar/zookeeper-*-fatjar.jar",
	"/usr/share/java/zookeeper-*.jar",
	"/usr/local/zookeeper-*/contrib/fatjar/zookeeper-*-fatjar.jar",
	"/usr/local/Cellar/zookeeper/*/libexec/contrib/fatjar/zookeeper-*-fatjar.jar",
***REMOVED***

func findZookeeperFatJar() string ***REMOVED***
	var paths []string
	zkPath := os.Getenv("ZOOKEEPER_PATH")
	if zkPath == "" ***REMOVED***
		paths = jarSearchPaths
	***REMOVED*** else ***REMOVED***
		paths = []string***REMOVED***filepath.Join(zkPath, "contrib/fatjar/zookeeper-*-fatjar.jar")***REMOVED***
	***REMOVED***
	for _, path := range paths ***REMOVED***
		matches, _ := filepath.Glob(path)
		// TODO: could sort by version and pick latest
		if len(matches) > 0 ***REMOVED***
			return matches[0]
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

type Server struct ***REMOVED***
	JarPath        string
	ConfigPath     string
	Stdout, Stderr io.Writer

	cmd *exec.Cmd
***REMOVED***

func (srv *Server) Start() error ***REMOVED***
	if srv.JarPath == "" ***REMOVED***
		srv.JarPath = findZookeeperFatJar()
		if srv.JarPath == "" ***REMOVED***
			return fmt.Errorf("zk: unable to find server jar")
		***REMOVED***
	***REMOVED***
	srv.cmd = exec.Command("java", "-jar", srv.JarPath, "server", srv.ConfigPath)
	srv.cmd.Stdout = srv.Stdout
	srv.cmd.Stderr = srv.Stderr
	return srv.cmd.Start()
***REMOVED***

func (srv *Server) Stop() error ***REMOVED***
	srv.cmd.Process.Signal(os.Kill)
	return srv.cmd.Wait()
***REMOVED***
