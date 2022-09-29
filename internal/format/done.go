package format

import (
	"bytes"

	"github.com/mitchellh/colorstring"
)

func Done(msg string) string {

	c := colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Reset:   true,
		Disable: NoColor,
	}

	var buf bytes.Buffer

	buf.WriteString(c.Color("[bold][green]Done: "))
	buf.WriteString(msg)

	return withLeftRule(&buf, "green")
}
