package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (r *Registry) registerInitiativeTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("list_initiatives",
		mcp.WithDescription("List all initiatives in a workspace."),
	), r.handleListInitiatives)

	s.AddTool(mcp.NewTool("create_initiative",
		mcp.WithDescription("Create a new initiative in the workspace."),
		mcp.WithString("name", mcp.Description("Initiative name"), mcp.Required()),
		mcp.WithString("description_html", mcp.Description("HTML description of the initiative")),
		mcp.WithString("start_date", mcp.Description("Initiative start date (ISO 8601)")),
		mcp.WithString("end_date", mcp.Description("Initiative end date (ISO 8601)")),
		mcp.WithString("state", mcp.Description("Initiative state"), mcp.Enum("DRAFT", "PLANNED", "ACTIVE", "COMPLETED", "CLOSED")),
		mcp.WithString("lead", mcp.Description("UUID of the user who leads the initiative")),
	), r.handleCreateInitiative)

	s.AddTool(mcp.NewTool("retrieve_initiative",
		mcp.WithDescription("Retrieve an initiative by ID."),
		mcp.WithString("initiative_id", mcp.Description("UUID of the initiative"), mcp.Required()),
	), r.handleRetrieveInitiative)

	s.AddTool(mcp.NewTool("update_initiative",
		mcp.WithDescription("Update an initiative by ID."),
		mcp.WithString("initiative_id", mcp.Description("UUID of the initiative"), mcp.Required()),
		mcp.WithString("name", mcp.Description("Initiative name")),
		mcp.WithString("description_html", mcp.Description("HTML description of the initiative")),
		mcp.WithString("start_date", mcp.Description("Initiative start date (ISO 8601)")),
		mcp.WithString("end_date", mcp.Description("Initiative end date (ISO 8601)")),
		mcp.WithString("state", mcp.Description("Initiative state"), mcp.Enum("DRAFT", "PLANNED", "ACTIVE", "COMPLETED", "CLOSED")),
		mcp.WithString("lead", mcp.Description("UUID of the user who leads the initiative")),
	), r.handleUpdateInitiative)

	s.AddTool(mcp.NewTool("delete_initiative",
		mcp.WithDescription("Delete an initiative by ID."),
		mcp.WithString("initiative_id", mcp.Description("UUID of the initiative"), mcp.Required()),
	), r.handleDeleteInitiative)
}

func (r *Registry) handleListInitiatives(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/initiatives", r.WorkspaceSlug), nil)
	if err != nil {
		return errorResult(err)
	}
	var data struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return rawJSONResult(resp)
	}
	return rawJSONResult(data.Results)
}

func (r *Registry) handleCreateInitiative(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	body := stripNulls(map[string]any{
		"name":             getString(req, "name", ""),
		"description_html": optionalString(req, "description_html"),
		"start_date":       optionalString(req, "start_date"),
		"end_date":         optionalString(req, "end_date"),
		"state":            optionalString(req, "state"),
		"lead":             optionalString(req, "lead"),
	})
	resp, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/initiatives", r.WorkspaceSlug), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleRetrieveInitiative(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	initiativeID := decodeID(req, "initiative_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/initiatives/%s", r.WorkspaceSlug, initiativeID), nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleUpdateInitiative(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	initiativeID := decodeID(req, "initiative_id")
	body := stripNulls(map[string]any{
		"name":             optionalString(req, "name"),
		"description_html": optionalString(req, "description_html"),
		"start_date":       optionalString(req, "start_date"),
		"end_date":         optionalString(req, "end_date"),
		"state":            optionalString(req, "state"),
		"lead":             optionalString(req, "lead"),
	})
	resp, err := r.Client.Patch(ctx, fmt.Sprintf("workspaces/%s/initiatives/%s", r.WorkspaceSlug, initiativeID), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleDeleteInitiative(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	initiativeID := decodeID(req, "initiative_id")
	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/initiatives/%s", r.WorkspaceSlug, initiativeID)); err != nil {
		return errorResult(err)
	}
	return textResult("Initiative deleted successfully")
}
