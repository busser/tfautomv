package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"

	"github.com/padok-team/tfautomv/internal/format"
	"github.com/padok-team/tfautomv/internal/terraform"
	"github.com/padok-team/tfautomv/internal/tfautomv"
	"github.com/padok-team/tfautomv/internal/tfautomv/ignore"
)

func main() {
	if err := run(); err != nil {
		os.Stderr.WriteString(format.Error(err))
		os.Exit(1)
	}
}

//go:embed VERSION
var version string

func run() error {
	parseFlags()

	if noColor {
		format.NoColor = true
	}

	if printVersion {
		fmt.Println(version)
		return nil
	}

	tf := terraform.NewRunner(".")

	// Some Terraform versions do not support some of tfautomv's output options.
	// Check that everything is OK early on, to avoid wasting time running a
	// plan for nothing.

	tfVer, err := tf.Version()
	if err != nil {
		return err
	}

	switch outputFormat {
	case "blocks":
		if tfVer.LessThan(semver.MustParse("1.1")) {
			return fmt.Errorf("terraform version %s does not support moved blocks", tfVer.String())
		}
	case "commands":
	default:
		return fmt.Errorf("unknown output format %q", outputFormat)
	}

	// Parse rules early on so that the user gets quick feedback in case of
	// syntax errors.
	var rules []ignore.Rule
	for _, raw := range ignoreRules {
		r, err := ignore.ParseRule(raw)
		if err != nil {
			return fmt.Errorf("invalid rule passed with -ignore flag %q: %w", raw, err)
		}
		rules = append(rules, r)
	}

	// Terraform's plan contains a lot of information. For now, this is all we
	// need. In the future, we may choose to use other sources of information.

	logln("Running \"terraform init\"...")
	err = tf.Init()
	if err != nil {
		return err
	}

	logln("Running \"terraform plan\"...")
	plan, err := tf.Plan()
	if err != nil {
		return err
	}

	analysis, err := tfautomv.AnalysisFromPlan(plan, rules)
	if err != nil {
		return err
	}
	if showAnalysis {
		fmt.Fprint(os.Stderr, format.Analysis(analysis))
	}

	moves := tfautomv.MovesFromAnalysis(analysis)
	if len(moves) == 0 {
		fmt.Fprint(os.Stderr, format.Done("Found no moves to make"))
		return nil
	}

	// At this point, we need to output the moves we found. The Terraform
	// community originally used `tf state mv` commands. Terraform 1.1+ supports
	// moved blocks as a replacement, but those remain incomplete for now.
	// Community tools like tfmigrate are also popular.

	if dryRun {
		fmt.Fprint(os.Stderr, format.Moves(moves))
		return nil
	}

	switch outputFormat {
	case "blocks":
		err = terraform.AppendMovesToFile(moves, "moves.tf")
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stderr, format.Done(fmt.Sprintf("Added %d moved blocks to \"moves.tf\".", len(moves))))

	case "commands":
		terraform.WriteMovesShellCommands(moves, os.Stdout)
		fmt.Fprint(os.Stderr, format.Done(fmt.Sprintf("Wrote %d commands to standard output.", len(moves))))

	default:
		return fmt.Errorf("unknown output format %q", outputFormat)
	}

	return nil
}

func logln(msg string) {
	fmt.Fprint(os.Stderr, format.Info(msg))
}

// Flags
var (
	dryRun       bool
	ignoreRules  []string
	noColor      bool
	outputFormat string
	printVersion bool
	showAnalysis bool
)

func parseFlags() {
	flag.BoolVar(&dryRun, "dry-run", false, "print moves instead of writing them to disk")
	flag.Var(stringSliceValue{&ignoreRules}, "ignore", "ignore differences based on a `rule`")
	flag.BoolVar(&noColor, "no-color", false, "disable color in output")
	flag.StringVar(&outputFormat, "output", "blocks", "output `format` of moves (\"blocks\" or \"commands\")")
	flag.BoolVar(&showAnalysis, "show-analysis", false, "show detailed analysis of Terraform plan")
	flag.BoolVar(&printVersion, "version", false, "print version and exit")

	flag.Parse()
}

type stringSliceValue struct {
	s *[]string
}

func (v stringSliceValue) String() string {
	if v.s == nil || *v.s == nil {
		return ""
	}
	return fmt.Sprintf("%q", *v.s)
}

func (v stringSliceValue) Set(raw string) error {
	*v.s = append(*v.s, raw)
	return nil
}
