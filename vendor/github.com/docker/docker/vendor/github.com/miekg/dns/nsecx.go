package dns

import (
	"crypto/sha1"
	"hash"
	"io"
	"strings"
)

type saltWireFmt struct ***REMOVED***
	Salt string `dns:"size-hex"`
***REMOVED***

// HashName hashes a string (label) according to RFC 5155. It returns the hashed string in
// uppercase.
func HashName(label string, ha uint8, iter uint16, salt string) string ***REMOVED***
	saltwire := new(saltWireFmt)
	saltwire.Salt = salt
	wire := make([]byte, DefaultMsgSize)
	n, err := PackStruct(saltwire, wire, 0)
	if err != nil ***REMOVED***
		return ""
	***REMOVED***
	wire = wire[:n]
	name := make([]byte, 255)
	off, err := PackDomainName(strings.ToLower(label), name, 0, nil, false)
	if err != nil ***REMOVED***
		return ""
	***REMOVED***
	name = name[:off]
	var s hash.Hash
	switch ha ***REMOVED***
	case SHA1:
		s = sha1.New()
	default:
		return ""
	***REMOVED***

	// k = 0
	name = append(name, wire...)
	io.WriteString(s, string(name))
	nsec3 := s.Sum(nil)
	// k > 0
	for k := uint16(0); k < iter; k++ ***REMOVED***
		s.Reset()
		nsec3 = append(nsec3, wire...)
		io.WriteString(s, string(nsec3))
		nsec3 = s.Sum(nil)
	***REMOVED***
	return toBase32(nsec3)
***REMOVED***

// Denialer is an interface that should be implemented by types that are used to denial
// answers in DNSSEC.
type Denialer interface ***REMOVED***
	// Cover will check if the (unhashed) name is being covered by this NSEC or NSEC3.
	Cover(name string) bool
	// Match will check if the ownername matches the (unhashed) name for this NSEC3 or NSEC3.
	Match(name string) bool
***REMOVED***

// Cover implements the Denialer interface.
func (rr *NSEC) Cover(name string) bool ***REMOVED***
	return true
***REMOVED***

// Match implements the Denialer interface.
func (rr *NSEC) Match(name string) bool ***REMOVED***
	return true
***REMOVED***

// Cover implements the Denialer interface.
func (rr *NSEC3) Cover(name string) bool ***REMOVED***
	// FIXME(miek): check if the zones match
	// FIXME(miek): check if we're not dealing with parent nsec3
	hname := HashName(name, rr.Hash, rr.Iterations, rr.Salt)
	labels := Split(rr.Hdr.Name)
	if len(labels) < 2 ***REMOVED***
		return false
	***REMOVED***
	hash := strings.ToUpper(rr.Hdr.Name[labels[0] : labels[1]-1]) // -1 to remove the dot
	if hash == rr.NextDomain ***REMOVED***
		return false // empty interval
	***REMOVED***
	if hash > rr.NextDomain ***REMOVED*** // last name, points to apex
		// hname > hash
		// hname > rr.NextDomain
		// TODO(miek)
	***REMOVED***
	if hname <= hash ***REMOVED***
		return false
	***REMOVED***
	if hname >= rr.NextDomain ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// Match implements the Denialer interface.
func (rr *NSEC3) Match(name string) bool ***REMOVED***
	// FIXME(miek): Check if we are in the same zone
	hname := HashName(name, rr.Hash, rr.Iterations, rr.Salt)
	labels := Split(rr.Hdr.Name)
	if len(labels) < 2 ***REMOVED***
		return false
	***REMOVED***
	hash := strings.ToUpper(rr.Hdr.Name[labels[0] : labels[1]-1]) // -1 to remove the .
	if hash == hname ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***
