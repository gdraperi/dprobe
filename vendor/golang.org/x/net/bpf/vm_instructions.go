// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf

import (
	"encoding/binary"
	"fmt"
)

func aluOpConstant(ins ALUOpConstant, regA uint32) uint32 ***REMOVED***
	return aluOpCommon(ins.Op, regA, ins.Val)
***REMOVED***

func aluOpX(ins ALUOpX, regA uint32, regX uint32) (uint32, bool) ***REMOVED***
	// Guard against division or modulus by zero by terminating
	// the program, as the OS BPF VM does
	if regX == 0 ***REMOVED***
		switch ins.Op ***REMOVED***
		case ALUOpDiv, ALUOpMod:
			return 0, false
		***REMOVED***
	***REMOVED***

	return aluOpCommon(ins.Op, regA, regX), true
***REMOVED***

func aluOpCommon(op ALUOp, regA uint32, value uint32) uint32 ***REMOVED***
	switch op ***REMOVED***
	case ALUOpAdd:
		return regA + value
	case ALUOpSub:
		return regA - value
	case ALUOpMul:
		return regA * value
	case ALUOpDiv:
		// Division by zero not permitted by NewVM and aluOpX checks
		return regA / value
	case ALUOpOr:
		return regA | value
	case ALUOpAnd:
		return regA & value
	case ALUOpShiftLeft:
		return regA << value
	case ALUOpShiftRight:
		return regA >> value
	case ALUOpMod:
		// Modulus by zero not permitted by NewVM and aluOpX checks
		return regA % value
	case ALUOpXor:
		return regA ^ value
	default:
		return regA
	***REMOVED***
***REMOVED***

func jumpIf(ins JumpIf, value uint32) int ***REMOVED***
	var ok bool
	inV := uint32(ins.Val)

	switch ins.Cond ***REMOVED***
	case JumpEqual:
		ok = value == inV
	case JumpNotEqual:
		ok = value != inV
	case JumpGreaterThan:
		ok = value > inV
	case JumpLessThan:
		ok = value < inV
	case JumpGreaterOrEqual:
		ok = value >= inV
	case JumpLessOrEqual:
		ok = value <= inV
	case JumpBitsSet:
		ok = (value & inV) != 0
	case JumpBitsNotSet:
		ok = (value & inV) == 0
	***REMOVED***

	if ok ***REMOVED***
		return int(ins.SkipTrue)
	***REMOVED***

	return int(ins.SkipFalse)
***REMOVED***

func loadAbsolute(ins LoadAbsolute, in []byte) (uint32, bool) ***REMOVED***
	offset := int(ins.Off)
	size := int(ins.Size)

	return loadCommon(in, offset, size)
***REMOVED***

func loadConstant(ins LoadConstant, regA uint32, regX uint32) (uint32, uint32) ***REMOVED***
	switch ins.Dst ***REMOVED***
	case RegA:
		regA = ins.Val
	case RegX:
		regX = ins.Val
	***REMOVED***

	return regA, regX
***REMOVED***

func loadExtension(ins LoadExtension, in []byte) uint32 ***REMOVED***
	switch ins.Num ***REMOVED***
	case ExtLen:
		return uint32(len(in))
	default:
		panic(fmt.Sprintf("unimplemented extension: %d", ins.Num))
	***REMOVED***
***REMOVED***

func loadIndirect(ins LoadIndirect, in []byte, regX uint32) (uint32, bool) ***REMOVED***
	offset := int(ins.Off) + int(regX)
	size := int(ins.Size)

	return loadCommon(in, offset, size)
***REMOVED***

func loadMemShift(ins LoadMemShift, in []byte) (uint32, bool) ***REMOVED***
	offset := int(ins.Off)

	if !inBounds(len(in), offset, 0) ***REMOVED***
		return 0, false
	***REMOVED***

	// Mask off high 4 bits and multiply low 4 bits by 4
	return uint32(in[offset]&0x0f) * 4, true
***REMOVED***

func inBounds(inLen int, offset int, size int) bool ***REMOVED***
	return offset+size <= inLen
***REMOVED***

func loadCommon(in []byte, offset int, size int) (uint32, bool) ***REMOVED***
	if !inBounds(len(in), offset, size) ***REMOVED***
		return 0, false
	***REMOVED***

	switch size ***REMOVED***
	case 1:
		return uint32(in[offset]), true
	case 2:
		return uint32(binary.BigEndian.Uint16(in[offset : offset+size])), true
	case 4:
		return uint32(binary.BigEndian.Uint32(in[offset : offset+size])), true
	default:
		panic(fmt.Sprintf("invalid load size: %d", size))
	***REMOVED***
***REMOVED***

func loadScratch(ins LoadScratch, regScratch [16]uint32, regA uint32, regX uint32) (uint32, uint32) ***REMOVED***
	switch ins.Dst ***REMOVED***
	case RegA:
		regA = regScratch[ins.N]
	case RegX:
		regX = regScratch[ins.N]
	***REMOVED***

	return regA, regX
***REMOVED***

func storeScratch(ins StoreScratch, regScratch [16]uint32, regA uint32, regX uint32) [16]uint32 ***REMOVED***
	switch ins.Src ***REMOVED***
	case RegA:
		regScratch[ins.N] = regA
	case RegX:
		regScratch[ins.N] = regX
	***REMOVED***

	return regScratch
***REMOVED***
