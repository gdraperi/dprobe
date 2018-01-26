// Package streamformatter provides helper functions to format a stream.
package streamformatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/progress"
)

const streamNewline = "\r\n"

type jsonProgressFormatter struct***REMOVED******REMOVED***

func appendNewline(source []byte) []byte ***REMOVED***
	return append(source, []byte(streamNewline)...)
***REMOVED***

// FormatStatus formats the specified objects according to the specified format (and id).
func FormatStatus(id, format string, a ...interface***REMOVED******REMOVED***) []byte ***REMOVED***
	str := fmt.Sprintf(format, a...)
	b, err := json.Marshal(&jsonmessage.JSONMessage***REMOVED***ID: id, Status: str***REMOVED***)
	if err != nil ***REMOVED***
		return FormatError(err)
	***REMOVED***
	return appendNewline(b)
***REMOVED***

// FormatError formats the error as a JSON object
func FormatError(err error) []byte ***REMOVED***
	jsonError, ok := err.(*jsonmessage.JSONError)
	if !ok ***REMOVED***
		jsonError = &jsonmessage.JSONError***REMOVED***Message: err.Error()***REMOVED***
	***REMOVED***
	if b, err := json.Marshal(&jsonmessage.JSONMessage***REMOVED***Error: jsonError, ErrorMessage: err.Error()***REMOVED***); err == nil ***REMOVED***
		return appendNewline(b)
	***REMOVED***
	return []byte(`***REMOVED***"error":"format error"***REMOVED***` + streamNewline)
***REMOVED***

func (sf *jsonProgressFormatter) formatStatus(id, format string, a ...interface***REMOVED******REMOVED***) []byte ***REMOVED***
	return FormatStatus(id, format, a...)
***REMOVED***

// formatProgress formats the progress information for a specified action.
func (sf *jsonProgressFormatter) formatProgress(id, action string, progress *jsonmessage.JSONProgress, aux interface***REMOVED******REMOVED***) []byte ***REMOVED***
	if progress == nil ***REMOVED***
		progress = &jsonmessage.JSONProgress***REMOVED******REMOVED***
	***REMOVED***
	var auxJSON *json.RawMessage
	if aux != nil ***REMOVED***
		auxJSONBytes, err := json.Marshal(aux)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		auxJSON = new(json.RawMessage)
		*auxJSON = auxJSONBytes
	***REMOVED***
	b, err := json.Marshal(&jsonmessage.JSONMessage***REMOVED***
		Status:          action,
		ProgressMessage: progress.String(),
		Progress:        progress,
		ID:              id,
		Aux:             auxJSON,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return appendNewline(b)
***REMOVED***

type rawProgressFormatter struct***REMOVED******REMOVED***

func (sf *rawProgressFormatter) formatStatus(id, format string, a ...interface***REMOVED******REMOVED***) []byte ***REMOVED***
	return []byte(fmt.Sprintf(format, a...) + streamNewline)
***REMOVED***

func (sf *rawProgressFormatter) formatProgress(id, action string, progress *jsonmessage.JSONProgress, aux interface***REMOVED******REMOVED***) []byte ***REMOVED***
	if progress == nil ***REMOVED***
		progress = &jsonmessage.JSONProgress***REMOVED******REMOVED***
	***REMOVED***
	endl := "\r"
	if progress.String() == "" ***REMOVED***
		endl += "\n"
	***REMOVED***
	return []byte(action + " " + progress.String() + endl)
***REMOVED***

// NewProgressOutput returns a progress.Output object that can be passed to
// progress.NewProgressReader.
func NewProgressOutput(out io.Writer) progress.Output ***REMOVED***
	return &progressOutput***REMOVED***sf: &rawProgressFormatter***REMOVED******REMOVED***, out: out, newLines: true***REMOVED***
***REMOVED***

// NewJSONProgressOutput returns a progress.Output that that formats output
// using JSON objects
func NewJSONProgressOutput(out io.Writer, newLines bool) progress.Output ***REMOVED***
	return &progressOutput***REMOVED***sf: &jsonProgressFormatter***REMOVED******REMOVED***, out: out, newLines: newLines***REMOVED***
***REMOVED***

type formatProgress interface ***REMOVED***
	formatStatus(id, format string, a ...interface***REMOVED******REMOVED***) []byte
	formatProgress(id, action string, progress *jsonmessage.JSONProgress, aux interface***REMOVED******REMOVED***) []byte
***REMOVED***

type progressOutput struct ***REMOVED***
	sf       formatProgress
	out      io.Writer
	newLines bool
***REMOVED***

// WriteProgress formats progress information from a ProgressReader.
func (out *progressOutput) WriteProgress(prog progress.Progress) error ***REMOVED***
	var formatted []byte
	if prog.Message != "" ***REMOVED***
		formatted = out.sf.formatStatus(prog.ID, prog.Message)
	***REMOVED*** else ***REMOVED***
		jsonProgress := jsonmessage.JSONProgress***REMOVED***Current: prog.Current, Total: prog.Total, HideCounts: prog.HideCounts, Units: prog.Units***REMOVED***
		formatted = out.sf.formatProgress(prog.ID, prog.Action, &jsonProgress, prog.Aux)
	***REMOVED***
	_, err := out.out.Write(formatted)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if out.newLines && prog.LastUpdate ***REMOVED***
		_, err = out.out.Write(out.sf.formatStatus("", ""))
		return err
	***REMOVED***

	return nil
***REMOVED***

// AuxFormatter is a streamFormatter that writes aux progress messages
type AuxFormatter struct ***REMOVED***
	io.Writer
***REMOVED***

// Emit emits the given interface as an aux progress message
func (sf *AuxFormatter) Emit(aux interface***REMOVED******REMOVED***) error ***REMOVED***
	auxJSONBytes, err := json.Marshal(aux)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	auxJSON := new(json.RawMessage)
	*auxJSON = auxJSONBytes
	msgJSON, err := json.Marshal(&jsonmessage.JSONMessage***REMOVED***Aux: auxJSON***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	msgJSON = appendNewline(msgJSON)
	n, err := sf.Writer.Write(msgJSON)
	if n != len(msgJSON) ***REMOVED***
		return io.ErrShortWrite
	***REMOVED***
	return err
***REMOVED***
