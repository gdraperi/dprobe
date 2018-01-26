package portmapper

import "net"

func newMockProxyCommand(proto string, hostIP net.IP, hostPort int, containerIP net.IP, containerPort int, userlandProxyPath string) (userlandProxy, error) ***REMOVED***
	return &mockProxyCommand***REMOVED******REMOVED***, nil
***REMOVED***

type mockProxyCommand struct ***REMOVED***
***REMOVED***

func (p *mockProxyCommand) Start() error ***REMOVED***
	return nil
***REMOVED***

func (p *mockProxyCommand) Stop() error ***REMOVED***
	return nil
***REMOVED***
