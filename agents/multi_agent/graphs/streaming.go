package multi_agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

// StreamingAgent is an implementation of the Agent interface with streaming features.
type StreamingAgent struct {
	Graph *Graph
	Tools []tools.Tool
}

var _ agents.Agent = (*StreamingAgent)(nil)

// NewStreamingAgent creates a new StreamingAgent.
func NewStreamingAgent() *StreamingAgent {
	return &StreamingAgent{
		Graph: NewGraph(),
		Tools: []tools.Tool{},
	}
}

// Plan decides what action to take or returns the final result of the input.
func (a *StreamingAgent) Plan(
	ctx context.Context,
	intermediateSteps []schema.AgentStep,
	inputs map[string]string,
) ([]schema.AgentAction, *schema.AgentFinish, error) {
	// Implement the logic to decide the next action or return the final result.
	// This is a placeholder implementation.
	if len(inputs) == 0 {
		return nil, &schema.AgentFinish{ReturnValues: map[string]any{"output": "no input"}}, nil
	}
	return nil, &schema.AgentFinish{ReturnValues: map[string]any{"output": inputs["input"]}}, nil
}

// GetInputKeys returns the input keys for the agent.
func (a *StreamingAgent) GetInputKeys() []string {
	// Implement the logic to return the input keys.
	// This is a placeholder implementation.
	return []string{"input"}
}

// GetOutputKeys returns the output keys for the agent.
func (a *StreamingAgent) GetOutputKeys() []string {
	// Implement the logic to return the output keys.
	// This is a placeholder implementation.
	return []string{"output"}
}

// GetTools returns the tools available to the agent.
func (a *StreamingAgent) GetTools() []tools.Tool {
	return a.Tools
}

// InitializeStreamingActions initializes the streaming actions for the agent.
func (a *StreamingAgent) InitializeStreamingActions(actions []schema.AgentAction) {
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

// ExecuteStreamingActions executes the streaming actions for the agent.
func (a *StreamingAgent) ExecuteStreamingActions() {
	var wg sync.WaitGroup
	for _, node := range a.Graph.Nodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			for _, action := range n.Actions {
				// Implement the logic for self-reflecting or human-in-the-loop tasks.
				// This is a placeholder implementation.
				n.mu.Lock()
				n.State = "executing"
				n.mu.Unlock()
				// Simulate action execution
				n.mu.Lock()
				n.State = "completed"
				n.mu.Unlock()
			}
		}(node)
	}
	wg.Wait()
	a.Graph.Execute()
}

// StreamOutput streams the output of the agent's actions.
func (a *StreamingAgent) StreamOutput(ctx context.Context, outputChan chan<- string) error {
	for _, node := range a.Graph.Nodes {
		for _, action := range node.Actions {
			// Implement the logic to stream the output of the action.
			// This is a placeholder implementation.
			outputChan <- fmt.Sprintf("Streaming output for action: %v", action)
		}
	}
	close(outputChan)
	return nil
}
