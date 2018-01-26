// +build linux

package gelf

import (
	"net"
	"testing"

	"github.com/docker/docker/daemon/logger"
)

// Validate parseAddress
func TestParseAddress(t *testing.T) ***REMOVED***
	url, err := parseAddress("udp://127.0.0.1:12201")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if url.String() != "udp://127.0.0.1:12201" ***REMOVED***
		t.Fatalf("Expected address udp://127.0.0.1:12201, got %s", url.String())
	***REMOVED***

	_, err = parseAddress("127.0.0.1:12201")
	if err == nil ***REMOVED***
		t.Fatal("Expected error requiring protocol")
	***REMOVED***

	_, err = parseAddress("http://127.0.0.1:12201")
	if err == nil ***REMOVED***
		t.Fatal("Expected error restricting protocol")
	***REMOVED***
***REMOVED***

// Validate TCP options
func TestTCPValidateLogOpt(t *testing.T) ***REMOVED***
	err := ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address": "tcp://127.0.0.1:12201",
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("Expected TCP to be supported")
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address":           "tcp://127.0.0.1:12201",
		"gelf-compression-level": "9",
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected TCP to reject compression level")
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address":          "tcp://127.0.0.1:12201",
		"gelf-compression-type": "gzip",
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected TCP to reject compression type")
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address":             "tcp://127.0.0.1:12201",
		"gelf-tcp-max-reconnect":   "5",
		"gelf-tcp-reconnect-delay": "10",
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal("Expected TCP reconnect to be a valid parameters")
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address":             "tcp://127.0.0.1:12201",
		"gelf-tcp-max-reconnect":   "-1",
		"gelf-tcp-reconnect-delay": "-3",
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected negative TCP reconnect to be rejected")
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address":             "tcp://127.0.0.1:12201",
		"gelf-tcp-max-reconnect":   "invalid",
		"gelf-tcp-reconnect-delay": "invalid",
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected TCP reconnect to be required to be an int")
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address":             "udp://127.0.0.1:12201",
		"gelf-tcp-max-reconnect":   "1",
		"gelf-tcp-reconnect-delay": "3",
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected TCP reconnect to be invalid for UDP")
	***REMOVED***
***REMOVED***

// Validate UDP options
func TestUDPValidateLogOpt(t *testing.T) ***REMOVED***
	err := ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address":           "udp://127.0.0.1:12201",
		"tag":                    "testtag",
		"labels":                 "testlabel",
		"env":                    "testenv",
		"env-regex":              "testenv-regex",
		"gelf-compression-level": "9",
		"gelf-compression-type":  "gzip",
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address":           "udp://127.0.0.1:12201",
		"gelf-compression-level": "ultra",
		"gelf-compression-type":  "zlib",
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected compression level error")
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"gelf-address":          "udp://127.0.0.1:12201",
		"gelf-compression-type": "rar",
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected compression type error")
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"invalid": "invalid",
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected unknown option error")
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED******REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected required parameter error")
	***REMOVED***
***REMOVED***

// Validate newGELFTCPWriter
func TestNewGELFTCPWriter(t *testing.T) ***REMOVED***
	address := "127.0.0.1:0"
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	url := "tcp://" + listener.Addr().String()
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			"gelf-address":             url,
			"gelf-tcp-max-reconnect":   "0",
			"gelf-tcp-reconnect-delay": "0",
			"tag": "***REMOVED******REMOVED***.ID***REMOVED******REMOVED***",
		***REMOVED***,
		ContainerID: "12345678901234567890",
	***REMOVED***

	writer, err := newGELFTCPWriter(listener.Addr().String(), info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = writer.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = listener.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Validate newGELFUDPWriter
func TestNewGELFUDPWriter(t *testing.T) ***REMOVED***
	address := "127.0.0.1:0"
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			"gelf-address":           "udp://127.0.0.1:0",
			"gelf-compression-level": "5",
			"gelf-compression-type":  "gzip",
		***REMOVED***,
	***REMOVED***

	writer, err := newGELFUDPWriter(address, info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	writer.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Validate New for TCP
func TestNewTCP(t *testing.T) ***REMOVED***
	address := "127.0.0.1:0"
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	url := "tcp://" + listener.Addr().String()
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			"gelf-address":             url,
			"gelf-tcp-max-reconnect":   "0",
			"gelf-tcp-reconnect-delay": "0",
		***REMOVED***,
		ContainerID: "12345678901234567890",
	***REMOVED***

	logger, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = logger.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = listener.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Validate New for UDP
func TestNewUDP(t *testing.T) ***REMOVED***
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			"gelf-address":           "udp://127.0.0.1:0",
			"gelf-compression-level": "5",
			"gelf-compression-type":  "gzip",
		***REMOVED***,
		ContainerID: "12345678901234567890",
	***REMOVED***

	logger, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = logger.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
