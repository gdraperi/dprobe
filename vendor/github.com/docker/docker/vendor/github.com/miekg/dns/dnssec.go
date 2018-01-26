package dns

import (
	"bytes"
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	_ "crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/asn1"
	"encoding/hex"
	"math/big"
	"sort"
	"strings"
	"time"
)

// DNSSEC encryption algorithm codes.
const (
	_ uint8 = iota
	RSAMD5
	DH
	DSA
	_ // Skip 4, RFC 6725, section 2.1
	RSASHA1
	DSANSEC3SHA1
	RSASHA1NSEC3SHA1
	RSASHA256
	_ // Skip 9, RFC 6725, section 2.1
	RSASHA512
	_ // Skip 11, RFC 6725, section 2.1
	ECCGOST
	ECDSAP256SHA256
	ECDSAP384SHA384
	INDIRECT   uint8 = 252
	PRIVATEDNS uint8 = 253 // Private (experimental keys)
	PRIVATEOID uint8 = 254
)

// Map for algorithm names.
var AlgorithmToString = map[uint8]string***REMOVED***
	RSAMD5:           "RSAMD5",
	DH:               "DH",
	DSA:              "DSA",
	RSASHA1:          "RSASHA1",
	DSANSEC3SHA1:     "DSA-NSEC3-SHA1",
	RSASHA1NSEC3SHA1: "RSASHA1-NSEC3-SHA1",
	RSASHA256:        "RSASHA256",
	RSASHA512:        "RSASHA512",
	ECCGOST:          "ECC-GOST",
	ECDSAP256SHA256:  "ECDSAP256SHA256",
	ECDSAP384SHA384:  "ECDSAP384SHA384",
	INDIRECT:         "INDIRECT",
	PRIVATEDNS:       "PRIVATEDNS",
	PRIVATEOID:       "PRIVATEOID",
***REMOVED***

// Map of algorithm strings.
var StringToAlgorithm = reverseInt8(AlgorithmToString)

// Map of algorithm crypto hashes.
var AlgorithmToHash = map[uint8]crypto.Hash***REMOVED***
	RSAMD5:           crypto.MD5, // Deprecated in RFC 6725
	RSASHA1:          crypto.SHA1,
	RSASHA1NSEC3SHA1: crypto.SHA1,
	RSASHA256:        crypto.SHA256,
	ECDSAP256SHA256:  crypto.SHA256,
	ECDSAP384SHA384:  crypto.SHA384,
	RSASHA512:        crypto.SHA512,
***REMOVED***

// DNSSEC hashing algorithm codes.
const (
	_      uint8 = iota
	SHA1         // RFC 4034
	SHA256       // RFC 4509
	GOST94       // RFC 5933
	SHA384       // Experimental
	SHA512       // Experimental
)

// Map for hash names.
var HashToString = map[uint8]string***REMOVED***
	SHA1:   "SHA1",
	SHA256: "SHA256",
	GOST94: "GOST94",
	SHA384: "SHA384",
	SHA512: "SHA512",
***REMOVED***

// Map of hash strings.
var StringToHash = reverseInt8(HashToString)

// DNSKEY flag values.
const (
	SEP    = 1
	REVOKE = 1 << 7
	ZONE   = 1 << 8
)

// The RRSIG needs to be converted to wireformat with some of
// the rdata (the signature) missing. Use this struct to ease
// the conversion (and re-use the pack/unpack functions).
type rrsigWireFmt struct ***REMOVED***
	TypeCovered uint16
	Algorithm   uint8
	Labels      uint8
	OrigTtl     uint32
	Expiration  uint32
	Inception   uint32
	KeyTag      uint16
	SignerName  string `dns:"domain-name"`
	/* No Signature */
***REMOVED***

// Used for converting DNSKEY's rdata to wirefmt.
type dnskeyWireFmt struct ***REMOVED***
	Flags     uint16
	Protocol  uint8
	Algorithm uint8
	PublicKey string `dns:"base64"`
	/* Nothing is left out */
***REMOVED***

func divRoundUp(a, b int) int ***REMOVED***
	return (a + b - 1) / b
***REMOVED***

// KeyTag calculates the keytag (or key-id) of the DNSKEY.
func (k *DNSKEY) KeyTag() uint16 ***REMOVED***
	if k == nil ***REMOVED***
		return 0
	***REMOVED***
	var keytag int
	switch k.Algorithm ***REMOVED***
	case RSAMD5:
		// Look at the bottom two bytes of the modules, which the last
		// item in the pubkey. We could do this faster by looking directly
		// at the base64 values. But I'm lazy.
		modulus, _ := fromBase64([]byte(k.PublicKey))
		if len(modulus) > 1 ***REMOVED***
			x, _ := unpackUint16(modulus, len(modulus)-2)
			keytag = int(x)
		***REMOVED***
	default:
		keywire := new(dnskeyWireFmt)
		keywire.Flags = k.Flags
		keywire.Protocol = k.Protocol
		keywire.Algorithm = k.Algorithm
		keywire.PublicKey = k.PublicKey
		wire := make([]byte, DefaultMsgSize)
		n, err := PackStruct(keywire, wire, 0)
		if err != nil ***REMOVED***
			return 0
		***REMOVED***
		wire = wire[:n]
		for i, v := range wire ***REMOVED***
			if i&1 != 0 ***REMOVED***
				keytag += int(v) // must be larger than uint32
			***REMOVED*** else ***REMOVED***
				keytag += int(v) << 8
			***REMOVED***
		***REMOVED***
		keytag += (keytag >> 16) & 0xFFFF
		keytag &= 0xFFFF
	***REMOVED***
	return uint16(keytag)
***REMOVED***

// ToDS converts a DNSKEY record to a DS record.
func (k *DNSKEY) ToDS(h uint8) *DS ***REMOVED***
	if k == nil ***REMOVED***
		return nil
	***REMOVED***
	ds := new(DS)
	ds.Hdr.Name = k.Hdr.Name
	ds.Hdr.Class = k.Hdr.Class
	ds.Hdr.Rrtype = TypeDS
	ds.Hdr.Ttl = k.Hdr.Ttl
	ds.Algorithm = k.Algorithm
	ds.DigestType = h
	ds.KeyTag = k.KeyTag()

	keywire := new(dnskeyWireFmt)
	keywire.Flags = k.Flags
	keywire.Protocol = k.Protocol
	keywire.Algorithm = k.Algorithm
	keywire.PublicKey = k.PublicKey
	wire := make([]byte, DefaultMsgSize)
	n, err := PackStruct(keywire, wire, 0)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	wire = wire[:n]

	owner := make([]byte, 255)
	off, err1 := PackDomainName(strings.ToLower(k.Hdr.Name), owner, 0, nil, false)
	if err1 != nil ***REMOVED***
		return nil
	***REMOVED***
	owner = owner[:off]
	// RFC4034:
	// digest = digest_algorithm( DNSKEY owner name | DNSKEY RDATA);
	// "|" denotes concatenation
	// DNSKEY RDATA = Flags | Protocol | Algorithm | Public Key.

	// digest buffer
	digest := append(owner, wire...) // another copy

	var hash crypto.Hash
	switch h ***REMOVED***
	case SHA1:
		hash = crypto.SHA1
	case SHA256:
		hash = crypto.SHA256
	case SHA384:
		hash = crypto.SHA384
	case SHA512:
		hash = crypto.SHA512
	default:
		return nil
	***REMOVED***

	s := hash.New()
	s.Write(digest)
	ds.Digest = hex.EncodeToString(s.Sum(nil))
	return ds
***REMOVED***

// ToCDNSKEY converts a DNSKEY record to a CDNSKEY record.
func (k *DNSKEY) ToCDNSKEY() *CDNSKEY ***REMOVED***
	c := &CDNSKEY***REMOVED***DNSKEY: *k***REMOVED***
	c.Hdr = *k.Hdr.copyHeader()
	c.Hdr.Rrtype = TypeCDNSKEY
	return c
***REMOVED***

// ToCDS converts a DS record to a CDS record.
func (d *DS) ToCDS() *CDS ***REMOVED***
	c := &CDS***REMOVED***DS: *d***REMOVED***
	c.Hdr = *d.Hdr.copyHeader()
	c.Hdr.Rrtype = TypeCDS
	return c
***REMOVED***

// Sign signs an RRSet. The signature needs to be filled in with the values:
// Inception, Expiration, KeyTag, SignerName and Algorithm.  The rest is copied
// from the RRset. Sign returns a non-nill error when the signing went OK.
// There is no check if RRSet is a proper (RFC 2181) RRSet.  If OrigTTL is non
// zero, it is used as-is, otherwise the TTL of the RRset is used as the
// OrigTTL.
func (rr *RRSIG) Sign(k crypto.Signer, rrset []RR) error ***REMOVED***
	if k == nil ***REMOVED***
		return ErrPrivKey
	***REMOVED***
	// s.Inception and s.Expiration may be 0 (rollover etc.), the rest must be set
	if rr.KeyTag == 0 || len(rr.SignerName) == 0 || rr.Algorithm == 0 ***REMOVED***
		return ErrKey
	***REMOVED***

	rr.Hdr.Rrtype = TypeRRSIG
	rr.Hdr.Name = rrset[0].Header().Name
	rr.Hdr.Class = rrset[0].Header().Class
	if rr.OrigTtl == 0 ***REMOVED*** // If set don't override
		rr.OrigTtl = rrset[0].Header().Ttl
	***REMOVED***
	rr.TypeCovered = rrset[0].Header().Rrtype
	rr.Labels = uint8(CountLabel(rrset[0].Header().Name))

	if strings.HasPrefix(rrset[0].Header().Name, "*") ***REMOVED***
		rr.Labels-- // wildcard, remove from label count
	***REMOVED***

	sigwire := new(rrsigWireFmt)
	sigwire.TypeCovered = rr.TypeCovered
	sigwire.Algorithm = rr.Algorithm
	sigwire.Labels = rr.Labels
	sigwire.OrigTtl = rr.OrigTtl
	sigwire.Expiration = rr.Expiration
	sigwire.Inception = rr.Inception
	sigwire.KeyTag = rr.KeyTag
	// For signing, lowercase this name
	sigwire.SignerName = strings.ToLower(rr.SignerName)

	// Create the desired binary blob
	signdata := make([]byte, DefaultMsgSize)
	n, err := PackStruct(sigwire, signdata, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	signdata = signdata[:n]
	wire, err := rawSignatureData(rrset, rr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	signdata = append(signdata, wire...)

	hash, ok := AlgorithmToHash[rr.Algorithm]
	if !ok ***REMOVED***
		return ErrAlg
	***REMOVED***

	h := hash.New()
	h.Write(signdata)

	signature, err := sign(k, h.Sum(nil), hash, rr.Algorithm)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	rr.Signature = toBase64(signature)

	return nil
***REMOVED***

func sign(k crypto.Signer, hashed []byte, hash crypto.Hash, alg uint8) ([]byte, error) ***REMOVED***
	signature, err := k.Sign(rand.Reader, hashed, hash)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch alg ***REMOVED***
	case RSASHA1, RSASHA1NSEC3SHA1, RSASHA256, RSASHA512:
		return signature, nil

	case ECDSAP256SHA256, ECDSAP384SHA384:
		ecdsaSignature := &struct ***REMOVED***
			R, S *big.Int
		***REMOVED******REMOVED******REMOVED***
		if _, err := asn1.Unmarshal(signature, ecdsaSignature); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		var intlen int
		switch alg ***REMOVED***
		case ECDSAP256SHA256:
			intlen = 32
		case ECDSAP384SHA384:
			intlen = 48
		***REMOVED***

		signature := intToBytes(ecdsaSignature.R, intlen)
		signature = append(signature, intToBytes(ecdsaSignature.S, intlen)...)
		return signature, nil

	// There is no defined interface for what a DSA backed crypto.Signer returns
	case DSA, DSANSEC3SHA1:
		// 	t := divRoundUp(divRoundUp(p.PublicKey.Y.BitLen(), 8)-64, 8)
		// 	signature := []byte***REMOVED***byte(t)***REMOVED***
		// 	signature = append(signature, intToBytes(r1, 20)...)
		// 	signature = append(signature, intToBytes(s1, 20)...)
		// 	rr.Signature = signature
	***REMOVED***

	return nil, ErrAlg
***REMOVED***

// Verify validates an RRSet with the signature and key. This is only the
// cryptographic test, the signature validity period must be checked separately.
// This function copies the rdata of some RRs (to lowercase domain names) for the validation to work.
func (rr *RRSIG) Verify(k *DNSKEY, rrset []RR) error ***REMOVED***
	// First the easy checks
	if !IsRRset(rrset) ***REMOVED***
		return ErrRRset
	***REMOVED***
	if rr.KeyTag != k.KeyTag() ***REMOVED***
		return ErrKey
	***REMOVED***
	if rr.Hdr.Class != k.Hdr.Class ***REMOVED***
		return ErrKey
	***REMOVED***
	if rr.Algorithm != k.Algorithm ***REMOVED***
		return ErrKey
	***REMOVED***
	if strings.ToLower(rr.SignerName) != strings.ToLower(k.Hdr.Name) ***REMOVED***
		return ErrKey
	***REMOVED***
	if k.Protocol != 3 ***REMOVED***
		return ErrKey
	***REMOVED***

	// IsRRset checked that we have at least one RR and that the RRs in
	// the set have consistent type, class, and name. Also check that type and
	// class matches the RRSIG record.
	if rrset[0].Header().Class != rr.Hdr.Class ***REMOVED***
		return ErrRRset
	***REMOVED***
	if rrset[0].Header().Rrtype != rr.TypeCovered ***REMOVED***
		return ErrRRset
	***REMOVED***

	// RFC 4035 5.3.2.  Reconstructing the Signed Data
	// Copy the sig, except the rrsig data
	sigwire := new(rrsigWireFmt)
	sigwire.TypeCovered = rr.TypeCovered
	sigwire.Algorithm = rr.Algorithm
	sigwire.Labels = rr.Labels
	sigwire.OrigTtl = rr.OrigTtl
	sigwire.Expiration = rr.Expiration
	sigwire.Inception = rr.Inception
	sigwire.KeyTag = rr.KeyTag
	sigwire.SignerName = strings.ToLower(rr.SignerName)
	// Create the desired binary blob
	signeddata := make([]byte, DefaultMsgSize)
	n, err := PackStruct(sigwire, signeddata, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	signeddata = signeddata[:n]
	wire, err := rawSignatureData(rrset, rr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	signeddata = append(signeddata, wire...)

	sigbuf := rr.sigBuf()           // Get the binary signature data
	if rr.Algorithm == PRIVATEDNS ***REMOVED*** // PRIVATEOID
		// TODO(miek)
		// remove the domain name and assume its ours?
	***REMOVED***

	hash, ok := AlgorithmToHash[rr.Algorithm]
	if !ok ***REMOVED***
		return ErrAlg
	***REMOVED***

	switch rr.Algorithm ***REMOVED***
	case RSASHA1, RSASHA1NSEC3SHA1, RSASHA256, RSASHA512, RSAMD5:
		// TODO(mg): this can be done quicker, ie. cache the pubkey data somewhere??
		pubkey := k.publicKeyRSA() // Get the key
		if pubkey == nil ***REMOVED***
			return ErrKey
		***REMOVED***

		h := hash.New()
		h.Write(signeddata)
		return rsa.VerifyPKCS1v15(pubkey, hash, h.Sum(nil), sigbuf)

	case ECDSAP256SHA256, ECDSAP384SHA384:
		pubkey := k.publicKeyECDSA()
		if pubkey == nil ***REMOVED***
			return ErrKey
		***REMOVED***

		// Split sigbuf into the r and s coordinates
		r := new(big.Int).SetBytes(sigbuf[:len(sigbuf)/2])
		s := new(big.Int).SetBytes(sigbuf[len(sigbuf)/2:])

		h := hash.New()
		h.Write(signeddata)
		if ecdsa.Verify(pubkey, h.Sum(nil), r, s) ***REMOVED***
			return nil
		***REMOVED***
		return ErrSig

	default:
		return ErrAlg
	***REMOVED***
***REMOVED***

// ValidityPeriod uses RFC1982 serial arithmetic to calculate
// if a signature period is valid. If t is the zero time, the
// current time is taken other t is. Returns true if the signature
// is valid at the given time, otherwise returns false.
func (rr *RRSIG) ValidityPeriod(t time.Time) bool ***REMOVED***
	var utc int64
	if t.IsZero() ***REMOVED***
		utc = time.Now().UTC().Unix()
	***REMOVED*** else ***REMOVED***
		utc = t.UTC().Unix()
	***REMOVED***
	modi := (int64(rr.Inception) - utc) / year68
	mode := (int64(rr.Expiration) - utc) / year68
	ti := int64(rr.Inception) + (modi * year68)
	te := int64(rr.Expiration) + (mode * year68)
	return ti <= utc && utc <= te
***REMOVED***

// Return the signatures base64 encodedig sigdata as a byte slice.
func (rr *RRSIG) sigBuf() []byte ***REMOVED***
	sigbuf, err := fromBase64([]byte(rr.Signature))
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return sigbuf
***REMOVED***

// publicKeyRSA returns the RSA public key from a DNSKEY record.
func (k *DNSKEY) publicKeyRSA() *rsa.PublicKey ***REMOVED***
	keybuf, err := fromBase64([]byte(k.PublicKey))
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	// RFC 2537/3110, section 2. RSA Public KEY Resource Records
	// Length is in the 0th byte, unless its zero, then it
	// it in bytes 1 and 2 and its a 16 bit number
	explen := uint16(keybuf[0])
	keyoff := 1
	if explen == 0 ***REMOVED***
		explen = uint16(keybuf[1])<<8 | uint16(keybuf[2])
		keyoff = 3
	***REMOVED***
	pubkey := new(rsa.PublicKey)

	pubkey.N = big.NewInt(0)
	shift := uint64((explen - 1) * 8)
	expo := uint64(0)
	for i := int(explen - 1); i > 0; i-- ***REMOVED***
		expo += uint64(keybuf[keyoff+i]) << shift
		shift -= 8
	***REMOVED***
	// Remainder
	expo += uint64(keybuf[keyoff])
	if expo > 2<<31 ***REMOVED***
		// Larger expo than supported.
		// println("dns: F5 primes (or larger) are not supported")
		return nil
	***REMOVED***
	pubkey.E = int(expo)

	pubkey.N.SetBytes(keybuf[keyoff+int(explen):])
	return pubkey
***REMOVED***

// publicKeyECDSA returns the Curve public key from the DNSKEY record.
func (k *DNSKEY) publicKeyECDSA() *ecdsa.PublicKey ***REMOVED***
	keybuf, err := fromBase64([]byte(k.PublicKey))
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	pubkey := new(ecdsa.PublicKey)
	switch k.Algorithm ***REMOVED***
	case ECDSAP256SHA256:
		pubkey.Curve = elliptic.P256()
		if len(keybuf) != 64 ***REMOVED***
			// wrongly encoded key
			return nil
		***REMOVED***
	case ECDSAP384SHA384:
		pubkey.Curve = elliptic.P384()
		if len(keybuf) != 96 ***REMOVED***
			// Wrongly encoded key
			return nil
		***REMOVED***
	***REMOVED***
	pubkey.X = big.NewInt(0)
	pubkey.X.SetBytes(keybuf[:len(keybuf)/2])
	pubkey.Y = big.NewInt(0)
	pubkey.Y.SetBytes(keybuf[len(keybuf)/2:])
	return pubkey
***REMOVED***

func (k *DNSKEY) publicKeyDSA() *dsa.PublicKey ***REMOVED***
	keybuf, err := fromBase64([]byte(k.PublicKey))
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	if len(keybuf) < 22 ***REMOVED***
		return nil
	***REMOVED***
	t, keybuf := int(keybuf[0]), keybuf[1:]
	size := 64 + t*8
	q, keybuf := keybuf[:20], keybuf[20:]
	if len(keybuf) != 3*size ***REMOVED***
		return nil
	***REMOVED***
	p, keybuf := keybuf[:size], keybuf[size:]
	g, y := keybuf[:size], keybuf[size:]
	pubkey := new(dsa.PublicKey)
	pubkey.Parameters.Q = big.NewInt(0).SetBytes(q)
	pubkey.Parameters.P = big.NewInt(0).SetBytes(p)
	pubkey.Parameters.G = big.NewInt(0).SetBytes(g)
	pubkey.Y = big.NewInt(0).SetBytes(y)
	return pubkey
***REMOVED***

type wireSlice [][]byte

func (p wireSlice) Len() int      ***REMOVED*** return len(p) ***REMOVED***
func (p wireSlice) Swap(i, j int) ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***
func (p wireSlice) Less(i, j int) bool ***REMOVED***
	_, ioff, _ := UnpackDomainName(p[i], 0)
	_, joff, _ := UnpackDomainName(p[j], 0)
	return bytes.Compare(p[i][ioff+10:], p[j][joff+10:]) < 0
***REMOVED***

// Return the raw signature data.
func rawSignatureData(rrset []RR, s *RRSIG) (buf []byte, err error) ***REMOVED***
	wires := make(wireSlice, len(rrset))
	for i, r := range rrset ***REMOVED***
		r1 := r.copy()
		r1.Header().Ttl = s.OrigTtl
		labels := SplitDomainName(r1.Header().Name)
		// 6.2. Canonical RR Form. (4) - wildcards
		if len(labels) > int(s.Labels) ***REMOVED***
			// Wildcard
			r1.Header().Name = "*." + strings.Join(labels[len(labels)-int(s.Labels):], ".") + "."
		***REMOVED***
		// RFC 4034: 6.2.  Canonical RR Form. (2) - domain name to lowercase
		r1.Header().Name = strings.ToLower(r1.Header().Name)
		// 6.2. Canonical RR Form. (3) - domain rdata to lowercase.
		//   NS, MD, MF, CNAME, SOA, MB, MG, MR, PTR,
		//   HINFO, MINFO, MX, RP, AFSDB, RT, SIG, PX, NXT, NAPTR, KX,
		//   SRV, DNAME, A6
		//
		// RFC 6840 - Clarifications and Implementation Notes for DNS Security (DNSSEC):
		//	Section 6.2 of [RFC4034] also erroneously lists HINFO as a record
		//	that needs conversion to lowercase, and twice at that.  Since HINFO
		//	records contain no domain names, they are not subject to case
		//	conversion.
		switch x := r1.(type) ***REMOVED***
		case *NS:
			x.Ns = strings.ToLower(x.Ns)
		case *CNAME:
			x.Target = strings.ToLower(x.Target)
		case *SOA:
			x.Ns = strings.ToLower(x.Ns)
			x.Mbox = strings.ToLower(x.Mbox)
		case *MB:
			x.Mb = strings.ToLower(x.Mb)
		case *MG:
			x.Mg = strings.ToLower(x.Mg)
		case *MR:
			x.Mr = strings.ToLower(x.Mr)
		case *PTR:
			x.Ptr = strings.ToLower(x.Ptr)
		case *MINFO:
			x.Rmail = strings.ToLower(x.Rmail)
			x.Email = strings.ToLower(x.Email)
		case *MX:
			x.Mx = strings.ToLower(x.Mx)
		case *NAPTR:
			x.Replacement = strings.ToLower(x.Replacement)
		case *KX:
			x.Exchanger = strings.ToLower(x.Exchanger)
		case *SRV:
			x.Target = strings.ToLower(x.Target)
		case *DNAME:
			x.Target = strings.ToLower(x.Target)
		***REMOVED***
		// 6.2. Canonical RR Form. (5) - origTTL
		wire := make([]byte, r1.len()+1) // +1 to be safe(r)
		off, err1 := PackRR(r1, wire, 0, nil, false)
		if err1 != nil ***REMOVED***
			return nil, err1
		***REMOVED***
		wire = wire[:off]
		wires[i] = wire
	***REMOVED***
	sort.Sort(wires)
	for i, wire := range wires ***REMOVED***
		if i > 0 && bytes.Equal(wire, wires[i-1]) ***REMOVED***
			continue
		***REMOVED***
		buf = append(buf, wire...)
	***REMOVED***
	return buf, nil
***REMOVED***
