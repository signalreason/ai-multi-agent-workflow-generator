package all

import (
	"ai-multi-agent-workflow-generator/internal/agents/analyticsintegrator"
	"ai-multi-agent-workflow-generator/internal/agents/apidesigner"
	"ai-multi-agent-workflow-generator/internal/agents/backupmanager"
	"ai-multi-agent-workflow-generator/internal/agents/benchmarktester"
	"ai-multi-agent-workflow-generator/internal/agents/buildcop"
	"ai-multi-agent-workflow-generator/internal/agents/chaosengineer"
	"ai-multi-agent-workflow-generator/internal/agents/codegenerator"
	"ai-multi-agent-workflow-generator/internal/agents/complianceofficer"
	"ai-multi-agent-workflow-generator/internal/agents/configmanager"
	"ai-multi-agent-workflow-generator/internal/agents/costoptimizer"
	"ai-multi-agent-workflow-generator/internal/agents/datamodeler"
	"ai-multi-agent-workflow-generator/internal/agents/dependencyauditor"
	"ai-multi-agent-workflow-generator/internal/agents/deploymentagent"
	"ai-multi-agent-workflow-generator/internal/agents/documentationwriter"
	"ai-multi-agent-workflow-generator/internal/agents/featurearchitect"
	"ai-multi-agent-workflow-generator/internal/agents/featureflagger"
	"ai-multi-agent-workflow-generator/internal/agents/ideadistiller"
	"ai-multi-agent-workflow-generator/internal/agents/incidentresponder"
	"ai-multi-agent-workflow-generator/internal/agents/infrastructurebuilder"
	"ai-multi-agent-workflow-generator/internal/agents/licensechecker"
	"ai-multi-agent-workflow-generator/internal/agents/localizationmanager"
	"ai-multi-agent-workflow-generator/internal/agents/loganalyzer"
	"ai-multi-agent-workflow-generator/internal/agents/migrationplanner"
	"ai-multi-agent-workflow-generator/internal/agents/performanceoptimizer"
	"ai-multi-agent-workflow-generator/internal/agents/qatriage"
	"ai-multi-agent-workflow-generator/internal/agents/refactorer"
	"ai-multi-agent-workflow-generator/internal/agents/releasenoteswriter"
	"ai-multi-agent-workflow-generator/internal/agents/releasepackager"
	"ai-multi-agent-workflow-generator/internal/agents/securityauditor"
	"ai-multi-agent-workflow-generator/internal/agents/tddenforcer"
	"ai-multi-agent-workflow-generator/internal/agents/testdatagenerator"
	"ai-multi-agent-workflow-generator/internal/agents/uxresearcher"

	"ai-multi-agent-workflow-generator/internal/agent"
)

// RegisterAll registers all available agents.
func RegisterAll(r *agent.Registry) {
	ideadistiller.Register(r)
	featurearchitect.Register(r)
	codegenerator.Register(r)
	refactorer.Register(r)
	tddenforcer.Register(r)
	documentationwriter.Register(r)
	infrastructurebuilder.Register(r)
	securityauditor.Register(r)
	benchmarktester.Register(r)
	performanceoptimizer.Register(r)
	releasepackager.Register(r)
	deploymentagent.Register(r)
	uxresearcher.Register(r)
	dependencyauditor.Register(r)
	apidesigner.Register(r)
	datamodeler.Register(r)
	configmanager.Register(r)
	loganalyzer.Register(r)
	testdatagenerator.Register(r)
	migrationplanner.Register(r)
	incidentresponder.Register(r)
	qatriage.Register(r)
	buildcop.Register(r)
	licensechecker.Register(r)
	complianceofficer.Register(r)
	localizationmanager.Register(r)
	analyticsintegrator.Register(r)
	chaosengineer.Register(r)
	costoptimizer.Register(r)
	backupmanager.Register(r)
	featureflagger.Register(r)
	releasenoteswriter.Register(r)
}
