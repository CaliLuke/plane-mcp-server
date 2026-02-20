package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (r *Registry) registerCycleTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("list_cycles",
		mcp.WithDescription("List all cycles in a project."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
	), r.handleListCycles)

	s.AddTool(mcp.NewTool("create_cycle",
		mcp.WithDescription("Create a new cycle."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("name", mcp.Description("Cycle name"), mcp.Required()),
		mcp.WithString("owned_by", mcp.Description("UUID of the user who owns the cycle"), mcp.Required()),
		mcp.WithString("description", mcp.Description("Cycle description")),
		mcp.WithString("start_date", mcp.Description("Cycle start date (ISO 8601)")),
		mcp.WithString("end_date", mcp.Description("Cycle end date (ISO 8601)")),
		mcp.WithString("timezone", mcp.Description("Cycle timezone")),
		mcp.WithString("external_source", mcp.Description("External system source name")),
		mcp.WithString("external_id", mcp.Description("External system identifier")),
	), r.handleCreateCycle)

	s.AddTool(mcp.NewTool("retrieve_cycle",
		mcp.WithDescription("Retrieve a cycle by ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("cycle_id", mcp.Description("UUID of the cycle"), mcp.Required()),
	), r.handleRetrieveCycle)

	s.AddTool(mcp.NewTool("update_cycle",
		mcp.WithDescription("Update a cycle by ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("cycle_id", mcp.Description("UUID of the cycle"), mcp.Required()),
		mcp.WithString("name", mcp.Description("Cycle name")),
		mcp.WithString("description", mcp.Description("Cycle description")),
		mcp.WithString("start_date", mcp.Description("Cycle start date (ISO 8601)")),
		mcp.WithString("end_date", mcp.Description("Cycle end date (ISO 8601)")),
		mcp.WithString("owned_by", mcp.Description("UUID of the cycle owner")),
		mcp.WithString("timezone", mcp.Description("Cycle timezone")),
		mcp.WithString("external_source", mcp.Description("External system source name")),
		mcp.WithString("external_id", mcp.Description("External system identifier")),
	), r.handleUpdateCycle)

	s.AddTool(mcp.NewTool("delete_cycle",
		mcp.WithDescription("Delete a cycle by ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("cycle_id", mcp.Description("UUID of the cycle"), mcp.Required()),
	), r.handleDeleteCycle)

	s.AddTool(mcp.NewTool("list_archived_cycles",
		mcp.WithDescription("List archived cycles in a project."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
	), r.handleListArchivedCycles)

	s.AddTool(mcp.NewTool("add_work_items_to_cycle",
		mcp.WithDescription("Add work items to a cycle."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("cycle_id", mcp.Description("UUID of the cycle"), mcp.Required()),
		mcp.WithArray("issue_ids", mcp.Description("List of work item IDs to add"), mcp.Required()),
	), r.handleAddWorkItemsToCycle)

	s.AddTool(mcp.NewTool("remove_work_item_from_cycle",
		mcp.WithDescription("Remove a work item from a cycle."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("cycle_id", mcp.Description("UUID of the cycle"), mcp.Required()),
		mcp.WithString("work_item_id", mcp.Description("UUID of the work item to remove"), mcp.Required()),
	), r.handleRemoveWorkItemFromCycle)

	s.AddTool(mcp.NewTool("list_cycle_work_items",
		mcp.WithDescription("List work items in a cycle."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("cycle_id", mcp.Description("UUID of the cycle"), mcp.Required()),
	), r.handleListCycleWorkItems)

	s.AddTool(mcp.NewTool("transfer_cycle_work_items",
		mcp.WithDescription("Transfer work items from one cycle to another."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("cycle_id", mcp.Description("UUID of the source cycle"), mcp.Required()),
		mcp.WithString("new_cycle_id", mcp.Description("UUID of the target cycle"), mcp.Required()),
	), r.handleTransferCycleWorkItems)

	s.AddTool(mcp.NewTool("archive_cycle",
		mcp.WithDescription("Archive a cycle."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("cycle_id", mcp.Description("UUID of the cycle"), mcp.Required()),
	), r.handleArchiveCycle)

	s.AddTool(mcp.NewTool("unarchive_cycle",
		mcp.WithDescription("Unarchive a cycle."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("cycle_id", mcp.Description("UUID of the cycle"), mcp.Required()),
	), r.handleUnarchiveCycle)
}

func (r *Registry) handleListCycles(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles", r.WorkspaceSlug, projectID), nil)
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

func (r *Registry) handleCreateCycle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	body := stripNulls(map[string]any{
		"name":            getString(req, "name", ""),
		"owned_by":        getString(req, "owned_by", ""),
		"project_id":      projectID,
		"description":     optionalString(req, "description"),
		"start_date":      optionalString(req, "start_date"),
		"end_date":        optionalString(req, "end_date"),
		"timezone":        optionalString(req, "timezone"),
		"external_source": optionalString(req, "external_source"),
		"external_id":     optionalString(req, "external_id"),
	})
	resp, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles", r.WorkspaceSlug, projectID), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleRetrieveCycle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	cycleID := decodeID(req, "cycle_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles/%s", r.WorkspaceSlug, projectID, cycleID), nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleUpdateCycle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	cycleID := decodeID(req, "cycle_id")
	body := stripNulls(map[string]any{
		"name":            optionalString(req, "name"),
		"description":     optionalString(req, "description"),
		"start_date":      optionalString(req, "start_date"),
		"end_date":        optionalString(req, "end_date"),
		"owned_by":        optionalString(req, "owned_by"),
		"timezone":        optionalString(req, "timezone"),
		"external_source": optionalString(req, "external_source"),
		"external_id":     optionalString(req, "external_id"),
	})
	resp, err := r.Client.Patch(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles/%s", r.WorkspaceSlug, projectID, cycleID), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleDeleteCycle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	cycleID := decodeID(req, "cycle_id")
	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles/%s", r.WorkspaceSlug, projectID, cycleID)); err != nil {
		return errorResult(err)
	}
	return textResult("Cycle deleted successfully")
}

func (r *Registry) handleListArchivedCycles(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/archived-cycles", r.WorkspaceSlug, projectID), nil)
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

func (r *Registry) handleAddWorkItemsToCycle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	cycleID := decodeID(req, "cycle_id")
	issueIDs := decodeIDSlice(req, "issue_ids")
	body := map[string]any{"issues": issueIDs}
	_, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles/%s/cycle-issues", r.WorkspaceSlug, projectID, cycleID), body)
	if err != nil {
		return errorResult(err)
	}
	return textResult("Work items added to cycle successfully")
}

func (r *Registry) handleRemoveWorkItemFromCycle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	cycleID := decodeID(req, "cycle_id")
	workItemID := decodeID(req, "work_item_id")
	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles/%s/cycle-issues/%s", r.WorkspaceSlug, projectID, cycleID, workItemID)); err != nil {
		return errorResult(err)
	}
	return textResult("Work item removed from cycle successfully")
}

func (r *Registry) handleListCycleWorkItems(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	cycleID := decodeID(req, "cycle_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles/%s/cycle-issues", r.WorkspaceSlug, projectID, cycleID), nil)
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

func (r *Registry) handleTransferCycleWorkItems(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	cycleID := decodeID(req, "cycle_id")
	newCycleID := decodeID(req, "new_cycle_id")
	body := map[string]any{"new_cycle_id": newCycleID}
	_, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles/%s/transfer-issues", r.WorkspaceSlug, projectID, cycleID), body)
	if err != nil {
		return errorResult(err)
	}
	return textResult("Work items transferred successfully")
}

func (r *Registry) handleArchiveCycle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	cycleID := decodeID(req, "cycle_id")
	_, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects/%s/cycles/%s/archive", r.WorkspaceSlug, projectID, cycleID), map[string]any{})
	if err != nil {
		return errorResult(err)
	}
	return textResult("Cycle archived successfully")
}

func (r *Registry) handleUnarchiveCycle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	cycleID := decodeID(req, "cycle_id")
	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/projects/%s/archived-cycles/%s/unarchive", r.WorkspaceSlug, projectID, cycleID)); err != nil {
		return errorResult(err)
	}
	return textResult("Cycle unarchived successfully")
}
