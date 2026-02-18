"""State-related tools for Plane MCP Server."""

from fastmcp import FastMCP
from plane.models.states import PaginatedStateResponse

from plane_mcp.client import get_plane_client_context
from plane_mcp.uid import ShortUUID


def register_state_tools(mcp: FastMCP) -> None:
    """Register all state-related tools with the MCP server."""

    @mcp.tool()
    def list_states(
        project_id: ShortUUID,
    ) -> list[str]:
        """
        List valid state names for a project.

        Use this to discover which state names can be passed to
        create_work_item / update_work_item.

        Args:
            project_id: UUID of the project

        Returns:
            List of state names grouped by workflow group, e.g. "Done (completed)".
        """
        client, workspace_slug = get_plane_client_context()

        response: PaginatedStateResponse = client.states.list(
            workspace_slug=workspace_slug,
            project_id=project_id,
        )

        def _fmt(s) -> str:
            group = s.group.value if hasattr(s.group, "value") else str(s.group or "")
            default = " *" if s.default else ""
            return f"{s.name} ({group}){default}"

        return [_fmt(s) for s in response.results]
