package endpoints

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type partitions []partition

func (ps partitions) EndpointFor(service, region string, opts ...func(*Options)) (ResolvedEndpoint, error) ***REMOVED***
	var opt Options
	opt.Set(opts...)

	for i := 0; i < len(ps); i++ ***REMOVED***
		if !ps[i].canResolveEndpoint(service, region, opt.StrictMatching) ***REMOVED***
			continue
		***REMOVED***

		return ps[i].EndpointFor(service, region, opts...)
	***REMOVED***

	// If loose matching fallback to first partition format to use
	// when resolving the endpoint.
	if !opt.StrictMatching && len(ps) > 0 ***REMOVED***
		return ps[0].EndpointFor(service, region, opts...)
	***REMOVED***

	return ResolvedEndpoint***REMOVED******REMOVED***, NewUnknownEndpointError("all partitions", service, region, []string***REMOVED******REMOVED***)
***REMOVED***

// Partitions satisfies the EnumPartitions interface and returns a list
// of Partitions representing each partition represented in the SDK's
// endpoints model.
func (ps partitions) Partitions() []Partition ***REMOVED***
	parts := make([]Partition, 0, len(ps))
	for i := 0; i < len(ps); i++ ***REMOVED***
		parts = append(parts, ps[i].Partition())
	***REMOVED***

	return parts
***REMOVED***

type partition struct ***REMOVED***
	ID          string      `json:"partition"`
	Name        string      `json:"partitionName"`
	DNSSuffix   string      `json:"dnsSuffix"`
	RegionRegex regionRegex `json:"regionRegex"`
	Defaults    endpoint    `json:"defaults"`
	Regions     regions     `json:"regions"`
	Services    services    `json:"services"`
***REMOVED***

func (p partition) Partition() Partition ***REMOVED***
	return Partition***REMOVED***
		id: p.ID,
		p:  &p,
	***REMOVED***
***REMOVED***

func (p partition) canResolveEndpoint(service, region string, strictMatch bool) bool ***REMOVED***
	s, hasService := p.Services[service]
	_, hasEndpoint := s.Endpoints[region]

	if hasEndpoint && hasService ***REMOVED***
		return true
	***REMOVED***

	if strictMatch ***REMOVED***
		return false
	***REMOVED***

	return p.RegionRegex.MatchString(region)
***REMOVED***

func (p partition) EndpointFor(service, region string, opts ...func(*Options)) (resolved ResolvedEndpoint, err error) ***REMOVED***
	var opt Options
	opt.Set(opts...)

	s, hasService := p.Services[service]
	if !(hasService || opt.ResolveUnknownService) ***REMOVED***
		// Only return error if the resolver will not fallback to creating
		// endpoint based on service endpoint ID passed in.
		return resolved, NewUnknownServiceError(p.ID, service, serviceList(p.Services))
	***REMOVED***

	e, hasEndpoint := s.endpointForRegion(region)
	if !hasEndpoint && opt.StrictMatching ***REMOVED***
		return resolved, NewUnknownEndpointError(p.ID, service, region, endpointList(s.Endpoints))
	***REMOVED***

	defs := []endpoint***REMOVED***p.Defaults, s.Defaults***REMOVED***
	return e.resolve(service, region, p.DNSSuffix, defs, opt), nil
***REMOVED***

func serviceList(ss services) []string ***REMOVED***
	list := make([]string, 0, len(ss))
	for k := range ss ***REMOVED***
		list = append(list, k)
	***REMOVED***
	return list
***REMOVED***
func endpointList(es endpoints) []string ***REMOVED***
	list := make([]string, 0, len(es))
	for k := range es ***REMOVED***
		list = append(list, k)
	***REMOVED***
	return list
***REMOVED***

type regionRegex struct ***REMOVED***
	*regexp.Regexp
***REMOVED***

func (rr *regionRegex) UnmarshalJSON(b []byte) (err error) ***REMOVED***
	// Strip leading and trailing quotes
	regex, err := strconv.Unquote(string(b))
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to strip quotes from regex, %v", err)
	***REMOVED***

	rr.Regexp, err = regexp.Compile(regex)
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to unmarshal region regex, %v", err)
	***REMOVED***
	return nil
***REMOVED***

type regions map[string]region

type region struct ***REMOVED***
	Description string `json:"description"`
***REMOVED***

type services map[string]service

type service struct ***REMOVED***
	PartitionEndpoint string    `json:"partitionEndpoint"`
	IsRegionalized    boxedBool `json:"isRegionalized,omitempty"`
	Defaults          endpoint  `json:"defaults"`
	Endpoints         endpoints `json:"endpoints"`
***REMOVED***

func (s *service) endpointForRegion(region string) (endpoint, bool) ***REMOVED***
	if s.IsRegionalized == boxedFalse ***REMOVED***
		return s.Endpoints[s.PartitionEndpoint], region == s.PartitionEndpoint
	***REMOVED***

	if e, ok := s.Endpoints[region]; ok ***REMOVED***
		return e, true
	***REMOVED***

	// Unable to find any matching endpoint, return
	// blank that will be used for generic endpoint creation.
	return endpoint***REMOVED******REMOVED***, false
***REMOVED***

type endpoints map[string]endpoint

type endpoint struct ***REMOVED***
	Hostname        string          `json:"hostname"`
	Protocols       []string        `json:"protocols"`
	CredentialScope credentialScope `json:"credentialScope"`

	// Custom fields not modeled
	HasDualStack      boxedBool `json:"-"`
	DualStackHostname string    `json:"-"`

	// Signature Version not used
	SignatureVersions []string `json:"signatureVersions"`

	// SSLCommonName not used.
	SSLCommonName string `json:"sslCommonName"`
***REMOVED***

const (
	defaultProtocol = "https"
	defaultSigner   = "v4"
)

var (
	protocolPriority = []string***REMOVED***"https", "http"***REMOVED***
	signerPriority   = []string***REMOVED***"v4", "v2"***REMOVED***
)

func getByPriority(s []string, p []string, def string) string ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return def
	***REMOVED***

	for i := 0; i < len(p); i++ ***REMOVED***
		for j := 0; j < len(s); j++ ***REMOVED***
			if s[j] == p[i] ***REMOVED***
				return s[j]
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return s[0]
***REMOVED***

func (e endpoint) resolve(service, region, dnsSuffix string, defs []endpoint, opts Options) ResolvedEndpoint ***REMOVED***
	var merged endpoint
	for _, def := range defs ***REMOVED***
		merged.mergeIn(def)
	***REMOVED***
	merged.mergeIn(e)
	e = merged

	hostname := e.Hostname

	// Offset the hostname for dualstack if enabled
	if opts.UseDualStack && e.HasDualStack == boxedTrue ***REMOVED***
		hostname = e.DualStackHostname
	***REMOVED***

	u := strings.Replace(hostname, "***REMOVED***service***REMOVED***", service, 1)
	u = strings.Replace(u, "***REMOVED***region***REMOVED***", region, 1)
	u = strings.Replace(u, "***REMOVED***dnsSuffix***REMOVED***", dnsSuffix, 1)

	scheme := getEndpointScheme(e.Protocols, opts.DisableSSL)
	u = fmt.Sprintf("%s://%s", scheme, u)

	signingRegion := e.CredentialScope.Region
	if len(signingRegion) == 0 ***REMOVED***
		signingRegion = region
	***REMOVED***
	signingName := e.CredentialScope.Service
	if len(signingName) == 0 ***REMOVED***
		signingName = service
	***REMOVED***

	return ResolvedEndpoint***REMOVED***
		URL:           u,
		SigningRegion: signingRegion,
		SigningName:   signingName,
		SigningMethod: getByPriority(e.SignatureVersions, signerPriority, defaultSigner),
	***REMOVED***
***REMOVED***

func getEndpointScheme(protocols []string, disableSSL bool) string ***REMOVED***
	if disableSSL ***REMOVED***
		return "http"
	***REMOVED***

	return getByPriority(protocols, protocolPriority, defaultProtocol)
***REMOVED***

func (e *endpoint) mergeIn(other endpoint) ***REMOVED***
	if len(other.Hostname) > 0 ***REMOVED***
		e.Hostname = other.Hostname
	***REMOVED***
	if len(other.Protocols) > 0 ***REMOVED***
		e.Protocols = other.Protocols
	***REMOVED***
	if len(other.SignatureVersions) > 0 ***REMOVED***
		e.SignatureVersions = other.SignatureVersions
	***REMOVED***
	if len(other.CredentialScope.Region) > 0 ***REMOVED***
		e.CredentialScope.Region = other.CredentialScope.Region
	***REMOVED***
	if len(other.CredentialScope.Service) > 0 ***REMOVED***
		e.CredentialScope.Service = other.CredentialScope.Service
	***REMOVED***
	if len(other.SSLCommonName) > 0 ***REMOVED***
		e.SSLCommonName = other.SSLCommonName
	***REMOVED***
	if other.HasDualStack != boxedBoolUnset ***REMOVED***
		e.HasDualStack = other.HasDualStack
	***REMOVED***
	if len(other.DualStackHostname) > 0 ***REMOVED***
		e.DualStackHostname = other.DualStackHostname
	***REMOVED***
***REMOVED***

type credentialScope struct ***REMOVED***
	Region  string `json:"region"`
	Service string `json:"service"`
***REMOVED***

type boxedBool int

func (b *boxedBool) UnmarshalJSON(buf []byte) error ***REMOVED***
	v, err := strconv.ParseBool(string(buf))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if v ***REMOVED***
		*b = boxedTrue
	***REMOVED*** else ***REMOVED***
		*b = boxedFalse
	***REMOVED***

	return nil
***REMOVED***

const (
	boxedBoolUnset boxedBool = iota
	boxedFalse
	boxedTrue
)
