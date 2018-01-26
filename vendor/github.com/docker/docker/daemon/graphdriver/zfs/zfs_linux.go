package zfs

import (
	"github.com/docker/docker/daemon/graphdriver"
	"github.com/sirupsen/logrus"
)

func checkRootdirFs(rootDir string) error ***REMOVED***
	fsMagic, err := graphdriver.GetFSMagic(rootDir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	backingFS := "unknown"
	if fsName, ok := graphdriver.FsNames[fsMagic]; ok ***REMOVED***
		backingFS = fsName
	***REMOVED***

	if fsMagic != graphdriver.FsMagicZfs ***REMOVED***
		logrus.WithField("root", rootDir).WithField("backingFS", backingFS).WithField("driver", "zfs").Error("No zfs dataset found for root")
		return graphdriver.ErrPrerequisites
	***REMOVED***

	return nil
***REMOVED***

func getMountpoint(id string) string ***REMOVED***
	return id
***REMOVED***
