package multi_agent

import (
	"context"
	"fmt"
	"time"
)

// Demonstrates creating parent-child context relationships
func main() {
	// Root context
	parentCtx := context.Background()

	// Create a child context with a 2-second timeout
	childCtx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()

	go agentTask(childCtx)

	time.Sleep(3 * time.Second) // Wait to see if the task finishes
}

// agentTask simulates a task that checks for context cancellation
func agentTask(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Println("Task cancelled or timeout:", ctx.Err()) // Runs when the context timeout occurs
		return
	default:
		fmt.Println("Executing task")
	}
}
