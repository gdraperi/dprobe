// +build windows

package winterm

// effectiveSr gets the current effective scroll region in buffer coordinates
func (h *windowsAnsiEventHandler) effectiveSr(window SMALL_RECT) scrollRegion ***REMOVED***
	top := addInRange(window.Top, h.sr.top, window.Top, window.Bottom)
	bottom := addInRange(window.Top, h.sr.bottom, window.Top, window.Bottom)
	if top >= bottom ***REMOVED***
		top = window.Top
		bottom = window.Bottom
	***REMOVED***
	return scrollRegion***REMOVED***top: top, bottom: bottom***REMOVED***
***REMOVED***

func (h *windowsAnsiEventHandler) scrollUp(param int) error ***REMOVED***
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	sr := h.effectiveSr(info.Window)
	return h.scroll(param, sr, info)
***REMOVED***

func (h *windowsAnsiEventHandler) scrollDown(param int) error ***REMOVED***
	return h.scrollUp(-param)
***REMOVED***

func (h *windowsAnsiEventHandler) deleteLines(param int) error ***REMOVED***
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	start := info.CursorPosition.Y
	sr := h.effectiveSr(info.Window)
	// Lines cannot be inserted or deleted outside the scrolling region.
	if start >= sr.top && start <= sr.bottom ***REMOVED***
		sr.top = start
		return h.scroll(param, sr, info)
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

func (h *windowsAnsiEventHandler) insertLines(param int) error ***REMOVED***
	return h.deleteLines(-param)
***REMOVED***

// scroll scrolls the provided scroll region by param lines. The scroll region is in buffer coordinates.
func (h *windowsAnsiEventHandler) scroll(param int, sr scrollRegion, info *CONSOLE_SCREEN_BUFFER_INFO) error ***REMOVED***
	h.logf("scroll: scrollTop: %d, scrollBottom: %d", sr.top, sr.bottom)
	h.logf("scroll: windowTop: %d, windowBottom: %d", info.Window.Top, info.Window.Bottom)

	// Copy from and clip to the scroll region (full buffer width)
	scrollRect := SMALL_RECT***REMOVED***
		Top:    sr.top,
		Bottom: sr.bottom,
		Left:   0,
		Right:  info.Size.X - 1,
	***REMOVED***

	// Origin to which area should be copied
	destOrigin := COORD***REMOVED***
		X: 0,
		Y: sr.top - int16(param),
	***REMOVED***

	char := CHAR_INFO***REMOVED***
		UnicodeChar: ' ',
		Attributes:  h.attributes,
	***REMOVED***

	if err := ScrollConsoleScreenBuffer(h.fd, scrollRect, scrollRect, destOrigin, char); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) deleteCharacters(param int) error ***REMOVED***
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return h.scrollLine(param, info.CursorPosition, info)
***REMOVED***

func (h *windowsAnsiEventHandler) insertCharacters(param int) error ***REMOVED***
	return h.deleteCharacters(-param)
***REMOVED***

// scrollLine scrolls a line horizontally starting at the provided position by a number of columns.
func (h *windowsAnsiEventHandler) scrollLine(columns int, position COORD, info *CONSOLE_SCREEN_BUFFER_INFO) error ***REMOVED***
	// Copy from and clip to the scroll region (full buffer width)
	scrollRect := SMALL_RECT***REMOVED***
		Top:    position.Y,
		Bottom: position.Y,
		Left:   position.X,
		Right:  info.Size.X - 1,
	***REMOVED***

	// Origin to which area should be copied
	destOrigin := COORD***REMOVED***
		X: position.X - int16(columns),
		Y: position.Y,
	***REMOVED***

	char := CHAR_INFO***REMOVED***
		UnicodeChar: ' ',
		Attributes:  h.attributes,
	***REMOVED***

	if err := ScrollConsoleScreenBuffer(h.fd, scrollRect, scrollRect, destOrigin, char); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
