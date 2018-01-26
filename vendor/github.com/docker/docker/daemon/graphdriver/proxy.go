package graphdriver

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/plugins"
)

type graphDriverProxy struct ***REMOVED***
	name string
	p    plugingetter.CompatPlugin
	caps Capabilities
***REMOVED***

type graphDriverRequest struct ***REMOVED***
	ID         string            `json:",omitempty"`
	Parent     string            `json:",omitempty"`
	MountLabel string            `json:",omitempty"`
	StorageOpt map[string]string `json:",omitempty"`
***REMOVED***

type graphDriverResponse struct ***REMOVED***
	Err          string            `json:",omitempty"`
	Dir          string            `json:",omitempty"`
	Exists       bool              `json:",omitempty"`
	Status       [][2]string       `json:",omitempty"`
	Changes      []archive.Change  `json:",omitempty"`
	Size         int64             `json:",omitempty"`
	Metadata     map[string]string `json:",omitempty"`
	Capabilities Capabilities      `json:",omitempty"`
***REMOVED***

type graphDriverInitRequest struct ***REMOVED***
	Home    string
	Opts    []string        `json:"Opts"`
	UIDMaps []idtools.IDMap `json:"UIDMaps"`
	GIDMaps []idtools.IDMap `json:"GIDMaps"`
***REMOVED***

func (d *graphDriverProxy) Init(home string, opts []string, uidMaps, gidMaps []idtools.IDMap) error ***REMOVED***
	if !d.p.IsV1() ***REMOVED***
		if cp, ok := d.p.(plugingetter.CountedPlugin); ok ***REMOVED***
			// always acquire here, it will be cleaned up on daemon shutdown
			cp.Acquire()
		***REMOVED***
	***REMOVED***
	args := &graphDriverInitRequest***REMOVED***
		Home:    home,
		Opts:    opts,
		UIDMaps: uidMaps,
		GIDMaps: gidMaps,
	***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.Init", args, &ret); err != nil ***REMOVED***
		return err
	***REMOVED***
	if ret.Err != "" ***REMOVED***
		return errors.New(ret.Err)
	***REMOVED***
	caps, err := d.fetchCaps()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	d.caps = caps
	return nil
***REMOVED***

func (d *graphDriverProxy) fetchCaps() (Capabilities, error) ***REMOVED***
	args := &graphDriverRequest***REMOVED******REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.Capabilities", args, &ret); err != nil ***REMOVED***
		if !plugins.IsNotFound(err) ***REMOVED***
			return Capabilities***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***
	return ret.Capabilities, nil
***REMOVED***

func (d *graphDriverProxy) String() string ***REMOVED***
	return d.name
***REMOVED***

func (d *graphDriverProxy) Capabilities() Capabilities ***REMOVED***
	return d.caps
***REMOVED***

func (d *graphDriverProxy) CreateReadWrite(id, parent string, opts *CreateOpts) error ***REMOVED***
	return d.create("GraphDriver.CreateReadWrite", id, parent, opts)
***REMOVED***

func (d *graphDriverProxy) Create(id, parent string, opts *CreateOpts) error ***REMOVED***
	return d.create("GraphDriver.Create", id, parent, opts)
***REMOVED***

func (d *graphDriverProxy) create(method, id, parent string, opts *CreateOpts) error ***REMOVED***
	args := &graphDriverRequest***REMOVED***
		ID:     id,
		Parent: parent,
	***REMOVED***
	if opts != nil ***REMOVED***
		args.MountLabel = opts.MountLabel
		args.StorageOpt = opts.StorageOpt
	***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call(method, args, &ret); err != nil ***REMOVED***
		return err
	***REMOVED***
	if ret.Err != "" ***REMOVED***
		return errors.New(ret.Err)
	***REMOVED***
	return nil
***REMOVED***

func (d *graphDriverProxy) Remove(id string) error ***REMOVED***
	args := &graphDriverRequest***REMOVED***ID: id***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.Remove", args, &ret); err != nil ***REMOVED***
		return err
	***REMOVED***
	if ret.Err != "" ***REMOVED***
		return errors.New(ret.Err)
	***REMOVED***
	return nil
***REMOVED***

func (d *graphDriverProxy) Get(id, mountLabel string) (containerfs.ContainerFS, error) ***REMOVED***
	args := &graphDriverRequest***REMOVED***
		ID:         id,
		MountLabel: mountLabel,
	***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.Get", args, &ret); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var err error
	if ret.Err != "" ***REMOVED***
		err = errors.New(ret.Err)
	***REMOVED***
	return containerfs.NewLocalContainerFS(filepath.Join(d.p.BasePath(), ret.Dir)), err
***REMOVED***

func (d *graphDriverProxy) Put(id string) error ***REMOVED***
	args := &graphDriverRequest***REMOVED***ID: id***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.Put", args, &ret); err != nil ***REMOVED***
		return err
	***REMOVED***
	if ret.Err != "" ***REMOVED***
		return errors.New(ret.Err)
	***REMOVED***
	return nil
***REMOVED***

func (d *graphDriverProxy) Exists(id string) bool ***REMOVED***
	args := &graphDriverRequest***REMOVED***ID: id***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.Exists", args, &ret); err != nil ***REMOVED***
		return false
	***REMOVED***
	return ret.Exists
***REMOVED***

func (d *graphDriverProxy) Status() [][2]string ***REMOVED***
	args := &graphDriverRequest***REMOVED******REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.Status", args, &ret); err != nil ***REMOVED***
		return nil
	***REMOVED***
	return ret.Status
***REMOVED***

func (d *graphDriverProxy) GetMetadata(id string) (map[string]string, error) ***REMOVED***
	args := &graphDriverRequest***REMOVED***
		ID: id,
	***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.GetMetadata", args, &ret); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if ret.Err != "" ***REMOVED***
		return nil, errors.New(ret.Err)
	***REMOVED***
	return ret.Metadata, nil
***REMOVED***

func (d *graphDriverProxy) Cleanup() error ***REMOVED***
	if !d.p.IsV1() ***REMOVED***
		if cp, ok := d.p.(plugingetter.CountedPlugin); ok ***REMOVED***
			// always release
			defer cp.Release()
		***REMOVED***
	***REMOVED***

	args := &graphDriverRequest***REMOVED******REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.Cleanup", args, &ret); err != nil ***REMOVED***
		return nil
	***REMOVED***
	if ret.Err != "" ***REMOVED***
		return errors.New(ret.Err)
	***REMOVED***
	return nil
***REMOVED***

func (d *graphDriverProxy) Diff(id, parent string) (io.ReadCloser, error) ***REMOVED***
	args := &graphDriverRequest***REMOVED***
		ID:     id,
		Parent: parent,
	***REMOVED***
	body, err := d.p.Client().Stream("GraphDriver.Diff", args)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return body, nil
***REMOVED***

func (d *graphDriverProxy) Changes(id, parent string) ([]archive.Change, error) ***REMOVED***
	args := &graphDriverRequest***REMOVED***
		ID:     id,
		Parent: parent,
	***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.Changes", args, &ret); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if ret.Err != "" ***REMOVED***
		return nil, errors.New(ret.Err)
	***REMOVED***

	return ret.Changes, nil
***REMOVED***

func (d *graphDriverProxy) ApplyDiff(id, parent string, diff io.Reader) (int64, error) ***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().SendFile(fmt.Sprintf("GraphDriver.ApplyDiff?id=%s&parent=%s", id, parent), diff, &ret); err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	if ret.Err != "" ***REMOVED***
		return -1, errors.New(ret.Err)
	***REMOVED***
	return ret.Size, nil
***REMOVED***

func (d *graphDriverProxy) DiffSize(id, parent string) (int64, error) ***REMOVED***
	args := &graphDriverRequest***REMOVED***
		ID:     id,
		Parent: parent,
	***REMOVED***
	var ret graphDriverResponse
	if err := d.p.Client().Call("GraphDriver.DiffSize", args, &ret); err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	if ret.Err != "" ***REMOVED***
		return -1, errors.New(ret.Err)
	***REMOVED***
	return ret.Size, nil
***REMOVED***
