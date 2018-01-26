// +build windows

package winterm

const (
	horizontal = iota
	vertical
)

func (h *windowsAnsiEventHandler) getCursorWindow(info *CONSOLE_SCREEN_BUFFER_INFO) SMALL_RECT ***REMOVED***
	if h.originMode ***REMOVED***
		sr := h.effectiveSr(info.Window)
		return SMALL_RECT***REMOVED***
			Top:    sr.top,
			Bottom: sr.bottom,
			Left:   0,
			Right:  info.Size.X - 1,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return SMALL_RECT***REMOVED***
			Top:    info.Window.Top,
			Bottom: info.Window.Bottom,
			Left:   0,
			Right:  info.Size.X - 1,
		***REMOVED***
	***REMOVED***
***REMOVED***

// setCursorPosition sets the cursor to the specified position, bounded to the screen size
func (h *windowsAnsiEventHandler) setCursorPosition(position COORD, window SMALL_RECT) error ***REMOVED***
	position.X = ensureInRange(position.X, window.Left, window.Right)
	position.Y = ensureInRange(position.Y, window.Top, window.Bottom)
	err := SetConsoleCursorPosition(h.fd, position)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("Cursor position set: (%d, %d)", position.X, position.Y)
	return err
***REMOVED***

func (h *windowsAnsiEventHandler) moveCursorVertical(param int) error ***REMOVED***
	return h.moveCursor(vertical, param)
***REMOVED***

func (h *windowsAnsiEventHandler) moveCursorHorizontal(param int) error ***REMOVED***
	return h.moveCursor(horizontal, param)
***REMOVED***

func (h *windowsAnsiEventHandler) moveCursor(moveMode int, param int) error ***REMOVED***
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	position := info.CursorPosition
	switch moveMode ***REMOVED***
	case horizontal:
		position.X += int16(param)
	case vertical:
		position.Y += int16(param)
	***REMOVED***

	if err = h.setCursorPosition(position, h.getCursorWindow(info)); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) moveCursorLine(param int) error ***REMOVED***
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	position := info.CursorPosition
	position.X = 0
	position.Y += int16(param)

	if err = h.setCursorPosition(position, h.getCursorWindow(info)); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) moveCursorColumn(param int) error ***REMOVED***
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	position := info.CursorPosition
	position.X = int16(param) - 1

	if err = h.setCursorPosition(position, h.getCursorWindow(info)); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
