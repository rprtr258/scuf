package scuf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	for name, test := range map[string]struct {
		modify   func(Buffer) Buffer
		expected string
	}{
		"ClearScreen":              {(Buffer).ClearScreen, "\x1b[2J\x1b[1;1H"},
		"SaveCursorPosition":       {(Buffer).SaveCursorPosition, "\x1b[s"},
		"RestoreCursorPosition":    {(Buffer).RestoreCursorPosition, "\x1b[u"},
		"ClearLine":                {(Buffer).ClearLine, "\x1b[2K"},
		"ClearLineLeft":            {(Buffer).ClearLineLeft, "\x1b[1K"},
		"ClearLineRight":           {(Buffer).ClearLineRight, "\x1b[0K"},
		"SaveScreen":               {func(b Buffer) Buffer { return b.Enable(SaveScreen) }, "\x1b[?47h"},
		"RestoreScreen":            {func(b Buffer) Buffer { return b.Disable(SaveScreen) }, "\x1b[?47l"},
		"AltScreen":                {func(b Buffer) Buffer { return b.Enable(AltScreen) }, "\x1b[?1049h"},
		"ExitAltScreen":            {func(b Buffer) Buffer { return b.Disable(AltScreen) }, "\x1b[?1049l"},
		"ShowCursor":               {func(b Buffer) Buffer { return b.Enable(Cursor) }, "\x1b[?25h"},
		"HideCursor":               {func(b Buffer) Buffer { return b.Disable(Cursor) }, "\x1b[?25l"},
		"EnableMousePress":         {func(b Buffer) Buffer { return b.Enable(MousePress) }, "\x1b[?9h"},
		"DisableMousePress":        {func(b Buffer) Buffer { return b.Disable(MousePress) }, "\x1b[?9l"},
		"EnableMouse":              {func(b Buffer) Buffer { return b.Enable(Mouse) }, "\x1b[?1000h"},
		"DisableMouse":             {func(b Buffer) Buffer { return b.Disable(Mouse) }, "\x1b[?1000l"},
		"EnableMouseHilite":        {func(b Buffer) Buffer { return b.Enable(MouseHilite) }, "\x1b[?1001h"},
		"DisableMouseHilite":       {func(b Buffer) Buffer { return b.Disable(MouseHilite) }, "\x1b[?1001l"},
		"EnableMouseCellMotion":    {func(b Buffer) Buffer { return b.Enable(MouseCellMotion) }, "\x1b[?1002h"},
		"DisableMouseCellMotion":   {func(b Buffer) Buffer { return b.Disable(MouseCellMotion) }, "\x1b[?1002l"},
		"EnableMouseAllMotion":     {func(b Buffer) Buffer { return b.Enable(MouseAllMotion) }, "\x1b[?1003h"},
		"DisableMouseAllMotion":    {func(b Buffer) Buffer { return b.Disable(MouseAllMotion) }, "\x1b[?1003l"},
		"EnableMouseExtendedMode":  {func(b Buffer) Buffer { return b.Enable(MouseExtendedMode) }, "\x1b[?1006h"},
		"DisableMouseExtendedMode": {func(b Buffer) Buffer { return b.Disable(MouseExtendedMode) }, "\x1b[?1006l"},
		"EnableMousePixelsMode":    {func(b Buffer) Buffer { return b.Enable(MousePixelsMode) }, "\x1b[?1016h"},
		"DisableMousePixelsMode":   {func(b Buffer) Buffer { return b.Disable(MousePixelsMode) }, "\x1b[?1016l"},
		"SetForegroundColor":       {func(b Buffer) Buffer { return b.SetForegroundColor("#000000") }, "\x1b]10;#000000\a"},
		"SetBackgroundColor":       {func(b Buffer) Buffer { return b.SetBackgroundColor("#000000") }, "\x1b]11;#000000\a"},
		"SetCursorColor":           {func(b Buffer) Buffer { return b.SetCursorColor("#000000") }, "\x1b]12;#000000\a"},
		"MoveCursor":               {func(b Buffer) Buffer { return b.MoveCursor(16, 8) }, "\x1b[16;8H"},
		"CursorUp":                 {func(b Buffer) Buffer { return b.CursorUp(8) }, "\x1b[8A"},
		"CursorDown":               {func(b Buffer) Buffer { return b.CursorDown(8) }, "\x1b[8B"},
		"CursorForward":            {func(b Buffer) Buffer { return b.CursorForward(8) }, "\x1b[8C"},
		"CursorBack":               {func(b Buffer) Buffer { return b.CursorBack(8) }, "\x1b[8D"},
		"CursorNextLine":           {func(b Buffer) Buffer { return b.CursorNextLine(8) }, "\x1b[8E"},
		"CursorPrevLine":           {func(b Buffer) Buffer { return b.CursorPrevLine(8) }, "\x1b[8F"},
		"ClearLines":               {func(b Buffer) Buffer { return b.ClearLines(8) }, "\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K"},
		"ChangeScrollingRegion":    {func(b Buffer) Buffer { return b.ChangeScrollingRegion(16, 8) }, "\x1b[16;8r"},
		"InsertLines":              {func(b Buffer) Buffer { return b.InsertLines(8) }, "\x1b[8L"},
		"DeleteLines":              {func(b Buffer) Buffer { return b.DeleteLines(8) }, "\x1b[8M"},
		"SetWindowTitle":           {func(b Buffer) Buffer { return b.SetWindowTitle("test") }, "\x1b]2;test\a"},
		"CopyClipboard":            {func(b Buffer) Buffer { return b.Copy("hello") }, "\x1b]52;c;aGVsbG8=\a"},
		"CopyPrimary":              {func(b Buffer) Buffer { return b.CopyPrimary("hello") }, "\x1b]52;p;aGVsbG8=\a"},
		"Hyperlink":                {func(b Buffer) Buffer { return b.Hyperlink("http://example.com", "example") }, "\x1b]8;;http://example.com\x1b\\example\x1b]8;;\x1b\\"},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, NewString(func(b Buffer) {
				test.modify(b)
			}))
		})
	}
}

func TestRendering(t *testing.T) {
	// "Unstyled strings should be returned as plain text"
	assert.Equal(t, "foobar", String("foobar"))

	assert.Equal(t, "\x1b[38;2;171;205;239;48;5;69;1;3;2;4;5mfoobar\x1b[0m", String("foobar",
		FgRGB(MustParseHexRGB("#abcdef")),
		BgANSI256(69),
		ModBold,
		ModItalic,
		ModFaint,
		ModUnderline,
		ModBlink))
}

func TestColorConversion(t *testing.T) {
	for name, test := range map[string]struct {
		hex string
		c   []byte
	}{
		"ANSI color":     {"#c0c0c0", BgWhite},
		"BgHiGreen":      {"#00ff00", BgHiGreen},
		"ANSI-256 color": {"#8700af", BgANSI256(91)},
		"hex color":      {"#abcdef", BgRGB(MustParseHexRGB("#abcdef"))},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.hex, ToHex(test.c))
		})
	}
}

func TestFromColor(t *testing.T) {
	assert.Equal(t, Modifier("38;2;255;128;0"), FgRGB(255, 128, 0))
}

func TestStyles(t *testing.T) {
	assert.Equal(t, "\x1b[32mfoobar\x1b[0m", String("foobar", FgGreen))
}

func TestString(t *testing.T) {
	for name, test := range map[string]struct {
		mods   []Modifier
		result string
	}{
		"blue on red": {
			[]Modifier{FgBlue, BgRed},
			"\x1b[34;41mtext\x1b[0m",
		},
		"magenta on white": {
			[]Modifier{FgMagenta, BgWhite},
			"\x1b[35;47mtext\x1b[0m",
		},
		"cyan": {
			[]Modifier{FgCyan},
			"\x1b[36mtext\x1b[0m",
		},
		"default on red": {
			[]Modifier{BgRed},
			"\x1b[41mtext\x1b[0m",
		},
		"default bold on yellow": {
			[]Modifier{ModBold, BgYellow},
			"\x1b[1;43mtext\x1b[0m",
		},
		"bold": {
			[]Modifier{ModBold},
			"\x1b[1mtext\x1b[0m",
		},
		"no color at all": {
			[]Modifier{},
			"text",
		},
		"nil modifiers": {
			[]Modifier{nil, nil, nil},
			"text",
		},
	} {
		t.Run(name, func(t *testing.T) {
			if output := String("text", test.mods...); output != test.result {
				t.Errorf("Expected %q, got %q", test.result, output)
			}
		})
	}
}
