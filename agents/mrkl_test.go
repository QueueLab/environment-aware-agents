package agents

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

func TestMRKLPlanEncapsulation(t *testing.T) {
	t.Parallel()
	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	llm, err := openai.New(openai.WithModel("gpt-4"))
	require.NoError(t, err)

	agent := NewOneShotAgent(llm, []tools.Tool{tools.Calculator{}})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	intermediateSteps := []schema.AgentStep{}
	inputs := map[string]string{"input": "What is 2 + 2?"}

	actions, finish, err := agent.Plan(ctx, intermediateSteps, inputs)
	require.NoError(t, err)
	require.Nil(t, finish)
	require.NotEmpty(t, actions)
}

func TestMRKLDynamicControlFlow(t *testing.T) {
	t.Parallel()
	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	llm, err := openai.New(openai.WithModel("gpt-4"))
	require.NoError(t, err)

	agent := NewOneShotAgent(llm, []tools.Tool{tools.Calculator{}})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	intermediateSteps := []schema.AgentStep{}
	inputs := map[string]string{"input": "Calculate the square root of 144"}

	actions, finish, err := agent.Plan(ctx, intermediateSteps, inputs)
	require.NoError(t, err)
	require.Nil(t, finish)
	require.NotEmpty(t, actions)
}

func TestMRKLErrorHandling(t *testing.T) {
	t.Parallel()
	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	llm, err := openai.New(openai.WithModel("gpt-4"))
	require.NoError(t, err)

	agent := NewOneShotAgent(llm, []tools.Tool{tools.Calculator{}})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	intermediateSteps := []schema.AgentStep{}
	inputs := map[string]string{"input": "What is 2 + 2?"}

	_, _, err = agent.Plan(ctx, intermediateSteps, inputs)
	require.NoError(t, err)

	// Simulate an error
	inputs["input"] = "invalid input"
	_, _, err = agent.Plan(ctx, intermediateSteps, inputs)
	require.Error(t, err)
}
