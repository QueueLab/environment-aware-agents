package multi_agent

import (
	"context"
	"fmt"
	"time"
)

// TaskWeight represents the weight and priority of a task
type TaskWeight struct {
	Priority int
	Deadline time.Time
}

// weightedContext creates a custom context to carry priority information
func weightedContext(ctx context.Context, weight TaskWeight) context.Context {
	return context.WithValue(ctx, "weight", weight)
}

// agentTask simulates a task that checks for context cancellation and prioritization
func agentTask(ctx context.Context) {
 weight, ok := ctx.Value("weight").(TaskWeight)
 if !ok {
     fmt.Println("Error: Invalid weight in context")
     return
 }
	fmt.Printf("Executing task with priority: %d\n", weight.Priority)

	// Handle timeouts, deadlines
	select {
	case <-ctx.Done():
		fmt.Println("Task canceled due to:", ctx.Err())
		return
	default:
		// Task continues based on weight and priority
	}
}

func main() {
	parentCtx := context.Background()

	// Create a child context with a task weight
	taskCtx := weightedContext(parentCtx, TaskWeight{Priority: 10, Deadline: time.Now().Add(5 * time.Second)})
	go agentTask(taskCtx)

	time.Sleep(2 * time.Second)
}
