// +build !windows

package libnetwork

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"

	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/libnetwork/iptables"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netns"
)

func init() ***REMOVED***
	reexec.Register("setup-resolver", reexecSetupResolver)
***REMOVED***

const (
	// outputChain used for docker embed dns
	outputChain = "DOCKER_OUTPUT"
	//postroutingchain used for docker embed dns
	postroutingchain = "DOCKER_POSTROUTING"
)

func reexecSetupResolver() ***REMOVED***
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if len(os.Args) < 4 ***REMOVED***
		logrus.Error("invalid number of arguments..")
		os.Exit(1)
	***REMOVED***

	resolverIP, ipPort, _ := net.SplitHostPort(os.Args[2])
	_, tcpPort, _ := net.SplitHostPort(os.Args[3])
	rules := [][]string***REMOVED***
		***REMOVED***"-t", "nat", "-I", outputChain, "-d", resolverIP, "-p", "udp", "--dport", dnsPort, "-j", "DNAT", "--to-destination", os.Args[2]***REMOVED***,
		***REMOVED***"-t", "nat", "-I", postroutingchain, "-s", resolverIP, "-p", "udp", "--sport", ipPort, "-j", "SNAT", "--to-source", ":" + dnsPort***REMOVED***,
		***REMOVED***"-t", "nat", "-I", outputChain, "-d", resolverIP, "-p", "tcp", "--dport", dnsPort, "-j", "DNAT", "--to-destination", os.Args[3]***REMOVED***,
		***REMOVED***"-t", "nat", "-I", postroutingchain, "-s", resolverIP, "-p", "tcp", "--sport", tcpPort, "-j", "SNAT", "--to-source", ":" + dnsPort***REMOVED***,
	***REMOVED***

	f, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)
	if err != nil ***REMOVED***
		logrus.Errorf("failed get network namespace %q: %v", os.Args[1], err)
		os.Exit(2)
	***REMOVED***
	defer f.Close()

	nsFD := f.Fd()
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil ***REMOVED***
		logrus.Errorf("setting into container net ns %v failed, %v", os.Args[1], err)
		os.Exit(3)
	***REMOVED***

	// insert outputChain and postroutingchain
	err = iptables.RawCombinedOutputNative("-t", "nat", "-C", "OUTPUT", "-d", resolverIP, "-j", outputChain)
	if err == nil ***REMOVED***
		iptables.RawCombinedOutputNative("-t", "nat", "-F", outputChain)
	***REMOVED*** else ***REMOVED***
		iptables.RawCombinedOutputNative("-t", "nat", "-N", outputChain)
		iptables.RawCombinedOutputNative("-t", "nat", "-I", "OUTPUT", "-d", resolverIP, "-j", outputChain)
	***REMOVED***

	err = iptables.RawCombinedOutputNative("-t", "nat", "-C", "POSTROUTING", "-d", resolverIP, "-j", postroutingchain)
	if err == nil ***REMOVED***
		iptables.RawCombinedOutputNative("-t", "nat", "-F", postroutingchain)
	***REMOVED*** else ***REMOVED***
		iptables.RawCombinedOutputNative("-t", "nat", "-N", postroutingchain)
		iptables.RawCombinedOutputNative("-t", "nat", "-I", "POSTROUTING", "-d", resolverIP, "-j", postroutingchain)
	***REMOVED***

	for _, rule := range rules ***REMOVED***
		if iptables.RawCombinedOutputNative(rule...) != nil ***REMOVED***
			logrus.Errorf("setting up rule failed, %v", rule)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *resolver) setupIPTable() error ***REMOVED***
	if r.err != nil ***REMOVED***
		return r.err
	***REMOVED***
	laddr := r.conn.LocalAddr().String()
	ltcpaddr := r.tcpListen.Addr().String()

	cmd := &exec.Cmd***REMOVED***
		Path:   reexec.Self(),
		Args:   append([]string***REMOVED***"setup-resolver"***REMOVED***, r.resolverKey, laddr, ltcpaddr),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	***REMOVED***
	if err := cmd.Run(); err != nil ***REMOVED***
		return fmt.Errorf("reexec failed: %v", err)
	***REMOVED***
	return nil
***REMOVED***
