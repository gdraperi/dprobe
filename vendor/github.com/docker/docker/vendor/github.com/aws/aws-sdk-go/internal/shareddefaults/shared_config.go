package shareddefaults

import (
	"os"
	"path/filepath"
	"runtime"
)

// SharedCredentialsFilename returns the SDK's default file path
// for the shared credentials file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.aws/credentials
//   - Windows: %USERPROFILE%\.aws\credentials
func SharedCredentialsFilename() string ***REMOVED***
	return filepath.Join(UserHomeDir(), ".aws", "credentials")
***REMOVED***

// SharedConfigFilename returns the SDK's default file path for
// the shared config file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.aws/config
//   - Windows: %USERPROFILE%\.aws\config
func SharedConfigFilename() string ***REMOVED***
	return filepath.Join(UserHomeDir(), ".aws", "config")
***REMOVED***

// UserHomeDir returns the home directory for the user the process is
// running under.
func UserHomeDir() string ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED*** // Windows
		return os.Getenv("USERPROFILE")
	***REMOVED***

	// *nix
	return os.Getenv("HOME")
***REMOVED***
