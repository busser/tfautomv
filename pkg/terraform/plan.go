package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

// GetPlan obtains a Terraform plan from the module in the given working
// directory. It does so by running a series of Terraform commands.
func GetPlan(ctx context.Context, opts ...Option) (*tfjson.Plan, error) {
	var settings settings

	settings.apply(append(defaultOptions(), opts...))

	err := settings.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	tf, err := tfexec.NewTerraform(settings.workdir, settings.terraformBin)
	if err != nil {
		return nil, fmt.Errorf("failed to create Terraform executor: %w", err)
	}

	if !settings.skipInit {
		err := tf.Init(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Terraform: %w", err)
		}
	}

	planFile, err := os.CreateTemp("", "tfautomv.*.plan")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file to store raw plan: %w", err)
	}
	defer os.Remove(planFile.Name())

	_, err = tf.Plan(ctx, tfexec.Out(planFile.Name()), tfexec.Refresh(!settings.skipRefresh))
	if err != nil {
		return nil, fmt.Errorf("failed to compute Terraform plan: %w", err)
	}

	plan, err := tf.ShowPlanFile(ctx, planFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read raw Terraform plan: %w", err)
	}

	return plan, nil
}

// GetPlanFromFile reads a Terraform plan from a file. The file can be either
// a binary plan file or a JSON plan file. If the file has a .json extension,
// it's treated as JSON. Otherwise, it's treated as binary and converted using
// terraform show.
func GetPlanFromFile(ctx context.Context, planPath string, opts ...Option) (*tfjson.Plan, error) {
	var settings settings
	settings.apply(append(defaultOptions(), opts...))

	err := settings.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	// Check if plan file exists
	if _, err := os.Stat(planPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plan file does not exist: %s", planPath)
	}

	// If file has .json extension, read it directly as JSON
	if strings.HasSuffix(strings.ToLower(planPath), ".json") {
		return readJSONPlanFile(planPath)
	}

	// Otherwise, treat as binary plan and convert using terraform show
	return convertBinaryPlanToJSON(ctx, planPath, settings)
}

// readJSONPlanFile reads a JSON plan file directly
func readJSONPlanFile(planPath string) (*tfjson.Plan, error) {
	data, err := os.ReadFile(planPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	var plan tfjson.Plan
	err = json.Unmarshal(data, &plan)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON plan file: %w", err)
	}

	return &plan, nil
}

// convertBinaryPlanToJSON converts a binary plan file to JSON using terraform show
func convertBinaryPlanToJSON(ctx context.Context, planPath string, settings settings) (*tfjson.Plan, error) {
	// Get the directory of the plan file to use as working directory
	planDir := filepath.Dir(planPath)
	if planDir == "." {
		planDir = settings.workdir
	}

	tf, err := tfexec.NewTerraform(planDir, settings.terraformBin)
	if err != nil {
		return nil, fmt.Errorf("failed to create Terraform executor: %w", err)
	}

	plan, err := tf.ShowPlanFile(ctx, planPath)
	if err != nil {
		return nil, fmt.Errorf("failed to convert binary plan to JSON: %w", err)
	}

	return plan, nil
}
