package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (r *Registry) registerUserTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_me",
			mcp.WithDescription("Get current user information."),
		),
		r.handleGetMe,
	)
}

func (r *Registry) handleGetMe(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resp, err := r.Client.Get(ctx, "users/me", nil)
	if err != nil {
		return errorResult(err)
	}
	return rawJSONResult(resp)
}
