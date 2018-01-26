package dns

// Dedup removes identical RRs from rrs. It preserves the original ordering.
// The lowest TTL of any duplicates is used in the remaining one. Dedup modifies
// rrs.
// m is used to store the RRs temporay. If it is nil a new map will be allocated.
func Dedup(rrs []RR, m map[string]RR) []RR ***REMOVED***
	if m == nil ***REMOVED***
		m = make(map[string]RR)
	***REMOVED***
	// Save the keys, so we don't have to call normalizedString twice.
	keys := make([]*string, 0, len(rrs))

	for _, r := range rrs ***REMOVED***
		key := normalizedString(r)
		keys = append(keys, &key)
		if _, ok := m[key]; ok ***REMOVED***
			// Shortest TTL wins.
			if m[key].Header().Ttl > r.Header().Ttl ***REMOVED***
				m[key].Header().Ttl = r.Header().Ttl
			***REMOVED***
			continue
		***REMOVED***

		m[key] = r
	***REMOVED***
	// If the length of the result map equals the amount of RRs we got,
	// it means they were all different. We can then just return the original rrset.
	if len(m) == len(rrs) ***REMOVED***
		return rrs
	***REMOVED***

	j := 0
	for i, r := range rrs ***REMOVED***
		// If keys[i] lives in the map, we should copy and remove it.
		if _, ok := m[*keys[i]]; ok ***REMOVED***
			delete(m, *keys[i])
			rrs[j] = r
			j++
		***REMOVED***

		if len(m) == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return rrs[:j]
***REMOVED***

// normalizedString returns a normalized string from r. The TTL
// is removed and the domain name is lowercased. We go from this:
// DomainName<TAB>TTL<TAB>CLASS<TAB>TYPE<TAB>RDATA to:
// lowercasename<TAB>CLASS<TAB>TYPE...
func normalizedString(r RR) string ***REMOVED***
	// A string Go DNS makes has: domainname<TAB>TTL<TAB>...
	b := []byte(r.String())

	// find the first non-escaped tab, then another, so we capture where the TTL lives.
	esc := false
	ttlStart, ttlEnd := 0, 0
	for i := 0; i < len(b) && ttlEnd == 0; i++ ***REMOVED***
		switch ***REMOVED***
		case b[i] == '\\':
			esc = !esc
		case b[i] == '\t' && !esc:
			if ttlStart == 0 ***REMOVED***
				ttlStart = i
				continue
			***REMOVED***
			if ttlEnd == 0 ***REMOVED***
				ttlEnd = i
			***REMOVED***
		case b[i] >= 'A' && b[i] <= 'Z' && !esc:
			b[i] += 32
		default:
			esc = false
		***REMOVED***
	***REMOVED***

	// remove TTL.
	copy(b[ttlStart:], b[ttlEnd:])
	cut := ttlEnd - ttlStart
	return string(b[:len(b)-cut])
***REMOVED***
