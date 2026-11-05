package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/go-version"
	flag "github.com/spf13/pflag"

	"github.com/busser/tfautomv/pkg/engine"
	"github.com/busser/tfautomv/pkg/engine/rules"
	"github.com/busser/tfautomv/pkg/pretty"
	"github.com/busser/tfautomv/pkg/terraform"
)

func main() {
	if err := run(); err != nil {
		os.Stderr.WriteString(pretty.Error(err) + "\n")
		os.Exit(1)
	}
}

//go:embed VERSION
var tfautomvVersion string

func run() error {
	parseFlags()

	workdirs := flag.Args()
	if len(workdirs) == 0 {
		workdirs = []string{"."}
	}

	if noColor {
		pretty.DisableColors()
	}

	if printVersion {
		fmt.Println(tfautomvVersion)
		return nil
	}

	ctx := context.TODO()

	/*
	 * Step 0: Smoke tests
	 *
	 * Detect any obvious issues with the user's configuration.
	 */

	if outputFormat != "auto" && outputFormat != "blocks" && outputFormat != "commands" && outputFormat != "import-commands" {
		return fmt.Errorf("unknown output format %q", outputFormat)
	}

	if usePreplanned && (skipInit || skipRefresh) {
		return fmt.Errorf("--preplanned cannot be used with --skip-init or --skip-refresh flags")
	}

	if !usePreplanned && flag.Lookup("preplanned-file").Changed {
		return fmt.Errorf("--preplanned-file can only be used with --preplanned")
	}

	tfVersion, err := terraform.GetVersion(ctx, terraform.WithTerraformBin(terraformBin))
	if err != nil {
		return fmt.Errorf("failed to get Terraform version: %w", err)
	}

	movedBlocksSupported := tfVersion.GreaterThanOrEqual(version.Must(version.NewSemver("1.1.0")))
	if outputFormat == "blocks" && !movedBlocksSupported {
		return fmt.Errorf("Terraform version %s does not support moved blocks", tfVersion)
	}

	crossModuleMovesSupported := tfVersion.GreaterThanOrEqual(version.Must(version.NewSemver("0.14.0")))
	if len(workdirs) > 1 && !crossModuleMovesSupported {
		return fmt.Errorf("Terraform version %s does not support moves across modules", tfVersion)
	}

	/*
	 * Step 1: Parse user-provided rules
	 *
	 * Best to parse the rules early on, so that any syntax errors are caught
	 * before we start running Terraform commands.
	 */

	var userRules []engine.Rule
	for _, raw := range ignoreRules {
		rule, err := rules.Parse(raw)
		if err != nil {
			return fmt.Errorf("invalid rule passed with -ignore flag %q: %w", raw, err)
		}

		userRules = append(userRules, rule)
	}

	/*
	 * Step 2: Obtain Terraform plan
	 *
	 * Either run `terraform plan` for each working directory provided by the user,
	 * or read existing plan files if --preplanned flag is used.
	 */

	terraformOptions := []terraform.Option{
		terraform.WithTerraformBin(terraformBin),
		terraform.WithSkipInit(skipInit),
		terraform.WithSkipRefresh(skipRefresh),
	}

	var plans []engine.Plan

	if usePreplanned {
		plans, err = getPreplannedPlans(ctx, workdirs, preplannedFile, terraformOptions)
	} else {
		plans, err = getPlans(ctx, workdirs, terraformOptions)
	}
	if err != nil {
		return err
	}

	/*
	 * Step 3: Use the tfautomv engine to determine moves to make
	 *
	 * We merge the plans from each working directory into a single plan
	 * containing all the resources. This doesn't make a difference to the
	 * engine.
	 */

	mergedPlan := engine.MergePlans(plans)
	comparisons := engine.CompareAll(mergedPlan, userRules)
	moves := engine.DetermineMoves(comparisons)
	var netZeroError error
	if requireNetZero {
		if len(mergedPlan.ToCreate) != len(mergedPlan.ToDelete) {
			netZeroError = fmt.Errorf("net-zero validation failed: %d creates and %d deletes", len(mergedPlan.ToCreate), len(mergedPlan.ToDelete))
		} else if len(moves) != len(mergedPlan.ToCreate) {
			netZeroError = fmt.Errorf("net-zero validation failed: only %d of %d planned creates have unique matches", len(moves), len(mergedPlan.ToCreate))
		}
	}

	/*
	 * Step 4: Print a human-readable summary for the user
	 *
	 * This summary allows the user to gain insight into what tfautomv has
	 * done, and why. This information allows the user to make an informed
	 * decision about what to do next.
	 */

	summarizer := pretty.NewSummarizer(moves, comparisons, verbosity)
	summary := summarizer.Summary()

	os.Stderr.WriteString("\n" + summary + "\n\n")

	/*
	 * Step 5: Write the moves found by the engine.
	 *
	 * Depending on the output format chosen by the user, we either write
	 * the moves to a file, print them to standard output, or a combination
	 * of both.
	 */

	terraformMoves := engineMovesToTerraformMoves(moves)
	sameModule, differentModule := categorizeMoves(terraformMoves)
	if dryRun {
		printMappingReport(terraformMoves)
		if netZeroError != nil {
			printMigrationDiagnostics(mergedPlan, comparisons, moves)
			return netZeroError
		}
		return nil
	}
	if netZeroError != nil {
		printMappingReport(terraformMoves)
		printMigrationDiagnostics(mergedPlan, comparisons, moves)
		return netZeroError
	}

	switch outputFormat {
	case "auto":
		switch {
		case movedBlocksSupported:
			if err := writeMovedBlocks(sameModule); err != nil {
				return err
			}
			if err := writeMoveCommands(differentModule, terraformOptions...); err != nil {
				return err
			}
		case !movedBlocksSupported:
			if err := writeMoveCommands(terraformMoves, terraformOptions...); err != nil {
				return err
			}
		}
	case "blocks":
		if err := writeMigrationBlocks(terraformMoves); err != nil {
			return err
		}
	case "commands":
		if err := writeMoveCommands(terraformMoves, terraformOptions...); err != nil {
			return err
		}
	case "import-commands":
		if err := writeImportCommands(terraformMoves); err != nil {
			return err
		}
	default:
		// This should have been caught by the smoke tests.
		return fmt.Errorf("unknown output format %q", outputFormat)
	}

	return nil
}

// Flags
var (
	ignoreRules    []string
	noColor        bool
	outputFormat   string
	printVersion   bool
	skipInit       bool
	skipRefresh    bool
	terraformBin   string
	verbosity      int
	preplannedFile string
	usePreplanned  bool
	dryRun         bool
	requireNetZero bool
)

func parseFlags() {
	flag.StringSliceVar(&ignoreRules, "ignore", nil, "ignore differences based on a `rule`")
	flag.BoolVar(&noColor, "no-color", false, "disable color in output")
	flag.StringVarP(&outputFormat, "output", "o", "auto", "output `format` of moves (\"auto\", \"blocks\", \"commands\" or \"import-commands\")")
	flag.BoolVarP(&printVersion, "version", "V", false, "print version and exit")
	flag.BoolVarP(&skipInit, "skip-init", "s", false, "skip running terraform init")
	flag.BoolVarP(&skipRefresh, "skip-refresh", "S", false, "skip running terraform refresh")
	flag.StringVar(&terraformBin, "terraform-bin", "terraform", "terraform binary to use")
	flag.CountVarP(&verbosity, "verbosity", "v", "increase verbosity (can be specified multiple times)")
	flag.BoolVar(&usePreplanned, "preplanned", false, "use existing plan files instead of running terraform plan")
	flag.StringVar(&preplannedFile, "preplanned-file", "tfplan.bin", "plan file name when using --preplanned")
	flag.BoolVar(&dryRun, "dry-run", false, "show resource mappings without writing files or state commands")
	flag.BoolVar(&requireNetZero, "require-net-zero", false, "fail unless planned creates equal deletes and every operation matches")

	flag.Parse()
}

func engineMovesToTerraformMoves(moves []engine.Move) []terraform.Move {
	var terraformMoves []terraform.Move

	for _, m := range moves {
		terraformMoves = append(terraformMoves, terraform.Move{
			FromWorkdir: m.SourceModule,
			ToWorkdir:   m.DestinationModule,
			FromAddress: m.SourceAddress,
			ToAddress:   m.DestinationAddress,
			ImportID:    m.ImportID,
		})
	}

	return terraformMoves
}

// printMappingReport gives dry-run users a complete mapping before any state
// or configuration output is changed.
func printMappingReport(moves []terraform.Move) {
	if len(moves) == 0 {
		fmt.Println("No resource mappings found.")
		return
	}
	for _, move := range moves {
		fmt.Printf("%s:%s -> %s:%s", move.FromWorkdir, move.FromAddress, move.ToWorkdir, move.ToAddress)
		if move.ImportID != "" {
			fmt.Printf(" (id=%s)", move.ImportID)
		}
		fmt.Println()
	}
}

// printMigrationDiagnostics explains why strict matching rejected a plan so
// users can distinguish a real configuration change from an ambiguous move.
func printMigrationDiagnostics(plan engine.Plan, comparisons []engine.ResourceComparison, moves []engine.Move) {
	matched := make(map[string]bool)
	for _, move := range moves {
		matched[move.SourceModule+":"+move.SourceAddress] = true
		matched[move.DestinationModule+":"+move.DestinationAddress] = true
	}
	fmt.Printf("Migration diagnostics: %d creates, %d deletes, %d unique mappings\n", len(plan.ToCreate), len(plan.ToDelete), len(moves))
	for _, resource := range plan.ToCreate {
		key := resource.ModuleID + ":" + resource.Address
		if !matched[key] {
			fmt.Printf("  unmatched create: %s:%s (%s)\n", resource.ModuleID, resource.Address, resource.Type)
		}
	}
	for _, resource := range plan.ToDelete {
		key := resource.ModuleID + ":" + resource.Address
		if !matched[key] {
			fmt.Printf("  unmatched delete: %s:%s (%s)\n", resource.ModuleID, resource.Address, resource.Type)
		}
	}
	for _, comparison := range comparisons {
		if comparison.IsMatch() {
			continue
		}
		fmt.Printf("  non-matching candidate: %s:%s -> %s:%s; differing attributes: %s\n", comparison.ToDelete.ModuleID, comparison.ToDelete.Address, comparison.ToCreate.ModuleID, comparison.ToCreate.Address, strings.Join(comparison.MismatchingAttributes, ", "))
	}
}

func getPlans(ctx context.Context, workdirs []string, options []terraform.Option) ([]engine.Plan, error) {
	type result struct {
		plan engine.Plan
		err  error
	}
	results := make([]result, len(workdirs))

	getPlan := func(i int) {
		workdir := workdirs[i]

		os.Stderr.WriteString(pretty.Colorf("getting Terraform plan for %s...", (*pretty.Summarizer).StyledModule(nil, workdir)) + "\n")

		workdirOptions := append(
			[]terraform.Option{terraform.WithWorkdir(workdir)},
			options...,
		)

		jsonPlan, err := terraform.GetPlan(ctx, workdirOptions...)
		if err != nil {
			results[i].err = fmt.Errorf("failed to get plan for workdir %q: %w", workdir, err)
			return
		}

		plan, err := engine.SummarizeJSONPlan(workdirs[i], jsonPlan)
		if err != nil {
			results[i].err = fmt.Errorf("failed to summarize plan for workdir %q: %w", workdir, err)
			return
		}

		results[i].plan = plan
	}

	var wg sync.WaitGroup
	for i := range workdirs {
		wg.Add(1)
		go func(i int) {
			getPlan(i)
			wg.Done()
		}(i)
	}

	wg.Wait()

	var errs []error
	for _, r := range results {
		if r.err != nil {
			errs = append(errs, r.err)
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	var plans []engine.Plan
	for _, r := range results {
		plans = append(plans, r.plan)
	}

	return plans, nil
}

func getPreplannedPlans(ctx context.Context, workdirs []string, planFilename string, options []terraform.Option) ([]engine.Plan, error) {
	type result struct {
		plan engine.Plan
		err  error
	}
	results := make([]result, len(workdirs))

	// First, validate that all directories have the plan file
	for _, workdir := range workdirs {
		planPath := filepath.Join(workdir, planFilename)
		if _, err := os.Stat(planPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("plan file not found: %s (all directories must have plan files when using --preplanned)", planPath)
		}
	}

	getPlan := func(i int) {
		workdir := workdirs[i]
		planPath := filepath.Join(workdir, planFilename)

		os.Stderr.WriteString(pretty.Colorf("reading Terraform plan from %s...", (*pretty.Summarizer).StyledModule(nil, planPath)) + "\n")

		workdirOptions := append(
			[]terraform.Option{terraform.WithWorkdir(workdir)},
			options...,
		)

		jsonPlan, err := terraform.GetPlanFromFile(ctx, planPath, workdirOptions...)
		if err != nil {
			results[i].err = fmt.Errorf("failed to read plan from %q: %w", planPath, err)
			return
		}

		plan, err := engine.SummarizeJSONPlan(workdir, jsonPlan)
		if err != nil {
			results[i].err = fmt.Errorf("failed to summarize plan from %q: %w", planPath, err)
			return
		}

		results[i].plan = plan
	}

	var wg sync.WaitGroup
	for i := range workdirs {
		wg.Add(1)
		go func(i int) {
			getPlan(i)
			wg.Done()
		}(i)
	}

	wg.Wait()

	var errs []error
	for _, r := range results {
		if r.err != nil {
			errs = append(errs, r.err)
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	var plans []engine.Plan
	for _, r := range results {
		plans = append(plans, r.plan)
	}

	return plans, nil
}

func categorizeMoves(moves []terraform.Move) (sameWorkdir, differentWorkdir []terraform.Move) {
	for _, m := range moves {
		if m.FromWorkdir == m.ToWorkdir {
			sameWorkdir = append(sameWorkdir, m)
		} else {
			differentWorkdir = append(differentWorkdir, m)
		}
	}

	return sameWorkdir, differentWorkdir
}

func writeMovedBlocks(moves []terraform.Move) error {
	if len(moves) == 0 {
		return nil
	}

	movesByWorkdir := make(map[string][]terraform.Move)
	for _, m := range moves {
		movesByWorkdir[m.FromWorkdir] = append(movesByWorkdir[m.FromWorkdir], m)
	}

	for workdir, moves := range movesByWorkdir {
		movesFilePath := filepath.Join(workdir, "moves.tf")
		movesFile, err := os.OpenFile(movesFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open %q: %w", movesFilePath, err)
		}

		err = terraform.WriteMovedBlocks(movesFile, moves)
		if err != nil {
			return fmt.Errorf("failed to write moved blocks: %w", err)
		}

		os.Stderr.WriteString(pretty.Colorf("%s written to [bold][green]%s", pretty.StyledNumMoves(len(moves)), movesFilePath))
		os.Stderr.WriteString("\n")
	}

	return nil
}

func writeMigrationBlocks(moves []terraform.Move) error {
	if len(moves) == 0 {
		return nil
	}

	for _, move := range moves {
		if move.FromWorkdir == move.ToWorkdir {
			if err := appendMigrationBlock(filepath.Join(move.FromWorkdir, "moves.tf"), move, false); err != nil {
				return err
			}
			continue
		}
		if err := appendMigrationBlock(filepath.Join(move.FromWorkdir, "state-migration.tf"), move, true); err != nil {
			return err
		}
		if err := appendMigrationBlock(filepath.Join(move.ToWorkdir, "state-migration.tf"), move, false); err != nil {
			return err
		}
	}
	return nil
}

func appendMigrationBlock(path string, move terraform.Move, source bool) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %q: %w", path, err)
	}
	if move.FromWorkdir == move.ToWorkdir {
		err = terraform.WriteMovedBlocks(file, []terraform.Move{move})
	} else if source {
		err = terraform.WriteRemovedBlock(file, move)
	} else {
		err = terraform.WriteImportBlock(file, move)
	}
	closeErr := file.Close()
	if err != nil {
		return fmt.Errorf("failed to write migration block: %w", err)
	}
	if closeErr != nil {
		return fmt.Errorf("failed to close %q: %w", path, closeErr)
	}
	os.Stderr.WriteString(pretty.Colorf("migration block written to [bold][green]%s", path) + "\n")
	return nil
}

func writeImportCommands(moves []terraform.Move) error {
	if len(moves) == 0 {
		return nil
	}
	if err := terraform.WriteImportCommands(os.Stdout, moves, terraformBin); err != nil {
		return fmt.Errorf("failed to write import commands: %w", err)
	}
	return nil
}

func writeMoveCommands(moves []terraform.Move, options ...terraform.Option) error {
	if len(moves) == 0 {
		return nil
	}

	err := terraform.WriteMoveCommands(os.Stdout, moves, options...)
	if err != nil {
		return fmt.Errorf("failed to write move commands: %w", err)
	}

	os.Stderr.WriteString(pretty.Colorf("%s written to [bold][green]standard output", pretty.StyledNumMoves(len(moves))))
	os.Stderr.WriteString("\n")

	return nil
}
