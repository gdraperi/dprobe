package hcsshim

import "github.com/sirupsen/logrus"

// LayerExists will return true if a layer with the given id exists and is known
// to the system.
func LayerExists(info DriverInfo, id string) (bool, error) ***REMOVED***
	title := "hcsshim::LayerExists "
	logrus.Debugf(title+"Flavour %d ID %s", info.Flavour, id)

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil ***REMOVED***
		logrus.Error(err)
		return false, err
	***REMOVED***

	// Call the procedure itself.
	var exists uint32

	err = layerExists(&infop, id, &exists)
	if err != nil ***REMOVED***
		err = makeErrorf(err, title, "id=%s flavour=%d", id, info.Flavour)
		logrus.Error(err)
		return false, err
	***REMOVED***

	logrus.Debugf(title+"succeeded flavour=%d id=%s exists=%d", info.Flavour, id, exists)
	return exists != 0, nil
***REMOVED***
