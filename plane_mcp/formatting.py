"""Compact single-line formatters for MCP list responses."""

from plane_mcp import uid


def project(p) -> str:
    """'IDENTIFIER · Name [short_id]'"""
    short = uid.encode(p.id) or "?"
    return f"{p.identifier} · {p.name} [{short}]"


def work_item(item) -> str:
    """'#seq · Name (priority) [short_id]'"""
    short = uid.encode(item.id) or "?"
    priority = _enum_str(item.priority)
    meta = f"({priority})" if priority and priority != "none" else ""
    parts = [f"#{item.sequence_id}", "·", item.name]
    if meta:
        parts.append(meta)
    parts.append(f"[{short}]")
    return " ".join(parts)


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
