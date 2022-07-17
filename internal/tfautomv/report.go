package tfautomv

import "github.com/padok-team/tfautomv/internal/terraform"

// Report summarizes the work done by tfautomv.
type Report struct {
	InitialPlan *terraform.Plan
	Analysis    *Analysis
	Moves       []terraform.Move
}
