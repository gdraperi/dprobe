// DNS packet assembly, see RFC 1035. Converting from - Unpack() -
// and to - Pack() - wire format.
// All the packers and unpackers take a (msg []byte, off int)
// and return (off1 int, ok bool).  If they return ok==false, they
// also return off1==len(msg), so that the next unpacker will
// also fail.  This lets us avoid checks of ok until the end of a
// packing sequence.

package dns

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"time"
)

const maxCompressionOffset = 2 << 13 // We have 14 bits for the compression pointer

var (
	// ErrAlg indicates an error with the (DNSSEC) algorithm.
	ErrAlg error = &Error***REMOVED***err: "bad algorithm"***REMOVED***
	// ErrAuth indicates an error in the TSIG authentication.
	ErrAuth error = &Error***REMOVED***err: "bad authentication"***REMOVED***
	// ErrBuf indicates that the buffer used it too small for the message.
	ErrBuf error = &Error***REMOVED***err: "buffer size too small"***REMOVED***
	// ErrConnEmpty indicates a connection is being uses before it is initialized.
	ErrConnEmpty error = &Error***REMOVED***err: "conn has no connection"***REMOVED***
	// ErrExtendedRcode ...
	ErrExtendedRcode error = &Error***REMOVED***err: "bad extended rcode"***REMOVED***
	// ErrFqdn indicates that a domain name does not have a closing dot.
	ErrFqdn error = &Error***REMOVED***err: "domain must be fully qualified"***REMOVED***
	// ErrId indicates there is a mismatch with the message's ID.
	ErrId error = &Error***REMOVED***err: "id mismatch"***REMOVED***
	// ErrKeyAlg indicates that the algorithm in the key is not valid.
	ErrKeyAlg    error = &Error***REMOVED***err: "bad key algorithm"***REMOVED***
	ErrKey       error = &Error***REMOVED***err: "bad key"***REMOVED***
	ErrKeySize   error = &Error***REMOVED***err: "bad key size"***REMOVED***
	ErrNoSig     error = &Error***REMOVED***err: "no signature found"***REMOVED***
	ErrPrivKey   error = &Error***REMOVED***err: "bad private key"***REMOVED***
	ErrRcode     error = &Error***REMOVED***err: "bad rcode"***REMOVED***
	ErrRdata     error = &Error***REMOVED***err: "bad rdata"***REMOVED***
	ErrRRset     error = &Error***REMOVED***err: "bad rrset"***REMOVED***
	ErrSecret    error = &Error***REMOVED***err: "no secrets defined"***REMOVED***
	ErrShortRead error = &Error***REMOVED***err: "short read"***REMOVED***
	// ErrSig indicates that a signature can not be cryptographically validated.
	ErrSig error = &Error***REMOVED***err: "bad signature"***REMOVED***
	// ErrSOA indicates that no SOA RR was seen when doing zone transfers.
	ErrSoa error = &Error***REMOVED***err: "no SOA"***REMOVED***
	// ErrTime indicates a timing error in TSIG authentication.
	ErrTime error = &Error***REMOVED***err: "bad time"***REMOVED***
	// ErrTruncated indicates that we failed to unpack a truncated message.
	// We unpacked as much as we had so Msg can still be used, if desired.
	ErrTruncated error = &Error***REMOVED***err: "failed to unpack truncated message"***REMOVED***
)

// Id, by default, returns a 16 bits random number to be used as a
// message id. The random provided should be good enough. This being a
// variable the function can be reassigned to a custom function.
// For instance, to make it return a static value:
//
//	dns.Id = func() uint16 ***REMOVED*** return 3 ***REMOVED***
var Id func() uint16 = id

// MsgHdr is a a manually-unpacked version of (id, bits).
type MsgHdr struct ***REMOVED***
	Id                 uint16
	Response           bool
	Opcode             int
	Authoritative      bool
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
	Zero               bool
	AuthenticatedData  bool
	CheckingDisabled   bool
	Rcode              int
***REMOVED***

// Msg contains the layout of a DNS message.
type Msg struct ***REMOVED***
	MsgHdr
	Compress bool       `json:"-"` // If true, the message will be compressed when converted to wire format. This not part of the official DNS packet format.
	Question []Question // Holds the RR(s) of the question section.
	Answer   []RR       // Holds the RR(s) of the answer section.
	Ns       []RR       // Holds the RR(s) of the authority section.
	Extra    []RR       // Holds the RR(s) of the additional section.
***REMOVED***

// StringToType is the reverse of TypeToString, needed for string parsing.
var StringToType = reverseInt16(TypeToString)

// StringToClass is the reverse of ClassToString, needed for string parsing.
var StringToClass = reverseInt16(ClassToString)

// Map of opcodes strings.
var StringToOpcode = reverseInt(OpcodeToString)

// Map of rcodes strings.
var StringToRcode = reverseInt(RcodeToString)

// ClassToString is a maps Classes to strings for each CLASS wire type.
var ClassToString = map[uint16]string***REMOVED***
	ClassINET:   "IN",
	ClassCSNET:  "CS",
	ClassCHAOS:  "CH",
	ClassHESIOD: "HS",
	ClassNONE:   "NONE",
	ClassANY:    "ANY",
***REMOVED***

// OpcodeToString maps Opcodes to strings.
var OpcodeToString = map[int]string***REMOVED***
	OpcodeQuery:  "QUERY",
	OpcodeIQuery: "IQUERY",
	OpcodeStatus: "STATUS",
	OpcodeNotify: "NOTIFY",
	OpcodeUpdate: "UPDATE",
***REMOVED***

// RcodeToString maps Rcodes to strings.
var RcodeToString = map[int]string***REMOVED***
	RcodeSuccess:        "NOERROR",
	RcodeFormatError:    "FORMERR",
	RcodeServerFailure:  "SERVFAIL",
	RcodeNameError:      "NXDOMAIN",
	RcodeNotImplemented: "NOTIMPL",
	RcodeRefused:        "REFUSED",
	RcodeYXDomain:       "YXDOMAIN", // From RFC 2136
	RcodeYXRrset:        "YXRRSET",
	RcodeNXRrset:        "NXRRSET",
	RcodeNotAuth:        "NOTAUTH",
	RcodeNotZone:        "NOTZONE",
	RcodeBadSig:         "BADSIG", // Also known as RcodeBadVers, see RFC 6891
	//	RcodeBadVers:        "BADVERS",
	RcodeBadKey:   "BADKEY",
	RcodeBadTime:  "BADTIME",
	RcodeBadMode:  "BADMODE",
	RcodeBadName:  "BADNAME",
	RcodeBadAlg:   "BADALG",
	RcodeBadTrunc: "BADTRUNC",
***REMOVED***

// Rather than write the usual handful of routines to pack and
// unpack every message that can appear on the wire, we use
// reflection to write a generic pack/unpack for structs and then
// use it. Thus, if in the future we need to define new message
// structs, no new pack/unpack/printing code needs to be written.

// Domain names are a sequence of counted strings
// split at the dots. They end with a zero-length string.

// PackDomainName packs a domain name s into msg[off:].
// If compression is wanted compress must be true and the compression
// map needs to hold a mapping between domain names and offsets
// pointing into msg.
func PackDomainName(s string, msg []byte, off int, compression map[string]int, compress bool) (off1 int, err error) ***REMOVED***
	off1, _, err = packDomainName(s, msg, off, compression, compress)
	return
***REMOVED***

func packDomainName(s string, msg []byte, off int, compression map[string]int, compress bool) (off1 int, labels int, err error) ***REMOVED***
	// special case if msg == nil
	lenmsg := 256
	if msg != nil ***REMOVED***
		lenmsg = len(msg)
	***REMOVED***
	ls := len(s)
	if ls == 0 ***REMOVED*** // Ok, for instance when dealing with update RR without any rdata.
		return off, 0, nil
	***REMOVED***
	// If not fully qualified, error out, but only if msg == nil #ugly
	switch ***REMOVED***
	case msg == nil:
		if s[ls-1] != '.' ***REMOVED***
			s += "."
			ls++
		***REMOVED***
	case msg != nil:
		if s[ls-1] != '.' ***REMOVED***
			return lenmsg, 0, ErrFqdn
		***REMOVED***
	***REMOVED***
	// Each dot ends a segment of the name.
	// We trade each dot byte for a length byte.
	// Except for escaped dots (\.), which are normal dots.
	// There is also a trailing zero.

	// Compression
	nameoffset := -1
	pointer := -1
	// Emit sequence of counted strings, chopping at dots.
	begin := 0
	bs := []byte(s)
	roBs, bsFresh, escapedDot := s, true, false
	for i := 0; i < ls; i++ ***REMOVED***
		if bs[i] == '\\' ***REMOVED***
			for j := i; j < ls-1; j++ ***REMOVED***
				bs[j] = bs[j+1]
			***REMOVED***
			ls--
			if off+1 > lenmsg ***REMOVED***
				return lenmsg, labels, ErrBuf
			***REMOVED***
			// check for \DDD
			if i+2 < ls && isDigit(bs[i]) && isDigit(bs[i+1]) && isDigit(bs[i+2]) ***REMOVED***
				bs[i] = dddToByte(bs[i:])
				for j := i + 1; j < ls-2; j++ ***REMOVED***
					bs[j] = bs[j+2]
				***REMOVED***
				ls -= 2
			***REMOVED*** else if bs[i] == 't' ***REMOVED***
				bs[i] = '\t'
			***REMOVED*** else if bs[i] == 'r' ***REMOVED***
				bs[i] = '\r'
			***REMOVED*** else if bs[i] == 'n' ***REMOVED***
				bs[i] = '\n'
			***REMOVED***
			escapedDot = bs[i] == '.'
			bsFresh = false
			continue
		***REMOVED***

		if bs[i] == '.' ***REMOVED***
			if i > 0 && bs[i-1] == '.' && !escapedDot ***REMOVED***
				// two dots back to back is not legal
				return lenmsg, labels, ErrRdata
			***REMOVED***
			if i-begin >= 1<<6 ***REMOVED*** // top two bits of length must be clear
				return lenmsg, labels, ErrRdata
			***REMOVED***
			// off can already (we're in a loop) be bigger than len(msg)
			// this happens when a name isn't fully qualified
			if off+1 > lenmsg ***REMOVED***
				return lenmsg, labels, ErrBuf
			***REMOVED***
			if msg != nil ***REMOVED***
				msg[off] = byte(i - begin)
			***REMOVED***
			offset := off
			off++
			for j := begin; j < i; j++ ***REMOVED***
				if off+1 > lenmsg ***REMOVED***
					return lenmsg, labels, ErrBuf
				***REMOVED***
				if msg != nil ***REMOVED***
					msg[off] = bs[j]
				***REMOVED***
				off++
			***REMOVED***
			if compress && !bsFresh ***REMOVED***
				roBs = string(bs)
				bsFresh = true
			***REMOVED***
			// Dont try to compress '.'
			if compress && roBs[begin:] != "." ***REMOVED***
				if p, ok := compression[roBs[begin:]]; !ok ***REMOVED***
					// Only offsets smaller than this can be used.
					if offset < maxCompressionOffset ***REMOVED***
						compression[roBs[begin:]] = offset
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					// The first hit is the longest matching dname
					// keep the pointer offset we get back and store
					// the offset of the current name, because that's
					// where we need to insert the pointer later

					// If compress is true, we're allowed to compress this dname
					if pointer == -1 && compress ***REMOVED***
						pointer = p         // Where to point to
						nameoffset = offset // Where to point from
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
			labels++
			begin = i + 1
		***REMOVED***
		escapedDot = false
	***REMOVED***
	// Root label is special
	if len(bs) == 1 && bs[0] == '.' ***REMOVED***
		return off, labels, nil
	***REMOVED***
	// If we did compression and we find something add the pointer here
	if pointer != -1 ***REMOVED***
		// We have two bytes (14 bits) to put the pointer in
		// if msg == nil, we will never do compression
		msg[nameoffset], msg[nameoffset+1] = packUint16(uint16(pointer ^ 0xC000))
		off = nameoffset + 1
		goto End
	***REMOVED***
	if msg != nil ***REMOVED***
		msg[off] = 0
	***REMOVED***
End:
	off++
	return off, labels, nil
***REMOVED***

// Unpack a domain name.
// In addition to the simple sequences of counted strings above,
// domain names are allowed to refer to strings elsewhere in the
// packet, to avoid repeating common suffixes when returning
// many entries in a single domain.  The pointers are marked
// by a length byte with the top two bits set.  Ignoring those
// two bits, that byte and the next give a 14 bit offset from msg[0]
// where we should pick up the trail.
// Note that if we jump elsewhere in the packet,
// we return off1 == the offset after the first pointer we found,
// which is where the next record will start.
// In theory, the pointers are only allowed to jump backward.
// We let them jump anywhere and stop jumping after a while.

// UnpackDomainName unpacks a domain name into a string.
func UnpackDomainName(msg []byte, off int) (string, int, error) ***REMOVED***
	s := make([]byte, 0, 64)
	off1 := 0
	lenmsg := len(msg)
	ptr := 0 // number of pointers followed
Loop:
	for ***REMOVED***
		if off >= lenmsg ***REMOVED***
			return "", lenmsg, ErrBuf
		***REMOVED***
		c := int(msg[off])
		off++
		switch c & 0xC0 ***REMOVED***
		case 0x00:
			if c == 0x00 ***REMOVED***
				// end of name
				break Loop
			***REMOVED***
			// literal string
			if off+c > lenmsg ***REMOVED***
				return "", lenmsg, ErrBuf
			***REMOVED***
			for j := off; j < off+c; j++ ***REMOVED***
				switch b := msg[j]; b ***REMOVED***
				case '.', '(', ')', ';', ' ', '@':
					fallthrough
				case '"', '\\':
					s = append(s, '\\', b)
				case '\t':
					s = append(s, '\\', 't')
				case '\r':
					s = append(s, '\\', 'r')
				default:
					if b < 32 || b >= 127 ***REMOVED*** // unprintable use \DDD
						var buf [3]byte
						bufs := strconv.AppendInt(buf[:0], int64(b), 10)
						s = append(s, '\\')
						for i := 0; i < 3-len(bufs); i++ ***REMOVED***
							s = append(s, '0')
						***REMOVED***
						for _, r := range bufs ***REMOVED***
							s = append(s, r)
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						s = append(s, b)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			s = append(s, '.')
			off += c
		case 0xC0:
			// pointer to somewhere else in msg.
			// remember location after first ptr,
			// since that's how many bytes we consumed.
			// also, don't follow too many pointers --
			// maybe there's a loop.
			if off >= lenmsg ***REMOVED***
				return "", lenmsg, ErrBuf
			***REMOVED***
			c1 := msg[off]
			off++
			if ptr == 0 ***REMOVED***
				off1 = off
			***REMOVED***
			if ptr++; ptr > 10 ***REMOVED***
				return "", lenmsg, &Error***REMOVED***err: "too many compression pointers"***REMOVED***
			***REMOVED***
			off = (c^0xC0)<<8 | int(c1)
		default:
			// 0x80 and 0x40 are reserved
			return "", lenmsg, ErrRdata
		***REMOVED***
	***REMOVED***
	if ptr == 0 ***REMOVED***
		off1 = off
	***REMOVED***
	if len(s) == 0 ***REMOVED***
		s = []byte(".")
	***REMOVED***
	return string(s), off1, nil
***REMOVED***

func packTxt(txt []string, msg []byte, offset int, tmp []byte) (int, error) ***REMOVED***
	var err error
	if len(txt) == 0 ***REMOVED***
		if offset >= len(msg) ***REMOVED***
			return offset, ErrBuf
		***REMOVED***
		msg[offset] = 0
		return offset, nil
	***REMOVED***
	for i := range txt ***REMOVED***
		if len(txt[i]) > len(tmp) ***REMOVED***
			return offset, ErrBuf
		***REMOVED***
		offset, err = packTxtString(txt[i], msg, offset, tmp)
		if err != nil ***REMOVED***
			return offset, err
		***REMOVED***
	***REMOVED***
	return offset, err
***REMOVED***

func packTxtString(s string, msg []byte, offset int, tmp []byte) (int, error) ***REMOVED***
	lenByteOffset := offset
	if offset >= len(msg) ***REMOVED***
		return offset, ErrBuf
	***REMOVED***
	offset++
	bs := tmp[:len(s)]
	copy(bs, s)
	for i := 0; i < len(bs); i++ ***REMOVED***
		if len(msg) <= offset ***REMOVED***
			return offset, ErrBuf
		***REMOVED***
		if bs[i] == '\\' ***REMOVED***
			i++
			if i == len(bs) ***REMOVED***
				break
			***REMOVED***
			// check for \DDD
			if i+2 < len(bs) && isDigit(bs[i]) && isDigit(bs[i+1]) && isDigit(bs[i+2]) ***REMOVED***
				msg[offset] = dddToByte(bs[i:])
				i += 2
			***REMOVED*** else if bs[i] == 't' ***REMOVED***
				msg[offset] = '\t'
			***REMOVED*** else if bs[i] == 'r' ***REMOVED***
				msg[offset] = '\r'
			***REMOVED*** else if bs[i] == 'n' ***REMOVED***
				msg[offset] = '\n'
			***REMOVED*** else ***REMOVED***
				msg[offset] = bs[i]
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			msg[offset] = bs[i]
		***REMOVED***
		offset++
	***REMOVED***
	l := offset - lenByteOffset - 1
	if l > 255 ***REMOVED***
		return offset, &Error***REMOVED***err: "string exceeded 255 bytes in txt"***REMOVED***
	***REMOVED***
	msg[lenByteOffset] = byte(l)
	return offset, nil
***REMOVED***

func packOctetString(s string, msg []byte, offset int, tmp []byte) (int, error) ***REMOVED***
	if offset >= len(msg) ***REMOVED***
		return offset, ErrBuf
	***REMOVED***
	bs := tmp[:len(s)]
	copy(bs, s)
	for i := 0; i < len(bs); i++ ***REMOVED***
		if len(msg) <= offset ***REMOVED***
			return offset, ErrBuf
		***REMOVED***
		if bs[i] == '\\' ***REMOVED***
			i++
			if i == len(bs) ***REMOVED***
				break
			***REMOVED***
			// check for \DDD
			if i+2 < len(bs) && isDigit(bs[i]) && isDigit(bs[i+1]) && isDigit(bs[i+2]) ***REMOVED***
				msg[offset] = dddToByte(bs[i:])
				i += 2
			***REMOVED*** else ***REMOVED***
				msg[offset] = bs[i]
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			msg[offset] = bs[i]
		***REMOVED***
		offset++
	***REMOVED***
	return offset, nil
***REMOVED***

func unpackTxt(msg []byte, off0 int) (ss []string, off int, err error) ***REMOVED***
	off = off0
	var s string
	for off < len(msg) && err == nil ***REMOVED***
		s, off, err = unpackTxtString(msg, off)
		if err == nil ***REMOVED***
			ss = append(ss, s)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func unpackTxtString(msg []byte, offset int) (string, int, error) ***REMOVED***
	if offset+1 > len(msg) ***REMOVED***
		return "", offset, &Error***REMOVED***err: "overflow unpacking txt"***REMOVED***
	***REMOVED***
	l := int(msg[offset])
	if offset+l+1 > len(msg) ***REMOVED***
		return "", offset, &Error***REMOVED***err: "overflow unpacking txt"***REMOVED***
	***REMOVED***
	s := make([]byte, 0, l)
	for _, b := range msg[offset+1 : offset+1+l] ***REMOVED***
		switch b ***REMOVED***
		case '"', '\\':
			s = append(s, '\\', b)
		case '\t':
			s = append(s, `\t`...)
		case '\r':
			s = append(s, `\r`...)
		case '\n':
			s = append(s, `\n`...)
		default:
			if b < 32 || b > 127 ***REMOVED*** // unprintable
				var buf [3]byte
				bufs := strconv.AppendInt(buf[:0], int64(b), 10)
				s = append(s, '\\')
				for i := 0; i < 3-len(bufs); i++ ***REMOVED***
					s = append(s, '0')
				***REMOVED***
				for _, r := range bufs ***REMOVED***
					s = append(s, r)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				s = append(s, b)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	offset += 1 + l
	return string(s), offset, nil
***REMOVED***

// Pack a reflect.StructValue into msg.  Struct members can only be uint8, uint16, uint32, string,
// slices and other (often anonymous) structs.
func packStructValue(val reflect.Value, msg []byte, off int, compression map[string]int, compress bool) (off1 int, err error) ***REMOVED***
	var txtTmp []byte
	lenmsg := len(msg)
	numfield := val.NumField()
	for i := 0; i < numfield; i++ ***REMOVED***
		typefield := val.Type().Field(i)
		if typefield.Tag == `dns:"-"` ***REMOVED***
			continue
		***REMOVED***
		switch fv := val.Field(i); fv.Kind() ***REMOVED***
		default:
			return lenmsg, &Error***REMOVED***err: "bad kind packing"***REMOVED***
		case reflect.Interface:
			// PrivateRR is the only RR implementation that has interface field.
			// therefore it's expected that this interface would be PrivateRdata
			switch data := fv.Interface().(type) ***REMOVED***
			case PrivateRdata:
				n, err := data.Pack(msg[off:])
				if err != nil ***REMOVED***
					return lenmsg, err
				***REMOVED***
				off += n
			default:
				return lenmsg, &Error***REMOVED***err: "bad kind interface packing"***REMOVED***
			***REMOVED***
		case reflect.Slice:
			switch typefield.Tag ***REMOVED***
			default:
				return lenmsg, &Error***REMOVED***"bad tag packing slice: " + typefield.Tag.Get("dns")***REMOVED***
			case `dns:"domain-name"`:
				for j := 0; j < val.Field(i).Len(); j++ ***REMOVED***
					element := val.Field(i).Index(j).String()
					off, err = PackDomainName(element, msg, off, compression, false && compress)
					if err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
				***REMOVED***
			case `dns:"txt"`:
				if txtTmp == nil ***REMOVED***
					txtTmp = make([]byte, 256*4+1)
				***REMOVED***
				off, err = packTxt(fv.Interface().([]string), msg, off, txtTmp)
				if err != nil ***REMOVED***
					return lenmsg, err
				***REMOVED***
			case `dns:"opt"`: // edns
				for j := 0; j < val.Field(i).Len(); j++ ***REMOVED***
					element := val.Field(i).Index(j).Interface()
					b, e := element.(EDNS0).pack()
					if e != nil ***REMOVED***
						return lenmsg, &Error***REMOVED***err: "overflow packing opt"***REMOVED***
					***REMOVED***
					// Option code
					msg[off], msg[off+1] = packUint16(element.(EDNS0).Option())
					// Length
					msg[off+2], msg[off+3] = packUint16(uint16(len(b)))
					off += 4
					if off+len(b) > lenmsg ***REMOVED***
						copy(msg[off:], b)
						off = lenmsg
						continue
					***REMOVED***
					// Actual data
					copy(msg[off:off+len(b)], b)
					off += len(b)
				***REMOVED***
			case `dns:"a"`:
				if val.Type().String() == "dns.IPSECKEY" ***REMOVED***
					// Field(2) is GatewayType, must be 1
					if val.Field(2).Uint() != 1 ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
				// It must be a slice of 4, even if it is 16, we encode
				// only the first 4
				if off+net.IPv4len > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow packing a"***REMOVED***
				***REMOVED***
				switch fv.Len() ***REMOVED***
				case net.IPv6len:
					msg[off] = byte(fv.Index(12).Uint())
					msg[off+1] = byte(fv.Index(13).Uint())
					msg[off+2] = byte(fv.Index(14).Uint())
					msg[off+3] = byte(fv.Index(15).Uint())
					off += net.IPv4len
				case net.IPv4len:
					msg[off] = byte(fv.Index(0).Uint())
					msg[off+1] = byte(fv.Index(1).Uint())
					msg[off+2] = byte(fv.Index(2).Uint())
					msg[off+3] = byte(fv.Index(3).Uint())
					off += net.IPv4len
				case 0:
					// Allowed, for dynamic updates
				default:
					return lenmsg, &Error***REMOVED***err: "overflow packing a"***REMOVED***
				***REMOVED***
			case `dns:"aaaa"`:
				if val.Type().String() == "dns.IPSECKEY" ***REMOVED***
					// Field(2) is GatewayType, must be 2
					if val.Field(2).Uint() != 2 ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
				if fv.Len() == 0 ***REMOVED***
					break
				***REMOVED***
				if fv.Len() > net.IPv6len || off+fv.Len() > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow packing aaaa"***REMOVED***
				***REMOVED***
				for j := 0; j < net.IPv6len; j++ ***REMOVED***
					msg[off] = byte(fv.Index(j).Uint())
					off++
				***REMOVED***
			case `dns:"wks"`:
				// TODO(miek): this is wrong should be lenrd
				if off == lenmsg ***REMOVED***
					break // dyn. updates
				***REMOVED***
				if val.Field(i).Len() == 0 ***REMOVED***
					break
				***REMOVED***
				off1 := off
				for j := 0; j < val.Field(i).Len(); j++ ***REMOVED***
					serv := int(fv.Index(j).Uint())
					if off+serv/8+1 > len(msg) ***REMOVED***
						return len(msg), &Error***REMOVED***err: "overflow packing wks"***REMOVED***
					***REMOVED***
					msg[off+serv/8] |= byte(1 << (7 - uint(serv%8)))
					if off+serv/8+1 > off1 ***REMOVED***
						off1 = off + serv/8 + 1
					***REMOVED***
				***REMOVED***
				off = off1
			case `dns:"nsec"`: // NSEC/NSEC3
				// This is the uint16 type bitmap
				if val.Field(i).Len() == 0 ***REMOVED***
					// Do absolutely nothing
					break
				***REMOVED***
				var lastwindow, lastlength uint16
				for j := 0; j < val.Field(i).Len(); j++ ***REMOVED***
					t := uint16(fv.Index(j).Uint())
					window := t / 256
					length := (t-window*256)/8 + 1
					if window > lastwindow && lastlength != 0 ***REMOVED***
						// New window, jump to the new offset
						off += int(lastlength) + 2
						lastlength = 0
					***REMOVED***
					if window < lastwindow || length < lastlength ***REMOVED***
						return len(msg), &Error***REMOVED***err: "nsec bits out of order"***REMOVED***
					***REMOVED***
					if off+2+int(length) > len(msg) ***REMOVED***
						return len(msg), &Error***REMOVED***err: "overflow packing nsec"***REMOVED***
					***REMOVED***
					// Setting the window #
					msg[off] = byte(window)
					// Setting the octets length
					msg[off+1] = byte(length)
					// Setting the bit value for the type in the right octet
					msg[off+1+int(length)] |= byte(1 << (7 - (t % 8)))
					lastwindow, lastlength = window, length
				***REMOVED***
				off += int(lastlength) + 2
			***REMOVED***
		case reflect.Struct:
			off, err = packStructValue(fv, msg, off, compression, compress)
			if err != nil ***REMOVED***
				return lenmsg, err
			***REMOVED***
		case reflect.Uint8:
			if off+1 > lenmsg ***REMOVED***
				return lenmsg, &Error***REMOVED***err: "overflow packing uint8"***REMOVED***
			***REMOVED***
			msg[off] = byte(fv.Uint())
			off++
		case reflect.Uint16:
			if off+2 > lenmsg ***REMOVED***
				return lenmsg, &Error***REMOVED***err: "overflow packing uint16"***REMOVED***
			***REMOVED***
			i := fv.Uint()
			msg[off] = byte(i >> 8)
			msg[off+1] = byte(i)
			off += 2
		case reflect.Uint32:
			if off+4 > lenmsg ***REMOVED***
				return lenmsg, &Error***REMOVED***err: "overflow packing uint32"***REMOVED***
			***REMOVED***
			i := fv.Uint()
			msg[off] = byte(i >> 24)
			msg[off+1] = byte(i >> 16)
			msg[off+2] = byte(i >> 8)
			msg[off+3] = byte(i)
			off += 4
		case reflect.Uint64:
			switch typefield.Tag ***REMOVED***
			default:
				if off+8 > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow packing uint64"***REMOVED***
				***REMOVED***
				i := fv.Uint()
				msg[off] = byte(i >> 56)
				msg[off+1] = byte(i >> 48)
				msg[off+2] = byte(i >> 40)
				msg[off+3] = byte(i >> 32)
				msg[off+4] = byte(i >> 24)
				msg[off+5] = byte(i >> 16)
				msg[off+6] = byte(i >> 8)
				msg[off+7] = byte(i)
				off += 8
			case `dns:"uint48"`:
				// Used in TSIG, where it stops at 48 bits, so we discard the upper 16
				if off+6 > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow packing uint64 as uint48"***REMOVED***
				***REMOVED***
				i := fv.Uint()
				msg[off] = byte(i >> 40)
				msg[off+1] = byte(i >> 32)
				msg[off+2] = byte(i >> 24)
				msg[off+3] = byte(i >> 16)
				msg[off+4] = byte(i >> 8)
				msg[off+5] = byte(i)
				off += 6
			***REMOVED***
		case reflect.String:
			// There are multiple string encodings.
			// The tag distinguishes ordinary strings from domain names.
			s := fv.String()
			switch typefield.Tag ***REMOVED***
			default:
				return lenmsg, &Error***REMOVED***"bad tag packing string: " + typefield.Tag.Get("dns")***REMOVED***
			case `dns:"base64"`:
				b64, e := fromBase64([]byte(s))
				if e != nil ***REMOVED***
					return lenmsg, e
				***REMOVED***
				copy(msg[off:off+len(b64)], b64)
				off += len(b64)
			case `dns:"domain-name"`:
				if val.Type().String() == "dns.IPSECKEY" ***REMOVED***
					// Field(2) is GatewayType, 1 and 2 or used for addresses
					x := val.Field(2).Uint()
					if x == 1 || x == 2 ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
				if off, err = PackDomainName(s, msg, off, compression, false && compress); err != nil ***REMOVED***
					return lenmsg, err
				***REMOVED***
			case `dns:"cdomain-name"`:
				if off, err = PackDomainName(s, msg, off, compression, true && compress); err != nil ***REMOVED***
					return lenmsg, err
				***REMOVED***
			case `dns:"size-base32"`:
				// This is purely for NSEC3 atm, the previous byte must
				// holds the length of the encoded string. As NSEC3
				// is only defined to SHA1, the hashlength is 20 (160 bits)
				msg[off-1] = 20
				fallthrough
			case `dns:"base32"`:
				b32, e := fromBase32([]byte(s))
				if e != nil ***REMOVED***
					return lenmsg, e
				***REMOVED***
				copy(msg[off:off+len(b32)], b32)
				off += len(b32)
			case `dns:"size-hex"`:
				fallthrough
			case `dns:"hex"`:
				// There is no length encoded here
				h, e := hex.DecodeString(s)
				if e != nil ***REMOVED***
					return lenmsg, e
				***REMOVED***
				if off+hex.DecodedLen(len(s)) > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow packing hex"***REMOVED***
				***REMOVED***
				copy(msg[off:off+hex.DecodedLen(len(s))], h)
				off += hex.DecodedLen(len(s))
			case `dns:"size"`:
				// the size is already encoded in the RR, we can safely use the
				// length of string. String is RAW (not encoded in hex, nor base64)
				copy(msg[off:off+len(s)], s)
				off += len(s)
			case `dns:"octet"`:
				bytesTmp := make([]byte, 256)
				off, err = packOctetString(fv.String(), msg, off, bytesTmp)
				if err != nil ***REMOVED***
					return lenmsg, err
				***REMOVED***
			case `dns:"txt"`:
				fallthrough
			case "":
				if txtTmp == nil ***REMOVED***
					txtTmp = make([]byte, 256*4+1)
				***REMOVED***
				off, err = packTxtString(fv.String(), msg, off, txtTmp)
				if err != nil ***REMOVED***
					return lenmsg, err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return off, nil
***REMOVED***

func structValue(any interface***REMOVED******REMOVED***) reflect.Value ***REMOVED***
	return reflect.ValueOf(any).Elem()
***REMOVED***

// PackStruct packs any structure to wire format.
func PackStruct(any interface***REMOVED******REMOVED***, msg []byte, off int) (off1 int, err error) ***REMOVED***
	off, err = packStructValue(structValue(any), msg, off, nil, false)
	return off, err
***REMOVED***

func packStructCompress(any interface***REMOVED******REMOVED***, msg []byte, off int, compression map[string]int, compress bool) (off1 int, err error) ***REMOVED***
	off, err = packStructValue(structValue(any), msg, off, compression, compress)
	return off, err
***REMOVED***

// Unpack a reflect.StructValue from msg.
// Same restrictions as packStructValue.
func unpackStructValue(val reflect.Value, msg []byte, off int) (off1 int, err error) ***REMOVED***
	lenmsg := len(msg)
	for i := 0; i < val.NumField(); i++ ***REMOVED***
		if off > lenmsg ***REMOVED***
			return lenmsg, &Error***REMOVED***"bad offset unpacking"***REMOVED***
		***REMOVED***
		switch fv := val.Field(i); fv.Kind() ***REMOVED***
		default:
			return lenmsg, &Error***REMOVED***err: "bad kind unpacking"***REMOVED***
		case reflect.Interface:
			// PrivateRR is the only RR implementation that has interface field.
			// therefore it's expected that this interface would be PrivateRdata
			switch data := fv.Interface().(type) ***REMOVED***
			case PrivateRdata:
				n, err := data.Unpack(msg[off:])
				if err != nil ***REMOVED***
					return lenmsg, err
				***REMOVED***
				off += n
			default:
				return lenmsg, &Error***REMOVED***err: "bad kind interface unpacking"***REMOVED***
			***REMOVED***
		case reflect.Slice:
			switch val.Type().Field(i).Tag ***REMOVED***
			default:
				return lenmsg, &Error***REMOVED***"bad tag unpacking slice: " + val.Type().Field(i).Tag.Get("dns")***REMOVED***
			case `dns:"domain-name"`:
				// HIP record slice of name (or none)
				var servers []string
				var s string
				for off < lenmsg ***REMOVED***
					s, off, err = UnpackDomainName(msg, off)
					if err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
					servers = append(servers, s)
				***REMOVED***
				fv.Set(reflect.ValueOf(servers))
			case `dns:"txt"`:
				if off == lenmsg ***REMOVED***
					break
				***REMOVED***
				var txt []string
				txt, off, err = unpackTxt(msg, off)
				if err != nil ***REMOVED***
					return lenmsg, err
				***REMOVED***
				fv.Set(reflect.ValueOf(txt))
			case `dns:"opt"`: // edns0
				if off == lenmsg ***REMOVED***
					// This is an EDNS0 (OPT Record) with no rdata
					// We can safely return here.
					break
				***REMOVED***
				var edns []EDNS0
			Option:
				code := uint16(0)
				if off+4 > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking opt"***REMOVED***
				***REMOVED***
				code, off = unpackUint16(msg, off)
				optlen, off1 := unpackUint16(msg, off)
				if off1+int(optlen) > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking opt"***REMOVED***
				***REMOVED***
				switch code ***REMOVED***
				case EDNS0NSID:
					e := new(EDNS0_NSID)
					if err := e.unpack(msg[off1 : off1+int(optlen)]); err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
					edns = append(edns, e)
					off = off1 + int(optlen)
				case EDNS0SUBNET, EDNS0SUBNETDRAFT:
					e := new(EDNS0_SUBNET)
					if err := e.unpack(msg[off1 : off1+int(optlen)]); err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
					edns = append(edns, e)
					off = off1 + int(optlen)
					if code == EDNS0SUBNETDRAFT ***REMOVED***
						e.DraftOption = true
					***REMOVED***
				case EDNS0UL:
					e := new(EDNS0_UL)
					if err := e.unpack(msg[off1 : off1+int(optlen)]); err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
					edns = append(edns, e)
					off = off1 + int(optlen)
				case EDNS0LLQ:
					e := new(EDNS0_LLQ)
					if err := e.unpack(msg[off1 : off1+int(optlen)]); err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
					edns = append(edns, e)
					off = off1 + int(optlen)
				case EDNS0DAU:
					e := new(EDNS0_DAU)
					if err := e.unpack(msg[off1 : off1+int(optlen)]); err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
					edns = append(edns, e)
					off = off1 + int(optlen)
				case EDNS0DHU:
					e := new(EDNS0_DHU)
					if err := e.unpack(msg[off1 : off1+int(optlen)]); err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
					edns = append(edns, e)
					off = off1 + int(optlen)
				case EDNS0N3U:
					e := new(EDNS0_N3U)
					if err := e.unpack(msg[off1 : off1+int(optlen)]); err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
					edns = append(edns, e)
					off = off1 + int(optlen)
				default:
					e := new(EDNS0_LOCAL)
					e.Code = code
					if err := e.unpack(msg[off1 : off1+int(optlen)]); err != nil ***REMOVED***
						return lenmsg, err
					***REMOVED***
					edns = append(edns, e)
					off = off1 + int(optlen)
				***REMOVED***
				if off < lenmsg ***REMOVED***
					goto Option
				***REMOVED***
				fv.Set(reflect.ValueOf(edns))
			case `dns:"a"`:
				if val.Type().String() == "dns.IPSECKEY" ***REMOVED***
					// Field(2) is GatewayType, must be 1
					if val.Field(2).Uint() != 1 ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
				if off == lenmsg ***REMOVED***
					break // dyn. update
				***REMOVED***
				if off+net.IPv4len > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking a"***REMOVED***
				***REMOVED***
				fv.Set(reflect.ValueOf(net.IPv4(msg[off], msg[off+1], msg[off+2], msg[off+3])))
				off += net.IPv4len
			case `dns:"aaaa"`:
				if val.Type().String() == "dns.IPSECKEY" ***REMOVED***
					// Field(2) is GatewayType, must be 2
					if val.Field(2).Uint() != 2 ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
				if off == lenmsg ***REMOVED***
					break
				***REMOVED***
				if off+net.IPv6len > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking aaaa"***REMOVED***
				***REMOVED***
				fv.Set(reflect.ValueOf(net.IP***REMOVED***msg[off], msg[off+1], msg[off+2], msg[off+3], msg[off+4],
					msg[off+5], msg[off+6], msg[off+7], msg[off+8], msg[off+9], msg[off+10],
					msg[off+11], msg[off+12], msg[off+13], msg[off+14], msg[off+15]***REMOVED***))
				off += net.IPv6len
			case `dns:"wks"`:
				// Rest of the record is the bitmap
				var serv []uint16
				j := 0
				for off < lenmsg ***REMOVED***
					if off+1 > lenmsg ***REMOVED***
						return lenmsg, &Error***REMOVED***err: "overflow unpacking wks"***REMOVED***
					***REMOVED***
					b := msg[off]
					// Check the bits one by one, and set the type
					if b&0x80 == 0x80 ***REMOVED***
						serv = append(serv, uint16(j*8+0))
					***REMOVED***
					if b&0x40 == 0x40 ***REMOVED***
						serv = append(serv, uint16(j*8+1))
					***REMOVED***
					if b&0x20 == 0x20 ***REMOVED***
						serv = append(serv, uint16(j*8+2))
					***REMOVED***
					if b&0x10 == 0x10 ***REMOVED***
						serv = append(serv, uint16(j*8+3))
					***REMOVED***
					if b&0x8 == 0x8 ***REMOVED***
						serv = append(serv, uint16(j*8+4))
					***REMOVED***
					if b&0x4 == 0x4 ***REMOVED***
						serv = append(serv, uint16(j*8+5))
					***REMOVED***
					if b&0x2 == 0x2 ***REMOVED***
						serv = append(serv, uint16(j*8+6))
					***REMOVED***
					if b&0x1 == 0x1 ***REMOVED***
						serv = append(serv, uint16(j*8+7))
					***REMOVED***
					j++
					off++
				***REMOVED***
				fv.Set(reflect.ValueOf(serv))
			case `dns:"nsec"`: // NSEC/NSEC3
				if off == len(msg) ***REMOVED***
					break
				***REMOVED***
				// Rest of the record is the type bitmap
				var nsec []uint16
				length := 0
				window := 0
				lastwindow := -1
				for off < len(msg) ***REMOVED***
					if off+2 > len(msg) ***REMOVED***
						return len(msg), &Error***REMOVED***err: "overflow unpacking nsecx"***REMOVED***
					***REMOVED***
					window = int(msg[off])
					length = int(msg[off+1])
					off += 2
					if window <= lastwindow ***REMOVED***
						// RFC 4034: Blocks are present in the NSEC RR RDATA in
						// increasing numerical order.
						return len(msg), &Error***REMOVED***err: "out of order NSEC block"***REMOVED***
					***REMOVED***
					if length == 0 ***REMOVED***
						// RFC 4034: Blocks with no types present MUST NOT be included.
						return len(msg), &Error***REMOVED***err: "empty NSEC block"***REMOVED***
					***REMOVED***
					if length > 32 ***REMOVED***
						return len(msg), &Error***REMOVED***err: "NSEC block too long"***REMOVED***
					***REMOVED***
					if off+length > len(msg) ***REMOVED***
						return len(msg), &Error***REMOVED***err: "overflowing NSEC block"***REMOVED***
					***REMOVED***

					// Walk the bytes in the window and extract the type bits
					for j := 0; j < length; j++ ***REMOVED***
						b := msg[off+j]
						// Check the bits one by one, and set the type
						if b&0x80 == 0x80 ***REMOVED***
							nsec = append(nsec, uint16(window*256+j*8+0))
						***REMOVED***
						if b&0x40 == 0x40 ***REMOVED***
							nsec = append(nsec, uint16(window*256+j*8+1))
						***REMOVED***
						if b&0x20 == 0x20 ***REMOVED***
							nsec = append(nsec, uint16(window*256+j*8+2))
						***REMOVED***
						if b&0x10 == 0x10 ***REMOVED***
							nsec = append(nsec, uint16(window*256+j*8+3))
						***REMOVED***
						if b&0x8 == 0x8 ***REMOVED***
							nsec = append(nsec, uint16(window*256+j*8+4))
						***REMOVED***
						if b&0x4 == 0x4 ***REMOVED***
							nsec = append(nsec, uint16(window*256+j*8+5))
						***REMOVED***
						if b&0x2 == 0x2 ***REMOVED***
							nsec = append(nsec, uint16(window*256+j*8+6))
						***REMOVED***
						if b&0x1 == 0x1 ***REMOVED***
							nsec = append(nsec, uint16(window*256+j*8+7))
						***REMOVED***
					***REMOVED***
					off += length
					lastwindow = window
				***REMOVED***
				fv.Set(reflect.ValueOf(nsec))
			***REMOVED***
		case reflect.Struct:
			off, err = unpackStructValue(fv, msg, off)
			if err != nil ***REMOVED***
				return lenmsg, err
			***REMOVED***
			if val.Type().Field(i).Name == "Hdr" ***REMOVED***
				lenrd := off + int(val.FieldByName("Hdr").FieldByName("Rdlength").Uint())
				if lenrd > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflowing header size"***REMOVED***
				***REMOVED***
				msg = msg[:lenrd]
				lenmsg = len(msg)
			***REMOVED***
		case reflect.Uint8:
			if off == lenmsg ***REMOVED***
				break
			***REMOVED***
			if off+1 > lenmsg ***REMOVED***
				return lenmsg, &Error***REMOVED***err: "overflow unpacking uint8"***REMOVED***
			***REMOVED***
			fv.SetUint(uint64(uint8(msg[off])))
			off++
		case reflect.Uint16:
			if off == lenmsg ***REMOVED***
				break
			***REMOVED***
			var i uint16
			if off+2 > lenmsg ***REMOVED***
				return lenmsg, &Error***REMOVED***err: "overflow unpacking uint16"***REMOVED***
			***REMOVED***
			i, off = unpackUint16(msg, off)
			fv.SetUint(uint64(i))
		case reflect.Uint32:
			if off == lenmsg ***REMOVED***
				break
			***REMOVED***
			if off+4 > lenmsg ***REMOVED***
				return lenmsg, &Error***REMOVED***err: "overflow unpacking uint32"***REMOVED***
			***REMOVED***
			fv.SetUint(uint64(uint32(msg[off])<<24 | uint32(msg[off+1])<<16 | uint32(msg[off+2])<<8 | uint32(msg[off+3])))
			off += 4
		case reflect.Uint64:
			if off == lenmsg ***REMOVED***
				break
			***REMOVED***
			switch val.Type().Field(i).Tag ***REMOVED***
			default:
				if off+8 > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking uint64"***REMOVED***
				***REMOVED***
				fv.SetUint(uint64(uint64(msg[off])<<56 | uint64(msg[off+1])<<48 | uint64(msg[off+2])<<40 |
					uint64(msg[off+3])<<32 | uint64(msg[off+4])<<24 | uint64(msg[off+5])<<16 | uint64(msg[off+6])<<8 | uint64(msg[off+7])))
				off += 8
			case `dns:"uint48"`:
				// Used in TSIG where the last 48 bits are occupied, so for now, assume a uint48 (6 bytes)
				if off+6 > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking uint64 as uint48"***REMOVED***
				***REMOVED***
				fv.SetUint(uint64(uint64(msg[off])<<40 | uint64(msg[off+1])<<32 | uint64(msg[off+2])<<24 | uint64(msg[off+3])<<16 |
					uint64(msg[off+4])<<8 | uint64(msg[off+5])))
				off += 6
			***REMOVED***
		case reflect.String:
			var s string
			if off == lenmsg ***REMOVED***
				break
			***REMOVED***
			switch val.Type().Field(i).Tag ***REMOVED***
			default:
				return lenmsg, &Error***REMOVED***"bad tag unpacking string: " + val.Type().Field(i).Tag.Get("dns")***REMOVED***
			case `dns:"octet"`:
				s = string(msg[off:])
				off = lenmsg
			case `dns:"hex"`:
				hexend := lenmsg
				if val.FieldByName("Hdr").FieldByName("Rrtype").Uint() == uint64(TypeHIP) ***REMOVED***
					hexend = off + int(val.FieldByName("HitLength").Uint())
				***REMOVED***
				if hexend > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking HIP hex"***REMOVED***
				***REMOVED***
				s = hex.EncodeToString(msg[off:hexend])
				off = hexend
			case `dns:"base64"`:
				// Rest of the RR is base64 encoded value
				b64end := lenmsg
				if val.FieldByName("Hdr").FieldByName("Rrtype").Uint() == uint64(TypeHIP) ***REMOVED***
					b64end = off + int(val.FieldByName("PublicKeyLength").Uint())
				***REMOVED***
				if b64end > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking HIP base64"***REMOVED***
				***REMOVED***
				s = toBase64(msg[off:b64end])
				off = b64end
			case `dns:"cdomain-name"`:
				fallthrough
			case `dns:"domain-name"`:
				if val.Type().String() == "dns.IPSECKEY" ***REMOVED***
					// Field(2) is GatewayType, 1 and 2 or used for addresses
					x := val.Field(2).Uint()
					if x == 1 || x == 2 ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
				if off == lenmsg && int(val.FieldByName("Hdr").FieldByName("Rdlength").Uint()) == 0 ***REMOVED***
					// zero rdata is ok for dyn updates, but only if rdlength is 0
					break
				***REMOVED***
				s, off, err = UnpackDomainName(msg, off)
				if err != nil ***REMOVED***
					return lenmsg, err
				***REMOVED***
			case `dns:"size-base32"`:
				var size int
				switch val.Type().Name() ***REMOVED***
				case "NSEC3":
					switch val.Type().Field(i).Name ***REMOVED***
					case "NextDomain":
						name := val.FieldByName("HashLength")
						size = int(name.Uint())
					***REMOVED***
				***REMOVED***
				if off+size > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking base32"***REMOVED***
				***REMOVED***
				s = toBase32(msg[off : off+size])
				off += size
			case `dns:"size-hex"`:
				// a "size" string, but it must be encoded in hex in the string
				var size int
				switch val.Type().Name() ***REMOVED***
				case "NSEC3":
					switch val.Type().Field(i).Name ***REMOVED***
					case "Salt":
						name := val.FieldByName("SaltLength")
						size = int(name.Uint())
					case "NextDomain":
						name := val.FieldByName("HashLength")
						size = int(name.Uint())
					***REMOVED***
				case "TSIG":
					switch val.Type().Field(i).Name ***REMOVED***
					case "MAC":
						name := val.FieldByName("MACSize")
						size = int(name.Uint())
					case "OtherData":
						name := val.FieldByName("OtherLen")
						size = int(name.Uint())
					***REMOVED***
				***REMOVED***
				if off+size > lenmsg ***REMOVED***
					return lenmsg, &Error***REMOVED***err: "overflow unpacking hex"***REMOVED***
				***REMOVED***
				s = hex.EncodeToString(msg[off : off+size])
				off += size
			case `dns:"txt"`:
				fallthrough
			case "":
				s, off, err = unpackTxtString(msg, off)
			***REMOVED***
			fv.SetString(s)
		***REMOVED***
	***REMOVED***
	return off, nil
***REMOVED***

// Helpers for dealing with escaped bytes
func isDigit(b byte) bool ***REMOVED*** return b >= '0' && b <= '9' ***REMOVED***

func dddToByte(s []byte) byte ***REMOVED***
	return byte((s[0]-'0')*100 + (s[1]-'0')*10 + (s[2] - '0'))
***REMOVED***

// UnpackStruct unpacks a binary message from offset off to the interface
// value given.
func UnpackStruct(any interface***REMOVED******REMOVED***, msg []byte, off int) (int, error) ***REMOVED***
	return unpackStructValue(structValue(any), msg, off)
***REMOVED***

// Helper function for packing and unpacking
func intToBytes(i *big.Int, length int) []byte ***REMOVED***
	buf := i.Bytes()
	if len(buf) < length ***REMOVED***
		b := make([]byte, length)
		copy(b[length-len(buf):], buf)
		return b
	***REMOVED***
	return buf
***REMOVED***

func unpackUint16(msg []byte, off int) (uint16, int) ***REMOVED***
	return uint16(msg[off])<<8 | uint16(msg[off+1]), off + 2
***REMOVED***

func packUint16(i uint16) (byte, byte) ***REMOVED***
	return byte(i >> 8), byte(i)
***REMOVED***

func toBase32(b []byte) string ***REMOVED***
	return base32.HexEncoding.EncodeToString(b)
***REMOVED***

func fromBase32(s []byte) (buf []byte, err error) ***REMOVED***
	buflen := base32.HexEncoding.DecodedLen(len(s))
	buf = make([]byte, buflen)
	n, err := base32.HexEncoding.Decode(buf, s)
	buf = buf[:n]
	return
***REMOVED***

func toBase64(b []byte) string ***REMOVED***
	return base64.StdEncoding.EncodeToString(b)
***REMOVED***

func fromBase64(s []byte) (buf []byte, err error) ***REMOVED***
	buflen := base64.StdEncoding.DecodedLen(len(s))
	buf = make([]byte, buflen)
	n, err := base64.StdEncoding.Decode(buf, s)
	buf = buf[:n]
	return
***REMOVED***

// PackRR packs a resource record rr into msg[off:].
// See PackDomainName for documentation about the compression.
func PackRR(rr RR, msg []byte, off int, compression map[string]int, compress bool) (off1 int, err error) ***REMOVED***
	if rr == nil ***REMOVED***
		return len(msg), &Error***REMOVED***err: "nil rr"***REMOVED***
	***REMOVED***

	off1, err = packStructCompress(rr, msg, off, compression, compress)
	if err != nil ***REMOVED***
		return len(msg), err
	***REMOVED***
	if rawSetRdlength(msg, off, off1) ***REMOVED***
		return off1, nil
	***REMOVED***
	return off, ErrRdata
***REMOVED***

// UnpackRR unpacks msg[off:] into an RR.
func UnpackRR(msg []byte, off int) (rr RR, off1 int, err error) ***REMOVED***
	// unpack just the header, to find the rr type and length
	var h RR_Header
	off0 := off
	if off, err = UnpackStruct(&h, msg, off); err != nil ***REMOVED***
		return nil, len(msg), err
	***REMOVED***
	end := off + int(h.Rdlength)
	// make an rr of that type and re-unpack.
	mk, known := TypeToRR[h.Rrtype]
	if !known ***REMOVED***
		rr = new(RFC3597)
	***REMOVED*** else ***REMOVED***
		rr = mk()
	***REMOVED***
	off, err = UnpackStruct(rr, msg, off0)
	if off != end ***REMOVED***
		return &h, end, &Error***REMOVED***err: "bad rdlength"***REMOVED***
	***REMOVED***
	return rr, off, err
***REMOVED***

// unpackRRslice unpacks msg[off:] into an []RR.
// If we cannot unpack the whole array, then it will return nil
func unpackRRslice(l int, msg []byte, off int) (dst1 []RR, off1 int, err error) ***REMOVED***
	var r RR
	// Optimistically make dst be the length that was sent
	dst := make([]RR, 0, l)
	for i := 0; i < l; i++ ***REMOVED***
		off1 := off
		r, off, err = UnpackRR(msg, off)
		if err != nil ***REMOVED***
			off = len(msg)
			break
		***REMOVED***
		// If offset does not increase anymore, l is a lie
		if off1 == off ***REMOVED***
			l = i
			break
		***REMOVED***
		dst = append(dst, r)
	***REMOVED***
	if err != nil && off == len(msg) ***REMOVED***
		dst = nil
	***REMOVED***
	return dst, off, err
***REMOVED***

// Reverse a map
func reverseInt8(m map[uint8]string) map[string]uint8 ***REMOVED***
	n := make(map[string]uint8)
	for u, s := range m ***REMOVED***
		n[s] = u
	***REMOVED***
	return n
***REMOVED***

func reverseInt16(m map[uint16]string) map[string]uint16 ***REMOVED***
	n := make(map[string]uint16)
	for u, s := range m ***REMOVED***
		n[s] = u
	***REMOVED***
	return n
***REMOVED***

func reverseInt(m map[int]string) map[string]int ***REMOVED***
	n := make(map[string]int)
	for u, s := range m ***REMOVED***
		n[s] = u
	***REMOVED***
	return n
***REMOVED***

// Convert a MsgHdr to a string, with dig-like headers:
//
//;; opcode: QUERY, status: NOERROR, id: 48404
//
//;; flags: qr aa rd ra;
func (h *MsgHdr) String() string ***REMOVED***
	if h == nil ***REMOVED***
		return "<nil> MsgHdr"
	***REMOVED***

	s := ";; opcode: " + OpcodeToString[h.Opcode]
	s += ", status: " + RcodeToString[h.Rcode]
	s += ", id: " + strconv.Itoa(int(h.Id)) + "\n"

	s += ";; flags:"
	if h.Response ***REMOVED***
		s += " qr"
	***REMOVED***
	if h.Authoritative ***REMOVED***
		s += " aa"
	***REMOVED***
	if h.Truncated ***REMOVED***
		s += " tc"
	***REMOVED***
	if h.RecursionDesired ***REMOVED***
		s += " rd"
	***REMOVED***
	if h.RecursionAvailable ***REMOVED***
		s += " ra"
	***REMOVED***
	if h.Zero ***REMOVED*** // Hmm
		s += " z"
	***REMOVED***
	if h.AuthenticatedData ***REMOVED***
		s += " ad"
	***REMOVED***
	if h.CheckingDisabled ***REMOVED***
		s += " cd"
	***REMOVED***

	s += ";"
	return s
***REMOVED***

// Pack packs a Msg: it is converted to to wire format.
// If the dns.Compress is true the message will be in compressed wire format.
func (dns *Msg) Pack() (msg []byte, err error) ***REMOVED***
	return dns.PackBuffer(nil)
***REMOVED***

// PackBuffer packs a Msg, using the given buffer buf. If buf is too small
// a new buffer is allocated.
func (dns *Msg) PackBuffer(buf []byte) (msg []byte, err error) ***REMOVED***
	var dh Header
	var compression map[string]int
	if dns.Compress ***REMOVED***
		compression = make(map[string]int) // Compression pointer mappings
	***REMOVED***

	if dns.Rcode < 0 || dns.Rcode > 0xFFF ***REMOVED***
		return nil, ErrRcode
	***REMOVED***
	if dns.Rcode > 0xF ***REMOVED***
		// Regular RCODE field is 4 bits
		opt := dns.IsEdns0()
		if opt == nil ***REMOVED***
			return nil, ErrExtendedRcode
		***REMOVED***
		opt.SetExtendedRcode(uint8(dns.Rcode >> 4))
		dns.Rcode &= 0xF
	***REMOVED***

	// Convert convenient Msg into wire-like Header.
	dh.Id = dns.Id
	dh.Bits = uint16(dns.Opcode)<<11 | uint16(dns.Rcode)
	if dns.Response ***REMOVED***
		dh.Bits |= _QR
	***REMOVED***
	if dns.Authoritative ***REMOVED***
		dh.Bits |= _AA
	***REMOVED***
	if dns.Truncated ***REMOVED***
		dh.Bits |= _TC
	***REMOVED***
	if dns.RecursionDesired ***REMOVED***
		dh.Bits |= _RD
	***REMOVED***
	if dns.RecursionAvailable ***REMOVED***
		dh.Bits |= _RA
	***REMOVED***
	if dns.Zero ***REMOVED***
		dh.Bits |= _Z
	***REMOVED***
	if dns.AuthenticatedData ***REMOVED***
		dh.Bits |= _AD
	***REMOVED***
	if dns.CheckingDisabled ***REMOVED***
		dh.Bits |= _CD
	***REMOVED***

	// Prepare variable sized arrays.
	question := dns.Question
	answer := dns.Answer
	ns := dns.Ns
	extra := dns.Extra

	dh.Qdcount = uint16(len(question))
	dh.Ancount = uint16(len(answer))
	dh.Nscount = uint16(len(ns))
	dh.Arcount = uint16(len(extra))

	// We need the uncompressed length here, because we first pack it and then compress it.
	msg = buf
	compress := dns.Compress
	dns.Compress = false
	if packLen := dns.Len() + 1; len(msg) < packLen ***REMOVED***
		msg = make([]byte, packLen)
	***REMOVED***
	dns.Compress = compress

	// Pack it in: header and then the pieces.
	off := 0
	off, err = packStructCompress(&dh, msg, off, compression, dns.Compress)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for i := 0; i < len(question); i++ ***REMOVED***
		off, err = packStructCompress(&question[i], msg, off, compression, dns.Compress)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	for i := 0; i < len(answer); i++ ***REMOVED***
		off, err = PackRR(answer[i], msg, off, compression, dns.Compress)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	for i := 0; i < len(ns); i++ ***REMOVED***
		off, err = PackRR(ns[i], msg, off, compression, dns.Compress)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	for i := 0; i < len(extra); i++ ***REMOVED***
		off, err = PackRR(extra[i], msg, off, compression, dns.Compress)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return msg[:off], nil
***REMOVED***

// Unpack unpacks a binary message to a Msg structure.
func (dns *Msg) Unpack(msg []byte) (err error) ***REMOVED***
	// Header.
	var dh Header
	off := 0
	if off, err = UnpackStruct(&dh, msg, off); err != nil ***REMOVED***
		return err
	***REMOVED***
	dns.Id = dh.Id
	dns.Response = (dh.Bits & _QR) != 0
	dns.Opcode = int(dh.Bits>>11) & 0xF
	dns.Authoritative = (dh.Bits & _AA) != 0
	dns.Truncated = (dh.Bits & _TC) != 0
	dns.RecursionDesired = (dh.Bits & _RD) != 0
	dns.RecursionAvailable = (dh.Bits & _RA) != 0
	dns.Zero = (dh.Bits & _Z) != 0
	dns.AuthenticatedData = (dh.Bits & _AD) != 0
	dns.CheckingDisabled = (dh.Bits & _CD) != 0
	dns.Rcode = int(dh.Bits & 0xF)

	// Optimistically use the count given to us in the header
	dns.Question = make([]Question, 0, int(dh.Qdcount))

	var q Question
	for i := 0; i < int(dh.Qdcount); i++ ***REMOVED***
		off1 := off
		off, err = UnpackStruct(&q, msg, off)
		if err != nil ***REMOVED***
			// Even if Truncated is set, we only will set ErrTruncated if we
			// actually got the questions
			return err
		***REMOVED***
		if off1 == off ***REMOVED*** // Offset does not increase anymore, dh.Qdcount is a lie!
			dh.Qdcount = uint16(i)
			break
		***REMOVED***
		dns.Question = append(dns.Question, q)
	***REMOVED***

	dns.Answer, off, err = unpackRRslice(int(dh.Ancount), msg, off)
	// The header counts might have been wrong so we need to update it
	dh.Ancount = uint16(len(dns.Answer))
	if err == nil ***REMOVED***
		dns.Ns, off, err = unpackRRslice(int(dh.Nscount), msg, off)
	***REMOVED***
	// The header counts might have been wrong so we need to update it
	dh.Nscount = uint16(len(dns.Ns))
	if err == nil ***REMOVED***
		dns.Extra, off, err = unpackRRslice(int(dh.Arcount), msg, off)
	***REMOVED***
	// The header counts might have been wrong so we need to update it
	dh.Arcount = uint16(len(dns.Extra))
	if off != len(msg) ***REMOVED***
		// TODO(miek) make this an error?
		// use PackOpt to let people tell how detailed the error reporting should be?
		// println("dns: extra bytes in dns packet", off, "<", len(msg))
	***REMOVED*** else if dns.Truncated ***REMOVED***
		// Whether we ran into a an error or not, we want to return that it
		// was truncated
		err = ErrTruncated
	***REMOVED***
	return err
***REMOVED***

// Convert a complete message to a string with dig-like output.
func (dns *Msg) String() string ***REMOVED***
	if dns == nil ***REMOVED***
		return "<nil> MsgHdr"
	***REMOVED***
	s := dns.MsgHdr.String() + " "
	s += "QUERY: " + strconv.Itoa(len(dns.Question)) + ", "
	s += "ANSWER: " + strconv.Itoa(len(dns.Answer)) + ", "
	s += "AUTHORITY: " + strconv.Itoa(len(dns.Ns)) + ", "
	s += "ADDITIONAL: " + strconv.Itoa(len(dns.Extra)) + "\n"
	if len(dns.Question) > 0 ***REMOVED***
		s += "\n;; QUESTION SECTION:\n"
		for i := 0; i < len(dns.Question); i++ ***REMOVED***
			s += dns.Question[i].String() + "\n"
		***REMOVED***
	***REMOVED***
	if len(dns.Answer) > 0 ***REMOVED***
		s += "\n;; ANSWER SECTION:\n"
		for i := 0; i < len(dns.Answer); i++ ***REMOVED***
			if dns.Answer[i] != nil ***REMOVED***
				s += dns.Answer[i].String() + "\n"
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(dns.Ns) > 0 ***REMOVED***
		s += "\n;; AUTHORITY SECTION:\n"
		for i := 0; i < len(dns.Ns); i++ ***REMOVED***
			if dns.Ns[i] != nil ***REMOVED***
				s += dns.Ns[i].String() + "\n"
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(dns.Extra) > 0 ***REMOVED***
		s += "\n;; ADDITIONAL SECTION:\n"
		for i := 0; i < len(dns.Extra); i++ ***REMOVED***
			if dns.Extra[i] != nil ***REMOVED***
				s += dns.Extra[i].String() + "\n"
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***

// Len returns the message length when in (un)compressed wire format.
// If dns.Compress is true compression it is taken into account. Len()
// is provided to be a faster way to get the size of the resulting packet,
// than packing it, measuring the size and discarding the buffer.
func (dns *Msg) Len() int ***REMOVED***
	// We always return one more than needed.
	l := 12 // Message header is always 12 bytes
	var compression map[string]int
	if dns.Compress ***REMOVED***
		compression = make(map[string]int)
	***REMOVED***
	for i := 0; i < len(dns.Question); i++ ***REMOVED***
		l += dns.Question[i].len()
		if dns.Compress ***REMOVED***
			compressionLenHelper(compression, dns.Question[i].Name)
		***REMOVED***
	***REMOVED***
	for i := 0; i < len(dns.Answer); i++ ***REMOVED***
		l += dns.Answer[i].len()
		if dns.Compress ***REMOVED***
			k, ok := compressionLenSearch(compression, dns.Answer[i].Header().Name)
			if ok ***REMOVED***
				l += 1 - k
			***REMOVED***
			compressionLenHelper(compression, dns.Answer[i].Header().Name)
			k, ok = compressionLenSearchType(compression, dns.Answer[i])
			if ok ***REMOVED***
				l += 1 - k
			***REMOVED***
			compressionLenHelperType(compression, dns.Answer[i])
		***REMOVED***
	***REMOVED***
	for i := 0; i < len(dns.Ns); i++ ***REMOVED***
		l += dns.Ns[i].len()
		if dns.Compress ***REMOVED***
			k, ok := compressionLenSearch(compression, dns.Ns[i].Header().Name)
			if ok ***REMOVED***
				l += 1 - k
			***REMOVED***
			compressionLenHelper(compression, dns.Ns[i].Header().Name)
			k, ok = compressionLenSearchType(compression, dns.Ns[i])
			if ok ***REMOVED***
				l += 1 - k
			***REMOVED***
			compressionLenHelperType(compression, dns.Ns[i])
		***REMOVED***
	***REMOVED***
	for i := 0; i < len(dns.Extra); i++ ***REMOVED***
		l += dns.Extra[i].len()
		if dns.Compress ***REMOVED***
			k, ok := compressionLenSearch(compression, dns.Extra[i].Header().Name)
			if ok ***REMOVED***
				l += 1 - k
			***REMOVED***
			compressionLenHelper(compression, dns.Extra[i].Header().Name)
			k, ok = compressionLenSearchType(compression, dns.Extra[i])
			if ok ***REMOVED***
				l += 1 - k
			***REMOVED***
			compressionLenHelperType(compression, dns.Extra[i])
		***REMOVED***
	***REMOVED***
	return l
***REMOVED***

// Put the parts of the name in the compression map.
func compressionLenHelper(c map[string]int, s string) ***REMOVED***
	pref := ""
	lbs := Split(s)
	for j := len(lbs) - 1; j >= 0; j-- ***REMOVED***
		pref = s[lbs[j]:]
		if _, ok := c[pref]; !ok ***REMOVED***
			c[pref] = len(pref)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Look for each part in the compression map and returns its length,
// keep on searching so we get the longest match.
func compressionLenSearch(c map[string]int, s string) (int, bool) ***REMOVED***
	off := 0
	end := false
	if s == "" ***REMOVED*** // don't bork on bogus data
		return 0, false
	***REMOVED***
	for ***REMOVED***
		if _, ok := c[s[off:]]; ok ***REMOVED***
			return len(s[off:]), true
		***REMOVED***
		if end ***REMOVED***
			break
		***REMOVED***
		off, end = NextLabel(s, off)
	***REMOVED***
	return 0, false
***REMOVED***

// TODO(miek): should add all types, because the all can be *used* for compression.
func compressionLenHelperType(c map[string]int, r RR) ***REMOVED***
	switch x := r.(type) ***REMOVED***
	case *NS:
		compressionLenHelper(c, x.Ns)
	case *MX:
		compressionLenHelper(c, x.Mx)
	case *CNAME:
		compressionLenHelper(c, x.Target)
	case *PTR:
		compressionLenHelper(c, x.Ptr)
	case *SOA:
		compressionLenHelper(c, x.Ns)
		compressionLenHelper(c, x.Mbox)
	case *MB:
		compressionLenHelper(c, x.Mb)
	case *MG:
		compressionLenHelper(c, x.Mg)
	case *MR:
		compressionLenHelper(c, x.Mr)
	case *MF:
		compressionLenHelper(c, x.Mf)
	case *MD:
		compressionLenHelper(c, x.Md)
	case *RT:
		compressionLenHelper(c, x.Host)
	case *MINFO:
		compressionLenHelper(c, x.Rmail)
		compressionLenHelper(c, x.Email)
	case *AFSDB:
		compressionLenHelper(c, x.Hostname)
	***REMOVED***
***REMOVED***

// Only search on compressing these types.
func compressionLenSearchType(c map[string]int, r RR) (int, bool) ***REMOVED***
	switch x := r.(type) ***REMOVED***
	case *NS:
		return compressionLenSearch(c, x.Ns)
	case *MX:
		return compressionLenSearch(c, x.Mx)
	case *CNAME:
		return compressionLenSearch(c, x.Target)
	case *PTR:
		return compressionLenSearch(c, x.Ptr)
	case *SOA:
		k, ok := compressionLenSearch(c, x.Ns)
		k1, ok1 := compressionLenSearch(c, x.Mbox)
		if !ok && !ok1 ***REMOVED***
			return 0, false
		***REMOVED***
		return k + k1, true
	case *MB:
		return compressionLenSearch(c, x.Mb)
	case *MG:
		return compressionLenSearch(c, x.Mg)
	case *MR:
		return compressionLenSearch(c, x.Mr)
	case *MF:
		return compressionLenSearch(c, x.Mf)
	case *MD:
		return compressionLenSearch(c, x.Md)
	case *RT:
		return compressionLenSearch(c, x.Host)
	case *MINFO:
		k, ok := compressionLenSearch(c, x.Rmail)
		k1, ok1 := compressionLenSearch(c, x.Email)
		if !ok && !ok1 ***REMOVED***
			return 0, false
		***REMOVED***
		return k + k1, true
	case *AFSDB:
		return compressionLenSearch(c, x.Hostname)
	***REMOVED***
	return 0, false
***REMOVED***

// id returns a 16 bits random number to be used as a
// message id. The random provided should be good enough.
func id() uint16 ***REMOVED***
	return uint16(rand.Int()) ^ uint16(time.Now().Nanosecond())
***REMOVED***

// Copy returns a new RR which is a deep-copy of r.
func Copy(r RR) RR ***REMOVED***
	r1 := r.copy()
	return r1
***REMOVED***

// Copy returns a new *Msg which is a deep-copy of dns.
func (dns *Msg) Copy() *Msg ***REMOVED***
	return dns.CopyTo(new(Msg))
***REMOVED***

// CopyTo copies the contents to the provided message using a deep-copy and returns the copy.
func (dns *Msg) CopyTo(r1 *Msg) *Msg ***REMOVED***
	r1.MsgHdr = dns.MsgHdr
	r1.Compress = dns.Compress

	if len(dns.Question) > 0 ***REMOVED***
		r1.Question = make([]Question, len(dns.Question))
		copy(r1.Question, dns.Question) // TODO(miek): Question is an immutable value, ok to do a shallow-copy
	***REMOVED***

	rrArr := make([]RR, len(dns.Answer)+len(dns.Ns)+len(dns.Extra))
	var rri int

	if len(dns.Answer) > 0 ***REMOVED***
		rrbegin := rri
		for i := 0; i < len(dns.Answer); i++ ***REMOVED***
			rrArr[rri] = dns.Answer[i].copy()
			rri++
		***REMOVED***
		r1.Answer = rrArr[rrbegin:rri:rri]
	***REMOVED***

	if len(dns.Ns) > 0 ***REMOVED***
		rrbegin := rri
		for i := 0; i < len(dns.Ns); i++ ***REMOVED***
			rrArr[rri] = dns.Ns[i].copy()
			rri++
		***REMOVED***
		r1.Ns = rrArr[rrbegin:rri:rri]
	***REMOVED***

	if len(dns.Extra) > 0 ***REMOVED***
		rrbegin := rri
		for i := 0; i < len(dns.Extra); i++ ***REMOVED***
			rrArr[rri] = dns.Extra[i].copy()
			rri++
		***REMOVED***
		r1.Extra = rrArr[rrbegin:rri:rri]
	***REMOVED***

	return r1
***REMOVED***
