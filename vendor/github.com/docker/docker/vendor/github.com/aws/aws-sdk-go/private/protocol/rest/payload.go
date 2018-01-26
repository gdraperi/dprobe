package rest

import "reflect"

// PayloadMember returns the payload field member of i if there is one, or nil.
func PayloadMember(i interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if i == nil ***REMOVED***
		return nil
	***REMOVED***

	v := reflect.ValueOf(i).Elem()
	if !v.IsValid() ***REMOVED***
		return nil
	***REMOVED***
	if field, ok := v.Type().FieldByName("_"); ok ***REMOVED***
		if payloadName := field.Tag.Get("payload"); payloadName != "" ***REMOVED***
			field, _ := v.Type().FieldByName(payloadName)
			if field.Tag.Get("type") != "structure" ***REMOVED***
				return nil
			***REMOVED***

			payload := v.FieldByName(payloadName)
			if payload.IsValid() || (payload.Kind() == reflect.Ptr && !payload.IsNil()) ***REMOVED***
				return payload.Interface()
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// PayloadType returns the type of a payload field member of i if there is one, or "".
func PayloadType(i interface***REMOVED******REMOVED***) string ***REMOVED***
	v := reflect.Indirect(reflect.ValueOf(i))
	if !v.IsValid() ***REMOVED***
		return ""
	***REMOVED***
	if field, ok := v.Type().FieldByName("_"); ok ***REMOVED***
		if payloadName := field.Tag.Get("payload"); payloadName != "" ***REMOVED***
			if member, ok := v.Type().FieldByName(payloadName); ok ***REMOVED***
				return member.Tag.Get("type")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***
