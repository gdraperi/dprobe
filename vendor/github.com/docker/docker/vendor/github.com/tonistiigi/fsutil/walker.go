package fsutil

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type WalkOpt struct ***REMOVED***
	IncludePatterns []string
	ExcludePatterns []string
	Map             func(*Stat) bool
***REMOVED***

func Walk(ctx context.Context, p string, opt *WalkOpt, fn filepath.WalkFunc) error ***REMOVED***
	root, err := filepath.EvalSymlinks(p)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to resolve %s", root)
	***REMOVED***
	fi, err := os.Stat(root)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to stat: %s", root)
	***REMOVED***
	if !fi.IsDir() ***REMOVED***
		return errors.Errorf("%s is not a directory", root)
	***REMOVED***

	var pm *fileutils.PatternMatcher
	if opt != nil && opt.ExcludePatterns != nil ***REMOVED***
		pm, err = fileutils.NewPatternMatcher(opt.ExcludePatterns)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "invalid excludepaths %s", opt.ExcludePatterns)
		***REMOVED***
	***REMOVED***

	seenFiles := make(map[uint64]string)
	return filepath.Walk(root, func(path string, fi os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				return filepath.SkipDir
			***REMOVED***
			return err
		***REMOVED***
		origpath := path
		path, err = filepath.Rel(root, path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// Skip root
		if path == "." ***REMOVED***
			return nil
		***REMOVED***

		if opt != nil ***REMOVED***
			if opt.IncludePatterns != nil ***REMOVED***
				matched := false
				for _, p := range opt.IncludePatterns ***REMOVED***
					if m, _ := filepath.Match(p, path); m ***REMOVED***
						matched = true
						break
					***REMOVED***
				***REMOVED***
				if !matched ***REMOVED***
					if fi.IsDir() ***REMOVED***
						return filepath.SkipDir
					***REMOVED***
					return nil
				***REMOVED***
			***REMOVED***
			if pm != nil ***REMOVED***
				m, err := pm.Matches(path)
				if err != nil ***REMOVED***
					return errors.Wrap(err, "failed to match excludepatterns")
				***REMOVED***

				if m ***REMOVED***
					if fi.IsDir() ***REMOVED***
						if !pm.Exclusions() ***REMOVED***
							return filepath.SkipDir
						***REMOVED***
						dirSlash := path + string(filepath.Separator)
						for _, pat := range pm.Patterns() ***REMOVED***
							if !pat.Exclusion() ***REMOVED***
								continue
							***REMOVED***
							patStr := pat.String() + string(filepath.Separator)
							if strings.HasPrefix(patStr, dirSlash) ***REMOVED***
								goto passedFilter
							***REMOVED***
						***REMOVED***
						return filepath.SkipDir
					***REMOVED***
					return nil
				***REMOVED***
			***REMOVED***
		***REMOVED***

	passedFilter:
		path = filepath.ToSlash(path)

		stat := &Stat***REMOVED***
			Path:    path,
			Mode:    uint32(fi.Mode()),
			Size_:   fi.Size(),
			ModTime: fi.ModTime().UnixNano(),
		***REMOVED***

		setUnixOpt(fi, stat, path, seenFiles)

		if !fi.IsDir() ***REMOVED***
			if fi.Mode()&os.ModeSymlink != 0 ***REMOVED***
				link, err := os.Readlink(origpath)
				if err != nil ***REMOVED***
					return errors.Wrapf(err, "failed to readlink %s", origpath)
				***REMOVED***
				stat.Linkname = link
			***REMOVED***
		***REMOVED***
		if err := loadXattr(origpath, stat); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to xattr %s", path)
		***REMOVED***

		if runtime.GOOS == "windows" ***REMOVED***
			permPart := stat.Mode & uint32(os.ModePerm)
			noPermPart := stat.Mode &^ uint32(os.ModePerm)
			// Add the x bit: make everything +x from windows
			permPart |= 0111
			permPart &= 0755
			stat.Mode = noPermPart | permPart
		***REMOVED***

		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		default:
			if opt != nil && opt.Map != nil ***REMOVED***
				if allowed := opt.Map(stat); !allowed ***REMOVED***
					return nil
				***REMOVED***
			***REMOVED***
			if err := fn(stat.Path, &StatInfo***REMOVED***stat***REMOVED***, nil); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
***REMOVED***

type StatInfo struct ***REMOVED***
	*Stat
***REMOVED***

func (s *StatInfo) Name() string ***REMOVED***
	return filepath.Base(s.Stat.Path)
***REMOVED***
func (s *StatInfo) Size() int64 ***REMOVED***
	return s.Stat.Size_
***REMOVED***
func (s *StatInfo) Mode() os.FileMode ***REMOVED***
	return os.FileMode(s.Stat.Mode)
***REMOVED***
func (s *StatInfo) ModTime() time.Time ***REMOVED***
	return time.Unix(s.Stat.ModTime/1e9, s.Stat.ModTime%1e9)
***REMOVED***
func (s *StatInfo) IsDir() bool ***REMOVED***
	return s.Mode().IsDir()
***REMOVED***
func (s *StatInfo) Sys() interface***REMOVED******REMOVED*** ***REMOVED***
	return s.Stat
***REMOVED***
