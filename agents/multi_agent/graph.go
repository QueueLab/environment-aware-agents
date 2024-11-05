package multi_agent

import (
	"sync"
)

// Node represents a node in the graph.
type Node struct {
	ID       int
	Value    interface{}
	Priority int
}

// Edge represents an edge in the graph.
type Edge struct {
	From   *Node
	To     *Node
	Action func()
	Weight int
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

// TopologicalSort performs a topological sort on the graph based on task planned complexity.
func (g *Graph) TopologicalSort() []*Node {
	inDegree := make(map[*Node]int)
	for _, node := range g.Nodes {
		inDegree[node] = 0
	}
	for _, edge := range g.Edges {
		inDegree[edge.To]++
	}

	var queue []*Node
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	var sortedNodes []*Node
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		sortedNodes = append(sortedNodes, node)

		for _, edge := range g.Edges {
			if edge.From == node {
				inDegree[edge.To]--
				if inDegree[edge.To] == 0 {
					queue = append(queue, edge.To)
				}
			}
		}
	}

	return sortedNodes
}

// AddWeightedContext adds weighted context to the graph nodes based on priority.
func (g *Graph) AddWeightedContext(nodes []*Node, priorities []int) {
	for i, node := range nodes {
		node.Priority = priorities[i]
		g.AddNode(node)
	}
}

// ExecuteWithPriority executes the graph nodes based on their priority.
func (g *Graph) ExecuteWithPriority() {
	sortedNodes := g.TopologicalSort()
	var wg sync.WaitGroup
	for _, node := range sortedNodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			for _, edge := range g.Edges {
				if edge.From == n {
					edge.Action()
				}
			}
		}(node)
	}
	wg.Wait()
}
