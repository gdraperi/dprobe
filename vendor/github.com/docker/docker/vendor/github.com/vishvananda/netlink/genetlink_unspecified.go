// +build !linux

package netlink

type GenlOp struct***REMOVED******REMOVED***

type GenlMulticastGroup struct***REMOVED******REMOVED***

type GenlFamily struct***REMOVED******REMOVED***

func (h *Handle) GenlFamilyList() ([]*GenlFamily, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func GenlFamilyList() ([]*GenlFamily, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func (h *Handle) GenlFamilyGet(name string) (*GenlFamily, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***

func GenlFamilyGet(name string) (*GenlFamily, error) ***REMOVED***
	return nil, ErrNotImplemented
***REMOVED***
