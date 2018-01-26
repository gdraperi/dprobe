package cgroups

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func NewNetCls(root string) *netclsController ***REMOVED***
	return &netclsController***REMOVED***
		root: filepath.Join(root, string(NetCLS)),
	***REMOVED***
***REMOVED***

type netclsController struct ***REMOVED***
	root string
***REMOVED***

func (n *netclsController) Name() Name ***REMOVED***
	return NetCLS
***REMOVED***

func (n *netclsController) Path(path string) string ***REMOVED***
	return filepath.Join(n.root, path)
***REMOVED***

func (n *netclsController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	if err := os.MkdirAll(n.Path(path), defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	if resources.Network != nil && resources.Network.ClassID != nil && *resources.Network.ClassID > 0 ***REMOVED***
		return ioutil.WriteFile(
			filepath.Join(n.Path(path), "net_cls.classid"),
			[]byte(strconv.FormatUint(uint64(*resources.Network.ClassID), 10)),
			defaultFilePerm,
		)
	***REMOVED***
	return nil
***REMOVED***
