package opts

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseHost(t *testing.T) ***REMOVED***
	invalid := []string***REMOVED***
		"something with spaces",
		"://",
		"unknown://",
		"tcp://:port",
		"tcp://invalid:port",
	***REMOVED***

	valid := map[string]string***REMOVED***
		"":                         DefaultHost,
		" ":                        DefaultHost,
		"  ":                       DefaultHost,
		"fd://":                    "fd://",
		"fd://something":           "fd://something",
		"tcp://host:":              fmt.Sprintf("tcp://host:%d", DefaultHTTPPort),
		"tcp://":                   DefaultTCPHost,
		"tcp://:2375":              fmt.Sprintf("tcp://%s:2375", DefaultHTTPHost),
		"tcp://:2376":              fmt.Sprintf("tcp://%s:2376", DefaultHTTPHost),
		"tcp://0.0.0.0:8080":       "tcp://0.0.0.0:8080",
		"tcp://192.168.0.0:12000":  "tcp://192.168.0.0:12000",
		"tcp://192.168:8080":       "tcp://192.168:8080",
		"tcp://0.0.0.0:1234567890": "tcp://0.0.0.0:1234567890", // yeah it's valid :P
		" tcp://:7777/path ":       fmt.Sprintf("tcp://%s:7777/path", DefaultHTTPHost),
		"tcp://docker.com:2375":    "tcp://docker.com:2375",
		"unix://":                  "unix://" + DefaultUnixSocket,
		"unix://path/to/socket":    "unix://path/to/socket",
		"npipe://":                 "npipe://" + DefaultNamedPipe,
		"npipe:////./pipe/foo":     "npipe:////./pipe/foo",
	***REMOVED***

	for _, value := range invalid ***REMOVED***
		if _, err := ParseHost(false, value); err == nil ***REMOVED***
			t.Errorf("Expected an error for %v, got [nil]", value)
		***REMOVED***
	***REMOVED***

	for value, expected := range valid ***REMOVED***
		if actual, err := ParseHost(false, value); err != nil || actual != expected ***REMOVED***
			t.Errorf("Expected for %v [%v], got [%v, %v]", value, expected, actual, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseDockerDaemonHost(t *testing.T) ***REMOVED***
	invalids := map[string]string***REMOVED***

		"tcp:a.b.c.d":                   "Invalid bind address format: tcp:a.b.c.d",
		"tcp:a.b.c.d/path":              "Invalid bind address format: tcp:a.b.c.d/path",
		"udp://127.0.0.1":               "Invalid bind address format: udp://127.0.0.1",
		"udp://127.0.0.1:2375":          "Invalid bind address format: udp://127.0.0.1:2375",
		"tcp://unix:///run/docker.sock": "Invalid proto, expected tcp: unix:///run/docker.sock",
		" tcp://:7777/path ":            "Invalid bind address format:  tcp://:7777/path ",
		"":                              "Invalid bind address format: ",
	***REMOVED***
	valids := map[string]string***REMOVED***
		"0.0.0.1:":                    "tcp://0.0.0.1:2375",
		"0.0.0.1:5555":                "tcp://0.0.0.1:5555",
		"0.0.0.1:5555/path":           "tcp://0.0.0.1:5555/path",
		"[::1]:":                      "tcp://[::1]:2375",
		"[::1]:5555/path":             "tcp://[::1]:5555/path",
		"[0:0:0:0:0:0:0:1]:":          "tcp://[0:0:0:0:0:0:0:1]:2375",
		"[0:0:0:0:0:0:0:1]:5555/path": "tcp://[0:0:0:0:0:0:0:1]:5555/path",
		":6666":                   fmt.Sprintf("tcp://%s:6666", DefaultHTTPHost),
		":6666/path":              fmt.Sprintf("tcp://%s:6666/path", DefaultHTTPHost),
		"tcp://":                  DefaultTCPHost,
		"tcp://:7777":             fmt.Sprintf("tcp://%s:7777", DefaultHTTPHost),
		"tcp://:7777/path":        fmt.Sprintf("tcp://%s:7777/path", DefaultHTTPHost),
		"unix:///run/docker.sock": "unix:///run/docker.sock",
		"unix://":                 "unix://" + DefaultUnixSocket,
		"fd://":                   "fd://",
		"fd://something":          "fd://something",
		"localhost:":              "tcp://localhost:2375",
		"localhost:5555":          "tcp://localhost:5555",
		"localhost:5555/path":     "tcp://localhost:5555/path",
	***REMOVED***
	for invalidAddr, expectedError := range invalids ***REMOVED***
		if addr, err := parseDaemonHost(invalidAddr); err == nil || err.Error() != expectedError ***REMOVED***
			t.Errorf("tcp %v address expected error %q return, got %q and addr %v", invalidAddr, expectedError, err, addr)
		***REMOVED***
	***REMOVED***
	for validAddr, expectedAddr := range valids ***REMOVED***
		if addr, err := parseDaemonHost(validAddr); err != nil || addr != expectedAddr ***REMOVED***
			t.Errorf("%v -> expected %v, got (%v) addr (%v)", validAddr, expectedAddr, err, addr)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseTCP(t *testing.T) ***REMOVED***
	var (
		defaultHTTPHost = "tcp://127.0.0.1:2376"
	)
	invalids := map[string]string***REMOVED***
		"tcp:a.b.c.d":          "Invalid bind address format: tcp:a.b.c.d",
		"tcp:a.b.c.d/path":     "Invalid bind address format: tcp:a.b.c.d/path",
		"udp://127.0.0.1":      "Invalid proto, expected tcp: udp://127.0.0.1",
		"udp://127.0.0.1:2375": "Invalid proto, expected tcp: udp://127.0.0.1:2375",
	***REMOVED***
	valids := map[string]string***REMOVED***
		"":                            defaultHTTPHost,
		"tcp://":                      defaultHTTPHost,
		"0.0.0.1:":                    "tcp://0.0.0.1:2376",
		"0.0.0.1:5555":                "tcp://0.0.0.1:5555",
		"0.0.0.1:5555/path":           "tcp://0.0.0.1:5555/path",
		":6666":                       "tcp://127.0.0.1:6666",
		":6666/path":                  "tcp://127.0.0.1:6666/path",
		"tcp://:7777":                 "tcp://127.0.0.1:7777",
		"tcp://:7777/path":            "tcp://127.0.0.1:7777/path",
		"[::1]:":                      "tcp://[::1]:2376",
		"[::1]:5555":                  "tcp://[::1]:5555",
		"[::1]:5555/path":             "tcp://[::1]:5555/path",
		"[0:0:0:0:0:0:0:1]:":          "tcp://[0:0:0:0:0:0:0:1]:2376",
		"[0:0:0:0:0:0:0:1]:5555":      "tcp://[0:0:0:0:0:0:0:1]:5555",
		"[0:0:0:0:0:0:0:1]:5555/path": "tcp://[0:0:0:0:0:0:0:1]:5555/path",
		"localhost:":                  "tcp://localhost:2376",
		"localhost:5555":              "tcp://localhost:5555",
		"localhost:5555/path":         "tcp://localhost:5555/path",
	***REMOVED***
	for invalidAddr, expectedError := range invalids ***REMOVED***
		if addr, err := ParseTCPAddr(invalidAddr, defaultHTTPHost); err == nil || err.Error() != expectedError ***REMOVED***
			t.Errorf("tcp %v address expected error %v return, got %s and addr %v", invalidAddr, expectedError, err, addr)
		***REMOVED***
	***REMOVED***
	for validAddr, expectedAddr := range valids ***REMOVED***
		if addr, err := ParseTCPAddr(validAddr, defaultHTTPHost); err != nil || addr != expectedAddr ***REMOVED***
			t.Errorf("%v -> expected %v, got %v and addr %v", validAddr, expectedAddr, err, addr)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseInvalidUnixAddrInvalid(t *testing.T) ***REMOVED***
	if _, err := parseSimpleProtoAddr("unix", "tcp://127.0.0.1", "unix:///var/run/docker.sock"); err == nil || err.Error() != "Invalid proto, expected unix: tcp://127.0.0.1" ***REMOVED***
		t.Fatalf("Expected an error, got %v", err)
	***REMOVED***
	if _, err := parseSimpleProtoAddr("unix", "unix://tcp://127.0.0.1", "/var/run/docker.sock"); err == nil || err.Error() != "Invalid proto, expected unix: tcp://127.0.0.1" ***REMOVED***
		t.Fatalf("Expected an error, got %v", err)
	***REMOVED***
	if v, err := parseSimpleProtoAddr("unix", "", "/var/run/docker.sock"); err != nil || v != "unix:///var/run/docker.sock" ***REMOVED***
		t.Fatalf("Expected an %v, got %v", v, "unix:///var/run/docker.sock")
	***REMOVED***
***REMOVED***

func TestValidateExtraHosts(t *testing.T) ***REMOVED***
	valid := []string***REMOVED***
		`myhost:192.168.0.1`,
		`thathost:10.0.2.1`,
		`anipv6host:2003:ab34:e::1`,
		`ipv6local:::1`,
	***REMOVED***

	invalid := map[string]string***REMOVED***
		`myhost:192.notanipaddress.1`:  `invalid IP`,
		`thathost-nosemicolon10.0.0.1`: `bad format`,
		`anipv6host:::::1`:             `invalid IP`,
		`ipv6local:::0::`:              `invalid IP`,
	***REMOVED***

	for _, extrahost := range valid ***REMOVED***
		if _, err := ValidateExtraHost(extrahost); err != nil ***REMOVED***
			t.Fatalf("ValidateExtraHost(`"+extrahost+"`) should succeed: error %v", err)
		***REMOVED***
	***REMOVED***

	for extraHost, expectedError := range invalid ***REMOVED***
		if _, err := ValidateExtraHost(extraHost); err == nil ***REMOVED***
			t.Fatalf("ValidateExtraHost(`%q`) should have failed validation", extraHost)
		***REMOVED*** else ***REMOVED***
			if !strings.Contains(err.Error(), expectedError) ***REMOVED***
				t.Fatalf("ValidateExtraHost(`%q`) error should contain %q", extraHost, expectedError)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
