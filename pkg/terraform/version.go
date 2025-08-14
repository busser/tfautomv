package terraform

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// GetVersion obtains the version of the configured Terraform binary.
func GetVersion(ctx context.Context, opts ...Option) (*version.Version, error) {
	var settings settings

	settings.apply(append(defaultOptions(), opts...))

	err := settings.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	// Use custom terragrunt version parsing if terragrunt mode is enabled
	if settings.useTerragrunt {
		return getTerragruntVersion(ctx, settings)
	}

	tf, err := tfexec.NewTerraform(settings.workdir, settings.terraformBin)
	if err != nil {
		return nil, fmt.Errorf("failed to create Terraform executor: %w", err)
	}

	version, _, err := tf.Version(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("\"terraform version\" failed: %w", err)
	}

	return version, nil
}

// getTerragruntVersion obtains the version of terragrunt using custom parsing
// that handles terragrunt's output format and deprecation warnings.
func getTerragruntVersion(ctx context.Context, settings settings) (*version.Version, error) {
	// Use "terragrunt run -- version" to avoid deprecation warning
	cmd := exec.CommandContext(ctx, settings.terraformBin, "run", "--", "version")
	cmd.Dir = settings.workdir

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run terragrunt version: %w", err)
	}

	// Parse the version from the output
	// Expected format: "Terraform v1.12.2\non darwin_arm64\n"
	versionRegex := regexp.MustCompile(`Terraform v(\d+\.\d+\.\d+)`)
	matches := versionRegex.FindStringSubmatch(string(output))

	if len(matches) < 2 {
		return nil, fmt.Errorf("unable to parse version from terragrunt output: %s", strings.TrimSpace(string(output)))
	}

	versionStr := matches[1]
	parsedVersion, err := version.NewSemver(versionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version %q: %w", versionStr, err)
	}

	return parsedVersion, nil
}
