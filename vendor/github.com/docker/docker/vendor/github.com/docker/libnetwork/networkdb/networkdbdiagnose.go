package networkdb

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/libnetwork/common"
	"github.com/docker/libnetwork/diagnose"
	"github.com/sirupsen/logrus"
)

const (
	missingParameter = "missing parameter"
	dbNotAvailable   = "database not available"
)

// NetDbPaths2Func TODO
var NetDbPaths2Func = map[string]diagnose.HTTPHandlerFunc***REMOVED***
	"/join":         dbJoin,
	"/networkpeers": dbPeers,
	"/clusterpeers": dbClusterPeers,
	"/joinnetwork":  dbJoinNetwork,
	"/leavenetwork": dbLeaveNetwork,
	"/createentry":  dbCreateEntry,
	"/updateentry":  dbUpdateEntry,
	"/deleteentry":  dbDeleteEntry,
	"/getentry":     dbGetEntry,
	"/gettable":     dbGetTable,
***REMOVED***

func dbJoin(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	_, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("join cluster")

	if len(r.Form["members"]) < 1 ***REMOVED***
		rsp := diagnose.WrongCommand(missingParameter, fmt.Sprintf("%s?members=ip1,ip2,...", r.URL.Path))
		log.Error("join cluster failed, wrong input")
		diagnose.HTTPReply(w, rsp, json)
		return
	***REMOVED***

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		err := nDB.Join(strings.Split(r.Form["members"][0], ","))
		if err != nil ***REMOVED***
			rsp := diagnose.FailCommand(fmt.Errorf("%s error in the DB join %s", r.URL.Path, err))
			log.WithError(err).Error("join cluster failed")
			diagnose.HTTPReply(w, rsp, json)
			return
		***REMOVED***

		log.Info("join cluster done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(nil), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***

func dbPeers(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	_, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("network peers")

	if len(r.Form["nid"]) < 1 ***REMOVED***
		rsp := diagnose.WrongCommand(missingParameter, fmt.Sprintf("%s?nid=test", r.URL.Path))
		log.Error("network peers failed, wrong input")
		diagnose.HTTPReply(w, rsp, json)
		return
	***REMOVED***

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		peers := nDB.Peers(r.Form["nid"][0])
		rsp := &diagnose.TableObj***REMOVED***Length: len(peers)***REMOVED***
		for i, peerInfo := range peers ***REMOVED***
			rsp.Elements = append(rsp.Elements, &diagnose.PeerEntryObj***REMOVED***Index: i, Name: peerInfo.Name, IP: peerInfo.IP***REMOVED***)
		***REMOVED***
		log.WithField("response", fmt.Sprintf("%+v", rsp)).Info("network peers done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(rsp), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***

func dbClusterPeers(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	_, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("cluster peers")

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		peers := nDB.ClusterPeers()
		rsp := &diagnose.TableObj***REMOVED***Length: len(peers)***REMOVED***
		for i, peerInfo := range peers ***REMOVED***
			rsp.Elements = append(rsp.Elements, &diagnose.PeerEntryObj***REMOVED***Index: i, Name: peerInfo.Name, IP: peerInfo.IP***REMOVED***)
		***REMOVED***
		log.WithField("response", fmt.Sprintf("%+v", rsp)).Info("cluster peers done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(rsp), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***

func dbCreateEntry(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	unsafe, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("create entry")

	if len(r.Form["tname"]) < 1 ||
		len(r.Form["nid"]) < 1 ||
		len(r.Form["key"]) < 1 ||
		len(r.Form["value"]) < 1 ***REMOVED***
		rsp := diagnose.WrongCommand(missingParameter, fmt.Sprintf("%s?tname=table_name&nid=network_id&key=k&value=v", r.URL.Path))
		log.Error("create entry failed, wrong input")
		diagnose.HTTPReply(w, rsp, json)
		return
	***REMOVED***

	tname := r.Form["tname"][0]
	nid := r.Form["nid"][0]
	key := r.Form["key"][0]
	value := r.Form["value"][0]
	decodedValue := []byte(value)
	if !unsafe ***REMOVED***
		var err error
		decodedValue, err = base64.StdEncoding.DecodeString(value)
		if err != nil ***REMOVED***
			log.WithError(err).Error("create entry failed")
			diagnose.HTTPReply(w, diagnose.FailCommand(err), json)
			return
		***REMOVED***
	***REMOVED***

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		if err := nDB.CreateEntry(tname, nid, key, decodedValue); err != nil ***REMOVED***
			rsp := diagnose.FailCommand(err)
			diagnose.HTTPReply(w, rsp, json)
			log.WithError(err).Error("create entry failed")
			return
		***REMOVED***
		log.Info("create entry done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(nil), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***

func dbUpdateEntry(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	unsafe, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("update entry")

	if len(r.Form["tname"]) < 1 ||
		len(r.Form["nid"]) < 1 ||
		len(r.Form["key"]) < 1 ||
		len(r.Form["value"]) < 1 ***REMOVED***
		rsp := diagnose.WrongCommand(missingParameter, fmt.Sprintf("%s?tname=table_name&nid=network_id&key=k&value=v", r.URL.Path))
		log.Error("update entry failed, wrong input")
		diagnose.HTTPReply(w, rsp, json)
		return
	***REMOVED***

	tname := r.Form["tname"][0]
	nid := r.Form["nid"][0]
	key := r.Form["key"][0]
	value := r.Form["value"][0]
	decodedValue := []byte(value)
	if !unsafe ***REMOVED***
		var err error
		decodedValue, err = base64.StdEncoding.DecodeString(value)
		if err != nil ***REMOVED***
			log.WithError(err).Error("update entry failed")
			diagnose.HTTPReply(w, diagnose.FailCommand(err), json)
			return
		***REMOVED***
	***REMOVED***

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		if err := nDB.UpdateEntry(tname, nid, key, decodedValue); err != nil ***REMOVED***
			log.WithError(err).Error("update entry failed")
			diagnose.HTTPReply(w, diagnose.FailCommand(err), json)
			return
		***REMOVED***
		log.Info("update entry done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(nil), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***

func dbDeleteEntry(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	_, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("delete entry")

	if len(r.Form["tname"]) < 1 ||
		len(r.Form["nid"]) < 1 ||
		len(r.Form["key"]) < 1 ***REMOVED***
		rsp := diagnose.WrongCommand(missingParameter, fmt.Sprintf("%s?tname=table_name&nid=network_id&key=k", r.URL.Path))
		log.Error("delete entry failed, wrong input")
		diagnose.HTTPReply(w, rsp, json)
		return
	***REMOVED***

	tname := r.Form["tname"][0]
	nid := r.Form["nid"][0]
	key := r.Form["key"][0]

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		err := nDB.DeleteEntry(tname, nid, key)
		if err != nil ***REMOVED***
			log.WithError(err).Error("delete entry failed")
			diagnose.HTTPReply(w, diagnose.FailCommand(err), json)
			return
		***REMOVED***
		log.Info("delete entry done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(nil), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***

func dbGetEntry(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	unsafe, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("get entry")

	if len(r.Form["tname"]) < 1 ||
		len(r.Form["nid"]) < 1 ||
		len(r.Form["key"]) < 1 ***REMOVED***
		rsp := diagnose.WrongCommand(missingParameter, fmt.Sprintf("%s?tname=table_name&nid=network_id&key=k", r.URL.Path))
		log.Error("get entry failed, wrong input")
		diagnose.HTTPReply(w, rsp, json)
		return
	***REMOVED***

	tname := r.Form["tname"][0]
	nid := r.Form["nid"][0]
	key := r.Form["key"][0]

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		value, err := nDB.GetEntry(tname, nid, key)
		if err != nil ***REMOVED***
			log.WithError(err).Error("get entry failed")
			diagnose.HTTPReply(w, diagnose.FailCommand(err), json)
			return
		***REMOVED***

		var encodedValue string
		if unsafe ***REMOVED***
			encodedValue = string(value)
		***REMOVED*** else ***REMOVED***
			encodedValue = base64.StdEncoding.EncodeToString(value)
		***REMOVED***

		rsp := &diagnose.TableEntryObj***REMOVED***Key: key, Value: encodedValue***REMOVED***
		log.WithField("response", fmt.Sprintf("%+v", rsp)).Info("update entry done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(rsp), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***

func dbJoinNetwork(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	_, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("join network")

	if len(r.Form["nid"]) < 1 ***REMOVED***
		rsp := diagnose.WrongCommand(missingParameter, fmt.Sprintf("%s?nid=network_id", r.URL.Path))
		log.Error("join network failed, wrong input")
		diagnose.HTTPReply(w, rsp, json)
		return
	***REMOVED***

	nid := r.Form["nid"][0]

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		if err := nDB.JoinNetwork(nid); err != nil ***REMOVED***
			log.WithError(err).Error("join network failed")
			diagnose.HTTPReply(w, diagnose.FailCommand(err), json)
			return
		***REMOVED***
		log.Info("join network done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(nil), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***

func dbLeaveNetwork(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	_, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("leave network")

	if len(r.Form["nid"]) < 1 ***REMOVED***
		rsp := diagnose.WrongCommand(missingParameter, fmt.Sprintf("%s?nid=network_id", r.URL.Path))
		log.Error("leave network failed, wrong input")
		diagnose.HTTPReply(w, rsp, json)
		return
	***REMOVED***

	nid := r.Form["nid"][0]

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		if err := nDB.LeaveNetwork(nid); err != nil ***REMOVED***
			log.WithError(err).Error("leave network failed")
			diagnose.HTTPReply(w, diagnose.FailCommand(err), json)
			return
		***REMOVED***
		log.Info("leave network done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(nil), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***

func dbGetTable(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	diagnose.DebugHTTPForm(r)
	unsafe, json := diagnose.ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("get table")

	if len(r.Form["tname"]) < 1 ||
		len(r.Form["nid"]) < 1 ***REMOVED***
		rsp := diagnose.WrongCommand(missingParameter, fmt.Sprintf("%s?tname=table_name&nid=network_id", r.URL.Path))
		log.Error("get table failed, wrong input")
		diagnose.HTTPReply(w, rsp, json)
		return
	***REMOVED***

	tname := r.Form["tname"][0]
	nid := r.Form["nid"][0]

	nDB, ok := ctx.(*NetworkDB)
	if ok ***REMOVED***
		table := nDB.GetTableByNetwork(tname, nid)
		rsp := &diagnose.TableObj***REMOVED***Length: len(table)***REMOVED***
		var i = 0
		for k, v := range table ***REMOVED***
			var encodedValue string
			if unsafe ***REMOVED***
				encodedValue = string(v.Value)
			***REMOVED*** else ***REMOVED***
				encodedValue = base64.StdEncoding.EncodeToString(v.Value)
			***REMOVED***
			rsp.Elements = append(rsp.Elements,
				&diagnose.TableEntryObj***REMOVED***
					Index: i,
					Key:   k,
					Value: encodedValue,
					Owner: v.owner,
				***REMOVED***)
			i++
		***REMOVED***
		log.WithField("response", fmt.Sprintf("%+v", rsp)).Info("get table done")
		diagnose.HTTPReply(w, diagnose.CommandSucceed(rsp), json)
		return
	***REMOVED***
	diagnose.HTTPReply(w, diagnose.FailCommand(fmt.Errorf("%s", dbNotAvailable)), json)
***REMOVED***
