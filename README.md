





This is the Go language implementation of [LangChain](https://github.com/langchain-ai/langchain).<img width="333" alt="Capture d’écran 2024-09-01 à 10 37 36" src="https://github.com/user-attachments/assets/803478f1-40b3-40b6-9559-f022c71ec888">




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




