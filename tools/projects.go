package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CaliLuke/plane-mcp-server/formatting"
	"github.com/CaliLuke/plane-mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (r *Registry) registerProjectTools(s *server.MCPServer) {
	// list_projects
	s.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("List all projects in a workspace."),
			mcp.WithString("cursor", mcp.Description("Pagination cursor for getting next set of results")),
			mcp.WithString("per_page", mcp.Description("Number of results per page (1-100)")),
			mcp.WithString("order_by", mcp.Description("Field to order results by. Prefix with '-' for descending")),
		),
		r.handleListProjects,
	)

	// create_project
	s.AddTool(
		mcp.NewTool("create_project",
			mcp.WithDescription("Create a new project."),
			mcp.WithString("name", mcp.Description("Project name"), mcp.Required()),
			mcp.WithString("identifier", mcp.Description("Project identifier (e.g., \"MP\")"), mcp.Required()),
			mcp.WithString("description", mcp.Description("Project description")),
			mcp.WithString("project_lead", mcp.Description("UUID of the project lead user")),
			mcp.WithString("default_assignee", mcp.Description("UUID of the default assignee user")),
			mcp.WithString("emoji", mcp.Description("Emoji for the project")),
			mcp.WithString("timezone", mcp.Description("Project timezone")),
			mcp.WithString("external_source", mcp.Description("External system source name")),
			mcp.WithString("external_id", mcp.Description("External system identifier")),
		),
		r.handleCreateProject,
	)

	// retrieve_project
	s.AddTool(
		mcp.NewTool("retrieve_project",
			mcp.WithDescription("Retrieve a project by ID."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		),
		r.handleRetrieveProject,
	)

	// update_project
	s.AddTool(
		mcp.NewTool("update_project",
			mcp.WithDescription("Update a project by ID."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
			mcp.WithString("name", mcp.Description("Project name")),
			mcp.WithString("description", mcp.Description("Project description")),
			mcp.WithString("identifier", mcp.Description("Project identifier")),
			mcp.WithString("project_lead", mcp.Description("UUID of the project lead user")),
			mcp.WithString("default_assignee", mcp.Description("UUID of the default assignee user")),
			mcp.WithString("emoji", mcp.Description("Emoji for the project")),
			mcp.WithString("timezone", mcp.Description("Project timezone")),
			mcp.WithString("external_source", mcp.Description("External system source name")),
			mcp.WithString("external_id", mcp.Description("External system identifier")),
		),
		r.handleUpdateProject,
	)

	// delete_project
	s.AddTool(
		mcp.NewTool("delete_project",
			mcp.WithDescription("Delete a project by ID."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		),
		r.handleDeleteProject,
	)

	// get_project_members
	s.AddTool(
		mcp.NewTool("get_project_members",
			mcp.WithDescription("Get all members of a project."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		),
		r.handleGetProjectMembers,
	)

	// get_project_features
	s.AddTool(
		mcp.NewTool("get_project_features",
			mcp.WithDescription("Get features of a project."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		),
		r.handleGetProjectFeatures,
	)

	// update_project_features
	s.AddTool(
		mcp.NewTool("update_project_features",
			mcp.WithDescription("Update features of a project."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
			mcp.WithBoolean("epics", mcp.Description("Enable/disable epics feature")),
			mcp.WithBoolean("modules", mcp.Description("Enable/disable modules feature")),
			mcp.WithBoolean("cycles", mcp.Description("Enable/disable cycles feature")),
			mcp.WithBoolean("views", mcp.Description("Enable/disable views feature")),
			mcp.WithBoolean("pages", mcp.Description("Enable/disable pages feature")),
			mcp.WithBoolean("intakes", mcp.Description("Enable/disable intakes feature")),
			mcp.WithBoolean("work_item_types", mcp.Description("Enable/disable work item types feature")),
		),
		r.handleUpdateProjectFeatures,
	)

	// get_project_worklog_summary
	s.AddTool(
		mcp.NewTool("get_project_worklog_summary",
			mcp.WithDescription("Get work log summary for a project."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		),
		r.handleGetProjectWorklogSummary,
	)
}

func (r *Registry) handleListProjects(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := buildParams(
		"cursor", getString(req, "cursor", ""),
		"per_page", getString(req, "per_page", ""),
		"order_by", getString(req, "order_by", ""),
	)

	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects", r.WorkspaceSlug), params)
	if err != nil {
		return errorResult(err)
	}

	var data struct {
		Results []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Identifier string `json:"identifier"`
		} `json:"results"`
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return errorResult(err)
	}

	lines := make([]string, len(data.Results))
	for i, p := range data.Results {
		lines[i] = formatting.Project(p.ID, p.Name, p.Identifier)
	}
	return listResult(lines)
}

func (r *Registry) handleCreateProject(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	body := stripNulls(map[string]any{
		"name":             getString(req, "name", ""),
		"identifier":       getString(req, "identifier", ""),
		"description":      optionalString(req, "description"),
		"project_lead":     optionalString(req, "project_lead"),
		"default_assignee": optionalString(req, "default_assignee"),
		"emoji":            optionalString(req, "emoji"),
		"timezone":         optionalString(req, "timezone"),
		"external_source":  optionalString(req, "external_source"),
		"external_id":      optionalString(req, "external_id"),
	})

	resp, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects", r.WorkspaceSlug), body)
	if err != nil {
		return errorResult(err)
	}

	var p struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
	}
	if err := json.Unmarshal(resp, &p); err != nil {
		return errorResult(err)
	}

	summary := models.ProjectSummary{ID: p.ID, Name: p.Name, Identifier: p.Identifier}
	return jsonResult(summary)
}

func (r *Registry) handleRetrieveProject(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}

	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s", r.WorkspaceSlug, projectID), nil)
	if err != nil {
		return errorResult(err)
	}

	var p struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
	}
	if err := json.Unmarshal(resp, &p); err != nil {
		return errorResult(err)
	}

	summary := models.ProjectSummary{ID: p.ID, Name: p.Name, Identifier: p.Identifier}
	return jsonResult(summary)
}

func (r *Registry) handleUpdateProject(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}

	body := stripNulls(map[string]any{
		"name":             optionalString(req, "name"),
		"description":      optionalString(req, "description"),
		"identifier":       optionalString(req, "identifier"),
		"project_lead":     optionalString(req, "project_lead"),
		"default_assignee": optionalString(req, "default_assignee"),
		"emoji":            optionalString(req, "emoji"),
		"timezone":         optionalString(req, "timezone"),
		"external_source":  optionalString(req, "external_source"),
		"external_id":      optionalString(req, "external_id"),
	})

	resp, err := r.Client.Patch(ctx, fmt.Sprintf("workspaces/%s/projects/%s", r.WorkspaceSlug, projectID), body)
	if err != nil {
		return errorResult(err)
	}

	var p struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
	}
	if err := json.Unmarshal(resp, &p); err != nil {
		return errorResult(err)
	}

	summary := models.ProjectSummary{ID: p.ID, Name: p.Name, Identifier: p.Identifier}
	return jsonResult(summary)
}

func (r *Registry) handleDeleteProject(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}

	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/projects/%s", r.WorkspaceSlug, projectID)); err != nil {
		return errorResult(err)
	}
	return textResult("Project deleted successfully")
}

func (r *Registry) handleGetProjectMembers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/members", r.WorkspaceSlug, projectID), nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleGetProjectFeatures(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/features", r.WorkspaceSlug, projectID), nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleUpdateProjectFeatures(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	args := getArguments(req)
	body := make(map[string]any)
	for _, key := range []string{"epics", "modules", "cycles", "views", "pages", "intakes", "work_item_types"} {
		if v, ok := args[key]; ok {
			body[key] = v
		}
	}

	resp, err := r.Client.Patch(ctx, fmt.Sprintf("workspaces/%s/projects/%s/features", r.WorkspaceSlug, projectID), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleGetProjectWorklogSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := decodeID(req, "project_id")
	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/total-worklogs", r.WorkspaceSlug, projectID), nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}
