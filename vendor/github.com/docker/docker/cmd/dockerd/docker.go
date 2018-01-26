package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/docker/pkg/term"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newDaemonCommand() *cobra.Command ***REMOVED***
	opts := newDaemonOptions(config.New())

	cmd := &cobra.Command***REMOVED***
		Use:           "dockerd [OPTIONS]",
		Short:         "A self-sufficient runtime for containers.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			opts.flags = cmd.Flags()
			return runDaemon(opts)
		***REMOVED***,
	***REMOVED***
	cli.SetupRootCommand(cmd)

	flags := cmd.Flags()
	flags.BoolVarP(&opts.version, "version", "v", false, "Print version information and quit")
	flags.StringVar(&opts.configFile, "config-file", defaultDaemonConfigFile, "Daemon configuration file")
	opts.InstallFlags(flags)
	installConfigFlags(opts.daemonConfig, flags)
	installServiceFlags(flags)

	return cmd
***REMOVED***

func runDaemon(opts *daemonOptions) error ***REMOVED***
	if opts.version ***REMOVED***
		showVersion()
		return nil
	***REMOVED***

	daemonCli := NewDaemonCli()

	// Windows specific settings as these are not defaulted.
	if runtime.GOOS == "windows" ***REMOVED***
		if opts.daemonConfig.Pidfile == "" ***REMOVED***
			opts.daemonConfig.Pidfile = filepath.Join(opts.daemonConfig.Root, "docker.pid")
		***REMOVED***
		if opts.configFile == "" ***REMOVED***
			opts.configFile = filepath.Join(opts.daemonConfig.Root, `config\daemon.json`)
		***REMOVED***
	***REMOVED***

	// On Windows, this may be launching as a service or with an option to
	// register the service.
	stop, runAsService, err := initService(daemonCli)
	if err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***

	if stop ***REMOVED***
		return nil
	***REMOVED***

	// If Windows SCM manages the service - no need for PID files
	if runAsService ***REMOVED***
		opts.daemonConfig.Pidfile = ""
	***REMOVED***

	err = daemonCli.start(opts)
	notifyShutdown(err)
	return err
***REMOVED***

func showVersion() ***REMOVED***
	fmt.Printf("Docker version %s, build %s\n", dockerversion.Version, dockerversion.GitCommit)
***REMOVED***

func main() ***REMOVED***
	if reexec.Init() ***REMOVED***
		return
	***REMOVED***

	// Set terminal emulation based on platform as required.
	_, stdout, stderr := term.StdStreams()

	// @jhowardmsft - maybe there is a historic reason why on non-Windows, stderr is used
	// here. However, on Windows it makes no sense and there is no need.
	if runtime.GOOS == "windows" ***REMOVED***
		logrus.SetOutput(stdout)
	***REMOVED*** else ***REMOVED***
		logrus.SetOutput(stderr)
	***REMOVED***

	cmd := newDaemonCommand()
	cmd.SetOutput(stdout)
	if err := cmd.Execute(); err != nil ***REMOVED***
		fmt.Fprintf(stderr, "%s\n", err)
		os.Exit(1)
	***REMOVED***
***REMOVED***
