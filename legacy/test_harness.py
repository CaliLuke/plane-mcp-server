#!/usr/bin/env python3
"""Quick test harness for Plane API calls without MCP."""

import os
os.environ.setdefault("PLANE_API_KEY", "plane_api_d4f6412ae79440a4a5633ea7fe20950d")
os.environ.setdefault("PLANE_BASE_URL", "https://plane.luckymethod.com")
os.environ.setdefault("PLANE_WORKSPACE_SLUG", "autok")

from plane import PlaneClient
from plane.models.query_params import RetrieveQueryParams, WorkItemQueryParams

client = PlaneClient(
    base_url=os.environ["PLANE_BASE_URL"],
    api_key=os.environ["PLANE_API_KEY"],
)
WS = os.environ["PLANE_WORKSPACE_SLUG"]
PROJECT_ID = "cRJFz4vGXTHBAKXV548YkV"  # SERVE


def test_expand(endpoint: str, expand: str, **kwargs):
    """Test an expand value and report success/failure."""
    try:
        if endpoint == "list":
            params = WorkItemQueryParams(expand=expand, per_page=2)
            r = client.work_items.list(workspace_slug=WS, project_id=PROJECT_ID, params=params)
            print(f"  list expand={expand!r}: OK ({len(r.results)} results)")
            if r.results:
                item = r.results[0]
                print(f"    state={item.state!r}  assignees={item.assignees!r}")
        elif endpoint == "retrieve":
            params = RetrieveQueryParams(expand=expand)
            r = client.work_items.retrieve_by_identifier(
                workspace_slug=WS, project_identifier="SERVE",
                issue_identifier=kwargs.get("seq", 46), params=params,
            )
            print(f"  retrieve expand={expand!r}: OK")
            print(f"    state={r.state!r}  assignees={r.assignees!r}")
    except Exception as e:
        print(f"  {endpoint} expand={expand!r}: FAILED — {e}")


if __name__ == "__main__":
    print("=== Testing list endpoint ===")
    test_expand("list", "state")
    test_expand("list", "assignees")
    test_expand("list", "state,assignees")

    print("\n=== Testing retrieve endpoint ===")
    test_expand("retrieve", "state")
    test_expand("retrieve", "assignees")
    test_expand("retrieve", "state,assignees")

    print("\n=== Testing search ===")
    try:
        r = client.work_items.search(workspace_slug=WS, query="SERVE-46")
        print(f"  search 'SERVE-46': {len(r.issues)} results")
        for i in r.issues[:3]:
            print(f"    {i.name}")
    except Exception as e:
        print(f"  search FAILED: {e}")

    try:
        r = client.work_items.search(workspace_slug=WS, query="DELETE Graph Node")
        print(f"  search 'DELETE Graph Node': {len(r.issues)} results")
        for i in r.issues[:3]:
            print(f"    {i.name}")
    except Exception as e:
        print(f"  search FAILED: {e}")
