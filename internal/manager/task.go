package manager

import "ai-multi-agent-workflow-generator/internal/agent"

// TaskFile describes the YAML task document consumed by the manager.
type TaskFile struct {
	Version string `yaml:"version"`
	Project string `yaml:"project"`
	Steps   []Step `yaml:"steps"`
}

// Step captures one agent execution entry.
type Step struct {
	ID    string      `yaml:"id"`
	Agent string      `yaml:"agent"`
	Input agent.Input `yaml:"input"`
}

// StepResult captures outputs per step.
type StepResult struct {
	ID     string
	Agent  string
	Output agent.Output
}
