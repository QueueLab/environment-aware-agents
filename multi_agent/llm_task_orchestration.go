package multi_agent

import (
	"context"
	"fmt"
	"time"
)

// batchRequests batches requests into a single context window, maximizing the token limit
func batchRequests(ctx context.Context, requests []string) {
	var batch []string
	tokenLimit := 100 // Simulated context window token limit

	for _, req := range requests {
		batch = append(batch, req)
		if len(batch) >= tokenLimit {
			processBatch(ctx, batch) // Send batch to LLM for processing
			batch = nil              // Reset batch
		}
	}
	if len(batch) > 0 {
		processBatch(ctx, batch) // Process any remaining tasks
	}
}

// processBatch processes a batch of requests
func processBatch(ctx context.Context, batch []string) {
	fmt.Println("Processing batch of requests:", batch)
	select {
	case <-ctx.Done():
		fmt.Println("Batch processing canceled:", ctx.Err())
		return
	default:
		// LLM processing logic here
	}
}

func main() {
	parentCtx := context.Background()

	// Create task with timeout or priority
	childCtx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	requests := []string{"task1", "task2", "task3", "task4"} // Simulated requests
	batchRequests(childCtx, requests)
}
