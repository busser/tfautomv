package format

import (
	"bytes"

	"github.com/mitchellh/colorstring"
)

func Done(msg string) string {
	var buf bytes.Buffer

	buf.WriteString(colorstring.Color("[bold][green]Done: "))
	buf.WriteString(msg)

	return WithLeftRule(&buf, "green")
}
