package funker

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"time"
)

// Call a Funker function
func Call(name string, args interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	argsJSON, err := json.Marshal(args)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	addr, err := net.ResolveTCPAddr("tcp", name+":9999")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Keepalive is a workaround for docker/docker#29655 .
	// The implementation of FIN_WAIT2 seems weird on Swarm-mode.
	// It seems always refuseing any packet after 60 seconds.
	//
	// TODO: remove this workaround if the issue gets resolved on the Docker side
	if err := conn.SetKeepAlive(true); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := conn.SetKeepAlivePeriod(30 * time.Second); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err = conn.Write(argsJSON); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = conn.CloseWrite(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	retJSON, err := ioutil.ReadAll(conn)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var ret interface***REMOVED******REMOVED***
	err = json.Unmarshal(retJSON, &ret)
	return ret, err
***REMOVED***
