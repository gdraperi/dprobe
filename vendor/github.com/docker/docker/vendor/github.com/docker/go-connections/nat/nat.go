// Package nat is a convenience package for manipulation of strings describing network ports.
package nat

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	// portSpecTemplate is the expected format for port specifications
	portSpecTemplate = "ip:hostPort:containerPort"
)

// PortBinding represents a binding between a Host IP address and a Host Port
type PortBinding struct ***REMOVED***
	// HostIP is the host IP Address
	HostIP string `json:"HostIp"`
	// HostPort is the host port number
	HostPort string
***REMOVED***

// PortMap is a collection of PortBinding indexed by Port
type PortMap map[Port][]PortBinding

// PortSet is a collection of structs indexed by Port
type PortSet map[Port]struct***REMOVED******REMOVED***

// Port is a string containing port number and protocol in the format "80/tcp"
type Port string

// NewPort creates a new instance of a Port given a protocol and port number or port range
func NewPort(proto, port string) (Port, error) ***REMOVED***
	// Check for parsing issues on "port" now so we can avoid having
	// to check it later on.

	portStartInt, portEndInt, err := ParsePortRangeToInt(port)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if portStartInt == portEndInt ***REMOVED***
		return Port(fmt.Sprintf("%d/%s", portStartInt, proto)), nil
	***REMOVED***
	return Port(fmt.Sprintf("%d-%d/%s", portStartInt, portEndInt, proto)), nil
***REMOVED***

// ParsePort parses the port number string and returns an int
func ParsePort(rawPort string) (int, error) ***REMOVED***
	if len(rawPort) == 0 ***REMOVED***
		return 0, nil
	***REMOVED***
	port, err := strconv.ParseUint(rawPort, 10, 16)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return int(port), nil
***REMOVED***

// ParsePortRangeToInt parses the port range string and returns start/end ints
func ParsePortRangeToInt(rawPort string) (int, int, error) ***REMOVED***
	if len(rawPort) == 0 ***REMOVED***
		return 0, 0, nil
	***REMOVED***
	start, end, err := ParsePortRange(rawPort)
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***
	return int(start), int(end), nil
***REMOVED***

// Proto returns the protocol of a Port
func (p Port) Proto() string ***REMOVED***
	proto, _ := SplitProtoPort(string(p))
	return proto
***REMOVED***

// Port returns the port number of a Port
func (p Port) Port() string ***REMOVED***
	_, port := SplitProtoPort(string(p))
	return port
***REMOVED***

// Int returns the port number of a Port as an int
func (p Port) Int() int ***REMOVED***
	portStr := p.Port()
	// We don't need to check for an error because we're going to
	// assume that any error would have been found, and reported, in NewPort()
	port, _ := ParsePort(portStr)
	return port
***REMOVED***

// Range returns the start/end port numbers of a Port range as ints
func (p Port) Range() (int, int, error) ***REMOVED***
	return ParsePortRangeToInt(p.Port())
***REMOVED***

// SplitProtoPort splits a port in the format of proto/port
func SplitProtoPort(rawPort string) (string, string) ***REMOVED***
	parts := strings.Split(rawPort, "/")
	l := len(parts)
	if len(rawPort) == 0 || l == 0 || len(parts[0]) == 0 ***REMOVED***
		return "", ""
	***REMOVED***
	if l == 1 ***REMOVED***
		return "tcp", rawPort
	***REMOVED***
	if len(parts[1]) == 0 ***REMOVED***
		return "tcp", parts[0]
	***REMOVED***
	return parts[1], parts[0]
***REMOVED***

func validateProto(proto string) bool ***REMOVED***
	for _, availableProto := range []string***REMOVED***"tcp", "udp"***REMOVED*** ***REMOVED***
		if availableProto == proto ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// ParsePortSpecs receives port specs in the format of ip:public:private/proto and parses
// these in to the internal types
func ParsePortSpecs(ports []string) (map[Port]struct***REMOVED******REMOVED***, map[Port][]PortBinding, error) ***REMOVED***
	var (
		exposedPorts = make(map[Port]struct***REMOVED******REMOVED***, len(ports))
		bindings     = make(map[Port][]PortBinding)
	)
	for _, rawPort := range ports ***REMOVED***
		portMappings, err := ParsePortSpec(rawPort)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***

		for _, portMapping := range portMappings ***REMOVED***
			port := portMapping.Port
			if _, exists := exposedPorts[port]; !exists ***REMOVED***
				exposedPorts[port] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
			bslice, exists := bindings[port]
			if !exists ***REMOVED***
				bslice = []PortBinding***REMOVED******REMOVED***
			***REMOVED***
			bindings[port] = append(bslice, portMapping.Binding)
		***REMOVED***
	***REMOVED***
	return exposedPorts, bindings, nil
***REMOVED***

// PortMapping is a data object mapping a Port to a PortBinding
type PortMapping struct ***REMOVED***
	Port    Port
	Binding PortBinding
***REMOVED***

func splitParts(rawport string) (string, string, string) ***REMOVED***
	parts := strings.Split(rawport, ":")
	n := len(parts)
	containerport := parts[n-1]

	switch n ***REMOVED***
	case 1:
		return "", "", containerport
	case 2:
		return "", parts[0], containerport
	case 3:
		return parts[0], parts[1], containerport
	default:
		return strings.Join(parts[:n-2], ":"), parts[n-2], containerport
	***REMOVED***
***REMOVED***

// ParsePortSpec parses a port specification string into a slice of PortMappings
func ParsePortSpec(rawPort string) ([]PortMapping, error) ***REMOVED***
	var proto string
	rawIP, hostPort, containerPort := splitParts(rawPort)
	proto, containerPort = SplitProtoPort(containerPort)

	// Strip [] from IPV6 addresses
	ip, _, err := net.SplitHostPort(rawIP + ":")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Invalid ip address %v: %s", rawIP, err)
	***REMOVED***
	if ip != "" && net.ParseIP(ip) == nil ***REMOVED***
		return nil, fmt.Errorf("Invalid ip address: %s", ip)
	***REMOVED***
	if containerPort == "" ***REMOVED***
		return nil, fmt.Errorf("No port specified: %s<empty>", rawPort)
	***REMOVED***

	startPort, endPort, err := ParsePortRange(containerPort)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Invalid containerPort: %s", containerPort)
	***REMOVED***

	var startHostPort, endHostPort uint64 = 0, 0
	if len(hostPort) > 0 ***REMOVED***
		startHostPort, endHostPort, err = ParsePortRange(hostPort)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Invalid hostPort: %s", hostPort)
		***REMOVED***
	***REMOVED***

	if hostPort != "" && (endPort-startPort) != (endHostPort-startHostPort) ***REMOVED***
		// Allow host port range iff containerPort is not a range.
		// In this case, use the host port range as the dynamic
		// host port range to allocate into.
		if endPort != startPort ***REMOVED***
			return nil, fmt.Errorf("Invalid ranges specified for container and host Ports: %s and %s", containerPort, hostPort)
		***REMOVED***
	***REMOVED***

	if !validateProto(strings.ToLower(proto)) ***REMOVED***
		return nil, fmt.Errorf("Invalid proto: %s", proto)
	***REMOVED***

	ports := []PortMapping***REMOVED******REMOVED***
	for i := uint64(0); i <= (endPort - startPort); i++ ***REMOVED***
		containerPort = strconv.FormatUint(startPort+i, 10)
		if len(hostPort) > 0 ***REMOVED***
			hostPort = strconv.FormatUint(startHostPort+i, 10)
		***REMOVED***
		// Set hostPort to a range only if there is a single container port
		// and a dynamic host port.
		if startPort == endPort && startHostPort != endHostPort ***REMOVED***
			hostPort = fmt.Sprintf("%s-%s", hostPort, strconv.FormatUint(endHostPort, 10))
		***REMOVED***
		port, err := NewPort(strings.ToLower(proto), containerPort)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		binding := PortBinding***REMOVED***
			HostIP:   ip,
			HostPort: hostPort,
		***REMOVED***
		ports = append(ports, PortMapping***REMOVED***Port: port, Binding: binding***REMOVED***)
	***REMOVED***
	return ports, nil
***REMOVED***
