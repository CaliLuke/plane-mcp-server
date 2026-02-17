"""Compact single-line formatters for MCP list responses."""

from plane_mcp import uid


def project(p) -> str:
    """'IDENTIFIER · Name [short_id]'"""
    short = uid.encode(p.id) or "?"
    return f"{p.identifier} · {p.name} [{short}]"


def work_item(item, project_identifier: str = "", state_name: str | None = None) -> str:
    """'PROJ-seq (state) · Name [short_id]'"""
    short = uid.encode(item.id) or "?"
    seq = f"{project_identifier}-{item.sequence_id}" if project_identifier else f"#{item.sequence_id}"
    state_part = f" ({state_name})" if state_name else ""
    return f"{seq}{state_part} · {item.name} [{short}]"


def state(s) -> str:
    """'Name (group) [short_id]' — marks default with *"""
    short = uid.encode(s.id) or "?"
    group = _enum_str(s.group) or ""
    default = " *" if s.default else ""
    return f"{s.name} ({group}){default} [{short}]"


def _enum_str(v) -> str | None:
    if v is None:
        return None
    return v.value if hasattr(v, "value") else str(v)
