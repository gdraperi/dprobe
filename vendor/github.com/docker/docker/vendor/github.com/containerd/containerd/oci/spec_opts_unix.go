// +build !windows

package oci

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/fs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/opencontainers/runc/libcontainer/user"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

// WithTTY sets the information on the spec as well as the environment variables for
// using a TTY
func WithTTY(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
	s.Process.Terminal = true
	s.Process.Env = append(s.Process.Env, "TERM=xterm")
	return nil
***REMOVED***

// WithHostNamespace allows a task to run inside the host's linux namespace
func WithHostNamespace(ns specs.LinuxNamespaceType) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		for i, n := range s.Linux.Namespaces ***REMOVED***
			if n.Type == ns ***REMOVED***
				s.Linux.Namespaces = append(s.Linux.Namespaces[:i], s.Linux.Namespaces[i+1:]...)
				return nil
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// WithLinuxNamespace uses the passed in namespace for the spec. If a namespace of the same type already exists in the
// spec, the existing namespace is replaced by the one provided.
func WithLinuxNamespace(ns specs.LinuxNamespace) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		for i, n := range s.Linux.Namespaces ***REMOVED***
			if n.Type == ns.Type ***REMOVED***
				before := s.Linux.Namespaces[:i]
				after := s.Linux.Namespaces[i+1:]
				s.Linux.Namespaces = append(before, ns)
				s.Linux.Namespaces = append(s.Linux.Namespaces, after...)
				return nil
			***REMOVED***
		***REMOVED***
		s.Linux.Namespaces = append(s.Linux.Namespaces, ns)
		return nil
	***REMOVED***
***REMOVED***

// WithImageConfig configures the spec to from the configuration of an Image
func WithImageConfig(image Image) SpecOpts ***REMOVED***
	return func(ctx context.Context, client Client, c *containers.Container, s *specs.Spec) error ***REMOVED***
		ic, err := image.Config(ctx)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		var (
			ociimage v1.Image
			config   v1.ImageConfig
		)
		switch ic.MediaType ***REMOVED***
		case v1.MediaTypeImageConfig, images.MediaTypeDockerSchema2Config:
			p, err := content.ReadBlob(ctx, image.ContentStore(), ic.Digest)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if err := json.Unmarshal(p, &ociimage); err != nil ***REMOVED***
				return err
			***REMOVED***
			config = ociimage.Config
		default:
			return fmt.Errorf("unknown image config media type %s", ic.MediaType)
		***REMOVED***

		if s.Process == nil ***REMOVED***
			s.Process = &specs.Process***REMOVED******REMOVED***
		***REMOVED***

		s.Process.Env = append(s.Process.Env, config.Env...)
		cmd := config.Cmd
		s.Process.Args = append(config.Entrypoint, cmd...)
		if config.User != "" ***REMOVED***
			parts := strings.Split(config.User, ":")
			switch len(parts) ***REMOVED***
			case 1:
				v, err := strconv.Atoi(parts[0])
				if err != nil ***REMOVED***
					// if we cannot parse as a uint they try to see if it is a username
					if err := WithUsername(config.User)(ctx, client, c, s); err != nil ***REMOVED***
						return err
					***REMOVED***
					return err
				***REMOVED***
				if err := WithUserID(uint32(v))(ctx, client, c, s); err != nil ***REMOVED***
					return err
				***REMOVED***
			case 2:
				v, err := strconv.Atoi(parts[0])
				if err != nil ***REMOVED***
					return errors.Wrapf(err, "parse uid %s", parts[0])
				***REMOVED***
				uid := uint32(v)
				if v, err = strconv.Atoi(parts[1]); err != nil ***REMOVED***
					return errors.Wrapf(err, "parse gid %s", parts[1])
				***REMOVED***
				gid := uint32(v)
				s.Process.User.UID, s.Process.User.GID = uid, gid
			default:
				return fmt.Errorf("invalid USER value %s", config.User)
			***REMOVED***
		***REMOVED***
		cwd := config.WorkingDir
		if cwd == "" ***REMOVED***
			cwd = "/"
		***REMOVED***
		s.Process.Cwd = cwd
		return nil
	***REMOVED***
***REMOVED***

// WithRootFSPath specifies unmanaged rootfs path.
func WithRootFSPath(path string) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		if s.Root == nil ***REMOVED***
			s.Root = &specs.Root***REMOVED******REMOVED***
		***REMOVED***
		s.Root.Path = path
		// Entrypoint is not set here (it's up to caller)
		return nil
	***REMOVED***
***REMOVED***

// WithRootFSReadonly sets specs.Root.Readonly to true
func WithRootFSReadonly() SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		if s.Root == nil ***REMOVED***
			s.Root = &specs.Root***REMOVED******REMOVED***
		***REMOVED***
		s.Root.Readonly = true
		return nil
	***REMOVED***
***REMOVED***

// WithNoNewPrivileges sets no_new_privileges on the process for the container
func WithNoNewPrivileges(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
	s.Process.NoNewPrivileges = true
	return nil
***REMOVED***

// WithHostHostsFile bind-mounts the host's /etc/hosts into the container as readonly
func WithHostHostsFile(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
	s.Mounts = append(s.Mounts, specs.Mount***REMOVED***
		Destination: "/etc/hosts",
		Type:        "bind",
		Source:      "/etc/hosts",
		Options:     []string***REMOVED***"rbind", "ro"***REMOVED***,
	***REMOVED***)
	return nil
***REMOVED***

// WithHostResolvconf bind-mounts the host's /etc/resolv.conf into the container as readonly
func WithHostResolvconf(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
	s.Mounts = append(s.Mounts, specs.Mount***REMOVED***
		Destination: "/etc/resolv.conf",
		Type:        "bind",
		Source:      "/etc/resolv.conf",
		Options:     []string***REMOVED***"rbind", "ro"***REMOVED***,
	***REMOVED***)
	return nil
***REMOVED***

// WithHostLocaltime bind-mounts the host's /etc/localtime into the container as readonly
func WithHostLocaltime(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
	s.Mounts = append(s.Mounts, specs.Mount***REMOVED***
		Destination: "/etc/localtime",
		Type:        "bind",
		Source:      "/etc/localtime",
		Options:     []string***REMOVED***"rbind", "ro"***REMOVED***,
	***REMOVED***)
	return nil
***REMOVED***

// WithUserNamespace sets the uid and gid mappings for the task
// this can be called multiple times to add more mappings to the generated spec
func WithUserNamespace(container, host, size uint32) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		var hasUserns bool
		for _, ns := range s.Linux.Namespaces ***REMOVED***
			if ns.Type == specs.UserNamespace ***REMOVED***
				hasUserns = true
				break
			***REMOVED***
		***REMOVED***
		if !hasUserns ***REMOVED***
			s.Linux.Namespaces = append(s.Linux.Namespaces, specs.LinuxNamespace***REMOVED***
				Type: specs.UserNamespace,
			***REMOVED***)
		***REMOVED***
		mapping := specs.LinuxIDMapping***REMOVED***
			ContainerID: container,
			HostID:      host,
			Size:        size,
		***REMOVED***
		s.Linux.UIDMappings = append(s.Linux.UIDMappings, mapping)
		s.Linux.GIDMappings = append(s.Linux.GIDMappings, mapping)
		return nil
	***REMOVED***
***REMOVED***

// WithCgroup sets the container's cgroup path
func WithCgroup(path string) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		s.Linux.CgroupsPath = path
		return nil
	***REMOVED***
***REMOVED***

// WithNamespacedCgroup uses the namespace set on the context to create a
// root directory for containers in the cgroup with the id as the subcgroup
func WithNamespacedCgroup() SpecOpts ***REMOVED***
	return func(ctx context.Context, _ Client, c *containers.Container, s *specs.Spec) error ***REMOVED***
		namespace, err := namespaces.NamespaceRequired(ctx)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		s.Linux.CgroupsPath = filepath.Join("/", namespace, c.ID)
		return nil
	***REMOVED***
***REMOVED***

// WithUIDGID allows the UID and GID for the Process to be set
func WithUIDGID(uid, gid uint32) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		s.Process.User.UID = uid
		s.Process.User.GID = gid
		return nil
	***REMOVED***
***REMOVED***

// WithUserID sets the correct UID and GID for the container based
// on the image's /etc/passwd contents. If /etc/passwd does not exist,
// or uid is not found in /etc/passwd, it sets gid to be the same with
// uid, and not returns error.
func WithUserID(uid uint32) SpecOpts ***REMOVED***
	return func(ctx context.Context, client Client, c *containers.Container, s *specs.Spec) (err error) ***REMOVED***
		if c.Snapshotter == "" ***REMOVED***
			return errors.Errorf("no snapshotter set for container")
		***REMOVED***
		if c.SnapshotKey == "" ***REMOVED***
			return errors.Errorf("rootfs snapshot not created for container")
		***REMOVED***
		snapshotter := client.SnapshotService(c.Snapshotter)
		mounts, err := snapshotter.Mounts(ctx, c.SnapshotKey)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		root, err := ioutil.TempDir("", "ctd-username")
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer os.Remove(root)
		for _, m := range mounts ***REMOVED***
			if err := m.Mount(root); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		defer func() ***REMOVED***
			if uerr := mount.Unmount(root, 0); uerr != nil ***REMOVED***
				if err == nil ***REMOVED***
					err = uerr
				***REMOVED***
			***REMOVED***
		***REMOVED***()
		ppath, err := fs.RootPath(root, "/etc/passwd")
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		f, err := os.Open(ppath)
		if err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				s.Process.User.UID, s.Process.User.GID = uid, uid
				return nil
			***REMOVED***
			return err
		***REMOVED***
		defer f.Close()
		users, err := user.ParsePasswdFilter(f, func(u user.User) bool ***REMOVED***
			return u.Uid == int(uid)
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if len(users) == 0 ***REMOVED***
			s.Process.User.UID, s.Process.User.GID = uid, uid
			return nil
		***REMOVED***
		u := users[0]
		s.Process.User.UID, s.Process.User.GID = uint32(u.Uid), uint32(u.Gid)
		return nil
	***REMOVED***
***REMOVED***

// WithUsername sets the correct UID and GID for the container
// based on the the image's /etc/passwd contents. If /etc/passwd
// does not exist, or the username is not found in /etc/passwd,
// it returns error.
func WithUsername(username string) SpecOpts ***REMOVED***
	return func(ctx context.Context, client Client, c *containers.Container, s *specs.Spec) (err error) ***REMOVED***
		if c.Snapshotter == "" ***REMOVED***
			return errors.Errorf("no snapshotter set for container")
		***REMOVED***
		if c.SnapshotKey == "" ***REMOVED***
			return errors.Errorf("rootfs snapshot not created for container")
		***REMOVED***
		snapshotter := client.SnapshotService(c.Snapshotter)
		mounts, err := snapshotter.Mounts(ctx, c.SnapshotKey)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		root, err := ioutil.TempDir("", "ctd-username")
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer os.Remove(root)
		for _, m := range mounts ***REMOVED***
			if err := m.Mount(root); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		defer func() ***REMOVED***
			if uerr := mount.Unmount(root, 0); uerr != nil ***REMOVED***
				if err == nil ***REMOVED***
					err = uerr
				***REMOVED***
			***REMOVED***
		***REMOVED***()
		ppath, err := fs.RootPath(root, "/etc/passwd")
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		f, err := os.Open(ppath)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer f.Close()
		users, err := user.ParsePasswdFilter(f, func(u user.User) bool ***REMOVED***
			return u.Name == username
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if len(users) == 0 ***REMOVED***
			return errors.Errorf("no users found for %s", username)
		***REMOVED***
		u := users[0]
		s.Process.User.UID, s.Process.User.GID = uint32(u.Uid), uint32(u.Gid)
		return nil
	***REMOVED***
***REMOVED***
