package remotecontext

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/containerd/continuity/driver"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/docker/docker/builder/dockerignore"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ClientSessionRemote is identifier for client-session context transport
const ClientSessionRemote = "client-session"

// Detect returns a context and dockerfile from remote location or local
// archive. progressReader is only used if remoteURL is actually a URL
// (not empty, and not a Git endpoint).
func Detect(config backend.BuildConfig) (remote builder.Source, dockerfile *parser.Result, err error) ***REMOVED***
	remoteURL := config.Options.RemoteContext
	dockerfilePath := config.Options.Dockerfile

	switch ***REMOVED***
	case remoteURL == "":
		remote, dockerfile, err = newArchiveRemote(config.Source, dockerfilePath)
	case remoteURL == ClientSessionRemote:
		res, err := parser.Parse(config.Source)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		return nil, res, nil
	case urlutil.IsGitURL(remoteURL):
		remote, dockerfile, err = newGitRemote(remoteURL, dockerfilePath)
	case urlutil.IsURL(remoteURL):
		remote, dockerfile, err = newURLRemote(remoteURL, dockerfilePath, config.ProgressWriter.ProgressReaderFunc)
	default:
		err = fmt.Errorf("remoteURL (%s) could not be recognized as URL", remoteURL)
	***REMOVED***
	return
***REMOVED***

func newArchiveRemote(rc io.ReadCloser, dockerfilePath string) (builder.Source, *parser.Result, error) ***REMOVED***
	defer rc.Close()
	c, err := FromArchive(rc)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return withDockerfileFromContext(c.(modifiableContext), dockerfilePath)
***REMOVED***

func withDockerfileFromContext(c modifiableContext, dockerfilePath string) (builder.Source, *parser.Result, error) ***REMOVED***
	df, err := openAt(c, dockerfilePath)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			if dockerfilePath == builder.DefaultDockerfileName ***REMOVED***
				lowercase := strings.ToLower(dockerfilePath)
				if _, err := StatAt(c, lowercase); err == nil ***REMOVED***
					return withDockerfileFromContext(c, lowercase)
				***REMOVED***
			***REMOVED***
			return nil, nil, errors.Errorf("Cannot locate specified Dockerfile: %s", dockerfilePath) // backwards compatible error
		***REMOVED***
		c.Close()
		return nil, nil, err
	***REMOVED***

	res, err := readAndParseDockerfile(dockerfilePath, df)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	df.Close()

	if err := removeDockerfile(c, dockerfilePath); err != nil ***REMOVED***
		c.Close()
		return nil, nil, err
	***REMOVED***

	return c, res, nil
***REMOVED***

func newGitRemote(gitURL string, dockerfilePath string) (builder.Source, *parser.Result, error) ***REMOVED***
	c, err := MakeGitContext(gitURL) // TODO: change this to NewLazySource
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return withDockerfileFromContext(c.(modifiableContext), dockerfilePath)
***REMOVED***

func newURLRemote(url string, dockerfilePath string, progressReader func(in io.ReadCloser) io.ReadCloser) (builder.Source, *parser.Result, error) ***REMOVED***
	contentType, content, err := downloadRemote(url)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer content.Close()

	switch contentType ***REMOVED***
	case mimeTypes.TextPlain:
		res, err := parser.Parse(progressReader(content))
		return nil, res, err
	default:
		source, err := FromArchive(progressReader(content))
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		return withDockerfileFromContext(source.(modifiableContext), dockerfilePath)
	***REMOVED***
***REMOVED***

func removeDockerfile(c modifiableContext, filesToRemove ...string) error ***REMOVED***
	f, err := openAt(c, ".dockerignore")
	// Note that a missing .dockerignore file isn't treated as an error
	switch ***REMOVED***
	case os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	***REMOVED***
	excludes, err := dockerignore.ReadAll(f)
	if err != nil ***REMOVED***
		f.Close()
		return err
	***REMOVED***
	f.Close()
	filesToRemove = append([]string***REMOVED***".dockerignore"***REMOVED***, filesToRemove...)
	for _, fileToRemove := range filesToRemove ***REMOVED***
		if rm, _ := fileutils.Matches(fileToRemove, excludes); rm ***REMOVED***
			if err := c.Remove(fileToRemove); err != nil ***REMOVED***
				logrus.Errorf("failed to remove %s: %v", fileToRemove, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func readAndParseDockerfile(name string, rc io.Reader) (*parser.Result, error) ***REMOVED***
	br := bufio.NewReader(rc)
	if _, err := br.Peek(1); err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			return nil, errors.Errorf("the Dockerfile (%s) cannot be empty", name)
		***REMOVED***
		return nil, errors.Wrap(err, "unexpected error reading Dockerfile")
	***REMOVED***
	return parser.Parse(br)
***REMOVED***

func openAt(remote builder.Source, path string) (driver.File, error) ***REMOVED***
	fullPath, err := FullPath(remote, path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return remote.Root().Open(fullPath)
***REMOVED***

// StatAt is a helper for calling Stat on a path from a source
func StatAt(remote builder.Source, path string) (os.FileInfo, error) ***REMOVED***
	fullPath, err := FullPath(remote, path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return remote.Root().Stat(fullPath)
***REMOVED***

// FullPath is a helper for getting a full path for a path from a source
func FullPath(remote builder.Source, path string) (string, error) ***REMOVED***
	fullPath, err := remote.Root().ResolveScopedPath(path, true)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("Forbidden path outside the build context: %s (%s)", path, fullPath) // backwards compat with old error
	***REMOVED***
	return fullPath, nil
***REMOVED***
