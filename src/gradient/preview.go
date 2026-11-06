package gradient

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mj6x/designer/src/colour"
)

// renderGradient renders a horizontal gradient and outputs it using the lipgloss colour package.
func renderGradient(colours []*colour.RGB) string {
	if len(colours) < 2 {
		return ""
	}

	var result strings.Builder
	const steps = 40

	for i := 0; i < steps; i++ {
		position := float64(i) / float64(steps-1)
		segment := position * float64(len(colours)-1)
		index := int(segment)

		if index >= len(colours)-1 {
			index = len(colours) - 2
		}

		t := segment - float64(index)

		c1 := colours[index]
		c2 := colours[index+1]

		r := int(float64(c1.R)*(1-t) + float64(c2.R)*t)
		g := int(float64(c1.G)*(1-t) + float64(c2.G)*t)
		b := int(float64(c1.B)*(1-t) + float64(c2.B)*t)

		color := lipgloss.Color(
			fmt.Sprintf("#%02X%02X%02X", r, g, b),
		)

		result.WriteString(
			lipgloss.NewStyle().Foreground(color).Render("██"),
		)
	}

	return result.String()
}
