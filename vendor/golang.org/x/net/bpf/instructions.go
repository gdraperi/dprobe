// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf

import "fmt"

// An Instruction is one instruction executed by the BPF virtual
// machine.
type Instruction interface ***REMOVED***
	// Assemble assembles the Instruction into a RawInstruction.
	Assemble() (RawInstruction, error)
***REMOVED***

// A RawInstruction is a raw BPF virtual machine instruction.
type RawInstruction struct ***REMOVED***
	// Operation to execute.
	Op uint16
	// For conditional jump instructions, the number of instructions
	// to skip if the condition is true/false.
	Jt uint8
	Jf uint8
	// Constant parameter. The meaning depends on the Op.
	K uint32
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (ri RawInstruction) Assemble() (RawInstruction, error) ***REMOVED*** return ri, nil ***REMOVED***

// Disassemble parses ri into an Instruction and returns it. If ri is
// not recognized by this package, ri itself is returned.
func (ri RawInstruction) Disassemble() Instruction ***REMOVED***
	switch ri.Op & opMaskCls ***REMOVED***
	case opClsLoadA, opClsLoadX:
		reg := Register(ri.Op & opMaskLoadDest)
		sz := 0
		switch ri.Op & opMaskLoadWidth ***REMOVED***
		case opLoadWidth4:
			sz = 4
		case opLoadWidth2:
			sz = 2
		case opLoadWidth1:
			sz = 1
		default:
			return ri
		***REMOVED***
		switch ri.Op & opMaskLoadMode ***REMOVED***
		case opAddrModeImmediate:
			if sz != 4 ***REMOVED***
				return ri
			***REMOVED***
			return LoadConstant***REMOVED***Dst: reg, Val: ri.K***REMOVED***
		case opAddrModeScratch:
			if sz != 4 || ri.K > 15 ***REMOVED***
				return ri
			***REMOVED***
			return LoadScratch***REMOVED***Dst: reg, N: int(ri.K)***REMOVED***
		case opAddrModeAbsolute:
			if ri.K > extOffset+0xffffffff ***REMOVED***
				return LoadExtension***REMOVED***Num: Extension(-extOffset + ri.K)***REMOVED***
			***REMOVED***
			return LoadAbsolute***REMOVED***Size: sz, Off: ri.K***REMOVED***
		case opAddrModeIndirect:
			return LoadIndirect***REMOVED***Size: sz, Off: ri.K***REMOVED***
		case opAddrModePacketLen:
			if sz != 4 ***REMOVED***
				return ri
			***REMOVED***
			return LoadExtension***REMOVED***Num: ExtLen***REMOVED***
		case opAddrModeMemShift:
			return LoadMemShift***REMOVED***Off: ri.K***REMOVED***
		default:
			return ri
		***REMOVED***

	case opClsStoreA:
		if ri.Op != opClsStoreA || ri.K > 15 ***REMOVED***
			return ri
		***REMOVED***
		return StoreScratch***REMOVED***Src: RegA, N: int(ri.K)***REMOVED***

	case opClsStoreX:
		if ri.Op != opClsStoreX || ri.K > 15 ***REMOVED***
			return ri
		***REMOVED***
		return StoreScratch***REMOVED***Src: RegX, N: int(ri.K)***REMOVED***

	case opClsALU:
		switch op := ALUOp(ri.Op & opMaskOperator); op ***REMOVED***
		case ALUOpAdd, ALUOpSub, ALUOpMul, ALUOpDiv, ALUOpOr, ALUOpAnd, ALUOpShiftLeft, ALUOpShiftRight, ALUOpMod, ALUOpXor:
			if ri.Op&opMaskOperandSrc != 0 ***REMOVED***
				return ALUOpX***REMOVED***Op: op***REMOVED***
			***REMOVED***
			return ALUOpConstant***REMOVED***Op: op, Val: ri.K***REMOVED***
		case aluOpNeg:
			return NegateA***REMOVED******REMOVED***
		default:
			return ri
		***REMOVED***

	case opClsJump:
		if ri.Op&opMaskJumpConst != opClsJump ***REMOVED***
			return ri
		***REMOVED***
		switch ri.Op & opMaskJumpCond ***REMOVED***
		case opJumpAlways:
			return Jump***REMOVED***Skip: ri.K***REMOVED***
		case opJumpEqual:
			if ri.Jt == 0 ***REMOVED***
				return JumpIf***REMOVED***
					Cond:      JumpNotEqual,
					Val:       ri.K,
					SkipTrue:  ri.Jf,
					SkipFalse: 0,
				***REMOVED***
			***REMOVED***
			return JumpIf***REMOVED***
				Cond:      JumpEqual,
				Val:       ri.K,
				SkipTrue:  ri.Jt,
				SkipFalse: ri.Jf,
			***REMOVED***
		case opJumpGT:
			if ri.Jt == 0 ***REMOVED***
				return JumpIf***REMOVED***
					Cond:      JumpLessOrEqual,
					Val:       ri.K,
					SkipTrue:  ri.Jf,
					SkipFalse: 0,
				***REMOVED***
			***REMOVED***
			return JumpIf***REMOVED***
				Cond:      JumpGreaterThan,
				Val:       ri.K,
				SkipTrue:  ri.Jt,
				SkipFalse: ri.Jf,
			***REMOVED***
		case opJumpGE:
			if ri.Jt == 0 ***REMOVED***
				return JumpIf***REMOVED***
					Cond:      JumpLessThan,
					Val:       ri.K,
					SkipTrue:  ri.Jf,
					SkipFalse: 0,
				***REMOVED***
			***REMOVED***
			return JumpIf***REMOVED***
				Cond:      JumpGreaterOrEqual,
				Val:       ri.K,
				SkipTrue:  ri.Jt,
				SkipFalse: ri.Jf,
			***REMOVED***
		case opJumpSet:
			return JumpIf***REMOVED***
				Cond:      JumpBitsSet,
				Val:       ri.K,
				SkipTrue:  ri.Jt,
				SkipFalse: ri.Jf,
			***REMOVED***
		default:
			return ri
		***REMOVED***

	case opClsReturn:
		switch ri.Op ***REMOVED***
		case opClsReturn | opRetSrcA:
			return RetA***REMOVED******REMOVED***
		case opClsReturn | opRetSrcConstant:
			return RetConstant***REMOVED***Val: ri.K***REMOVED***
		default:
			return ri
		***REMOVED***

	case opClsMisc:
		switch ri.Op ***REMOVED***
		case opClsMisc | opMiscTAX:
			return TAX***REMOVED******REMOVED***
		case opClsMisc | opMiscTXA:
			return TXA***REMOVED******REMOVED***
		default:
			return ri
		***REMOVED***

	default:
		panic("unreachable") // switch is exhaustive on the bit pattern
	***REMOVED***
***REMOVED***

// LoadConstant loads Val into register Dst.
type LoadConstant struct ***REMOVED***
	Dst Register
	Val uint32
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a LoadConstant) Assemble() (RawInstruction, error) ***REMOVED***
	return assembleLoad(a.Dst, 4, opAddrModeImmediate, a.Val)
***REMOVED***

// String returns the the instruction in assembler notation.
func (a LoadConstant) String() string ***REMOVED***
	switch a.Dst ***REMOVED***
	case RegA:
		return fmt.Sprintf("ld #%d", a.Val)
	case RegX:
		return fmt.Sprintf("ldx #%d", a.Val)
	default:
		return fmt.Sprintf("unknown instruction: %#v", a)
	***REMOVED***
***REMOVED***

// LoadScratch loads scratch[N] into register Dst.
type LoadScratch struct ***REMOVED***
	Dst Register
	N   int // 0-15
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a LoadScratch) Assemble() (RawInstruction, error) ***REMOVED***
	if a.N < 0 || a.N > 15 ***REMOVED***
		return RawInstruction***REMOVED******REMOVED***, fmt.Errorf("invalid scratch slot %d", a.N)
	***REMOVED***
	return assembleLoad(a.Dst, 4, opAddrModeScratch, uint32(a.N))
***REMOVED***

// String returns the the instruction in assembler notation.
func (a LoadScratch) String() string ***REMOVED***
	switch a.Dst ***REMOVED***
	case RegA:
		return fmt.Sprintf("ld M[%d]", a.N)
	case RegX:
		return fmt.Sprintf("ldx M[%d]", a.N)
	default:
		return fmt.Sprintf("unknown instruction: %#v", a)
	***REMOVED***
***REMOVED***

// LoadAbsolute loads packet[Off:Off+Size] as an integer value into
// register A.
type LoadAbsolute struct ***REMOVED***
	Off  uint32
	Size int // 1, 2 or 4
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a LoadAbsolute) Assemble() (RawInstruction, error) ***REMOVED***
	return assembleLoad(RegA, a.Size, opAddrModeAbsolute, a.Off)
***REMOVED***

// String returns the the instruction in assembler notation.
func (a LoadAbsolute) String() string ***REMOVED***
	switch a.Size ***REMOVED***
	case 1: // byte
		return fmt.Sprintf("ldb [%d]", a.Off)
	case 2: // half word
		return fmt.Sprintf("ldh [%d]", a.Off)
	case 4: // word
		if a.Off > extOffset+0xffffffff ***REMOVED***
			return LoadExtension***REMOVED***Num: Extension(a.Off + 0x1000)***REMOVED***.String()
		***REMOVED***
		return fmt.Sprintf("ld [%d]", a.Off)
	default:
		return fmt.Sprintf("unknown instruction: %#v", a)
	***REMOVED***
***REMOVED***

// LoadIndirect loads packet[X+Off:X+Off+Size] as an integer value
// into register A.
type LoadIndirect struct ***REMOVED***
	Off  uint32
	Size int // 1, 2 or 4
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a LoadIndirect) Assemble() (RawInstruction, error) ***REMOVED***
	return assembleLoad(RegA, a.Size, opAddrModeIndirect, a.Off)
***REMOVED***

// String returns the the instruction in assembler notation.
func (a LoadIndirect) String() string ***REMOVED***
	switch a.Size ***REMOVED***
	case 1: // byte
		return fmt.Sprintf("ldb [x + %d]", a.Off)
	case 2: // half word
		return fmt.Sprintf("ldh [x + %d]", a.Off)
	case 4: // word
		return fmt.Sprintf("ld [x + %d]", a.Off)
	default:
		return fmt.Sprintf("unknown instruction: %#v", a)
	***REMOVED***
***REMOVED***

// LoadMemShift multiplies the first 4 bits of the byte at packet[Off]
// by 4 and stores the result in register X.
//
// This instruction is mainly useful to load into X the length of an
// IPv4 packet header in a single instruction, rather than have to do
// the arithmetic on the header's first byte by hand.
type LoadMemShift struct ***REMOVED***
	Off uint32
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a LoadMemShift) Assemble() (RawInstruction, error) ***REMOVED***
	return assembleLoad(RegX, 1, opAddrModeMemShift, a.Off)
***REMOVED***

// String returns the the instruction in assembler notation.
func (a LoadMemShift) String() string ***REMOVED***
	return fmt.Sprintf("ldx 4*([%d]&0xf)", a.Off)
***REMOVED***

// LoadExtension invokes a linux-specific extension and stores the
// result in register A.
type LoadExtension struct ***REMOVED***
	Num Extension
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a LoadExtension) Assemble() (RawInstruction, error) ***REMOVED***
	if a.Num == ExtLen ***REMOVED***
		return assembleLoad(RegA, 4, opAddrModePacketLen, 0)
	***REMOVED***
	return assembleLoad(RegA, 4, opAddrModeAbsolute, uint32(extOffset+a.Num))
***REMOVED***

// String returns the the instruction in assembler notation.
func (a LoadExtension) String() string ***REMOVED***
	switch a.Num ***REMOVED***
	case ExtLen:
		return "ld #len"
	case ExtProto:
		return "ld #proto"
	case ExtType:
		return "ld #type"
	case ExtPayloadOffset:
		return "ld #poff"
	case ExtInterfaceIndex:
		return "ld #ifidx"
	case ExtNetlinkAttr:
		return "ld #nla"
	case ExtNetlinkAttrNested:
		return "ld #nlan"
	case ExtMark:
		return "ld #mark"
	case ExtQueue:
		return "ld #queue"
	case ExtLinkLayerType:
		return "ld #hatype"
	case ExtRXHash:
		return "ld #rxhash"
	case ExtCPUID:
		return "ld #cpu"
	case ExtVLANTag:
		return "ld #vlan_tci"
	case ExtVLANTagPresent:
		return "ld #vlan_avail"
	case ExtVLANProto:
		return "ld #vlan_tpid"
	case ExtRand:
		return "ld #rand"
	default:
		return fmt.Sprintf("unknown instruction: %#v", a)
	***REMOVED***
***REMOVED***

// StoreScratch stores register Src into scratch[N].
type StoreScratch struct ***REMOVED***
	Src Register
	N   int // 0-15
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a StoreScratch) Assemble() (RawInstruction, error) ***REMOVED***
	if a.N < 0 || a.N > 15 ***REMOVED***
		return RawInstruction***REMOVED******REMOVED***, fmt.Errorf("invalid scratch slot %d", a.N)
	***REMOVED***
	var op uint16
	switch a.Src ***REMOVED***
	case RegA:
		op = opClsStoreA
	case RegX:
		op = opClsStoreX
	default:
		return RawInstruction***REMOVED******REMOVED***, fmt.Errorf("invalid source register %v", a.Src)
	***REMOVED***

	return RawInstruction***REMOVED***
		Op: op,
		K:  uint32(a.N),
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a StoreScratch) String() string ***REMOVED***
	switch a.Src ***REMOVED***
	case RegA:
		return fmt.Sprintf("st M[%d]", a.N)
	case RegX:
		return fmt.Sprintf("stx M[%d]", a.N)
	default:
		return fmt.Sprintf("unknown instruction: %#v", a)
	***REMOVED***
***REMOVED***

// ALUOpConstant executes A = A <Op> Val.
type ALUOpConstant struct ***REMOVED***
	Op  ALUOp
	Val uint32
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a ALUOpConstant) Assemble() (RawInstruction, error) ***REMOVED***
	return RawInstruction***REMOVED***
		Op: opClsALU | opALUSrcConstant | uint16(a.Op),
		K:  a.Val,
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a ALUOpConstant) String() string ***REMOVED***
	switch a.Op ***REMOVED***
	case ALUOpAdd:
		return fmt.Sprintf("add #%d", a.Val)
	case ALUOpSub:
		return fmt.Sprintf("sub #%d", a.Val)
	case ALUOpMul:
		return fmt.Sprintf("mul #%d", a.Val)
	case ALUOpDiv:
		return fmt.Sprintf("div #%d", a.Val)
	case ALUOpMod:
		return fmt.Sprintf("mod #%d", a.Val)
	case ALUOpAnd:
		return fmt.Sprintf("and #%d", a.Val)
	case ALUOpOr:
		return fmt.Sprintf("or #%d", a.Val)
	case ALUOpXor:
		return fmt.Sprintf("xor #%d", a.Val)
	case ALUOpShiftLeft:
		return fmt.Sprintf("lsh #%d", a.Val)
	case ALUOpShiftRight:
		return fmt.Sprintf("rsh #%d", a.Val)
	default:
		return fmt.Sprintf("unknown instruction: %#v", a)
	***REMOVED***
***REMOVED***

// ALUOpX executes A = A <Op> X
type ALUOpX struct ***REMOVED***
	Op ALUOp
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a ALUOpX) Assemble() (RawInstruction, error) ***REMOVED***
	return RawInstruction***REMOVED***
		Op: opClsALU | opALUSrcX | uint16(a.Op),
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a ALUOpX) String() string ***REMOVED***
	switch a.Op ***REMOVED***
	case ALUOpAdd:
		return "add x"
	case ALUOpSub:
		return "sub x"
	case ALUOpMul:
		return "mul x"
	case ALUOpDiv:
		return "div x"
	case ALUOpMod:
		return "mod x"
	case ALUOpAnd:
		return "and x"
	case ALUOpOr:
		return "or x"
	case ALUOpXor:
		return "xor x"
	case ALUOpShiftLeft:
		return "lsh x"
	case ALUOpShiftRight:
		return "rsh x"
	default:
		return fmt.Sprintf("unknown instruction: %#v", a)
	***REMOVED***
***REMOVED***

// NegateA executes A = -A.
type NegateA struct***REMOVED******REMOVED***

// Assemble implements the Instruction Assemble method.
func (a NegateA) Assemble() (RawInstruction, error) ***REMOVED***
	return RawInstruction***REMOVED***
		Op: opClsALU | uint16(aluOpNeg),
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a NegateA) String() string ***REMOVED***
	return fmt.Sprintf("neg")
***REMOVED***

// Jump skips the following Skip instructions in the program.
type Jump struct ***REMOVED***
	Skip uint32
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a Jump) Assemble() (RawInstruction, error) ***REMOVED***
	return RawInstruction***REMOVED***
		Op: opClsJump | opJumpAlways,
		K:  a.Skip,
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a Jump) String() string ***REMOVED***
	return fmt.Sprintf("ja %d", a.Skip)
***REMOVED***

// JumpIf skips the following Skip instructions in the program if A
// <Cond> Val is true.
type JumpIf struct ***REMOVED***
	Cond      JumpTest
	Val       uint32
	SkipTrue  uint8
	SkipFalse uint8
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a JumpIf) Assemble() (RawInstruction, error) ***REMOVED***
	var (
		cond uint16
		flip bool
	)
	switch a.Cond ***REMOVED***
	case JumpEqual:
		cond = opJumpEqual
	case JumpNotEqual:
		cond, flip = opJumpEqual, true
	case JumpGreaterThan:
		cond = opJumpGT
	case JumpLessThan:
		cond, flip = opJumpGE, true
	case JumpGreaterOrEqual:
		cond = opJumpGE
	case JumpLessOrEqual:
		cond, flip = opJumpGT, true
	case JumpBitsSet:
		cond = opJumpSet
	case JumpBitsNotSet:
		cond, flip = opJumpSet, true
	default:
		return RawInstruction***REMOVED******REMOVED***, fmt.Errorf("unknown JumpTest %v", a.Cond)
	***REMOVED***
	jt, jf := a.SkipTrue, a.SkipFalse
	if flip ***REMOVED***
		jt, jf = jf, jt
	***REMOVED***
	return RawInstruction***REMOVED***
		Op: opClsJump | cond,
		Jt: jt,
		Jf: jf,
		K:  a.Val,
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a JumpIf) String() string ***REMOVED***
	switch a.Cond ***REMOVED***
	// K == A
	case JumpEqual:
		return conditionalJump(a, "jeq", "jneq")
	// K != A
	case JumpNotEqual:
		return fmt.Sprintf("jneq #%d,%d", a.Val, a.SkipTrue)
	// K > A
	case JumpGreaterThan:
		return conditionalJump(a, "jgt", "jle")
	// K < A
	case JumpLessThan:
		return fmt.Sprintf("jlt #%d,%d", a.Val, a.SkipTrue)
	// K >= A
	case JumpGreaterOrEqual:
		return conditionalJump(a, "jge", "jlt")
	// K <= A
	case JumpLessOrEqual:
		return fmt.Sprintf("jle #%d,%d", a.Val, a.SkipTrue)
	// K & A != 0
	case JumpBitsSet:
		if a.SkipFalse > 0 ***REMOVED***
			return fmt.Sprintf("jset #%d,%d,%d", a.Val, a.SkipTrue, a.SkipFalse)
		***REMOVED***
		return fmt.Sprintf("jset #%d,%d", a.Val, a.SkipTrue)
	// K & A == 0, there is no assembler instruction for JumpBitNotSet, use JumpBitSet and invert skips
	case JumpBitsNotSet:
		return JumpIf***REMOVED***Cond: JumpBitsSet, SkipTrue: a.SkipFalse, SkipFalse: a.SkipTrue, Val: a.Val***REMOVED***.String()
	default:
		return fmt.Sprintf("unknown instruction: %#v", a)
	***REMOVED***
***REMOVED***

func conditionalJump(inst JumpIf, positiveJump, negativeJump string) string ***REMOVED***
	if inst.SkipTrue > 0 ***REMOVED***
		if inst.SkipFalse > 0 ***REMOVED***
			return fmt.Sprintf("%s #%d,%d,%d", positiveJump, inst.Val, inst.SkipTrue, inst.SkipFalse)
		***REMOVED***
		return fmt.Sprintf("%s #%d,%d", positiveJump, inst.Val, inst.SkipTrue)
	***REMOVED***
	return fmt.Sprintf("%s #%d,%d", negativeJump, inst.Val, inst.SkipFalse)
***REMOVED***

// RetA exits the BPF program, returning the value of register A.
type RetA struct***REMOVED******REMOVED***

// Assemble implements the Instruction Assemble method.
func (a RetA) Assemble() (RawInstruction, error) ***REMOVED***
	return RawInstruction***REMOVED***
		Op: opClsReturn | opRetSrcA,
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a RetA) String() string ***REMOVED***
	return fmt.Sprintf("ret a")
***REMOVED***

// RetConstant exits the BPF program, returning a constant value.
type RetConstant struct ***REMOVED***
	Val uint32
***REMOVED***

// Assemble implements the Instruction Assemble method.
func (a RetConstant) Assemble() (RawInstruction, error) ***REMOVED***
	return RawInstruction***REMOVED***
		Op: opClsReturn | opRetSrcConstant,
		K:  a.Val,
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a RetConstant) String() string ***REMOVED***
	return fmt.Sprintf("ret #%d", a.Val)
***REMOVED***

// TXA copies the value of register X to register A.
type TXA struct***REMOVED******REMOVED***

// Assemble implements the Instruction Assemble method.
func (a TXA) Assemble() (RawInstruction, error) ***REMOVED***
	return RawInstruction***REMOVED***
		Op: opClsMisc | opMiscTXA,
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a TXA) String() string ***REMOVED***
	return fmt.Sprintf("txa")
***REMOVED***

// TAX copies the value of register A to register X.
type TAX struct***REMOVED******REMOVED***

// Assemble implements the Instruction Assemble method.
func (a TAX) Assemble() (RawInstruction, error) ***REMOVED***
	return RawInstruction***REMOVED***
		Op: opClsMisc | opMiscTAX,
	***REMOVED***, nil
***REMOVED***

// String returns the the instruction in assembler notation.
func (a TAX) String() string ***REMOVED***
	return fmt.Sprintf("tax")
***REMOVED***

func assembleLoad(dst Register, loadSize int, mode uint16, k uint32) (RawInstruction, error) ***REMOVED***
	var (
		cls uint16
		sz  uint16
	)
	switch dst ***REMOVED***
	case RegA:
		cls = opClsLoadA
	case RegX:
		cls = opClsLoadX
	default:
		return RawInstruction***REMOVED******REMOVED***, fmt.Errorf("invalid target register %v", dst)
	***REMOVED***
	switch loadSize ***REMOVED***
	case 1:
		sz = opLoadWidth1
	case 2:
		sz = opLoadWidth2
	case 4:
		sz = opLoadWidth4
	default:
		return RawInstruction***REMOVED******REMOVED***, fmt.Errorf("invalid load byte length %d", sz)
	***REMOVED***
	return RawInstruction***REMOVED***
		Op: cls | sz | mode,
		K:  k,
	***REMOVED***, nil
***REMOVED***
