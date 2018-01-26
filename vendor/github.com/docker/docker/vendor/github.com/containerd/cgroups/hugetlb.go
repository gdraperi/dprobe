package cgroups

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func NewHugetlb(root string) (*hugetlbController, error) ***REMOVED***
	sizes, err := hugePageSizes()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &hugetlbController***REMOVED***
		root:  filepath.Join(root, string(Hugetlb)),
		sizes: sizes,
	***REMOVED***, nil
***REMOVED***

type hugetlbController struct ***REMOVED***
	root  string
	sizes []string
***REMOVED***

func (h *hugetlbController) Name() Name ***REMOVED***
	return Hugetlb
***REMOVED***

func (h *hugetlbController) Path(path string) string ***REMOVED***
	return filepath.Join(h.root, path)
***REMOVED***

func (h *hugetlbController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	if err := os.MkdirAll(h.Path(path), defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, limit := range resources.HugepageLimits ***REMOVED***
		if err := ioutil.WriteFile(
			filepath.Join(h.Path(path), strings.Join([]string***REMOVED***"hugetlb", limit.Pagesize, "limit_in_bytes"***REMOVED***, ".")),
			[]byte(strconv.FormatUint(limit.Limit, 10)),
			defaultFilePerm,
		); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (h *hugetlbController) Stat(path string, stats *Metrics) error ***REMOVED***
	for _, size := range h.sizes ***REMOVED***
		s, err := h.readSizeStat(path, size)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		stats.Hugetlb = append(stats.Hugetlb, s)
	***REMOVED***
	return nil
***REMOVED***

func (h *hugetlbController) readSizeStat(path, size string) (*HugetlbStat, error) ***REMOVED***
	s := HugetlbStat***REMOVED***
		Pagesize: size,
	***REMOVED***
	for _, t := range []struct ***REMOVED***
		name  string
		value *uint64
	***REMOVED******REMOVED***
		***REMOVED***
			name:  "usage_in_bytes",
			value: &s.Usage,
		***REMOVED***,
		***REMOVED***
			name:  "max_usage_in_bytes",
			value: &s.Max,
		***REMOVED***,
		***REMOVED***
			name:  "failcnt",
			value: &s.Failcnt,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		v, err := readUint(filepath.Join(h.Path(path), strings.Join([]string***REMOVED***"hugetlb", size, t.name***REMOVED***, ".")))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		*t.value = v
	***REMOVED***
	return &s, nil
***REMOVED***
