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
	// The provider ID recorded in the source state for cross-state imports.
	ImportID string
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
	if len(moves) == 0 {
		return nil
	}

	var blocks []string

	for _, move := range moves {
		if !move.isWithinSameWorkdir() {
			return fmt.Errorf("cannot write blocks for moves between different working directories")
		}

		blocks = append(blocks, move.block())
	}

	_, err := w.Write([]byte(strings.Join(blocks, "\n") + "\n"))

	return err
}

// WriteMigrationBlocks writes moved blocks for same-workdir moves and paired
// removed/import blocks for moves between workdirs. Terraform recommends the
// latter for new migrations between separate state files or remote backends.
func WriteMigrationBlocks(w io.Writer, moves []Move) error {
	var blocks []string
	for _, move := range moves {
		if move.isWithinSameWorkdir() {
			blocks = append(blocks, move.block())
			continue
		}
		if move.ImportID == "" {
			return fmt.Errorf("cannot write import block for %s: source state has no resource ID", move.FromAddress)
		}
		blocks = append(blocks, fmt.Sprintf("removed {\n  from = %s\n\n  lifecycle {\n    destroy = false\n  }\n}", move.FromAddress))
		blocks = append(blocks, fmt.Sprintf("import {\n  to = %s\n  id = %q\n}", move.ToAddress, move.ImportID))
	}
	if len(blocks) == 0 {
		return nil
	}
	_, err := w.Write([]byte(strings.Join(blocks, "\n\n") + "\n"))
	return err
}

// WriteRemovedBlock writes the source-side block for a cross-state move.
func WriteRemovedBlock(w io.Writer, move Move) error {
	if move.isWithinSameWorkdir() {
		return fmt.Errorf("removed blocks are only for cross-workdir moves")
	}
	_, err := fmt.Fprintf(w, "removed {\n  from = %s\n\n  lifecycle {\n    destroy = false\n  }\n}\n", move.FromAddress)
	return err
}

// WriteImportBlock writes the destination-side block for a cross-state move.
func WriteImportBlock(w io.Writer, move Move) error {
	if move.ImportID == "" {
		return fmt.Errorf("cannot write import block for %s: source state has no resource ID", move.FromAddress)
	}
	_, err := fmt.Fprintf(w, "import {\n  to = %s\n  id = %q\n}\n", move.ToAddress, move.ImportID)
	return err
}

// WriteImportCommands emits a conservative cross-backend transfer. It imports
// the destination first and removes the source binding second; same-workdir
// moves remain state mv commands.
func WriteImportCommands(w io.Writer, moves []Move, terraformBin string) error {
	var commands []string
	for _, move := range moves {
		if !move.isWithinSameWorkdir() && move.ImportID == "" {
			return fmt.Errorf("cannot write import command for %s: source state has no resource ID", move.FromAddress)
		}
		commands = append(commands,
			fmt.Sprintf("%s -chdir=%q import %q %q", terraformBin, move.ToWorkdir, move.ToAddress, move.ImportID),
		)
		if move.isWithinSameWorkdir() {
			commands[len(commands)-1] = fmt.Sprintf("%s -chdir=%q state mv %q %q", terraformBin, move.FromWorkdir, move.FromAddress, move.ToAddress)
		} else {
			// Import first so a failed destination operation does not lose the
			// source binding.
			commands = append(commands, fmt.Sprintf("%s -chdir=%q state rm %q", terraformBin, move.FromWorkdir, move.FromAddress))
		}
	}
	_, err := fmt.Fprintln(w, strings.Join(commands, "\n"))
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
func WriteMoveCommands(w io.Writer, moves []Move, opts ...Option) error {
	var settings settings

	settings.apply(append(defaultOptions(), opts...))

	err := settings.validate()
	if err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	return writeMoveCommands(w, moves, settings.terraformBin)
}

func writeMoveCommands(w io.Writer, moves []Move, terraformBin string) error {
	var commands []string

	// Start with moves within the same module.

	for _, m := range moves {
		if m.FromWorkdir == m.ToWorkdir {
			var chdirFlag string
			if m.FromWorkdir != "." {
				chdirFlag = fmt.Sprintf(" -chdir=%q", m.FromWorkdir)
			}

			commands = append(commands,
				fmt.Sprintf("%s%s state mv %q %q",
					terraformBin,
					chdirFlag,
					m.FromAddress,
					m.ToAddress,
				),
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
			fmt.Sprintf("%s -chdir=%q state pull > %q",
				terraformBin,
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
			fmt.Sprintf("%s state mv -state=%q -state-out=%q %q %q",
				terraformBin,
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
			fmt.Sprintf("%s -chdir=%q state push %q",
				terraformBin,
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
