package format

import (
	"bytes"

	"github.com/mitchellh/colorstring"
	"github.com/padok-team/tfautomv/internal/terraform"
)

func Moves(moves []terraform.Move) string {

	c := colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Reset:   true,
		Disable: NoColor,
	}

	var buf bytes.Buffer

	buf.WriteString(c.Color("[bold][green]Moves"))
	buf.WriteByte('\n')

	for _, move := range moves {
		var moveBuf bytes.Buffer

		moveBuf.WriteString(c.Color("[bold]From: "))
		moveBuf.WriteString(move.From)
		moveBuf.WriteByte('\n')

		moveBuf.WriteString(c.Color("[bold]To:   "))
		moveBuf.WriteString(move.To)
		moveBuf.WriteByte('\n')

		buf.WriteString(withLeftRule(&moveBuf, "white"))
	}

	return withLeftRule(&buf, "green")
}
