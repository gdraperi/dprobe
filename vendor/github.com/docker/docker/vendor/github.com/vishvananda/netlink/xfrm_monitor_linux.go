package netlink

import (
	"fmt"
	"syscall"

	"github.com/vishvananda/netns"

	"github.com/vishvananda/netlink/nl"
)

type XfrmMsg interface ***REMOVED***
	Type() nl.XfrmMsgType
***REMOVED***

type XfrmMsgExpire struct ***REMOVED***
	XfrmState *XfrmState
	Hard      bool
***REMOVED***

func (ue *XfrmMsgExpire) Type() nl.XfrmMsgType ***REMOVED***
	return nl.XFRM_MSG_EXPIRE
***REMOVED***

func parseXfrmMsgExpire(b []byte) *XfrmMsgExpire ***REMOVED***
	var e XfrmMsgExpire

	msg := nl.DeserializeXfrmUserExpire(b)
	e.XfrmState = xfrmStateFromXfrmUsersaInfo(&msg.XfrmUsersaInfo)
	e.Hard = msg.Hard == 1

	return &e
***REMOVED***

func XfrmMonitor(ch chan<- XfrmMsg, done <-chan struct***REMOVED******REMOVED***, errorChan chan<- error,
	types ...nl.XfrmMsgType) error ***REMOVED***

	groups, err := xfrmMcastGroups(types)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	s, err := nl.SubscribeAt(netns.None(), netns.None(), syscall.NETLINK_XFRM, groups...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if done != nil ***REMOVED***
		go func() ***REMOVED***
			<-done
			s.Close()
		***REMOVED***()

	***REMOVED***

	go func() ***REMOVED***
		defer close(ch)
		for ***REMOVED***
			msgs, err := s.Receive()
			if err != nil ***REMOVED***
				errorChan <- err
				return
			***REMOVED***
			for _, m := range msgs ***REMOVED***
				switch m.Header.Type ***REMOVED***
				case nl.XFRM_MSG_EXPIRE:
					ch <- parseXfrmMsgExpire(m.Data)
				default:
					errorChan <- fmt.Errorf("unsupported msg type: %x", m.Header.Type)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return nil
***REMOVED***

func xfrmMcastGroups(types []nl.XfrmMsgType) ([]uint, error) ***REMOVED***
	groups := make([]uint, 0)

	if len(types) == 0 ***REMOVED***
		return nil, fmt.Errorf("no xfrm msg type specified")
	***REMOVED***

	for _, t := range types ***REMOVED***
		var group uint

		switch t ***REMOVED***
		case nl.XFRM_MSG_EXPIRE:
			group = nl.XFRMNLGRP_EXPIRE
		default:
			return nil, fmt.Errorf("unsupported group: %x", t)
		***REMOVED***

		groups = append(groups, group)
	***REMOVED***

	return groups, nil
***REMOVED***
