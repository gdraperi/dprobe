package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	cliconfig "github.com/docker/docker/cli/config"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/integration-cli/fixtures/plugin"
	"github.com/docker/docker/integration-cli/request"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

var notaryBinary = "notary"
var notaryServerBinary = "notary-server"

type keyPair struct ***REMOVED***
	Public  string
	Private string
***REMOVED***

type testNotary struct ***REMOVED***
	cmd  *exec.Cmd
	dir  string
	keys []keyPair
***REMOVED***

const notaryHost = "localhost:4443"
const notaryURL = "https://" + notaryHost

var SuccessTagging = icmd.Expected***REMOVED***
	Out: "Tagging",
***REMOVED***

var SuccessSigningAndPushing = icmd.Expected***REMOVED***
	Out: "Signing and pushing trust metadata",
***REMOVED***

var SuccessDownloaded = icmd.Expected***REMOVED***
	Out: "Status: Downloaded",
***REMOVED***

var SuccessDownloadedOnStderr = icmd.Expected***REMOVED***
	Err: "Status: Downloaded",
***REMOVED***

func newTestNotary(c *check.C) (*testNotary, error) ***REMOVED***
	// generate server config
	template := `***REMOVED***
	"server": ***REMOVED***
		"http_addr": "%s",
		"tls_key_file": "%s",
		"tls_cert_file": "%s"
	***REMOVED***,
	"trust_service": ***REMOVED***
		"type": "local",
		"hostname": "",
		"port": "",
		"key_algorithm": "ed25519"
	***REMOVED***,
	"logging": ***REMOVED***
		"level": "debug"
	***REMOVED***,
	"storage": ***REMOVED***
        "backend": "memory"
***REMOVED***
***REMOVED***`
	tmp, err := ioutil.TempDir("", "notary-test-")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	confPath := filepath.Join(tmp, "config.json")
	config, err := os.Create(confPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer config.Close()

	workingDir, err := os.Getwd()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err := fmt.Fprintf(config, template, notaryHost, filepath.Join(workingDir, "fixtures/notary/localhost.key"), filepath.Join(workingDir, "fixtures/notary/localhost.cert")); err != nil ***REMOVED***
		os.RemoveAll(tmp)
		return nil, err
	***REMOVED***

	// generate client config
	clientConfPath := filepath.Join(tmp, "client-config.json")
	clientConfig, err := os.Create(clientConfPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer clientConfig.Close()

	template = `***REMOVED***
	"trust_dir" : "%s",
	"remote_server": ***REMOVED***
		"url": "%s",
		"skipTLSVerify": true
	***REMOVED***
***REMOVED***`
	if _, err = fmt.Fprintf(clientConfig, template, filepath.Join(cliconfig.Dir(), "trust"), notaryURL); err != nil ***REMOVED***
		os.RemoveAll(tmp)
		return nil, err
	***REMOVED***

	// load key fixture filenames
	var keys []keyPair
	for i := 1; i < 5; i++ ***REMOVED***
		keys = append(keys, keyPair***REMOVED***
			Public:  filepath.Join(workingDir, fmt.Sprintf("fixtures/notary/delgkey%v.crt", i)),
			Private: filepath.Join(workingDir, fmt.Sprintf("fixtures/notary/delgkey%v.key", i)),
		***REMOVED***)
	***REMOVED***

	// run notary-server
	cmd := exec.Command(notaryServerBinary, "-config", confPath)
	if err := cmd.Start(); err != nil ***REMOVED***
		os.RemoveAll(tmp)
		if os.IsNotExist(err) ***REMOVED***
			c.Skip(err.Error())
		***REMOVED***
		return nil, err
	***REMOVED***

	testNotary := &testNotary***REMOVED***
		cmd:  cmd,
		dir:  tmp,
		keys: keys,
	***REMOVED***

	// Wait for notary to be ready to serve requests.
	for i := 1; i <= 20; i++ ***REMOVED***
		if err = testNotary.Ping(); err == nil ***REMOVED***
			break
		***REMOVED***
		time.Sleep(10 * time.Millisecond * time.Duration(i*i))
	***REMOVED***

	if err != nil ***REMOVED***
		c.Fatalf("Timeout waiting for test notary to become available: %s", err)
	***REMOVED***

	return testNotary, nil
***REMOVED***

func (t *testNotary) Ping() error ***REMOVED***
	tlsConfig := tlsconfig.ClientDefault()
	tlsConfig.InsecureSkipVerify = true
	client := http.Client***REMOVED***
		Transport: &http.Transport***REMOVED***
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer***REMOVED***
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			***REMOVED***).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     tlsConfig,
		***REMOVED***,
	***REMOVED***
	resp, err := client.Get(fmt.Sprintf("%s/v2/", notaryURL))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if resp.StatusCode != http.StatusOK ***REMOVED***
		return fmt.Errorf("notary ping replied with an unexpected status code %d", resp.StatusCode)
	***REMOVED***
	return nil
***REMOVED***

func (t *testNotary) Close() ***REMOVED***
	t.cmd.Process.Kill()
	t.cmd.Process.Wait()
	os.RemoveAll(t.dir)
***REMOVED***

func trustedCmd(cmd *icmd.Cmd) func() ***REMOVED***
	pwd := "12345678"
	cmd.Env = append(cmd.Env, trustEnv(notaryURL, pwd, pwd)...)
	return nil
***REMOVED***

func trustedCmdWithServer(server string) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		pwd := "12345678"
		cmd.Env = append(cmd.Env, trustEnv(server, pwd, pwd)...)
		return nil
	***REMOVED***
***REMOVED***

func trustedCmdWithPassphrases(rootPwd, repositoryPwd string) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Env = append(cmd.Env, trustEnv(notaryURL, rootPwd, repositoryPwd)...)
		return nil
	***REMOVED***
***REMOVED***

func trustEnv(server, rootPwd, repositoryPwd string) []string ***REMOVED***
	env := append(os.Environ(), []string***REMOVED***
		"DOCKER_CONTENT_TRUST=1",
		fmt.Sprintf("DOCKER_CONTENT_TRUST_SERVER=%s", server),
		fmt.Sprintf("DOCKER_CONTENT_TRUST_ROOT_PASSPHRASE=%s", rootPwd),
		fmt.Sprintf("DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE=%s", repositoryPwd),
	***REMOVED***...)
	return env
***REMOVED***

func (s *DockerTrustSuite) setupTrustedImage(c *check.C, name string) string ***REMOVED***
	repoName := fmt.Sprintf("%v/dockercli/%s:latest", privateRegistryURL, name)
	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", repoName)
	cli.Docker(cli.Args("push", repoName), trustedCmd).Assert(c, SuccessSigningAndPushing)
	cli.DockerCmd(c, "rmi", repoName)
	return repoName
***REMOVED***

func (s *DockerTrustSuite) setupTrustedplugin(c *check.C, source, name string) string ***REMOVED***
	repoName := fmt.Sprintf("%v/dockercli/%s:latest", privateRegistryURL, name)

	client, err := request.NewClient()
	c.Assert(err, checker.IsNil, check.Commentf("could not create test client"))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	err = plugin.Create(ctx, client, repoName)
	cancel()
	c.Assert(err, checker.IsNil, check.Commentf("could not create test plugin"))

	// tag the image and upload it to the private registry
	// TODO: shouldn't need to use the CLI to do trust
	cli.Docker(cli.Args("plugin", "push", repoName), trustedCmd).Assert(c, SuccessSigningAndPushing)

	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	err = client.PluginRemove(ctx, repoName, types.PluginRemoveOptions***REMOVED***Force: true***REMOVED***)
	cancel()
	c.Assert(err, checker.IsNil, check.Commentf("failed to cleanup test plugin for trust suite"))
	return repoName
***REMOVED***

func (s *DockerTrustSuite) notaryCmd(c *check.C, args ...string) string ***REMOVED***
	pwd := "12345678"
	env := []string***REMOVED***
		fmt.Sprintf("NOTARY_ROOT_PASSPHRASE=%s", pwd),
		fmt.Sprintf("NOTARY_TARGETS_PASSPHRASE=%s", pwd),
		fmt.Sprintf("NOTARY_SNAPSHOT_PASSPHRASE=%s", pwd),
		fmt.Sprintf("NOTARY_DELEGATION_PASSPHRASE=%s", pwd),
	***REMOVED***
	result := icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: append([]string***REMOVED***notaryBinary, "-c", filepath.Join(s.not.dir, "client-config.json")***REMOVED***, args...),
		Env:     append(os.Environ(), env...),
	***REMOVED***)
	result.Assert(c, icmd.Success)
	return result.Combined()
***REMOVED***

func (s *DockerTrustSuite) notaryInitRepo(c *check.C, repoName string) ***REMOVED***
	s.notaryCmd(c, "init", repoName)
***REMOVED***

func (s *DockerTrustSuite) notaryCreateDelegation(c *check.C, repoName, role string, pubKey string, paths ...string) ***REMOVED***
	pathsArg := "--all-paths"
	if len(paths) > 0 ***REMOVED***
		pathsArg = "--paths=" + strings.Join(paths, ",")
	***REMOVED***

	s.notaryCmd(c, "delegation", "add", repoName, role, pubKey, pathsArg)
***REMOVED***

func (s *DockerTrustSuite) notaryPublish(c *check.C, repoName string) ***REMOVED***
	s.notaryCmd(c, "publish", repoName)
***REMOVED***

func (s *DockerTrustSuite) notaryImportKey(c *check.C, repoName, role string, privKey string) ***REMOVED***
	s.notaryCmd(c, "key", "import", privKey, "-g", repoName, "-r", role)
***REMOVED***

func (s *DockerTrustSuite) notaryListTargetsInRole(c *check.C, repoName, role string) map[string]string ***REMOVED***
	out := s.notaryCmd(c, "list", repoName, "-r", role)

	// should look something like:
	//    NAME                                 DIGEST                                SIZE (BYTES)    ROLE
	// ------------------------------------------------------------------------------------------------------
	//   latest   24a36bbc059b1345b7e8be0df20f1b23caa3602e85d42fff7ecd9d0bd255de56   1377           targets

	targets := make(map[string]string)

	// no target
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 1 && strings.Contains(out, "No targets present in this repository.") ***REMOVED***
		return targets
	***REMOVED***

	// otherwise, there is at least one target
	c.Assert(len(lines), checker.GreaterOrEqualThan, 3)

	for _, line := range lines[2:] ***REMOVED***
		tokens := strings.Fields(line)
		c.Assert(tokens, checker.HasLen, 4)
		targets[tokens[0]] = tokens[3]
	***REMOVED***

	return targets
***REMOVED***

func (s *DockerTrustSuite) assertTargetInRoles(c *check.C, repoName, target string, roles ...string) ***REMOVED***
	// check all the roles
	for _, role := range roles ***REMOVED***
		targets := s.notaryListTargetsInRole(c, repoName, role)
		roleName, ok := targets[target]
		c.Assert(ok, checker.True)
		c.Assert(roleName, checker.Equals, role)
	***REMOVED***
***REMOVED***

func (s *DockerTrustSuite) assertTargetNotInRoles(c *check.C, repoName, target string, roles ...string) ***REMOVED***
	targets := s.notaryListTargetsInRole(c, repoName, "targets")

	roleName, ok := targets[target]
	if ok ***REMOVED***
		for _, role := range roles ***REMOVED***
			c.Assert(roleName, checker.Not(checker.Equals), role)
		***REMOVED***
	***REMOVED***
***REMOVED***
