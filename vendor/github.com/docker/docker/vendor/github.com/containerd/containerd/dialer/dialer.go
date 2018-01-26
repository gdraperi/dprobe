package dialer

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

type dialResult struct ***REMOVED***
	c   net.Conn
	err error
***REMOVED***

// Dialer returns a GRPC net.Conn connected to the provided address
func Dialer(address string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	var (
		stopC = make(chan struct***REMOVED******REMOVED***)
		synC  = make(chan *dialResult)
	)
	go func() ***REMOVED***
		defer close(synC)
		for ***REMOVED***
			select ***REMOVED***
			case <-stopC:
				return
			default:
				c, err := dialer(address, timeout)
				if isNoent(err) ***REMOVED***
					<-time.After(10 * time.Millisecond)
					continue
				***REMOVED***
				synC <- &dialResult***REMOVED***c, err***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	select ***REMOVED***
	case dr := <-synC:
		return dr.c, dr.err
	case <-time.After(timeout):
		close(stopC)
		go func() ***REMOVED***
			dr := <-synC
			if dr != nil ***REMOVED***
				dr.c.Close()
			***REMOVED***
		***REMOVED***()
		return nil, errors.Errorf("dial %s: timeout", address)
	***REMOVED***
***REMOVED***
