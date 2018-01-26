package store

import (
	"sync"

	"github.com/boltdb/bolt"
	"github.com/docker/docker/volume"
	"github.com/docker/docker/volume/drivers"
	"github.com/sirupsen/logrus"
)

// restore is called when a new volume store is created.
// It's primary purpose is to ensure that all drivers' refcounts are set based
// on known volumes after a restart.
// This only attempts to track volumes that are actually stored in the on-disk db.
// It does not probe the available drivers to find anything that may have been added
// out of band.
func (s *VolumeStore) restore() ***REMOVED***
	var ls []volumeMetadata
	s.db.View(func(tx *bolt.Tx) error ***REMOVED***
		ls = listMeta(tx)
		return nil
	***REMOVED***)

	chRemove := make(chan *volumeMetadata, len(ls))
	var wg sync.WaitGroup
	for _, meta := range ls ***REMOVED***
		wg.Add(1)
		// this is potentially a very slow operation, so do it in a goroutine
		go func(meta volumeMetadata) ***REMOVED***
			defer wg.Done()

			var v volume.Volume
			var err error
			if meta.Driver != "" ***REMOVED***
				v, err = lookupVolume(meta.Driver, meta.Name)
				if err != nil && err != errNoSuchVolume ***REMOVED***
					logrus.WithError(err).WithField("driver", meta.Driver).WithField("volume", meta.Name).Warn("Error restoring volume")
					return
				***REMOVED***
				if v == nil ***REMOVED***
					// doesn't exist in the driver, remove it from the db
					chRemove <- &meta
					return
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				v, err = s.getVolume(meta.Name)
				if err != nil ***REMOVED***
					if err == errNoSuchVolume ***REMOVED***
						chRemove <- &meta
					***REMOVED***
					return
				***REMOVED***

				meta.Driver = v.DriverName()
				if err := s.setMeta(v.Name(), meta); err != nil ***REMOVED***
					logrus.WithError(err).WithField("driver", meta.Driver).WithField("volume", v.Name()).Warn("Error updating volume metadata on restore")
				***REMOVED***
			***REMOVED***

			// increment driver refcount
			volumedrivers.CreateDriver(meta.Driver)

			// cache the volume
			s.globalLock.Lock()
			s.options[v.Name()] = meta.Options
			s.labels[v.Name()] = meta.Labels
			s.names[v.Name()] = v
			s.globalLock.Unlock()
		***REMOVED***(meta)
	***REMOVED***

	wg.Wait()
	close(chRemove)
	s.db.Update(func(tx *bolt.Tx) error ***REMOVED***
		for meta := range chRemove ***REMOVED***
			if err := removeMeta(tx, meta.Name); err != nil ***REMOVED***
				logrus.WithField("volume", meta.Name).Warnf("Error removing stale entry from volume db: %v", err)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
***REMOVED***
