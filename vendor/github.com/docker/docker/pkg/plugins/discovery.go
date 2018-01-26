package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// ErrNotFound plugin not found
	ErrNotFound = errors.New("plugin not found")
	socketsPath = "/run/docker/plugins"
)

// localRegistry defines a registry that is local (using unix socket).
type localRegistry struct***REMOVED******REMOVED***

func newLocalRegistry() localRegistry ***REMOVED***
	return localRegistry***REMOVED******REMOVED***
***REMOVED***

// Scan scans all the plugin paths and returns all the names it found
func Scan() ([]string, error) ***REMOVED***
	var names []string
	if err := filepath.Walk(socketsPath, func(path string, fi os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return nil
		***REMOVED***

		if fi.Mode()&os.ModeSocket != 0 ***REMOVED***
			name := strings.TrimSuffix(fi.Name(), filepath.Ext(fi.Name()))
			names = append(names, name)
		***REMOVED***
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, path := range specsPaths ***REMOVED***
		if err := filepath.Walk(path, func(p string, fi os.FileInfo, err error) error ***REMOVED***
			if err != nil || fi.IsDir() ***REMOVED***
				return nil
			***REMOVED***
			name := strings.TrimSuffix(fi.Name(), filepath.Ext(fi.Name()))
			names = append(names, name)
			return nil
		***REMOVED***); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return names, nil
***REMOVED***

// Plugin returns the plugin registered with the given name (or returns an error).
func (l *localRegistry) Plugin(name string) (*Plugin, error) ***REMOVED***
	socketpaths := pluginPaths(socketsPath, name, ".sock")

	for _, p := range socketpaths ***REMOVED***
		if fi, err := os.Stat(p); err == nil && fi.Mode()&os.ModeSocket != 0 ***REMOVED***
			return NewLocalPlugin(name, "unix://"+p), nil
		***REMOVED***
	***REMOVED***

	var txtspecpaths []string
	for _, p := range specsPaths ***REMOVED***
		txtspecpaths = append(txtspecpaths, pluginPaths(p, name, ".spec")...)
		txtspecpaths = append(txtspecpaths, pluginPaths(p, name, ".json")...)
	***REMOVED***

	for _, p := range txtspecpaths ***REMOVED***
		if _, err := os.Stat(p); err == nil ***REMOVED***
			if strings.HasSuffix(p, ".json") ***REMOVED***
				return readPluginJSONInfo(name, p)
			***REMOVED***
			return readPluginInfo(name, p)
		***REMOVED***
	***REMOVED***
	return nil, ErrNotFound
***REMOVED***

func readPluginInfo(name, path string) (*Plugin, error) ***REMOVED***
	content, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	addr := strings.TrimSpace(string(content))

	u, err := url.Parse(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(u.Scheme) == 0 ***REMOVED***
		return nil, fmt.Errorf("Unknown protocol")
	***REMOVED***

	return NewLocalPlugin(name, addr), nil
***REMOVED***

func readPluginJSONInfo(name, path string) (*Plugin, error) ***REMOVED***
	f, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	var p Plugin
	if err := json.NewDecoder(f).Decode(&p); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	p.name = name
	if p.TLSConfig != nil && len(p.TLSConfig.CAFile) == 0 ***REMOVED***
		p.TLSConfig.InsecureSkipVerify = true
	***REMOVED***
	p.activateWait = sync.NewCond(&sync.Mutex***REMOVED******REMOVED***)

	return &p, nil
***REMOVED***

func pluginPaths(base, name, ext string) []string ***REMOVED***
	return []string***REMOVED***
		filepath.Join(base, name+ext),
		filepath.Join(base, name, name+ext),
	***REMOVED***
***REMOVED***
