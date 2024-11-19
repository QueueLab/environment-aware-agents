package graph_orchestration

import (
	"context"
	"errors"
	"sync"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

// Node represents a single node in the graph.
type Node struct {
	ID       string
	Children []*Node
	Action   schema.AgentAction
}

// Graph represents a directional graph.
type Graph struct {
	Nodes map[string]*Node
	mu    sync.Mutex
}

// NewGraph creates a new Graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

// AddNode adds a new node to the graph.
func (g *Graph) AddNode(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.Nodes[id]; !exists {
		g.Nodes[id] = &Node{ID: id}
	}
}

// AddEdge adds a directional edge between two nodes.
func (g *Graph) AddEdge(fromID, toID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	fromNode, fromExists := g.Nodes[fromID]
	toNode, toExists := g.Nodes[toID]

	if !fromExists || !toExists {
		return errors.New("one or both nodes do not exist")
	}

	fromNode.Children = append(fromNode.Children, toNode)
	return nil
}

// Execute executes the graph starting from the given node ID.
func (g *Graph) Execute(ctx context.Context, startID string, agent *GraphAgent) error {
	g.mu.Lock()
	startNode, exists := g.Nodes[startID]
	g.mu.Unlock()

	if !exists {
		return errors.New("start node does not exist")
	}

	return g.executeNode(ctx, startNode, agent)
}

func (g *Graph) executeNode(ctx context.Context, node *Node, agent *GraphAgent) error {
	// Execute the action associated with the node
	_, _, err := agent.Plan(ctx, nil, map[string]string{"input": node.ID})
	if err != nil {
		return err
	}

	// Execute child nodes
	for _, child := range node.Children {
		if err := g.executeNode(ctx, child, agent); err != nil {
			return err
		}
	}
	return nil
}

// GetTools returns the tools available to the graph.
func (g *Graph) GetTools() []tools.Tool {

	g.mu.Lock()
	defer g.mu.Unlock()

	var allTools []tools.Tool
	for _, node := range g.Nodes {
		if node.Agent != nil && node.Agent.Tools != nil {
			allTools = append(allTools, node.Agent.Tools...)
		}
	}
	return allTools
}
