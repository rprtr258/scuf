package main

import (
	"os"

	"github.com/rprtr258/scuf"
)

func main() {
	out := scuf.New(os.Stdout)
	out.TAB().
		String("bold", scuf.ModBold).
		SPC().String("faint", scuf.ModFaint).
		SPC().String("italic", scuf.ModItalic).
		SPC().String("underline", scuf.ModUnderline).
		SPC().String("crossout", scuf.ModCrossout).
		SPC().String("overline", scuf.ModOverline).
		SPC().String("reverse", scuf.ModReverse).
		SPC().String("blink", scuf.ModBlink).
		NL().TAB().
		String("red", scuf.FgRGB(scuf.MustParseHexRGB("#E88388"))).
		SPC().String("green", scuf.FgRGB(scuf.MustParseHexRGB("#A8CC8C"))).
		SPC().String("yellow", scuf.FgRGB(scuf.MustParseHexRGB("#DBAB79"))).
		SPC().String("blue", scuf.FgRGB(scuf.MustParseHexRGB("#71BEF2"))).
		SPC().String("magenta", scuf.FgRGB(scuf.MustParseHexRGB("#D290E4"))).
		SPC().String("cyan", scuf.FgRGB(scuf.MustParseHexRGB("#66C2CD"))).
		SPC().String("gray", scuf.FgRGB(scuf.MustParseHexRGB("#B9BFCA"))).
		NL().TAB().
		String("red", scuf.FgBlack, scuf.BgRGB(scuf.MustParseHexRGB("#E88388"))).
		SPC().String("green", scuf.FgBlack, scuf.BgRGB(scuf.MustParseHexRGB("#A8CC8C"))).
		SPC().String("yellow", scuf.FgBlack, scuf.BgRGB(scuf.MustParseHexRGB("#DBAB79"))).
		SPC().String("blue", scuf.FgBlack, scuf.BgRGB(scuf.MustParseHexRGB("#71BEF2"))).
		SPC().String("magenta", scuf.FgBlack, scuf.BgRGB(scuf.MustParseHexRGB("#D290E4"))).
		SPC().String("cyan", scuf.FgBlack, scuf.BgRGB(scuf.MustParseHexRGB("#66C2CD"))).
		SPC().String("gray", scuf.FgBlack, scuf.BgRGB(scuf.MustParseHexRGB("#B9BFCA"))).
		NL().NL()

	hw := "Hello, world!"
	out.
		Copy(hw).
		TAB().Printf("%q", hw).String(" copied to clipboard").
		NL().NL().
		SetWindowTitle(hw).
		TAB().Printf("%q", hw).String(" set as window title").
		NL().NL().
		Notify("Termenv", hw).
		TAB().String("Triggered a notification").
		NL().
		TAB().Hyperlink("http://example.com", "This is a link").
		NL()
}
