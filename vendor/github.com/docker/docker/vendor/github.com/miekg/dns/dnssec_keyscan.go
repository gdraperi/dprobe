package dns

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"io"
	"math/big"
	"strconv"
	"strings"
)

// NewPrivateKey returns a PrivateKey by parsing the string s.
// s should be in the same form of the BIND private key files.
func (k *DNSKEY) NewPrivateKey(s string) (crypto.PrivateKey, error) ***REMOVED***
	if s[len(s)-1] != '\n' ***REMOVED*** // We need a closing newline
		return k.ReadPrivateKey(strings.NewReader(s+"\n"), "")
	***REMOVED***
	return k.ReadPrivateKey(strings.NewReader(s), "")
***REMOVED***

// ReadPrivateKey reads a private key from the io.Reader q. The string file is
// only used in error reporting.
// The public key must be known, because some cryptographic algorithms embed
// the public inside the privatekey.
func (k *DNSKEY) ReadPrivateKey(q io.Reader, file string) (crypto.PrivateKey, error) ***REMOVED***
	m, e := parseKey(q, file)
	if m == nil ***REMOVED***
		return nil, e
	***REMOVED***
	if _, ok := m["private-key-format"]; !ok ***REMOVED***
		return nil, ErrPrivKey
	***REMOVED***
	if m["private-key-format"] != "v1.2" && m["private-key-format"] != "v1.3" ***REMOVED***
		return nil, ErrPrivKey
	***REMOVED***
	// TODO(mg): check if the pubkey matches the private key
	algo, err := strconv.Atoi(strings.SplitN(m["algorithm"], " ", 2)[0])
	if err != nil ***REMOVED***
		return nil, ErrPrivKey
	***REMOVED***
	switch uint8(algo) ***REMOVED***
	case DSA:
		priv, e := readPrivateKeyDSA(m)
		if e != nil ***REMOVED***
			return nil, e
		***REMOVED***
		pub := k.publicKeyDSA()
		if pub == nil ***REMOVED***
			return nil, ErrKey
		***REMOVED***
		priv.PublicKey = *pub
		return priv, e
	case RSAMD5:
		fallthrough
	case RSASHA1:
		fallthrough
	case RSASHA1NSEC3SHA1:
		fallthrough
	case RSASHA256:
		fallthrough
	case RSASHA512:
		priv, e := readPrivateKeyRSA(m)
		if e != nil ***REMOVED***
			return nil, e
		***REMOVED***
		pub := k.publicKeyRSA()
		if pub == nil ***REMOVED***
			return nil, ErrKey
		***REMOVED***
		priv.PublicKey = *pub
		return priv, e
	case ECCGOST:
		return nil, ErrPrivKey
	case ECDSAP256SHA256:
		fallthrough
	case ECDSAP384SHA384:
		priv, e := readPrivateKeyECDSA(m)
		if e != nil ***REMOVED***
			return nil, e
		***REMOVED***
		pub := k.publicKeyECDSA()
		if pub == nil ***REMOVED***
			return nil, ErrKey
		***REMOVED***
		priv.PublicKey = *pub
		return priv, e
	default:
		return nil, ErrPrivKey
	***REMOVED***
***REMOVED***

// Read a private key (file) string and create a public key. Return the private key.
func readPrivateKeyRSA(m map[string]string) (*rsa.PrivateKey, error) ***REMOVED***
	p := new(rsa.PrivateKey)
	p.Primes = []*big.Int***REMOVED***nil, nil***REMOVED***
	for k, v := range m ***REMOVED***
		switch k ***REMOVED***
		case "modulus", "publicexponent", "privateexponent", "prime1", "prime2":
			v1, err := fromBase64([]byte(v))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			switch k ***REMOVED***
			case "modulus":
				p.PublicKey.N = big.NewInt(0)
				p.PublicKey.N.SetBytes(v1)
			case "publicexponent":
				i := big.NewInt(0)
				i.SetBytes(v1)
				p.PublicKey.E = int(i.Int64()) // int64 should be large enough
			case "privateexponent":
				p.D = big.NewInt(0)
				p.D.SetBytes(v1)
			case "prime1":
				p.Primes[0] = big.NewInt(0)
				p.Primes[0].SetBytes(v1)
			case "prime2":
				p.Primes[1] = big.NewInt(0)
				p.Primes[1].SetBytes(v1)
			***REMOVED***
		case "exponent1", "exponent2", "coefficient":
			// not used in Go (yet)
		case "created", "publish", "activate":
			// not used in Go (yet)
		***REMOVED***
	***REMOVED***
	return p, nil
***REMOVED***

func readPrivateKeyDSA(m map[string]string) (*dsa.PrivateKey, error) ***REMOVED***
	p := new(dsa.PrivateKey)
	p.X = big.NewInt(0)
	for k, v := range m ***REMOVED***
		switch k ***REMOVED***
		case "private_value(x)":
			v1, err := fromBase64([]byte(v))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			p.X.SetBytes(v1)
		case "created", "publish", "activate":
			/* not used in Go (yet) */
		***REMOVED***
	***REMOVED***
	return p, nil
***REMOVED***

func readPrivateKeyECDSA(m map[string]string) (*ecdsa.PrivateKey, error) ***REMOVED***
	p := new(ecdsa.PrivateKey)
	p.D = big.NewInt(0)
	// TODO: validate that the required flags are present
	for k, v := range m ***REMOVED***
		switch k ***REMOVED***
		case "privatekey":
			v1, err := fromBase64([]byte(v))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			p.D.SetBytes(v1)
		case "created", "publish", "activate":
			/* not used in Go (yet) */
		***REMOVED***
	***REMOVED***
	return p, nil
***REMOVED***

// parseKey reads a private key from r. It returns a map[string]string,
// with the key-value pairs, or an error when the file is not correct.
func parseKey(r io.Reader, file string) (map[string]string, error) ***REMOVED***
	s := scanInit(r)
	m := make(map[string]string)
	c := make(chan lex)
	k := ""
	// Start the lexer
	go klexer(s, c)
	for l := range c ***REMOVED***
		// It should alternate
		switch l.value ***REMOVED***
		case zKey:
			k = l.token
		case zValue:
			if k == "" ***REMOVED***
				return nil, &ParseError***REMOVED***file, "no private key seen", l***REMOVED***
			***REMOVED***
			//println("Setting", strings.ToLower(k), "to", l.token, "b")
			m[strings.ToLower(k)] = l.token
			k = ""
		***REMOVED***
	***REMOVED***
	return m, nil
***REMOVED***

// klexer scans the sourcefile and returns tokens on the channel c.
func klexer(s *scan, c chan lex) ***REMOVED***
	var l lex
	str := "" // Hold the current read text
	commt := false
	key := true
	x, err := s.tokenText()
	defer close(c)
	for err == nil ***REMOVED***
		l.column = s.position.Column
		l.line = s.position.Line
		switch x ***REMOVED***
		case ':':
			if commt ***REMOVED***
				break
			***REMOVED***
			l.token = str
			if key ***REMOVED***
				l.value = zKey
				c <- l
				// Next token is a space, eat it
				s.tokenText()
				key = false
				str = ""
			***REMOVED*** else ***REMOVED***
				l.value = zValue
			***REMOVED***
		case ';':
			commt = true
		case '\n':
			if commt ***REMOVED***
				// Reset a comment
				commt = false
			***REMOVED***
			l.value = zValue
			l.token = str
			c <- l
			str = ""
			commt = false
			key = true
		default:
			if commt ***REMOVED***
				break
			***REMOVED***
			str += string(x)
		***REMOVED***
		x, err = s.tokenText()
	***REMOVED***
	if len(str) > 0 ***REMOVED***
		// Send remainder
		l.token = str
		l.value = zValue
		c <- l
	***REMOVED***
***REMOVED***
