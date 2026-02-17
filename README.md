# plane-mcp-server

I think the original Plane MCP server was quite bad — verbose, poorly thought out, and returning oceans of tokens for no reason. So I made my own. Feel free to use it without any guarantees whatsoever.

## What's different

- Compact response models: list calls return only what you need (id, name, identifier for projects; id, sequence_id, name, priority, state for work items — no HTML blobs)
- Detail calls strip HTML descriptions to plain text
- `PLANE_TOOLS` env var to exclude tool groups you don't need (e.g. `!cycles,!modules,!initiatives`)
- fastmcp v3

## Setup

```json
{
  "mcpServers": {
    "plane": {
      "command": "uvx",
      "args": [
        "--from",
        "git+https://github.com/CaliLuke/plane-mcp-server",
        "plane-mcp-server",
        "stdio"
      ],
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

## License

MIT
