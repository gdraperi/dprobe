package hcsshim

import "github.com/sirupsen/logrus"

// ExpandSandboxSize expands the size of a layer to at least size bytes.
func ExpandSandboxSize(info DriverInfo, layerId string, size uint64) error ***REMOVED***
	title := "hcsshim::ExpandSandboxSize "
	logrus.Debugf(title+"layerId=%s size=%d", layerId, size)

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil ***REMOVED***
		logrus.Error(err)
		return err
	***REMOVED***

	err = expandSandboxSize(&infop, layerId, size)
	if err != nil ***REMOVED***
		err = makeErrorf(err, title, "layerId=%s  size=%d", layerId, size)
		logrus.Error(err)
		return err
	***REMOVED***

	logrus.Debugf(title+"- succeeded layerId=%s size=%d", layerId, size)
	return nil
***REMOVED***
