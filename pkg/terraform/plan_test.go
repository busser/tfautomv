package terraform

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadJSONPlanFile(t *testing.T) {
	// Create a minimal valid JSON plan
	jsonPlan := `{
		"format_version": "1.0",
		"terraform_version": "1.0.0",
		"planned_values": {},
		"resource_changes": [],
		"configuration": {}
	}`

	// Create a temporary file
	tmpDir := t.TempDir()
	planFile := filepath.Join(tmpDir, "test.json")

	err := os.WriteFile(planFile, []byte(jsonPlan), 0644)
	if err != nil {
		t.Fatalf("Failed to write test plan file: %v", err)
	}

	// Test reading the JSON plan file
	plan, err := readJSONPlanFile(planFile)
	if err != nil {
		t.Fatalf("readJSONPlanFile failed: %v", err)
	}

	if plan.FormatVersion != "1.0" {
		t.Errorf("Expected format_version 1.0, got %s", plan.FormatVersion)
	}

	if plan.TerraformVersion != "1.0.0" {
		t.Errorf("Expected terraform_version 1.0.0, got %s", plan.TerraformVersion)
	}
}

func TestGetPlanFromFile_JSONFile(t *testing.T) {
	// Create a minimal valid JSON plan
	jsonPlan := `{
		"format_version": "1.0",
		"terraform_version": "1.0.0",
		"planned_values": {},
		"resource_changes": [],
		"configuration": {}
	}`

	// Create a temporary file
	tmpDir := t.TempDir()
	planFile := filepath.Join(tmpDir, "test.json")

	err := os.WriteFile(planFile, []byte(jsonPlan), 0644)
	if err != nil {
		t.Fatalf("Failed to write test plan file: %v", err)
	}

	ctx := context.Background()

	// Test reading the JSON plan file through GetPlanFromFile
	plan, err := GetPlanFromFile(ctx, planFile)
	if err != nil {
		t.Fatalf("GetPlanFromFile failed: %v", err)
	}

	if plan.FormatVersion != "1.0" {
		t.Errorf("Expected format_version 1.0, got %s", plan.FormatVersion)
	}
}

func TestGetPlanFromFile_NonexistentFile(t *testing.T) {
	ctx := context.Background()

	_, err := GetPlanFromFile(ctx, "/nonexistent/plan.json")
	if err == nil {
		t.Fatal("Expected error for nonexistent file, got nil")
	}

	if !strings.Contains(err.Error(), "plan file does not exist") {
		t.Errorf("Expected 'plan file does not exist' error, got: %v", err)
	}
}
