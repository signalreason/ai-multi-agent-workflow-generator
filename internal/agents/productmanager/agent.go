package productmanager

import (
	"context"
	"fmt"
	"strings"

	"ai-multi-agent-workflow-generator/internal/agent"
)

const agentName = "product-manager"

var spec = agent.Spec{
	Name:         agentName,
	Purpose:      "Interviews stakeholders to turn limited input into complete, testable product specifications with constraints, edge cases, and acceptance criteria.",
	Interface:    "Input: task:string, artifacts:map[string]string, parameters:map[string]string -> Output: result:string, artifacts:map[string]string, diagnostics:[]string",
	InputSpec:    "task: initial request; artifacts: known fields (objective, users, constraints, etc); parameters: interview_depth, mode, and field values",
	OutputSpec:   "result: spec status; artifacts: spec, questions, acceptance_criteria, risks, open_questions, completion_criteria; diagnostics: assumptions or gaps",
	FailureModes: []string{"missing task input", "conflicting constraints", "insufficient information", "validation failed"},
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

	answers, derived := collectAnswers(in, fieldCatalog)
	if !isAnswered(answers["objective"]) {
		answers["objective"] = normalizeWhitespace(in.Task)
		derived["objective"] = true
	}

	missingRequired, missingOptional := missingFields(answers, derived, fieldCatalog)
	interviewDepth := normalizeWhitespace(in.Parameters["interview_depth"])
	if interviewDepth == "" {
		interviewDepth = "deep"
	}
	mode := strings.ToLower(normalizeWhitespace(in.Parameters["mode"]))

	questions := buildQuestions(in.Task, answers, derived, missingRequired, missingOptional, interviewDepth)
	openQuestions := buildOpenQuestions(missingRequired, missingOptional, interviewDepth)
	acceptanceCriteria := buildAcceptanceCriteria(answers)
	risks := buildRisks(missingRequired)
	completionCriteria := buildCompletionCriteria(missingRequired)
	specBody := buildSpec(in.Task, answers, derived, acceptanceCriteria, risks, completionCriteria)

	status := "draft"
	if len(missingRequired) == 0 {
		status = "complete"
	}
	if mode == "draft-only" {
		status = "draft"
	}

	result := fmt.Sprintf("%s %s spec produced: %s", agentName, status, summarizeTask(in.Task))
	diagnostics := buildDiagnostics(missingRequired, derived)

	return agent.Output{
		Result: result,
		Artifacts: map[string]string{
			"spec":                specBody,
			"questions":           questions,
			"acceptance_criteria": acceptanceCriteria,
			"risks":               risks,
			"open_questions":      openQuestions,
			"completion_criteria": completionCriteria,
		},
		Diagnostics: diagnostics,
	}, nil
}

// Register wires the agent into a registry.
func Register(r *agent.Registry) {
	r.Register(agentName, func() agent.Agent { return &Agent{} })
}

type field struct {
	Key       string
	Label     string
	Required  bool
	Question  string
	ShortHint string
}

var fieldCatalog = []field{
	{
		Key:       "objective",
		Label:     "Objective",
		Required:  true,
		Question:  "What is the primary objective in one sentence (who benefits, what outcome, why now)?",
		ShortHint: "primary outcome in one sentence",
	},
	{
		Key:       "problem_statement",
		Label:     "Problem Statement",
		Required:  false,
		Question:  "What problem are we solving and why now?",
		ShortHint: "problem statement and urgency",
	},
	{
		Key:       "users",
		Label:     "Primary Users",
		Required:  true,
		Question:  "Who are the primary user roles and what are their permissions?",
		ShortHint: "user roles and permissions",
	},
	{
		Key:       "success_metrics",
		Label:     "Success Metrics",
		Required:  true,
		Question:  "What measurable success metrics define done (baseline, target, timeframe)?",
		ShortHint: "baseline, target, timeframe",
	},
	{
		Key:       "in_scope",
		Label:     "In Scope",
		Required:  true,
		Question:  "What is in scope for v1 (capabilities, platforms, user groups)?",
		ShortHint: "what is in scope for v1",
	},
	{
		Key:       "out_of_scope",
		Label:     "Out of Scope",
		Required:  true,
		Question:  "What is explicitly out of scope to prevent scope creep?",
		ShortHint: "explicit non-goals",
	},
	{
		Key:       "workflows",
		Label:     "Core Workflows",
		Required:  true,
		Question:  "Describe the core workflows (happy path steps).",
		ShortHint: "happy path steps",
	},
	{
		Key:       "data_entities",
		Label:     "Data Entities",
		Required:  true,
		Question:  "What are the key data entities and fields?",
		ShortHint: "entities and key fields",
	},
	{
		Key:       "system_of_record",
		Label:     "System of Record",
		Required:  true,
		Question:  "What is the system of record for each entity (or is this new)?",
		ShortHint: "system of record per entity",
	},
	{
		Key:       "constraints",
		Label:     "Constraints",
		Required:  true,
		Question:  "What constraints apply (deadline, budget, tech, policy, staffing)?",
		ShortHint: "deadlines, budget, policies",
	},
	{
		Key:       "edge_cases",
		Label:     "Edge Cases",
		Required:  true,
		Question:  "List the top 5 edge cases or exceptions that must be handled.",
		ShortHint: "top 5 exceptions",
	},
	{
		Key:       "non_functional",
		Label:     "Non-Functional Requirements",
		Required:  true,
		Question:  "What are the non-functional requirements (latency, availability, scale, reliability)?",
		ShortHint: "performance, availability, scale",
	},
	{
		Key:       "security_compliance",
		Label:     "Security and Compliance",
		Required:  true,
		Question:  "What security/compliance requirements apply (auth, PII, audit, retention)?",
		ShortHint: "auth, PII, audit, retention",
	},
	{
		Key:       "stakeholders",
		Label:     "Stakeholders",
		Required:  false,
		Question:  "Who are the stakeholders/approvers and what are their priorities?",
		ShortHint: "stakeholders and priorities",
	},
	{
		Key:       "integrations",
		Label:     "Integrations",
		Required:  false,
		Question:  "What integrations or external systems are required?",
		ShortHint: "external systems and APIs",
	},
	{
		Key:       "dependencies",
		Label:     "Dependencies",
		Required:  false,
		Question:  "What internal dependencies or teams are involved?",
		ShortHint: "teams, services, approvals",
	},
	{
		Key:       "timeline",
		Label:     "Timeline",
		Required:  false,
		Question:  "What timeline or milestones are expected?",
		ShortHint: "dates and milestones",
	},
	{
		Key:       "budget",
		Label:     "Budget",
		Required:  false,
		Question:  "What budget or staffing constraints exist?",
		ShortHint: "budget, staffing, procurement",
	},
	{
		Key:       "rollout",
		Label:     "Rollout Plan",
		Required:  false,
		Question:  "What rollout plan is desired (pilot, phased, full)?",
		ShortHint: "pilot/phased/full",
	},
	{
		Key:       "operations",
		Label:     "Operations and Support",
		Required:  false,
		Question:  "Who owns support and what are the SLA/support expectations?",
		ShortHint: "support ownership and SLA",
	},
	{
		Key:       "analytics",
		Label:     "Analytics and Reporting",
		Required:  false,
		Question:  "What analytics, dashboards, or reporting are required?",
		ShortHint: "dashboards and reporting",
	},
	{
		Key:       "requirements_must",
		Label:     "Must Requirements",
		Required:  false,
		Question:  "List must-have requirements.",
		ShortHint: "must-have requirements",
	},
	{
		Key:       "requirements_should",
		Label:     "Should Requirements",
		Required:  false,
		Question:  "List should-have requirements.",
		ShortHint: "should-have requirements",
	},
	{
		Key:       "requirements_could",
		Label:     "Could Requirements",
		Required:  false,
		Question:  "List could-have requirements.",
		ShortHint: "could-have requirements",
	},
	{
		Key:       "assumptions",
		Label:     "Assumptions",
		Required:  false,
		Question:  "Any assumptions to document?",
		ShortHint: "assumptions",
	},
}

func collectAnswers(in agent.Input, catalog []field) (map[string]string, map[string]bool) {
	answers := map[string]string{}
	derived := map[string]bool{}
	keys := catalogKeys(catalog)

	for key, value := range in.Parameters {
		normalized := strings.ToLower(strings.TrimSpace(key))
		if keys[normalized] {
			answers[normalized] = normalizeWhitespace(value)
		}
	}

	for key, value := range in.Artifacts {
		if strings.HasPrefix(key, "pipeline.") {
			continue
		}
		normalized := strings.ToLower(strings.TrimSpace(key))
		if keys[normalized] {
			if _, exists := answers[normalized]; !exists {
				answers[normalized] = normalizeWhitespace(value)
			}
		}
	}

	if raw, ok := in.Parameters["answers"]; ok {
		parseKeyValueLines(raw, answers, keys)
	}
	if raw, ok := in.Artifacts["answers"]; ok {
		parseKeyValueLines(raw, answers, keys)
	}

	return answers, derived
}

func catalogKeys(catalog []field) map[string]bool {
	keys := map[string]bool{}
	for _, entry := range catalog {
		keys[entry.Key] = true
	}
	return keys
}

func parseKeyValueLines(raw string, answers map[string]string, keys map[string]bool) {
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		if !keys[key] {
			continue
		}
		value := normalizeWhitespace(parts[1])
		if value == "" {
			continue
		}
		if _, exists := answers[key]; !exists {
			answers[key] = value
		}
	}
}

func missingFields(answers map[string]string, derived map[string]bool, catalog []field) ([]field, []field) {
	var missingRequired []field
	var missingOptional []field
	for _, entry := range catalog {
		if !isAnswered(answers[entry.Key]) || derived[entry.Key] {
			if entry.Required {
				missingRequired = append(missingRequired, entry)
			} else {
				missingOptional = append(missingOptional, entry)
			}
		}
	}
	return missingRequired, missingOptional
}

func isAnswered(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}
	lower := strings.ToLower(trimmed)
	switch lower {
	case "tbd", "unknown", "todo", "unspecified":
		return false
	default:
		return true
	}
}

func buildQuestions(task string, answers map[string]string, derived map[string]bool, missingRequired []field, missingOptional []field, depth string) string {
	var builder strings.Builder

	builder.WriteString("Priority Questions\n")
	builder.WriteString("------------------\n")
	priority := questionsFromFields(missingRequired)
	if derived["objective"] && answers["objective"] != "" {
		priority = append([]string{fmt.Sprintf("Confirm objective derived from task: %q", summarizeTask(task))}, priority...)
	}
	if len(priority) == 0 {
		priority = append(priority, "Confirm the spec aligns with stakeholder intent and there are no missing constraints.")
	}
	for i, q := range priority {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, q))
	}

	if strings.ToLower(depth) == "light" {
		return builder.String()
	}

	builder.WriteString("\nSecondary Questions\n")
	builder.WriteString("-------------------\n")
	secondary := questionsFromFields(missingOptional)
	secondary = append(secondary,
		"What are the top 3 risks you are worried about?",
		"Who must sign off before build starts?",
		"What does success look like 30/90 days after launch?",
		"What are the non-negotiable constraints the team must not violate?",
	)
	for i, q := range secondary {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, q))
	}

	return builder.String()
}

func questionsFromFields(fields []field) []string {
	out := make([]string, 0, len(fields))
	for _, entry := range fields {
		out = append(out, fmt.Sprintf("%s: %s", entry.Label, entry.Question))
	}
	return out
}

func buildOpenQuestions(missingRequired []field, missingOptional []field, depth string) string {
	var builder strings.Builder
	builder.WriteString("Open Questions\n")
	builder.WriteString("--------------\n")
	if len(missingRequired) == 0 && (strings.ToLower(depth) == "light" || len(missingOptional) == 0) {
		builder.WriteString("- None. Spec is ready for implementation.\n")
		return builder.String()
	}

	for _, entry := range missingRequired {
		builder.WriteString(fmt.Sprintf("- %s (%s)\n", entry.Label, entry.ShortHint))
	}
	if strings.ToLower(depth) != "light" {
		for _, entry := range missingOptional {
			builder.WriteString(fmt.Sprintf("- %s (%s)\n", entry.Label, entry.ShortHint))
		}
	}

	return builder.String()
}

func buildAcceptanceCriteria(answers map[string]string) string {
	criteria := []string{
		"All must-have requirements have testable acceptance criteria.",
		"Core workflow completes end-to-end without manual intervention.",
		"Role-based access is enforced for all user roles.",
		"Data validation prevents invalid states and logs failures.",
		"Edge cases are covered by tests or documented manual steps.",
		"Performance meets the stated non-functional requirements.",
		"Audit/logging meets compliance and support needs.",
	}

	nonFunctional := answers["non_functional"]
	if isAnswered(nonFunctional) {
		criteria[5] = fmt.Sprintf("Performance meets non-functional requirements: %s.", nonFunctional)
	}
	security := answers["security_compliance"]
	if isAnswered(security) {
		criteria[6] = fmt.Sprintf("Security/compliance requirements are satisfied: %s.", security)
	}

	return formatList(criteria)
}

func buildRisks(missingRequired []field) string {
	if len(missingRequired) == 0 {
		return "- None identified beyond normal delivery risk.\n"
	}
	var builder strings.Builder
	for _, entry := range missingRequired {
		builder.WriteString(fmt.Sprintf("- Missing %s may cause scope drift or rework; confirm before build.\n", entry.Label))
	}
	return builder.String()
}

func buildCompletionCriteria(missingRequired []field) string {
	if len(missingRequired) == 0 {
		return "- All required fields are provided and confirmed by stakeholders.\n"
	}
	var builder strings.Builder
	builder.WriteString("- Provide and confirm the missing required fields:\n")
	for _, entry := range missingRequired {
		builder.WriteString(fmt.Sprintf("  - %s\n", entry.Label))
	}
	return builder.String()
}

func buildDiagnostics(missingRequired []field, derived map[string]bool) []string {
	var diagnostics []string
	if len(missingRequired) > 0 {
		labels := make([]string, 0, len(missingRequired))
		for _, entry := range missingRequired {
			labels = append(labels, entry.Label)
		}
		diagnostics = append(diagnostics, fmt.Sprintf("incomplete spec; missing: %s", strings.Join(labels, ", ")))
	}
	if derived["objective"] {
		diagnostics = append(diagnostics, "objective derived from task; confirm with stakeholders")
	}
	if len(diagnostics) == 0 {
		diagnostics = append(diagnostics, "spec complete")
	}
	return diagnostics
}

func buildSpec(task string, answers map[string]string, derived map[string]bool, acceptanceCriteria string, risks string, completionCriteria string) string {
	status := "Draft - Needs Interview"
	if !derived["objective"] && isAnswered(answers["objective"]) {
		status = "Draft - Pending Confirmation"
	}
	if len(strings.TrimSpace(completionCriteria)) == 0 || strings.HasPrefix(completionCriteria, "- All required fields") {
		status = "Complete - Ready for Implementation"
	}

	var builder strings.Builder
	builder.WriteString("# Product Spec\n\n")
	builder.WriteString(fmt.Sprintf("Status: %s\n\n", status))
	builder.WriteString("## Initial Brief\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", normalizeWhitespace(task)))

	builder.WriteString("## Objective\n")
	builder.WriteString(formatValue("objective", answers, derived))

	builder.WriteString("## Problem Statement\n")
	builder.WriteString(formatValueWithFallback("problem_statement", answers, derived, "TBD - define the problem being solved"))

	builder.WriteString("## Users and Stakeholders\n")
	builder.WriteString(formatListSection(
		[]string{
			fmt.Sprintf("Primary users: %s", valueOrPlaceholder(answers["users"], "TBD - define user roles and permissions")),
			fmt.Sprintf("Stakeholders: %s", valueOrPlaceholder(answers["stakeholders"], "TBD - list stakeholders/approvers")),
		},
	))

	builder.WriteString("## Scope\n")
	builder.WriteString(formatListSection(
		[]string{
			fmt.Sprintf("In scope: %s", valueOrPlaceholder(answers["in_scope"], "TBD - in-scope capabilities")),
			fmt.Sprintf("Out of scope: %s", valueOrPlaceholder(answers["out_of_scope"], "TBD - explicit non-goals")),
		},
	))

	builder.WriteString("## Core Workflows\n")
	builder.WriteString(formatValueWithFallback("workflows", answers, derived, "TBD - describe happy path steps"))

	builder.WriteString("## Data and Systems\n")
	builder.WriteString(formatListSection(
		[]string{
			fmt.Sprintf("Data entities: %s", valueOrPlaceholder(answers["data_entities"], "TBD - entities and key fields")),
			fmt.Sprintf("System of record: %s", valueOrPlaceholder(answers["system_of_record"], "TBD - source of truth")),
			fmt.Sprintf("Integrations: %s", valueOrPlaceholder(answers["integrations"], "TBD - external systems/APIs")),
		},
	))

	builder.WriteString("## Requirements\n")
	builder.WriteString("Must:\n")
	builder.WriteString(formatListOrPlaceholder(answers["requirements_must"], "TBD - must-have requirements"))
	builder.WriteString("\nShould:\n")
	builder.WriteString(formatListOrPlaceholder(answers["requirements_should"], "TBD - should-have requirements"))
	builder.WriteString("\nCould:\n")
	builder.WriteString(formatListOrPlaceholder(answers["requirements_could"], "TBD - could-have requirements"))
	builder.WriteString("\n")

	builder.WriteString("## Non-Functional Requirements\n")
	builder.WriteString(formatValueWithFallback("non_functional", answers, derived, "TBD - performance, availability, scale, reliability"))

	builder.WriteString("## Security and Compliance\n")
	builder.WriteString(formatValueWithFallback("security_compliance", answers, derived, "TBD - auth, PII, audit, retention"))

	builder.WriteString("## Constraints\n")
	builder.WriteString(formatValueWithFallback("constraints", answers, derived, "TBD - deadlines, budget, tech, policies"))

	builder.WriteString("## Budget and Resourcing\n")
	builder.WriteString(formatValueWithFallback("budget", answers, derived, "TBD - budget, staffing, procurement"))

	builder.WriteString("## Edge Cases and Exceptions\n")
	builder.WriteString(formatValueWithFallback("edge_cases", answers, derived, "TBD - top 5 exceptions"))

	builder.WriteString("## Success Metrics\n")
	builder.WriteString(formatValueWithFallback("success_metrics", answers, derived, "TBD - baseline, target, timeframe"))

	builder.WriteString("## Dependencies\n")
	builder.WriteString(formatValueWithFallback("dependencies", answers, derived, "TBD - teams/services/approvals"))

	builder.WriteString("## Timeline and Rollout\n")
	builder.WriteString(formatListSection(
		[]string{
			fmt.Sprintf("Timeline: %s", valueOrPlaceholder(answers["timeline"], "TBD - milestones")),
			fmt.Sprintf("Rollout: %s", valueOrPlaceholder(answers["rollout"], "TBD - pilot/phased/full")),
			fmt.Sprintf("Operations: %s", valueOrPlaceholder(answers["operations"], "TBD - support ownership/SLA")),
		},
	))

	builder.WriteString("## Analytics and Reporting\n")
	builder.WriteString(formatValueWithFallback("analytics", answers, derived, "TBD - dashboards, reports, instrumentation"))

	builder.WriteString("## Assumptions\n")
	builder.WriteString(formatValueWithFallback("assumptions", answers, derived, "TBD - assumptions pending confirmation"))

	builder.WriteString("## Acceptance Criteria\n")
	builder.WriteString(acceptanceCriteria)

	builder.WriteString("## Risks\n")
	builder.WriteString(risks)

	builder.WriteString("## Completion Criteria\n")
	builder.WriteString(completionCriteria)

	builder.WriteString("\n## Verification Plan\n")
	builder.WriteString(formatListSection([]string{
		"Stakeholder sign-off on scope and acceptance criteria.",
		"Demo of core workflows with representative data.",
		"Metrics instrumentation validated against success criteria.",
		"Edge cases verified via tests or documented manual checks.",
	}))

	return builder.String()
}

func formatValue(key string, answers map[string]string, derived map[string]bool) string {
	value := answers[key]
	if isAnswered(value) && !derived[key] {
		return fmt.Sprintf("%s\n\n", value)
	}
	if derived[key] && isAnswered(value) {
		return fmt.Sprintf("%s (derived from task, confirm)\n\n", value)
	}
	return "TBD - provide objective\n\n"
}

func formatValueWithFallback(key string, answers map[string]string, derived map[string]bool, fallback string) string {
	value := answers[key]
	if isAnswered(value) && !derived[key] {
		return fmt.Sprintf("%s\n\n", value)
	}
	if derived[key] && isAnswered(value) {
		return fmt.Sprintf("%s (derived from task, confirm)\n\n", value)
	}
	return fmt.Sprintf("%s\n\n", fallback)
}

func formatListSection(items []string) string {
	var builder strings.Builder
	for _, item := range items {
		builder.WriteString(fmt.Sprintf("- %s\n", item))
	}
	builder.WriteString("\n")
	return builder.String()
}

func valueOrPlaceholder(value string, placeholder string) string {
	if isAnswered(value) {
		return value
	}
	return placeholder
}

func formatList(items []string) string {
	var builder strings.Builder
	for _, item := range items {
		builder.WriteString(fmt.Sprintf("- %s\n", item))
	}
	return builder.String()
}

func formatListOrPlaceholder(value string, placeholder string) string {
	if isAnswered(value) {
		trimmed := strings.TrimSpace(value)
		if strings.HasPrefix(trimmed, "-") || strings.Contains(trimmed, "\n") {
			return fmt.Sprintf("%s\n", trimmed)
		}
		return fmt.Sprintf("- %s\n", trimmed)
	}
	return fmt.Sprintf("- %s\n", placeholder)
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func summarizeTask(task string) string {
	task = normalizeWhitespace(task)
	if len(task) <= 120 {
		return task
	}
	return task[:117] + "..."
}
