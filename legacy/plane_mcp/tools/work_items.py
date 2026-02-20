"""Work item-related tools for Plane MCP Server."""

from fastmcp import FastMCP
from plane.models.enums import PriorityEnum
from plane.models.query_params import RetrieveQueryParams, WorkItemQueryParams
from plane.models.work_items import (
    CreateWorkItem,
    PaginatedWorkItemResponse,
    UpdateWorkItem,
    WorkItemDetail,
    WorkItemSearch,
)

from plane_mcp.client import get_plane_client_context
from plane_mcp.uid import ShortUUID, _UUID_RE
from plane_mcp import formatting
from plane_mcp.models import (
    AssigneeSummary,
    LabelSummary,
    WorkItemFull,
    WorkItemSummary,
    strip_html,
)


def _is_uuid_like(s: str) -> bool:
    """Return True if s looks like a full UUID or a shortUUID (base57, ~22 chars)."""
    if _UUID_RE.match(s):
        return True
    # shortUUIDs are 22 alphanumeric chars (base57)
    if len(s) >= 20 and s.isalnum():
        return True
    return False


def _resolve_state(client, workspace_slug: str, project_id: str, state: str | None) -> str | None:
    """Resolve a state name like 'Done' to its UUID. Pass through UUIDs unchanged."""
    if state is None:
        return None
    if _is_uuid_like(state):
        return state
    # Fetch states and match by name (case-insensitive)
    from plane.models.states import PaginatedStateResponse
    response: PaginatedStateResponse = client.states.list(
        workspace_slug=workspace_slug, project_id=project_id,
    )
    lower = state.lower()
    for s in response.results:
        if s.name and s.name.lower() == lower:
            return s.id
    valid = [s.name for s in response.results if s.name]
    raise ValueError(f"Unknown state {state!r}. Valid states: {', '.join(valid)}")


def _ensure_expand(expand: str | None, *fields: str) -> str:
    """Merge *fields* into a comma-separated expand string, avoiding duplicates."""
    parts = [p.strip() for p in (expand or "").split(",") if p.strip()]
    for f in fields:
        if f not in parts:
            parts.append(f)
    return ",".join(parts)


def register_work_item_tools(mcp: FastMCP) -> None:
    """Register all work item-related tools with the MCP server."""

    @mcp.tool()
    def list_work_items(
        project_id: ShortUUID,
        cursor: str | None = None,
        per_page: int | None = None,
        order_by: str | None = None,
        include_completed: bool = False,
        external_id: str | None = None,
        external_source: str | None = None,
    ) -> list[str]:
        """
        List work items in a project. By default excludes completed/cancelled items.

        Args:
            project_id: UUID of the project
            cursor: Pagination cursor for getting next set of results
            per_page: Number of results per page (1-100)
            order_by: Field to order results by. Prefix with '-' for descending order
            include_completed: Include items in Done/Cancelled states (default: false)
            external_id: External system identifier for filtering or lookup
            external_source: External system source name for filtering or lookup

        Returns:
            List of work item summaries (use retrieve_work_item for full detail).
        """
        client, workspace_slug = get_plane_client_context()

        project = client.projects.retrieve(workspace_slug=workspace_slug, project_id=project_id)
        project_identifier = project.identifier or ""

        params = WorkItemQueryParams(
            cursor=cursor,
            per_page=per_page,
            order_by=order_by,
            external_id=external_id,
            external_source=external_source,
            expand="state",
        )

        response: PaginatedWorkItemResponse = client.work_items.list(
            workspace_slug=workspace_slug,
            project_id=project_id,
            params=params,
        )

        _CLOSED_GROUPS = {"completed", "cancelled"}

        def _state_name(item) -> str | None:
            s = item.state
            if s is None:
                return None
            if isinstance(s, str):
                return None
            return getattr(s, "name", None)

        def _state_group(item) -> str | None:
            s = item.state
            if s is None or isinstance(s, str):
                return None
            return getattr(s, "group", None)

        return [
            formatting.work_item(item, project_identifier, _state_name(item))
            for item in response.results
            if include_completed or _state_group(item) not in _CLOSED_GROUPS
        ]

    @mcp.tool()
    def create_work_item(
        project_id: ShortUUID,
        name: str,
        assignees: list[ShortUUID] | None = None,
        labels: list[ShortUUID] | None = None,
        type_id: str | None = None,
        point: int | None = None,
        description_html: str | None = None,
        description_stripped: str | None = None,
        priority: PriorityEnum | str | None = None,
        start_date: str | None = None,
        target_date: str | None = None,
        sort_order: float | None = None,
        is_draft: bool | None = None,
        external_source: str | None = None,
        external_id: str | None = None,
        parent: ShortUUID | None = None,
        state: str | None = None,
        estimate_point: str | None = None,
        type: str | None = None,
    ) -> dict:
        """
        Create a new work item.

        Args:
            project_id: UUID of the project
            name: Work item name (required)
            assignees: List of user IDs to assign to the work item
            labels: List of label IDs to attach to the work item
            type_id: UUID of the work item type
            point: Story point value
            description_html: HTML description of the work item
            description_stripped: Plain text description (stripped of HTML)
            priority: Priority level (urgent, high, medium, low, none)
            start_date: Start date (ISO 8601 format)
            target_date: Target/end date (ISO 8601 format)
            sort_order: Sort order value
            is_draft: Whether the work item is a draft
            external_source: External system source name
            external_id: External system identifier
            parent: UUID of the parent work item
            state: State name (e.g. "Done", "In Progress") or UUID
            estimate_point: Estimate point value
            type: Work item type identifier

        Returns:
            Created WorkItemSummary object
        """
        client, workspace_slug = get_plane_client_context()
        state = _resolve_state(client, workspace_slug, project_id, state)

        data = CreateWorkItem(
            name=name,
            assignees=assignees,
            labels=labels,
            type_id=type_id,
            point=point,
            description_html=description_html,
            description_stripped=description_stripped,
            priority=priority,
            start_date=start_date,
            target_date=target_date,
            sort_order=sort_order,
            is_draft=is_draft,
            external_source=external_source,
            external_id=external_id,
            parent=parent,
            state=state,
            estimate_point=estimate_point,
            type=type,
        )

        item = client.work_items.create(
            workspace_slug=workspace_slug, project_id=project_id, data=data
        )
        project = client.projects.retrieve(workspace_slug=workspace_slug, project_id=project_id)
        identifier = f"{project.identifier or '?'}-{item.sequence_id}"
        return {"message": f"{identifier} created successfully"}

    @mcp.tool()
    def retrieve_work_item(
        project_id: ShortUUID,
        work_item_id: ShortUUID,
        expand: str | None = None,
        fields: str | None = None,
        external_id: str | None = None,
        external_source: str | None = None,
        order_by: str | None = None,
    ) -> dict:
        """
        Retrieve a work item by ID.

        Args:
            project_id: UUID of the project
            work_item_id: UUID of the work item
            expand: Comma-separated fields to expand (e.g., "assignees,labels,state")
            fields: Comma-separated fields to include in response
            external_id: External system identifier for filtering
            external_source: External system source name for filtering
            order_by: Field to order results by (typically not used for single item retrieval)

        Returns:
            WorkItemFull with plain-text description and expanded assignees/labels.
        """
        client, workspace_slug = get_plane_client_context()

        params = RetrieveQueryParams(
            expand=_ensure_expand(expand, "state", "assignees"),
            fields=fields,
            external_id=external_id,
            external_source=external_source,
            order_by=order_by,
        )

        detail: WorkItemDetail = client.work_items.retrieve(
            workspace_slug=workspace_slug,
            project_id=project_id,
            work_item_id=work_item_id,
            params=params,
        )
        return _to_work_item_full(detail).slim()

    @mcp.tool()
    def retrieve_work_item_by_identifier(
        project_identifier: str,
        issue_identifier: int,
        expand: str | None = None,
        fields: str | None = None,
        external_id: str | None = None,
        external_source: str | None = None,
        order_by: str | None = None,
    ) -> dict:
        """
        Retrieve a work item by project identifier and issue sequence number.

        Args:
            project_identifier: Project identifier string (e.g., "MP" for "My Project")
            issue_identifier: Issue sequence number (e.g., 1, 2, 3)
            expand: Comma-separated fields to expand (e.g., "assignees,labels,state")
            fields: Comma-separated list of fields to include in response
            external_id: External system identifier for filtering
            external_source: External system source name for filtering
            order_by: Field to order results by (typically not used for single item retrieval)

        Returns:
            WorkItemFull with plain-text description and expanded assignees/labels.
        """
        client, workspace_slug = get_plane_client_context()

        params = RetrieveQueryParams(
            expand=_ensure_expand(expand, "state", "assignees"),
            fields=fields,
            external_id=external_id,
            external_source=external_source,
            order_by=order_by,
        )

        detail: WorkItemDetail = client.work_items.retrieve_by_identifier(
            workspace_slug=workspace_slug,
            project_identifier=project_identifier,
            issue_identifier=issue_identifier,
            params=params,
        )
        return _to_work_item_full(detail).slim()

    @mcp.tool()
    def update_work_item(
        project_id: ShortUUID,
        work_item_id: ShortUUID,
        name: str | None = None,
        assignees: list[ShortUUID] | None = None,
        labels: list[ShortUUID] | None = None,
        type_id: str | None = None,
        point: int | None = None,
        description_html: str | None = None,
        description_stripped: str | None = None,
        priority: PriorityEnum | str | None = None,
        start_date: str | None = None,
        target_date: str | None = None,
        sort_order: float | None = None,
        is_draft: bool | None = None,
        external_source: str | None = None,
        external_id: str | None = None,
        parent: ShortUUID | None = None,
        state: str | None = None,
        estimate_point: str | None = None,
        type: str | None = None,
    ) -> dict:
        """
        Update a work item by ID.

        Args:
            project_id: UUID of the project
            work_item_id: UUID of the work item
            name: Work item name
            assignees: List of user IDs to assign to the work item
            labels: List of label IDs to attach to the work item
            type_id: UUID of the work item type
            point: Story point value
            description_html: HTML description of the work item
            description_stripped: Plain text description (stripped of HTML)
            priority: Priority level (urgent, high, medium, low, none)
            start_date: Start date (ISO 8601 format)
            target_date: Target/end date (ISO 8601 format)
            sort_order: Sort order value
            is_draft: Whether the work item is a draft
            external_source: External system source name
            external_id: External system identifier
            parent: UUID of the parent work item
            state: State name (e.g. "Done", "In Progress") or UUID
            estimate_point: Estimate point value
            type: Work item type identifier

        Returns:
            Updated WorkItemSummary object
        """
        client, workspace_slug = get_plane_client_context()
        state = _resolve_state(client, workspace_slug, project_id, state)

        data = UpdateWorkItem(
            name=name,
            assignees=assignees,
            labels=labels,
            type_id=type_id,
            point=point,
            description_html=description_html,
            description_stripped=description_stripped,
            priority=priority,
            start_date=start_date,
            target_date=target_date,
            sort_order=sort_order,
            is_draft=is_draft,
            external_source=external_source,
            external_id=external_id,
            parent=parent,
            state=state,
            estimate_point=estimate_point,
            type=type,
        )

        item = client.work_items.update(
            workspace_slug=workspace_slug,
            project_id=project_id,
            work_item_id=work_item_id,
            data=data,
        )
        project = client.projects.retrieve(workspace_slug=workspace_slug, project_id=project_id)
        identifier = f"{project.identifier or '?'}-{item.sequence_id}"
        return {"message": f"{identifier} updated successfully"}

    @mcp.tool()
    def delete_work_item(project_id: ShortUUID, work_item_id: ShortUUID) -> None:
        """
        Delete a work item by ID.

        Args:
            workspace_slug: The workspace slug identifier
            project_id: UUID of the project
            work_item_id: UUID of the work item
        """
        client, workspace_slug = get_plane_client_context()
        client.work_items.delete(
            workspace_slug=workspace_slug, project_id=project_id, work_item_id=work_item_id
        )

    @mcp.tool()
    def search_work_items(
        query: str,
        expand: str | None = None,
        fields: str | None = None,
        external_id: str | None = None,
        external_source: str | None = None,
        order_by: str | None = None,
    ) -> WorkItemSearch:
        """
        Search work items across a workspace.

        Args:
            workspace_slug: The workspace slug identifier
            query: This is a free-form text search and will be used to search the work items
                    by name, description etc.
            expand: Comma-separated list of related fields to expand in response
            fields: Comma-separated list of fields to include in response
            external_id: External system identifier for filtering
            external_source: External system source name for filtering
            order_by: Field to order results by. Prefix with '-' for descending order

        Returns:
            WorkItemSearch object containing search results
        """
        client, workspace_slug = get_plane_client_context()

        params = RetrieveQueryParams(
            expand=expand,
            fields=fields,
            external_id=external_id,
            external_source=external_source,
            order_by=order_by,
        )

        return client.work_items.search(workspace_slug=workspace_slug, query=query, params=params)


def _retrieve_expanded(client, workspace_slug: str, project_id: str, work_item_id: str) -> dict:
    """Re-fetch a work item with state+assignees expanded, return slim summary."""
    detail = client.work_items.retrieve(
        workspace_slug=workspace_slug,
        project_id=project_id,
        work_item_id=work_item_id,
        params=RetrieveQueryParams(expand="state,assignees"),
    )
    return _to_work_item_summary(detail).slim()


def _enum_str(v) -> str | None:
    if v is None:
        return None
    return v.value if hasattr(v, "value") else str(v)


def _to_work_item_summary(item) -> WorkItemSummary:
    assignees = getattr(item, "assignees", None) or []
    labels = getattr(item, "labels", None) or []
    assignee_names = [
        a if isinstance(a, str) else (getattr(a, "display_name", None) or a.id)
        for a in assignees if a
    ]
    label_ids = [lb if isinstance(lb, str) else lb.id for lb in labels if lb]
    return WorkItemSummary(
        id=item.id,
        sequence_id=item.sequence_id,
        name=item.name,
        priority=_enum_str(item.priority),
        state=item.state if isinstance(item.state, str) else getattr(item.state, "name", None) or getattr(item.state, "id", None),
        assignees=[a for a in assignee_names if a],
        labels=[lb for lb in label_ids if lb],
    )


def _to_work_item_full(detail: WorkItemDetail) -> WorkItemFull:
    assignees = [
        AssigneeSummary(id=a, display_name=None) if isinstance(a, str)
        else AssigneeSummary(
            id=None if a.display_name else a.id,
            display_name=a.display_name,
        )
        for a in (detail.assignees or [])
    ]
    labels = [
        LabelSummary(id=lb, name=lb) if isinstance(lb, str)
        else LabelSummary(id=lb.id, name=lb.name, color=lb.color)
        for lb in (detail.labels or [])
    ]
    state = detail.state if isinstance(detail.state, str) else getattr(detail.state, "name", None) or getattr(detail.state, "id", None)
    return WorkItemFull(
        id=detail.id,
        sequence_id=detail.sequence_id,
        name=detail.name,
        description=strip_html(detail.description_html),
        priority=_enum_str(detail.priority),
        state=state,
        assignees=assignees,
        labels=labels,
        parent=detail.parent,
        start_date=detail.start_date,
        target_date=detail.target_date,
    )
