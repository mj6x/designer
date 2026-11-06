package animation

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

// ANSI converts an image into ANSI true-colour escape sequences, rendering each pair of vertical pixels as a single terminal character.
func ANSI(i image.Image) string {
	var out bytes.Buffer
	bounds := i.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := color.NRGBAModel.Convert(i.At(x, y)).RGBA()
			r2, g2, b2, _ := color.NRGBAModel.Convert(i.At(x, y+1)).RGBA()

			fmt.Fprintf(&out, "\x1b[38;2;%d;%d;%d;48;2;%d;%d;%dm▀",
				r1>>8, g1>>8, b1>>8,
				r2>>8, g2>>8, b2>>8,
			)
		}

		out.WriteString("\x1b[0m\r\n")
	}

	return out.String()
}
