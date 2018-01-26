package request

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

const (
	// InvalidParameterErrCode is the error code for invalid parameters errors
	InvalidParameterErrCode = "InvalidParameter"
	// ParamRequiredErrCode is the error code for required parameter errors
	ParamRequiredErrCode = "ParamRequiredError"
	// ParamMinValueErrCode is the error code for fields with too low of a
	// number value.
	ParamMinValueErrCode = "ParamMinValueError"
	// ParamMinLenErrCode is the error code for fields without enough elements.
	ParamMinLenErrCode = "ParamMinLenError"
)

// Validator provides a way for types to perform validation logic on their
// input values that external code can use to determine if a type's values
// are valid.
type Validator interface ***REMOVED***
	Validate() error
***REMOVED***

// An ErrInvalidParams provides wrapping of invalid parameter errors found when
// validating API operation input parameters.
type ErrInvalidParams struct ***REMOVED***
	// Context is the base context of the invalid parameter group.
	Context string
	errs    []ErrInvalidParam
***REMOVED***

// Add adds a new invalid parameter error to the collection of invalid
// parameters. The context of the invalid parameter will be updated to reflect
// this collection.
func (e *ErrInvalidParams) Add(err ErrInvalidParam) ***REMOVED***
	err.SetContext(e.Context)
	e.errs = append(e.errs, err)
***REMOVED***

// AddNested adds the invalid parameter errors from another ErrInvalidParams
// value into this collection. The nested errors will have their nested context
// updated and base context to reflect the merging.
//
// Use for nested validations errors.
func (e *ErrInvalidParams) AddNested(nestedCtx string, nested ErrInvalidParams) ***REMOVED***
	for _, err := range nested.errs ***REMOVED***
		err.SetContext(e.Context)
		err.AddNestedContext(nestedCtx)
		e.errs = append(e.errs, err)
	***REMOVED***
***REMOVED***

// Len returns the number of invalid parameter errors
func (e ErrInvalidParams) Len() int ***REMOVED***
	return len(e.errs)
***REMOVED***

// Code returns the code of the error
func (e ErrInvalidParams) Code() string ***REMOVED***
	return InvalidParameterErrCode
***REMOVED***

// Message returns the message of the error
func (e ErrInvalidParams) Message() string ***REMOVED***
	return fmt.Sprintf("%d validation error(s) found.", len(e.errs))
***REMOVED***

// Error returns the string formatted form of the invalid parameters.
func (e ErrInvalidParams) Error() string ***REMOVED***
	w := &bytes.Buffer***REMOVED******REMOVED***
	fmt.Fprintf(w, "%s: %s\n", e.Code(), e.Message())

	for _, err := range e.errs ***REMOVED***
		fmt.Fprintf(w, "- %s\n", err.Message())
	***REMOVED***

	return w.String()
***REMOVED***

// OrigErr returns the invalid parameters as a awserr.BatchedErrors value
func (e ErrInvalidParams) OrigErr() error ***REMOVED***
	return awserr.NewBatchError(
		InvalidParameterErrCode, e.Message(), e.OrigErrs())
***REMOVED***

// OrigErrs returns a slice of the invalid parameters
func (e ErrInvalidParams) OrigErrs() []error ***REMOVED***
	errs := make([]error, len(e.errs))
	for i := 0; i < len(errs); i++ ***REMOVED***
		errs[i] = e.errs[i]
	***REMOVED***

	return errs
***REMOVED***

// An ErrInvalidParam represents an invalid parameter error type.
type ErrInvalidParam interface ***REMOVED***
	awserr.Error

	// Field name the error occurred on.
	Field() string

	// SetContext updates the context of the error.
	SetContext(string)

	// AddNestedContext updates the error's context to include a nested level.
	AddNestedContext(string)
***REMOVED***

type errInvalidParam struct ***REMOVED***
	context       string
	nestedContext string
	field         string
	code          string
	msg           string
***REMOVED***

// Code returns the error code for the type of invalid parameter.
func (e *errInvalidParam) Code() string ***REMOVED***
	return e.code
***REMOVED***

// Message returns the reason the parameter was invalid, and its context.
func (e *errInvalidParam) Message() string ***REMOVED***
	return fmt.Sprintf("%s, %s.", e.msg, e.Field())
***REMOVED***

// Error returns the string version of the invalid parameter error.
func (e *errInvalidParam) Error() string ***REMOVED***
	return fmt.Sprintf("%s: %s", e.code, e.Message())
***REMOVED***

// OrigErr returns nil, Implemented for awserr.Error interface.
func (e *errInvalidParam) OrigErr() error ***REMOVED***
	return nil
***REMOVED***

// Field Returns the field and context the error occurred.
func (e *errInvalidParam) Field() string ***REMOVED***
	field := e.context
	if len(field) > 0 ***REMOVED***
		field += "."
	***REMOVED***
	if len(e.nestedContext) > 0 ***REMOVED***
		field += fmt.Sprintf("%s.", e.nestedContext)
	***REMOVED***
	field += e.field

	return field
***REMOVED***

// SetContext updates the base context of the error.
func (e *errInvalidParam) SetContext(ctx string) ***REMOVED***
	e.context = ctx
***REMOVED***

// AddNestedContext prepends a context to the field's path.
func (e *errInvalidParam) AddNestedContext(ctx string) ***REMOVED***
	if len(e.nestedContext) == 0 ***REMOVED***
		e.nestedContext = ctx
	***REMOVED*** else ***REMOVED***
		e.nestedContext = fmt.Sprintf("%s.%s", ctx, e.nestedContext)
	***REMOVED***

***REMOVED***

// An ErrParamRequired represents an required parameter error.
type ErrParamRequired struct ***REMOVED***
	errInvalidParam
***REMOVED***

// NewErrParamRequired creates a new required parameter error.
func NewErrParamRequired(field string) *ErrParamRequired ***REMOVED***
	return &ErrParamRequired***REMOVED***
		errInvalidParam***REMOVED***
			code:  ParamRequiredErrCode,
			field: field,
			msg:   fmt.Sprintf("missing required field"),
		***REMOVED***,
	***REMOVED***
***REMOVED***

// An ErrParamMinValue represents a minimum value parameter error.
type ErrParamMinValue struct ***REMOVED***
	errInvalidParam
	min float64
***REMOVED***

// NewErrParamMinValue creates a new minimum value parameter error.
func NewErrParamMinValue(field string, min float64) *ErrParamMinValue ***REMOVED***
	return &ErrParamMinValue***REMOVED***
		errInvalidParam: errInvalidParam***REMOVED***
			code:  ParamMinValueErrCode,
			field: field,
			msg:   fmt.Sprintf("minimum field value of %v", min),
		***REMOVED***,
		min: min,
	***REMOVED***
***REMOVED***

// MinValue returns the field's require minimum value.
//
// float64 is returned for both int and float min values.
func (e *ErrParamMinValue) MinValue() float64 ***REMOVED***
	return e.min
***REMOVED***

// An ErrParamMinLen represents a minimum length parameter error.
type ErrParamMinLen struct ***REMOVED***
	errInvalidParam
	min int
***REMOVED***

// NewErrParamMinLen creates a new minimum length parameter error.
func NewErrParamMinLen(field string, min int) *ErrParamMinLen ***REMOVED***
	return &ErrParamMinLen***REMOVED***
		errInvalidParam: errInvalidParam***REMOVED***
			code:  ParamMinLenErrCode,
			field: field,
			msg:   fmt.Sprintf("minimum field size of %v", min),
		***REMOVED***,
		min: min,
	***REMOVED***
***REMOVED***

// MinLen returns the field's required minimum length.
func (e *ErrParamMinLen) MinLen() int ***REMOVED***
	return e.min
***REMOVED***
