package libnetwork

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/libnetwork/iptables"
	"github.com/docker/libnetwork/ipvs"
	"github.com/docker/libnetwork/ns"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink/nl"
	"github.com/vishvananda/netns"
)

func init() ***REMOVED***
	reexec.Register("fwmarker", fwMarker)
	reexec.Register("redirecter", redirecter)
***REMOVED***

// Get all loadbalancers on this network that is currently discovered
// on this node.
func (n *network) connectedLoadbalancers() []*loadBalancer ***REMOVED***
	c := n.getController()

	c.Lock()
	serviceBindings := make([]*service, 0, len(c.serviceBindings))
	for _, s := range c.serviceBindings ***REMOVED***
		serviceBindings = append(serviceBindings, s)
	***REMOVED***
	c.Unlock()

	var lbs []*loadBalancer
	for _, s := range serviceBindings ***REMOVED***
		s.Lock()
		// Skip the serviceBindings that got deleted
		if s.deleted ***REMOVED***
			s.Unlock()
			continue
		***REMOVED***
		if lb, ok := s.loadBalancers[n.ID()]; ok ***REMOVED***
			lbs = append(lbs, lb)
		***REMOVED***
		s.Unlock()
	***REMOVED***

	return lbs
***REMOVED***

// Populate all loadbalancers on the network that the passed endpoint
// belongs to, into this sandbox.
func (sb *sandbox) populateLoadbalancers(ep *endpoint) ***REMOVED***
	var gwIP net.IP

	// This is an interface less endpoint. Nothing to do.
	if ep.Iface() == nil ***REMOVED***
		return
	***REMOVED***

	n := ep.getNetwork()
	eIP := ep.Iface().Address()

	if n.ingress ***REMOVED***
		if err := addRedirectRules(sb.Key(), eIP, ep.ingressPorts); err != nil ***REMOVED***
			logrus.Errorf("Failed to add redirect rules for ep %s (%s): %v", ep.Name(), ep.ID()[0:7], err)
		***REMOVED***
	***REMOVED***

	if sb.ingress ***REMOVED***
		// For the ingress sandbox if this is not gateway
		// endpoint do nothing.
		if ep != sb.getGatewayEndpoint() ***REMOVED***
			return
		***REMOVED***

		// This is the gateway endpoint. Now get the ingress
		// network and plumb the loadbalancers.
		gwIP = ep.Iface().Address().IP
		for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
			if !ep.endpointInGWNetwork() ***REMOVED***
				n = ep.getNetwork()
				eIP = ep.Iface().Address()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, lb := range n.connectedLoadbalancers() ***REMOVED***
		// Skip if vip is not valid.
		if len(lb.vip) == 0 ***REMOVED***
			continue
		***REMOVED***

		lb.service.Lock()
		for _, ip := range lb.backEnds ***REMOVED***
			sb.addLBBackend(ip, lb.vip, lb.fwMark, lb.service.ingressPorts, eIP, gwIP, n.ingress)
		***REMOVED***
		lb.service.Unlock()
	***REMOVED***
***REMOVED***

// Add loadbalancer backend to all sandboxes which has a connection to
// this network. If needed add the service as well.
func (n *network) addLBBackend(ip, vip net.IP, lb *loadBalancer, ingressPorts []*PortConfig) ***REMOVED***
	n.WalkEndpoints(func(e Endpoint) bool ***REMOVED***
		ep := e.(*endpoint)
		if sb, ok := ep.getSandbox(); ok ***REMOVED***
			if !sb.isEndpointPopulated(ep) ***REMOVED***
				return false
			***REMOVED***

			var gwIP net.IP
			if ep := sb.getGatewayEndpoint(); ep != nil ***REMOVED***
				gwIP = ep.Iface().Address().IP
			***REMOVED***

			sb.addLBBackend(ip, vip, lb.fwMark, ingressPorts, ep.Iface().Address(), gwIP, n.ingress)
		***REMOVED***

		return false
	***REMOVED***)
***REMOVED***

// Remove loadbalancer backend from all sandboxes which has a
// connection to this network. If needed remove the service entry as
// well, as specified by the rmService bool.
func (n *network) rmLBBackend(ip, vip net.IP, lb *loadBalancer, ingressPorts []*PortConfig, rmService bool) ***REMOVED***
	n.WalkEndpoints(func(e Endpoint) bool ***REMOVED***
		ep := e.(*endpoint)
		if sb, ok := ep.getSandbox(); ok ***REMOVED***
			if !sb.isEndpointPopulated(ep) ***REMOVED***
				return false
			***REMOVED***

			var gwIP net.IP
			if ep := sb.getGatewayEndpoint(); ep != nil ***REMOVED***
				gwIP = ep.Iface().Address().IP
			***REMOVED***

			sb.rmLBBackend(ip, vip, lb.fwMark, ingressPorts, ep.Iface().Address(), gwIP, rmService, n.ingress)
		***REMOVED***

		return false
	***REMOVED***)
***REMOVED***

// Add loadbalancer backend into one connected sandbox.
func (sb *sandbox) addLBBackend(ip, vip net.IP, fwMark uint32, ingressPorts []*PortConfig, eIP *net.IPNet, gwIP net.IP, isIngressNetwork bool) ***REMOVED***
	if sb.osSbox == nil ***REMOVED***
		return
	***REMOVED***

	if isIngressNetwork && !sb.ingress ***REMOVED***
		return
	***REMOVED***

	i, err := ipvs.New(sb.Key())
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to create an ipvs handle for sbox %s (%s,%s) for lb addition: %v", sb.ID()[0:7], sb.ContainerID()[0:7], sb.Key(), err)
		return
	***REMOVED***
	defer i.Close()

	s := &ipvs.Service***REMOVED***
		AddressFamily: nl.FAMILY_V4,
		FWMark:        fwMark,
		SchedName:     ipvs.RoundRobin,
	***REMOVED***

	if !i.IsServicePresent(s) ***REMOVED***
		var filteredPorts []*PortConfig
		if sb.ingress ***REMOVED***
			filteredPorts = filterPortConfigs(ingressPorts, false)
			if err := programIngress(gwIP, filteredPorts, false); err != nil ***REMOVED***
				logrus.Errorf("Failed to add ingress: %v", err)
				return
			***REMOVED***
		***REMOVED***

		logrus.Debugf("Creating service for vip %s fwMark %d ingressPorts %#v in sbox %s (%s)", vip, fwMark, ingressPorts, sb.ID()[0:7], sb.ContainerID()[0:7])
		if err := invokeFWMarker(sb.Key(), vip, fwMark, ingressPorts, eIP, false); err != nil ***REMOVED***
			logrus.Errorf("Failed to add firewall mark rule in sbox %s (%s): %v", sb.ID()[0:7], sb.ContainerID()[0:7], err)
			return
		***REMOVED***

		if err := i.NewService(s); err != nil && err != syscall.EEXIST ***REMOVED***
			logrus.Errorf("Failed to create a new service for vip %s fwmark %d in sbox %s (%s): %v", vip, fwMark, sb.ID()[0:7], sb.ContainerID()[0:7], err)
			return
		***REMOVED***
	***REMOVED***

	d := &ipvs.Destination***REMOVED***
		AddressFamily: nl.FAMILY_V4,
		Address:       ip,
		Weight:        1,
	***REMOVED***

	// Remove the sched name before using the service to add
	// destination.
	s.SchedName = ""
	if err := i.NewDestination(s, d); err != nil && err != syscall.EEXIST ***REMOVED***
		logrus.Errorf("Failed to create real server %s for vip %s fwmark %d in sbox %s (%s): %v", ip, vip, fwMark, sb.ID()[0:7], sb.ContainerID()[0:7], err)
	***REMOVED***
***REMOVED***

// Remove loadbalancer backend from one connected sandbox.
func (sb *sandbox) rmLBBackend(ip, vip net.IP, fwMark uint32, ingressPorts []*PortConfig, eIP *net.IPNet, gwIP net.IP, rmService bool, isIngressNetwork bool) ***REMOVED***
	if sb.osSbox == nil ***REMOVED***
		return
	***REMOVED***

	if isIngressNetwork && !sb.ingress ***REMOVED***
		return
	***REMOVED***

	i, err := ipvs.New(sb.Key())
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to create an ipvs handle for sbox %s (%s,%s) for lb removal: %v", sb.ID()[0:7], sb.ContainerID()[0:7], sb.Key(), err)
		return
	***REMOVED***
	defer i.Close()

	s := &ipvs.Service***REMOVED***
		AddressFamily: nl.FAMILY_V4,
		FWMark:        fwMark,
	***REMOVED***

	d := &ipvs.Destination***REMOVED***
		AddressFamily: nl.FAMILY_V4,
		Address:       ip,
		Weight:        1,
	***REMOVED***

	if err := i.DelDestination(s, d); err != nil && err != syscall.ENOENT ***REMOVED***
		logrus.Errorf("Failed to delete real server %s for vip %s fwmark %d in sbox %s (%s): %v", ip, vip, fwMark, sb.ID()[0:7], sb.ContainerID()[0:7], err)
	***REMOVED***

	if rmService ***REMOVED***
		s.SchedName = ipvs.RoundRobin
		if err := i.DelService(s); err != nil && err != syscall.ENOENT ***REMOVED***
			logrus.Errorf("Failed to delete service for vip %s fwmark %d in sbox %s (%s): %v", vip, fwMark, sb.ID()[0:7], sb.ContainerID()[0:7], err)
		***REMOVED***

		var filteredPorts []*PortConfig
		if sb.ingress ***REMOVED***
			filteredPorts = filterPortConfigs(ingressPorts, true)
			if err := programIngress(gwIP, filteredPorts, true); err != nil ***REMOVED***
				logrus.Errorf("Failed to delete ingress: %v", err)
			***REMOVED***
		***REMOVED***

		if err := invokeFWMarker(sb.Key(), vip, fwMark, ingressPorts, eIP, true); err != nil ***REMOVED***
			logrus.Errorf("Failed to delete firewall mark rule in sbox %s (%s): %v", sb.ID()[0:7], sb.ContainerID()[0:7], err)
		***REMOVED***
	***REMOVED***
***REMOVED***

const ingressChain = "DOCKER-INGRESS"

var (
	ingressOnce     sync.Once
	ingressProxyMu  sync.Mutex
	ingressProxyTbl = make(map[string]io.Closer)
	portConfigMu    sync.Mutex
	portConfigTbl   = make(map[PortConfig]int)
)

func filterPortConfigs(ingressPorts []*PortConfig, isDelete bool) []*PortConfig ***REMOVED***
	portConfigMu.Lock()
	iPorts := make([]*PortConfig, 0, len(ingressPorts))
	for _, pc := range ingressPorts ***REMOVED***
		if isDelete ***REMOVED***
			if cnt, ok := portConfigTbl[*pc]; ok ***REMOVED***
				// This is the last reference to this
				// port config. Delete the port config
				// and add it to filtered list to be
				// plumbed.
				if cnt == 1 ***REMOVED***
					delete(portConfigTbl, *pc)
					iPorts = append(iPorts, pc)
					continue
				***REMOVED***

				portConfigTbl[*pc] = cnt - 1
			***REMOVED***

			continue
		***REMOVED***

		if cnt, ok := portConfigTbl[*pc]; ok ***REMOVED***
			portConfigTbl[*pc] = cnt + 1
			continue
		***REMOVED***

		// We are adding it for the first time. Add it to the
		// filter list to be plumbed.
		portConfigTbl[*pc] = 1
		iPorts = append(iPorts, pc)
	***REMOVED***
	portConfigMu.Unlock()

	return iPorts
***REMOVED***

func programIngress(gwIP net.IP, ingressPorts []*PortConfig, isDelete bool) error ***REMOVED***
	addDelOpt := "-I"
	if isDelete ***REMOVED***
		addDelOpt = "-D"
	***REMOVED***

	chainExists := iptables.ExistChain(ingressChain, iptables.Nat)
	filterChainExists := iptables.ExistChain(ingressChain, iptables.Filter)

	ingressOnce.Do(func() ***REMOVED***
		// Flush nat table and filter table ingress chain rules during init if it
		// exists. It might contain stale rules from previous life.
		if chainExists ***REMOVED***
			if err := iptables.RawCombinedOutput("-t", "nat", "-F", ingressChain); err != nil ***REMOVED***
				logrus.Errorf("Could not flush nat table ingress chain rules during init: %v", err)
			***REMOVED***
		***REMOVED***
		if filterChainExists ***REMOVED***
			if err := iptables.RawCombinedOutput("-F", ingressChain); err != nil ***REMOVED***
				logrus.Errorf("Could not flush filter table ingress chain rules during init: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	if !isDelete ***REMOVED***
		if !chainExists ***REMOVED***
			if err := iptables.RawCombinedOutput("-t", "nat", "-N", ingressChain); err != nil ***REMOVED***
				return fmt.Errorf("failed to create ingress chain: %v", err)
			***REMOVED***
		***REMOVED***
		if !filterChainExists ***REMOVED***
			if err := iptables.RawCombinedOutput("-N", ingressChain); err != nil ***REMOVED***
				return fmt.Errorf("failed to create filter table ingress chain: %v", err)
			***REMOVED***
		***REMOVED***

		if !iptables.Exists(iptables.Nat, ingressChain, "-j", "RETURN") ***REMOVED***
			if err := iptables.RawCombinedOutput("-t", "nat", "-A", ingressChain, "-j", "RETURN"); err != nil ***REMOVED***
				return fmt.Errorf("failed to add return rule in nat table ingress chain: %v", err)
			***REMOVED***
		***REMOVED***

		if !iptables.Exists(iptables.Filter, ingressChain, "-j", "RETURN") ***REMOVED***
			if err := iptables.RawCombinedOutput("-A", ingressChain, "-j", "RETURN"); err != nil ***REMOVED***
				return fmt.Errorf("failed to add return rule to filter table ingress chain: %v", err)
			***REMOVED***
		***REMOVED***

		for _, chain := range []string***REMOVED***"OUTPUT", "PREROUTING"***REMOVED*** ***REMOVED***
			if !iptables.Exists(iptables.Nat, chain, "-m", "addrtype", "--dst-type", "LOCAL", "-j", ingressChain) ***REMOVED***
				if err := iptables.RawCombinedOutput("-t", "nat", "-I", chain, "-m", "addrtype", "--dst-type", "LOCAL", "-j", ingressChain); err != nil ***REMOVED***
					return fmt.Errorf("failed to add jump rule in %s to ingress chain: %v", chain, err)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if !iptables.Exists(iptables.Filter, "FORWARD", "-j", ingressChain) ***REMOVED***
			if err := iptables.RawCombinedOutput("-I", "FORWARD", "-j", ingressChain); err != nil ***REMOVED***
				return fmt.Errorf("failed to add jump rule to %s in filter table forward chain: %v", ingressChain, err)
			***REMOVED***
			arrangeUserFilterRule()
		***REMOVED***

		oifName, err := findOIFName(gwIP)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to find gateway bridge interface name for %s: %v", gwIP, err)
		***REMOVED***

		path := filepath.Join("/proc/sys/net/ipv4/conf", oifName, "route_localnet")
		if err := ioutil.WriteFile(path, []byte***REMOVED***'1', '\n'***REMOVED***, 0644); err != nil ***REMOVED***
			return fmt.Errorf("could not write to %s: %v", path, err)
		***REMOVED***

		ruleArgs := strings.Fields(fmt.Sprintf("-m addrtype --src-type LOCAL -o %s -j MASQUERADE", oifName))
		if !iptables.Exists(iptables.Nat, "POSTROUTING", ruleArgs...) ***REMOVED***
			if err := iptables.RawCombinedOutput(append([]string***REMOVED***"-t", "nat", "-I", "POSTROUTING"***REMOVED***, ruleArgs...)...); err != nil ***REMOVED***
				return fmt.Errorf("failed to add ingress localhost POSTROUTING rule for %s: %v", oifName, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, iPort := range ingressPorts ***REMOVED***
		if iptables.ExistChain(ingressChain, iptables.Nat) ***REMOVED***
			rule := strings.Fields(fmt.Sprintf("-t nat %s %s -p %s --dport %d -j DNAT --to-destination %s:%d",
				addDelOpt, ingressChain, strings.ToLower(PortConfig_Protocol_name[int32(iPort.Protocol)]), iPort.PublishedPort, gwIP, iPort.PublishedPort))
			if err := iptables.RawCombinedOutput(rule...); err != nil ***REMOVED***
				errStr := fmt.Sprintf("setting up rule failed, %v: %v", rule, err)
				if !isDelete ***REMOVED***
					return fmt.Errorf("%s", errStr)
				***REMOVED***

				logrus.Infof("%s", errStr)
			***REMOVED***
		***REMOVED***

		// Filter table rules to allow a published service to be accessible in the local node from..
		// 1) service tasks attached to other networks
		// 2) unmanaged containers on bridge networks
		rule := strings.Fields(fmt.Sprintf("%s %s -m state -p %s --sport %d --state ESTABLISHED,RELATED -j ACCEPT",
			addDelOpt, ingressChain, strings.ToLower(PortConfig_Protocol_name[int32(iPort.Protocol)]), iPort.PublishedPort))
		if err := iptables.RawCombinedOutput(rule...); err != nil ***REMOVED***
			errStr := fmt.Sprintf("setting up rule failed, %v: %v", rule, err)
			if !isDelete ***REMOVED***
				return fmt.Errorf("%s", errStr)
			***REMOVED***
			logrus.Warnf("%s", errStr)
		***REMOVED***

		rule = strings.Fields(fmt.Sprintf("%s %s -p %s --dport %d -j ACCEPT",
			addDelOpt, ingressChain, strings.ToLower(PortConfig_Protocol_name[int32(iPort.Protocol)]), iPort.PublishedPort))
		if err := iptables.RawCombinedOutput(rule...); err != nil ***REMOVED***
			errStr := fmt.Sprintf("setting up rule failed, %v: %v", rule, err)
			if !isDelete ***REMOVED***
				return fmt.Errorf("%s", errStr)
			***REMOVED***

			logrus.Warnf("%s", errStr)
		***REMOVED***

		if err := plumbProxy(iPort, isDelete); err != nil ***REMOVED***
			logrus.Warnf("failed to create proxy for port %d: %v", iPort.PublishedPort, err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// In the filter table FORWARD chain the first rule should be to jump to
// DOCKER-USER so the user is able to filter packet first.
// The second rule should be jump to INGRESS-CHAIN.
// This chain has the rules to allow access to the published ports for swarm tasks
// from local bridge networks and docker_gwbridge (ie:taks on other swarm netwroks)
func arrangeIngressFilterRule() ***REMOVED***
	if iptables.ExistChain(ingressChain, iptables.Filter) ***REMOVED***
		if iptables.Exists(iptables.Filter, "FORWARD", "-j", ingressChain) ***REMOVED***
			if err := iptables.RawCombinedOutput("-D", "FORWARD", "-j", ingressChain); err != nil ***REMOVED***
				logrus.Warnf("failed to delete jump rule to ingressChain in filter table: %v", err)
			***REMOVED***
		***REMOVED***
		if err := iptables.RawCombinedOutput("-I", "FORWARD", "-j", ingressChain); err != nil ***REMOVED***
			logrus.Warnf("failed to add jump rule to ingressChain in filter table: %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func findOIFName(ip net.IP) (string, error) ***REMOVED***
	nlh := ns.NlHandle()

	routes, err := nlh.RouteGet(ip)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if len(routes) == 0 ***REMOVED***
		return "", fmt.Errorf("no route to %s", ip)
	***REMOVED***

	// Pick the first route(typically there is only one route). We
	// don't support multipath.
	link, err := nlh.LinkByIndex(routes[0].LinkIndex)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return link.Attrs().Name, nil
***REMOVED***

func plumbProxy(iPort *PortConfig, isDelete bool) error ***REMOVED***
	var (
		err error
		l   io.Closer
	)

	portSpec := fmt.Sprintf("%d/%s", iPort.PublishedPort, strings.ToLower(PortConfig_Protocol_name[int32(iPort.Protocol)]))
	if isDelete ***REMOVED***
		ingressProxyMu.Lock()
		if listener, ok := ingressProxyTbl[portSpec]; ok ***REMOVED***
			if listener != nil ***REMOVED***
				listener.Close()
			***REMOVED***
		***REMOVED***
		ingressProxyMu.Unlock()

		return nil
	***REMOVED***

	switch iPort.Protocol ***REMOVED***
	case ProtocolTCP:
		l, err = net.ListenTCP("tcp", &net.TCPAddr***REMOVED***Port: int(iPort.PublishedPort)***REMOVED***)
	case ProtocolUDP:
		l, err = net.ListenUDP("udp", &net.UDPAddr***REMOVED***Port: int(iPort.PublishedPort)***REMOVED***)
	***REMOVED***

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	ingressProxyMu.Lock()
	ingressProxyTbl[portSpec] = l
	ingressProxyMu.Unlock()

	return nil
***REMOVED***

func writePortsToFile(ports []*PortConfig) (string, error) ***REMOVED***
	f, err := ioutil.TempFile("", "port_configs")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer f.Close()

	buf, _ := proto.Marshal(&EndpointRecord***REMOVED***
		IngressPorts: ports,
	***REMOVED***)

	n, err := f.Write(buf)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if n < len(buf) ***REMOVED***
		return "", io.ErrShortWrite
	***REMOVED***

	return f.Name(), nil
***REMOVED***

func readPortsFromFile(fileName string) ([]*PortConfig, error) ***REMOVED***
	buf, err := ioutil.ReadFile(fileName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var epRec EndpointRecord
	err = proto.Unmarshal(buf, &epRec)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return epRec.IngressPorts, nil
***REMOVED***

// Invoke fwmarker reexec routine to mark vip destined packets with
// the passed firewall mark.
func invokeFWMarker(path string, vip net.IP, fwMark uint32, ingressPorts []*PortConfig, eIP *net.IPNet, isDelete bool) error ***REMOVED***
	var ingressPortsFile string

	if len(ingressPorts) != 0 ***REMOVED***
		var err error
		ingressPortsFile, err = writePortsToFile(ingressPorts)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		defer os.Remove(ingressPortsFile)
	***REMOVED***

	addDelOpt := "-A"
	if isDelete ***REMOVED***
		addDelOpt = "-D"
	***REMOVED***

	cmd := &exec.Cmd***REMOVED***
		Path:   reexec.Self(),
		Args:   append([]string***REMOVED***"fwmarker"***REMOVED***, path, vip.String(), fmt.Sprintf("%d", fwMark), addDelOpt, ingressPortsFile, eIP.String()),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	***REMOVED***

	if err := cmd.Run(); err != nil ***REMOVED***
		return fmt.Errorf("reexec failed: %v", err)
	***REMOVED***

	return nil
***REMOVED***

// Firewall marker reexec function.
func fwMarker() ***REMOVED***
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if len(os.Args) < 7 ***REMOVED***
		logrus.Error("invalid number of arguments..")
		os.Exit(1)
	***REMOVED***

	var ingressPorts []*PortConfig
	if os.Args[5] != "" ***REMOVED***
		var err error
		ingressPorts, err = readPortsFromFile(os.Args[5])
		if err != nil ***REMOVED***
			logrus.Errorf("Failed reading ingress ports file: %v", err)
			os.Exit(6)
		***REMOVED***
	***REMOVED***

	vip := os.Args[2]
	fwMark, err := strconv.ParseUint(os.Args[3], 10, 32)
	if err != nil ***REMOVED***
		logrus.Errorf("bad fwmark value(%s) passed: %v", os.Args[3], err)
		os.Exit(2)
	***REMOVED***
	addDelOpt := os.Args[4]

	rules := [][]string***REMOVED******REMOVED***
	for _, iPort := range ingressPorts ***REMOVED***
		rule := strings.Fields(fmt.Sprintf("-t mangle %s PREROUTING -p %s --dport %d -j MARK --set-mark %d",
			addDelOpt, strings.ToLower(PortConfig_Protocol_name[int32(iPort.Protocol)]), iPort.PublishedPort, fwMark))
		rules = append(rules, rule)
	***REMOVED***

	ns, err := netns.GetFromPath(os.Args[1])
	if err != nil ***REMOVED***
		logrus.Errorf("failed get network namespace %q: %v", os.Args[1], err)
		os.Exit(3)
	***REMOVED***
	defer ns.Close()

	if err := netns.Set(ns); err != nil ***REMOVED***
		logrus.Errorf("setting into container net ns %v failed, %v", os.Args[1], err)
		os.Exit(4)
	***REMOVED***

	if addDelOpt == "-A" ***REMOVED***
		eIP, subnet, err := net.ParseCIDR(os.Args[6])
		if err != nil ***REMOVED***
			logrus.Errorf("Failed to parse endpoint IP %s: %v", os.Args[6], err)
			os.Exit(9)
		***REMOVED***

		ruleParams := strings.Fields(fmt.Sprintf("-m ipvs --ipvs -d %s -j SNAT --to-source %s", subnet, eIP))
		if !iptables.Exists("nat", "POSTROUTING", ruleParams...) ***REMOVED***
			rule := append(strings.Fields("-t nat -A POSTROUTING"), ruleParams...)
			rules = append(rules, rule)

			err := ioutil.WriteFile("/proc/sys/net/ipv4/vs/conntrack", []byte***REMOVED***'1', '\n'***REMOVED***, 0644)
			if err != nil ***REMOVED***
				logrus.Errorf("Failed to write to /proc/sys/net/ipv4/vs/conntrack: %v", err)
				os.Exit(8)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	rule := strings.Fields(fmt.Sprintf("-t mangle %s OUTPUT -d %s/32 -j MARK --set-mark %d", addDelOpt, vip, fwMark))
	rules = append(rules, rule)

	rule = strings.Fields(fmt.Sprintf("-t nat %s OUTPUT -p icmp --icmp echo-request -d %s -j DNAT --to 127.0.0.1", addDelOpt, vip))
	rules = append(rules, rule)

	for _, rule := range rules ***REMOVED***
		if err := iptables.RawCombinedOutputNative(rule...); err != nil ***REMOVED***
			logrus.Errorf("setting up rule failed, %v: %v", rule, err)
			os.Exit(5)
		***REMOVED***
	***REMOVED***
***REMOVED***

func addRedirectRules(path string, eIP *net.IPNet, ingressPorts []*PortConfig) error ***REMOVED***
	var ingressPortsFile string

	if len(ingressPorts) != 0 ***REMOVED***
		var err error
		ingressPortsFile, err = writePortsToFile(ingressPorts)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer os.Remove(ingressPortsFile)
	***REMOVED***

	cmd := &exec.Cmd***REMOVED***
		Path:   reexec.Self(),
		Args:   append([]string***REMOVED***"redirecter"***REMOVED***, path, eIP.String(), ingressPortsFile),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	***REMOVED***

	if err := cmd.Run(); err != nil ***REMOVED***
		return fmt.Errorf("reexec failed: %v", err)
	***REMOVED***

	return nil
***REMOVED***

// Redirecter reexec function.
func redirecter() ***REMOVED***
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if len(os.Args) < 4 ***REMOVED***
		logrus.Error("invalid number of arguments..")
		os.Exit(1)
	***REMOVED***

	var ingressPorts []*PortConfig
	if os.Args[3] != "" ***REMOVED***
		var err error
		ingressPorts, err = readPortsFromFile(os.Args[3])
		if err != nil ***REMOVED***
			logrus.Errorf("Failed reading ingress ports file: %v", err)
			os.Exit(2)
		***REMOVED***
	***REMOVED***

	eIP, _, err := net.ParseCIDR(os.Args[2])
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to parse endpoint IP %s: %v", os.Args[2], err)
		os.Exit(3)
	***REMOVED***

	rules := [][]string***REMOVED******REMOVED***
	for _, iPort := range ingressPorts ***REMOVED***
		rule := strings.Fields(fmt.Sprintf("-t nat -A PREROUTING -d %s -p %s --dport %d -j REDIRECT --to-port %d",
			eIP.String(), strings.ToLower(PortConfig_Protocol_name[int32(iPort.Protocol)]), iPort.PublishedPort, iPort.TargetPort))
		rules = append(rules, rule)
		// Allow only incoming connections to exposed ports
		iRule := strings.Fields(fmt.Sprintf("-I INPUT -d %s -p %s --dport %d -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT",
			eIP.String(), strings.ToLower(PortConfig_Protocol_name[int32(iPort.Protocol)]), iPort.TargetPort))
		rules = append(rules, iRule)
		// Allow only outgoing connections from exposed ports
		oRule := strings.Fields(fmt.Sprintf("-I OUTPUT -s %s -p %s --sport %d -m conntrack --ctstate ESTABLISHED -j ACCEPT",
			eIP.String(), strings.ToLower(PortConfig_Protocol_name[int32(iPort.Protocol)]), iPort.TargetPort))
		rules = append(rules, oRule)
	***REMOVED***

	ns, err := netns.GetFromPath(os.Args[1])
	if err != nil ***REMOVED***
		logrus.Errorf("failed get network namespace %q: %v", os.Args[1], err)
		os.Exit(4)
	***REMOVED***
	defer ns.Close()

	if err := netns.Set(ns); err != nil ***REMOVED***
		logrus.Errorf("setting into container net ns %v failed, %v", os.Args[1], err)
		os.Exit(5)
	***REMOVED***

	for _, rule := range rules ***REMOVED***
		if err := iptables.RawCombinedOutputNative(rule...); err != nil ***REMOVED***
			logrus.Errorf("setting up rule failed, %v: %v", rule, err)
			os.Exit(6)
		***REMOVED***
	***REMOVED***

	if len(ingressPorts) == 0 ***REMOVED***
		return
	***REMOVED***

	// Ensure blocking rules for anything else in/to ingress network
	for _, rule := range [][]string***REMOVED***
		***REMOVED***"-d", eIP.String(), "-p", "udp", "-j", "DROP"***REMOVED***,
		***REMOVED***"-d", eIP.String(), "-p", "tcp", "-j", "DROP"***REMOVED***,
	***REMOVED*** ***REMOVED***
		if !iptables.ExistsNative(iptables.Filter, "INPUT", rule...) ***REMOVED***
			if err := iptables.RawCombinedOutputNative(append([]string***REMOVED***"-A", "INPUT"***REMOVED***, rule...)...); err != nil ***REMOVED***
				logrus.Errorf("setting up rule failed, %v: %v", rule, err)
				os.Exit(7)
			***REMOVED***
		***REMOVED***
		rule[0] = "-s"
		if !iptables.ExistsNative(iptables.Filter, "OUTPUT", rule...) ***REMOVED***
			if err := iptables.RawCombinedOutputNative(append([]string***REMOVED***"-A", "OUTPUT"***REMOVED***, rule...)...); err != nil ***REMOVED***
				logrus.Errorf("setting up rule failed, %v: %v", rule, err)
				os.Exit(8)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
