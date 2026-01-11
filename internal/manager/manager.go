package manager

import (
	"context"
	"fmt"
	"os"

	"ai-multi-agent-workflow-generator/internal/agent"
	"gopkg.in/yaml.v3"
)

// Manager orchestrates agent execution based on a YAML task file.
type Manager struct {
	registry *agent.Registry
}

// New creates a new manager with a registry.
func New(registry *agent.Registry) *Manager {
	return &Manager{registry: registry}
}

// LoadTaskFile reads a YAML task file from disk.
func (m *Manager) LoadTaskFile(path string) (TaskFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TaskFile{}, fmt.Errorf("read task file: %w", err)
	}

	var task TaskFile
	if err := yaml.Unmarshal(data, &task); err != nil {
		return TaskFile{}, fmt.Errorf("parse task file: %w", err)
	}

	return task, nil
}

// Run executes steps sequentially and returns results.
func (m *Manager) Run(ctx context.Context, task TaskFile) ([]StepResult, error) {
	results := make([]StepResult, 0, len(task.Steps))
	var previous agent.Output

	for _, step := range task.Steps {
		a, ok := m.registry.Create(step.Agent)
		if !ok {
			return nil, fmt.Errorf("unknown agent: %s", step.Agent)
		}

		if step.Input.Artifacts == nil {
			step.Input.Artifacts = map[string]string{}
		}
		if previous.Result != "" {
			step.Input.Artifacts["pipeline.previous_result"] = previous.Result
		}

		output, err := a.Run(ctx, step.Input)
		if err != nil {
			return nil, fmt.Errorf("step %s (%s) failed: %w", step.ID, step.Agent, err)
		}

		results = append(results, StepResult{ID: step.ID, Agent: step.Agent, Output: output})
		previous = output
	}

	return results, nil
}
