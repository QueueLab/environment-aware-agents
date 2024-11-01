package graph_orchestration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

func TestAddNode(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")

	assert.Contains(t, graph.Nodes, "A")
}

func TestAddEdge(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	graph.AddNode("B")

	err := graph.AddEdge("A", "B")
	assert.NoError(t, err)
	assert.Contains(t, graph.Nodes["A"].Children, graph.Nodes["B"])
}

func TestAddEdgeNonExistentNode(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")

	err := graph.AddEdge("A", "B")
	assert.Error(t, err)
}

func TestExecute(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	graph.AddNode("B")
	graph.AddEdge("A", "B")

	llm := &llms.MockModel{}
	tools := []tools.Tool{&tools.MockTool{}}
	agent := NewGraphAgent(llm, tools)

	err := graph.Execute(context.Background(), "A", agent)
	assert.NoError(t, err)
}

func TestExecuteNonExistentNode(t *testing.T) {
	graph := NewGraph()

	llm := &llms.MockModel{}
	tools := []tools.Tool{&tools.MockTool{}}
	agent := NewGraphAgent(llm, tools)

	err := graph.Execute(context.Background(), "A", agent)
	assert.Error(t, err)
}
