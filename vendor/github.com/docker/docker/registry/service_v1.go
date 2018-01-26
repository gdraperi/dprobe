package registry

import "net/url"

func (s *DefaultService) lookupV1Endpoints(hostname string) (endpoints []APIEndpoint, err error) ***REMOVED***
	if hostname == DefaultNamespace || hostname == DefaultV2Registry.Host || hostname == IndexHostname ***REMOVED***
		return []APIEndpoint***REMOVED******REMOVED***, nil
	***REMOVED***

	tlsConfig, err := s.tlsConfig(hostname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	endpoints = []APIEndpoint***REMOVED***
		***REMOVED***
			URL: &url.URL***REMOVED***
				Scheme: "https",
				Host:   hostname,
			***REMOVED***,
			Version:      APIVersion1,
			TrimHostname: true,
			TLSConfig:    tlsConfig,
		***REMOVED***,
	***REMOVED***

	if tlsConfig.InsecureSkipVerify ***REMOVED***
		endpoints = append(endpoints, APIEndpoint***REMOVED*** // or this
			URL: &url.URL***REMOVED***
				Scheme: "http",
				Host:   hostname,
			***REMOVED***,
			Version:      APIVersion1,
			TrimHostname: true,
			// used to check if supposed to be secure via InsecureSkipVerify
			TLSConfig: tlsConfig,
		***REMOVED***)
	***REMOVED***
	return endpoints, nil
***REMOVED***
