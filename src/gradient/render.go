package gradient

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/mj6x/designer/src/colour"
)

// longestLine works out the longest line of the banner minus the gradient tags.
func longestLine(content []string) int {
	longest := 0

	tagRegex := regexp.MustCompile(`<[^>]*>`)

	for _, line := range content {
		line = tagRegex.ReplaceAllString(line, "")

		length := utf8.RuneCountInString(line)

		if length > longest {
			longest = length
		}
	}

	return longest
}

// Render text using the supplied RGB colours and gradient directions.
func render(data string, output io.Writer, colours []*colour.RGB) {
	rows := strings.Split(data, "\n")
	width := longestLine(rows)
	height := len(rows)

	gradients := map[string][][]*colour.RGB{
		"radial":     colour.Gradient("radial", colours, width, height),
		"diagonal":   colour.Gradient("diagonal", colours, width, height),
		"horizontal": colour.Gradient("horizontal", colours, width, height),
		"vertical":   colour.Gradient("vertical", colours, width, height),
	}

	var (
		current [][]*colour.RGB
		active  bool
	)

	for y, row := range rows {
		var line strings.Builder
		position := 0

		for len(row) > 0 {
			if strings.HasPrefix(row, "<") {
				end := strings.Index(row, ">")
				if end == -1 {
					break
				}

				tag := row[1:end]
				row = row[end+1:]

				if gradient, ok := gradients[tag]; ok {
					current = gradient
					active = true
				} else if tag == "/" || tag == "stop" {
					active = false
					current = nil
				}

				continue
			}

			char, size := utf8.DecodeRuneInString(row)
			row = row[size:]

			if active && char != ' ' && position < width {
				line.WriteString(current[y][position].ANSI())
			}

			line.WriteRune(char)
			position++
		}

		if y > 0 {
			fmt.Fprintln(output)
		}

		fmt.Fprint(output, line.String())
	}
}
