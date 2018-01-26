package libtrust

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
	"unicode"
)

var (
	// ErrInvalidSignContent is used when the content to be signed is invalid.
	ErrInvalidSignContent = errors.New("invalid sign content")

	// ErrInvalidJSONContent is used when invalid json is encountered.
	ErrInvalidJSONContent = errors.New("invalid json content")

	// ErrMissingSignatureKey is used when the specified signature key
	// does not exist in the JSON content.
	ErrMissingSignatureKey = errors.New("missing signature key")
)

type jsHeader struct ***REMOVED***
	JWK       PublicKey `json:"jwk,omitempty"`
	Algorithm string    `json:"alg"`
	Chain     []string  `json:"x5c,omitempty"`
***REMOVED***

type jsSignature struct ***REMOVED***
	Header    jsHeader `json:"header"`
	Signature string   `json:"signature"`
	Protected string   `json:"protected,omitempty"`
***REMOVED***

type jsSignaturesSorted []jsSignature

func (jsbkid jsSignaturesSorted) Swap(i, j int) ***REMOVED*** jsbkid[i], jsbkid[j] = jsbkid[j], jsbkid[i] ***REMOVED***
func (jsbkid jsSignaturesSorted) Len() int      ***REMOVED*** return len(jsbkid) ***REMOVED***

func (jsbkid jsSignaturesSorted) Less(i, j int) bool ***REMOVED***
	ki, kj := jsbkid[i].Header.JWK.KeyID(), jsbkid[j].Header.JWK.KeyID()
	si, sj := jsbkid[i].Signature, jsbkid[j].Signature

	if ki == kj ***REMOVED***
		return si < sj
	***REMOVED***

	return ki < kj
***REMOVED***

type signKey struct ***REMOVED***
	PrivateKey
	Chain []*x509.Certificate
***REMOVED***

// JSONSignature represents a signature of a json object.
type JSONSignature struct ***REMOVED***
	payload      string
	signatures   []jsSignature
	indent       string
	formatLength int
	formatTail   []byte
***REMOVED***

func newJSONSignature() *JSONSignature ***REMOVED***
	return &JSONSignature***REMOVED***
		signatures: make([]jsSignature, 0, 1),
	***REMOVED***
***REMOVED***

// Payload returns the encoded payload of the signature. This
// payload should not be signed directly
func (js *JSONSignature) Payload() ([]byte, error) ***REMOVED***
	return joseBase64UrlDecode(js.payload)
***REMOVED***

func (js *JSONSignature) protectedHeader() (string, error) ***REMOVED***
	protected := map[string]interface***REMOVED******REMOVED******REMOVED***
		"formatLength": js.formatLength,
		"formatTail":   joseBase64UrlEncode(js.formatTail),
		"time":         time.Now().UTC().Format(time.RFC3339),
	***REMOVED***
	protectedBytes, err := json.Marshal(protected)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return joseBase64UrlEncode(protectedBytes), nil
***REMOVED***

func (js *JSONSignature) signBytes(protectedHeader string) ([]byte, error) ***REMOVED***
	buf := make([]byte, len(js.payload)+len(protectedHeader)+1)
	copy(buf, protectedHeader)
	buf[len(protectedHeader)] = '.'
	copy(buf[len(protectedHeader)+1:], js.payload)
	return buf, nil
***REMOVED***

// Sign adds a signature using the given private key.
func (js *JSONSignature) Sign(key PrivateKey) error ***REMOVED***
	protected, err := js.protectedHeader()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	signBytes, err := js.signBytes(protected)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	sigBytes, algorithm, err := key.Sign(bytes.NewReader(signBytes), crypto.SHA256)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	js.signatures = append(js.signatures, jsSignature***REMOVED***
		Header: jsHeader***REMOVED***
			JWK:       key.PublicKey(),
			Algorithm: algorithm,
		***REMOVED***,
		Signature: joseBase64UrlEncode(sigBytes),
		Protected: protected,
	***REMOVED***)

	return nil
***REMOVED***

// SignWithChain adds a signature using the given private key
// and setting the x509 chain. The public key of the first element
// in the chain must be the public key corresponding with the sign key.
func (js *JSONSignature) SignWithChain(key PrivateKey, chain []*x509.Certificate) error ***REMOVED***
	// Ensure key.Chain[0] is public key for key
	//key.Chain.PublicKey
	//key.PublicKey().CryptoPublicKey()

	// Verify chain
	protected, err := js.protectedHeader()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	signBytes, err := js.signBytes(protected)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	sigBytes, algorithm, err := key.Sign(bytes.NewReader(signBytes), crypto.SHA256)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	header := jsHeader***REMOVED***
		Chain:     make([]string, len(chain)),
		Algorithm: algorithm,
	***REMOVED***

	for i, cert := range chain ***REMOVED***
		header.Chain[i] = base64.StdEncoding.EncodeToString(cert.Raw)
	***REMOVED***

	js.signatures = append(js.signatures, jsSignature***REMOVED***
		Header:    header,
		Signature: joseBase64UrlEncode(sigBytes),
		Protected: protected,
	***REMOVED***)

	return nil
***REMOVED***

// Verify verifies all the signatures and returns the list of
// public keys used to sign. Any x509 chains are not checked.
func (js *JSONSignature) Verify() ([]PublicKey, error) ***REMOVED***
	keys := make([]PublicKey, len(js.signatures))
	for i, signature := range js.signatures ***REMOVED***
		signBytes, err := js.signBytes(signature.Protected)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var publicKey PublicKey
		if len(signature.Header.Chain) > 0 ***REMOVED***
			certBytes, err := base64.StdEncoding.DecodeString(signature.Header.Chain[0])
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			cert, err := x509.ParseCertificate(certBytes)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			publicKey, err = FromCryptoPublicKey(cert.PublicKey)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED*** else if signature.Header.JWK != nil ***REMOVED***
			publicKey = signature.Header.JWK
		***REMOVED*** else ***REMOVED***
			return nil, errors.New("missing public key")
		***REMOVED***

		sigBytes, err := joseBase64UrlDecode(signature.Signature)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		err = publicKey.Verify(bytes.NewReader(signBytes), signature.Header.Algorithm, sigBytes)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		keys[i] = publicKey
	***REMOVED***
	return keys, nil
***REMOVED***

// VerifyChains verifies all the signatures and the chains associated
// with each signature and returns the list of verified chains.
// Signatures without an x509 chain are not checked.
func (js *JSONSignature) VerifyChains(ca *x509.CertPool) ([][]*x509.Certificate, error) ***REMOVED***
	chains := make([][]*x509.Certificate, 0, len(js.signatures))
	for _, signature := range js.signatures ***REMOVED***
		signBytes, err := js.signBytes(signature.Protected)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var publicKey PublicKey
		if len(signature.Header.Chain) > 0 ***REMOVED***
			certBytes, err := base64.StdEncoding.DecodeString(signature.Header.Chain[0])
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			cert, err := x509.ParseCertificate(certBytes)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			publicKey, err = FromCryptoPublicKey(cert.PublicKey)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			intermediates := x509.NewCertPool()
			if len(signature.Header.Chain) > 1 ***REMOVED***
				intermediateChain := signature.Header.Chain[1:]
				for i := range intermediateChain ***REMOVED***
					certBytes, err := base64.StdEncoding.DecodeString(intermediateChain[i])
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					intermediate, err := x509.ParseCertificate(certBytes)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					intermediates.AddCert(intermediate)
				***REMOVED***
			***REMOVED***

			verifyOptions := x509.VerifyOptions***REMOVED***
				Intermediates: intermediates,
				Roots:         ca,
			***REMOVED***

			verifiedChains, err := cert.Verify(verifyOptions)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			chains = append(chains, verifiedChains...)

			sigBytes, err := joseBase64UrlDecode(signature.Signature)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			err = publicKey.Verify(bytes.NewReader(signBytes), signature.Header.Algorithm, sigBytes)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***

	***REMOVED***
	return chains, nil
***REMOVED***

// JWS returns JSON serialized JWS according to
// http://tools.ietf.org/html/draft-ietf-jose-json-web-signature-31#section-7.2
func (js *JSONSignature) JWS() ([]byte, error) ***REMOVED***
	if len(js.signatures) == 0 ***REMOVED***
		return nil, errors.New("missing signature")
	***REMOVED***

	sort.Sort(jsSignaturesSorted(js.signatures))

	jsonMap := map[string]interface***REMOVED******REMOVED******REMOVED***
		"payload":    js.payload,
		"signatures": js.signatures,
	***REMOVED***

	return json.MarshalIndent(jsonMap, "", "   ")
***REMOVED***

func notSpace(r rune) bool ***REMOVED***
	return !unicode.IsSpace(r)
***REMOVED***

func detectJSONIndent(jsonContent []byte) (indent string) ***REMOVED***
	if len(jsonContent) > 2 && jsonContent[0] == '***REMOVED***' && jsonContent[1] == '\n' ***REMOVED***
		quoteIndex := bytes.IndexRune(jsonContent[1:], '"')
		if quoteIndex > 0 ***REMOVED***
			indent = string(jsonContent[2 : quoteIndex+1])
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

type jsParsedHeader struct ***REMOVED***
	JWK       json.RawMessage `json:"jwk"`
	Algorithm string          `json:"alg"`
	Chain     []string        `json:"x5c"`
***REMOVED***

type jsParsedSignature struct ***REMOVED***
	Header    jsParsedHeader `json:"header"`
	Signature string         `json:"signature"`
	Protected string         `json:"protected"`
***REMOVED***

// ParseJWS parses a JWS serialized JSON object into a Json Signature.
func ParseJWS(content []byte) (*JSONSignature, error) ***REMOVED***
	type jsParsed struct ***REMOVED***
		Payload    string              `json:"payload"`
		Signatures []jsParsedSignature `json:"signatures"`
	***REMOVED***
	parsed := &jsParsed***REMOVED******REMOVED***
	err := json.Unmarshal(content, parsed)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(parsed.Signatures) == 0 ***REMOVED***
		return nil, errors.New("missing signatures")
	***REMOVED***
	payload, err := joseBase64UrlDecode(parsed.Payload)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	js, err := NewJSONSignature(payload)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	js.signatures = make([]jsSignature, len(parsed.Signatures))
	for i, signature := range parsed.Signatures ***REMOVED***
		header := jsHeader***REMOVED***
			Algorithm: signature.Header.Algorithm,
		***REMOVED***
		if signature.Header.Chain != nil ***REMOVED***
			header.Chain = signature.Header.Chain
		***REMOVED***
		if signature.Header.JWK != nil ***REMOVED***
			publicKey, err := UnmarshalPublicKeyJWK([]byte(signature.Header.JWK))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			header.JWK = publicKey
		***REMOVED***
		js.signatures[i] = jsSignature***REMOVED***
			Header:    header,
			Signature: signature.Signature,
			Protected: signature.Protected,
		***REMOVED***
	***REMOVED***

	return js, nil
***REMOVED***

// NewJSONSignature returns a new unsigned JWS from a json byte array.
// JSONSignature will need to be signed before serializing or storing.
// Optionally, one or more signatures can be provided as byte buffers,
// containing serialized JWS signatures, to assemble a fully signed JWS
// package. It is the callers responsibility to ensure uniqueness of the
// provided signatures.
func NewJSONSignature(content []byte, signatures ...[]byte) (*JSONSignature, error) ***REMOVED***
	var dataMap map[string]interface***REMOVED******REMOVED***
	err := json.Unmarshal(content, &dataMap)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	js := newJSONSignature()
	js.indent = detectJSONIndent(content)

	js.payload = joseBase64UrlEncode(content)

	// Find trailing ***REMOVED*** and whitespace, put in protected header
	closeIndex := bytes.LastIndexFunc(content, notSpace)
	if content[closeIndex] != '***REMOVED***' ***REMOVED***
		return nil, ErrInvalidJSONContent
	***REMOVED***
	lastRuneIndex := bytes.LastIndexFunc(content[:closeIndex], notSpace)
	if content[lastRuneIndex] == ',' ***REMOVED***
		return nil, ErrInvalidJSONContent
	***REMOVED***
	js.formatLength = lastRuneIndex + 1
	js.formatTail = content[js.formatLength:]

	if len(signatures) > 0 ***REMOVED***
		for _, signature := range signatures ***REMOVED***
			var parsedJSig jsParsedSignature

			if err := json.Unmarshal(signature, &parsedJSig); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			// TODO(stevvooe): A lot of the code below is repeated in
			// ParseJWS. It will require more refactoring to fix that.
			jsig := jsSignature***REMOVED***
				Header: jsHeader***REMOVED***
					Algorithm: parsedJSig.Header.Algorithm,
				***REMOVED***,
				Signature: parsedJSig.Signature,
				Protected: parsedJSig.Protected,
			***REMOVED***

			if parsedJSig.Header.Chain != nil ***REMOVED***
				jsig.Header.Chain = parsedJSig.Header.Chain
			***REMOVED***

			if parsedJSig.Header.JWK != nil ***REMOVED***
				publicKey, err := UnmarshalPublicKeyJWK([]byte(parsedJSig.Header.JWK))
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				jsig.Header.JWK = publicKey
			***REMOVED***

			js.signatures = append(js.signatures, jsig)
		***REMOVED***
	***REMOVED***

	return js, nil
***REMOVED***

// NewJSONSignatureFromMap returns a new unsigned JSONSignature from a map or
// struct. JWS will need to be signed before serializing or storing.
func NewJSONSignatureFromMap(content interface***REMOVED******REMOVED***) (*JSONSignature, error) ***REMOVED***
	switch content.(type) ***REMOVED***
	case map[string]interface***REMOVED******REMOVED***:
	case struct***REMOVED******REMOVED***:
	default:
		return nil, errors.New("invalid data type")
	***REMOVED***

	js := newJSONSignature()
	js.indent = "   "

	payload, err := json.MarshalIndent(content, "", js.indent)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	js.payload = joseBase64UrlEncode(payload)

	// Remove '\n***REMOVED***' from formatted section, put in protected header
	js.formatLength = len(payload) - 2
	js.formatTail = payload[js.formatLength:]

	return js, nil
***REMOVED***

func readIntFromMap(key string, m map[string]interface***REMOVED******REMOVED***) (int, bool) ***REMOVED***
	value, ok := m[key]
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	switch v := value.(type) ***REMOVED***
	case int:
		return v, true
	case float64:
		return int(v), true
	default:
		return 0, false
	***REMOVED***
***REMOVED***

func readStringFromMap(key string, m map[string]interface***REMOVED******REMOVED***) (v string, ok bool) ***REMOVED***
	value, ok := m[key]
	if !ok ***REMOVED***
		return "", false
	***REMOVED***
	v, ok = value.(string)
	return
***REMOVED***

// ParsePrettySignature parses a formatted signature into a
// JSON signature. If the signatures are missing the format information
// an error is thrown. The formatted signature must be created by
// the same method as format signature.
func ParsePrettySignature(content []byte, signatureKey string) (*JSONSignature, error) ***REMOVED***
	var contentMap map[string]json.RawMessage
	err := json.Unmarshal(content, &contentMap)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error unmarshalling content: %s", err)
	***REMOVED***
	sigMessage, ok := contentMap[signatureKey]
	if !ok ***REMOVED***
		return nil, ErrMissingSignatureKey
	***REMOVED***

	var signatureBlocks []jsParsedSignature
	err = json.Unmarshal([]byte(sigMessage), &signatureBlocks)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error unmarshalling signatures: %s", err)
	***REMOVED***

	js := newJSONSignature()
	js.signatures = make([]jsSignature, len(signatureBlocks))

	for i, signatureBlock := range signatureBlocks ***REMOVED***
		protectedBytes, err := joseBase64UrlDecode(signatureBlock.Protected)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("base64 decode error: %s", err)
		***REMOVED***
		var protectedHeader map[string]interface***REMOVED******REMOVED***
		err = json.Unmarshal(protectedBytes, &protectedHeader)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("error unmarshalling protected header: %s", err)
		***REMOVED***

		formatLength, ok := readIntFromMap("formatLength", protectedHeader)
		if !ok ***REMOVED***
			return nil, errors.New("missing formatted length")
		***REMOVED***
		encodedTail, ok := readStringFromMap("formatTail", protectedHeader)
		if !ok ***REMOVED***
			return nil, errors.New("missing formatted tail")
		***REMOVED***
		formatTail, err := joseBase64UrlDecode(encodedTail)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("base64 decode error on tail: %s", err)
		***REMOVED***
		if js.formatLength == 0 ***REMOVED***
			js.formatLength = formatLength
		***REMOVED*** else if js.formatLength != formatLength ***REMOVED***
			return nil, errors.New("conflicting format length")
		***REMOVED***
		if len(js.formatTail) == 0 ***REMOVED***
			js.formatTail = formatTail
		***REMOVED*** else if bytes.Compare(js.formatTail, formatTail) != 0 ***REMOVED***
			return nil, errors.New("conflicting format tail")
		***REMOVED***

		header := jsHeader***REMOVED***
			Algorithm: signatureBlock.Header.Algorithm,
			Chain:     signatureBlock.Header.Chain,
		***REMOVED***
		if signatureBlock.Header.JWK != nil ***REMOVED***
			publicKey, err := UnmarshalPublicKeyJWK([]byte(signatureBlock.Header.JWK))
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("error unmarshalling public key: %s", err)
			***REMOVED***
			header.JWK = publicKey
		***REMOVED***
		js.signatures[i] = jsSignature***REMOVED***
			Header:    header,
			Signature: signatureBlock.Signature,
			Protected: signatureBlock.Protected,
		***REMOVED***
	***REMOVED***
	if js.formatLength > len(content) ***REMOVED***
		return nil, errors.New("invalid format length")
	***REMOVED***
	formatted := make([]byte, js.formatLength+len(js.formatTail))
	copy(formatted, content[:js.formatLength])
	copy(formatted[js.formatLength:], js.formatTail)
	js.indent = detectJSONIndent(formatted)
	js.payload = joseBase64UrlEncode(formatted)

	return js, nil
***REMOVED***

// PrettySignature formats a json signature into an easy to read
// single json serialized object.
func (js *JSONSignature) PrettySignature(signatureKey string) ([]byte, error) ***REMOVED***
	if len(js.signatures) == 0 ***REMOVED***
		return nil, errors.New("no signatures")
	***REMOVED***
	payload, err := joseBase64UrlDecode(js.payload)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	payload = payload[:js.formatLength]

	sort.Sort(jsSignaturesSorted(js.signatures))

	var marshalled []byte
	var marshallErr error
	if js.indent != "" ***REMOVED***
		marshalled, marshallErr = json.MarshalIndent(js.signatures, js.indent, js.indent)
	***REMOVED*** else ***REMOVED***
		marshalled, marshallErr = json.Marshal(js.signatures)
	***REMOVED***
	if marshallErr != nil ***REMOVED***
		return nil, marshallErr
	***REMOVED***

	buf := bytes.NewBuffer(make([]byte, 0, len(payload)+len(marshalled)+34))
	buf.Write(payload)
	buf.WriteByte(',')
	if js.indent != "" ***REMOVED***
		buf.WriteByte('\n')
		buf.WriteString(js.indent)
		buf.WriteByte('"')
		buf.WriteString(signatureKey)
		buf.WriteString("\": ")
		buf.Write(marshalled)
		buf.WriteByte('\n')
	***REMOVED*** else ***REMOVED***
		buf.WriteByte('"')
		buf.WriteString(signatureKey)
		buf.WriteString("\":")
		buf.Write(marshalled)
	***REMOVED***
	buf.WriteByte('***REMOVED***')

	return buf.Bytes(), nil
***REMOVED***

// Signatures provides the signatures on this JWS as opaque blobs, sorted by
// keyID. These blobs can be stored and reassembled with payloads. Internally,
// they are simply marshaled json web signatures but implementations should
// not rely on this.
func (js *JSONSignature) Signatures() ([][]byte, error) ***REMOVED***
	sort.Sort(jsSignaturesSorted(js.signatures))

	var sb [][]byte
	for _, jsig := range js.signatures ***REMOVED***
		p, err := json.Marshal(jsig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		sb = append(sb, p)
	***REMOVED***

	return sb, nil
***REMOVED***

// Merge combines the signatures from one or more other signatures into the
// method receiver. If the payloads differ for any argument, an error will be
// returned and the receiver will not be modified.
func (js *JSONSignature) Merge(others ...*JSONSignature) error ***REMOVED***
	merged := js.signatures
	for _, other := range others ***REMOVED***
		if js.payload != other.payload ***REMOVED***
			return fmt.Errorf("payloads differ from merge target")
		***REMOVED***
		merged = append(merged, other.signatures...)
	***REMOVED***

	js.signatures = merged
	return nil
***REMOVED***
