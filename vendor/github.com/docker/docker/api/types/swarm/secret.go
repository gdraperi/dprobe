package swarm

import "os"

// Secret represents a secret.
type Secret struct ***REMOVED***
	ID string
	Meta
	Spec SecretSpec
***REMOVED***

// SecretSpec represents a secret specification from a secret in swarm
type SecretSpec struct ***REMOVED***
	Annotations
	Data   []byte  `json:",omitempty"`
	Driver *Driver `json:",omitempty"` // name of the secrets driver used to fetch the secret's value from an external secret store
***REMOVED***

// SecretReferenceFileTarget is a file target in a secret reference
type SecretReferenceFileTarget struct ***REMOVED***
	Name string
	UID  string
	GID  string
	Mode os.FileMode
***REMOVED***

// SecretReference is a reference to a secret in swarm
type SecretReference struct ***REMOVED***
	File       *SecretReferenceFileTarget
	SecretID   string
	SecretName string
***REMOVED***
