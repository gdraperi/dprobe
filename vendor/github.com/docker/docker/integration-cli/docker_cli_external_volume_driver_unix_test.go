// +build !windows

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/volume"
	"github.com/go-check/check"
)

const volumePluginName = "test-external-volume-driver"

func init() ***REMOVED***
	check.Suite(&DockerExternalVolumeSuite***REMOVED***
		ds: &DockerSuite***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

type eventCounter struct ***REMOVED***
	activations int
	creations   int
	removals    int
	mounts      int
	unmounts    int
	paths       int
	lists       int
	gets        int
	caps        int
***REMOVED***

type DockerExternalVolumeSuite struct ***REMOVED***
	ds *DockerSuite
	d  *daemon.Daemon
	*volumePlugin
***REMOVED***

func (s *DockerExternalVolumeSuite) SetUpTest(c *check.C) ***REMOVED***
	testRequires(c, SameHostDaemon)
	s.d = daemon.New(c, dockerBinary, dockerdBinary, daemon.Config***REMOVED***
		Experimental: testEnv.DaemonInfo.ExperimentalBuild,
	***REMOVED***)
	s.ec = &eventCounter***REMOVED******REMOVED***
***REMOVED***

func (s *DockerExternalVolumeSuite) TearDownTest(c *check.C) ***REMOVED***
	if s.d != nil ***REMOVED***
		s.d.Stop(c)
		s.ds.TearDownTest(c)
	***REMOVED***
***REMOVED***

func (s *DockerExternalVolumeSuite) SetUpSuite(c *check.C) ***REMOVED***
	s.volumePlugin = newVolumePlugin(c, volumePluginName)
***REMOVED***

type volumePlugin struct ***REMOVED***
	ec *eventCounter
	*httptest.Server
	vols map[string]vol
***REMOVED***

type vol struct ***REMOVED***
	Name       string
	Mountpoint string
	Ninja      bool // hack used to trigger a null volume return on `Get`
	Status     map[string]interface***REMOVED******REMOVED***
	Options    map[string]string
***REMOVED***

func (p *volumePlugin) Close() ***REMOVED***
	p.Server.Close()
***REMOVED***

func newVolumePlugin(c *check.C, name string) *volumePlugin ***REMOVED***
	mux := http.NewServeMux()
	s := &volumePlugin***REMOVED***Server: httptest.NewServer(mux), ec: &eventCounter***REMOVED******REMOVED***, vols: make(map[string]vol)***REMOVED***

	type pluginRequest struct ***REMOVED***
		Name string
		Opts map[string]string
		ID   string
	***REMOVED***

	type pluginResp struct ***REMOVED***
		Mountpoint string `json:",omitempty"`
		Err        string `json:",omitempty"`
	***REMOVED***

	read := func(b io.ReadCloser) (pluginRequest, error) ***REMOVED***
		defer b.Close()
		var pr pluginRequest
		err := json.NewDecoder(b).Decode(&pr)
		return pr, err
	***REMOVED***

	send := func(w http.ResponseWriter, data interface***REMOVED******REMOVED***) ***REMOVED***
		switch t := data.(type) ***REMOVED***
		case error:
			http.Error(w, t.Error(), 500)
		case string:
			w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
			fmt.Fprintln(w, t)
		default:
			w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
			json.NewEncoder(w).Encode(&data)
		***REMOVED***
	***REMOVED***

	mux.HandleFunc("/Plugin.Activate", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s.ec.activations++
		send(w, `***REMOVED***"Implements": ["VolumeDriver"]***REMOVED***`)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Create", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s.ec.creations++
		pr, err := read(r.Body)
		if err != nil ***REMOVED***
			send(w, err)
			return
		***REMOVED***
		_, isNinja := pr.Opts["ninja"]
		status := map[string]interface***REMOVED******REMOVED******REMOVED***"Hello": "world"***REMOVED***
		s.vols[pr.Name] = vol***REMOVED***Name: pr.Name, Ninja: isNinja, Status: status, Options: pr.Opts***REMOVED***
		send(w, nil)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.List", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s.ec.lists++
		vols := make([]vol, 0, len(s.vols))
		for _, v := range s.vols ***REMOVED***
			if v.Ninja ***REMOVED***
				continue
			***REMOVED***
			vols = append(vols, v)
		***REMOVED***
		send(w, map[string][]vol***REMOVED***"Volumes": vols***REMOVED***)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Get", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s.ec.gets++
		pr, err := read(r.Body)
		if err != nil ***REMOVED***
			send(w, err)
			return
		***REMOVED***

		v, exists := s.vols[pr.Name]
		if !exists ***REMOVED***
			send(w, `***REMOVED***"Err": "no such volume"***REMOVED***`)
		***REMOVED***

		if v.Ninja ***REMOVED***
			send(w, map[string]vol***REMOVED******REMOVED***)
			return
		***REMOVED***

		v.Mountpoint = hostVolumePath(pr.Name)
		send(w, map[string]vol***REMOVED***"Volume": v***REMOVED***)
		return
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Remove", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s.ec.removals++
		pr, err := read(r.Body)
		if err != nil ***REMOVED***
			send(w, err)
			return
		***REMOVED***

		v, ok := s.vols[pr.Name]
		if !ok ***REMOVED***
			send(w, nil)
			return
		***REMOVED***

		if err := os.RemoveAll(hostVolumePath(v.Name)); err != nil ***REMOVED***
			send(w, &pluginResp***REMOVED***Err: err.Error()***REMOVED***)
			return
		***REMOVED***
		delete(s.vols, v.Name)
		send(w, nil)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Path", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s.ec.paths++

		pr, err := read(r.Body)
		if err != nil ***REMOVED***
			send(w, err)
			return
		***REMOVED***
		p := hostVolumePath(pr.Name)
		send(w, &pluginResp***REMOVED***Mountpoint: p***REMOVED***)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Mount", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s.ec.mounts++

		pr, err := read(r.Body)
		if err != nil ***REMOVED***
			send(w, err)
			return
		***REMOVED***

		if v, exists := s.vols[pr.Name]; exists ***REMOVED***
			// Use this to simulate a mount failure
			if _, exists := v.Options["invalidOption"]; exists ***REMOVED***
				send(w, fmt.Errorf("invalid argument"))
				return
			***REMOVED***
		***REMOVED***

		p := hostVolumePath(pr.Name)
		if err := os.MkdirAll(p, 0755); err != nil ***REMOVED***
			send(w, &pluginResp***REMOVED***Err: err.Error()***REMOVED***)
			return
		***REMOVED***

		if err := ioutil.WriteFile(filepath.Join(p, "test"), []byte(s.Server.URL), 0644); err != nil ***REMOVED***
			send(w, err)
			return
		***REMOVED***

		if err := ioutil.WriteFile(filepath.Join(p, "mountID"), []byte(pr.ID), 0644); err != nil ***REMOVED***
			send(w, err)
			return
		***REMOVED***

		send(w, &pluginResp***REMOVED***Mountpoint: p***REMOVED***)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Unmount", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s.ec.unmounts++

		_, err := read(r.Body)
		if err != nil ***REMOVED***
			send(w, err)
			return
		***REMOVED***

		send(w, nil)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Capabilities", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s.ec.caps++

		_, err := read(r.Body)
		if err != nil ***REMOVED***
			send(w, err)
			return
		***REMOVED***

		send(w, `***REMOVED***"Capabilities": ***REMOVED*** "Scope": "global" ***REMOVED******REMOVED***`)
	***REMOVED***)

	err := os.MkdirAll("/etc/docker/plugins", 0755)
	c.Assert(err, checker.IsNil)

	err = ioutil.WriteFile("/etc/docker/plugins/"+name+".spec", []byte(s.Server.URL), 0644)
	c.Assert(err, checker.IsNil)
	return s
***REMOVED***

func (s *DockerExternalVolumeSuite) TearDownSuite(c *check.C) ***REMOVED***
	s.volumePlugin.Close()

	err := os.RemoveAll("/etc/docker/plugins")
	c.Assert(err, checker.IsNil)
***REMOVED***

func (s *DockerExternalVolumeSuite) TestVolumeCLICreateOptionConflict(c *check.C) ***REMOVED***
	dockerCmd(c, "volume", "create", "test")

	out, _, err := dockerCmdWithError("volume", "create", "test", "--driver", volumePluginName)
	c.Assert(err, check.NotNil, check.Commentf("volume create exception name already in use with another driver"))
	c.Assert(out, checker.Contains, "must be unique")

	out, _ = dockerCmd(c, "volume", "inspect", "--format=***REMOVED******REMOVED*** .Driver ***REMOVED******REMOVED***", "test")
	_, _, err = dockerCmdWithError("volume", "create", "test", "--driver", strings.TrimSpace(out))
	c.Assert(err, check.IsNil)
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverNamed(c *check.C) ***REMOVED***
	s.d.StartWithBusybox(c)

	out, err := s.d.Cmd("run", "--rm", "--name", "test-data", "-v", "external-volume-test:/tmp/external-volume-test", "--volume-driver", volumePluginName, "busybox:latest", "cat", "/tmp/external-volume-test/test")
	c.Assert(err, checker.IsNil, check.Commentf(out))
	c.Assert(out, checker.Contains, s.Server.URL)

	_, err = s.d.Cmd("volume", "rm", "external-volume-test")
	c.Assert(err, checker.IsNil)

	p := hostVolumePath("external-volume-test")
	_, err = os.Lstat(p)
	c.Assert(err, checker.NotNil)
	c.Assert(os.IsNotExist(err), checker.True, check.Commentf("Expected volume path in host to not exist: %s, %v\n", p, err))

	c.Assert(s.ec.activations, checker.Equals, 1)
	c.Assert(s.ec.creations, checker.Equals, 1)
	c.Assert(s.ec.removals, checker.Equals, 1)
	c.Assert(s.ec.mounts, checker.Equals, 1)
	c.Assert(s.ec.unmounts, checker.Equals, 1)
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverUnnamed(c *check.C) ***REMOVED***
	s.d.StartWithBusybox(c)

	out, err := s.d.Cmd("run", "--rm", "--name", "test-data", "-v", "/tmp/external-volume-test", "--volume-driver", volumePluginName, "busybox:latest", "cat", "/tmp/external-volume-test/test")
	c.Assert(err, checker.IsNil, check.Commentf(out))
	c.Assert(out, checker.Contains, s.Server.URL)

	c.Assert(s.ec.activations, checker.Equals, 1)
	c.Assert(s.ec.creations, checker.Equals, 1)
	c.Assert(s.ec.removals, checker.Equals, 1)
	c.Assert(s.ec.mounts, checker.Equals, 1)
	c.Assert(s.ec.unmounts, checker.Equals, 1)
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverVolumesFrom(c *check.C) ***REMOVED***
	s.d.StartWithBusybox(c)

	out, err := s.d.Cmd("run", "--name", "vol-test1", "-v", "/foo", "--volume-driver", volumePluginName, "busybox:latest")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = s.d.Cmd("run", "--rm", "--volumes-from", "vol-test1", "--name", "vol-test2", "busybox", "ls", "/tmp")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = s.d.Cmd("rm", "-fv", "vol-test1")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	c.Assert(s.ec.activations, checker.Equals, 1)
	c.Assert(s.ec.creations, checker.Equals, 1)
	c.Assert(s.ec.removals, checker.Equals, 1)
	c.Assert(s.ec.mounts, checker.Equals, 2)
	c.Assert(s.ec.unmounts, checker.Equals, 2)
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverDeleteContainer(c *check.C) ***REMOVED***
	s.d.StartWithBusybox(c)

	out, err := s.d.Cmd("run", "--name", "vol-test1", "-v", "/foo", "--volume-driver", volumePluginName, "busybox:latest")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = s.d.Cmd("rm", "-fv", "vol-test1")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	c.Assert(s.ec.activations, checker.Equals, 1)
	c.Assert(s.ec.creations, checker.Equals, 1)
	c.Assert(s.ec.removals, checker.Equals, 1)
	c.Assert(s.ec.mounts, checker.Equals, 1)
	c.Assert(s.ec.unmounts, checker.Equals, 1)
***REMOVED***

func hostVolumePath(name string) string ***REMOVED***
	return fmt.Sprintf("/var/lib/docker/volumes/%s", name)
***REMOVED***

// Make sure a request to use a down driver doesn't block other requests
func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverLookupNotBlocked(c *check.C) ***REMOVED***
	specPath := "/etc/docker/plugins/down-driver.spec"
	err := ioutil.WriteFile(specPath, []byte("tcp://127.0.0.7:9999"), 0644)
	c.Assert(err, check.IsNil)
	defer os.RemoveAll(specPath)

	chCmd1 := make(chan struct***REMOVED******REMOVED***)
	chCmd2 := make(chan error)
	cmd1 := exec.Command(dockerBinary, "volume", "create", "-d", "down-driver")
	cmd2 := exec.Command(dockerBinary, "volume", "create")

	c.Assert(cmd1.Start(), checker.IsNil)
	defer cmd1.Process.Kill()
	time.Sleep(100 * time.Millisecond) // ensure API has been called
	c.Assert(cmd2.Start(), checker.IsNil)

	go func() ***REMOVED***
		cmd1.Wait()
		close(chCmd1)
	***REMOVED***()
	go func() ***REMOVED***
		chCmd2 <- cmd2.Wait()
	***REMOVED***()

	select ***REMOVED***
	case <-chCmd1:
		cmd2.Process.Kill()
		c.Fatalf("volume create with down driver finished unexpectedly")
	case err := <-chCmd2:
		c.Assert(err, checker.IsNil)
	case <-time.After(5 * time.Second):
		cmd2.Process.Kill()
		c.Fatal("volume creates are blocked by previous create requests when previous driver is down")
	***REMOVED***
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverRetryNotImmediatelyExists(c *check.C) ***REMOVED***
	s.d.StartWithBusybox(c)
	driverName := "test-external-volume-driver-retry"

	errchan := make(chan error)
	started := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		close(started)
		if out, err := s.d.Cmd("run", "--rm", "--name", "test-data-retry", "-v", "external-volume-test:/tmp/external-volume-test", "--volume-driver", driverName, "busybox:latest"); err != nil ***REMOVED***
			errchan <- fmt.Errorf("%v:\n%s", err, out)
		***REMOVED***
		close(errchan)
	***REMOVED***()

	<-started
	// wait for a retry to occur, then create spec to allow plugin to register
	time.Sleep(2 * time.Second)
	p := newVolumePlugin(c, driverName)
	defer p.Close()

	select ***REMOVED***
	case err := <-errchan:
		c.Assert(err, checker.IsNil)
	case <-time.After(8 * time.Second):
		c.Fatal("volume creates fail when plugin not immediately available")
	***REMOVED***

	_, err := s.d.Cmd("volume", "rm", "external-volume-test")
	c.Assert(err, checker.IsNil)

	c.Assert(p.ec.activations, checker.Equals, 1)
	c.Assert(p.ec.creations, checker.Equals, 1)
	c.Assert(p.ec.removals, checker.Equals, 1)
	c.Assert(p.ec.mounts, checker.Equals, 1)
	c.Assert(p.ec.unmounts, checker.Equals, 1)
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverBindExternalVolume(c *check.C) ***REMOVED***
	dockerCmd(c, "volume", "create", "-d", volumePluginName, "foo")
	dockerCmd(c, "run", "-d", "--name", "testing", "-v", "foo:/bar", "busybox", "top")

	var mounts []struct ***REMOVED***
		Name   string
		Driver string
	***REMOVED***
	out := inspectFieldJSON(c, "testing", "Mounts")
	c.Assert(json.NewDecoder(strings.NewReader(out)).Decode(&mounts), checker.IsNil)
	c.Assert(len(mounts), checker.Equals, 1, check.Commentf(out))
	c.Assert(mounts[0].Name, checker.Equals, "foo")
	c.Assert(mounts[0].Driver, checker.Equals, volumePluginName)
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverList(c *check.C) ***REMOVED***
	dockerCmd(c, "volume", "create", "-d", volumePluginName, "abc3")
	out, _ := dockerCmd(c, "volume", "ls")
	ls := strings.Split(strings.TrimSpace(out), "\n")
	c.Assert(len(ls), check.Equals, 2, check.Commentf("\n%s", out))

	vol := strings.Fields(ls[len(ls)-1])
	c.Assert(len(vol), check.Equals, 2, check.Commentf("%v", vol))
	c.Assert(vol[0], check.Equals, volumePluginName)
	c.Assert(vol[1], check.Equals, "abc3")

	c.Assert(s.ec.lists, check.Equals, 1)
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverGet(c *check.C) ***REMOVED***
	out, _, err := dockerCmdWithError("volume", "inspect", "dummy")
	c.Assert(err, check.NotNil, check.Commentf(out))
	c.Assert(out, checker.Contains, "No such volume")
	c.Assert(s.ec.gets, check.Equals, 1)

	dockerCmd(c, "volume", "create", "test", "-d", volumePluginName)
	out, _ = dockerCmd(c, "volume", "inspect", "test")

	type vol struct ***REMOVED***
		Status map[string]string
	***REMOVED***
	var st []vol

	c.Assert(json.Unmarshal([]byte(out), &st), checker.IsNil)
	c.Assert(st, checker.HasLen, 1)
	c.Assert(st[0].Status, checker.HasLen, 1, check.Commentf("%v", st[0]))
	c.Assert(st[0].Status["Hello"], checker.Equals, "world", check.Commentf("%v", st[0].Status))
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverWithDaemonRestart(c *check.C) ***REMOVED***
	dockerCmd(c, "volume", "create", "-d", volumePluginName, "abc1")
	s.d.Restart(c)

	dockerCmd(c, "run", "--name=test", "-v", "abc1:/foo", "busybox", "true")
	var mounts []types.MountPoint
	inspectFieldAndUnmarshall(c, "test", "Mounts", &mounts)
	c.Assert(mounts, checker.HasLen, 1)
	c.Assert(mounts[0].Driver, checker.Equals, volumePluginName)
***REMOVED***

// Ensures that the daemon handles when the plugin responds to a `Get` request with a null volume and a null error.
// Prior the daemon would panic in this scenario.
func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverGetEmptyResponse(c *check.C) ***REMOVED***
	s.d.Start(c)

	out, err := s.d.Cmd("volume", "create", "-d", volumePluginName, "abc2", "--opt", "ninja=1")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = s.d.Cmd("volume", "inspect", "abc2")
	c.Assert(err, checker.NotNil, check.Commentf(out))
	c.Assert(out, checker.Contains, "No such volume")
***REMOVED***

// Ensure only cached paths are used in volume list to prevent N+1 calls to `VolumeDriver.Path`
func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverPathCalls(c *check.C) ***REMOVED***
	s.d.Start(c)
	c.Assert(s.ec.paths, checker.Equals, 0)

	out, err := s.d.Cmd("volume", "create", "test", "--driver=test-external-volume-driver")
	c.Assert(err, checker.IsNil, check.Commentf(out))
	c.Assert(s.ec.paths, checker.Equals, 1)

	out, err = s.d.Cmd("volume", "ls")
	c.Assert(err, checker.IsNil, check.Commentf(out))
	c.Assert(s.ec.paths, checker.Equals, 1)

	out, err = s.d.Cmd("volume", "inspect", "--format='***REMOVED******REMOVED***.Mountpoint***REMOVED******REMOVED***'", "test")
	c.Assert(err, checker.IsNil, check.Commentf(out))
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Equals), "")
	c.Assert(s.ec.paths, checker.Equals, 1)
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverMountID(c *check.C) ***REMOVED***
	s.d.StartWithBusybox(c)

	out, err := s.d.Cmd("run", "--rm", "-v", "external-volume-test:/tmp/external-volume-test", "--volume-driver", volumePluginName, "busybox:latest", "cat", "/tmp/external-volume-test/test")
	c.Assert(err, checker.IsNil, check.Commentf(out))
	c.Assert(strings.TrimSpace(out), checker.Not(checker.Equals), "")
***REMOVED***

// Check that VolumeDriver.Capabilities gets called, and only called once
func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverCapabilities(c *check.C) ***REMOVED***
	s.d.Start(c)
	c.Assert(s.ec.caps, checker.Equals, 0)

	for i := 0; i < 3; i++ ***REMOVED***
		out, err := s.d.Cmd("volume", "create", "-d", volumePluginName, fmt.Sprintf("test%d", i))
		c.Assert(err, checker.IsNil, check.Commentf(out))
		c.Assert(s.ec.caps, checker.Equals, 1)
		out, err = s.d.Cmd("volume", "inspect", "--format=***REMOVED******REMOVED***.Scope***REMOVED******REMOVED***", fmt.Sprintf("test%d", i))
		c.Assert(err, checker.IsNil)
		c.Assert(strings.TrimSpace(out), checker.Equals, volume.GlobalScope)
	***REMOVED***
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverOutOfBandDelete(c *check.C) ***REMOVED***
	driverName := stringid.GenerateNonCryptoID()
	p := newVolumePlugin(c, driverName)
	defer p.Close()

	s.d.StartWithBusybox(c)

	out, err := s.d.Cmd("volume", "create", "-d", driverName, "--name", "test")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = s.d.Cmd("volume", "create", "-d", "local", "--name", "test")
	c.Assert(err, checker.NotNil, check.Commentf(out))
	c.Assert(out, checker.Contains, "must be unique")

	// simulate out of band volume deletion on plugin level
	delete(p.vols, "test")

	// test re-create with same driver
	out, err = s.d.Cmd("volume", "create", "-d", driverName, "--opt", "foo=bar", "--name", "test")
	c.Assert(err, checker.IsNil, check.Commentf(out))
	out, err = s.d.Cmd("volume", "inspect", "test")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	var vs []types.Volume
	err = json.Unmarshal([]byte(out), &vs)
	c.Assert(err, checker.IsNil)
	c.Assert(vs, checker.HasLen, 1)
	c.Assert(vs[0].Driver, checker.Equals, driverName)
	c.Assert(vs[0].Options, checker.NotNil)
	c.Assert(vs[0].Options["foo"], checker.Equals, "bar")
	c.Assert(vs[0].Driver, checker.Equals, driverName)

	// simulate out of band volume deletion on plugin level
	delete(p.vols, "test")

	// test create with different driver
	out, err = s.d.Cmd("volume", "create", "-d", "local", "--name", "test")
	c.Assert(err, checker.IsNil, check.Commentf(out))

	out, err = s.d.Cmd("volume", "inspect", "test")
	c.Assert(err, checker.IsNil, check.Commentf(out))
	vs = nil
	err = json.Unmarshal([]byte(out), &vs)
	c.Assert(err, checker.IsNil)
	c.Assert(vs, checker.HasLen, 1)
	c.Assert(vs[0].Options, checker.HasLen, 0)
	c.Assert(vs[0].Driver, checker.Equals, "local")
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverUnmountOnMountFail(c *check.C) ***REMOVED***
	s.d.StartWithBusybox(c)
	s.d.Cmd("volume", "create", "-d", "test-external-volume-driver", "--opt=invalidOption=1", "--name=testumount")

	out, _ := s.d.Cmd("run", "-v", "testumount:/foo", "busybox", "true")
	c.Assert(s.ec.unmounts, checker.Equals, 0, check.Commentf(out))
	out, _ = s.d.Cmd("run", "-w", "/foo", "-v", "testumount:/foo", "busybox", "true")
	c.Assert(s.ec.unmounts, checker.Equals, 0, check.Commentf(out))
***REMOVED***

func (s *DockerExternalVolumeSuite) TestExternalVolumeDriverUnmountOnCp(c *check.C) ***REMOVED***
	s.d.StartWithBusybox(c)
	s.d.Cmd("volume", "create", "-d", "test-external-volume-driver", "--name=test")

	out, _ := s.d.Cmd("run", "-d", "--name=test", "-v", "test:/foo", "busybox", "/bin/sh", "-c", "touch /test && top")
	c.Assert(s.ec.mounts, checker.Equals, 1, check.Commentf(out))

	out, _ = s.d.Cmd("cp", "test:/test", "/tmp/test")
	c.Assert(s.ec.mounts, checker.Equals, 2, check.Commentf(out))
	c.Assert(s.ec.unmounts, checker.Equals, 1, check.Commentf(out))

	out, _ = s.d.Cmd("kill", "test")
	c.Assert(s.ec.unmounts, checker.Equals, 2, check.Commentf(out))
***REMOVED***
