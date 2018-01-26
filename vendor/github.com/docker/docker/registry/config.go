package registry

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/distribution/reference"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ServiceOptions holds command line options.
type ServiceOptions struct ***REMOVED***
	AllowNondistributableArtifacts []string `json:"allow-nondistributable-artifacts,omitempty"`
	Mirrors                        []string `json:"registry-mirrors,omitempty"`
	InsecureRegistries             []string `json:"insecure-registries,omitempty"`

	// V2Only controls access to legacy registries.  If it is set to true via the
	// command line flag the daemon will not attempt to contact v1 legacy registries
	V2Only bool `json:"disable-legacy-registry,omitempty"`
***REMOVED***

// serviceConfig holds daemon configuration for the registry service.
type serviceConfig struct ***REMOVED***
	registrytypes.ServiceConfig
	V2Only bool
***REMOVED***

var (
	// DefaultNamespace is the default namespace
	DefaultNamespace = "docker.io"
	// DefaultRegistryVersionHeader is the name of the default HTTP header
	// that carries Registry version info
	DefaultRegistryVersionHeader = "Docker-Distribution-Api-Version"

	// IndexHostname is the index hostname
	IndexHostname = "index.docker.io"
	// IndexServer is used for user auth and image search
	IndexServer = "https://" + IndexHostname + "/v1/"
	// IndexName is the name of the index
	IndexName = "docker.io"

	// NotaryServer is the endpoint serving the Notary trust server
	NotaryServer = "https://notary.docker.io"

	// DefaultV2Registry is the URI of the default v2 registry
	DefaultV2Registry = &url.URL***REMOVED***
		Scheme: "https",
		Host:   "registry-1.docker.io",
	***REMOVED***
)

var (
	// ErrInvalidRepositoryName is an error returned if the repository name did
	// not have the correct form
	ErrInvalidRepositoryName = errors.New("Invalid repository name (ex: \"registry.domain.tld/myrepos\")")

	emptyServiceConfig, _ = newServiceConfig(ServiceOptions***REMOVED******REMOVED***)
)

var (
	validHostPortRegex = regexp.MustCompile(`^` + reference.DomainRegexp.String() + `$`)
)

// for mocking in unit tests
var lookupIP = net.LookupIP

// newServiceConfig returns a new instance of ServiceConfig
func newServiceConfig(options ServiceOptions) (*serviceConfig, error) ***REMOVED***
	config := &serviceConfig***REMOVED***
		ServiceConfig: registrytypes.ServiceConfig***REMOVED***
			InsecureRegistryCIDRs: make([]*registrytypes.NetIPNet, 0),
			IndexConfigs:          make(map[string]*registrytypes.IndexInfo),
			// Hack: Bypass setting the mirrors to IndexConfigs since they are going away
			// and Mirrors are only for the official registry anyways.
		***REMOVED***,
		V2Only: options.V2Only,
	***REMOVED***
	if err := config.LoadAllowNondistributableArtifacts(options.AllowNondistributableArtifacts); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := config.LoadMirrors(options.Mirrors); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := config.LoadInsecureRegistries(options.InsecureRegistries); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return config, nil
***REMOVED***

// LoadAllowNondistributableArtifacts loads allow-nondistributable-artifacts registries into config.
func (config *serviceConfig) LoadAllowNondistributableArtifacts(registries []string) error ***REMOVED***
	cidrs := map[string]*registrytypes.NetIPNet***REMOVED******REMOVED***
	hostnames := map[string]bool***REMOVED******REMOVED***

	for _, r := range registries ***REMOVED***
		if _, err := ValidateIndexName(r); err != nil ***REMOVED***
			return err
		***REMOVED***
		if validateNoScheme(r) != nil ***REMOVED***
			return fmt.Errorf("allow-nondistributable-artifacts registry %s should not contain '://'", r)
		***REMOVED***

		if _, ipnet, err := net.ParseCIDR(r); err == nil ***REMOVED***
			// Valid CIDR.
			cidrs[ipnet.String()] = (*registrytypes.NetIPNet)(ipnet)
		***REMOVED*** else if err := validateHostPort(r); err == nil ***REMOVED***
			// Must be `host:port` if not CIDR.
			hostnames[r] = true
		***REMOVED*** else ***REMOVED***
			return fmt.Errorf("allow-nondistributable-artifacts registry %s is not valid: %v", r, err)
		***REMOVED***
	***REMOVED***

	config.AllowNondistributableArtifactsCIDRs = make([]*(registrytypes.NetIPNet), 0)
	for _, c := range cidrs ***REMOVED***
		config.AllowNondistributableArtifactsCIDRs = append(config.AllowNondistributableArtifactsCIDRs, c)
	***REMOVED***

	config.AllowNondistributableArtifactsHostnames = make([]string, 0)
	for h := range hostnames ***REMOVED***
		config.AllowNondistributableArtifactsHostnames = append(config.AllowNondistributableArtifactsHostnames, h)
	***REMOVED***

	return nil
***REMOVED***

// LoadMirrors loads mirrors to config, after removing duplicates.
// Returns an error if mirrors contains an invalid mirror.
func (config *serviceConfig) LoadMirrors(mirrors []string) error ***REMOVED***
	mMap := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	unique := []string***REMOVED******REMOVED***

	for _, mirror := range mirrors ***REMOVED***
		m, err := ValidateMirror(mirror)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, exist := mMap[m]; !exist ***REMOVED***
			mMap[m] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			unique = append(unique, m)
		***REMOVED***
	***REMOVED***

	config.Mirrors = unique

	// Configure public registry since mirrors may have changed.
	config.IndexConfigs[IndexName] = &registrytypes.IndexInfo***REMOVED***
		Name:     IndexName,
		Mirrors:  config.Mirrors,
		Secure:   true,
		Official: true,
	***REMOVED***

	return nil
***REMOVED***

// LoadInsecureRegistries loads insecure registries to config
func (config *serviceConfig) LoadInsecureRegistries(registries []string) error ***REMOVED***
	// Localhost is by default considered as an insecure registry
	// This is a stop-gap for people who are running a private registry on localhost (especially on Boot2docker).
	//
	// TODO: should we deprecate this once it is easier for people to set up a TLS registry or change
	// daemon flags on boot2docker?
	registries = append(registries, "127.0.0.0/8")

	// Store original InsecureRegistryCIDRs and IndexConfigs
	// Clean InsecureRegistryCIDRs and IndexConfigs in config, as passed registries has all insecure registry info.
	originalCIDRs := config.ServiceConfig.InsecureRegistryCIDRs
	originalIndexInfos := config.ServiceConfig.IndexConfigs

	config.ServiceConfig.InsecureRegistryCIDRs = make([]*registrytypes.NetIPNet, 0)
	config.ServiceConfig.IndexConfigs = make(map[string]*registrytypes.IndexInfo)

skip:
	for _, r := range registries ***REMOVED***
		// validate insecure registry
		if _, err := ValidateIndexName(r); err != nil ***REMOVED***
			// before returning err, roll back to original data
			config.ServiceConfig.InsecureRegistryCIDRs = originalCIDRs
			config.ServiceConfig.IndexConfigs = originalIndexInfos
			return err
		***REMOVED***
		if strings.HasPrefix(strings.ToLower(r), "http://") ***REMOVED***
			logrus.Warnf("insecure registry %s should not contain 'http://' and 'http://' has been removed from the insecure registry config", r)
			r = r[7:]
		***REMOVED*** else if strings.HasPrefix(strings.ToLower(r), "https://") ***REMOVED***
			logrus.Warnf("insecure registry %s should not contain 'https://' and 'https://' has been removed from the insecure registry config", r)
			r = r[8:]
		***REMOVED*** else if validateNoScheme(r) != nil ***REMOVED***
			// Insecure registry should not contain '://'
			// before returning err, roll back to original data
			config.ServiceConfig.InsecureRegistryCIDRs = originalCIDRs
			config.ServiceConfig.IndexConfigs = originalIndexInfos
			return fmt.Errorf("insecure registry %s should not contain '://'", r)
		***REMOVED***
		// Check if CIDR was passed to --insecure-registry
		_, ipnet, err := net.ParseCIDR(r)
		if err == nil ***REMOVED***
			// Valid CIDR. If ipnet is already in config.InsecureRegistryCIDRs, skip.
			data := (*registrytypes.NetIPNet)(ipnet)
			for _, value := range config.InsecureRegistryCIDRs ***REMOVED***
				if value.IP.String() == data.IP.String() && value.Mask.String() == data.Mask.String() ***REMOVED***
					continue skip
				***REMOVED***
			***REMOVED***
			// ipnet is not found, add it in config.InsecureRegistryCIDRs
			config.InsecureRegistryCIDRs = append(config.InsecureRegistryCIDRs, data)

		***REMOVED*** else ***REMOVED***
			if err := validateHostPort(r); err != nil ***REMOVED***
				config.ServiceConfig.InsecureRegistryCIDRs = originalCIDRs
				config.ServiceConfig.IndexConfigs = originalIndexInfos
				return fmt.Errorf("insecure registry %s is not valid: %v", r, err)

			***REMOVED***
			// Assume `host:port` if not CIDR.
			config.IndexConfigs[r] = &registrytypes.IndexInfo***REMOVED***
				Name:     r,
				Mirrors:  make([]string, 0),
				Secure:   false,
				Official: false,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Configure public registry.
	config.IndexConfigs[IndexName] = &registrytypes.IndexInfo***REMOVED***
		Name:     IndexName,
		Mirrors:  config.Mirrors,
		Secure:   true,
		Official: true,
	***REMOVED***

	return nil
***REMOVED***

// allowNondistributableArtifacts returns true if the provided hostname is part of the list of registries
// that allow push of nondistributable artifacts.
//
// The list can contain elements with CIDR notation to specify a whole subnet. If the subnet contains an IP
// of the registry specified by hostname, true is returned.
//
// hostname should be a URL.Host (`host:port` or `host`) where the `host` part can be either a domain name
// or an IP address. If it is a domain name, then it will be resolved to IP addresses for matching. If
// resolution fails, CIDR matching is not performed.
func allowNondistributableArtifacts(config *serviceConfig, hostname string) bool ***REMOVED***
	for _, h := range config.AllowNondistributableArtifactsHostnames ***REMOVED***
		if h == hostname ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return isCIDRMatch(config.AllowNondistributableArtifactsCIDRs, hostname)
***REMOVED***

// isSecureIndex returns false if the provided indexName is part of the list of insecure registries
// Insecure registries accept HTTP and/or accept HTTPS with certificates from unknown CAs.
//
// The list of insecure registries can contain an element with CIDR notation to specify a whole subnet.
// If the subnet contains one of the IPs of the registry specified by indexName, the latter is considered
// insecure.
//
// indexName should be a URL.Host (`host:port` or `host`) where the `host` part can be either a domain name
// or an IP address. If it is a domain name, then it will be resolved in order to check if the IP is contained
// in a subnet. If the resolving is not successful, isSecureIndex will only try to match hostname to any element
// of insecureRegistries.
func isSecureIndex(config *serviceConfig, indexName string) bool ***REMOVED***
	// Check for configured index, first.  This is needed in case isSecureIndex
	// is called from anything besides newIndexInfo, in order to honor per-index configurations.
	if index, ok := config.IndexConfigs[indexName]; ok ***REMOVED***
		return index.Secure
	***REMOVED***

	return !isCIDRMatch(config.InsecureRegistryCIDRs, indexName)
***REMOVED***

// isCIDRMatch returns true if URLHost matches an element of cidrs. URLHost is a URL.Host (`host:port` or `host`)
// where the `host` part can be either a domain name or an IP address. If it is a domain name, then it will be
// resolved to IP addresses for matching. If resolution fails, false is returned.
func isCIDRMatch(cidrs []*registrytypes.NetIPNet, URLHost string) bool ***REMOVED***
	host, _, err := net.SplitHostPort(URLHost)
	if err != nil ***REMOVED***
		// Assume URLHost is of the form `host` without the port and go on.
		host = URLHost
	***REMOVED***

	addrs, err := lookupIP(host)
	if err != nil ***REMOVED***
		ip := net.ParseIP(host)
		if ip != nil ***REMOVED***
			addrs = []net.IP***REMOVED***ip***REMOVED***
		***REMOVED***

		// if ip == nil, then `host` is neither an IP nor it could be looked up,
		// either because the index is unreachable, or because the index is behind an HTTP proxy.
		// So, len(addrs) == 0 and we're not aborting.
	***REMOVED***

	// Try CIDR notation only if addrs has any elements, i.e. if `host`'s IP could be determined.
	for _, addr := range addrs ***REMOVED***
		for _, ipnet := range cidrs ***REMOVED***
			// check if the addr falls in the subnet
			if (*net.IPNet)(ipnet).Contains(addr) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// ValidateMirror validates an HTTP(S) registry mirror
func ValidateMirror(val string) (string, error) ***REMOVED***
	uri, err := url.Parse(val)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("invalid mirror: %q is not a valid URI", val)
	***REMOVED***
	if uri.Scheme != "http" && uri.Scheme != "https" ***REMOVED***
		return "", fmt.Errorf("invalid mirror: unsupported scheme %q in %q", uri.Scheme, uri)
	***REMOVED***
	if (uri.Path != "" && uri.Path != "/") || uri.RawQuery != "" || uri.Fragment != "" ***REMOVED***
		return "", fmt.Errorf("invalid mirror: path, query, or fragment at end of the URI %q", uri)
	***REMOVED***
	if uri.User != nil ***REMOVED***
		// strip password from output
		uri.User = url.UserPassword(uri.User.Username(), "xxxxx")
		return "", fmt.Errorf("invalid mirror: username/password not allowed in URI %q", uri)
	***REMOVED***
	return strings.TrimSuffix(val, "/") + "/", nil
***REMOVED***

// ValidateIndexName validates an index name.
func ValidateIndexName(val string) (string, error) ***REMOVED***
	// TODO: upstream this to check to reference package
	if val == "index.docker.io" ***REMOVED***
		val = "docker.io"
	***REMOVED***
	if strings.HasPrefix(val, "-") || strings.HasSuffix(val, "-") ***REMOVED***
		return "", fmt.Errorf("invalid index name (%s). Cannot begin or end with a hyphen", val)
	***REMOVED***
	return val, nil
***REMOVED***

func validateNoScheme(reposName string) error ***REMOVED***
	if strings.Contains(reposName, "://") ***REMOVED***
		// It cannot contain a scheme!
		return ErrInvalidRepositoryName
	***REMOVED***
	return nil
***REMOVED***

func validateHostPort(s string) error ***REMOVED***
	// Split host and port, and in case s can not be splitted, assume host only
	host, port, err := net.SplitHostPort(s)
	if err != nil ***REMOVED***
		host = s
		port = ""
	***REMOVED***
	// If match against the `host:port` pattern fails,
	// it might be `IPv6:port`, which will be captured by net.ParseIP(host)
	if !validHostPortRegex.MatchString(s) && net.ParseIP(host) == nil ***REMOVED***
		return fmt.Errorf("invalid host %q", host)
	***REMOVED***
	if port != "" ***REMOVED***
		v, err := strconv.Atoi(port)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if v < 0 || v > 65535 ***REMOVED***
			return fmt.Errorf("invalid port %q", port)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// newIndexInfo returns IndexInfo configuration from indexName
func newIndexInfo(config *serviceConfig, indexName string) (*registrytypes.IndexInfo, error) ***REMOVED***
	var err error
	indexName, err = ValidateIndexName(indexName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Return any configured index info, first.
	if index, ok := config.IndexConfigs[indexName]; ok ***REMOVED***
		return index, nil
	***REMOVED***

	// Construct a non-configured index info.
	index := &registrytypes.IndexInfo***REMOVED***
		Name:     indexName,
		Mirrors:  make([]string, 0),
		Official: false,
	***REMOVED***
	index.Secure = isSecureIndex(config, indexName)
	return index, nil
***REMOVED***

// GetAuthConfigKey special-cases using the full index address of the official
// index as the AuthConfig key, and uses the (host)name[:port] for private indexes.
func GetAuthConfigKey(index *registrytypes.IndexInfo) string ***REMOVED***
	if index.Official ***REMOVED***
		return IndexServer
	***REMOVED***
	return index.Name
***REMOVED***

// newRepositoryInfo validates and breaks down a repository name into a RepositoryInfo
func newRepositoryInfo(config *serviceConfig, name reference.Named) (*RepositoryInfo, error) ***REMOVED***
	index, err := newIndexInfo(config, reference.Domain(name))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	official := !strings.ContainsRune(reference.FamiliarName(name), '/')

	return &RepositoryInfo***REMOVED***
		Name:     reference.TrimNamed(name),
		Index:    index,
		Official: official,
	***REMOVED***, nil
***REMOVED***

// ParseRepositoryInfo performs the breakdown of a repository name into a RepositoryInfo, but
// lacks registry configuration.
func ParseRepositoryInfo(reposName reference.Named) (*RepositoryInfo, error) ***REMOVED***
	return newRepositoryInfo(emptyServiceConfig, reposName)
***REMOVED***

// ParseSearchIndexInfo will use repository name to get back an indexInfo.
func ParseSearchIndexInfo(reposName string) (*registrytypes.IndexInfo, error) ***REMOVED***
	indexName, _ := splitReposSearchTerm(reposName)

	indexInfo, err := newIndexInfo(emptyServiceConfig, indexName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return indexInfo, nil
***REMOVED***
