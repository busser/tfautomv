package terraform

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/busser/tfautomv/pkg/engine"
)

// WriteMigrationBlocks writes the declarative blocks needed for state
// migration. Same-plan moves use moved blocks; cross-plan moves use a
// destroy-free removed block and an import block in the destination plan.
func WriteMigrationBlocks(w io.Writer, migrations []engine.Migration) error {
	return writeMigrationBlocks(w, migrations, "")
}

// WriteMigrationBlocksForPlan writes only blocks belonging to one named plan,
// allowing callers to generate separate configuration files for each state.
func WriteMigrationBlocksForPlan(w io.Writer, migrations []engine.Migration, planID string) error {
	return writeMigrationBlocks(w, migrations, planID)
}

func writeMigrationBlocks(w io.Writer, migrations []engine.Migration, planID string) error {
	var blocks []string
	for _, migration := range migrations {
		if migration.SourcePlan == migration.DestinationPlan {
			if planID != "" && migration.SourcePlan != planID {
				continue
			}
			blocks = append(blocks, fmt.Sprintf("moved {\n  from = %s\n  to   = %s\n}", migration.SourceAddress, migration.DestinationAddr))
			continue
		}
		if migration.ImportID == "" {
			return fmt.Errorf("resource %s has no state ID for cross-plan import", migration.SourceAddress)
		}
		if planID == "" || migration.SourcePlan == planID {
			blocks = append(blocks, fmt.Sprintf("removed {\n  from = %s\n\n  lifecycle {\n    destroy = false\n  }\n}", migration.SourceAddress))
		}
		if planID == "" || migration.DestinationPlan == planID {
			blocks = append(blocks, fmt.Sprintf("import {\n  to = %s\n  id = %s\n}", migration.DestinationAddr, strconv.Quote(migration.ImportID)))
		}
	}
	if len(blocks) == 0 {
		return nil
	}
	_, err := io.WriteString(w, strings.Join(blocks, "\n\n")+"\n")
	return err
}
