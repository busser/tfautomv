# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

- **Build**: `make build` - Builds the tfautomv binary to `bin/tfautomv`
- **Format**: `make fmt` - Formats Go source code using `go fmt`
- **Vet**: `make vet` - Runs `go vet` for static analysis
- **Unit Tests**: `make test` - Runs unit tests in `pkg/` with coverage
- **E2E Tests**: `make test-e2e` - Runs end-to-end tests in `test/e2e/`
- **Full Test**: `make test test-e2e` - Runs both unit and e2e tests
- **Help**: `make help` - Shows available make targets

## Architecture Overview

tfautomv is a Terraform refactoring tool that automatically generates `moved` blocks and `terraform state mv` commands when resources are moved in code.

### Core Components

1. **Engine** (`pkg/engine/`): Core business logic
   - `Plan`: Represents Terraform plans with resources to create/delete
   - `Resource`: Represents a Terraform resource with flattened attributes
   - `Move`: Represents a resource move between addresses/modules
   - `ResourceComparison`: Compares create/delete pairs to detect moves
   - Rules system for ignoring specific attribute differences

2. **Terraform Integration** (`pkg/terraform/`): Terraform CLI interaction
   - Executes `terraform init`, `refresh`, and `plan` commands
   - Parses JSON plan output using hashicorp/terraform-json
   - Generates HCL `moved` blocks and shell `terraform state mv` commands

3. **Rules System** (`pkg/engine/rules/`): Configurable difference ignoring
   - `everything`: Ignore any difference in an attribute
   - `whitespace`: Ignore whitespace differences 
   - `prefix`: Ignore specific prefixes in attribute values
   - Extensible rule parsing and application system

4. **Pretty Printing** (`pkg/pretty/`): User output formatting
   - Colored terminal output with box formatting
   - Verbosity levels (0-3) for debugging move decisions
   - Summary generation showing matched/unmatched resources

### Data Flow

1. Parse CLI flags and user-provided ignore rules
2. Run `terraform plan` for each working directory
3. Extract create/delete resources from JSON plan using `SummarizeJSONPlan`
4. Flatten resource attributes using `flatmap` package
5. Compare all create/delete pairs using rules to find matches
6. Generate moves for uniquely matched resource pairs
7. Output as `moved` blocks (same module) or `terraform state mv` commands (cross-module)

### Testing

- Unit tests use table-driven patterns with golden files for expected outputs
- Golden files in `pkg/*/testdata/` directories store expected test outputs
- E2E tests in `test/e2e/` test full CLI workflows
- Use `github.com/busser/tfautomv/pkg/golden` for golden file management

## Key Implementation Details

- Resources are compared using flattened attribute maps for efficient diffing
- Moves are only generated for 1:1 matches to avoid ambiguous state operations
- Cross-module moves require Terraform v0.14+ and generate shell commands
- `moved` blocks require Terraform v1.1+ and are written to `moves.tf` files
- Concurrent plan fetching for multiple working directories using goroutines

## Release Process

tfautomv uses a structured release process with organized release notes and automated tooling.

### Release Commands

- **Release Dry Run**: `make release-dry-run` - Test the complete release process without publishing
- **Release**: `make release` - Create and publish a new release (requires `gh` CLI authentication)

### Release Workflow

1. **Prepare Release Notes**: Create `docs/release-notes/vX.Y.Z.md` with comprehensive release notes
2. **Create Release Branch**: `git checkout -b release/vX.Y.Z`
3. **Update VERSION**: Change `VERSION` file to target version (e.g., `v0.7.0`)
4. **Commit and PR**: `git commit -m "release vX.Y.Z"` and open PR
5. **Merge and Release**: After PR merge, checkout main, pull, and run `make release`

### Release Notes Format

Release notes in `docs/release-notes/` follow this structure:
- **Emoji-prefixed sections** (🔥, 📋, 📚, 🔧)
- **Brief descriptions** with **code examples**
- **Links to documentation** for detailed information
- **User-focused language** highlighting benefits

### Automated Release Features

- **Multi-platform binaries**: Linux, macOS, Windows (amd64, arm64, 386)
- **Package management**: Automatic Homebrew and AUR updates
- **GitHub integration**: Automated release creation with binaries and checksums
- **Release notes**: Automatically included from `docs/release-notes/$(VERSION).md`

### Dependencies

- **GoReleaser v2**: Handles cross-compilation and publishing
- **GitHub CLI (`gh`)**: Provides authentication token for releases
- **Git tags**: Version tags trigger GoReleaser's release process