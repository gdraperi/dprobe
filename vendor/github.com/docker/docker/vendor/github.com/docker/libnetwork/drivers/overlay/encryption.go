package overlay

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"net"
	"sync"
	"syscall"

	"strconv"

	"github.com/docker/libnetwork/iptables"
	"github.com/docker/libnetwork/ns"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

const (
	r            = 0xD0C4E3
	pktExpansion = 26 // SPI(4) + SeqN(4) + IV(8) + PadLength(1) + NextHeader(1) + ICV(8)
)

const (
	forward = iota + 1
	reverse
	bidir
)

var spMark = netlink.XfrmMark***REMOVED***Value: uint32(r), Mask: 0xffffffff***REMOVED***

type key struct ***REMOVED***
	value []byte
	tag   uint32
***REMOVED***

func (k *key) String() string ***REMOVED***
	if k != nil ***REMOVED***
		return fmt.Sprintf("(key: %s, tag: 0x%x)", hex.EncodeToString(k.value)[0:5], k.tag)
	***REMOVED***
	return ""
***REMOVED***

type spi struct ***REMOVED***
	forward int
	reverse int
***REMOVED***

func (s *spi) String() string ***REMOVED***
	return fmt.Sprintf("SPI(FWD: 0x%x, REV: 0x%x)", uint32(s.forward), uint32(s.reverse))
***REMOVED***

type encrMap struct ***REMOVED***
	nodes map[string][]*spi
	sync.Mutex
***REMOVED***

func (e *encrMap) String() string ***REMOVED***
	e.Lock()
	defer e.Unlock()
	b := new(bytes.Buffer)
	for k, v := range e.nodes ***REMOVED***
		b.WriteString("\n")
		b.WriteString(k)
		b.WriteString(":")
		b.WriteString("[")
		for _, s := range v ***REMOVED***
			b.WriteString(s.String())
			b.WriteString(",")
		***REMOVED***
		b.WriteString("]")

	***REMOVED***
	return b.String()
***REMOVED***

func (d *driver) checkEncryption(nid string, rIP net.IP, vxlanID uint32, isLocal, add bool) error ***REMOVED***
	logrus.Debugf("checkEncryption(%s, %v, %d, %t)", nid[0:7], rIP, vxlanID, isLocal)

	n := d.network(nid)
	if n == nil || !n.secure ***REMOVED***
		return nil
	***REMOVED***

	if len(d.keys) == 0 ***REMOVED***
		return types.ForbiddenErrorf("encryption key is not present")
	***REMOVED***

	lIP := net.ParseIP(d.bindAddress)
	aIP := net.ParseIP(d.advertiseAddress)
	nodes := map[string]net.IP***REMOVED******REMOVED***

	switch ***REMOVED***
	case isLocal:
		if err := d.peerDbNetworkWalk(nid, func(pKey *peerKey, pEntry *peerEntry) bool ***REMOVED***
			if !aIP.Equal(pEntry.vtep) ***REMOVED***
				nodes[pEntry.vtep.String()] = pEntry.vtep
			***REMOVED***
			return false
		***REMOVED***); err != nil ***REMOVED***
			logrus.Warnf("Failed to retrieve list of participating nodes in overlay network %s: %v", nid[0:5], err)
		***REMOVED***
	default:
		if len(d.network(nid).endpoints) > 0 ***REMOVED***
			nodes[rIP.String()] = rIP
		***REMOVED***
	***REMOVED***

	logrus.Debugf("List of nodes: %s", nodes)

	if add ***REMOVED***
		for _, rIP := range nodes ***REMOVED***
			if err := setupEncryption(lIP, aIP, rIP, vxlanID, d.secMap, d.keys); err != nil ***REMOVED***
				logrus.Warnf("Failed to program network encryption between %s and %s: %v", lIP, rIP, err)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if len(nodes) == 0 ***REMOVED***
			if err := removeEncryption(lIP, rIP, d.secMap); err != nil ***REMOVED***
				logrus.Warnf("Failed to remove network encryption between %s and %s: %v", lIP, rIP, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func setupEncryption(localIP, advIP, remoteIP net.IP, vni uint32, em *encrMap, keys []*key) error ***REMOVED***
	logrus.Debugf("Programming encryption for vxlan %d between %s and %s", vni, localIP, remoteIP)
	rIPs := remoteIP.String()

	indices := make([]*spi, 0, len(keys))

	err := programMangle(vni, true)
	if err != nil ***REMOVED***
		logrus.Warn(err)
	***REMOVED***

	err = programInput(vni, true)
	if err != nil ***REMOVED***
		logrus.Warn(err)
	***REMOVED***

	for i, k := range keys ***REMOVED***
		spis := &spi***REMOVED***buildSPI(advIP, remoteIP, k.tag), buildSPI(remoteIP, advIP, k.tag)***REMOVED***
		dir := reverse
		if i == 0 ***REMOVED***
			dir = bidir
		***REMOVED***
		fSA, rSA, err := programSA(localIP, remoteIP, spis, k, dir, true)
		if err != nil ***REMOVED***
			logrus.Warn(err)
		***REMOVED***
		indices = append(indices, spis)
		if i != 0 ***REMOVED***
			continue
		***REMOVED***
		err = programSP(fSA, rSA, true)
		if err != nil ***REMOVED***
			logrus.Warn(err)
		***REMOVED***
	***REMOVED***

	em.Lock()
	em.nodes[rIPs] = indices
	em.Unlock()

	return nil
***REMOVED***

func removeEncryption(localIP, remoteIP net.IP, em *encrMap) error ***REMOVED***
	em.Lock()
	indices, ok := em.nodes[remoteIP.String()]
	em.Unlock()
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	for i, idxs := range indices ***REMOVED***
		dir := reverse
		if i == 0 ***REMOVED***
			dir = bidir
		***REMOVED***
		fSA, rSA, err := programSA(localIP, remoteIP, idxs, nil, dir, false)
		if err != nil ***REMOVED***
			logrus.Warn(err)
		***REMOVED***
		if i != 0 ***REMOVED***
			continue
		***REMOVED***
		err = programSP(fSA, rSA, false)
		if err != nil ***REMOVED***
			logrus.Warn(err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func programMangle(vni uint32, add bool) (err error) ***REMOVED***
	var (
		p      = strconv.FormatUint(uint64(vxlanPort), 10)
		c      = fmt.Sprintf("0>>22&0x3C@12&0xFFFFFF00=%d", int(vni)<<8)
		m      = strconv.FormatUint(uint64(r), 10)
		chain  = "OUTPUT"
		rule   = []string***REMOVED***"-p", "udp", "--dport", p, "-m", "u32", "--u32", c, "-j", "MARK", "--set-mark", m***REMOVED***
		a      = "-A"
		action = "install"
	)

	if add == iptables.Exists(iptables.Mangle, chain, rule...) ***REMOVED***
		return
	***REMOVED***

	if !add ***REMOVED***
		a = "-D"
		action = "remove"
	***REMOVED***

	if err = iptables.RawCombinedOutput(append([]string***REMOVED***"-t", string(iptables.Mangle), a, chain***REMOVED***, rule...)...); err != nil ***REMOVED***
		logrus.Warnf("could not %s mangle rule: %v", action, err)
	***REMOVED***

	return
***REMOVED***

func programInput(vni uint32, add bool) (err error) ***REMOVED***
	var (
		port       = strconv.FormatUint(uint64(vxlanPort), 10)
		vniMatch   = fmt.Sprintf("0>>22&0x3C@12&0xFFFFFF00=%d", int(vni)<<8)
		plainVxlan = []string***REMOVED***"-p", "udp", "--dport", port, "-m", "u32", "--u32", vniMatch, "-j"***REMOVED***
		ipsecVxlan = append([]string***REMOVED***"-m", "policy", "--dir", "in", "--pol", "ipsec"***REMOVED***, plainVxlan...)
		block      = append(plainVxlan, "DROP")
		accept     = append(ipsecVxlan, "ACCEPT")
		chain      = "INPUT"
		action     = iptables.Append
		msg        = "add"
	)

	if !add ***REMOVED***
		action = iptables.Delete
		msg = "remove"
	***REMOVED***

	if err := iptables.ProgramRule(iptables.Filter, chain, action, accept); err != nil ***REMOVED***
		logrus.Errorf("could not %s input rule: %v. Please do it manually.", msg, err)
	***REMOVED***

	if err := iptables.ProgramRule(iptables.Filter, chain, action, block); err != nil ***REMOVED***
		logrus.Errorf("could not %s input rule: %v. Please do it manually.", msg, err)
	***REMOVED***

	return
***REMOVED***

func programSA(localIP, remoteIP net.IP, spi *spi, k *key, dir int, add bool) (fSA *netlink.XfrmState, rSA *netlink.XfrmState, err error) ***REMOVED***
	var (
		action      = "Removing"
		xfrmProgram = ns.NlHandle().XfrmStateDel
	)

	if add ***REMOVED***
		action = "Adding"
		xfrmProgram = ns.NlHandle().XfrmStateAdd
	***REMOVED***

	if dir&reverse > 0 ***REMOVED***
		rSA = &netlink.XfrmState***REMOVED***
			Src:   remoteIP,
			Dst:   localIP,
			Proto: netlink.XFRM_PROTO_ESP,
			Spi:   spi.reverse,
			Mode:  netlink.XFRM_MODE_TRANSPORT,
			Reqid: r,
		***REMOVED***
		if add ***REMOVED***
			rSA.Aead = buildAeadAlgo(k, spi.reverse)
		***REMOVED***

		exists, err := saExists(rSA)
		if err != nil ***REMOVED***
			exists = !add
		***REMOVED***

		if add != exists ***REMOVED***
			logrus.Debugf("%s: rSA***REMOVED***%s***REMOVED***", action, rSA)
			if err := xfrmProgram(rSA); err != nil ***REMOVED***
				logrus.Warnf("Failed %s rSA***REMOVED***%s***REMOVED***: %v", action, rSA, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if dir&forward > 0 ***REMOVED***
		fSA = &netlink.XfrmState***REMOVED***
			Src:   localIP,
			Dst:   remoteIP,
			Proto: netlink.XFRM_PROTO_ESP,
			Spi:   spi.forward,
			Mode:  netlink.XFRM_MODE_TRANSPORT,
			Reqid: r,
		***REMOVED***
		if add ***REMOVED***
			fSA.Aead = buildAeadAlgo(k, spi.forward)
		***REMOVED***

		exists, err := saExists(fSA)
		if err != nil ***REMOVED***
			exists = !add
		***REMOVED***

		if add != exists ***REMOVED***
			logrus.Debugf("%s fSA***REMOVED***%s***REMOVED***", action, fSA)
			if err := xfrmProgram(fSA); err != nil ***REMOVED***
				logrus.Warnf("Failed %s fSA***REMOVED***%s***REMOVED***: %v.", action, fSA, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func programSP(fSA *netlink.XfrmState, rSA *netlink.XfrmState, add bool) error ***REMOVED***
	action := "Removing"
	xfrmProgram := ns.NlHandle().XfrmPolicyDel
	if add ***REMOVED***
		action = "Adding"
		xfrmProgram = ns.NlHandle().XfrmPolicyAdd
	***REMOVED***

	// Create a congruent cidr
	s := types.GetMinimalIP(fSA.Src)
	d := types.GetMinimalIP(fSA.Dst)
	fullMask := net.CIDRMask(8*len(s), 8*len(s))

	fPol := &netlink.XfrmPolicy***REMOVED***
		Src:     &net.IPNet***REMOVED***IP: s, Mask: fullMask***REMOVED***,
		Dst:     &net.IPNet***REMOVED***IP: d, Mask: fullMask***REMOVED***,
		Dir:     netlink.XFRM_DIR_OUT,
		Proto:   17,
		DstPort: 4789,
		Mark:    &spMark,
		Tmpls: []netlink.XfrmPolicyTmpl***REMOVED***
			***REMOVED***
				Src:   fSA.Src,
				Dst:   fSA.Dst,
				Proto: netlink.XFRM_PROTO_ESP,
				Mode:  netlink.XFRM_MODE_TRANSPORT,
				Spi:   fSA.Spi,
				Reqid: r,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	exists, err := spExists(fPol)
	if err != nil ***REMOVED***
		exists = !add
	***REMOVED***

	if add != exists ***REMOVED***
		logrus.Debugf("%s fSP***REMOVED***%s***REMOVED***", action, fPol)
		if err := xfrmProgram(fPol); err != nil ***REMOVED***
			logrus.Warnf("%s fSP***REMOVED***%s***REMOVED***: %v", action, fPol, err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func saExists(sa *netlink.XfrmState) (bool, error) ***REMOVED***
	_, err := ns.NlHandle().XfrmStateGet(sa)
	switch err ***REMOVED***
	case nil:
		return true, nil
	case syscall.ESRCH:
		return false, nil
	default:
		err = fmt.Errorf("Error while checking for SA existence: %v", err)
		logrus.Warn(err)
		return false, err
	***REMOVED***
***REMOVED***

func spExists(sp *netlink.XfrmPolicy) (bool, error) ***REMOVED***
	_, err := ns.NlHandle().XfrmPolicyGet(sp)
	switch err ***REMOVED***
	case nil:
		return true, nil
	case syscall.ENOENT:
		return false, nil
	default:
		err = fmt.Errorf("Error while checking for SP existence: %v", err)
		logrus.Warn(err)
		return false, err
	***REMOVED***
***REMOVED***

func buildSPI(src, dst net.IP, st uint32) int ***REMOVED***
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, st)
	h := fnv.New32a()
	h.Write(src)
	h.Write(b)
	h.Write(dst)
	return int(binary.BigEndian.Uint32(h.Sum(nil)))
***REMOVED***

func buildAeadAlgo(k *key, s int) *netlink.XfrmStateAlgo ***REMOVED***
	salt := make([]byte, 4)
	binary.BigEndian.PutUint32(salt, uint32(s))
	return &netlink.XfrmStateAlgo***REMOVED***
		Name:   "rfc4106(gcm(aes))",
		Key:    append(k.value, salt...),
		ICVLen: 64,
	***REMOVED***
***REMOVED***

func (d *driver) secMapWalk(f func(string, []*spi) ([]*spi, bool)) error ***REMOVED***
	d.secMap.Lock()
	for node, indices := range d.secMap.nodes ***REMOVED***
		idxs, stop := f(node, indices)
		if idxs != nil ***REMOVED***
			d.secMap.nodes[node] = idxs
		***REMOVED***
		if stop ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	d.secMap.Unlock()
	return nil
***REMOVED***

func (d *driver) setKeys(keys []*key) error ***REMOVED***
	// Remove any stale policy, state
	clearEncryptionStates()
	// Accept the encryption keys and clear any stale encryption map
	d.Lock()
	d.keys = keys
	d.secMap = &encrMap***REMOVED***nodes: map[string][]*spi***REMOVED******REMOVED******REMOVED***
	d.Unlock()
	logrus.Debugf("Initial encryption keys: %v", d.keys)
	return nil
***REMOVED***

// updateKeys allows to add a new key and/or change the primary key and/or prune an existing key
// The primary key is the key used in transmission and will go in first position in the list.
func (d *driver) updateKeys(newKey, primary, pruneKey *key) error ***REMOVED***
	logrus.Debugf("Updating Keys. New: %v, Primary: %v, Pruned: %v", newKey, primary, pruneKey)

	logrus.Debugf("Current: %v", d.keys)

	var (
		newIdx = -1
		priIdx = -1
		delIdx = -1
		lIP    = net.ParseIP(d.bindAddress)
		aIP    = net.ParseIP(d.advertiseAddress)
	)

	d.Lock()
	// add new
	if newKey != nil ***REMOVED***
		d.keys = append(d.keys, newKey)
		newIdx += len(d.keys)
	***REMOVED***
	for i, k := range d.keys ***REMOVED***
		if primary != nil && k.tag == primary.tag ***REMOVED***
			priIdx = i
		***REMOVED***
		if pruneKey != nil && k.tag == pruneKey.tag ***REMOVED***
			delIdx = i
		***REMOVED***
	***REMOVED***
	d.Unlock()

	if (newKey != nil && newIdx == -1) ||
		(primary != nil && priIdx == -1) ||
		(pruneKey != nil && delIdx == -1) ***REMOVED***
		return types.BadRequestErrorf("cannot find proper key indices while processing key update:"+
			"(newIdx,priIdx,delIdx):(%d, %d, %d)", newIdx, priIdx, delIdx)
	***REMOVED***

	d.secMapWalk(func(rIPs string, spis []*spi) ([]*spi, bool) ***REMOVED***
		rIP := net.ParseIP(rIPs)
		return updateNodeKey(lIP, aIP, rIP, spis, d.keys, newIdx, priIdx, delIdx), false
	***REMOVED***)

	d.Lock()
	// swap primary
	if priIdx != -1 ***REMOVED***
		swp := d.keys[0]
		d.keys[0] = d.keys[priIdx]
		d.keys[priIdx] = swp
	***REMOVED***
	// prune
	if delIdx != -1 ***REMOVED***
		if delIdx == 0 ***REMOVED***
			delIdx = priIdx
		***REMOVED***
		d.keys = append(d.keys[:delIdx], d.keys[delIdx+1:]...)
	***REMOVED***
	d.Unlock()

	logrus.Debugf("Updated: %v", d.keys)

	return nil
***REMOVED***

/********************************************************
 * Steady state: rSA0, rSA1, rSA2, fSA1, fSP1
 * Rotation --> -rSA0, +rSA3, +fSA2, +fSP2/-fSP1, -fSA1
 * Steady state: rSA1, rSA2, rSA3, fSA2, fSP2
 *********************************************************/

// Spis and keys are sorted in such away the one in position 0 is the primary
func updateNodeKey(lIP, aIP, rIP net.IP, idxs []*spi, curKeys []*key, newIdx, priIdx, delIdx int) []*spi ***REMOVED***
	logrus.Debugf("Updating keys for node: %s (%d,%d,%d)", rIP, newIdx, priIdx, delIdx)

	spis := idxs
	logrus.Debugf("Current: %v", spis)

	// add new
	if newIdx != -1 ***REMOVED***
		spis = append(spis, &spi***REMOVED***
			forward: buildSPI(aIP, rIP, curKeys[newIdx].tag),
			reverse: buildSPI(rIP, aIP, curKeys[newIdx].tag),
		***REMOVED***)
	***REMOVED***

	if delIdx != -1 ***REMOVED***
		// -rSA0
		programSA(lIP, rIP, spis[delIdx], nil, reverse, false)
	***REMOVED***

	if newIdx > -1 ***REMOVED***
		// +rSA2
		programSA(lIP, rIP, spis[newIdx], curKeys[newIdx], reverse, true)
	***REMOVED***

	if priIdx > 0 ***REMOVED***
		// +fSA2
		fSA2, _, _ := programSA(lIP, rIP, spis[priIdx], curKeys[priIdx], forward, true)

		// +fSP2, -fSP1
		s := types.GetMinimalIP(fSA2.Src)
		d := types.GetMinimalIP(fSA2.Dst)
		fullMask := net.CIDRMask(8*len(s), 8*len(s))

		fSP1 := &netlink.XfrmPolicy***REMOVED***
			Src:     &net.IPNet***REMOVED***IP: s, Mask: fullMask***REMOVED***,
			Dst:     &net.IPNet***REMOVED***IP: d, Mask: fullMask***REMOVED***,
			Dir:     netlink.XFRM_DIR_OUT,
			Proto:   17,
			DstPort: 4789,
			Mark:    &spMark,
			Tmpls: []netlink.XfrmPolicyTmpl***REMOVED***
				***REMOVED***
					Src:   fSA2.Src,
					Dst:   fSA2.Dst,
					Proto: netlink.XFRM_PROTO_ESP,
					Mode:  netlink.XFRM_MODE_TRANSPORT,
					Spi:   fSA2.Spi,
					Reqid: r,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		logrus.Debugf("Updating fSP***REMOVED***%s***REMOVED***", fSP1)
		if err := ns.NlHandle().XfrmPolicyUpdate(fSP1); err != nil ***REMOVED***
			logrus.Warnf("Failed to update fSP***REMOVED***%s***REMOVED***: %v", fSP1, err)
		***REMOVED***

		// -fSA1
		programSA(lIP, rIP, spis[0], nil, forward, false)
	***REMOVED***

	// swap
	if priIdx > 0 ***REMOVED***
		swp := spis[0]
		spis[0] = spis[priIdx]
		spis[priIdx] = swp
	***REMOVED***
	// prune
	if delIdx != -1 ***REMOVED***
		if delIdx == 0 ***REMOVED***
			delIdx = priIdx
		***REMOVED***
		spis = append(spis[:delIdx], spis[delIdx+1:]...)
	***REMOVED***

	logrus.Debugf("Updated: %v", spis)

	return spis
***REMOVED***

func (n *network) maxMTU() int ***REMOVED***
	mtu := 1500
	if n.mtu != 0 ***REMOVED***
		mtu = n.mtu
	***REMOVED***
	mtu -= vxlanEncap
	if n.secure ***REMOVED***
		// In case of encryption account for the
		// esp packet espansion and padding
		mtu -= pktExpansion
		mtu -= (mtu % 4)
	***REMOVED***
	return mtu
***REMOVED***

func clearEncryptionStates() ***REMOVED***
	nlh := ns.NlHandle()
	spList, err := nlh.XfrmPolicyList(netlink.FAMILY_ALL)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to retrieve SP list for cleanup: %v", err)
	***REMOVED***
	saList, err := nlh.XfrmStateList(netlink.FAMILY_ALL)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to retrieve SA list for cleanup: %v", err)
	***REMOVED***
	for _, sp := range spList ***REMOVED***
		if sp.Mark != nil && sp.Mark.Value == spMark.Value ***REMOVED***
			if err := nlh.XfrmPolicyDel(&sp); err != nil ***REMOVED***
				logrus.Warnf("Failed to delete stale SP %s: %v", sp, err)
				continue
			***REMOVED***
			logrus.Debugf("Removed stale SP: %s", sp)
		***REMOVED***
	***REMOVED***
	for _, sa := range saList ***REMOVED***
		if sa.Reqid == r ***REMOVED***
			if err := nlh.XfrmStateDel(&sa); err != nil ***REMOVED***
				logrus.Warnf("Failed to delete stale SA %s: %v", sa, err)
				continue
			***REMOVED***
			logrus.Debugf("Removed stale SA: %s", sa)
		***REMOVED***
	***REMOVED***
***REMOVED***
