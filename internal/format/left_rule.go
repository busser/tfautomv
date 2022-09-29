package format

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/colorstring"
)

func withLeftRule(r io.Reader, lineColor string) string {

	c := colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Reset:   true,
		Disable: NoColor,
	}

	// these leftRule* variables are markers for the beginning of the lines
	// containing the error that are intended to help sighted users
	// better understand the information hierarchy when errors appear
	// alongside other information or alongside other diagnostics.
	//
	// Without this, it seems (based on folks sharing incomplete messages when
	// asking questions, or including extra content that's not part of the
	// error) that some readers have trouble easily identifying which
	// text belongs to the error and which does not.
	var (
		leftRuleStart = c.Color(fmt.Sprintf("[%s]╷[reset]", lineColor))
		leftRuleLine  = c.Color(fmt.Sprintf("[%s]│[reset]", lineColor))
		leftRuleEnd   = c.Color(fmt.Sprintf("[%s]╵[reset]", lineColor))
	)

	// We add the left rule prefixes to each line so that the overall message is
	// visually delimited from what's around it. We'll do that by scanning over
	// what we already generated and adding the prefix for each line.
	var sb strings.Builder
	sc := bufio.NewScanner(r)
	sb.WriteString(leftRuleStart)
	sb.WriteByte('\n')
	for sc.Scan() {
		line := sc.Text()
		sb.WriteString(leftRuleLine)
		if line != "" {
			sb.WriteByte(' ')
			sb.WriteString(line)
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(leftRuleEnd)
	sb.WriteByte('\n')

	return sb.String()
}
