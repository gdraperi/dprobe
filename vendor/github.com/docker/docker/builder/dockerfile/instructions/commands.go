package instructions

import (
	"errors"

	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
)

// KeyValuePair represent an arbitrary named value (useful in slice insted of map[string] string to preserve ordering)
type KeyValuePair struct ***REMOVED***
	Key   string
	Value string
***REMOVED***

func (kvp *KeyValuePair) String() string ***REMOVED***
	return kvp.Key + "=" + kvp.Value
***REMOVED***

// Command is implemented by every command present in a dockerfile
type Command interface ***REMOVED***
	Name() string
***REMOVED***

// KeyValuePairs is a slice of KeyValuePair
type KeyValuePairs []KeyValuePair

// withNameAndCode is the base of every command in a Dockerfile (String() returns its source code)
type withNameAndCode struct ***REMOVED***
	code string
	name string
***REMOVED***

func (c *withNameAndCode) String() string ***REMOVED***
	return c.code
***REMOVED***

// Name of the command
func (c *withNameAndCode) Name() string ***REMOVED***
	return c.name
***REMOVED***

func newWithNameAndCode(req parseRequest) withNameAndCode ***REMOVED***
	return withNameAndCode***REMOVED***code: strings.TrimSpace(req.original), name: req.command***REMOVED***
***REMOVED***

// SingleWordExpander is a provider for variable expansion where 1 word => 1 output
type SingleWordExpander func(word string) (string, error)

// SupportsSingleWordExpansion interface marks a command as supporting variable expansion
type SupportsSingleWordExpansion interface ***REMOVED***
	Expand(expander SingleWordExpander) error
***REMOVED***

// PlatformSpecific adds platform checks to a command
type PlatformSpecific interface ***REMOVED***
	CheckPlatform(platform string) error
***REMOVED***

func expandKvp(kvp KeyValuePair, expander SingleWordExpander) (KeyValuePair, error) ***REMOVED***
	key, err := expander(kvp.Key)
	if err != nil ***REMOVED***
		return KeyValuePair***REMOVED******REMOVED***, err
	***REMOVED***
	value, err := expander(kvp.Value)
	if err != nil ***REMOVED***
		return KeyValuePair***REMOVED******REMOVED***, err
	***REMOVED***
	return KeyValuePair***REMOVED***Key: key, Value: value***REMOVED***, nil
***REMOVED***
func expandKvpsInPlace(kvps KeyValuePairs, expander SingleWordExpander) error ***REMOVED***
	for i, kvp := range kvps ***REMOVED***
		newKvp, err := expandKvp(kvp, expander)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		kvps[i] = newKvp
	***REMOVED***
	return nil
***REMOVED***

func expandSliceInPlace(values []string, expander SingleWordExpander) error ***REMOVED***
	for i, v := range values ***REMOVED***
		newValue, err := expander(v)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		values[i] = newValue
	***REMOVED***
	return nil
***REMOVED***

// EnvCommand : ENV key1 value1 [keyN valueN...]
type EnvCommand struct ***REMOVED***
	withNameAndCode
	Env KeyValuePairs // kvp slice instead of map to preserve ordering
***REMOVED***

// Expand variables
func (c *EnvCommand) Expand(expander SingleWordExpander) error ***REMOVED***
	return expandKvpsInPlace(c.Env, expander)
***REMOVED***

// MaintainerCommand : MAINTAINER maintainer_name
type MaintainerCommand struct ***REMOVED***
	withNameAndCode
	Maintainer string
***REMOVED***

// LabelCommand : LABEL some json data describing the image
//
// Sets the Label variable foo to bar,
//
type LabelCommand struct ***REMOVED***
	withNameAndCode
	Labels KeyValuePairs // kvp slice instead of map to preserve ordering
***REMOVED***

// Expand variables
func (c *LabelCommand) Expand(expander SingleWordExpander) error ***REMOVED***
	return expandKvpsInPlace(c.Labels, expander)
***REMOVED***

// SourcesAndDest represent a list of source files and a destination
type SourcesAndDest []string

// Sources list the source paths
func (s SourcesAndDest) Sources() []string ***REMOVED***
	res := make([]string, len(s)-1)
	copy(res, s[:len(s)-1])
	return res
***REMOVED***

// Dest path of the operation
func (s SourcesAndDest) Dest() string ***REMOVED***
	return s[len(s)-1]
***REMOVED***

// AddCommand : ADD foo /path
//
// Add the file 'foo' to '/path'. Tarball and Remote URL (git, http) handling
// exist here. If you do not wish to have this automatic handling, use COPY.
//
type AddCommand struct ***REMOVED***
	withNameAndCode
	SourcesAndDest
	Chown string
***REMOVED***

// Expand variables
func (c *AddCommand) Expand(expander SingleWordExpander) error ***REMOVED***
	return expandSliceInPlace(c.SourcesAndDest, expander)
***REMOVED***

// CopyCommand : COPY foo /path
//
// Same as 'ADD' but without the tar and remote url handling.
//
type CopyCommand struct ***REMOVED***
	withNameAndCode
	SourcesAndDest
	From  string
	Chown string
***REMOVED***

// Expand variables
func (c *CopyCommand) Expand(expander SingleWordExpander) error ***REMOVED***
	return expandSliceInPlace(c.SourcesAndDest, expander)
***REMOVED***

// OnbuildCommand : ONBUILD <some other command>
type OnbuildCommand struct ***REMOVED***
	withNameAndCode
	Expression string
***REMOVED***

// WorkdirCommand : WORKDIR /tmp
//
// Set the working directory for future RUN/CMD/etc statements.
//
type WorkdirCommand struct ***REMOVED***
	withNameAndCode
	Path string
***REMOVED***

// Expand variables
func (c *WorkdirCommand) Expand(expander SingleWordExpander) error ***REMOVED***
	p, err := expander(c.Path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.Path = p
	return nil
***REMOVED***

// ShellDependantCmdLine represents a cmdline optionaly prepended with the shell
type ShellDependantCmdLine struct ***REMOVED***
	CmdLine      strslice.StrSlice
	PrependShell bool
***REMOVED***

// RunCommand : RUN some command yo
//
// run a command and commit the image. Args are automatically prepended with
// the current SHELL which defaults to 'sh -c' under linux or 'cmd /S /C' under
// Windows, in the event there is only one argument The difference in processing:
//
// RUN echo hi          # sh -c echo hi       (Linux)
// RUN echo hi          # cmd /S /C echo hi   (Windows)
// RUN [ "echo", "hi" ] # echo hi
//
type RunCommand struct ***REMOVED***
	withNameAndCode
	ShellDependantCmdLine
***REMOVED***

// CmdCommand : CMD foo
//
// Set the default command to run in the container (which may be empty).
// Argument handling is the same as RUN.
//
type CmdCommand struct ***REMOVED***
	withNameAndCode
	ShellDependantCmdLine
***REMOVED***

// HealthCheckCommand : HEALTHCHECK foo
//
// Set the default healthcheck command to run in the container (which may be empty).
// Argument handling is the same as RUN.
//
type HealthCheckCommand struct ***REMOVED***
	withNameAndCode
	Health *container.HealthConfig
***REMOVED***

// EntrypointCommand : ENTRYPOINT /usr/sbin/nginx
//
// Set the entrypoint to /usr/sbin/nginx. Will accept the CMD as the arguments
// to /usr/sbin/nginx. Uses the default shell if not in JSON format.
//
// Handles command processing similar to CMD and RUN, only req.runConfig.Entrypoint
// is initialized at newBuilder time instead of through argument parsing.
//
type EntrypointCommand struct ***REMOVED***
	withNameAndCode
	ShellDependantCmdLine
***REMOVED***

// ExposeCommand : EXPOSE 6667/tcp 7000/tcp
//
// Expose ports for links and port mappings. This all ends up in
// req.runConfig.ExposedPorts for runconfig.
//
type ExposeCommand struct ***REMOVED***
	withNameAndCode
	Ports []string
***REMOVED***

// UserCommand : USER foo
//
// Set the user to 'foo' for future commands and when running the
// ENTRYPOINT/CMD at container run time.
//
type UserCommand struct ***REMOVED***
	withNameAndCode
	User string
***REMOVED***

// Expand variables
func (c *UserCommand) Expand(expander SingleWordExpander) error ***REMOVED***
	p, err := expander(c.User)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.User = p
	return nil
***REMOVED***

// VolumeCommand : VOLUME /foo
//
// Expose the volume /foo for use. Will also accept the JSON array form.
//
type VolumeCommand struct ***REMOVED***
	withNameAndCode
	Volumes []string
***REMOVED***

// Expand variables
func (c *VolumeCommand) Expand(expander SingleWordExpander) error ***REMOVED***
	return expandSliceInPlace(c.Volumes, expander)
***REMOVED***

// StopSignalCommand : STOPSIGNAL signal
//
// Set the signal that will be used to kill the container.
type StopSignalCommand struct ***REMOVED***
	withNameAndCode
	Signal string
***REMOVED***

// Expand variables
func (c *StopSignalCommand) Expand(expander SingleWordExpander) error ***REMOVED***
	p, err := expander(c.Signal)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.Signal = p
	return nil
***REMOVED***

// CheckPlatform checks that the command is supported in the target platform
func (c *StopSignalCommand) CheckPlatform(platform string) error ***REMOVED***
	if platform == "windows" ***REMOVED***
		return errors.New("The daemon on this platform does not support the command stopsignal")
	***REMOVED***
	return nil
***REMOVED***

// ArgCommand : ARG name[=value]
//
// Adds the variable foo to the trusted list of variables that can be passed
// to builder using the --build-arg flag for expansion/substitution or passing to 'run'.
// Dockerfile author may optionally set a default value of this variable.
type ArgCommand struct ***REMOVED***
	withNameAndCode
	Key   string
	Value *string
***REMOVED***

// Expand variables
func (c *ArgCommand) Expand(expander SingleWordExpander) error ***REMOVED***
	p, err := expander(c.Key)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.Key = p
	if c.Value != nil ***REMOVED***
		p, err = expander(*c.Value)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.Value = &p
	***REMOVED***
	return nil
***REMOVED***

// ShellCommand : SHELL powershell -command
//
// Set the non-default shell to use.
type ShellCommand struct ***REMOVED***
	withNameAndCode
	Shell strslice.StrSlice
***REMOVED***

// Stage represents a single stage in a multi-stage build
type Stage struct ***REMOVED***
	Name       string
	Commands   []Command
	BaseName   string
	SourceCode string
***REMOVED***

// AddCommand to the stage
func (s *Stage) AddCommand(cmd Command) ***REMOVED***
	// todo: validate cmd type
	s.Commands = append(s.Commands, cmd)
***REMOVED***

// IsCurrentStage check if the stage name is the current stage
func IsCurrentStage(s []Stage, name string) bool ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return false
	***REMOVED***
	return s[len(s)-1].Name == name
***REMOVED***

// CurrentStage return the last stage in a slice
func CurrentStage(s []Stage) (*Stage, error) ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return nil, errors.New("No build stage in current context")
	***REMOVED***
	return &s[len(s)-1], nil
***REMOVED***

// HasStage looks for the presence of a given stage name
func HasStage(s []Stage, name string) (int, bool) ***REMOVED***
	for i, stage := range s ***REMOVED***
		if stage.Name == name ***REMOVED***
			return i, true
		***REMOVED***
	***REMOVED***
	return -1, false
***REMOVED***
