# Contributing to squall

Thanks for contributing. `squall` follows the same contributor flow as
[core-agent](https://github.com/go-steer/core-agent),
[shipmind](https://github.com/go-steer/shipmind), and
[purser](https://github.com/go-steer/purser) — a maintainer of one `go-steer`
repository should recognize the others.

## Workflow

- Single long-lived branch: `main`. Work on short-lived feature branches
  (`feat/…`, `fix/…`, `chore/…`, `docs/…`) → PR against `main` → merge
  once CI's required checks are green.
- **Rebase, don't merge.** Keep feature branches rebased on `main`;
  `git push --force-with-lease` on your own branch is normal. Never
  force-push `main`.
- **Feature worktrees live under `.worktrees/`.** Use
  `git worktree add .worktrees/<slug> -b <branch> origin/main` rather than
  switching branches in the root checkout.

## Commits

- **Conventional Commits** subject lines: `feat:`, `fix:`, `docs:`,
  `chore:`, `refactor:`, `test:`, `ci:`, `build:`. Bodies explain *why*
  and call out the verification done.
- **Design requirement trailers** when implementing a `{#PREFIX-N}` section in
  [`docs/DESIGN.md`](./docs/DESIGN.md) (or another design doc):
  ```text
  Design-Item: docs/DESIGN.md#<ID>
  Design-Rev: <40-char-commit-sha-of-doc>
  ```
  Use `Design-Removed: <doc>#<ID>` when removing a retired requirement.
- **DCO sign-off** on every commit: `git commit -s` (adds a
  `Signed-off-by:` trailer certifying the [Developer Certificate of
  Origin](https://developercertificate.org/)).
- **No `Co-Authored-By` trailer, and no assistant attribution anywhere.**
  Author the work under your own name. Do not add `Co-Authored-By:` lines,
  "Generated with …" footers, or any tool/assistant credit to commits, PR
  titles/bodies, or other committed or published artifacts.

## Before you push

- Run `dev/tools/ci` — the full presubmit sweep, the same scripts CI runs
  (`dev/ci/presubmits/*`). A green local run is the green remote run.
- **Adversarial-review gate** on any PR touching Go code — see
  [`AGENTS.md`](./AGENTS.md) "How we develop". Record the outcome under an
  `## Adversarial review` heading in the PR body; the `review-gate`
  required CI check enforces it.

## Tests & coverage floors

Every new package ships with unit tests, and every bug fix ships with a
regression test verified to fail before the fix. Unit tests must require no
network, no cluster, and no external credentials.

`dev/tools/verify-coverage` holds every package to a statement-coverage floor
in [`dev/coverage-floors.txt`](./dev/coverage-floors.txt). Every package in the
coverage profile must have an entry; a new package without a floor fails CI.

## License

By contributing you agree your contributions are licensed under the
project's [Apache 2.0](./LICENSE) license. Every Go / shell / YAML / Python
source file carries the Apache 2.0 header attributed to Google LLC.
