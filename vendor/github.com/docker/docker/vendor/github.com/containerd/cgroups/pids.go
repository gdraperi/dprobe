package cgroups

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func NewPids(root string) *pidsController ***REMOVED***
	return &pidsController***REMOVED***
		root: filepath.Join(root, string(Pids)),
	***REMOVED***
***REMOVED***

type pidsController struct ***REMOVED***
	root string
***REMOVED***

func (p *pidsController) Name() Name ***REMOVED***
	return Pids
***REMOVED***

func (p *pidsController) Path(path string) string ***REMOVED***
	return filepath.Join(p.root, path)
***REMOVED***

func (p *pidsController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	if err := os.MkdirAll(p.Path(path), defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	if resources.Pids != nil && resources.Pids.Limit > 0 ***REMOVED***
		return ioutil.WriteFile(
			filepath.Join(p.Path(path), "pids.max"),
			[]byte(strconv.FormatInt(resources.Pids.Limit, 10)),
			defaultFilePerm,
		)
	***REMOVED***
	return nil
***REMOVED***

func (p *pidsController) Update(path string, resources *specs.LinuxResources) error ***REMOVED***
	return p.Create(path, resources)
***REMOVED***

func (p *pidsController) Stat(path string, stats *Metrics) error ***REMOVED***
	current, err := readUint(filepath.Join(p.Path(path), "pids.current"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var max uint64
	maxData, err := ioutil.ReadFile(filepath.Join(p.Path(path), "pids.max"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if maxS := strings.TrimSpace(string(maxData)); maxS != "max" ***REMOVED***
		if max, err = parseUint(maxS, 10, 64); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	stats.Pids = &PidsStat***REMOVED***
		Current: current,
		Limit:   max,
	***REMOVED***
	return nil
***REMOVED***
