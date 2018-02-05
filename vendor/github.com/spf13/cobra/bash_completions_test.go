package cobra

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func checkOmit(t *testing.T, found, unexpected string) ***REMOVED***
	if strings.Contains(found, unexpected) ***REMOVED***
		t.Errorf("Got: %q\nBut should not have!\n", unexpected)
	***REMOVED***
***REMOVED***

func check(t *testing.T, found, expected string) ***REMOVED***
	if !strings.Contains(found, expected) ***REMOVED***
		t.Errorf("Expecting to contain: \n %q\nGot:\n %q\n", expected, found)
	***REMOVED***
***REMOVED***

func runShellCheck(s string) error ***REMOVED***
	excluded := []string***REMOVED***
		"SC2034", // PREFIX appears unused. Verify it or export it.
	***REMOVED***
	cmd := exec.Command("shellcheck", "-s", "bash", "-", "-e", strings.Join(excluded, ","))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	stdin, err := cmd.StdinPipe()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	go func() ***REMOVED***
		stdin.Write([]byte(s))
		stdin.Close()
	***REMOVED***()

	return cmd.Run()
***REMOVED***

// World worst custom function, just keep telling you to enter hello!
const bashCompletionFunc = `__custom_func() ***REMOVED***
	COMPREPLY=( "hello" )
***REMOVED***
`

func TestBashCompletions(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***
		Use:                    "root",
		ArgAliases:             []string***REMOVED***"pods", "nodes", "services", "replicationcontrollers", "po", "no", "svc", "rc"***REMOVED***,
		ValidArgs:              []string***REMOVED***"pod", "node", "service", "replicationcontroller"***REMOVED***,
		BashCompletionFunction: bashCompletionFunc,
		Run: emptyRun,
	***REMOVED***
	rootCmd.Flags().IntP("introot", "i", -1, "help message for flag introot")
	rootCmd.MarkFlagRequired("introot")

	// Filename.
	rootCmd.Flags().String("filename", "", "Enter a filename")
	rootCmd.MarkFlagFilename("filename", "json", "yaml", "yml")

	// Persistent filename.
	rootCmd.PersistentFlags().String("persistent-filename", "", "Enter a filename")
	rootCmd.MarkPersistentFlagFilename("persistent-filename")
	rootCmd.MarkPersistentFlagRequired("persistent-filename")

	// Filename extensions.
	rootCmd.Flags().String("filename-ext", "", "Enter a filename (extension limited)")
	rootCmd.MarkFlagFilename("filename-ext")
	rootCmd.Flags().String("custom", "", "Enter a filename (extension limited)")
	rootCmd.MarkFlagCustom("custom", "__complete_custom")

	// Subdirectories in a given directory.
	rootCmd.Flags().String("theme", "", "theme to use (located in /themes/THEMENAME/)")
	rootCmd.Flags().SetAnnotation("theme", BashCompSubdirsInDir, []string***REMOVED***"themes"***REMOVED***)

	echoCmd := &Command***REMOVED***
		Use:     "echo [string to echo]",
		Aliases: []string***REMOVED***"say"***REMOVED***,
		Short:   "Echo anything to the screen",
		Long:    "an utterly useless command for testing.",
		Example: "Just run cobra-test echo",
		Run:     emptyRun,
	***REMOVED***

	printCmd := &Command***REMOVED***
		Use:   "print [string to print]",
		Args:  MinimumNArgs(1),
		Short: "Print anything to the screen",
		Long:  "an absolutely utterly useless command for testing.",
		Run:   emptyRun,
	***REMOVED***

	deprecatedCmd := &Command***REMOVED***
		Use:        "deprecated [can't do anything here]",
		Args:       NoArgs,
		Short:      "A command which is deprecated",
		Long:       "an absolutely utterly useless command for testing deprecation!.",
		Deprecated: "Please use echo instead",
		Run:        emptyRun,
	***REMOVED***

	colonCmd := &Command***REMOVED***
		Use: "cmd:colon",
		Run: emptyRun,
	***REMOVED***

	timesCmd := &Command***REMOVED***
		Use:        "times [# times] [string to echo]",
		SuggestFor: []string***REMOVED***"counts"***REMOVED***,
		Args:       OnlyValidArgs,
		ValidArgs:  []string***REMOVED***"one", "two", "three", "four"***REMOVED***,
		Short:      "Echo anything to the screen more times",
		Long:       "a slightly useless command for testing.",
		Run:        emptyRun,
	***REMOVED***

	echoCmd.AddCommand(timesCmd)
	rootCmd.AddCommand(echoCmd, printCmd, deprecatedCmd, colonCmd)

	buf := new(bytes.Buffer)
	rootCmd.GenBashCompletion(buf)
	output := buf.String()

	check(t, output, "_root")
	check(t, output, "_root_echo")
	check(t, output, "_root_echo_times")
	check(t, output, "_root_print")
	check(t, output, "_root_cmd__colon")

	// check for required flags
	check(t, output, `must_have_one_flag+=("--introot=")`)
	check(t, output, `must_have_one_flag+=("--persistent-filename=")`)
	// check for custom completion function
	check(t, output, `COMPREPLY=( "hello" )`)
	// check for required nouns
	check(t, output, `must_have_one_noun+=("pod")`)
	// check for noun aliases
	check(t, output, `noun_aliases+=("pods")`)
	check(t, output, `noun_aliases+=("rc")`)
	checkOmit(t, output, `must_have_one_noun+=("pods")`)
	// check for filename extension flags
	check(t, output, `flags_completion+=("_filedir")`)
	// check for filename extension flags
	check(t, output, `must_have_one_noun+=("three")`)
	// check for filename extension flags
	check(t, output, `flags_completion+=("__handle_filename_extension_flag json|yaml|yml")`)
	// check for custom flags
	check(t, output, `flags_completion+=("__complete_custom")`)
	// check for subdirs_in_dir flags
	check(t, output, `flags_completion+=("__handle_subdirs_in_dir_flag themes")`)

	checkOmit(t, output, deprecatedCmd.Name())

	// If available, run shellcheck against the script.
	if err := exec.Command("which", "shellcheck").Run(); err != nil ***REMOVED***
		return
	***REMOVED***
	if err := runShellCheck(output); err != nil ***REMOVED***
		t.Fatalf("shellcheck failed: %v", err)
	***REMOVED***
***REMOVED***

func TestBashCompletionHiddenFlag(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***

	const flagName = "hiddenFlag"
	c.Flags().Bool(flagName, false, "")
	c.Flags().MarkHidden(flagName)

	buf := new(bytes.Buffer)
	c.GenBashCompletion(buf)
	output := buf.String()

	if strings.Contains(output, flagName) ***REMOVED***
		t.Errorf("Expected completion to not include %q flag: Got %v", flagName, output)
	***REMOVED***
***REMOVED***

func TestBashCompletionDeprecatedFlag(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Run: emptyRun***REMOVED***

	const flagName = "deprecated-flag"
	c.Flags().Bool(flagName, false, "")
	c.Flags().MarkDeprecated(flagName, "use --not-deprecated instead")

	buf := new(bytes.Buffer)
	c.GenBashCompletion(buf)
	output := buf.String()

	if strings.Contains(output, flagName) ***REMOVED***
		t.Errorf("expected completion to not include %q flag: Got %v", flagName, output)
	***REMOVED***
***REMOVED***
