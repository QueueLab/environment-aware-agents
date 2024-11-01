package graph_orchestration

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

// PersistenceManager handles the persistence of graph states.
type PersistenceManager struct {
	mu       sync.Mutex
	filePath string
}

// NewPersistenceManager creates a new PersistenceManager.
func NewPersistenceManager(filePath string) *PersistenceManager {
	return &PersistenceManager{
		filePath: filePath,
	}
}

// SaveState saves the current state of the graph to a file.
func (pm *PersistenceManager) SaveState(graph *Graph) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	data, err := json.Marshal(graph)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(pm.filePath, data, 0644)
}

// LoadState loads the state of the graph from a file.
func (pm *PersistenceManager) LoadState(graph *Graph) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	data, err := ioutil.ReadFile(pm.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, graph)
}

// PauseExecution pauses the execution of the graph.
func (pm *PersistenceManager) PauseExecution(graph *Graph) error {
	return pm.SaveState(graph)
}

// ResumeExecution resumes the execution of the graph.
func (pm *PersistenceManager) ResumeExecution(graph *Graph) error {
	return pm.LoadState(graph)
}

// ErrorRecovery attempts to recover from an error by loading the last saved state.
func (pm *PersistenceManager) ErrorRecovery(graph *Graph) error {
	return pm.LoadState(graph)
}

// TimeTravel allows traveling back to a previous state of the graph.
func (pm *PersistenceManager) TimeTravel(graph *Graph, stateFilePath string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	data, err := ioutil.ReadFile(stateFilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, graph)
}
