// Package config contains the configuration logic for CFSSL.
package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cloudflare/cfssl/auth"
	cferr "github.com/cloudflare/cfssl/errors"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/log"
	ocspConfig "github.com/cloudflare/cfssl/ocsp/config"
)

// A CSRWhitelist stores booleans for fields in the CSR. If a CSRWhitelist is
// not present in a SigningProfile, all of these fields may be copied from the
// CSR into the signed certificate. If a CSRWhitelist *is* present in a
// SigningProfile, only those fields with a `true` value in the CSRWhitelist may
// be copied from the CSR to the signed certificate. Note that some of these
// fields, like Subject, can be provided or partially provided through the API.
// Since API clients are expected to be trusted, but CSRs are not, fields
// provided through the API are not subject to whitelisting through this
// mechanism.
type CSRWhitelist struct ***REMOVED***
	Subject, PublicKeyAlgorithm, PublicKey, SignatureAlgorithm bool
	DNSNames, IPAddresses, EmailAddresses                      bool
***REMOVED***

// OID is our own version of asn1's ObjectIdentifier, so we can define a custom
// JSON marshal / unmarshal.
type OID asn1.ObjectIdentifier

// CertificatePolicy represents the ASN.1 PolicyInformation structure from
// https://tools.ietf.org/html/rfc3280.html#page-106.
// Valid values of Type are "id-qt-unotice" and "id-qt-cps"
type CertificatePolicy struct ***REMOVED***
	ID         OID
	Qualifiers []CertificatePolicyQualifier
***REMOVED***

// CertificatePolicyQualifier represents a single qualifier from an ASN.1
// PolicyInformation structure.
type CertificatePolicyQualifier struct ***REMOVED***
	Type  string
	Value string
***REMOVED***

// AuthRemote is an authenticated remote signer.
type AuthRemote struct ***REMOVED***
	RemoteName  string `json:"remote"`
	AuthKeyName string `json:"auth_key"`
***REMOVED***

// CAConstraint specifies various CA constraints on the signed certificate.
// CAConstraint would verify against (and override) the CA
// extensions in the given CSR.
type CAConstraint struct ***REMOVED***
	IsCA           bool `json:"is_ca"`
	MaxPathLen     int  `json:"max_path_len"`
	MaxPathLenZero bool `json:"max_path_len_zero"`
***REMOVED***

// A SigningProfile stores information that the CA needs to store
// signature policy.
type SigningProfile struct ***REMOVED***
	Usage               []string     `json:"usages"`
	IssuerURL           []string     `json:"issuer_urls"`
	OCSP                string       `json:"ocsp_url"`
	CRL                 string       `json:"crl_url"`
	CAConstraint        CAConstraint `json:"ca_constraint"`
	OCSPNoCheck         bool         `json:"ocsp_no_check"`
	ExpiryString        string       `json:"expiry"`
	BackdateString      string       `json:"backdate"`
	AuthKeyName         string       `json:"auth_key"`
	RemoteName          string       `json:"remote"`
	NotBefore           time.Time    `json:"not_before"`
	NotAfter            time.Time    `json:"not_after"`
	NameWhitelistString string       `json:"name_whitelist"`
	AuthRemote          AuthRemote   `json:"auth_remote"`
	CTLogServers        []string     `json:"ct_log_servers"`
	AllowedExtensions   []OID        `json:"allowed_extensions"`
	CertStore           string       `json:"cert_store"`

	Policies                    []CertificatePolicy
	Expiry                      time.Duration
	Backdate                    time.Duration
	Provider                    auth.Provider
	RemoteProvider              auth.Provider
	RemoteServer                string
	RemoteCAs                   *x509.CertPool
	ClientCert                  *tls.Certificate
	CSRWhitelist                *CSRWhitelist
	NameWhitelist               *regexp.Regexp
	ExtensionWhitelist          map[string]bool
	ClientProvidesSerialNumbers bool
***REMOVED***

// UnmarshalJSON unmarshals a JSON string into an OID.
func (oid *OID) UnmarshalJSON(data []byte) (err error) ***REMOVED***
	if data[0] != '"' || data[len(data)-1] != '"' ***REMOVED***
		return errors.New("OID JSON string not wrapped in quotes." + string(data))
	***REMOVED***
	data = data[1 : len(data)-1]
	parsedOid, err := parseObjectIdentifier(string(data))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*oid = OID(parsedOid)
	return
***REMOVED***

// MarshalJSON marshals an oid into a JSON string.
func (oid OID) MarshalJSON() ([]byte, error) ***REMOVED***
	return []byte(fmt.Sprintf(`"%v"`, asn1.ObjectIdentifier(oid))), nil
***REMOVED***

func parseObjectIdentifier(oidString string) (oid asn1.ObjectIdentifier, err error) ***REMOVED***
	validOID, err := regexp.MatchString("\\d(\\.\\d+)*", oidString)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if !validOID ***REMOVED***
		err = errors.New("Invalid OID")
		return
	***REMOVED***

	segments := strings.Split(oidString, ".")
	oid = make(asn1.ObjectIdentifier, len(segments))
	for i, intString := range segments ***REMOVED***
		oid[i], err = strconv.Atoi(intString)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

const timeFormat = "2006-01-02T15:04:05"

// populate is used to fill in the fields that are not in JSON
//
// First, the ExpiryString parameter is needed to parse
// expiration timestamps from JSON. The JSON decoder is not able to
// decode a string time duration to a time.Duration, so this is called
// when loading the configuration to properly parse and fill out the
// Expiry parameter.
// This function is also used to create references to the auth key
// and default remote for the profile.
// It returns true if ExpiryString is a valid representation of a
// time.Duration, and the AuthKeyString and RemoteName point to
// valid objects. It returns false otherwise.
func (p *SigningProfile) populate(cfg *Config) error ***REMOVED***
	if p == nil ***REMOVED***
		return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy, errors.New("can't parse nil profile"))
	***REMOVED***

	var err error
	if p.RemoteName == "" && p.AuthRemote.RemoteName == "" ***REMOVED***
		log.Debugf("parse expiry in profile")
		if p.ExpiryString == "" ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy, errors.New("empty expiry string"))
		***REMOVED***

		dur, err := time.ParseDuration(p.ExpiryString)
		if err != nil ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy, err)
		***REMOVED***

		log.Debugf("expiry is valid")
		p.Expiry = dur

		if p.BackdateString != "" ***REMOVED***
			dur, err = time.ParseDuration(p.BackdateString)
			if err != nil ***REMOVED***
				return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy, err)
			***REMOVED***

			p.Backdate = dur
		***REMOVED***

		if !p.NotBefore.IsZero() && !p.NotAfter.IsZero() && p.NotAfter.Before(p.NotBefore) ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy, err)
		***REMOVED***

		if len(p.Policies) > 0 ***REMOVED***
			for _, policy := range p.Policies ***REMOVED***
				for _, qualifier := range policy.Qualifiers ***REMOVED***
					if qualifier.Type != "" && qualifier.Type != "id-qt-unotice" && qualifier.Type != "id-qt-cps" ***REMOVED***
						return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
							errors.New("invalid policy qualifier type"))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if p.RemoteName != "" ***REMOVED***
		log.Debug("match remote in profile to remotes section")
		if p.AuthRemote.RemoteName != "" ***REMOVED***
			log.Error("profile has both a remote and an auth remote specified")
			return cferr.New(cferr.PolicyError, cferr.InvalidPolicy)
		***REMOVED***
		if remote := cfg.Remotes[p.RemoteName]; remote != "" ***REMOVED***
			if err := p.updateRemote(remote); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
				errors.New("failed to find remote in remotes section"))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		log.Debug("match auth remote in profile to remotes section")
		if remote := cfg.Remotes[p.AuthRemote.RemoteName]; remote != "" ***REMOVED***
			if err := p.updateRemote(remote); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
				errors.New("failed to find remote in remotes section"))
		***REMOVED***
	***REMOVED***

	if p.AuthKeyName != "" ***REMOVED***
		log.Debug("match auth key in profile to auth_keys section")
		if key, ok := cfg.AuthKeys[p.AuthKeyName]; ok == true ***REMOVED***
			if key.Type == "standard" ***REMOVED***
				p.Provider, err = auth.New(key.Key, nil)
				if err != nil ***REMOVED***
					log.Debugf("failed to create new standard auth provider: %v", err)
					return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
						errors.New("failed to create new standard auth provider"))
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				log.Debugf("unknown authentication type %v", key.Type)
				return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
					errors.New("unknown authentication type"))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
				errors.New("failed to find auth_key in auth_keys section"))
		***REMOVED***
	***REMOVED***

	if p.AuthRemote.AuthKeyName != "" ***REMOVED***
		log.Debug("match auth remote key in profile to auth_keys section")
		if key, ok := cfg.AuthKeys[p.AuthRemote.AuthKeyName]; ok == true ***REMOVED***
			if key.Type == "standard" ***REMOVED***
				p.RemoteProvider, err = auth.New(key.Key, nil)
				if err != nil ***REMOVED***
					log.Debugf("failed to create new standard auth provider: %v", err)
					return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
						errors.New("failed to create new standard auth provider"))
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				log.Debugf("unknown authentication type %v", key.Type)
				return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
					errors.New("unknown authentication type"))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
				errors.New("failed to find auth_remote's auth_key in auth_keys section"))
		***REMOVED***
	***REMOVED***

	if p.NameWhitelistString != "" ***REMOVED***
		log.Debug("compiling whitelist regular expression")
		rule, err := regexp.Compile(p.NameWhitelistString)
		if err != nil ***REMOVED***
			return cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
				errors.New("failed to compile name whitelist section"))
		***REMOVED***
		p.NameWhitelist = rule
	***REMOVED***

	p.ExtensionWhitelist = map[string]bool***REMOVED******REMOVED***
	for _, oid := range p.AllowedExtensions ***REMOVED***
		p.ExtensionWhitelist[asn1.ObjectIdentifier(oid).String()] = true
	***REMOVED***

	return nil
***REMOVED***

// updateRemote takes a signing profile and initializes the remote server object
// to the hostname:port combination sent by remote.
func (p *SigningProfile) updateRemote(remote string) error ***REMOVED***
	if remote != "" ***REMOVED***
		p.RemoteServer = remote
	***REMOVED***
	return nil
***REMOVED***

// OverrideRemotes takes a signing configuration and updates the remote server object
// to the hostname:port combination sent by remote
func (p *Signing) OverrideRemotes(remote string) error ***REMOVED***
	if remote != "" ***REMOVED***
		var err error
		for _, profile := range p.Profiles ***REMOVED***
			err = profile.updateRemote(remote)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		err = p.Default.updateRemote(remote)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// SetClientCertKeyPairFromFile updates the properties to set client certificates for mutual
// authenticated TLS remote requests
func (p *Signing) SetClientCertKeyPairFromFile(certFile string, keyFile string) error ***REMOVED***
	if certFile != "" && keyFile != "" ***REMOVED***
		cert, err := helpers.LoadClientCertificate(certFile, keyFile)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, profile := range p.Profiles ***REMOVED***
			profile.ClientCert = cert
		***REMOVED***
		p.Default.ClientCert = cert
	***REMOVED***
	return nil
***REMOVED***

// SetRemoteCAsFromFile reads root CAs from file and updates the properties to set remote CAs for TLS
// remote requests
func (p *Signing) SetRemoteCAsFromFile(caFile string) error ***REMOVED***
	if caFile != "" ***REMOVED***
		remoteCAs, err := helpers.LoadPEMCertPool(caFile)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		p.SetRemoteCAs(remoteCAs)
	***REMOVED***
	return nil
***REMOVED***

// SetRemoteCAs updates the properties to set remote CAs for TLS
// remote requests
func (p *Signing) SetRemoteCAs(remoteCAs *x509.CertPool) ***REMOVED***
	for _, profile := range p.Profiles ***REMOVED***
		profile.RemoteCAs = remoteCAs
	***REMOVED***
	p.Default.RemoteCAs = remoteCAs
***REMOVED***

// NeedsRemoteSigner returns true if one of the profiles has a remote set
func (p *Signing) NeedsRemoteSigner() bool ***REMOVED***
	for _, profile := range p.Profiles ***REMOVED***
		if profile.RemoteServer != "" ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	if p.Default.RemoteServer != "" ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// NeedsLocalSigner returns true if one of the profiles doe not have a remote set
func (p *Signing) NeedsLocalSigner() bool ***REMOVED***
	for _, profile := range p.Profiles ***REMOVED***
		if profile.RemoteServer == "" ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	if p.Default.RemoteServer == "" ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// Usages parses the list of key uses in the profile, translating them
// to a list of X.509 key usages and extended key usages.  The unknown
// uses are collected into a slice that is also returned.
func (p *SigningProfile) Usages() (ku x509.KeyUsage, eku []x509.ExtKeyUsage, unk []string) ***REMOVED***
	for _, keyUse := range p.Usage ***REMOVED***
		if kuse, ok := KeyUsage[keyUse]; ok ***REMOVED***
			ku |= kuse
		***REMOVED*** else if ekuse, ok := ExtKeyUsage[keyUse]; ok ***REMOVED***
			eku = append(eku, ekuse)
		***REMOVED*** else ***REMOVED***
			unk = append(unk, keyUse)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// A valid profile must be a valid local profile or a valid remote profile.
// A valid local profile has defined at least key usages to be used, and a
// valid local default profile has defined at least a default expiration.
// A valid remote profile (default or not) has remote signer initialized.
// In addition, a remote profile must has a valid auth provider if auth
// key defined.
func (p *SigningProfile) validProfile(isDefault bool) bool ***REMOVED***
	if p == nil ***REMOVED***
		return false
	***REMOVED***

	if p.AuthRemote.RemoteName == "" && p.AuthRemote.AuthKeyName != "" ***REMOVED***
		log.Debugf("invalid auth remote profile: no remote signer specified")
		return false
	***REMOVED***

	if p.RemoteName != "" ***REMOVED***
		log.Debugf("validate remote profile")

		if p.RemoteServer == "" ***REMOVED***
			log.Debugf("invalid remote profile: no remote signer specified")
			return false
		***REMOVED***

		if p.AuthKeyName != "" && p.Provider == nil ***REMOVED***
			log.Debugf("invalid remote profile: auth key name is defined but no auth provider is set")
			return false
		***REMOVED***

		if p.AuthRemote.RemoteName != "" ***REMOVED***
			log.Debugf("invalid remote profile: auth remote is also specified")
			return false
		***REMOVED***
	***REMOVED*** else if p.AuthRemote.RemoteName != "" ***REMOVED***
		log.Debugf("validate auth remote profile")
		if p.RemoteServer == "" ***REMOVED***
			log.Debugf("invalid auth remote profile: no remote signer specified")
			return false
		***REMOVED***

		if p.AuthRemote.AuthKeyName == "" || p.RemoteProvider == nil ***REMOVED***
			log.Debugf("invalid auth remote profile: no auth key is defined")
			return false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		log.Debugf("validate local profile")
		if !isDefault ***REMOVED***
			if len(p.Usage) == 0 ***REMOVED***
				log.Debugf("invalid local profile: no usages specified")
				return false
			***REMOVED*** else if _, _, unk := p.Usages(); len(unk) == len(p.Usage) ***REMOVED***
				log.Debugf("invalid local profile: no valid usages")
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if p.Expiry == 0 ***REMOVED***
				log.Debugf("invalid local profile: no expiry set")
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	log.Debugf("profile is valid")
	return true
***REMOVED***

// This checks if the SigningProfile object contains configurations that are only effective with a local signer
// which has access to CA private key.
func (p *SigningProfile) hasLocalConfig() bool ***REMOVED***
	if p.Usage != nil ||
		p.IssuerURL != nil ||
		p.OCSP != "" ||
		p.ExpiryString != "" ||
		p.BackdateString != "" ||
		p.CAConstraint.IsCA != false ||
		!p.NotBefore.IsZero() ||
		!p.NotAfter.IsZero() ||
		p.NameWhitelistString != "" ||
		len(p.CTLogServers) != 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// warnSkippedSettings prints a log warning message about skipped settings
// in a SigningProfile, usually due to remote signer.
func (p *Signing) warnSkippedSettings() ***REMOVED***
	const warningMessage = `The configuration value by "usages", "issuer_urls", "ocsp_url", "crl_url", "ca_constraint", "expiry", "backdate", "not_before", "not_after", "cert_store" and "ct_log_servers" are skipped`
	if p == nil ***REMOVED***
		return
	***REMOVED***

	if (p.Default.RemoteName != "" || p.Default.AuthRemote.RemoteName != "") && p.Default.hasLocalConfig() ***REMOVED***
		log.Warning("default profile points to a remote signer: ", warningMessage)
	***REMOVED***

	for name, profile := range p.Profiles ***REMOVED***
		if (profile.RemoteName != "" || profile.AuthRemote.RemoteName != "") && profile.hasLocalConfig() ***REMOVED***
			log.Warningf("Profiles[%s] points to a remote signer: %s", name, warningMessage)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Signing codifies the signature configuration policy for a CA.
type Signing struct ***REMOVED***
	Profiles map[string]*SigningProfile `json:"profiles"`
	Default  *SigningProfile            `json:"default"`
***REMOVED***

// Config stores configuration information for the CA.
type Config struct ***REMOVED***
	Signing  *Signing           `json:"signing"`
	OCSP     *ocspConfig.Config `json:"ocsp"`
	AuthKeys map[string]AuthKey `json:"auth_keys,omitempty"`
	Remotes  map[string]string  `json:"remotes,omitempty"`
***REMOVED***

// Valid ensures that Config is a valid configuration. It should be
// called immediately after parsing a configuration file.
func (c *Config) Valid() bool ***REMOVED***
	return c.Signing.Valid()
***REMOVED***

// Valid checks the signature policies, ensuring they are valid
// policies. A policy is valid if it has defined at least key usages
// to be used, and a valid default profile has defined at least a
// default expiration.
func (p *Signing) Valid() bool ***REMOVED***
	if p == nil ***REMOVED***
		return false
	***REMOVED***

	log.Debugf("validating configuration")
	if !p.Default.validProfile(true) ***REMOVED***
		log.Debugf("default profile is invalid")
		return false
	***REMOVED***

	for _, sp := range p.Profiles ***REMOVED***
		if !sp.validProfile(false) ***REMOVED***
			log.Debugf("invalid profile")
			return false
		***REMOVED***
	***REMOVED***

	p.warnSkippedSettings()

	return true
***REMOVED***

// KeyUsage contains a mapping of string names to key usages.
var KeyUsage = map[string]x509.KeyUsage***REMOVED***
	"signing":             x509.KeyUsageDigitalSignature,
	"digital signature":   x509.KeyUsageDigitalSignature,
	"content committment": x509.KeyUsageContentCommitment,
	"key encipherment":    x509.KeyUsageKeyEncipherment,
	"key agreement":       x509.KeyUsageKeyAgreement,
	"data encipherment":   x509.KeyUsageDataEncipherment,
	"cert sign":           x509.KeyUsageCertSign,
	"crl sign":            x509.KeyUsageCRLSign,
	"encipher only":       x509.KeyUsageEncipherOnly,
	"decipher only":       x509.KeyUsageDecipherOnly,
***REMOVED***

// ExtKeyUsage contains a mapping of string names to extended key
// usages.
var ExtKeyUsage = map[string]x509.ExtKeyUsage***REMOVED***
	"any":              x509.ExtKeyUsageAny,
	"server auth":      x509.ExtKeyUsageServerAuth,
	"client auth":      x509.ExtKeyUsageClientAuth,
	"code signing":     x509.ExtKeyUsageCodeSigning,
	"email protection": x509.ExtKeyUsageEmailProtection,
	"s/mime":           x509.ExtKeyUsageEmailProtection,
	"ipsec end system": x509.ExtKeyUsageIPSECEndSystem,
	"ipsec tunnel":     x509.ExtKeyUsageIPSECTunnel,
	"ipsec user":       x509.ExtKeyUsageIPSECUser,
	"timestamping":     x509.ExtKeyUsageTimeStamping,
	"ocsp signing":     x509.ExtKeyUsageOCSPSigning,
	"microsoft sgc":    x509.ExtKeyUsageMicrosoftServerGatedCrypto,
	"netscape sgc":     x509.ExtKeyUsageNetscapeServerGatedCrypto,
***REMOVED***

// An AuthKey contains an entry for a key used for authentication.
type AuthKey struct ***REMOVED***
	// Type contains information needed to select the appropriate
	// constructor. For example, "standard" for HMAC-SHA-256,
	// "standard-ip" for HMAC-SHA-256 incorporating the client's
	// IP.
	Type string `json:"type"`
	// Key contains the key information, such as a hex-encoded
	// HMAC key.
	Key string `json:"key"`
***REMOVED***

// DefaultConfig returns a default configuration specifying basic key
// usage and a 1 year expiration time. The key usages chosen are
// signing, key encipherment, client auth and server auth.
func DefaultConfig() *SigningProfile ***REMOVED***
	d := helpers.OneYear
	return &SigningProfile***REMOVED***
		Usage:        []string***REMOVED***"signing", "key encipherment", "server auth", "client auth"***REMOVED***,
		Expiry:       d,
		ExpiryString: "8760h",
	***REMOVED***
***REMOVED***

// LoadFile attempts to load the configuration file stored at the path
// and returns the configuration. On error, it returns nil.
func LoadFile(path string) (*Config, error) ***REMOVED***
	log.Debugf("loading configuration file from %s", path)
	if path == "" ***REMOVED***
		return nil, cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy, errors.New("invalid path"))
	***REMOVED***

	body, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return nil, cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy, errors.New("could not read configuration file"))
	***REMOVED***

	return LoadConfig(body)
***REMOVED***

// LoadConfig attempts to load the configuration from a byte slice.
// On error, it returns nil.
func LoadConfig(config []byte) (*Config, error) ***REMOVED***
	var cfg = &Config***REMOVED******REMOVED***
	err := json.Unmarshal(config, &cfg)
	if err != nil ***REMOVED***
		return nil, cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy,
			errors.New("failed to unmarshal configuration: "+err.Error()))
	***REMOVED***

	if cfg.Signing == nil ***REMOVED***
		return nil, errors.New("No \"signing\" field present")
	***REMOVED***

	if cfg.Signing.Default == nil ***REMOVED***
		log.Debugf("no default given: using default config")
		cfg.Signing.Default = DefaultConfig()
	***REMOVED*** else ***REMOVED***
		if err := cfg.Signing.Default.populate(cfg); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	for k := range cfg.Signing.Profiles ***REMOVED***
		if err := cfg.Signing.Profiles[k].populate(cfg); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if !cfg.Valid() ***REMOVED***
		return nil, cferr.Wrap(cferr.PolicyError, cferr.InvalidPolicy, errors.New("invalid configuration"))
	***REMOVED***

	log.Debugf("configuration ok")
	return cfg, nil
***REMOVED***
