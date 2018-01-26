package cli

import (
	"fmt"
	"strings"

	"github.com/docker/docker/pkg/term"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// SetupRootCommand sets default usage, help, and error handling for the
// root command.
func SetupRootCommand(rootCmd *cobra.Command) ***REMOVED***
	cobra.AddTemplateFunc("hasSubCommands", hasSubCommands)
	cobra.AddTemplateFunc("hasManagementSubCommands", hasManagementSubCommands)
	cobra.AddTemplateFunc("operationSubCommands", operationSubCommands)
	cobra.AddTemplateFunc("managementSubCommands", managementSubCommands)
	cobra.AddTemplateFunc("wrappedFlagUsages", wrappedFlagUsages)

	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetHelpTemplate(helpTemplate)
	rootCmd.SetFlagErrorFunc(FlagErrorFunc)
	rootCmd.SetHelpCommand(helpCommand)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	rootCmd.PersistentFlags().MarkShorthandDeprecated("help", "please use --help")
***REMOVED***

// FlagErrorFunc prints an error message which matches the format of the
// docker/docker/cli error messages
func FlagErrorFunc(cmd *cobra.Command, err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***

	usage := ""
	if cmd.HasSubCommands() ***REMOVED***
		usage = "\n\n" + cmd.UsageString()
	***REMOVED***
	return StatusError***REMOVED***
		Status:     fmt.Sprintf("%s\nSee '%s --help'.%s", err, cmd.CommandPath(), usage),
		StatusCode: 125,
	***REMOVED***
***REMOVED***

var helpCommand = &cobra.Command***REMOVED***
	Use:               "help [command]",
	Short:             "Help about the command",
	PersistentPreRun:  func(cmd *cobra.Command, args []string) ***REMOVED******REMOVED***,
	PersistentPostRun: func(cmd *cobra.Command, args []string) ***REMOVED******REMOVED***,
	RunE: func(c *cobra.Command, args []string) error ***REMOVED***
		cmd, args, e := c.Root().Find(args)
		if cmd == nil || e != nil || len(args) > 0 ***REMOVED***
			return errors.Errorf("unknown help topic: %v", strings.Join(args, " "))
		***REMOVED***

		helpFunc := cmd.HelpFunc()
		helpFunc(cmd, args)
		return nil
	***REMOVED***,
***REMOVED***

func hasSubCommands(cmd *cobra.Command) bool ***REMOVED***
	return len(operationSubCommands(cmd)) > 0
***REMOVED***

func hasManagementSubCommands(cmd *cobra.Command) bool ***REMOVED***
	return len(managementSubCommands(cmd)) > 0
***REMOVED***

func operationSubCommands(cmd *cobra.Command) []*cobra.Command ***REMOVED***
	cmds := []*cobra.Command***REMOVED******REMOVED***
	for _, sub := range cmd.Commands() ***REMOVED***
		if sub.IsAvailableCommand() && !sub.HasSubCommands() ***REMOVED***
			cmds = append(cmds, sub)
		***REMOVED***
	***REMOVED***
	return cmds
***REMOVED***

func wrappedFlagUsages(cmd *cobra.Command) string ***REMOVED***
	width := 80
	if ws, err := term.GetWinsize(0); err == nil ***REMOVED***
		width = int(ws.Width)
	***REMOVED***
	return cmd.Flags().FlagUsagesWrapped(width - 1)
***REMOVED***

func managementSubCommands(cmd *cobra.Command) []*cobra.Command ***REMOVED***
	cmds := []*cobra.Command***REMOVED******REMOVED***
	for _, sub := range cmd.Commands() ***REMOVED***
		if sub.IsAvailableCommand() && sub.HasSubCommands() ***REMOVED***
			cmds = append(cmds, sub)
		***REMOVED***
	***REMOVED***
	return cmds
***REMOVED***

var usageTemplate = `Usage:

***REMOVED******REMOVED***- if not .HasSubCommands***REMOVED******REMOVED***	***REMOVED******REMOVED***.UseLine***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
***REMOVED******REMOVED***- if .HasSubCommands***REMOVED******REMOVED***	***REMOVED******REMOVED*** .CommandPath***REMOVED******REMOVED*** COMMAND***REMOVED******REMOVED***end***REMOVED******REMOVED***

***REMOVED******REMOVED*** .Short | trim ***REMOVED******REMOVED***

***REMOVED******REMOVED***- if gt .Aliases 0***REMOVED******REMOVED***

Aliases:
  ***REMOVED******REMOVED***.NameAndAliases***REMOVED******REMOVED***

***REMOVED******REMOVED***- end***REMOVED******REMOVED***
***REMOVED******REMOVED***- if .HasExample***REMOVED******REMOVED***

Examples:
***REMOVED******REMOVED*** .Example ***REMOVED******REMOVED***

***REMOVED******REMOVED***- end***REMOVED******REMOVED***
***REMOVED******REMOVED***- if .HasFlags***REMOVED******REMOVED***

Options:
***REMOVED******REMOVED*** wrappedFlagUsages . | trimRightSpace***REMOVED******REMOVED***

***REMOVED******REMOVED***- end***REMOVED******REMOVED***
***REMOVED******REMOVED***- if hasManagementSubCommands . ***REMOVED******REMOVED***

Management Commands:

***REMOVED******REMOVED***- range managementSubCommands . ***REMOVED******REMOVED***
  ***REMOVED******REMOVED***rpad .Name .NamePadding ***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Short***REMOVED******REMOVED***
***REMOVED******REMOVED***- end***REMOVED******REMOVED***

***REMOVED******REMOVED***- end***REMOVED******REMOVED***
***REMOVED******REMOVED***- if hasSubCommands .***REMOVED******REMOVED***

Commands:

***REMOVED******REMOVED***- range operationSubCommands . ***REMOVED******REMOVED***
  ***REMOVED******REMOVED***rpad .Name .NamePadding ***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Short***REMOVED******REMOVED***
***REMOVED******REMOVED***- end***REMOVED******REMOVED***
***REMOVED******REMOVED***- end***REMOVED******REMOVED***

***REMOVED******REMOVED***- if .HasSubCommands ***REMOVED******REMOVED***

Run '***REMOVED******REMOVED***.CommandPath***REMOVED******REMOVED*** COMMAND --help' for more information on a command.
***REMOVED******REMOVED***- end***REMOVED******REMOVED***
`

var helpTemplate = `
***REMOVED******REMOVED***if or .Runnable .HasSubCommands***REMOVED******REMOVED******REMOVED******REMOVED***.UsageString***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***`
