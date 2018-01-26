package ansiterm

import (
	"strconv"
)

func parseParams(bytes []byte) ([]string, error) ***REMOVED***
	paramBuff := make([]byte, 0, 0)
	params := []string***REMOVED******REMOVED***

	for _, v := range bytes ***REMOVED***
		if v == ';' ***REMOVED***
			if len(paramBuff) > 0 ***REMOVED***
				// Completed parameter, append it to the list
				s := string(paramBuff)
				params = append(params, s)
				paramBuff = make([]byte, 0, 0)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			paramBuff = append(paramBuff, v)
		***REMOVED***
	***REMOVED***

	// Last parameter may not be terminated with ';'
	if len(paramBuff) > 0 ***REMOVED***
		s := string(paramBuff)
		params = append(params, s)
	***REMOVED***

	return params, nil
***REMOVED***

func parseCmd(context ansiContext) (string, error) ***REMOVED***
	return string(context.currentChar), nil
***REMOVED***

func getInt(params []string, dflt int) int ***REMOVED***
	i := getInts(params, 1, dflt)[0]
	return i
***REMOVED***

func getInts(params []string, minCount int, dflt int) []int ***REMOVED***
	ints := []int***REMOVED******REMOVED***

	for _, v := range params ***REMOVED***
		i, _ := strconv.Atoi(v)
		// Zero is mapped to the default value in VT100.
		if i == 0 ***REMOVED***
			i = dflt
		***REMOVED***
		ints = append(ints, i)
	***REMOVED***

	if len(ints) < minCount ***REMOVED***
		remaining := minCount - len(ints)
		for i := 0; i < remaining; i++ ***REMOVED***
			ints = append(ints, dflt)
		***REMOVED***
	***REMOVED***

	return ints
***REMOVED***

func (ap *AnsiParser) modeDispatch(param string, set bool) error ***REMOVED***
	switch param ***REMOVED***
	case "?3":
		return ap.eventHandler.DECCOLM(set)
	case "?6":
		return ap.eventHandler.DECOM(set)
	case "?25":
		return ap.eventHandler.DECTCEM(set)
	***REMOVED***
	return nil
***REMOVED***

func (ap *AnsiParser) hDispatch(params []string) error ***REMOVED***
	if len(params) == 1 ***REMOVED***
		return ap.modeDispatch(params[0], true)
	***REMOVED***

	return nil
***REMOVED***

func (ap *AnsiParser) lDispatch(params []string) error ***REMOVED***
	if len(params) == 1 ***REMOVED***
		return ap.modeDispatch(params[0], false)
	***REMOVED***

	return nil
***REMOVED***

func getEraseParam(params []string) int ***REMOVED***
	param := getInt(params, 0)
	if param < 0 || 3 < param ***REMOVED***
		param = 0
	***REMOVED***

	return param
***REMOVED***
