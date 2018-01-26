// +build linux

package seccomp

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/opencontainers/runtime-spec/specs-go"
	libseccomp "github.com/seccomp/libseccomp-golang"
)

//go:generate go run -tags 'seccomp' generate.go

// GetDefaultProfile returns the default seccomp profile.
func GetDefaultProfile(rs *specs.Spec) (*specs.LinuxSeccomp, error) ***REMOVED***
	return setupSeccomp(DefaultProfile(), rs)
***REMOVED***

// LoadProfile takes a json string and decodes the seccomp profile.
func LoadProfile(body string, rs *specs.Spec) (*specs.LinuxSeccomp, error) ***REMOVED***
	var config types.Seccomp
	if err := json.Unmarshal([]byte(body), &config); err != nil ***REMOVED***
		return nil, fmt.Errorf("Decoding seccomp profile failed: %v", err)
	***REMOVED***
	return setupSeccomp(&config, rs)
***REMOVED***

var nativeToSeccomp = map[string]types.Arch***REMOVED***
	"amd64":       types.ArchX86_64,
	"arm64":       types.ArchAARCH64,
	"mips64":      types.ArchMIPS64,
	"mips64n32":   types.ArchMIPS64N32,
	"mipsel64":    types.ArchMIPSEL64,
	"mipsel64n32": types.ArchMIPSEL64N32,
	"s390x":       types.ArchS390X,
***REMOVED***

// inSlice tests whether a string is contained in a slice of strings or not.
// Comparison is case sensitive
func inSlice(slice []string, s string) bool ***REMOVED***
	for _, ss := range slice ***REMOVED***
		if s == ss ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func setupSeccomp(config *types.Seccomp, rs *specs.Spec) (*specs.LinuxSeccomp, error) ***REMOVED***
	if config == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	// No default action specified, no syscalls listed, assume seccomp disabled
	if config.DefaultAction == "" && len(config.Syscalls) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	newConfig := &specs.LinuxSeccomp***REMOVED******REMOVED***

	var arch string
	var native, err = libseccomp.GetNativeArch()
	if err == nil ***REMOVED***
		arch = native.String()
	***REMOVED***

	if len(config.Architectures) != 0 && len(config.ArchMap) != 0 ***REMOVED***
		return nil, errors.New("'architectures' and 'archMap' were specified in the seccomp profile, use either 'architectures' or 'archMap'")
	***REMOVED***

	// if config.Architectures == 0 then libseccomp will figure out the architecture to use
	if len(config.Architectures) != 0 ***REMOVED***
		for _, a := range config.Architectures ***REMOVED***
			newConfig.Architectures = append(newConfig.Architectures, specs.Arch(a))
		***REMOVED***
	***REMOVED***

	if len(config.ArchMap) != 0 ***REMOVED***
		for _, a := range config.ArchMap ***REMOVED***
			seccompArch, ok := nativeToSeccomp[arch]
			if ok ***REMOVED***
				if a.Arch == seccompArch ***REMOVED***
					newConfig.Architectures = append(newConfig.Architectures, specs.Arch(a.Arch))
					for _, sa := range a.SubArches ***REMOVED***
						newConfig.Architectures = append(newConfig.Architectures, specs.Arch(sa))
					***REMOVED***
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	newConfig.DefaultAction = specs.LinuxSeccompAction(config.DefaultAction)

Loop:
	// Loop through all syscall blocks and convert them to libcontainer format after filtering them
	for _, call := range config.Syscalls ***REMOVED***
		if len(call.Excludes.Arches) > 0 ***REMOVED***
			if inSlice(call.Excludes.Arches, arch) ***REMOVED***
				continue Loop
			***REMOVED***
		***REMOVED***
		if len(call.Excludes.Caps) > 0 ***REMOVED***
			for _, c := range call.Excludes.Caps ***REMOVED***
				if inSlice(rs.Process.Capabilities.Effective, c) ***REMOVED***
					continue Loop
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if len(call.Includes.Arches) > 0 ***REMOVED***
			if !inSlice(call.Includes.Arches, arch) ***REMOVED***
				continue Loop
			***REMOVED***
		***REMOVED***
		if len(call.Includes.Caps) > 0 ***REMOVED***
			for _, c := range call.Includes.Caps ***REMOVED***
				if !inSlice(rs.Process.Capabilities.Effective, c) ***REMOVED***
					continue Loop
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if call.Name != "" && len(call.Names) != 0 ***REMOVED***
			return nil, errors.New("'name' and 'names' were specified in the seccomp profile, use either 'name' or 'names'")
		***REMOVED***

		if call.Name != "" ***REMOVED***
			newConfig.Syscalls = append(newConfig.Syscalls, createSpecsSyscall(call.Name, call.Action, call.Args))
		***REMOVED***

		for _, n := range call.Names ***REMOVED***
			newConfig.Syscalls = append(newConfig.Syscalls, createSpecsSyscall(n, call.Action, call.Args))
		***REMOVED***
	***REMOVED***

	return newConfig, nil
***REMOVED***

func createSpecsSyscall(name string, action types.Action, args []*types.Arg) specs.LinuxSyscall ***REMOVED***
	newCall := specs.LinuxSyscall***REMOVED***
		Names:  []string***REMOVED***name***REMOVED***,
		Action: specs.LinuxSeccompAction(action),
	***REMOVED***

	// Loop through all the arguments of the syscall and convert them
	for _, arg := range args ***REMOVED***
		newArg := specs.LinuxSeccompArg***REMOVED***
			Index:    arg.Index,
			Value:    arg.Value,
			ValueTwo: arg.ValueTwo,
			Op:       specs.LinuxSeccompOperator(arg.Op),
		***REMOVED***

		newCall.Args = append(newCall.Args, newArg)
	***REMOVED***
	return newCall
***REMOVED***
