package ec2rolecreds

import (
	"bufio"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

// ProviderName provides a name of EC2Role provider
const ProviderName = "EC2RoleProvider"

// A EC2RoleProvider retrieves credentials from the EC2 service, and keeps track if
// those credentials are expired.
//
// Example how to configure the EC2RoleProvider with custom http Client, Endpoint
// or ExpiryWindow
//
//     p := &ec2rolecreds.EC2RoleProvider***REMOVED***
//         // Pass in a custom timeout to be used when requesting
//         // IAM EC2 Role credentials.
//         Client: ec2metadata.New(sess, aws.Config***REMOVED***
//             HTTPClient: &http.Client***REMOVED***Timeout: 10 * time.Second***REMOVED***,
//     ***REMOVED***),
//
//         // Do not use early expiry of credentials. If a non zero value is
//         // specified the credentials will be expired early
//         ExpiryWindow: 0,
// ***REMOVED***
type EC2RoleProvider struct ***REMOVED***
	credentials.Expiry

	// Required EC2Metadata client to use when connecting to EC2 metadata service.
	Client *ec2metadata.EC2Metadata

	// ExpiryWindow will allow the credentials to trigger refreshing prior to
	// the credentials actually expiring. This is beneficial so race conditions
	// with expiring credentials do not cause request to fail unexpectedly
	// due to ExpiredTokenException exceptions.
	//
	// So a ExpiryWindow of 10s would cause calls to IsExpired() to return true
	// 10 seconds before the credentials are actually expired.
	//
	// If ExpiryWindow is 0 or less it will be ignored.
	ExpiryWindow time.Duration
***REMOVED***

// NewCredentials returns a pointer to a new Credentials object wrapping
// the EC2RoleProvider. Takes a ConfigProvider to create a EC2Metadata client.
// The ConfigProvider is satisfied by the session.Session type.
func NewCredentials(c client.ConfigProvider, options ...func(*EC2RoleProvider)) *credentials.Credentials ***REMOVED***
	p := &EC2RoleProvider***REMOVED***
		Client: ec2metadata.New(c),
	***REMOVED***

	for _, option := range options ***REMOVED***
		option(p)
	***REMOVED***

	return credentials.NewCredentials(p)
***REMOVED***

// NewCredentialsWithClient returns a pointer to a new Credentials object wrapping
// the EC2RoleProvider. Takes a EC2Metadata client to use when connecting to EC2
// metadata service.
func NewCredentialsWithClient(client *ec2metadata.EC2Metadata, options ...func(*EC2RoleProvider)) *credentials.Credentials ***REMOVED***
	p := &EC2RoleProvider***REMOVED***
		Client: client,
	***REMOVED***

	for _, option := range options ***REMOVED***
		option(p)
	***REMOVED***

	return credentials.NewCredentials(p)
***REMOVED***

// Retrieve retrieves credentials from the EC2 service.
// Error will be returned if the request fails, or unable to extract
// the desired credentials.
func (m *EC2RoleProvider) Retrieve() (credentials.Value, error) ***REMOVED***
	credsList, err := requestCredList(m.Client)
	if err != nil ***REMOVED***
		return credentials.Value***REMOVED***ProviderName: ProviderName***REMOVED***, err
	***REMOVED***

	if len(credsList) == 0 ***REMOVED***
		return credentials.Value***REMOVED***ProviderName: ProviderName***REMOVED***, awserr.New("EmptyEC2RoleList", "empty EC2 Role list", nil)
	***REMOVED***
	credsName := credsList[0]

	roleCreds, err := requestCred(m.Client, credsName)
	if err != nil ***REMOVED***
		return credentials.Value***REMOVED***ProviderName: ProviderName***REMOVED***, err
	***REMOVED***

	m.SetExpiration(roleCreds.Expiration, m.ExpiryWindow)

	return credentials.Value***REMOVED***
		AccessKeyID:     roleCreds.AccessKeyID,
		SecretAccessKey: roleCreds.SecretAccessKey,
		SessionToken:    roleCreds.Token,
		ProviderName:    ProviderName,
	***REMOVED***, nil
***REMOVED***

// A ec2RoleCredRespBody provides the shape for unmarshaling credential
// request responses.
type ec2RoleCredRespBody struct ***REMOVED***
	// Success State
	Expiration      time.Time
	AccessKeyID     string
	SecretAccessKey string
	Token           string

	// Error state
	Code    string
	Message string
***REMOVED***

const iamSecurityCredsPath = "/iam/security-credentials"

// requestCredList requests a list of credentials from the EC2 service.
// If there are no credentials, or there is an error making or receiving the request
func requestCredList(client *ec2metadata.EC2Metadata) ([]string, error) ***REMOVED***
	resp, err := client.GetMetadata(iamSecurityCredsPath)
	if err != nil ***REMOVED***
		return nil, awserr.New("EC2RoleRequestError", "no EC2 instance role found", err)
	***REMOVED***

	credsList := []string***REMOVED******REMOVED***
	s := bufio.NewScanner(strings.NewReader(resp))
	for s.Scan() ***REMOVED***
		credsList = append(credsList, s.Text())
	***REMOVED***

	if err := s.Err(); err != nil ***REMOVED***
		return nil, awserr.New("SerializationError", "failed to read EC2 instance role from metadata service", err)
	***REMOVED***

	return credsList, nil
***REMOVED***

// requestCred requests the credentials for a specific credentials from the EC2 service.
//
// If the credentials cannot be found, or there is an error reading the response
// and error will be returned.
func requestCred(client *ec2metadata.EC2Metadata, credsName string) (ec2RoleCredRespBody, error) ***REMOVED***
	resp, err := client.GetMetadata(path.Join(iamSecurityCredsPath, credsName))
	if err != nil ***REMOVED***
		return ec2RoleCredRespBody***REMOVED******REMOVED***,
			awserr.New("EC2RoleRequestError",
				fmt.Sprintf("failed to get %s EC2 instance role credentials", credsName),
				err)
	***REMOVED***

	respCreds := ec2RoleCredRespBody***REMOVED******REMOVED***
	if err := json.NewDecoder(strings.NewReader(resp)).Decode(&respCreds); err != nil ***REMOVED***
		return ec2RoleCredRespBody***REMOVED******REMOVED***,
			awserr.New("SerializationError",
				fmt.Sprintf("failed to decode %s EC2 instance role credentials", credsName),
				err)
	***REMOVED***

	if respCreds.Code != "Success" ***REMOVED***
		// If an error code was returned something failed requesting the role.
		return ec2RoleCredRespBody***REMOVED******REMOVED***, awserr.New(respCreds.Code, respCreds.Message, nil)
	***REMOVED***

	return respCreds, nil
***REMOVED***
