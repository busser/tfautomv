package terraform

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

// A Move represents a an object in Terraform's state that should be moved to
// another address and possibly to another working directory.
// This is analogous to a moved block or a state mv command.
type Move struct {
	// The working directory the resource is being moved from. Equal to
	// ToWorkdir when the resource is being moved within the same working
	// directory.
	FromWorkdir string
	// The working directory the resource is being moved to. Equal to
	// FromWorkdir when the resource is being moved within the same working
	// directory.
	ToWorkdir string

	// The resource's address before the move.
	FromAddress string
	// The resource's address after the move.
	ToAddress string
}

func (m Move) block() string {
	return fmt.Sprintf("moved {\n  from = %s\n  to   = %s\n}", m.FromAddress, m.ToAddress)
}

func (m Move) isWithinSameWorkdir() bool {
	return m.FromWorkdir == m.ToWorkdir
}

// WriteMovedBlocks encodes the given moves as a series of Terraform moved
// blocks, in HCL, and writes them to the given writer.
//
// Currently, moved blocks cannot be used to move resources between different
// working directories. If the given moves contain such a move, WriteMovedBlocks
// returns an error.
func WriteMovedBlocks(w io.Writer, moves []Move) error {
	var blocks []string

	for _, move := range moves {
		if !move.isWithinSameWorkdir() {
			return fmt.Errorf("cannot write blocks for moves between different working directories")
		}

		blocks = append(blocks, move.block())
	}

	_, err := w.Write([]byte(strings.Join(blocks, "\n")))

	return err
}

// WriteMoveCommands encodes the given moves as a series of Terraform CLI
// commands and writes them to the given writer.
//
// Moving resources across working directories requires a bit of extra work.
// In those cases, WriteMoveCommands will write a series of commands that
// copies the state of the working directories involved to the local filesystem,
// performs the moves, and pushes the states back to configured backends.
//
// Since Terraform's state can contain sensitive information, the copies are
// written in the same directory as the working directories. This allows the
// user to know exactly where their state is stored. The files are not deleted
// automatically, in case the user wants to use them (to revert any changes, for
// instance).
func WriteMoveCommands(w io.Writer, moves []Move) error {
	var commands []string

	// Start with moves within the same module.

	for _, m := range moves {
		if m.FromWorkdir == m.ToWorkdir {
			var chdirFlag string
			if m.FromWorkdir != "." {
				chdirFlag = fmt.Sprintf("-chdir=%q", m.FromWorkdir)
			}

			commands = append(commands,
				fmt.Sprintf("terraform %s state mv %q %q",
					chdirFlag, m.FromAddress, m.ToAddress),
			)
		}
	}

	// Then, pull the states of all working directories that require
	// cross-directory moves.

	var workdirs []string
	for _, m := range moves {
		if m.FromWorkdir != m.ToWorkdir {
			workdirs = append(workdirs, m.FromWorkdir, m.ToWorkdir)
		}
	}
	workdirs = unique(workdirs)
	sort.Strings(workdirs)

	const localCopyFileName = ".tfautomv.tfstate"

	for _, workdir := range workdirs {
		commands = append(commands,
			fmt.Sprintf("terraform -chdir=%q state pull > %q",
				workdir,
				filepath.Join(workdir, localCopyFileName),
			),
		)
	}

	// Next, perform all the moves.

	for _, move := range moves {
		if move.FromWorkdir == move.ToWorkdir {
			// Already handled above.
			continue
		}

		commands = append(commands,
			fmt.Sprintf("terraform state mv -state=%q -state-out=%q %q %q",
				filepath.Join(move.FromWorkdir, localCopyFileName),
				filepath.Join(move.ToWorkdir, localCopyFileName),
				move.FromAddress,
				move.ToAddress,
			),
		)
	}

	// Then, push the states of all modules we manipulated.

	for _, workdir := range workdirs {
		commands = append(commands,
			fmt.Sprintf("terraform -chdir=%q state push %q",
				workdir,
				localCopyFileName,
			),
		)
	}

	// And we're done.

	_, err := fmt.Fprintln(w, strings.Join(commands, "\n"))

	return err
}

func unique(s []string) []string {
	seen := make(map[string]struct{})
	var unique []string
	for _, e := range s {
		if _, ok := seen[e]; !ok {
			unique = append(unique, e)
			seen[e] = struct{}{}
		}
	}
	return unique
}
