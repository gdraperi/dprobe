package dockerfile

import (
	"archive/tar"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/remotecontext"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/pkg/errors"
)

const unnamedFilename = "__unnamed__"

type pathCache interface ***REMOVED***
	Load(key interface***REMOVED******REMOVED***) (value interface***REMOVED******REMOVED***, ok bool)
	Store(key, value interface***REMOVED******REMOVED***)
***REMOVED***

// copyInfo is a data object which stores the metadata about each source file in
// a copyInstruction
type copyInfo struct ***REMOVED***
	root         containerfs.ContainerFS
	path         string
	hash         string
	noDecompress bool
***REMOVED***

func (c copyInfo) fullPath() (string, error) ***REMOVED***
	return c.root.ResolveScopedPath(c.path, true)
***REMOVED***

func newCopyInfoFromSource(source builder.Source, path string, hash string) copyInfo ***REMOVED***
	return copyInfo***REMOVED***root: source.Root(), path: path, hash: hash***REMOVED***
***REMOVED***

func newCopyInfos(copyInfos ...copyInfo) []copyInfo ***REMOVED***
	return copyInfos
***REMOVED***

// copyInstruction is a fully parsed COPY or ADD command that is passed to
// Builder.performCopy to copy files into the image filesystem
type copyInstruction struct ***REMOVED***
	cmdName                 string
	infos                   []copyInfo
	dest                    string
	chownStr                string
	allowLocalDecompression bool
***REMOVED***

// copier reads a raw COPY or ADD command, fetches remote sources using a downloader,
// and creates a copyInstruction
type copier struct ***REMOVED***
	imageSource *imageMount
	source      builder.Source
	pathCache   pathCache
	download    sourceDownloader
	tmpPaths    []string
	platform    string
***REMOVED***

func copierFromDispatchRequest(req dispatchRequest, download sourceDownloader, imageSource *imageMount) copier ***REMOVED***
	return copier***REMOVED***
		source:      req.source,
		pathCache:   req.builder.pathCache,
		download:    download,
		imageSource: imageSource,
		platform:    req.builder.options.Platform,
	***REMOVED***
***REMOVED***

func (o *copier) createCopyInstruction(args []string, cmdName string) (copyInstruction, error) ***REMOVED***
	inst := copyInstruction***REMOVED***cmdName: cmdName***REMOVED***
	last := len(args) - 1

	// Work in platform-specific filepath semantics
	inst.dest = fromSlash(args[last], o.platform)
	separator := string(separator(o.platform))
	infos, err := o.getCopyInfosForSourcePaths(args[0:last], inst.dest)
	if err != nil ***REMOVED***
		return inst, errors.Wrapf(err, "%s failed", cmdName)
	***REMOVED***
	if len(infos) > 1 && !strings.HasSuffix(inst.dest, separator) ***REMOVED***
		return inst, errors.Errorf("When using %s with more than one source file, the destination must be a directory and end with a /", cmdName)
	***REMOVED***
	inst.infos = infos
	return inst, nil
***REMOVED***

// getCopyInfosForSourcePaths iterates over the source files and calculate the info
// needed to copy (e.g. hash value if cached)
// The dest is used in case source is URL (and ends with "/")
func (o *copier) getCopyInfosForSourcePaths(sources []string, dest string) ([]copyInfo, error) ***REMOVED***
	var infos []copyInfo
	for _, orig := range sources ***REMOVED***
		subinfos, err := o.getCopyInfoForSourcePath(orig, dest)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		infos = append(infos, subinfos...)
	***REMOVED***

	if len(infos) == 0 ***REMOVED***
		return nil, errors.New("no source files were specified")
	***REMOVED***
	return infos, nil
***REMOVED***

func (o *copier) getCopyInfoForSourcePath(orig, dest string) ([]copyInfo, error) ***REMOVED***
	if !urlutil.IsURL(orig) ***REMOVED***
		return o.calcCopyInfo(orig, true)
	***REMOVED***

	remote, path, err := o.download(orig)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// If path == "" then we are unable to determine filename from src
	// We have to make sure dest is available
	if path == "" ***REMOVED***
		if strings.HasSuffix(dest, "/") ***REMOVED***
			return nil, errors.Errorf("cannot determine filename for source %s", orig)
		***REMOVED***
		path = unnamedFilename
	***REMOVED***
	o.tmpPaths = append(o.tmpPaths, remote.Root().Path())

	hash, err := remote.Hash(path)
	ci := newCopyInfoFromSource(remote, path, hash)
	ci.noDecompress = true // data from http shouldn't be extracted even on ADD
	return newCopyInfos(ci), err
***REMOVED***

// Cleanup removes any temporary directories created as part of downloading
// remote files.
func (o *copier) Cleanup() ***REMOVED***
	for _, path := range o.tmpPaths ***REMOVED***
		os.RemoveAll(path)
	***REMOVED***
	o.tmpPaths = []string***REMOVED******REMOVED***
***REMOVED***

// TODO: allowWildcards can probably be removed by refactoring this function further.
func (o *copier) calcCopyInfo(origPath string, allowWildcards bool) ([]copyInfo, error) ***REMOVED***
	imageSource := o.imageSource

	// TODO: do this when creating copier. Requires validateCopySourcePath
	// (and other below) to be aware of the difference sources. Why is it only
	// done on image Source?
	if imageSource != nil ***REMOVED***
		var err error
		o.source, err = imageSource.Source()
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to copy from %s", imageSource.ImageID())
		***REMOVED***
	***REMOVED***

	if o.source == nil ***REMOVED***
		return nil, errors.Errorf("missing build context")
	***REMOVED***

	root := o.source.Root()

	if err := validateCopySourcePath(imageSource, origPath, root.OS()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Work in source OS specific filepath semantics
	// For LCOW, this is NOT the daemon OS.
	origPath = root.FromSlash(origPath)
	origPath = strings.TrimPrefix(origPath, string(root.Separator()))
	origPath = strings.TrimPrefix(origPath, "."+string(root.Separator()))

	// Deal with wildcards
	if allowWildcards && containsWildcards(origPath, root.OS()) ***REMOVED***
		return o.copyWithWildcards(origPath)
	***REMOVED***

	if imageSource != nil && imageSource.ImageID() != "" ***REMOVED***
		// return a cached copy if one exists
		if h, ok := o.pathCache.Load(imageSource.ImageID() + origPath); ok ***REMOVED***
			return newCopyInfos(newCopyInfoFromSource(o.source, origPath, h.(string))), nil
		***REMOVED***
	***REMOVED***

	// Deal with the single file case
	copyInfo, err := copyInfoForFile(o.source, origPath)
	switch ***REMOVED***
	case err != nil:
		return nil, err
	case copyInfo.hash != "":
		o.storeInPathCache(imageSource, origPath, copyInfo.hash)
		return newCopyInfos(copyInfo), err
	***REMOVED***

	// TODO: remove, handle dirs in Hash()
	subfiles, err := walkSource(o.source, origPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	hash := hashStringSlice("dir", subfiles)
	o.storeInPathCache(imageSource, origPath, hash)
	return newCopyInfos(newCopyInfoFromSource(o.source, origPath, hash)), nil
***REMOVED***

func containsWildcards(name, platform string) bool ***REMOVED***
	isWindows := platform == "windows"
	for i := 0; i < len(name); i++ ***REMOVED***
		ch := name[i]
		if ch == '\\' && !isWindows ***REMOVED***
			i++
		***REMOVED*** else if ch == '*' || ch == '?' || ch == '[' ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (o *copier) storeInPathCache(im *imageMount, path string, hash string) ***REMOVED***
	if im != nil ***REMOVED***
		o.pathCache.Store(im.ImageID()+path, hash)
	***REMOVED***
***REMOVED***

func (o *copier) copyWithWildcards(origPath string) ([]copyInfo, error) ***REMOVED***
	root := o.source.Root()
	var copyInfos []copyInfo
	if err := root.Walk(root.Path(), func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		rel, err := remotecontext.Rel(root, path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if rel == "." ***REMOVED***
			return nil
		***REMOVED***
		if match, _ := root.Match(origPath, rel); !match ***REMOVED***
			return nil
		***REMOVED***

		// Note we set allowWildcards to false in case the name has
		// a * in it
		subInfos, err := o.calcCopyInfo(rel, false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		copyInfos = append(copyInfos, subInfos...)
		return nil
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return copyInfos, nil
***REMOVED***

func copyInfoForFile(source builder.Source, path string) (copyInfo, error) ***REMOVED***
	fi, err := remotecontext.StatAt(source, path)
	if err != nil ***REMOVED***
		return copyInfo***REMOVED******REMOVED***, err
	***REMOVED***

	if fi.IsDir() ***REMOVED***
		return copyInfo***REMOVED******REMOVED***, nil
	***REMOVED***
	hash, err := source.Hash(path)
	if err != nil ***REMOVED***
		return copyInfo***REMOVED******REMOVED***, err
	***REMOVED***
	return newCopyInfoFromSource(source, path, "file:"+hash), nil
***REMOVED***

// TODO: dedupe with copyWithWildcards()
func walkSource(source builder.Source, origPath string) ([]string, error) ***REMOVED***
	fp, err := remotecontext.FullPath(source, origPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Must be a dir
	var subfiles []string
	err = source.Root().Walk(fp, func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		rel, err := remotecontext.Rel(source.Root(), path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if rel == "." ***REMOVED***
			return nil
		***REMOVED***
		hash, err := source.Hash(rel)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		// we already checked handleHash above
		subfiles = append(subfiles, hash)
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sort.Strings(subfiles)
	return subfiles, nil
***REMOVED***

type sourceDownloader func(string) (builder.Source, string, error)

func newRemoteSourceDownloader(output, stdout io.Writer) sourceDownloader ***REMOVED***
	return func(url string) (builder.Source, string, error) ***REMOVED***
		return downloadSource(output, stdout, url)
	***REMOVED***
***REMOVED***

func errOnSourceDownload(_ string) (builder.Source, string, error) ***REMOVED***
	return nil, "", errors.New("source can't be a URL for COPY")
***REMOVED***

func getFilenameForDownload(path string, resp *http.Response) string ***REMOVED***
	// Guess filename based on source
	if path != "" && !strings.HasSuffix(path, "/") ***REMOVED***
		if filename := filepath.Base(filepath.FromSlash(path)); filename != "" ***REMOVED***
			return filename
		***REMOVED***
	***REMOVED***

	// Guess filename based on Content-Disposition
	if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" ***REMOVED***
		if _, params, err := mime.ParseMediaType(contentDisposition); err == nil ***REMOVED***
			if params["filename"] != "" && !strings.HasSuffix(params["filename"], "/") ***REMOVED***
				if filename := filepath.Base(filepath.FromSlash(params["filename"])); filename != "" ***REMOVED***
					return filename
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func downloadSource(output io.Writer, stdout io.Writer, srcURL string) (remote builder.Source, p string, err error) ***REMOVED***
	u, err := url.Parse(srcURL)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	resp, err := remotecontext.GetWithStatusError(srcURL)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	filename := getFilenameForDownload(u.Path, resp)

	// Prepare file in a tmp dir
	tmpDir, err := ioutils.TempDir("", "docker-remote")
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			os.RemoveAll(tmpDir)
		***REMOVED***
	***REMOVED***()
	// If filename is empty, the returned filename will be "" but
	// the tmp filename will be created as "__unnamed__"
	tmpFileName := filename
	if filename == "" ***REMOVED***
		tmpFileName = unnamedFilename
	***REMOVED***
	tmpFileName = filepath.Join(tmpDir, tmpFileName)
	tmpFile, err := os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	progressOutput := streamformatter.NewJSONProgressOutput(output, true)
	progressReader := progress.NewProgressReader(resp.Body, progressOutput, resp.ContentLength, "", "Downloading")
	// Download and dump result to tmp file
	// TODO: add filehash directly
	if _, err = io.Copy(tmpFile, progressReader); err != nil ***REMOVED***
		tmpFile.Close()
		return
	***REMOVED***
	// TODO: how important is this random blank line to the output?
	fmt.Fprintln(stdout)

	// Set the mtime to the Last-Modified header value if present
	// Otherwise just remove atime and mtime
	mTime := time.Time***REMOVED******REMOVED***

	lastMod := resp.Header.Get("Last-Modified")
	if lastMod != "" ***REMOVED***
		// If we can't parse it then just let it default to 'zero'
		// otherwise use the parsed time value
		if parsedMTime, err := http.ParseTime(lastMod); err == nil ***REMOVED***
			mTime = parsedMTime
		***REMOVED***
	***REMOVED***

	tmpFile.Close()

	if err = system.Chtimes(tmpFileName, mTime, mTime); err != nil ***REMOVED***
		return
	***REMOVED***

	lc, err := remotecontext.NewLazySource(containerfs.NewLocalContainerFS(tmpDir))
	return lc, filename, err
***REMOVED***

type copyFileOptions struct ***REMOVED***
	decompress bool
	chownPair  idtools.IDPair
	archiver   Archiver
***REMOVED***

type copyEndpoint struct ***REMOVED***
	driver containerfs.Driver
	path   string
***REMOVED***

func performCopyForInfo(dest copyInfo, source copyInfo, options copyFileOptions) error ***REMOVED***
	srcPath, err := source.fullPath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	destPath, err := dest.fullPath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	archiver := options.archiver

	srcEndpoint := &copyEndpoint***REMOVED***driver: source.root, path: srcPath***REMOVED***
	destEndpoint := &copyEndpoint***REMOVED***driver: dest.root, path: destPath***REMOVED***

	src, err := source.root.Stat(srcPath)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "source path not found")
	***REMOVED***
	if src.IsDir() ***REMOVED***
		return copyDirectory(archiver, srcEndpoint, destEndpoint, options.chownPair)
	***REMOVED***
	if options.decompress && isArchivePath(source.root, srcPath) && !source.noDecompress ***REMOVED***
		return archiver.UntarPath(srcPath, destPath)
	***REMOVED***

	destExistsAsDir, err := isExistingDirectory(destEndpoint)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// dest.path must be used because destPath has already been cleaned of any
	// trailing slash
	if endsInSlash(dest.root, dest.path) || destExistsAsDir ***REMOVED***
		// source.path must be used to get the correct filename when the source
		// is a symlink
		destPath = dest.root.Join(destPath, source.root.Base(source.path))
		destEndpoint = &copyEndpoint***REMOVED***driver: dest.root, path: destPath***REMOVED***
	***REMOVED***
	return copyFile(archiver, srcEndpoint, destEndpoint, options.chownPair)
***REMOVED***

func isArchivePath(driver containerfs.ContainerFS, path string) bool ***REMOVED***
	file, err := driver.Open(path)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	defer file.Close()
	rdr, err := archive.DecompressStream(file)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	r := tar.NewReader(rdr)
	_, err = r.Next()
	return err == nil
***REMOVED***

func copyDirectory(archiver Archiver, source, dest *copyEndpoint, chownPair idtools.IDPair) error ***REMOVED***
	destExists, err := isExistingDirectory(dest)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to query destination path")
	***REMOVED***

	if err := archiver.CopyWithTar(source.path, dest.path); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to copy directory")
	***REMOVED***
	// TODO: @gupta-ak. Investigate how LCOW permission mappings will work.
	return fixPermissions(source.path, dest.path, chownPair, !destExists)
***REMOVED***

func copyFile(archiver Archiver, source, dest *copyEndpoint, chownPair idtools.IDPair) error ***REMOVED***
	if runtime.GOOS == "windows" && dest.driver.OS() == "linux" ***REMOVED***
		// LCOW
		if err := dest.driver.MkdirAll(dest.driver.Dir(dest.path), 0755); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to create new directory")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := idtools.MkdirAllAndChownNew(filepath.Dir(dest.path), 0755, chownPair); err != nil ***REMOVED***
			// Normal containers
			return errors.Wrapf(err, "failed to create new directory")
		***REMOVED***
	***REMOVED***

	if err := archiver.CopyFileWithTar(source.path, dest.path); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to copy file")
	***REMOVED***
	// TODO: @gupta-ak. Investigate how LCOW permission mappings will work.
	return fixPermissions(source.path, dest.path, chownPair, false)
***REMOVED***

func endsInSlash(driver containerfs.Driver, path string) bool ***REMOVED***
	return strings.HasSuffix(path, string(driver.Separator()))
***REMOVED***

// isExistingDirectory returns true if the path exists and is a directory
func isExistingDirectory(point *copyEndpoint) (bool, error) ***REMOVED***
	destStat, err := point.driver.Stat(point.path)
	switch ***REMOVED***
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	***REMOVED***
	return destStat.IsDir(), nil
***REMOVED***
