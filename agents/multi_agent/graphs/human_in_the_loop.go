package multi_agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

// HumanInTheLoopAgent is an implementation of the Agent interface with human-in-the-loop features.
type HumanInTheLoopAgent struct {
	Graph *Graph
	Tools []tools.Tool
}

var _ agents.Agent = (*HumanInTheLoopAgent)(nil)

// NewHumanInTheLoopAgent creates a new HumanInTheLoopAgent.
func NewHumanInTheLoopAgent() *HumanInTheLoopAgent {
	return &HumanInTheLoopAgent{
		Graph: NewGraph(),
		Tools: []tools.Tool{},
	}
}

// Plan decides what action to take or returns the final result of the input.
func (a *HumanInTheLoopAgent) Plan(
	ctx context.Context,
	intermediateSteps []schema.AgentStep,
	inputs map[string]string,
) ([]schema.AgentAction, *schema.AgentFinish, error) {
	if len(intermediateSteps) == 0 {
		// Initial step - create actions based on input
		actions := []schema.AgentAction{
			{Tool: "analyze", ToolInput: inputs["input"]},
		}
		return actions, nil, nil
	}

	// Analyze previous steps and decide next action
	lastStep := intermediateSteps[len(intermediateSteps)-1]
	if lastStep.Observation != "" {
		// If we have an observation, process it and potentially finish
		if len(intermediateSteps) >= 3 { // Limit steps to prevent infinite loops
			return nil, &schema.AgentFinish{
				ReturnValues: map[string]any{
					"output": fmt.Sprintf("Final result after analysis: %s", lastStep.Observation),
				},
			}, nil
		}
		// Continue with next action
		return []schema.AgentAction{
			{Tool: "process", ToolInput: lastStep.Observation},
		}, nil, nil
	}

	return nil, &schema.AgentFinish{
		ReturnValues: map[string]any{
			"output": "Unable to determine next action",
		},
	}, nil
}

// GetInputKeys returns the input keys for the agent.
func (a *HumanInTheLoopAgent) GetInputKeys() []string {
	return []string{"input", "context", "parameters"}
}

// GetOutputKeys returns the output keys for the agent.
func (a *HumanInTheLoopAgent) GetOutputKeys() []string {
	return []string{"output", "reasoning", "actions"}
}

// GetTools returns the tools available to the agent.
func (a *HumanInTheLoopAgent) GetTools() []tools.Tool {
	return a.Tools
}

// InitializeHumanInTheLoopActions initializes the human-in-the-loop actions for the agent.
func (a *HumanInTheLoopAgent) InitializeHumanInTheLoopActions(actions []schema.AgentAction) {
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

// ExecuteHumanInTheLoopActions executes the human-in-the-loop actions for the agent.
func (a *HumanInTheLoopAgent) ExecuteHumanInTheLoopActions() {
	var wg sync.WaitGroup
	for _, node := range a.Graph.Nodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			for _, action := range n.Actions {
				n.mu.Lock()
				n.State = "executing"
				n.mu.Unlock()

				if agentAction, ok := action.(schema.AgentAction); ok {
					n.mu.Lock()
					// Execute the tool
					for _, tool := range a.Tools {
						if tool.Name() == agentAction.Tool {
							result, err := tool.Call(agentAction.ToolInput)
							if err == nil {
								n.Result = result
							} else {
								n.Result = fmt.Sprintf("Error: %v", err)
							}
							break
						}
					}
					n.State = "completed"
					n.mu.Unlock()
				}
			}
		}(node)
	}
	wg.Wait()
	a.Graph.Execute()
}

// HumanFeedback collects feedback from a human user.
func (a *HumanInTheLoopAgent) HumanFeedback(ctx context.Context, feedback string) error {
	// Process the feedback and update the graph
	feedbackNode := &Node{
		ID:    len(a.Graph.Nodes),
		Value: feedback,
		State: "feedback",
	}
	a.Graph.AddNode(feedbackNode)

	// Update the state of related nodes based on feedback
	for _, node := range a.Graph.Nodes {
		if node.State == "completed" {
			node.mu.Lock()
			node.Feedback = feedback
			// Potentially modify the node's state or trigger new actions based on feedback
			if feedback == "reject" {
				node.State = "rejected"
			} else if feedback == "approve" {
				node.State = "approved"
			}
			node.mu.Unlock()
		}
	}

	return nil
}
