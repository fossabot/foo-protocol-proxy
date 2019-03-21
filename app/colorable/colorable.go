package colorable

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type (
	// Colorable interface for console colors.
	Colorable struct {
		output io.Writer
	}
	// ColorValue type.
	ColorValue int
)

const (
	// Escape value at the console output.
	Escape = "\x1b"
)

// Font effects.
const (
	Reset ColorValue = iota
	Bold
	Faint
	Italic
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors.
const (
	FGBlack ColorValue = iota + 30
	FGRed
	FGGreen
	FGYellow
	FGBlue
	Magenta
	FGCyan
	FGWhite
)

// Background text colors.
const (
	BGBlack ColorValue = iota + 40
	BGRed
	BGGReen
	BGYellow
	BGBlue
	BGMagenta
	BGCyan
	BGWhite
)

// NewColorable allocates and returns a new Colorable.
func NewColorable(output io.Writer) *Colorable {
	return &Colorable{
		output: output,
	}
}

// Set sets a group of color values for the next output operations.
func (c *Colorable) Set(color ...ColorValue) *Colorable {
	fmt.Fprintf(c.output, c.format(color...))

	return c
}

// Reset resets the color value to the default.
func (c *Colorable) Reset() {
	fmt.Fprintf(c.output, "%s[%dm", Escape, Reset)
}

// Wrap wraps a passed a string with a color values.
func (c *Colorable) Wrap(str string, color ...ColorValue) string {
	return fmt.Sprintf("%s%s%s", c.format(color...), str, c.resetFormat())
}

func (c *Colorable) format(color ...ColorValue) string {
	return fmt.Sprintf("%s[%sm", Escape, c.sequence(color...))
}

func (c *Colorable) resetFormat() string {
	return fmt.Sprintf("%s[%dm", Escape, Reset)
}

func (c *Colorable) sequence(color ...ColorValue) string {
	format := make([]string, len(color))

	for key, value := range color {
		format[key] = strconv.Itoa(int(value))
	}

	return strings.Join(format, ";")
}
