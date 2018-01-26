package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func NewNetPrio(root string) *netprioController ***REMOVED***
	return &netprioController***REMOVED***
		root: filepath.Join(root, string(NetPrio)),
	***REMOVED***
***REMOVED***

type netprioController struct ***REMOVED***
	root string
***REMOVED***

func (n *netprioController) Name() Name ***REMOVED***
	return NetPrio
***REMOVED***

func (n *netprioController) Path(path string) string ***REMOVED***
	return filepath.Join(n.root, path)
***REMOVED***

func (n *netprioController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	if err := os.MkdirAll(n.Path(path), defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	if resources.Network != nil ***REMOVED***
		for _, prio := range resources.Network.Priorities ***REMOVED***
			if err := ioutil.WriteFile(
				filepath.Join(n.Path(path), "net_prio_ifpriomap"),
				formatPrio(prio.Name, prio.Priority),
				defaultFilePerm,
			); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func formatPrio(name string, prio uint32) []byte ***REMOVED***
	return []byte(fmt.Sprintf("%s %d", name, prio))
***REMOVED***
