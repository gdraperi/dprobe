package session

import (
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func serve(ctx context.Context, grpcServer *grpc.Server, conn net.Conn) ***REMOVED***
	go func() ***REMOVED***
		<-ctx.Done()
		conn.Close()
	***REMOVED***()
	logrus.Debugf("serving grpc connection")
	(&http2.Server***REMOVED******REMOVED***).ServeConn(conn, &http2.ServeConnOpts***REMOVED***Handler: grpcServer***REMOVED***)
***REMOVED***

func grpcClientConn(ctx context.Context, conn net.Conn) (context.Context, *grpc.ClientConn, error) ***REMOVED***
	dialOpt := grpc.WithDialer(func(addr string, d time.Duration) (net.Conn, error) ***REMOVED***
		return conn, nil
	***REMOVED***)

	cc, err := grpc.DialContext(ctx, "", dialOpt, grpc.WithInsecure())
	if err != nil ***REMOVED***
		return nil, nil, errors.Wrap(err, "failed to create grpc client")
	***REMOVED***

	ctx, cancel := context.WithCancel(ctx)
	go monitorHealth(ctx, cc, cancel)

	return ctx, cc, nil
***REMOVED***

func monitorHealth(ctx context.Context, cc *grpc.ClientConn, cancelConn func()) ***REMOVED***
	defer cancelConn()
	defer cc.Close()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	healthClient := grpc_health_v1.NewHealthClient(cc)

	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return
		case <-ticker.C:
			<-ticker.C
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest***REMOVED******REMOVED***)
			cancel()
			if err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
