# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

A Go MCP (Model Context Protocol) server for the [Plane](https://plane.so) project management API. It's a leaner fork of the official Plane MCP server, optimized for AI agent consumption with compact responses, short UUIDs, and inline state/user name resolution.

## Build & Run

```sh
go build -o plane-mcp-server .   # build binary
go run .                          # run from source
go install github.com/CaliLuke/plane-mcp-server@latest  # install globally
```

No test suite exists yet. No Makefile or CI at the Go level.

## Required Environment Variables

- `PLANE_API_KEY` — Plane API key (required)
- `PLANE_WORKSPACE_SLUG` — Workspace slug (required)
- `PLANE_BASE_URL` — Defaults to `https://api.plane.so`
- `PLANE_TOOLS` — Tool group filtering: `!cycles,!modules` (exclusion) or `projects,work_items` (inclusion)

## Architecture

The server communicates over **stdio** using the `mark3labs/mcp-go` SDK. Flow: `main.go` reads env → creates Plane HTTP client → creates MCP server → registers tool groups → serves via stdio.

### Key Packages

- **`plane/`** — HTTP client wrapping `net/http`. Auth via `X-Api-Key` header. Methods: `Get`, `Post`, `Patch`, `Delete` returning `json.RawMessage`.
- **`tools/`** — MCP tool registrations and handlers. Each file covers one domain (projects, work_items, cycles, etc.).
- **`tools/registry.go`** — Tool group registry. `toolGroups` map defines 9 groups. `PLANE_TOOLS` env var controls which groups are enabled.
- **`tools/helpers.go`** — Shared utilities: `getString`, `requireString`, `requireID`, `decodeID`, `jsonResult`, `textResult`, `listResult`, `errorResult`.
- **`models/`** — Slim response types (`ProjectSummary`, `WorkItemSummary`, `WorkItemFull`, `StateSummary`). List calls return compact summaries; detail calls return full objects.
- **`uid/`** — Short UUID encoding/decoding via `shortuuid` (base57, ~22 chars). `EncodeValue()` recursively walks JSON to encode all UUIDs. All tools accept both short and full UUID formats.
- **`formatting/`** — Compact single-line display formatters. Pattern: `"PROJ-42 (In Progress) · Name [short_id]"`. Also has `StripHTML()` for cleaning rich text.

### Tool Handler Pattern

Every tool file follows this structure:

1. `registerXxxTools(s *server.MCPServer)` method on `Registry` adds tools via `s.AddTool(mcp.NewTool(...), r.handleXxx)`
2. Handler receives `(ctx context.Context, req mcp.CallToolRequest)`
3. Parameters extracted via helpers (`requireString`, `getString`, `getBool`, `requireID`)
4. `requireID`/`decodeID` transparently accept both short (~22-char) and full UUIDs
5. API call via `r.Client.Get/Post/Patch/Delete(ctx, endpoint, body)`
6. Response returned via `jsonResult()`, `textResult()`, `listResult()`, or `errorResult()`

### Adding a New Tool Group

1. Create `tools/new_domain.go` with a `registerNewDomainTools` method on `Registry`
2. Add the group to `toolGroups` map in `tools/registry.go`
3. Follow existing handler patterns — use helpers from `tools/helpers.go`
4. Add slim response models to `models/models.go` if needed
5. Add compact formatters to `formatting/formatting.go` if needed

## Design Principles

- **Token efficiency** — Compact responses by default; full detail only on retrieve
- **Short UUIDs everywhere** — All outgoing UUIDs encoded to ~22-char base57; inputs accept both formats
- **Inline resolution** — State names and user display names resolved in responses, no extra lookups needed
- **HTML stripped** — Rich text descriptions converted to plain text via `StripHTML()`
- **Progressive disclosure** — List calls return one-line summaries; retrieve calls return full objects
