// +build !windows

package runconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/sysinfo"
	"github.com/stretchr/testify/assert"
)

// TODO Windows: This will need addressing for a Windows daemon.
func TestNetworkModeTest(t *testing.T) ***REMOVED***
	networkModes := map[container.NetworkMode][]bool***REMOVED***
		// private, bridge, host, container, none, default
		"":                         ***REMOVED***true, false, false, false, false, false***REMOVED***,
		"something:weird":          ***REMOVED***true, false, false, false, false, false***REMOVED***,
		"bridge":                   ***REMOVED***true, true, false, false, false, false***REMOVED***,
		DefaultDaemonNetworkMode(): ***REMOVED***true, true, false, false, false, false***REMOVED***,
		"host":           ***REMOVED***false, false, true, false, false, false***REMOVED***,
		"container:name": ***REMOVED***false, false, false, true, false, false***REMOVED***,
		"none":           ***REMOVED***true, false, false, false, true, false***REMOVED***,
		"default":        ***REMOVED***true, false, false, false, false, true***REMOVED***,
	***REMOVED***
	networkModeNames := map[container.NetworkMode]string***REMOVED***
		"":                         "",
		"something:weird":          "something:weird",
		"bridge":                   "bridge",
		DefaultDaemonNetworkMode(): "bridge",
		"host":           "host",
		"container:name": "container",
		"none":           "none",
		"default":        "default",
	***REMOVED***
	for networkMode, state := range networkModes ***REMOVED***
		if networkMode.IsPrivate() != state[0] ***REMOVED***
			t.Fatalf("NetworkMode.IsPrivate for %v should have been %v but was %v", networkMode, state[0], networkMode.IsPrivate())
		***REMOVED***
		if networkMode.IsBridge() != state[1] ***REMOVED***
			t.Fatalf("NetworkMode.IsBridge for %v should have been %v but was %v", networkMode, state[1], networkMode.IsBridge())
		***REMOVED***
		if networkMode.IsHost() != state[2] ***REMOVED***
			t.Fatalf("NetworkMode.IsHost for %v should have been %v but was %v", networkMode, state[2], networkMode.IsHost())
		***REMOVED***
		if networkMode.IsContainer() != state[3] ***REMOVED***
			t.Fatalf("NetworkMode.IsContainer for %v should have been %v but was %v", networkMode, state[3], networkMode.IsContainer())
		***REMOVED***
		if networkMode.IsNone() != state[4] ***REMOVED***
			t.Fatalf("NetworkMode.IsNone for %v should have been %v but was %v", networkMode, state[4], networkMode.IsNone())
		***REMOVED***
		if networkMode.IsDefault() != state[5] ***REMOVED***
			t.Fatalf("NetworkMode.IsDefault for %v should have been %v but was %v", networkMode, state[5], networkMode.IsDefault())
		***REMOVED***
		if networkMode.NetworkName() != networkModeNames[networkMode] ***REMOVED***
			t.Fatalf("Expected name %v, got %v", networkModeNames[networkMode], networkMode.NetworkName())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIpcModeTest(t *testing.T) ***REMOVED***
	ipcModes := map[container.IpcMode]struct ***REMOVED***
		private   bool
		host      bool
		container bool
		shareable bool
		valid     bool
		ctrName   string
	***REMOVED******REMOVED***
		"":                      ***REMOVED***valid: true***REMOVED***,
		"private":               ***REMOVED***private: true, valid: true***REMOVED***,
		"something:weird":       ***REMOVED******REMOVED***,
		":weird":                ***REMOVED******REMOVED***,
		"host":                  ***REMOVED***host: true, valid: true***REMOVED***,
		"container":             ***REMOVED******REMOVED***,
		"container:":            ***REMOVED***container: true, valid: true, ctrName: ""***REMOVED***,
		"container:name":        ***REMOVED***container: true, valid: true, ctrName: "name"***REMOVED***,
		"container:name1:name2": ***REMOVED***container: true, valid: true, ctrName: "name1:name2"***REMOVED***,
		"shareable":             ***REMOVED***shareable: true, valid: true***REMOVED***,
	***REMOVED***

	for ipcMode, state := range ipcModes ***REMOVED***
		assert.Equal(t, state.private, ipcMode.IsPrivate(), "IpcMode.IsPrivate() parsing failed for %q", ipcMode)
		assert.Equal(t, state.host, ipcMode.IsHost(), "IpcMode.IsHost()  parsing failed for %q", ipcMode)
		assert.Equal(t, state.container, ipcMode.IsContainer(), "IpcMode.IsContainer()  parsing failed for %q", ipcMode)
		assert.Equal(t, state.shareable, ipcMode.IsShareable(), "IpcMode.IsShareable()  parsing failed for %q", ipcMode)
		assert.Equal(t, state.valid, ipcMode.Valid(), "IpcMode.Valid()  parsing failed for %q", ipcMode)
		assert.Equal(t, state.ctrName, ipcMode.Container(), "IpcMode.Container() parsing failed for %q", ipcMode)
	***REMOVED***
***REMOVED***

func TestUTSModeTest(t *testing.T) ***REMOVED***
	utsModes := map[container.UTSMode][]bool***REMOVED***
		// private, host, valid
		"":                ***REMOVED***true, false, true***REMOVED***,
		"something:weird": ***REMOVED***true, false, false***REMOVED***,
		"host":            ***REMOVED***false, true, true***REMOVED***,
		"host:name":       ***REMOVED***true, false, true***REMOVED***,
	***REMOVED***
	for utsMode, state := range utsModes ***REMOVED***
		if utsMode.IsPrivate() != state[0] ***REMOVED***
			t.Fatalf("UtsMode.IsPrivate for %v should have been %v but was %v", utsMode, state[0], utsMode.IsPrivate())
		***REMOVED***
		if utsMode.IsHost() != state[1] ***REMOVED***
			t.Fatalf("UtsMode.IsHost for %v should have been %v but was %v", utsMode, state[1], utsMode.IsHost())
		***REMOVED***
		if utsMode.Valid() != state[2] ***REMOVED***
			t.Fatalf("UtsMode.Valid for %v should have been %v but was %v", utsMode, state[2], utsMode.Valid())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUsernsModeTest(t *testing.T) ***REMOVED***
	usrensMode := map[container.UsernsMode][]bool***REMOVED***
		// private, host, valid
		"":                ***REMOVED***true, false, true***REMOVED***,
		"something:weird": ***REMOVED***true, false, false***REMOVED***,
		"host":            ***REMOVED***false, true, true***REMOVED***,
		"host:name":       ***REMOVED***true, false, true***REMOVED***,
	***REMOVED***
	for usernsMode, state := range usrensMode ***REMOVED***
		if usernsMode.IsPrivate() != state[0] ***REMOVED***
			t.Fatalf("UsernsMode.IsPrivate for %v should have been %v but was %v", usernsMode, state[0], usernsMode.IsPrivate())
		***REMOVED***
		if usernsMode.IsHost() != state[1] ***REMOVED***
			t.Fatalf("UsernsMode.IsHost for %v should have been %v but was %v", usernsMode, state[1], usernsMode.IsHost())
		***REMOVED***
		if usernsMode.Valid() != state[2] ***REMOVED***
			t.Fatalf("UsernsMode.Valid for %v should have been %v but was %v", usernsMode, state[2], usernsMode.Valid())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPidModeTest(t *testing.T) ***REMOVED***
	pidModes := map[container.PidMode][]bool***REMOVED***
		// private, host, valid
		"":                ***REMOVED***true, false, true***REMOVED***,
		"something:weird": ***REMOVED***true, false, false***REMOVED***,
		"host":            ***REMOVED***false, true, true***REMOVED***,
		"host:name":       ***REMOVED***true, false, true***REMOVED***,
	***REMOVED***
	for pidMode, state := range pidModes ***REMOVED***
		if pidMode.IsPrivate() != state[0] ***REMOVED***
			t.Fatalf("PidMode.IsPrivate for %v should have been %v but was %v", pidMode, state[0], pidMode.IsPrivate())
		***REMOVED***
		if pidMode.IsHost() != state[1] ***REMOVED***
			t.Fatalf("PidMode.IsHost for %v should have been %v but was %v", pidMode, state[1], pidMode.IsHost())
		***REMOVED***
		if pidMode.Valid() != state[2] ***REMOVED***
			t.Fatalf("PidMode.Valid for %v should have been %v but was %v", pidMode, state[2], pidMode.Valid())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRestartPolicy(t *testing.T) ***REMOVED***
	restartPolicies := map[container.RestartPolicy][]bool***REMOVED***
		// none, always, failure
		***REMOVED******REMOVED***: ***REMOVED***true, false, false***REMOVED***,
		***REMOVED***Name: "something", MaximumRetryCount: 0***REMOVED***:  ***REMOVED***false, false, false***REMOVED***,
		***REMOVED***Name: "no", MaximumRetryCount: 0***REMOVED***:         ***REMOVED***true, false, false***REMOVED***,
		***REMOVED***Name: "always", MaximumRetryCount: 0***REMOVED***:     ***REMOVED***false, true, false***REMOVED***,
		***REMOVED***Name: "on-failure", MaximumRetryCount: 0***REMOVED***: ***REMOVED***false, false, true***REMOVED***,
	***REMOVED***
	for restartPolicy, state := range restartPolicies ***REMOVED***
		if restartPolicy.IsNone() != state[0] ***REMOVED***
			t.Fatalf("RestartPolicy.IsNone for %v should have been %v but was %v", restartPolicy, state[0], restartPolicy.IsNone())
		***REMOVED***
		if restartPolicy.IsAlways() != state[1] ***REMOVED***
			t.Fatalf("RestartPolicy.IsAlways for %v should have been %v but was %v", restartPolicy, state[1], restartPolicy.IsAlways())
		***REMOVED***
		if restartPolicy.IsOnFailure() != state[2] ***REMOVED***
			t.Fatalf("RestartPolicy.IsOnFailure for %v should have been %v but was %v", restartPolicy, state[2], restartPolicy.IsOnFailure())
		***REMOVED***
	***REMOVED***
***REMOVED***
func TestDecodeHostConfig(t *testing.T) ***REMOVED***
	fixtures := []struct ***REMOVED***
		file string
	***REMOVED******REMOVED***
		***REMOVED***"fixtures/unix/container_hostconfig_1_14.json"***REMOVED***,
		***REMOVED***"fixtures/unix/container_hostconfig_1_19.json"***REMOVED***,
	***REMOVED***

	for _, f := range fixtures ***REMOVED***
		b, err := ioutil.ReadFile(f.file)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		c, err := decodeHostConfig(bytes.NewReader(b))
		if err != nil ***REMOVED***
			t.Fatal(fmt.Errorf("Error parsing %s: %v", f, err))
		***REMOVED***

		assert.False(t, c.Privileged)

		if l := len(c.Binds); l != 1 ***REMOVED***
			t.Fatalf("Expected 1 bind, found %d\n", l)
		***REMOVED***

		if len(c.CapAdd) != 1 && c.CapAdd[0] != "NET_ADMIN" ***REMOVED***
			t.Fatalf("Expected CapAdd NET_ADMIN, got %v", c.CapAdd)
		***REMOVED***

		if len(c.CapDrop) != 1 && c.CapDrop[0] != "NET_ADMIN" ***REMOVED***
			t.Fatalf("Expected CapDrop NET_ADMIN, got %v", c.CapDrop)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestValidateResources(t *testing.T) ***REMOVED***
	type resourceTest struct ***REMOVED***
		ConfigCPURealtimePeriod   int64
		ConfigCPURealtimeRuntime  int64
		SysInfoCPURealtimePeriod  bool
		SysInfoCPURealtimeRuntime bool
		ErrorExpected             bool
		FailureMsg                string
	***REMOVED***

	tests := []resourceTest***REMOVED***
		***REMOVED***
			ConfigCPURealtimePeriod:   1000,
			ConfigCPURealtimeRuntime:  1000,
			SysInfoCPURealtimePeriod:  true,
			SysInfoCPURealtimeRuntime: true,
			ErrorExpected:             false,
			FailureMsg:                "Expected valid configuration",
		***REMOVED***,
		***REMOVED***
			ConfigCPURealtimePeriod:   5000,
			ConfigCPURealtimeRuntime:  5000,
			SysInfoCPURealtimePeriod:  false,
			SysInfoCPURealtimeRuntime: true,
			ErrorExpected:             true,
			FailureMsg:                "Expected failure when cpu-rt-period is set but kernel doesn't support it",
		***REMOVED***,
		***REMOVED***
			ConfigCPURealtimePeriod:   5000,
			ConfigCPURealtimeRuntime:  5000,
			SysInfoCPURealtimePeriod:  true,
			SysInfoCPURealtimeRuntime: false,
			ErrorExpected:             true,
			FailureMsg:                "Expected failure when cpu-rt-runtime is set but kernel doesn't support it",
		***REMOVED***,
		***REMOVED***
			ConfigCPURealtimePeriod:   5000,
			ConfigCPURealtimeRuntime:  10000,
			SysInfoCPURealtimePeriod:  true,
			SysInfoCPURealtimeRuntime: false,
			ErrorExpected:             true,
			FailureMsg:                "Expected failure when cpu-rt-runtime is greater than cpu-rt-period",
		***REMOVED***,
	***REMOVED***

	for _, rt := range tests ***REMOVED***
		var hc container.HostConfig
		hc.Resources.CPURealtimePeriod = rt.ConfigCPURealtimePeriod
		hc.Resources.CPURealtimeRuntime = rt.ConfigCPURealtimeRuntime

		var si sysinfo.SysInfo
		si.CPURealtimePeriod = rt.SysInfoCPURealtimePeriod
		si.CPURealtimeRuntime = rt.SysInfoCPURealtimeRuntime

		if err := validateResources(&hc, &si); (err != nil) != rt.ErrorExpected ***REMOVED***
			t.Fatal(rt.FailureMsg, err)
		***REMOVED***
	***REMOVED***
***REMOVED***
