package daemon

import (
	"fmt"
	"strings"

	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/names"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/docker/pkg/stringid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	validContainerNameChars   = names.RestrictedNameChars
	validContainerNamePattern = names.RestrictedNamePattern
)

func (daemon *Daemon) registerName(container *container.Container) error ***REMOVED***
	if daemon.Exists(container.ID) ***REMOVED***
		return fmt.Errorf("Container is already loaded")
	***REMOVED***
	if err := validateID(container.ID); err != nil ***REMOVED***
		return err
	***REMOVED***
	if container.Name == "" ***REMOVED***
		name, err := daemon.generateNewName(container.ID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		container.Name = name
	***REMOVED***
	return daemon.containersReplica.ReserveName(container.Name, container.ID)
***REMOVED***

func (daemon *Daemon) generateIDAndName(name string) (string, string, error) ***REMOVED***
	var (
		err error
		id  = stringid.GenerateNonCryptoID()
	)

	if name == "" ***REMOVED***
		if name, err = daemon.generateNewName(id); err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
		return id, name, nil
	***REMOVED***

	if name, err = daemon.reserveName(id, name); err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	return id, name, nil
***REMOVED***

func (daemon *Daemon) reserveName(id, name string) (string, error) ***REMOVED***
	if !validContainerNamePattern.MatchString(strings.TrimPrefix(name, "/")) ***REMOVED***
		return "", errdefs.InvalidParameter(errors.Errorf("Invalid container name (%s), only %s are allowed", name, validContainerNameChars))
	***REMOVED***
	if name[0] != '/' ***REMOVED***
		name = "/" + name
	***REMOVED***

	if err := daemon.containersReplica.ReserveName(name, id); err != nil ***REMOVED***
		if err == container.ErrNameReserved ***REMOVED***
			id, err := daemon.containersReplica.Snapshot().GetID(name)
			if err != nil ***REMOVED***
				logrus.Errorf("got unexpected error while looking up reserved name: %v", err)
				return "", err
			***REMOVED***
			return "", nameConflictError***REMOVED***id: id, name: name***REMOVED***
		***REMOVED***
		return "", errors.Wrapf(err, "error reserving name: %q", name)
	***REMOVED***
	return name, nil
***REMOVED***

func (daemon *Daemon) releaseName(name string) ***REMOVED***
	daemon.containersReplica.ReleaseName(name)
***REMOVED***

func (daemon *Daemon) generateNewName(id string) (string, error) ***REMOVED***
	var name string
	for i := 0; i < 6; i++ ***REMOVED***
		name = namesgenerator.GetRandomName(i)
		if name[0] != '/' ***REMOVED***
			name = "/" + name
		***REMOVED***

		if err := daemon.containersReplica.ReserveName(name, id); err != nil ***REMOVED***
			if err == container.ErrNameReserved ***REMOVED***
				continue
			***REMOVED***
			return "", err
		***REMOVED***
		return name, nil
	***REMOVED***

	name = "/" + stringid.TruncateID(id)
	if err := daemon.containersReplica.ReserveName(name, id); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return name, nil
***REMOVED***

func validateID(id string) error ***REMOVED***
	if id == "" ***REMOVED***
		return fmt.Errorf("Invalid empty id")
	***REMOVED***
	return nil
***REMOVED***
