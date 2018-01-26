package hcsshim

import "github.com/sirupsen/logrus"

// DestroyLayer will remove the on-disk files representing the layer with the given
// id, including that layer's containing folder, if any.
func DestroyLayer(info DriverInfo, id string) error ***REMOVED***
	title := "hcsshim::DestroyLayer "
	logrus.Debugf(title+"Flavour %d ID %s", info.Flavour, id)

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil ***REMOVED***
		logrus.Error(err)
		return err
	***REMOVED***

	err = destroyLayer(&infop, id)
	if err != nil ***REMOVED***
		err = makeErrorf(err, title, "id=%s flavour=%d", id, info.Flavour)
		logrus.Error(err)
		return err
	***REMOVED***

	logrus.Debugf(title+"succeeded flavour=%d id=%s", info.Flavour, id)
	return nil
***REMOVED***
