package endpoints

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// Options provide the configuration needed to direct how the
// endpoints will be resolved.
type Options struct ***REMOVED***
	// DisableSSL forces the endpoint to be resolved as HTTP.
	// instead of HTTPS if the service supports it.
	DisableSSL bool

	// Sets the resolver to resolve the endpoint as a dualstack endpoint
	// for the service. If dualstack support for a service is not known and
	// StrictMatching is not enabled a dualstack endpoint for the service will
	// be returned. This endpoint may not be valid. If StrictMatching is
	// enabled only services that are known to support dualstack will return
	// dualstack endpoints.
	UseDualStack bool

	// Enables strict matching of services and regions resolved endpoints.
	// If the partition doesn't enumerate the exact service and region an
	// error will be returned. This option will prevent returning endpoints
	// that look valid, but may not resolve to any real endpoint.
	StrictMatching bool

	// Enables resolving a service endpoint based on the region provided if the
	// service does not exist. The service endpoint ID will be used as the service
	// domain name prefix. By default the endpoint resolver requires the service
	// to be known when resolving endpoints.
	//
	// If resolving an endpoint on the partition list the provided region will
	// be used to determine which partition's domain name pattern to the service
	// endpoint ID with. If both the service and region are unkonwn and resolving
	// the endpoint on partition list an UnknownEndpointError error will be returned.
	//
	// If resolving and endpoint on a partition specific resolver that partition's
	// domain name pattern will be used with the service endpoint ID. If both
	// region and service do not exist when resolving an endpoint on a specific
	// partition the partition's domain pattern will be used to combine the
	// endpoint and region together.
	//
	// This option is ignored if StrictMatching is enabled.
	ResolveUnknownService bool
***REMOVED***

// Set combines all of the option functions together.
func (o *Options) Set(optFns ...func(*Options)) ***REMOVED***
	for _, fn := range optFns ***REMOVED***
		fn(o)
	***REMOVED***
***REMOVED***

// DisableSSLOption sets the DisableSSL options. Can be used as a functional
// option when resolving endpoints.
func DisableSSLOption(o *Options) ***REMOVED***
	o.DisableSSL = true
***REMOVED***

// UseDualStackOption sets the UseDualStack option. Can be used as a functional
// option when resolving endpoints.
func UseDualStackOption(o *Options) ***REMOVED***
	o.UseDualStack = true
***REMOVED***

// StrictMatchingOption sets the StrictMatching option. Can be used as a functional
// option when resolving endpoints.
func StrictMatchingOption(o *Options) ***REMOVED***
	o.StrictMatching = true
***REMOVED***

// ResolveUnknownServiceOption sets the ResolveUnknownService option. Can be used
// as a functional option when resolving endpoints.
func ResolveUnknownServiceOption(o *Options) ***REMOVED***
	o.ResolveUnknownService = true
***REMOVED***

// A Resolver provides the interface for functionality to resolve endpoints.
// The build in Partition and DefaultResolver return value satisfy this interface.
type Resolver interface ***REMOVED***
	EndpointFor(service, region string, opts ...func(*Options)) (ResolvedEndpoint, error)
***REMOVED***

// ResolverFunc is a helper utility that wraps a function so it satisfies the
// Resolver interface. This is useful when you want to add additional endpoint
// resolving logic, or stub out specific endpoints with custom values.
type ResolverFunc func(service, region string, opts ...func(*Options)) (ResolvedEndpoint, error)

// EndpointFor wraps the ResolverFunc function to satisfy the Resolver interface.
func (fn ResolverFunc) EndpointFor(service, region string, opts ...func(*Options)) (ResolvedEndpoint, error) ***REMOVED***
	return fn(service, region, opts...)
***REMOVED***

var schemeRE = regexp.MustCompile("^([^:]+)://")

// AddScheme adds the HTTP or HTTPS schemes to a endpoint URL if there is no
// scheme. If disableSSL is true HTTP will set HTTP instead of the default HTTPS.
//
// If disableSSL is set, it will only set the URL's scheme if the URL does not
// contain a scheme.
func AddScheme(endpoint string, disableSSL bool) string ***REMOVED***
	if !schemeRE.MatchString(endpoint) ***REMOVED***
		scheme := "https"
		if disableSSL ***REMOVED***
			scheme = "http"
		***REMOVED***
		endpoint = fmt.Sprintf("%s://%s", scheme, endpoint)
	***REMOVED***

	return endpoint
***REMOVED***

// EnumPartitions a provides a way to retrieve the underlying partitions that
// make up the SDK's default Resolver, or any resolver decoded from a model
// file.
//
// Use this interface with DefaultResolver and DecodeModels to get the list of
// Partitions.
type EnumPartitions interface ***REMOVED***
	Partitions() []Partition
***REMOVED***

// RegionsForService returns a map of regions for the partition and service.
// If either the partition or service does not exist false will be returned
// as the second parameter.
//
// This example shows how  to get the regions for DynamoDB in the AWS partition.
//    rs, exists := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, endpoints.DynamodbServiceID)
//
// This is equivalent to using the partition directly.
//    rs := endpoints.AwsPartition().Services()[endpoints.DynamodbServiceID].Regions()
func RegionsForService(ps []Partition, partitionID, serviceID string) (map[string]Region, bool) ***REMOVED***
	for _, p := range ps ***REMOVED***
		if p.ID() != partitionID ***REMOVED***
			continue
		***REMOVED***
		if _, ok := p.p.Services[serviceID]; !ok ***REMOVED***
			break
		***REMOVED***

		s := Service***REMOVED***
			id: serviceID,
			p:  p.p,
		***REMOVED***
		return s.Regions(), true
	***REMOVED***

	return map[string]Region***REMOVED******REMOVED***, false
***REMOVED***

// PartitionForRegion returns the first partition which includes the region
// passed in. This includes both known regions and regions which match
// a pattern supported by the partition which may include regions that are
// not explicitly known by the partition. Use the Regions method of the
// returned Partition if explicit support is needed.
func PartitionForRegion(ps []Partition, regionID string) (Partition, bool) ***REMOVED***
	for _, p := range ps ***REMOVED***
		if _, ok := p.p.Regions[regionID]; ok || p.p.RegionRegex.MatchString(regionID) ***REMOVED***
			return p, true
		***REMOVED***
	***REMOVED***

	return Partition***REMOVED******REMOVED***, false
***REMOVED***

// A Partition provides the ability to enumerate the partition's regions
// and services.
type Partition struct ***REMOVED***
	id string
	p  *partition
***REMOVED***

// ID returns the identifier of the partition.
func (p Partition) ID() string ***REMOVED*** return p.id ***REMOVED***

// EndpointFor attempts to resolve the endpoint based on service and region.
// See Options for information on configuring how the endpoint is resolved.
//
// If the service cannot be found in the metadata the UnknownServiceError
// error will be returned. This validation will occur regardless if
// StrictMatching is enabled. To enable resolving unknown services set the
// "ResolveUnknownService" option to true. When StrictMatching is disabled
// this option allows the partition resolver to resolve a endpoint based on
// the service endpoint ID provided.
//
// When resolving endpoints you can choose to enable StrictMatching. This will
// require the provided service and region to be known by the partition.
// If the endpoint cannot be strictly resolved an error will be returned. This
// mode is useful to ensure the endpoint resolved is valid. Without
// StrictMatching enabled the endpoint returned my look valid but may not work.
// StrictMatching requires the SDK to be updated if you want to take advantage
// of new regions and services expansions.
//
// Errors that can be returned.
//   * UnknownServiceError
//   * UnknownEndpointError
func (p Partition) EndpointFor(service, region string, opts ...func(*Options)) (ResolvedEndpoint, error) ***REMOVED***
	return p.p.EndpointFor(service, region, opts...)
***REMOVED***

// Regions returns a map of Regions indexed by their ID. This is useful for
// enumerating over the regions in a partition.
func (p Partition) Regions() map[string]Region ***REMOVED***
	rs := map[string]Region***REMOVED******REMOVED***
	for id := range p.p.Regions ***REMOVED***
		rs[id] = Region***REMOVED***
			id: id,
			p:  p.p,
		***REMOVED***
	***REMOVED***

	return rs
***REMOVED***

// Services returns a map of Service indexed by their ID. This is useful for
// enumerating over the services in a partition.
func (p Partition) Services() map[string]Service ***REMOVED***
	ss := map[string]Service***REMOVED******REMOVED***
	for id := range p.p.Services ***REMOVED***
		ss[id] = Service***REMOVED***
			id: id,
			p:  p.p,
		***REMOVED***
	***REMOVED***

	return ss
***REMOVED***

// A Region provides information about a region, and ability to resolve an
// endpoint from the context of a region, given a service.
type Region struct ***REMOVED***
	id, desc string
	p        *partition
***REMOVED***

// ID returns the region's identifier.
func (r Region) ID() string ***REMOVED*** return r.id ***REMOVED***

// ResolveEndpoint resolves an endpoint from the context of the region given
// a service. See Partition.EndpointFor for usage and errors that can be returned.
func (r Region) ResolveEndpoint(service string, opts ...func(*Options)) (ResolvedEndpoint, error) ***REMOVED***
	return r.p.EndpointFor(service, r.id, opts...)
***REMOVED***

// Services returns a list of all services that are known to be in this region.
func (r Region) Services() map[string]Service ***REMOVED***
	ss := map[string]Service***REMOVED******REMOVED***
	for id, s := range r.p.Services ***REMOVED***
		if _, ok := s.Endpoints[r.id]; ok ***REMOVED***
			ss[id] = Service***REMOVED***
				id: id,
				p:  r.p,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return ss
***REMOVED***

// A Service provides information about a service, and ability to resolve an
// endpoint from the context of a service, given a region.
type Service struct ***REMOVED***
	id string
	p  *partition
***REMOVED***

// ID returns the identifier for the service.
func (s Service) ID() string ***REMOVED*** return s.id ***REMOVED***

// ResolveEndpoint resolves an endpoint from the context of a service given
// a region. See Partition.EndpointFor for usage and errors that can be returned.
func (s Service) ResolveEndpoint(region string, opts ...func(*Options)) (ResolvedEndpoint, error) ***REMOVED***
	return s.p.EndpointFor(s.id, region, opts...)
***REMOVED***

// Regions returns a map of Regions that the service is present in.
//
// A region is the AWS region the service exists in. Whereas a Endpoint is
// an URL that can be resolved to a instance of a service.
func (s Service) Regions() map[string]Region ***REMOVED***
	rs := map[string]Region***REMOVED******REMOVED***
	for id := range s.p.Services[s.id].Endpoints ***REMOVED***
		if _, ok := s.p.Regions[id]; ok ***REMOVED***
			rs[id] = Region***REMOVED***
				id: id,
				p:  s.p,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return rs
***REMOVED***

// Endpoints returns a map of Endpoints indexed by their ID for all known
// endpoints for a service.
//
// A region is the AWS region the service exists in. Whereas a Endpoint is
// an URL that can be resolved to a instance of a service.
func (s Service) Endpoints() map[string]Endpoint ***REMOVED***
	es := map[string]Endpoint***REMOVED******REMOVED***
	for id := range s.p.Services[s.id].Endpoints ***REMOVED***
		es[id] = Endpoint***REMOVED***
			id:        id,
			serviceID: s.id,
			p:         s.p,
		***REMOVED***
	***REMOVED***

	return es
***REMOVED***

// A Endpoint provides information about endpoints, and provides the ability
// to resolve that endpoint for the service, and the region the endpoint
// represents.
type Endpoint struct ***REMOVED***
	id        string
	serviceID string
	p         *partition
***REMOVED***

// ID returns the identifier for an endpoint.
func (e Endpoint) ID() string ***REMOVED*** return e.id ***REMOVED***

// ServiceID returns the identifier the endpoint belongs to.
func (e Endpoint) ServiceID() string ***REMOVED*** return e.serviceID ***REMOVED***

// ResolveEndpoint resolves an endpoint from the context of a service and
// region the endpoint represents. See Partition.EndpointFor for usage and
// errors that can be returned.
func (e Endpoint) ResolveEndpoint(opts ...func(*Options)) (ResolvedEndpoint, error) ***REMOVED***
	return e.p.EndpointFor(e.serviceID, e.id, opts...)
***REMOVED***

// A ResolvedEndpoint is an endpoint that has been resolved based on a partition
// service, and region.
type ResolvedEndpoint struct ***REMOVED***
	// The endpoint URL
	URL string

	// The region that should be used for signing requests.
	SigningRegion string

	// The service name that should be used for signing requests.
	SigningName string

	// The signing method that should be used for signing requests.
	SigningMethod string
***REMOVED***

// So that the Error interface type can be included as an anonymous field
// in the requestError struct and not conflict with the error.Error() method.
type awsError awserr.Error

// A EndpointNotFoundError is returned when in StrictMatching mode, and the
// endpoint for the service and region cannot be found in any of the partitions.
type EndpointNotFoundError struct ***REMOVED***
	awsError
	Partition string
	Service   string
	Region    string
***REMOVED***

// A UnknownServiceError is returned when the service does not resolve to an
// endpoint. Includes a list of all known services for the partition. Returned
// when a partition does not support the service.
type UnknownServiceError struct ***REMOVED***
	awsError
	Partition string
	Service   string
	Known     []string
***REMOVED***

// NewUnknownServiceError builds and returns UnknownServiceError.
func NewUnknownServiceError(p, s string, known []string) UnknownServiceError ***REMOVED***
	return UnknownServiceError***REMOVED***
		awsError: awserr.New("UnknownServiceError",
			"could not resolve endpoint for unknown service", nil),
		Partition: p,
		Service:   s,
		Known:     known,
	***REMOVED***
***REMOVED***

// String returns the string representation of the error.
func (e UnknownServiceError) Error() string ***REMOVED***
	extra := fmt.Sprintf("partition: %q, service: %q",
		e.Partition, e.Service)
	if len(e.Known) > 0 ***REMOVED***
		extra += fmt.Sprintf(", known: %v", e.Known)
	***REMOVED***
	return awserr.SprintError(e.Code(), e.Message(), extra, e.OrigErr())
***REMOVED***

// String returns the string representation of the error.
func (e UnknownServiceError) String() string ***REMOVED***
	return e.Error()
***REMOVED***

// A UnknownEndpointError is returned when in StrictMatching mode and the
// service is valid, but the region does not resolve to an endpoint. Includes
// a list of all known endpoints for the service.
type UnknownEndpointError struct ***REMOVED***
	awsError
	Partition string
	Service   string
	Region    string
	Known     []string
***REMOVED***

// NewUnknownEndpointError builds and returns UnknownEndpointError.
func NewUnknownEndpointError(p, s, r string, known []string) UnknownEndpointError ***REMOVED***
	return UnknownEndpointError***REMOVED***
		awsError: awserr.New("UnknownEndpointError",
			"could not resolve endpoint", nil),
		Partition: p,
		Service:   s,
		Region:    r,
		Known:     known,
	***REMOVED***
***REMOVED***

// String returns the string representation of the error.
func (e UnknownEndpointError) Error() string ***REMOVED***
	extra := fmt.Sprintf("partition: %q, service: %q, region: %q",
		e.Partition, e.Service, e.Region)
	if len(e.Known) > 0 ***REMOVED***
		extra += fmt.Sprintf(", known: %v", e.Known)
	***REMOVED***
	return awserr.SprintError(e.Code(), e.Message(), extra, e.OrigErr())
***REMOVED***

// String returns the string representation of the error.
func (e UnknownEndpointError) String() string ***REMOVED***
	return e.Error()
***REMOVED***
