package validation

import "fmt"

// MaxSecretSize is the maximum byte length of the `Secret.Spec.Data` field.
const MaxSecretSize = 500 * 1024 // 500KB

// ValidateSecretPayload validates the secret payload size
func ValidateSecretPayload(data []byte) error ***REMOVED***
	if len(data) >= MaxSecretSize || len(data) < 1 ***REMOVED***
		return fmt.Errorf("secret data must be larger than 0 and less than %d bytes", MaxSecretSize)
	***REMOVED***
	return nil
***REMOVED***
