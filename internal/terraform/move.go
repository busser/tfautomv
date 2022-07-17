package terraform

import (
	"fmt"
	"os"
)

type Move struct {
	From string
	To   string
}

func (m Move) Block() string {
	return fmt.Sprintf("moved {\n  from = %s\n  to   = %s\n}", m.From, m.To)
}

func AppendMovesToFile(moves []Move, path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, m := range moves {
		fmt.Fprintln(f, m.Block())
	}

	return nil
}
