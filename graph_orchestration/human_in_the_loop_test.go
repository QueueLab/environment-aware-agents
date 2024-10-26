package graph_orchestration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterruptExecution(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	hm := NewHumanInTheLoopManager()

	err := hm.InterruptExecution(context.Background(), graph, "A")
	assert.NoError(t, err)
}

func TestInterruptExecutionNonExistentNode(t *testing.T) {
	graph := NewGraph()
	hm := NewHumanInTheLoopManager()

	err := hm.InterruptExecution(context.Background(), graph, "A")
	assert.Error(t, err)
}

func TestApproveAction(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	hm := NewHumanInTheLoopManager()

	err := hm.ApproveAction(context.Background(), graph, "A")
	assert.NoError(t, err)
}

func TestApproveActionNonExistentNode(t *testing.T) {
	graph := NewGraph()
	hm := NewHumanInTheLoopManager()

	err := hm.ApproveAction(context.Background(), graph, "A")
	assert.Error(t, err)
}

func TestEditAction(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	hm := NewHumanInTheLoopManager()

	err := hm.EditAction(context.Background(), graph, "A", "newAction")
	assert.NoError(t, err)
}

func TestEditActionNonExistentNode(t *testing.T) {
	graph := NewGraph()
	hm := NewHumanInTheLoopManager()

	err := hm.EditAction(context.Background(), graph, "A", "newAction")
	assert.Error(t, err)
}
