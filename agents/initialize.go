package agents

import (
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

const _defaultMaxIterations = 5

// AgentType is a string type representing the type of agent to create.
type AgentType string

const (
	// ZeroShotReactDescription is an AgentType constant that represents
	// the "zeroShotReactDescription" agent type.
	ZeroShotReactDescription AgentType = "zeroShotReactDescription"
	// ConversationalReactDescription is an AgentType constant that represents
	// the "conversationalReactDescription" agent type.
	ConversationalReactDescription AgentType = "conversationalReactDescription"
	// DynamicControlFlowAgent is an AgentType constant that represents
	// the "dynamicControlFlowAgent" agent type.
	DynamicControlFlowAgent AgentType = "dynamicControlFlowAgent"
)

// Deprecated: This may be removed in the future; please use NewExecutor instead.
// Initialize is a function that creates a new executor with the specified LLM
// model, tools, agent type, and options. It returns an Executor or an error
// if there is any issues during the creation process.
func Initialize(
	llm llms.Model,
	tools []tools.Tool,
	agentType AgentType,
	opts ...Option,
) (*Executor, error) {
	var agent Agent
	switch agentType {
	case ZeroShotReactDescription:
		agent = NewOneShotAgent(llm, tools, opts...)
	case ConversationalReactDescription:
		agent = NewConversationalAgent(llm, tools, opts...)
	case DynamicControlFlowAgent:
		agent = NewDynamicControlFlowAgent(llm, tools, opts...)
	default:
		return &Executor{}, ErrUnknownAgentType
	}
	return NewExecutor(agent, opts...), nil
}

// NewDynamicControlFlowAgent creates a new agent with dynamic control flow capabilities.
func NewDynamicControlFlowAgent(llm llms.Model, tools []tools.Tool, opts ...Option) *OneShotZeroAgent {
	options := mrklDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &OneShotZeroAgent{
		Chain: chains.NewLLMChain(
			llm,
			options.getMrklPrompt(tools),
			chains.WithCallback(options.callbacksHandler),
		),
		Tools:            tools,
		OutputKey:        options.outputKey,
		CallbacksHandler: options.callbacksHandler,
	}
}
