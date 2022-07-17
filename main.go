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

	fmt.Fprint(os.Stderr, format.Info("Analysing Terraform plan..."))

	report, err := tfautomv.GenerateReport(".")
	if err != nil {
		return err
	}

	if *showAnalysis {
		fmt.Fprint(os.Stderr, format.Analysis(report.Analysis))
	}

	err = terraform.AppendMovesToFile(report.Moves, "moves.tf")
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stderr, format.Done(fmt.Sprintf("Generated %d moved blocks.", len(report.Moves))))

	return nil
}
