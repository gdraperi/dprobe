package instructions

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/builder/dockerfile/command"
	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/pkg/errors"
)

type parseRequest struct ***REMOVED***
	command    string
	args       []string
	attributes map[string]bool
	flags      *BFlags
	original   string
***REMOVED***

func nodeArgs(node *parser.Node) []string ***REMOVED***
	result := []string***REMOVED******REMOVED***
	for ; node.Next != nil; node = node.Next ***REMOVED***
		arg := node.Next
		if len(arg.Children) == 0 ***REMOVED***
			result = append(result, arg.Value)
		***REMOVED*** else if len(arg.Children) == 1 ***REMOVED***
			//sub command
			result = append(result, arg.Children[0].Value)
			result = append(result, nodeArgs(arg.Children[0])...)
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

func newParseRequestFromNode(node *parser.Node) parseRequest ***REMOVED***
	return parseRequest***REMOVED***
		command:    node.Value,
		args:       nodeArgs(node),
		attributes: node.Attributes,
		original:   node.Original,
		flags:      NewBFlagsWithArgs(node.Flags),
	***REMOVED***
***REMOVED***

// ParseInstruction converts an AST to a typed instruction (either a command or a build stage beginning when encountering a `FROM` statement)
func ParseInstruction(node *parser.Node) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	req := newParseRequestFromNode(node)
	switch node.Value ***REMOVED***
	case command.Env:
		return parseEnv(req)
	case command.Maintainer:
		return parseMaintainer(req)
	case command.Label:
		return parseLabel(req)
	case command.Add:
		return parseAdd(req)
	case command.Copy:
		return parseCopy(req)
	case command.From:
		return parseFrom(req)
	case command.Onbuild:
		return parseOnBuild(req)
	case command.Workdir:
		return parseWorkdir(req)
	case command.Run:
		return parseRun(req)
	case command.Cmd:
		return parseCmd(req)
	case command.Healthcheck:
		return parseHealthcheck(req)
	case command.Entrypoint:
		return parseEntrypoint(req)
	case command.Expose:
		return parseExpose(req)
	case command.User:
		return parseUser(req)
	case command.Volume:
		return parseVolume(req)
	case command.StopSignal:
		return parseStopSignal(req)
	case command.Arg:
		return parseArg(req)
	case command.Shell:
		return parseShell(req)
	***REMOVED***

	return nil, &UnknownInstruction***REMOVED***Instruction: node.Value, Line: node.StartLine***REMOVED***
***REMOVED***

// ParseCommand converts an AST to a typed Command
func ParseCommand(node *parser.Node) (Command, error) ***REMOVED***
	s, err := ParseInstruction(node)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if c, ok := s.(Command); ok ***REMOVED***
		return c, nil
	***REMOVED***
	return nil, errors.Errorf("%T is not a command type", s)
***REMOVED***

// UnknownInstruction represents an error occurring when a command is unresolvable
type UnknownInstruction struct ***REMOVED***
	Line        int
	Instruction string
***REMOVED***

func (e *UnknownInstruction) Error() string ***REMOVED***
	return fmt.Sprintf("unknown instruction: %s", strings.ToUpper(e.Instruction))
***REMOVED***

// IsUnknownInstruction checks if the error is an UnknownInstruction or a parseError containing an UnknownInstruction
func IsUnknownInstruction(err error) bool ***REMOVED***
	_, ok := err.(*UnknownInstruction)
	if !ok ***REMOVED***
		var pe *parseError
		if pe, ok = err.(*parseError); ok ***REMOVED***
			_, ok = pe.inner.(*UnknownInstruction)
		***REMOVED***
	***REMOVED***
	return ok
***REMOVED***

type parseError struct ***REMOVED***
	inner error
	node  *parser.Node
***REMOVED***

func (e *parseError) Error() string ***REMOVED***
	return fmt.Sprintf("Dockerfile parse error line %d: %v", e.node.StartLine, e.inner.Error())
***REMOVED***

// Parse a docker file into a collection of buildable stages
func Parse(ast *parser.Node) (stages []Stage, metaArgs []ArgCommand, err error) ***REMOVED***
	for _, n := range ast.Children ***REMOVED***
		cmd, err := ParseInstruction(n)
		if err != nil ***REMOVED***
			return nil, nil, &parseError***REMOVED***inner: err, node: n***REMOVED***
		***REMOVED***
		if len(stages) == 0 ***REMOVED***
			// meta arg case
			if a, isArg := cmd.(*ArgCommand); isArg ***REMOVED***
				metaArgs = append(metaArgs, *a)
				continue
			***REMOVED***
		***REMOVED***
		switch c := cmd.(type) ***REMOVED***
		case *Stage:
			stages = append(stages, *c)
		case Command:
			stage, err := CurrentStage(stages)
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
			stage.AddCommand(c)
		default:
			return nil, nil, errors.Errorf("%T is not a command type", cmd)
		***REMOVED***

	***REMOVED***
	return stages, metaArgs, nil
***REMOVED***

func parseKvps(args []string, cmdName string) (KeyValuePairs, error) ***REMOVED***
	if len(args) == 0 ***REMOVED***
		return nil, errAtLeastOneArgument(cmdName)
	***REMOVED***
	if len(args)%2 != 0 ***REMOVED***
		// should never get here, but just in case
		return nil, errTooManyArguments(cmdName)
	***REMOVED***
	var res KeyValuePairs
	for j := 0; j < len(args); j += 2 ***REMOVED***
		if len(args[j]) == 0 ***REMOVED***
			return nil, errBlankCommandNames(cmdName)
		***REMOVED***
		name := args[j]
		value := args[j+1]
		res = append(res, KeyValuePair***REMOVED***Key: name, Value: value***REMOVED***)
	***REMOVED***
	return res, nil
***REMOVED***

func parseEnv(req parseRequest) (*EnvCommand, error) ***REMOVED***

	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	envs, err := parseKvps(req.args, "ENV")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &EnvCommand***REMOVED***
		Env:             envs,
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***, nil
***REMOVED***

func parseMaintainer(req parseRequest) (*MaintainerCommand, error) ***REMOVED***
	if len(req.args) != 1 ***REMOVED***
		return nil, errExactlyOneArgument("MAINTAINER")
	***REMOVED***

	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &MaintainerCommand***REMOVED***
		Maintainer:      req.args[0],
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***, nil
***REMOVED***

func parseLabel(req parseRequest) (*LabelCommand, error) ***REMOVED***

	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	labels, err := parseKvps(req.args, "LABEL")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &LabelCommand***REMOVED***
		Labels:          labels,
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***, nil
***REMOVED***

func parseAdd(req parseRequest) (*AddCommand, error) ***REMOVED***
	if len(req.args) < 2 ***REMOVED***
		return nil, errNoDestinationArgument("ADD")
	***REMOVED***
	flChown := req.flags.AddString("chown", "")
	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &AddCommand***REMOVED***
		SourcesAndDest:  SourcesAndDest(req.args),
		withNameAndCode: newWithNameAndCode(req),
		Chown:           flChown.Value,
	***REMOVED***, nil
***REMOVED***

func parseCopy(req parseRequest) (*CopyCommand, error) ***REMOVED***
	if len(req.args) < 2 ***REMOVED***
		return nil, errNoDestinationArgument("COPY")
	***REMOVED***
	flChown := req.flags.AddString("chown", "")
	flFrom := req.flags.AddString("from", "")
	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &CopyCommand***REMOVED***
		SourcesAndDest:  SourcesAndDest(req.args),
		From:            flFrom.Value,
		withNameAndCode: newWithNameAndCode(req),
		Chown:           flChown.Value,
	***REMOVED***, nil
***REMOVED***

func parseFrom(req parseRequest) (*Stage, error) ***REMOVED***
	stageName, err := parseBuildStageName(req.args)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	code := strings.TrimSpace(req.original)

	return &Stage***REMOVED***
		BaseName:   req.args[0],
		Name:       stageName,
		SourceCode: code,
		Commands:   []Command***REMOVED******REMOVED***,
	***REMOVED***, nil

***REMOVED***

func parseBuildStageName(args []string) (string, error) ***REMOVED***
	stageName := ""
	switch ***REMOVED***
	case len(args) == 3 && strings.EqualFold(args[1], "as"):
		stageName = strings.ToLower(args[2])
		if ok, _ := regexp.MatchString("^[a-z][a-z0-9-_\\.]*$", stageName); !ok ***REMOVED***
			return "", errors.Errorf("invalid name for build stage: %q, name can't start with a number or contain symbols", stageName)
		***REMOVED***
	case len(args) != 1:
		return "", errors.New("FROM requires either one or three arguments")
	***REMOVED***

	return stageName, nil
***REMOVED***

func parseOnBuild(req parseRequest) (*OnbuildCommand, error) ***REMOVED***
	if len(req.args) == 0 ***REMOVED***
		return nil, errAtLeastOneArgument("ONBUILD")
	***REMOVED***
	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	triggerInstruction := strings.ToUpper(strings.TrimSpace(req.args[0]))
	switch strings.ToUpper(triggerInstruction) ***REMOVED***
	case "ONBUILD":
		return nil, errors.New("Chaining ONBUILD via `ONBUILD ONBUILD` isn't allowed")
	case "MAINTAINER", "FROM":
		return nil, fmt.Errorf("%s isn't allowed as an ONBUILD trigger", triggerInstruction)
	***REMOVED***

	original := regexp.MustCompile(`(?i)^\s*ONBUILD\s*`).ReplaceAllString(req.original, "")
	return &OnbuildCommand***REMOVED***
		Expression:      original,
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***, nil

***REMOVED***

func parseWorkdir(req parseRequest) (*WorkdirCommand, error) ***REMOVED***
	if len(req.args) != 1 ***REMOVED***
		return nil, errExactlyOneArgument("WORKDIR")
	***REMOVED***

	err := req.flags.Parse()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &WorkdirCommand***REMOVED***
		Path:            req.args[0],
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***, nil

***REMOVED***

func parseShellDependentCommand(req parseRequest, emptyAsNil bool) ShellDependantCmdLine ***REMOVED***
	args := handleJSONArgs(req.args, req.attributes)
	cmd := strslice.StrSlice(args)
	if emptyAsNil && len(cmd) == 0 ***REMOVED***
		cmd = nil
	***REMOVED***
	return ShellDependantCmdLine***REMOVED***
		CmdLine:      cmd,
		PrependShell: !req.attributes["json"],
	***REMOVED***
***REMOVED***

func parseRun(req parseRequest) (*RunCommand, error) ***REMOVED***

	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &RunCommand***REMOVED***
		ShellDependantCmdLine: parseShellDependentCommand(req, false),
		withNameAndCode:       newWithNameAndCode(req),
	***REMOVED***, nil

***REMOVED***

func parseCmd(req parseRequest) (*CmdCommand, error) ***REMOVED***
	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &CmdCommand***REMOVED***
		ShellDependantCmdLine: parseShellDependentCommand(req, false),
		withNameAndCode:       newWithNameAndCode(req),
	***REMOVED***, nil

***REMOVED***

func parseEntrypoint(req parseRequest) (*EntrypointCommand, error) ***REMOVED***
	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cmd := &EntrypointCommand***REMOVED***
		ShellDependantCmdLine: parseShellDependentCommand(req, true),
		withNameAndCode:       newWithNameAndCode(req),
	***REMOVED***

	return cmd, nil
***REMOVED***

// parseOptInterval(flag) is the duration of flag.Value, or 0 if
// empty. An error is reported if the value is given and less than minimum duration.
func parseOptInterval(f *Flag) (time.Duration, error) ***REMOVED***
	s := f.Value
	if s == "" ***REMOVED***
		return 0, nil
	***REMOVED***
	d, err := time.ParseDuration(s)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if d < container.MinimumDuration ***REMOVED***
		return 0, fmt.Errorf("Interval %#v cannot be less than %s", f.name, container.MinimumDuration)
	***REMOVED***
	return d, nil
***REMOVED***
func parseHealthcheck(req parseRequest) (*HealthCheckCommand, error) ***REMOVED***
	if len(req.args) == 0 ***REMOVED***
		return nil, errAtLeastOneArgument("HEALTHCHECK")
	***REMOVED***
	cmd := &HealthCheckCommand***REMOVED***
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***

	typ := strings.ToUpper(req.args[0])
	args := req.args[1:]
	if typ == "NONE" ***REMOVED***
		if len(args) != 0 ***REMOVED***
			return nil, errors.New("HEALTHCHECK NONE takes no arguments")
		***REMOVED***
		test := strslice.StrSlice***REMOVED***typ***REMOVED***
		cmd.Health = &container.HealthConfig***REMOVED***
			Test: test,
		***REMOVED***
	***REMOVED*** else ***REMOVED***

		healthcheck := container.HealthConfig***REMOVED******REMOVED***

		flInterval := req.flags.AddString("interval", "")
		flTimeout := req.flags.AddString("timeout", "")
		flStartPeriod := req.flags.AddString("start-period", "")
		flRetries := req.flags.AddString("retries", "")

		if err := req.flags.Parse(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch typ ***REMOVED***
		case "CMD":
			cmdSlice := handleJSONArgs(args, req.attributes)
			if len(cmdSlice) == 0 ***REMOVED***
				return nil, errors.New("Missing command after HEALTHCHECK CMD")
			***REMOVED***

			if !req.attributes["json"] ***REMOVED***
				typ = "CMD-SHELL"
			***REMOVED***

			healthcheck.Test = strslice.StrSlice(append([]string***REMOVED***typ***REMOVED***, cmdSlice...))
		default:
			return nil, fmt.Errorf("Unknown type %#v in HEALTHCHECK (try CMD)", typ)
		***REMOVED***

		interval, err := parseOptInterval(flInterval)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		healthcheck.Interval = interval

		timeout, err := parseOptInterval(flTimeout)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		healthcheck.Timeout = timeout

		startPeriod, err := parseOptInterval(flStartPeriod)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		healthcheck.StartPeriod = startPeriod

		if flRetries.Value != "" ***REMOVED***
			retries, err := strconv.ParseInt(flRetries.Value, 10, 32)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if retries < 1 ***REMOVED***
				return nil, fmt.Errorf("--retries must be at least 1 (not %d)", retries)
			***REMOVED***
			healthcheck.Retries = int(retries)
		***REMOVED*** else ***REMOVED***
			healthcheck.Retries = 0
		***REMOVED***

		cmd.Health = &healthcheck
	***REMOVED***
	return cmd, nil
***REMOVED***

func parseExpose(req parseRequest) (*ExposeCommand, error) ***REMOVED***
	portsTab := req.args

	if len(req.args) == 0 ***REMOVED***
		return nil, errAtLeastOneArgument("EXPOSE")
	***REMOVED***

	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sort.Strings(portsTab)
	return &ExposeCommand***REMOVED***
		Ports:           portsTab,
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***, nil
***REMOVED***

func parseUser(req parseRequest) (*UserCommand, error) ***REMOVED***
	if len(req.args) != 1 ***REMOVED***
		return nil, errExactlyOneArgument("USER")
	***REMOVED***

	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &UserCommand***REMOVED***
		User:            req.args[0],
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***, nil
***REMOVED***

func parseVolume(req parseRequest) (*VolumeCommand, error) ***REMOVED***
	if len(req.args) == 0 ***REMOVED***
		return nil, errAtLeastOneArgument("VOLUME")
	***REMOVED***

	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cmd := &VolumeCommand***REMOVED***
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***

	for _, v := range req.args ***REMOVED***
		v = strings.TrimSpace(v)
		if v == "" ***REMOVED***
			return nil, errors.New("VOLUME specified can not be an empty string")
		***REMOVED***
		cmd.Volumes = append(cmd.Volumes, v)
	***REMOVED***
	return cmd, nil

***REMOVED***

func parseStopSignal(req parseRequest) (*StopSignalCommand, error) ***REMOVED***
	if len(req.args) != 1 ***REMOVED***
		return nil, errExactlyOneArgument("STOPSIGNAL")
	***REMOVED***
	sig := req.args[0]

	cmd := &StopSignalCommand***REMOVED***
		Signal:          sig,
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***
	return cmd, nil

***REMOVED***

func parseArg(req parseRequest) (*ArgCommand, error) ***REMOVED***
	if len(req.args) != 1 ***REMOVED***
		return nil, errExactlyOneArgument("ARG")
	***REMOVED***

	var (
		name     string
		newValue *string
	)

	arg := req.args[0]
	// 'arg' can just be a name or name-value pair. Note that this is different
	// from 'env' that handles the split of name and value at the parser level.
	// The reason for doing it differently for 'arg' is that we support just
	// defining an arg and not assign it a value (while 'env' always expects a
	// name-value pair). If possible, it will be good to harmonize the two.
	if strings.Contains(arg, "=") ***REMOVED***
		parts := strings.SplitN(arg, "=", 2)
		if len(parts[0]) == 0 ***REMOVED***
			return nil, errBlankCommandNames("ARG")
		***REMOVED***

		name = parts[0]
		newValue = &parts[1]
	***REMOVED*** else ***REMOVED***
		name = arg
	***REMOVED***

	return &ArgCommand***REMOVED***
		Key:             name,
		Value:           newValue,
		withNameAndCode: newWithNameAndCode(req),
	***REMOVED***, nil
***REMOVED***

func parseShell(req parseRequest) (*ShellCommand, error) ***REMOVED***
	if err := req.flags.Parse(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	shellSlice := handleJSONArgs(req.args, req.attributes)
	switch ***REMOVED***
	case len(shellSlice) == 0:
		// SHELL []
		return nil, errAtLeastOneArgument("SHELL")
	case req.attributes["json"]:
		// SHELL ["powershell", "-command"]

		return &ShellCommand***REMOVED***
			Shell:           strslice.StrSlice(shellSlice),
			withNameAndCode: newWithNameAndCode(req),
		***REMOVED***, nil
	default:
		// SHELL powershell -command - not JSON
		return nil, errNotJSON("SHELL", req.original)
	***REMOVED***
***REMOVED***

func errAtLeastOneArgument(command string) error ***REMOVED***
	return errors.Errorf("%s requires at least one argument", command)
***REMOVED***

func errExactlyOneArgument(command string) error ***REMOVED***
	return errors.Errorf("%s requires exactly one argument", command)
***REMOVED***

func errNoDestinationArgument(command string) error ***REMOVED***
	return errors.Errorf("%s requires at least two arguments, but only one was provided. Destination could not be determined.", command)
***REMOVED***

func errBlankCommandNames(command string) error ***REMOVED***
	return errors.Errorf("%s names can not be blank", command)
***REMOVED***

func errTooManyArguments(command string) error ***REMOVED***
	return errors.Errorf("Bad input to %s, too many arguments", command)
***REMOVED***
