package swarm

import "os"

// Config represents a config.
type Config struct ***REMOVED***
	ID string
	Meta
	Spec ConfigSpec
***REMOVED***

// ConfigSpec represents a config specification from a config in swarm
type ConfigSpec struct ***REMOVED***
	Annotations
	Data []byte `json:",omitempty"`
***REMOVED***

// ConfigReferenceFileTarget is a file target in a config reference
type ConfigReferenceFileTarget struct ***REMOVED***
	Name string
	UID  string
	GID  string
	Mode os.FileMode
***REMOVED***

// ConfigReference is a reference to a config in swarm
type ConfigReference struct ***REMOVED***
	File       *ConfigReferenceFileTarget
	ConfigID   string
	ConfigName string
***REMOVED***
