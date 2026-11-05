package colour

import (
	"math"
	"strings"
)

// Gradient generates a 2D RGB colour gradient using the specified direction and colours.
func Gradient(direction string, colours []RGB, width, height int) [][]RGB {
	if width <= 0 || height <= 0 || len(colours) == 0 {
		return nil
	}

	// interpolate calculates the RGB colour at a given position between the supplied colours.
	interpolate := func(t float64) RGB {
		if len(colours) == 1 {
			return colours[0]
		}

		t = math.Max(0, math.Min(t, 1))

		position := t * float64(len(colours)-1)
		index := int(position)

		if index >= len(colours)-1 {
			return colours[len(colours)-1]
		}

		fraction := position - float64(index)
		a, b := colours[index], colours[index+1]

		return RGB{
			R: uint8(float64(a.R) + (float64(b.R)-float64(a.R))*fraction),
			G: uint8(float64(a.G) + (float64(b.G)-float64(a.G))*fraction),
			B: uint8(float64(a.B) + (float64(b.B)-float64(a.B))*fraction),
		}
	}

	direction = strings.ToLower(direction)

	output := make([][]RGB, height)

	maxDiagonal := math.Hypot(float64(width-1), float64(height-1))
	maxRadial := math.Hypot(float64(width)/2, float64(height)/2)

	for y := 0; y < height; y++ {
		output[y] = make([]RGB, width)

		for x := 0; x < width; x++ {
			var t float64

			switch direction {
			case "vertical":
				if height > 1 {
					t = float64(y) / float64(height-1)
				}

			case "diagonal":
				if maxDiagonal > 0 {
					t = math.Hypot(float64(x), float64(y)) / maxDiagonal
				}

			case "radial":
				if maxRadial > 0 {
					distance := math.Hypot(
						float64(x)-float64(width-1)/2,
						float64(y)-float64(height-1)/2,
					)
					t = distance / maxRadial
				}

			default:
				if width > 1 {
					t = float64(x) / float64(width-1)
				}
			}

			output[y][x] = interpolate(t)
		}
	}

	return output
}
