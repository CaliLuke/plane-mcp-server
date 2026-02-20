package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (r *Registry) registerIntakeTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("list_intake_work_items",
		mcp.WithDescription("List all intake work items in a project."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
	), r.handleListIntakeWorkItems)

	s.AddTool(mcp.NewTool("create_intake_work_item",
		mcp.WithDescription("Create a new intake work item in a project."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithObject("data", mcp.Description("Intake work item data as a JSON object"), mcp.Required()),
	), r.handleCreateIntakeWorkItem)

	s.AddTool(mcp.NewTool("retrieve_intake_work_item",
		mcp.WithDescription("Retrieve an intake work item by work item ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("work_item_id", mcp.Description("UUID of the work item"), mcp.Required()),
	), r.handleRetrieveIntakeWorkItem)

	s.AddTool(mcp.NewTool("update_intake_work_item",
		mcp.WithDescription("Update an intake work item by work item ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("work_item_id", mcp.Description("UUID of the work item"), mcp.Required()),
		mcp.WithObject("data", mcp.Description("Updated intake work item data as a JSON object"), mcp.Required()),
	), r.handleUpdateIntakeWorkItem)

	s.AddTool(mcp.NewTool("delete_intake_work_item",
		mcp.WithDescription("Delete an intake work item by work item ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("work_item_id", mcp.Description("UUID of the work item"), mcp.Required()),
	), r.handleDeleteIntakeWorkItem)
}

func (r *Registry) handleListIntakeWorkItems(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/intake-issues", r.WorkspaceSlug, projectID), nil)
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

func (r *Registry) handleCreateIntakeWorkItem(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	args := getArguments(req)
	body, _ := args["data"].(map[string]any)
	if body == nil {
		return mcp.NewToolResultError("data is required"), nil
	}
	resp, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects/%s/intake-issues", r.WorkspaceSlug, projectID), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleRetrieveIntakeWorkItem(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	workItemID, err := requireID(req, "work_item_id")
	if err != nil {
		return errorResult(err)
	}
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/intake-issues/%s", r.WorkspaceSlug, projectID, workItemID), nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleUpdateIntakeWorkItem(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	workItemID, err := requireID(req, "work_item_id")
	if err != nil {
		return errorResult(err)
	}
	args := getArguments(req)
	body, _ := args["data"].(map[string]any)
	if body == nil {
		return mcp.NewToolResultError("data is required"), nil
	}
	resp, err := r.Client.Patch(ctx, fmt.Sprintf("workspaces/%s/projects/%s/intake-issues/%s", r.WorkspaceSlug, projectID, workItemID), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleDeleteIntakeWorkItem(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	workItemID, err := requireID(req, "work_item_id")
	if err != nil {
		return errorResult(err)
	}
	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/projects/%s/intake-issues/%s", r.WorkspaceSlug, projectID, workItemID)); err != nil {
		return errorResult(err)
	}
	return textResult("Intake work item deleted successfully")
}
