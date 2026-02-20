package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (r *Registry) registerStateTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("list_states",
			mcp.WithDescription("List valid state names for a project. Use this to discover which state names can be passed to create_work_item / update_work_item."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		),
		r.handleListStates,
	)
}

func (r *Registry) handleListStates(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}

	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/states", r.WorkspaceSlug, projectID), nil)
	if err != nil {
		return errorResult(err)
	}

	var data struct {
		Results []struct {
			Name    string `json:"name"`
			Group   string `json:"group"`
			Default bool   `json:"default"`
		} `json:"results"`
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return errorResult(err)
	}

	lines := make([]string, len(data.Results))
	for i, s := range data.Results {
		def := ""
		if s.Default {
			def = " *"
		}
		lines[i] = fmt.Sprintf("%s (%s)%s", s.Name, s.Group, def)
	}
	return listResult(lines)
}
