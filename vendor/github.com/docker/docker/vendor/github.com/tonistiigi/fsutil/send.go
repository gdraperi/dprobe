package fsutil

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

var bufPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return make([]byte, 32*1<<10)
	***REMOVED***,
***REMOVED***

type Stream interface ***REMOVED***
	RecvMsg(interface***REMOVED******REMOVED***) error
	SendMsg(m interface***REMOVED******REMOVED***) error
	Context() context.Context
***REMOVED***

func Send(ctx context.Context, conn Stream, root string, opt *WalkOpt, progressCb func(int, bool)) error ***REMOVED***
	s := &sender***REMOVED***
		conn:         &syncStream***REMOVED***Stream: conn***REMOVED***,
		root:         root,
		opt:          opt,
		files:        make(map[uint32]string),
		progressCb:   progressCb,
		sendpipeline: make(chan *sendHandle, 128),
	***REMOVED***
	return s.run(ctx)
***REMOVED***

type sendHandle struct ***REMOVED***
	id   uint32
	path string
***REMOVED***

type sender struct ***REMOVED***
	conn            Stream
	opt             *WalkOpt
	root            string
	files           map[uint32]string
	mu              sync.RWMutex
	progressCb      func(int, bool)
	progressCurrent int
	sendpipeline    chan *sendHandle
***REMOVED***

func (s *sender) run(ctx context.Context) error ***REMOVED***
	g, ctx := errgroup.WithContext(ctx)

	defer s.updateProgress(0, true)

	g.Go(func() error ***REMOVED***
		return s.walk(ctx)
	***REMOVED***)

	for i := 0; i < 4; i++ ***REMOVED***
		g.Go(func() error ***REMOVED***
			for h := range s.sendpipeline ***REMOVED***
				select ***REMOVED***
				case <-ctx.Done():
					return ctx.Err()
				default:
				***REMOVED***
				if err := s.sendFile(h); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***)
	***REMOVED***

	g.Go(func() error ***REMOVED***
		defer close(s.sendpipeline)

		for ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				return ctx.Err()
			default:
			***REMOVED***
			var p Packet
			if err := s.conn.RecvMsg(&p); err != nil ***REMOVED***
				return err
			***REMOVED***
			switch p.Type ***REMOVED***
			case PACKET_REQ:
				if err := s.queue(p.ID); err != nil ***REMOVED***
					return err
				***REMOVED***
			case PACKET_FIN:
				return s.conn.SendMsg(&Packet***REMOVED***Type: PACKET_FIN***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	return g.Wait()
***REMOVED***

func (s *sender) updateProgress(size int, last bool) ***REMOVED***
	if s.progressCb != nil ***REMOVED***
		s.progressCurrent += size
		s.progressCb(s.progressCurrent, last)
	***REMOVED***
***REMOVED***

func (s *sender) queue(id uint32) error ***REMOVED***
	s.mu.Lock()
	p, ok := s.files[id]
	if !ok ***REMOVED***
		s.mu.Unlock()
		return errors.Errorf("invalid file id %d", id)
	***REMOVED***
	delete(s.files, id)
	s.mu.Unlock()
	s.sendpipeline <- &sendHandle***REMOVED***id, p***REMOVED***
	return nil
***REMOVED***

func (s *sender) sendFile(h *sendHandle) error ***REMOVED***
	f, err := os.Open(filepath.Join(s.root, h.path))
	if err == nil ***REMOVED***
		buf := bufPool.Get().([]byte)
		defer bufPool.Put(buf)
		if _, err := io.CopyBuffer(&fileSender***REMOVED***sender: s, id: h.id***REMOVED***, f, buf); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return s.conn.SendMsg(&Packet***REMOVED***ID: h.id, Type: PACKET_DATA***REMOVED***)
***REMOVED***

func (s *sender) walk(ctx context.Context) error ***REMOVED***
	var i uint32 = 0
	err := Walk(ctx, s.root, s.opt, func(path string, fi os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		stat, ok := fi.Sys().(*Stat)
		if !ok ***REMOVED***
			return errors.Wrapf(err, "invalid fileinfo without stat info: %s", path)
		***REMOVED***

		p := &Packet***REMOVED***
			Type: PACKET_STAT,
			Stat: stat,
		***REMOVED***
		if fileCanRequestData(os.FileMode(stat.Mode)) ***REMOVED***
			s.mu.Lock()
			s.files[i] = stat.Path
			s.mu.Unlock()
		***REMOVED***
		i++
		s.updateProgress(p.Size(), false)
		return errors.Wrapf(s.conn.SendMsg(p), "failed to send stat %s", path)
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return errors.Wrapf(s.conn.SendMsg(&Packet***REMOVED***Type: PACKET_STAT***REMOVED***), "failed to send last stat")
***REMOVED***

func fileCanRequestData(m os.FileMode) bool ***REMOVED***
	// avoid updating this function as it needs to match between sender/receiver.
	// version if needed
	return m&os.ModeType == 0
***REMOVED***

type fileSender struct ***REMOVED***
	sender *sender
	id     uint32
***REMOVED***

func (fs *fileSender) Write(dt []byte) (int, error) ***REMOVED***
	if len(dt) == 0 ***REMOVED***
		return 0, nil
	***REMOVED***
	p := &Packet***REMOVED***Type: PACKET_DATA, ID: fs.id, Data: dt***REMOVED***
	if err := fs.sender.conn.SendMsg(p); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	fs.sender.updateProgress(p.Size(), false)
	return len(dt), nil
***REMOVED***

type syncStream struct ***REMOVED***
	Stream
	mu sync.Mutex
***REMOVED***

func (ss *syncStream) SendMsg(m interface***REMOVED******REMOVED***) error ***REMOVED***
	ss.mu.Lock()
	err := ss.Stream.SendMsg(m)
	ss.mu.Unlock()
	return err
***REMOVED***
