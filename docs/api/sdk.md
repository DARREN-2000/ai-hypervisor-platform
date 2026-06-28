# SDK Guide

While the AI Hypervisor Platform can be interacted with directly via the REST API, we recommend using one of our official SDKs for a smoother developer experience.

## Supported Languages

- **Go** (Primary)
- **Python**
- **TypeScript**

*Note: The Python and TypeScript SDKs are currently auto-generated from our OpenAPI specification and are in beta.*

## Go SDK

The Go SDK is part of the main repository and is used internally by our CLI and control plane components.

### Installation

```bash
go get github.com/ai-hypervisor/platform/pkg/client
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ai-hypervisor/platform/pkg/client"
)

func main() {
    // Initialize the client
    c, err := client.NewClient("https://api.ai-hypervisor.example.com/api/v1", "your-api-token")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // List VMs
    vms, err := c.ListVMs(context.Background())
    if err != nil {
        log.Fatalf("Failed to list VMs: %v", err)
    }

    for _, vm := range vms {
        fmt.Printf("VM ID: %s, Status: %s\n", vm.ID, vm.Status)
    }
}
```

## Future Roadmaps

We are actively working on dedicated, idiomatic SDKs for Python and TypeScript. In the meantime, you can use standard HTTP libraries (`requests`, `fetch`) as demonstrated in the [API Reference](endpoints.md).
