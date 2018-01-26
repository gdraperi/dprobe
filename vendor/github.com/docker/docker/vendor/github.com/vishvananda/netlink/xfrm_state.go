package netlink

import (
	"fmt"
	"net"
)

// XfrmStateAlgo represents the algorithm to use for the ipsec encryption.
type XfrmStateAlgo struct ***REMOVED***
	Name        string
	Key         []byte
	TruncateLen int // Auth only
	ICVLen      int // AEAD only
***REMOVED***

func (a XfrmStateAlgo) String() string ***REMOVED***
	base := fmt.Sprintf("***REMOVED***Name: %s, Key: 0x%x", a.Name, a.Key)
	if a.TruncateLen != 0 ***REMOVED***
		base = fmt.Sprintf("%s, Truncate length: %d", base, a.TruncateLen)
	***REMOVED***
	if a.ICVLen != 0 ***REMOVED***
		base = fmt.Sprintf("%s, ICV length: %d", base, a.ICVLen)
	***REMOVED***
	return fmt.Sprintf("%s***REMOVED***", base)
***REMOVED***

// EncapType is an enum representing the optional packet encapsulation.
type EncapType uint8

const (
	XFRM_ENCAP_ESPINUDP_NONIKE EncapType = iota + 1
	XFRM_ENCAP_ESPINUDP
)

func (e EncapType) String() string ***REMOVED***
	switch e ***REMOVED***
	case XFRM_ENCAP_ESPINUDP_NONIKE:
		return "espinudp-non-ike"
	case XFRM_ENCAP_ESPINUDP:
		return "espinudp"
	***REMOVED***
	return "unknown"
***REMOVED***

// XfrmStateEncap represents the encapsulation to use for the ipsec encryption.
type XfrmStateEncap struct ***REMOVED***
	Type            EncapType
	SrcPort         int
	DstPort         int
	OriginalAddress net.IP
***REMOVED***

func (e XfrmStateEncap) String() string ***REMOVED***
	return fmt.Sprintf("***REMOVED***Type: %s, Srcport: %d, DstPort: %d, OriginalAddress: %v***REMOVED***",
		e.Type, e.SrcPort, e.DstPort, e.OriginalAddress)
***REMOVED***

// XfrmStateLimits represents the configured limits for the state.
type XfrmStateLimits struct ***REMOVED***
	ByteSoft    uint64
	ByteHard    uint64
	PacketSoft  uint64
	PacketHard  uint64
	TimeSoft    uint64
	TimeHard    uint64
	TimeUseSoft uint64
	TimeUseHard uint64
***REMOVED***

// XfrmState represents the state of an ipsec policy. It optionally
// contains an XfrmStateAlgo for encryption and one for authentication.
type XfrmState struct ***REMOVED***
	Dst          net.IP
	Src          net.IP
	Proto        Proto
	Mode         Mode
	Spi          int
	Reqid        int
	ReplayWindow int
	Limits       XfrmStateLimits
	Mark         *XfrmMark
	Auth         *XfrmStateAlgo
	Crypt        *XfrmStateAlgo
	Aead         *XfrmStateAlgo
	Encap        *XfrmStateEncap
	ESN          bool
***REMOVED***

func (sa XfrmState) String() string ***REMOVED***
	return fmt.Sprintf("Dst: %v, Src: %v, Proto: %s, Mode: %s, SPI: 0x%x, ReqID: 0x%x, ReplayWindow: %d, Mark: %v, Auth: %v, Crypt: %v, Aead: %v, Encap: %v, ESN: %t",
		sa.Dst, sa.Src, sa.Proto, sa.Mode, sa.Spi, sa.Reqid, sa.ReplayWindow, sa.Mark, sa.Auth, sa.Crypt, sa.Aead, sa.Encap, sa.ESN)
***REMOVED***
func (sa XfrmState) Print(stats bool) string ***REMOVED***
	if !stats ***REMOVED***
		return sa.String()
	***REMOVED***

	return fmt.Sprintf("%s, ByteSoft: %s, ByteHard: %s, PacketSoft: %s, PacketHard: %s, TimeSoft: %d, TimeHard: %d, TimeUseSoft: %d, TimeUseHard: %d",
		sa.String(), printLimit(sa.Limits.ByteSoft), printLimit(sa.Limits.ByteHard), printLimit(sa.Limits.PacketSoft), printLimit(sa.Limits.PacketHard),
		sa.Limits.TimeSoft, sa.Limits.TimeHard, sa.Limits.TimeUseSoft, sa.Limits.TimeUseHard)
***REMOVED***

func printLimit(lmt uint64) string ***REMOVED***
	if lmt == ^uint64(0) ***REMOVED***
		return "(INF)"
	***REMOVED***
	return fmt.Sprintf("%d", lmt)
***REMOVED***
