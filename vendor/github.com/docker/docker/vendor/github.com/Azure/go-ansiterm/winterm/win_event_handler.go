// +build windows

package winterm

import (
	"bytes"
	"log"
	"os"
	"strconv"

	"github.com/Azure/go-ansiterm"
)

type windowsAnsiEventHandler struct ***REMOVED***
	fd             uintptr
	file           *os.File
	infoReset      *CONSOLE_SCREEN_BUFFER_INFO
	sr             scrollRegion
	buffer         bytes.Buffer
	attributes     uint16
	inverted       bool
	wrapNext       bool
	drewMarginByte bool
	originMode     bool
	marginByte     byte
	curInfo        *CONSOLE_SCREEN_BUFFER_INFO
	curPos         COORD
	logf           func(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type Option func(*windowsAnsiEventHandler)

func WithLogf(f func(string, ...interface***REMOVED******REMOVED***)) Option ***REMOVED***
	return func(w *windowsAnsiEventHandler) ***REMOVED***
		w.logf = f
	***REMOVED***
***REMOVED***

func CreateWinEventHandler(fd uintptr, file *os.File, opts ...Option) ansiterm.AnsiEventHandler ***REMOVED***
	infoReset, err := GetConsoleScreenBufferInfo(fd)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	h := &windowsAnsiEventHandler***REMOVED***
		fd:         fd,
		file:       file,
		infoReset:  infoReset,
		attributes: infoReset.Attributes,
	***REMOVED***
	for _, o := range opts ***REMOVED***
		o(h)
	***REMOVED***

	if isDebugEnv := os.Getenv(ansiterm.LogEnv); isDebugEnv == "1" ***REMOVED***
		logFile, _ := os.Create("winEventHandler.log")
		logger := log.New(logFile, "", log.LstdFlags)
		if h.logf != nil ***REMOVED***
			l := h.logf
			h.logf = func(s string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
				l(s, v...)
				logger.Printf(s, v...)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			h.logf = logger.Printf
		***REMOVED***
	***REMOVED***

	if h.logf == nil ***REMOVED***
		h.logf = func(string, ...interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***
	***REMOVED***

	return h
***REMOVED***

type scrollRegion struct ***REMOVED***
	top    int16
	bottom int16
***REMOVED***

// simulateLF simulates a LF or CR+LF by scrolling if necessary to handle the
// current cursor position and scroll region settings, in which case it returns
// true. If no special handling is necessary, then it does nothing and returns
// false.
//
// In the false case, the caller should ensure that a carriage return
// and line feed are inserted or that the text is otherwise wrapped.
func (h *windowsAnsiEventHandler) simulateLF(includeCR bool) (bool, error) ***REMOVED***
	if h.wrapNext ***REMOVED***
		if err := h.Flush(); err != nil ***REMOVED***
			return false, err
		***REMOVED***
		h.clearWrap()
	***REMOVED***
	pos, info, err := h.getCurrentInfo()
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	sr := h.effectiveSr(info.Window)
	if pos.Y == sr.bottom ***REMOVED***
		// Scrolling is necessary. Let Windows automatically scroll if the scrolling region
		// is the full window.
		if sr.top == info.Window.Top && sr.bottom == info.Window.Bottom ***REMOVED***
			if includeCR ***REMOVED***
				pos.X = 0
				h.updatePos(pos)
			***REMOVED***
			return false, nil
		***REMOVED***

		// A custom scroll region is active. Scroll the window manually to simulate
		// the LF.
		if err := h.Flush(); err != nil ***REMOVED***
			return false, err
		***REMOVED***
		h.logf("Simulating LF inside scroll region")
		if err := h.scrollUp(1); err != nil ***REMOVED***
			return false, err
		***REMOVED***
		if includeCR ***REMOVED***
			pos.X = 0
			if err := SetConsoleCursorPosition(h.fd, pos); err != nil ***REMOVED***
				return false, err
			***REMOVED***
		***REMOVED***
		return true, nil

	***REMOVED*** else if pos.Y < info.Window.Bottom ***REMOVED***
		// Let Windows handle the LF.
		pos.Y++
		if includeCR ***REMOVED***
			pos.X = 0
		***REMOVED***
		h.updatePos(pos)
		return false, nil
	***REMOVED*** else ***REMOVED***
		// The cursor is at the bottom of the screen but outside the scroll
		// region. Skip the LF.
		h.logf("Simulating LF outside scroll region")
		if includeCR ***REMOVED***
			if err := h.Flush(); err != nil ***REMOVED***
				return false, err
			***REMOVED***
			pos.X = 0
			if err := SetConsoleCursorPosition(h.fd, pos); err != nil ***REMOVED***
				return false, err
			***REMOVED***
		***REMOVED***
		return true, nil
	***REMOVED***
***REMOVED***

// executeLF executes a LF without a CR.
func (h *windowsAnsiEventHandler) executeLF() error ***REMOVED***
	handled, err := h.simulateLF(false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !handled ***REMOVED***
		// Windows LF will reset the cursor column position. Write the LF
		// and restore the cursor position.
		pos, _, err := h.getCurrentInfo()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		h.buffer.WriteByte(ansiterm.ANSI_LINE_FEED)
		if pos.X != 0 ***REMOVED***
			if err := h.Flush(); err != nil ***REMOVED***
				return err
			***REMOVED***
			h.logf("Resetting cursor position for LF without CR")
			if err := SetConsoleCursorPosition(h.fd, pos); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) Print(b byte) error ***REMOVED***
	if h.wrapNext ***REMOVED***
		h.buffer.WriteByte(h.marginByte)
		h.clearWrap()
		if _, err := h.simulateLF(true); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	pos, info, err := h.getCurrentInfo()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if pos.X == info.Size.X-1 ***REMOVED***
		h.wrapNext = true
		h.marginByte = b
	***REMOVED*** else ***REMOVED***
		pos.X++
		h.updatePos(pos)
		h.buffer.WriteByte(b)
	***REMOVED***
	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) Execute(b byte) error ***REMOVED***
	switch b ***REMOVED***
	case ansiterm.ANSI_TAB:
		h.logf("Execute(TAB)")
		// Move to the next tab stop, but preserve auto-wrap if already set.
		if !h.wrapNext ***REMOVED***
			pos, info, err := h.getCurrentInfo()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			pos.X = (pos.X + 8) - pos.X%8
			if pos.X >= info.Size.X ***REMOVED***
				pos.X = info.Size.X - 1
			***REMOVED***
			if err := h.Flush(); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := SetConsoleCursorPosition(h.fd, pos); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil

	case ansiterm.ANSI_BEL:
		h.buffer.WriteByte(ansiterm.ANSI_BEL)
		return nil

	case ansiterm.ANSI_BACKSPACE:
		if h.wrapNext ***REMOVED***
			if err := h.Flush(); err != nil ***REMOVED***
				return err
			***REMOVED***
			h.clearWrap()
		***REMOVED***
		pos, _, err := h.getCurrentInfo()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if pos.X > 0 ***REMOVED***
			pos.X--
			h.updatePos(pos)
			h.buffer.WriteByte(ansiterm.ANSI_BACKSPACE)
		***REMOVED***
		return nil

	case ansiterm.ANSI_VERTICAL_TAB, ansiterm.ANSI_FORM_FEED:
		// Treat as true LF.
		return h.executeLF()

	case ansiterm.ANSI_LINE_FEED:
		// Simulate a CR and LF for now since there is no way in go-ansiterm
		// to tell if the LF should include CR (and more things break when it's
		// missing than when it's incorrectly added).
		handled, err := h.simulateLF(true)
		if handled || err != nil ***REMOVED***
			return err
		***REMOVED***
		return h.buffer.WriteByte(ansiterm.ANSI_LINE_FEED)

	case ansiterm.ANSI_CARRIAGE_RETURN:
		if h.wrapNext ***REMOVED***
			if err := h.Flush(); err != nil ***REMOVED***
				return err
			***REMOVED***
			h.clearWrap()
		***REMOVED***
		pos, _, err := h.getCurrentInfo()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if pos.X != 0 ***REMOVED***
			pos.X = 0
			h.updatePos(pos)
			h.buffer.WriteByte(ansiterm.ANSI_CARRIAGE_RETURN)
		***REMOVED***
		return nil

	default:
		return nil
	***REMOVED***
***REMOVED***

func (h *windowsAnsiEventHandler) CUU(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("CUU: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()
	return h.moveCursorVertical(-param)
***REMOVED***

func (h *windowsAnsiEventHandler) CUD(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("CUD: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()
	return h.moveCursorVertical(param)
***REMOVED***

func (h *windowsAnsiEventHandler) CUF(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("CUF: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()
	return h.moveCursorHorizontal(param)
***REMOVED***

func (h *windowsAnsiEventHandler) CUB(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("CUB: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()
	return h.moveCursorHorizontal(-param)
***REMOVED***

func (h *windowsAnsiEventHandler) CNL(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("CNL: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()
	return h.moveCursorLine(param)
***REMOVED***

func (h *windowsAnsiEventHandler) CPL(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("CPL: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()
	return h.moveCursorLine(-param)
***REMOVED***

func (h *windowsAnsiEventHandler) CHA(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("CHA: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()
	return h.moveCursorColumn(param)
***REMOVED***

func (h *windowsAnsiEventHandler) VPA(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("VPA: [[%d]]", param)
	h.clearWrap()
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	window := h.getCursorWindow(info)
	position := info.CursorPosition
	position.Y = window.Top + int16(param) - 1
	return h.setCursorPosition(position, window)
***REMOVED***

func (h *windowsAnsiEventHandler) CUP(row int, col int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("CUP: [[%d %d]]", row, col)
	h.clearWrap()
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	window := h.getCursorWindow(info)
	position := COORD***REMOVED***window.Left + int16(col) - 1, window.Top + int16(row) - 1***REMOVED***
	return h.setCursorPosition(position, window)
***REMOVED***

func (h *windowsAnsiEventHandler) HVP(row int, col int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("HVP: [[%d %d]]", row, col)
	h.clearWrap()
	return h.CUP(row, col)
***REMOVED***

func (h *windowsAnsiEventHandler) DECTCEM(visible bool) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("DECTCEM: [%v]", []string***REMOVED***strconv.FormatBool(visible)***REMOVED***)
	h.clearWrap()
	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) DECOM(enable bool) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("DECOM: [%v]", []string***REMOVED***strconv.FormatBool(enable)***REMOVED***)
	h.clearWrap()
	h.originMode = enable
	return h.CUP(1, 1)
***REMOVED***

func (h *windowsAnsiEventHandler) DECCOLM(use132 bool) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("DECCOLM: [%v]", []string***REMOVED***strconv.FormatBool(use132)***REMOVED***)
	h.clearWrap()
	if err := h.ED(2); err != nil ***REMOVED***
		return err
	***REMOVED***
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	targetWidth := int16(80)
	if use132 ***REMOVED***
		targetWidth = 132
	***REMOVED***
	if info.Size.X < targetWidth ***REMOVED***
		if err := SetConsoleScreenBufferSize(h.fd, COORD***REMOVED***targetWidth, info.Size.Y***REMOVED***); err != nil ***REMOVED***
			h.logf("set buffer failed: %v", err)
			return err
		***REMOVED***
	***REMOVED***
	window := info.Window
	window.Left = 0
	window.Right = targetWidth - 1
	if err := SetConsoleWindowInfo(h.fd, true, window); err != nil ***REMOVED***
		h.logf("set window failed: %v", err)
		return err
	***REMOVED***
	if info.Size.X > targetWidth ***REMOVED***
		if err := SetConsoleScreenBufferSize(h.fd, COORD***REMOVED***targetWidth, info.Size.Y***REMOVED***); err != nil ***REMOVED***
			h.logf("set buffer failed: %v", err)
			return err
		***REMOVED***
	***REMOVED***
	return SetConsoleCursorPosition(h.fd, COORD***REMOVED***0, 0***REMOVED***)
***REMOVED***

func (h *windowsAnsiEventHandler) ED(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("ED: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()

	// [J  -- Erases from the cursor to the end of the screen, including the cursor position.
	// [1J -- Erases from the beginning of the screen to the cursor, including the cursor position.
	// [2J -- Erases the complete display. The cursor does not move.
	// Notes:
	// -- Clearing the entire buffer, versus just the Window, works best for Windows Consoles

	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var start COORD
	var end COORD

	switch param ***REMOVED***
	case 0:
		start = info.CursorPosition
		end = COORD***REMOVED***info.Size.X - 1, info.Size.Y - 1***REMOVED***

	case 1:
		start = COORD***REMOVED***0, 0***REMOVED***
		end = info.CursorPosition

	case 2:
		start = COORD***REMOVED***0, 0***REMOVED***
		end = COORD***REMOVED***info.Size.X - 1, info.Size.Y - 1***REMOVED***
	***REMOVED***

	err = h.clearRange(h.attributes, start, end)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If the whole buffer was cleared, move the window to the top while preserving
	// the window-relative cursor position.
	if param == 2 ***REMOVED***
		pos := info.CursorPosition
		window := info.Window
		pos.Y -= window.Top
		window.Bottom -= window.Top
		window.Top = 0
		if err := SetConsoleCursorPosition(h.fd, pos); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := SetConsoleWindowInfo(h.fd, true, window); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) EL(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("EL: [%v]", strconv.Itoa(param))
	h.clearWrap()

	// [K  -- Erases from the cursor to the end of the line, including the cursor position.
	// [1K -- Erases from the beginning of the line to the cursor, including the cursor position.
	// [2K -- Erases the complete line.

	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var start COORD
	var end COORD

	switch param ***REMOVED***
	case 0:
		start = info.CursorPosition
		end = COORD***REMOVED***info.Size.X, info.CursorPosition.Y***REMOVED***

	case 1:
		start = COORD***REMOVED***0, info.CursorPosition.Y***REMOVED***
		end = info.CursorPosition

	case 2:
		start = COORD***REMOVED***0, info.CursorPosition.Y***REMOVED***
		end = COORD***REMOVED***info.Size.X, info.CursorPosition.Y***REMOVED***
	***REMOVED***

	err = h.clearRange(h.attributes, start, end)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) IL(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("IL: [%v]", strconv.Itoa(param))
	h.clearWrap()
	return h.insertLines(param)
***REMOVED***

func (h *windowsAnsiEventHandler) DL(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("DL: [%v]", strconv.Itoa(param))
	h.clearWrap()
	return h.deleteLines(param)
***REMOVED***

func (h *windowsAnsiEventHandler) ICH(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("ICH: [%v]", strconv.Itoa(param))
	h.clearWrap()
	return h.insertCharacters(param)
***REMOVED***

func (h *windowsAnsiEventHandler) DCH(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("DCH: [%v]", strconv.Itoa(param))
	h.clearWrap()
	return h.deleteCharacters(param)
***REMOVED***

func (h *windowsAnsiEventHandler) SGR(params []int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	strings := []string***REMOVED******REMOVED***
	for _, v := range params ***REMOVED***
		strings = append(strings, strconv.Itoa(v))
	***REMOVED***

	h.logf("SGR: [%v]", strings)

	if len(params) <= 0 ***REMOVED***
		h.attributes = h.infoReset.Attributes
		h.inverted = false
	***REMOVED*** else ***REMOVED***
		for _, attr := range params ***REMOVED***

			if attr == ansiterm.ANSI_SGR_RESET ***REMOVED***
				h.attributes = h.infoReset.Attributes
				h.inverted = false
				continue
			***REMOVED***

			h.attributes, h.inverted = collectAnsiIntoWindowsAttributes(h.attributes, h.inverted, h.infoReset.Attributes, int16(attr))
		***REMOVED***
	***REMOVED***

	attributes := h.attributes
	if h.inverted ***REMOVED***
		attributes = invertAttributes(attributes)
	***REMOVED***
	err := SetConsoleTextAttribute(h.fd, attributes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) SU(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("SU: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()
	return h.scrollUp(param)
***REMOVED***

func (h *windowsAnsiEventHandler) SD(param int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("SD: [%v]", []string***REMOVED***strconv.Itoa(param)***REMOVED***)
	h.clearWrap()
	return h.scrollDown(param)
***REMOVED***

func (h *windowsAnsiEventHandler) DA(params []string) error ***REMOVED***
	h.logf("DA: [%v]", params)
	// DA cannot be implemented because it must send data on the VT100 input stream,
	// which is not available to go-ansiterm.
	return nil
***REMOVED***

func (h *windowsAnsiEventHandler) DECSTBM(top int, bottom int) error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("DECSTBM: [%d, %d]", top, bottom)

	// Windows is 0 indexed, Linux is 1 indexed
	h.sr.top = int16(top - 1)
	h.sr.bottom = int16(bottom - 1)

	// This command also moves the cursor to the origin.
	h.clearWrap()
	return h.CUP(1, 1)
***REMOVED***

func (h *windowsAnsiEventHandler) RI() error ***REMOVED***
	if err := h.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.logf("RI: []")
	h.clearWrap()

	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	sr := h.effectiveSr(info.Window)
	if info.CursorPosition.Y == sr.top ***REMOVED***
		return h.scrollDown(1)
	***REMOVED***

	return h.moveCursorVertical(-1)
***REMOVED***

func (h *windowsAnsiEventHandler) IND() error ***REMOVED***
	h.logf("IND: []")
	return h.executeLF()
***REMOVED***

func (h *windowsAnsiEventHandler) Flush() error ***REMOVED***
	h.curInfo = nil
	if h.buffer.Len() > 0 ***REMOVED***
		h.logf("Flush: [%s]", h.buffer.Bytes())
		if _, err := h.buffer.WriteTo(h.file); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if h.wrapNext && !h.drewMarginByte ***REMOVED***
		h.logf("Flush: drawing margin byte '%c'", h.marginByte)

		info, err := GetConsoleScreenBufferInfo(h.fd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		charInfo := []CHAR_INFO***REMOVED******REMOVED***UnicodeChar: uint16(h.marginByte), Attributes: info.Attributes***REMOVED******REMOVED***
		size := COORD***REMOVED***1, 1***REMOVED***
		position := COORD***REMOVED***0, 0***REMOVED***
		region := SMALL_RECT***REMOVED***Left: info.CursorPosition.X, Top: info.CursorPosition.Y, Right: info.CursorPosition.X, Bottom: info.CursorPosition.Y***REMOVED***
		if err := WriteConsoleOutput(h.fd, charInfo, size, position, &region); err != nil ***REMOVED***
			return err
		***REMOVED***
		h.drewMarginByte = true
	***REMOVED***
	return nil
***REMOVED***

// cacheConsoleInfo ensures that the current console screen information has been queried
// since the last call to Flush(). It must be called before accessing h.curInfo or h.curPos.
func (h *windowsAnsiEventHandler) getCurrentInfo() (COORD, *CONSOLE_SCREEN_BUFFER_INFO, error) ***REMOVED***
	if h.curInfo == nil ***REMOVED***
		info, err := GetConsoleScreenBufferInfo(h.fd)
		if err != nil ***REMOVED***
			return COORD***REMOVED******REMOVED***, nil, err
		***REMOVED***
		h.curInfo = info
		h.curPos = info.CursorPosition
	***REMOVED***
	return h.curPos, h.curInfo, nil
***REMOVED***

func (h *windowsAnsiEventHandler) updatePos(pos COORD) ***REMOVED***
	if h.curInfo == nil ***REMOVED***
		panic("failed to call getCurrentInfo before calling updatePos")
	***REMOVED***
	h.curPos = pos
***REMOVED***

// clearWrap clears the state where the cursor is in the margin
// waiting for the next character before wrapping the line. This must
// be done before most operations that act on the cursor.
func (h *windowsAnsiEventHandler) clearWrap() ***REMOVED***
	h.wrapNext = false
	h.drewMarginByte = false
***REMOVED***
