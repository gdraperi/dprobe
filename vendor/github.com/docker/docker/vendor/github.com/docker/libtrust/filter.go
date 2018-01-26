package libtrust

import (
	"path/filepath"
)

// FilterByHosts filters the list of PublicKeys to only those which contain a
// 'hosts' pattern which matches the given host. If *includeEmpty* is true,
// then keys which do not specify any hosts are also returned.
func FilterByHosts(keys []PublicKey, host string, includeEmpty bool) ([]PublicKey, error) ***REMOVED***
	filtered := make([]PublicKey, 0, len(keys))

	for _, pubKey := range keys ***REMOVED***
		var hosts []string
		switch v := pubKey.GetExtendedField("hosts").(type) ***REMOVED***
		case []string:
			hosts = v
		case []interface***REMOVED******REMOVED***:
			for _, value := range v ***REMOVED***
				h, ok := value.(string)
				if !ok ***REMOVED***
					continue
				***REMOVED***
				hosts = append(hosts, h)
			***REMOVED***
		***REMOVED***

		if len(hosts) == 0 ***REMOVED***
			if includeEmpty ***REMOVED***
				filtered = append(filtered, pubKey)
			***REMOVED***
			continue
		***REMOVED***

		// Check if any hosts match pattern
		for _, hostPattern := range hosts ***REMOVED***
			match, err := filepath.Match(hostPattern, host)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if match ***REMOVED***
				filtered = append(filtered, pubKey)
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return filtered, nil
***REMOVED***
