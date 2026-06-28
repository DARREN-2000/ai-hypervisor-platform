# CLI Documentation

The AI Hypervisor Platform provides a command-line interface (`ai-hyp`) for managing cluster resources directly from your terminal.

## Installation

Currently, the CLI must be built from source.

```bash
go build -o ai-hyp ./cmd/cli
sudo mv ai-hyp /usr/local/bin/
```

## Configuration

The CLI requires an API endpoint and an authentication token. You can provide these via environment variables:

```bash
export AI_HYP_ENDPOINT="https://api.ai-hypervisor.example.com"
export AI_HYP_TOKEN="your-secret-token"
```

## Common Commands

### Virtual Machines

List all VMs:
```bash
ai-hyp vm list
```

Create a new VM from a JSON configuration file:
```bash
ai-hyp vm create -f vm-config.json
```

Stop a VM:
```bash
ai-hyp vm stop <vm-id>
```

### Hosts

List all hosts and their status:
```bash
ai-hyp host list
```

Get live metrics for a specific host:
```bash
ai-hyp host metrics <host-id>
```

### GPUs

List available GPUs:
```bash
ai-hyp gpu list
```

## Help

For a full list of commands and options, use the `--help` flag:

```bash
ai-hyp --help
ai-hyp vm --help
```
