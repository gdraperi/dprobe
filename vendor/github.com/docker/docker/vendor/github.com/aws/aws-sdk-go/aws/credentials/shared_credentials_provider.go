package credentials

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/internal/shareddefaults"
)

// SharedCredsProviderName provides a name of SharedCreds provider
const SharedCredsProviderName = "SharedCredentialsProvider"

var (
	// ErrSharedCredentialsHomeNotFound is emitted when the user directory cannot be found.
	ErrSharedCredentialsHomeNotFound = awserr.New("UserHomeNotFound", "user home directory not found.", nil)
)

// A SharedCredentialsProvider retrieves credentials from the current user's home
// directory, and keeps track if those credentials are expired.
//
// Profile ini file example: $HOME/.aws/credentials
type SharedCredentialsProvider struct ***REMOVED***
	// Path to the shared credentials file.
	//
	// If empty will look for "AWS_SHARED_CREDENTIALS_FILE" env variable. If the
	// env value is empty will default to current user's home directory.
	// Linux/OSX: "$HOME/.aws/credentials"
	// Windows:   "%USERPROFILE%\.aws\credentials"
	Filename string

	// AWS Profile to extract credentials from the shared credentials file. If empty
	// will default to environment variable "AWS_PROFILE" or "default" if
	// environment variable is also not set.
	Profile string

	// retrieved states if the credentials have been successfully retrieved.
	retrieved bool
***REMOVED***

// NewSharedCredentials returns a pointer to a new Credentials object
// wrapping the Profile file provider.
func NewSharedCredentials(filename, profile string) *Credentials ***REMOVED***
	return NewCredentials(&SharedCredentialsProvider***REMOVED***
		Filename: filename,
		Profile:  profile,
	***REMOVED***)
***REMOVED***

// Retrieve reads and extracts the shared credentials from the current
// users home directory.
func (p *SharedCredentialsProvider) Retrieve() (Value, error) ***REMOVED***
	p.retrieved = false

	filename, err := p.filename()
	if err != nil ***REMOVED***
		return Value***REMOVED***ProviderName: SharedCredsProviderName***REMOVED***, err
	***REMOVED***

	creds, err := loadProfile(filename, p.profile())
	if err != nil ***REMOVED***
		return Value***REMOVED***ProviderName: SharedCredsProviderName***REMOVED***, err
	***REMOVED***

	p.retrieved = true
	return creds, nil
***REMOVED***

// IsExpired returns if the shared credentials have expired.
func (p *SharedCredentialsProvider) IsExpired() bool ***REMOVED***
	return !p.retrieved
***REMOVED***

// loadProfiles loads from the file pointed to by shared credentials filename for profile.
// The credentials retrieved from the profile will be returned or error. Error will be
// returned if it fails to read from the file, or the data is invalid.
func loadProfile(filename, profile string) (Value, error) ***REMOVED***
	config, err := ini.Load(filename)
	if err != nil ***REMOVED***
		return Value***REMOVED***ProviderName: SharedCredsProviderName***REMOVED***, awserr.New("SharedCredsLoad", "failed to load shared credentials file", err)
	***REMOVED***
	iniProfile, err := config.GetSection(profile)
	if err != nil ***REMOVED***
		return Value***REMOVED***ProviderName: SharedCredsProviderName***REMOVED***, awserr.New("SharedCredsLoad", "failed to get profile", err)
	***REMOVED***

	id, err := iniProfile.GetKey("aws_access_key_id")
	if err != nil ***REMOVED***
		return Value***REMOVED***ProviderName: SharedCredsProviderName***REMOVED***, awserr.New("SharedCredsAccessKey",
			fmt.Sprintf("shared credentials %s in %s did not contain aws_access_key_id", profile, filename),
			err)
	***REMOVED***

	secret, err := iniProfile.GetKey("aws_secret_access_key")
	if err != nil ***REMOVED***
		return Value***REMOVED***ProviderName: SharedCredsProviderName***REMOVED***, awserr.New("SharedCredsSecret",
			fmt.Sprintf("shared credentials %s in %s did not contain aws_secret_access_key", profile, filename),
			nil)
	***REMOVED***

	// Default to empty string if not found
	token := iniProfile.Key("aws_session_token")

	return Value***REMOVED***
		AccessKeyID:     id.String(),
		SecretAccessKey: secret.String(),
		SessionToken:    token.String(),
		ProviderName:    SharedCredsProviderName,
	***REMOVED***, nil
***REMOVED***

// filename returns the filename to use to read AWS shared credentials.
//
// Will return an error if the user's home directory path cannot be found.
func (p *SharedCredentialsProvider) filename() (string, error) ***REMOVED***
	if len(p.Filename) != 0 ***REMOVED***
		return p.Filename, nil
	***REMOVED***

	if p.Filename = os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); len(p.Filename) != 0 ***REMOVED***
		return p.Filename, nil
	***REMOVED***

	if home := shareddefaults.UserHomeDir(); len(home) == 0 ***REMOVED***
		// Backwards compatibility of home directly not found error being returned.
		// This error is too verbose, failure when opening the file would of been
		// a better error to return.
		return "", ErrSharedCredentialsHomeNotFound
	***REMOVED***

	p.Filename = shareddefaults.SharedCredentialsFilename()

	return p.Filename, nil
***REMOVED***

// profile returns the AWS shared credentials profile.  If empty will read
// environment variable "AWS_PROFILE". If that is not set profile will
// return "default".
func (p *SharedCredentialsProvider) profile() string ***REMOVED***
	if p.Profile == "" ***REMOVED***
		p.Profile = os.Getenv("AWS_PROFILE")
	***REMOVED***
	if p.Profile == "" ***REMOVED***
		p.Profile = "default"
	***REMOVED***

	return p.Profile
***REMOVED***
