package cobra

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func emptyRun(*Command, []string) ***REMOVED******REMOVED***

func executeCommand(root *Command, args ...string) (output string, err error) ***REMOVED***
	_, output, err = executeCommandC(root, args...)
	return output, err
***REMOVED***

func executeCommandC(root *Command, args ...string) (c *Command, output string, err error) ***REMOVED***
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
***REMOVED***

func resetCommandLineFlagSet() ***REMOVED***
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
***REMOVED***

func checkStringContains(t *testing.T, got, expected string) ***REMOVED***
	if !strings.Contains(got, expected) ***REMOVED***
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expected, got)
	***REMOVED***
***REMOVED***

func checkStringOmits(t *testing.T, got, expected string) ***REMOVED***
	if strings.Contains(got, expected) ***REMOVED***
		t.Errorf("Expected to not contain: \n %v\nGot: %v", expected, got)
	***REMOVED***
***REMOVED***

func TestSingleCommand(t *testing.T) ***REMOVED***
	var rootCmdArgs []string
	rootCmd := &Command***REMOVED***
		Use:  "root",
		Args: ExactArgs(2),
		Run:  func(_ *Command, args []string) ***REMOVED*** rootCmdArgs = args ***REMOVED***,
	***REMOVED***
	aCmd := &Command***REMOVED***Use: "a", Args: NoArgs, Run: emptyRun***REMOVED***
	bCmd := &Command***REMOVED***Use: "b", Args: NoArgs, Run: emptyRun***REMOVED***
	rootCmd.AddCommand(aCmd, bCmd)

	output, err := executeCommand(rootCmd, "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	got := strings.Join(rootCmdArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("rootCmdArgs expected: %q, got: %q", expected, got)
	***REMOVED***
***REMOVED***

func TestChildCommand(t *testing.T) ***REMOVED***
	var child1CmdArgs []string
	rootCmd := &Command***REMOVED***Use: "root", Args: NoArgs, Run: emptyRun***REMOVED***
	child1Cmd := &Command***REMOVED***
		Use:  "child1",
		Args: ExactArgs(2),
		Run:  func(_ *Command, args []string) ***REMOVED*** child1CmdArgs = args ***REMOVED***,
	***REMOVED***
	child2Cmd := &Command***REMOVED***Use: "child2", Args: NoArgs, Run: emptyRun***REMOVED***
	rootCmd.AddCommand(child1Cmd, child2Cmd)

	output, err := executeCommand(rootCmd, "child1", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	got := strings.Join(child1CmdArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("child1CmdArgs expected: %q, got: %q", expected, got)
	***REMOVED***
***REMOVED***

func TestCallCommandWithoutSubcommands(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Args: NoArgs, Run: emptyRun***REMOVED***
	_, err := executeCommand(rootCmd)
	if err != nil ***REMOVED***
		t.Errorf("Calling command without subcommands should not have error: %v", err)
	***REMOVED***
***REMOVED***

func TestRootExecuteUnknownCommand(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(&Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***)

	output, _ := executeCommand(rootCmd, "unknown")

	expected := "Error: unknown command \"unknown\" for \"root\"\nRun 'root --help' for usage.\n"

	if output != expected ***REMOVED***
		t.Errorf("Expected:\n %q\nGot:\n %q\n", expected, output)
	***REMOVED***
***REMOVED***

func TestSubcommandExecuteC(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	c, output, err := executeCommandC(rootCmd, "child")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if c.Name() != "child" ***REMOVED***
		t.Errorf(`invalid command returned from ExecuteC: expected "child"', got %q`, c.Name())
	***REMOVED***
***REMOVED***

func TestRootUnknownCommandSilenced(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.AddCommand(&Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***)

	output, _ := executeCommand(rootCmd, "unknown")
	if output != "" ***REMOVED***
		t.Errorf("Expected blank output, because of silenced usage.\nGot:\n %q\n", output)
	***REMOVED***
***REMOVED***

func TestCommandAlias(t *testing.T) ***REMOVED***
	var timesCmdArgs []string
	rootCmd := &Command***REMOVED***Use: "root", Args: NoArgs, Run: emptyRun***REMOVED***
	echoCmd := &Command***REMOVED***
		Use:     "echo",
		Aliases: []string***REMOVED***"say", "tell"***REMOVED***,
		Args:    NoArgs,
		Run:     emptyRun,
	***REMOVED***
	timesCmd := &Command***REMOVED***
		Use:  "times",
		Args: ExactArgs(2),
		Run:  func(_ *Command, args []string) ***REMOVED*** timesCmdArgs = args ***REMOVED***,
	***REMOVED***
	echoCmd.AddCommand(timesCmd)
	rootCmd.AddCommand(echoCmd)

	output, err := executeCommand(rootCmd, "tell", "times", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	got := strings.Join(timesCmdArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("timesCmdArgs expected: %v, got: %v", expected, got)
	***REMOVED***
***REMOVED***

func TestEnablePrefixMatching(t *testing.T) ***REMOVED***
	EnablePrefixMatching = true

	var aCmdArgs []string
	rootCmd := &Command***REMOVED***Use: "root", Args: NoArgs, Run: emptyRun***REMOVED***
	aCmd := &Command***REMOVED***
		Use:  "aCmd",
		Args: ExactArgs(2),
		Run:  func(_ *Command, args []string) ***REMOVED*** aCmdArgs = args ***REMOVED***,
	***REMOVED***
	bCmd := &Command***REMOVED***Use: "bCmd", Args: NoArgs, Run: emptyRun***REMOVED***
	rootCmd.AddCommand(aCmd, bCmd)

	output, err := executeCommand(rootCmd, "a", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	got := strings.Join(aCmdArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("aCmdArgs expected: %q, got: %q", expected, got)
	***REMOVED***

	EnablePrefixMatching = false
***REMOVED***

func TestAliasPrefixMatching(t *testing.T) ***REMOVED***
	EnablePrefixMatching = true

	var timesCmdArgs []string
	rootCmd := &Command***REMOVED***Use: "root", Args: NoArgs, Run: emptyRun***REMOVED***
	echoCmd := &Command***REMOVED***
		Use:     "echo",
		Aliases: []string***REMOVED***"say", "tell"***REMOVED***,
		Args:    NoArgs,
		Run:     emptyRun,
	***REMOVED***
	timesCmd := &Command***REMOVED***
		Use:  "times",
		Args: ExactArgs(2),
		Run:  func(_ *Command, args []string) ***REMOVED*** timesCmdArgs = args ***REMOVED***,
	***REMOVED***
	echoCmd.AddCommand(timesCmd)
	rootCmd.AddCommand(echoCmd)

	output, err := executeCommand(rootCmd, "sa", "times", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	got := strings.Join(timesCmdArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("timesCmdArgs expected: %v, got: %v", expected, got)
	***REMOVED***

	EnablePrefixMatching = false
***REMOVED***

// TestChildSameName checks the correct behaviour of cobra in cases,
// when an application with name "foo" and with subcommand "foo"
// is executed with args "foo foo".
func TestChildSameName(t *testing.T) ***REMOVED***
	var fooCmdArgs []string
	rootCmd := &Command***REMOVED***Use: "foo", Args: NoArgs, Run: emptyRun***REMOVED***
	fooCmd := &Command***REMOVED***
		Use:  "foo",
		Args: ExactArgs(2),
		Run:  func(_ *Command, args []string) ***REMOVED*** fooCmdArgs = args ***REMOVED***,
	***REMOVED***
	barCmd := &Command***REMOVED***Use: "bar", Args: NoArgs, Run: emptyRun***REMOVED***
	rootCmd.AddCommand(fooCmd, barCmd)

	output, err := executeCommand(rootCmd, "foo", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	got := strings.Join(fooCmdArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("fooCmdArgs expected: %v, got: %v", expected, got)
	***REMOVED***
***REMOVED***

// TestGrandChildSameName checks the correct behaviour of cobra in cases,
// when user has a root command and a grand child
// with the same name.
func TestGrandChildSameName(t *testing.T) ***REMOVED***
	var fooCmdArgs []string
	rootCmd := &Command***REMOVED***Use: "foo", Args: NoArgs, Run: emptyRun***REMOVED***
	barCmd := &Command***REMOVED***Use: "bar", Args: NoArgs, Run: emptyRun***REMOVED***
	fooCmd := &Command***REMOVED***
		Use:  "foo",
		Args: ExactArgs(2),
		Run:  func(_ *Command, args []string) ***REMOVED*** fooCmdArgs = args ***REMOVED***,
	***REMOVED***
	barCmd.AddCommand(fooCmd)
	rootCmd.AddCommand(barCmd)

	output, err := executeCommand(rootCmd, "bar", "foo", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	got := strings.Join(fooCmdArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("fooCmdArgs expected: %v, got: %v", expected, got)
	***REMOVED***
***REMOVED***

func TestFlagLong(t *testing.T) ***REMOVED***
	var cArgs []string
	c := &Command***REMOVED***
		Use:  "c",
		Args: ArbitraryArgs,
		Run:  func(_ *Command, args []string) ***REMOVED*** cArgs = args ***REMOVED***,
	***REMOVED***

	var intFlagValue int
	var stringFlagValue string
	c.Flags().IntVar(&intFlagValue, "intf", -1, "")
	c.Flags().StringVar(&stringFlagValue, "sf", "", "")

	output, err := executeCommand(c, "--intf=7", "--sf=abc", "one", "--", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", err)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if c.ArgsLenAtDash() != 1 ***REMOVED***
		t.Errorf("Expected ArgsLenAtDash: %v but got %v", 1, c.ArgsLenAtDash())
	***REMOVED***
	if intFlagValue != 7 ***REMOVED***
		t.Errorf("Expected intFlagValue: %v, got %v", 7, intFlagValue)
	***REMOVED***
	if stringFlagValue != "abc" ***REMOVED***
		t.Errorf("Expected stringFlagValue: %q, got %q", "abc", stringFlagValue)
	***REMOVED***

	got := strings.Join(cArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("Expected arguments: %q, got %q", expected, got)
	***REMOVED***
***REMOVED***

func TestFlagShort(t *testing.T) ***REMOVED***
	var cArgs []string
	c := &Command***REMOVED***
		Use:  "c",
		Args: ArbitraryArgs,
		Run:  func(_ *Command, args []string) ***REMOVED*** cArgs = args ***REMOVED***,
	***REMOVED***

	var intFlagValue int
	var stringFlagValue string
	c.Flags().IntVarP(&intFlagValue, "intf", "i", -1, "")
	c.Flags().StringVarP(&stringFlagValue, "sf", "s", "", "")

	output, err := executeCommand(c, "-i", "7", "-sabc", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", err)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if intFlagValue != 7 ***REMOVED***
		t.Errorf("Expected flag value: %v, got %v", 7, intFlagValue)
	***REMOVED***
	if stringFlagValue != "abc" ***REMOVED***
		t.Errorf("Expected stringFlagValue: %q, got %q", "abc", stringFlagValue)
	***REMOVED***

	got := strings.Join(cArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("Expected arguments: %q, got %q", expected, got)
	***REMOVED***
***REMOVED***

func TestChildFlag(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	var intFlagValue int
	childCmd.Flags().IntVarP(&intFlagValue, "intf", "i", -1, "")

	output, err := executeCommand(rootCmd, "child", "-i7")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", err)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if intFlagValue != 7 ***REMOVED***
		t.Errorf("Expected flag value: %v, got %v", 7, intFlagValue)
	***REMOVED***
***REMOVED***

func TestChildFlagWithParentLocalFlag(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	var intFlagValue int
	rootCmd.Flags().StringP("sf", "s", "", "")
	childCmd.Flags().IntVarP(&intFlagValue, "intf", "i", -1, "")

	_, err := executeCommand(rootCmd, "child", "-i7", "-sabc")
	if err == nil ***REMOVED***
		t.Errorf("Invalid flag should generate error")
	***REMOVED***

	checkStringContains(t, err.Error(), "unknown shorthand")

	if intFlagValue != 7 ***REMOVED***
		t.Errorf("Expected flag value: %v, got %v", 7, intFlagValue)
	***REMOVED***
***REMOVED***

func TestFlagInvalidInput(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	rootCmd.Flags().IntP("intf", "i", -1, "")

	_, err := executeCommand(rootCmd, "-iabc")
	if err == nil ***REMOVED***
		t.Errorf("Invalid flag value should generate error")
	***REMOVED***

	checkStringContains(t, err.Error(), "invalid syntax")
***REMOVED***

func TestFlagBeforeCommand(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	var flagValue int
	childCmd.Flags().IntVarP(&flagValue, "intf", "i", -1, "")

	// With short flag.
	_, err := executeCommand(rootCmd, "-i7", "child")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
	if flagValue != 7 ***REMOVED***
		t.Errorf("Expected flag value: %v, got %v", 7, flagValue)
	***REMOVED***

	// With long flag.
	_, err = executeCommand(rootCmd, "--intf=8", "child")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
	if flagValue != 8 ***REMOVED***
		t.Errorf("Expected flag value: %v, got %v", 9, flagValue)
	***REMOVED***
***REMOVED***

func TestStripFlags(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		input  []string
		output []string
	***REMOVED******REMOVED***
		***REMOVED***
			[]string***REMOVED***"foo", "bar"***REMOVED***,
			[]string***REMOVED***"foo", "bar"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"foo", "--str", "-s"***REMOVED***,
			[]string***REMOVED***"foo"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"-s", "foo", "--str", "bar"***REMOVED***,
			[]string***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"-i10", "echo"***REMOVED***,
			[]string***REMOVED***"echo"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"-i=10", "echo"***REMOVED***,
			[]string***REMOVED***"echo"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"--int=100", "echo"***REMOVED***,
			[]string***REMOVED***"echo"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"-ib", "echo", "-sfoo", "baz"***REMOVED***,
			[]string***REMOVED***"echo", "baz"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"-i=baz", "bar", "-i", "foo", "blah"***REMOVED***,
			[]string***REMOVED***"bar", "blah"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"--int=baz", "-sbar", "-i", "foo", "blah"***REMOVED***,
			[]string***REMOVED***"blah"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"--bool", "bar", "-i", "foo", "blah"***REMOVED***,
			[]string***REMOVED***"bar", "blah"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"-b", "bar", "-i", "foo", "blah"***REMOVED***,
			[]string***REMOVED***"bar", "blah"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"--persist", "bar"***REMOVED***,
			[]string***REMOVED***"bar"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***"-p", "bar"***REMOVED***,
			[]string***REMOVED***"bar"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***
	c.PersistentFlags().BoolP("persist", "p", false, "")
	c.Flags().IntP("int", "i", -1, "")
	c.Flags().StringP("str", "s", "", "")
	c.Flags().BoolP("bool", "b", false, "")

	for i, test := range tests ***REMOVED***
		got := stripFlags(test.input, c)
		if !reflect.DeepEqual(test.output, got) ***REMOVED***
			t.Errorf("(%v) Expected: %v, got: %v", i, test.output, got)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDisableFlagParsing(t *testing.T) ***REMOVED***
	var cArgs []string
	c := &Command***REMOVED***
		Use:                "c",
		DisableFlagParsing: true,
		Run: func(_ *Command, args []string) ***REMOVED***
			cArgs = args
		***REMOVED***,
	***REMOVED***

	args := []string***REMOVED***"cmd", "-v", "-race", "-file", "foo.go"***REMOVED***
	output, err := executeCommand(c, args...)
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if !reflect.DeepEqual(args, cArgs) ***REMOVED***
		t.Errorf("Expected: %v, got: %v", args, cArgs)
	***REMOVED***
***REMOVED***

func TestPersistentFlagsOnSameCommand(t *testing.T) ***REMOVED***
	var rootCmdArgs []string
	rootCmd := &Command***REMOVED***
		Use:  "root",
		Args: ArbitraryArgs,
		Run:  func(_ *Command, args []string) ***REMOVED*** rootCmdArgs = args ***REMOVED***,
	***REMOVED***

	var flagValue int
	rootCmd.PersistentFlags().IntVarP(&flagValue, "intf", "i", -1, "")

	output, err := executeCommand(rootCmd, "-i7", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	got := strings.Join(rootCmdArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("rootCmdArgs expected: %q, got %q", expected, got)
	***REMOVED***
	if flagValue != 7 ***REMOVED***
		t.Errorf("flagValue expected: %v, got %v", 7, flagValue)
	***REMOVED***
***REMOVED***

// TestEmptyInputs checks,
// if flags correctly parsed with blank strings in args.
func TestEmptyInputs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***

	var flagValue int
	c.Flags().IntVarP(&flagValue, "intf", "i", -1, "")

	output, err := executeCommand(c, "", "-i7", "")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if flagValue != 7 ***REMOVED***
		t.Errorf("flagValue expected: %v, got %v", 7, flagValue)
	***REMOVED***
***REMOVED***

func TestOverwrittenFlag(t *testing.T) ***REMOVED***
	// TODO: This test fails, but should work.
	t.Skip()

	parent := &Command***REMOVED***Use: "parent", Run: emptyRun***REMOVED***
	child := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***

	parent.PersistentFlags().Bool("boolf", false, "")
	parent.PersistentFlags().Int("intf", -1, "")
	child.Flags().String("strf", "", "")
	child.Flags().Int("intf", -1, "")

	parent.AddCommand(child)

	childInherited := child.InheritedFlags()
	childLocal := child.LocalFlags()

	if childLocal.Lookup("strf") == nil ***REMOVED***
		t.Error(`LocalFlags expected to contain "strf", got "nil"`)
	***REMOVED***
	if childInherited.Lookup("boolf") == nil ***REMOVED***
		t.Error(`InheritedFlags expected to contain "boolf", got "nil"`)
	***REMOVED***

	if childInherited.Lookup("intf") != nil ***REMOVED***
		t.Errorf(`InheritedFlags should not contain overwritten flag "intf"`)
	***REMOVED***
	if childLocal.Lookup("intf") == nil ***REMOVED***
		t.Error(`LocalFlags expected to contain "intf", got "nil"`)
	***REMOVED***
***REMOVED***

func TestPersistentFlagsOnChild(t *testing.T) ***REMOVED***
	var childCmdArgs []string
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***
		Use:  "child",
		Args: ArbitraryArgs,
		Run:  func(_ *Command, args []string) ***REMOVED*** childCmdArgs = args ***REMOVED***,
	***REMOVED***
	rootCmd.AddCommand(childCmd)

	var parentFlagValue int
	var childFlagValue int
	rootCmd.PersistentFlags().IntVarP(&parentFlagValue, "parentf", "p", -1, "")
	childCmd.Flags().IntVarP(&childFlagValue, "childf", "c", -1, "")

	output, err := executeCommand(rootCmd, "child", "-c7", "-p8", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	got := strings.Join(childCmdArgs, " ")
	expected := "one two"
	if got != expected ***REMOVED***
		t.Errorf("childCmdArgs expected: %q, got %q", expected, got)
	***REMOVED***
	if parentFlagValue != 8 ***REMOVED***
		t.Errorf("parentFlagValue expected: %v, got %v", 8, parentFlagValue)
	***REMOVED***
	if childFlagValue != 7 ***REMOVED***
		t.Errorf("childFlagValue expected: %v, got %v", 7, childFlagValue)
	***REMOVED***
***REMOVED***

func TestRequiredFlags(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***
	c.Flags().String("foo1", "", "")
	c.MarkFlagRequired("foo1")
	c.Flags().String("foo2", "", "")
	c.MarkFlagRequired("foo2")
	c.Flags().String("bar", "", "")

	expected := fmt.Sprintf("required flag(s) %q, %q not set", "foo1", "foo2")

	_, err := executeCommand(c)
	got := err.Error()

	if got != expected ***REMOVED***
		t.Errorf("Expected error: %q, got: %q", expected, got)
	***REMOVED***
***REMOVED***

func TestPersistentRequiredFlags(t *testing.T) ***REMOVED***
	parent := &Command***REMOVED***Use: "parent", Run: emptyRun***REMOVED***
	parent.PersistentFlags().String("foo1", "", "")
	parent.MarkPersistentFlagRequired("foo1")
	parent.PersistentFlags().String("foo2", "", "")
	parent.MarkPersistentFlagRequired("foo2")
	parent.Flags().String("foo3", "", "")

	child := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***
	child.Flags().String("bar1", "", "")
	child.MarkFlagRequired("bar1")
	child.Flags().String("bar2", "", "")
	child.MarkFlagRequired("bar2")
	child.Flags().String("bar3", "", "")

	parent.AddCommand(child)

	expected := fmt.Sprintf("required flag(s) %q, %q, %q, %q not set", "bar1", "bar2", "foo1", "foo2")

	_, err := executeCommand(parent, "child")
	if err.Error() != expected ***REMOVED***
		t.Errorf("Expected %q, got %q", expected, err.Error())
	***REMOVED***
***REMOVED***

func TestInitHelpFlagMergesFlags(t *testing.T) ***REMOVED***
	usage := "custom flag"
	rootCmd := &Command***REMOVED***Use: "root"***REMOVED***
	rootCmd.PersistentFlags().Bool("help", false, "custom flag")
	childCmd := &Command***REMOVED***Use: "child"***REMOVED***
	rootCmd.AddCommand(childCmd)

	childCmd.InitDefaultHelpFlag()
	got := childCmd.Flags().Lookup("help").Usage
	if got != usage ***REMOVED***
		t.Errorf("Expected the help flag from the root command with usage: %v\nGot the default with usage: %v", usage, got)
	***REMOVED***
***REMOVED***

func TestHelpCommandExecuted(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Long: "Long description", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(&Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***)

	output, err := executeCommand(rootCmd, "help")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, rootCmd.Long)
***REMOVED***

func TestHelpCommandExecutedOnChild(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Long: "Long description", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	output, err := executeCommand(rootCmd, "help", "child")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, childCmd.Long)
***REMOVED***

func TestSetHelpCommand(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***
	c.AddCommand(&Command***REMOVED***Use: "empty", Run: emptyRun***REMOVED***)

	expected := "WORKS"
	c.SetHelpCommand(&Command***REMOVED***
		Use:   "help [command]",
		Short: "Help about any command",
		Long: `Help provides help for any command in the application.
	Simply type ` + c.Name() + ` help [path to command] for full details.`,
		Run: func(c *Command, _ []string) ***REMOVED*** c.Print(expected) ***REMOVED***,
	***REMOVED***)

	got, err := executeCommand(c, "help")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if got != expected ***REMOVED***
		t.Errorf("Expected to contain %q, got %q", expected, got)
	***REMOVED***
***REMOVED***

func TestHelpFlagExecuted(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Long: "Long description", Run: emptyRun***REMOVED***

	output, err := executeCommand(rootCmd, "--help")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, rootCmd.Long)
***REMOVED***

func TestHelpFlagExecutedOnChild(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Long: "Long description", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	output, err := executeCommand(rootCmd, "child", "--help")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, childCmd.Long)
***REMOVED***

// TestHelpFlagInHelp checks,
// if '--help' flag is shown in help for child (executing `parent help child`),
// that has no other flags.
// Related to https://github.com/spf13/cobra/issues/302.
func TestHelpFlagInHelp(t *testing.T) ***REMOVED***
	parentCmd := &Command***REMOVED***Use: "parent", Run: func(*Command, []string) ***REMOVED******REMOVED******REMOVED***

	childCmd := &Command***REMOVED***Use: "child", Run: func(*Command, []string) ***REMOVED******REMOVED******REMOVED***
	parentCmd.AddCommand(childCmd)

	output, err := executeCommand(parentCmd, "help", "child")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, "[flags]")
***REMOVED***

func TestFlagsInUsage(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Args: NoArgs, Run: func(*Command, []string) ***REMOVED******REMOVED******REMOVED***
	output, err := executeCommand(rootCmd, "--help")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, "[flags]")
***REMOVED***

func TestHelpExecutedOnNonRunnableChild(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Long: "Long description"***REMOVED***
	rootCmd.AddCommand(childCmd)

	output, err := executeCommand(rootCmd, "child")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, childCmd.Long)
***REMOVED***

func TestVersionFlagExecuted(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Version: "1.0.0", Run: emptyRun***REMOVED***

	output, err := executeCommand(rootCmd, "--version", "arg1")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, "root version 1.0.0")
***REMOVED***

func TestVersionTemplate(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Version: "1.0.0", Run: emptyRun***REMOVED***
	rootCmd.SetVersionTemplate(`customized version: ***REMOVED******REMOVED***.Version***REMOVED******REMOVED***`)

	output, err := executeCommand(rootCmd, "--version", "arg1")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, "customized version: 1.0.0")
***REMOVED***

func TestVersionFlagExecutedOnSubcommand(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Version: "1.0.0"***REMOVED***
	rootCmd.AddCommand(&Command***REMOVED***Use: "sub", Run: emptyRun***REMOVED***)

	output, err := executeCommand(rootCmd, "--version", "sub")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, "root version 1.0.0")
***REMOVED***

func TestVersionFlagOnlyAddedToRoot(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Version: "1.0.0", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(&Command***REMOVED***Use: "sub", Run: emptyRun***REMOVED***)

	_, err := executeCommand(rootCmd, "sub", "--version")
	if err == nil ***REMOVED***
		t.Errorf("Expected error")
	***REMOVED***

	checkStringContains(t, err.Error(), "unknown flag: --version")
***REMOVED***

func TestVersionFlagOnlyExistsIfVersionNonEmpty(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***

	_, err := executeCommand(rootCmd, "--version")
	if err == nil ***REMOVED***
		t.Errorf("Expected error")
	***REMOVED***
	checkStringContains(t, err.Error(), "unknown flag: --version")
***REMOVED***

func TestUsageIsNotPrintedTwice(t *testing.T) ***REMOVED***
	var cmd = &Command***REMOVED***Use: "root"***REMOVED***
	var sub = &Command***REMOVED***Use: "sub"***REMOVED***
	cmd.AddCommand(sub)

	output, _ := executeCommand(cmd, "")
	if strings.Count(output, "Usage:") != 1 ***REMOVED***
		t.Error("Usage output is not printed exactly once")
	***REMOVED***
***REMOVED***

func TestVisitParents(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "app"***REMOVED***
	sub := &Command***REMOVED***Use: "sub"***REMOVED***
	dsub := &Command***REMOVED***Use: "dsub"***REMOVED***
	sub.AddCommand(dsub)
	c.AddCommand(sub)

	total := 0
	add := func(x *Command) ***REMOVED***
		total++
	***REMOVED***
	sub.VisitParents(add)
	if total != 1 ***REMOVED***
		t.Errorf("Should have visited 1 parent but visited %d", total)
	***REMOVED***

	total = 0
	dsub.VisitParents(add)
	if total != 2 ***REMOVED***
		t.Errorf("Should have visited 2 parents but visited %d", total)
	***REMOVED***

	total = 0
	c.VisitParents(add)
	if total != 0 ***REMOVED***
		t.Errorf("Should have visited no parents but visited %d", total)
	***REMOVED***
***REMOVED***

func TestSuggestions(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	timesCmd := &Command***REMOVED***
		Use:        "times",
		SuggestFor: []string***REMOVED***"counts"***REMOVED***,
		Run:        emptyRun,
	***REMOVED***
	rootCmd.AddCommand(timesCmd)

	templateWithSuggestions := "Error: unknown command \"%s\" for \"root\"\n\nDid you mean this?\n\t%s\n\nRun 'root --help' for usage.\n"
	templateWithoutSuggestions := "Error: unknown command \"%s\" for \"root\"\nRun 'root --help' for usage.\n"

	tests := map[string]string***REMOVED***
		"time":     "times",
		"tiems":    "times",
		"tims":     "times",
		"timeS":    "times",
		"rimes":    "times",
		"ti":       "times",
		"t":        "times",
		"timely":   "times",
		"ri":       "",
		"timezone": "",
		"foo":      "",
		"counts":   "times",
	***REMOVED***

	for typo, suggestion := range tests ***REMOVED***
		for _, suggestionsDisabled := range []bool***REMOVED***true, false***REMOVED*** ***REMOVED***
			rootCmd.DisableSuggestions = suggestionsDisabled

			var expected string
			output, _ := executeCommand(rootCmd, typo)

			if suggestion == "" || suggestionsDisabled ***REMOVED***
				expected = fmt.Sprintf(templateWithoutSuggestions, typo)
			***REMOVED*** else ***REMOVED***
				expected = fmt.Sprintf(templateWithSuggestions, typo, suggestion)
			***REMOVED***

			if output != expected ***REMOVED***
				t.Errorf("Unexpected response.\nExpected:\n %q\nGot:\n %q\n", expected, output)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRemoveCommand(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Args: NoArgs, Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)
	rootCmd.RemoveCommand(childCmd)

	_, err := executeCommand(rootCmd, "child")
	if err == nil ***REMOVED***
		t.Error("Expected error on calling removed command. Got nil.")
	***REMOVED***
***REMOVED***

func TestReplaceCommandWithRemove(t *testing.T) ***REMOVED***
	childUsed := 0
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	child1Cmd := &Command***REMOVED***
		Use: "child",
		Run: func(*Command, []string) ***REMOVED*** childUsed = 1 ***REMOVED***,
	***REMOVED***
	child2Cmd := &Command***REMOVED***
		Use: "child",
		Run: func(*Command, []string) ***REMOVED*** childUsed = 2 ***REMOVED***,
	***REMOVED***
	rootCmd.AddCommand(child1Cmd)
	rootCmd.RemoveCommand(child1Cmd)
	rootCmd.AddCommand(child2Cmd)

	output, err := executeCommand(rootCmd, "child")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if childUsed == 1 ***REMOVED***
		t.Error("Removed command shouldn't be called")
	***REMOVED***
	if childUsed != 2 ***REMOVED***
		t.Error("Replacing command should have been called but didn't")
	***REMOVED***
***REMOVED***

func TestDeprecatedCommand(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	deprecatedCmd := &Command***REMOVED***
		Use:        "deprecated",
		Deprecated: "This command is deprecated",
		Run:        emptyRun,
	***REMOVED***
	rootCmd.AddCommand(deprecatedCmd)

	output, err := executeCommand(rootCmd, "deprecated")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	checkStringContains(t, output, deprecatedCmd.Deprecated)
***REMOVED***

func TestHooks(t *testing.T) ***REMOVED***
	var (
		persPreArgs  string
		preArgs      string
		runArgs      string
		postArgs     string
		persPostArgs string
	)

	c := &Command***REMOVED***
		Use: "c",
		PersistentPreRun: func(_ *Command, args []string) ***REMOVED***
			persPreArgs = strings.Join(args, " ")
		***REMOVED***,
		PreRun: func(_ *Command, args []string) ***REMOVED***
			preArgs = strings.Join(args, " ")
		***REMOVED***,
		Run: func(_ *Command, args []string) ***REMOVED***
			runArgs = strings.Join(args, " ")
		***REMOVED***,
		PostRun: func(_ *Command, args []string) ***REMOVED***
			postArgs = strings.Join(args, " ")
		***REMOVED***,
		PersistentPostRun: func(_ *Command, args []string) ***REMOVED***
			persPostArgs = strings.Join(args, " ")
		***REMOVED***,
	***REMOVED***

	output, err := executeCommand(c, "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if persPreArgs != "one two" ***REMOVED***
		t.Errorf("Expected persPreArgs %q, got %q", "one two", persPreArgs)
	***REMOVED***
	if preArgs != "one two" ***REMOVED***
		t.Errorf("Expected preArgs %q, got %q", "one two", preArgs)
	***REMOVED***
	if runArgs != "one two" ***REMOVED***
		t.Errorf("Expected runArgs %q, got %q", "one two", runArgs)
	***REMOVED***
	if postArgs != "one two" ***REMOVED***
		t.Errorf("Expected postArgs %q, got %q", "one two", postArgs)
	***REMOVED***
	if persPostArgs != "one two" ***REMOVED***
		t.Errorf("Expected persPostArgs %q, got %q", "one two", persPostArgs)
	***REMOVED***
***REMOVED***

func TestPersistentHooks(t *testing.T) ***REMOVED***
	var (
		parentPersPreArgs  string
		parentPreArgs      string
		parentRunArgs      string
		parentPostArgs     string
		parentPersPostArgs string
	)

	var (
		childPersPreArgs  string
		childPreArgs      string
		childRunArgs      string
		childPostArgs     string
		childPersPostArgs string
	)

	parentCmd := &Command***REMOVED***
		Use: "parent",
		PersistentPreRun: func(_ *Command, args []string) ***REMOVED***
			parentPersPreArgs = strings.Join(args, " ")
		***REMOVED***,
		PreRun: func(_ *Command, args []string) ***REMOVED***
			parentPreArgs = strings.Join(args, " ")
		***REMOVED***,
		Run: func(_ *Command, args []string) ***REMOVED***
			parentRunArgs = strings.Join(args, " ")
		***REMOVED***,
		PostRun: func(_ *Command, args []string) ***REMOVED***
			parentPostArgs = strings.Join(args, " ")
		***REMOVED***,
		PersistentPostRun: func(_ *Command, args []string) ***REMOVED***
			parentPersPostArgs = strings.Join(args, " ")
		***REMOVED***,
	***REMOVED***

	childCmd := &Command***REMOVED***
		Use: "child",
		PersistentPreRun: func(_ *Command, args []string) ***REMOVED***
			childPersPreArgs = strings.Join(args, " ")
		***REMOVED***,
		PreRun: func(_ *Command, args []string) ***REMOVED***
			childPreArgs = strings.Join(args, " ")
		***REMOVED***,
		Run: func(_ *Command, args []string) ***REMOVED***
			childRunArgs = strings.Join(args, " ")
		***REMOVED***,
		PostRun: func(_ *Command, args []string) ***REMOVED***
			childPostArgs = strings.Join(args, " ")
		***REMOVED***,
		PersistentPostRun: func(_ *Command, args []string) ***REMOVED***
			childPersPostArgs = strings.Join(args, " ")
		***REMOVED***,
	***REMOVED***
	parentCmd.AddCommand(childCmd)

	output, err := executeCommand(parentCmd, "child", "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	// TODO: This test fails, but should not.
	// Related to https://github.com/spf13/cobra/issues/252.
	//
	// if parentPersPreArgs != "one two" ***REMOVED***
	// 	t.Errorf("Expected parentPersPreArgs %q, got %q", "one two", parentPersPreArgs)
	// ***REMOVED***
	if parentPreArgs != "" ***REMOVED***
		t.Errorf("Expected blank parentPreArgs, got %q", parentPreArgs)
	***REMOVED***
	if parentRunArgs != "" ***REMOVED***
		t.Errorf("Expected blank parentRunArgs, got %q", parentRunArgs)
	***REMOVED***
	if parentPostArgs != "" ***REMOVED***
		t.Errorf("Expected blank parentPostArgs, got %q", parentPostArgs)
	***REMOVED***
	// TODO: This test fails, but should not.
	// Related to https://github.com/spf13/cobra/issues/252.
	//
	// if parentPersPostArgs != "one two" ***REMOVED***
	// 	t.Errorf("Expected parentPersPostArgs %q, got %q", "one two", parentPersPostArgs)
	// ***REMOVED***

	if childPersPreArgs != "one two" ***REMOVED***
		t.Errorf("Expected childPersPreArgs %q, got %q", "one two", childPersPreArgs)
	***REMOVED***
	if childPreArgs != "one two" ***REMOVED***
		t.Errorf("Expected childPreArgs %q, got %q", "one two", childPreArgs)
	***REMOVED***
	if childRunArgs != "one two" ***REMOVED***
		t.Errorf("Expected childRunArgs %q, got %q", "one two", childRunArgs)
	***REMOVED***
	if childPostArgs != "one two" ***REMOVED***
		t.Errorf("Expected childPostArgs %q, got %q", "one two", childPostArgs)
	***REMOVED***
	if childPersPostArgs != "one two" ***REMOVED***
		t.Errorf("Expected childPersPostArgs %q, got %q", "one two", childPersPostArgs)
	***REMOVED***
***REMOVED***

// Related to https://github.com/spf13/cobra/issues/521.
func TestGlobalNormFuncPropagation(t *testing.T) ***REMOVED***
	normFunc := func(f *pflag.FlagSet, name string) pflag.NormalizedName ***REMOVED***
		return pflag.NormalizedName(name)
	***REMOVED***

	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	rootCmd.SetGlobalNormalizationFunc(normFunc)
	if reflect.ValueOf(normFunc).Pointer() != reflect.ValueOf(rootCmd.GlobalNormalizationFunc()).Pointer() ***REMOVED***
		t.Error("rootCmd seems to have a wrong normalization function")
	***REMOVED***

	if reflect.ValueOf(normFunc).Pointer() != reflect.ValueOf(childCmd.GlobalNormalizationFunc()).Pointer() ***REMOVED***
		t.Error("childCmd should have had the normalization function of rootCmd")
	***REMOVED***
***REMOVED***

// Related to https://github.com/spf13/cobra/issues/521.
func TestNormPassedOnLocal(t *testing.T) ***REMOVED***
	toUpper := func(f *pflag.FlagSet, name string) pflag.NormalizedName ***REMOVED***
		return pflag.NormalizedName(strings.ToUpper(name))
	***REMOVED***

	c := &Command***REMOVED******REMOVED***
	c.Flags().Bool("flagname", true, "this is a dummy flag")
	c.SetGlobalNormalizationFunc(toUpper)
	if c.LocalFlags().Lookup("flagname") != c.LocalFlags().Lookup("FLAGNAME") ***REMOVED***
		t.Error("Normalization function should be passed on to Local flag set")
	***REMOVED***
***REMOVED***

// Related to https://github.com/spf13/cobra/issues/521.
func TestNormPassedOnInherited(t *testing.T) ***REMOVED***
	toUpper := func(f *pflag.FlagSet, name string) pflag.NormalizedName ***REMOVED***
		return pflag.NormalizedName(strings.ToUpper(name))
	***REMOVED***

	c := &Command***REMOVED******REMOVED***
	c.SetGlobalNormalizationFunc(toUpper)

	child1 := &Command***REMOVED******REMOVED***
	c.AddCommand(child1)

	c.PersistentFlags().Bool("flagname", true, "")

	child2 := &Command***REMOVED******REMOVED***
	c.AddCommand(child2)

	inherited := child1.InheritedFlags()
	if inherited.Lookup("flagname") == nil || inherited.Lookup("flagname") != inherited.Lookup("FLAGNAME") ***REMOVED***
		t.Error("Normalization function should be passed on to inherited flag set in command added before flag")
	***REMOVED***

	inherited = child2.InheritedFlags()
	if inherited.Lookup("flagname") == nil || inherited.Lookup("flagname") != inherited.Lookup("FLAGNAME") ***REMOVED***
		t.Error("Normalization function should be passed on to inherited flag set in command added after flag")
	***REMOVED***
***REMOVED***

// Related to https://github.com/spf13/cobra/issues/521.
func TestConsistentNormalizedName(t *testing.T) ***REMOVED***
	toUpper := func(f *pflag.FlagSet, name string) pflag.NormalizedName ***REMOVED***
		return pflag.NormalizedName(strings.ToUpper(name))
	***REMOVED***
	n := func(f *pflag.FlagSet, name string) pflag.NormalizedName ***REMOVED***
		return pflag.NormalizedName(name)
	***REMOVED***

	c := &Command***REMOVED******REMOVED***
	c.Flags().Bool("flagname", true, "")
	c.SetGlobalNormalizationFunc(toUpper)
	c.SetGlobalNormalizationFunc(n)

	if c.LocalFlags().Lookup("flagname") == c.LocalFlags().Lookup("FLAGNAME") ***REMOVED***
		t.Error("Normalizing flag names should not result in duplicate flags")
	***REMOVED***
***REMOVED***

func TestFlagOnPflagCommandLine(t *testing.T) ***REMOVED***
	flagName := "flagOnCommandLine"
	pflag.String(flagName, "", "about my flag")

	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***
	c.AddCommand(&Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***)

	output, _ := executeCommand(c, "--help")
	checkStringContains(t, output, flagName)

	resetCommandLineFlagSet()
***REMOVED***

// TestHiddenCommandExecutes checks,
// if hidden commands run as intended.
func TestHiddenCommandExecutes(t *testing.T) ***REMOVED***
	executed := false
	c := &Command***REMOVED***
		Use:    "c",
		Hidden: true,
		Run:    func(*Command, []string) ***REMOVED*** executed = true ***REMOVED***,
	***REMOVED***

	output, err := executeCommand(c)
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***

	if !executed ***REMOVED***
		t.Error("Hidden command should have been executed")
	***REMOVED***
***REMOVED***

// test to ensure hidden commands do not show up in usage/help text
func TestHiddenCommandIsHidden(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Hidden: true, Run: emptyRun***REMOVED***
	if c.IsAvailableCommand() ***REMOVED***
		t.Errorf("Hidden command should be unavailable")
	***REMOVED***
***REMOVED***

func TestCommandsAreSorted(t *testing.T) ***REMOVED***
	EnableCommandSorting = true

	originalNames := []string***REMOVED***"middle", "zlast", "afirst"***REMOVED***
	expectedNames := []string***REMOVED***"afirst", "middle", "zlast"***REMOVED***

	var rootCmd = &Command***REMOVED***Use: "root"***REMOVED***

	for _, name := range originalNames ***REMOVED***
		rootCmd.AddCommand(&Command***REMOVED***Use: name***REMOVED***)
	***REMOVED***

	for i, c := range rootCmd.Commands() ***REMOVED***
		got := c.Name()
		if expectedNames[i] != got ***REMOVED***
			t.Errorf("Expected: %s, got: %s", expectedNames[i], got)
		***REMOVED***
	***REMOVED***

	EnableCommandSorting = true
***REMOVED***

func TestEnableCommandSortingIsDisabled(t *testing.T) ***REMOVED***
	EnableCommandSorting = false

	originalNames := []string***REMOVED***"middle", "zlast", "afirst"***REMOVED***

	var rootCmd = &Command***REMOVED***Use: "root"***REMOVED***

	for _, name := range originalNames ***REMOVED***
		rootCmd.AddCommand(&Command***REMOVED***Use: name***REMOVED***)
	***REMOVED***

	for i, c := range rootCmd.Commands() ***REMOVED***
		got := c.Name()
		if originalNames[i] != got ***REMOVED***
			t.Errorf("expected: %s, got: %s", originalNames[i], got)
		***REMOVED***
	***REMOVED***

	EnableCommandSorting = true
***REMOVED***

func TestSetOutput(t *testing.T) ***REMOVED***
	c := &Command***REMOVED******REMOVED***
	c.SetOutput(nil)
	if out := c.OutOrStdout(); out != os.Stdout ***REMOVED***
		t.Errorf("Expected setting output to nil to revert back to stdout")
	***REMOVED***
***REMOVED***

func TestFlagErrorFunc(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***

	expectedFmt := "This is expected: %v"
	c.SetFlagErrorFunc(func(_ *Command, err error) error ***REMOVED***
		return fmt.Errorf(expectedFmt, err)
	***REMOVED***)

	_, err := executeCommand(c, "--unknown-flag")

	got := err.Error()
	expected := fmt.Sprintf(expectedFmt, "unknown flag: --unknown-flag")
	if got != expected ***REMOVED***
		t.Errorf("Expected %v, got %v", expected, got)
	***REMOVED***
***REMOVED***

// TestSortedFlags checks,
// if cmd.LocalFlags() is unsorted when cmd.Flags().SortFlags set to false.
// Related to https://github.com/spf13/cobra/issues/404.
func TestSortedFlags(t *testing.T) ***REMOVED***
	c := &Command***REMOVED******REMOVED***
	c.Flags().SortFlags = false
	names := []string***REMOVED***"C", "B", "A", "D"***REMOVED***
	for _, name := range names ***REMOVED***
		c.Flags().Bool(name, false, "")
	***REMOVED***

	i := 0
	c.LocalFlags().VisitAll(func(f *pflag.Flag) ***REMOVED***
		if i == len(names) ***REMOVED***
			return
		***REMOVED***
		if stringInSlice(f.Name, names) ***REMOVED***
			if names[i] != f.Name ***REMOVED***
				t.Errorf("Incorrect order. Expected %v, got %v", names[i], f.Name)
			***REMOVED***
			i++
		***REMOVED***
	***REMOVED***)
***REMOVED***

// TestMergeCommandLineToFlags checks,
// if pflag.CommandLine is correctly merged to c.Flags() after first call
// of c.mergePersistentFlags.
// Related to https://github.com/spf13/cobra/issues/443.
func TestMergeCommandLineToFlags(t *testing.T) ***REMOVED***
	pflag.Bool("boolflag", false, "")
	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***
	c.mergePersistentFlags()
	if c.Flags().Lookup("boolflag") == nil ***REMOVED***
		t.Fatal("Expecting to have flag from CommandLine in c.Flags()")
	***REMOVED***

	resetCommandLineFlagSet()
***REMOVED***

// TestUseDeprecatedFlags checks,
// if cobra.Execute() prints a message, if a deprecated flag is used.
// Related to https://github.com/spf13/cobra/issues/463.
func TestUseDeprecatedFlags(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***
	c.Flags().BoolP("deprecated", "d", false, "deprecated flag")
	c.Flags().MarkDeprecated("deprecated", "This flag is deprecated")

	output, err := executeCommand(c, "c", "-d")
	if err != nil ***REMOVED***
		t.Error("Unexpected error:", err)
	***REMOVED***
	checkStringContains(t, output, "This flag is deprecated")
***REMOVED***

func TestTraverseWithParentFlags(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", TraverseChildren: true***REMOVED***
	rootCmd.Flags().String("str", "", "")
	rootCmd.Flags().BoolP("bool", "b", false, "")

	childCmd := &Command***REMOVED***Use: "child"***REMOVED***
	childCmd.Flags().Int("int", -1, "")

	rootCmd.AddCommand(childCmd)

	c, args, err := rootCmd.Traverse([]string***REMOVED***"-b", "--str", "ok", "child", "--int"***REMOVED***)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
	if len(args) != 1 && args[0] != "--add" ***REMOVED***
		t.Errorf("Wrong args: %v", args)
	***REMOVED***
	if c.Name() != childCmd.Name() ***REMOVED***
		t.Errorf("Expected command: %q, got: %q", childCmd.Name(), c.Name())
	***REMOVED***
***REMOVED***

func TestTraverseNoParentFlags(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", TraverseChildren: true***REMOVED***
	rootCmd.Flags().String("foo", "", "foo things")

	childCmd := &Command***REMOVED***Use: "child"***REMOVED***
	childCmd.Flags().String("str", "", "")
	rootCmd.AddCommand(childCmd)

	c, args, err := rootCmd.Traverse([]string***REMOVED***"child"***REMOVED***)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
	if len(args) != 0 ***REMOVED***
		t.Errorf("Wrong args %v", args)
	***REMOVED***
	if c.Name() != childCmd.Name() ***REMOVED***
		t.Errorf("Expected command: %q, got: %q", childCmd.Name(), c.Name())
	***REMOVED***
***REMOVED***

func TestTraverseWithBadParentFlags(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", TraverseChildren: true***REMOVED***

	childCmd := &Command***REMOVED***Use: "child"***REMOVED***
	childCmd.Flags().String("str", "", "")
	rootCmd.AddCommand(childCmd)

	expected := "unknown flag: --str"

	c, _, err := rootCmd.Traverse([]string***REMOVED***"--str", "ok", "child"***REMOVED***)
	if err == nil || !strings.Contains(err.Error(), expected) ***REMOVED***
		t.Errorf("Expected error, %q, got %q", expected, err)
	***REMOVED***
	if c != nil ***REMOVED***
		t.Errorf("Expected nil command")
	***REMOVED***
***REMOVED***

func TestTraverseWithBadChildFlag(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", TraverseChildren: true***REMOVED***
	rootCmd.Flags().String("str", "", "")

	childCmd := &Command***REMOVED***Use: "child"***REMOVED***
	rootCmd.AddCommand(childCmd)

	// Expect no error because the last commands args shouldn't be parsed in
	// Traverse.
	c, args, err := rootCmd.Traverse([]string***REMOVED***"child", "--str"***REMOVED***)
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
	if len(args) != 1 && args[0] != "--str" ***REMOVED***
		t.Errorf("Wrong args: %v", args)
	***REMOVED***
	if c.Name() != childCmd.Name() ***REMOVED***
		t.Errorf("Expected command %q, got: %q", childCmd.Name(), c.Name())
	***REMOVED***
***REMOVED***

func TestTraverseWithTwoSubcommands(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", TraverseChildren: true***REMOVED***

	subCmd := &Command***REMOVED***Use: "sub", TraverseChildren: true***REMOVED***
	rootCmd.AddCommand(subCmd)

	subsubCmd := &Command***REMOVED***
		Use: "subsub",
	***REMOVED***
	subCmd.AddCommand(subsubCmd)

	c, _, err := rootCmd.Traverse([]string***REMOVED***"sub", "subsub"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %v", err)
	***REMOVED***
	if c.Name() != subsubCmd.Name() ***REMOVED***
		t.Fatalf("Expected command: %q, got %q", subsubCmd.Name(), c.Name())
	***REMOVED***
***REMOVED***

// TestUpdateName checks if c.Name() updates on changed c.Use.
// Related to https://github.com/spf13/cobra/pull/422#discussion_r143918343.
func TestUpdateName(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "name xyz"***REMOVED***
	originalName := c.Name()

	c.Use = "changedName abc"
	if originalName == c.Name() || c.Name() != "changedName" ***REMOVED***
		t.Error("c.Name() should be updated on changed c.Use")
	***REMOVED***
***REMOVED***

type calledAsTestcase struct ***REMOVED***
	args []string
	call string
	want string
	epm  bool
	tc   bool
***REMOVED***

func (tc *calledAsTestcase) test(t *testing.T) ***REMOVED***
	defer func(ov bool) ***REMOVED*** EnablePrefixMatching = ov ***REMOVED***(EnablePrefixMatching)
	EnablePrefixMatching = tc.epm

	var called *Command
	run := func(c *Command, _ []string) ***REMOVED*** t.Logf("called: %q", c.Name()); called = c ***REMOVED***

	parent := &Command***REMOVED***Use: "parent", Run: run***REMOVED***
	child1 := &Command***REMOVED***Use: "child1", Run: run, Aliases: []string***REMOVED***"this"***REMOVED******REMOVED***
	child2 := &Command***REMOVED***Use: "child2", Run: run, Aliases: []string***REMOVED***"that"***REMOVED******REMOVED***

	parent.AddCommand(child1)
	parent.AddCommand(child2)
	parent.SetArgs(tc.args)

	output := new(bytes.Buffer)
	parent.SetOutput(output)

	parent.Execute()

	if called == nil ***REMOVED***
		if tc.call != "" ***REMOVED***
			t.Errorf("missing expected call to command: %s", tc.call)
		***REMOVED***
		return
	***REMOVED***

	if called.Name() != tc.call ***REMOVED***
		t.Errorf("called command == %q; Wanted %q", called.Name(), tc.call)
	***REMOVED*** else if got := called.CalledAs(); got != tc.want ***REMOVED***
		t.Errorf("%s.CalledAs() == %q; Wanted: %q", tc.call, got, tc.want)
	***REMOVED***
***REMOVED***

func TestCalledAs(t *testing.T) ***REMOVED***
	tests := map[string]calledAsTestcase***REMOVED***
		"find/no-args":            ***REMOVED***nil, "parent", "parent", false, false***REMOVED***,
		"find/real-name":          ***REMOVED***[]string***REMOVED***"child1"***REMOVED***, "child1", "child1", false, false***REMOVED***,
		"find/full-alias":         ***REMOVED***[]string***REMOVED***"that"***REMOVED***, "child2", "that", false, false***REMOVED***,
		"find/part-no-prefix":     ***REMOVED***[]string***REMOVED***"thi"***REMOVED***, "", "", false, false***REMOVED***,
		"find/part-alias":         ***REMOVED***[]string***REMOVED***"thi"***REMOVED***, "child1", "this", true, false***REMOVED***,
		"find/conflict":           ***REMOVED***[]string***REMOVED***"th"***REMOVED***, "", "", true, false***REMOVED***,
		"traverse/no-args":        ***REMOVED***nil, "parent", "parent", false, true***REMOVED***,
		"traverse/real-name":      ***REMOVED***[]string***REMOVED***"child1"***REMOVED***, "child1", "child1", false, true***REMOVED***,
		"traverse/full-alias":     ***REMOVED***[]string***REMOVED***"that"***REMOVED***, "child2", "that", false, true***REMOVED***,
		"traverse/part-no-prefix": ***REMOVED***[]string***REMOVED***"thi"***REMOVED***, "", "", false, true***REMOVED***,
		"traverse/part-alias":     ***REMOVED***[]string***REMOVED***"thi"***REMOVED***, "child1", "this", true, true***REMOVED***,
		"traverse/conflict":       ***REMOVED***[]string***REMOVED***"th"***REMOVED***, "", "", true, true***REMOVED***,
	***REMOVED***

	for name, tc := range tests ***REMOVED***
		t.Run(name, tc.test)
	***REMOVED***
***REMOVED***
