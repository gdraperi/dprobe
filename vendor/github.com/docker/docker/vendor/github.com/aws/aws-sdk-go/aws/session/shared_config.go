package session

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-ini/ini"
)

const (
	// Static Credentials group
	accessKeyIDKey  = `aws_access_key_id`     // group required
	secretAccessKey = `aws_secret_access_key` // group required
	sessionTokenKey = `aws_session_token`     // optional

	// Assume Role Credentials group
	roleArnKey         = `role_arn`          // group required
	sourceProfileKey   = `source_profile`    // group required
	externalIDKey      = `external_id`       // optional
	mfaSerialKey       = `mfa_serial`        // optional
	roleSessionNameKey = `role_session_name` // optional

	// Additional Config fields
	regionKey = `region`

	// DefaultSharedConfigProfile is the default profile to be used when
	// loading configuration from the config files if another profile name
	// is not provided.
	DefaultSharedConfigProfile = `default`
)

type assumeRoleConfig struct ***REMOVED***
	RoleARN         string
	SourceProfile   string
	ExternalID      string
	MFASerial       string
	RoleSessionName string
***REMOVED***

// sharedConfig represents the configuration fields of the SDK config files.
type sharedConfig struct ***REMOVED***
	// Credentials values from the config file. Both aws_access_key_id
	// and aws_secret_access_key must be provided together in the same file
	// to be considered valid. The values will be ignored if not a complete group.
	// aws_session_token is an optional field that can be provided if both of the
	// other two fields are also provided.
	//
	//	aws_access_key_id
	//	aws_secret_access_key
	//	aws_session_token
	Creds credentials.Value

	AssumeRole       assumeRoleConfig
	AssumeRoleSource *sharedConfig

	// Region is the region the SDK should use for looking up AWS service endpoints
	// and signing requests.
	//
	//	region
	Region string
***REMOVED***

type sharedConfigFile struct ***REMOVED***
	Filename string
	IniData  *ini.File
***REMOVED***

// loadSharedConfig retrieves the configuration from the list of files
// using the profile provided. The order the files are listed will determine
// precedence. Values in subsequent files will overwrite values defined in
// earlier files.
//
// For example, given two files A and B. Both define credentials. If the order
// of the files are A then B, B's credential values will be used instead of A's.
//
// See sharedConfig.setFromFile for information how the config files
// will be loaded.
func loadSharedConfig(profile string, filenames []string) (sharedConfig, error) ***REMOVED***
	if len(profile) == 0 ***REMOVED***
		profile = DefaultSharedConfigProfile
	***REMOVED***

	files, err := loadSharedConfigIniFiles(filenames)
	if err != nil ***REMOVED***
		return sharedConfig***REMOVED******REMOVED***, err
	***REMOVED***

	cfg := sharedConfig***REMOVED******REMOVED***
	if err = cfg.setFromIniFiles(profile, files); err != nil ***REMOVED***
		return sharedConfig***REMOVED******REMOVED***, err
	***REMOVED***

	if len(cfg.AssumeRole.SourceProfile) > 0 ***REMOVED***
		if err := cfg.setAssumeRoleSource(profile, files); err != nil ***REMOVED***
			return sharedConfig***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***

	return cfg, nil
***REMOVED***

func loadSharedConfigIniFiles(filenames []string) ([]sharedConfigFile, error) ***REMOVED***
	files := make([]sharedConfigFile, 0, len(filenames))

	for _, filename := range filenames ***REMOVED***
		b, err := ioutil.ReadFile(filename)
		if err != nil ***REMOVED***
			// Skip files which can't be opened and read for whatever reason
			continue
		***REMOVED***

		f, err := ini.Load(b)
		if err != nil ***REMOVED***
			return nil, SharedConfigLoadError***REMOVED***Filename: filename, Err: err***REMOVED***
		***REMOVED***

		files = append(files, sharedConfigFile***REMOVED***
			Filename: filename, IniData: f,
		***REMOVED***)
	***REMOVED***

	return files, nil
***REMOVED***

func (cfg *sharedConfig) setAssumeRoleSource(origProfile string, files []sharedConfigFile) error ***REMOVED***
	var assumeRoleSrc sharedConfig

	// Multiple level assume role chains are not support
	if cfg.AssumeRole.SourceProfile == origProfile ***REMOVED***
		assumeRoleSrc = *cfg
		assumeRoleSrc.AssumeRole = assumeRoleConfig***REMOVED******REMOVED***
	***REMOVED*** else ***REMOVED***
		err := assumeRoleSrc.setFromIniFiles(cfg.AssumeRole.SourceProfile, files)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if len(assumeRoleSrc.Creds.AccessKeyID) == 0 ***REMOVED***
		return SharedConfigAssumeRoleError***REMOVED***RoleARN: cfg.AssumeRole.RoleARN***REMOVED***
	***REMOVED***

	cfg.AssumeRoleSource = &assumeRoleSrc

	return nil
***REMOVED***

func (cfg *sharedConfig) setFromIniFiles(profile string, files []sharedConfigFile) error ***REMOVED***
	// Trim files from the list that don't exist.
	for _, f := range files ***REMOVED***
		if err := cfg.setFromIniFile(profile, f); err != nil ***REMOVED***
			if _, ok := err.(SharedConfigProfileNotExistsError); ok ***REMOVED***
				// Ignore proviles missings
				continue
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// setFromFile loads the configuration from the file using
// the profile provided. A sharedConfig pointer type value is used so that
// multiple config file loadings can be chained.
//
// Only loads complete logically grouped values, and will not set fields in cfg
// for incomplete grouped values in the config. Such as credentials. For example
// if a config file only includes aws_access_key_id but no aws_secret_access_key
// the aws_access_key_id will be ignored.
func (cfg *sharedConfig) setFromIniFile(profile string, file sharedConfigFile) error ***REMOVED***
	section, err := file.IniData.GetSection(profile)
	if err != nil ***REMOVED***
		// Fallback to to alternate profile name: profile <name>
		section, err = file.IniData.GetSection(fmt.Sprintf("profile %s", profile))
		if err != nil ***REMOVED***
			return SharedConfigProfileNotExistsError***REMOVED***Profile: profile, Err: err***REMOVED***
		***REMOVED***
	***REMOVED***

	// Shared Credentials
	akid := section.Key(accessKeyIDKey).String()
	secret := section.Key(secretAccessKey).String()
	if len(akid) > 0 && len(secret) > 0 ***REMOVED***
		cfg.Creds = credentials.Value***REMOVED***
			AccessKeyID:     akid,
			SecretAccessKey: secret,
			SessionToken:    section.Key(sessionTokenKey).String(),
			ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", file.Filename),
		***REMOVED***
	***REMOVED***

	// Assume Role
	roleArn := section.Key(roleArnKey).String()
	srcProfile := section.Key(sourceProfileKey).String()
	if len(roleArn) > 0 && len(srcProfile) > 0 ***REMOVED***
		cfg.AssumeRole = assumeRoleConfig***REMOVED***
			RoleARN:         roleArn,
			SourceProfile:   srcProfile,
			ExternalID:      section.Key(externalIDKey).String(),
			MFASerial:       section.Key(mfaSerialKey).String(),
			RoleSessionName: section.Key(roleSessionNameKey).String(),
		***REMOVED***
	***REMOVED***

	// Region
	if v := section.Key(regionKey).String(); len(v) > 0 ***REMOVED***
		cfg.Region = v
	***REMOVED***

	return nil
***REMOVED***

// SharedConfigLoadError is an error for the shared config file failed to load.
type SharedConfigLoadError struct ***REMOVED***
	Filename string
	Err      error
***REMOVED***

// Code is the short id of the error.
func (e SharedConfigLoadError) Code() string ***REMOVED***
	return "SharedConfigLoadError"
***REMOVED***

// Message is the description of the error
func (e SharedConfigLoadError) Message() string ***REMOVED***
	return fmt.Sprintf("failed to load config file, %s", e.Filename)
***REMOVED***

// OrigErr is the underlying error that caused the failure.
func (e SharedConfigLoadError) OrigErr() error ***REMOVED***
	return e.Err
***REMOVED***

// Error satisfies the error interface.
func (e SharedConfigLoadError) Error() string ***REMOVED***
	return awserr.SprintError(e.Code(), e.Message(), "", e.Err)
***REMOVED***

// SharedConfigProfileNotExistsError is an error for the shared config when
// the profile was not find in the config file.
type SharedConfigProfileNotExistsError struct ***REMOVED***
	Profile string
	Err     error
***REMOVED***

// Code is the short id of the error.
func (e SharedConfigProfileNotExistsError) Code() string ***REMOVED***
	return "SharedConfigProfileNotExistsError"
***REMOVED***

// Message is the description of the error
func (e SharedConfigProfileNotExistsError) Message() string ***REMOVED***
	return fmt.Sprintf("failed to get profile, %s", e.Profile)
***REMOVED***

// OrigErr is the underlying error that caused the failure.
func (e SharedConfigProfileNotExistsError) OrigErr() error ***REMOVED***
	return e.Err
***REMOVED***

// Error satisfies the error interface.
func (e SharedConfigProfileNotExistsError) Error() string ***REMOVED***
	return awserr.SprintError(e.Code(), e.Message(), "", e.Err)
***REMOVED***

// SharedConfigAssumeRoleError is an error for the shared config when the
// profile contains assume role information, but that information is invalid
// or not complete.
type SharedConfigAssumeRoleError struct ***REMOVED***
	RoleARN string
***REMOVED***

// Code is the short id of the error.
func (e SharedConfigAssumeRoleError) Code() string ***REMOVED***
	return "SharedConfigAssumeRoleError"
***REMOVED***

// Message is the description of the error
func (e SharedConfigAssumeRoleError) Message() string ***REMOVED***
	return fmt.Sprintf("failed to load assume role for %s, source profile has no shared credentials",
		e.RoleARN)
***REMOVED***

// OrigErr is the underlying error that caused the failure.
func (e SharedConfigAssumeRoleError) OrigErr() error ***REMOVED***
	return nil
***REMOVED***

// Error satisfies the error interface.
func (e SharedConfigAssumeRoleError) Error() string ***REMOVED***
	return awserr.SprintError(e.Code(), e.Message(), "", nil)
***REMOVED***
