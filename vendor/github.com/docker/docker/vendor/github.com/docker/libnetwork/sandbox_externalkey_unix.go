// +build linux freebsd

package libnetwork

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/docker/libnetwork/types"
	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/sirupsen/logrus"
)

const udsBase = "/run/docker/libnetwork/"
const success = "success"

// processSetKeyReexec is a private function that must be called only on an reexec path
// It expects 3 args ***REMOVED*** [0] = "libnetwork-setkey", [1] = <container-id>, [2] = <controller-id> ***REMOVED***
// It also expects configs.HookState as a json string in <stdin>
// Refer to https://github.com/opencontainers/runc/pull/160/ for more information
func processSetKeyReexec() ***REMOVED***
	var err error

	// Return a failure to the calling process via ExitCode
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			logrus.Fatalf("%v", err)
		***REMOVED***
	***REMOVED***()

	// expecting 3 args ***REMOVED***[0]="libnetwork-setkey", [1]=<container-id>, [2]=<controller-id> ***REMOVED***
	if len(os.Args) < 3 ***REMOVED***
		err = fmt.Errorf("Re-exec expects 3 args, received : %d", len(os.Args))
		return
	***REMOVED***
	containerID := os.Args[1]

	// We expect configs.HookState as a json string in <stdin>
	stateBuf, err := ioutil.ReadAll(os.Stdin)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	var state configs.HookState
	if err = json.Unmarshal(stateBuf, &state); err != nil ***REMOVED***
		return
	***REMOVED***

	controllerID := os.Args[2]

	err = SetExternalKey(controllerID, containerID, fmt.Sprintf("/proc/%d/ns/net", state.Pid))
***REMOVED***

// SetExternalKey provides a convenient way to set an External key to a sandbox
func SetExternalKey(controllerID string, containerID string, key string) error ***REMOVED***
	keyData := setKeyData***REMOVED***
		ContainerID: containerID,
		Key:         key***REMOVED***

	c, err := net.Dial("unix", udsBase+controllerID+".sock")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer c.Close()

	if err = sendKey(c, keyData); err != nil ***REMOVED***
		return fmt.Errorf("sendKey failed with : %v", err)
	***REMOVED***
	return processReturn(c)
***REMOVED***

func sendKey(c net.Conn, data setKeyData) error ***REMOVED***
	var err error
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			c.Close()
		***REMOVED***
	***REMOVED***()

	var b []byte
	if b, err = json.Marshal(data); err != nil ***REMOVED***
		return err
	***REMOVED***

	_, err = c.Write(b)
	return err
***REMOVED***

func processReturn(r io.Reader) error ***REMOVED***
	buf := make([]byte, 1024)
	n, err := r.Read(buf[:])
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to read buf in processReturn : %v", err)
	***REMOVED***
	if string(buf[0:n]) != success ***REMOVED***
		return fmt.Errorf(string(buf[0:n]))
	***REMOVED***
	return nil
***REMOVED***

func (c *controller) startExternalKeyListener() error ***REMOVED***
	if err := os.MkdirAll(udsBase, 0600); err != nil ***REMOVED***
		return err
	***REMOVED***
	uds := udsBase + c.id + ".sock"
	l, err := net.Listen("unix", uds)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.Chmod(uds, 0600); err != nil ***REMOVED***
		l.Close()
		return err
	***REMOVED***
	c.Lock()
	c.extKeyListener = l
	c.Unlock()

	go c.acceptClientConnections(uds, l)
	return nil
***REMOVED***

func (c *controller) acceptClientConnections(sock string, l net.Listener) ***REMOVED***
	for ***REMOVED***
		conn, err := l.Accept()
		if err != nil ***REMOVED***
			if _, err1 := os.Stat(sock); os.IsNotExist(err1) ***REMOVED***
				logrus.Debugf("Unix socket %s doesn't exist. cannot accept client connections", sock)
				return
			***REMOVED***
			logrus.Errorf("Error accepting connection %v", err)
			continue
		***REMOVED***
		go func() ***REMOVED***
			defer conn.Close()

			err := c.processExternalKey(conn)
			ret := success
			if err != nil ***REMOVED***
				ret = err.Error()
			***REMOVED***

			_, err = conn.Write([]byte(ret))
			if err != nil ***REMOVED***
				logrus.Errorf("Error returning to the client %v", err)
			***REMOVED***
		***REMOVED***()
	***REMOVED***
***REMOVED***

func (c *controller) processExternalKey(conn net.Conn) error ***REMOVED***
	buf := make([]byte, 1280)
	nr, err := conn.Read(buf)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var s setKeyData
	if err = json.Unmarshal(buf[0:nr], &s); err != nil ***REMOVED***
		return err
	***REMOVED***

	var sandbox Sandbox
	search := SandboxContainerWalker(&sandbox, s.ContainerID)
	c.WalkSandboxes(search)
	if sandbox == nil ***REMOVED***
		return types.BadRequestErrorf("no sandbox present for %s", s.ContainerID)
	***REMOVED***

	return sandbox.SetKey(s.Key)
***REMOVED***

func (c *controller) stopExternalKeyListener() ***REMOVED***
	c.extKeyListener.Close()
***REMOVED***
