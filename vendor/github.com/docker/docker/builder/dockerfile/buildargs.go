package dockerfile

import (
	"fmt"
	"io"

	"github.com/docker/docker/runconfig/opts"
)

// builtinAllowedBuildArgs is list of built-in allowed build args
// these args are considered transparent and are excluded from the image history.
// Filtering from history is implemented in dispatchers.go
var builtinAllowedBuildArgs = map[string]bool***REMOVED***
	"HTTP_PROXY":  true,
	"http_proxy":  true,
	"HTTPS_PROXY": true,
	"https_proxy": true,
	"FTP_PROXY":   true,
	"ftp_proxy":   true,
	"NO_PROXY":    true,
	"no_proxy":    true,
***REMOVED***

// buildArgs manages arguments used by the builder
type buildArgs struct ***REMOVED***
	// args that are allowed for expansion/substitution and passing to commands in 'run'.
	allowedBuildArgs map[string]*string
	// args defined before the first `FROM` in a Dockerfile
	allowedMetaArgs map[string]*string
	// args referenced by the Dockerfile
	referencedArgs map[string]struct***REMOVED******REMOVED***
	// args provided by the user on the command line
	argsFromOptions map[string]*string
***REMOVED***

func newBuildArgs(argsFromOptions map[string]*string) *buildArgs ***REMOVED***
	return &buildArgs***REMOVED***
		allowedBuildArgs: make(map[string]*string),
		allowedMetaArgs:  make(map[string]*string),
		referencedArgs:   make(map[string]struct***REMOVED******REMOVED***),
		argsFromOptions:  argsFromOptions,
	***REMOVED***
***REMOVED***

func (b *buildArgs) Clone() *buildArgs ***REMOVED***
	result := newBuildArgs(b.argsFromOptions)
	for k, v := range b.allowedBuildArgs ***REMOVED***
		result.allowedBuildArgs[k] = v
	***REMOVED***
	for k, v := range b.allowedMetaArgs ***REMOVED***
		result.allowedMetaArgs[k] = v
	***REMOVED***
	for k := range b.referencedArgs ***REMOVED***
		result.referencedArgs[k] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	return result
***REMOVED***

func (b *buildArgs) MergeReferencedArgs(other *buildArgs) ***REMOVED***
	for k := range other.referencedArgs ***REMOVED***
		b.referencedArgs[k] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// WarnOnUnusedBuildArgs checks if there are any leftover build-args that were
// passed but not consumed during build. Print a warning, if there are any.
func (b *buildArgs) WarnOnUnusedBuildArgs(out io.Writer) ***REMOVED***
	leftoverArgs := []string***REMOVED******REMOVED***
	for arg := range b.argsFromOptions ***REMOVED***
		_, isReferenced := b.referencedArgs[arg]
		_, isBuiltin := builtinAllowedBuildArgs[arg]
		if !isBuiltin && !isReferenced ***REMOVED***
			leftoverArgs = append(leftoverArgs, arg)
		***REMOVED***
	***REMOVED***
	if len(leftoverArgs) > 0 ***REMOVED***
		fmt.Fprintf(out, "[Warning] One or more build-args %v were not consumed\n", leftoverArgs)
	***REMOVED***
***REMOVED***

// ResetAllowed clears the list of args that are allowed to be used by a
// directive
func (b *buildArgs) ResetAllowed() ***REMOVED***
	b.allowedBuildArgs = make(map[string]*string)
***REMOVED***

// AddMetaArg adds a new meta arg that can be used by FROM directives
func (b *buildArgs) AddMetaArg(key string, value *string) ***REMOVED***
	b.allowedMetaArgs[key] = value
***REMOVED***

// AddArg adds a new arg that can be used by directives
func (b *buildArgs) AddArg(key string, value *string) ***REMOVED***
	b.allowedBuildArgs[key] = value
	b.referencedArgs[key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***

// IsReferencedOrNotBuiltin checks if the key is a built-in arg, or if it has been
// referenced by the Dockerfile. Returns true if the arg is not a builtin or
// if the builtin has been referenced in the Dockerfile.
func (b *buildArgs) IsReferencedOrNotBuiltin(key string) bool ***REMOVED***
	_, isBuiltin := builtinAllowedBuildArgs[key]
	_, isAllowed := b.allowedBuildArgs[key]
	return isAllowed || !isBuiltin
***REMOVED***

// GetAllAllowed returns a mapping with all the allowed args
func (b *buildArgs) GetAllAllowed() map[string]string ***REMOVED***
	return b.getAllFromMapping(b.allowedBuildArgs)
***REMOVED***

// GetAllMeta returns a mapping with all the meta meta args
func (b *buildArgs) GetAllMeta() map[string]string ***REMOVED***
	return b.getAllFromMapping(b.allowedMetaArgs)
***REMOVED***

func (b *buildArgs) getAllFromMapping(source map[string]*string) map[string]string ***REMOVED***
	m := make(map[string]string)

	keys := keysFromMaps(source, builtinAllowedBuildArgs)
	for _, key := range keys ***REMOVED***
		v, ok := b.getBuildArg(key, source)
		if ok ***REMOVED***
			m[key] = v
		***REMOVED***
	***REMOVED***
	return m
***REMOVED***

// FilterAllowed returns all allowed args without the filtered args
func (b *buildArgs) FilterAllowed(filter []string) []string ***REMOVED***
	envs := []string***REMOVED******REMOVED***
	configEnv := opts.ConvertKVStringsToMap(filter)

	for key, val := range b.GetAllAllowed() ***REMOVED***
		if _, ok := configEnv[key]; !ok ***REMOVED***
			envs = append(envs, fmt.Sprintf("%s=%s", key, val))
		***REMOVED***
	***REMOVED***
	return envs
***REMOVED***

func (b *buildArgs) getBuildArg(key string, mapping map[string]*string) (string, bool) ***REMOVED***
	defaultValue, exists := mapping[key]
	// Return override from options if one is defined
	if v, ok := b.argsFromOptions[key]; ok && v != nil ***REMOVED***
		return *v, ok
	***REMOVED***

	if defaultValue == nil ***REMOVED***
		if v, ok := b.allowedMetaArgs[key]; ok && v != nil ***REMOVED***
			return *v, ok
		***REMOVED***
		return "", false
	***REMOVED***
	return *defaultValue, exists
***REMOVED***

func keysFromMaps(source map[string]*string, builtin map[string]bool) []string ***REMOVED***
	keys := []string***REMOVED******REMOVED***
	for key := range source ***REMOVED***
		keys = append(keys, key)
	***REMOVED***
	for key := range builtin ***REMOVED***
		keys = append(keys, key)
	***REMOVED***
	return keys
***REMOVED***
