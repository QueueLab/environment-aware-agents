package multi_agent

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// GraphPersister is responsible for persisting the state of the graph.
type GraphPersister struct {
	mu    sync.Mutex
	graph *Graph
}

// NewGraphPersister creates a new GraphPersister.
func NewGraphPersister(graph *Graph) *GraphPersister {
	return &GraphPersister{
		graph: graph,
	}
}

// Save saves the state of the graph to a file.
func (p *GraphPersister) Save(filename string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, err := json.Marshal(p.graph)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// Load loads the state of the graph from a file.
func (p *GraphPersister) Load(filename string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, p.graph)
}

// SaveNode saves the state of a single node to a file.
func (p *GraphPersister) SaveNode(filename string, nodeID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	node, exists := p.graph.Nodes[nodeID]
	if !exists {
		return os.ErrNotExist
	}

	data, err := json.Marshal(node)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// LoadNode loads the state of a single node from a file.
func (p *GraphPersister) LoadNode(filename string, nodeID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	node := &Node{}
	if err := json.Unmarshal(data, node); err != nil {
		return err
	}

	p.graph.Nodes[nodeID] = node
	return nil
}
