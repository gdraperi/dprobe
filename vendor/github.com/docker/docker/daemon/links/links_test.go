package links

import (
	"fmt"
	"strings"
	"testing"

	"github.com/docker/go-connections/nat"
)

// Just to make life easier
func newPortNoError(proto, port string) nat.Port ***REMOVED***
	p, _ := nat.NewPort(proto, port)
	return p
***REMOVED***

func TestLinkNaming(t *testing.T) ***REMOVED***
	ports := make(nat.PortSet)
	ports[newPortNoError("tcp", "6379")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	link := NewLink("172.0.17.3", "172.0.17.2", "/db/docker-1", nil, ports)

	rawEnv := link.ToEnv()
	env := make(map[string]string, len(rawEnv))
	for _, e := range rawEnv ***REMOVED***
		parts := strings.Split(e, "=")
		if len(parts) != 2 ***REMOVED***
			t.FailNow()
		***REMOVED***
		env[parts[0]] = parts[1]
	***REMOVED***

	value, ok := env["DOCKER_1_PORT"]

	if !ok ***REMOVED***
		t.Fatal("DOCKER_1_PORT not found in env")
	***REMOVED***

	if value != "tcp://172.0.17.2:6379" ***REMOVED***
		t.Fatalf("Expected 172.0.17.2:6379, got %s", env["DOCKER_1_PORT"])
	***REMOVED***
***REMOVED***

func TestLinkNew(t *testing.T) ***REMOVED***
	ports := make(nat.PortSet)
	ports[newPortNoError("tcp", "6379")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	link := NewLink("172.0.17.3", "172.0.17.2", "/db/docker", nil, ports)

	if link.Name != "/db/docker" ***REMOVED***
		t.Fail()
	***REMOVED***
	if link.ParentIP != "172.0.17.3" ***REMOVED***
		t.Fail()
	***REMOVED***
	if link.ChildIP != "172.0.17.2" ***REMOVED***
		t.Fail()
	***REMOVED***
	for _, p := range link.Ports ***REMOVED***
		if p != newPortNoError("tcp", "6379") ***REMOVED***
			t.Fail()
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestLinkEnv(t *testing.T) ***REMOVED***
	ports := make(nat.PortSet)
	ports[newPortNoError("tcp", "6379")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	link := NewLink("172.0.17.3", "172.0.17.2", "/db/docker", []string***REMOVED***"PASSWORD=gordon"***REMOVED***, ports)

	rawEnv := link.ToEnv()
	env := make(map[string]string, len(rawEnv))
	for _, e := range rawEnv ***REMOVED***
		parts := strings.Split(e, "=")
		if len(parts) != 2 ***REMOVED***
			t.FailNow()
		***REMOVED***
		env[parts[0]] = parts[1]
	***REMOVED***
	if env["DOCKER_PORT"] != "tcp://172.0.17.2:6379" ***REMOVED***
		t.Fatalf("Expected 172.0.17.2:6379, got %s", env["DOCKER_PORT"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP"] != "tcp://172.0.17.2:6379" ***REMOVED***
		t.Fatalf("Expected tcp://172.0.17.2:6379, got %s", env["DOCKER_PORT_6379_TCP"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_PROTO"] != "tcp" ***REMOVED***
		t.Fatalf("Expected tcp, got %s", env["DOCKER_PORT_6379_TCP_PROTO"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_ADDR"] != "172.0.17.2" ***REMOVED***
		t.Fatalf("Expected 172.0.17.2, got %s", env["DOCKER_PORT_6379_TCP_ADDR"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_PORT"] != "6379" ***REMOVED***
		t.Fatalf("Expected 6379, got %s", env["DOCKER_PORT_6379_TCP_PORT"])
	***REMOVED***
	if env["DOCKER_NAME"] != "/db/docker" ***REMOVED***
		t.Fatalf("Expected /db/docker, got %s", env["DOCKER_NAME"])
	***REMOVED***
	if env["DOCKER_ENV_PASSWORD"] != "gordon" ***REMOVED***
		t.Fatalf("Expected gordon, got %s", env["DOCKER_ENV_PASSWORD"])
	***REMOVED***
***REMOVED***

func TestLinkMultipleEnv(t *testing.T) ***REMOVED***
	ports := make(nat.PortSet)
	ports[newPortNoError("tcp", "6379")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports[newPortNoError("tcp", "6380")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports[newPortNoError("tcp", "6381")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	link := NewLink("172.0.17.3", "172.0.17.2", "/db/docker", []string***REMOVED***"PASSWORD=gordon"***REMOVED***, ports)

	rawEnv := link.ToEnv()
	env := make(map[string]string, len(rawEnv))
	for _, e := range rawEnv ***REMOVED***
		parts := strings.Split(e, "=")
		if len(parts) != 2 ***REMOVED***
			t.FailNow()
		***REMOVED***
		env[parts[0]] = parts[1]
	***REMOVED***
	if env["DOCKER_PORT"] != "tcp://172.0.17.2:6379" ***REMOVED***
		t.Fatalf("Expected 172.0.17.2:6379, got %s", env["DOCKER_PORT"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_START"] != "tcp://172.0.17.2:6379" ***REMOVED***
		t.Fatalf("Expected tcp://172.0.17.2:6379, got %s", env["DOCKER_PORT_6379_TCP_START"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_END"] != "tcp://172.0.17.2:6381" ***REMOVED***
		t.Fatalf("Expected tcp://172.0.17.2:6381, got %s", env["DOCKER_PORT_6379_TCP_END"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_PROTO"] != "tcp" ***REMOVED***
		t.Fatalf("Expected tcp, got %s", env["DOCKER_PORT_6379_TCP_PROTO"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_ADDR"] != "172.0.17.2" ***REMOVED***
		t.Fatalf("Expected 172.0.17.2, got %s", env["DOCKER_PORT_6379_TCP_ADDR"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_PORT_START"] != "6379" ***REMOVED***
		t.Fatalf("Expected 6379, got %s", env["DOCKER_PORT_6379_TCP_PORT_START"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_PORT_END"] != "6381" ***REMOVED***
		t.Fatalf("Expected 6381, got %s", env["DOCKER_PORT_6379_TCP_PORT_END"])
	***REMOVED***
	if env["DOCKER_NAME"] != "/db/docker" ***REMOVED***
		t.Fatalf("Expected /db/docker, got %s", env["DOCKER_NAME"])
	***REMOVED***
	if env["DOCKER_ENV_PASSWORD"] != "gordon" ***REMOVED***
		t.Fatalf("Expected gordon, got %s", env["DOCKER_ENV_PASSWORD"])
	***REMOVED***
***REMOVED***

func TestLinkPortRangeEnv(t *testing.T) ***REMOVED***
	ports := make(nat.PortSet)
	ports[newPortNoError("tcp", "6379")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports[newPortNoError("tcp", "6380")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	ports[newPortNoError("tcp", "6381")] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	link := NewLink("172.0.17.3", "172.0.17.2", "/db/docker", []string***REMOVED***"PASSWORD=gordon"***REMOVED***, ports)

	rawEnv := link.ToEnv()
	env := make(map[string]string, len(rawEnv))
	for _, e := range rawEnv ***REMOVED***
		parts := strings.Split(e, "=")
		if len(parts) != 2 ***REMOVED***
			t.FailNow()
		***REMOVED***
		env[parts[0]] = parts[1]
	***REMOVED***

	if env["DOCKER_PORT"] != "tcp://172.0.17.2:6379" ***REMOVED***
		t.Fatalf("Expected 172.0.17.2:6379, got %s", env["DOCKER_PORT"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_START"] != "tcp://172.0.17.2:6379" ***REMOVED***
		t.Fatalf("Expected tcp://172.0.17.2:6379, got %s", env["DOCKER_PORT_6379_TCP_START"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_END"] != "tcp://172.0.17.2:6381" ***REMOVED***
		t.Fatalf("Expected tcp://172.0.17.2:6381, got %s", env["DOCKER_PORT_6379_TCP_END"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_PROTO"] != "tcp" ***REMOVED***
		t.Fatalf("Expected tcp, got %s", env["DOCKER_PORT_6379_TCP_PROTO"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_ADDR"] != "172.0.17.2" ***REMOVED***
		t.Fatalf("Expected 172.0.17.2, got %s", env["DOCKER_PORT_6379_TCP_ADDR"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_PORT_START"] != "6379" ***REMOVED***
		t.Fatalf("Expected 6379, got %s", env["DOCKER_PORT_6379_TCP_PORT_START"])
	***REMOVED***
	if env["DOCKER_PORT_6379_TCP_PORT_END"] != "6381" ***REMOVED***
		t.Fatalf("Expected 6381, got %s", env["DOCKER_PORT_6379_TCP_PORT_END"])
	***REMOVED***
	if env["DOCKER_NAME"] != "/db/docker" ***REMOVED***
		t.Fatalf("Expected /db/docker, got %s", env["DOCKER_NAME"])
	***REMOVED***
	if env["DOCKER_ENV_PASSWORD"] != "gordon" ***REMOVED***
		t.Fatalf("Expected gordon, got %s", env["DOCKER_ENV_PASSWORD"])
	***REMOVED***
	for _, i := range []int***REMOVED***6379, 6380, 6381***REMOVED*** ***REMOVED***
		tcpaddr := fmt.Sprintf("DOCKER_PORT_%d_TCP_ADDR", i)
		tcpport := fmt.Sprintf("DOCKER_PORT_%d_TCP_PORT", i)
		tcpproto := fmt.Sprintf("DOCKER_PORT_%d_TCP_PROTO", i)
		tcp := fmt.Sprintf("DOCKER_PORT_%d_TCP", i)
		if env[tcpaddr] != "172.0.17.2" ***REMOVED***
			t.Fatalf("Expected env %s  = 172.0.17.2, got %s", tcpaddr, env[tcpaddr])
		***REMOVED***
		if env[tcpport] != fmt.Sprintf("%d", i) ***REMOVED***
			t.Fatalf("Expected env %s  = %d, got %s", tcpport, i, env[tcpport])
		***REMOVED***
		if env[tcpproto] != "tcp" ***REMOVED***
			t.Fatalf("Expected env %s  = tcp, got %s", tcpproto, env[tcpproto])
		***REMOVED***
		if env[tcp] != fmt.Sprintf("tcp://172.0.17.2:%d", i) ***REMOVED***
			t.Fatalf("Expected env %s  = tcp://172.0.17.2:%d, got %s", tcp, i, env[tcp])
		***REMOVED***
	***REMOVED***
***REMOVED***
