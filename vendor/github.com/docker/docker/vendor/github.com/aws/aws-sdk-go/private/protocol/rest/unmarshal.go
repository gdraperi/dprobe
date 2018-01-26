package rest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol"
)

// UnmarshalHandler is a named request handler for unmarshaling rest protocol requests
var UnmarshalHandler = request.NamedHandler***REMOVED***Name: "awssdk.rest.Unmarshal", Fn: Unmarshal***REMOVED***

// UnmarshalMetaHandler is a named request handler for unmarshaling rest protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler***REMOVED***Name: "awssdk.rest.UnmarshalMeta", Fn: UnmarshalMeta***REMOVED***

// Unmarshal unmarshals the REST component of a response in a REST service.
func Unmarshal(r *request.Request) ***REMOVED***
	if r.DataFilled() ***REMOVED***
		v := reflect.Indirect(reflect.ValueOf(r.Data))
		unmarshalBody(r, v)
	***REMOVED***
***REMOVED***

// UnmarshalMeta unmarshals the REST metadata of a response in a REST service
func UnmarshalMeta(r *request.Request) ***REMOVED***
	r.RequestID = r.HTTPResponse.Header.Get("X-Amzn-Requestid")
	if r.RequestID == "" ***REMOVED***
		// Alternative version of request id in the header
		r.RequestID = r.HTTPResponse.Header.Get("X-Amz-Request-Id")
	***REMOVED***
	if r.DataFilled() ***REMOVED***
		v := reflect.Indirect(reflect.ValueOf(r.Data))
		unmarshalLocationElements(r, v)
	***REMOVED***
***REMOVED***

func unmarshalBody(r *request.Request, v reflect.Value) ***REMOVED***
	if field, ok := v.Type().FieldByName("_"); ok ***REMOVED***
		if payloadName := field.Tag.Get("payload"); payloadName != "" ***REMOVED***
			pfield, _ := v.Type().FieldByName(payloadName)
			if ptag := pfield.Tag.Get("type"); ptag != "" && ptag != "structure" ***REMOVED***
				payload := v.FieldByName(payloadName)
				if payload.IsValid() ***REMOVED***
					switch payload.Interface().(type) ***REMOVED***
					case []byte:
						defer r.HTTPResponse.Body.Close()
						b, err := ioutil.ReadAll(r.HTTPResponse.Body)
						if err != nil ***REMOVED***
							r.Error = awserr.New("SerializationError", "failed to decode REST response", err)
						***REMOVED*** else ***REMOVED***
							payload.Set(reflect.ValueOf(b))
						***REMOVED***
					case *string:
						defer r.HTTPResponse.Body.Close()
						b, err := ioutil.ReadAll(r.HTTPResponse.Body)
						if err != nil ***REMOVED***
							r.Error = awserr.New("SerializationError", "failed to decode REST response", err)
						***REMOVED*** else ***REMOVED***
							str := string(b)
							payload.Set(reflect.ValueOf(&str))
						***REMOVED***
					default:
						switch payload.Type().String() ***REMOVED***
						case "io.ReadCloser":
							payload.Set(reflect.ValueOf(r.HTTPResponse.Body))
						case "io.ReadSeeker":
							b, err := ioutil.ReadAll(r.HTTPResponse.Body)
							if err != nil ***REMOVED***
								r.Error = awserr.New("SerializationError",
									"failed to read response body", err)
								return
							***REMOVED***
							payload.Set(reflect.ValueOf(ioutil.NopCloser(bytes.NewReader(b))))
						default:
							io.Copy(ioutil.Discard, r.HTTPResponse.Body)
							defer r.HTTPResponse.Body.Close()
							r.Error = awserr.New("SerializationError",
								"failed to decode REST response",
								fmt.Errorf("unknown payload type %s", payload.Type()))
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func unmarshalLocationElements(r *request.Request, v reflect.Value) ***REMOVED***
	for i := 0; i < v.NumField(); i++ ***REMOVED***
		m, field := v.Field(i), v.Type().Field(i)
		if n := field.Name; n[0:1] == strings.ToLower(n[0:1]) ***REMOVED***
			continue
		***REMOVED***

		if m.IsValid() ***REMOVED***
			name := field.Tag.Get("locationName")
			if name == "" ***REMOVED***
				name = field.Name
			***REMOVED***

			switch field.Tag.Get("location") ***REMOVED***
			case "statusCode":
				unmarshalStatusCode(m, r.HTTPResponse.StatusCode)
			case "header":
				err := unmarshalHeader(m, r.HTTPResponse.Header.Get(name), field.Tag)
				if err != nil ***REMOVED***
					r.Error = awserr.New("SerializationError", "failed to decode REST response", err)
					break
				***REMOVED***
			case "headers":
				prefix := field.Tag.Get("locationName")
				err := unmarshalHeaderMap(m, r.HTTPResponse.Header, prefix)
				if err != nil ***REMOVED***
					r.Error = awserr.New("SerializationError", "failed to decode REST response", err)
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if r.Error != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func unmarshalStatusCode(v reflect.Value, statusCode int) ***REMOVED***
	if !v.IsValid() ***REMOVED***
		return
	***REMOVED***

	switch v.Interface().(type) ***REMOVED***
	case *int64:
		s := int64(statusCode)
		v.Set(reflect.ValueOf(&s))
	***REMOVED***
***REMOVED***

func unmarshalHeaderMap(r reflect.Value, headers http.Header, prefix string) error ***REMOVED***
	switch r.Interface().(type) ***REMOVED***
	case map[string]*string: // we only support string map value types
		out := map[string]*string***REMOVED******REMOVED***
		for k, v := range headers ***REMOVED***
			k = http.CanonicalHeaderKey(k)
			if strings.HasPrefix(strings.ToLower(k), strings.ToLower(prefix)) ***REMOVED***
				out[k[len(prefix):]] = &v[0]
			***REMOVED***
		***REMOVED***
		r.Set(reflect.ValueOf(out))
	***REMOVED***
	return nil
***REMOVED***

func unmarshalHeader(v reflect.Value, header string, tag reflect.StructTag) error ***REMOVED***
	isJSONValue := tag.Get("type") == "jsonvalue"
	if isJSONValue ***REMOVED***
		if len(header) == 0 ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED*** else if !v.IsValid() || (header == "" && v.Elem().Kind() != reflect.String) ***REMOVED***
		return nil
	***REMOVED***

	switch v.Interface().(type) ***REMOVED***
	case *string:
		v.Set(reflect.ValueOf(&header))
	case []byte:
		b, err := base64.StdEncoding.DecodeString(header)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v.Set(reflect.ValueOf(&b))
	case *bool:
		b, err := strconv.ParseBool(header)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v.Set(reflect.ValueOf(&b))
	case *int64:
		i, err := strconv.ParseInt(header, 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v.Set(reflect.ValueOf(&i))
	case *float64:
		f, err := strconv.ParseFloat(header, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v.Set(reflect.ValueOf(&f))
	case *time.Time:
		t, err := time.Parse(RFC822, header)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v.Set(reflect.ValueOf(&t))
	case aws.JSONValue:
		escaping := protocol.NoEscape
		if tag.Get("location") == "header" ***REMOVED***
			escaping = protocol.Base64Escape
		***REMOVED***
		m, err := protocol.DecodeJSONValue(header, escaping)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v.Set(reflect.ValueOf(m))
	default:
		err := fmt.Errorf("Unsupported value for param %v (%s)", v.Interface(), v.Type())
		return err
	***REMOVED***
	return nil
***REMOVED***
