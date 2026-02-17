"""State-related tools for Plane MCP Server."""

from fastmcp import FastMCP
from plane.models.states import PaginatedStateResponse

from plane_mcp.client import get_plane_client_context
from plane_mcp.uid import ShortUUID
from plane_mcp import formatting


def register_state_tools(mcp: FastMCP) -> None:
    """Register all state-related tools with the MCP server."""

    @mcp.tool()
    def list_states(
        project_id: ShortUUID,
    ) -> list[str]:
        """
        List all states for a project.

        Use this to get state UUIDs needed when creating or updating work items.

        Args:
            project_id: UUID of the project

        Returns:
            List of StateSummary objects containing id, name, group, color, default, sequence.
        """
        client, workspace_slug = get_plane_client_context()

        response: PaginatedStateResponse = client.states.list(
            workspace_slug=workspace_slug,
            project_id=project_id,
        )

        return [formatting.state(s) for s in response.results]
