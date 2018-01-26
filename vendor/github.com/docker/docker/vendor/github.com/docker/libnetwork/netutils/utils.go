// Network utility functions.

package netutils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/docker/libnetwork/types"
)

var (
	// ErrNetworkOverlapsWithNameservers preformatted error
	ErrNetworkOverlapsWithNameservers = errors.New("requested network overlaps with nameserver")
	// ErrNetworkOverlaps preformatted error
	ErrNetworkOverlaps = errors.New("requested network overlaps with existing network")
	// ErrNoDefaultRoute preformatted error
	ErrNoDefaultRoute = errors.New("no default route")
)

// CheckNameserverOverlaps checks whether the passed network overlaps with any of the nameservers
func CheckNameserverOverlaps(nameservers []string, toCheck *net.IPNet) error ***REMOVED***
	if len(nameservers) > 0 ***REMOVED***
		for _, ns := range nameservers ***REMOVED***
			_, nsNetwork, err := net.ParseCIDR(ns)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if NetworkOverlaps(toCheck, nsNetwork) ***REMOVED***
				return ErrNetworkOverlapsWithNameservers
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// NetworkOverlaps detects overlap between one IPNet and another
func NetworkOverlaps(netX *net.IPNet, netY *net.IPNet) bool ***REMOVED***
	return netX.Contains(netY.IP) || netY.Contains(netX.IP)
***REMOVED***

// NetworkRange calculates the first and last IP addresses in an IPNet
func NetworkRange(network *net.IPNet) (net.IP, net.IP) ***REMOVED***
	if network == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	firstIP := network.IP.Mask(network.Mask)
	lastIP := types.GetIPCopy(firstIP)
	for i := 0; i < len(firstIP); i++ ***REMOVED***
		lastIP[i] = firstIP[i] | ^network.Mask[i]
	***REMOVED***

	if network.IP.To4() != nil ***REMOVED***
		firstIP = firstIP.To4()
		lastIP = lastIP.To4()
	***REMOVED***

	return firstIP, lastIP
***REMOVED***

// GetIfaceAddr returns the first IPv4 address and slice of IPv6 addresses for the specified network interface
func GetIfaceAddr(name string) (net.Addr, []net.Addr, error) ***REMOVED***
	iface, err := net.InterfaceByName(name)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	addrs, err := iface.Addrs()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	var addrs4 []net.Addr
	var addrs6 []net.Addr
	for _, addr := range addrs ***REMOVED***
		ip := (addr.(*net.IPNet)).IP
		if ip4 := ip.To4(); ip4 != nil ***REMOVED***
			addrs4 = append(addrs4, addr)
		***REMOVED*** else if ip6 := ip.To16(); len(ip6) == net.IPv6len ***REMOVED***
			addrs6 = append(addrs6, addr)
		***REMOVED***
	***REMOVED***
	switch ***REMOVED***
	case len(addrs4) == 0:
		return nil, nil, fmt.Errorf("Interface %v has no IPv4 addresses", name)
	case len(addrs4) > 1:
		fmt.Printf("Interface %v has more than 1 IPv4 address. Defaulting to using %v\n",
			name, (addrs4[0].(*net.IPNet)).IP)
	***REMOVED***
	return addrs4[0], addrs6, nil
***REMOVED***

func genMAC(ip net.IP) net.HardwareAddr ***REMOVED***
	hw := make(net.HardwareAddr, 6)
	// The first byte of the MAC address has to comply with these rules:
	// 1. Unicast: Set the least-significant bit to 0.
	// 2. Address is locally administered: Set the second-least-significant bit (U/L) to 1.
	hw[0] = 0x02
	// The first 24 bits of the MAC represent the Organizationally Unique Identifier (OUI).
	// Since this address is locally administered, we can do whatever we want as long as
	// it doesn't conflict with other addresses.
	hw[1] = 0x42
	// Fill the remaining 4 bytes based on the input
	if ip == nil ***REMOVED***
		rand.Read(hw[2:])
	***REMOVED*** else ***REMOVED***
		copy(hw[2:], ip.To4())
	***REMOVED***
	return hw
***REMOVED***

// GenerateRandomMAC returns a new 6-byte(48-bit) hardware address (MAC)
func GenerateRandomMAC() net.HardwareAddr ***REMOVED***
	return genMAC(nil)
***REMOVED***

// GenerateMACFromIP returns a locally administered MAC address where the 4 least
// significant bytes are derived from the IPv4 address.
func GenerateMACFromIP(ip net.IP) net.HardwareAddr ***REMOVED***
	return genMAC(ip)
***REMOVED***

// GenerateRandomName returns a new name joined with a prefix.  This size
// specified is used to truncate the randomly generated value
func GenerateRandomName(prefix string, size int) (string, error) ***REMOVED***
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return prefix + hex.EncodeToString(id)[:size], nil
***REMOVED***

// ReverseIP accepts a V4 or V6 IP string in the canonical form and returns a reversed IP in
// the dotted decimal form . This is used to setup the IP to service name mapping in the optimal
// way for the DNS PTR queries.
func ReverseIP(IP string) string ***REMOVED***
	var reverseIP []string

	if net.ParseIP(IP).To4() != nil ***REMOVED***
		reverseIP = strings.Split(IP, ".")
		l := len(reverseIP)
		for i, j := 0, l-1; i < l/2; i, j = i+1, j-1 ***REMOVED***
			reverseIP[i], reverseIP[j] = reverseIP[j], reverseIP[i]
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		reverseIP = strings.Split(IP, ":")

		// Reversed IPv6 is represented in dotted decimal instead of the typical
		// colon hex notation
		for key := range reverseIP ***REMOVED***
			if len(reverseIP[key]) == 0 ***REMOVED*** // expand the compressed 0s
				reverseIP[key] = strings.Repeat("0000", 8-strings.Count(IP, ":"))
			***REMOVED*** else if len(reverseIP[key]) < 4 ***REMOVED*** // 0-padding needed
				reverseIP[key] = strings.Repeat("0", 4-len(reverseIP[key])) + reverseIP[key]
			***REMOVED***
		***REMOVED***

		reverseIP = strings.Split(strings.Join(reverseIP, ""), "")

		l := len(reverseIP)
		for i, j := 0, l-1; i < l/2; i, j = i+1, j-1 ***REMOVED***
			reverseIP[i], reverseIP[j] = reverseIP[j], reverseIP[i]
		***REMOVED***
	***REMOVED***

	return strings.Join(reverseIP, ".")
***REMOVED***

// ParseAlias parses and validates the specified string as an alias format (name:alias)
func ParseAlias(val string) (string, string, error) ***REMOVED***
	if val == "" ***REMOVED***
		return "", "", errors.New("empty string specified for alias")
	***REMOVED***
	arr := strings.Split(val, ":")
	if len(arr) > 2 ***REMOVED***
		return "", "", fmt.Errorf("bad format for alias: %s", val)
	***REMOVED***
	if len(arr) == 1 ***REMOVED***
		return val, val, nil
	***REMOVED***
	return arr[0], arr[1], nil
***REMOVED***

// ValidateAlias validates that the specified string has a valid alias format (containerName:alias).
func ValidateAlias(val string) (string, error) ***REMOVED***
	if _, _, err := ParseAlias(val); err != nil ***REMOVED***
		return val, err
	***REMOVED***
	return val, nil
***REMOVED***
