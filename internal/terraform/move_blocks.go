package terraform

import (
	"fmt"
	"os"
	"strings"
)

type MoveBlocks struct {
	Moves []Moved `hcl:"moved,block"`
}

type Moved struct {
	From string `hcl:"from"`
	To   string `hcl:"to"`
}

func (blocks MoveBlocks) String() string {
	var sb strings.Builder
	for _, m := range blocks.Moves {
		fmt.Fprintln(&sb, m)
	}
	return sb.String()
}

func (m Moved) String() string {
	return fmt.Sprintf("moved {\n  from = %s\n  to   = %s\n}", m.From, m.To)
}

func (blocks MoveBlocks) AppendTo(path string) error {
	movesFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer movesFile.Close()

	fmt.Fprint(movesFile, blocks)

	return nil
}
