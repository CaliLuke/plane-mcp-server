# plane-mcp-server

The official Plane MCP server returns way more data than any AI agent needs, burning tokens on noise. This is a leaner fork with compact response models and a few quality-of-life improvements. Use it if it's useful to you — no guarantees.

## What's different

- Compact response models: list calls return only what you need (id, name, identifier for projects; id, sequence_id, name, priority, state for work items — no HTML blobs)
- Detail calls strip HTML descriptions to plain text
- `PLANE_TOOLS` env var to exclude tool groups you don't need (e.g. `!cycles,!modules,!initiatives`)
- State names and user display names resolved inline — no extra lookups needed
- Short UUIDs throughout for smaller payloads

## Setup

### Pre-built binary

Download a release binary or build from source:

```sh
go build -o plane-mcp-server .
```

Then add to your MCP config:

```json
{
  "mcpServers": {
    "plane": {
      "command": "/path/to/plane-mcp-server",
      "env": {
        "PLANE_API_KEY": "<your-api-key>",
        "PLANE_WORKSPACE_SLUG": "<your-workspace-slug>",
        "PLANE_BASE_URL": "https://api.plane.so",
        "PLANE_TOOLS": "!cycles,!modules,!initiatives"
      }
    }
  }
}
```

### From source (go run)

```json
{
  "mcpServers": {
    "plane": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/path/to/plane-mcp-server",
      "env": {
        "PLANE_API_KEY": "<your-api-key>",
        "PLANE_WORKSPACE_SLUG": "<your-workspace-slug>",
        "PLANE_BASE_URL": "https://api.plane.so",
        "PLANE_TOOLS": "!cycles,!modules,!initiatives"
      }
    }
  }
}
```

## Legacy Python version

The original Python implementation lives in `legacy/`. It is no longer actively maintained.

## License

MIT
