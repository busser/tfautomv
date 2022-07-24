package format

import (
	"bytes"

	"github.com/mitchellh/colorstring"
	"github.com/padok-team/tfautomv/internal/terraform"
)

func Moves(moves []terraform.Move) string {
	var buf bytes.Buffer

	buf.WriteString(colorstring.Color("[bold][green]Moves"))
	buf.WriteByte('\n')

	for _, move := range moves {
		var moveBuf bytes.Buffer

		moveBuf.WriteString(colorstring.Color("[bold]From: "))
		moveBuf.WriteString(move.From)
		moveBuf.WriteByte('\n')

		moveBuf.WriteString(colorstring.Color("[bold]To:   "))
		moveBuf.WriteString(move.To)
		moveBuf.WriteByte('\n')

		buf.WriteString(WithLeftRule(&moveBuf, "white"))
	}

	return WithLeftRule(&buf, "green")
}
