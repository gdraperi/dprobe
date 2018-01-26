// Package resolvconf provides utility code to query and update DNS configuration in /etc/resolv.conf
package resolvconf

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/libnetwork/resolvconf/dns"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

var (
	// Note: the default IPv4 & IPv6 resolvers are set to Google's Public DNS
	defaultIPv4Dns = []string***REMOVED***"nameserver 8.8.8.8", "nameserver 8.8.4.4"***REMOVED***
	defaultIPv6Dns = []string***REMOVED***"nameserver 2001:4860:4860::8888", "nameserver 2001:4860:4860::8844"***REMOVED***
	ipv4NumBlock   = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	ipv4Address    = `(` + ipv4NumBlock + `\.)***REMOVED***3***REMOVED***` + ipv4NumBlock
	// This is not an IPv6 address verifier as it will accept a super-set of IPv6, and also
	// will *not match* IPv4-Embedded IPv6 Addresses (RFC6052), but that and other variants
	// -- e.g. other link-local types -- either won't work in containers or are unnecessary.
	// For readability and sufficiency for Docker purposes this seemed more reasonable than a
	// 1000+ character regexp with exact and complete IPv6 validation
	ipv6Address = `([0-9A-Fa-f]***REMOVED***0,4***REMOVED***:)***REMOVED***2,7***REMOVED***([0-9A-Fa-f]***REMOVED***0,4***REMOVED***)(%\w+)?`

	localhostNSRegexp = regexp.MustCompile(`(?m)^nameserver\s+` + dns.IPLocalhost + `\s*\n*`)
	nsIPv6Regexp      = regexp.MustCompile(`(?m)^nameserver\s+` + ipv6Address + `\s*\n*`)
	nsRegexp          = regexp.MustCompile(`^\s*nameserver\s*((` + ipv4Address + `)|(` + ipv6Address + `))\s*$`)
	nsIPv6Regexpmatch = regexp.MustCompile(`^\s*nameserver\s*((` + ipv6Address + `))\s*$`)
	nsIPv4Regexpmatch = regexp.MustCompile(`^\s*nameserver\s*((` + ipv4Address + `))\s*$`)
	searchRegexp      = regexp.MustCompile(`^\s*search\s*(([^\s]+\s*)*)$`)
	optionsRegexp     = regexp.MustCompile(`^\s*options\s*(([^\s]+\s*)*)$`)
)

var lastModified struct ***REMOVED***
	sync.Mutex
	sha256   string
	contents []byte
***REMOVED***

// File contains the resolv.conf content and its hash
type File struct ***REMOVED***
	Content []byte
	Hash    string
***REMOVED***

// Get returns the contents of /etc/resolv.conf and its hash
func Get() (*File, error) ***REMOVED***
	resolv, err := ioutil.ReadFile("/etc/resolv.conf")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	hash, err := ioutils.HashData(bytes.NewReader(resolv))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &File***REMOVED***Content: resolv, Hash: hash***REMOVED***, nil
***REMOVED***

// GetSpecific returns the contents of the user specified resolv.conf file and its hash
func GetSpecific(path string) (*File, error) ***REMOVED***
	resolv, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	hash, err := ioutils.HashData(bytes.NewReader(resolv))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &File***REMOVED***Content: resolv, Hash: hash***REMOVED***, nil
***REMOVED***

// GetIfChanged retrieves the host /etc/resolv.conf file, checks against the last hash
// and, if modified since last check, returns the bytes and new hash.
// This feature is used by the resolv.conf updater for containers
func GetIfChanged() (*File, error) ***REMOVED***
	lastModified.Lock()
	defer lastModified.Unlock()

	resolv, err := ioutil.ReadFile("/etc/resolv.conf")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	newHash, err := ioutils.HashData(bytes.NewReader(resolv))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if lastModified.sha256 != newHash ***REMOVED***
		lastModified.sha256 = newHash
		lastModified.contents = resolv
		return &File***REMOVED***Content: resolv, Hash: newHash***REMOVED***, nil
	***REMOVED***
	// nothing changed, so return no data
	return nil, nil
***REMOVED***

// GetLastModified retrieves the last used contents and hash of the host resolv.conf.
// Used by containers updating on restart
func GetLastModified() *File ***REMOVED***
	lastModified.Lock()
	defer lastModified.Unlock()

	return &File***REMOVED***Content: lastModified.contents, Hash: lastModified.sha256***REMOVED***
***REMOVED***

// FilterResolvDNS cleans up the config in resolvConf.  It has two main jobs:
// 1. It looks for localhost (127.*|::1) entries in the provided
//    resolv.conf, removing local nameserver entries, and, if the resulting
//    cleaned config has no defined nameservers left, adds default DNS entries
// 2. Given the caller provides the enable/disable state of IPv6, the filter
//    code will remove all IPv6 nameservers if it is not enabled for containers
//
func FilterResolvDNS(resolvConf []byte, ipv6Enabled bool) (*File, error) ***REMOVED***
	cleanedResolvConf := localhostNSRegexp.ReplaceAll(resolvConf, []byte***REMOVED******REMOVED***)
	// if IPv6 is not enabled, also clean out any IPv6 address nameserver
	if !ipv6Enabled ***REMOVED***
		cleanedResolvConf = nsIPv6Regexp.ReplaceAll(cleanedResolvConf, []byte***REMOVED******REMOVED***)
	***REMOVED***
	// if the resulting resolvConf has no more nameservers defined, add appropriate
	// default DNS servers for IPv4 and (optionally) IPv6
	if len(GetNameservers(cleanedResolvConf, types.IP)) == 0 ***REMOVED***
		logrus.Infof("No non-localhost DNS nameservers are left in resolv.conf. Using default external servers: %v", defaultIPv4Dns)
		dns := defaultIPv4Dns
		if ipv6Enabled ***REMOVED***
			logrus.Infof("IPv6 enabled; Adding default IPv6 external servers: %v", defaultIPv6Dns)
			dns = append(dns, defaultIPv6Dns...)
		***REMOVED***
		cleanedResolvConf = append(cleanedResolvConf, []byte("\n"+strings.Join(dns, "\n"))...)
	***REMOVED***
	hash, err := ioutils.HashData(bytes.NewReader(cleanedResolvConf))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &File***REMOVED***Content: cleanedResolvConf, Hash: hash***REMOVED***, nil
***REMOVED***

// getLines parses input into lines and strips away comments.
func getLines(input []byte, commentMarker []byte) [][]byte ***REMOVED***
	lines := bytes.Split(input, []byte("\n"))
	var output [][]byte
	for _, currentLine := range lines ***REMOVED***
		var commentIndex = bytes.Index(currentLine, commentMarker)
		if commentIndex == -1 ***REMOVED***
			output = append(output, currentLine)
		***REMOVED*** else ***REMOVED***
			output = append(output, currentLine[:commentIndex])
		***REMOVED***
	***REMOVED***
	return output
***REMOVED***

// GetNameservers returns nameservers (if any) listed in /etc/resolv.conf
func GetNameservers(resolvConf []byte, kind int) []string ***REMOVED***
	nameservers := []string***REMOVED******REMOVED***
	for _, line := range getLines(resolvConf, []byte("#")) ***REMOVED***
		var ns [][]byte
		if kind == types.IP ***REMOVED***
			ns = nsRegexp.FindSubmatch(line)
		***REMOVED*** else if kind == types.IPv4 ***REMOVED***
			ns = nsIPv4Regexpmatch.FindSubmatch(line)
		***REMOVED*** else if kind == types.IPv6 ***REMOVED***
			ns = nsIPv6Regexpmatch.FindSubmatch(line)
		***REMOVED***
		if len(ns) > 0 ***REMOVED***
			nameservers = append(nameservers, string(ns[1]))
		***REMOVED***
	***REMOVED***
	return nameservers
***REMOVED***

// GetNameserversAsCIDR returns nameservers (if any) listed in
// /etc/resolv.conf as CIDR blocks (e.g., "1.2.3.4/32")
// This function's output is intended for net.ParseCIDR
func GetNameserversAsCIDR(resolvConf []byte) []string ***REMOVED***
	nameservers := []string***REMOVED******REMOVED***
	for _, nameserver := range GetNameservers(resolvConf, types.IP) ***REMOVED***
		var address string
		// If IPv6, strip zone if present
		if strings.Contains(nameserver, ":") ***REMOVED***
			address = strings.Split(nameserver, "%")[0] + "/128"
		***REMOVED*** else ***REMOVED***
			address = nameserver + "/32"
		***REMOVED***
		nameservers = append(nameservers, address)
	***REMOVED***
	return nameservers
***REMOVED***

// GetSearchDomains returns search domains (if any) listed in /etc/resolv.conf
// If more than one search line is encountered, only the contents of the last
// one is returned.
func GetSearchDomains(resolvConf []byte) []string ***REMOVED***
	domains := []string***REMOVED******REMOVED***
	for _, line := range getLines(resolvConf, []byte("#")) ***REMOVED***
		match := searchRegexp.FindSubmatch(line)
		if match == nil ***REMOVED***
			continue
		***REMOVED***
		domains = strings.Fields(string(match[1]))
	***REMOVED***
	return domains
***REMOVED***

// GetOptions returns options (if any) listed in /etc/resolv.conf
// If more than one options line is encountered, only the contents of the last
// one is returned.
func GetOptions(resolvConf []byte) []string ***REMOVED***
	options := []string***REMOVED******REMOVED***
	for _, line := range getLines(resolvConf, []byte("#")) ***REMOVED***
		match := optionsRegexp.FindSubmatch(line)
		if match == nil ***REMOVED***
			continue
		***REMOVED***
		options = strings.Fields(string(match[1]))
	***REMOVED***
	return options
***REMOVED***

// Build writes a configuration file to path containing a "nameserver" entry
// for every element in dns, a "search" entry for every element in
// dnsSearch, and an "options" entry for every element in dnsOptions.
func Build(path string, dns, dnsSearch, dnsOptions []string) (*File, error) ***REMOVED***
	content := bytes.NewBuffer(nil)
	if len(dnsSearch) > 0 ***REMOVED***
		if searchString := strings.Join(dnsSearch, " "); strings.Trim(searchString, " ") != "." ***REMOVED***
			if _, err := content.WriteString("search " + searchString + "\n"); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, dns := range dns ***REMOVED***
		if _, err := content.WriteString("nameserver " + dns + "\n"); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if len(dnsOptions) > 0 ***REMOVED***
		if optsString := strings.Join(dnsOptions, " "); strings.Trim(optsString, " ") != "" ***REMOVED***
			if _, err := content.WriteString("options " + optsString + "\n"); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	hash, err := ioutils.HashData(bytes.NewReader(content.Bytes()))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &File***REMOVED***Content: content.Bytes(), Hash: hash***REMOVED***, ioutil.WriteFile(path, content.Bytes(), 0644)
***REMOVED***
