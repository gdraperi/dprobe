package client

import (
	"net"
	"net/http"

	"golang.org/x/net/context"
)

// DialSession returns a connection that can be used communication with daemon
func (cli *Client) DialSession(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) ***REMOVED***
	req, err := http.NewRequest("POST", "/session", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req = cli.addHeaders(req, meta)

	return cli.setupHijackConn(req, proto)
***REMOVED***
