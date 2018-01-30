// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpf

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// This is a direct translation of the program in
// testdata/all_instructions.txt.
var allInstructions = []Instruction***REMOVED***
	LoadConstant***REMOVED***Dst: RegA, Val: 42***REMOVED***,
	LoadConstant***REMOVED***Dst: RegX, Val: 42***REMOVED***,

	LoadScratch***REMOVED***Dst: RegA, N: 3***REMOVED***,
	LoadScratch***REMOVED***Dst: RegX, N: 3***REMOVED***,

	LoadAbsolute***REMOVED***Off: 42, Size: 1***REMOVED***,
	LoadAbsolute***REMOVED***Off: 42, Size: 2***REMOVED***,
	LoadAbsolute***REMOVED***Off: 42, Size: 4***REMOVED***,

	LoadIndirect***REMOVED***Off: 42, Size: 1***REMOVED***,
	LoadIndirect***REMOVED***Off: 42, Size: 2***REMOVED***,
	LoadIndirect***REMOVED***Off: 42, Size: 4***REMOVED***,

	LoadMemShift***REMOVED***Off: 42***REMOVED***,

	LoadExtension***REMOVED***Num: ExtLen***REMOVED***,
	LoadExtension***REMOVED***Num: ExtProto***REMOVED***,
	LoadExtension***REMOVED***Num: ExtType***REMOVED***,
	LoadExtension***REMOVED***Num: ExtRand***REMOVED***,

	StoreScratch***REMOVED***Src: RegA, N: 3***REMOVED***,
	StoreScratch***REMOVED***Src: RegX, N: 3***REMOVED***,

	ALUOpConstant***REMOVED***Op: ALUOpAdd, Val: 42***REMOVED***,
	ALUOpConstant***REMOVED***Op: ALUOpSub, Val: 42***REMOVED***,
	ALUOpConstant***REMOVED***Op: ALUOpMul, Val: 42***REMOVED***,
	ALUOpConstant***REMOVED***Op: ALUOpDiv, Val: 42***REMOVED***,
	ALUOpConstant***REMOVED***Op: ALUOpOr, Val: 42***REMOVED***,
	ALUOpConstant***REMOVED***Op: ALUOpAnd, Val: 42***REMOVED***,
	ALUOpConstant***REMOVED***Op: ALUOpShiftLeft, Val: 42***REMOVED***,
	ALUOpConstant***REMOVED***Op: ALUOpShiftRight, Val: 42***REMOVED***,
	ALUOpConstant***REMOVED***Op: ALUOpMod, Val: 42***REMOVED***,
	ALUOpConstant***REMOVED***Op: ALUOpXor, Val: 42***REMOVED***,

	ALUOpX***REMOVED***Op: ALUOpAdd***REMOVED***,
	ALUOpX***REMOVED***Op: ALUOpSub***REMOVED***,
	ALUOpX***REMOVED***Op: ALUOpMul***REMOVED***,
	ALUOpX***REMOVED***Op: ALUOpDiv***REMOVED***,
	ALUOpX***REMOVED***Op: ALUOpOr***REMOVED***,
	ALUOpX***REMOVED***Op: ALUOpAnd***REMOVED***,
	ALUOpX***REMOVED***Op: ALUOpShiftLeft***REMOVED***,
	ALUOpX***REMOVED***Op: ALUOpShiftRight***REMOVED***,
	ALUOpX***REMOVED***Op: ALUOpMod***REMOVED***,
	ALUOpX***REMOVED***Op: ALUOpXor***REMOVED***,

	NegateA***REMOVED******REMOVED***,

	Jump***REMOVED***Skip: 10***REMOVED***,
	JumpIf***REMOVED***Cond: JumpEqual, Val: 42, SkipTrue: 8, SkipFalse: 9***REMOVED***,
	JumpIf***REMOVED***Cond: JumpNotEqual, Val: 42, SkipTrue: 8***REMOVED***,
	JumpIf***REMOVED***Cond: JumpLessThan, Val: 42, SkipTrue: 7***REMOVED***,
	JumpIf***REMOVED***Cond: JumpLessOrEqual, Val: 42, SkipTrue: 6***REMOVED***,
	JumpIf***REMOVED***Cond: JumpGreaterThan, Val: 42, SkipTrue: 4, SkipFalse: 5***REMOVED***,
	JumpIf***REMOVED***Cond: JumpGreaterOrEqual, Val: 42, SkipTrue: 3, SkipFalse: 4***REMOVED***,
	JumpIf***REMOVED***Cond: JumpBitsSet, Val: 42, SkipTrue: 2, SkipFalse: 3***REMOVED***,

	TAX***REMOVED******REMOVED***,
	TXA***REMOVED******REMOVED***,

	RetA***REMOVED******REMOVED***,
	RetConstant***REMOVED***Val: 42***REMOVED***,
***REMOVED***
var allInstructionsExpected = "testdata/all_instructions.bpf"

// Check that we produce the same output as the canonical bpf_asm
// linux kernel tool.
func TestInterop(t *testing.T) ***REMOVED***
	out, err := Assemble(allInstructions)
	if err != nil ***REMOVED***
		t.Fatalf("assembly of allInstructions program failed: %s", err)
	***REMOVED***
	t.Logf("Assembled program is %d instructions long", len(out))

	bs, err := ioutil.ReadFile(allInstructionsExpected)
	if err != nil ***REMOVED***
		t.Fatalf("reading %s: %s", allInstructionsExpected, err)
	***REMOVED***
	// First statement is the number of statements, last statement is
	// empty. We just ignore both and rely on slice length.
	stmts := strings.Split(string(bs), ",")
	if len(stmts)-2 != len(out) ***REMOVED***
		t.Fatalf("test program lengths don't match: %s has %d, Go implementation has %d", allInstructionsExpected, len(stmts)-2, len(allInstructions))
	***REMOVED***

	for i, stmt := range stmts[1 : len(stmts)-2] ***REMOVED***
		nums := strings.Split(stmt, " ")
		if len(nums) != 4 ***REMOVED***
			t.Fatalf("malformed instruction %d in %s: %s", i+1, allInstructionsExpected, stmt)
		***REMOVED***

		actual := out[i]

		op, err := strconv.ParseUint(nums[0], 10, 16)
		if err != nil ***REMOVED***
			t.Fatalf("malformed opcode %s in instruction %d of %s", nums[0], i+1, allInstructionsExpected)
		***REMOVED***
		if actual.Op != uint16(op) ***REMOVED***
			t.Errorf("opcode mismatch on instruction %d (%#v): got 0x%02x, want 0x%02x", i+1, allInstructions[i], actual.Op, op)
		***REMOVED***

		jt, err := strconv.ParseUint(nums[1], 10, 8)
		if err != nil ***REMOVED***
			t.Fatalf("malformed jt offset %s in instruction %d of %s", nums[1], i+1, allInstructionsExpected)
		***REMOVED***
		if actual.Jt != uint8(jt) ***REMOVED***
			t.Errorf("jt mismatch on instruction %d (%#v): got %d, want %d", i+1, allInstructions[i], actual.Jt, jt)
		***REMOVED***

		jf, err := strconv.ParseUint(nums[2], 10, 8)
		if err != nil ***REMOVED***
			t.Fatalf("malformed jf offset %s in instruction %d of %s", nums[2], i+1, allInstructionsExpected)
		***REMOVED***
		if actual.Jf != uint8(jf) ***REMOVED***
			t.Errorf("jf mismatch on instruction %d (%#v): got %d, want %d", i+1, allInstructions[i], actual.Jf, jf)
		***REMOVED***

		k, err := strconv.ParseUint(nums[3], 10, 32)
		if err != nil ***REMOVED***
			t.Fatalf("malformed constant %s in instruction %d of %s", nums[3], i+1, allInstructionsExpected)
		***REMOVED***
		if actual.K != uint32(k) ***REMOVED***
			t.Errorf("constant mismatch on instruction %d (%#v): got %d, want %d", i+1, allInstructions[i], actual.K, k)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Check that assembly and disassembly match each other.
func TestAsmDisasm(t *testing.T) ***REMOVED***
	prog1, err := Assemble(allInstructions)
	if err != nil ***REMOVED***
		t.Fatalf("assembly of allInstructions program failed: %s", err)
	***REMOVED***
	t.Logf("Assembled program is %d instructions long", len(prog1))

	got, allDecoded := Disassemble(prog1)
	if !allDecoded ***REMOVED***
		t.Errorf("Disassemble(Assemble(allInstructions)) produced unrecognized instructions:")
		for i, inst := range got ***REMOVED***
			if r, ok := inst.(RawInstruction); ok ***REMOVED***
				t.Logf("  insn %d, %#v --> %#v", i+1, allInstructions[i], r)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if len(allInstructions) != len(got) ***REMOVED***
		t.Fatalf("disassembly changed program size: %d insns before, %d insns after", len(allInstructions), len(got))
	***REMOVED***
	if !reflect.DeepEqual(allInstructions, got) ***REMOVED***
		t.Errorf("program mutated by disassembly:")
		for i := range got ***REMOVED***
			if !reflect.DeepEqual(allInstructions[i], got[i]) ***REMOVED***
				t.Logf("  insn %d, s: %#v, p1: %#v, got: %#v", i+1, allInstructions[i], prog1[i], got[i])
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type InvalidInstruction struct***REMOVED******REMOVED***

func (a InvalidInstruction) Assemble() (RawInstruction, error) ***REMOVED***
	return RawInstruction***REMOVED******REMOVED***, fmt.Errorf("Invalid Instruction")
***REMOVED***

func (a InvalidInstruction) String() string ***REMOVED***
	return fmt.Sprintf("unknown instruction: %#v", a)
***REMOVED***

func TestString(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		instruction Instruction
		assembler   string
	***REMOVED******REMOVED***
		***REMOVED***
			instruction: LoadConstant***REMOVED***Dst: RegA, Val: 42***REMOVED***,
			assembler:   "ld #42",
		***REMOVED***,
		***REMOVED***
			instruction: LoadConstant***REMOVED***Dst: RegX, Val: 42***REMOVED***,
			assembler:   "ldx #42",
		***REMOVED***,
		***REMOVED***
			instruction: LoadConstant***REMOVED***Dst: 0xffff, Val: 42***REMOVED***,
			assembler:   "unknown instruction: bpf.LoadConstant***REMOVED***Dst:0xffff, Val:0x2a***REMOVED***",
		***REMOVED***,
		***REMOVED***
			instruction: LoadScratch***REMOVED***Dst: RegA, N: 3***REMOVED***,
			assembler:   "ld M[3]",
		***REMOVED***,
		***REMOVED***
			instruction: LoadScratch***REMOVED***Dst: RegX, N: 3***REMOVED***,
			assembler:   "ldx M[3]",
		***REMOVED***,
		***REMOVED***
			instruction: LoadScratch***REMOVED***Dst: 0xffff, N: 3***REMOVED***,
			assembler:   "unknown instruction: bpf.LoadScratch***REMOVED***Dst:0xffff, N:3***REMOVED***",
		***REMOVED***,
		***REMOVED***
			instruction: LoadAbsolute***REMOVED***Off: 42, Size: 1***REMOVED***,
			assembler:   "ldb [42]",
		***REMOVED***,
		***REMOVED***
			instruction: LoadAbsolute***REMOVED***Off: 42, Size: 2***REMOVED***,
			assembler:   "ldh [42]",
		***REMOVED***,
		***REMOVED***
			instruction: LoadAbsolute***REMOVED***Off: 42, Size: 4***REMOVED***,
			assembler:   "ld [42]",
		***REMOVED***,
		***REMOVED***
			instruction: LoadAbsolute***REMOVED***Off: 42, Size: -1***REMOVED***,
			assembler:   "unknown instruction: bpf.LoadAbsolute***REMOVED***Off:0x2a, Size:-1***REMOVED***",
		***REMOVED***,
		***REMOVED***
			instruction: LoadIndirect***REMOVED***Off: 42, Size: 1***REMOVED***,
			assembler:   "ldb [x + 42]",
		***REMOVED***,
		***REMOVED***
			instruction: LoadIndirect***REMOVED***Off: 42, Size: 2***REMOVED***,
			assembler:   "ldh [x + 42]",
		***REMOVED***,
		***REMOVED***
			instruction: LoadIndirect***REMOVED***Off: 42, Size: 4***REMOVED***,
			assembler:   "ld [x + 42]",
		***REMOVED***,
		***REMOVED***
			instruction: LoadIndirect***REMOVED***Off: 42, Size: -1***REMOVED***,
			assembler:   "unknown instruction: bpf.LoadIndirect***REMOVED***Off:0x2a, Size:-1***REMOVED***",
		***REMOVED***,
		***REMOVED***
			instruction: LoadMemShift***REMOVED***Off: 42***REMOVED***,
			assembler:   "ldx 4*([42]&0xf)",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtLen***REMOVED***,
			assembler:   "ld #len",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtProto***REMOVED***,
			assembler:   "ld #proto",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtType***REMOVED***,
			assembler:   "ld #type",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtPayloadOffset***REMOVED***,
			assembler:   "ld #poff",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtInterfaceIndex***REMOVED***,
			assembler:   "ld #ifidx",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtNetlinkAttr***REMOVED***,
			assembler:   "ld #nla",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtNetlinkAttrNested***REMOVED***,
			assembler:   "ld #nlan",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtMark***REMOVED***,
			assembler:   "ld #mark",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtQueue***REMOVED***,
			assembler:   "ld #queue",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtLinkLayerType***REMOVED***,
			assembler:   "ld #hatype",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtRXHash***REMOVED***,
			assembler:   "ld #rxhash",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtCPUID***REMOVED***,
			assembler:   "ld #cpu",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtVLANTag***REMOVED***,
			assembler:   "ld #vlan_tci",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtVLANTagPresent***REMOVED***,
			assembler:   "ld #vlan_avail",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtVLANProto***REMOVED***,
			assembler:   "ld #vlan_tpid",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: ExtRand***REMOVED***,
			assembler:   "ld #rand",
		***REMOVED***,
		***REMOVED***
			instruction: LoadAbsolute***REMOVED***Off: 0xfffff038, Size: 4***REMOVED***,
			assembler:   "ld #rand",
		***REMOVED***,
		***REMOVED***
			instruction: LoadExtension***REMOVED***Num: 0xfff***REMOVED***,
			assembler:   "unknown instruction: bpf.LoadExtension***REMOVED***Num:4095***REMOVED***",
		***REMOVED***,
		***REMOVED***
			instruction: StoreScratch***REMOVED***Src: RegA, N: 3***REMOVED***,
			assembler:   "st M[3]",
		***REMOVED***,
		***REMOVED***
			instruction: StoreScratch***REMOVED***Src: RegX, N: 3***REMOVED***,
			assembler:   "stx M[3]",
		***REMOVED***,
		***REMOVED***
			instruction: StoreScratch***REMOVED***Src: 0xffff, N: 3***REMOVED***,
			assembler:   "unknown instruction: bpf.StoreScratch***REMOVED***Src:0xffff, N:3***REMOVED***",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpAdd, Val: 42***REMOVED***,
			assembler:   "add #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpSub, Val: 42***REMOVED***,
			assembler:   "sub #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpMul, Val: 42***REMOVED***,
			assembler:   "mul #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpDiv, Val: 42***REMOVED***,
			assembler:   "div #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpOr, Val: 42***REMOVED***,
			assembler:   "or #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpAnd, Val: 42***REMOVED***,
			assembler:   "and #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpShiftLeft, Val: 42***REMOVED***,
			assembler:   "lsh #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpShiftRight, Val: 42***REMOVED***,
			assembler:   "rsh #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpMod, Val: 42***REMOVED***,
			assembler:   "mod #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: ALUOpXor, Val: 42***REMOVED***,
			assembler:   "xor #42",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpConstant***REMOVED***Op: 0xffff, Val: 42***REMOVED***,
			assembler:   "unknown instruction: bpf.ALUOpConstant***REMOVED***Op:0xffff, Val:0x2a***REMOVED***",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpAdd***REMOVED***,
			assembler:   "add x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpSub***REMOVED***,
			assembler:   "sub x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpMul***REMOVED***,
			assembler:   "mul x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpDiv***REMOVED***,
			assembler:   "div x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpOr***REMOVED***,
			assembler:   "or x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpAnd***REMOVED***,
			assembler:   "and x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpShiftLeft***REMOVED***,
			assembler:   "lsh x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpShiftRight***REMOVED***,
			assembler:   "rsh x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpMod***REMOVED***,
			assembler:   "mod x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: ALUOpXor***REMOVED***,
			assembler:   "xor x",
		***REMOVED***,
		***REMOVED***
			instruction: ALUOpX***REMOVED***Op: 0xffff***REMOVED***,
			assembler:   "unknown instruction: bpf.ALUOpX***REMOVED***Op:0xffff***REMOVED***",
		***REMOVED***,
		***REMOVED***
			instruction: NegateA***REMOVED******REMOVED***,
			assembler:   "neg",
		***REMOVED***,
		***REMOVED***
			instruction: Jump***REMOVED***Skip: 10***REMOVED***,
			assembler:   "ja 10",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpEqual, Val: 42, SkipTrue: 8, SkipFalse: 9***REMOVED***,
			assembler:   "jeq #42,8,9",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpEqual, Val: 42, SkipTrue: 8***REMOVED***,
			assembler:   "jeq #42,8",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpEqual, Val: 42, SkipFalse: 8***REMOVED***,
			assembler:   "jneq #42,8",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpNotEqual, Val: 42, SkipTrue: 8***REMOVED***,
			assembler:   "jneq #42,8",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpLessThan, Val: 42, SkipTrue: 7***REMOVED***,
			assembler:   "jlt #42,7",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpLessOrEqual, Val: 42, SkipTrue: 6***REMOVED***,
			assembler:   "jle #42,6",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpGreaterThan, Val: 42, SkipTrue: 4, SkipFalse: 5***REMOVED***,
			assembler:   "jgt #42,4,5",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpGreaterThan, Val: 42, SkipTrue: 4***REMOVED***,
			assembler:   "jgt #42,4",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpGreaterOrEqual, Val: 42, SkipTrue: 3, SkipFalse: 4***REMOVED***,
			assembler:   "jge #42,3,4",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpGreaterOrEqual, Val: 42, SkipTrue: 3***REMOVED***,
			assembler:   "jge #42,3",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpBitsSet, Val: 42, SkipTrue: 2, SkipFalse: 3***REMOVED***,
			assembler:   "jset #42,2,3",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpBitsSet, Val: 42, SkipTrue: 2***REMOVED***,
			assembler:   "jset #42,2",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpBitsNotSet, Val: 42, SkipTrue: 2, SkipFalse: 3***REMOVED***,
			assembler:   "jset #42,3,2",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: JumpBitsNotSet, Val: 42, SkipTrue: 2***REMOVED***,
			assembler:   "jset #42,0,2",
		***REMOVED***,
		***REMOVED***
			instruction: JumpIf***REMOVED***Cond: 0xffff, Val: 42, SkipTrue: 1, SkipFalse: 2***REMOVED***,
			assembler:   "unknown instruction: bpf.JumpIf***REMOVED***Cond:0xffff, Val:0x2a, SkipTrue:0x1, SkipFalse:0x2***REMOVED***",
		***REMOVED***,
		***REMOVED***
			instruction: TAX***REMOVED******REMOVED***,
			assembler:   "tax",
		***REMOVED***,
		***REMOVED***
			instruction: TXA***REMOVED******REMOVED***,
			assembler:   "txa",
		***REMOVED***,
		***REMOVED***
			instruction: RetA***REMOVED******REMOVED***,
			assembler:   "ret a",
		***REMOVED***,
		***REMOVED***
			instruction: RetConstant***REMOVED***Val: 42***REMOVED***,
			assembler:   "ret #42",
		***REMOVED***,
		// Invalid instruction
		***REMOVED***
			instruction: InvalidInstruction***REMOVED******REMOVED***,
			assembler:   "unknown instruction: bpf.InvalidInstruction***REMOVED******REMOVED***",
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		if input, ok := testCase.instruction.(fmt.Stringer); ok ***REMOVED***
			got := input.String()
			if got != testCase.assembler ***REMOVED***
				t.Errorf("String did not return expected assembler notation, expected: %s, got: %s", testCase.assembler, got)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			t.Errorf("Instruction %#v is not a fmt.Stringer", testCase.instruction)
		***REMOVED***
	***REMOVED***
***REMOVED***
