package ec2metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
)

// GetMetadata uses the path provided to request information from the EC2
// instance metdata service. The content will be returned as a string, or
// error if the request failed.
func (c *EC2Metadata) GetMetadata(p string) (string, error) ***REMOVED***
	op := &request.Operation***REMOVED***
		Name:       "GetMetadata",
		HTTPMethod: "GET",
		HTTPPath:   path.Join("/", "meta-data", p),
	***REMOVED***

	output := &metadataOutput***REMOVED******REMOVED***
	req := c.NewRequest(op, nil, output)

	return output.Content, req.Send()
***REMOVED***

// GetUserData returns the userdata that was configured for the service. If
// there is no user-data setup for the EC2 instance a "NotFoundError" error
// code will be returned.
func (c *EC2Metadata) GetUserData() (string, error) ***REMOVED***
	op := &request.Operation***REMOVED***
		Name:       "GetUserData",
		HTTPMethod: "GET",
		HTTPPath:   path.Join("/", "user-data"),
	***REMOVED***

	output := &metadataOutput***REMOVED******REMOVED***
	req := c.NewRequest(op, nil, output)
	req.Handlers.UnmarshalError.PushBack(func(r *request.Request) ***REMOVED***
		if r.HTTPResponse.StatusCode == http.StatusNotFound ***REMOVED***
			r.Error = awserr.New("NotFoundError", "user-data not found", r.Error)
		***REMOVED***
	***REMOVED***)

	return output.Content, req.Send()
***REMOVED***

// GetDynamicData uses the path provided to request information from the EC2
// instance metadata service for dynamic data. The content will be returned
// as a string, or error if the request failed.
func (c *EC2Metadata) GetDynamicData(p string) (string, error) ***REMOVED***
	op := &request.Operation***REMOVED***
		Name:       "GetDynamicData",
		HTTPMethod: "GET",
		HTTPPath:   path.Join("/", "dynamic", p),
	***REMOVED***

	output := &metadataOutput***REMOVED******REMOVED***
	req := c.NewRequest(op, nil, output)

	return output.Content, req.Send()
***REMOVED***

// GetInstanceIdentityDocument retrieves an identity document describing an
// instance. Error is returned if the request fails or is unable to parse
// the response.
func (c *EC2Metadata) GetInstanceIdentityDocument() (EC2InstanceIdentityDocument, error) ***REMOVED***
	resp, err := c.GetDynamicData("instance-identity/document")
	if err != nil ***REMOVED***
		return EC2InstanceIdentityDocument***REMOVED******REMOVED***,
			awserr.New("EC2MetadataRequestError",
				"failed to get EC2 instance identity document", err)
	***REMOVED***

	doc := EC2InstanceIdentityDocument***REMOVED******REMOVED***
	if err := json.NewDecoder(strings.NewReader(resp)).Decode(&doc); err != nil ***REMOVED***
		return EC2InstanceIdentityDocument***REMOVED******REMOVED***,
			awserr.New("SerializationError",
				"failed to decode EC2 instance identity document", err)
	***REMOVED***

	return doc, nil
***REMOVED***

// IAMInfo retrieves IAM info from the metadata API
func (c *EC2Metadata) IAMInfo() (EC2IAMInfo, error) ***REMOVED***
	resp, err := c.GetMetadata("iam/info")
	if err != nil ***REMOVED***
		return EC2IAMInfo***REMOVED******REMOVED***,
			awserr.New("EC2MetadataRequestError",
				"failed to get EC2 IAM info", err)
	***REMOVED***

	info := EC2IAMInfo***REMOVED******REMOVED***
	if err := json.NewDecoder(strings.NewReader(resp)).Decode(&info); err != nil ***REMOVED***
		return EC2IAMInfo***REMOVED******REMOVED***,
			awserr.New("SerializationError",
				"failed to decode EC2 IAM info", err)
	***REMOVED***

	if info.Code != "Success" ***REMOVED***
		errMsg := fmt.Sprintf("failed to get EC2 IAM Info (%s)", info.Code)
		return EC2IAMInfo***REMOVED******REMOVED***,
			awserr.New("EC2MetadataError", errMsg, nil)
	***REMOVED***

	return info, nil
***REMOVED***

// Region returns the region the instance is running in.
func (c *EC2Metadata) Region() (string, error) ***REMOVED***
	resp, err := c.GetMetadata("placement/availability-zone")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// returns region without the suffix. Eg: us-west-2a becomes us-west-2
	return resp[:len(resp)-1], nil
***REMOVED***

// Available returns if the application has access to the EC2 Metadata service.
// Can be used to determine if application is running within an EC2 Instance and
// the metadata service is available.
func (c *EC2Metadata) Available() bool ***REMOVED***
	if _, err := c.GetMetadata("instance-id"); err != nil ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

// An EC2IAMInfo provides the shape for unmarshaling
// an IAM info from the metadata API
type EC2IAMInfo struct ***REMOVED***
	Code               string
	LastUpdated        time.Time
	InstanceProfileArn string
	InstanceProfileID  string
***REMOVED***

// An EC2InstanceIdentityDocument provides the shape for unmarshaling
// an instance identity document
type EC2InstanceIdentityDocument struct ***REMOVED***
	DevpayProductCodes []string  `json:"devpayProductCodes"`
	AvailabilityZone   string    `json:"availabilityZone"`
	PrivateIP          string    `json:"privateIp"`
	Version            string    `json:"version"`
	Region             string    `json:"region"`
	InstanceID         string    `json:"instanceId"`
	BillingProducts    []string  `json:"billingProducts"`
	InstanceType       string    `json:"instanceType"`
	AccountID          string    `json:"accountId"`
	PendingTime        time.Time `json:"pendingTime"`
	ImageID            string    `json:"imageId"`
	KernelID           string    `json:"kernelId"`
	RamdiskID          string    `json:"ramdiskId"`
	Architecture       string    `json:"architecture"`
***REMOVED***
