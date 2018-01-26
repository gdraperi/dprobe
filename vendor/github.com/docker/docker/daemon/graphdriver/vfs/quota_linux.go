package vfs

import (
	"github.com/docker/docker/daemon/graphdriver/quota"
	"github.com/sirupsen/logrus"
)

type driverQuota struct ***REMOVED***
	quotaCtl *quota.Control
***REMOVED***

func setupDriverQuota(driver *Driver) ***REMOVED***
	if quotaCtl, err := quota.NewControl(driver.home); err == nil ***REMOVED***
		driver.quotaCtl = quotaCtl
	***REMOVED*** else if err != quota.ErrQuotaNotSupported ***REMOVED***
		logrus.Warnf("Unable to setup quota: %v\n", err)
	***REMOVED***
***REMOVED***

func (d *Driver) setupQuota(dir string, size uint64) error ***REMOVED***
	return d.quotaCtl.SetQuota(dir, quota.Quota***REMOVED***Size: size***REMOVED***)
***REMOVED***

func (d *Driver) quotaSupported() bool ***REMOVED***
	return d.quotaCtl != nil
***REMOVED***
