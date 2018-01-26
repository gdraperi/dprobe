package libnetwork

import (
	"fmt"
	"net"
	"sync"

	"github.com/docker/libnetwork/common"
)

var (
	// A global monotonic counter to assign firewall marks to
	// services.
	fwMarkCtr   uint32 = 256
	fwMarkCtrMu sync.Mutex
)

type portConfigs []*PortConfig

func (p portConfigs) String() string ***REMOVED***
	if len(p) == 0 ***REMOVED***
		return ""
	***REMOVED***

	pc := p[0]
	str := fmt.Sprintf("%d:%d/%s", pc.PublishedPort, pc.TargetPort, PortConfig_Protocol_name[int32(pc.Protocol)])
	for _, pc := range p[1:] ***REMOVED***
		str = str + fmt.Sprintf(",%d:%d/%s", pc.PublishedPort, pc.TargetPort, PortConfig_Protocol_name[int32(pc.Protocol)])
	***REMOVED***

	return str
***REMOVED***

type serviceKey struct ***REMOVED***
	id    string
	ports string
***REMOVED***

type service struct ***REMOVED***
	name string // Service Name
	id   string // Service ID

	// Map of loadbalancers for the service one-per attached
	// network. It is keyed with network ID.
	loadBalancers map[string]*loadBalancer

	// List of ingress ports exposed by the service
	ingressPorts portConfigs

	// Service aliases
	aliases []string

	// This maps tracks for each IP address the list of endpoints ID
	// associated with it. At stable state the endpoint ID expected is 1
	// but during transition and service change it is possible to have
	// temporary more than 1
	ipToEndpoint common.SetMatrix

	deleted bool

	sync.Mutex
***REMOVED***

// assignIPToEndpoint inserts the mapping between the IP and the endpoint identifier
// returns true if the mapping was not present, false otherwise
// returns also the number of endpoints associated to the IP
func (s *service) assignIPToEndpoint(ip, eID string) (bool, int) ***REMOVED***
	return s.ipToEndpoint.Insert(ip, eID)
***REMOVED***

// removeIPToEndpoint removes the mapping between the IP and the endpoint identifier
// returns true if the mapping was deleted, false otherwise
// returns also the number of endpoints associated to the IP
func (s *service) removeIPToEndpoint(ip, eID string) (bool, int) ***REMOVED***
	return s.ipToEndpoint.Remove(ip, eID)
***REMOVED***

func (s *service) printIPToEndpoint(ip string) (string, bool) ***REMOVED***
	return s.ipToEndpoint.String(ip)
***REMOVED***

type loadBalancer struct ***REMOVED***
	vip    net.IP
	fwMark uint32

	// Map of backend IPs backing this loadbalancer on this
	// network. It is keyed with endpoint ID.
	backEnds map[string]net.IP

	// Back pointer to service to which the loadbalancer belongs.
	service *service
	sync.Mutex
***REMOVED***
