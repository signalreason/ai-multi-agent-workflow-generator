# Self-Modifying Improvement Loop

The framework can run a guarded improvement loop to evolve agent behavior based on results and test feedback.

## Loop Stages

1. **Observe**: Manager aggregates outputs, metrics, and test results per task.
2. **Diagnose**: Diagnostic agent flags gaps, regressions, and low-confidence outputs.
3. **Plan**: Refactorer and TDD enforcer propose targeted improvements.
4. **Change**: Code generator writes new changes into a staging branch.
5. **Verify**: Benchmark tester and test suites must pass with regression checks.
6. **Promote**: Release packager tags a new version once verified.

## Guardrails

- All self-modifying edits require a staging branch and explicit approval.
- Changes must pass unit tests and a curated risk checklist.
- Rollback artifacts are stored for each successful cycle.

## Suggested YAML Loop Task

```
version: "1"
project: "agent-framework"
steps:
  - id: observe
    agent: log-analyzer
    input:
      task: "Collect result quality signals and failure rates."
  - id: diagnose
    agent: security-auditor
    input:
      task: "Flag risky behavior in recent changes."
  - id: plan
    agent: refactorer
    input:
      task: "Propose refactors for agents with high failure rates."
  - id: change
    agent: code-generator
    input:
      task: "Implement refactor proposals into a staging branch."
  - id: verify
    agent: tdd-enforcer
    input:
      task: "Run the test plan for modified agents."
  - id: package
    agent: release-packager
    input:
      task: "Prepare a versioned artifact if tests pass."
```
