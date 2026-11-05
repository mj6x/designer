package colour

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// RGB represents the red, green, and blue values of a colour.
type RGB struct {
	R, G, B uint8 // we use uint8 as it is all signed integers from 0 through 255. (rgb colour spectrum)
}

// ANSI returns the true-colour ANSI code.
func (r *RGB) ANSI() string {
	return fmt.Sprintf("\x1b[48;2;0;0;0;38;2;%d;%d:%dm", r.R, r.G, r.B)
}

// hexToRGB converts a hex code into an RGB value.
func hexToRGB(c string) (*RGB, error) {
	c = strings.TrimPrefix(c, "#")

	if len(c) == 3 {
		c = string([]byte{
			c[0], c[0],
			c[1], c[1],
			c[2], c[2],
		})
	}

	if len(c) != 6 {
		return nil, errors.New("invalid hex colour provided")
	}

	value, err := strconv.ParseUint(c, 16, 24)
	if err != nil {
		return nil, errors.New("invalid hex colour provided")
	}

	return &RGB{
		R: uint8(value >> 16),
		G: uint8(value >> 8),
		B: uint8(value),
	}, nil
}

// StringToRGB converts a string into an RGB value.
func StringToRGB(c string) (*RGB, error) {
	if !strings.Contains(c, ",") {
		return hexToRGB(c)
	}

	parts := strings.Split(c, ",")
	if len(parts) != 3 {
		return nil, errors.New("invalid RGB colour provided")
	}

	values := [3]uint8{}

	for i, part := range parts {
		value, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || value < 0 || value > 255 {
			return nil, errors.New("invalid RGB colour provided")
		}

		values[i] = uint8(value)
	}

	return &RGB{
		R: values[0],
		G: values[1],
		B: values[2],
	}, nil
}
