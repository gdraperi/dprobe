package dockerfile

import "github.com/docker/docker/pkg/idtools"

func parseChownFlag(chown, ctrRootPath string, idMappings *idtools.IDMappings) (idtools.IDPair, error) ***REMOVED***
	return idMappings.RootPair(), nil
***REMOVED***
