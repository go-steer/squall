# dev/

Build- and test-tooling. The same scripts power both local development and
GitHub Actions CI, so a green local run is the same green run as remote.

## Quickstart

```bash
# Run every CI check locally (fast-fail order).
dev/tools/ci

# Run all checks even after a failure (collect every problem at once).
dev/tools/ci --keep-going

# Auto-fix formatting (gofmt + goimports).
dev/tools/fix-go-format
```

Missing tools (`golangci-lint`, `goimports`, `govulncheck`) auto-install into
`$GOBIN` (or `$(go env GOPATH)/bin`) on first use. No setup needed beyond a Go
toolchain.

## Layout

```
dev/
├── coverage-floors.txt    # per-package statement-coverage floors
├── claude/                # opt-in Claude Code config sample (not live config)
├── tools/                 # entry points you run locally
│   ├── ci                 # aggregator — runs every check below
│   ├── verify-no-agent-attribution # no assistant attribution in commits
│   ├── verify-go-format   # gofmt -s + goimports check (read-only)
│   ├── fix-go-format      # gofmt -s -w + goimports -w (auto-fix)
│   ├── vet                # go vet ./...
│   ├── build              # go build ./...
│   ├── lint-go            # golangci-lint (auto-installs v2.12.1)
│   ├── verify-mod-tidy    # `go mod tidy` clean check
│   ├── test-unit          # go test -race -coverprofile
│   ├── verify-coverage    # each package at or above its floor
│   ├── verify-vuln        # govulncheck ./...
│   ├── add-license-headers # bulk-applier for the Apache 2.0 header
│   ├── common.sh          # shared bash helpers (ensure_tool, run_step)
│   └── .golangci.yml      # linter config
└── ci/
    └── presubmits/        # thin delegators called by .github/workflows/
        ├── verify-no-agent-attribution
        ├── verify-go-format
        ├── vet
        ├── build
        ├── lint-go
        ├── verify-mod-tidy
        ├── verify-vuln
        ├── verify-coverage
        └── test-unit
```

## Adding a check

1. Drop a new script under `dev/tools/<name>` (executable, `set -euo pipefail`,
   sources `common.sh`).
2. Add it to the `STEPS` array in `dev/tools/ci`.
3. Add a one-line delegator under `dev/ci/presubmits/<name>` that `exec`s the
   tool script.

The `ci` workflow runs `dev/tools/ci`, so a new step is picked up without
touching any YAML.

## CI on PRs

`.github/workflows/ci.yml` runs the whole sweep in one `ci` job on every push
to `main` and every PR. Branch protection on `main` requires:

- `ci`
- `review-gate`
- `agent-attribution`

## License headers

Every source file carries the full Apache 2.0 header at the top, attributed to
Google LLC (`//`-prefixed for Go, `#`-prefixed for shell, YAML, Python, and
Dockerfiles; never on Markdown). `goheader` inside `dev/tools/lint-go` enforces
it on `.go` files; run `dev/tools/add-license-headers` for the rest.
