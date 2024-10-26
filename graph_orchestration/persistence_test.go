package graph_orchestration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveState(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	pm := NewPersistenceManager("test_state.json")

	err := pm.SaveState(graph)
	assert.NoError(t, err)

	loadedGraph := NewGraph()
	err = pm.LoadState(loadedGraph)
	assert.NoError(t, err)
	assert.Contains(t, loadedGraph.Nodes, "A")
}

func TestLoadState(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	pm := NewPersistenceManager("test_state.json")

	err := pm.SaveState(graph)
	assert.NoError(t, err)

	loadedGraph := NewGraph()
	err = pm.LoadState(loadedGraph)
	assert.NoError(t, err)
	assert.Contains(t, loadedGraph.Nodes, "A")
}

func TestPauseAndResumeExecution(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	pm := NewPersistenceManager("test_state.json")

	err := pm.PauseExecution(graph)
	assert.NoError(t, err)

	loadedGraph := NewGraph()
	err = pm.ResumeExecution(loadedGraph)
	assert.NoError(t, err)
	assert.Contains(t, loadedGraph.Nodes, "A")
}

func TestErrorRecovery(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	pm := NewPersistenceManager("test_state.json")

	err := pm.SaveState(graph)
	assert.NoError(t, err)

	loadedGraph := NewGraph()
	err = pm.ErrorRecovery(loadedGraph)
	assert.NoError(t, err)
	assert.Contains(t, loadedGraph.Nodes, "A")
}

func TestTimeTravel(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A")
	pm := NewPersistenceManager("test_state.json")

	err := pm.SaveState(graph)
	assert.NoError(t, err)

	loadedGraph := NewGraph()
	err = pm.TimeTravel(loadedGraph, "test_state.json")
	assert.NoError(t, err)
	assert.Contains(t, loadedGraph.Nodes, "A")
}
