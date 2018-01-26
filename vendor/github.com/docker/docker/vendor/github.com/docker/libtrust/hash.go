package libtrust

import (
	"crypto"
	_ "crypto/sha256" // Registrer SHA224 and SHA256
	_ "crypto/sha512" // Registrer SHA384 and SHA512
	"fmt"
)

type signatureAlgorithm struct ***REMOVED***
	algHeaderParam string
	hashID         crypto.Hash
***REMOVED***

func (h *signatureAlgorithm) HeaderParam() string ***REMOVED***
	return h.algHeaderParam
***REMOVED***

func (h *signatureAlgorithm) HashID() crypto.Hash ***REMOVED***
	return h.hashID
***REMOVED***

var (
	rs256 = &signatureAlgorithm***REMOVED***"RS256", crypto.SHA256***REMOVED***
	rs384 = &signatureAlgorithm***REMOVED***"RS384", crypto.SHA384***REMOVED***
	rs512 = &signatureAlgorithm***REMOVED***"RS512", crypto.SHA512***REMOVED***
	es256 = &signatureAlgorithm***REMOVED***"ES256", crypto.SHA256***REMOVED***
	es384 = &signatureAlgorithm***REMOVED***"ES384", crypto.SHA384***REMOVED***
	es512 = &signatureAlgorithm***REMOVED***"ES512", crypto.SHA512***REMOVED***
)

func rsaSignatureAlgorithmByName(alg string) (*signatureAlgorithm, error) ***REMOVED***
	switch ***REMOVED***
	case alg == "RS256":
		return rs256, nil
	case alg == "RS384":
		return rs384, nil
	case alg == "RS512":
		return rs512, nil
	default:
		return nil, fmt.Errorf("RSA Digital Signature Algorithm %q not supported", alg)
	***REMOVED***
***REMOVED***

func rsaPKCS1v15SignatureAlgorithmForHashID(hashID crypto.Hash) *signatureAlgorithm ***REMOVED***
	switch ***REMOVED***
	case hashID == crypto.SHA512:
		return rs512
	case hashID == crypto.SHA384:
		return rs384
	case hashID == crypto.SHA256:
		fallthrough
	default:
		return rs256
	***REMOVED***
***REMOVED***
