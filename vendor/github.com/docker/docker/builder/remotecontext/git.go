package remotecontext

import (
	"os"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/remotecontext/git"
	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
)

// MakeGitContext returns a Context from gitURL that is cloned in a temporary directory.
func MakeGitContext(gitURL string) (builder.Source, error) ***REMOVED***
	root, err := git.Clone(gitURL)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c, err := archive.Tar(root, archive.Uncompressed)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer func() ***REMOVED***
		err := c.Close()
		if err != nil ***REMOVED***
			logrus.WithField("action", "MakeGitContext").WithField("module", "builder").WithField("url", gitURL).WithError(err).Error("error while closing git context")
		***REMOVED***
		err = os.RemoveAll(root)
		if err != nil ***REMOVED***
			logrus.WithField("action", "MakeGitContext").WithField("module", "builder").WithField("url", gitURL).WithError(err).Error("error while removing path and children of root")
		***REMOVED***
	***REMOVED***()
	return FromArchive(c)
***REMOVED***
