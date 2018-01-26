package network

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/errdefs"
	"github.com/docker/libnetwork"
	netconst "github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/networkdb"
	"github.com/pkg/errors"
)

var (
	// acceptedNetworkFilters is a list of acceptable filters
	acceptedNetworkFilters = map[string]bool***REMOVED***
		"driver": true,
		"type":   true,
		"name":   true,
		"id":     true,
		"label":  true,
		"scope":  true,
	***REMOVED***
)

func (n *networkRouter) getNetworksList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	filter := r.Form.Get("filters")
	netFilters, err := filters.FromJSON(filter)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := netFilters.Validate(acceptedNetworkFilters); err != nil ***REMOVED***
		return err
	***REMOVED***

	list := []types.NetworkResource***REMOVED******REMOVED***

	if nr, err := n.cluster.GetNetworks(); err == nil ***REMOVED***
		list = append(list, nr...)
	***REMOVED***

	// Combine the network list returned by Docker daemon if it is not already
	// returned by the cluster manager
SKIP:
	for _, nw := range n.backend.GetNetworks() ***REMOVED***
		for _, nl := range list ***REMOVED***
			if nl.ID == nw.ID() ***REMOVED***
				continue SKIP
			***REMOVED***
		***REMOVED***

		var nr *types.NetworkResource
		// Versions < 1.28 fetches all the containers attached to a network
		// in a network list api call. It is a heavy weight operation when
		// run across all the networks. Starting API version 1.28, this detailed
		// info is available for network specific GET API (equivalent to inspect)
		if versions.LessThan(httputils.VersionFromContext(ctx), "1.28") ***REMOVED***
			nr = n.buildDetailedNetworkResources(nw, false)
		***REMOVED*** else ***REMOVED***
			nr = n.buildNetworkResource(nw)
		***REMOVED***
		list = append(list, *nr)
	***REMOVED***

	list, err = filterNetworks(list, netFilters)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, list)
***REMOVED***

type invalidRequestError struct ***REMOVED***
	cause error
***REMOVED***

func (e invalidRequestError) Error() string ***REMOVED***
	return e.cause.Error()
***REMOVED***

func (e invalidRequestError) InvalidParameter() ***REMOVED******REMOVED***

type ambigousResultsError string

func (e ambigousResultsError) Error() string ***REMOVED***
	return "network " + string(e) + " is ambiguous"
***REMOVED***

func (ambigousResultsError) InvalidParameter() ***REMOVED******REMOVED***

func nameConflict(name string) error ***REMOVED***
	return errdefs.Conflict(libnetwork.NetworkNameError(name))
***REMOVED***

func (n *networkRouter) getNetwork(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	term := vars["id"]
	var (
		verbose bool
		err     error
	)
	if v := r.URL.Query().Get("verbose"); v != "" ***REMOVED***
		if verbose, err = strconv.ParseBool(v); err != nil ***REMOVED***
			return errors.Wrapf(invalidRequestError***REMOVED***err***REMOVED***, "invalid value for verbose: %s", v)
		***REMOVED***
	***REMOVED***
	scope := r.URL.Query().Get("scope")

	isMatchingScope := func(scope, term string) bool ***REMOVED***
		if term != "" ***REMOVED***
			return scope == term
		***REMOVED***
		return true
	***REMOVED***

	// In case multiple networks have duplicate names, return error.
	// TODO (yongtang): should we wrap with version here for backward compatibility?

	// First find based on full ID, return immediately once one is found.
	// If a network appears both in swarm and local, assume it is in local first

	// For full name and partial ID, save the result first, and process later
	// in case multiple records was found based on the same term
	listByFullName := map[string]types.NetworkResource***REMOVED******REMOVED***
	listByPartialID := map[string]types.NetworkResource***REMOVED******REMOVED***

	nw := n.backend.GetNetworks()
	for _, network := range nw ***REMOVED***
		if network.ID() == term && isMatchingScope(network.Info().Scope(), scope) ***REMOVED***
			return httputils.WriteJSON(w, http.StatusOK, *n.buildDetailedNetworkResources(network, verbose))
		***REMOVED***
		if network.Name() == term && isMatchingScope(network.Info().Scope(), scope) ***REMOVED***
			// No need to check the ID collision here as we are still in
			// local scope and the network ID is unique in this scope.
			listByFullName[network.ID()] = *n.buildDetailedNetworkResources(network, verbose)
		***REMOVED***
		if strings.HasPrefix(network.ID(), term) && isMatchingScope(network.Info().Scope(), scope) ***REMOVED***
			// No need to check the ID collision here as we are still in
			// local scope and the network ID is unique in this scope.
			listByPartialID[network.ID()] = *n.buildDetailedNetworkResources(network, verbose)
		***REMOVED***
	***REMOVED***

	nwk, err := n.cluster.GetNetwork(term)
	if err == nil ***REMOVED***
		// If the get network is passed with a specific network ID / partial network ID
		// or if the get network was passed with a network name and scope as swarm
		// return the network. Skipped using isMatchingScope because it is true if the scope
		// is not set which would be case if the client API v1.30
		if strings.HasPrefix(nwk.ID, term) || (netconst.SwarmScope == scope) ***REMOVED***
			// If we have a previous match "backend", return it, we need verbose when enabled
			// ex: overlay/partial_ID or name/swarm_scope
			if nwv, ok := listByPartialID[nwk.ID]; ok ***REMOVED***
				nwk = nwv
			***REMOVED*** else if nwv, ok := listByFullName[nwk.ID]; ok ***REMOVED***
				nwk = nwv
			***REMOVED***
			return httputils.WriteJSON(w, http.StatusOK, nwk)
		***REMOVED***
	***REMOVED***

	nr, _ := n.cluster.GetNetworks()
	for _, network := range nr ***REMOVED***
		if network.ID == term && isMatchingScope(network.Scope, scope) ***REMOVED***
			return httputils.WriteJSON(w, http.StatusOK, network)
		***REMOVED***
		if network.Name == term && isMatchingScope(network.Scope, scope) ***REMOVED***
			// Check the ID collision as we are in swarm scope here, and
			// the map (of the listByFullName) may have already had a
			// network with the same ID (from local scope previously)
			if _, ok := listByFullName[network.ID]; !ok ***REMOVED***
				listByFullName[network.ID] = network
			***REMOVED***
		***REMOVED***
		if strings.HasPrefix(network.ID, term) && isMatchingScope(network.Scope, scope) ***REMOVED***
			// Check the ID collision as we are in swarm scope here, and
			// the map (of the listByPartialID) may have already had a
			// network with the same ID (from local scope previously)
			if _, ok := listByPartialID[network.ID]; !ok ***REMOVED***
				listByPartialID[network.ID] = network
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Find based on full name, returns true only if no duplicates
	if len(listByFullName) == 1 ***REMOVED***
		for _, v := range listByFullName ***REMOVED***
			return httputils.WriteJSON(w, http.StatusOK, v)
		***REMOVED***
	***REMOVED***
	if len(listByFullName) > 1 ***REMOVED***
		return errors.Wrapf(ambigousResultsError(term), "%d matches found based on name", len(listByFullName))
	***REMOVED***

	// Find based on partial ID, returns true only if no duplicates
	if len(listByPartialID) == 1 ***REMOVED***
		for _, v := range listByPartialID ***REMOVED***
			return httputils.WriteJSON(w, http.StatusOK, v)
		***REMOVED***
	***REMOVED***
	if len(listByPartialID) > 1 ***REMOVED***
		return errors.Wrapf(ambigousResultsError(term), "%d matches found based on ID prefix", len(listByPartialID))
	***REMOVED***

	return libnetwork.ErrNoSuchNetwork(term)
***REMOVED***

func (n *networkRouter) postNetworkCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var create types.NetworkCreateRequest

	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := json.NewDecoder(r.Body).Decode(&create); err != nil ***REMOVED***
		return err
	***REMOVED***

	if nws, err := n.cluster.GetNetworksByName(create.Name); err == nil && len(nws) > 0 ***REMOVED***
		return nameConflict(create.Name)
	***REMOVED***

	nw, err := n.backend.CreateNetwork(create)
	if err != nil ***REMOVED***
		var warning string
		if _, ok := err.(libnetwork.NetworkNameError); ok ***REMOVED***
			// check if user defined CheckDuplicate, if set true, return err
			// otherwise prepare a warning message
			if create.CheckDuplicate ***REMOVED***
				return nameConflict(create.Name)
			***REMOVED***
			warning = libnetwork.NetworkNameError(create.Name).Error()
		***REMOVED***

		if _, ok := err.(libnetwork.ManagerRedirectError); !ok ***REMOVED***
			return err
		***REMOVED***
		id, err := n.cluster.CreateNetwork(create)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		nw = &types.NetworkCreateResponse***REMOVED***
			ID:      id,
			Warning: warning,
		***REMOVED***
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusCreated, nw)
***REMOVED***

func (n *networkRouter) postNetworkConnect(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var connect types.NetworkConnect
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := json.NewDecoder(r.Body).Decode(&connect); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Unlike other operations, we does not check ambiguity of the name/ID here.
	// The reason is that, In case of attachable network in swarm scope, the actual local network
	// may not be available at the time. At the same time, inside daemon `ConnectContainerToNetwork`
	// does the ambiguity check anyway. Therefore, passing the name to daemon would be enough.
	return n.backend.ConnectContainerToNetwork(connect.Container, vars["id"], connect.EndpointConfig)
***REMOVED***

func (n *networkRouter) postNetworkDisconnect(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var disconnect types.NetworkDisconnect
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := json.NewDecoder(r.Body).Decode(&disconnect); err != nil ***REMOVED***
		return err
	***REMOVED***

	return n.backend.DisconnectContainerFromNetwork(disconnect.Container, vars["id"], disconnect.Force)
***REMOVED***

func (n *networkRouter) deleteNetwork(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	nw, err := n.findUniqueNetwork(vars["id"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if nw.Scope == "swarm" ***REMOVED***
		if err = n.cluster.RemoveNetwork(nw.ID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := n.backend.DeleteNetwork(nw.ID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	w.WriteHeader(http.StatusNoContent)
	return nil
***REMOVED***

func (n *networkRouter) buildNetworkResource(nw libnetwork.Network) *types.NetworkResource ***REMOVED***
	r := &types.NetworkResource***REMOVED******REMOVED***
	if nw == nil ***REMOVED***
		return r
	***REMOVED***

	info := nw.Info()
	r.Name = nw.Name()
	r.ID = nw.ID()
	r.Created = info.Created()
	r.Scope = info.Scope()
	r.Driver = nw.Type()
	r.EnableIPv6 = info.IPv6Enabled()
	r.Internal = info.Internal()
	r.Attachable = info.Attachable()
	r.Ingress = info.Ingress()
	r.Options = info.DriverOptions()
	r.Containers = make(map[string]types.EndpointResource)
	buildIpamResources(r, info)
	r.Labels = info.Labels()
	r.ConfigOnly = info.ConfigOnly()

	if cn := info.ConfigFrom(); cn != "" ***REMOVED***
		r.ConfigFrom = network.ConfigReference***REMOVED***Network: cn***REMOVED***
	***REMOVED***

	peers := info.Peers()
	if len(peers) != 0 ***REMOVED***
		r.Peers = buildPeerInfoResources(peers)
	***REMOVED***

	return r
***REMOVED***

func (n *networkRouter) buildDetailedNetworkResources(nw libnetwork.Network, verbose bool) *types.NetworkResource ***REMOVED***
	if nw == nil ***REMOVED***
		return &types.NetworkResource***REMOVED******REMOVED***
	***REMOVED***

	r := n.buildNetworkResource(nw)
	epl := nw.Endpoints()
	for _, e := range epl ***REMOVED***
		ei := e.Info()
		if ei == nil ***REMOVED***
			continue
		***REMOVED***
		sb := ei.Sandbox()
		tmpID := e.ID()
		key := "ep-" + tmpID
		if sb != nil ***REMOVED***
			key = sb.ContainerID()
		***REMOVED***

		r.Containers[key] = buildEndpointResource(tmpID, e.Name(), ei)
	***REMOVED***
	if !verbose ***REMOVED***
		return r
	***REMOVED***
	services := nw.Info().Services()
	r.Services = make(map[string]network.ServiceInfo)
	for name, service := range services ***REMOVED***
		tasks := []network.Task***REMOVED******REMOVED***
		for _, t := range service.Tasks ***REMOVED***
			tasks = append(tasks, network.Task***REMOVED***
				Name:       t.Name,
				EndpointID: t.EndpointID,
				EndpointIP: t.EndpointIP,
				Info:       t.Info,
			***REMOVED***)
		***REMOVED***
		r.Services[name] = network.ServiceInfo***REMOVED***
			VIP:          service.VIP,
			Ports:        service.Ports,
			Tasks:        tasks,
			LocalLBIndex: service.LocalLBIndex,
		***REMOVED***
	***REMOVED***
	return r
***REMOVED***

func buildPeerInfoResources(peers []networkdb.PeerInfo) []network.PeerInfo ***REMOVED***
	peerInfo := make([]network.PeerInfo, 0, len(peers))
	for _, peer := range peers ***REMOVED***
		peerInfo = append(peerInfo, network.PeerInfo***REMOVED***
			Name: peer.Name,
			IP:   peer.IP,
		***REMOVED***)
	***REMOVED***
	return peerInfo
***REMOVED***

func buildIpamResources(r *types.NetworkResource, nwInfo libnetwork.NetworkInfo) ***REMOVED***
	id, opts, ipv4conf, ipv6conf := nwInfo.IpamConfig()

	ipv4Info, ipv6Info := nwInfo.IpamInfo()

	r.IPAM.Driver = id

	r.IPAM.Options = opts

	r.IPAM.Config = []network.IPAMConfig***REMOVED******REMOVED***
	for _, ip4 := range ipv4conf ***REMOVED***
		if ip4.PreferredPool == "" ***REMOVED***
			continue
		***REMOVED***
		iData := network.IPAMConfig***REMOVED******REMOVED***
		iData.Subnet = ip4.PreferredPool
		iData.IPRange = ip4.SubPool
		iData.Gateway = ip4.Gateway
		iData.AuxAddress = ip4.AuxAddresses
		r.IPAM.Config = append(r.IPAM.Config, iData)
	***REMOVED***

	if len(r.IPAM.Config) == 0 ***REMOVED***
		for _, ip4Info := range ipv4Info ***REMOVED***
			iData := network.IPAMConfig***REMOVED******REMOVED***
			iData.Subnet = ip4Info.IPAMData.Pool.String()
			if ip4Info.IPAMData.Gateway != nil ***REMOVED***
				iData.Gateway = ip4Info.IPAMData.Gateway.IP.String()
			***REMOVED***
			r.IPAM.Config = append(r.IPAM.Config, iData)
		***REMOVED***
	***REMOVED***

	hasIpv6Conf := false
	for _, ip6 := range ipv6conf ***REMOVED***
		if ip6.PreferredPool == "" ***REMOVED***
			continue
		***REMOVED***
		hasIpv6Conf = true
		iData := network.IPAMConfig***REMOVED******REMOVED***
		iData.Subnet = ip6.PreferredPool
		iData.IPRange = ip6.SubPool
		iData.Gateway = ip6.Gateway
		iData.AuxAddress = ip6.AuxAddresses
		r.IPAM.Config = append(r.IPAM.Config, iData)
	***REMOVED***

	if !hasIpv6Conf ***REMOVED***
		for _, ip6Info := range ipv6Info ***REMOVED***
			if ip6Info.IPAMData.Pool == nil ***REMOVED***
				continue
			***REMOVED***
			iData := network.IPAMConfig***REMOVED******REMOVED***
			iData.Subnet = ip6Info.IPAMData.Pool.String()
			iData.Gateway = ip6Info.IPAMData.Gateway.String()
			r.IPAM.Config = append(r.IPAM.Config, iData)
		***REMOVED***
	***REMOVED***
***REMOVED***

func buildEndpointResource(id string, name string, info libnetwork.EndpointInfo) types.EndpointResource ***REMOVED***
	er := types.EndpointResource***REMOVED******REMOVED***

	er.EndpointID = id
	er.Name = name
	ei := info
	if ei == nil ***REMOVED***
		return er
	***REMOVED***

	if iface := ei.Iface(); iface != nil ***REMOVED***
		if mac := iface.MacAddress(); mac != nil ***REMOVED***
			er.MacAddress = mac.String()
		***REMOVED***
		if ip := iface.Address(); ip != nil && len(ip.IP) > 0 ***REMOVED***
			er.IPv4Address = ip.String()
		***REMOVED***

		if ipv6 := iface.AddressIPv6(); ipv6 != nil && len(ipv6.IP) > 0 ***REMOVED***
			er.IPv6Address = ipv6.String()
		***REMOVED***
	***REMOVED***
	return er
***REMOVED***

func (n *networkRouter) postNetworksPrune(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	pruneFilters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pruneReport, err := n.backend.NetworksPrune(ctx, pruneFilters)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, pruneReport)
***REMOVED***

// findUniqueNetwork will search network across different scopes (both local and swarm).
// NOTE: This findUniqueNetwork is different from FindNetwork in the daemon.
// In case multiple networks have duplicate names, return error.
// First find based on full ID, return immediately once one is found.
// If a network appears both in swarm and local, assume it is in local first
// For full name and partial ID, save the result first, and process later
// in case multiple records was found based on the same term
// TODO (yongtang): should we wrap with version here for backward compatibility?
func (n *networkRouter) findUniqueNetwork(term string) (types.NetworkResource, error) ***REMOVED***
	listByFullName := map[string]types.NetworkResource***REMOVED******REMOVED***
	listByPartialID := map[string]types.NetworkResource***REMOVED******REMOVED***

	nw := n.backend.GetNetworks()
	for _, network := range nw ***REMOVED***
		if network.ID() == term ***REMOVED***
			return *n.buildDetailedNetworkResources(network, false), nil

		***REMOVED***
		if network.Name() == term && !network.Info().Ingress() ***REMOVED***
			// No need to check the ID collision here as we are still in
			// local scope and the network ID is unique in this scope.
			listByFullName[network.ID()] = *n.buildDetailedNetworkResources(network, false)
		***REMOVED***
		if strings.HasPrefix(network.ID(), term) ***REMOVED***
			// No need to check the ID collision here as we are still in
			// local scope and the network ID is unique in this scope.
			listByPartialID[network.ID()] = *n.buildDetailedNetworkResources(network, false)
		***REMOVED***
	***REMOVED***

	nr, _ := n.cluster.GetNetworks()
	for _, network := range nr ***REMOVED***
		if network.ID == term ***REMOVED***
			return network, nil
		***REMOVED***
		if network.Name == term ***REMOVED***
			// Check the ID collision as we are in swarm scope here, and
			// the map (of the listByFullName) may have already had a
			// network with the same ID (from local scope previously)
			if _, ok := listByFullName[network.ID]; !ok ***REMOVED***
				listByFullName[network.ID] = network
			***REMOVED***
		***REMOVED***
		if strings.HasPrefix(network.ID, term) ***REMOVED***
			// Check the ID collision as we are in swarm scope here, and
			// the map (of the listByPartialID) may have already had a
			// network with the same ID (from local scope previously)
			if _, ok := listByPartialID[network.ID]; !ok ***REMOVED***
				listByPartialID[network.ID] = network
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Find based on full name, returns true only if no duplicates
	if len(listByFullName) == 1 ***REMOVED***
		for _, v := range listByFullName ***REMOVED***
			return v, nil
		***REMOVED***
	***REMOVED***
	if len(listByFullName) > 1 ***REMOVED***
		return types.NetworkResource***REMOVED******REMOVED***, errdefs.InvalidParameter(errors.Errorf("network %s is ambiguous (%d matches found based on name)", term, len(listByFullName)))
	***REMOVED***

	// Find based on partial ID, returns true only if no duplicates
	if len(listByPartialID) == 1 ***REMOVED***
		for _, v := range listByPartialID ***REMOVED***
			return v, nil
		***REMOVED***
	***REMOVED***
	if len(listByPartialID) > 1 ***REMOVED***
		return types.NetworkResource***REMOVED******REMOVED***, errdefs.InvalidParameter(errors.Errorf("network %s is ambiguous (%d matches found based on ID prefix)", term, len(listByPartialID)))
	***REMOVED***

	return types.NetworkResource***REMOVED******REMOVED***, errdefs.NotFound(libnetwork.ErrNoSuchNetwork(term))
***REMOVED***
