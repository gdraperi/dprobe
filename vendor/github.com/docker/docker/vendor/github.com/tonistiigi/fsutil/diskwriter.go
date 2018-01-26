package fsutil

import (
	"hash"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type WriteToFunc func(context.Context, string, io.WriteCloser) error

type DiskWriterOpt struct ***REMOVED***
	AsyncDataCb   WriteToFunc
	SyncDataCb    WriteToFunc
	NotifyCb      func(ChangeKind, string, os.FileInfo, error) error
	ContentHasher ContentHasher
	Filter        FilterFunc
***REMOVED***

type FilterFunc func(*Stat) bool

type DiskWriter struct ***REMOVED***
	opt  DiskWriterOpt
	dest string

	wg     sync.WaitGroup
	ctx    context.Context
	cancel func()
	eg     *errgroup.Group
	filter FilterFunc
***REMOVED***

func NewDiskWriter(ctx context.Context, dest string, opt DiskWriterOpt) (*DiskWriter, error) ***REMOVED***
	if opt.SyncDataCb == nil && opt.AsyncDataCb == nil ***REMOVED***
		return nil, errors.New("no data callback specified")
	***REMOVED***
	if opt.SyncDataCb != nil && opt.AsyncDataCb != nil ***REMOVED***
		return nil, errors.New("can't specify both sync and async data callbacks")
	***REMOVED***

	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)

	return &DiskWriter***REMOVED***
		opt:    opt,
		dest:   dest,
		eg:     eg,
		ctx:    ctx,
		cancel: cancel,
	***REMOVED***, nil
***REMOVED***

func (dw *DiskWriter) Wait(ctx context.Context) error ***REMOVED***
	return dw.eg.Wait()
***REMOVED***

func (dw *DiskWriter) HandleChange(kind ChangeKind, p string, fi os.FileInfo, err error) (retErr error) ***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	select ***REMOVED***
	case <-dw.ctx.Done():
		return dw.ctx.Err()
	default:
	***REMOVED***

	defer func() ***REMOVED***
		if retErr != nil ***REMOVED***
			dw.cancel()
		***REMOVED***
	***REMOVED***()

	p = filepath.FromSlash(p)

	destPath := filepath.Join(dw.dest, p)

	if kind == ChangeKindDelete ***REMOVED***
		// todo: no need to validate if diff is trusted but is it always?
		if err := os.RemoveAll(destPath); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to remove: %s", destPath)
		***REMOVED***
		if dw.opt.NotifyCb != nil ***REMOVED***
			if err := dw.opt.NotifyCb(kind, p, nil, nil); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	stat, ok := fi.Sys().(*Stat)
	if !ok ***REMOVED***
		return errors.Errorf("%s invalid change without stat information", p)
	***REMOVED***

	if dw.filter != nil ***REMOVED***
		if ok := dw.filter(stat); !ok ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	rename := true
	oldFi, err := os.Lstat(destPath)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			if kind != ChangeKindAdd ***REMOVED***
				return errors.Wrapf(err, "invalid addition: %s", destPath)
			***REMOVED***
			rename = false
		***REMOVED*** else ***REMOVED***
			return errors.Wrapf(err, "failed to stat %s", destPath)
		***REMOVED***
	***REMOVED***

	if oldFi != nil && fi.IsDir() && oldFi.IsDir() ***REMOVED***
		if err := rewriteMetadata(destPath, stat); err != nil ***REMOVED***
			return errors.Wrapf(err, "error setting dir metadata for %s", destPath)
		***REMOVED***
		return nil
	***REMOVED***

	newPath := destPath
	if rename ***REMOVED***
		newPath = filepath.Join(filepath.Dir(destPath), ".tmp."+nextSuffix())
	***REMOVED***

	isRegularFile := false

	switch ***REMOVED***
	case fi.IsDir():
		if err := os.Mkdir(newPath, fi.Mode()); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to create dir %s", newPath)
		***REMOVED***
	case fi.Mode()&os.ModeDevice != 0 || fi.Mode()&os.ModeNamedPipe != 0:
		if err := handleTarTypeBlockCharFifo(newPath, stat); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to create device %s", newPath)
		***REMOVED***
	case fi.Mode()&os.ModeSymlink != 0:
		if err := os.Symlink(stat.Linkname, newPath); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to symlink %s", newPath)
		***REMOVED***
	case stat.Linkname != "":
		if err := os.Link(filepath.Join(dw.dest, stat.Linkname), newPath); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to link %s to %s", newPath, stat.Linkname)
		***REMOVED***
	default:
		isRegularFile = true
		file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY, fi.Mode()) //todo: windows
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to create %s", newPath)
		***REMOVED***
		if dw.opt.SyncDataCb != nil ***REMOVED***
			if err := dw.processChange(ChangeKindAdd, p, fi, file); err != nil ***REMOVED***
				file.Close()
				return err
			***REMOVED***
			break
		***REMOVED***
		if err := file.Close(); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to close %s", newPath)
		***REMOVED***
	***REMOVED***

	if err := rewriteMetadata(newPath, stat); err != nil ***REMOVED***
		return errors.Wrapf(err, "error setting metadata for %s", newPath)
	***REMOVED***

	if rename ***REMOVED***
		if err := os.Rename(newPath, destPath); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to rename %s to %s", newPath, destPath)
		***REMOVED***
	***REMOVED***

	if isRegularFile ***REMOVED***
		if dw.opt.AsyncDataCb != nil ***REMOVED***
			dw.requestAsyncFileData(p, destPath, fi)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return dw.processChange(kind, p, fi, nil)
	***REMOVED***

	return nil
***REMOVED***

func (dw *DiskWriter) requestAsyncFileData(p, dest string, fi os.FileInfo) ***REMOVED***
	// todo: limit worker threads
	dw.eg.Go(func() error ***REMOVED***
		if err := dw.processChange(ChangeKindAdd, p, fi, &lazyFileWriter***REMOVED***
			dest: dest,
		***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***
		return chtimes(dest, fi.ModTime().UnixNano()) // TODO: parent dirs
	***REMOVED***)
***REMOVED***

func (dw *DiskWriter) processChange(kind ChangeKind, p string, fi os.FileInfo, w io.WriteCloser) error ***REMOVED***
	origw := w
	var hw *hashedWriter
	if dw.opt.NotifyCb != nil ***REMOVED***
		var err error
		if hw, err = newHashWriter(dw.opt.ContentHasher, fi, w); err != nil ***REMOVED***
			return err
		***REMOVED***
		w = hw
	***REMOVED***
	if origw != nil ***REMOVED***
		fn := dw.opt.SyncDataCb
		if fn == nil && dw.opt.AsyncDataCb != nil ***REMOVED***
			fn = dw.opt.AsyncDataCb
		***REMOVED***
		if err := fn(dw.ctx, p, w); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if hw != nil ***REMOVED***
			hw.Close()
		***REMOVED***
	***REMOVED***
	if hw != nil ***REMOVED***
		return dw.opt.NotifyCb(kind, p, hw, nil)
	***REMOVED***
	return nil
***REMOVED***

type hashedWriter struct ***REMOVED***
	os.FileInfo
	io.Writer
	h    hash.Hash
	w    io.WriteCloser
	dgst digest.Digest
***REMOVED***

func newHashWriter(ch ContentHasher, fi os.FileInfo, w io.WriteCloser) (*hashedWriter, error) ***REMOVED***
	stat, ok := fi.Sys().(*Stat)
	if !ok ***REMOVED***
		return nil, errors.Errorf("invalid change without stat information")
	***REMOVED***

	h, err := ch(stat)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	hw := &hashedWriter***REMOVED***
		FileInfo: fi,
		Writer:   io.MultiWriter(w, h),
		h:        h,
		w:        w,
	***REMOVED***
	return hw, nil
***REMOVED***

func (hw *hashedWriter) Close() error ***REMOVED***
	hw.dgst = digest.NewDigest(digest.SHA256, hw.h)
	if hw.w != nil ***REMOVED***
		return hw.w.Close()
	***REMOVED***
	return nil
***REMOVED***

func (hw *hashedWriter) Digest() digest.Digest ***REMOVED***
	return hw.dgst
***REMOVED***

type lazyFileWriter struct ***REMOVED***
	dest string
	ctx  context.Context
	f    *os.File
***REMOVED***

func (lfw *lazyFileWriter) Write(dt []byte) (int, error) ***REMOVED***
	if lfw.f == nil ***REMOVED***
		file, err := os.OpenFile(lfw.dest, os.O_WRONLY, 0) //todo: windows
		if err != nil ***REMOVED***
			return 0, errors.Wrapf(err, "failed to open %s", lfw.dest)
		***REMOVED***
		lfw.f = file
	***REMOVED***
	return lfw.f.Write(dt)
***REMOVED***

func (lfw *lazyFileWriter) Close() error ***REMOVED***
	if lfw.f != nil ***REMOVED***
		return lfw.f.Close()
	***REMOVED***
	return nil
***REMOVED***

func mkdev(major int64, minor int64) uint32 ***REMOVED***
	return uint32(((minor & 0xfff00) << 12) | ((major & 0xfff) << 8) | (minor & 0xff))
***REMOVED***

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var rand uint32
var randmu sync.Mutex

func reseed() uint32 ***REMOVED***
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
***REMOVED***

func nextSuffix() string ***REMOVED***
	randmu.Lock()
	r := rand
	if r == 0 ***REMOVED***
		r = reseed()
	***REMOVED***
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
***REMOVED***
