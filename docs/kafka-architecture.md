# Kafka-Based Agent Messaging Architecture

This framework can scale beyond in-process orchestration by adopting Kafka for agent messaging. The manager becomes a producer/consumer that publishes task steps and collects results.

## Topics

- `agent.tasks`: manager publishes task requests keyed by agent name
- `agent.results`: agents publish results keyed by task id
- `agent.telemetry`: agents emit metrics, timing, and diagnostics
- `agent.deadletter`: failed messages that exceeded retry policy

## Message Schema (JSON)

### Task message

```
{
  "task_id": "distill-001",
  "agent": "idea-distiller",
  "input": {
    "task": "Distill the goal for a minimal Note API.",
    "artifacts": {},
    "parameters": {}
  },
  "trace": {
    "project": "note-api",
    "attempt": 1,
    "created_at": "2025-01-10T18:00:00Z"
  }
}
```

### Result message

```
{
  "task_id": "distill-001",
  "agent": "idea-distiller",
  "status": "success",
  "output": {
    "result": "idea-distiller completed: Distill the goal for a minimal Note API.",
    "artifacts": {"summary": "..."},
    "diagnostics": []
  },
  "duration_ms": 120,
  "trace": {
    "project": "note-api",
    "attempt": 1,
    "completed_at": "2025-01-10T18:00:00Z"
  }
}
```

## Orchestration Flow

1. Manager publishes ordered steps to `agent.tasks` with a shared correlation id.
2. Agents subscribe only to their `agent` name and claim matching tasks.
3. Results are published to `agent.results`; the manager aggregates in order.
4. Retries publish the same task id with incremented attempt.
5. Failures over policy are sent to `agent.deadletter` for review.

## Reliability

- Use idempotent producers for the manager.
- Add a retry policy with exponential backoff in each agent.
- Commit offsets only after results are published.

## Security

- Enforce ACLs per topic and agent consumer group.
- Encrypt traffic with TLS and rotate credentials regularly.
