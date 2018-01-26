package configs

import (
	"fmt"
)

type IfPrioMap struct ***REMOVED***
	Interface string `json:"interface"`
	Priority  int64  `json:"priority"`
***REMOVED***

func (i *IfPrioMap) CgroupString() string ***REMOVED***
	return fmt.Sprintf("%s %d", i.Interface, i.Priority)
***REMOVED***
