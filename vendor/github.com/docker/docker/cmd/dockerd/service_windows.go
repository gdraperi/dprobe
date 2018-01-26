package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	flServiceName       *string
	flRegisterService   *bool
	flUnregisterService *bool
	flRunService        *bool

	setStdHandle = windows.NewLazySystemDLL("kernel32.dll").NewProc("SetStdHandle")
	oldStderr    windows.Handle
	panicFile    *os.File

	service *handler
)

const (
	// These should match the values in event_messages.mc.
	eventInfo  = 1
	eventWarn  = 1
	eventError = 1
	eventDebug = 2
	eventPanic = 3
	eventFatal = 4

	eventExtraOffset = 10 // Add this to any event to get a string that supports extended data
)

func installServiceFlags(flags *pflag.FlagSet) ***REMOVED***
	flServiceName = flags.String("service-name", "docker", "Set the Windows service name")
	flRegisterService = flags.Bool("register-service", false, "Register the service and exit")
	flUnregisterService = flags.Bool("unregister-service", false, "Unregister the service and exit")
	flRunService = flags.Bool("run-service", false, "")
	flags.MarkHidden("run-service")
***REMOVED***

type handler struct ***REMOVED***
	tosvc     chan bool
	fromsvc   chan error
	daemonCli *DaemonCli
***REMOVED***

type etwHook struct ***REMOVED***
	log *eventlog.Log
***REMOVED***

func (h *etwHook) Levels() []logrus.Level ***REMOVED***
	return []logrus.Level***REMOVED***
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	***REMOVED***
***REMOVED***

func (h *etwHook) Fire(e *logrus.Entry) error ***REMOVED***
	var (
		etype uint16
		eid   uint32
	)

	switch e.Level ***REMOVED***
	case logrus.PanicLevel:
		etype = windows.EVENTLOG_ERROR_TYPE
		eid = eventPanic
	case logrus.FatalLevel:
		etype = windows.EVENTLOG_ERROR_TYPE
		eid = eventFatal
	case logrus.ErrorLevel:
		etype = windows.EVENTLOG_ERROR_TYPE
		eid = eventError
	case logrus.WarnLevel:
		etype = windows.EVENTLOG_WARNING_TYPE
		eid = eventWarn
	case logrus.InfoLevel:
		etype = windows.EVENTLOG_INFORMATION_TYPE
		eid = eventInfo
	case logrus.DebugLevel:
		etype = windows.EVENTLOG_INFORMATION_TYPE
		eid = eventDebug
	default:
		return errors.New("unknown level")
	***REMOVED***

	// If there is additional data, include it as a second string.
	exts := ""
	if len(e.Data) > 0 ***REMOVED***
		fs := bytes.Buffer***REMOVED******REMOVED***
		for k, v := range e.Data ***REMOVED***
			fs.WriteString(k)
			fs.WriteByte('=')
			fmt.Fprint(&fs, v)
			fs.WriteByte(' ')
		***REMOVED***

		exts = fs.String()[:fs.Len()-1]
		eid += eventExtraOffset
	***REMOVED***

	if h.log == nil ***REMOVED***
		fmt.Fprintf(os.Stderr, "%s [%s]\n", e.Message, exts)
		return nil
	***REMOVED***

	var (
		ss  [2]*uint16
		err error
	)

	ss[0], err = windows.UTF16PtrFromString(e.Message)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	count := uint16(1)
	if exts != "" ***REMOVED***
		ss[1], err = windows.UTF16PtrFromString(exts)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		count++
	***REMOVED***

	return windows.ReportEvent(h.log.Handle, etype, 0, eid, 0, count, 0, &ss[0], nil)
***REMOVED***

func getServicePath() (string, error) ***REMOVED***
	p, err := exec.LookPath(os.Args[0])
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return filepath.Abs(p)
***REMOVED***

func registerService() error ***REMOVED***
	p, err := getServicePath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer m.Disconnect()

	depends := []string***REMOVED******REMOVED***

	// This dependency is required on build 14393 (RS1)
	// it is added to the platform in newer builds
	if system.GetOSVersion().Build == 14393 ***REMOVED***
		depends = append(depends, "ConDrv")
	***REMOVED***

	c := mgr.Config***REMOVED***
		ServiceType:  windows.SERVICE_WIN32_OWN_PROCESS,
		StartType:    mgr.StartAutomatic,
		ErrorControl: mgr.ErrorNormal,
		Dependencies: depends,
		DisplayName:  "Docker Engine",
	***REMOVED***

	// Configure the service to launch with the arguments that were just passed.
	args := []string***REMOVED***"--run-service"***REMOVED***
	for _, a := range os.Args[1:] ***REMOVED***
		if a != "--register-service" && a != "--unregister-service" ***REMOVED***
			args = append(args, a)
		***REMOVED***
	***REMOVED***

	s, err := m.CreateService(*flServiceName, p, c, args...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer s.Close()

	// See http://stackoverflow.com/questions/35151052/how-do-i-configure-failure-actions-of-a-windows-service-written-in-go
	const (
		scActionNone       = 0
		scActionRestart    = 1
		scActionReboot     = 2
		scActionRunCommand = 3

		serviceConfigFailureActions = 2
	)

	type serviceFailureActions struct ***REMOVED***
		ResetPeriod  uint32
		RebootMsg    *uint16
		Command      *uint16
		ActionsCount uint32
		Actions      uintptr
	***REMOVED***

	type scAction struct ***REMOVED***
		Type  uint32
		Delay uint32
	***REMOVED***
	t := []scAction***REMOVED***
		***REMOVED***Type: scActionRestart, Delay: uint32(60 * time.Second / time.Millisecond)***REMOVED***,
		***REMOVED***Type: scActionRestart, Delay: uint32(60 * time.Second / time.Millisecond)***REMOVED***,
		***REMOVED***Type: scActionNone***REMOVED***,
	***REMOVED***
	lpInfo := serviceFailureActions***REMOVED***ResetPeriod: uint32(24 * time.Hour / time.Second), ActionsCount: uint32(3), Actions: uintptr(unsafe.Pointer(&t[0]))***REMOVED***
	err = windows.ChangeServiceConfig2(s.Handle, serviceConfigFailureActions, (*byte)(unsafe.Pointer(&lpInfo)))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return eventlog.Install(*flServiceName, p, false, eventlog.Info|eventlog.Warning|eventlog.Error)
***REMOVED***

func unregisterService() error ***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer m.Disconnect()

	s, err := m.OpenService(*flServiceName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer s.Close()

	eventlog.Remove(*flServiceName)
	err = s.Delete()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// initService is the entry point for running the daemon as a Windows
// service. It returns an indication to stop (if registering/un-registering);
// an indication of whether it is running as a service; and an error.
func initService(daemonCli *DaemonCli) (bool, bool, error) ***REMOVED***
	if *flUnregisterService ***REMOVED***
		if *flRegisterService ***REMOVED***
			return true, false, errors.New("--register-service and --unregister-service cannot be used together")
		***REMOVED***
		return true, false, unregisterService()
	***REMOVED***

	if *flRegisterService ***REMOVED***
		return true, false, registerService()
	***REMOVED***

	if !*flRunService ***REMOVED***
		return false, false, nil
	***REMOVED***

	interactive, err := svc.IsAnInteractiveSession()
	if err != nil ***REMOVED***
		return false, false, err
	***REMOVED***

	h := &handler***REMOVED***
		tosvc:     make(chan bool),
		fromsvc:   make(chan error),
		daemonCli: daemonCli,
	***REMOVED***

	var log *eventlog.Log
	if !interactive ***REMOVED***
		log, err = eventlog.Open(*flServiceName)
		if err != nil ***REMOVED***
			return false, false, err
		***REMOVED***
	***REMOVED***

	logrus.AddHook(&etwHook***REMOVED***log***REMOVED***)
	logrus.SetOutput(ioutil.Discard)

	service = h
	go func() ***REMOVED***
		if interactive ***REMOVED***
			err = debug.Run(*flServiceName, h)
		***REMOVED*** else ***REMOVED***
			err = svc.Run(*flServiceName, h)
		***REMOVED***

		h.fromsvc <- err
	***REMOVED***()

	// Wait for the first signal from the service handler.
	err = <-h.fromsvc
	if err != nil ***REMOVED***
		return false, false, err
	***REMOVED***
	return false, true, nil
***REMOVED***

func (h *handler) started() error ***REMOVED***
	// This must be delayed until daemonCli initializes Config.Root
	err := initPanicFile(filepath.Join(h.daemonCli.Config.Root, "panic.log"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	h.tosvc <- false
	return nil
***REMOVED***

func (h *handler) stopped(err error) ***REMOVED***
	logrus.Debugf("Stopping service: %v", err)
	h.tosvc <- err != nil
	<-h.fromsvc
***REMOVED***

func (h *handler) Execute(_ []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) ***REMOVED***
	s <- svc.Status***REMOVED***State: svc.StartPending, Accepts: 0***REMOVED***
	// Unblock initService()
	h.fromsvc <- nil

	// Wait for initialization to complete.
	failed := <-h.tosvc
	if failed ***REMOVED***
		logrus.Debug("Aborting service start due to failure during initialization")
		return true, 1
	***REMOVED***

	s <- svc.Status***REMOVED***State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown | svc.Accepted(windows.SERVICE_ACCEPT_PARAMCHANGE)***REMOVED***
	logrus.Debug("Service running")
Loop:
	for ***REMOVED***
		select ***REMOVED***
		case failed = <-h.tosvc:
			break Loop
		case c := <-r:
			switch c.Cmd ***REMOVED***
			case svc.Cmd(windows.SERVICE_CONTROL_PARAMCHANGE):
				h.daemonCli.reloadConfig()
			case svc.Interrogate:
				s <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s <- svc.Status***REMOVED***State: svc.StopPending, Accepts: 0***REMOVED***
				h.daemonCli.stop()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	removePanicFile()
	if failed ***REMOVED***
		return true, 1
	***REMOVED***
	return false, 0
***REMOVED***

func initPanicFile(path string) error ***REMOVED***
	var err error
	panicFile, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	st, err := panicFile.Stat()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If there are contents in the file already, move the file out of the way
	// and replace it.
	if st.Size() > 0 ***REMOVED***
		panicFile.Close()
		os.Rename(path, path+".old")
		panicFile, err = os.Create(path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Update STD_ERROR_HANDLE to point to the panic file so that Go writes to
	// it when it panics. Remember the old stderr to restore it before removing
	// the panic file.
	sh := windows.STD_ERROR_HANDLE
	h, err := windows.GetStdHandle(uint32(sh))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	oldStderr = h

	r, _, err := setStdHandle.Call(uintptr(sh), uintptr(panicFile.Fd()))
	if r == 0 && err != nil ***REMOVED***
		return err
	***REMOVED***

	// Reset os.Stderr to the panic file (so fmt.Fprintf(os.Stderr,...) actually gets redirected)
	os.Stderr = os.NewFile(uintptr(panicFile.Fd()), "/dev/stderr")

	// Force threads that panic to write to stderr (the panicFile handle now), otherwise it will go into the ether
	log.SetOutput(os.Stderr)

	return nil
***REMOVED***

func removePanicFile() ***REMOVED***
	if st, err := panicFile.Stat(); err == nil ***REMOVED***
		if st.Size() == 0 ***REMOVED***
			sh := windows.STD_ERROR_HANDLE
			setStdHandle.Call(uintptr(sh), uintptr(oldStderr))
			panicFile.Close()
			os.Remove(panicFile.Name())
		***REMOVED***
	***REMOVED***
***REMOVED***
