package registry

import "testing"

func TestLookupV1Endpoints(t *testing.T) ***REMOVED***
	s, err := NewService(ServiceOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	cases := []struct ***REMOVED***
		hostname    string
		expectedLen int
	***REMOVED******REMOVED***
		***REMOVED***"example.com", 1***REMOVED***,
		***REMOVED***DefaultNamespace, 0***REMOVED***,
		***REMOVED***DefaultV2Registry.Host, 0***REMOVED***,
		***REMOVED***IndexHostname, 0***REMOVED***,
	***REMOVED***

	for _, c := range cases ***REMOVED***
		if ret, err := s.lookupV1Endpoints(c.hostname); err != nil || len(ret) != c.expectedLen ***REMOVED***
			t.Errorf("lookupV1Endpoints(`"+c.hostname+"`) returned %+v and %+v", ret, err)
		***REMOVED***
	***REMOVED***
***REMOVED***
