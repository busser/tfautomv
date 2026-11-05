package engine

import "testing"

func TestAnalyzeMigrationsMatchesAcrossPlansAndWithinOnePlan(t *testing.T) {
	deleteResource := Resource{ModuleID: "source", Type: "aws_instance", Address: "aws_instance.old", Attributes: map[string]any{"id": "i-123", "ami": "ami-1"}}
	createResource := Resource{ModuleID: "destination", Type: "aws_instance", Address: "aws_instance.new", Attributes: map[string]any{"ami": "ami-1"}}
	withinDelete := Resource{ModuleID: "same", Type: "aws_instance", Address: "aws_instance.before", Attributes: map[string]any{"id": "i-456", "ami": "ami-2"}}
	withinCreate := Resource{ModuleID: "same", Type: "aws_instance", Address: "aws_instance.after", Attributes: map[string]any{"ami": "ami-2"}}

	analysis := AnalyzeMigrations([]Plan{
		{ToDelete: []Resource{deleteResource}, ToCreate: []Resource{withinCreate}},
		{ToDelete: []Resource{withinDelete}, ToCreate: []Resource{createResource}},
	}, nil)

	if len(analysis.Migrations) != 2 {
		t.Fatalf("got %d migrations, want 2", len(analysis.Migrations))
	}
	for _, migration := range analysis.Migrations {
		if migration.SourcePlan == "source" && migration.ImportID != "i-123" {
			t.Fatalf("got import ID %q", migration.ImportID)
		}
	}
}

func TestAnalyzeMigrationsRejectsAmbiguousMatch(t *testing.T) {
	deleted := Resource{ModuleID: "source", Type: "aws_instance", Address: "aws_instance.old", Attributes: map[string]any{"id": "i-123", "ami": "ami-1"}}
	created := func(plan, address string) Resource {
		return Resource{ModuleID: plan, Type: "aws_instance", Address: address, Attributes: map[string]any{"ami": "ami-1"}}
	}
	analysis := AnalyzeMigrations([]Plan{
		{ToDelete: []Resource{deleted}},
		{ToCreate: []Resource{created("one", "aws_instance.one")}},
		{ToCreate: []Resource{created("two", "aws_instance.two")}},
	}, nil)
	if len(analysis.Migrations) != 0 {
		t.Fatal("ambiguous match produced a migration")
	}
	if len(analysis.AmbiguousDeletes) != 1 || len(analysis.AmbiguousCreates) != 2 {
		t.Fatalf("unexpected ambiguity counts: %#v", analysis)
	}
}

func TestMigrationAnalysisValidateNetZero(t *testing.T) {
	analysis := MigrationAnalysis{}
	if err := analysis.ValidateNetZero([]Plan{{ToCreate: []Resource{{}}, ToDelete: nil}}); err == nil {
		t.Fatal("expected net-zero validation error")
	}
}
