package dns

// NameUsed sets the RRs in the prereq section to
// "Name is in use" RRs. RFC 2136 section 2.4.4.
func (u *Msg) NameUsed(rr []RR) ***REMOVED***
	u.Answer = make([]RR, len(rr))
	for i, r := range rr ***REMOVED***
		u.Answer[i] = &ANY***REMOVED***Hdr: RR_Header***REMOVED***Name: r.Header().Name, Ttl: 0, Rrtype: TypeANY, Class: ClassANY***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// NameNotUsed sets the RRs in the prereq section to
// "Name is in not use" RRs. RFC 2136 section 2.4.5.
func (u *Msg) NameNotUsed(rr []RR) ***REMOVED***
	u.Answer = make([]RR, len(rr))
	for i, r := range rr ***REMOVED***
		u.Answer[i] = &ANY***REMOVED***Hdr: RR_Header***REMOVED***Name: r.Header().Name, Ttl: 0, Rrtype: TypeANY, Class: ClassNONE***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// Used sets the RRs in the prereq section to
// "RRset exists (value dependent -- with rdata)" RRs. RFC 2136 section 2.4.2.
func (u *Msg) Used(rr []RR) ***REMOVED***
	if len(u.Question) == 0 ***REMOVED***
		panic("dns: empty question section")
	***REMOVED***
	u.Answer = make([]RR, len(rr))
	for i, r := range rr ***REMOVED***
		u.Answer[i] = r
		u.Answer[i].Header().Class = u.Question[0].Qclass
	***REMOVED***
***REMOVED***

// RRsetUsed sets the RRs in the prereq section to
// "RRset exists (value independent -- no rdata)" RRs. RFC 2136 section 2.4.1.
func (u *Msg) RRsetUsed(rr []RR) ***REMOVED***
	u.Answer = make([]RR, len(rr))
	for i, r := range rr ***REMOVED***
		u.Answer[i] = r
		u.Answer[i].Header().Class = ClassANY
		u.Answer[i].Header().Ttl = 0
		u.Answer[i].Header().Rdlength = 0
	***REMOVED***
***REMOVED***

// RRsetNotUsed sets the RRs in the prereq section to
// "RRset does not exist" RRs. RFC 2136 section 2.4.3.
func (u *Msg) RRsetNotUsed(rr []RR) ***REMOVED***
	u.Answer = make([]RR, len(rr))
	for i, r := range rr ***REMOVED***
		u.Answer[i] = r
		u.Answer[i].Header().Class = ClassNONE
		u.Answer[i].Header().Rdlength = 0
		u.Answer[i].Header().Ttl = 0
	***REMOVED***
***REMOVED***

// Insert creates a dynamic update packet that adds an complete RRset, see RFC 2136 section 2.5.1.
func (u *Msg) Insert(rr []RR) ***REMOVED***
	if len(u.Question) == 0 ***REMOVED***
		panic("dns: empty question section")
	***REMOVED***
	u.Ns = make([]RR, len(rr))
	for i, r := range rr ***REMOVED***
		u.Ns[i] = r
		u.Ns[i].Header().Class = u.Question[0].Qclass
	***REMOVED***
***REMOVED***

// RemoveRRset creates a dynamic update packet that deletes an RRset, see RFC 2136 section 2.5.2.
func (u *Msg) RemoveRRset(rr []RR) ***REMOVED***
	u.Ns = make([]RR, len(rr))
	for i, r := range rr ***REMOVED***
		u.Ns[i] = &ANY***REMOVED***Hdr: RR_Header***REMOVED***Name: r.Header().Name, Ttl: 0, Rrtype: r.Header().Rrtype, Class: ClassANY***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// RemoveName creates a dynamic update packet that deletes all RRsets of a name, see RFC 2136 section 2.5.3
func (u *Msg) RemoveName(rr []RR) ***REMOVED***
	u.Ns = make([]RR, len(rr))
	for i, r := range rr ***REMOVED***
		u.Ns[i] = &ANY***REMOVED***Hdr: RR_Header***REMOVED***Name: r.Header().Name, Ttl: 0, Rrtype: TypeANY, Class: ClassANY***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// Remove creates a dynamic update packet deletes RR from the RRSset, see RFC 2136 section 2.5.4
func (u *Msg) Remove(rr []RR) ***REMOVED***
	u.Ns = make([]RR, len(rr))
	for i, r := range rr ***REMOVED***
		u.Ns[i] = r
		u.Ns[i].Header().Class = ClassNONE
		u.Ns[i].Header().Ttl = 0
	***REMOVED***
***REMOVED***
