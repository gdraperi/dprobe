package logger

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Info provides enough information for a logging driver to do its function.
type Info struct ***REMOVED***
	Config              map[string]string
	ContainerID         string
	ContainerName       string
	ContainerEntrypoint string
	ContainerArgs       []string
	ContainerImageID    string
	ContainerImageName  string
	ContainerCreated    time.Time
	ContainerEnv        []string
	ContainerLabels     map[string]string
	LogPath             string
	DaemonName          string
***REMOVED***

// ExtraAttributes returns the user-defined extra attributes (labels,
// environment variables) in key-value format. This can be used by log drivers
// that support metadata to add more context to a log.
func (info *Info) ExtraAttributes(keyMod func(string) string) (map[string]string, error) ***REMOVED***
	extra := make(map[string]string)
	labels, ok := info.Config["labels"]
	if ok && len(labels) > 0 ***REMOVED***
		for _, l := range strings.Split(labels, ",") ***REMOVED***
			if v, ok := info.ContainerLabels[l]; ok ***REMOVED***
				if keyMod != nil ***REMOVED***
					l = keyMod(l)
				***REMOVED***
				extra[l] = v
			***REMOVED***
		***REMOVED***
	***REMOVED***

	envMapping := make(map[string]string)
	for _, e := range info.ContainerEnv ***REMOVED***
		if kv := strings.SplitN(e, "=", 2); len(kv) == 2 ***REMOVED***
			envMapping[kv[0]] = kv[1]
		***REMOVED***
	***REMOVED***

	env, ok := info.Config["env"]
	if ok && len(env) > 0 ***REMOVED***
		for _, l := range strings.Split(env, ",") ***REMOVED***
			if v, ok := envMapping[l]; ok ***REMOVED***
				if keyMod != nil ***REMOVED***
					l = keyMod(l)
				***REMOVED***
				extra[l] = v
			***REMOVED***
		***REMOVED***
	***REMOVED***

	envRegex, ok := info.Config["env-regex"]
	if ok && len(envRegex) > 0 ***REMOVED***
		re, err := regexp.Compile(envRegex)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for k, v := range envMapping ***REMOVED***
			if re.MatchString(k) ***REMOVED***
				if keyMod != nil ***REMOVED***
					k = keyMod(k)
				***REMOVED***
				extra[k] = v
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return extra, nil
***REMOVED***

// Hostname returns the hostname from the underlying OS.
func (info *Info) Hostname() (string, error) ***REMOVED***
	hostname, err := os.Hostname()
	if err != nil ***REMOVED***
		return "", fmt.Errorf("logger: can not resolve hostname: %v", err)
	***REMOVED***
	return hostname, nil
***REMOVED***

// Command returns the command that the container being logged was
// started with. The Entrypoint is prepended to the container
// arguments.
func (info *Info) Command() string ***REMOVED***
	terms := []string***REMOVED***info.ContainerEntrypoint***REMOVED***
	terms = append(terms, info.ContainerArgs...)
	command := strings.Join(terms, " ")
	return command
***REMOVED***

// ID Returns the Container ID shortened to 12 characters.
func (info *Info) ID() string ***REMOVED***
	return info.ContainerID[:12]
***REMOVED***

// FullID is an alias of ContainerID.
func (info *Info) FullID() string ***REMOVED***
	return info.ContainerID
***REMOVED***

// Name returns the ContainerName without a preceding '/'.
func (info *Info) Name() string ***REMOVED***
	return strings.TrimPrefix(info.ContainerName, "/")
***REMOVED***

// ImageID returns the ContainerImageID shortened to 12 characters.
func (info *Info) ImageID() string ***REMOVED***
	return info.ContainerImageID[:12]
***REMOVED***

// ImageFullID is an alias of ContainerImageID.
func (info *Info) ImageFullID() string ***REMOVED***
	return info.ContainerImageID
***REMOVED***

// ImageName is an alias of ContainerImageName
func (info *Info) ImageName() string ***REMOVED***
	return info.ContainerImageName
***REMOVED***
