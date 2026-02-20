package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CaliLuke/plane-mcp-server/formatting"
	"github.com/CaliLuke/plane-mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (r *Registry) registerWorkItemTools(s *server.MCPServer) {
	// list_work_items
	s.AddTool(
		mcp.NewTool("list_work_items",
			mcp.WithDescription("List work items in a project. By default excludes completed/cancelled items."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
			mcp.WithString("cursor", mcp.Description("Pagination cursor for getting next set of results")),
			mcp.WithString("per_page", mcp.Description("Number of results per page (1-100)")),
			mcp.WithString("order_by", mcp.Description("Field to order results by. Prefix with '-' for descending")),
			mcp.WithBoolean("include_completed", mcp.Description("Include items in Done/Cancelled states (default: false)")),
			mcp.WithString("external_id", mcp.Description("External system identifier for filtering")),
			mcp.WithString("external_source", mcp.Description("External system source name for filtering")),
		),
		r.handleListWorkItems,
	)

	// create_work_item
	s.AddTool(
		mcp.NewTool("create_work_item",
			mcp.WithDescription("Create a new work item."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
			mcp.WithString("name", mcp.Description("Work item name"), mcp.Required()),
			mcp.WithArray("assignees", mcp.Description("List of user IDs to assign")),
			mcp.WithArray("labels", mcp.Description("List of label IDs to attach")),
			mcp.WithString("type_id", mcp.Description("UUID of the work item type")),
			mcp.WithNumber("point", mcp.Description("Story point value")),
			mcp.WithString("description_html", mcp.Description("HTML description")),
			mcp.WithString("priority", mcp.Description("Priority level (urgent, high, medium, low, none)"),
				mcp.Enum("urgent", "high", "medium", "low", "none")),
			mcp.WithString("start_date", mcp.Description("Start date (ISO 8601)")),
			mcp.WithString("target_date", mcp.Description("Target/end date (ISO 8601)")),
			mcp.WithString("parent", mcp.Description("UUID of the parent work item")),
			mcp.WithString("state", mcp.Description("State name (e.g. \"Done\", \"In Progress\") or UUID")),
			mcp.WithString("external_source", mcp.Description("External system source name")),
			mcp.WithString("external_id", mcp.Description("External system identifier")),
		),
		r.handleCreateWorkItem,
	)

	// retrieve_work_item
	s.AddTool(
		mcp.NewTool("retrieve_work_item",
			mcp.WithDescription("Retrieve a work item by ID."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
			mcp.WithString("work_item_id", mcp.Description("UUID of the work item"), mcp.Required()),
			mcp.WithString("expand", mcp.Description("Comma-separated fields to expand (e.g., \"assignees,labels,state\")")),
			mcp.WithString("fields", mcp.Description("Comma-separated fields to include")),
		),
		r.handleRetrieveWorkItem,
	)

	// retrieve_work_item_by_identifier
	s.AddTool(
		mcp.NewTool("retrieve_work_item_by_identifier",
			mcp.WithDescription("Retrieve a work item by project identifier and issue sequence number."),
			mcp.WithString("project_identifier", mcp.Description("Project identifier string (e.g., \"MP\")"), mcp.Required()),
			mcp.WithNumber("issue_identifier", mcp.Description("Issue sequence number (e.g., 1, 2, 3)"), mcp.Required()),
			mcp.WithString("expand", mcp.Description("Comma-separated fields to expand")),
			mcp.WithString("fields", mcp.Description("Comma-separated fields to include")),
		),
		r.handleRetrieveWorkItemByIdentifier,
	)

	// update_work_item
	s.AddTool(
		mcp.NewTool("update_work_item",
			mcp.WithDescription("Update a work item by ID."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
			mcp.WithString("work_item_id", mcp.Description("UUID of the work item"), mcp.Required()),
			mcp.WithString("name", mcp.Description("Work item name")),
			mcp.WithArray("assignees", mcp.Description("List of user IDs to assign")),
			mcp.WithArray("labels", mcp.Description("List of label IDs to attach")),
			mcp.WithString("priority", mcp.Description("Priority level (urgent, high, medium, low, none)"),
				mcp.Enum("urgent", "high", "medium", "low", "none")),
			mcp.WithString("start_date", mcp.Description("Start date (ISO 8601)")),
			mcp.WithString("target_date", mcp.Description("Target/end date (ISO 8601)")),
			mcp.WithString("parent", mcp.Description("UUID of the parent work item")),
			mcp.WithString("state", mcp.Description("State name (e.g. \"Done\", \"In Progress\") or UUID")),
			mcp.WithString("description_html", mcp.Description("HTML description")),
			mcp.WithString("external_source", mcp.Description("External system source name")),
			mcp.WithString("external_id", mcp.Description("External system identifier")),
		),
		r.handleUpdateWorkItem,
	)

	// delete_work_item
	s.AddTool(
		mcp.NewTool("delete_work_item",
			mcp.WithDescription("Delete a work item by ID."),
			mcp.WithString("project_id", mcp.Description("UUID of the project"), mcp.Required()),
			mcp.WithString("work_item_id", mcp.Description("UUID of the work item"), mcp.Required()),
		),
		r.handleDeleteWorkItem,
	)

	// search_work_items
	s.AddTool(
		mcp.NewTool("search_work_items",
			mcp.WithDescription("Search work items across a workspace."),
			mcp.WithString("query", mcp.Description("Free-form text search query"), mcp.Required()),
			mcp.WithString("expand", mcp.Description("Comma-separated fields to expand")),
			mcp.WithString("fields", mcp.Description("Comma-separated fields to include")),
			mcp.WithString("order_by", mcp.Description("Field to order results by")),
		),
		r.handleSearchWorkItems,
	)
}

func (r *Registry) handleListWorkItems(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}

	includeCompleted := getBool(req, "include_completed", false)

	// Fetch project to get identifier
	projResp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s", r.WorkspaceSlug, projectID), nil)
	if err != nil {
		return errorResult(err)
	}
	var proj struct {
		Identifier string `json:"identifier"`
	}
	if err := json.Unmarshal(projResp, &proj); err != nil {
		return errorResult(err)
	}

	params := buildParams(
		"cursor", getString(req, "cursor", ""),
		"per_page", getString(req, "per_page", ""),
		"order_by", getString(req, "order_by", ""),
		"external_id", getString(req, "external_id", ""),
		"external_source", getString(req, "external_source", ""),
		"expand", "state",
	)

	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/work-items", r.WorkspaceSlug, projectID), params)
	if err != nil {
		return errorResult(err)
	}

	var data struct {
		Results []json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return errorResult(err)
	}

	closedGroups := map[string]bool{"completed": true, "cancelled": true}
	var lines []string

	for _, raw := range data.Results {
		var item struct {
			ID         string `json:"id"`
			SequenceID int    `json:"sequence_id"`
			Name       string `json:"name"`
			State      any    `json:"state"`
		}
		if err := json.Unmarshal(raw, &item); err != nil {
			continue
		}

		stateName, stateGroup := extractStateInfo(item.State)

		if !includeCompleted && closedGroups[stateGroup] {
			continue
		}

		lines = append(lines, formatting.WorkItem(item.ID, item.SequenceID, item.Name, proj.Identifier, stateName))
	}

	return listResult(lines)
}

func (r *Registry) handleCreateWorkItem(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}

	// Resolve state name to UUID
	stateVal := getString(req, "state", "")
	if stateVal != "" {
		resolved, err := resolveState(ctx, r.Client, r.WorkspaceSlug, projectID, stateVal)
		if err != nil {
			return errorResult(err)
		}
		stateVal = resolved
	}

	body := stripNulls(map[string]any{
		"name":             getString(req, "name", ""),
		"assignees":        decodeIDSlice(req, "assignees"),
		"labels":           decodeIDSlice(req, "labels"),
		"type_id":          optionalString(req, "type_id"),
		"description_html": optionalString(req, "description_html"),
		"priority":         optionalString(req, "priority"),
		"start_date":       optionalString(req, "start_date"),
		"target_date":      optionalString(req, "target_date"),
		"parent":           maybeDecodeID(optionalString(req, "parent")),
		"external_source":  optionalString(req, "external_source"),
		"external_id":      optionalString(req, "external_id"),
	})
	if stateVal != "" {
		body["state"] = stateVal
	}
	if v := getFloat(req, "point", 0); v != 0 {
		body["point"] = int(v)
	}

	resp, err := r.Client.Post(ctx, fmt.Sprintf("workspaces/%s/projects/%s/work-items", r.WorkspaceSlug, projectID), body)
	if err != nil {
		return errorResult(err)
	}

	var item struct {
		SequenceID int `json:"sequence_id"`
	}
	if err := json.Unmarshal(resp, &item); err != nil {
		return errorResult(err)
	}

	// Fetch project identifier
	projResp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s", r.WorkspaceSlug, projectID), nil)
	if err != nil {
		return errorResult(err)
	}
	var proj struct {
		Identifier string `json:"identifier"`
	}
	json.Unmarshal(projResp, &proj)

	identifier := fmt.Sprintf("%s-%d", proj.Identifier, item.SequenceID)
	return jsonResult(map[string]string{"message": fmt.Sprintf("%s created successfully", identifier)})
}

func (r *Registry) handleRetrieveWorkItem(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	workItemID, err := requireID(req, "work_item_id")
	if err != nil {
		return errorResult(err)
	}

	expand := ensureExpand(getString(req, "expand", ""), "state", "assignees")
	params := buildParams(
		"expand", expand,
		"fields", getString(req, "fields", ""),
	)

	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/work-items/%s", r.WorkspaceSlug, projectID, workItemID), params)
	if err != nil {
		return errorResult(err)
	}

	full := parseWorkItemFull(resp)
	return jsonResult(full)
}

func (r *Registry) handleRetrieveWorkItemByIdentifier(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projIdent := getString(req, "project_identifier", "")
	issueIdent := getInt(req, "issue_identifier", 0)

	expand := ensureExpand(getString(req, "expand", ""), "state", "assignees")
	params := buildParams(
		"expand", expand,
		"fields", getString(req, "fields", ""),
	)

	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/work-items/%s-%d", r.WorkspaceSlug, projIdent, issueIdent), params)
	if err != nil {
		return errorResult(err)
	}

	full := parseWorkItemFull(resp)
	return jsonResult(full)
}

func (r *Registry) handleUpdateWorkItem(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	workItemID, err := requireID(req, "work_item_id")
	if err != nil {
		return errorResult(err)
	}

	stateVal := getString(req, "state", "")
	if stateVal != "" {
		resolved, err := resolveState(ctx, r.Client, r.WorkspaceSlug, projectID, stateVal)
		if err != nil {
			return errorResult(err)
		}
		stateVal = resolved
	}

	body := stripNulls(map[string]any{
		"name":             optionalString(req, "name"),
		"assignees":        decodeIDSlice(req, "assignees"),
		"labels":           decodeIDSlice(req, "labels"),
		"description_html": optionalString(req, "description_html"),
		"priority":         optionalString(req, "priority"),
		"start_date":       optionalString(req, "start_date"),
		"target_date":      optionalString(req, "target_date"),
		"parent":           maybeDecodeID(optionalString(req, "parent")),
		"external_source":  optionalString(req, "external_source"),
		"external_id":      optionalString(req, "external_id"),
	})
	if stateVal != "" {
		body["state"] = stateVal
	}

	resp, err := r.Client.Patch(ctx, fmt.Sprintf("workspaces/%s/projects/%s/work-items/%s", r.WorkspaceSlug, projectID, workItemID), body)
	if err != nil {
		return errorResult(err)
	}

	var item struct {
		SequenceID int `json:"sequence_id"`
	}
	json.Unmarshal(resp, &item)

	projResp, _ := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s", r.WorkspaceSlug, projectID), nil)
	var proj struct {
		Identifier string `json:"identifier"`
	}
	json.Unmarshal(projResp, &proj)

	identifier := fmt.Sprintf("%s-%d", proj.Identifier, item.SequenceID)
	return jsonResult(map[string]string{"message": fmt.Sprintf("%s updated successfully", identifier)})
}

func (r *Registry) handleDeleteWorkItem(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := requireID(req, "project_id")
	if err != nil {
		return errorResult(err)
	}
	workItemID, err := requireID(req, "work_item_id")
	if err != nil {
		return errorResult(err)
	}

	if err := r.Client.Delete(ctx, fmt.Sprintf("workspaces/%s/projects/%s/work-items/%s", r.WorkspaceSlug, projectID, workItemID)); err != nil {
		return errorResult(err)
	}
	return textResult("Work item deleted successfully")
}

func (r *Registry) handleSearchWorkItems(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := getString(req, "query", "")
	params := buildParams(
		"search", query,
		"expand", getString(req, "expand", ""),
		"fields", getString(req, "fields", ""),
		"order_by", getString(req, "order_by", ""),
	)

	resp, err := r.Client.Get(ctx, fmt.Sprintf("workspaces/%s/work-items/search", r.WorkspaceSlug), params)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}

// extractStateInfo extracts the state name and group from a state field,
// which can be either a string (UUID) or an expanded object.
func extractStateInfo(state any) (name, group string) {
	switch s := state.(type) {
	case string:
		return "", ""
	case map[string]any:
		if n, ok := s["name"].(string); ok {
			name = n
		}
		if g, ok := s["group"].(string); ok {
			group = g
		}
		return name, group
	}
	return "", ""
}

// parseWorkItemFull parses a raw API response into a WorkItemFull.
func parseWorkItemFull(data json.RawMessage) *models.WorkItemFull {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return &models.WorkItemFull{}
	}

	full := &models.WorkItemFull{
		ID:         jsonStr(raw, "id"),
		SequenceID: jsonInt(raw, "sequence_id"),
		Name:       jsonStr(raw, "name"),
		Description: formatting.StripHTML(jsonStr(raw, "description_html")),
		Priority:   jsonStr(raw, "priority"),
		Parent:     jsonStr(raw, "parent"),
		StartDate:  jsonStr(raw, "start_date"),
		TargetDate: jsonStr(raw, "target_date"),
	}

	// Extract state (expanded object or string)
	if stateVal, ok := raw["state"]; ok {
		switch s := stateVal.(type) {
		case string:
			full.State = s
		case map[string]any:
			if n, ok := s["name"].(string); ok {
				full.State = n
			} else if id, ok := s["id"].(string); ok {
				full.State = id
			}
		}
	}

	// Extract assignees (expanded objects or string IDs)
	if assigneesVal, ok := raw["assignees"]; ok {
		if arr, ok := assigneesVal.([]any); ok {
			for _, a := range arr {
				switch av := a.(type) {
				case string:
					full.Assignees = append(full.Assignees, models.AssigneeSummary{ID: av})
				case map[string]any:
					as := models.AssigneeSummary{}
					if dn, ok := av["display_name"].(string); ok && dn != "" {
						as.DisplayName = dn
					} else if id, ok := av["id"].(string); ok {
						as.ID = id
					}
					full.Assignees = append(full.Assignees, as)
				}
			}
		}
	}

	// Extract labels (expanded objects or string IDs)
	if labelsVal, ok := raw["labels"]; ok {
		if arr, ok := labelsVal.([]any); ok {
			for _, l := range arr {
				switch lv := l.(type) {
				case string:
					full.Labels = append(full.Labels, models.LabelSummary{ID: lv, Name: lv})
				case map[string]any:
					ls := models.LabelSummary{
						ID:    jsonStrMap(lv, "id"),
						Name:  jsonStrMap(lv, "name"),
						Color: jsonStrMap(lv, "color"),
					}
					full.Labels = append(full.Labels, ls)
				}
			}
		}
	}

	return full
}

func maybeDecodeID(s *string) *string {
	if s == nil {
		return nil
	}
	decoded := decodeRawID(*s)
	return &decoded
}

func decodeRawID(s string) string {
	if s == "" {
		return ""
	}
	return strings.TrimSpace(s)
}

func jsonStr(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func jsonStrMap(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func jsonInt(m map[string]any, key string) int {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return 0
}
