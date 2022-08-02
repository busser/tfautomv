package terraform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/Masterminds/semver/v3"
)

type runner struct {
	workdir   string
	stateFile string
}

func NewRunner(workdir string) *runner {
	r := runner{
		workdir: workdir,
	}

	return &r
}

func (r *runner) Init() error {
	return r.runCommand([]string{"init"}, nil)
}

func (r *runner) Plan() (*Plan, error) {
	planFile, err := ioutil.TempFile("", "tfautomv.*.plan")
	if err != nil {
		return nil, err
	}
	defer os.Remove(planFile.Name())

	if err := r.runCommand([]string{"plan", "-out", planFile.Name()}, nil); err != nil {
		return nil, err
	}

	var jsonPlan bytes.Buffer
	if err := r.runCommand([]string{"show", "-json", planFile.Name()}, &jsonPlan); err != nil {
		return nil, err
	}

	var plan Plan
	if err := json.Unmarshal(jsonPlan.Bytes(), &plan); err != nil {
		return nil, fmt.Errorf("could not parse plan: %w", err)
	}

	return &plan, nil
}

func (r *runner) Apply() error {
	return r.runCommand([]string{"apply", "-auto-approve"}, nil)
}

func (r *runner) Version() (*semver.Version, error) {
	var jsonVersion bytes.Buffer
	if err := r.runCommand([]string{"version", "-json"}, &jsonVersion); err != nil {
		return nil, err
	}

	var version Version
	if err := json.Unmarshal(jsonVersion.Bytes(), &version); err != nil {
		return nil, fmt.Errorf("could not parse version: %w", err)
	}

	return semver.NewVersion(version.TerraformVersion)
}

func (r *runner) runCommand(args []string, out io.Writer) error {
	cmd := exec.Command("terraform", args...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	if out != nil {
		cmd.Stdout = out
	}
	cmd.Stderr = &buf
	cmd.Dir = r.workdir

	if err := cmd.Run(); err != nil {
		return Error{cmd, &buf, err}
	}

	return nil
}
