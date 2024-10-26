package graph_orchestration

import (
	"errors"
	"sync"
)

// Node represents a single node in the graph.
type Node struct {
	ID       string
	Children []*Node
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
func (g *Graph) Execute(startID string) error {
	g.mu.Lock()
	startNode, exists := g.Nodes[startID]
	g.mu.Unlock()

	if !exists {
		return errors.New("start node does not exist")
	}

	return g.executeNode(startNode)
}

func (g *Graph) executeNode(node *Node) error {
	// Placeholder for node execution logic
	// This is where you would integrate the agent's actions
	for _, child := range node.Children {
		if err := g.executeNode(child); err != nil {
			return err
		}
	}
	return nil
}
