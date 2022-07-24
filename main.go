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
	showAnalysis = flag.Bool("show-analysis", false, "show detailed analysis of Terraform plan")
	printVersion = flag.Bool("version", false, "print version and exit")
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

	logln("Analysing Terraform plan...")
	analysis, err := tfautomv.AnalysisFromPlan(plan)
	if err != nil {
		return err
	}
	if *showAnalysis {
		fmt.Fprint(os.Stderr, format.Analysis(analysis))
	}

	logln("Determining valid moves...")
	moves := tfautomv.MovesFromAnalysis(analysis)

	logln("Adding moves to \"moves.tf\"...")
	err = terraform.AppendMovesToFile(moves, "moves.tf")
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stderr, format.Done(fmt.Sprintf("Generated %d moved blocks.", len(moves))))

	return nil
}

func logln(msg string) {
	fmt.Fprint(os.Stderr, format.Info(msg))
}
