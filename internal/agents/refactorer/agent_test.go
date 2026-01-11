package refactorer

import (
	"testing"

	"ai-multi-agent-workflow-generator/internal/agent/agenttest"
)

func TestSpec(t *testing.T) {
	agenttest.VerifySpec(t, &Agent{})
}

func TestRun(t *testing.T) {
	agenttest.VerifyRun(t, &Agent{})
}
