package format

import (
	"fmt"

	"github.com/mitchellh/colorstring"
)

func Info(msg string) string {
	return colorstring.Color(fmt.Sprintf("[bold]%s\n", msg))
}
