// Package tools registers MCP tools for the Plane server.
package tools

import (
	"log"
	"strings"

	"github.com/CaliLuke/plane-mcp-server/plane"
	"github.com/mark3labs/mcp-go/server"
)

// Registry holds shared state for all tool handlers.
type Registry struct {
	Client        *plane.Client
	WorkspaceSlug string
}

// NewRegistry creates a new tool registry.
func NewRegistry(client *plane.Client, workspaceSlug string) *Registry {
	return &Registry{Client: client, WorkspaceSlug: workspaceSlug}
}

type registerFunc func(r *Registry, s *server.MCPServer)

var toolGroups = map[string]registerFunc{
	"projects":             (*Registry).registerProjectTools,
	"work_items":           (*Registry).registerWorkItemTools,
	"states":               (*Registry).registerStateTools,
	"cycles":               (*Registry).registerCycleTools,
	"users":                (*Registry).registerUserTools,
	"modules":              (*Registry).registerModuleTools,
	"initiatives":          (*Registry).registerInitiativeTools,
	"intake":               (*Registry).registerIntakeTools,
	"work_item_properties": (*Registry).registerWorkItemPropertyTools,
}

// Register registers tool groups with the MCP server based on the PLANE_TOOLS config string.
func (r *Registry) Register(s *server.MCPServer, toolsConfig string) {
	enabled := resolveEnabledGroups(toolsConfig)
	for name, fn := range toolGroups {
		if enabled[name] {
			fn(r, s)
		}
	}
}

func resolveEnabledGroups(config string) map[string]bool {
	config = strings.TrimSpace(config)
	allGroups := make(map[string]bool)
	for name := range toolGroups {
		allGroups[name] = true
	}

	if config == "" {
		return allGroups
	}

	if strings.HasPrefix(config, "!") {
		// Exclusion mode
		items := strings.Split(config, ",")
		excluded := make(map[string]bool)
		for _, item := range items {
			item = strings.TrimSpace(item)
			if strings.HasPrefix(item, "!") {
				name := item[1:]
				if !allGroups[name] {
					log.Printf("WARNING: Unknown tool group in PLANE_TOOLS: %s", name)
				}
				excluded[name] = true
			}
		}
		for name := range excluded {
			delete(allGroups, name)
		}
		return allGroups
	}

	// Inclusion mode
	enabled := make(map[string]bool)
	for _, item := range strings.Split(config, ",") {
		name := strings.TrimSpace(item)
		if name == "" {
			continue
		}
		if _, ok := toolGroups[name]; !ok {
			log.Printf("WARNING: Unknown tool group in PLANE_TOOLS: %s", name)
			continue
		}
		enabled[name] = true
	}
	return enabled
}
