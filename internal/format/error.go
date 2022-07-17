package format

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/colorstring"
)

func Error(err error) string {
	var buf bytes.Buffer

	buf.WriteString(colorstring.Color("[bold][red]Error: [reset]\n\n"))
	fmt.Fprintf(&buf, "%s\n", err.Error())

	return WithLeftRule(&buf, "red")
}
