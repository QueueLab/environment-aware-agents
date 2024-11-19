package multi_agent

import (
	"context"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

// GraphAgent is an implementation of the Agent interface with graph features.
type GraphAgent struct {
	Graph *Graph
	Tools []tools.Tool
}

var _ agents.Agent = (*GraphAgent)(nil)

// NewGraphAgent creates a new GraphAgent.
func NewGraphAgent() *GraphAgent {
	return &GraphAgent{
		Graph: NewGraph(),
		Tools: []tools.Tool{},
	}
}

// Plan decides what action to take or returns the final result of the input.
func (a *GraphAgent) Plan(
	ctx context.Context,
	intermediateSteps []schema.AgentStep,
	inputs map[string]string,
) ([]schema.AgentAction, *schema.AgentFinish, error) {
	// Analyze the current state of the graph
	currentState := a.Graph.GetCurrentState()

	// Check if we've reached a terminal state
	if a.Graph.IsTerminalState(currentState) {
		return nil, &schema.AgentFinish{
			ReturnValues: map[string]string{"result": currentState.String()},
			Log:          "Reached terminal state",
		}, nil
	}

	// Determine next actions based on graph traversal
	actions := make([]schema.AgentAction, 0)
	nextStates := a.Graph.GetNextStates(currentState)

	for _, state := range nextStates {
		action := schema.AgentAction{
			Tool: "graph_traversal",
			ToolInput: map[string]string{
				"target_state": state.String(),
			},
			Log: "Transitioning to next state",
		}
		actions = append(actions, action)
	}

	return actions, nil, nil
}

// GetInputKeys returns the input keys for the agent.
func (a *GraphAgent) GetInputKeys() []string {
	return []string{"initial_state", "goal_state", "constraints"}
}

// GetOutputKeys returns the output keys for the agent.
func (a *GraphAgent) GetOutputKeys() []string {
	return []string{"final_state", "path", "metrics"}
}

// GetTools returns the tools available to the agent.
func (a *GraphAgent) GetTools() []tools.Tool {
	return a.Tools
}
