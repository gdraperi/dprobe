package libnetwork

import (
	"net"

	"github.com/Microsoft/hcsshim"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)

type policyLists struct ***REMOVED***
	ilb *hcsshim.PolicyList
	elb *hcsshim.PolicyList
***REMOVED***

var lbPolicylistMap map[*loadBalancer]*policyLists

func init() ***REMOVED***
	lbPolicylistMap = make(map[*loadBalancer]*policyLists)
***REMOVED***

func (n *network) addLBBackend(ip, vip net.IP, lb *loadBalancer, ingressPorts []*PortConfig) ***REMOVED***

	if system.GetOSVersion().Build > 16236 ***REMOVED***
		lb.Lock()
		defer lb.Unlock()
		//find the load balancer IP for the network.
		var sourceVIP string
		for _, e := range n.Endpoints() ***REMOVED***
			epInfo := e.Info()
			if epInfo == nil ***REMOVED***
				continue
			***REMOVED***
			if epInfo.LoadBalancer() ***REMOVED***
				sourceVIP = epInfo.Iface().Address().IP.String()
				break
			***REMOVED***
		***REMOVED***

		if sourceVIP == "" ***REMOVED***
			logrus.Errorf("Failed to find load balancer IP for network %s", n.Name())
			return
		***REMOVED***

		var endpoints []hcsshim.HNSEndpoint

		for eid := range lb.backEnds ***REMOVED***
			//Call HNS to get back ID (GUID) corresponding to the endpoint.
			hnsEndpoint, err := hcsshim.GetHNSEndpointByName(eid)
			if err != nil ***REMOVED***
				logrus.Errorf("Failed to find HNS ID for endpoint %v: %v", eid, err)
				return
			***REMOVED***

			endpoints = append(endpoints, *hnsEndpoint)
		***REMOVED***

		if policies, ok := lbPolicylistMap[lb]; ok ***REMOVED***

			if policies.ilb != nil ***REMOVED***
				policies.ilb.Delete()
				policies.ilb = nil
			***REMOVED***

			if policies.elb != nil ***REMOVED***
				policies.elb.Delete()
				policies.elb = nil
			***REMOVED***
			delete(lbPolicylistMap, lb)
		***REMOVED***

		ilbPolicy, err := hcsshim.AddLoadBalancer(endpoints, true, sourceVIP, vip.String(), 0, 0, 0)
		if err != nil ***REMOVED***
			logrus.Errorf("Failed to add ILB policy for service %s (%s) with endpoints %v using load balancer IP %s on network %s: %v",
				lb.service.name, vip.String(), endpoints, sourceVIP, n.Name(), err)
			return
		***REMOVED***

		lbPolicylistMap[lb] = &policyLists***REMOVED***
			ilb: ilbPolicy,
		***REMOVED***

		publishedPorts := make(map[uint32]uint32)

		for i, port := range ingressPorts ***REMOVED***
			protocol := uint16(6)

			// Skip already published port
			if publishedPorts[port.PublishedPort] == port.TargetPort ***REMOVED***
				continue
			***REMOVED***

			if port.Protocol == ProtocolUDP ***REMOVED***
				protocol = 17
			***REMOVED***

			// check if already has udp matching to add wild card publishing
			for j := i + 1; j < len(ingressPorts); j++ ***REMOVED***
				if ingressPorts[j].TargetPort == port.TargetPort &&
					ingressPorts[j].PublishedPort == port.PublishedPort ***REMOVED***
					protocol = 0
				***REMOVED***
			***REMOVED***

			publishedPorts[port.PublishedPort] = port.TargetPort

			lbPolicylistMap[lb].elb, err = hcsshim.AddLoadBalancer(endpoints, false, sourceVIP, "", protocol, uint16(port.TargetPort), uint16(port.PublishedPort))
			if err != nil ***REMOVED***
				logrus.Errorf("Failed to add ELB policy for service %s (ip:%s target port:%v published port:%v) with endpoints %v using load balancer IP %s on network %s: %v",
					lb.service.name, vip.String(), uint16(port.TargetPort), uint16(port.PublishedPort), endpoints, sourceVIP, n.Name(), err)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (n *network) rmLBBackend(ip, vip net.IP, lb *loadBalancer, ingressPorts []*PortConfig, rmService bool) ***REMOVED***
	if system.GetOSVersion().Build > 16236 ***REMOVED***
		if len(lb.backEnds) > 0 ***REMOVED***
			//Reprogram HNS (actually VFP) with the existing backends.
			n.addLBBackend(ip, vip, lb, ingressPorts)
		***REMOVED*** else ***REMOVED***
			lb.Lock()
			defer lb.Unlock()
			logrus.Debugf("No more backends for service %s (ip:%s).  Removing all policies", lb.service.name, lb.vip.String())

			if policyLists, ok := lbPolicylistMap[lb]; ok ***REMOVED***
				if policyLists.ilb != nil ***REMOVED***
					policyLists.ilb.Delete()
					policyLists.ilb = nil
				***REMOVED***

				if policyLists.elb != nil ***REMOVED***
					policyLists.elb.Delete()
					policyLists.elb = nil
				***REMOVED***
				delete(lbPolicylistMap, lb)

			***REMOVED*** else ***REMOVED***
				logrus.Errorf("Failed to find policies for service %s (%s)", lb.service.name, lb.vip.String())
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (sb *sandbox) populateLoadbalancers(ep *endpoint) ***REMOVED***
***REMOVED***

func arrangeIngressFilterRule() ***REMOVED***
***REMOVED***
