// +build !linux

package vfs

import "github.com/docker/docker/daemon/graphdriver/quota"

type driverQuota struct ***REMOVED***
***REMOVED***

func setupDriverQuota(driver *Driver) error ***REMOVED***
	return nil
***REMOVED***

func (d *Driver) setupQuota(dir string, size uint64) error ***REMOVED***
	return quota.ErrQuotaNotSupported
***REMOVED***

func (d *Driver) quotaSupported() bool ***REMOVED***
	return false
***REMOVED***
