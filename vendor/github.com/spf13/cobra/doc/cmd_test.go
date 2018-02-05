package doc

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func emptyRun(*cobra.Command, []string) ***REMOVED******REMOVED***

func init() ***REMOVED***
	rootCmd.PersistentFlags().StringP("rootflag", "r", "two", "")
	rootCmd.PersistentFlags().StringP("strtwo", "t", "two", "help message for parent flag strtwo")

	echoCmd.PersistentFlags().StringP("strone", "s", "one", "help message for flag strone")
	echoCmd.PersistentFlags().BoolP("persistentbool", "p", false, "help message for flag persistentbool")
	echoCmd.Flags().IntP("intone", "i", 123, "help message for flag intone")
	echoCmd.Flags().BoolP("boolone", "b", true, "help message for flag boolone")

	timesCmd.PersistentFlags().StringP("strtwo", "t", "2", "help message for child flag strtwo")
	timesCmd.Flags().IntP("inttwo", "j", 234, "help message for flag inttwo")
	timesCmd.Flags().BoolP("booltwo", "c", false, "help message for flag booltwo")

	printCmd.PersistentFlags().StringP("strthree", "s", "three", "help message for flag strthree")
	printCmd.Flags().IntP("intthree", "i", 345, "help message for flag intthree")
	printCmd.Flags().BoolP("boolthree", "b", true, "help message for flag boolthree")

	echoCmd.AddCommand(timesCmd, echoSubCmd, deprecatedCmd)
	rootCmd.AddCommand(printCmd, echoCmd)
***REMOVED***

var rootCmd = &cobra.Command***REMOVED***
	Use:   "root",
	Short: "Root short description",
	Long:  "Root long description",
	Run:   emptyRun,
***REMOVED***

var echoCmd = &cobra.Command***REMOVED***
	Use:     "echo [string to echo]",
	Aliases: []string***REMOVED***"say"***REMOVED***,
	Short:   "Echo anything to the screen",
	Long:    "an utterly useless command for testing",
	Example: "Just run cobra-test echo",
***REMOVED***

var echoSubCmd = &cobra.Command***REMOVED***
	Use:   "echosub [string to print]",
	Short: "second sub command for echo",
	Long:  "an absolutely utterly useless command for testing gendocs!.",
	Run:   emptyRun,
***REMOVED***

var timesCmd = &cobra.Command***REMOVED***
	Use:        "times [# times] [string to echo]",
	SuggestFor: []string***REMOVED***"counts"***REMOVED***,
	Short:      "Echo anything to the screen more times",
	Long:       `a slightly useless command for testing.`,
	Run:        emptyRun,
***REMOVED***

var deprecatedCmd = &cobra.Command***REMOVED***
	Use:        "deprecated [can't do anything here]",
	Short:      "A command which is deprecated",
	Long:       `an absolutely utterly useless command for testing deprecation!.`,
	Deprecated: "Please use echo instead",
***REMOVED***

var printCmd = &cobra.Command***REMOVED***
	Use:   "print [string to print]",
	Short: "Print anything to the screen",
	Long:  `an absolutely utterly useless command for testing.`,
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
