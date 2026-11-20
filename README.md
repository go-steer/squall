# squall

Agent-driven chaos engineering for Kubernetes: a deterministic fault engine
(`squall`), a self-reverting node agent (`squalld`), and an intelligent chaos
planner and scenario coder (`squall-agent`) that inject bounded faults into
opt-in arenas without leaving a readable answer in the cluster.

## Why squall exists

Existing Kubernetes chaos operators reconcile custom resources and emit core
`v1` `Event`s into the target namespace. When an autonomous SRE agent with
broad read access investigates the resulting incident, listing namespace Events
hands it the fault's kind, target, and parameters instead of testing real
diagnosis.

`squall` closes that leak at the root:

- **Zero per-fault API footprint for runtime faults:** Pod loss uses a plain
  Kubernetes `pods delete`. Network delay, packet loss, L3/L4 partition,
  container kill, and CPU/memory pressure run through `squalld` — a node
  DaemonSet (`hostPID`, four capabilities, no container runtime socket, no
  Kubernetes API token) driven from outside the target cluster via `purser`
  mTLS and Ed25519-signed command leases.
- **Self-reverting by construction:** Every fault carries a lease deadline.
  `squalld` persists active leases on node-local `tmpfs`, arms a local revert
  timer, and backs network partitions with `nftables` kernel set timeouts so
  partitions clear even if both the engine and `squalld` die mid-fault.
- **Unified Engine API, MCP & SDKs:** One protobuf-defined Engine contract
  (asynchronous submit, observe, wait, structured dry-run, clear, and
  heartbeated scenario sessions) powers both MCP tools and generated Go and
  Python SDKs for adaptive, multi-fault scenario programs.
- **Agent on `core-agent/v3` + `core-models`:** Interactive prompt-driven
  chaos and unattended periodic campaigns run on
  [`go-steer/core-agent`](https://github.com/go-steer/core-agent) and
  [`go-steer/core-models`](https://github.com/go-steer/core-models), with
  agent-authored scenario code compiled and executed inside network-fenced
  gVisor sandboxes ([`sigs.k8s.io/agent-sandbox`](https://github.com/kubernetes-sigs/agent-sandbox))
  and startup enforcement of model-family diversity against the SRE agent under
  test.

## Status

Pre-Milestone 1 — founding design under review in
[`docs/DESIGN.md`](./docs/DESIGN.md). Decision log:
[`docs/decisions.md`](./docs/decisions.md).

## Development

```bash
dev/tools/ci              # run every CI presubmit check locally
dev/tools/fix-go-format   # gofmt -s + goimports
```

See [`AGENTS.md`](./AGENTS.md) and [`CONTRIBUTING.md`](./CONTRIBUTING.md) for
invariants, design requirement trailers, and the adversarial review gate.

## License

[Apache License 2.0](./LICENSE).
