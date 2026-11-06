package animation

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

// ASCII converts an image into coloured ASCII art using ANSI true-colour escape sequences.
func ASCII(i image.Image) string {
	var out bytes.Buffer

	const chars = " .,:;i1tfLCG08@&"

	bounds := i.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := color.NRGBAModel.Convert(i.At(x, y)).RGBA()

			brightness := (r>>8 + g>>8 + b>>8) * (a >> 8) / (255 * 3)
			index := int(brightness * uint32(len(chars)-1) / 255)

			fmt.Fprintf(&out, "\x1b[38;2;%d;%d;%dm%c",
				r>>8, g>>8, b>>8, chars[index],
			)
		}

		out.WriteString("\x1b[0m\r\n")
	}

	return out.String()
}
