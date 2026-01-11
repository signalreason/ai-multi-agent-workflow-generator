package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"ai-multi-agent-workflow-generator/internal/agent"
	"ai-multi-agent-workflow-generator/internal/agents/all"
	"ai-multi-agent-workflow-generator/internal/manager"
)

func main() {
	taskPath := flag.String("task", "", "path to YAML task file")
	listAgents := flag.Bool("list-agents", false, "list available agents")
	flag.Parse()

	registry := agent.NewRegistry()
	all.RegisterAll(registry)

	if *listAgents {
		for _, name := range registry.Names() {
			fmt.Println(name)
		}
		return
	}

	if *taskPath == "" {
		fmt.Fprintln(os.Stderr, "-task is required")
		os.Exit(2)
	}

	mgr := manager.New(registry)
	task, err := mgr.LoadTaskFile(*taskPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	results, err := mgr.Run(context.Background(), task)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	payload, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(payload))
}
