package listeners

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/Microsoft/go-winio"
	"github.com/docker/go-connections/sockets"
)

// Init creates new listeners for the server.
func Init(proto, addr, socketGroup string, tlsConfig *tls.Config) ([]net.Listener, error) ***REMOVED***
	ls := []net.Listener***REMOVED******REMOVED***

	switch proto ***REMOVED***
	case "tcp":
		l, err := sockets.NewTCPSocket(addr, tlsConfig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ls = append(ls, l)

	case "npipe":
		// allow Administrators and SYSTEM, plus whatever additional users or groups were specified
		sddl := "D:P(A;;GA;;;BA)(A;;GA;;;SY)"
		if socketGroup != "" ***REMOVED***
			for _, g := range strings.Split(socketGroup, ",") ***REMOVED***
				sid, err := winio.LookupSidByName(g)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				sddl += fmt.Sprintf("(A;;GRGW;;;%s)", sid)
			***REMOVED***
		***REMOVED***
		c := winio.PipeConfig***REMOVED***
			SecurityDescriptor: sddl,
			MessageMode:        true,  // Use message mode so that CloseWrite() is supported
			InputBufferSize:    65536, // Use 64KB buffers to improve performance
			OutputBufferSize:   65536,
		***REMOVED***
		l, err := winio.ListenPipe(addr, &c)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ls = append(ls, l)

	default:
		return nil, fmt.Errorf("invalid protocol format: windows only supports tcp and npipe")
	***REMOVED***

	return ls, nil
***REMOVED***
