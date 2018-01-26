package credentials

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// EnvProviderName provides a name of Env provider
const EnvProviderName = "EnvProvider"

var (
	// ErrAccessKeyIDNotFound is returned when the AWS Access Key ID can't be
	// found in the process's environment.
	//
	// @readonly
	ErrAccessKeyIDNotFound = awserr.New("EnvAccessKeyNotFound", "AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY not found in environment", nil)

	// ErrSecretAccessKeyNotFound is returned when the AWS Secret Access Key
	// can't be found in the process's environment.
	//
	// @readonly
	ErrSecretAccessKeyNotFound = awserr.New("EnvSecretNotFound", "AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY not found in environment", nil)
)

// A EnvProvider retrieves credentials from the environment variables of the
// running process. Environment credentials never expire.
//
// Environment variables used:
//
// * Access Key ID:     AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
//
// * Secret Access Key: AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
type EnvProvider struct ***REMOVED***
	retrieved bool
***REMOVED***

// NewEnvCredentials returns a pointer to a new Credentials object
// wrapping the environment variable provider.
func NewEnvCredentials() *Credentials ***REMOVED***
	return NewCredentials(&EnvProvider***REMOVED******REMOVED***)
***REMOVED***

// Retrieve retrieves the keys from the environment.
func (e *EnvProvider) Retrieve() (Value, error) ***REMOVED***
	e.retrieved = false

	id := os.Getenv("AWS_ACCESS_KEY_ID")
	if id == "" ***REMOVED***
		id = os.Getenv("AWS_ACCESS_KEY")
	***REMOVED***

	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secret == "" ***REMOVED***
		secret = os.Getenv("AWS_SECRET_KEY")
	***REMOVED***

	if id == "" ***REMOVED***
		return Value***REMOVED***ProviderName: EnvProviderName***REMOVED***, ErrAccessKeyIDNotFound
	***REMOVED***

	if secret == "" ***REMOVED***
		return Value***REMOVED***ProviderName: EnvProviderName***REMOVED***, ErrSecretAccessKeyNotFound
	***REMOVED***

	e.retrieved = true
	return Value***REMOVED***
		AccessKeyID:     id,
		SecretAccessKey: secret,
		SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
		ProviderName:    EnvProviderName,
	***REMOVED***, nil
***REMOVED***

// IsExpired returns if the credentials have been retrieved.
func (e *EnvProvider) IsExpired() bool ***REMOVED***
	return !e.retrieved
***REMOVED***
