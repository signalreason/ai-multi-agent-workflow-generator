package localizationmanager

import (
	"context"
	"fmt"
	"strings"

	"ai-multi-agent-workflow-generator/internal/agent"
)

const agentName = "localization-manager"

var spec = agent.Spec{
	Name:         agentName,
	Purpose:      "Coordinates localization assets and review workflows.",
	Interface:    "Input: task:string, artifacts:map[string]string, parameters:map[string]string -> Output: result:string, artifacts:map[string]string, diagnostics:[]string",
	InputSpec:    "task: instruction for the agent; artifacts: upstream outputs; parameters: tuning knobs",
	OutputSpec:   "result: primary outcome; artifacts: derived assets; diagnostics: warnings or notes",
	FailureModes: []string{"missing task input", "invalid parameters", "required artifact not provided", "validation failed"},
	Tests:        []string{"spec completeness validation", "run returns result with summary", "run rejects empty task"},
	Version:      "v1",
}

// Agent implements the agent.Agent interface.
type Agent struct{}

func (a *Agent) Spec() agent.Spec {
	return spec
}

func (a *Agent) Run(ctx context.Context, in agent.Input) (agent.Output, error) {
	_ = ctx
	if strings.TrimSpace(in.Task) == "" {
		return agent.Output{}, agent.AgentError{
			Code:    "INVALID_INPUT",
			Message: "task is required",
		}
	}

	summary := fmt.Sprintf("%s completed: %s", agentName, in.Task)
	return agent.Output{
		Result:    summary,
		Artifacts: map[string]string{"summary": summary},
	}, nil
}

// Register wires the agent into a registry.
func Register(r *agent.Registry) {
	r.Register(agentName, func() agent.Agent { return &Agent{} })
}
