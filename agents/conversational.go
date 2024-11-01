package agents

import (
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

const (
	_conversationalFinalAnswerAction = "AI:"
)

// ConversationalAgent is a struct that represents an agent responsible for deciding
// what to do or give the final output if the task is finished given a set of inputs
// and previous steps taken.
//
// Other agents are often optimized for using tools to figure out the best response,
// which is not ideal in a conversational setting where you may want the agent to be
// able to chat with the user as well.
type ConversationalAgent struct {
	// Chain is the chain used to call with the values. The chain should have an
	// input called "agent_scratchpad" for the agent to put its thoughts in.
	Chain chains.Chain
	// Tools is a list of the tools the agent can use.
	Tools []tools.Tool
	// Output key is the key where the final output is placed.
	OutputKey string
	// CallbacksHandler is the handler for callbacks.
	CallbacksHandler callbacks.Handler
}

var _ Agent = (*ConversationalAgent)(nil)

func NewConversationalAgent(llm llms.Model, tools []tools.Tool, opts ...Option) *ConversationalAgent {
	options := conversationalDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &ConversationalAgent{
		Chain: chains.NewLLMChain(
			llm,
			options.getConversationalPrompt(tools),
			chains.WithCallback(options.callbacksHandler),
		),
		Tools:            tools,
		OutputKey:        options.outputKey,
		CallbacksHandler: options.callbacksHandler,
	}
}

// Plan decides what action to take or returns the final result of the input.
func (a *ConversationalAgent) Plan(
	ctx context.Context,
	intermediateSteps []schema.AgentStep,
	inputs map[string]string,
) ([]schema.AgentAction, *schema.AgentFinish, error) {
	fullInputs := make(map[string]any, len(inputs))
	for key, value := range inputs {
		fullInputs[key] = value
	}

	fullInputs["agent_scratchpad"] = constructScratchPad(intermediateSteps)

	var stream func(ctx context.Context, chunk []byte) error

	if a.CallbacksHandler != nil {
		stream = func(ctx context.Context, chunk []byte) error {
			a.CallbacksHandler.HandleStreamingFunc(ctx, chunk)
			return nil
		}
	}

	// Use context.WithTimeout to manage task execution time
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	output, err := chains.Predict(
		ctx,
		a.Chain,
		fullInputs,
		chains.WithStopWords([]string{"\nObservation:", "\n\tObservation:"}),
		chains.WithStreamingFunc(stream),
	)
	if err != nil {
		return nil, nil, err
	}

 if action, ok := handleDynamicControlFlow(output); ok {
     return []schema.AgentAction{action}, nil, nil
 }

 // ... (elsewhere in the file)
 func handleDynamicControlFlow(output string) (schema.AgentAction, bool) {
     if strings.Contains(output, "dynamic_control_flow") {
         nextTool := extractNextTool(output)
         if nextTool != "" {
             return schema.AgentAction{
                 Tool:      nextTool,
                 ToolInput: "",
                 Log:       output,
             }, true
         }
     }
     return schema.AgentAction{}, false
 }

	return a.parseOutput(output)
}

func (a *ConversationalAgent) GetInputKeys() []string {
	chainInputs := a.Chain.GetInputKeys()

	// Remove inputs given in plan.
	agentInput := make([]string, 0, len(chainInputs))
	for _, v := range chainInputs {
		if v == "agent_scratchpad" {
			continue
		}
		agentInput = append(agentInput, v)
	}

	return agentInput
}

func (a *ConversationalAgent) GetOutputKeys() []string {
	return []string{a.OutputKey}
}

func (a *ConversationalAgent) GetTools() []tools.Tool {
	return a.Tools
}

func constructScratchPad(steps []schema.AgentStep) string {
	var scratchPad string
	if len(steps) > 0 {
		for _, step := range steps {
			scratchPad += step.Action.Log
			scratchPad += "\nObservation: " + step.Observation
		}
		scratchPad += "\n" + "Thought:"
	}

	return scratchPad
}

func (a *ConversationalAgent) parseOutput(output string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	if strings.Contains(output, _conversationalFinalAnswerAction) {
		splits := strings.Split(output, _conversationalFinalAnswerAction)

		finishAction := &schema.AgentFinish{
			ReturnValues: map[string]any{
				a.OutputKey: splits[len(splits)-1],
			},
			Log: output,
		}

		return nil, finishAction, nil
	}

	r := regexp.MustCompile(`Action: (.*?)[\n]*Action Input: (.*)`)
	matches := r.FindStringSubmatch(output)
	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("%w: %s", ErrUnableToParseOutput, output)
	}

	// Handle dynamic control flow
	if strings.Contains(output, "dynamic_control_flow") {
		// Extract the next tool to call from the output
		nextTool := extractNextTool(output)
		if nextTool != "" {
			action := schema.AgentAction{
				Tool:      nextTool,
				ToolInput: "", // You may need to set the appropriate input for the next tool
				Log:       output,
			}
			return []schema.AgentAction{action}, nil, nil
		}
	}

	return []schema.AgentAction{
		{Tool: strings.TrimSpace(matches[1]), ToolInput: strings.TrimSpace(matches[2]), Log: output},
	}, nil, nil
}

//go:embed prompts/conversational_prefix.txt
var _defaultConversationalPrefix string //nolint:gochecknoglobals

//go:embed prompts/conversational_format_instructions.txt
var _defaultConversationalFormatInstructions string //nolint:gochecknoglobals

//go:embed prompts/conversational_suffix.txt
var _defaultConversationalSuffix string //nolint:gochecknoglobals

func createConversationalPrompt(tools []tools.Tool, prefix, instructions, suffix string) prompts.PromptTemplate {
	template := strings.Join([]string{prefix, instructions, suffix}, "\n\n")

	return prompts.PromptTemplate{
		Template:       template,
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		InputVariables: []string{"input", "agent_scratchpad"},
		PartialVariables: map[string]any{
			"tool_names":        toolNames(tools),
			"tool_descriptions": toolDescriptions(tools),
			"history":           "",
		},
	}
}

// extractNextTool is a helper function to extract the next tool to call from the output.
func extractNextTool(output string) string {
	// Implement the logic to extract the next tool from the output.
	// This is a placeholder implementation and should be replaced with actual logic.
	if strings.Contains(output, "next_tool:") {
		parts := strings.Split(output, "next_tool:")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}
