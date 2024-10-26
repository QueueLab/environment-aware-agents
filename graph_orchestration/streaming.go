package graph_orchestration

import (
	"context"
	"sync"
)

// StreamingManager handles streaming of outputs as they are produced by each node.
type StreamingManager struct {
	mu sync.Mutex
}

// NewStreamingManager creates a new StreamingManager.
func NewStreamingManager() *StreamingManager {
	return &StreamingManager{}
}

// StreamOutput streams the output as it is produced by each node.
func (sm *StreamingManager) StreamOutput(ctx context.Context, graph *Graph, nodeID string, outputChan chan<- string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	node, exists := graph.Nodes[nodeID]
	if !exists {
		return errors.New("node does not exist")
	}

	// Placeholder for streaming logic
	// This is where you would integrate the streaming of outputs
	go func() {
		for _, child := range node.Children {
			outputChan <- "Streaming output from node: " + child.ID
		}
		close(outputChan)
	}()

	return nil
}

// StreamToken streams tokens as they are produced by each node.
func (sm *StreamingManager) StreamToken(ctx context.Context, graph *Graph, nodeID string, tokenChan chan<- string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	node, exists := graph.Nodes[nodeID]
	if !exists {
		return errors.New("node does not exist")
	}

	// Placeholder for token streaming logic
	// This is where you would integrate the streaming of tokens
	go func() {
		for _, child := range node.Children {
			tokenChan <- "Streaming token from node: " + child.ID
		}
		close(tokenChan)
	}()

	return nil
}
