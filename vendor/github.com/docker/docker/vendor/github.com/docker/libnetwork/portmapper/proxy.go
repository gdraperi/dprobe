package portmapper

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"time"
)

var userlandProxyCommandName = "docker-proxy"

type userlandProxy interface ***REMOVED***
	Start() error
	Stop() error
***REMOVED***

// proxyCommand wraps an exec.Cmd to run the userland TCP and UDP
// proxies as separate processes.
type proxyCommand struct ***REMOVED***
	cmd *exec.Cmd
***REMOVED***

func (p *proxyCommand) Start() error ***REMOVED***
	r, w, err := os.Pipe()
	if err != nil ***REMOVED***
		return fmt.Errorf("proxy unable to open os.Pipe %s", err)
	***REMOVED***
	defer r.Close()
	p.cmd.ExtraFiles = []*os.File***REMOVED***w***REMOVED***
	if err := p.cmd.Start(); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.Close()

	errchan := make(chan error, 1)
	go func() ***REMOVED***
		buf := make([]byte, 2)
		r.Read(buf)

		if string(buf) != "0\n" ***REMOVED***
			errStr, err := ioutil.ReadAll(r)
			if err != nil ***REMOVED***
				errchan <- fmt.Errorf("Error reading exit status from userland proxy: %v", err)
				return
			***REMOVED***

			errchan <- fmt.Errorf("Error starting userland proxy: %s", errStr)
			return
		***REMOVED***
		errchan <- nil
	***REMOVED***()

	select ***REMOVED***
	case err := <-errchan:
		return err
	case <-time.After(16 * time.Second):
		return fmt.Errorf("Timed out proxy starting the userland proxy")
	***REMOVED***
***REMOVED***

func (p *proxyCommand) Stop() error ***REMOVED***
	if p.cmd.Process != nil ***REMOVED***
		if err := p.cmd.Process.Signal(os.Interrupt); err != nil ***REMOVED***
			return err
		***REMOVED***
		return p.cmd.Wait()
	***REMOVED***
	return nil
***REMOVED***

// dummyProxy just listen on some port, it is needed to prevent accidental
// port allocations on bound port, because without userland proxy we using
// iptables rules and not net.Listen
type dummyProxy struct ***REMOVED***
	listener io.Closer
	addr     net.Addr
***REMOVED***

func newDummyProxy(proto string, hostIP net.IP, hostPort int) userlandProxy ***REMOVED***
	switch proto ***REMOVED***
	case "tcp":
		addr := &net.TCPAddr***REMOVED***IP: hostIP, Port: hostPort***REMOVED***
		return &dummyProxy***REMOVED***addr: addr***REMOVED***
	case "udp":
		addr := &net.UDPAddr***REMOVED***IP: hostIP, Port: hostPort***REMOVED***
		return &dummyProxy***REMOVED***addr: addr***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (p *dummyProxy) Start() error ***REMOVED***
	switch addr := p.addr.(type) ***REMOVED***
	case *net.TCPAddr:
		l, err := net.ListenTCP("tcp", addr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		p.listener = l
	case *net.UDPAddr:
		l, err := net.ListenUDP("udp", addr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		p.listener = l
	default:
		return fmt.Errorf("Unknown addr type: %T", p.addr)
	***REMOVED***
	return nil
***REMOVED***

func (p *dummyProxy) Stop() error ***REMOVED***
	if p.listener != nil ***REMOVED***
		return p.listener.Close()
	***REMOVED***
	return nil
***REMOVED***
