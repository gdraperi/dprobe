// +build windows

package winterm

import "github.com/Azure/go-ansiterm"

func (h *windowsAnsiEventHandler) clearRange(attributes uint16, fromCoord COORD, toCoord COORD) error ***REMOVED***
	// Ignore an invalid (negative area) request
	if toCoord.Y < fromCoord.Y ***REMOVED***
		return nil
	***REMOVED***

	var err error

	var coordStart = COORD***REMOVED******REMOVED***
	var coordEnd = COORD***REMOVED******REMOVED***

	xCurrent, yCurrent := fromCoord.X, fromCoord.Y
	xEnd, yEnd := toCoord.X, toCoord.Y

	// Clear any partial initial line
	if xCurrent > 0 ***REMOVED***
		coordStart.X, coordStart.Y = xCurrent, yCurrent
		coordEnd.X, coordEnd.Y = xEnd, yCurrent

		err = h.clearRect(attributes, coordStart, coordEnd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		xCurrent = 0
		yCurrent += 1
	***REMOVED***

	// Clear intervening rectangular section
	if yCurrent < yEnd ***REMOVED***
		coordStart.X, coordStart.Y = xCurrent, yCurrent
		coordEnd.X, coordEnd.Y = xEnd, yEnd-1

		err = h.clearRect(attributes, coordStart, coordEnd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		xCurrent = 0
		yCurrent = yEnd
	***REMOVED***

	// Clear remaining partial ending line
	coordStart.X, coordStart.Y = xCurrent, yCurrent
	coordEnd.X, coordEnd.Y = xEnd, yEnd

	err = h.clearRect(attributes, coordStart, coordEnd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) clearRect(attributes uint16, fromCoord COORD, toCoord COORD) error ***REMOVED***
	region := SMALL_RECT***REMOVED***Top: fromCoord.Y, Left: fromCoord.X, Bottom: toCoord.Y, Right: toCoord.X***REMOVED***
	width := toCoord.X - fromCoord.X + 1
	height := toCoord.Y - fromCoord.Y + 1
	size := uint32(width) * uint32(height)

	if size <= 0 ***REMOVED***
		return nil
	***REMOVED***

	buffer := make([]CHAR_INFO, size)

	char := CHAR_INFO***REMOVED***ansiterm.FILL_CHARACTER, attributes***REMOVED***
	for i := 0; i < int(size); i++ ***REMOVED***
		buffer[i] = char
	***REMOVED***

	err := WriteConsoleOutput(h.fd, buffer, COORD***REMOVED***X: width, Y: height***REMOVED***, COORD***REMOVED***X: 0, Y: 0***REMOVED***, &region)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
