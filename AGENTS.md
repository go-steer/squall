# AGENTS.md — squall

Instructions for AI agents and contributors working in this repository.

## Start here

- [`README.md`](./README.md) for what `squall` is.
- [`docs/DESIGN.md`](./docs/DESIGN.md) for the design of record: leak budget
  (`T1`–`T9`), executor invariants, `squalld` node agent, unified Engine API +
  sessions + SDKs, `core-agent/v3` planner + gVisor scenario sandbox, and the
  `k8s-lookout` coverage matrix.
- [`docs/decisions.md`](./docs/decisions.md) for architectural decisions; record
  any design change as a new entry in the same PR.

## Hard rules (violations are bugs)

- **Every fault traverses `pkg/executor`.** (`validate -> safety -> reserve ->
  lease -> narrow -> precheck -> driver -> confirm -> settle -> recovery ->
  audit`). Never add a CLI flag, MCP tool, SDK call, or scenario-sandbox path
  that mutates a cluster or calls `squalld` around the executor.
- **Deterministic engine independence.** `squall` (the engine) and `squalld`
  (the node agent) must never import `core-agent` or `core-models`. Ground
  truth in the audit trail comes strictly from applied and confirmed engine
  state, never from LLM claims.
- **Zero `K` and `F` leak budget at wildcard-read RBAC.** Presence of chaos
  furniture (`P`) is allowed; fault kind/target (`K`) and parameters (`F`) in
  API objects, core/`events.k8s.io` `Event`s, target pod metadata/specs, and
  reachable `pods/log` streams must be zero (`T1`–`T9` in
  [`docs/DESIGN.md`](./docs/DESIGN.md)). `squalld` logs opaque `fault_id`s
  only, and the engine audit stream never goes to target-cluster stdout.
- **Every fault is reversible and bounded.** Every fault requires a lease with
  a hard deadline, an idempotent driver `Clear`, `ListLive` adoption across
  engine restarts, and node-local timer/kernel backstops on `squalld`.
- **Ask before creating billable cloud resources or touching a live cluster.**
  Keep `KUBECONFIG` scoped explicitly; never echo, log, or commit credentials,
  tokens, or signing keys.

## Build & test

```bash
dev/tools/ci                                  # full local presubmit sweep
dev/ci/presubmits/build                       # go build ./...
dev/ci/presubmits/test-unit                   # go test -race ./...
dev/ci/presubmits/vet                         # go vet ./...
dev/ci/presubmits/lint-go                     # golangci-lint v2.12.1
dev/ci/presubmits/verify-go-format            # gofmt + goimports (fix: dev/tools/fix-go-format)
dev/ci/presubmits/verify-mod-tidy             # go.mod/go.sum are tidy
dev/ci/presubmits/verify-coverage             # per-package floors (dev/coverage-floors.txt)
dev/ci/presubmits/verify-vuln                 # govulncheck ./...
dev/ci/presubmits/verify-no-agent-attribution # no assistant attribution
```

Unit tests must run offline with no cluster and no external credentials.

## Conventions

- **License headers everywhere.** Every Go, shell, YAML, Python, and Dockerfile
  source file carries the 13-line Google LLC Apache 2.0 header at the top
  (never Markdown). `goheader` enforces it on `.go`; run
  `dev/tools/add-license-headers` for the rest.
- **Feature worktrees under `.worktrees/`.** Create isolated worktrees via
  `git worktree add .worktrees/<slug> -b <branch> origin/main`.
- **Design requirement trailers (`{#PREFIX-N}`).** Commits implementing an
  anchored requirement in [`docs/DESIGN.md`](./docs/DESIGN.md) carry:
  ```text
  Design-Item: docs/DESIGN.md#<ID>
  Design-Rev: <40-char-commit-sha-of-doc>
  ```
  (and `Design-Removed: docs/DESIGN.md#<ID>` when retiring one).
- **Conventional Commits + DCO sign-off (`git commit -s`).**
- **No `Co-Authored-By` trailer and no assistant attribution** in commits, PR
  titles/bodies, or committed files (`agent-attribution` CI enforces this).
- **Per-package coverage floors.** Every Go package requires an entry in
  [`dev/coverage-floors.txt`](./dev/coverage-floors.txt).

## How we develop

Single long-lived branch `main`; short-lived `feat/…`, `fix/…`, `chore/…`,
`docs/…` branches → PR → merge once required CI checks (`ci`, `review-gate`,
`agent-attribution`) are green. Never autonomously merge PRs that modify
`.github/`, `dev/ci/`, `docs/DESIGN.md`, `docs/*-design.md`, `docs/design/`,
or `docs/decisions.md`; leave them open for human review.

- **Run `dev/tools/ci` before every push.**
- **Adversarial-review gate before every PR touching Go code.** Run a read-only
  skeptic subagent review over the diff (checking design-of-record conformance,
  beyond-the-diff callers, races, and error paths; for bug fixes, verify the
  regression test fails on pre-fix code), resolve every finding, and record the
  result under `## Adversarial review` in the PR body.
