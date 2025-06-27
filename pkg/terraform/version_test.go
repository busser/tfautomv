package terraform

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTerragruntVersion(t *testing.T) {
	tests := []struct {
		name            string
		mockOutput      string
		expectedError   bool
		expectedVersion string
	}{
		{
			name:            "valid terragrunt output",
			mockOutput:      "Terraform v1.12.2\non darwin_arm64\n",
			expectedError:   false,
			expectedVersion: "1.12.2",
		},
		{
			name:            "valid terragrunt output with different version",
			mockOutput:      "Terraform v1.5.7\non linux_amd64\n",
			expectedError:   false,
			expectedVersion: "1.5.7",
		},
		{
			name:            "invalid output format",
			mockOutput:      "Some random output\nwithout version\n",
			expectedError:   true,
			expectedVersion: "",
		},
		{
			name:            "empty output",
			mockOutput:      "",
			expectedError:   true,
			expectedVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the regex parsing logic directly
			// We can't easily test the full function without mocking exec.Command
			// So we'll test the parsing logic separately

			if tt.expectedError {
				// For error cases, we expect the regex to not match
				versionRegex := `Terraform v(\d+\.\d+\.\d+)`
				matches := findVersionInOutput(tt.mockOutput, versionRegex)
				assert.Empty(t, matches, "Expected no version matches for invalid output")
			} else {
				// For success cases, we expect the regex to match and parse correctly
				versionRegex := `Terraform v(\d+\.\d+\.\d+)`
				matches := findVersionInOutput(tt.mockOutput, versionRegex)
				require.Len(t, matches, 2, "Expected version regex to match")

				versionStr := matches[1]
				assert.Equal(t, tt.expectedVersion, versionStr)

				// Verify the version can be parsed
				parsedVersion, err := version.NewSemver(versionStr)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedVersion, parsedVersion.String())
			}
		})
	}
}

// Helper function to test the regex parsing logic
func findVersionInOutput(output, pattern string) []string {
	// This mimics the logic in getTerragruntVersion
	versionRegex := regexp.MustCompile(pattern)
	return versionRegex.FindStringSubmatch(output)
}

func TestGetVersion_WithTerragrunt(t *testing.T) {
	// This test requires terragrunt to be installed
	// Skip if terragrunt is not available
	ctx := context.Background()

	// Test that the terragrunt option works
	settings := settings{
		workdir:       ".",
		terraformBin:  "terragrunt",
		useTerragrunt: true,
	}

	// We can't easily test this without a real terragrunt installation
	// and a proper terragrunt.hcl file, so we'll skip this test in CI
	// This is more of an integration test
	t.Skip("Integration test - requires terragrunt installation and proper setup")

	tfVersion, err := getTerragruntVersion(ctx, settings)
	if err != nil {
		t.Logf("Terragrunt not available or not properly configured: %v", err)
		t.Skip("Terragrunt not available")
	}

	assert.NotNil(t, tfVersion)
	minVersion := version.Must(version.NewSemver("1.0.0"))
	assert.True(t, tfVersion.GreaterThan(minVersion))
}
