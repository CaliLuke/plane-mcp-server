package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CaliLuke/plane-mcp-server/plane"
	"github.com/CaliLuke/plane-mcp-server/uid"
	"github.com/mark3labs/mcp-go/mcp"
)

// getString extracts a string argument with a default.
func getString(req mcp.CallToolRequest, key, defaultVal string) string {
	return mcp.ParseString(req, key, defaultVal)
}

// requireString extracts a required string argument.
func requireString(req mcp.CallToolRequest, key string) (string, error) {
	s := mcp.ParseString(req, key, "")
	if s == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return s, nil
}

// getBool extracts a boolean argument with a default.
func getBool(req mcp.CallToolRequest, key string, defaultVal bool) bool {
	return mcp.ParseBoolean(req, key, defaultVal)
}

// getFloat extracts a float argument with a default.
func getFloat(req mcp.CallToolRequest, key string, defaultVal float64) float64 {
	return mcp.ParseFloat64(req, key, defaultVal)
}

// getInt extracts an int argument with a default.
func getInt(req mcp.CallToolRequest, key string, defaultVal int) int {
	return mcp.ParseInt(req, key, defaultVal)
}

// getArguments returns the raw arguments map.
func getArguments(req mcp.CallToolRequest) map[string]any {
	return req.Params.Arguments
}

// decodeID decodes a short or full UUID from a tool argument.
func decodeID(req mcp.CallToolRequest, key string) string {
	s := getString(req, key, "")
	if s == "" {
		return ""
	}
	return uid.Decode(s)
}

// requireID decodes and requires a UUID argument.
func requireID(req mcp.CallToolRequest, key string) (string, error) {
	s, err := requireString(req, key)
	if err != nil {
		return "", err
	}
	return uid.Decode(s), nil
}

// getStringSlice extracts a string slice from arguments.
func getStringSlice(req mcp.CallToolRequest, key string) []string {
	args := req.Params.Arguments
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// decodeIDSlice decodes a slice of short/full UUIDs.
func decodeIDSlice(req mcp.CallToolRequest, key string) []string {
	raw := getStringSlice(req, key)
	if raw == nil {
		return nil
	}
	out := make([]string, len(raw))
	for i, s := range raw {
		out[i] = uid.Decode(s)
	}
	return out
}

// optionalString returns a *string from a request argument, or nil if empty.
func optionalString(req mcp.CallToolRequest, key string) *string {
	s := getString(req, key, "")
	if s == "" {
		return nil
	}
	return &s
}

// textResult returns a text tool result.
func textResult(text string) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(text), nil
}

// jsonResult returns a JSON tool result with UUID encoding.
func jsonResult(v any) (*mcp.CallToolResult, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	encoded := uid.EncodeValue(raw)
	result, err := json.MarshalIndent(encoded, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(result)), nil
}

// listResult returns a list of formatted strings as a tool result.
func listResult(items []string) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(strings.Join(items, "\n")), nil
}

// rawJSONResult returns a raw JSON response with UUID encoding.
func rawJSONResult(data json.RawMessage) (*mcp.CallToolResult, error) {
	if data == nil {
		return mcp.NewToolResultText("OK"), nil
	}
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return mcp.NewToolResultText(string(data)), nil
	}
	encoded := uid.EncodeValue(raw)
	result, err := json.MarshalIndent(encoded, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(string(data)), nil
	}
	return mcp.NewToolResultText(string(result)), nil
}

// errorResult wraps an error as a tool error result.
func errorResult(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(err.Error()), nil
}

// resolveState resolves a state name like "Done" to its UUID. Passes through UUIDs unchanged.
func resolveState(ctx context.Context, client *plane.Client, workspaceSlug, projectID, state string) (string, error) {
	if state == "" {
		return "", nil
	}
	if uid.IsUUIDLike(state) {
		return state, nil
	}
	resp, err := client.Get(ctx, fmt.Sprintf("workspaces/%s/projects/%s/states", workspaceSlug, projectID), nil)
	if err != nil {
		return "", fmt.Errorf("fetch states: %w", err)
	}
	var stateResp struct {
		Results []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"results"`
	}
	if err := json.Unmarshal(resp, &stateResp); err != nil {
		return "", fmt.Errorf("parse states: %w", err)
	}
	lower := strings.ToLower(state)
	var validNames []string
	for _, s := range stateResp.Results {
		if strings.ToLower(s.Name) == lower {
			return s.ID, nil
		}
		validNames = append(validNames, s.Name)
	}
	return "", fmt.Errorf("unknown state %q. Valid states: %s", state, strings.Join(validNames, ", "))
}

// buildParams builds a query params map, skipping empty values.
func buildParams(pairs ...string) map[string]string {
	if len(pairs)%2 != 0 {
		panic("buildParams requires even number of arguments")
	}
	params := make(map[string]string)
	for i := 0; i < len(pairs); i += 2 {
		if pairs[i+1] != "" {
			params[pairs[i]] = pairs[i+1]
		}
	}
	if len(params) == 0 {
		return nil
	}
	return params
}

// ensureExpand merges fields into a comma-separated expand string.
func ensureExpand(expand string, fields ...string) string {
	parts := make(map[string]bool)
	for _, p := range strings.Split(expand, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			parts[p] = true
		}
	}
	for _, f := range fields {
		parts[f] = true
	}
	var result []string
	for p := range parts {
		result = append(result, p)
	}
	return strings.Join(result, ",")
}

// stripNulls removes nil values from a map (for JSON request bodies).
func stripNulls(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if v != nil {
			out[k] = v
		}
	}
	return out
}
