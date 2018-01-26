// Package credentials provides credential retrieval and management
//
// The Credentials is the primary method of getting access to and managing
// credentials Values. Using dependency injection retrieval of the credential
// values is handled by a object which satisfies the Provider interface.
//
// By default the Credentials.Get() will cache the successful result of a
// Provider's Retrieve() until Provider.IsExpired() returns true. At which
// point Credentials will call Provider's Retrieve() to get new credential Value.
//
// The Provider is responsible for determining when credentials Value have expired.
// It is also important to note that Credentials will always call Retrieve the
// first time Credentials.Get() is called.
//
// Example of using the environment variable credentials.
//
//     creds := credentials.NewEnvCredentials()
//
//     // Retrieve the credentials value
//     credValue, err := creds.Get()
//     if err != nil ***REMOVED***
//         // handle error
// ***REMOVED***
//
// Example of forcing credentials to expire and be refreshed on the next Get().
// This may be helpful to proactively expire credentials and refresh them sooner
// than they would naturally expire on their own.
//
//     creds := credentials.NewCredentials(&ec2rolecreds.EC2RoleProvider***REMOVED******REMOVED***)
//     creds.Expire()
//     credsValue, err := creds.Get()
//     // New credentials will be retrieved instead of from cache.
//
//
// Custom Provider
//
// Each Provider built into this package also provides a helper method to generate
// a Credentials pointer setup with the provider. To use a custom Provider just
// create a type which satisfies the Provider interface and pass it to the
// NewCredentials method.
//
//     type MyProvider struct***REMOVED******REMOVED***
//     func (m *MyProvider) Retrieve() (Value, error) ***REMOVED***...***REMOVED***
//     func (m *MyProvider) IsExpired() bool ***REMOVED***...***REMOVED***
//
//     creds := credentials.NewCredentials(&MyProvider***REMOVED******REMOVED***)
//     credValue, err := creds.Get()
//
package credentials

import (
	"sync"
	"time"
)

// AnonymousCredentials is an empty Credential object that can be used as
// dummy placeholder credentials for requests that do not need signed.
//
// This Credentials can be used to configure a service to not sign requests
// when making service API calls. For example, when accessing public
// s3 buckets.
//
//     svc := s3.New(session.Must(session.NewSession(&aws.Config***REMOVED***
//       Credentials: credentials.AnonymousCredentials,
// ***REMOVED***)))
//     // Access public S3 buckets.
//
// @readonly
var AnonymousCredentials = NewStaticCredentials("", "", "")

// A Value is the AWS credentials value for individual credential fields.
type Value struct ***REMOVED***
	// AWS Access key ID
	AccessKeyID string

	// AWS Secret Access Key
	SecretAccessKey string

	// AWS Session Token
	SessionToken string

	// Provider used to get credentials
	ProviderName string
***REMOVED***

// A Provider is the interface for any component which will provide credentials
// Value. A provider is required to manage its own Expired state, and what to
// be expired means.
//
// The Provider should not need to implement its own mutexes, because
// that will be managed by Credentials.
type Provider interface ***REMOVED***
	// Retrieve returns nil if it successfully retrieved the value.
	// Error is returned if the value were not obtainable, or empty.
	Retrieve() (Value, error)

	// IsExpired returns if the credentials are no longer valid, and need
	// to be retrieved.
	IsExpired() bool
***REMOVED***

// An ErrorProvider is a stub credentials provider that always returns an error
// this is used by the SDK when construction a known provider is not possible
// due to an error.
type ErrorProvider struct ***REMOVED***
	// The error to be returned from Retrieve
	Err error

	// The provider name to set on the Retrieved returned Value
	ProviderName string
***REMOVED***

// Retrieve will always return the error that the ErrorProvider was created with.
func (p ErrorProvider) Retrieve() (Value, error) ***REMOVED***
	return Value***REMOVED***ProviderName: p.ProviderName***REMOVED***, p.Err
***REMOVED***

// IsExpired will always return not expired.
func (p ErrorProvider) IsExpired() bool ***REMOVED***
	return false
***REMOVED***

// A Expiry provides shared expiration logic to be used by credentials
// providers to implement expiry functionality.
//
// The best method to use this struct is as an anonymous field within the
// provider's struct.
//
// Example:
//     type EC2RoleProvider struct ***REMOVED***
//         Expiry
//         ...
// ***REMOVED***
type Expiry struct ***REMOVED***
	// The date/time when to expire on
	expiration time.Time

	// If set will be used by IsExpired to determine the current time.
	// Defaults to time.Now if CurrentTime is not set.  Available for testing
	// to be able to mock out the current time.
	CurrentTime func() time.Time
***REMOVED***

// SetExpiration sets the expiration IsExpired will check when called.
//
// If window is greater than 0 the expiration time will be reduced by the
// window value.
//
// Using a window is helpful to trigger credentials to expire sooner than
// the expiration time given to ensure no requests are made with expired
// tokens.
func (e *Expiry) SetExpiration(expiration time.Time, window time.Duration) ***REMOVED***
	e.expiration = expiration
	if window > 0 ***REMOVED***
		e.expiration = e.expiration.Add(-window)
	***REMOVED***
***REMOVED***

// IsExpired returns if the credentials are expired.
func (e *Expiry) IsExpired() bool ***REMOVED***
	if e.CurrentTime == nil ***REMOVED***
		e.CurrentTime = time.Now
	***REMOVED***
	return e.expiration.Before(e.CurrentTime())
***REMOVED***

// A Credentials provides synchronous safe retrieval of AWS credentials Value.
// Credentials will cache the credentials value until they expire. Once the value
// expires the next Get will attempt to retrieve valid credentials.
//
// Credentials is safe to use across multiple goroutines and will manage the
// synchronous state so the Providers do not need to implement their own
// synchronization.
//
// The first Credentials.Get() will always call Provider.Retrieve() to get the
// first instance of the credentials Value. All calls to Get() after that
// will return the cached credentials Value until IsExpired() returns true.
type Credentials struct ***REMOVED***
	creds        Value
	forceRefresh bool
	m            sync.Mutex

	provider Provider
***REMOVED***

// NewCredentials returns a pointer to a new Credentials with the provider set.
func NewCredentials(provider Provider) *Credentials ***REMOVED***
	return &Credentials***REMOVED***
		provider:     provider,
		forceRefresh: true,
	***REMOVED***
***REMOVED***

// Get returns the credentials value, or error if the credentials Value failed
// to be retrieved.
//
// Will return the cached credentials Value if it has not expired. If the
// credentials Value has expired the Provider's Retrieve() will be called
// to refresh the credentials.
//
// If Credentials.Expire() was called the credentials Value will be force
// expired, and the next call to Get() will cause them to be refreshed.
func (c *Credentials) Get() (Value, error) ***REMOVED***
	c.m.Lock()
	defer c.m.Unlock()

	if c.isExpired() ***REMOVED***
		creds, err := c.provider.Retrieve()
		if err != nil ***REMOVED***
			return Value***REMOVED******REMOVED***, err
		***REMOVED***
		c.creds = creds
		c.forceRefresh = false
	***REMOVED***

	return c.creds, nil
***REMOVED***

// Expire expires the credentials and forces them to be retrieved on the
// next call to Get().
//
// This will override the Provider's expired state, and force Credentials
// to call the Provider's Retrieve().
func (c *Credentials) Expire() ***REMOVED***
	c.m.Lock()
	defer c.m.Unlock()

	c.forceRefresh = true
***REMOVED***

// IsExpired returns if the credentials are no longer valid, and need
// to be retrieved.
//
// If the Credentials were forced to be expired with Expire() this will
// reflect that override.
func (c *Credentials) IsExpired() bool ***REMOVED***
	c.m.Lock()
	defer c.m.Unlock()

	return c.isExpired()
***REMOVED***

// isExpired helper method wrapping the definition of expired credentials.
func (c *Credentials) isExpired() bool ***REMOVED***
	return c.forceRefresh || c.provider.IsExpired()
***REMOVED***
