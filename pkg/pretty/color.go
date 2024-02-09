package pretty

import (
	"fmt"
	"os"

	"github.com/mitchellh/colorstring"
)

var colorsDisabled = false

func init() {
	// https://no-color.org/
	if os.Getenv("NO_COLOR") != "" {
		DisableColors()
	}
}

// DisableColors disables colors and other formatting.
func DisableColors() {
	colorsDisabled = true
}

// EnableColors enables colors and other formatting.
func EnableColors() {
	colorsDisabled = false
}

// Color returns a string with color and other formatting applied.
// Under the hood, this function uses the colorstring package:
// github.com/mitchellh/colorstring
func Color(s string) string {
	c := colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Reset:   true,
		Disable: colorsDisabled,
	}

	return c.Color(s)
}

// Colorf returns a string with color and other formatting applied.
// When inserting user-provided values into a format string, use this
// function instead of Color. This is to avoid mistakenly interpreting user
// data as formatting directives.
func Colorf(format string, a ...interface{}) string {
	return fmt.Sprintf(Color(format), a...)
}
