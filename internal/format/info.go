package format

import (
	"fmt"

	"github.com/mitchellh/colorstring"
)

func Info(msg string) string {

	c := colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Reset:   true,
		Disable: NoColor,
	}

	return c.Color(fmt.Sprintf("[bold]%s\n", msg))
}
