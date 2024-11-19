package multi_agent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

func TestConcurrentAgent_Plan(t *testing.T) {
	agent := NewConcurrentAgent()

	// Mock inputs and expected outputs
	inputs := map[string]string{"input": "test"}
	expectedActions := []schema.AgentAction{}
	expectedFinish := &schema.AgentFinish{ReturnValues: map[string]any{"output": "test"}}

	actions, finish, err := agent.Plan(context.Background(), nil, inputs)
	assert.NoError(t, err)
	assert.Equal(t, expectedActions, actions)
	assert.Equal(t, expectedFinish, finish)
}

func TestConcurrentAgent_GetInputKeys(t *testing.T) {
	agent := NewConcurrentAgent()
	expectedKeys := []string{}
	assert.Equal(t, expectedKeys, agent.GetInputKeys())
}

func TestConcurrentAgent_GetOutputKeys(t *testing.T) {
	agent := NewConcurrentAgent()
	expectedKeys := []string{}
	assert.Equal(t, expectedKeys, agent.GetOutputKeys())
}

func TestConcurrentAgent_GetTools(t *testing.T) {
	agent := NewConcurrentAgent()
	expectedTools := []tools.Tool{}
	assert.Equal(t, expectedTools, agent.GetTools())
}

func TestConcurrentAgent_InitializeConcurrentActions(t *testing.T) {
	agent := NewConcurrentAgent()
	actions := []schema.AgentAction{
		{Tool: "testTool", ToolInput: "testInput"},
	}
	agent.InitializeConcurrentActions(actions)
	assert.Equal(t, 1, len(agent.Graph.Nodes))
	assert.Equal(t, "initialized", agent.Graph.Nodes[0].State)
}

func TestConcurrentAgent_ExecuteConcurrentActions(t *testing.T) {
	agent := NewConcurrentAgent()
	actions := []schema.AgentAction{
		{Tool: "testTool", ToolInput: "testInput"},
	}
	agent.InitializeConcurrentActions(actions)
	agent.ExecuteConcurrentActions()
	assert.Equal(t, "completed", agent.Graph.Nodes[0].State)
}
