package schema1

import (
	"crypto/x509"

	"github.com/docker/libtrust"
	"github.com/sirupsen/logrus"
)

// Verify verifies the signature of the signed manifest returning the public
// keys used during signing.
func Verify(sm *SignedManifest) ([]libtrust.PublicKey, error) ***REMOVED***
	js, err := libtrust.ParsePrettySignature(sm.all, "signatures")
	if err != nil ***REMOVED***
		logrus.WithField("err", err).Debugf("(*SignedManifest).Verify")
		return nil, err
	***REMOVED***

	return js.Verify()
***REMOVED***

// VerifyChains verifies the signature of the signed manifest against the
// certificate pool returning the list of verified chains. Signatures without
// an x509 chain are not checked.
func VerifyChains(sm *SignedManifest, ca *x509.CertPool) ([][]*x509.Certificate, error) ***REMOVED***
	js, err := libtrust.ParsePrettySignature(sm.all, "signatures")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return js.VerifyChains(ca)
***REMOVED***
