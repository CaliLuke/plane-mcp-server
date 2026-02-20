package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (r *Registry) registerModuleTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("list_modules",
		mcp.WithDescription("List all modules in a project."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
	), r.handleListModules)

	s.AddTool(mcp.NewTool("create_module",
		mcp.WithDescription("Create a new module."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("name", mcp.Description("Module name"), mcp.Required()),
		mcp.WithString("description", mcp.Description("Module description")),
		mcp.WithString("start_date", mcp.Description("Module start date (ISO 8601)")),
		mcp.WithString("target_date", mcp.Description("Module target/end date (ISO 8601)")),
		mcp.WithString("status", mcp.Description("Module status"), mcp.Enum("backlog", "planned", "in-progress", "paused", "completed", "cancelled")),
		mcp.WithString("lead", mcp.Description("UUID of the user who leads the module")),
		mcp.WithArray("members", mcp.Description("List of user IDs who are members")),
		mcp.WithString("external_source", mcp.Description("External system source name")),
		mcp.WithString("external_id", mcp.Description("External system identifier")),
	), r.handleCreateModule)

	s.AddTool(mcp.NewTool("retrieve_module",
		mcp.WithDescription("Retrieve a module by ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("module_id", mcp.Description("UUID of the module"), mcp.Required()),
	), r.handleRetrieveModule)

	s.AddTool(mcp.NewTool("update_module",
		mcp.WithDescription("Update a module by ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("module_id", mcp.Description("UUID of the module"), mcp.Required()),
		mcp.WithString("name", mcp.Description("Module name")),
		mcp.WithString("description", mcp.Description("Module description")),
		mcp.WithString("start_date", mcp.Description("Module start date (ISO 8601)")),
		mcp.WithString("target_date", mcp.Description("Module target/end date (ISO 8601)")),
		mcp.WithString("status", mcp.Description("Module status"), mcp.Enum("backlog", "planned", "in-progress", "paused", "completed", "cancelled")),
		mcp.WithString("lead", mcp.Description("UUID of the module lead")),
		mcp.WithArray("members", mcp.Description("List of user IDs who are members")),
		mcp.WithString("external_source", mcp.Description("External system source name")),
		mcp.WithString("external_id", mcp.Description("External system identifier")),
	), r.handleUpdateModule)

	s.AddTool(mcp.NewTool("delete_module",
		mcp.WithDescription("Delete a module by ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("module_id", mcp.Description("UUID of the module"), mcp.Required()),
	), r.handleDeleteModule)

	s.AddTool(mcp.NewTool("list_archived_modules",
		mcp.WithDescription("List archived modules in a project."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
	), r.handleListArchivedModules)

	s.AddTool(mcp.NewTool("add_work_items_to_module",
		mcp.WithDescription("Add work items to a module."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("module_id", mcp.Description("UUID of the module"), mcp.Required()),
		mcp.WithArray("issue_ids", mcp.Description("List of work item IDs to add"), mcp.Required()),
	), r.handleAddWorkItemsToModule)

	s.AddTool(mcp.NewTool("remove_work_item_from_module",
		mcp.WithDescription("Remove a work item from a module."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("module_id", mcp.Description("UUID of the module"), mcp.Required()),
		mcp.WithString("work_item_id", mcp.Description("UUID of the work item to remove"), mcp.Required()),
	), r.handleRemoveWorkItemFromModule)

	s.AddTool(mcp.NewTool("list_module_work_items",
		mcp.WithDescription("List work items in a module."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("module_id", mcp.Description("UUID of the module"), mcp.Required()),
	), r.handleListModuleWorkItems)

	s.AddTool(mcp.NewTool("archive_module",
		mcp.WithDescription("Archive a module."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("module_id", mcp.Description("UUID of the module"), mcp.Required()),
	), r.handleArchiveModule)

	s.AddTool(mcp.NewTool("unarchive_module",
		mcp.WithDescription("Unarchive a module."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("module_id", mcp.Description("UUID of the module"), mcp.Required()),
	), r.handleUnarchiveModule)
}

func (r *Registry) handleListModules(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/modules", r.WorkspaceSlug, projectID), nil)
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

func (r *Registry) handleCreateModule(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	body := stripNulls(map[string]any{
		"name":            getString(req, "name", ""),
		"description":     optionalString(req, "description"),
		"start_date":      optionalString(req, "start_date"),
		"target_date":     optionalString(req, "target_date"),
		"status":          optionalString(req, "status"),
		"lead":            optionalString(req, "lead"),
		"members":         decodeIDSlice(req, "members"),
		"external_source": optionalString(req, "external_source"),
		"external_id":     optionalString(req, "external_id"),
	})
	resp, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects/%s/modules", r.WorkspaceSlug, projectID), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleRetrieveModule(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	moduleID := decodeID(req, "module_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/modules/%s", r.WorkspaceSlug, projectID, moduleID), nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleUpdateModule(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	moduleID := decodeID(req, "module_id")
	body := stripNulls(map[string]any{
		"name":            optionalString(req, "name"),
		"description":     optionalString(req, "description"),
		"start_date":      optionalString(req, "start_date"),
		"target_date":     optionalString(req, "target_date"),
		"status":          optionalString(req, "status"),
		"lead":            optionalString(req, "lead"),
		"members":         decodeIDSlice(req, "members"),
		"external_source": optionalString(req, "external_source"),
		"external_id":     optionalString(req, "external_id"),
	})
	resp, err := r.Client.Patch(ctx, fmt.Sprintf("workspaces/%s/projects/%s/modules/%s", r.WorkspaceSlug, projectID, moduleID), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleDeleteModule(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	moduleID := decodeID(req, "module_id")
	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/projects/%s/modules/%s", r.WorkspaceSlug, projectID, moduleID)); err != nil {
		return errorResult(err)
	}
	return textResult("Module deleted successfully")
}

func (r *Registry) handleListArchivedModules(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/archived-modules", r.WorkspaceSlug, projectID), nil)
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

func (r *Registry) handleAddWorkItemsToModule(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	moduleID := decodeID(req, "module_id")
	issueIDs := decodeIDSlice(req, "issue_ids")
	body := map[string]any{"issues": issueIDs}
	_, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects/%s/modules/%s/module-issues", r.WorkspaceSlug, projectID, moduleID), body)
	if err != nil {
		return errorResult(err)
	}
	return textResult("Work items added to module successfully")
}

func (r *Registry) handleRemoveWorkItemFromModule(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	moduleID := decodeID(req, "module_id")
	workItemID := decodeID(req, "work_item_id")
	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/projects/%s/modules/%s/module-issues/%s", r.WorkspaceSlug, projectID, moduleID, workItemID)); err != nil {
		return errorResult(err)
	}
	return textResult("Work item removed from module successfully")
}

func (r *Registry) handleListModuleWorkItems(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	moduleID := decodeID(req, "module_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/modules/%s/module-issues", r.WorkspaceSlug, projectID, moduleID), nil)
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

func (r *Registry) handleArchiveModule(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	moduleID := decodeID(req, "module_id")
	_, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects/%s/modules/%s/archive", r.WorkspaceSlug, projectID, moduleID), map[string]any{})
	if err != nil {
		return errorResult(err)
	}
	return textResult("Module archived successfully")
}

func (r *Registry) handleUnarchiveModule(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	moduleID := decodeID(req, "module_id")
	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/projects/%s/archived-modules/%s/unarchive", r.WorkspaceSlug, projectID, moduleID)); err != nil {
		return errorResult(err)
	}
	return textResult("Module unarchived successfully")
}
