// +build !windows

package shim

import (
	"context"
	"path/filepath"

	shimapi "github.com/containerd/containerd/linux/shim/v1"
	"github.com/containerd/containerd/mount"
	ptypes "github.com/gogo/protobuf/types"
)

// NewLocal returns a shim client implementation for issue commands to a shim
func NewLocal(s *Service) shimapi.ShimService ***REMOVED***
	return &local***REMOVED***
		s: s,
	***REMOVED***
***REMOVED***

type local struct ***REMOVED***
	s *Service
***REMOVED***

func (c *local) Create(ctx context.Context, in *shimapi.CreateTaskRequest) (*shimapi.CreateTaskResponse, error) ***REMOVED***
	return c.s.Create(ctx, in)
***REMOVED***

func (c *local) Start(ctx context.Context, in *shimapi.StartRequest) (*shimapi.StartResponse, error) ***REMOVED***
	return c.s.Start(ctx, in)
***REMOVED***

func (c *local) Delete(ctx context.Context, in *ptypes.Empty) (*shimapi.DeleteResponse, error) ***REMOVED***
	// make sure we unmount the containers rootfs for this local
	if err := mount.Unmount(filepath.Join(c.s.config.Path, "rootfs"), 0); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return c.s.Delete(ctx, in)
***REMOVED***

func (c *local) DeleteProcess(ctx context.Context, in *shimapi.DeleteProcessRequest) (*shimapi.DeleteResponse, error) ***REMOVED***
	return c.s.DeleteProcess(ctx, in)
***REMOVED***

func (c *local) Exec(ctx context.Context, in *shimapi.ExecProcessRequest) (*ptypes.Empty, error) ***REMOVED***
	return c.s.Exec(ctx, in)
***REMOVED***

func (c *local) ResizePty(ctx context.Context, in *shimapi.ResizePtyRequest) (*ptypes.Empty, error) ***REMOVED***
	return c.s.ResizePty(ctx, in)
***REMOVED***

func (c *local) State(ctx context.Context, in *shimapi.StateRequest) (*shimapi.StateResponse, error) ***REMOVED***
	return c.s.State(ctx, in)
***REMOVED***

func (c *local) Pause(ctx context.Context, in *ptypes.Empty) (*ptypes.Empty, error) ***REMOVED***
	return c.s.Pause(ctx, in)
***REMOVED***

func (c *local) Resume(ctx context.Context, in *ptypes.Empty) (*ptypes.Empty, error) ***REMOVED***
	return c.s.Resume(ctx, in)
***REMOVED***

func (c *local) Kill(ctx context.Context, in *shimapi.KillRequest) (*ptypes.Empty, error) ***REMOVED***
	return c.s.Kill(ctx, in)
***REMOVED***

func (c *local) ListPids(ctx context.Context, in *shimapi.ListPidsRequest) (*shimapi.ListPidsResponse, error) ***REMOVED***
	return c.s.ListPids(ctx, in)
***REMOVED***

func (c *local) CloseIO(ctx context.Context, in *shimapi.CloseIORequest) (*ptypes.Empty, error) ***REMOVED***
	return c.s.CloseIO(ctx, in)
***REMOVED***

func (c *local) Checkpoint(ctx context.Context, in *shimapi.CheckpointTaskRequest) (*ptypes.Empty, error) ***REMOVED***
	return c.s.Checkpoint(ctx, in)
***REMOVED***

func (c *local) ShimInfo(ctx context.Context, in *ptypes.Empty) (*shimapi.ShimInfoResponse, error) ***REMOVED***
	return c.s.ShimInfo(ctx, in)
***REMOVED***

func (c *local) Update(ctx context.Context, in *shimapi.UpdateTaskRequest) (*ptypes.Empty, error) ***REMOVED***
	return c.s.Update(ctx, in)
***REMOVED***

func (c *local) Wait(ctx context.Context, in *shimapi.WaitRequest) (*shimapi.WaitResponse, error) ***REMOVED***
	return c.s.Wait(ctx, in)
***REMOVED***
