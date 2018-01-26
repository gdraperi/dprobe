package dns

import (
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type (
	// Type is a DNS type.
	Type uint16
	// Class is a DNS class.
	Class uint16
	// Name is a DNS domain name.
	Name string
)

// Packet formats

// Wire constants and supported types.
const (
	// valid RR_Header.Rrtype and Question.qtype

	TypeNone       uint16 = 0
	TypeA          uint16 = 1
	TypeNS         uint16 = 2
	TypeMD         uint16 = 3
	TypeMF         uint16 = 4
	TypeCNAME      uint16 = 5
	TypeSOA        uint16 = 6
	TypeMB         uint16 = 7
	TypeMG         uint16 = 8
	TypeMR         uint16 = 9
	TypeNULL       uint16 = 10
	TypeWKS        uint16 = 11
	TypePTR        uint16 = 12
	TypeHINFO      uint16 = 13
	TypeMINFO      uint16 = 14
	TypeMX         uint16 = 15
	TypeTXT        uint16 = 16
	TypeRP         uint16 = 17
	TypeAFSDB      uint16 = 18
	TypeX25        uint16 = 19
	TypeISDN       uint16 = 20
	TypeRT         uint16 = 21
	TypeNSAPPTR    uint16 = 23
	TypeSIG        uint16 = 24
	TypeKEY        uint16 = 25
	TypePX         uint16 = 26
	TypeGPOS       uint16 = 27
	TypeAAAA       uint16 = 28
	TypeLOC        uint16 = 29
	TypeNXT        uint16 = 30
	TypeEID        uint16 = 31
	TypeNIMLOC     uint16 = 32
	TypeSRV        uint16 = 33
	TypeATMA       uint16 = 34
	TypeNAPTR      uint16 = 35
	TypeKX         uint16 = 36
	TypeCERT       uint16 = 37
	TypeDNAME      uint16 = 39
	TypeOPT        uint16 = 41 // EDNS
	TypeDS         uint16 = 43
	TypeSSHFP      uint16 = 44
	TypeIPSECKEY   uint16 = 45
	TypeRRSIG      uint16 = 46
	TypeNSEC       uint16 = 47
	TypeDNSKEY     uint16 = 48
	TypeDHCID      uint16 = 49
	TypeNSEC3      uint16 = 50
	TypeNSEC3PARAM uint16 = 51
	TypeTLSA       uint16 = 52
	TypeHIP        uint16 = 55
	TypeNINFO      uint16 = 56
	TypeRKEY       uint16 = 57
	TypeTALINK     uint16 = 58
	TypeCDS        uint16 = 59
	TypeCDNSKEY    uint16 = 60
	TypeOPENPGPKEY uint16 = 61
	TypeSPF        uint16 = 99
	TypeUINFO      uint16 = 100
	TypeUID        uint16 = 101
	TypeGID        uint16 = 102
	TypeUNSPEC     uint16 = 103
	TypeNID        uint16 = 104
	TypeL32        uint16 = 105
	TypeL64        uint16 = 106
	TypeLP         uint16 = 107
	TypeEUI48      uint16 = 108
	TypeEUI64      uint16 = 109
	TypeURI        uint16 = 256
	TypeCAA        uint16 = 257

	TypeTKEY uint16 = 249
	TypeTSIG uint16 = 250

	// valid Question.Qtype only
	TypeIXFR  uint16 = 251
	TypeAXFR  uint16 = 252
	TypeMAILB uint16 = 253
	TypeMAILA uint16 = 254
	TypeANY   uint16 = 255

	TypeTA       uint16 = 32768
	TypeDLV      uint16 = 32769
	TypeReserved uint16 = 65535

	// valid Question.Qclass
	ClassINET   = 1
	ClassCSNET  = 2
	ClassCHAOS  = 3
	ClassHESIOD = 4
	ClassNONE   = 254
	ClassANY    = 255

	// Message Response Codes.
	RcodeSuccess        = 0
	RcodeFormatError    = 1
	RcodeServerFailure  = 2
	RcodeNameError      = 3
	RcodeNotImplemented = 4
	RcodeRefused        = 5
	RcodeYXDomain       = 6
	RcodeYXRrset        = 7
	RcodeNXRrset        = 8
	RcodeNotAuth        = 9
	RcodeNotZone        = 10
	RcodeBadSig         = 16 // TSIG
	RcodeBadVers        = 16 // EDNS0
	RcodeBadKey         = 17
	RcodeBadTime        = 18
	RcodeBadMode        = 19 // TKEY
	RcodeBadName        = 20
	RcodeBadAlg         = 21
	RcodeBadTrunc       = 22 // TSIG

	// Message Opcodes. There is no 3.
	OpcodeQuery  = 0
	OpcodeIQuery = 1
	OpcodeStatus = 2
	OpcodeNotify = 4
	OpcodeUpdate = 5
)

// Headers is the wire format for the DNS packet header.
type Header struct ***REMOVED***
	Id                                 uint16
	Bits                               uint16
	Qdcount, Ancount, Nscount, Arcount uint16
***REMOVED***

const (
	headerSize = 12

	// Header.Bits
	_QR = 1 << 15 // query/response (response=1)
	_AA = 1 << 10 // authoritative
	_TC = 1 << 9  // truncated
	_RD = 1 << 8  // recursion desired
	_RA = 1 << 7  // recursion available
	_Z  = 1 << 6  // Z
	_AD = 1 << 5  // authticated data
	_CD = 1 << 4  // checking disabled

	LOC_EQUATOR       = 1 << 31 // RFC 1876, Section 2.
	LOC_PRIMEMERIDIAN = 1 << 31 // RFC 1876, Section 2.

	LOC_HOURS   = 60 * 1000
	LOC_DEGREES = 60 * LOC_HOURS

	LOC_ALTITUDEBASE = 100000
)

// Different Certificate Types, see RFC 4398, Section 2.1
const (
	CertPKIX = 1 + iota
	CertSPKI
	CertPGP
	CertIPIX
	CertISPKI
	CertIPGP
	CertACPKIX
	CertIACPKIX
	CertURI = 253
	CertOID = 254
)

// CertTypeToString converts the Cert Type to its string representation.
// See RFC 4398 and RFC 6944.
var CertTypeToString = map[uint16]string***REMOVED***
	CertPKIX:    "PKIX",
	CertSPKI:    "SPKI",
	CertPGP:     "PGP",
	CertIPIX:    "IPIX",
	CertISPKI:   "ISPKI",
	CertIPGP:    "IPGP",
	CertACPKIX:  "ACPKIX",
	CertIACPKIX: "IACPKIX",
	CertURI:     "URI",
	CertOID:     "OID",
***REMOVED***

// StringToCertType is the reverseof CertTypeToString.
var StringToCertType = reverseInt16(CertTypeToString)

//go:generate go run types_generate.go

// Question holds a DNS question. There can be multiple questions in the
// question section of a message. Usually there is just one.
type Question struct ***REMOVED***
	Name   string `dns:"cdomain-name"` // "cdomain-name" specifies encoding (and may be compressed)
	Qtype  uint16
	Qclass uint16
***REMOVED***

func (q *Question) len() int ***REMOVED***
	return len(q.Name) + 1 + 2 + 2
***REMOVED***

func (q *Question) String() (s string) ***REMOVED***
	// prefix with ; (as in dig)
	s = ";" + sprintName(q.Name) + "\t"
	s += Class(q.Qclass).String() + "\t"
	s += " " + Type(q.Qtype).String()
	return s
***REMOVED***

// ANY is a wildcard record. See RFC 1035, Section 3.2.3. ANY
// is named "*" there.
type ANY struct ***REMOVED***
	Hdr RR_Header
	// Does not have any rdata
***REMOVED***

func (rr *ANY) String() string ***REMOVED*** return rr.Hdr.String() ***REMOVED***

type CNAME struct ***REMOVED***
	Hdr    RR_Header
	Target string `dns:"cdomain-name"`
***REMOVED***

func (rr *CNAME) String() string ***REMOVED*** return rr.Hdr.String() + sprintName(rr.Target) ***REMOVED***

type HINFO struct ***REMOVED***
	Hdr RR_Header
	Cpu string
	Os  string
***REMOVED***

func (rr *HINFO) String() string ***REMOVED***
	return rr.Hdr.String() + sprintTxt([]string***REMOVED***rr.Cpu, rr.Os***REMOVED***)
***REMOVED***

type MB struct ***REMOVED***
	Hdr RR_Header
	Mb  string `dns:"cdomain-name"`
***REMOVED***

func (rr *MB) String() string ***REMOVED*** return rr.Hdr.String() + sprintName(rr.Mb) ***REMOVED***

type MG struct ***REMOVED***
	Hdr RR_Header
	Mg  string `dns:"cdomain-name"`
***REMOVED***

func (rr *MG) String() string ***REMOVED*** return rr.Hdr.String() + sprintName(rr.Mg) ***REMOVED***

type MINFO struct ***REMOVED***
	Hdr   RR_Header
	Rmail string `dns:"cdomain-name"`
	Email string `dns:"cdomain-name"`
***REMOVED***

func (rr *MINFO) String() string ***REMOVED***
	return rr.Hdr.String() + sprintName(rr.Rmail) + " " + sprintName(rr.Email)
***REMOVED***

type MR struct ***REMOVED***
	Hdr RR_Header
	Mr  string `dns:"cdomain-name"`
***REMOVED***

func (rr *MR) String() string ***REMOVED***
	return rr.Hdr.String() + sprintName(rr.Mr)
***REMOVED***

type MF struct ***REMOVED***
	Hdr RR_Header
	Mf  string `dns:"cdomain-name"`
***REMOVED***

func (rr *MF) String() string ***REMOVED***
	return rr.Hdr.String() + sprintName(rr.Mf)
***REMOVED***

type MD struct ***REMOVED***
	Hdr RR_Header
	Md  string `dns:"cdomain-name"`
***REMOVED***

func (rr *MD) String() string ***REMOVED***
	return rr.Hdr.String() + sprintName(rr.Md)
***REMOVED***

type MX struct ***REMOVED***
	Hdr        RR_Header
	Preference uint16
	Mx         string `dns:"cdomain-name"`
***REMOVED***

func (rr *MX) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Preference)) + " " + sprintName(rr.Mx)
***REMOVED***

type AFSDB struct ***REMOVED***
	Hdr      RR_Header
	Subtype  uint16
	Hostname string `dns:"cdomain-name"`
***REMOVED***

func (rr *AFSDB) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Subtype)) + " " + sprintName(rr.Hostname)
***REMOVED***

type X25 struct ***REMOVED***
	Hdr         RR_Header
	PSDNAddress string
***REMOVED***

func (rr *X25) String() string ***REMOVED***
	return rr.Hdr.String() + rr.PSDNAddress
***REMOVED***

type RT struct ***REMOVED***
	Hdr        RR_Header
	Preference uint16
	Host       string `dns:"cdomain-name"`
***REMOVED***

func (rr *RT) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Preference)) + " " + sprintName(rr.Host)
***REMOVED***

type NS struct ***REMOVED***
	Hdr RR_Header
	Ns  string `dns:"cdomain-name"`
***REMOVED***

func (rr *NS) String() string ***REMOVED***
	return rr.Hdr.String() + sprintName(rr.Ns)
***REMOVED***

type PTR struct ***REMOVED***
	Hdr RR_Header
	Ptr string `dns:"cdomain-name"`
***REMOVED***

func (rr *PTR) String() string ***REMOVED***
	return rr.Hdr.String() + sprintName(rr.Ptr)
***REMOVED***

type RP struct ***REMOVED***
	Hdr  RR_Header
	Mbox string `dns:"domain-name"`
	Txt  string `dns:"domain-name"`
***REMOVED***

func (rr *RP) String() string ***REMOVED***
	return rr.Hdr.String() + rr.Mbox + " " + sprintTxt([]string***REMOVED***rr.Txt***REMOVED***)
***REMOVED***

type SOA struct ***REMOVED***
	Hdr     RR_Header
	Ns      string `dns:"cdomain-name"`
	Mbox    string `dns:"cdomain-name"`
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minttl  uint32
***REMOVED***

func (rr *SOA) String() string ***REMOVED***
	return rr.Hdr.String() + sprintName(rr.Ns) + " " + sprintName(rr.Mbox) +
		" " + strconv.FormatInt(int64(rr.Serial), 10) +
		" " + strconv.FormatInt(int64(rr.Refresh), 10) +
		" " + strconv.FormatInt(int64(rr.Retry), 10) +
		" " + strconv.FormatInt(int64(rr.Expire), 10) +
		" " + strconv.FormatInt(int64(rr.Minttl), 10)
***REMOVED***

type TXT struct ***REMOVED***
	Hdr RR_Header
	Txt []string `dns:"txt"`
***REMOVED***

func (rr *TXT) String() string ***REMOVED*** return rr.Hdr.String() + sprintTxt(rr.Txt) ***REMOVED***

func sprintName(s string) string ***REMOVED***
	src := []byte(s)
	dst := make([]byte, 0, len(src))
	for i := 0; i < len(src); ***REMOVED***
		if i+1 < len(src) && src[i] == '\\' && src[i+1] == '.' ***REMOVED***
			dst = append(dst, src[i:i+2]...)
			i += 2
		***REMOVED*** else ***REMOVED***
			b, n := nextByte(src, i)
			if n == 0 ***REMOVED***
				i++ // dangling back slash
			***REMOVED*** else if b == '.' ***REMOVED***
				dst = append(dst, b)
			***REMOVED*** else ***REMOVED***
				dst = appendDomainNameByte(dst, b)
			***REMOVED***
			i += n
		***REMOVED***
	***REMOVED***
	return string(dst)
***REMOVED***

func sprintTxtOctet(s string) string ***REMOVED***
	src := []byte(s)
	dst := make([]byte, 0, len(src))
	dst = append(dst, '"')
	for i := 0; i < len(src); ***REMOVED***
		if i+1 < len(src) && src[i] == '\\' && src[i+1] == '.' ***REMOVED***
			dst = append(dst, src[i:i+2]...)
			i += 2
		***REMOVED*** else ***REMOVED***
			b, n := nextByte(src, i)
			if n == 0 ***REMOVED***
				i++ // dangling back slash
			***REMOVED*** else if b == '.' ***REMOVED***
				dst = append(dst, b)
			***REMOVED*** else ***REMOVED***
				if b < ' ' || b > '~' ***REMOVED***
					dst = appendByte(dst, b)
				***REMOVED*** else ***REMOVED***
					dst = append(dst, b)
				***REMOVED***
			***REMOVED***
			i += n
		***REMOVED***
	***REMOVED***
	dst = append(dst, '"')
	return string(dst)
***REMOVED***

func sprintTxt(txt []string) string ***REMOVED***
	var out []byte
	for i, s := range txt ***REMOVED***
		if i > 0 ***REMOVED***
			out = append(out, ` "`...)
		***REMOVED*** else ***REMOVED***
			out = append(out, '"')
		***REMOVED***
		bs := []byte(s)
		for j := 0; j < len(bs); ***REMOVED***
			b, n := nextByte(bs, j)
			if n == 0 ***REMOVED***
				break
			***REMOVED***
			out = appendTXTStringByte(out, b)
			j += n
		***REMOVED***
		out = append(out, '"')
	***REMOVED***
	return string(out)
***REMOVED***

func appendDomainNameByte(s []byte, b byte) []byte ***REMOVED***
	switch b ***REMOVED***
	case '.', ' ', '\'', '@', ';', '(', ')': // additional chars to escape
		return append(s, '\\', b)
	***REMOVED***
	return appendTXTStringByte(s, b)
***REMOVED***

func appendTXTStringByte(s []byte, b byte) []byte ***REMOVED***
	switch b ***REMOVED***
	case '\t':
		return append(s, '\\', 't')
	case '\r':
		return append(s, '\\', 'r')
	case '\n':
		return append(s, '\\', 'n')
	case '"', '\\':
		return append(s, '\\', b)
	***REMOVED***
	if b < ' ' || b > '~' ***REMOVED***
		return appendByte(s, b)
	***REMOVED***
	return append(s, b)
***REMOVED***

func appendByte(s []byte, b byte) []byte ***REMOVED***
	var buf [3]byte
	bufs := strconv.AppendInt(buf[:0], int64(b), 10)
	s = append(s, '\\')
	for i := 0; i < 3-len(bufs); i++ ***REMOVED***
		s = append(s, '0')
	***REMOVED***
	for _, r := range bufs ***REMOVED***
		s = append(s, r)
	***REMOVED***
	return s
***REMOVED***

func nextByte(b []byte, offset int) (byte, int) ***REMOVED***
	if offset >= len(b) ***REMOVED***
		return 0, 0
	***REMOVED***
	if b[offset] != '\\' ***REMOVED***
		// not an escape sequence
		return b[offset], 1
	***REMOVED***
	switch len(b) - offset ***REMOVED***
	case 1: // dangling escape
		return 0, 0
	case 2, 3: // too short to be \ddd
	default: // maybe \ddd
		if isDigit(b[offset+1]) && isDigit(b[offset+2]) && isDigit(b[offset+3]) ***REMOVED***
			return dddToByte(b[offset+1:]), 4
		***REMOVED***
	***REMOVED***
	// not \ddd, maybe a control char
	switch b[offset+1] ***REMOVED***
	case 't':
		return '\t', 2
	case 'r':
		return '\r', 2
	case 'n':
		return '\n', 2
	default:
		return b[offset+1], 2
	***REMOVED***
***REMOVED***

type SPF struct ***REMOVED***
	Hdr RR_Header
	Txt []string `dns:"txt"`
***REMOVED***

func (rr *SPF) String() string ***REMOVED*** return rr.Hdr.String() + sprintTxt(rr.Txt) ***REMOVED***

type SRV struct ***REMOVED***
	Hdr      RR_Header
	Priority uint16
	Weight   uint16
	Port     uint16
	Target   string `dns:"domain-name"`
***REMOVED***

func (rr *SRV) String() string ***REMOVED***
	return rr.Hdr.String() +
		strconv.Itoa(int(rr.Priority)) + " " +
		strconv.Itoa(int(rr.Weight)) + " " +
		strconv.Itoa(int(rr.Port)) + " " + sprintName(rr.Target)
***REMOVED***

type NAPTR struct ***REMOVED***
	Hdr         RR_Header
	Order       uint16
	Preference  uint16
	Flags       string
	Service     string
	Regexp      string
	Replacement string `dns:"domain-name"`
***REMOVED***

func (rr *NAPTR) String() string ***REMOVED***
	return rr.Hdr.String() +
		strconv.Itoa(int(rr.Order)) + " " +
		strconv.Itoa(int(rr.Preference)) + " " +
		"\"" + rr.Flags + "\" " +
		"\"" + rr.Service + "\" " +
		"\"" + rr.Regexp + "\" " +
		rr.Replacement
***REMOVED***

// The CERT resource record, see RFC 4398.
type CERT struct ***REMOVED***
	Hdr         RR_Header
	Type        uint16
	KeyTag      uint16
	Algorithm   uint8
	Certificate string `dns:"base64"`
***REMOVED***

func (rr *CERT) String() string ***REMOVED***
	var (
		ok                  bool
		certtype, algorithm string
	)
	if certtype, ok = CertTypeToString[rr.Type]; !ok ***REMOVED***
		certtype = strconv.Itoa(int(rr.Type))
	***REMOVED***
	if algorithm, ok = AlgorithmToString[rr.Algorithm]; !ok ***REMOVED***
		algorithm = strconv.Itoa(int(rr.Algorithm))
	***REMOVED***
	return rr.Hdr.String() + certtype +
		" " + strconv.Itoa(int(rr.KeyTag)) +
		" " + algorithm +
		" " + rr.Certificate
***REMOVED***

// The DNAME resource record, see RFC 2672.
type DNAME struct ***REMOVED***
	Hdr    RR_Header
	Target string `dns:"domain-name"`
***REMOVED***

func (rr *DNAME) String() string ***REMOVED***
	return rr.Hdr.String() + sprintName(rr.Target)
***REMOVED***

type A struct ***REMOVED***
	Hdr RR_Header
	A   net.IP `dns:"a"`
***REMOVED***

func (rr *A) String() string ***REMOVED***
	if rr.A == nil ***REMOVED***
		return rr.Hdr.String()
	***REMOVED***
	return rr.Hdr.String() + rr.A.String()
***REMOVED***

type AAAA struct ***REMOVED***
	Hdr  RR_Header
	AAAA net.IP `dns:"aaaa"`
***REMOVED***

func (rr *AAAA) String() string ***REMOVED***
	if rr.AAAA == nil ***REMOVED***
		return rr.Hdr.String()
	***REMOVED***
	return rr.Hdr.String() + rr.AAAA.String()
***REMOVED***

type PX struct ***REMOVED***
	Hdr        RR_Header
	Preference uint16
	Map822     string `dns:"domain-name"`
	Mapx400    string `dns:"domain-name"`
***REMOVED***

func (rr *PX) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Preference)) + " " + sprintName(rr.Map822) + " " + sprintName(rr.Mapx400)
***REMOVED***

type GPOS struct ***REMOVED***
	Hdr       RR_Header
	Longitude string
	Latitude  string
	Altitude  string
***REMOVED***

func (rr *GPOS) String() string ***REMOVED***
	return rr.Hdr.String() + rr.Longitude + " " + rr.Latitude + " " + rr.Altitude
***REMOVED***

type LOC struct ***REMOVED***
	Hdr       RR_Header
	Version   uint8
	Size      uint8
	HorizPre  uint8
	VertPre   uint8
	Latitude  uint32
	Longitude uint32
	Altitude  uint32
***REMOVED***

// cmToM takes a cm value expressed in RFC1876 SIZE mantissa/exponent
// format and returns a string in m (two decimals for the cm)
func cmToM(m, e uint8) string ***REMOVED***
	if e < 2 ***REMOVED***
		if e == 1 ***REMOVED***
			m *= 10
		***REMOVED***

		return fmt.Sprintf("0.%02d", m)
	***REMOVED***

	s := fmt.Sprintf("%d", m)
	for e > 2 ***REMOVED***
		s += "0"
		e--
	***REMOVED***
	return s
***REMOVED***

func (rr *LOC) String() string ***REMOVED***
	s := rr.Hdr.String()

	lat := rr.Latitude
	ns := "N"
	if lat > LOC_EQUATOR ***REMOVED***
		lat = lat - LOC_EQUATOR
	***REMOVED*** else ***REMOVED***
		ns = "S"
		lat = LOC_EQUATOR - lat
	***REMOVED***
	h := lat / LOC_DEGREES
	lat = lat % LOC_DEGREES
	m := lat / LOC_HOURS
	lat = lat % LOC_HOURS
	s += fmt.Sprintf("%02d %02d %0.3f %s ", h, m, (float64(lat) / 1000), ns)

	lon := rr.Longitude
	ew := "E"
	if lon > LOC_PRIMEMERIDIAN ***REMOVED***
		lon = lon - LOC_PRIMEMERIDIAN
	***REMOVED*** else ***REMOVED***
		ew = "W"
		lon = LOC_PRIMEMERIDIAN - lon
	***REMOVED***
	h = lon / LOC_DEGREES
	lon = lon % LOC_DEGREES
	m = lon / LOC_HOURS
	lon = lon % LOC_HOURS
	s += fmt.Sprintf("%02d %02d %0.3f %s ", h, m, (float64(lon) / 1000), ew)

	var alt = float64(rr.Altitude) / 100
	alt -= LOC_ALTITUDEBASE
	if rr.Altitude%100 != 0 ***REMOVED***
		s += fmt.Sprintf("%.2fm ", alt)
	***REMOVED*** else ***REMOVED***
		s += fmt.Sprintf("%.0fm ", alt)
	***REMOVED***

	s += cmToM((rr.Size&0xf0)>>4, rr.Size&0x0f) + "m "
	s += cmToM((rr.HorizPre&0xf0)>>4, rr.HorizPre&0x0f) + "m "
	s += cmToM((rr.VertPre&0xf0)>>4, rr.VertPre&0x0f) + "m"

	return s
***REMOVED***

// SIG is identical to RRSIG and nowadays only used for SIG(0), RFC2931.
type SIG struct ***REMOVED***
	RRSIG
***REMOVED***

type RRSIG struct ***REMOVED***
	Hdr         RR_Header
	TypeCovered uint16
	Algorithm   uint8
	Labels      uint8
	OrigTtl     uint32
	Expiration  uint32
	Inception   uint32
	KeyTag      uint16
	SignerName  string `dns:"domain-name"`
	Signature   string `dns:"base64"`
***REMOVED***

func (rr *RRSIG) String() string ***REMOVED***
	s := rr.Hdr.String()
	s += Type(rr.TypeCovered).String()
	s += " " + strconv.Itoa(int(rr.Algorithm)) +
		" " + strconv.Itoa(int(rr.Labels)) +
		" " + strconv.FormatInt(int64(rr.OrigTtl), 10) +
		" " + TimeToString(rr.Expiration) +
		" " + TimeToString(rr.Inception) +
		" " + strconv.Itoa(int(rr.KeyTag)) +
		" " + sprintName(rr.SignerName) +
		" " + rr.Signature
	return s
***REMOVED***

type NSEC struct ***REMOVED***
	Hdr        RR_Header
	NextDomain string   `dns:"domain-name"`
	TypeBitMap []uint16 `dns:"nsec"`
***REMOVED***

func (rr *NSEC) String() string ***REMOVED***
	s := rr.Hdr.String() + sprintName(rr.NextDomain)
	for i := 0; i < len(rr.TypeBitMap); i++ ***REMOVED***
		s += " " + Type(rr.TypeBitMap[i]).String()
	***REMOVED***
	return s
***REMOVED***

func (rr *NSEC) len() int ***REMOVED***
	l := rr.Hdr.len() + len(rr.NextDomain) + 1
	lastwindow := uint32(2 ^ 32 + 1)
	for _, t := range rr.TypeBitMap ***REMOVED***
		window := t / 256
		if uint32(window) != lastwindow ***REMOVED***
			l += 1 + 32
		***REMOVED***
		lastwindow = uint32(window)
	***REMOVED***
	return l
***REMOVED***

type DLV struct ***REMOVED***
	DS
***REMOVED***

type CDS struct ***REMOVED***
	DS
***REMOVED***

type DS struct ***REMOVED***
	Hdr        RR_Header
	KeyTag     uint16
	Algorithm  uint8
	DigestType uint8
	Digest     string `dns:"hex"`
***REMOVED***

func (rr *DS) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.KeyTag)) +
		" " + strconv.Itoa(int(rr.Algorithm)) +
		" " + strconv.Itoa(int(rr.DigestType)) +
		" " + strings.ToUpper(rr.Digest)
***REMOVED***

type KX struct ***REMOVED***
	Hdr        RR_Header
	Preference uint16
	Exchanger  string `dns:"domain-name"`
***REMOVED***

func (rr *KX) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Preference)) +
		" " + sprintName(rr.Exchanger)
***REMOVED***

type TA struct ***REMOVED***
	Hdr        RR_Header
	KeyTag     uint16
	Algorithm  uint8
	DigestType uint8
	Digest     string `dns:"hex"`
***REMOVED***

func (rr *TA) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.KeyTag)) +
		" " + strconv.Itoa(int(rr.Algorithm)) +
		" " + strconv.Itoa(int(rr.DigestType)) +
		" " + strings.ToUpper(rr.Digest)
***REMOVED***

type TALINK struct ***REMOVED***
	Hdr          RR_Header
	PreviousName string `dns:"domain-name"`
	NextName     string `dns:"domain-name"`
***REMOVED***

func (rr *TALINK) String() string ***REMOVED***
	return rr.Hdr.String() +
		sprintName(rr.PreviousName) + " " + sprintName(rr.NextName)
***REMOVED***

type SSHFP struct ***REMOVED***
	Hdr         RR_Header
	Algorithm   uint8
	Type        uint8
	FingerPrint string `dns:"hex"`
***REMOVED***

func (rr *SSHFP) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Algorithm)) +
		" " + strconv.Itoa(int(rr.Type)) +
		" " + strings.ToUpper(rr.FingerPrint)
***REMOVED***

type IPSECKEY struct ***REMOVED***
	Hdr        RR_Header
	Precedence uint8
	// GatewayType: 1: A record, 2: AAAA record, 3: domainname.
	// 0 is use for no type and GatewayName should be "." then.
	GatewayType uint8
	Algorithm   uint8
	// Gateway can be an A record, AAAA record or a domain name.
	GatewayA    net.IP `dns:"a"`
	GatewayAAAA net.IP `dns:"aaaa"`
	GatewayName string `dns:"domain-name"`
	PublicKey   string `dns:"base64"`
***REMOVED***

func (rr *IPSECKEY) String() string ***REMOVED***
	s := rr.Hdr.String() + strconv.Itoa(int(rr.Precedence)) +
		" " + strconv.Itoa(int(rr.GatewayType)) +
		" " + strconv.Itoa(int(rr.Algorithm))
	switch rr.GatewayType ***REMOVED***
	case 0:
		fallthrough
	case 3:
		s += " " + rr.GatewayName
	case 1:
		s += " " + rr.GatewayA.String()
	case 2:
		s += " " + rr.GatewayAAAA.String()
	default:
		s += " ."
	***REMOVED***
	s += " " + rr.PublicKey
	return s
***REMOVED***

func (rr *IPSECKEY) len() int ***REMOVED***
	l := rr.Hdr.len() + 3 + 1
	switch rr.GatewayType ***REMOVED***
	default:
		fallthrough
	case 0:
		fallthrough
	case 3:
		l += len(rr.GatewayName)
	case 1:
		l += 4
	case 2:
		l += 16
	***REMOVED***
	return l + base64.StdEncoding.DecodedLen(len(rr.PublicKey))
***REMOVED***

type KEY struct ***REMOVED***
	DNSKEY
***REMOVED***

type CDNSKEY struct ***REMOVED***
	DNSKEY
***REMOVED***

type DNSKEY struct ***REMOVED***
	Hdr       RR_Header
	Flags     uint16
	Protocol  uint8
	Algorithm uint8
	PublicKey string `dns:"base64"`
***REMOVED***

func (rr *DNSKEY) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Flags)) +
		" " + strconv.Itoa(int(rr.Protocol)) +
		" " + strconv.Itoa(int(rr.Algorithm)) +
		" " + rr.PublicKey
***REMOVED***

type RKEY struct ***REMOVED***
	Hdr       RR_Header
	Flags     uint16
	Protocol  uint8
	Algorithm uint8
	PublicKey string `dns:"base64"`
***REMOVED***

func (rr *RKEY) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Flags)) +
		" " + strconv.Itoa(int(rr.Protocol)) +
		" " + strconv.Itoa(int(rr.Algorithm)) +
		" " + rr.PublicKey
***REMOVED***

type NSAPPTR struct ***REMOVED***
	Hdr RR_Header
	Ptr string `dns:"domain-name"`
***REMOVED***

func (rr *NSAPPTR) String() string ***REMOVED*** return rr.Hdr.String() + sprintName(rr.Ptr) ***REMOVED***

type NSEC3 struct ***REMOVED***
	Hdr        RR_Header
	Hash       uint8
	Flags      uint8
	Iterations uint16
	SaltLength uint8
	Salt       string `dns:"size-hex"`
	HashLength uint8
	NextDomain string   `dns:"size-base32"`
	TypeBitMap []uint16 `dns:"nsec"`
***REMOVED***

func (rr *NSEC3) String() string ***REMOVED***
	s := rr.Hdr.String()
	s += strconv.Itoa(int(rr.Hash)) +
		" " + strconv.Itoa(int(rr.Flags)) +
		" " + strconv.Itoa(int(rr.Iterations)) +
		" " + saltToString(rr.Salt) +
		" " + rr.NextDomain
	for i := 0; i < len(rr.TypeBitMap); i++ ***REMOVED***
		s += " " + Type(rr.TypeBitMap[i]).String()
	***REMOVED***
	return s
***REMOVED***

func (rr *NSEC3) len() int ***REMOVED***
	l := rr.Hdr.len() + 6 + len(rr.Salt)/2 + 1 + len(rr.NextDomain) + 1
	lastwindow := uint32(2 ^ 32 + 1)
	for _, t := range rr.TypeBitMap ***REMOVED***
		window := t / 256
		if uint32(window) != lastwindow ***REMOVED***
			l += 1 + 32
		***REMOVED***
		lastwindow = uint32(window)
	***REMOVED***
	return l
***REMOVED***

type NSEC3PARAM struct ***REMOVED***
	Hdr        RR_Header
	Hash       uint8
	Flags      uint8
	Iterations uint16
	SaltLength uint8
	Salt       string `dns:"hex"`
***REMOVED***

func (rr *NSEC3PARAM) String() string ***REMOVED***
	s := rr.Hdr.String()
	s += strconv.Itoa(int(rr.Hash)) +
		" " + strconv.Itoa(int(rr.Flags)) +
		" " + strconv.Itoa(int(rr.Iterations)) +
		" " + saltToString(rr.Salt)
	return s
***REMOVED***

type TKEY struct ***REMOVED***
	Hdr        RR_Header
	Algorithm  string `dns:"domain-name"`
	Inception  uint32
	Expiration uint32
	Mode       uint16
	Error      uint16
	KeySize    uint16
	Key        string
	OtherLen   uint16
	OtherData  string
***REMOVED***

func (rr *TKEY) String() string ***REMOVED***
	// It has no presentation format
	return ""
***REMOVED***

// RFC3597 represents an unknown/generic RR.
type RFC3597 struct ***REMOVED***
	Hdr   RR_Header
	Rdata string `dns:"hex"`
***REMOVED***

func (rr *RFC3597) String() string ***REMOVED***
	// Let's call it a hack
	s := rfc3597Header(rr.Hdr)

	s += "\\# " + strconv.Itoa(len(rr.Rdata)/2) + " " + rr.Rdata
	return s
***REMOVED***

func rfc3597Header(h RR_Header) string ***REMOVED***
	var s string

	s += sprintName(h.Name) + "\t"
	s += strconv.FormatInt(int64(h.Ttl), 10) + "\t"
	s += "CLASS" + strconv.Itoa(int(h.Class)) + "\t"
	s += "TYPE" + strconv.Itoa(int(h.Rrtype)) + "\t"
	return s
***REMOVED***

type URI struct ***REMOVED***
	Hdr      RR_Header
	Priority uint16
	Weight   uint16
	Target   string `dns:"octet"`
***REMOVED***

func (rr *URI) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Priority)) +
		" " + strconv.Itoa(int(rr.Weight)) + " " + sprintTxtOctet(rr.Target)
***REMOVED***

type DHCID struct ***REMOVED***
	Hdr    RR_Header
	Digest string `dns:"base64"`
***REMOVED***

func (rr *DHCID) String() string ***REMOVED*** return rr.Hdr.String() + rr.Digest ***REMOVED***

type TLSA struct ***REMOVED***
	Hdr          RR_Header
	Usage        uint8
	Selector     uint8
	MatchingType uint8
	Certificate  string `dns:"hex"`
***REMOVED***

func (rr *TLSA) String() string ***REMOVED***
	return rr.Hdr.String() +
		strconv.Itoa(int(rr.Usage)) +
		" " + strconv.Itoa(int(rr.Selector)) +
		" " + strconv.Itoa(int(rr.MatchingType)) +
		" " + rr.Certificate
***REMOVED***

type HIP struct ***REMOVED***
	Hdr                RR_Header
	HitLength          uint8
	PublicKeyAlgorithm uint8
	PublicKeyLength    uint16
	Hit                string   `dns:"hex"`
	PublicKey          string   `dns:"base64"`
	RendezvousServers  []string `dns:"domain-name"`
***REMOVED***

func (rr *HIP) String() string ***REMOVED***
	s := rr.Hdr.String() +
		strconv.Itoa(int(rr.PublicKeyAlgorithm)) +
		" " + rr.Hit +
		" " + rr.PublicKey
	for _, d := range rr.RendezvousServers ***REMOVED***
		s += " " + sprintName(d)
	***REMOVED***
	return s
***REMOVED***

type NINFO struct ***REMOVED***
	Hdr    RR_Header
	ZSData []string `dns:"txt"`
***REMOVED***

func (rr *NINFO) String() string ***REMOVED*** return rr.Hdr.String() + sprintTxt(rr.ZSData) ***REMOVED***

type WKS struct ***REMOVED***
	Hdr      RR_Header
	Address  net.IP `dns:"a"`
	Protocol uint8
	BitMap   []uint16 `dns:"wks"`
***REMOVED***

func (rr *WKS) len() int ***REMOVED***
	// TODO: this is missing something...
	return rr.Hdr.len() + net.IPv4len + 1
***REMOVED***

func (rr *WKS) String() (s string) ***REMOVED***
	s = rr.Hdr.String()
	if rr.Address != nil ***REMOVED***
		s += rr.Address.String()
	***REMOVED***
	// TODO(miek): missing protocol here, see /etc/protocols
	for i := 0; i < len(rr.BitMap); i++ ***REMOVED***
		// should lookup the port
		s += " " + strconv.Itoa(int(rr.BitMap[i]))
	***REMOVED***
	return s
***REMOVED***

type NID struct ***REMOVED***
	Hdr        RR_Header
	Preference uint16
	NodeID     uint64
***REMOVED***

func (rr *NID) String() string ***REMOVED***
	s := rr.Hdr.String() + strconv.Itoa(int(rr.Preference))
	node := fmt.Sprintf("%0.16x", rr.NodeID)
	s += " " + node[0:4] + ":" + node[4:8] + ":" + node[8:12] + ":" + node[12:16]
	return s
***REMOVED***

type L32 struct ***REMOVED***
	Hdr        RR_Header
	Preference uint16
	Locator32  net.IP `dns:"a"`
***REMOVED***

func (rr *L32) String() string ***REMOVED***
	if rr.Locator32 == nil ***REMOVED***
		return rr.Hdr.String() + strconv.Itoa(int(rr.Preference))
	***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Preference)) +
		" " + rr.Locator32.String()
***REMOVED***

type L64 struct ***REMOVED***
	Hdr        RR_Header
	Preference uint16
	Locator64  uint64
***REMOVED***

func (rr *L64) String() string ***REMOVED***
	s := rr.Hdr.String() + strconv.Itoa(int(rr.Preference))
	node := fmt.Sprintf("%0.16X", rr.Locator64)
	s += " " + node[0:4] + ":" + node[4:8] + ":" + node[8:12] + ":" + node[12:16]
	return s
***REMOVED***

type LP struct ***REMOVED***
	Hdr        RR_Header
	Preference uint16
	Fqdn       string `dns:"domain-name"`
***REMOVED***

func (rr *LP) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Preference)) + " " + sprintName(rr.Fqdn)
***REMOVED***

type EUI48 struct ***REMOVED***
	Hdr     RR_Header
	Address uint64 `dns:"uint48"`
***REMOVED***

func (rr *EUI48) String() string ***REMOVED*** return rr.Hdr.String() + euiToString(rr.Address, 48) ***REMOVED***

type EUI64 struct ***REMOVED***
	Hdr     RR_Header
	Address uint64
***REMOVED***

func (rr *EUI64) String() string ***REMOVED*** return rr.Hdr.String() + euiToString(rr.Address, 64) ***REMOVED***

type CAA struct ***REMOVED***
	Hdr   RR_Header
	Flag  uint8
	Tag   string
	Value string `dns:"octet"`
***REMOVED***

func (rr *CAA) String() string ***REMOVED***
	return rr.Hdr.String() + strconv.Itoa(int(rr.Flag)) + " " + rr.Tag + " " + sprintTxtOctet(rr.Value)
***REMOVED***

type UID struct ***REMOVED***
	Hdr RR_Header
	Uid uint32
***REMOVED***

func (rr *UID) String() string ***REMOVED*** return rr.Hdr.String() + strconv.FormatInt(int64(rr.Uid), 10) ***REMOVED***

type GID struct ***REMOVED***
	Hdr RR_Header
	Gid uint32
***REMOVED***

func (rr *GID) String() string ***REMOVED*** return rr.Hdr.String() + strconv.FormatInt(int64(rr.Gid), 10) ***REMOVED***

type UINFO struct ***REMOVED***
	Hdr   RR_Header
	Uinfo string
***REMOVED***

func (rr *UINFO) String() string ***REMOVED*** return rr.Hdr.String() + sprintTxt([]string***REMOVED***rr.Uinfo***REMOVED***) ***REMOVED***

type EID struct ***REMOVED***
	Hdr      RR_Header
	Endpoint string `dns:"hex"`
***REMOVED***

func (rr *EID) String() string ***REMOVED*** return rr.Hdr.String() + strings.ToUpper(rr.Endpoint) ***REMOVED***

type NIMLOC struct ***REMOVED***
	Hdr     RR_Header
	Locator string `dns:"hex"`
***REMOVED***

func (rr *NIMLOC) String() string ***REMOVED*** return rr.Hdr.String() + strings.ToUpper(rr.Locator) ***REMOVED***

type OPENPGPKEY struct ***REMOVED***
	Hdr       RR_Header
	PublicKey string `dns:"base64"`
***REMOVED***

func (rr *OPENPGPKEY) String() string ***REMOVED*** return rr.Hdr.String() + rr.PublicKey ***REMOVED***

// TimeToString translates the RRSIG's incep. and expir. times to the
// string representation used when printing the record.
// It takes serial arithmetic (RFC 1982) into account.
func TimeToString(t uint32) string ***REMOVED***
	mod := ((int64(t) - time.Now().Unix()) / year68) - 1
	if mod < 0 ***REMOVED***
		mod = 0
	***REMOVED***
	ti := time.Unix(int64(t)-(mod*year68), 0).UTC()
	return ti.Format("20060102150405")
***REMOVED***

// StringToTime translates the RRSIG's incep. and expir. times from
// string values like "20110403154150" to an 32 bit integer.
// It takes serial arithmetic (RFC 1982) into account.
func StringToTime(s string) (uint32, error) ***REMOVED***
	t, e := time.Parse("20060102150405", s)
	if e != nil ***REMOVED***
		return 0, e
	***REMOVED***
	mod := (t.Unix() / year68) - 1
	if mod < 0 ***REMOVED***
		mod = 0
	***REMOVED***
	return uint32(t.Unix() - (mod * year68)), nil
***REMOVED***

// saltToString converts a NSECX salt to uppercase and
// returns "-" when it is empty
func saltToString(s string) string ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return "-"
	***REMOVED***
	return strings.ToUpper(s)
***REMOVED***

func euiToString(eui uint64, bits int) (hex string) ***REMOVED***
	switch bits ***REMOVED***
	case 64:
		hex = fmt.Sprintf("%16.16x", eui)
		hex = hex[0:2] + "-" + hex[2:4] + "-" + hex[4:6] + "-" + hex[6:8] +
			"-" + hex[8:10] + "-" + hex[10:12] + "-" + hex[12:14] + "-" + hex[14:16]
	case 48:
		hex = fmt.Sprintf("%12.12x", eui)
		hex = hex[0:2] + "-" + hex[2:4] + "-" + hex[4:6] + "-" + hex[6:8] +
			"-" + hex[8:10] + "-" + hex[10:12]
	***REMOVED***
	return
***REMOVED***

// copyIP returns a copy of ip.
func copyIP(ip net.IP) net.IP ***REMOVED***
	p := make(net.IP, len(ip))
	copy(p, ip)
	return p
***REMOVED***
