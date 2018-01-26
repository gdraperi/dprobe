// Copyright © 2013 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//Package cobra is a commander providing a simple interface to create powerful modern CLI interfaces.
//In addition to providing an interface, Cobra simultaneously provides a controller to organize your application code.
package cobra

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	flag "github.com/spf13/pflag"
)

// Command is just that, a command for your application.
// eg.  'go run' ... 'run' is the command. Cobra requires
// you to define the usage and description as part of your command
// definition to ensure usability.
type Command struct ***REMOVED***
	// Name is the command name, usually the executable's name.
	name string
	// The one-line usage message.
	Use string
	// An array of aliases that can be used instead of the first word in Use.
	Aliases []string
	// An array of command names for which this command will be suggested - similar to aliases but only suggests.
	SuggestFor []string
	// The short description shown in the 'help' output.
	Short string
	// The long message shown in the 'help <this-command>' output.
	Long string
	// Examples of how to use the command
	Example string
	// List of all valid non-flag arguments that are accepted in bash completions
	ValidArgs []string
	// List of aliases for ValidArgs. These are not suggested to the user in the bash
	// completion, but accepted if entered manually.
	ArgAliases []string
	// Expected arguments
	Args PositionalArgs
	// Custom functions used by the bash autocompletion generator
	BashCompletionFunction string
	// Is this command deprecated and should print this string when used?
	Deprecated string
	// Is this command hidden and should NOT show up in the list of available commands?
	Hidden bool
	// Tags are key/value pairs that can be used by applications to identify or
	// group commands
	Tags map[string]string
	// Full set of flags
	flags *flag.FlagSet
	// Set of flags childrens of this command will inherit
	pflags *flag.FlagSet
	// Flags that are declared specifically by this command (not inherited).
	lflags *flag.FlagSet
	// SilenceErrors is an option to quiet errors down stream
	SilenceErrors bool
	// Silence Usage is an option to silence usage when an error occurs.
	SilenceUsage bool
	// The *Run functions are executed in the following order:
	//   * PersistentPreRun()
	//   * PreRun()
	//   * Run()
	//   * PostRun()
	//   * PersistentPostRun()
	// All functions get the same args, the arguments after the command name
	// PersistentPreRun: children of this command will inherit and execute
	PersistentPreRun func(cmd *Command, args []string)
	// PersistentPreRunE: PersistentPreRun but returns an error
	PersistentPreRunE func(cmd *Command, args []string) error
	// PreRun: children of this command will not inherit.
	PreRun func(cmd *Command, args []string)
	// PreRunE: PreRun but returns an error
	PreRunE func(cmd *Command, args []string) error
	// Run: Typically the actual work function. Most commands will only implement this
	Run func(cmd *Command, args []string)
	// RunE: Run but returns an error
	RunE func(cmd *Command, args []string) error
	// PostRun: run after the Run command.
	PostRun func(cmd *Command, args []string)
	// PostRunE: PostRun but returns an error
	PostRunE func(cmd *Command, args []string) error
	// PersistentPostRun: children of this command will inherit and execute after PostRun
	PersistentPostRun func(cmd *Command, args []string)
	// PersistentPostRunE: PersistentPostRun but returns an error
	PersistentPostRunE func(cmd *Command, args []string) error
	// DisableAutoGenTag remove
	DisableAutoGenTag bool
	// Commands is the list of commands supported by this program.
	commands []*Command
	// Parent Command for this command
	parent *Command
	// max lengths of commands' string lengths for use in padding
	commandsMaxUseLen         int
	commandsMaxCommandPathLen int
	commandsMaxNameLen        int
	// is commands slice are sorted or not
	commandsAreSorted bool

	flagErrorBuf *bytes.Buffer

	args          []string             // actual args parsed from flags
	output        *io.Writer           // nil means stderr; use Out() method instead
	usageFunc     func(*Command) error // Usage can be defined by application
	usageTemplate string               // Can be defined by Application
	flagErrorFunc func(*Command, error) error
	helpTemplate  string                   // Can be defined by Application
	helpFunc      func(*Command, []string) // Help can be defined by application
	helpCommand   *Command                 // The help command
	// The global normalization function that we can use on every pFlag set and children commands
	globNormFunc func(f *flag.FlagSet, name string) flag.NormalizedName

	// Disable the suggestions based on Levenshtein distance that go along with 'unknown command' messages
	DisableSuggestions bool
	// If displaying suggestions, allows to set the minimum levenshtein distance to display, must be > 0
	SuggestionsMinimumDistance int

	// Disable the flag parsing. If this is true all flags will be passed to the command as arguments.
	DisableFlagParsing bool

	// TraverseChildren parses flags on all parents before executing child command
	TraverseChildren bool
***REMOVED***

// os.Args[1:] by default, if desired, can be overridden
// particularly useful when testing.
func (c *Command) SetArgs(a []string) ***REMOVED***
	c.args = a
***REMOVED***

func (c *Command) getOut(def io.Writer) io.Writer ***REMOVED***
	if c.output != nil ***REMOVED***
		return *c.output
	***REMOVED***

	if c.HasParent() ***REMOVED***
		return c.parent.Out()
	***REMOVED***
	return def
***REMOVED***

func (c *Command) Out() io.Writer ***REMOVED***
	return c.getOut(os.Stderr)
***REMOVED***

func (c *Command) getOutOrStdout() io.Writer ***REMOVED***
	return c.getOut(os.Stdout)
***REMOVED***

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (c *Command) SetOutput(output io.Writer) ***REMOVED***
	c.output = &output
***REMOVED***

// Usage can be defined by application
func (c *Command) SetUsageFunc(f func(*Command) error) ***REMOVED***
	c.usageFunc = f
***REMOVED***

// Can be defined by Application
func (c *Command) SetUsageTemplate(s string) ***REMOVED***
	c.usageTemplate = s
***REMOVED***

// SetFlagErrorFunc sets a function to generate an error when flag parsing
// fails
func (c *Command) SetFlagErrorFunc(f func(*Command, error) error) ***REMOVED***
	c.flagErrorFunc = f
***REMOVED***

// Can be defined by Application
func (c *Command) SetHelpFunc(f func(*Command, []string)) ***REMOVED***
	c.helpFunc = f
***REMOVED***

func (c *Command) SetHelpCommand(cmd *Command) ***REMOVED***
	c.helpCommand = cmd
***REMOVED***

// Can be defined by Application
func (c *Command) SetHelpTemplate(s string) ***REMOVED***
	c.helpTemplate = s
***REMOVED***

// SetGlobalNormalizationFunc sets a normalization function to all flag sets and also to child commands.
// The user should not have a cyclic dependency on commands.
func (c *Command) SetGlobalNormalizationFunc(n func(f *flag.FlagSet, name string) flag.NormalizedName) ***REMOVED***
	c.Flags().SetNormalizeFunc(n)
	c.PersistentFlags().SetNormalizeFunc(n)
	c.globNormFunc = n

	for _, command := range c.commands ***REMOVED***
		command.SetGlobalNormalizationFunc(n)
	***REMOVED***
***REMOVED***

func (c *Command) UsageFunc() (f func(*Command) error) ***REMOVED***
	if c.usageFunc != nil ***REMOVED***
		return c.usageFunc
	***REMOVED***

	if c.HasParent() ***REMOVED***
		return c.parent.UsageFunc()
	***REMOVED***
	return func(c *Command) error ***REMOVED***
		err := tmpl(c.Out(), c.UsageTemplate(), c)
		if err != nil ***REMOVED***
			fmt.Print(err)
		***REMOVED***
		return err
	***REMOVED***
***REMOVED***

// HelpFunc returns either the function set by SetHelpFunc for this command
// or a parent, or it returns a function which calls c.Help()
func (c *Command) HelpFunc() func(*Command, []string) ***REMOVED***
	cmd := c
	for cmd != nil ***REMOVED***
		if cmd.helpFunc != nil ***REMOVED***
			return cmd.helpFunc
		***REMOVED***
		cmd = cmd.parent
	***REMOVED***
	return func(*Command, []string) ***REMOVED***
		err := c.Help()
		if err != nil ***REMOVED***
			c.Println(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// FlagErrorFunc returns either the function set by SetFlagErrorFunc for this
// command or a parent, or it returns a function which returns the original
// error.
func (c *Command) FlagErrorFunc() (f func(*Command, error) error) ***REMOVED***
	if c.flagErrorFunc != nil ***REMOVED***
		return c.flagErrorFunc
	***REMOVED***

	if c.HasParent() ***REMOVED***
		return c.parent.FlagErrorFunc()
	***REMOVED***
	return func(c *Command, err error) error ***REMOVED***
		return err
	***REMOVED***
***REMOVED***

var minUsagePadding = 25

func (c *Command) UsagePadding() int ***REMOVED***
	if c.parent == nil || minUsagePadding > c.parent.commandsMaxUseLen ***REMOVED***
		return minUsagePadding
	***REMOVED***
	return c.parent.commandsMaxUseLen
***REMOVED***

var minCommandPathPadding = 11

//
func (c *Command) CommandPathPadding() int ***REMOVED***
	if c.parent == nil || minCommandPathPadding > c.parent.commandsMaxCommandPathLen ***REMOVED***
		return minCommandPathPadding
	***REMOVED***
	return c.parent.commandsMaxCommandPathLen
***REMOVED***

var minNamePadding = 11

func (c *Command) NamePadding() int ***REMOVED***
	if c.parent == nil || minNamePadding > c.parent.commandsMaxNameLen ***REMOVED***
		return minNamePadding
	***REMOVED***
	return c.parent.commandsMaxNameLen
***REMOVED***

func (c *Command) UsageTemplate() string ***REMOVED***
	if c.usageTemplate != "" ***REMOVED***
		return c.usageTemplate
	***REMOVED***

	if c.HasParent() ***REMOVED***
		return c.parent.UsageTemplate()
	***REMOVED***
	return `Usage:***REMOVED******REMOVED***if .Runnable***REMOVED******REMOVED***
  ***REMOVED******REMOVED***if .HasAvailableFlags***REMOVED******REMOVED******REMOVED******REMOVED***appendIfNotPresent .UseLine "[flags]"***REMOVED******REMOVED******REMOVED******REMOVED***else***REMOVED******REMOVED******REMOVED******REMOVED***.UseLine***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasAvailableSubCommands***REMOVED******REMOVED***
  ***REMOVED******REMOVED*** .CommandPath***REMOVED******REMOVED*** [command]***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if gt .Aliases 0***REMOVED******REMOVED***

Aliases:
  ***REMOVED******REMOVED***.NameAndAliases***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasExample***REMOVED******REMOVED***

Examples:
***REMOVED******REMOVED*** .Example ***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED*** if .HasAvailableSubCommands***REMOVED******REMOVED***

Available Commands:***REMOVED******REMOVED***range .Commands***REMOVED******REMOVED******REMOVED******REMOVED***if .IsAvailableCommand***REMOVED******REMOVED***
  ***REMOVED******REMOVED***rpad .Name .NamePadding ***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Short***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED*** if .HasAvailableLocalFlags***REMOVED******REMOVED***

Flags:
***REMOVED******REMOVED***.LocalFlags.FlagUsages | trimRightSpace***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED*** if .HasAvailableInheritedFlags***REMOVED******REMOVED***

Global Flags:
***REMOVED******REMOVED***.InheritedFlags.FlagUsages | trimRightSpace***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasHelpSubCommands***REMOVED******REMOVED***

Additional help topics:***REMOVED******REMOVED***range .Commands***REMOVED******REMOVED******REMOVED******REMOVED***if .IsHelpCommand***REMOVED******REMOVED***
  ***REMOVED******REMOVED***rpad .CommandPath .CommandPathPadding***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Short***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED*** if .HasAvailableSubCommands ***REMOVED******REMOVED***

Use "***REMOVED******REMOVED***.CommandPath***REMOVED******REMOVED*** [command] --help" for more information about a command.***REMOVED******REMOVED***end***REMOVED******REMOVED***
`
***REMOVED***

func (c *Command) HelpTemplate() string ***REMOVED***
	if c.helpTemplate != "" ***REMOVED***
		return c.helpTemplate
	***REMOVED***

	if c.HasParent() ***REMOVED***
		return c.parent.HelpTemplate()
	***REMOVED***
	return `***REMOVED******REMOVED***with or .Long .Short ***REMOVED******REMOVED******REMOVED******REMOVED***. | trim***REMOVED******REMOVED***

***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if or .Runnable .HasSubCommands***REMOVED******REMOVED******REMOVED******REMOVED***.UsageString***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***`
***REMOVED***

// Really only used when casting a command to a commander
func (c *Command) resetChildrensParents() ***REMOVED***
	for _, x := range c.commands ***REMOVED***
		x.parent = c
	***REMOVED***
***REMOVED***

// Test if the named flag is a boolean flag.
func isBooleanFlag(name string, f *flag.FlagSet) bool ***REMOVED***
	flag := f.Lookup(name)
	if flag == nil ***REMOVED***
		return false
	***REMOVED***
	return flag.Value.Type() == "bool"
***REMOVED***

// Test if the named flag is a boolean flag.
func isBooleanShortFlag(name string, f *flag.FlagSet) bool ***REMOVED***
	result := false
	f.VisitAll(func(f *flag.Flag) ***REMOVED***
		if f.Shorthand == name && f.Value.Type() == "bool" ***REMOVED***
			result = true
		***REMOVED***
	***REMOVED***)
	return result
***REMOVED***

func stripFlags(args []string, c *Command) []string ***REMOVED***
	if len(args) < 1 ***REMOVED***
		return args
	***REMOVED***
	c.mergePersistentFlags()

	commands := []string***REMOVED******REMOVED***

	inQuote := false
	inFlag := false
	for _, y := range args ***REMOVED***
		if !inQuote ***REMOVED***
			switch ***REMOVED***
			case strings.HasPrefix(y, "\""):
				inQuote = true
			case strings.Contains(y, "=\""):
				inQuote = true
			case strings.HasPrefix(y, "--") && !strings.Contains(y, "="):
				// TODO: this isn't quite right, we should really check ahead for 'true' or 'false'
				inFlag = !isBooleanFlag(y[2:], c.Flags())
			case strings.HasPrefix(y, "-") && !strings.Contains(y, "=") && len(y) == 2 && !isBooleanShortFlag(y[1:], c.Flags()):
				inFlag = true
			case inFlag:
				inFlag = false
			case y == "":
				// strip empty commands, as the go tests expect this to be ok....
			case !strings.HasPrefix(y, "-"):
				commands = append(commands, y)
				inFlag = false
			***REMOVED***
		***REMOVED***

		if strings.HasSuffix(y, "\"") && !strings.HasSuffix(y, "\\\"") ***REMOVED***
			inQuote = false
		***REMOVED***
	***REMOVED***

	return commands
***REMOVED***

// argsMinusFirstX removes only the first x from args.  Otherwise, commands that look like
// openshift admin policy add-role-to-user admin my-user, lose the admin argument (arg[4]).
func argsMinusFirstX(args []string, x string) []string ***REMOVED***
	for i, y := range args ***REMOVED***
		if x == y ***REMOVED***
			ret := []string***REMOVED******REMOVED***
			ret = append(ret, args[:i]...)
			ret = append(ret, args[i+1:]...)
			return ret
		***REMOVED***
	***REMOVED***
	return args
***REMOVED***

func isFlagArg(arg string) bool ***REMOVED***
	return ((len(arg) >= 3 && arg[1] == '-') ||
		(len(arg) >= 2 && arg[0] == '-' && arg[1] != '-'))
***REMOVED***

// Find the target command given the args and command tree
// Meant to be run on the highest node. Only searches down.
func (c *Command) Find(args []string) (*Command, []string, error) ***REMOVED***
	var innerfind func(*Command, []string) (*Command, []string)

	innerfind = func(c *Command, innerArgs []string) (*Command, []string) ***REMOVED***
		argsWOflags := stripFlags(innerArgs, c)
		if len(argsWOflags) == 0 ***REMOVED***
			return c, innerArgs
		***REMOVED***
		nextSubCmd := argsWOflags[0]

		cmd := c.findNext(nextSubCmd)
		if cmd != nil ***REMOVED***
			return innerfind(cmd, argsMinusFirstX(innerArgs, nextSubCmd))
		***REMOVED***
		return c, innerArgs
	***REMOVED***

	commandFound, a := innerfind(c, args)
	if commandFound.Args == nil ***REMOVED***
		return commandFound, a, legacyArgs(commandFound, stripFlags(a, commandFound))
	***REMOVED***
	return commandFound, a, nil
***REMOVED***

func (c *Command) findNext(next string) *Command ***REMOVED***
	matches := make([]*Command, 0)
	for _, cmd := range c.commands ***REMOVED***
		if cmd.Name() == next || cmd.HasAlias(next) ***REMOVED***
			return cmd
		***REMOVED***
		if EnablePrefixMatching && cmd.HasNameOrAliasPrefix(next) ***REMOVED***
			matches = append(matches, cmd)
		***REMOVED***
	***REMOVED***

	if len(matches) == 1 ***REMOVED***
		return matches[0]
	***REMOVED***
	return nil
***REMOVED***

// Traverse the command tree to find the command, and parse args for
// each parent.
func (c *Command) Traverse(args []string) (*Command, []string, error) ***REMOVED***
	flags := []string***REMOVED******REMOVED***
	inFlag := false

	for i, arg := range args ***REMOVED***
		switch ***REMOVED***
		// A long flag with a space separated value
		case strings.HasPrefix(arg, "--") && !strings.Contains(arg, "="):
			// TODO: this isn't quite right, we should really check ahead for 'true' or 'false'
			inFlag = !isBooleanFlag(arg[2:], c.Flags())
			flags = append(flags, arg)
			continue
		// A short flag with a space separated value
		case strings.HasPrefix(arg, "-") && !strings.Contains(arg, "=") && len(arg) == 2 && !isBooleanShortFlag(arg[1:], c.Flags()):
			inFlag = true
			flags = append(flags, arg)
			continue
		// The value for a flag
		case inFlag:
			inFlag = false
			flags = append(flags, arg)
			continue
		// A flag without a value, or with an `=` separated value
		case isFlagArg(arg):
			flags = append(flags, arg)
			continue
		***REMOVED***

		cmd := c.findNext(arg)
		if cmd == nil ***REMOVED***
			return c, args, nil
		***REMOVED***

		if err := c.ParseFlags(flags); err != nil ***REMOVED***
			return nil, args, err
		***REMOVED***
		return cmd.Traverse(args[i+1:])
	***REMOVED***
	return c, args, nil
***REMOVED***

func (c *Command) findSuggestions(arg string) string ***REMOVED***
	if c.DisableSuggestions ***REMOVED***
		return ""
	***REMOVED***
	if c.SuggestionsMinimumDistance <= 0 ***REMOVED***
		c.SuggestionsMinimumDistance = 2
	***REMOVED***
	suggestionsString := ""
	if suggestions := c.SuggestionsFor(arg); len(suggestions) > 0 ***REMOVED***
		suggestionsString += "\n\nDid you mean this?\n"
		for _, s := range suggestions ***REMOVED***
			suggestionsString += fmt.Sprintf("\t%v\n", s)
		***REMOVED***
	***REMOVED***
	return suggestionsString
***REMOVED***

func (c *Command) SuggestionsFor(typedName string) []string ***REMOVED***
	suggestions := []string***REMOVED******REMOVED***
	for _, cmd := range c.commands ***REMOVED***
		if cmd.IsAvailableCommand() ***REMOVED***
			levenshteinDistance := ld(typedName, cmd.Name(), true)
			suggestByLevenshtein := levenshteinDistance <= c.SuggestionsMinimumDistance
			suggestByPrefix := strings.HasPrefix(strings.ToLower(cmd.Name()), strings.ToLower(typedName))
			if suggestByLevenshtein || suggestByPrefix ***REMOVED***
				suggestions = append(suggestions, cmd.Name())
			***REMOVED***
			for _, explicitSuggestion := range cmd.SuggestFor ***REMOVED***
				if strings.EqualFold(typedName, explicitSuggestion) ***REMOVED***
					suggestions = append(suggestions, cmd.Name())
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return suggestions
***REMOVED***

func (c *Command) VisitParents(fn func(*Command)) ***REMOVED***
	var traverse func(*Command) *Command

	traverse = func(x *Command) *Command ***REMOVED***
		if x != c ***REMOVED***
			fn(x)
		***REMOVED***
		if x.HasParent() ***REMOVED***
			return traverse(x.parent)
		***REMOVED***
		return x
	***REMOVED***
	traverse(c)
***REMOVED***

func (c *Command) Root() *Command ***REMOVED***
	var findRoot func(*Command) *Command

	findRoot = func(x *Command) *Command ***REMOVED***
		if x.HasParent() ***REMOVED***
			return findRoot(x.parent)
		***REMOVED***
		return x
	***REMOVED***

	return findRoot(c)
***REMOVED***

// ArgsLenAtDash will return the length of f.Args at the moment when a -- was
// found during arg parsing. This allows your program to know which args were
// before the -- and which came after. (Description from
// https://godoc.org/github.com/spf13/pflag#FlagSet.ArgsLenAtDash).
func (c *Command) ArgsLenAtDash() int ***REMOVED***
	return c.Flags().ArgsLenAtDash()
***REMOVED***

func (c *Command) execute(a []string) (err error) ***REMOVED***
	if c == nil ***REMOVED***
		return fmt.Errorf("Called Execute() on a nil Command")
	***REMOVED***

	if len(c.Deprecated) > 0 ***REMOVED***
		c.Printf("Command %q is deprecated, %s\n", c.Name(), c.Deprecated)
	***REMOVED***

	// initialize help flag as the last point possible to allow for user
	// overriding
	c.initHelpFlag()

	err = c.ParseFlags(a)
	if err != nil ***REMOVED***
		return c.FlagErrorFunc()(c, err)
	***REMOVED***
	// If help is called, regardless of other flags, return we want help
	// Also say we need help if the command isn't runnable.
	helpVal, err := c.Flags().GetBool("help")
	if err != nil ***REMOVED***
		// should be impossible to get here as we always declare a help
		// flag in initHelpFlag()
		c.Println("\"help\" flag declared as non-bool. Please correct your code")
		return err
	***REMOVED***
	if helpVal || !c.Runnable() ***REMOVED***
		return flag.ErrHelp
	***REMOVED***

	c.preRun()

	argWoFlags := c.Flags().Args()
	if c.DisableFlagParsing ***REMOVED***
		argWoFlags = a
	***REMOVED***

	if err := c.ValidateArgs(argWoFlags); err != nil ***REMOVED***
		return err
	***REMOVED***

	for p := c; p != nil; p = p.Parent() ***REMOVED***
		if p.PersistentPreRunE != nil ***REMOVED***
			if err := p.PersistentPreRunE(c, argWoFlags); err != nil ***REMOVED***
				return err
			***REMOVED***
			break
		***REMOVED*** else if p.PersistentPreRun != nil ***REMOVED***
			p.PersistentPreRun(c, argWoFlags)
			break
		***REMOVED***
	***REMOVED***
	if c.PreRunE != nil ***REMOVED***
		if err := c.PreRunE(c, argWoFlags); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else if c.PreRun != nil ***REMOVED***
		c.PreRun(c, argWoFlags)
	***REMOVED***

	if c.RunE != nil ***REMOVED***
		if err := c.RunE(c, argWoFlags); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.Run(c, argWoFlags)
	***REMOVED***
	if c.PostRunE != nil ***REMOVED***
		if err := c.PostRunE(c, argWoFlags); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else if c.PostRun != nil ***REMOVED***
		c.PostRun(c, argWoFlags)
	***REMOVED***
	for p := c; p != nil; p = p.Parent() ***REMOVED***
		if p.PersistentPostRunE != nil ***REMOVED***
			if err := p.PersistentPostRunE(c, argWoFlags); err != nil ***REMOVED***
				return err
			***REMOVED***
			break
		***REMOVED*** else if p.PersistentPostRun != nil ***REMOVED***
			p.PersistentPostRun(c, argWoFlags)
			break
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *Command) preRun() ***REMOVED***
	for _, x := range initializers ***REMOVED***
		x()
	***REMOVED***
***REMOVED***

func (c *Command) errorMsgFromParse() string ***REMOVED***
	s := c.flagErrorBuf.String()

	x := strings.Split(s, "\n")

	if len(x) > 0 ***REMOVED***
		return x[0]
	***REMOVED***
	return ""
***REMOVED***

// Call execute to use the args (os.Args[1:] by default)
// and run through the command tree finding appropriate matches
// for commands and then corresponding flags.
func (c *Command) Execute() error ***REMOVED***
	_, err := c.ExecuteC()
	return err
***REMOVED***

func (c *Command) ExecuteC() (cmd *Command, err error) ***REMOVED***

	// Regardless of what command execute is called on, run on Root only
	if c.HasParent() ***REMOVED***
		return c.Root().ExecuteC()
	***REMOVED***

	// windows hook
	if preExecHookFn != nil ***REMOVED***
		preExecHookFn(c)
	***REMOVED***

	// initialize help as the last point possible to allow for user
	// overriding
	c.initHelpCmd()

	var args []string

	// Workaround FAIL with "go test -v" or "cobra.test -test.v", see #155
	if c.args == nil && filepath.Base(os.Args[0]) != "cobra.test" ***REMOVED***
		args = os.Args[1:]
	***REMOVED*** else ***REMOVED***
		args = c.args
	***REMOVED***

	var flags []string
	if c.TraverseChildren ***REMOVED***
		cmd, flags, err = c.Traverse(args)
	***REMOVED*** else ***REMOVED***
		cmd, flags, err = c.Find(args)
	***REMOVED***
	if err != nil ***REMOVED***
		// If found parse to a subcommand and then failed, talk about the subcommand
		if cmd != nil ***REMOVED***
			c = cmd
		***REMOVED***
		if !c.SilenceErrors ***REMOVED***
			c.Println("Error:", err.Error())
			c.Printf("Run '%v --help' for usage.\n", c.CommandPath())
		***REMOVED***
		return c, err
	***REMOVED***

	err = cmd.execute(flags)
	if err != nil ***REMOVED***
		// Always show help if requested, even if SilenceErrors is in
		// effect
		if err == flag.ErrHelp ***REMOVED***
			cmd.HelpFunc()(cmd, args)
			return cmd, nil
		***REMOVED***

		// If root command has SilentErrors flagged,
		// all subcommands should respect it
		if !cmd.SilenceErrors && !c.SilenceErrors ***REMOVED***
			c.Println("Error:", err.Error())
		***REMOVED***

		// If root command has SilentUsage flagged,
		// all subcommands should respect it
		if !cmd.SilenceUsage && !c.SilenceUsage ***REMOVED***
			c.Println(cmd.UsageString())
		***REMOVED***
		return cmd, err
	***REMOVED***
	return cmd, nil
***REMOVED***

func (c *Command) ValidateArgs(args []string) error ***REMOVED***
	if c.Args == nil ***REMOVED***
		return nil
	***REMOVED***
	return c.Args(c, args)
***REMOVED***

func (c *Command) initHelpFlag() ***REMOVED***
	c.mergePersistentFlags()
	if c.Flags().Lookup("help") == nil ***REMOVED***
		c.Flags().BoolP("help", "h", false, "help for "+c.Name())
	***REMOVED***
***REMOVED***

func (c *Command) initHelpCmd() ***REMOVED***
	if c.helpCommand == nil ***REMOVED***
		if !c.HasSubCommands() ***REMOVED***
			return
		***REMOVED***

		c.helpCommand = &Command***REMOVED***
			Use:   "help [command]",
			Short: "Help about any command",
			Long: `Help provides help for any command in the application.
    Simply type ` + c.Name() + ` help [path to command] for full details.`,
			PersistentPreRun:  func(cmd *Command, args []string) ***REMOVED******REMOVED***,
			PersistentPostRun: func(cmd *Command, args []string) ***REMOVED******REMOVED***,

			Run: func(c *Command, args []string) ***REMOVED***
				cmd, _, e := c.Root().Find(args)
				if cmd == nil || e != nil ***REMOVED***
					c.Printf("Unknown help topic %#q.", args)
					c.Root().Usage()
				***REMOVED*** else ***REMOVED***
					helpFunc := cmd.HelpFunc()
					helpFunc(cmd, args)
				***REMOVED***
			***REMOVED***,
		***REMOVED***
	***REMOVED***
	c.AddCommand(c.helpCommand)
***REMOVED***

// Used for testing
func (c *Command) ResetCommands() ***REMOVED***
	c.commands = nil
	c.helpCommand = nil
***REMOVED***

// Sorts commands by their names
type commandSorterByName []*Command

func (c commandSorterByName) Len() int           ***REMOVED*** return len(c) ***REMOVED***
func (c commandSorterByName) Swap(i, j int)      ***REMOVED*** c[i], c[j] = c[j], c[i] ***REMOVED***
func (c commandSorterByName) Less(i, j int) bool ***REMOVED*** return c[i].Name() < c[j].Name() ***REMOVED***

// Commands returns a sorted slice of child commands.
func (c *Command) Commands() []*Command ***REMOVED***
	// do not sort commands if it already sorted or sorting was disabled
	if EnableCommandSorting && !c.commandsAreSorted ***REMOVED***
		sort.Sort(commandSorterByName(c.commands))
		c.commandsAreSorted = true
	***REMOVED***
	return c.commands
***REMOVED***

// AddCommand adds one or more commands to this parent command.
func (c *Command) AddCommand(cmds ...*Command) ***REMOVED***
	for i, x := range cmds ***REMOVED***
		if cmds[i] == c ***REMOVED***
			panic("Command can't be a child of itself")
		***REMOVED***
		cmds[i].parent = c
		// update max lengths
		usageLen := len(x.Use)
		if usageLen > c.commandsMaxUseLen ***REMOVED***
			c.commandsMaxUseLen = usageLen
		***REMOVED***
		commandPathLen := len(x.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen ***REMOVED***
			c.commandsMaxCommandPathLen = commandPathLen
		***REMOVED***
		nameLen := len(x.Name())
		if nameLen > c.commandsMaxNameLen ***REMOVED***
			c.commandsMaxNameLen = nameLen
		***REMOVED***
		// If global normalization function exists, update all children
		if c.globNormFunc != nil ***REMOVED***
			x.SetGlobalNormalizationFunc(c.globNormFunc)
		***REMOVED***
		c.commands = append(c.commands, x)
		c.commandsAreSorted = false
	***REMOVED***
***REMOVED***

// RemoveCommand removes one or more commands from a parent command.
func (c *Command) RemoveCommand(cmds ...*Command) ***REMOVED***
	commands := []*Command***REMOVED******REMOVED***
main:
	for _, command := range c.commands ***REMOVED***
		for _, cmd := range cmds ***REMOVED***
			if command == cmd ***REMOVED***
				command.parent = nil
				continue main
			***REMOVED***
		***REMOVED***
		commands = append(commands, command)
	***REMOVED***
	c.commands = commands
	// recompute all lengths
	c.commandsMaxUseLen = 0
	c.commandsMaxCommandPathLen = 0
	c.commandsMaxNameLen = 0
	for _, command := range c.commands ***REMOVED***
		usageLen := len(command.Use)
		if usageLen > c.commandsMaxUseLen ***REMOVED***
			c.commandsMaxUseLen = usageLen
		***REMOVED***
		commandPathLen := len(command.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen ***REMOVED***
			c.commandsMaxCommandPathLen = commandPathLen
		***REMOVED***
		nameLen := len(command.Name())
		if nameLen > c.commandsMaxNameLen ***REMOVED***
			c.commandsMaxNameLen = nameLen
		***REMOVED***
	***REMOVED***
***REMOVED***

// Print is a convenience method to Print to the defined output
func (c *Command) Print(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	fmt.Fprint(c.Out(), i...)
***REMOVED***

// Println is a convenience method to Println to the defined output
func (c *Command) Println(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	str := fmt.Sprintln(i...)
	c.Print(str)
***REMOVED***

// Printf is a convenience method to Printf to the defined output
func (c *Command) Printf(format string, i ...interface***REMOVED******REMOVED***) ***REMOVED***
	str := fmt.Sprintf(format, i...)
	c.Print(str)
***REMOVED***

// Output the usage for the command
// Used when a user provides invalid input
// Can be defined by user by overriding UsageFunc
func (c *Command) Usage() error ***REMOVED***
	c.mergePersistentFlags()
	err := c.UsageFunc()(c)
	return err
***REMOVED***

// Output the help for the command
// Used when a user calls help [command]
// by the default HelpFunc in the commander
func (c *Command) Help() error ***REMOVED***
	c.mergePersistentFlags()
	err := tmpl(c.getOutOrStdout(), c.HelpTemplate(), c)
	return err
***REMOVED***

func (c *Command) UsageString() string ***REMOVED***
	tmpOutput := c.output
	bb := new(bytes.Buffer)
	c.SetOutput(bb)
	c.Usage()
	c.output = tmpOutput
	return bb.String()
***REMOVED***

// CommandPath returns the full path to this command.
func (c *Command) CommandPath() string ***REMOVED***
	str := c.Name()
	x := c
	for x.HasParent() ***REMOVED***
		str = x.parent.Name() + " " + str
		x = x.parent
	***REMOVED***
	return str
***REMOVED***

//The full usage for a given command (including parents)
func (c *Command) UseLine() string ***REMOVED***
	str := ""
	if c.HasParent() ***REMOVED***
		str = c.parent.CommandPath() + " "
	***REMOVED***
	return str + c.Use
***REMOVED***

// For use in determining which flags have been assigned to which commands
// and which persist
func (c *Command) DebugFlags() ***REMOVED***
	c.Println("DebugFlags called on", c.Name())
	var debugflags func(*Command)

	debugflags = func(x *Command) ***REMOVED***
		if x.HasFlags() || x.HasPersistentFlags() ***REMOVED***
			c.Println(x.Name())
		***REMOVED***
		if x.HasFlags() ***REMOVED***
			x.flags.VisitAll(func(f *flag.Flag) ***REMOVED***
				if x.HasPersistentFlags() ***REMOVED***
					if x.persistentFlag(f.Name) == nil ***REMOVED***
						c.Println("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value, "  [L]")
					***REMOVED*** else ***REMOVED***
						c.Println("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value, "  [LP]")
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					c.Println("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value, "  [L]")
				***REMOVED***
			***REMOVED***)
		***REMOVED***
		if x.HasPersistentFlags() ***REMOVED***
			x.pflags.VisitAll(func(f *flag.Flag) ***REMOVED***
				if x.HasFlags() ***REMOVED***
					if x.flags.Lookup(f.Name) == nil ***REMOVED***
						c.Println("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value, "  [P]")
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					c.Println("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value, "  [P]")
				***REMOVED***
			***REMOVED***)
		***REMOVED***
		c.Println(x.flagErrorBuf)
		if x.HasSubCommands() ***REMOVED***
			for _, y := range x.commands ***REMOVED***
				debugflags(y)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	debugflags(c)
***REMOVED***

// Name returns the command's name: the first word in the use line.
func (c *Command) Name() string ***REMOVED***
	if c.name != "" ***REMOVED***
		return c.name
	***REMOVED***
	name := c.Use
	i := strings.Index(name, " ")
	if i >= 0 ***REMOVED***
		name = name[:i]
	***REMOVED***
	return name
***REMOVED***

// HasAlias determines if a given string is an alias of the command.
func (c *Command) HasAlias(s string) bool ***REMOVED***
	for _, a := range c.Aliases ***REMOVED***
		if a == s ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// HasNameOrAliasPrefix returns true if the Name or any of aliases start
// with prefix
func (c *Command) HasNameOrAliasPrefix(prefix string) bool ***REMOVED***
	if strings.HasPrefix(c.Name(), prefix) ***REMOVED***
		return true
	***REMOVED***
	for _, alias := range c.Aliases ***REMOVED***
		if strings.HasPrefix(alias, prefix) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (c *Command) NameAndAliases() string ***REMOVED***
	return strings.Join(append([]string***REMOVED***c.Name()***REMOVED***, c.Aliases...), ", ")
***REMOVED***

func (c *Command) HasExample() bool ***REMOVED***
	return len(c.Example) > 0
***REMOVED***

// Runnable determines if the command is itself runnable
func (c *Command) Runnable() bool ***REMOVED***
	return c.Run != nil || c.RunE != nil
***REMOVED***

// HasSubCommands determines if the command has children commands
func (c *Command) HasSubCommands() bool ***REMOVED***
	return len(c.commands) > 0
***REMOVED***

// IsAvailableCommand determines if a command is available as a non-help command
// (this includes all non deprecated/hidden commands)
func (c *Command) IsAvailableCommand() bool ***REMOVED***
	if len(c.Deprecated) != 0 || c.Hidden ***REMOVED***
		return false
	***REMOVED***

	if c.HasParent() && c.Parent().helpCommand == c ***REMOVED***
		return false
	***REMOVED***

	if c.Runnable() || c.HasAvailableSubCommands() ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// IsHelpCommand determines if a command is a 'help' command; a help command is
// determined by the fact that it is NOT runnable/hidden/deprecated, and has no
// sub commands that are runnable/hidden/deprecated
func (c *Command) IsHelpCommand() bool ***REMOVED***

	// if a command is runnable, deprecated, or hidden it is not a 'help' command
	if c.Runnable() || len(c.Deprecated) != 0 || c.Hidden ***REMOVED***
		return false
	***REMOVED***

	// if any non-help sub commands are found, the command is not a 'help' command
	for _, sub := range c.commands ***REMOVED***
		if !sub.IsHelpCommand() ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// the command either has no sub commands, or no non-help sub commands
	return true
***REMOVED***

// HasHelpSubCommands determines if a command has any avilable 'help' sub commands
// that need to be shown in the usage/help default template under 'additional help
// topics'
func (c *Command) HasHelpSubCommands() bool ***REMOVED***

	// return true on the first found available 'help' sub command
	for _, sub := range c.commands ***REMOVED***
		if sub.IsHelpCommand() ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	// the command either has no sub commands, or no available 'help' sub commands
	return false
***REMOVED***

// HasAvailableSubCommands determines if a command has available sub commands that
// need to be shown in the usage/help default template under 'available commands'
func (c *Command) HasAvailableSubCommands() bool ***REMOVED***

	// return true on the first found available (non deprecated/help/hidden)
	// sub command
	for _, sub := range c.commands ***REMOVED***
		if sub.IsAvailableCommand() ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	// the command either has no sub comamnds, or no available (non deprecated/help/hidden)
	// sub commands
	return false
***REMOVED***

// Determine if the command is a child command
func (c *Command) HasParent() bool ***REMOVED***
	return c.parent != nil
***REMOVED***

// GlobalNormalizationFunc returns the global normalization function or nil if doesn't exists
func (c *Command) GlobalNormalizationFunc() func(f *flag.FlagSet, name string) flag.NormalizedName ***REMOVED***
	return c.globNormFunc
***REMOVED***

// Get the complete FlagSet that applies to this command (local and persistent declared here and by all parents)
func (c *Command) Flags() *flag.FlagSet ***REMOVED***
	if c.flags == nil ***REMOVED***
		c.flags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil ***REMOVED***
			c.flagErrorBuf = new(bytes.Buffer)
		***REMOVED***
		c.flags.SetOutput(c.flagErrorBuf)
	***REMOVED***
	return c.flags
***REMOVED***

// LocalNonPersistentFlags are flags specific to this command which will NOT persist to subcommands
func (c *Command) LocalNonPersistentFlags() *flag.FlagSet ***REMOVED***
	persistentFlags := c.PersistentFlags()

	out := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.LocalFlags().VisitAll(func(f *flag.Flag) ***REMOVED***
		if persistentFlags.Lookup(f.Name) == nil ***REMOVED***
			out.AddFlag(f)
		***REMOVED***
	***REMOVED***)
	return out
***REMOVED***

// Get the local FlagSet specifically set in the current command
func (c *Command) LocalFlags() *flag.FlagSet ***REMOVED***
	c.mergePersistentFlags()

	local := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.lflags.VisitAll(func(f *flag.Flag) ***REMOVED***
		local.AddFlag(f)
	***REMOVED***)
	if !c.HasParent() ***REMOVED***
		flag.CommandLine.VisitAll(func(f *flag.Flag) ***REMOVED***
			if local.Lookup(f.Name) == nil ***REMOVED***
				local.AddFlag(f)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
	return local
***REMOVED***

// All Flags which were inherited from parents commands
func (c *Command) InheritedFlags() *flag.FlagSet ***REMOVED***
	c.mergePersistentFlags()

	inherited := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	local := c.LocalFlags()

	var rmerge func(x *Command)

	rmerge = func(x *Command) ***REMOVED***
		if x.HasPersistentFlags() ***REMOVED***
			x.PersistentFlags().VisitAll(func(f *flag.Flag) ***REMOVED***
				if inherited.Lookup(f.Name) == nil && local.Lookup(f.Name) == nil ***REMOVED***
					inherited.AddFlag(f)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
		if x.HasParent() ***REMOVED***
			rmerge(x.parent)
		***REMOVED***
	***REMOVED***

	if c.HasParent() ***REMOVED***
		rmerge(c.parent)
	***REMOVED***

	return inherited
***REMOVED***

// All Flags which were not inherited from parent commands
func (c *Command) NonInheritedFlags() *flag.FlagSet ***REMOVED***
	return c.LocalFlags()
***REMOVED***

// Get the Persistent FlagSet specifically set in the current command
func (c *Command) PersistentFlags() *flag.FlagSet ***REMOVED***
	if c.pflags == nil ***REMOVED***
		c.pflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil ***REMOVED***
			c.flagErrorBuf = new(bytes.Buffer)
		***REMOVED***
		c.pflags.SetOutput(c.flagErrorBuf)
	***REMOVED***
	return c.pflags
***REMOVED***

// For use in testing
func (c *Command) ResetFlags() ***REMOVED***
	c.flagErrorBuf = new(bytes.Buffer)
	c.flagErrorBuf.Reset()
	c.flags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.flags.SetOutput(c.flagErrorBuf)
	c.pflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.pflags.SetOutput(c.flagErrorBuf)
***REMOVED***

// Does the command contain any flags (local plus persistent from the entire structure)
func (c *Command) HasFlags() bool ***REMOVED***
	return c.Flags().HasFlags()
***REMOVED***

// Does the command contain persistent flags
func (c *Command) HasPersistentFlags() bool ***REMOVED***
	return c.PersistentFlags().HasFlags()
***REMOVED***

// Does the command has flags specifically declared locally
func (c *Command) HasLocalFlags() bool ***REMOVED***
	return c.LocalFlags().HasFlags()
***REMOVED***

// Does the command have flags inherited from its parent command
func (c *Command) HasInheritedFlags() bool ***REMOVED***
	return c.InheritedFlags().HasFlags()
***REMOVED***

// Does the command contain any flags (local plus persistent from the entire
// structure) which are not hidden or deprecated
func (c *Command) HasAvailableFlags() bool ***REMOVED***
	return c.Flags().HasAvailableFlags()
***REMOVED***

// Does the command contain persistent flags which are not hidden or deprecated
func (c *Command) HasAvailablePersistentFlags() bool ***REMOVED***
	return c.PersistentFlags().HasAvailableFlags()
***REMOVED***

// Does the command has flags specifically declared locally which are not hidden
// or deprecated
func (c *Command) HasAvailableLocalFlags() bool ***REMOVED***
	return c.LocalFlags().HasAvailableFlags()
***REMOVED***

// Does the command have flags inherited from its parent command which are
// not hidden or deprecated
func (c *Command) HasAvailableInheritedFlags() bool ***REMOVED***
	return c.InheritedFlags().HasAvailableFlags()
***REMOVED***

// Flag climbs up the command tree looking for matching flag
func (c *Command) Flag(name string) (flag *flag.Flag) ***REMOVED***
	flag = c.Flags().Lookup(name)

	if flag == nil ***REMOVED***
		flag = c.persistentFlag(name)
	***REMOVED***

	return
***REMOVED***

// recursively find matching persistent flag
func (c *Command) persistentFlag(name string) (flag *flag.Flag) ***REMOVED***
	if c.HasPersistentFlags() ***REMOVED***
		flag = c.PersistentFlags().Lookup(name)
	***REMOVED***

	if flag == nil && c.HasParent() ***REMOVED***
		flag = c.parent.persistentFlag(name)
	***REMOVED***
	return
***REMOVED***

// ParseFlags parses persistent flag tree & local flags
func (c *Command) ParseFlags(args []string) (err error) ***REMOVED***
	if c.DisableFlagParsing ***REMOVED***
		return nil
	***REMOVED***
	c.mergePersistentFlags()
	err = c.Flags().Parse(args)
	return
***REMOVED***

// Parent returns a commands parent command
func (c *Command) Parent() *Command ***REMOVED***
	return c.parent
***REMOVED***

func (c *Command) mergePersistentFlags() ***REMOVED***
	var rmerge func(x *Command)

	// Save the set of local flags
	if c.lflags == nil ***REMOVED***
		c.lflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil ***REMOVED***
			c.flagErrorBuf = new(bytes.Buffer)
		***REMOVED***
		c.lflags.SetOutput(c.flagErrorBuf)
		addtolocal := func(f *flag.Flag) ***REMOVED***
			c.lflags.AddFlag(f)
		***REMOVED***
		c.Flags().VisitAll(addtolocal)
		c.PersistentFlags().VisitAll(addtolocal)
	***REMOVED***
	rmerge = func(x *Command) ***REMOVED***
		if !x.HasParent() ***REMOVED***
			flag.CommandLine.VisitAll(func(f *flag.Flag) ***REMOVED***
				if x.PersistentFlags().Lookup(f.Name) == nil ***REMOVED***
					x.PersistentFlags().AddFlag(f)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
		if x.HasPersistentFlags() ***REMOVED***
			x.PersistentFlags().VisitAll(func(f *flag.Flag) ***REMOVED***
				if c.Flags().Lookup(f.Name) == nil ***REMOVED***
					c.Flags().AddFlag(f)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
		if x.HasParent() ***REMOVED***
			rmerge(x.parent)
		***REMOVED***
	***REMOVED***

	rmerge(c)
***REMOVED***
