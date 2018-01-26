package fsutil

import (
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type ReceiveOpt struct ***REMOVED***
	NotifyHashed  ChangeFunc
	ContentHasher ContentHasher
	ProgressCb    func(int, bool)
	Merge         bool
	Filter        FilterFunc
***REMOVED***

func Receive(ctx context.Context, conn Stream, dest string, opt ReceiveOpt) error ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := &receiver***REMOVED***
		conn:          &syncStream***REMOVED***Stream: conn***REMOVED***,
		dest:          dest,
		files:         make(map[string]uint32),
		pipes:         make(map[uint32]io.WriteCloser),
		notifyHashed:  opt.NotifyHashed,
		contentHasher: opt.ContentHasher,
		progressCb:    opt.ProgressCb,
		merge:         opt.Merge,
		filter:        opt.Filter,
	***REMOVED***
	return r.run(ctx)
***REMOVED***

type receiver struct ***REMOVED***
	dest       string
	conn       Stream
	files      map[string]uint32
	pipes      map[uint32]io.WriteCloser
	mu         sync.RWMutex
	muPipes    sync.RWMutex
	progressCb func(int, bool)
	merge      bool
	filter     FilterFunc

	notifyHashed   ChangeFunc
	contentHasher  ContentHasher
	orderValidator Validator
	hlValidator    Hardlinks
***REMOVED***

type dynamicWalker struct ***REMOVED***
	walkChan chan *currentPath
	closed   bool
***REMOVED***

func newDynamicWalker() *dynamicWalker ***REMOVED***
	return &dynamicWalker***REMOVED***
		walkChan: make(chan *currentPath, 128),
	***REMOVED***
***REMOVED***

func (w *dynamicWalker) update(p *currentPath) error ***REMOVED***
	if w.closed ***REMOVED***
		return errors.New("walker is closed")
	***REMOVED***
	if p == nil ***REMOVED***
		close(w.walkChan)
		return nil
	***REMOVED***
	w.walkChan <- p
	return nil
***REMOVED***

func (w *dynamicWalker) fill(ctx context.Context, pathC chan<- *currentPath) error ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case p, ok := <-w.walkChan:
			if !ok ***REMOVED***
				return nil
			***REMOVED***
			pathC <- p
		case <-ctx.Done():
			return ctx.Err()
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *receiver) run(ctx context.Context) error ***REMOVED***
	g, ctx := errgroup.WithContext(ctx)

	dw, err := NewDiskWriter(ctx, r.dest, DiskWriterOpt***REMOVED***
		AsyncDataCb:   r.asyncDataFunc,
		NotifyCb:      r.notifyHashed,
		ContentHasher: r.contentHasher,
		Filter:        r.filter,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	w := newDynamicWalker()

	g.Go(func() error ***REMOVED***
		destWalker := emptyWalker
		if !r.merge ***REMOVED***
			destWalker = GetWalkerFn(r.dest)
		***REMOVED***
		err := doubleWalkDiff(ctx, dw.HandleChange, destWalker, w.fill)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := dw.Wait(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
		r.conn.SendMsg(&Packet***REMOVED***Type: PACKET_FIN***REMOVED***)
		return nil
	***REMOVED***)

	g.Go(func() error ***REMOVED***
		var i uint32 = 0

		size := 0
		if r.progressCb != nil ***REMOVED***
			defer func() ***REMOVED***
				r.progressCb(size, true)
			***REMOVED***()
		***REMOVED***
		var p Packet
		for ***REMOVED***
			p = Packet***REMOVED***Data: p.Data[:0]***REMOVED***
			if err := r.conn.RecvMsg(&p); err != nil ***REMOVED***
				return err
			***REMOVED***
			if r.progressCb != nil ***REMOVED***
				size += p.Size()
				r.progressCb(size, false)
			***REMOVED***

			switch p.Type ***REMOVED***
			case PACKET_STAT:
				if p.Stat == nil ***REMOVED***
					if err := w.update(nil); err != nil ***REMOVED***
						return err
					***REMOVED***
					break
				***REMOVED***
				if fileCanRequestData(os.FileMode(p.Stat.Mode)) ***REMOVED***
					r.mu.Lock()
					r.files[p.Stat.Path] = i
					r.mu.Unlock()
				***REMOVED***
				i++
				cp := &currentPath***REMOVED***path: p.Stat.Path, f: &StatInfo***REMOVED***p.Stat***REMOVED******REMOVED***
				if err := r.orderValidator.HandleChange(ChangeKindAdd, cp.path, cp.f, nil); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := r.hlValidator.HandleChange(ChangeKindAdd, cp.path, cp.f, nil); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := w.update(cp); err != nil ***REMOVED***
					return err
				***REMOVED***
			case PACKET_DATA:
				r.muPipes.Lock()
				pw, ok := r.pipes[p.ID]
				r.muPipes.Unlock()
				if !ok ***REMOVED***
					return errors.Errorf("invalid file request %s", p.ID)
				***REMOVED***
				if len(p.Data) == 0 ***REMOVED***
					if err := pw.Close(); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if _, err := pw.Write(p.Data); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			case PACKET_FIN:
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	return g.Wait()
***REMOVED***

func (r *receiver) asyncDataFunc(ctx context.Context, p string, wc io.WriteCloser) error ***REMOVED***
	r.mu.Lock()
	id, ok := r.files[p]
	if !ok ***REMOVED***
		r.mu.Unlock()
		return errors.Errorf("invalid file request %s", p)
	***REMOVED***
	delete(r.files, p)
	r.mu.Unlock()

	wwc := newWrappedWriteCloser(wc)
	r.muPipes.Lock()
	r.pipes[id] = wwc
	r.muPipes.Unlock()
	if err := r.conn.SendMsg(&Packet***REMOVED***Type: PACKET_REQ, ID: id***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	err := wwc.Wait(ctx)
	r.muPipes.Lock()
	delete(r.pipes, id)
	r.muPipes.Unlock()
	return err
***REMOVED***

type wrappedWriteCloser struct ***REMOVED***
	io.WriteCloser
	err  error
	once sync.Once
	done chan struct***REMOVED******REMOVED***
***REMOVED***

func newWrappedWriteCloser(wc io.WriteCloser) *wrappedWriteCloser ***REMOVED***
	return &wrappedWriteCloser***REMOVED***WriteCloser: wc, done: make(chan struct***REMOVED******REMOVED***)***REMOVED***
***REMOVED***

func (w *wrappedWriteCloser) Close() error ***REMOVED***
	w.err = w.WriteCloser.Close()
	w.once.Do(func() ***REMOVED*** close(w.done) ***REMOVED***)
	return w.err
***REMOVED***

func (w *wrappedWriteCloser) Wait(ctx context.Context) error ***REMOVED***
	select ***REMOVED***
	case <-ctx.Done():
		return ctx.Err()
	case <-w.done:
		return w.err
	***REMOVED***
***REMOVED***
