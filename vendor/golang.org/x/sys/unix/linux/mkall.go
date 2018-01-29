// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// linux/mkall.go - Generates all Linux zsysnum, zsyscall, zerror, and ztype
// files for all 11 linux architectures supported by the go compiler. See
// README.md for more information about the build system.

// To run it you must have a git checkout of the Linux kernel and glibc. Once
// the appropriate sources are ready, the program is run as:
//     go run linux/mkall.go <linux_dir> <glibc_dir>

// +build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

// These will be paths to the appropriate source directories.
var LinuxDir string
var GlibcDir string

const TempDir = "/tmp"
const IncludeDir = TempDir + "/include" // To hold our C headers
const BuildDir = TempDir + "/build"     // To hold intermediate build files

const GOOS = "linux"       // Only for Linux targets
const BuildArch = "amd64"  // Must be built on this architecture
const MinKernel = "2.6.23" // https://golang.org/doc/install#requirements

type target struct ***REMOVED***
	GoArch     string // Architecture name according to Go
	LinuxArch  string // Architecture name according to the Linux Kernel
	GNUArch    string // Architecture name according to GNU tools (https://wiki.debian.org/Multiarch/Tuples)
	BigEndian  bool   // Default Little Endian
	SignedChar bool   // Is -fsigned-char needed (default no)
	Bits       int
***REMOVED***

// List of the 11 Linux targets supported by the go compiler. sparc64 is not
// currently supported, though a port is in progress.
var targets = []target***REMOVED***
	***REMOVED***
		GoArch:    "386",
		LinuxArch: "x86",
		GNUArch:   "i686-linux-gnu", // Note "i686" not "i386"
		Bits:      32,
	***REMOVED***,
	***REMOVED***
		GoArch:    "amd64",
		LinuxArch: "x86",
		GNUArch:   "x86_64-linux-gnu",
		Bits:      64,
	***REMOVED***,
	***REMOVED***
		GoArch:     "arm64",
		LinuxArch:  "arm64",
		GNUArch:    "aarch64-linux-gnu",
		SignedChar: true,
		Bits:       64,
	***REMOVED***,
	***REMOVED***
		GoArch:    "arm",
		LinuxArch: "arm",
		GNUArch:   "arm-linux-gnueabi",
		Bits:      32,
	***REMOVED***,
	***REMOVED***
		GoArch:    "mips",
		LinuxArch: "mips",
		GNUArch:   "mips-linux-gnu",
		BigEndian: true,
		Bits:      32,
	***REMOVED***,
	***REMOVED***
		GoArch:    "mipsle",
		LinuxArch: "mips",
		GNUArch:   "mipsel-linux-gnu",
		Bits:      32,
	***REMOVED***,
	***REMOVED***
		GoArch:    "mips64",
		LinuxArch: "mips",
		GNUArch:   "mips64-linux-gnuabi64",
		BigEndian: true,
		Bits:      64,
	***REMOVED***,
	***REMOVED***
		GoArch:    "mips64le",
		LinuxArch: "mips",
		GNUArch:   "mips64el-linux-gnuabi64",
		Bits:      64,
	***REMOVED***,
	***REMOVED***
		GoArch:    "ppc64",
		LinuxArch: "powerpc",
		GNUArch:   "powerpc64-linux-gnu",
		BigEndian: true,
		Bits:      64,
	***REMOVED***,
	***REMOVED***
		GoArch:    "ppc64le",
		LinuxArch: "powerpc",
		GNUArch:   "powerpc64le-linux-gnu",
		Bits:      64,
	***REMOVED***,
	***REMOVED***
		GoArch:     "s390x",
		LinuxArch:  "s390",
		GNUArch:    "s390x-linux-gnu",
		BigEndian:  true,
		SignedChar: true,
		Bits:       64,
	***REMOVED***,
	// ***REMOVED***
	// 	GoArch:    "sparc64",
	// 	LinuxArch: "sparc",
	// 	GNUArch:   "sparc64-linux-gnu",
	// 	BigEndian: true,
	// 	Bits:      64,
	// ***REMOVED***,
***REMOVED***

// ptracePairs is a list of pairs of targets that can, in some cases,
// run each other's binaries.
var ptracePairs = []struct***REMOVED*** a1, a2 string ***REMOVED******REMOVED***
	***REMOVED***"386", "amd64"***REMOVED***,
	***REMOVED***"arm", "arm64"***REMOVED***,
	***REMOVED***"mips", "mips64"***REMOVED***,
	***REMOVED***"mipsle", "mips64le"***REMOVED***,
***REMOVED***

func main() ***REMOVED***
	if runtime.GOOS != GOOS || runtime.GOARCH != BuildArch ***REMOVED***
		fmt.Printf("Build system has GOOS_GOARCH = %s_%s, need %s_%s\n",
			runtime.GOOS, runtime.GOARCH, GOOS, BuildArch)
		return
	***REMOVED***

	// Check that we are using the new build system if we should
	if os.Getenv("GOLANG_SYS_BUILD") != "docker" ***REMOVED***
		fmt.Println("In the new build system, mkall.go should not be called directly.")
		fmt.Println("See README.md")
		return
	***REMOVED***

	// Parse the command line options
	if len(os.Args) != 3 ***REMOVED***
		fmt.Println("USAGE: go run linux/mkall.go <linux_dir> <glibc_dir>")
		return
	***REMOVED***
	LinuxDir = os.Args[1]
	GlibcDir = os.Args[2]

	for _, t := range targets ***REMOVED***
		fmt.Printf("----- GENERATING: %s -----\n", t.GoArch)
		if err := t.generateFiles(); err != nil ***REMOVED***
			fmt.Printf("%v\n***** FAILURE:    %s *****\n\n", err, t.GoArch)
		***REMOVED*** else ***REMOVED***
			fmt.Printf("----- SUCCESS:    %s -----\n\n", t.GoArch)
		***REMOVED***
	***REMOVED***

	fmt.Printf("----- GENERATING ptrace pairs -----\n")
	ok := true
	for _, p := range ptracePairs ***REMOVED***
		if err := generatePtracePair(p.a1, p.a2); err != nil ***REMOVED***
			fmt.Printf("%v\n***** FAILURE: %s/%s *****\n\n", err, p.a1, p.a2)
			ok = false
		***REMOVED***
	***REMOVED***
	if ok ***REMOVED***
		fmt.Printf("----- SUCCESS ptrace pairs    -----\n\n")
	***REMOVED***
***REMOVED***

// Makes an exec.Cmd with Stderr attached to os.Stderr
func makeCommand(name string, args ...string) *exec.Cmd ***REMOVED***
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	return cmd
***REMOVED***

// Runs the command, pipes output to a formatter, pipes that to an output file.
func (t *target) commandFormatOutput(formatter string, outputFile string,
	name string, args ...string) (err error) ***REMOVED***
	mainCmd := makeCommand(name, args...)

	fmtCmd := makeCommand(formatter)
	if formatter == "mkpost" ***REMOVED***
		fmtCmd = makeCommand("go", "run", "mkpost.go")
		// Set GOARCH_TARGET so mkpost knows what GOARCH is..
		fmtCmd.Env = append(os.Environ(), "GOARCH_TARGET="+t.GoArch)
		// Set GOARCH to host arch for mkpost, so it can run natively.
		for i, s := range fmtCmd.Env ***REMOVED***
			if strings.HasPrefix(s, "GOARCH=") ***REMOVED***
				fmtCmd.Env[i] = "GOARCH=" + BuildArch
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// mainCmd | fmtCmd > outputFile
	if fmtCmd.Stdin, err = mainCmd.StdoutPipe(); err != nil ***REMOVED***
		return
	***REMOVED***
	if fmtCmd.Stdout, err = os.Create(outputFile); err != nil ***REMOVED***
		return
	***REMOVED***

	// Make sure the formatter eventually closes
	if err = fmtCmd.Start(); err != nil ***REMOVED***
		return
	***REMOVED***
	defer func() ***REMOVED***
		fmtErr := fmtCmd.Wait()
		if err == nil ***REMOVED***
			err = fmtErr
		***REMOVED***
	***REMOVED***()

	return mainCmd.Run()
***REMOVED***

// Generates all the files for a Linux target
func (t *target) generateFiles() error ***REMOVED***
	// Setup environment variables
	os.Setenv("GOOS", GOOS)
	os.Setenv("GOARCH", t.GoArch)

	// Get appropriate compiler and emulator (unless on x86)
	if t.LinuxArch != "x86" ***REMOVED***
		// Check/Setup cross compiler
		compiler := t.GNUArch + "-gcc"
		if _, err := exec.LookPath(compiler); err != nil ***REMOVED***
			return err
		***REMOVED***
		os.Setenv("CC", compiler)

		// Check/Setup emulator (usually first component of GNUArch)
		qemuArchName := t.GNUArch[:strings.Index(t.GNUArch, "-")]
		if t.LinuxArch == "powerpc" ***REMOVED***
			qemuArchName = t.GoArch
		***REMOVED***
		os.Setenv("GORUN", "qemu-"+qemuArchName)
	***REMOVED*** else ***REMOVED***
		os.Setenv("CC", "gcc")
	***REMOVED***

	// Make the include directory and fill it with headers
	if err := os.MkdirAll(IncludeDir, os.ModePerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.RemoveAll(IncludeDir)
	if err := t.makeHeaders(); err != nil ***REMOVED***
		return fmt.Errorf("could not make header files: %v", err)
	***REMOVED***
	fmt.Println("header files generated")

	// Make each of the four files
	if err := t.makeZSysnumFile(); err != nil ***REMOVED***
		return fmt.Errorf("could not make zsysnum file: %v", err)
	***REMOVED***
	fmt.Println("zsysnum file generated")

	if err := t.makeZSyscallFile(); err != nil ***REMOVED***
		return fmt.Errorf("could not make zsyscall file: %v", err)
	***REMOVED***
	fmt.Println("zsyscall file generated")

	if err := t.makeZTypesFile(); err != nil ***REMOVED***
		return fmt.Errorf("could not make ztypes file: %v", err)
	***REMOVED***
	fmt.Println("ztypes file generated")

	if err := t.makeZErrorsFile(); err != nil ***REMOVED***
		return fmt.Errorf("could not make zerrors file: %v", err)
	***REMOVED***
	fmt.Println("zerrors file generated")

	return nil
***REMOVED***

// Create the Linux and glibc headers in the include directory.
func (t *target) makeHeaders() error ***REMOVED***
	// Make the Linux headers we need for this architecture
	linuxMake := makeCommand("make", "headers_install", "ARCH="+t.LinuxArch, "INSTALL_HDR_PATH="+TempDir)
	linuxMake.Dir = LinuxDir
	if err := linuxMake.Run(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// A Temporary build directory for glibc
	if err := os.MkdirAll(BuildDir, os.ModePerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.RemoveAll(BuildDir)

	// Make the glibc headers we need for this architecture
	confScript := filepath.Join(GlibcDir, "configure")
	glibcConf := makeCommand(confScript, "--prefix="+TempDir, "--host="+t.GNUArch, "--enable-kernel="+MinKernel)
	glibcConf.Dir = BuildDir
	if err := glibcConf.Run(); err != nil ***REMOVED***
		return err
	***REMOVED***
	glibcMake := makeCommand("make", "install-headers")
	glibcMake.Dir = BuildDir
	if err := glibcMake.Run(); err != nil ***REMOVED***
		return err
	***REMOVED***
	// We only need an empty stubs file
	stubsFile := filepath.Join(IncludeDir, "gnu/stubs.h")
	if file, err := os.Create(stubsFile); err != nil ***REMOVED***
		return err
	***REMOVED*** else ***REMOVED***
		file.Close()
	***REMOVED***

	return nil
***REMOVED***

// makes the zsysnum_linux_$GOARCH.go file
func (t *target) makeZSysnumFile() error ***REMOVED***
	zsysnumFile := fmt.Sprintf("zsysnum_linux_%s.go", t.GoArch)
	unistdFile := filepath.Join(IncludeDir, "asm/unistd.h")

	args := append(t.cFlags(), unistdFile)
	return t.commandFormatOutput("gofmt", zsysnumFile, "linux/mksysnum.pl", args...)
***REMOVED***

// makes the zsyscall_linux_$GOARCH.go file
func (t *target) makeZSyscallFile() error ***REMOVED***
	zsyscallFile := fmt.Sprintf("zsyscall_linux_%s.go", t.GoArch)
	// Find the correct architecture syscall file (might end with x.go)
	archSyscallFile := fmt.Sprintf("syscall_linux_%s.go", t.GoArch)
	if _, err := os.Stat(archSyscallFile); os.IsNotExist(err) ***REMOVED***
		shortArch := strings.TrimSuffix(t.GoArch, "le")
		archSyscallFile = fmt.Sprintf("syscall_linux_%sx.go", shortArch)
	***REMOVED***

	args := append(t.mksyscallFlags(), "-tags", "linux,"+t.GoArch,
		"syscall_linux.go", archSyscallFile)
	return t.commandFormatOutput("gofmt", zsyscallFile, "./mksyscall.pl", args...)
***REMOVED***

// makes the zerrors_linux_$GOARCH.go file
func (t *target) makeZErrorsFile() error ***REMOVED***
	zerrorsFile := fmt.Sprintf("zerrors_linux_%s.go", t.GoArch)

	return t.commandFormatOutput("gofmt", zerrorsFile, "./mkerrors.sh", t.cFlags()...)
***REMOVED***

// makes the ztypes_linux_$GOARCH.go file
func (t *target) makeZTypesFile() error ***REMOVED***
	ztypesFile := fmt.Sprintf("ztypes_linux_%s.go", t.GoArch)

	args := []string***REMOVED***"tool", "cgo", "-godefs", "--"***REMOVED***
	args = append(args, t.cFlags()...)
	args = append(args, "linux/types.go")
	return t.commandFormatOutput("mkpost", ztypesFile, "go", args...)
***REMOVED***

// Flags that should be given to gcc and cgo for this target
func (t *target) cFlags() []string ***REMOVED***
	// Compile statically to avoid cross-architecture dynamic linking.
	flags := []string***REMOVED***"-Wall", "-Werror", "-static", "-I" + IncludeDir***REMOVED***

	// Architecture-specific flags
	if t.SignedChar ***REMOVED***
		flags = append(flags, "-fsigned-char")
	***REMOVED***
	if t.LinuxArch == "x86" ***REMOVED***
		flags = append(flags, fmt.Sprintf("-m%d", t.Bits))
	***REMOVED***

	return flags
***REMOVED***

// Flags that should be given to mksyscall for this target
func (t *target) mksyscallFlags() (flags []string) ***REMOVED***
	if t.Bits == 32 ***REMOVED***
		if t.BigEndian ***REMOVED***
			flags = append(flags, "-b32")
		***REMOVED*** else ***REMOVED***
			flags = append(flags, "-l32")
		***REMOVED***
	***REMOVED***

	// This flag menas a 64-bit value should use (even, odd)-pair.
	if t.GoArch == "arm" || (t.LinuxArch == "mips" && t.Bits == 32) ***REMOVED***
		flags = append(flags, "-arm")
	***REMOVED***
	return
***REMOVED***

// generatePtracePair takes a pair of GOARCH values that can run each
// other's binaries, such as 386 and amd64. It extracts the PtraceRegs
// type for each one. It writes a new file defining the types
// PtraceRegsArch1 and PtraceRegsArch2 and the corresponding functions
// Ptrace***REMOVED***Get,Set***REMOVED***Regs***REMOVED***arch1,arch2***REMOVED***. This permits debugging the other
// binary on a native system.
func generatePtracePair(arch1, arch2 string) error ***REMOVED***
	def1, err := ptraceDef(arch1)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	def2, err := ptraceDef(arch2)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f, err := os.Create(fmt.Sprintf("zptrace%s_linux.go", arch1))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	buf := bufio.NewWriter(f)
	fmt.Fprintf(buf, "// Code generated by linux/mkall.go generatePtracePair(%s, %s). DO NOT EDIT.\n", arch1, arch2)
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "// +build linux\n")
	fmt.Fprintf(buf, "// +build %s %s\n", arch1, arch2)
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "package unix\n")
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "%s\n", `import "unsafe"`)
	fmt.Fprintf(buf, "\n")
	writeOnePtrace(buf, arch1, def1)
	fmt.Fprintf(buf, "\n")
	writeOnePtrace(buf, arch2, def2)
	if err := buf.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := f.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// ptraceDef returns the definition of PtraceRegs for arch.
func ptraceDef(arch string) (string, error) ***REMOVED***
	filename := fmt.Sprintf("ztypes_linux_%s.go", arch)
	data, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("reading %s: %v", filename, err)
	***REMOVED***
	start := bytes.Index(data, []byte("type PtraceRegs struct"))
	if start < 0 ***REMOVED***
		return "", fmt.Errorf("%s: no definition of PtraceRegs", filename)
	***REMOVED***
	data = data[start:]
	end := bytes.Index(data, []byte("\n***REMOVED***\n"))
	if end < 0 ***REMOVED***
		return "", fmt.Errorf("%s: can't find end of PtraceRegs definition", filename)
	***REMOVED***
	return string(data[:end+2]), nil
***REMOVED***

// writeOnePtrace writes out the ptrace definitions for arch.
func writeOnePtrace(w io.Writer, arch, def string) ***REMOVED***
	uarch := string(unicode.ToUpper(rune(arch[0]))) + arch[1:]
	fmt.Fprintf(w, "// PtraceRegs%s is the registers used by %s binaries.\n", uarch, arch)
	fmt.Fprintf(w, "%s\n", strings.Replace(def, "PtraceRegs", "PtraceRegs"+uarch, 1))
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "// PtraceGetRegs%s fetches the registers used by %s binaries.\n", uarch, arch)
	fmt.Fprintf(w, "func PtraceGetRegs%s(pid int, regsout *PtraceRegs%s) error ***REMOVED***\n", uarch, uarch)
	fmt.Fprintf(w, "\treturn ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))\n")
	fmt.Fprintf(w, "***REMOVED***\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "// PtraceSetRegs%s sets the registers used by %s binaries.\n", uarch, arch)
	fmt.Fprintf(w, "func PtraceSetRegs%s(pid int, regs *PtraceRegs%s) error ***REMOVED***\n", uarch, uarch)
	fmt.Fprintf(w, "\treturn ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))\n")
	fmt.Fprintf(w, "***REMOVED***\n")
***REMOVED***
