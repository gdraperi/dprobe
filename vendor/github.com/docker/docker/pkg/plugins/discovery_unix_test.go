// +build !windows

package plugins

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalSocket(t *testing.T) ***REMOVED***
	// TODO Windows: Enable a similar version for Windows named pipes
	tmpdir, unregister := Setup(t)
	defer unregister()

	cases := []string***REMOVED***
		filepath.Join(tmpdir, "echo.sock"),
		filepath.Join(tmpdir, "echo", "echo.sock"),
	***REMOVED***

	for _, c := range cases ***REMOVED***
		if err := os.MkdirAll(filepath.Dir(c), 0755); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		l, err := net.Listen("unix", c)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		r := newLocalRegistry()
		p, err := r.Plugin("echo")
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		pp, err := r.Plugin("echo")
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if !reflect.DeepEqual(p, pp) ***REMOVED***
			t.Fatalf("Expected %v, was %v\n", p, pp)
		***REMOVED***

		if p.name != "echo" ***REMOVED***
			t.Fatalf("Expected plugin `echo`, got %s\n", p.name)
		***REMOVED***

		addr := fmt.Sprintf("unix://%s", c)
		if p.Addr != addr ***REMOVED***
			t.Fatalf("Expected plugin addr `%s`, got %s\n", addr, p.Addr)
		***REMOVED***
		if !p.TLSConfig.InsecureSkipVerify ***REMOVED***
			t.Fatalf("Expected TLS verification to be skipped")
		***REMOVED***
		l.Close()
	***REMOVED***
***REMOVED***

func TestScan(t *testing.T) ***REMOVED***
	tmpdir, unregister := Setup(t)
	defer unregister()

	pluginNames, err := Scan()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if pluginNames != nil ***REMOVED***
		t.Fatal("Plugin names should be empty.")
	***REMOVED***

	path := filepath.Join(tmpdir, "echo.spec")
	addr := "unix://var/lib/docker/plugins/echo.sock"
	name := "echo"

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = ioutil.WriteFile(path, []byte(addr), 0644)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	r := newLocalRegistry()
	p, err := r.Plugin(name)
	require.NoError(t, err)

	pluginNamesNotEmpty, err := Scan()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if p.Name() != pluginNamesNotEmpty[0] ***REMOVED***
		t.Fatalf("Unable to scan plugin with name %s", p.name)
	***REMOVED***
***REMOVED***
