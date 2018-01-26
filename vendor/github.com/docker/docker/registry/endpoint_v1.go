package registry

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/docker/distribution/registry/client/transport"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/sirupsen/logrus"
)

// V1Endpoint stores basic information about a V1 registry endpoint.
type V1Endpoint struct ***REMOVED***
	client   *http.Client
	URL      *url.URL
	IsSecure bool
***REMOVED***

// NewV1Endpoint parses the given address to return a registry endpoint.
func NewV1Endpoint(index *registrytypes.IndexInfo, userAgent string, metaHeaders http.Header) (*V1Endpoint, error) ***REMOVED***
	tlsConfig, err := newTLSConfig(index.Name, index.Secure)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	endpoint, err := newV1EndpointFromStr(GetAuthConfigKey(index), tlsConfig, userAgent, metaHeaders)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := validateEndpoint(endpoint); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return endpoint, nil
***REMOVED***

func validateEndpoint(endpoint *V1Endpoint) error ***REMOVED***
	logrus.Debugf("pinging registry endpoint %s", endpoint)

	// Try HTTPS ping to registry
	endpoint.URL.Scheme = "https"
	if _, err := endpoint.Ping(); err != nil ***REMOVED***
		if endpoint.IsSecure ***REMOVED***
			// If registry is secure and HTTPS failed, show user the error and tell them about `--insecure-registry`
			// in case that's what they need. DO NOT accept unknown CA certificates, and DO NOT fallback to HTTP.
			return fmt.Errorf("invalid registry endpoint %s: %v. If this private registry supports only HTTP or HTTPS with an unknown CA certificate, please add `--insecure-registry %s` to the daemon's arguments. In the case of HTTPS, if you have access to the registry's CA certificate, no need for the flag; simply place the CA certificate at /etc/docker/certs.d/%s/ca.crt", endpoint, err, endpoint.URL.Host, endpoint.URL.Host)
		***REMOVED***

		// If registry is insecure and HTTPS failed, fallback to HTTP.
		logrus.Debugf("Error from registry %q marked as insecure: %v. Insecurely falling back to HTTP", endpoint, err)
		endpoint.URL.Scheme = "http"

		var err2 error
		if _, err2 = endpoint.Ping(); err2 == nil ***REMOVED***
			return nil
		***REMOVED***

		return fmt.Errorf("invalid registry endpoint %q. HTTPS attempt: %v. HTTP attempt: %v", endpoint, err, err2)
	***REMOVED***

	return nil
***REMOVED***

func newV1Endpoint(address url.URL, tlsConfig *tls.Config, userAgent string, metaHeaders http.Header) *V1Endpoint ***REMOVED***
	endpoint := &V1Endpoint***REMOVED***
		IsSecure: (tlsConfig == nil || !tlsConfig.InsecureSkipVerify),
		URL:      new(url.URL),
	***REMOVED***

	*endpoint.URL = address

	// TODO(tiborvass): make sure a ConnectTimeout transport is used
	tr := NewTransport(tlsConfig)
	endpoint.client = HTTPClient(transport.NewTransport(tr, Headers(userAgent, metaHeaders)...))
	return endpoint
***REMOVED***

// trimV1Address trims the version off the address and returns the
// trimmed address or an error if there is a non-V1 version.
func trimV1Address(address string) (string, error) ***REMOVED***
	var (
		chunks        []string
		apiVersionStr string
	)

	if strings.HasSuffix(address, "/") ***REMOVED***
		address = address[:len(address)-1]
	***REMOVED***

	chunks = strings.Split(address, "/")
	apiVersionStr = chunks[len(chunks)-1]
	if apiVersionStr == "v1" ***REMOVED***
		return strings.Join(chunks[:len(chunks)-1], "/"), nil
	***REMOVED***

	for k, v := range apiVersions ***REMOVED***
		if k != APIVersion1 && apiVersionStr == v ***REMOVED***
			return "", fmt.Errorf("unsupported V1 version path %s", apiVersionStr)
		***REMOVED***
	***REMOVED***

	return address, nil
***REMOVED***

func newV1EndpointFromStr(address string, tlsConfig *tls.Config, userAgent string, metaHeaders http.Header) (*V1Endpoint, error) ***REMOVED***
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") ***REMOVED***
		address = "https://" + address
	***REMOVED***

	address, err := trimV1Address(address)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	uri, err := url.Parse(address)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	endpoint := newV1Endpoint(*uri, tlsConfig, userAgent, metaHeaders)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return endpoint, nil
***REMOVED***

// Get the formatted URL for the root of this registry Endpoint
func (e *V1Endpoint) String() string ***REMOVED***
	return e.URL.String() + "/v1/"
***REMOVED***

// Path returns a formatted string for the URL
// of this endpoint with the given path appended.
func (e *V1Endpoint) Path(path string) string ***REMOVED***
	return e.URL.String() + "/v1/" + path
***REMOVED***

// Ping returns a PingResult which indicates whether the registry is standalone or not.
func (e *V1Endpoint) Ping() (PingResult, error) ***REMOVED***
	logrus.Debugf("attempting v1 ping for registry endpoint %s", e)

	if e.String() == IndexServer ***REMOVED***
		// Skip the check, we know this one is valid
		// (and we never want to fallback to http in case of error)
		return PingResult***REMOVED***Standalone: false***REMOVED***, nil
	***REMOVED***

	req, err := http.NewRequest("GET", e.Path("_ping"), nil)
	if err != nil ***REMOVED***
		return PingResult***REMOVED***Standalone: false***REMOVED***, err
	***REMOVED***

	resp, err := e.client.Do(req)
	if err != nil ***REMOVED***
		return PingResult***REMOVED***Standalone: false***REMOVED***, err
	***REMOVED***

	defer resp.Body.Close()

	jsonString, err := ioutil.ReadAll(resp.Body)
	if err != nil ***REMOVED***
		return PingResult***REMOVED***Standalone: false***REMOVED***, fmt.Errorf("error while reading the http response: %s", err)
	***REMOVED***

	// If the header is absent, we assume true for compatibility with earlier
	// versions of the registry. default to true
	info := PingResult***REMOVED***
		Standalone: true,
	***REMOVED***
	if err := json.Unmarshal(jsonString, &info); err != nil ***REMOVED***
		logrus.Debugf("Error unmarshaling the _ping PingResult: %s", err)
		// don't stop here. Just assume sane defaults
	***REMOVED***
	if hdr := resp.Header.Get("X-Docker-Registry-Version"); hdr != "" ***REMOVED***
		logrus.Debugf("Registry version header: '%s'", hdr)
		info.Version = hdr
	***REMOVED***
	logrus.Debugf("PingResult.Version: %q", info.Version)

	standalone := resp.Header.Get("X-Docker-Registry-Standalone")
	logrus.Debugf("Registry standalone header: '%s'", standalone)
	// Accepted values are "true" (case-insensitive) and "1".
	if strings.EqualFold(standalone, "true") || standalone == "1" ***REMOVED***
		info.Standalone = true
	***REMOVED*** else if len(standalone) > 0 ***REMOVED***
		// there is a header set, and it is not "true" or "1", so assume fails
		info.Standalone = false
	***REMOVED***
	logrus.Debugf("PingResult.Standalone: %t", info.Standalone)
	return info, nil
***REMOVED***
