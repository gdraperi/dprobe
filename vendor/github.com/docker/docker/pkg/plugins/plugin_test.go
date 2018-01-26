package plugins

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/pkg/plugins/transport"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/stretchr/testify/assert"
)

const (
	fruitPlugin     = "fruit"
	fruitImplements = "apple"
)

// regression test for deadlock in handlers
func TestPluginAddHandler(t *testing.T) ***REMOVED***
	// make a plugin which is pre-activated
	p := &Plugin***REMOVED***activateWait: sync.NewCond(&sync.Mutex***REMOVED******REMOVED***)***REMOVED***
	p.Manifest = &Manifest***REMOVED***Implements: []string***REMOVED***"bananas"***REMOVED******REMOVED***
	storage.plugins["qwerty"] = p

	testActive(t, p)
	Handle("bananas", func(_ string, _ *Client) ***REMOVED******REMOVED***)
	testActive(t, p)
***REMOVED***

func TestPluginWaitBadPlugin(t *testing.T) ***REMOVED***
	p := &Plugin***REMOVED***activateWait: sync.NewCond(&sync.Mutex***REMOVED******REMOVED***)***REMOVED***
	p.activateErr = errors.New("some junk happened")
	testActive(t, p)
***REMOVED***

func testActive(t *testing.T, p *Plugin) ***REMOVED***
	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		p.waitActive()
		close(done)
	***REMOVED***()

	select ***REMOVED***
	case <-time.After(100 * time.Millisecond):
		_, f, l, _ := runtime.Caller(1)
		t.Fatalf("%s:%d: deadlock in waitActive", filepath.Base(f), l)
	case <-done:
	***REMOVED***

***REMOVED***

func TestGet(t *testing.T) ***REMOVED***
	p := &Plugin***REMOVED***name: fruitPlugin, activateWait: sync.NewCond(&sync.Mutex***REMOVED******REMOVED***)***REMOVED***
	p.Manifest = &Manifest***REMOVED***Implements: []string***REMOVED***fruitImplements***REMOVED******REMOVED***
	storage.plugins[fruitPlugin] = p

	plugin, err := Get(fruitPlugin, fruitImplements)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if p.Name() != plugin.Name() ***REMOVED***
		t.Fatalf("No matching plugin with name %s found", plugin.Name())
	***REMOVED***
	if plugin.Client() != nil ***REMOVED***
		t.Fatal("expected nil Client but found one")
	***REMOVED***
	if !plugin.IsV1() ***REMOVED***
		t.Fatal("Expected true for V1 plugin")
	***REMOVED***

	// check negative case where plugin fruit doesn't implement banana
	_, err = Get("fruit", "banana")
	assert.Equal(t, err, ErrNotImplements)

	// check negative case where plugin vegetable doesn't exist
	_, err = Get("vegetable", "potato")
	assert.Equal(t, err, ErrNotFound)

***REMOVED***

func TestPluginWithNoManifest(t *testing.T) ***REMOVED***
	addr := setupRemotePluginServer()
	defer teardownRemotePluginServer()

	m := Manifest***REMOVED***[]string***REMOVED***fruitImplements***REMOVED******REMOVED***
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	mux.HandleFunc("/Plugin.Activate", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.Method != "POST" ***REMOVED***
			t.Fatalf("Expected POST, got %s\n", r.Method)
		***REMOVED***

		header := w.Header()
		header.Set("Content-Type", transport.VersionMimetype)

		io.Copy(w, &buf)
	***REMOVED***)

	p := &Plugin***REMOVED***
		name:         fruitPlugin,
		activateWait: sync.NewCond(&sync.Mutex***REMOVED******REMOVED***),
		Addr:         addr,
		TLSConfig:    &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***,
	***REMOVED***
	storage.plugins[fruitPlugin] = p

	plugin, err := Get(fruitPlugin, fruitImplements)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if p.Name() != plugin.Name() ***REMOVED***
		t.Fatalf("No matching plugin with name %s found", plugin.Name())
	***REMOVED***
***REMOVED***

func TestGetAll(t *testing.T) ***REMOVED***
	tmpdir, unregister := Setup(t)
	defer unregister()

	p := filepath.Join(tmpdir, "example.json")
	spec := `***REMOVED***
	"Name": "example",
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
	plugin.Manifest = &Manifest***REMOVED***Implements: []string***REMOVED***"apple"***REMOVED******REMOVED***
	storage.plugins["example"] = plugin

	fetchedPlugins, err := GetAll("apple")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if fetchedPlugins[0].Name() != plugin.Name() ***REMOVED***
		t.Fatalf("Expected to get plugin with name %s", plugin.Name())
	***REMOVED***
***REMOVED***
