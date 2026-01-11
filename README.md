# Multi-Agent Development Framework

A Go-based, Unix-like multi-agent framework where each agent performs a single job with clear input/output contracts. The manager orchestrates agents using YAML task files.

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

- Kafka messaging architecture: `docs/kafka-architecture.md`
- Self-modifying improvement loop: `docs/self-modifying-loop.md`
- Additional agents list: `docs/additional-agents.md`
