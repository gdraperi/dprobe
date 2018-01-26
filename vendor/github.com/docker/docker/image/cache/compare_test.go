package cache

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

// Just to make life easier
func newPortNoError(proto, port string) nat.Port ***REMOVED***
	p, _ := nat.NewPort(proto, port)
	return p
***REMOVED***

func TestCompare(t *testing.T) ***REMOVED***
	ports1 := make(nat.PortSet)
	ports1[newPortNoError("tcp", "1111")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports1[newPortNoError("tcp", "2222")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports2 := make(nat.PortSet)
	ports2[newPortNoError("tcp", "3333")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports2[newPortNoError("tcp", "4444")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports3 := make(nat.PortSet)
	ports3[newPortNoError("tcp", "1111")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports3[newPortNoError("tcp", "2222")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports3[newPortNoError("tcp", "5555")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	volumes1 := make(map[string]struct***REMOVED******REMOVED***)
	volumes1["/test1"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	volumes2 := make(map[string]struct***REMOVED******REMOVED***)
	volumes2["/test2"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	volumes3 := make(map[string]struct***REMOVED******REMOVED***)
	volumes3["/test1"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	volumes3["/test3"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	envs1 := []string***REMOVED***"ENV1=value1", "ENV2=value2"***REMOVED***
	envs2 := []string***REMOVED***"ENV1=value1", "ENV3=value3"***REMOVED***
	entrypoint1 := strslice.StrSlice***REMOVED***"/bin/sh", "-c"***REMOVED***
	entrypoint2 := strslice.StrSlice***REMOVED***"/bin/sh", "-d"***REMOVED***
	entrypoint3 := strslice.StrSlice***REMOVED***"/bin/sh", "-c", "echo"***REMOVED***
	cmd1 := strslice.StrSlice***REMOVED***"/bin/sh", "-c"***REMOVED***
	cmd2 := strslice.StrSlice***REMOVED***"/bin/sh", "-d"***REMOVED***
	cmd3 := strslice.StrSlice***REMOVED***"/bin/sh", "-c", "echo"***REMOVED***
	labels1 := map[string]string***REMOVED***"LABEL1": "value1", "LABEL2": "value2"***REMOVED***
	labels2 := map[string]string***REMOVED***"LABEL1": "value1", "LABEL2": "value3"***REMOVED***
	labels3 := map[string]string***REMOVED***"LABEL1": "value1", "LABEL2": "value2", "LABEL3": "value3"***REMOVED***

	sameConfigs := map[*container.Config]*container.Config***REMOVED***
		// Empty config
		***REMOVED******REMOVED***: ***REMOVED******REMOVED***,
		// Does not compare hostname, domainname & image
		***REMOVED***
			Hostname:   "host1",
			Domainname: "domain1",
			Image:      "image1",
			User:       "user",
		***REMOVED***: ***REMOVED***
			Hostname:   "host2",
			Domainname: "domain2",
			Image:      "image2",
			User:       "user",
		***REMOVED***,
		// only OpenStdin
		***REMOVED***OpenStdin: false***REMOVED***: ***REMOVED***OpenStdin: false***REMOVED***,
		// only env
		***REMOVED***Env: envs1***REMOVED***: ***REMOVED***Env: envs1***REMOVED***,
		// only cmd
		***REMOVED***Cmd: cmd1***REMOVED***: ***REMOVED***Cmd: cmd1***REMOVED***,
		// only labels
		***REMOVED***Labels: labels1***REMOVED***: ***REMOVED***Labels: labels1***REMOVED***,
		// only exposedPorts
		***REMOVED***ExposedPorts: ports1***REMOVED***: ***REMOVED***ExposedPorts: ports1***REMOVED***,
		// only entrypoints
		***REMOVED***Entrypoint: entrypoint1***REMOVED***: ***REMOVED***Entrypoint: entrypoint1***REMOVED***,
		// only volumes
		***REMOVED***Volumes: volumes1***REMOVED***: ***REMOVED***Volumes: volumes1***REMOVED***,
	***REMOVED***
	differentConfigs := map[*container.Config]*container.Config***REMOVED***
		nil: nil,
		***REMOVED***
			Hostname:   "host1",
			Domainname: "domain1",
			Image:      "image1",
			User:       "user1",
		***REMOVED***: ***REMOVED***
			Hostname:   "host1",
			Domainname: "domain1",
			Image:      "image1",
			User:       "user2",
		***REMOVED***,
		// only OpenStdin
		***REMOVED***OpenStdin: false***REMOVED***: ***REMOVED***OpenStdin: true***REMOVED***,
		***REMOVED***OpenStdin: true***REMOVED***:  ***REMOVED***OpenStdin: false***REMOVED***,
		// only env
		***REMOVED***Env: envs1***REMOVED***: ***REMOVED***Env: envs2***REMOVED***,
		// only cmd
		***REMOVED***Cmd: cmd1***REMOVED***: ***REMOVED***Cmd: cmd2***REMOVED***,
		// not the same number of parts
		***REMOVED***Cmd: cmd1***REMOVED***: ***REMOVED***Cmd: cmd3***REMOVED***,
		// only labels
		***REMOVED***Labels: labels1***REMOVED***: ***REMOVED***Labels: labels2***REMOVED***,
		// not the same number of labels
		***REMOVED***Labels: labels1***REMOVED***: ***REMOVED***Labels: labels3***REMOVED***,
		// only exposedPorts
		***REMOVED***ExposedPorts: ports1***REMOVED***: ***REMOVED***ExposedPorts: ports2***REMOVED***,
		// not the same number of ports
		***REMOVED***ExposedPorts: ports1***REMOVED***: ***REMOVED***ExposedPorts: ports3***REMOVED***,
		// only entrypoints
		***REMOVED***Entrypoint: entrypoint1***REMOVED***: ***REMOVED***Entrypoint: entrypoint2***REMOVED***,
		// not the same number of parts
		***REMOVED***Entrypoint: entrypoint1***REMOVED***: ***REMOVED***Entrypoint: entrypoint3***REMOVED***,
		// only volumes
		***REMOVED***Volumes: volumes1***REMOVED***: ***REMOVED***Volumes: volumes2***REMOVED***,
		// not the same number of labels
		***REMOVED***Volumes: volumes1***REMOVED***: ***REMOVED***Volumes: volumes3***REMOVED***,
	***REMOVED***
	for config1, config2 := range sameConfigs ***REMOVED***
		if !compare(config1, config2) ***REMOVED***
			t.Fatalf("Compare should be true for [%v] and [%v]", config1, config2)
		***REMOVED***
	***REMOVED***
	for config1, config2 := range differentConfigs ***REMOVED***
		if compare(config1, config2) ***REMOVED***
			t.Fatalf("Compare should be false for [%v] and [%v]", config1, config2)
		***REMOVED***
	***REMOVED***
***REMOVED***
