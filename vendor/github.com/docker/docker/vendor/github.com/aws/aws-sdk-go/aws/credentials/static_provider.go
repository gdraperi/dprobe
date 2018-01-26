package credentials

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// StaticProviderName provides a name of Static provider
const StaticProviderName = "StaticProvider"

var (
	// ErrStaticCredentialsEmpty is emitted when static credentials are empty.
	//
	// @readonly
	ErrStaticCredentialsEmpty = awserr.New("EmptyStaticCreds", "static credentials are empty", nil)
)

// A StaticProvider is a set of credentials which are set programmatically,
// and will never expire.
type StaticProvider struct ***REMOVED***
	Value
***REMOVED***

// NewStaticCredentials returns a pointer to a new Credentials object
// wrapping a static credentials value provider.
func NewStaticCredentials(id, secret, token string) *Credentials ***REMOVED***
	return NewCredentials(&StaticProvider***REMOVED***Value: Value***REMOVED***
		AccessKeyID:     id,
		SecretAccessKey: secret,
		SessionToken:    token,
	***REMOVED******REMOVED***)
***REMOVED***

// NewStaticCredentialsFromCreds returns a pointer to a new Credentials object
// wrapping the static credentials value provide. Same as NewStaticCredentials
// but takes the creds Value instead of individual fields
func NewStaticCredentialsFromCreds(creds Value) *Credentials ***REMOVED***
	return NewCredentials(&StaticProvider***REMOVED***Value: creds***REMOVED***)
***REMOVED***

// Retrieve returns the credentials or error if the credentials are invalid.
func (s *StaticProvider) Retrieve() (Value, error) ***REMOVED***
	if s.AccessKeyID == "" || s.SecretAccessKey == "" ***REMOVED***
		return Value***REMOVED***ProviderName: StaticProviderName***REMOVED***, ErrStaticCredentialsEmpty
	***REMOVED***

	if len(s.Value.ProviderName) == 0 ***REMOVED***
		s.Value.ProviderName = StaticProviderName
	***REMOVED***
	return s.Value, nil
***REMOVED***

// IsExpired returns if the credentials are expired.
//
// For StaticProvider, the credentials never expired.
func (s *StaticProvider) IsExpired() bool ***REMOVED***
	return false
***REMOVED***
