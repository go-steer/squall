# Design & Documentation Index

This directory holds `squall`'s architectural design of record and decision log.

## Lifecycle rules

Every design document carries `**Status:** proposed | shipped | superseded | abandoned`
near the top, a **Settled decisions (do not relitigate)** section, and an
**Out of scope** list. Buildable requirements carry stable `{#PREFIX-N}`
heading anchors tracked in git commit history via `Design-Item:` and
`Design-Rev:` trailers.

## Active documents

| Document | Status | Scope |
|---|---|---|
| [`DESIGN.md`](./DESIGN.md) | `proposed` | Founding architecture: zero-leak bar (`T1`–`T9`), external `squall` engine, self-reverting `squalld` node agent, unified protobuf Engine API + sessions + Go/Python SDKs, `core-agent/v3` planner + gVisor scenario sandbox, and `k8s-lookout` coverage matrix. |
| [`decisions.md`](./decisions.md) | `active` | Chronological log of settled architectural decisions and trade-offs. |
