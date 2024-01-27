package main

import (
	"fmt"
	"os"

	"github.com/rprtr258/scuf"
)

func fgStep(i, step int) scuf.Modifier {
	if i >= step {
		return scuf.FgBlack
	}
	return scuf.FgWhite
}

func main() {
	b := scuf.New(os.Stdout)

	b.String("Basic ANSI colors", scuf.ModBold).NL()
	for i := 0; i < 16; i++ {
		if i%8 == 0 {
			b.NL()
		}

		bg := scuf.BgANSI(i)
		b.Styled(func(b scuf.Buffer) {
			b.Printf(" %2d %s ", i, scuf.ToHex(bg))
		}, fgStep(i, 5), bg)
	}
	b.NL().NL()

	fmt.Println(scuf.String("Extended ANSI colors", scuf.ModBold))
	for i := 16; i < 232; i++ {
		if (i-16)%6 == 0 {
			b.NL()
		}

		bg := scuf.BgANSI(i)
		b.Styled(func(b scuf.Buffer) {
			b.Printf(" %3d %s ", i, scuf.ToHex(bg))
		}, fgStep(i, 28), bg)
	}
	b.NL().NL()

	b.String("Extended ANSI Grayscale", scuf.ModBold).NL()
	for i := 232; i < 256; i++ {
		if (i-232)%6 == 0 {
			b.NL()
		}

		bg := scuf.BgANSI(i)
		b.Styled(func(b scuf.Buffer) {
			b.Printf(" %3d %s ", i, scuf.ToHex(bg))
		}, fgStep(i, 244), bg)
	}
	b.NL().NL()
}
