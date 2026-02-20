package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (r *Registry) registerWorkItemPropertyTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("list_work_item_properties",
		mcp.WithDescription("List work item properties for a work item type."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("type_id", mcp.Description("UUID of the work item type"), mcp.Required()),
	), r.handleListWorkItemProperties)

	s.AddTool(mcp.NewTool("create_work_item_property",
		mcp.WithDescription("Create a new work item property."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("type_id", mcp.Description("UUID of the work item type"), mcp.Required()),
		mcp.WithString("display_name", mcp.Description("Display name for the property"), mcp.Required()),
		mcp.WithString("property_type", mcp.Description("Type of property"),
			mcp.Enum("TEXT", "DATETIME", "DECIMAL", "BOOLEAN", "OPTION", "RELATION", "URL", "EMAIL", "FILE"), mcp.Required()),
		mcp.WithString("relation_type", mcp.Description("Relation type (ISSUE, USER) - required for RELATION properties"),
			mcp.Enum("ISSUE", "USER")),
		mcp.WithString("description", mcp.Description("Property description")),
		mcp.WithBoolean("is_required", mcp.Description("Whether the property is required")),
		mcp.WithBoolean("is_active", mcp.Description("Whether the property is active")),
		mcp.WithBoolean("is_multi", mcp.Description("Whether the property supports multiple values")),
		mcp.WithObject("settings", mcp.Description("Settings dict for TEXT/DATETIME properties")),
		mcp.WithString("external_source", mcp.Description("External system source name")),
		mcp.WithString("external_id", mcp.Description("External system identifier")),
	), r.handleCreateWorkItemProperty)

	s.AddTool(mcp.NewTool("retrieve_work_item_property",
		mcp.WithDescription("Retrieve a work item property by ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("type_id", mcp.Description("UUID of the work item type"), mcp.Required()),
		mcp.WithString("work_item_property_id", mcp.Description("UUID of the property"), mcp.Required()),
	), r.handleRetrieveWorkItemProperty)

	s.AddTool(mcp.NewTool("update_work_item_property",
		mcp.WithDescription("Update a work item property by ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("type_id", mcp.Description("UUID of the work item type"), mcp.Required()),
		mcp.WithString("work_item_property_id", mcp.Description("UUID of the property"), mcp.Required()),
		mcp.WithString("display_name", mcp.Description("Display name for the property")),
		mcp.WithString("property_type", mcp.Description("Type of property"),
			mcp.Enum("TEXT", "DATETIME", "DECIMAL", "BOOLEAN", "OPTION", "RELATION", "URL", "EMAIL", "FILE")),
		mcp.WithString("relation_type", mcp.Description("Relation type (ISSUE, USER)"), mcp.Enum("ISSUE", "USER")),
		mcp.WithString("description", mcp.Description("Property description")),
		mcp.WithBoolean("is_required", mcp.Description("Whether the property is required")),
		mcp.WithBoolean("is_active", mcp.Description("Whether the property is active")),
		mcp.WithBoolean("is_multi", mcp.Description("Whether the property supports multiple values")),
		mcp.WithObject("settings", mcp.Description("Settings dict for TEXT/DATETIME properties")),
		mcp.WithString("external_source", mcp.Description("External system source name")),
		mcp.WithString("external_id", mcp.Description("External system identifier")),
	), r.handleUpdateWorkItemProperty)

	s.AddTool(mcp.NewTool("delete_work_item_property",
		mcp.WithDescription("Delete a work item property by ID."),
		mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
		mcp.WithString("type_id", mcp.Description("UUID of the work item type"), mcp.Required()),
		mcp.WithString("work_item_property_id", mcp.Description("UUID of the property"), mcp.Required()),
	), r.handleDeleteWorkItemProperty)
}

func (r *Registry) workItemPropertyBasePath(projectID, typeID string) string {
	return fmt.Sprintf("workspaces/%s/projects/%s/work-item-types/%s/work-item-properties", r.WorkspaceSlug, projectID, typeID)
}

func (r *Registry) handleListWorkItemProperties(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	typeID, err := requireID(req, "type_id")
	if err != nil {
		return errorResult(err)
	}
	resp, err := r.Client.Get(ctx, r.workItemPropertyBasePath(projectID, typeID), nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleCreateWorkItemProperty(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	typeID, err := requireID(req, "type_id")
	if err != nil {
		return errorResult(err)
	}

	args := getArguments(req)
	body := stripNulls(map[string]any{
		"display_name":    getString(req, "display_name", ""),
		"property_type":   getString(req, "property_type", ""),
		"relation_type":   optionalString(req, "relation_type"),
		"description":     optionalString(req, "description"),
		"external_source": optionalString(req, "external_source"),
		"external_id":     optionalString(req, "external_id"),
	})
	if v, ok := args["is_required"]; ok {
		body["is_required"] = v
	}
	if v, ok := args["is_active"]; ok {
		body["is_active"] = v
	}
	if v, ok := args["is_multi"]; ok {
		body["is_multi"] = v
	}
	if v, ok := args["settings"]; ok {
		body["settings"] = v
	}

	resp, err := r.Client.Post(ctx, r.workItemPropertyBasePath(projectID, typeID), body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleRetrieveWorkItemProperty(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	typeID, err := requireID(req, "type_id")
	if err != nil {
		return errorResult(err)
	}
	propID, err := requireID(req, "work_item_property_id")
	if err != nil {
		return errorResult(err)
	}
	resp, err := r.Client.Get(ctx, r.workItemPropertyBasePath(projectID, typeID)+"/"+propID, nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleUpdateWorkItemProperty(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	typeID, err := requireID(req, "type_id")
	if err != nil {
		return errorResult(err)
	}
	propID, err := requireID(req, "work_item_property_id")
	if err != nil {
		return errorResult(err)
	}

	args := getArguments(req)
	body := stripNulls(map[string]any{
		"display_name":    optionalString(req, "display_name"),
		"property_type":   optionalString(req, "property_type"),
		"relation_type":   optionalString(req, "relation_type"),
		"description":     optionalString(req, "description"),
		"external_source": optionalString(req, "external_source"),
		"external_id":     optionalString(req, "external_id"),
	})
	if v, ok := args["is_required"]; ok {
		body["is_required"] = v
	}
	if v, ok := args["is_active"]; ok {
		body["is_active"] = v
	}
	if v, ok := args["is_multi"]; ok {
		body["is_multi"] = v
	}
	if v, ok := args["settings"]; ok {
		body["settings"] = v
	}

	resp, err := r.Client.Patch(ctx, r.workItemPropertyBasePath(projectID, typeID)+"/"+propID, body)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

func (r *Registry) handleDeleteWorkItemProperty(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	typeID, err := requireID(req, "type_id")
	if err != nil {
		return errorResult(err)
	}
	propID, err := requireID(req, "work_item_property_id")
	if err != nil {
		return errorResult(err)
	}
	if err := r.Client.Delete(ctx, r.workItemPropertyBasePath(projectID, typeID)+"/"+propID); err != nil {
		return errorResult(err)
	}
	return textResult("Work item property deleted successfully")
}
