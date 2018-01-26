package dns

import (
	"fmt"
	"strings"
)

// PrivateRdata is an interface used for implementing "Private Use" RR types, see
// RFC 6895. This allows one to experiment with new RR types, without requesting an
// official type code. Also see dns.PrivateHandle and dns.PrivateHandleRemove.
type PrivateRdata interface ***REMOVED***
	// String returns the text presentaton of the Rdata of the Private RR.
	String() string
	// Parse parses the Rdata of the private RR.
	Parse([]string) error
	// Pack is used when packing a private RR into a buffer.
	Pack([]byte) (int, error)
	// Unpack is used when unpacking a private RR from a buffer.
	// TODO(miek): diff. signature than Pack, see edns0.go for instance.
	Unpack([]byte) (int, error)
	// Copy copies the Rdata.
	Copy(PrivateRdata) error
	// Len returns the length in octets of the Rdata.
	Len() int
***REMOVED***

// PrivateRR represents an RR that uses a PrivateRdata user-defined type.
// It mocks normal RRs and implements dns.RR interface.
type PrivateRR struct ***REMOVED***
	Hdr  RR_Header
	Data PrivateRdata
***REMOVED***

func mkPrivateRR(rrtype uint16) *PrivateRR ***REMOVED***
	// Panics if RR is not an instance of PrivateRR.
	rrfunc, ok := TypeToRR[rrtype]
	if !ok ***REMOVED***
		panic(fmt.Sprintf("dns: invalid operation with Private RR type %d", rrtype))
	***REMOVED***

	anyrr := rrfunc()
	switch rr := anyrr.(type) ***REMOVED***
	case *PrivateRR:
		return rr
	***REMOVED***
	panic(fmt.Sprintf("dns: RR is not a PrivateRR, TypeToRR[%d] generator returned %T", rrtype, anyrr))
***REMOVED***

// Header return the RR header of r.
func (r *PrivateRR) Header() *RR_Header ***REMOVED*** return &r.Hdr ***REMOVED***

func (r *PrivateRR) String() string ***REMOVED*** return r.Hdr.String() + r.Data.String() ***REMOVED***

// Private len and copy parts to satisfy RR interface.
func (r *PrivateRR) len() int ***REMOVED*** return r.Hdr.len() + r.Data.Len() ***REMOVED***
func (r *PrivateRR) copy() RR ***REMOVED***
	// make new RR like this:
	rr := mkPrivateRR(r.Hdr.Rrtype)
	newh := r.Hdr.copyHeader()
	rr.Hdr = *newh

	err := r.Data.Copy(rr.Data)
	if err != nil ***REMOVED***
		panic("dns: got value that could not be used to copy Private rdata")
	***REMOVED***
	return rr
***REMOVED***

// PrivateHandle registers a private resource record type. It requires
// string and numeric representation of private RR type and generator function as argument.
func PrivateHandle(rtypestr string, rtype uint16, generator func() PrivateRdata) ***REMOVED***
	rtypestr = strings.ToUpper(rtypestr)

	TypeToRR[rtype] = func() RR ***REMOVED*** return &PrivateRR***REMOVED***RR_Header***REMOVED******REMOVED***, generator()***REMOVED*** ***REMOVED***
	TypeToString[rtype] = rtypestr
	StringToType[rtypestr] = rtype

	setPrivateRR := func(h RR_Header, c chan lex, o, f string) (RR, *ParseError, string) ***REMOVED***
		rr := mkPrivateRR(h.Rrtype)
		rr.Hdr = h

		var l lex
		text := make([]string, 0, 2) // could be 0..N elements, median is probably 1
	FETCH:
		for ***REMOVED***
			// TODO(miek): we could also be returning _QUOTE, this might or might not
			// be an issue (basically parsing TXT becomes hard)
			switch l = <-c; l.value ***REMOVED***
			case zNewline, zEOF:
				break FETCH
			case zString:
				text = append(text, l.token)
			***REMOVED***
		***REMOVED***

		err := rr.Data.Parse(text)
		if err != nil ***REMOVED***
			return nil, &ParseError***REMOVED***f, err.Error(), l***REMOVED***, ""
		***REMOVED***

		return rr, nil, ""
	***REMOVED***

	typeToparserFunc[rtype] = parserFunc***REMOVED***setPrivateRR, true***REMOVED***
***REMOVED***

// PrivateHandleRemove removes defenitions required to support private RR type.
func PrivateHandleRemove(rtype uint16) ***REMOVED***
	rtypestr, ok := TypeToString[rtype]
	if ok ***REMOVED***
		delete(TypeToRR, rtype)
		delete(TypeToString, rtype)
		delete(typeToparserFunc, rtype)
		delete(StringToType, rtypestr)
	***REMOVED***
	return
***REMOVED***
