package terraform

import (
	"context"
	"fmt"

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
