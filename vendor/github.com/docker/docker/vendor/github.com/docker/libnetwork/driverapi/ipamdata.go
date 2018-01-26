package driverapi

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/docker/libnetwork/types"
)

// MarshalJSON encodes IPAMData into json message
func (i *IPAMData) MarshalJSON() ([]byte, error) ***REMOVED***
	m := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	m["AddressSpace"] = i.AddressSpace
	if i.Pool != nil ***REMOVED***
		m["Pool"] = i.Pool.String()
	***REMOVED***
	if i.Gateway != nil ***REMOVED***
		m["Gateway"] = i.Gateway.String()
	***REMOVED***
	if i.AuxAddresses != nil ***REMOVED***
		am := make(map[string]string, len(i.AuxAddresses))
		for k, v := range i.AuxAddresses ***REMOVED***
			am[k] = v.String()
		***REMOVED***
		m["AuxAddresses"] = am
	***REMOVED***
	return json.Marshal(m)
***REMOVED***

// UnmarshalJSON decodes a json message into IPAMData
func (i *IPAMData) UnmarshalJSON(data []byte) error ***REMOVED***
	var (
		m   map[string]interface***REMOVED******REMOVED***
		err error
	)
	if err := json.Unmarshal(data, &m); err != nil ***REMOVED***
		return err
	***REMOVED***
	i.AddressSpace = m["AddressSpace"].(string)
	if v, ok := m["Pool"]; ok ***REMOVED***
		if i.Pool, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if v, ok := m["Gateway"]; ok ***REMOVED***
		if i.Gateway, err = types.ParseCIDR(v.(string)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if v, ok := m["AuxAddresses"]; ok ***REMOVED***
		b, _ := json.Marshal(v)
		var am map[string]string
		if err = json.Unmarshal(b, &am); err != nil ***REMOVED***
			return err
		***REMOVED***
		i.AuxAddresses = make(map[string]*net.IPNet, len(am))
		for k, v := range am ***REMOVED***
			if i.AuxAddresses[k], err = types.ParseCIDR(v); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Validate checks whether the IPAMData structure contains congruent data
func (i *IPAMData) Validate() error ***REMOVED***
	var isV6 bool
	if i.Pool == nil ***REMOVED***
		return types.BadRequestErrorf("invalid pool")
	***REMOVED***
	if i.Gateway == nil ***REMOVED***
		return types.BadRequestErrorf("invalid gateway address")
	***REMOVED***
	isV6 = i.IsV6()
	if isV6 && i.Gateway.IP.To4() != nil || !isV6 && i.Gateway.IP.To4() == nil ***REMOVED***
		return types.BadRequestErrorf("incongruent ip versions for pool and gateway")
	***REMOVED***
	for k, sip := range i.AuxAddresses ***REMOVED***
		if isV6 && sip.IP.To4() != nil || !isV6 && sip.IP.To4() == nil ***REMOVED***
			return types.BadRequestErrorf("incongruent ip versions for pool and secondary ip address %s", k)
		***REMOVED***
	***REMOVED***
	if !i.Pool.Contains(i.Gateway.IP) ***REMOVED***
		return types.BadRequestErrorf("invalid gateway address (%s) does not belong to the pool (%s)", i.Gateway, i.Pool)
	***REMOVED***
	for k, sip := range i.AuxAddresses ***REMOVED***
		if !i.Pool.Contains(sip.IP) ***REMOVED***
			return types.BadRequestErrorf("invalid secondary address %s (%s) does not belong to the pool (%s)", k, i.Gateway, i.Pool)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// IsV6 returns whether this is an IPv6 IPAMData structure
func (i *IPAMData) IsV6() bool ***REMOVED***
	return nil == i.Pool.IP.To4()
***REMOVED***

func (i *IPAMData) String() string ***REMOVED***
	return fmt.Sprintf("AddressSpace: %s\nPool: %v\nGateway: %v\nAddresses: %v", i.AddressSpace, i.Pool, i.Gateway, i.AuxAddresses)
***REMOVED***
