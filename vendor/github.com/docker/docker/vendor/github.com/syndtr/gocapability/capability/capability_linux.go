// Copyright (c) 2013, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package capability

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
)

var errUnknownVers = errors.New("unknown capability version")

const (
	linuxCapVer1 = 0x19980330
	linuxCapVer2 = 0x20071026
	linuxCapVer3 = 0x20080522
)

var (
	capVers    uint32
	capLastCap Cap
)

func init() ***REMOVED***
	var hdr capHeader
	capget(&hdr, nil)
	capVers = hdr.version

	if initLastCap() == nil ***REMOVED***
		CAP_LAST_CAP = capLastCap
		if capLastCap > 31 ***REMOVED***
			capUpperMask = (uint32(1) << (uint(capLastCap) - 31)) - 1
		***REMOVED*** else ***REMOVED***
			capUpperMask = 0
		***REMOVED***
	***REMOVED***
***REMOVED***

func initLastCap() error ***REMOVED***
	if capLastCap != 0 ***REMOVED***
		return nil
	***REMOVED***

	f, err := os.Open("/proc/sys/kernel/cap_last_cap")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	var b []byte = make([]byte, 11)
	_, err = f.Read(b)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	fmt.Sscanf(string(b), "%d", &capLastCap)

	return nil
***REMOVED***

func mkStringCap(c Capabilities, which CapType) (ret string) ***REMOVED***
	for i, first := Cap(0), true; i <= CAP_LAST_CAP; i++ ***REMOVED***
		if !c.Get(which, i) ***REMOVED***
			continue
		***REMOVED***
		if first ***REMOVED***
			first = false
		***REMOVED*** else ***REMOVED***
			ret += ", "
		***REMOVED***
		ret += i.String()
	***REMOVED***
	return
***REMOVED***

func mkString(c Capabilities, max CapType) (ret string) ***REMOVED***
	ret = "***REMOVED***"
	for i := CapType(1); i <= max; i <<= 1 ***REMOVED***
		ret += " " + i.String() + "=\""
		if c.Empty(i) ***REMOVED***
			ret += "empty"
		***REMOVED*** else if c.Full(i) ***REMOVED***
			ret += "full"
		***REMOVED*** else ***REMOVED***
			ret += c.StringCap(i)
		***REMOVED***
		ret += "\""
	***REMOVED***
	ret += " ***REMOVED***"
	return
***REMOVED***

func newPid(pid int) (c Capabilities, err error) ***REMOVED***
	switch capVers ***REMOVED***
	case linuxCapVer1:
		p := new(capsV1)
		p.hdr.version = capVers
		p.hdr.pid = pid
		c = p
	case linuxCapVer2, linuxCapVer3:
		p := new(capsV3)
		p.hdr.version = capVers
		p.hdr.pid = pid
		c = p
	default:
		err = errUnknownVers
		return
	***REMOVED***
	err = c.Load()
	if err != nil ***REMOVED***
		c = nil
	***REMOVED***
	return
***REMOVED***

type capsV1 struct ***REMOVED***
	hdr  capHeader
	data capData
***REMOVED***

func (c *capsV1) Get(which CapType, what Cap) bool ***REMOVED***
	if what > 32 ***REMOVED***
		return false
	***REMOVED***

	switch which ***REMOVED***
	case EFFECTIVE:
		return (1<<uint(what))&c.data.effective != 0
	case PERMITTED:
		return (1<<uint(what))&c.data.permitted != 0
	case INHERITABLE:
		return (1<<uint(what))&c.data.inheritable != 0
	***REMOVED***

	return false
***REMOVED***

func (c *capsV1) getData(which CapType) (ret uint32) ***REMOVED***
	switch which ***REMOVED***
	case EFFECTIVE:
		ret = c.data.effective
	case PERMITTED:
		ret = c.data.permitted
	case INHERITABLE:
		ret = c.data.inheritable
	***REMOVED***
	return
***REMOVED***

func (c *capsV1) Empty(which CapType) bool ***REMOVED***
	return c.getData(which) == 0
***REMOVED***

func (c *capsV1) Full(which CapType) bool ***REMOVED***
	return (c.getData(which) & 0x7fffffff) == 0x7fffffff
***REMOVED***

func (c *capsV1) Set(which CapType, caps ...Cap) ***REMOVED***
	for _, what := range caps ***REMOVED***
		if what > 32 ***REMOVED***
			continue
		***REMOVED***

		if which&EFFECTIVE != 0 ***REMOVED***
			c.data.effective |= 1 << uint(what)
		***REMOVED***
		if which&PERMITTED != 0 ***REMOVED***
			c.data.permitted |= 1 << uint(what)
		***REMOVED***
		if which&INHERITABLE != 0 ***REMOVED***
			c.data.inheritable |= 1 << uint(what)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *capsV1) Unset(which CapType, caps ...Cap) ***REMOVED***
	for _, what := range caps ***REMOVED***
		if what > 32 ***REMOVED***
			continue
		***REMOVED***

		if which&EFFECTIVE != 0 ***REMOVED***
			c.data.effective &= ^(1 << uint(what))
		***REMOVED***
		if which&PERMITTED != 0 ***REMOVED***
			c.data.permitted &= ^(1 << uint(what))
		***REMOVED***
		if which&INHERITABLE != 0 ***REMOVED***
			c.data.inheritable &= ^(1 << uint(what))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *capsV1) Fill(kind CapType) ***REMOVED***
	if kind&CAPS == CAPS ***REMOVED***
		c.data.effective = 0x7fffffff
		c.data.permitted = 0x7fffffff
		c.data.inheritable = 0
	***REMOVED***
***REMOVED***

func (c *capsV1) Clear(kind CapType) ***REMOVED***
	if kind&CAPS == CAPS ***REMOVED***
		c.data.effective = 0
		c.data.permitted = 0
		c.data.inheritable = 0
	***REMOVED***
***REMOVED***

func (c *capsV1) StringCap(which CapType) (ret string) ***REMOVED***
	return mkStringCap(c, which)
***REMOVED***

func (c *capsV1) String() (ret string) ***REMOVED***
	return mkString(c, BOUNDING)
***REMOVED***

func (c *capsV1) Load() (err error) ***REMOVED***
	return capget(&c.hdr, &c.data)
***REMOVED***

func (c *capsV1) Apply(kind CapType) error ***REMOVED***
	if kind&CAPS == CAPS ***REMOVED***
		return capset(&c.hdr, &c.data)
	***REMOVED***
	return nil
***REMOVED***

type capsV3 struct ***REMOVED***
	hdr    capHeader
	data   [2]capData
	bounds [2]uint32
***REMOVED***

func (c *capsV3) Get(which CapType, what Cap) bool ***REMOVED***
	var i uint
	if what > 31 ***REMOVED***
		i = uint(what) >> 5
		what %= 32
	***REMOVED***

	switch which ***REMOVED***
	case EFFECTIVE:
		return (1<<uint(what))&c.data[i].effective != 0
	case PERMITTED:
		return (1<<uint(what))&c.data[i].permitted != 0
	case INHERITABLE:
		return (1<<uint(what))&c.data[i].inheritable != 0
	case BOUNDING:
		return (1<<uint(what))&c.bounds[i] != 0
	***REMOVED***

	return false
***REMOVED***

func (c *capsV3) getData(which CapType, dest []uint32) ***REMOVED***
	switch which ***REMOVED***
	case EFFECTIVE:
		dest[0] = c.data[0].effective
		dest[1] = c.data[1].effective
	case PERMITTED:
		dest[0] = c.data[0].permitted
		dest[1] = c.data[1].permitted
	case INHERITABLE:
		dest[0] = c.data[0].inheritable
		dest[1] = c.data[1].inheritable
	case BOUNDING:
		dest[0] = c.bounds[0]
		dest[1] = c.bounds[1]
	***REMOVED***
***REMOVED***

func (c *capsV3) Empty(which CapType) bool ***REMOVED***
	var data [2]uint32
	c.getData(which, data[:])
	return data[0] == 0 && data[1] == 0
***REMOVED***

func (c *capsV3) Full(which CapType) bool ***REMOVED***
	var data [2]uint32
	c.getData(which, data[:])
	if (data[0] & 0xffffffff) != 0xffffffff ***REMOVED***
		return false
	***REMOVED***
	return (data[1] & capUpperMask) == capUpperMask
***REMOVED***

func (c *capsV3) Set(which CapType, caps ...Cap) ***REMOVED***
	for _, what := range caps ***REMOVED***
		var i uint
		if what > 31 ***REMOVED***
			i = uint(what) >> 5
			what %= 32
		***REMOVED***

		if which&EFFECTIVE != 0 ***REMOVED***
			c.data[i].effective |= 1 << uint(what)
		***REMOVED***
		if which&PERMITTED != 0 ***REMOVED***
			c.data[i].permitted |= 1 << uint(what)
		***REMOVED***
		if which&INHERITABLE != 0 ***REMOVED***
			c.data[i].inheritable |= 1 << uint(what)
		***REMOVED***
		if which&BOUNDING != 0 ***REMOVED***
			c.bounds[i] |= 1 << uint(what)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *capsV3) Unset(which CapType, caps ...Cap) ***REMOVED***
	for _, what := range caps ***REMOVED***
		var i uint
		if what > 31 ***REMOVED***
			i = uint(what) >> 5
			what %= 32
		***REMOVED***

		if which&EFFECTIVE != 0 ***REMOVED***
			c.data[i].effective &= ^(1 << uint(what))
		***REMOVED***
		if which&PERMITTED != 0 ***REMOVED***
			c.data[i].permitted &= ^(1 << uint(what))
		***REMOVED***
		if which&INHERITABLE != 0 ***REMOVED***
			c.data[i].inheritable &= ^(1 << uint(what))
		***REMOVED***
		if which&BOUNDING != 0 ***REMOVED***
			c.bounds[i] &= ^(1 << uint(what))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *capsV3) Fill(kind CapType) ***REMOVED***
	if kind&CAPS == CAPS ***REMOVED***
		c.data[0].effective = 0xffffffff
		c.data[0].permitted = 0xffffffff
		c.data[0].inheritable = 0
		c.data[1].effective = 0xffffffff
		c.data[1].permitted = 0xffffffff
		c.data[1].inheritable = 0
	***REMOVED***

	if kind&BOUNDS == BOUNDS ***REMOVED***
		c.bounds[0] = 0xffffffff
		c.bounds[1] = 0xffffffff
	***REMOVED***
***REMOVED***

func (c *capsV3) Clear(kind CapType) ***REMOVED***
	if kind&CAPS == CAPS ***REMOVED***
		c.data[0].effective = 0
		c.data[0].permitted = 0
		c.data[0].inheritable = 0
		c.data[1].effective = 0
		c.data[1].permitted = 0
		c.data[1].inheritable = 0
	***REMOVED***

	if kind&BOUNDS == BOUNDS ***REMOVED***
		c.bounds[0] = 0
		c.bounds[1] = 0
	***REMOVED***
***REMOVED***

func (c *capsV3) StringCap(which CapType) (ret string) ***REMOVED***
	return mkStringCap(c, which)
***REMOVED***

func (c *capsV3) String() (ret string) ***REMOVED***
	return mkString(c, BOUNDING)
***REMOVED***

func (c *capsV3) Load() (err error) ***REMOVED***
	err = capget(&c.hdr, &c.data[0])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	var status_path string

	if c.hdr.pid == 0 ***REMOVED***
		status_path = fmt.Sprintf("/proc/self/status")
	***REMOVED*** else ***REMOVED***
		status_path = fmt.Sprintf("/proc/%d/status", c.hdr.pid)
	***REMOVED***

	f, err := os.Open(status_path)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	b := bufio.NewReader(f)
	for ***REMOVED***
		line, e := b.ReadString('\n')
		if e != nil ***REMOVED***
			if e != io.EOF ***REMOVED***
				err = e
			***REMOVED***
			break
		***REMOVED***
		if strings.HasPrefix(line, "CapB") ***REMOVED***
			fmt.Sscanf(line[4:], "nd:  %08x%08x", &c.bounds[1], &c.bounds[0])
			break
		***REMOVED***
	***REMOVED***
	f.Close()

	return
***REMOVED***

func (c *capsV3) Apply(kind CapType) (err error) ***REMOVED***
	if kind&BOUNDS == BOUNDS ***REMOVED***
		var data [2]capData
		err = capget(&c.hdr, &data[0])
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if (1<<uint(CAP_SETPCAP))&data[0].effective != 0 ***REMOVED***
			for i := Cap(0); i <= CAP_LAST_CAP; i++ ***REMOVED***
				if c.Get(BOUNDING, i) ***REMOVED***
					continue
				***REMOVED***
				err = prctl(syscall.PR_CAPBSET_DROP, uintptr(i), 0, 0, 0)
				if err != nil ***REMOVED***
					// Ignore EINVAL since the capability may not be supported in this system.
					if errno, ok := err.(syscall.Errno); ok && errno == syscall.EINVAL ***REMOVED***
						err = nil
						continue
					***REMOVED***
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if kind&CAPS == CAPS ***REMOVED***
		return capset(&c.hdr, &c.data[0])
	***REMOVED***

	return
***REMOVED***

func newFile(path string) (c Capabilities, err error) ***REMOVED***
	c = &capsFile***REMOVED***path: path***REMOVED***
	err = c.Load()
	if err != nil ***REMOVED***
		c = nil
	***REMOVED***
	return
***REMOVED***

type capsFile struct ***REMOVED***
	path string
	data vfscapData
***REMOVED***

func (c *capsFile) Get(which CapType, what Cap) bool ***REMOVED***
	var i uint
	if what > 31 ***REMOVED***
		if c.data.version == 1 ***REMOVED***
			return false
		***REMOVED***
		i = uint(what) >> 5
		what %= 32
	***REMOVED***

	switch which ***REMOVED***
	case EFFECTIVE:
		return (1<<uint(what))&c.data.effective[i] != 0
	case PERMITTED:
		return (1<<uint(what))&c.data.data[i].permitted != 0
	case INHERITABLE:
		return (1<<uint(what))&c.data.data[i].inheritable != 0
	***REMOVED***

	return false
***REMOVED***

func (c *capsFile) getData(which CapType, dest []uint32) ***REMOVED***
	switch which ***REMOVED***
	case EFFECTIVE:
		dest[0] = c.data.effective[0]
		dest[1] = c.data.effective[1]
	case PERMITTED:
		dest[0] = c.data.data[0].permitted
		dest[1] = c.data.data[1].permitted
	case INHERITABLE:
		dest[0] = c.data.data[0].inheritable
		dest[1] = c.data.data[1].inheritable
	***REMOVED***
***REMOVED***

func (c *capsFile) Empty(which CapType) bool ***REMOVED***
	var data [2]uint32
	c.getData(which, data[:])
	return data[0] == 0 && data[1] == 0
***REMOVED***

func (c *capsFile) Full(which CapType) bool ***REMOVED***
	var data [2]uint32
	c.getData(which, data[:])
	if c.data.version == 0 ***REMOVED***
		return (data[0] & 0x7fffffff) == 0x7fffffff
	***REMOVED***
	if (data[0] & 0xffffffff) != 0xffffffff ***REMOVED***
		return false
	***REMOVED***
	return (data[1] & capUpperMask) == capUpperMask
***REMOVED***

func (c *capsFile) Set(which CapType, caps ...Cap) ***REMOVED***
	for _, what := range caps ***REMOVED***
		var i uint
		if what > 31 ***REMOVED***
			if c.data.version == 1 ***REMOVED***
				continue
			***REMOVED***
			i = uint(what) >> 5
			what %= 32
		***REMOVED***

		if which&EFFECTIVE != 0 ***REMOVED***
			c.data.effective[i] |= 1 << uint(what)
		***REMOVED***
		if which&PERMITTED != 0 ***REMOVED***
			c.data.data[i].permitted |= 1 << uint(what)
		***REMOVED***
		if which&INHERITABLE != 0 ***REMOVED***
			c.data.data[i].inheritable |= 1 << uint(what)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *capsFile) Unset(which CapType, caps ...Cap) ***REMOVED***
	for _, what := range caps ***REMOVED***
		var i uint
		if what > 31 ***REMOVED***
			if c.data.version == 1 ***REMOVED***
				continue
			***REMOVED***
			i = uint(what) >> 5
			what %= 32
		***REMOVED***

		if which&EFFECTIVE != 0 ***REMOVED***
			c.data.effective[i] &= ^(1 << uint(what))
		***REMOVED***
		if which&PERMITTED != 0 ***REMOVED***
			c.data.data[i].permitted &= ^(1 << uint(what))
		***REMOVED***
		if which&INHERITABLE != 0 ***REMOVED***
			c.data.data[i].inheritable &= ^(1 << uint(what))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *capsFile) Fill(kind CapType) ***REMOVED***
	if kind&CAPS == CAPS ***REMOVED***
		c.data.effective[0] = 0xffffffff
		c.data.data[0].permitted = 0xffffffff
		c.data.data[0].inheritable = 0
		if c.data.version == 2 ***REMOVED***
			c.data.effective[1] = 0xffffffff
			c.data.data[1].permitted = 0xffffffff
			c.data.data[1].inheritable = 0
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *capsFile) Clear(kind CapType) ***REMOVED***
	if kind&CAPS == CAPS ***REMOVED***
		c.data.effective[0] = 0
		c.data.data[0].permitted = 0
		c.data.data[0].inheritable = 0
		if c.data.version == 2 ***REMOVED***
			c.data.effective[1] = 0
			c.data.data[1].permitted = 0
			c.data.data[1].inheritable = 0
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *capsFile) StringCap(which CapType) (ret string) ***REMOVED***
	return mkStringCap(c, which)
***REMOVED***

func (c *capsFile) String() (ret string) ***REMOVED***
	return mkString(c, INHERITABLE)
***REMOVED***

func (c *capsFile) Load() (err error) ***REMOVED***
	return getVfsCap(c.path, &c.data)
***REMOVED***

func (c *capsFile) Apply(kind CapType) (err error) ***REMOVED***
	if kind&CAPS == CAPS ***REMOVED***
		return setVfsCap(c.path, &c.data)
	***REMOVED***
	return
***REMOVED***
