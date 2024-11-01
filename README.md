# Environment-Aware Agents

This repository contains ongoing research on multi-agent systems. The goal is to develop agents that are aware of their environment and can interact with it in meaningful ways.

## Documentation

- [Documentation Site](https://tmc.github.io/langchaingo/docs/)
- [API Reference](https://pkg.go.dev/github.com/tmc/langchaingo)

## Getting Started

### Prerequisites

- Go 1.22 or later

### Installation

To install the package, run:

```shell
go get github.com/tmc/langchaingo
```

### Usage

Here is a simple example to get you started:

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

To run the example, save it to a file (e.g., `main.go`) and execute:

```shell
go run main.go
```

## Running Tests

To run tests, use the following command:

```shell
go test ./...
```

This will run all the tests in the repository.

## Contributing

We welcome contributions from the community! Please read the [CONTRIBUTING.md](./CONTRIBUTING.md) file for more information on how to get started.

## Examples

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
