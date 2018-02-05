// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package remote integrates the remote features of Viper.
package remote

import (
	"bytes"
	"io"
	"os"

	"github.com/spf13/viper"
	crypt "github.com/xordataexchange/crypt/config"
)

type remoteConfigProvider struct***REMOVED******REMOVED***

func (rc remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) ***REMOVED***
	cm, err := getConfigManager(rp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	b, err := cm.Get(rp.Path())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return bytes.NewReader(b), nil
***REMOVED***

func (rc remoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) ***REMOVED***
	cm, err := getConfigManager(rp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	resp, err := cm.Get(rp.Path())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return bytes.NewReader(resp), nil
***REMOVED***

func (rc remoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) ***REMOVED***
	cm, err := getConfigManager(rp)
	if err != nil ***REMOVED***
		return nil, nil
	***REMOVED***
	quit := make(chan bool)
	quitwc := make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)
	cryptoResponseCh := cm.Watch(rp.Path(), quit)
	// need this function to convert the Channel response form crypt.Response to viper.Response
	go func(cr <-chan *crypt.Response, vr chan<- *viper.RemoteResponse, quitwc <-chan bool, quit chan<- bool) ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-quitwc:
				quit <- true
				return
			case resp := <-cr:
				vr <- &viper.RemoteResponse***REMOVED***
					Error: resp.Error,
					Value: resp.Value,
				***REMOVED***

			***REMOVED***

		***REMOVED***
	***REMOVED***(cryptoResponseCh, viperResponsCh, quitwc, quit)

	return viperResponsCh, quitwc
***REMOVED***

func getConfigManager(rp viper.RemoteProvider) (crypt.ConfigManager, error) ***REMOVED***
	var cm crypt.ConfigManager
	var err error

	if rp.SecretKeyring() != "" ***REMOVED***
		kr, err := os.Open(rp.SecretKeyring())
		defer kr.Close()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if rp.Provider() == "etcd" ***REMOVED***
			cm, err = crypt.NewEtcdConfigManager([]string***REMOVED***rp.Endpoint()***REMOVED***, kr)
		***REMOVED*** else ***REMOVED***
			cm, err = crypt.NewConsulConfigManager([]string***REMOVED***rp.Endpoint()***REMOVED***, kr)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if rp.Provider() == "etcd" ***REMOVED***
			cm, err = crypt.NewStandardEtcdConfigManager([]string***REMOVED***rp.Endpoint()***REMOVED***)
		***REMOVED*** else ***REMOVED***
			cm, err = crypt.NewStandardConsulConfigManager([]string***REMOVED***rp.Endpoint()***REMOVED***)
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return cm, nil
***REMOVED***

func init() ***REMOVED***
	viper.RemoteConfig = &remoteConfigProvider***REMOVED******REMOVED***
***REMOVED***
