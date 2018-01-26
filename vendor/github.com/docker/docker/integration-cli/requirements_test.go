package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/requirement"
)

func ArchitectureIsNot(arch string) bool ***REMOVED***
	return os.Getenv("DOCKER_ENGINE_GOARCH") != arch
***REMOVED***

func DaemonIsWindows() bool ***REMOVED***
	return testEnv.OSType == "windows"
***REMOVED***

func DaemonIsWindowsAtLeastBuild(buildNumber int) func() bool ***REMOVED***
	return func() bool ***REMOVED***
		if testEnv.OSType != "windows" ***REMOVED***
			return false
		***REMOVED***
		version := testEnv.DaemonInfo.KernelVersion
		numVersion, _ := strconv.Atoi(strings.Split(version, " ")[1])
		return numVersion >= buildNumber
	***REMOVED***
***REMOVED***

func DaemonIsLinux() bool ***REMOVED***
	return testEnv.OSType == "linux"
***REMOVED***

func OnlyDefaultNetworks() bool ***REMOVED***
	cli, err := client.NewEnvClient()
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	networks, err := cli.NetworkList(context.TODO(), types.NetworkListOptions***REMOVED******REMOVED***)
	if err != nil || len(networks) > 0 ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// Deprecated: use skip.IfCondition(t, !testEnv.DaemonInfo.ExperimentalBuild)
func ExperimentalDaemon() bool ***REMOVED***
	return testEnv.DaemonInfo.ExperimentalBuild
***REMOVED***

func IsAmd64() bool ***REMOVED***
	return os.Getenv("DOCKER_ENGINE_GOARCH") == "amd64"
***REMOVED***

func NotArm() bool ***REMOVED***
	return ArchitectureIsNot("arm")
***REMOVED***

func NotArm64() bool ***REMOVED***
	return ArchitectureIsNot("arm64")
***REMOVED***

func NotPpc64le() bool ***REMOVED***
	return ArchitectureIsNot("ppc64le")
***REMOVED***

func NotS390X() bool ***REMOVED***
	return ArchitectureIsNot("s390x")
***REMOVED***

func SameHostDaemon() bool ***REMOVED***
	return testEnv.IsLocalDaemon()
***REMOVED***

func UnixCli() bool ***REMOVED***
	return isUnixCli
***REMOVED***

func ExecSupport() bool ***REMOVED***
	return supportsExec
***REMOVED***

func Network() bool ***REMOVED***
	// Set a timeout on the GET at 15s
	var timeout = time.Duration(15 * time.Second)
	var url = "https://hub.docker.com"

	client := http.Client***REMOVED***
		Timeout: timeout,
	***REMOVED***

	resp, err := client.Get(url)
	if err != nil && strings.Contains(err.Error(), "use of closed network connection") ***REMOVED***
		panic(fmt.Sprintf("Timeout for GET request on %s", url))
	***REMOVED***
	if resp != nil ***REMOVED***
		resp.Body.Close()
	***REMOVED***
	return err == nil
***REMOVED***

func Apparmor() bool ***REMOVED***
	buf, err := ioutil.ReadFile("/sys/module/apparmor/parameters/enabled")
	return err == nil && len(buf) > 1 && buf[0] == 'Y'
***REMOVED***

func NotaryHosting() bool ***REMOVED***
	// for now notary binary is built only if we're running inside
	// container through `make test`. Figure that out by testing if
	// notary-server binary is in PATH.
	_, err := exec.LookPath(notaryServerBinary)
	return err == nil
***REMOVED***

func NotaryServerHosting() bool ***REMOVED***
	// for now notary-server binary is built only if we're running inside
	// container through `make test`. Figure that out by testing if
	// notary-server binary is in PATH.
	_, err := exec.LookPath(notaryServerBinary)
	return err == nil
***REMOVED***

func Devicemapper() bool ***REMOVED***
	return strings.HasPrefix(testEnv.DaemonInfo.Driver, "devicemapper")
***REMOVED***

func IPv6() bool ***REMOVED***
	cmd := exec.Command("test", "-f", "/proc/net/if_inet6")
	return cmd.Run() != nil
***REMOVED***

func UserNamespaceROMount() bool ***REMOVED***
	// quick case--userns not enabled in this test run
	if os.Getenv("DOCKER_REMAP_ROOT") == "" ***REMOVED***
		return true
	***REMOVED***
	if _, _, err := dockerCmdWithError("run", "--rm", "--read-only", "busybox", "date"); err != nil ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func NotUserNamespace() bool ***REMOVED***
	root := os.Getenv("DOCKER_REMAP_ROOT")
	return root == ""
***REMOVED***

func UserNamespaceInKernel() bool ***REMOVED***
	if _, err := os.Stat("/proc/self/uid_map"); os.IsNotExist(err) ***REMOVED***
		/*
		 * This kernel-provided file only exists if user namespaces are
		 * supported
		 */
		return false
	***REMOVED***

	// We need extra check on redhat based distributions
	if f, err := os.Open("/sys/module/user_namespace/parameters/enable"); err == nil ***REMOVED***
		defer f.Close()
		b := make([]byte, 1)
		_, _ = f.Read(b)
		return string(b) != "N"
	***REMOVED***

	return true
***REMOVED***

func IsPausable() bool ***REMOVED***
	if testEnv.OSType == "windows" ***REMOVED***
		return testEnv.DaemonInfo.Isolation == "hyperv"
	***REMOVED***
	return true
***REMOVED***

func NotPausable() bool ***REMOVED***
	if testEnv.OSType == "windows" ***REMOVED***
		return testEnv.DaemonInfo.Isolation == "process"
	***REMOVED***
	return false
***REMOVED***

func IsolationIs(expectedIsolation string) bool ***REMOVED***
	return testEnv.OSType == "windows" && string(testEnv.DaemonInfo.Isolation) == expectedIsolation
***REMOVED***

func IsolationIsHyperv() bool ***REMOVED***
	return IsolationIs("hyperv")
***REMOVED***

func IsolationIsProcess() bool ***REMOVED***
	return IsolationIs("process")
***REMOVED***

// testRequires checks if the environment satisfies the requirements
// for the test to run or skips the tests.
func testRequires(c requirement.SkipT, requirements ...requirement.Test) ***REMOVED***
	requirement.Is(c, requirements...)
***REMOVED***
