package agenttest

import (
	"context"
	"testing"

	"ai-multi-agent-workflow-generator/internal/agent"
)

// VerifySpec ensures each agent declares required specification fields.
func VerifySpec(t *testing.T, a agent.Agent) {
	t.Helper()
	spec := a.Spec()

	if spec.Name == "" {
		t.Fatal("spec name is required")
	}
	if spec.Purpose == "" {
		t.Fatal("spec purpose is required")
	}
	if spec.Interface == "" {
		t.Fatal("spec interface is required")
	}
	if spec.InputSpec == "" || spec.OutputSpec == "" {
		t.Fatal("input/output specs are required")
	}
	if len(spec.FailureModes) == 0 {
		t.Fatal("failure modes are required")
	}
	if len(spec.Tests) == 0 {
		t.Fatal("test suite entries are required")
	}
}

// VerifyRun ensures agent execution returns structured output without panic.
func VerifyRun(t *testing.T, a agent.Agent) {
	t.Helper()
	out, err := a.Run(context.Background(), agent.Input{Task: "test run"})
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if out.Result == "" {
		t.Fatal("result is required")
	}
}
