package scuf

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/aymanbagabas/go-osc52/v2"
)

type Modifier []byte

// hex values of ANSI colors
var ansiHex = [...]string{
	"#000000", "#800000", "#008000", "#808000", "#000080", "#800080", "#008080", "#c0c0c0",
	"#808080", "#ff0000", "#00ff00", "#ffff00", "#0000ff", "#ff00ff", "#00ffff", "#ffffff",
	"#000000", "#00005f", "#000087", "#0000af", "#0000d7", "#0000ff", "#005f00", "#005f5f",
	"#005f87", "#005faf", "#005fd7", "#005fff", "#008700", "#00875f", "#008787", "#0087af",
	"#0087d7", "#0087ff", "#00af00", "#00af5f", "#00af87", "#00afaf", "#00afd7", "#00afff",
	"#00d700", "#00d75f", "#00d787", "#00d7af", "#00d7d7", "#00d7ff", "#00ff00", "#00ff5f",
	"#00ff87", "#00ffaf", "#00ffd7", "#00ffff", "#5f0000", "#5f005f", "#5f0087", "#5f00af",
	"#5f00d7", "#5f00ff", "#5f5f00", "#5f5f5f", "#5f5f87", "#5f5faf", "#5f5fd7", "#5f5fff",
	"#5f8700", "#5f875f", "#5f8787", "#5f87af", "#5f87d7", "#5f87ff", "#5faf00", "#5faf5f",
	"#5faf87", "#5fafaf", "#5fafd7", "#5fafff", "#5fd700", "#5fd75f", "#5fd787", "#5fd7af",
	"#5fd7d7", "#5fd7ff", "#5fff00", "#5fff5f", "#5fff87", "#5fffaf", "#5fffd7", "#5fffff",
	"#870000", "#87005f", "#870087", "#8700af", "#8700d7", "#8700ff", "#875f00", "#875f5f",
	"#875f87", "#875faf", "#875fd7", "#875fff", "#878700", "#87875f", "#878787", "#8787af",
	"#8787d7", "#8787ff", "#87af00", "#87af5f", "#87af87", "#87afaf", "#87afd7", "#87afff",
	"#87d700", "#87d75f", "#87d787", "#87d7af", "#87d7d7", "#87d7ff", "#87ff00", "#87ff5f",
	"#87ff87", "#87ffaf", "#87ffd7", "#87ffff", "#af0000", "#af005f", "#af0087", "#af00af",
	"#af00d7", "#af00ff", "#af5f00", "#af5f5f", "#af5f87", "#af5faf", "#af5fd7", "#af5fff",
	"#af8700", "#af875f", "#af8787", "#af87af", "#af87d7", "#af87ff", "#afaf00", "#afaf5f",
	"#afaf87", "#afafaf", "#afafd7", "#afafff", "#afd700", "#afd75f", "#afd787", "#afd7af",
	"#afd7d7", "#afd7ff", "#afff00", "#afff5f", "#afff87", "#afffaf", "#afffd7", "#afffff",
	"#d70000", "#d7005f", "#d70087", "#d700af", "#d700d7", "#d700ff", "#d75f00", "#d75f5f",
	"#d75f87", "#d75faf", "#d75fd7", "#d75fff", "#d78700", "#d7875f", "#d78787", "#d787af",
	"#d787d7", "#d787ff", "#d7af00", "#d7af5f", "#d7af87", "#d7afaf", "#d7afd7", "#d7afff",
	"#d7d700", "#d7d75f", "#d7d787", "#d7d7af", "#d7d7d7", "#d7d7ff", "#d7ff00", "#d7ff5f",
	"#d7ff87", "#d7ffaf", "#d7ffd7", "#d7ffff", "#ff0000", "#ff005f", "#ff0087", "#ff00af",
	"#ff00d7", "#ff00ff", "#ff5f00", "#ff5f5f", "#ff5f87", "#ff5faf", "#ff5fd7", "#ff5fff",
	"#ff8700", "#ff875f", "#ff8787", "#ff87af", "#ff87d7", "#ff87ff", "#ffaf00", "#ffaf5f",
	"#ffaf87", "#ffafaf", "#ffafd7", "#ffafff", "#ffd700", "#ffd75f", "#ffd787", "#ffd7af",
	"#ffd7d7", "#ffd7ff", "#ffff00", "#ffff5f", "#ffff87", "#ffffaf", "#ffffd7", "#ffffff",
	"#080808", "#121212", "#1c1c1c", "#262626", "#303030", "#3a3a3a", "#444444", "#4e4e4e",
	"#585858", "#626262", "#6c6c6c", "#767676", "#808080", "#8a8a8a", "#949494", "#9e9e9e",
	"#a8a8a8", "#b2b2b2", "#bcbcbc", "#c6c6c6", "#d0d0d0", "#dadada", "#e4e4e4", "#eeeeee",
}

func ternary[T any](b bool, x, y T) T {
	if b {
		return x
	}
	return y
}

func FgANSI(col int) Modifier {
	// 0-7  -> 30-37
	// 8-15 -> 90-97
	return []byte(strconv.Itoa(col + ternary(col < 8, 30, 82)))
}

func BgANSI(col int) Modifier {
	// 0-7  -> 40-47
	// 8-15 -> 100-107
	return []byte(strconv.Itoa(col + ternary(col < 8, 40, 92)))
}

// c is 16-255
func FgANSI256(c int) Modifier {
	return []byte(fmt.Sprintf("38;5;%d", c))
}

// c is 16-255
func BgANSI256(c int) Modifier {
	return []byte(fmt.Sprintf("48;5;%d", c))
}

// MustParseHexRGB parses hex color string, in form "#f0c" or "#ff1034".
// If color is invalid, returns black (or junk)
func MustParseHexRGB(hex string) (r, g, b uint8) {
	isShort := len(hex) == 4

	format := ternary(isShort, "#%1x%1x%1x", "#%02x%02x%02x")
	n, err := fmt.Sscanf(hex, format, &r, &g, &b)
	if err != nil || n != 3 {
		// err = fmt.Errorf("color: %v is not a hex-color", hex)
		return
	}

	factor := ternary(isShort, 15, uint8(1))
	return r * factor, g * factor, b * factor
}

func FgRGB(r, g, b uint8) Modifier {
	return []byte(fmt.Sprintf("38;2;%d;%d;%d", r, g, b))
}

func BgRGB(r, g, b uint8) Modifier {
	return []byte(fmt.Sprintf("48;2;%d;%d;%d", r, g, b))
}

func ToHex(color Modifier) string {
	switch {
	case !bytes.ContainsRune(color, ';'):
		i, _ := strconv.Atoi(string(color))

		// just some magic formula
		j := i / 10
		return ansiHex[i-10*j+j>>3<<3]
	case color[3] == '5':
		var skip, i int
		fmt.Sscanf(string(color), "%d;5;%d", &skip, &i)
		return ansiHex[i]
	case color[3] == '2':
		var skip, r, g, b int
		fmt.Sscanf(string(color), "%d;2;%d;%d;%d", &skip, &r, &g, &b)
		return fmt.Sprintf("#%2x%2x%2x", r, g, b)
	default:
		return "invalid color"
	}
}

var (
	// Foreground colors
	FgBlack   = FgANSI(0)
	FgRed     = FgANSI(1)
	FgGreen   = FgANSI(2)
	FgYellow  = FgANSI(3)
	FgBlue    = FgANSI(4)
	FgMagenta = FgANSI(5)
	FgCyan    = FgANSI(6)
	FgWhite   = FgANSI(7)
	// Foreground bright colors
	FgHiBlack   = FgANSI(8)
	FgHiRed     = FgANSI(9)
	FgHiGreen   = FgANSI(10)
	FgHiYellow  = FgANSI(11)
	FgHiBlue    = FgANSI(12)
	FgHiMagenta = FgANSI(13)
	FgHiCyan    = FgANSI(14)
	FgHiWhite   = FgANSI(15)

	// Background colors
	BgBlack   = BgANSI(0)
	BgRed     = BgANSI(1)
	BgGreen   = BgANSI(2)
	BgYellow  = BgANSI(3)
	BgBlue    = BgANSI(4)
	BgMagenta = BgANSI(5)
	BgCyan    = BgANSI(6)
	BgWhite   = BgANSI(7)
	// Background bright colors
	BgHiBlack   = BgANSI(8)
	BgHiRed     = BgANSI(9)
	BgHiGreen   = BgANSI(10)
	BgHiYellow  = BgANSI(11)
	BgHiBlue    = BgANSI(12)
	BgHiMagenta = BgANSI(13)
	BgHiCyan    = BgANSI(14)
	BgHiWhite   = BgANSI(15)

	// Common consts
	_esc              byte = '\x1b'               // Escape character
	_csi                   = Modifier{_esc, '['}  // Control Sequence Introducer
	_osc                   = Modifier{_esc, ']'}  // Operating System Command
	_stringTerminator      = Modifier{_esc, '\\'} // String Terminator
	_modReset              = Modifier("0")
	ModBold                = Modifier("1")
	ModFaint               = Modifier("2")
	ModItalic              = Modifier("3")
	ModUnderline           = Modifier("4")
	ModBlink               = Modifier("5")
	ModReverse             = Modifier("7")
	ModCrossout            = Modifier("9")
	ModOverline            = Modifier("53")
)

type Buffer struct {
	w io.Writer
}

func New(out io.Writer) Buffer {
	return Buffer{out}
}

func (b Buffer) write(bs ...byte) Buffer {
	b.w.Write(bs) //nolint:errcheck // fuck you
	return b
}

// Bytes writes bytes to buffer
func (b Buffer) Bytes(bs ...byte) Buffer {
	return b.write(bs...)
}

// RepeatByte repeats byte n times
func (b Buffer) RepeatByte(c byte, n int) Buffer {
	return b.write(bytes.Repeat([]byte{c}, n)...)
}

// Printf writes formatted data to buffer
func (b Buffer) Printf(format string, args ...any) Buffer {
	fmt.Fprintf(b.w, format, args...)
	return b
}

func (b Buffer) writeMods(mods ...Modifier) {
	b.write(_csi...)
	for i, mod := range mods {
		if i > 0 {
			b.write(';')
		}
		b.write(mod...)
	}
	b.write('m')
}

// Styled write things in callback using modifiers. Don't use Styled inside Styled.
func (b Buffer) Styled(f func(Buffer), mods ...Modifier) Buffer {
	if len(mods) == 0 {
		f(b)
		return b
	}

	b.writeMods(mods...)
	f(b)
	b.writeMods(_modReset)
	return b
}

// String writes string to buffer with given modifiers
func (b Buffer) String(s string, mods ...Modifier) Buffer {
	return b.Styled(func(b Buffer) {
		io.WriteString(b.w, s) //nolint:errcheck // fuck you
	}, mods...)
}

// NL writes newline
func (b Buffer) NL() Buffer {
	return b.write('\n')
}

// TAB writes tab
func (b Buffer) TAB() Buffer {
	return b.write('\t')
}

// SPC writes space
func (b Buffer) SPC() Buffer {
	return b.write(' ')
}

// InBytePair writes callback inside given byte pair, e.g. parentheses or quotes
func (b Buffer) InBytePair(start, end byte, f func(Buffer)) Buffer {
	return b.Styled(func(b Buffer) {
		b.Bytes(start)
		f(b)
		b.Bytes(end)
	})
}

// Iter iterates callback with buffer modifier collbacks
func (b Buffer) Iter(seq func(yield func(func(Buffer)) bool) bool) Buffer {
	seq(func(f func(Buffer)) bool {
		f(b)
		return true
	})
	return b
}

// Hyperlink writes hyperlink using OSC8
func (b Buffer) Hyperlink(link, name string) Buffer {
	return b.
		write(_osc...).
		String("8;;").
		String(link).
		write(_stringTerminator...).
		String(name).
		write(_osc...).
		String("8;;").
		write(_stringTerminator...)
}

// Notify triggers a notification using OSC777
func (b Buffer) Notify(title, body string) Buffer {
	return b.
		write(_osc...).
		String("777;notify;").
		String(title).
		Bytes(';').
		String(body).
		Bytes(_stringTerminator...)
}

// Copy text to clipboard using OSC 52 escape sequence
func (b Buffer) Copy(str string) Buffer {
	s := osc52.New(str)
	s.WriteTo(b.w) //nolint:errcheck // fuck you
	return b
}

// CopyPrimary text to primary clipboard (X11) using OSC 52 escape sequence
func (b Buffer) CopyPrimary(str string) Buffer {
	s := osc52.New(str).Primary()
	s.WriteTo(b.w) //nolint:errcheck // fuck you
	return b
}

// SetForegroundColor set default foreground color
func (b Buffer) SetForegroundColor(hex string) Buffer {
	return b.write(_osc...).Printf("10;%s\a", hex)
}

// SetBackgroundColor set default background color
func (b Buffer) SetBackgroundColor(hex string) Buffer {
	return b.write(_osc...).Printf("11;%s\a", hex)
}

// SetCursorColor set cursor color
func (b Buffer) SetCursorColor(hex string) Buffer {
	return b.write(_osc...).Printf("12;%s\a", hex)
}

// MoveCursor moves the cursor to a given position.
func (b Buffer) MoveCursor(row, column int) Buffer {
	return b.write(_csi...).Printf("%d;%dH", row, column)
}

// ClearScreen clears the visible portion of the terminal.
func (b Buffer) ClearScreen() Buffer {
	return b.
		write(_csi...).
		Printf("%dJ", 2).
		MoveCursor(1, 1)
}

// SaveCursorPosition saves the cursor position.
func (b Buffer) SaveCursorPosition() Buffer {
	return b.write(_csi...).String("s")
}

// RestoreCursorPosition restores a saved cursor position.
func (b Buffer) RestoreCursorPosition() Buffer {
	return b.write(_csi...).String("u")
}

// CursorUp moves the cursor up a given number of lines.
func (b Buffer) CursorUp(n int) Buffer {
	return b.write(_csi...).Printf("%dA", n)
}

// CursorDown moves the cursor down a given number of lines.
func (b Buffer) CursorDown(n int) Buffer {
	return b.write(_csi...).Printf("%dB", n)
}

// CursorForward moves the cursor up a given number of lines.
func (b Buffer) CursorForward(n int) Buffer {
	return b.write(_csi...).Printf("%dC", n)
}

// CursorBack moves the cursor backwards a given number of cells.
func (b Buffer) CursorBack(n int) Buffer {
	return b.write(_csi...).Printf("%dD", n)
}

// CursorNextLine moves the cursor down a given number of lines and places it at
// the beginning of the line.
func (b Buffer) CursorNextLine(n int) Buffer {
	return b.write(_csi...).Printf("%dE", n)
}

// CursorPrevLine moves the cursor up a given number of lines and places it at
// the beginning of the line.
func (b Buffer) CursorPrevLine(n int) Buffer {
	return b.write(_csi...).Printf("%dF", n)
}

// ClearLineRight clears the line to the right of the cursor.
func (b Buffer) ClearLineRight() Buffer {
	return b.write(_csi...).String("0K")
}

// ClearLineLeft clears the line to the left of the cursor.
func (b Buffer) ClearLineLeft() Buffer {
	return b.write(_csi...).String("1K")
}

// ClearLine clears the current line.
func (b Buffer) ClearLine() Buffer {
	return b.write(_csi...).String("2K")
}

// ClearLines clears a given number of lines.
func (b Buffer) ClearLines(n int) Buffer {
	b.write(_csi...).Printf("2K")
	for i := 0; i < n; i++ {
		b.CursorUp(1).write(_csi...).Printf("2K")
	}
	return b
}

// ChangeScrollingRegion sets the scrolling region of the terminal
func (b Buffer) ChangeScrollingRegion(top, bottom int) Buffer {
	return b.write(_csi...).Printf("%d;%dr", top, bottom)
}

// InsertLines inserts the given number of lines at the top of the scrollable
// region, pushing lines below down.
func (b Buffer) InsertLines(n int) Buffer {
	return b.write(_csi...).Printf("%dL", n)
}

// DeleteLines deletes the given number of lines, pulling any lines in the scrollable region below up
func (b Buffer) DeleteLines(n int) Buffer {
	return b.write(_csi...).Printf("%dM", n)
}

// SetWindowTitle sets the terminal window title
func (b Buffer) SetWindowTitle(title string) Buffer {
	return b.write(_osc...).Printf("2;%s\a", title)
}

// Device things that can be enabled and disabled e.g. for handling events,
// changing to alt screen, showing/hiding cursor, etc.
type Device int

const (
	// MousePress enables X10 mouse mode. Button press events are sent only. press only (X10)
	MousePress Device = 9
	// Cursor shows cursor on enable, hides on disable
	Cursor Device = 25
	// SaveScreen saves the screen state on enable, on disable restors previously saved screen state
	SaveScreen Device = 47
	// Mouse enables Mouse Tracking mode. press, release, wheel
	Mouse Device = 1000
	// MouseHilite enables Hilite Mouse Tracking mode.
	MouseHilite Device = 1001
	// MouseCellMotion enables Cell Motion Mouse Tracking mode. press, release, move on pressed, wheel
	MouseCellMotion Device = 1002
	// MouseAllMotion enables All Motion Mouse mode. press, release, move, wheel
	MouseAllMotion Device = 1003
	// MouseExtendedMotion enables Extended Mouse mode (SGR). This should be enabled in conjunction with
	// MouseCellMotion, and MouseAllMotion. press, release, move, wheel, extended coordinates
	MouseExtendedMode Device = 1006
	// MousePixelsMotion enables Pixel Motion Mouse mode (SGR-Pixels). This should be enabled in conjunction with
	// MouseCellMotion, and MouseAllMotion. press, release, move, wheel, extended pixel coordinates
	MousePixelsMode Device = 1016
	// AltScreen switches to the alternate screen buffer on enable. On disable former view is restored
	AltScreen Device = 1049
	// BracketedPaste enables bracketed paste
	BracketedPaste Device = 2004
)

// Enable device. See devices for details
func (b Buffer) Enable(d Device) Buffer {
	return b.write(_csi...).Printf("?%dh", d)
}

// Disable device. See devices for details
func (b Buffer) Disable(d Device) Buffer {
	return b.write(_csi...).Printf("?%dl", d)
}

// NewString creates a string from a function modifying buffer
func NewString(f func(Buffer)) string {
	var bb bytes.Buffer
	f(New(&bb))
	return bb.String()
}

// String creates a string with given modifiers
func String(s string, mods ...Modifier) string {
	return NewString(func(b Buffer) {
		b.String(s, mods...)
	})
}
