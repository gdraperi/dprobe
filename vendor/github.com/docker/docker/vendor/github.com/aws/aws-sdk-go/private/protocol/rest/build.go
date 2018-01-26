// Package rest provides RESTful serialization of AWS requests and responses.
package rest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol"
)

// RFC822 returns an RFC822 formatted timestamp for AWS protocols
const RFC822 = "Mon, 2 Jan 2006 15:04:05 GMT"

// Whether the byte value can be sent without escaping in AWS URLs
var noEscape [256]bool

var errValueNotSet = fmt.Errorf("value not set")

func init() ***REMOVED***
	for i := 0; i < len(noEscape); i++ ***REMOVED***
		// AWS expects every character except these to be escaped
		noEscape[i] = (i >= 'A' && i <= 'Z') ||
			(i >= 'a' && i <= 'z') ||
			(i >= '0' && i <= '9') ||
			i == '-' ||
			i == '.' ||
			i == '_' ||
			i == '~'
	***REMOVED***
***REMOVED***

// BuildHandler is a named request handler for building rest protocol requests
var BuildHandler = request.NamedHandler***REMOVED***Name: "awssdk.rest.Build", Fn: Build***REMOVED***

// Build builds the REST component of a service request.
func Build(r *request.Request) ***REMOVED***
	if r.ParamsFilled() ***REMOVED***
		v := reflect.ValueOf(r.Params).Elem()
		buildLocationElements(r, v, false)
		buildBody(r, v)
	***REMOVED***
***REMOVED***

// BuildAsGET builds the REST component of a service request with the ability to hoist
// data from the body.
func BuildAsGET(r *request.Request) ***REMOVED***
	if r.ParamsFilled() ***REMOVED***
		v := reflect.ValueOf(r.Params).Elem()
		buildLocationElements(r, v, true)
		buildBody(r, v)
	***REMOVED***
***REMOVED***

func buildLocationElements(r *request.Request, v reflect.Value, buildGETQuery bool) ***REMOVED***
	query := r.HTTPRequest.URL.Query()

	// Setup the raw path to match the base path pattern. This is needed
	// so that when the path is mutated a custom escaped version can be
	// stored in RawPath that will be used by the Go client.
	r.HTTPRequest.URL.RawPath = r.HTTPRequest.URL.Path

	for i := 0; i < v.NumField(); i++ ***REMOVED***
		m := v.Field(i)
		if n := v.Type().Field(i).Name; n[0:1] == strings.ToLower(n[0:1]) ***REMOVED***
			continue
		***REMOVED***

		if m.IsValid() ***REMOVED***
			field := v.Type().Field(i)
			name := field.Tag.Get("locationName")
			if name == "" ***REMOVED***
				name = field.Name
			***REMOVED***
			if kind := m.Kind(); kind == reflect.Ptr ***REMOVED***
				m = m.Elem()
			***REMOVED*** else if kind == reflect.Interface ***REMOVED***
				if !m.Elem().IsValid() ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			if !m.IsValid() ***REMOVED***
				continue
			***REMOVED***
			if field.Tag.Get("ignore") != "" ***REMOVED***
				continue
			***REMOVED***

			var err error
			switch field.Tag.Get("location") ***REMOVED***
			case "headers": // header maps
				err = buildHeaderMap(&r.HTTPRequest.Header, m, field.Tag)
			case "header":
				err = buildHeader(&r.HTTPRequest.Header, m, name, field.Tag)
			case "uri":
				err = buildURI(r.HTTPRequest.URL, m, name, field.Tag)
			case "querystring":
				err = buildQueryString(query, m, name, field.Tag)
			default:
				if buildGETQuery ***REMOVED***
					err = buildQueryString(query, m, name, field.Tag)
				***REMOVED***
			***REMOVED***
			r.Error = err
		***REMOVED***
		if r.Error != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	r.HTTPRequest.URL.RawQuery = query.Encode()
	if !aws.BoolValue(r.Config.DisableRestProtocolURICleaning) ***REMOVED***
		cleanPath(r.HTTPRequest.URL)
	***REMOVED***
***REMOVED***

func buildBody(r *request.Request, v reflect.Value) ***REMOVED***
	if field, ok := v.Type().FieldByName("_"); ok ***REMOVED***
		if payloadName := field.Tag.Get("payload"); payloadName != "" ***REMOVED***
			pfield, _ := v.Type().FieldByName(payloadName)
			if ptag := pfield.Tag.Get("type"); ptag != "" && ptag != "structure" ***REMOVED***
				payload := reflect.Indirect(v.FieldByName(payloadName))
				if payload.IsValid() && payload.Interface() != nil ***REMOVED***
					switch reader := payload.Interface().(type) ***REMOVED***
					case io.ReadSeeker:
						r.SetReaderBody(reader)
					case []byte:
						r.SetBufferBody(reader)
					case string:
						r.SetStringBody(reader)
					default:
						r.Error = awserr.New("SerializationError",
							"failed to encode REST request",
							fmt.Errorf("unknown payload type %s", payload.Type()))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func buildHeader(header *http.Header, v reflect.Value, name string, tag reflect.StructTag) error ***REMOVED***
	str, err := convertType(v, tag)
	if err == errValueNotSet ***REMOVED***
		return nil
	***REMOVED*** else if err != nil ***REMOVED***
		return awserr.New("SerializationError", "failed to encode REST request", err)
	***REMOVED***

	header.Add(name, str)

	return nil
***REMOVED***

func buildHeaderMap(header *http.Header, v reflect.Value, tag reflect.StructTag) error ***REMOVED***
	prefix := tag.Get("locationName")
	for _, key := range v.MapKeys() ***REMOVED***
		str, err := convertType(v.MapIndex(key), tag)
		if err == errValueNotSet ***REMOVED***
			continue
		***REMOVED*** else if err != nil ***REMOVED***
			return awserr.New("SerializationError", "failed to encode REST request", err)

		***REMOVED***

		header.Add(prefix+key.String(), str)
	***REMOVED***
	return nil
***REMOVED***

func buildURI(u *url.URL, v reflect.Value, name string, tag reflect.StructTag) error ***REMOVED***
	value, err := convertType(v, tag)
	if err == errValueNotSet ***REMOVED***
		return nil
	***REMOVED*** else if err != nil ***REMOVED***
		return awserr.New("SerializationError", "failed to encode REST request", err)
	***REMOVED***

	u.Path = strings.Replace(u.Path, "***REMOVED***"+name+"***REMOVED***", value, -1)
	u.Path = strings.Replace(u.Path, "***REMOVED***"+name+"+***REMOVED***", value, -1)

	u.RawPath = strings.Replace(u.RawPath, "***REMOVED***"+name+"***REMOVED***", EscapePath(value, true), -1)
	u.RawPath = strings.Replace(u.RawPath, "***REMOVED***"+name+"+***REMOVED***", EscapePath(value, false), -1)

	return nil
***REMOVED***

func buildQueryString(query url.Values, v reflect.Value, name string, tag reflect.StructTag) error ***REMOVED***
	switch value := v.Interface().(type) ***REMOVED***
	case []*string:
		for _, item := range value ***REMOVED***
			query.Add(name, *item)
		***REMOVED***
	case map[string]*string:
		for key, item := range value ***REMOVED***
			query.Add(key, *item)
		***REMOVED***
	case map[string][]*string:
		for key, items := range value ***REMOVED***
			for _, item := range items ***REMOVED***
				query.Add(key, *item)
			***REMOVED***
		***REMOVED***
	default:
		str, err := convertType(v, tag)
		if err == errValueNotSet ***REMOVED***
			return nil
		***REMOVED*** else if err != nil ***REMOVED***
			return awserr.New("SerializationError", "failed to encode REST request", err)
		***REMOVED***
		query.Set(name, str)
	***REMOVED***

	return nil
***REMOVED***

func cleanPath(u *url.URL) ***REMOVED***
	hasSlash := strings.HasSuffix(u.Path, "/")

	// clean up path, removing duplicate `/`
	u.Path = path.Clean(u.Path)
	u.RawPath = path.Clean(u.RawPath)

	if hasSlash && !strings.HasSuffix(u.Path, "/") ***REMOVED***
		u.Path += "/"
		u.RawPath += "/"
	***REMOVED***
***REMOVED***

// EscapePath escapes part of a URL path in Amazon style
func EscapePath(path string, encodeSep bool) string ***REMOVED***
	var buf bytes.Buffer
	for i := 0; i < len(path); i++ ***REMOVED***
		c := path[i]
		if noEscape[c] || (c == '/' && !encodeSep) ***REMOVED***
			buf.WriteByte(c)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(&buf, "%%%02X", c)
		***REMOVED***
	***REMOVED***
	return buf.String()
***REMOVED***

func convertType(v reflect.Value, tag reflect.StructTag) (str string, err error) ***REMOVED***
	v = reflect.Indirect(v)
	if !v.IsValid() ***REMOVED***
		return "", errValueNotSet
	***REMOVED***

	switch value := v.Interface().(type) ***REMOVED***
	case string:
		str = value
	case []byte:
		str = base64.StdEncoding.EncodeToString(value)
	case bool:
		str = strconv.FormatBool(value)
	case int64:
		str = strconv.FormatInt(value, 10)
	case float64:
		str = strconv.FormatFloat(value, 'f', -1, 64)
	case time.Time:
		str = value.UTC().Format(RFC822)
	case aws.JSONValue:
		if len(value) == 0 ***REMOVED***
			return "", errValueNotSet
		***REMOVED***
		escaping := protocol.NoEscape
		if tag.Get("location") == "header" ***REMOVED***
			escaping = protocol.Base64Escape
		***REMOVED***
		str, err = protocol.EncodeJSONValue(value, escaping)
		if err != nil ***REMOVED***
			return "", fmt.Errorf("unable to encode JSONValue, %v", err)
		***REMOVED***
	default:
		err := fmt.Errorf("unsupported value for param %v (%s)", v.Interface(), v.Type())
		return "", err
	***REMOVED***
	return str, nil
***REMOVED***
