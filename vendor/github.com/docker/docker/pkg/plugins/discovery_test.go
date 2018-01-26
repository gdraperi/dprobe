package plugins

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T) (string, func()) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "docker-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	backup := socketsPath
	socketsPath = tmpdir
	specsPaths = []string***REMOVED***tmpdir***REMOVED***

	return tmpdir, func() ***REMOVED***
		socketsPath = backup
		os.RemoveAll(tmpdir)
	***REMOVED***
***REMOVED***

func TestFileSpecPlugin(t *testing.T) ***REMOVED***
	tmpdir, unregister := Setup(t)
	defer unregister()

	cases := []struct ***REMOVED***
		path string
		name string
		addr string
		fail bool
	***REMOVED******REMOVED***
		// TODO Windows: Factor out the unix:// variants.
		***REMOVED***filepath.Join(tmpdir, "echo.spec"), "echo", "unix://var/lib/docker/plugins/echo.sock", false***REMOVED***,
		***REMOVED***filepath.Join(tmpdir, "echo", "echo.spec"), "echo", "unix://var/lib/docker/plugins/echo.sock", false***REMOVED***,
		***REMOVED***filepath.Join(tmpdir, "foo.spec"), "foo", "tcp://localhost:8080", false***REMOVED***,
		***REMOVED***filepath.Join(tmpdir, "foo", "foo.spec"), "foo", "tcp://localhost:8080", false***REMOVED***,
		***REMOVED***filepath.Join(tmpdir, "bar.spec"), "bar", "localhost:8080", true***REMOVED***, // unknown transport
	***REMOVED***

	for _, c := range cases ***REMOVED***
		if err := os.MkdirAll(filepath.Dir(c.path), 0755); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if err := ioutil.WriteFile(c.path, []byte(c.addr), 0644); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		r := newLocalRegistry()
		p, err := r.Plugin(c.name)
		if c.fail && err == nil ***REMOVED***
			continue
		***REMOVED***

		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		if p.name != c.name ***REMOVED***
			t.Fatalf("Expected plugin `%s`, got %s\n", c.name, p.name)
		***REMOVED***

		if p.Addr != c.addr ***REMOVED***
			t.Fatalf("Expected plugin addr `%s`, got %s\n", c.addr, p.Addr)
		***REMOVED***

		if !p.TLSConfig.InsecureSkipVerify ***REMOVED***
			t.Fatalf("Expected TLS verification to be skipped")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFileJSONSpecPlugin(t *testing.T) ***REMOVED***
	tmpdir, unregister := Setup(t)
	defer unregister()

	p := filepath.Join(tmpdir, "example.json")
	spec := `***REMOVED***
  "Name": "plugin-example",
  "Addr": "https://example.com/docker/plugin",
  "TLSConfig": ***REMOVED***
    "CAFile": "/usr/shared/docker/certs/example-ca.pem",
    "CertFile": "/usr/shared/docker/certs/example-cert.pem",
    "KeyFile": "/usr/shared/docker/certs/example-key.pem"
	***REMOVED***
***REMOVED***`

	if err := ioutil.WriteFile(p, []byte(spec), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	r := newLocalRegistry()
	plugin, err := r.Plugin("example")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected, actual := "example", plugin.name; expected != actual ***REMOVED***
		t.Fatalf("Expected plugin %q, got %s\n", expected, actual)
	***REMOVED***

	if plugin.Addr != "https://example.com/docker/plugin" ***REMOVED***
		t.Fatalf("Expected plugin addr `https://example.com/docker/plugin`, got %s\n", plugin.Addr)
	***REMOVED***

	if plugin.TLSConfig.CAFile != "/usr/shared/docker/certs/example-ca.pem" ***REMOVED***
		t.Fatalf("Expected plugin CA `/usr/shared/docker/certs/example-ca.pem`, got %s\n", plugin.TLSConfig.CAFile)
	***REMOVED***

	if plugin.TLSConfig.CertFile != "/usr/shared/docker/certs/example-cert.pem" ***REMOVED***
		t.Fatalf("Expected plugin Certificate `/usr/shared/docker/certs/example-cert.pem`, got %s\n", plugin.TLSConfig.CertFile)
	***REMOVED***

	if plugin.TLSConfig.KeyFile != "/usr/shared/docker/certs/example-key.pem" ***REMOVED***
		t.Fatalf("Expected plugin Key `/usr/shared/docker/certs/example-key.pem`, got %s\n", plugin.TLSConfig.KeyFile)
	***REMOVED***
***REMOVED***

func TestFileJSONSpecPluginWithoutTLSConfig(t *testing.T) ***REMOVED***
	tmpdir, unregister := Setup(t)
	defer unregister()

	p := filepath.Join(tmpdir, "example.json")
	spec := `***REMOVED***
  "Name": "plugin-example",
  "Addr": "https://example.com/docker/plugin"
***REMOVED***`

	if err := ioutil.WriteFile(p, []byte(spec), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	r := newLocalRegistry()
	plugin, err := r.Plugin("example")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected, actual := "example", plugin.name; expected != actual ***REMOVED***
		t.Fatalf("Expected plugin %q, got %s\n", expected, actual)
	***REMOVED***

	if plugin.Addr != "https://example.com/docker/plugin" ***REMOVED***
		t.Fatalf("Expected plugin addr `https://example.com/docker/plugin`, got %s\n", plugin.Addr)
	***REMOVED***

	if plugin.TLSConfig != nil ***REMOVED***
		t.Fatalf("Expected plugin TLSConfig nil, got %v\n", plugin.TLSConfig)
	***REMOVED***
***REMOVED***
