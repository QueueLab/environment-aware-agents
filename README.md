This is research in semi autonomous highly concurrent multi agent system with varying levels of memory implemented across the stack. This pushes the boundary of in context learning and alignment for human in the loop processes. This research involves creating nature inspired workflows that mimic highly concurrent processes to enable solutions for high dimensional and dynamical problems.

## Documentation

- [Documentation Site](https://tmc.github.io/langchaingo/docs/)
- [API Reference](https://pkg.go.dev/github.com/tmc/langchaingo)


##  Examples

See [./examples](./examples) for example usage.

```go
package main

import (
  "context"
  "fmt"
  "log"

  "github.com/tmc/langchaingo/llms"
  "github.com/tmc/langchaingo/llms/openai"
)

func main() {
  ctx := context.Background()
  llm, err := openai.New()
  if err != nil {
    log.Fatal(err)
  }
  prompt := "What would be a good company name for a company that makes colorful socks?"
  completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(completion)
}
```

```shell
$ go run .
Socktastic
```




