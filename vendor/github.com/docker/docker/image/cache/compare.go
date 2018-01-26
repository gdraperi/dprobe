package cache

import (
	"github.com/docker/docker/api/types/container"
)

// compare two Config struct. Do not compare the "Image" nor "Hostname" fields
// If OpenStdin is set, then it differs
func compare(a, b *container.Config) bool ***REMOVED***
	if a == nil || b == nil ||
		a.OpenStdin || b.OpenStdin ***REMOVED***
		return false
	***REMOVED***
	if a.AttachStdout != b.AttachStdout ||
		a.AttachStderr != b.AttachStderr ||
		a.User != b.User ||
		a.OpenStdin != b.OpenStdin ||
		a.Tty != b.Tty ***REMOVED***
		return false
	***REMOVED***

	if len(a.Cmd) != len(b.Cmd) ||
		len(a.Env) != len(b.Env) ||
		len(a.Labels) != len(b.Labels) ||
		len(a.ExposedPorts) != len(b.ExposedPorts) ||
		len(a.Entrypoint) != len(b.Entrypoint) ||
		len(a.Volumes) != len(b.Volumes) ***REMOVED***
		return false
	***REMOVED***

	for i := 0; i < len(a.Cmd); i++ ***REMOVED***
		if a.Cmd[i] != b.Cmd[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	for i := 0; i < len(a.Env); i++ ***REMOVED***
		if a.Env[i] != b.Env[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	for k, v := range a.Labels ***REMOVED***
		if v != b.Labels[k] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	for k := range a.ExposedPorts ***REMOVED***
		if _, exists := b.ExposedPorts[k]; !exists ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	for i := 0; i < len(a.Entrypoint); i++ ***REMOVED***
		if a.Entrypoint[i] != b.Entrypoint[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	for key := range a.Volumes ***REMOVED***
		if _, exists := b.Volumes[key]; !exists ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
