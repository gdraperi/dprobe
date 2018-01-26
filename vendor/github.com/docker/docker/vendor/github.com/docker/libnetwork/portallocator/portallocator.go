package portallocator

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

const (
	// DefaultPortRangeStart indicates the first port in port range
	DefaultPortRangeStart = 49153
	// DefaultPortRangeEnd indicates the last port in port range
	DefaultPortRangeEnd = 65535
)

type ipMapping map[string]protoMap

var (
	// ErrAllPortsAllocated is returned when no more ports are available
	ErrAllPortsAllocated = errors.New("all ports are allocated")
	// ErrUnknownProtocol is returned when an unknown protocol was specified
	ErrUnknownProtocol = errors.New("unknown protocol")
	defaultIP          = net.ParseIP("0.0.0.0")
	once               sync.Once
	instance           *PortAllocator
	createInstance     = func() ***REMOVED*** instance = newInstance() ***REMOVED***
)

// ErrPortAlreadyAllocated is the returned error information when a requested port is already being used
type ErrPortAlreadyAllocated struct ***REMOVED***
	ip   string
	port int
***REMOVED***

func newErrPortAlreadyAllocated(ip string, port int) ErrPortAlreadyAllocated ***REMOVED***
	return ErrPortAlreadyAllocated***REMOVED***
		ip:   ip,
		port: port,
	***REMOVED***
***REMOVED***

// IP returns the address to which the used port is associated
func (e ErrPortAlreadyAllocated) IP() string ***REMOVED***
	return e.ip
***REMOVED***

// Port returns the value of the already used port
func (e ErrPortAlreadyAllocated) Port() int ***REMOVED***
	return e.port
***REMOVED***

// IPPort returns the address and the port in the form ip:port
func (e ErrPortAlreadyAllocated) IPPort() string ***REMOVED***
	return fmt.Sprintf("%s:%d", e.ip, e.port)
***REMOVED***

// Error is the implementation of error.Error interface
func (e ErrPortAlreadyAllocated) Error() string ***REMOVED***
	return fmt.Sprintf("Bind for %s:%d failed: port is already allocated", e.ip, e.port)
***REMOVED***

type (
	// PortAllocator manages the transport ports database
	PortAllocator struct ***REMOVED***
		mutex sync.Mutex
		ipMap ipMapping
		Begin int
		End   int
	***REMOVED***
	portRange struct ***REMOVED***
		begin int
		end   int
		last  int
	***REMOVED***
	portMap struct ***REMOVED***
		p            map[int]struct***REMOVED******REMOVED***
		defaultRange string
		portRanges   map[string]*portRange
	***REMOVED***
	protoMap map[string]*portMap
)

// Get returns the default instance of PortAllocator
func Get() *PortAllocator ***REMOVED***
	// Port Allocator is a singleton
	// Note: Long term solution will be each PortAllocator will have access to
	// the OS so that it can have up to date view of the OS port allocation.
	// When this happens singleton behavior will be removed. Clients do not
	// need to worry about this, they will not see a change in behavior.
	once.Do(createInstance)
	return instance
***REMOVED***

func newInstance() *PortAllocator ***REMOVED***
	start, end, err := getDynamicPortRange()
	if err != nil ***REMOVED***
		start, end = DefaultPortRangeStart, DefaultPortRangeEnd
	***REMOVED***
	return &PortAllocator***REMOVED***
		ipMap: ipMapping***REMOVED******REMOVED***,
		Begin: start,
		End:   end,
	***REMOVED***
***REMOVED***

// RequestPort requests new port from global ports pool for specified ip and proto.
// If port is 0 it returns first free port. Otherwise it checks port availability
// in proto's pool and returns that port or error if port is already busy.
func (p *PortAllocator) RequestPort(ip net.IP, proto string, port int) (int, error) ***REMOVED***
	return p.RequestPortInRange(ip, proto, port, port)
***REMOVED***

// RequestPortInRange requests new port from global ports pool for specified ip and proto.
// If portStart and portEnd are 0 it returns the first free port in the default ephemeral range.
// If portStart != portEnd it returns the first free port in the requested range.
// Otherwise (portStart == portEnd) it checks port availability in the requested proto's port-pool
// and returns that port or error if port is already busy.
func (p *PortAllocator) RequestPortInRange(ip net.IP, proto string, portStart, portEnd int) (int, error) ***REMOVED***
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if proto != "tcp" && proto != "udp" ***REMOVED***
		return 0, ErrUnknownProtocol
	***REMOVED***

	if ip == nil ***REMOVED***
		ip = defaultIP
	***REMOVED***
	ipstr := ip.String()
	protomap, ok := p.ipMap[ipstr]
	if !ok ***REMOVED***
		protomap = protoMap***REMOVED***
			"tcp": p.newPortMap(),
			"udp": p.newPortMap(),
		***REMOVED***

		p.ipMap[ipstr] = protomap
	***REMOVED***
	mapping := protomap[proto]
	if portStart > 0 && portStart == portEnd ***REMOVED***
		if _, ok := mapping.p[portStart]; !ok ***REMOVED***
			mapping.p[portStart] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			return portStart, nil
		***REMOVED***
		return 0, newErrPortAlreadyAllocated(ipstr, portStart)
	***REMOVED***

	port, err := mapping.findPort(portStart, portEnd)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return port, nil
***REMOVED***

// ReleasePort releases port from global ports pool for specified ip and proto.
func (p *PortAllocator) ReleasePort(ip net.IP, proto string, port int) error ***REMOVED***
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if ip == nil ***REMOVED***
		ip = defaultIP
	***REMOVED***
	protomap, ok := p.ipMap[ip.String()]
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	delete(protomap[proto].p, port)
	return nil
***REMOVED***

func (p *PortAllocator) newPortMap() *portMap ***REMOVED***
	defaultKey := getRangeKey(p.Begin, p.End)
	pm := &portMap***REMOVED***
		p:            map[int]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		defaultRange: defaultKey,
		portRanges: map[string]*portRange***REMOVED***
			defaultKey: newPortRange(p.Begin, p.End),
		***REMOVED***,
	***REMOVED***
	return pm
***REMOVED***

// ReleaseAll releases all ports for all ips.
func (p *PortAllocator) ReleaseAll() error ***REMOVED***
	p.mutex.Lock()
	p.ipMap = ipMapping***REMOVED******REMOVED***
	p.mutex.Unlock()
	return nil
***REMOVED***

func getRangeKey(portStart, portEnd int) string ***REMOVED***
	return fmt.Sprintf("%d-%d", portStart, portEnd)
***REMOVED***

func newPortRange(portStart, portEnd int) *portRange ***REMOVED***
	return &portRange***REMOVED***
		begin: portStart,
		end:   portEnd,
		last:  portEnd,
	***REMOVED***
***REMOVED***

func (pm *portMap) getPortRange(portStart, portEnd int) (*portRange, error) ***REMOVED***
	var key string
	if portStart == 0 && portEnd == 0 ***REMOVED***
		key = pm.defaultRange
	***REMOVED*** else ***REMOVED***
		key = getRangeKey(portStart, portEnd)
		if portStart == portEnd ||
			portStart == 0 || portEnd == 0 ||
			portEnd < portStart ***REMOVED***
			return nil, fmt.Errorf("invalid port range: %s", key)
		***REMOVED***
	***REMOVED***

	// Return existing port range, if already known.
	if pr, exists := pm.portRanges[key]; exists ***REMOVED***
		return pr, nil
	***REMOVED***

	// Otherwise create a new port range.
	pr := newPortRange(portStart, portEnd)
	pm.portRanges[key] = pr
	return pr, nil
***REMOVED***

func (pm *portMap) findPort(portStart, portEnd int) (int, error) ***REMOVED***
	pr, err := pm.getPortRange(portStart, portEnd)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	port := pr.last

	for i := 0; i <= pr.end-pr.begin; i++ ***REMOVED***
		port++
		if port > pr.end ***REMOVED***
			port = pr.begin
		***REMOVED***

		if _, ok := pm.p[port]; !ok ***REMOVED***
			pm.p[port] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			pr.last = port
			return port, nil
		***REMOVED***
	***REMOVED***
	return 0, ErrAllPortsAllocated
***REMOVED***
