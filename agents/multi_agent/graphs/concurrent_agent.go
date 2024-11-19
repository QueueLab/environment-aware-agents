package multi_agent

import (
	"context"
	"sync"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

// ConcurrentAgent is an implementation of the Agent interface with concurrency features.
type ConcurrentAgent struct {
	Graph *Graph
	Tools []tools.Tool
}

var _ agents.Agent = (*ConcurrentAgent)(nil)
var _ agents.ConcurrentAgent = (*ConcurrentAgent)(nil)

// NewConcurrentAgent creates a new ConcurrentAgent.
func NewConcurrentAgent() *ConcurrentAgent {
	return &ConcurrentAgent{
		Graph: NewGraph(),
		Tools: []tools.Tool{},
	}
}

// Plan decides what action to take or returns the final result of the input.
func (a *ConcurrentAgent) Plan(
	ctx context.Context,
	intermediateSteps []schema.AgentStep,
	inputs map[string]string,
) ([]schema.AgentAction, *schema.AgentFinish, error) {
	if len(intermediateSteps) == 0 {
		// Initial planning phase
		actions := make([]schema.AgentAction, 0)
		for toolName, input := range inputs {
			action := schema.AgentAction{
				Tool:  toolName,
				Input: input,
				Log:   "Initial action from input",
			}
			actions = append(actions, action)
		}
		return actions, nil, nil
	}

	// Analyze intermediate steps and decide next actions
	allCompleted := true
	for _, step := range intermediateSteps {
		if step.Action != nil && !step.Observation.IsEmpty() {
			continue
		}
		allCompleted = false
		break
	}

	if allCompleted {
		// All steps completed, return final result
		result := make(map[string]any)
		for _, step := range intermediateSteps {
			result[step.Action.Tool] = step.Observation.String()
		}
		return nil, &schema.AgentFinish{ReturnValues: result}, nil
	}

	// Continue with more actions if needed
	return []schema.AgentAction{}, nil, nil
}

// GetInputKeys returns the input keys for the agent.
func (a *ConcurrentAgent) GetInputKeys() []string {
	keys := make([]string, 0)
	for _, tool := range a.Tools {
		keys = append(keys, tool.Name())
	}
	return keys
}

// GetOutputKeys returns the output keys for the agent.
func (a *ConcurrentAgent) GetOutputKeys() []string {
	keys := make([]string, 0)
	for _, tool := range a.Tools {
		keys = append(keys, tool.Name()+"_result")
	}
	return keys
}

// GetTools returns the tools available to the agent.
func (a *ConcurrentAgent) GetTools() []tools.Tool {
	return a.Tools
}

// InitializeConcurrentActions initializes the concurrent actions for the agent.
func (a *ConcurrentAgent) InitializeConcurrentActions(actions []schema.AgentAction) {
	for _, action := range actions {
		node := &Node{
			ID:      len(a.Graph.Nodes),
			Value:   action,
			State:   "initialized",
			Actions: []interface{}{action},
		}
		a.Graph.AddNode(node)
	}
}

// ExecuteConcurrentActions executes the concurrent actions for the agent.
func (a *ConcurrentAgent) ExecuteConcurrentActions() {
	var wg sync.WaitGroup
	for _, node := range a.Graph.Nodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			for _, action := range n.Actions {
				agentAction, ok := action.(schema.AgentAction)
				if !ok {
					n.mu.Lock()
					n.State = "error"
					n.mu.Unlock()
					continue
				}

				n.mu.Lock()
				n.State = "executing"
				n.mu.Unlock()

				// Execute the tool
				for _, tool := range a.Tools {
					if tool.Name() == agentAction.Tool {
						result, err := tool.Call(agentAction.Input)
						n.mu.Lock()
						if err != nil {
							n.State = "error"
							n.Result = err.Error()
						} else {
							n.State = "completed"
							n.Result = result
						}
						n.mu.Unlock()
						break
					}
				}
			}
		}(node)
	}
	wg.Wait()
	a.Graph.Execute()
}
