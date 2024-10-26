package graph_orchestration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamOutput(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	graph.AddNode("B")
	graph.AddEdge("A", "B")

	streamingManager := NewStreamingManager()
	outputChan := make(chan string)

	go func() {
		err := streamingManager.StreamOutput(context.Background(), graph, "A", outputChan)
		assert.NoError(t, err)
	}()

	var outputs []string
	for output := range outputChan {
		outputs = append(outputs, output)
	}

	assert.Contains(t, outputs, "Streaming output from node: B")
}

func TestStreamToken(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	graph.AddNode("B")
	graph.AddEdge("A", "B")

	streamingManager := NewStreamingManager()
	tokenChan := make(chan string)

	go func() {
		err := streamingManager.StreamToken(context.Background(), graph, "A", tokenChan)
		assert.NoError(t, err)
	}()

	var tokens []string
	for token := range tokenChan {
		tokens = append(tokens, token)
	}

	assert.Contains(t, tokens, "Streaming token from node: B")
}
