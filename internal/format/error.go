package format

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/colorstring"
)

func Error(err error) string {

	c := colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Reset:   true,
		Disable: NoColor,
	}

	var buf bytes.Buffer

	buf.WriteString(c.Color("[bold][red]Error: [reset]\n\n"))
	fmt.Fprintf(&buf, "%s\n", err.Error())

	return withLeftRule(&buf, "red")
}
