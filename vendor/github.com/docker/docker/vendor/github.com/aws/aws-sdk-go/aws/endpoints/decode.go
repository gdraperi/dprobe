package endpoints

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

type modelDefinition map[string]json.RawMessage

// A DecodeModelOptions are the options for how the endpoints model definition
// are decoded.
type DecodeModelOptions struct ***REMOVED***
	SkipCustomizations bool
***REMOVED***

// Set combines all of the option functions together.
func (d *DecodeModelOptions) Set(optFns ...func(*DecodeModelOptions)) ***REMOVED***
	for _, fn := range optFns ***REMOVED***
		fn(d)
	***REMOVED***
***REMOVED***

// DecodeModel unmarshals a Regions and Endpoint model definition file into
// a endpoint Resolver. If the file format is not supported, or an error occurs
// when unmarshaling the model an error will be returned.
//
// Casting the return value of this func to a EnumPartitions will
// allow you to get a list of the partitions in the order the endpoints
// will be resolved in.
//
//    resolver, err := endpoints.DecodeModel(reader)
//
//    partitions := resolver.(endpoints.EnumPartitions).Partitions()
//    for _, p := range partitions ***REMOVED***
//        // ... inspect partitions
//***REMOVED***
func DecodeModel(r io.Reader, optFns ...func(*DecodeModelOptions)) (Resolver, error) ***REMOVED***
	var opts DecodeModelOptions
	opts.Set(optFns...)

	// Get the version of the partition file to determine what
	// unmarshaling model to use.
	modelDef := modelDefinition***REMOVED******REMOVED***
	if err := json.NewDecoder(r).Decode(&modelDef); err != nil ***REMOVED***
		return nil, newDecodeModelError("failed to decode endpoints model", err)
	***REMOVED***

	var version string
	if b, ok := modelDef["version"]; ok ***REMOVED***
		version = string(b)
	***REMOVED*** else ***REMOVED***
		return nil, newDecodeModelError("endpoints version not found in model", nil)
	***REMOVED***

	if version == "3" ***REMOVED***
		return decodeV3Endpoints(modelDef, opts)
	***REMOVED***

	return nil, newDecodeModelError(
		fmt.Sprintf("endpoints version %s, not supported", version), nil)
***REMOVED***

func decodeV3Endpoints(modelDef modelDefinition, opts DecodeModelOptions) (Resolver, error) ***REMOVED***
	b, ok := modelDef["partitions"]
	if !ok ***REMOVED***
		return nil, newDecodeModelError("endpoints model missing partitions", nil)
	***REMOVED***

	ps := partitions***REMOVED******REMOVED***
	if err := json.Unmarshal(b, &ps); err != nil ***REMOVED***
		return nil, newDecodeModelError("failed to decode endpoints model", err)
	***REMOVED***

	if opts.SkipCustomizations ***REMOVED***
		return ps, nil
	***REMOVED***

	// Customization
	for i := 0; i < len(ps); i++ ***REMOVED***
		p := &ps[i]
		custAddEC2Metadata(p)
		custAddS3DualStack(p)
		custRmIotDataService(p)
	***REMOVED***

	return ps, nil
***REMOVED***

func custAddS3DualStack(p *partition) ***REMOVED***
	if p.ID != "aws" ***REMOVED***
		return
	***REMOVED***

	s, ok := p.Services["s3"]
	if !ok ***REMOVED***
		return
	***REMOVED***

	s.Defaults.HasDualStack = boxedTrue
	s.Defaults.DualStackHostname = "***REMOVED***service***REMOVED***.dualstack.***REMOVED***region***REMOVED***.***REMOVED***dnsSuffix***REMOVED***"

	p.Services["s3"] = s
***REMOVED***

func custAddEC2Metadata(p *partition) ***REMOVED***
	p.Services["ec2metadata"] = service***REMOVED***
		IsRegionalized:    boxedFalse,
		PartitionEndpoint: "aws-global",
		Endpoints: endpoints***REMOVED***
			"aws-global": endpoint***REMOVED***
				Hostname:  "169.254.169.254/latest",
				Protocols: []string***REMOVED***"http"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func custRmIotDataService(p *partition) ***REMOVED***
	delete(p.Services, "data.iot")
***REMOVED***

type decodeModelError struct ***REMOVED***
	awsError
***REMOVED***

func newDecodeModelError(msg string, err error) decodeModelError ***REMOVED***
	return decodeModelError***REMOVED***
		awsError: awserr.New("DecodeEndpointsModelError", msg, err),
	***REMOVED***
***REMOVED***
