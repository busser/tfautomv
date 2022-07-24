package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/padok-team/tfautomv/internal/format"
	"github.com/padok-team/tfautomv/internal/terraform"
	"github.com/padok-team/tfautomv/internal/tfautomv"
)

var (
	dryRun       = flag.Bool("dry-run", false, "print moves instead of writing them to disk")
	printVersion = flag.Bool("version", false, "print version and exit")
	showAnalysis = flag.Bool("show-analysis", false, "show detailed analysis of Terraform plan")
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
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		return nil
	}

	tf := terraform.NewRunner(".")

	logln("Running \"terraform init\"...")
	err := tf.Init()
	if err != nil {
		return err
	}

	logln("Running \"terraform plan\"...")
	plan, err := tf.Plan()
	if err != nil {
		return err
	}

	analysis, err := tfautomv.AnalysisFromPlan(plan)
	if err != nil {
		return err
	}
	if *showAnalysis {
		fmt.Fprint(os.Stderr, format.Analysis(analysis))
	}

	moves := tfautomv.MovesFromAnalysis(analysis)

	if *dryRun {
		fmt.Fprint(os.Stderr, format.Moves(moves))
		return nil
	}

	err = terraform.AppendMovesToFile(moves, "moves.tf")
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stderr, format.Done(fmt.Sprintf("Added %d moved blocks to \"moves.tf\".", len(moves))))

	return nil
}

func logln(msg string) {
	fmt.Fprint(os.Stderr, format.Info(msg))
}
