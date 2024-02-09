package terraform

import (
	"fmt"
	"os"
	"os/exec"
)

type settings struct {
	workdir      string
	terraformBin string
	skipInit     bool
	skipRefresh  bool
}

// An Option configures how Terraform commands are run.
type Option func(*settings)

func defaultOptions() []Option {
	return []Option{
		WithWorkdir("."),
		WithTerraformBin("terraform"),
	}
}

// WithWorkdir sets the directory Terraform commands will run. If this option
// is not provided, it defaults to the current working directory.
func WithWorkdir(workdir string) Option {
	return func(s *settings) {
		s.workdir = workdir
	}
}

// WithTerraformBin replaces the Terraform binary with the provided executable.
// If this option is not provided, it defaults to the `terraform` executable in
// the current PATH.
func WithTerraformBin(path string) Option {
	return func(s *settings) {
		s.terraformBin = path
	}
}

// WithSkipInit configures whether the init step should be skipped. By default,
// this step is not skipped.
//
// Skipping the init step can save time, but subsequent steps may fail if the
// module was not initialized beforehand.
func WithSkipInit(skipInit bool) Option {
	return func(s *settings) {
		s.skipInit = skipInit
	}
}

// WithSkipRefresh configures whether the refresh step should be skipped. By
// default, this step is not skipped.
//
// Skipping the refresh step can save time, but can result in Terraform basing
// its plan on stale data.
func WithSkipRefresh(skipRefresh bool) Option {
	return func(s *settings) {
		s.skipRefresh = skipRefresh
	}
}

func (s *settings) apply(opts []Option) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *settings) validate() error {
	if !isDirectory(s.workdir) {
		return fmt.Errorf("target directory %q not found", s.workdir)
	}

	if !isInPath(s.terraformBin) {
		return fmt.Errorf("executable %q not found in PATH", s.terraformBin)
	}

	return nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func isInPath(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}
