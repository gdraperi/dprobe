package dns

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"math/big"
	"strings"
	"time"
)

// Sign signs a dns.Msg. It fills the signature with the appropriate data.
// The SIG record should have the SignerName, KeyTag, Algorithm, Inception
// and Expiration set.
func (rr *SIG) Sign(k crypto.Signer, m *Msg) ([]byte, error) ***REMOVED***
	if k == nil ***REMOVED***
		return nil, ErrPrivKey
	***REMOVED***
	if rr.KeyTag == 0 || len(rr.SignerName) == 0 || rr.Algorithm == 0 ***REMOVED***
		return nil, ErrKey
	***REMOVED***
	rr.Header().Rrtype = TypeSIG
	rr.Header().Class = ClassANY
	rr.Header().Ttl = 0
	rr.Header().Name = "."
	rr.OrigTtl = 0
	rr.TypeCovered = 0
	rr.Labels = 0

	buf := make([]byte, m.Len()+rr.len())
	mbuf, err := m.PackBuffer(buf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if &buf[0] != &mbuf[0] ***REMOVED***
		return nil, ErrBuf
	***REMOVED***
	off, err := PackRR(rr, buf, len(mbuf), nil, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	buf = buf[:off:cap(buf)]

	hash, ok := AlgorithmToHash[rr.Algorithm]
	if !ok ***REMOVED***
		return nil, ErrAlg
	***REMOVED***

	hasher := hash.New()
	// Write SIG rdata
	hasher.Write(buf[len(mbuf)+1+2+2+4+2:])
	// Write message
	hasher.Write(buf[:len(mbuf)])

	signature, err := sign(k, hasher.Sum(nil), hash, rr.Algorithm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rr.Signature = toBase64(signature)
	sig := string(signature)

	buf = append(buf, sig...)
	if len(buf) > int(^uint16(0)) ***REMOVED***
		return nil, ErrBuf
	***REMOVED***
	// Adjust sig data length
	rdoff := len(mbuf) + 1 + 2 + 2 + 4
	rdlen, _ := unpackUint16(buf, rdoff)
	rdlen += uint16(len(sig))
	buf[rdoff], buf[rdoff+1] = packUint16(rdlen)
	// Adjust additional count
	adc, _ := unpackUint16(buf, 10)
	adc++
	buf[10], buf[11] = packUint16(adc)
	return buf, nil
***REMOVED***

// Verify validates the message buf using the key k.
// It's assumed that buf is a valid message from which rr was unpacked.
func (rr *SIG) Verify(k *KEY, buf []byte) error ***REMOVED***
	if k == nil ***REMOVED***
		return ErrKey
	***REMOVED***
	if rr.KeyTag == 0 || len(rr.SignerName) == 0 || rr.Algorithm == 0 ***REMOVED***
		return ErrKey
	***REMOVED***

	var hash crypto.Hash
	switch rr.Algorithm ***REMOVED***
	case DSA, RSASHA1:
		hash = crypto.SHA1
	case RSASHA256, ECDSAP256SHA256:
		hash = crypto.SHA256
	case ECDSAP384SHA384:
		hash = crypto.SHA384
	case RSASHA512:
		hash = crypto.SHA512
	default:
		return ErrAlg
	***REMOVED***
	hasher := hash.New()

	buflen := len(buf)
	qdc, _ := unpackUint16(buf, 4)
	anc, _ := unpackUint16(buf, 6)
	auc, _ := unpackUint16(buf, 8)
	adc, offset := unpackUint16(buf, 10)
	var err error
	for i := uint16(0); i < qdc && offset < buflen; i++ ***REMOVED***
		_, offset, err = UnpackDomainName(buf, offset)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// Skip past Type and Class
		offset += 2 + 2
	***REMOVED***
	for i := uint16(1); i < anc+auc+adc && offset < buflen; i++ ***REMOVED***
		_, offset, err = UnpackDomainName(buf, offset)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// Skip past Type, Class and TTL
		offset += 2 + 2 + 4
		if offset+1 >= buflen ***REMOVED***
			continue
		***REMOVED***
		var rdlen uint16
		rdlen, offset = unpackUint16(buf, offset)
		offset += int(rdlen)
	***REMOVED***
	if offset >= buflen ***REMOVED***
		return &Error***REMOVED***err: "overflowing unpacking signed message"***REMOVED***
	***REMOVED***

	// offset should be just prior to SIG
	bodyend := offset
	// owner name SHOULD be root
	_, offset, err = UnpackDomainName(buf, offset)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Skip Type, Class, TTL, RDLen
	offset += 2 + 2 + 4 + 2
	sigstart := offset
	// Skip Type Covered, Algorithm, Labels, Original TTL
	offset += 2 + 1 + 1 + 4
	if offset+4+4 >= buflen ***REMOVED***
		return &Error***REMOVED***err: "overflow unpacking signed message"***REMOVED***
	***REMOVED***
	expire := uint32(buf[offset])<<24 | uint32(buf[offset+1])<<16 | uint32(buf[offset+2])<<8 | uint32(buf[offset+3])
	offset += 4
	incept := uint32(buf[offset])<<24 | uint32(buf[offset+1])<<16 | uint32(buf[offset+2])<<8 | uint32(buf[offset+3])
	offset += 4
	now := uint32(time.Now().Unix())
	if now < incept || now > expire ***REMOVED***
		return ErrTime
	***REMOVED***
	// Skip key tag
	offset += 2
	var signername string
	signername, offset, err = UnpackDomainName(buf, offset)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// If key has come from the DNS name compression might
	// have mangled the case of the name
	if strings.ToLower(signername) != strings.ToLower(k.Header().Name) ***REMOVED***
		return &Error***REMOVED***err: "signer name doesn't match key name"***REMOVED***
	***REMOVED***
	sigend := offset
	hasher.Write(buf[sigstart:sigend])
	hasher.Write(buf[:10])
	hasher.Write([]byte***REMOVED***
		byte((adc - 1) << 8),
		byte(adc - 1),
	***REMOVED***)
	hasher.Write(buf[12:bodyend])

	hashed := hasher.Sum(nil)
	sig := buf[sigend:]
	switch k.Algorithm ***REMOVED***
	case DSA:
		pk := k.publicKeyDSA()
		sig = sig[1:]
		r := big.NewInt(0)
		r.SetBytes(sig[:len(sig)/2])
		s := big.NewInt(0)
		s.SetBytes(sig[len(sig)/2:])
		if pk != nil ***REMOVED***
			if dsa.Verify(pk, hashed, r, s) ***REMOVED***
				return nil
			***REMOVED***
			return ErrSig
		***REMOVED***
	case RSASHA1, RSASHA256, RSASHA512:
		pk := k.publicKeyRSA()
		if pk != nil ***REMOVED***
			return rsa.VerifyPKCS1v15(pk, hash, hashed, sig)
		***REMOVED***
	case ECDSAP256SHA256, ECDSAP384SHA384:
		pk := k.publicKeyECDSA()
		r := big.NewInt(0)
		r.SetBytes(sig[:len(sig)/2])
		s := big.NewInt(0)
		s.SetBytes(sig[len(sig)/2:])
		if pk != nil ***REMOVED***
			if ecdsa.Verify(pk, hashed, r, s) ***REMOVED***
				return nil
			***REMOVED***
			return ErrSig
		***REMOVED***
	***REMOVED***
	return ErrKeyAlg
***REMOVED***
