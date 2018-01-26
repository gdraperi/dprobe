package filesync

import (
	"fmt"
	"os"
	"strings"

	"github.com/moby/buildkit/session"
	"github.com/pkg/errors"
	"github.com/tonistiigi/fsutil"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	keyOverrideExcludes = "override-excludes"
	keyIncludePatterns  = "include-patterns"
	keyDirName          = "dir-name"
)

type fsSyncProvider struct ***REMOVED***
	dirs   map[string]SyncedDir
	p      progressCb
	doneCh chan error
***REMOVED***

type SyncedDir struct ***REMOVED***
	Name     string
	Dir      string
	Excludes []string
	Map      func(*fsutil.Stat) bool
***REMOVED***

// NewFSSyncProvider creates a new provider for sending files from client
func NewFSSyncProvider(dirs []SyncedDir) session.Attachable ***REMOVED***
	p := &fsSyncProvider***REMOVED***
		dirs: map[string]SyncedDir***REMOVED******REMOVED***,
	***REMOVED***
	for _, d := range dirs ***REMOVED***
		p.dirs[d.Name] = d
	***REMOVED***
	return p
***REMOVED***

func (sp *fsSyncProvider) Register(server *grpc.Server) ***REMOVED***
	RegisterFileSyncServer(server, sp)
***REMOVED***

func (sp *fsSyncProvider) DiffCopy(stream FileSync_DiffCopyServer) error ***REMOVED***
	return sp.handle("diffcopy", stream)
***REMOVED***
func (sp *fsSyncProvider) TarStream(stream FileSync_TarStreamServer) error ***REMOVED***
	return sp.handle("tarstream", stream)
***REMOVED***

func (sp *fsSyncProvider) handle(method string, stream grpc.ServerStream) error ***REMOVED***
	var pr *protocol
	for _, p := range supportedProtocols ***REMOVED***
		if method == p.name && isProtoSupported(p.name) ***REMOVED***
			pr = &p
			break
		***REMOVED***
	***REMOVED***
	if pr == nil ***REMOVED***
		return errors.New("failed to negotiate protocol")
	***REMOVED***

	opts, _ := metadata.FromContext(stream.Context()) // if no metadata continue with empty object

	name, ok := opts[keyDirName]
	if !ok || len(name) != 1 ***REMOVED***
		return errors.New("no dir name in request")
	***REMOVED***

	dir, ok := sp.dirs[name[0]]
	if !ok ***REMOVED***
		return errors.Errorf("no access allowed to dir %q", name[0])
	***REMOVED***

	var excludes []string
	if len(opts[keyOverrideExcludes]) == 0 || opts[keyOverrideExcludes][0] != "true" ***REMOVED***
		excludes = dir.Excludes
	***REMOVED***
	includes := opts[keyIncludePatterns]

	var progress progressCb
	if sp.p != nil ***REMOVED***
		progress = sp.p
		sp.p = nil
	***REMOVED***

	var doneCh chan error
	if sp.doneCh != nil ***REMOVED***
		doneCh = sp.doneCh
		sp.doneCh = nil
	***REMOVED***
	err := pr.sendFn(stream, dir.Dir, includes, excludes, progress, dir.Map)
	if doneCh != nil ***REMOVED***
		if err != nil ***REMOVED***
			doneCh <- err
		***REMOVED***
		close(doneCh)
	***REMOVED***
	return err
***REMOVED***

func (sp *fsSyncProvider) SetNextProgressCallback(f func(int, bool), doneCh chan error) ***REMOVED***
	sp.p = f
	sp.doneCh = doneCh
***REMOVED***

type progressCb func(int, bool)

type protocol struct ***REMOVED***
	name   string
	sendFn func(stream grpc.Stream, srcDir string, includes, excludes []string, progress progressCb, _map func(*fsutil.Stat) bool) error
	recvFn func(stream grpc.Stream, destDir string, cu CacheUpdater, progress progressCb) error
***REMOVED***

func isProtoSupported(p string) bool ***REMOVED***
	// TODO: this should be removed after testing if stability is confirmed
	if override := os.Getenv("BUILD_STREAM_PROTOCOL"); override != "" ***REMOVED***
		return strings.EqualFold(p, override)
	***REMOVED***
	return true
***REMOVED***

var supportedProtocols = []protocol***REMOVED***
	***REMOVED***
		name:   "diffcopy",
		sendFn: sendDiffCopy,
		recvFn: recvDiffCopy,
	***REMOVED***,
***REMOVED***

// FSSendRequestOpt defines options for FSSend request
type FSSendRequestOpt struct ***REMOVED***
	Name             string
	IncludePatterns  []string
	OverrideExcludes bool
	DestDir          string
	CacheUpdater     CacheUpdater
	ProgressCb       func(int, bool)
***REMOVED***

// CacheUpdater is an object capable of sending notifications for the cache hash changes
type CacheUpdater interface ***REMOVED***
	MarkSupported(bool)
	HandleChange(fsutil.ChangeKind, string, os.FileInfo, error) error
	ContentHasher() fsutil.ContentHasher
***REMOVED***

// FSSync initializes a transfer of files
func FSSync(ctx context.Context, c session.Caller, opt FSSendRequestOpt) error ***REMOVED***
	var pr *protocol
	for _, p := range supportedProtocols ***REMOVED***
		if isProtoSupported(p.name) && c.Supports(session.MethodURL(_FileSync_serviceDesc.ServiceName, p.name)) ***REMOVED***
			pr = &p
			break
		***REMOVED***
	***REMOVED***
	if pr == nil ***REMOVED***
		return errors.New("no fssync handlers")
	***REMOVED***

	opts := make(map[string][]string)
	if opt.OverrideExcludes ***REMOVED***
		opts[keyOverrideExcludes] = []string***REMOVED***"true"***REMOVED***
	***REMOVED***

	if opt.IncludePatterns != nil ***REMOVED***
		opts[keyIncludePatterns] = opt.IncludePatterns
	***REMOVED***

	opts[keyDirName] = []string***REMOVED***opt.Name***REMOVED***

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := NewFileSyncClient(c.Conn())

	var stream grpc.ClientStream

	ctx = metadata.NewContext(ctx, opts)

	switch pr.name ***REMOVED***
	case "tarstream":
		cc, err := client.TarStream(ctx)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		stream = cc
	case "diffcopy":
		cc, err := client.DiffCopy(ctx)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		stream = cc
	default:
		panic(fmt.Sprintf("invalid protocol: %q", pr.name))
	***REMOVED***

	return pr.recvFn(stream, opt.DestDir, opt.CacheUpdater, opt.ProgressCb)
***REMOVED***

// NewFSSyncTarget allows writing into a directory
func NewFSSyncTarget(outdir string) session.Attachable ***REMOVED***
	p := &fsSyncTarget***REMOVED***
		outdir: outdir,
	***REMOVED***
	return p
***REMOVED***

type fsSyncTarget struct ***REMOVED***
	outdir string
***REMOVED***

func (sp *fsSyncTarget) Register(server *grpc.Server) ***REMOVED***
	RegisterFileSendServer(server, sp)
***REMOVED***

func (sp *fsSyncTarget) DiffCopy(stream FileSend_DiffCopyServer) error ***REMOVED***
	return syncTargetDiffCopy(stream, sp.outdir)
***REMOVED***

func CopyToCaller(ctx context.Context, srcPath string, c session.Caller, progress func(int, bool)) error ***REMOVED***
	method := session.MethodURL(_FileSend_serviceDesc.ServiceName, "diffcopy")
	if !c.Supports(method) ***REMOVED***
		return errors.Errorf("method %s not supported by the client", method)
	***REMOVED***

	client := NewFileSendClient(c.Conn())

	cc, err := client.DiffCopy(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return sendDiffCopy(cc, srcPath, nil, nil, progress, nil)
***REMOVED***
