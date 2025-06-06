# go-agents

A basic library to create AI agents

---

## Overview

go-agents is a lightweight library built for developers looking to integrate simple AI agents into their applications. With a focus on modularity and ease-of-use, this library serves as a foundation to build, test, and deploy AI-driven functionalities.

## Features

- **Modular Design:** Easily customize and extend agent functionalities.
- **Simple Integration:** Minimal setup required to start building intelligent agents.
- **Scalable:** Suitable for projects of all sizes, from prototypes to production-grade applications.

## Installation

Install go-agents via go get:

```bash
go get github.com/jlrosende/go-agents
```

Or add the dependency to your project’s `go.mod`:

```go
require github.com/jlrosende/go-agents vX.Y.Z
```

## Quick Start

Below is a basic example to demonstrate how to create and run an agent:

```go
package main

import (
	"fmt"
	"github.com/jlrosende/go-agents/agent"
)

func main() {
	// Create a new agent
	swarm := controller.AgentController()

	// Configure agent (add behaviors, intents, etc.)
	// agentContrller.AddAgent(...)

	// Run the agent
	swarm.Run()

}
```

For more comprehensive examples, please refer to the [examples](./examples) directory.

## Documentation

For detailed documentation and API reference, please visit the [Wiki](https://github.com/jlrosende/go-agents/wiki) or view the inline code comments.

## Contributing

Contributions are welcome! Please check out our [Contributing Guidelines](./CONTRIBUTING.md) for more details on how to get started, and read our [Code of Conduct](./CODE_OF_CONDUCT.md) to ensure a welcoming environment for all.

## Roadmap

- [ ] Implement advanced agent behaviors
- [ ] Integrate with external AI services
- [ ] Improve documentation with more examples

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## Acknowledgements

- Hat tip to all the contributors and the open-source community for their invaluable support.
- Thanks to the developers behind similar projects that inspired this library.
