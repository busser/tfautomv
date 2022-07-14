package terraform

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Error struct {
	cmd    *exec.Cmd
	output io.Reader
	err    error
}

func (e Error) Error() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Command %q with args %q failed: %s\n", e.cmd.Path, e.cmd.Args, e.err.Error())
	fmt.Fprintln(&sb, "Command output:")
	io.Copy(&sb, e.output)

	return sb.String()
}

func (e Error) Unwrap() error {
	return e.err
}
