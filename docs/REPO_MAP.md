# Repo Map: ai-multi-agent-workflow-generator

## Purpose and scope
Go-based multi-agent framework where a manager runs YAML task files to orchestrate agents with explicit input/output contracts.

## Quickstart commands
- List agents: `go run ./cmd/manager -list-agents`
- Run a task: `go run ./cmd/manager -task examples/tasks/example.yaml`
- Tests: `go test ./...`

## Top-level map
- `cmd/` - CLI entry points.
  - `cmd/manager/` - manager CLI main.
- `internal/` - core agent framework, contracts, and manager logic.
- `examples/` - sample task files and demo project.
  - `examples/tasks/` - YAML tasks.
  - `examples/project/` - demo Note API.
- `docs/` - architecture notes and agent references.
- `go.mod`, `go.sum` - Go module definition.
- `README.md` - usage and docs index.

## Key entry points
- `cmd/manager/` - main CLI for running tasks and listing agents.
- `internal/agent/agent.go` - agent interface contract.
- `examples/tasks/example.yaml` - sample task definition.
- `examples/project/cmd/noteapi` - demo service entry point.

## Core flows and data movement
- Manager loads a YAML task file and executes steps sequentially.
- Each agent runs `Run(ctx, Input) -> Output` and passes outputs to later steps.
- Optional Kafka-based messaging architecture described in docs.

## External integrations
- Optional Kafka messaging (see `docs/kafka-architecture.md`).

## Configuration and deployment
- YAML task files in `examples/tasks/` show expected task shape.
- No deployment manifests included; this repo is CLI-first.

## Common workflows (build/test/release)
- `go run ./cmd/manager -list-agents`
- `go run ./cmd/manager -task <task.yaml>`
- `go test ./...`

## Read-next list
- `README.md`
- `cmd/manager/`
- `internal/agent/agent.go`
- `docs/kafka-architecture.md`
- `docs/self-modifying-loop.md`
- `examples/tasks/`

## Unknowns and follow-ups
- No formal release or packaging process described.
