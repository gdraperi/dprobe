package blkiodev

import "fmt"

// WeightDevice is a structure that holds device:weight pair
type WeightDevice struct ***REMOVED***
	Path   string
	Weight uint16
***REMOVED***

func (w *WeightDevice) String() string ***REMOVED***
	return fmt.Sprintf("%s:%d", w.Path, w.Weight)
***REMOVED***

// ThrottleDevice is a structure that holds device:rate_per_second pair
type ThrottleDevice struct ***REMOVED***
	Path string
	Rate uint64
***REMOVED***

func (t *ThrottleDevice) String() string ***REMOVED***
	return fmt.Sprintf("%s:%d", t.Path, t.Rate)
***REMOVED***
