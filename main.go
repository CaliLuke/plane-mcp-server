package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CaliLuke/plane-mcp-server/plane"
	"github.com/CaliLuke/plane-mcp-server/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	apiKey := os.Getenv("PLANE_API_KEY")
	workspaceSlug := os.Getenv("PLANE_WORKSPACE_SLUG")
	baseURL := os.Getenv("PLANE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.plane.so"
	}

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "PLANE_API_KEY is not set")
		os.Exit(1)
	}
	if workspaceSlug == "" {
		fmt.Fprintln(os.Stderr, "PLANE_WORKSPACE_SLUG is not set")
		os.Exit(1)
	}

	client := plane.NewClient(baseURL, apiKey, "")

	s := server.NewMCPServer(
		"Plane MCP Server (stdio/go)",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	toolsConfig := os.Getenv("PLANE_TOOLS")
	reg := tools.NewRegistry(client, workspaceSlug)
	reg.Register(s, toolsConfig)

	// Serve via stdio
	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

// newTool is a helper to create a tool with description.
func newTool(name, description string, opts ...mcp.ToolOption) mcp.Tool {
	allOpts := append([]mcp.ToolOption{mcp.WithDescription(description)}, opts...)
	return mcp.NewTool(name, allOpts...)
}
