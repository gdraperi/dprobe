package procfs

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

// IPVSStats holds IPVS statistics, as exposed by the kernel in `/proc/net/ip_vs_stats`.
type IPVSStats struct ***REMOVED***
	// Total count of connections.
	Connections uint64
	// Total incoming packages processed.
	IncomingPackets uint64
	// Total outgoing packages processed.
	OutgoingPackets uint64
	// Total incoming traffic.
	IncomingBytes uint64
	// Total outgoing traffic.
	OutgoingBytes uint64
***REMOVED***

// IPVSBackendStatus holds current metrics of one virtual / real address pair.
type IPVSBackendStatus struct ***REMOVED***
	// The local (virtual) IP address.
	LocalAddress net.IP
	// The local (virtual) port.
	LocalPort uint16
	// The transport protocol (TCP, UDP).
	Proto string
	// The remote (real) IP address.
	RemoteAddress net.IP
	// The remote (real) port.
	RemotePort uint16
	// The current number of active connections for this virtual/real address pair.
	ActiveConn uint64
	// The current number of inactive connections for this virtual/real address pair.
	InactConn uint64
	// The current weight of this virtual/real address pair.
	Weight uint64
***REMOVED***

// NewIPVSStats reads the IPVS statistics.
func NewIPVSStats() (IPVSStats, error) ***REMOVED***
	fs, err := NewFS(DefaultMountPoint)
	if err != nil ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, err
	***REMOVED***

	return fs.NewIPVSStats()
***REMOVED***

// NewIPVSStats reads the IPVS statistics from the specified `proc` filesystem.
func (fs FS) NewIPVSStats() (IPVSStats, error) ***REMOVED***
	file, err := os.Open(fs.Path("net/ip_vs_stats"))
	if err != nil ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, err
	***REMOVED***
	defer file.Close()

	return parseIPVSStats(file)
***REMOVED***

// parseIPVSStats performs the actual parsing of `ip_vs_stats`.
func parseIPVSStats(file io.Reader) (IPVSStats, error) ***REMOVED***
	var (
		statContent []byte
		statLines   []string
		statFields  []string
		stats       IPVSStats
	)

	statContent, err := ioutil.ReadAll(file)
	if err != nil ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, err
	***REMOVED***

	statLines = strings.SplitN(string(statContent), "\n", 4)
	if len(statLines) != 4 ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, errors.New("ip_vs_stats corrupt: too short")
	***REMOVED***

	statFields = strings.Fields(statLines[2])
	if len(statFields) != 5 ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, errors.New("ip_vs_stats corrupt: unexpected number of fields")
	***REMOVED***

	stats.Connections, err = strconv.ParseUint(statFields[0], 16, 64)
	if err != nil ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, err
	***REMOVED***
	stats.IncomingPackets, err = strconv.ParseUint(statFields[1], 16, 64)
	if err != nil ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, err
	***REMOVED***
	stats.OutgoingPackets, err = strconv.ParseUint(statFields[2], 16, 64)
	if err != nil ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, err
	***REMOVED***
	stats.IncomingBytes, err = strconv.ParseUint(statFields[3], 16, 64)
	if err != nil ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, err
	***REMOVED***
	stats.OutgoingBytes, err = strconv.ParseUint(statFields[4], 16, 64)
	if err != nil ***REMOVED***
		return IPVSStats***REMOVED******REMOVED***, err
	***REMOVED***

	return stats, nil
***REMOVED***

// NewIPVSBackendStatus reads and returns the status of all (virtual,real) server pairs.
func NewIPVSBackendStatus() ([]IPVSBackendStatus, error) ***REMOVED***
	fs, err := NewFS(DefaultMountPoint)
	if err != nil ***REMOVED***
		return []IPVSBackendStatus***REMOVED******REMOVED***, err
	***REMOVED***

	return fs.NewIPVSBackendStatus()
***REMOVED***

// NewIPVSBackendStatus reads and returns the status of all (virtual,real) server pairs from the specified `proc` filesystem.
func (fs FS) NewIPVSBackendStatus() ([]IPVSBackendStatus, error) ***REMOVED***
	file, err := os.Open(fs.Path("net/ip_vs"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer file.Close()

	return parseIPVSBackendStatus(file)
***REMOVED***

func parseIPVSBackendStatus(file io.Reader) ([]IPVSBackendStatus, error) ***REMOVED***
	var (
		status       []IPVSBackendStatus
		scanner      = bufio.NewScanner(file)
		proto        string
		localAddress net.IP
		localPort    uint16
		err          error
	)

	for scanner.Scan() ***REMOVED***
		fields := strings.Fields(string(scanner.Text()))
		if len(fields) == 0 ***REMOVED***
			continue
		***REMOVED***
		switch ***REMOVED***
		case fields[0] == "IP" || fields[0] == "Prot" || fields[1] == "RemoteAddress:Port":
			continue
		case fields[0] == "TCP" || fields[0] == "UDP":
			if len(fields) < 2 ***REMOVED***
				continue
			***REMOVED***
			proto = fields[0]
			localAddress, localPort, err = parseIPPort(fields[1])
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		case fields[0] == "->":
			if len(fields) < 6 ***REMOVED***
				continue
			***REMOVED***
			remoteAddress, remotePort, err := parseIPPort(fields[1])
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			weight, err := strconv.ParseUint(fields[3], 10, 64)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			activeConn, err := strconv.ParseUint(fields[4], 10, 64)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			inactConn, err := strconv.ParseUint(fields[5], 10, 64)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			status = append(status, IPVSBackendStatus***REMOVED***
				LocalAddress:  localAddress,
				LocalPort:     localPort,
				RemoteAddress: remoteAddress,
				RemotePort:    remotePort,
				Proto:         proto,
				Weight:        weight,
				ActiveConn:    activeConn,
				InactConn:     inactConn,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return status, nil
***REMOVED***

func parseIPPort(s string) (net.IP, uint16, error) ***REMOVED***
	tmp := strings.SplitN(s, ":", 2)

	if len(tmp) != 2 ***REMOVED***
		return nil, 0, fmt.Errorf("invalid IP:Port: %s", s)
	***REMOVED***

	if len(tmp[0]) != 8 && len(tmp[0]) != 32 ***REMOVED***
		return nil, 0, fmt.Errorf("invalid IP: %s", tmp[0])
	***REMOVED***

	ip, err := hex.DecodeString(tmp[0])
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	port, err := strconv.ParseUint(tmp[1], 16, 16)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	return ip, uint16(port), nil
***REMOVED***
