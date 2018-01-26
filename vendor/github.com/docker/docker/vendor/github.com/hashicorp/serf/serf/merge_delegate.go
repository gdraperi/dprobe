package serf

import (
	"net"

	"github.com/hashicorp/memberlist"
)

type MergeDelegate interface ***REMOVED***
	NotifyMerge([]*Member) error
***REMOVED***

type mergeDelegate struct ***REMOVED***
	serf *Serf
***REMOVED***

func (m *mergeDelegate) NotifyMerge(nodes []*memberlist.Node) error ***REMOVED***
	members := make([]*Member, len(nodes))
	for idx, n := range nodes ***REMOVED***
		members[idx] = m.nodeToMember(n)
	***REMOVED***
	return m.serf.config.Merge.NotifyMerge(members)
***REMOVED***

func (m *mergeDelegate) NotifyAlive(peer *memberlist.Node) error ***REMOVED***
	member := m.nodeToMember(peer)
	return m.serf.config.Merge.NotifyMerge([]*Member***REMOVED***member***REMOVED***)
***REMOVED***

func (m *mergeDelegate) nodeToMember(n *memberlist.Node) *Member ***REMOVED***
	return &Member***REMOVED***
		Name:        n.Name,
		Addr:        net.IP(n.Addr),
		Port:        n.Port,
		Tags:        m.serf.decodeTags(n.Meta),
		Status:      StatusNone,
		ProtocolMin: n.PMin,
		ProtocolMax: n.PMax,
		ProtocolCur: n.PCur,
		DelegateMin: n.DMin,
		DelegateMax: n.DMax,
		DelegateCur: n.DCur,
	***REMOVED***
***REMOVED***
