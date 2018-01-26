package filesync

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tonistiigi/fsutil"
	"google.golang.org/grpc"
)

func sendDiffCopy(stream grpc.Stream, dir string, includes, excludes []string, progress progressCb, _map func(*fsutil.Stat) bool) error ***REMOVED***
	return fsutil.Send(stream.Context(), stream, dir, &fsutil.WalkOpt***REMOVED***
		ExcludePatterns: excludes,
		IncludePatterns: includes,
		Map:             _map,
	***REMOVED***, progress)
***REMOVED***

func recvDiffCopy(ds grpc.Stream, dest string, cu CacheUpdater, progress progressCb) error ***REMOVED***
	st := time.Now()
	defer func() ***REMOVED***
		logrus.Debugf("diffcopy took: %v", time.Since(st))
	***REMOVED***()
	var cf fsutil.ChangeFunc
	var ch fsutil.ContentHasher
	if cu != nil ***REMOVED***
		cu.MarkSupported(true)
		cf = cu.HandleChange
		ch = cu.ContentHasher()
	***REMOVED***
	return fsutil.Receive(ds.Context(), ds, dest, fsutil.ReceiveOpt***REMOVED***
		NotifyHashed:  cf,
		ContentHasher: ch,
		ProgressCb:    progress,
	***REMOVED***)
***REMOVED***

func syncTargetDiffCopy(ds grpc.Stream, dest string) error ***REMOVED***
	if err := os.MkdirAll(dest, 0700); err != nil ***REMOVED***
		return err
	***REMOVED***
	return fsutil.Receive(ds.Context(), ds, dest, fsutil.ReceiveOpt***REMOVED***
		Merge: true,
		Filter: func() func(*fsutil.Stat) bool ***REMOVED***
			uid := os.Getuid()
			gid := os.Getgid()
			return func(st *fsutil.Stat) bool ***REMOVED***
				st.Uid = uint32(uid)
				st.Gid = uint32(gid)
				return true
			***REMOVED***
		***REMOVED***(),
	***REMOVED***)
***REMOVED***
