# Architectural Decisions

Chronological record of settled decisions for `squall`. New decisions or
amendments to [`DESIGN.md`](./DESIGN.md) add an entry here in the same PR.

## 2026-11-20 — Standalone repository (`go-steer/squall`), leaving `simian-agent` intact

**Decision:** Build `squall` in a new public repository (`github.com/go-steer/squall`)
rather than retrofitting `go-steer/simian-agent` in place.

**Why:** `simian-agent` serves an established role around Chaos Mesh and its
`oracle_read` leak-detection benchmark ("does an SRE agent read the answer when
it is sitting in the cluster?"). `squall` makes incompatible structural moves:
no Chaos Mesh dependency, a privileged self-reverting node DaemonSet
(`squalld`), controller placement outside the target cluster, a protobuf Engine
API with session-scoped Go/Python SDKs, and an agent layer built on
`core-agent/v3` + `core-models`. Proven non-Chaos-Mesh packages from
`simian-agent` (`executor`, `arena`, `lease`, `probe`, `audit`, `kubestate`,
`networkpolicy`) transfer cleanly with attribution.

## 2026-11-20 — Own self-reverting node agent (`squalld`), not direct `chaos-daemon` RPCs

**Decision:** Replace Chaos Mesh runtime injection with plain API `pods delete`
plus `squalld`, a `squall`-owned DaemonSet (`hostPID`, `CAP_NET_ADMIN`,
`CAP_SYS_ADMIN`, `CAP_SYS_PTRACE`, `CAP_KILL`, no container runtime socket, no
Kubernetes API token).

**Why:** A live spike on `shipmind-stage` (2026-11-19) proved that driving
`chaos-daemon` directly injects netem delay and iptables/nftables partition with
zero Kubernetes objects and zero Events across 169 resource types. However,
`chaos-daemon`'s 16 gRPC methods carry no TTL or deadline (violating bounded
self-reversion if the controller dies), mount the node's container runtime
socket, and print full fault parameters to stdout (`pods/log`). `squalld`
enforces node-local `tmpfs` lease timers, `nftables` kernel set timeouts, and
parameter-free logs.

## 2026-11-20 — Leak bar measured at wildcard-read RBAC (`R`), controller outside target cluster

**Decision:** Hold runtime fault injection to `K = 0` (no fault kind or target)
and `F = 0` (no fault parameters) against an investigator with read access to
every Kubernetes resource kind (`R`), and run the `squall` engine outside the
target cluster.

**Why:** In the 2026-11-19 `shipmind-stage` benchmark (27 faults, 12 detected,
5 correct diagnoses), 3 of the 5 correct diagnoses were contaminated by Chaos
Mesh leaks — and in every leaked session, the agent discovered the fault via
ordinary core `v1` `Event` listings before querying any `chaos-mesh.org` object.
Scoping RBAC to hide custom resources does not hide Events or `pods/log`.

## 2026-11-20 — Build `squall-agent` on `core-agent/v3` + `core-models`, never `mast` or Yaegi

**Decision:** Embed `github.com/go-steer/core-agent/v3` and
`github.com/go-steer/core-models` under four independence rules (zero LLM
imports in the engine, engine-only ground truth, fixed scenario packs alongside
agent faults, and startup enforcement of model-family diversity against the
SUT). Execute agent-authored scenario code in prewarmed gVisor sandboxes
(`sigs.k8s.io/agent-sandbox`) rather than an in-process Go interpreter.

**Why:** `go-steer/mast` was frozen at `v0.10.0` and folded into
`core-agent/v3`, which `shipmind` now runs in production. Yaegi was tested in
`kode-gopher` and rejected after 16/32 differential failures (7 silent wrong
outputs); prewarmed `$GOCACHE` gVisor sandboxes compile and run real Go against
the `squall` SDK in ~4 seconds with zero Kubernetes API access.

## 2026-11-20 — Keep every fault in the catalog; score detection separately from diagnosis

**Decision:** Retain faults that `k8s-lookout`'s sentinel does not currently
detect (such as pure network delay), report no-incident faults as
`not_detected` rather than counting them against diagnosis accuracy, and expose
per-kind detection status to the chaos planner.

**Why:** Pruning undetected faults hides observability blind spots instead of
measuring them, and those faults still exercise workload resilience and cascade
scenarios.
