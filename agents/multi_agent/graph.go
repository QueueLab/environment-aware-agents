package multi_agent

import (
	"sync"
)

// Node represents a node in the graph.
type Node struct {
	ID      int
	Value   interface{}
	State   string
	Actions []interface{}
	mu      sync.Mutex
}

// Edge represents an edge in the graph.
type Edge struct {
	From   *Node
	To     *Node
	Action func()
}

// Graph represents a concurrent graph structure.
type Graph struct {
	Nodes []*Node
	Edges []*Edge
	mu    sync.Mutex
}

// NewGraph creates a new graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: []*Node{},
		Edges: []*Edge{},
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(node *Node) {
	g.mu.Lock()
	defer g.mu.Unlock()
	node.State = "initialized"
	node.Actions = []interface{}{}
	g.Nodes = append(g.Nodes, node)
}

// AddEdge adds an edge to the graph.
func (g *Graph) AddEdge(edge *Edge) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Edges = append(g.Edges, edge)
}

// Execute executes the graph concurrently.
func (g *Graph) Execute() {
	var wg sync.WaitGroup
	for _, edge := range g.Edges {
		wg.Add(1)
		go func(e *Edge) {
			defer wg.Done()
			e.Action()
		}(edge)
	}
	wg.Wait()
}
