package credentials

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var (
	// ErrNoValidProvidersFoundInChain Is returned when there are no valid
	// providers in the ChainProvider.
	//
	// This has been deprecated. For verbose error messaging set
	// aws.Config.CredentialsChainVerboseErrors to true
	//
	// @readonly
	ErrNoValidProvidersFoundInChain = awserr.New("NoCredentialProviders",
		`no valid providers in chain. Deprecated.
	For verbose messaging see aws.Config.CredentialsChainVerboseErrors`,
		nil)
)

// A ChainProvider will search for a provider which returns credentials
// and cache that provider until Retrieve is called again.
//
// The ChainProvider provides a way of chaining multiple providers together
// which will pick the first available using priority order of the Providers
// in the list.
//
// If none of the Providers retrieve valid credentials Value, ChainProvider's
// Retrieve() will return the error ErrNoValidProvidersFoundInChain.
//
// If a Provider is found which returns valid credentials Value ChainProvider
// will cache that Provider for all calls to IsExpired(), until Retrieve is
// called again.
//
// Example of ChainProvider to be used with an EnvProvider and EC2RoleProvider.
// In this example EnvProvider will first check if any credentials are available
// via the environment variables. If there are none ChainProvider will check
// the next Provider in the list, EC2RoleProvider in this case. If EC2RoleProvider
// does not return any credentials ChainProvider will return the error
// ErrNoValidProvidersFoundInChain
//
//     creds := credentials.NewChainCredentials(
//         []credentials.Provider***REMOVED***
//             &credentials.EnvProvider***REMOVED******REMOVED***,
//             &ec2rolecreds.EC2RoleProvider***REMOVED***
//                 Client: ec2metadata.New(sess),
//         ***REMOVED***,
//     ***REMOVED***)
//
//     // Usage of ChainCredentials with aws.Config
//     svc := ec2.New(session.Must(session.NewSession(&aws.Config***REMOVED***
//       Credentials: creds,
// ***REMOVED***)))
//
type ChainProvider struct ***REMOVED***
	Providers     []Provider
	curr          Provider
	VerboseErrors bool
***REMOVED***

// NewChainCredentials returns a pointer to a new Credentials object
// wrapping a chain of providers.
func NewChainCredentials(providers []Provider) *Credentials ***REMOVED***
	return NewCredentials(&ChainProvider***REMOVED***
		Providers: append([]Provider***REMOVED******REMOVED***, providers...),
	***REMOVED***)
***REMOVED***

// Retrieve returns the credentials value or error if no provider returned
// without error.
//
// If a provider is found it will be cached and any calls to IsExpired()
// will return the expired state of the cached provider.
func (c *ChainProvider) Retrieve() (Value, error) ***REMOVED***
	var errs []error
	for _, p := range c.Providers ***REMOVED***
		creds, err := p.Retrieve()
		if err == nil ***REMOVED***
			c.curr = p
			return creds, nil
		***REMOVED***
		errs = append(errs, err)
	***REMOVED***
	c.curr = nil

	var err error
	err = ErrNoValidProvidersFoundInChain
	if c.VerboseErrors ***REMOVED***
		err = awserr.NewBatchError("NoCredentialProviders", "no valid providers in chain", errs)
	***REMOVED***
	return Value***REMOVED******REMOVED***, err
***REMOVED***

// IsExpired will returned the expired state of the currently cached provider
// if there is one.  If there is no current provider, true will be returned.
func (c *ChainProvider) IsExpired() bool ***REMOVED***
	if c.curr != nil ***REMOVED***
		return c.curr.IsExpired()
	***REMOVED***

	return true
***REMOVED***
