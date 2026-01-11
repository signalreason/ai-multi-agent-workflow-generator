package agent

import (
	"context"
	"fmt"
	"sort"
)

// Input defines the standardized request contract for all agents.
type Input struct {
	Task       string            `yaml:"task"`
	Artifacts  map[string]string `yaml:"artifacts"`
	Parameters map[string]string `yaml:"parameters"`
}

// Output defines the standardized response contract for all agents.
type Output struct {
	Result      string            `yaml:"result"`
	Artifacts   map[string]string `yaml:"artifacts"`
	Diagnostics []string          `yaml:"diagnostics"`
}

// Spec captures the agent specification and documentation.
type Spec struct {
	Name         string
	Purpose      string
	Interface    string
	InputSpec    string
	OutputSpec   string
	FailureModes []string
	Tests        []string
	Version      string
}

// Agent defines the executable contract for all agents.
type Agent interface {
	Spec() Spec
	Run(ctx context.Context, in Input) (Output, error)
}

// AgentError makes failures explicit and machine readable.
type AgentError struct {
	Code      string
	Message   string
	Retryable bool
}

func (e AgentError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Registry holds agent factories keyed by name.
type Registry struct {
	factories map[string]func() Agent
}

// NewRegistry creates a registry for agent factories.
func NewRegistry() *Registry {
	return &Registry{factories: map[string]func() Agent{}}
}

// Register adds an agent factory.
func (r *Registry) Register(name string, factory func() Agent) {
	r.factories[name] = factory
}

// Create instantiates an agent by name.
func (r *Registry) Create(name string) (Agent, bool) {
	factory, ok := r.factories[name]
	if !ok {
		return nil, false
	}
	return factory(), true
}

// Names returns a sorted list of registered agent names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
