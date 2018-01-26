// +build !windows

package daemon

import (
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/plugins"
	metrics "github.com/docker/go-metrics"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func (daemon *Daemon) listenMetricsSock() (string, error) ***REMOVED***
	path := filepath.Join(daemon.configStore.ExecRoot, "metrics.sock")
	unix.Unlink(path)
	l, err := net.Listen("unix", path)
	if err != nil ***REMOVED***
		return "", errors.Wrap(err, "error setting up metrics plugin listener")
	***REMOVED***

	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.Handler())
	go func() ***REMOVED***
		http.Serve(l, mux)
	***REMOVED***()
	daemon.metricsPluginListener = l
	return path, nil
***REMOVED***

func registerMetricsPluginCallback(getter plugingetter.PluginGetter, sockPath string) ***REMOVED***
	getter.Handle(metricsPluginType, func(name string, client *plugins.Client) ***REMOVED***
		// Use lookup since nothing in the system can really reference it, no need
		// to protect against removal
		p, err := getter.Get(name, metricsPluginType, plugingetter.Lookup)
		if err != nil ***REMOVED***
			return
		***REMOVED***

		mp := metricsPlugin***REMOVED***p***REMOVED***
		sockBase := mp.sockBase()
		if err := os.MkdirAll(sockBase, 0755); err != nil ***REMOVED***
			logrus.WithError(err).WithField("name", name).WithField("path", sockBase).Error("error creating metrics plugin base path")
			return
		***REMOVED***

		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				os.RemoveAll(sockBase)
			***REMOVED***
		***REMOVED***()

		pluginSockPath := filepath.Join(sockBase, mp.sock())
		_, err = os.Stat(pluginSockPath)
		if err == nil ***REMOVED***
			mount.Unmount(pluginSockPath)
		***REMOVED*** else ***REMOVED***
			logrus.WithField("path", pluginSockPath).Debugf("creating plugin socket")
			f, err := os.OpenFile(pluginSockPath, os.O_CREATE, 0600)
			if err != nil ***REMOVED***
				return
			***REMOVED***
			f.Close()
		***REMOVED***

		if err := mount.Mount(sockPath, pluginSockPath, "none", "bind,ro"); err != nil ***REMOVED***
			logrus.WithError(err).WithField("name", name).Error("could not mount metrics socket to plugin")
			return
		***REMOVED***

		if err := pluginStartMetricsCollection(p); err != nil ***REMOVED***
			if err := mount.Unmount(pluginSockPath); err != nil ***REMOVED***
				if mounted, _ := mount.Mounted(pluginSockPath); mounted ***REMOVED***
					logrus.WithError(err).WithField("sock_path", pluginSockPath).Error("error unmounting metrics socket from plugin during cleanup")
				***REMOVED***
			***REMOVED***
			logrus.WithError(err).WithField("name", name).Error("error while initializing metrics plugin")
		***REMOVED***
	***REMOVED***)
***REMOVED***
