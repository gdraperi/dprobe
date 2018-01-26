package convert

import (
	"errors"
	"fmt"
	"strings"

	container "github.com/docker/docker/api/types/container"
	mounttypes "github.com/docker/docker/api/types/mount"
	types "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/sirupsen/logrus"
)

func containerSpecFromGRPC(c *swarmapi.ContainerSpec) *types.ContainerSpec ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	containerSpec := &types.ContainerSpec***REMOVED***
		Image:      c.Image,
		Labels:     c.Labels,
		Command:    c.Command,
		Args:       c.Args,
		Hostname:   c.Hostname,
		Env:        c.Env,
		Dir:        c.Dir,
		User:       c.User,
		Groups:     c.Groups,
		StopSignal: c.StopSignal,
		TTY:        c.TTY,
		OpenStdin:  c.OpenStdin,
		ReadOnly:   c.ReadOnly,
		Hosts:      c.Hosts,
		Secrets:    secretReferencesFromGRPC(c.Secrets),
		Configs:    configReferencesFromGRPC(c.Configs),
		Isolation:  IsolationFromGRPC(c.Isolation),
	***REMOVED***

	if c.DNSConfig != nil ***REMOVED***
		containerSpec.DNSConfig = &types.DNSConfig***REMOVED***
			Nameservers: c.DNSConfig.Nameservers,
			Search:      c.DNSConfig.Search,
			Options:     c.DNSConfig.Options,
		***REMOVED***
	***REMOVED***

	// Privileges
	if c.Privileges != nil ***REMOVED***
		containerSpec.Privileges = &types.Privileges***REMOVED******REMOVED***

		if c.Privileges.CredentialSpec != nil ***REMOVED***
			containerSpec.Privileges.CredentialSpec = &types.CredentialSpec***REMOVED******REMOVED***
			switch c.Privileges.CredentialSpec.Source.(type) ***REMOVED***
			case *swarmapi.Privileges_CredentialSpec_File:
				containerSpec.Privileges.CredentialSpec.File = c.Privileges.CredentialSpec.GetFile()
			case *swarmapi.Privileges_CredentialSpec_Registry:
				containerSpec.Privileges.CredentialSpec.Registry = c.Privileges.CredentialSpec.GetRegistry()
			***REMOVED***
		***REMOVED***

		if c.Privileges.SELinuxContext != nil ***REMOVED***
			containerSpec.Privileges.SELinuxContext = &types.SELinuxContext***REMOVED***
				Disable: c.Privileges.SELinuxContext.Disable,
				User:    c.Privileges.SELinuxContext.User,
				Type:    c.Privileges.SELinuxContext.Type,
				Role:    c.Privileges.SELinuxContext.Role,
				Level:   c.Privileges.SELinuxContext.Level,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Mounts
	for _, m := range c.Mounts ***REMOVED***
		mount := mounttypes.Mount***REMOVED***
			Target:   m.Target,
			Source:   m.Source,
			Type:     mounttypes.Type(strings.ToLower(swarmapi.Mount_MountType_name[int32(m.Type)])),
			ReadOnly: m.ReadOnly,
		***REMOVED***

		if m.BindOptions != nil ***REMOVED***
			mount.BindOptions = &mounttypes.BindOptions***REMOVED***
				Propagation: mounttypes.Propagation(strings.ToLower(swarmapi.Mount_BindOptions_MountPropagation_name[int32(m.BindOptions.Propagation)])),
			***REMOVED***
		***REMOVED***

		if m.VolumeOptions != nil ***REMOVED***
			mount.VolumeOptions = &mounttypes.VolumeOptions***REMOVED***
				NoCopy: m.VolumeOptions.NoCopy,
				Labels: m.VolumeOptions.Labels,
			***REMOVED***
			if m.VolumeOptions.DriverConfig != nil ***REMOVED***
				mount.VolumeOptions.DriverConfig = &mounttypes.Driver***REMOVED***
					Name:    m.VolumeOptions.DriverConfig.Name,
					Options: m.VolumeOptions.DriverConfig.Options,
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if m.TmpfsOptions != nil ***REMOVED***
			mount.TmpfsOptions = &mounttypes.TmpfsOptions***REMOVED***
				SizeBytes: m.TmpfsOptions.SizeBytes,
				Mode:      m.TmpfsOptions.Mode,
			***REMOVED***
		***REMOVED***
		containerSpec.Mounts = append(containerSpec.Mounts, mount)
	***REMOVED***

	if c.StopGracePeriod != nil ***REMOVED***
		grace, _ := gogotypes.DurationFromProto(c.StopGracePeriod)
		containerSpec.StopGracePeriod = &grace
	***REMOVED***

	if c.Healthcheck != nil ***REMOVED***
		containerSpec.Healthcheck = healthConfigFromGRPC(c.Healthcheck)
	***REMOVED***

	return containerSpec
***REMOVED***

func secretReferencesToGRPC(sr []*types.SecretReference) []*swarmapi.SecretReference ***REMOVED***
	refs := make([]*swarmapi.SecretReference, 0, len(sr))
	for _, s := range sr ***REMOVED***
		ref := &swarmapi.SecretReference***REMOVED***
			SecretID:   s.SecretID,
			SecretName: s.SecretName,
		***REMOVED***
		if s.File != nil ***REMOVED***
			ref.Target = &swarmapi.SecretReference_File***REMOVED***
				File: &swarmapi.FileTarget***REMOVED***
					Name: s.File.Name,
					UID:  s.File.UID,
					GID:  s.File.GID,
					Mode: s.File.Mode,
				***REMOVED***,
			***REMOVED***
		***REMOVED***

		refs = append(refs, ref)
	***REMOVED***

	return refs
***REMOVED***

func secretReferencesFromGRPC(sr []*swarmapi.SecretReference) []*types.SecretReference ***REMOVED***
	refs := make([]*types.SecretReference, 0, len(sr))
	for _, s := range sr ***REMOVED***
		target := s.GetFile()
		if target == nil ***REMOVED***
			// not a file target
			logrus.Warnf("secret target not a file: secret=%s", s.SecretID)
			continue
		***REMOVED***
		refs = append(refs, &types.SecretReference***REMOVED***
			File: &types.SecretReferenceFileTarget***REMOVED***
				Name: target.Name,
				UID:  target.UID,
				GID:  target.GID,
				Mode: target.Mode,
			***REMOVED***,
			SecretID:   s.SecretID,
			SecretName: s.SecretName,
		***REMOVED***)
	***REMOVED***

	return refs
***REMOVED***

func configReferencesToGRPC(sr []*types.ConfigReference) []*swarmapi.ConfigReference ***REMOVED***
	refs := make([]*swarmapi.ConfigReference, 0, len(sr))
	for _, s := range sr ***REMOVED***
		ref := &swarmapi.ConfigReference***REMOVED***
			ConfigID:   s.ConfigID,
			ConfigName: s.ConfigName,
		***REMOVED***
		if s.File != nil ***REMOVED***
			ref.Target = &swarmapi.ConfigReference_File***REMOVED***
				File: &swarmapi.FileTarget***REMOVED***
					Name: s.File.Name,
					UID:  s.File.UID,
					GID:  s.File.GID,
					Mode: s.File.Mode,
				***REMOVED***,
			***REMOVED***
		***REMOVED***

		refs = append(refs, ref)
	***REMOVED***

	return refs
***REMOVED***

func configReferencesFromGRPC(sr []*swarmapi.ConfigReference) []*types.ConfigReference ***REMOVED***
	refs := make([]*types.ConfigReference, 0, len(sr))
	for _, s := range sr ***REMOVED***
		target := s.GetFile()
		if target == nil ***REMOVED***
			// not a file target
			logrus.Warnf("config target not a file: config=%s", s.ConfigID)
			continue
		***REMOVED***
		refs = append(refs, &types.ConfigReference***REMOVED***
			File: &types.ConfigReferenceFileTarget***REMOVED***
				Name: target.Name,
				UID:  target.UID,
				GID:  target.GID,
				Mode: target.Mode,
			***REMOVED***,
			ConfigID:   s.ConfigID,
			ConfigName: s.ConfigName,
		***REMOVED***)
	***REMOVED***

	return refs
***REMOVED***

func containerToGRPC(c *types.ContainerSpec) (*swarmapi.ContainerSpec, error) ***REMOVED***
	containerSpec := &swarmapi.ContainerSpec***REMOVED***
		Image:      c.Image,
		Labels:     c.Labels,
		Command:    c.Command,
		Args:       c.Args,
		Hostname:   c.Hostname,
		Env:        c.Env,
		Dir:        c.Dir,
		User:       c.User,
		Groups:     c.Groups,
		StopSignal: c.StopSignal,
		TTY:        c.TTY,
		OpenStdin:  c.OpenStdin,
		ReadOnly:   c.ReadOnly,
		Hosts:      c.Hosts,
		Secrets:    secretReferencesToGRPC(c.Secrets),
		Configs:    configReferencesToGRPC(c.Configs),
		Isolation:  isolationToGRPC(c.Isolation),
	***REMOVED***

	if c.DNSConfig != nil ***REMOVED***
		containerSpec.DNSConfig = &swarmapi.ContainerSpec_DNSConfig***REMOVED***
			Nameservers: c.DNSConfig.Nameservers,
			Search:      c.DNSConfig.Search,
			Options:     c.DNSConfig.Options,
		***REMOVED***
	***REMOVED***

	if c.StopGracePeriod != nil ***REMOVED***
		containerSpec.StopGracePeriod = gogotypes.DurationProto(*c.StopGracePeriod)
	***REMOVED***

	// Privileges
	if c.Privileges != nil ***REMOVED***
		containerSpec.Privileges = &swarmapi.Privileges***REMOVED******REMOVED***

		if c.Privileges.CredentialSpec != nil ***REMOVED***
			containerSpec.Privileges.CredentialSpec = &swarmapi.Privileges_CredentialSpec***REMOVED******REMOVED***

			if c.Privileges.CredentialSpec.File != "" && c.Privileges.CredentialSpec.Registry != "" ***REMOVED***
				return nil, errors.New("cannot specify both \"file\" and \"registry\" credential specs")
			***REMOVED***
			if c.Privileges.CredentialSpec.File != "" ***REMOVED***
				containerSpec.Privileges.CredentialSpec.Source = &swarmapi.Privileges_CredentialSpec_File***REMOVED***
					File: c.Privileges.CredentialSpec.File,
				***REMOVED***
			***REMOVED*** else if c.Privileges.CredentialSpec.Registry != "" ***REMOVED***
				containerSpec.Privileges.CredentialSpec.Source = &swarmapi.Privileges_CredentialSpec_Registry***REMOVED***
					Registry: c.Privileges.CredentialSpec.Registry,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				return nil, errors.New("must either provide \"file\" or \"registry\" for credential spec")
			***REMOVED***
		***REMOVED***

		if c.Privileges.SELinuxContext != nil ***REMOVED***
			containerSpec.Privileges.SELinuxContext = &swarmapi.Privileges_SELinuxContext***REMOVED***
				Disable: c.Privileges.SELinuxContext.Disable,
				User:    c.Privileges.SELinuxContext.User,
				Type:    c.Privileges.SELinuxContext.Type,
				Role:    c.Privileges.SELinuxContext.Role,
				Level:   c.Privileges.SELinuxContext.Level,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Mounts
	for _, m := range c.Mounts ***REMOVED***
		mount := swarmapi.Mount***REMOVED***
			Target:   m.Target,
			Source:   m.Source,
			ReadOnly: m.ReadOnly,
		***REMOVED***

		if mountType, ok := swarmapi.Mount_MountType_value[strings.ToUpper(string(m.Type))]; ok ***REMOVED***
			mount.Type = swarmapi.Mount_MountType(mountType)
		***REMOVED*** else if string(m.Type) != "" ***REMOVED***
			return nil, fmt.Errorf("invalid MountType: %q", m.Type)
		***REMOVED***

		if m.BindOptions != nil ***REMOVED***
			if mountPropagation, ok := swarmapi.Mount_BindOptions_MountPropagation_value[strings.ToUpper(string(m.BindOptions.Propagation))]; ok ***REMOVED***
				mount.BindOptions = &swarmapi.Mount_BindOptions***REMOVED***Propagation: swarmapi.Mount_BindOptions_MountPropagation(mountPropagation)***REMOVED***
			***REMOVED*** else if string(m.BindOptions.Propagation) != "" ***REMOVED***
				return nil, fmt.Errorf("invalid MountPropagation: %q", m.BindOptions.Propagation)
			***REMOVED***
		***REMOVED***

		if m.VolumeOptions != nil ***REMOVED***
			mount.VolumeOptions = &swarmapi.Mount_VolumeOptions***REMOVED***
				NoCopy: m.VolumeOptions.NoCopy,
				Labels: m.VolumeOptions.Labels,
			***REMOVED***
			if m.VolumeOptions.DriverConfig != nil ***REMOVED***
				mount.VolumeOptions.DriverConfig = &swarmapi.Driver***REMOVED***
					Name:    m.VolumeOptions.DriverConfig.Name,
					Options: m.VolumeOptions.DriverConfig.Options,
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if m.TmpfsOptions != nil ***REMOVED***
			mount.TmpfsOptions = &swarmapi.Mount_TmpfsOptions***REMOVED***
				SizeBytes: m.TmpfsOptions.SizeBytes,
				Mode:      m.TmpfsOptions.Mode,
			***REMOVED***
		***REMOVED***

		containerSpec.Mounts = append(containerSpec.Mounts, mount)
	***REMOVED***

	if c.Healthcheck != nil ***REMOVED***
		containerSpec.Healthcheck = healthConfigToGRPC(c.Healthcheck)
	***REMOVED***

	return containerSpec, nil
***REMOVED***

func healthConfigFromGRPC(h *swarmapi.HealthConfig) *container.HealthConfig ***REMOVED***
	interval, _ := gogotypes.DurationFromProto(h.Interval)
	timeout, _ := gogotypes.DurationFromProto(h.Timeout)
	startPeriod, _ := gogotypes.DurationFromProto(h.StartPeriod)
	return &container.HealthConfig***REMOVED***
		Test:        h.Test,
		Interval:    interval,
		Timeout:     timeout,
		Retries:     int(h.Retries),
		StartPeriod: startPeriod,
	***REMOVED***
***REMOVED***

func healthConfigToGRPC(h *container.HealthConfig) *swarmapi.HealthConfig ***REMOVED***
	return &swarmapi.HealthConfig***REMOVED***
		Test:        h.Test,
		Interval:    gogotypes.DurationProto(h.Interval),
		Timeout:     gogotypes.DurationProto(h.Timeout),
		Retries:     int32(h.Retries),
		StartPeriod: gogotypes.DurationProto(h.StartPeriod),
	***REMOVED***
***REMOVED***

// IsolationFromGRPC converts a swarm api container isolation to a moby isolation representation
func IsolationFromGRPC(i swarmapi.ContainerSpec_Isolation) container.Isolation ***REMOVED***
	switch i ***REMOVED***
	case swarmapi.ContainerIsolationHyperV:
		return container.IsolationHyperV
	case swarmapi.ContainerIsolationProcess:
		return container.IsolationProcess
	case swarmapi.ContainerIsolationDefault:
		return container.IsolationDefault
	***REMOVED***
	return container.IsolationEmpty
***REMOVED***

func isolationToGRPC(i container.Isolation) swarmapi.ContainerSpec_Isolation ***REMOVED***
	if i.IsHyperV() ***REMOVED***
		return swarmapi.ContainerIsolationHyperV
	***REMOVED***
	if i.IsProcess() ***REMOVED***
		return swarmapi.ContainerIsolationProcess
	***REMOVED***
	return swarmapi.ContainerIsolationDefault
***REMOVED***
