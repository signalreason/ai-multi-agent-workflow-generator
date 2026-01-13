# Multi-Agent Development Framework

A Go-based, Unix-like multi-agent framework where each agent performs a single job with clear input/output contracts. The manager orchestrates agents using YAML task files.

## Purpose
- Provide a Go-based framework for running multi-agent workflows defined in YAML.

## Goals
- Keep agent interfaces explicit and testable with clear input/output contracts.
- Orchestrate sequential task execution while passing prior results forward.
- Offer example projects and docs that demonstrate extensible agent patterns.

## Highest-Impact Next Step
- Add a formal YAML task schema validator with clear error messages and unit tests.

## Checks
- Status: none (no GitHub Actions runs found).
- TODO: Add CI for `go test ./...` and `go vet ./...`.
- TODO: Add a formatting check using `gofmt -l` on Go sources.

## Core Concepts

- **Unix-like agents**: each agent performs one job and emits structured output.
- **Explicit specs**: every agent exposes a spec with interface definition, failure modes, and test suite.
- **Composable workflows**: the manager executes YAML task steps sequentially and injects previous results.

## CLI Usage

List available agents:

```
go run ./cmd/manager -list-agents
```

Run a task file:

```
go run ./cmd/manager -task examples/tasks/example.yaml
```

## Agent Interface

All agents implement:

```
Run(ctx, Input) -> Output
```

See `internal/agent/agent.go` for the shared contract, and each agent package for its spec.

## Example Project

A demo Note API lives at `examples/project`.

```
cd examples/project

go run ./cmd/noteapi
```

## Tests

```
go test ./...
```

## Docs

- Repo map: `docs/REPO_MAP.md`
- Kafka messaging architecture: `docs/kafka-architecture.md`
- Self-modifying improvement loop: `docs/self-modifying-loop.md`
- Additional agents list: `docs/additional-agents.md`
