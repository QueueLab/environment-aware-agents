package graph_orchestration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

func TestGraphAgent_Plan(t *testing.T) {
	llm := &llms.MockModel{}
	tools := []tools.Tool{&tools.MockTool{}}
	agent := NewGraphAgent(llm, tools)

	ctx := context.Background()
	intermediateSteps := []schema.AgentStep{}
	inputs := map[string]string{"input": "test"}

	actions, finish, err := agent.Plan(ctx, intermediateSteps, inputs)
	assert.NoError(t, err)
	assert.Nil(t, actions)
	assert.Nil(t, finish)
}

func TestGraphAgent_GetInputKeys(t *testing.T) {
	llm := &llms.MockModel{}
	tools := []tools.Tool{&tools.MockTool{}}
	agent := NewGraphAgent(llm, tools)

	inputKeys := agent.GetInputKeys()
	assert.Nil(t, inputKeys)
}

func TestGraphAgent_GetOutputKeys(t *testing.T) {
	llm := &llms.MockModel{}
	tools := []tools.Tool{&tools.MockTool{}}
	agent := NewGraphAgent(llm, tools)

	outputKeys := agent.GetOutputKeys()
	assert.Nil(t, outputKeys)
}

func TestGraphAgent_GetTools(t *testing.T) {
	llm := &llms.MockModel{}
	tools := []tools.Tool{&tools.MockTool{}}
	agent := NewGraphAgent(llm, tools)

	agentTools := agent.GetTools()
	assert.Equal(t, tools, agentTools)
}
