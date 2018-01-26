package registry

import (
	"net/url"
	"strings"

	"github.com/docker/go-connections/tlsconfig"
)

func (s *DefaultService) lookupV2Endpoints(hostname string) (endpoints []APIEndpoint, err error) ***REMOVED***
	tlsConfig := tlsconfig.ServerDefault()
	if hostname == DefaultNamespace || hostname == IndexHostname ***REMOVED***
		// v2 mirrors
		for _, mirror := range s.config.Mirrors ***REMOVED***
			if !strings.HasPrefix(mirror, "http://") && !strings.HasPrefix(mirror, "https://") ***REMOVED***
				mirror = "https://" + mirror
			***REMOVED***
			mirrorURL, err := url.Parse(mirror)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			mirrorTLSConfig, err := s.tlsConfigForMirror(mirrorURL)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			endpoints = append(endpoints, APIEndpoint***REMOVED***
				URL: mirrorURL,
				// guess mirrors are v2
				Version:      APIVersion2,
				Mirror:       true,
				TrimHostname: true,
				TLSConfig:    mirrorTLSConfig,
			***REMOVED***)
		***REMOVED***
		// v2 registry
		endpoints = append(endpoints, APIEndpoint***REMOVED***
			URL:          DefaultV2Registry,
			Version:      APIVersion2,
			Official:     true,
			TrimHostname: true,
			TLSConfig:    tlsConfig,
		***REMOVED***)

		return endpoints, nil
	***REMOVED***

	ana := allowNondistributableArtifacts(s.config, hostname)

	tlsConfig, err = s.tlsConfig(hostname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	endpoints = []APIEndpoint***REMOVED***
		***REMOVED***
			URL: &url.URL***REMOVED***
				Scheme: "https",
				Host:   hostname,
			***REMOVED***,
			Version: APIVersion2,
			AllowNondistributableArtifacts: ana,
			TrimHostname:                   true,
			TLSConfig:                      tlsConfig,
		***REMOVED***,
	***REMOVED***

	if tlsConfig.InsecureSkipVerify ***REMOVED***
		endpoints = append(endpoints, APIEndpoint***REMOVED***
			URL: &url.URL***REMOVED***
				Scheme: "http",
				Host:   hostname,
			***REMOVED***,
			Version: APIVersion2,
			AllowNondistributableArtifacts: ana,
			TrimHostname:                   true,
			// used to check if supposed to be secure via InsecureSkipVerify
			TLSConfig: tlsConfig,
		***REMOVED***)
	***REMOVED***

	return endpoints, nil
***REMOVED***
