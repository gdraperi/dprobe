package raft

import (
	"time"

	"golang.org/x/net/context"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/state"
	"github.com/docker/swarmkit/manager/state/store"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// dial returns a grpc client connection
func dial(addr string, protocol string, creds credentials.TransportCredentials, timeout time.Duration) (*grpc.ClientConn, error) ***REMOVED***
	grpcOptions := []grpc.DialOption***REMOVED***
		grpc.WithBackoffMaxDelay(2 * time.Second),
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
	***REMOVED***

	if timeout != 0 ***REMOVED***
		grpcOptions = append(grpcOptions, grpc.WithTimeout(timeout))
	***REMOVED***

	return grpc.Dial(addr, grpcOptions...)
***REMOVED***

// Register registers the node raft server
func Register(server *grpc.Server, node *Node) ***REMOVED***
	api.RegisterRaftServer(server, node)
	api.RegisterRaftMembershipServer(server, node)
***REMOVED***

// WaitForLeader waits until node observe some leader in cluster. It returns
// error if ctx was cancelled before leader appeared.
func WaitForLeader(ctx context.Context, n *Node) error ***REMOVED***
	_, err := n.Leader()
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for err != nil ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		***REMOVED***
		_, err = n.Leader()
	***REMOVED***
	return nil
***REMOVED***

// WaitForCluster waits until node observes that the cluster wide config is
// committed to raft. This ensures that we can see and serve informations
// related to the cluster.
func WaitForCluster(ctx context.Context, n *Node) (cluster *api.Cluster, err error) ***REMOVED***
	watch, cancel := state.Watch(n.MemoryStore().WatchQueue(), api.EventCreateCluster***REMOVED******REMOVED***)
	defer cancel()

	var clusters []*api.Cluster
	n.MemoryStore().View(func(readTx store.ReadTx) ***REMOVED***
		clusters, err = store.FindClusters(readTx, store.ByName(store.DefaultClusterName))
	***REMOVED***)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(clusters) == 1 ***REMOVED***
		cluster = clusters[0]
	***REMOVED*** else ***REMOVED***
		select ***REMOVED***
		case e := <-watch:
			cluster = e.(api.EventCreateCluster).Cluster
		case <-ctx.Done():
			return nil, ctx.Err()
		***REMOVED***
	***REMOVED***

	return cluster, nil
***REMOVED***
